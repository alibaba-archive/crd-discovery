package sync

import (
	"encoding/json"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"net/http"
)

type Syncer struct {
	Logger        logrus.FieldLogger
	DynamicClient dynamic.Interface
}

func NewSyncerOrDie(logger logrus.FieldLogger) *Syncer {
	config := utils.LoadKubeConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)
	return &Syncer{
		Logger:        logger,
		DynamicClient: dynamicClient,
	}
}

func (syncer *Syncer) WithGVR(gvr schema.GroupVersionResource) (logrus.FieldLogger, dynamic.NamespaceableResourceInterface) {
	return syncer.Logger.WithFields(logrus.Fields{
		"Group":    gvr.Group,
		"Version":  gvr.Version,
		"Resource": gvr.Resource,
	}), syncer.DynamicClient.Resource(gvr)
}

func (syncer *Syncer) Fetch(gvr schema.GroupVersionResource) FetchResult {
	logger, client := syncer.WithGVR(gvr)
	unstructuredList, err := client.List(v1.ListOptions{})
	if err != nil {
		logger.Errorf("list resource failed: %s\n", err.Error())
		return FetchResult{Code: http.StatusNotFound, Err: err.Error()}
	}
	return FetchResult{Code: http.StatusOK, Objects: unstructuredList.Items}
}

func (syncer *Syncer) Pull(gvr schema.GroupVersionResource, remoteReader io.Reader) PullResult {
	logger, client := syncer.WithGVR(gvr)

	var remoteObjects []unstructured.Unstructured
	bytes, err := ioutil.ReadAll(remoteReader)
	if err != nil {
		syncer.Logger.Errorf("read body failed: %s\n", err.Error())
		return PullResult{Code: http.StatusBadRequest, Err: err.Error()}
	}
	if err := json.Unmarshal(bytes, &remoteObjects); err != nil {
		syncer.Logger.Errorf("unmarshal body failed: %s\n", err.Error())
		return PullResult{Code: http.StatusBadRequest, Err: err.Error()}
	}

	unstructuredList, err := client.List(v1.ListOptions{})
	if err != nil {
		return PullResult{Code: http.StatusBadRequest, Err: err.Error()}
	}
	localObjectMap := make(map[string]unstructured.Unstructured)
	for _, object := range unstructuredList.Items {
		localObjectMap[utils.GetNamespacedName(&object)] = object
	}

	result := PullResult{Code: http.StatusOK}
	for _, object := range remoteObjects {
		newObject := object.DeepCopy()
		nn := utils.GetNamespacedName(newObject)
		if oldObject, ok := localObjectMap[nn]; ok { // existing
			newObject.SetUID(oldObject.GetUID())
			newObject.SetResourceVersion(oldObject.GetResourceVersion())
			if _, err := client.Namespace(oldObject.GetNamespace()).Update(newObject, v1.UpdateOptions{}); err != nil {
				logger.Errorf("update failed: %s\n", err.Error())
				return PullResult{Code: http.StatusInternalServerError, Err: err.Error()}
			}
			result.Updated = append(result.Updated, nn)
		} else { // new
			newObject.SetUID("")
			newObject.SetResourceVersion("")
			if _, err := client.Namespace(newObject.GetNamespace()).Create(newObject, v1.CreateOptions{}); err != nil {
				logger.Errorf("create failed: %s\n", err.Error())
				return PullResult{Code: http.StatusInternalServerError, Err: err.Error()}
			}
			result.Created = append(result.Created, nn)
		}
		delete(localObjectMap, nn)
	}
	for nn, oldObject := range localObjectMap { // outdated
		if err := client.Namespace(oldObject.GetNamespace()).Delete(oldObject.GetName(), &v1.DeleteOptions{}); err != nil {
			logger.Errorf("delete failed: %s\n", err.Error())
			return PullResult{Code: http.StatusInternalServerError, Err: err.Error()}
		}
		result.Deleted = append(result.Deleted, nn)
	}
	return result
}
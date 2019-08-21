package main

import (
	"encoding/json"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"net/http"
)

type Server struct {
	logger        logrus.FieldLogger
	dynamicClient dynamic.Interface
}

func NewServer(logger logrus.FieldLogger) *Server {
	config := utils.LoadKubeConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)
	return &Server{
		logger:        logger,
		dynamicClient: dynamicClient,
	}
}

func (s *Server) pull(w http.ResponseWriter, r *http.Request) {
	resource, logger := s.extractGVR(r)
	unstructuredList, err := resource.List(v1.ListOptions{})
	if err != nil {
		logger.Errorf("list resource failed: %s\n", err.Error())
		http.Error(w, err.Error(), 404)
		return
	}
	s.writeResponse(w, logger, unstructuredList.Items)
}

func (s *Server) push(w http.ResponseWriter, r *http.Request) {
	resource, logger := s.extractGVR(r)
	unstructuredList, err := resource.List(v1.ListOptions{})
	if err != nil {
		logger.Errorf("list resource failed: %s\n", err.Error())
		http.Error(w, err.Error(), 404)
		return
	}
	itemMap := make(map[string]unstructured.Unstructured)
	for _, item := range unstructuredList.Items {
		itemMap[utils.GetNamespacedName(&item)] = item
	}
	var slaveItems []unstructured.Unstructured
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("read body failed: %s\n", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
	if err := json.Unmarshal(bytes, &slaveItems); err != nil {
		logger.Errorf("unmarshal body failed: %s\n", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
	var created, updated, deleted []string
	for _, _item := range slaveItems {
		item := _item.DeepCopy()
		nn := utils.GetNamespacedName(item)
		if oldItem, ok := itemMap[nn]; ok { // existing
			item.SetUID(oldItem.GetUID())
			item.SetResourceVersion(oldItem.GetResourceVersion())
			if _, err := resource.Namespace(item.GetNamespace()).Update(item, v1.UpdateOptions{}); err != nil {
				logger.Errorf("update failed: %s\n", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			updated = append(updated, nn)
		} else { // new
			item.SetUID("")
			item.SetResourceVersion("")
			if _, err := resource.Namespace(item.GetNamespace()).Create(item, v1.CreateOptions{}); err != nil {
				logger.Errorf("create failed: %s\n", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			created = append(created, nn)
		}
		delete(itemMap, nn)
	}
	for nn, item := range itemMap {
		if err := resource.Namespace(item.GetNamespace()).Delete(item.GetName(), &v1.DeleteOptions{}); err != nil {
			logger.Errorf("delete failed: %s\n", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		deleted = append(deleted, nn)
	}
	s.writeResponse(w, logger, map[string][]string{"created": created, "updated": updated, "deleted": deleted})
}

func (s *Server) extractGVR(r *http.Request) (dynamic.NamespaceableResourceInterface, logrus.FieldLogger) {
	vars := mux.Vars(r)
	gvr := schema.GroupVersionResource{
		Group: vars["group"],
		Version: vars["version"],
		Resource: vars["resource"],
	}
	logger := s.logger.WithFields(logrus.Fields{
		"group": gvr.Group,
		"version": gvr.Version,
		"resource": gvr.Resource,
	})
	return s.dynamicClient.Resource(gvr), logger
}

func (s *Server) writeResponse(w http.ResponseWriter, logger logrus.FieldLogger, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(result)
	if err != nil {
		logger.Errorf("marshal failed: %s\n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	if _, err = w.Write(bytes); err != nil {
		logger.Errorf("write response failed: %s\n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

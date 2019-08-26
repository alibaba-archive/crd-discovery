package sync

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	"net/http"
	"testing"
)

func TestNewSyncerOrDie(t *testing.T) {
	_ = NewSyncerOrDie(logrus.StandardLogger())
}

func getFakeSyncer() *Syncer {
	return &Syncer{
		Logger:        logrus.StandardLogger(),
		DynamicClient: fake.NewSimpleDynamicClient(runtime.NewScheme()),
	}
}

func getFakeGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "g",
		Version:  "v",
		Resource: "r",
	}
}

func TestSyncer_Fetch(t *testing.T) {
	syncer := getFakeSyncer()
	gvr := getFakeGVR()
	name, namespace := "obj-name", "obj-namespace"
	obj := unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": name,
			"namespace": namespace,
		},
	}}
	_, err := syncer.DynamicClient.Resource(gvr).Namespace(namespace).Create(&obj, v1.CreateOptions{})
	require.NoError(t, err)
	result := syncer.Fetch(gvr)
	require.Equal(t, http.StatusOK, result.Code)
	require.Equal(t, 1, len(result.Objects))
	require.Equal(t, name, result.Objects[0].GetName())
	require.Equal(t, namespace, result.Objects[0].GetNamespace())
}

func TestSyncer_Pull(t *testing.T) {
	syncer := getFakeSyncer()
	gvr := getFakeGVR()
	namespace := "obj-namespace"
	obj0 := unstructured.Unstructured{Object: map[string]interface{}{
		"kind": "R",
		"metadata": map[string]interface{}{
			"name": "obj-name-0",
			"namespace": namespace,
		},
	}}
	_, err := syncer.DynamicClient.Resource(gvr).Namespace(namespace).Create(&obj0, v1.CreateOptions{})
	require.NoError(t, err)
	obj1 := unstructured.Unstructured{Object: map[string]interface{}{
		"kind": "R",
		"metadata": map[string]interface{}{
			"name": "obj-name-1",
			"namespace": namespace,
		},
	}}
	_, err = syncer.DynamicClient.Resource(gvr).Namespace(namespace).Create(&obj1, v1.CreateOptions{})
	require.NoError(t, err)
	obj2 := unstructured.Unstructured{Object: map[string]interface{}{
		"kind": "R",
		"metadata": map[string]interface{}{
			"name": "obj-name-2",
			"namespace": namespace,
		},
	}}
	bs, err := json.Marshal(&[]unstructured.Unstructured{obj1, obj2})
	require.NoError(t, err)
	remoteReader := bytes.NewReader(bs)
	result := syncer.Pull(gvr, remoteReader)
	require.Equal(t, http.StatusOK, result.Code)
	require.Equal(t, "", result.Err)
	require.Equal(t, 1, len(result.Created))
	require.Equal(t, 1, len(result.Updated))
	require.Equal(t, 1, len(result.Deleted))
	unstructuredList, err := syncer.DynamicClient.Resource(gvr).List(v1.ListOptions{})
	require.NoError(t, err)
	require.Equal(t, 2, len(unstructuredList.Items))
}
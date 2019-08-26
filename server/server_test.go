package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/sync"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	_ = NewServer(logrus.StandardLogger())
}

func TestServer(t *testing.T) {
	masterServer := NewFakeServer()
	namespace := "obj-namespace"
	gvr := sync.NewFakeGVR()
	obj0 := sync.NewFakeObject("obj-name-0", namespace)
	_, err := masterServer.syncer.DynamicClient.Resource(gvr).Namespace(namespace).Create(&obj0, v1.CreateOptions{})
	require.NoError(t, err)
	router := mux.NewRouter()
	router.HandleFunc("/sync/pull/{group}/{version}/{resource}", masterServer.pull)
	router.HandleFunc("/sync/push/{group}/{version}/{resource}", masterServer.push)
	go func() {
		if err := http.ListenAndServe(":18080", router); err != nil {
			fmt.Println(err.Error())
		}
	}()
	time.Sleep(time.Second)
	resp, err := http.Get("http://:18080/sync/pull/g/v/r")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	syncer := sync.NewFakeSyncer()
	pr := syncer.Pull(gvr, resp.Body)
	require.Equal(t, http.StatusOK, pr.Code)
	unstructuredList, err := syncer.DynamicClient.Resource(gvr).List(v1.ListOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, len(unstructuredList.Items))
	require.Equal(t, "obj-name-0", unstructuredList.Items[0].GetName())
	require.Equal(t, namespace, unstructuredList.Items[0].GetNamespace())
	obj1 := sync.NewFakeObject("obj-name-1", namespace)
	_, err = syncer.DynamicClient.Resource(gvr).Namespace(namespace).Create(&obj1, v1.CreateOptions{})
	require.NoError(t, err)
	err = syncer.DynamicClient.Resource(gvr).Namespace(obj0.GetNamespace()).Delete(obj0.GetName(), &v1.DeleteOptions{})
	require.NoError(t, err)
	fr := syncer.Fetch(gvr)
	require.Equal(t, http.StatusOK, fr.Code)
	bs, err := json.Marshal(&fr.Objects)
	require.NoError(t, err)
	resp, err = http.Post("http://:18080/sync/push/g/v/r", "application/json", bytes.NewReader(bs))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	unstructuredList, err = masterServer.syncer.DynamicClient.Resource(gvr).List(v1.ListOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, len(unstructuredList.Items))
	require.Equal(t, "obj-name-1", unstructuredList.Items[0].GetName())
	require.Equal(t, namespace, unstructuredList.Items[0].GetNamespace())
}

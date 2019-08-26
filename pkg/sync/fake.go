package sync

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func NewFakeGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "g",
		Version:  "v",
		Resource: "r",
	}
}

func NewFakeObject(name, namespace string) unstructured.Unstructured {
	return unstructured.Unstructured{Object: map[string]interface{}{
		"kind": "R",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
	}}
}

func NewFakeSyncer() *Syncer {
	return &Syncer{
		Logger:        logrus.StandardLogger(),
		DynamicClient: fake.NewSimpleDynamicClient(runtime.NewScheme()),
	}
}

package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func TestLoadKubeConfigOrDie(t *testing.T) {
	_ = LoadKubeConfigOrDie()
}

func TestGetNamespacedName(t *testing.T) {
	name := "obj-name"
	namespace := "obj-namespace"
	obj := unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
				"namespace": namespace,
			},
		},
	}
	require.Equal(t, GetNamespacedName(&obj), fmt.Sprintf("%s/%s", namespace, name))
}
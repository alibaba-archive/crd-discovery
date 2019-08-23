package cmd

import (
	"github.com/Somefive/crd-discovery/pkg/utils"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var crdGVR = schema.GroupVersionResource{
	Group:    "apiextensions.k8s.io",
	Version:  "v1beta1",
	Resource: "customresourcedefinitions",
}

func getGVRs(kinds []string) []schema.GroupVersionResource {
	config := utils.LoadKubeConfigOrDie()
	client := apixv1beta1client.NewForConfigOrDie(config)
	refs, err := client.CustomResourceDefinitions().List(v1.ListOptions{})
	utils.ErrExit("get crds failed", err)
	var gvrs []schema.GroupVersionResource
	for _, item := range refs.Items {
		selected := len(kinds) == 0
		for _, kind := range kinds {
			if kind == item.Spec.Names.Kind {
				selected = true
				break
			}
		}
		if selected {
			gvrs = append(gvrs, schema.GroupVersionResource{
				Group:    item.Spec.Group,
				Version:  item.Spec.Versions[0].Name,
				Resource: item.Spec.Names.Plural,
			})
		}
	}
	return gvrs
}
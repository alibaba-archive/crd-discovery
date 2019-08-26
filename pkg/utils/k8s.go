package utils

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/user"
	"path/filepath"
)

func LoadKubeConfigOrDie() *rest.Config {
	usr, err := user.Current()
	if err != nil {
		config, err := rest.InClusterConfig()
		ErrExit("load in cluster config failed", err)
		return config
	}
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(usr.HomeDir, ".kube", "config"))
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "")
		ErrExit("load local config failed", err)
	}
	return config
}

func ErrExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}

func GetNamespacedName(object *unstructured.Unstructured) string {
	return object.GetNamespace() + "/" + object.GetName()
}

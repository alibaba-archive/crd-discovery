package main

import (
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/spf13/cobra"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cmdSync = &cobra.Command{
	Use:   "sync <crds...>",
	Short: "sync crds",
	Long:  "sync crds from master k8s",
	Run:   sync,
}

func init() {
	rootCmd.AddCommand(cmdSync)
}

func sync(cmd *cobra.Command, args []string) {
	crds, err := getCRDs(args)
	utils.ErrExit("Fetch CRD from remote failed", err)
	client := apixv1beta1client.NewForConfigOrDie(utils.LoadKubeConfigOrDie())
	for _, crd := range crds {
		obj := crd.DeepCopy()
		_obj, err := client.CustomResourceDefinitions().Get(obj.Name, v1.GetOptions{})
		if err != nil {
			obj.SetUID("")
			obj.SetResourceVersion("")
			_, err = client.CustomResourceDefinitions().Create(obj)
			utils.ErrExit("create " + obj.Name + " failed", err)
			fmt.Println("create " + obj.Name + " succeed")
		} else {
			obj.SetUID(_obj.GetUID())
			obj.SetResourceVersion(_obj.GetResourceVersion())
			_, err = client.CustomResourceDefinitions().Update(obj)
			utils.ErrExit("update " + obj.Name + " failed", err)
			fmt.Println("update " + obj.Name + " succeed")
		}
	}
	fmt.Println("sync succeed")
}

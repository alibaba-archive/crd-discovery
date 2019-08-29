package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

var cmdUpdate = &cobra.Command{
	Use:   "update <crds...>",
	Short: "update crd",
	Long:  "update crd from local k8s to master k8s",
	Run:   update,
}

func init() {
	rootCmd.AddCommand(cmdUpdate)
}

func update(cmd *cobra.Command, args []string) {
	client := apixv1beta1client.NewForConfigOrDie(utils.LoadKubeConfigOrDie())
	var crds []v1beta1.CustomResourceDefinition
	if len(args) == 0 {
		definitionList, err := client.CustomResourceDefinitions().List(v1.ListOptions{})
		utils.ErrExit("load local crds failed", err)
		crds = definitionList.Items
	} else {
		for _, arg := range args {
			crd, err := client.CustomResourceDefinitions().Get(arg, v1.GetOptions{})
			utils.ErrExit("load crd " + arg + " failed", err)
			crds = append(crds, *crd)
		}
	}
	for _, crd := range crds {
		_crds, err := getCRDs([]string{crd.Name})
		var resp *http.Response
		if err == nil {
			_crd := _crds[0]
			crd.SetUID(_crd.GetUID())
			crd.SetResourceVersion(_crd.GetResourceVersion())
			bs, err := json.Marshal(crd)
			utils.ErrExit("marshal crd " + crd.Name + " failed", err)
			url := fmt.Sprintf("%s://%s/update", getProtocol(), masterURL)
			resp, err = http.Post(url, "application/json", bytes.NewReader(bs))
			utils.ErrExit("master:update crd " + crd.Name + " failed", err)
		} else {
			crd.SetUID("")
			crd.SetResourceVersion("")
			bs, err := json.Marshal(crd)
			utils.ErrExit("marshal crd " + crd.Name + " failed", err)
			url := fmt.Sprintf("%s://%s/create", getProtocol(), masterURL)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bs))
			utils.ErrExit("Construct http request failed", err)
			resp, err = http.DefaultClient.Do(req)
			utils.ErrExit("master:create crd " + crd.Name + " failed", err)
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("update crd " + crd.Name + " succeed")
		} else {
			all, err := ioutil.ReadAll(resp.Body)
			utils.ErrExit("read response body failed", err)
			fmt.Println("update crd " + crd.Name + " failed: " + string(all))
		}
	}
	fmt.Println("update completed")
}

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
	"net/http"
	"sigs.k8s.io/yaml"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "create crd",
	Long:  "create crd in both local and master k8s",
	Run:   create,
}

const (
	FlagFile = "file"
	ShorthandFile = "f"
)

var fileName string

func init() {
	rootCmd.PersistentFlags().StringVarP(&fileName, FlagFile, ShorthandFile, "", "The yaml file for crd to create")
	rootCmd.MarkFlagRequired(FlagFile)
	rootCmd.AddCommand(cmdCreate)
}

func create(cmd *cobra.Command, args []string) {
	bs, err := ioutil.ReadFile(fileName)
	utils.ErrExit("Read CRD from " + fileName + " failed", err)
	crd := &v1beta1.CustomResourceDefinition{}
	err = yaml.Unmarshal(bs, crd)
	utils.ErrExit("Unmarshal CRD in file failed", err)
	client := apixv1beta1client.NewForConfigOrDie(utils.LoadKubeConfigOrDie())
	_, err = client.CustomResourceDefinitions().Create(crd)
	utils.ErrExit("Create CRD locally failed", err)
	bs, err = json.Marshal(crd)
	utils.ErrExit("Marshal CRD to json failed", err)
	url := fmt.Sprintf("%s://%s/create", getProtocol(), masterURL)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bs))
	utils.ErrExit("Construct http request failed", err)
	resp, err := http.DefaultClient.Do(req)
	utils.ErrExit("Send http request failed", err)
	if resp.StatusCode == http.StatusOK {
		fmt.Println("create succeed")
	} else {
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("create failed: ", err.Error())
		} else {
			fmt.Println("create failed: ", string(all))
		}
	}
}

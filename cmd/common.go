package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"net/http"
)

func getCRDs(args []string) ([]v1beta1.CustomResourceDefinition, error) {
	var results []v1beta1.CustomResourceDefinition
	if len(args) == 0 {
		url := fmt.Sprintf("%s://%s/list", getProtocol(), masterURL)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		crds := &v1beta1.CustomResourceDefinitionList{}
		err = json.Unmarshal(bs, crds)
		if err != nil {
			return nil, err
		}
		results = crds.Items
	} else {
		for _, arg := range args {
			url := fmt.Sprintf("%s://%s/list/%s", getProtocol(), masterURL, arg)
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			crd := &v1beta1.CustomResourceDefinition{}
			err = json.Unmarshal(bs, crd)
			if err != nil {
				return nil, err
			}
			results = append(results, *crd)
		}
	}
	return results, nil
}

func getProtocol() string {
	if useHTTPS {
		return "https"
	}
	return "http"
}

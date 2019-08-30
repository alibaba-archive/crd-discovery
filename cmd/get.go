package main

import (
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var cmdGet = &cobra.Command{
	Use:   "get <optional: crds...>",
	Short: "get crds",
	Long:  "get crds on master k8s",
	Run:   get,
}

func init() {
	rootCmd.AddCommand(cmdGet)
}

func get(cmd *cobra.Command, args []string) {
	crds, err := getCRDs(args)
	utils.ErrExit("Fetch CRD from remote failed", err)
	for index, crd := range crds {
		bs, err := yaml.Marshal(crd)
		utils.ErrExit("Marshal CRD failed", err)
		if index > 0 {
			fmt.Println("---")
		}
		fmt.Println(string(bs))
	}
}

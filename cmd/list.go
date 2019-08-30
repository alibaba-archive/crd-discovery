package main

import (
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list crds",
	Long:  "list crds on master k8s",
	Run:   list,
}

func init() {
	rootCmd.AddCommand(cmdList)
}

func list(cmd *cobra.Command, args []string) {
	crds, err := getCRDs(args)
	utils.ErrExit("Fetch CRD from remote failed", err)
	fmt.Printf("%-30s%-30s\n", "NAME", "CREATED AT")
	for _, crd := range crds {
		fmt.Printf("%-30s%-30s\n", crd.Name, crd.CreationTimestamp)
	}
}

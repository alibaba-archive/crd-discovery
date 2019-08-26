package main

import (
	"github.com/spf13/cobra"
)

var cmdPull = &cobra.Command{
	Use:   "pull <kinds...>",
	Short: "manually pull objects",
	Long:  "manually pull objects from master k8s to current k8s",
	Run:   pull,
}

func init() {
	rootCmd.AddCommand(cmdPull)
}

func pull(cmd *cobra.Command, args []string) {
	client.pull(crdGVR)
	for _, gvr := range getGVRs(args) {
		client.pull(gvr)
	}
}

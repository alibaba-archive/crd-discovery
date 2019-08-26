package main

import (
	"github.com/spf13/cobra"
)

var cmdPush = &cobra.Command{
	Use: "push <kinds...>",
	Short: "manually push objects",
	Long: "manually push objects from current k8s to master k8s",
	Run: push,
}

func init() {
	rootCmd.AddCommand(cmdPush)
}

func push(cmd *cobra.Command, args []string) {
	client.pull(crdGVR)
	for _, gvr := range getGVRs(args) {
		client.push(gvr)
	}
}

package main

import (
	"github.com/spf13/cobra"
	"net/http"
)

var cmdSync = &cobra.Command{Use: "sync", Run: func(cmd *cobra.Command, args []string) {
	sync()
}}

func init() {
	rootCmd.AddCommand(cmdSync)
}

func getCRD() {

}

func sync() {
	http.Get(masterURL+"/")
}
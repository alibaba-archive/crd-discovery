package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

const (
	FlagMasterURL = "masterURL"
)

var masterURL string

var rootCmd = &cobra.Command{Use: "syncrd"}

func init() {
	rootCmd.PersistentFlags().StringVar(&masterURL, FlagMasterURL, "", "The url of server on master k8s ")
	rootCmd.MarkFlagRequired(FlagMasterURL)
	viper.BindPFlag(FlagMasterURL, rootCmd.PersistentFlags().Lookup(FlagMasterURL))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execute error: %s\n", err.Error())
	}
}

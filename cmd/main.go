package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

const (
	FlagMasterURL = "masterURL"
	FlagUseHTTPS  = "useHTTPS"
)

var masterURL string
var useHTTPS bool

var rootCmd = &cobra.Command{Use: "syncrd"}

func init() {
	rootCmd.PersistentFlags().StringVar(&masterURL, FlagMasterURL, "localhost:8088", "The url of server on master k8s")
	rootCmd.MarkFlagRequired(FlagMasterURL)
	viper.BindPFlag(FlagMasterURL, rootCmd.PersistentFlags().Lookup(FlagMasterURL))

	rootCmd.PersistentFlags().BoolVar(&useHTTPS, FlagUseHTTPS, false, "identify whether to enable https or not")
	viper.BindPFlag(FlagUseHTTPS, rootCmd.PersistentFlags().Lookup(FlagUseHTTPS))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execute error: %s\n", err.Error())
	}
}

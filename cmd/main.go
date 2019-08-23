package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

const (
	FlagMasterURL = "masterURL"
)

var masterURL string

var rootCmd = &cobra.Command{Use: "sync"}
var client *Client

func init() {
	rootCmd.PersistentFlags().StringVar(&masterURL, FlagMasterURL, "", "The url of server on master k8s ")
	rootCmd.MarkFlagRequired(FlagMasterURL)
	viper.BindPFlag(FlagMasterURL, rootCmd.PersistentFlags().Lookup(FlagMasterURL))

	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)
	client = NewClient(logger)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execute error: %s\n", err.Error())
	}
}

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

const (
	FlagAddr       = "addr"
	FlagEnablePull = "enable-pull"
	FlagEnablePush = "enable-push"
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server to sync",
	Long:  "run web server to serve client sync requests",
	Run:   serve,
}

var addr string
var enablePull bool
var enablePush bool

func init() {
	rootCmd.PersistentFlags().StringVar(&addr, FlagAddr, ":8080", "The address web server will listen")
	viper.BindPFlag(FlagAddr, rootCmd.PersistentFlags().Lookup(FlagAddr))

	rootCmd.PersistentFlags().BoolVar(&enablePull, FlagEnablePull, true, "Enable client pull")
	viper.BindPFlag(FlagEnablePull, rootCmd.PersistentFlags().Lookup(FlagEnablePull))

	rootCmd.PersistentFlags().BoolVar(&enablePush, FlagEnablePush, true, "Enable client push")
	viper.BindPFlag(FlagEnablePush, rootCmd.PersistentFlags().Lookup(FlagEnablePush))
}

func serve(cmd *cobra.Command, args []string) {
	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)
	server := NewServer(logger)

	router := mux.NewRouter()
	if enablePull {
		router.HandleFunc("/sync/pull/{group}/{version}/{resource}", server.pull)
	}
	if enablePush {
		router.HandleFunc("/sync/push/{group}/{version}/{resource}", server.push)
	}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello\n"))
	})

	fmt.Println("start listening at " + addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execute error: %s\n", err.Error())
	}
}

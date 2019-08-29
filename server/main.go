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
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server to sync",
	Long:  "run web server to serve client sync requests",
	Run:   serve,
}

var addr string

func init() {
	rootCmd.PersistentFlags().StringVar(&addr, FlagAddr, ":8080", "The address web server will listen")
	viper.BindPFlag(FlagAddr, rootCmd.PersistentFlags().Lookup(FlagAddr))
}

func serve(cmd *cobra.Command, args []string) {
	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)
	server := NewServer(logger)

	router := mux.NewRouter()
	router.HandleFunc("/list", server.list).Methods(http.MethodGet)
	router.HandleFunc("/list/{crd}", server.list).Methods(http.MethodGet)
	router.HandleFunc("/create", server.upsert).Methods(http.MethodPut)
	router.HandleFunc("/update", server.upsert).Methods(http.MethodPost)
	router.HandleFunc("/", server.index).Methods(http.MethodGet)

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

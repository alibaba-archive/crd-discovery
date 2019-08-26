package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)
	server := NewServer(logger)

	router := mux.NewRouter()
	router.HandleFunc("/sync/pull/{group}/{version}/{resource}", server.pull)
	router.HandleFunc("/sync/push/{group}/{version}/{resource}", server.push)

	addr := ":8080"
	fmt.Println("start listening at " + addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println(err.Error())
	}
}

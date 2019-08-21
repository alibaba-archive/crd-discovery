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

	fmt.Println("start listening")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err.Error())
	}
}

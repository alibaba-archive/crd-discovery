package main

import (
	"encoding/json"
	"github.com/Somefive/crd-discovery/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Server struct {
	logger logrus.FieldLogger
	client apixv1beta1client.ApiextensionsV1beta1Interface
}

func NewServer(logger logrus.FieldLogger) *Server {
	return &Server{
		logger: logger,
		client: apixv1beta1client.NewForConfigOrDie(utils.LoadKubeConfigOrDie()),
	}
}

func (s *Server) writeResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, err := json.Marshal(body)
	if err != nil {
		s.writeError(w, "marshal failed", http.StatusInternalServerError, err)
		return
	}
	if _, err = w.Write(bytes); err != nil {
		s.writeError(w, "write response failed", http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) writeError(w http.ResponseWriter, reason string, code int, err error) {
	s.logger.Errorln(reason, ": ", err.Error())
	http.Error(w, err.Error(), code)
	return
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	crd := mux.Vars(r)["crd"]
	var obj interface{}
	var err error
	if crd == "" {
		obj, err = s.client.CustomResourceDefinitions().List(v1.ListOptions{})
	} else {
		obj, err = s.client.CustomResourceDefinitions().Get(crd, v1.GetOptions{})
	}
	if err != nil {
		s.writeError(w, "fetch crd failed", http.StatusServiceUnavailable, err)
		return
	}
	s.writeResponse(w, http.StatusOK, obj)
}

func (s *Server) upsert(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, "read body failed", http.StatusBadRequest, err)
		return
	}
	crd := &v1beta1.CustomResourceDefinition{}
	if err = json.Unmarshal(bs, crd); err != nil {
		s.writeError(w, "unmarshal failed", http.StatusBadRequest, err)
		return
	}
	switch r.Method {
	case http.MethodPut: crd, err = s.client.CustomResourceDefinitions().Create(crd)
	case http.MethodPost: crd, err = s.client.CustomResourceDefinitions().Update(crd)
	}
	if err != nil {
		s.writeError(w, "create crd failed", http.StatusBadRequest, err)
		return
	}
	s.writeResponse(w, http.StatusOK, crd)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	// TODO output version information
	s.writeResponse(w, http.StatusOK, "hello, this is syncrd server")
}

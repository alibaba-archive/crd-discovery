package main

import (
	"encoding/json"
	"github.com/Somefive/crd-discovery/pkg/sync"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
)

type Server struct {
	syncer *sync.Syncer
}

func NewServer(logger logrus.FieldLogger) *Server {
	return &Server{
		syncer: sync.NewSyncerOrDie(logger),
	}
}

func (s *Server) pull(w http.ResponseWriter, r *http.Request) {
	gvr := s.extractGVR(r)
	result := s.syncer.Fetch(gvr)
	s.writeResponse(w, result.Code, &result.Objects)
}

func (s *Server) push(w http.ResponseWriter, r *http.Request) {
	gvr := s.extractGVR(r)
	result := s.syncer.Pull(gvr, r.Body)
	s.writeResponse(w, result.Code, &result)
}

func (s *Server) extractGVR(r *http.Request) schema.GroupVersionResource {
	vars := mux.Vars(r)
	return schema.GroupVersionResource{
		Group: vars["group"],
		Version: vars["version"],
		Resource: vars["resource"],
	}
}

func (s *Server) writeResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, err := json.Marshal(body)
	if err != nil {
		s.syncer.Logger.Errorf("marshal failed: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(bytes); err != nil {
		s.syncer.Logger.Errorf("write response failed: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

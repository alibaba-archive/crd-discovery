package main

import (
	"github.com/Somefive/crd-discovery/pkg/sync"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
)

func NewFakeServer() *Server {
	return &Server{
		syncer: &sync.Syncer{
			Logger:        logrus.StandardLogger(),
			DynamicClient: fake.NewSimpleDynamicClient(runtime.NewScheme()),
		},
	}
}

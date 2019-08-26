package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Somefive/crd-discovery/pkg/sync"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
)

type Client struct {
	syncer *sync.Syncer
}

func NewClient(logger logrus.FieldLogger) *Client {
	return &Client{
		syncer: sync.NewSyncerOrDie(logger),
	}
}

func (client *Client) pull(gvr schema.GroupVersionResource) {
	logger, _ := client.syncer.WithGVR(gvr)
	url := fmt.Sprintf("%s://%s/sync/pull/%s/%s/%s", getProtocol(), masterURL, gvr.Group, gvr.Version, gvr.Resource)
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("get from remote failed: %s\n", err.Error())
		return
	}
	pr := client.syncer.Pull(gvr, resp.Body)
	if pr.Code != http.StatusOK {
		return
	}
	client.syncer.Logger.Infof("local %s/%s:%s updated: %d created: %d deleted: %d\n", gvr.Group, gvr.Version, gvr.Resource, len(pr.Updated), len(pr.Created), len(pr.Deleted))
}

func (client *Client) push(gvr schema.GroupVersionResource) {
	logger, _ := client.syncer.WithGVR(gvr)
	result := client.syncer.Fetch(gvr)
	if result.Code != http.StatusOK {
		return
	}
	bs, err := json.Marshal(result.Objects)
	if err != nil {
		logger.Errorf("marshal failed: %s\n", err.Error())
		return
	}
	url := fmt.Sprintf("%s://%s/sync/push/%s/%s/%s", getProtocol(), masterURL, gvr.Group, gvr.Version, gvr.Resource)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		logger.Error("post to remote failed: %s\n", err.Error())
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read response failed: %s\n", err.Error())
		return
	}
	pr := sync.PullResult{}
	if err := json.Unmarshal(data, &pr); err != nil {
		logger.Errorf("unmarshal failed: %s\n", err.Error())
		return
	}
	client.syncer.Logger.Infof("remote %s/%s:%s updated: %d created: %d deleted: %d\n", gvr.Group, gvr.Version, gvr.Resource, len(pr.Updated), len(pr.Created), len(pr.Deleted))
}
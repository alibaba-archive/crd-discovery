package sync

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type FetchResult struct {
	Objects []unstructured.Unstructured `json:"objects"`
	Code    int                         `json:"code"`
	Err     string                      `json:"error"`
}

type PullResult struct {
	Created []string `json:"created"`
	Updated []string `json:"updated"`
	Deleted []string `json:"deleted"`
	Code    int      `json:"code"`
	Err     string   `json:"error"`
}

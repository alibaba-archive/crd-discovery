package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func ReadResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(bs))
	}
	return bs, nil
}

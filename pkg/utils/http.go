package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func GetStringOrElse(r *http.Request, key string, el string) string {
	values, ok := r.URL.Query()[key]
	if !ok || len(values) < 1 {
		return el
	}
	return values[0]
}

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
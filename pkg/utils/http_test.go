package utils

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestReadResponse_200(t *testing.T) {
	msg := []byte("hello")
	resp := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(msg)),
	}
	data, err := ReadResponse(&resp, nil)
	require.NoError(t, err)
	require.Equal(t, true, bytes.Equal(msg, data))
}

func TestReadResponse_500(t *testing.T) {
	msg := []byte("server failed")
	resp := http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader(msg)),
	}
	data, err := ReadResponse(&resp, nil)
	require.EqualError(t, err, string(msg))
	require.True(t, nil == data)
}

func TestReadResponse_err(t *testing.T) {
	msg := []byte("nop")
	resp := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(msg)),
	}
	e := errors.New("error happened")
	data, err := ReadResponse(&resp, e)
	require.EqualError(t, err, e.Error())
	require.True(t, nil == data)
}

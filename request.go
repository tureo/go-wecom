package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func HttpPost(url string, contentType string, body io.Reader) (content []byte, err error) {
	r, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	content, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return
}

func HttpGet(url string) (content []byte, err error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	content, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return
}

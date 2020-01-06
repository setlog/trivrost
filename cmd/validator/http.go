package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getFile(url string) ([]byte, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
	if strings.HasPrefix(url, "file://") {
		url = url[7:]
	}
	return ioutil.ReadFile(url)
}

func getHttpHeadResult(url string) (responseCode int, err error) {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	var response *http.Response
	response, err = client.Head(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	return response.StatusCode, err
}

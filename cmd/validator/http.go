package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func getFile(fileUrlString string) ([]byte, error) {
	fileUrl, err := url.Parse(fileUrlString)
	if err == nil && fileUrl.Scheme == "file" {
		fileUrl.Scheme = ""
		fileUrlString = fileUrl.String()
		return ioutil.ReadFile(fileUrlString)
	} else if err != nil || fileUrl.Scheme == "" {
		return ioutil.ReadFile(fileUrlString)
	}

	client := &http.Client{}
	client.Timeout = time.Second * 30
	resp, err := client.Get(fileUrlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received bad status code %s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
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

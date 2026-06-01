package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getFile(fileUrlString string) ([]byte, error) {
	fileUrl, err := url.Parse(fileUrlString)
	if err == nil && fileUrl.Scheme == "file" {
		fileUrl.Scheme = ""
		return os.ReadFile(fileUrl.String())
	} else if err != nil || fileUrl.Scheme == "" {
		return os.ReadFile(fileUrlString)
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
	return io.ReadAll(resp.Body)
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

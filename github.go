package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetFile(url string) (b []byte, err error) {
	if match := regexp.MustCompile(`^https://github.com/([\w_-]+)/([\w_-]+)/blob/[\w_-]+/([\w/._-]+)$`).FindAllStringSubmatch(url, -1); match != nil {
		return getFile(match[0][1], match[0][2], match[0][3])
	}
	return []byte{}, fmt.Errorf("[Error] Invalid URL format: %s", url)
}

func getFile(org string, repo string, path string) (b []byte, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", org, repo, path), nil)
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

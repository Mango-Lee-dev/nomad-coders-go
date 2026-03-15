package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request failed")

type results struct {
	url string
	status string
}

func main() {
	r := make(map[string]string)
	channel := make(chan results)
	urls := []string{
		"https://www.google.com",
		"https://www.naver.com",
		"https://www.daum.net",
		"https://www.yahoo.com",
		"https://www.bing.com",
		"https://www.ask.com",
		"https://www.duckduckgo.com",
		"https://www.yahoo.com",
		"https://www.bing.com",
	}

	for _, url := range urls {
		go hitURL(url, channel)
	}
	for i := 0; i < len(urls); i++ {
		result := <-channel
		r[result.url] = result.status
	}

	for url, status := range r {
		fmt.Println(url, status)
	}
}	

func hitURL(url string, channel chan<- results) {
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	channel <- results{url: url, status: status}
}
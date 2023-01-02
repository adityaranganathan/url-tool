package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

func main() {
	urls := os.Args[1:]

	DisplayURLInfo(GetURLInfo(urls))
}

type URLInfo struct {
	URL      string
	BodySize int
	err      error
}

func GetURLInfo(urls []string) []URLInfo {
	channel := make(chan URLInfo)
	for _, url := range urls {
		go GetResponseBodySize(url, channel)
	}

	var result []URLInfo
	for i := 0; i < len(urls); i++ {
		info := <-channel
		if info.err != nil {
			log.Printf(info.err.Error())
			continue
		}
		result = append(result, info)
	}

	return result
}

func GetResponseBodySize(url string, channel chan URLInfo) {
	resp, err := http.Get(url)
	if err != nil {
		channel <- URLInfo{err: fmt.Errorf("fetching url response of %s: %w", url, err)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		channel <- URLInfo{err: fmt.Errorf("reading response body: %s", err)}
		return
	}

	channel <- URLInfo{
		URL:      url,
		BodySize: len(body),
	}
}

func DisplayURLInfo(data []URLInfo) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].BodySize < data[j].BodySize
	})

	fmt.Println("Size(bytes) \t URL")
	for _, info := range data {
		fmt.Println(fmt.Sprintf("%d \t\t %s", info.BodySize, info.URL))
	}
}

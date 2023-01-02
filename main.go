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
}

func GetURLInfo(urls []string) []URLInfo {
	var result []URLInfo
	for _, url := range urls {
		size, err := GetResponseBodySize(url)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		result = append(result, URLInfo{
			URL:      url,
			BodySize: size,
		})
	}

	return result
}

func GetResponseBodySize(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("fetching url response of %s: %w", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response body: %w", err)
	}

	return len(body), nil
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

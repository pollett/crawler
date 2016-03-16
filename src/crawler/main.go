package main

import (
	"crawler/linkparser"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	startUrl  *url.URL
	visited   map[string]bool
	close     = "#CLOSE#"
	userAgent = "Crawler Test, see: https://github.com/pollett/crawler"
)

func usage() {
	fmt.Println("usage: crawler https://example.com")
	flag.PrintDefaults()
	os.Exit(5)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	var err error
	startUrl, err = url.Parse(args[0])
	if err != nil {
		fmt.Println("Bad start url: ", startUrl)
		fmt.Println(err.Error())
		os.Exit(7)
	}

	visited = map[string]bool{
		startUrl.Path: true,
	}

	linkQueue := make(chan string)

	go func() { linkQueue <- startUrl.String() }()

	go waitForEmpty(linkQueue)

	for link := range linkQueue {
		if link == close {
			fmt.Println("Crawl complete")
			return
		}
		found := crawlUri(link)
		found = linkparser.ProcessLinks(found, startUrl.String())
		processNewLinks(found, linkQueue)

	}
}

func crawlUri(uri string) []string {
	fmt.Println("Crawling, ", uri)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Failed to parse ", uri, err.Error())
		os.Exit(6)
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to parse ", uri, err.Error())
		os.Exit(7)
	}
	defer response.Body.Close()

	return linkparser.Parse(response.Body)
}

func processNewLinks(links []string, queue chan string) {
	for _, link := range links {
		linkobj, err := url.Parse(link)
		if !visited[linkobj.Path] {
			visited[linkobj.Path] = true
			if err != nil {
				fmt.Println("Bad link, ", link)
				continue
			}

			if linkobj.Host == startUrl.Host {
				req, err := http.NewRequest("HEAD", link, nil)
				if err != nil {
					fmt.Println("Failed to parse ", link, err.Error())
					os.Exit(6)
				}
				req.Header.Set("User-Agent", userAgent)

				client := &http.Client{}
				headers, err := client.Do(req)
				if err != nil {
					fmt.Println("Failed to parse ", link, err.Error())
					os.Exit(7)
				}
				if err != nil {
					fmt.Println("Unable to check link, ", link)
					continue
				}

				contentType := headers.Header.Get("Content-Type")

				if strings.Contains(contentType, "text/html") {
					go func() { queue <- linkobj.String() }()
				} else {
					fmt.Println("Resource, ", link)
				}
			} else {
				fmt.Println("External, ", link)
			}
		}
	}
}

func waitForEmpty(queue chan string) {
	empty := 0
	tick := time.Tick(time.Second)
	for range tick {
		if len(queue) == 0 {
			empty++
		} else {
			empty = 0
		}
		if empty >= 3 {
			queue <- close
		}
	}
}

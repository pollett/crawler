package main

import (
	"crawler/linkparser"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	startUrl *url.URL
	visited  map[string]bool
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

	for link := range linkQueue {
		found := crawlUri(link)
		found = linkparser.ProcessLinks(found, startUrl.String())
		processNewLinks(found, linkQueue)

	}
}

func crawlUri(uri string) []string {
	fmt.Println("Crawling ", uri)
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println("Failed to parse ", uri, err.Error())
		os.Exit(6)
	}
	defer response.Body.Close()

	return linkparser.Parse(response.Body)
}

func processNewLinks(links []string, queue chan string) {
	for _, link := range links {
		linkobj, err := url.Parse(link)
		if !visited[linkobj.Path] {
			if err == nil && linkobj.Host == startUrl.Host {
				go func() { queue <- linkobj.String() }()
			}
			visited[linkobj.Path] = true
		}
	}
}

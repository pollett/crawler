package main

import (
	"crawler/linkparser"
	"flag"
	"fmt"
	"net/http"
	"os"
	"net/url"
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
	crawlStartUri, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("Bad start url: ", crawlStartUri)
		fmt.Println(err.Error())
		os.Exit(7)
	}

	linkQueue := make(chan string)

	go func() { linkQueue <- crawlStartUri.String() }()

	for link := range linkQueue {
		found := crawlUri(link)
		found = linkparser.ProcessLinks(found, crawlStartUri.String())
		processNewLinks(found, linkQueue)

	}
}

func crawlUri(uri string) []string {
	fmt.Println("Crawling ",uri)
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
		go func() { queue <- link }()
	}
}
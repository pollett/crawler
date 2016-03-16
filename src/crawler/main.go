package main

import (
	"crawler/linkparser"
	"flag"
	"fmt"
	"net/http"
	"os"
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

	links := crawlUri(args[0])

	fmt.Println(links)
}

func crawlUri(uri string) []string {
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println("Failed to parse ", uri, err.Error())
		os.Exit(6)
	}
	defer response.Body.Close()

	return linkparser.Parse(response.Body)
}

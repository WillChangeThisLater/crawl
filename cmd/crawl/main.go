package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/WillChangeThisLater/crawl"
)

func main() {
	// Define the flags
	maxDepth := flag.Int("d", 1, "Maximum depth to crawl (use -1 for no limit)")
	silence := flag.Bool("s", false, "Silence stderr output")
	maxConns := flag.Int("c", 25, "Maximum number of concurrent requests")
	timeout := flag.Int("t", 5, "Request timeout")

	// Parse the flags
	flag.Parse()

	// Ensure that exactly one non-flag argument is provided, which is the URL
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s [flags] <url>\n", os.Args[0])
	}

	// Retrieve the URL from the arguments
	url := flag.Arg(0)

	// Silence stderr if the -s flag is set
	if *silence {
		log.SetOutput(ioutil.Discard)
	}

	// Crawl the URL
	linksChan := crawl.CrawlSiteForLinks(url, *maxConns, *maxDepth, *timeout)

	// Print discovered links
	for link := range linksChan {
		fmt.Println(link)
	}
}

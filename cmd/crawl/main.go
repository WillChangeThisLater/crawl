package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/WillChangeThisLater/crawl"
)

func main() {
	// Define the input flag
	urlFile := flag.String("input", "", "Input file containing URLs (one per line)")

	// Parse the CLI flags
	flag.Parse()

	var urls []string

	// Check if the input flag is provided
	if *urlFile != "" {
		file, err := os.Open(*urlFile)
		if err != nil {
			log.Fatalf("Failed to open file: %s\n", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := strings.TrimSpace(scanner.Text())
			if url != "" {
				urls = append(urls, url)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Failed to read file: %s\n", err)
		}
	} else {
		// If no input flag, read URLs from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			url := strings.TrimSpace(scanner.Text())
			if url != "" {
				urls = append(urls, url)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Failed to read stdin: %s\n", err)
		}
	}

	// Set max concurrency
	maxConns := 10

	// Crawl each URL
	for _, url := range urls {
		linksChan := crawl.CrawlSiteForLinks(url, maxConns)
		// Print discovered links
		for link := range linksChan {
			fmt.Println(link)
		}
	}
}

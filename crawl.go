package crawl

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	seenLinks      = make(map[string]struct{})
	normalizedSeen = make(map[string]struct{})
	mu             sync.Mutex
)

func getLinksFromHTML(htmlContent io.Reader) []string {
	links := make([]string, 0)
	tokenizer := html.NewTokenizer(htmlContent)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

func normalizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Error parsing URL %s: %v", rawURL, err)
		return rawURL
	}
	parsedURL.Fragment = "" // remove the fragment
	return parsedURL.String()
}

func crawlLinks(wg *sync.WaitGroup, semaphore chan struct{}, link string, discoveredLinks chan<- string, depth, maxDepth int, timeoutSecs int) {
	defer wg.Done()

	mu.Lock()
	if _, ok := seenLinks[link]; ok {
		mu.Unlock()
		return
	}
	seenLinks[link] = struct{}{}
	mu.Unlock()

	if maxDepth != -1 && depth > maxDepth {
		return
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutSecs) * time.Second,
	}

	semaphore <- struct{}{}
	resp, err := client.Get(link)
	<-semaphore

	if err != nil {
		log.Printf("Error fetching %s: %v\n", link, err)
		return
	} else if resp.StatusCode >= 400 {
		log.Printf("%s bad status code %d\n", link, resp.StatusCode)
	}

	finalURL := resp.Request.URL.String()
	mu.Lock()
	normalizedURL := normalizeURL(finalURL)
	if _, exists := normalizedSeen[normalizedURL]; !exists {
		discoveredLinks <- finalURL
		normalizedSeen[normalizedURL] = struct{}{}
	}
	mu.Unlock()

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body of %s: %v\n", finalURL, err)
		return
	}

	baseURL, err := url.Parse(finalURL)
	if err != nil {
		log.Printf("Error parsing base URL %s: %v", finalURL, err)
		return
	}
	childLinks := getLinksFromHTML(bytes.NewReader(bodyBytes))
	for _, childLink := range childLinks {
		parsedChildLink, err := url.Parse(childLink)
		if err != nil {
			log.Printf("Error parsing child link %s: %v", childLink, err)
			continue
		}

		resolvedURL := baseURL.ResolveReference(parsedChildLink)

		if resolvedURL.Host == baseURL.Host {
			wg.Add(1)
			go crawlLinks(wg, semaphore, resolvedURL.String(), discoveredLinks, depth+1, maxDepth, timeoutSecs)
		}
	}
}

func CrawlSiteForLinks(startURL string, maxConns, maxDepth int, timeoutSecs int) <-chan string {
	links := make(chan string)

	go func() {
		var waitGroup sync.WaitGroup
		done := make(chan bool)
		discoveredLinks := make(chan string)
		semaphore := make(chan struct{}, maxConns)

		waitGroup.Add(1)
		go crawlLinks(&waitGroup, semaphore, startURL, discoveredLinks, 0, maxDepth, timeoutSecs)
		go func() {
			for link := range discoveredLinks {
				links <- link
			}
			done <- true
			close(links)
		}()

		waitGroup.Wait()
		close(discoveredLinks)
		<-done
		log.Printf("%s crawler done\n", startURL)
	}()
	return links
}

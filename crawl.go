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
	seenLinks = make(map[string]struct{})
	mu        = new(sync.Mutex)
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

func crawlLinks(wg *sync.WaitGroup, semaphore chan struct{}, link string, discoveredLinks chan<- string, depth, maxDepth int, timeoutSecs int) {
	defer wg.Done()

	mu.Lock()
	if _, ok := seenLinks[link]; ok {
		mu.Unlock()
		return
	} else {
		seenLinks[link] = struct{}{}
	}
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

	if err != nil {
		log.Printf("%s error %s\n", link, err)
		return
	} else if resp.StatusCode >= 400 {
		log.Printf("%s bad status code %d\n", link, resp.StatusCode)
	}

	discoveredLinks <- link

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	baseURL, err := url.Parse(link)
	if err != nil {
		log.Printf("Error parsing base URL %s: %v", link, err)
		return
	}
	childLinks := getLinksFromHTML(bytes.NewReader(bodyBytes))
	for _, childLink := range childLinks {
		parsedChildLink, err := url.Parse(childLink)
		if err != nil {
			log.Printf("Error parsing child link %s: %v", childLink, err)
			return
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

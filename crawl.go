package crawl

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

var (
	seenLinks = make(map[string]struct{})
	mu        = new(sync.Mutex)
)

//func getSHA256Sum(content io.Reader) string {
//	bodyBytes, err := io.ReadAll(content)
//	if err != nil {
//		panic(err)
//	}
//	return fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
//}

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

func crawlLinks(wg *sync.WaitGroup, semaphore chan struct{}, link string, discoveredLinks chan<- string) {
	defer wg.Done()

	mu.Lock()
	if _, ok := seenLinks[link]; ok {
		log.Printf("%s link seen before - skipping\n", link)
		mu.Unlock()
		return
	} else {
		seenLinks[link] = struct{}{}
	}
	mu.Unlock()

	semaphore <- struct{}{}
	log.Printf("%s - GET request\n", link)
	resp, err := http.Get(link)
	<-semaphore

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
			log.Printf("%s Kicking off crawl for %s\n", link, resolvedURL.String())
			wg.Add(1)
			go crawlLinks(wg, semaphore, resolvedURL.String(), discoveredLinks)
		}
	}
}

func CrawlSiteForLinks(startURL string, maxConns int) []string {

	var waitGroup sync.WaitGroup
	done := make(chan bool)
	semaphore := make(chan struct{}, maxConns)
	discoveredLinks := make(chan string)
	linksToReturn := make([]string, 0)

	waitGroup.Add(1)
	go crawlLinks(&waitGroup, semaphore, startURL, discoveredLinks)
	go func() {
		for link := range discoveredLinks {
			linksToReturn = append(linksToReturn, link)
		}
		done <- true
	}()

	waitGroup.Wait()
	close(discoveredLinks)
	<-done
	log.Printf("%s crawler done\n", startURL)
	return linksToReturn
}

package crawl

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetLinksFromHTML(t *testing.T) {
	testCases := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "No links",
			html:     "<p>hey there</p>",
			expected: []string{},
		},
		{
			name:     "One link",
			html:     "<a href='www.google.com'>Google</a>",
			expected: []string{"www.google.com"},
		},
		{
			name:     "Multiple links",
			html:     "<a href='www.google.com'>Google</a><a href='www.example.com'>Example</a>",
			expected: []string{"www.google.com", "www.example.com"},
		},
		{
			name:     "Bad href tag",
			html:     "<a hrefs='www.google.com'>Google</p>",
			expected: []string{},
		},
		{
			name:     "Nested links",
			html:     "<div><a href='www.google.com'>Google</a></div>",
			expected: []string{"www.google.com"},
		},
		{
			name:     "Different types of links",
			html:     "<a href='www.google.com'>Google</a><a href='/relative/link'>Relative Link</a>",
			expected: []string{"www.google.com", "/relative/link"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			html := bytes.NewBufferString(tc.html)
			links := getLinksFromHTML(html)
			if !reflect.DeepEqual(links, tc.expected) {
				t.Errorf("getLinksFromHTML(%s) = %v, want %v", tc.html, links, tc.expected)
			}
		})
	}
}

// https://stackoverflow.com/questions/56861677/synchronizing-a-test-server-during-tests
func waitForServerToStart(port int) {
	backoff := 50 * time.Millisecond

	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), 1*time.Second)
		if err != nil {
			time.Sleep(backoff)
			continue
		}
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Fatalf("Server on port %d not up after 10 attempts", port)
}

func setupServer(site string) *http.Server {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.Dir(site)).ServeHTTP(w, r)
		}),
		Addr: listener.Addr().String(),
	}
	go srv.Serve(listener)
	return srv
}

func ChanToArr(channel <-chan string) []string {
	arr := make([]string, 0)
	for e := range channel {
		arr = append(arr, e)
	}
	return arr
}

func TestCrawlLinksBigSite(t *testing.T) {
	site := "./test-sites/eli.thegreenplace.net/"
	srv := setupServer(site)
	_, portStr, err := net.SplitHostPort(srv.Addr)
	if err != nil {
		log.Fatalf("Failed to parse server address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Failed to convert port string to int: %v", err)
	}
	defer srv.Close()

	fmt.Printf("test site running on port %d\n", port)
	waitForServerToStart(port)

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	linksChannel := CrawlSiteForLinks(baseURL, 10)
	links := ChanToArr(linksChannel)

	seenLinks := make(map[string]struct{})

	if len(links) == 0 {
		t.Errorf("No links found for %s", baseURL)
	}
	for _, link := range links {
		if _, ok := seenLinks[link]; ok {
			t.Errorf("Link %s is listed at least twice", link)
		} else {
			seenLinks[link] = struct{}{}
		}
	}
}

func TestCrawlLinksSmallSite(t *testing.T) {
	site := "./test-sites/sample-site/"
	srv := setupServer(site)
	_, portStr, err := net.SplitHostPort(srv.Addr)
	if err != nil {
		log.Fatalf("Failed to parse server address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Failed to convert port string to int: %v", err)
	}
	defer srv.Close()

	fmt.Printf("test site running on port %d\n", port)
	waitForServerToStart(port)

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	linksChannel := CrawlSiteForLinks(baseURL, 10)
	links := ChanToArr(linksChannel)
	if len(links) != 4 {
		t.Errorf("Expected exactly four links from %s; got %d", baseURL, len(links))
	}
}

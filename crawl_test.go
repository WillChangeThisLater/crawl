package crawl

import (
	"bytes"
	"reflect"
	"testing"
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

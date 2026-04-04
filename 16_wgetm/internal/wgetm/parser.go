package wgetm

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type link struct {
	url    string
	isPage bool
}

var tagAttrs = []struct {
	tag    string
	attr   string
	isPage bool
}{
	{"a", "href", true},
	{"link", "href", false},
	{"script", "src", false},
	{"img", "src", false},
	{"source", "src", false},
	{"video", "src", false},
	{"audio", "src", false},
	{"iframe", "src", false},
	{"embed", "src", false},
	{"object", "data", false},
}

func extractLinks(body []byte, base *url.URL) []link {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var links []link

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, rule := range tagAttrs {
				if n.Data != rule.tag {
					continue
				}
				for _, a := range n.Attr {
					if a.Key != rule.attr {
						continue
					}
					abs := resolveURL(base, a.Val)
					if abs != "" && !seen[abs] {
						seen[abs] = true
						links = append(links, link{url: abs, isPage: rule.isPage})
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

func rewriteLinks(body []byte, base *url.URL, urlToLocal map[string]string) []byte {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return body
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, a := range n.Attr {
				for _, rule := range tagAttrs {
					if n.Data != rule.tag || a.Key != rule.attr {
						continue
					}
					abs := resolveURL(base, a.Val)
					if local, ok := urlToLocal[abs]; ok {
						n.Attr[i].Val = local
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return body
	}
	return buf.Bytes()
}

func resolveURL(base *url.URL, rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" ||
		strings.HasPrefix(rawURL, "#") ||
		strings.HasPrefix(rawURL, "javascript:") ||
		strings.HasPrefix(rawURL, "mailto:") ||
		strings.HasPrefix(rawURL, "tel:") ||
		strings.HasPrefix(rawURL, "data:") {
		return ""
	}
	ref, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	abs := base.ResolveReference(ref)
	abs.Fragment = ""
	if abs.Scheme != "http" && abs.Scheme != "https" {
		return ""
	}
	return abs.String()
}

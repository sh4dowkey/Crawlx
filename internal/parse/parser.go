package parse

import (
	"log"
	"net/url"

	"golang.org/x/net/html"
)

// ExtractLinks recursively finds all links from <a> tags.
func ExtractLinks(n *html.Node) map[string]bool {
	var links = make(map[string]bool)
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links[a.Val] = false
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		linksFromChild := ExtractLinks(c)
		for key := range linksFromChild {
			links[key] = false
		}
	}
	return links
}

// ResolveLinks converts a map of raw links into a slice of absolute URLs.
func ResolveLinks(link map[string]bool, originalURL string) []string {
	var absoluteLinks []string
	baseURL, err := url.Parse(originalURL)
	if err != nil {
		log.Printf("[WARNING] Failed to parse base URL %s: %v\n", originalURL, err)
		return nil
	}
	for l := range link {
		parsedURL, err := url.Parse(l)
		if err != nil {
			log.Printf("[WARNING] Failed to parse URL: %s, Error: %v\n", l, err)
			continue
		}
		resolvedURL := baseURL.ResolveReference(parsedURL)
		absoluteLinks = append(absoluteLinks, resolvedURL.String())
	}
	return absoluteLinks
}

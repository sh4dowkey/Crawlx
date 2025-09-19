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
func ResolveLinks(link map[string]bool, originalURL string) ([]string, []string) {

	var absoluteLinks []string
	var failedLinks []string

	baseURL, err := url.Parse(originalURL)

	if err != nil {
		log.Printf("[WARNING] Failed to parse base URL %s: %v\n", originalURL, err)
		failedLinks = append(failedLinks, originalURL)
	}

	for l := range link {
		parsedURL, parseErr := url.Parse(l)
		if parseErr != nil {
			failedLinks = append(failedLinks, l)
			continue
		}

		// If the link is absolute (has its own host), we can use it directly.
		if parsedURL.Host != "" {
			absoluteLinks = append(absoluteLinks, parsedURL.String())
			continue
		}

		// If the link is relative, we can only resolve it if the baseURL was valid.
		if err == nil {
			resolvedURL := baseURL.ResolveReference(parsedURL)
			absoluteLinks = append(absoluteLinks, resolvedURL.String())
		} else {
			// The baseURL was invalid, so we cannot resolve this relative link.
			failedLinks = append(failedLinks, l)
		}
	}

	return absoluteLinks, failedLinks
}

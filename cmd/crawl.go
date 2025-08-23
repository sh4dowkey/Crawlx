package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// ANSI escape codes for coloring
const (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorReset  = "\033[0m"
)

var (
	visited       = make(map[string]bool)
	mu            sync.Mutex
	verbose       bool
	totalCrawled  int
	totalWarnings int
	totalErrors   int
)

// Global http.Client for reusable connections with timeout.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// init() is a special function that runs before `main`. It's a good place for flag definitions.
func init() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output with more details")
}

// parseAndValidateFlags handles all flag-related logic.
func parseAndValidateFlags() (string, int) {
	urlPtr := flag.String("url", "", "The URL where the crawler will start crawling")
	flag.StringVar(urlPtr, "u", "", "The URL where the crawler will start crawling")

	depthPtr := flag.Int("depth", 2, "The depth the crawler crawls the website ")
	flag.IntVar(depthPtr, "d", 2, "The depth the crawler crawls the website")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Error: The --url flag is required.")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	if *depthPtr < 0 {
		fmt.Println("Error: Crawl depth cannot be a negative value. Exiting...............................")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	return *urlPtr, *depthPtr
}

// checkAndResolveURL fetches a web page and returns its links and HTTP status code.
func checkAndResolveURL(targetURL string) ([]string, int) {
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		log.Printf("[WARNING] Unsupported protocol scheme: %s\n", targetURL)
		return nil, 0
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Printf("[WARNING] Failed to create request for URL %s: %v\n", targetURL, err)
		return nil, 0
	}
	// Set a User-Agent header
	req.Header.Set("User-Agent", "Go-Web-Crawler/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[WARNING] Failed to fetch URL %s: %v\n", targetURL, err)
		return nil, 0
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 && res.StatusCode < 400 {
		log.Printf("[WARNING] Redirect detected for %s (Status: %d)\n", targetURL, res.StatusCode)
		return nil, res.StatusCode
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		log.Printf("[WARNING] Failed to parse HTML from %s: %v\n", targetURL, err)
		return nil, res.StatusCode
	}
	linkMap := extractLinks(body)
	links := resolveLinks(linkMap, targetURL)
	return links, res.StatusCode
}

// extractLinks recursively finds all links from <a> tags.
func extractLinks(n *html.Node) map[string]bool {
	var links = make(map[string]bool)
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links[a.Val] = false
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		linksFromChild := extractLinks(c)
		for key := range linksFromChild {
			links[key] = false
		}
	}
	return links
}

// resolveLinks converts a map of raw links into a slice of absolute URLs.
func resolveLinks(link map[string]bool, originalURL string) []string {
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

// isSameDomain checks if the host of a link belongs to the base domain.
func isSameDomain(linkHost, baseHost string) bool {
	// The base host is a subdomain of itself.
	if linkHost == baseHost {
		return true
	}
	// Check if the link host ends with the base host.
	// The `.` is important to prevent matching `sub.domain.com` with `maindomain.com`
	return strings.HasSuffix(linkHost, "."+baseHost)
}

// Crawl is the main recursive function that performs the web crawl.
// It now takes the original startURL as an argument for the host check.
func Crawl(targetURL *url.URL, startURL *url.URL, initialDepth, currentDepth int) {
	mu.Lock()
	if visited[targetURL.String()] || currentDepth < 0 {
		mu.Unlock()
		return
	}
	visited[targetURL.String()] = true
	mu.Unlock()

	indent := strings.Repeat("    ", initialDepth-currentDepth)

	links, statusCode := checkAndResolveURL(targetURL.String())

	if statusCode >= 400 {
		totalErrors++
	} else if statusCode >= 300 {
		totalWarnings++
	} else {
		totalCrawled++
	}

	var statusColor string
	if statusCode >= 400 {
		statusColor = ColorRed
	} else if statusCode >= 300 {
		statusColor = ColorYellow
	} else {
		statusColor = ColorGreen
	}

	fmt.Printf("%s[+] Crawling: %s (Depth %d)\n", indent, targetURL, initialDepth-currentDepth)
	fmt.Printf("%s    %s[%d %s]%s Found %d links.\n", indent, statusColor, statusCode, http.StatusText(statusCode), ColorReset, len(links))

	if verbose {
		for _, link := range links {
			fmt.Printf("%s        Found Link: %s\n", indent, link)
		}
	}

	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		// Corrected check: use the new isSameDomain function.
		if isSameDomain(parsedLink.Host, startURL.Host) {
			Crawl(parsedLink, startURL, initialDepth, currentDepth-1)
		}
	}
}

// printSummaryAndLinks displays the final summary and the list of visited links.
func printSummaryAndLinks(duration time.Duration) {
	fmt.Println("\n--- Crawl Complete ---")
	fmt.Println()
	fmt.Println("--- Crawl Summary ---")
	fmt.Printf("- Crawled %d pages in %s.\n", totalCrawled, duration.Round(time.Millisecond))
	fmt.Printf("- %sSuccess: %d pages%s\n", ColorGreen, totalCrawled, ColorReset)
	fmt.Printf("- %sWarnings: %d pages%s\n", ColorYellow, totalWarnings, ColorReset)
	fmt.Printf("- %sErrors: %d pages%s\n", ColorRed, totalErrors, ColorReset)

	fmt.Println("\n--- Visited Links ---")
	for url := range visited {
		fmt.Printf("- Visited: %s\n", url)
	}
}

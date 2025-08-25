package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"
)

var (
	visited         = make(map[string]bool)
	mu              sync.Mutex
	verbose         bool
	totalCrawled    int
	totalWarnings   int
	totalErrors     int
	redirects       = make(map[string]bool)
	externalLinks   = make(map[string]bool)
	successfulPages = make(map[string]time.Duration)
	allVisitedLinks = make(map[string]bool)
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// checkAndResolveURL fetches a web page and returns its links and HTTP status code.
func checkAndResolveURL(targetURL string) ([]string, int, time.Duration) {
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		log.Printf("[WARNING] Unsupported protocol scheme: %s\n", targetURL)
		return nil, 0, 0
	}

	startTime := time.Now()
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Printf("[WARNING] Failed to create request for URL %s: %v\n", targetURL, err)
		return nil, 0, 0
	}
	req.Header.Set("User-Agent", "Go-Web-Crawler/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[WARNING] Failed to fetch URL %s: %v\n", targetURL, err)
		return nil, 0, 0
	}
	defer res.Body.Close()
	duration := time.Since(startTime)

	if res.StatusCode >= 300 && res.StatusCode < 400 {
		return nil, res.StatusCode, duration
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return nil, res.StatusCode, duration
	}
	linkMap := extractLinks(body)
	links := resolveLinks(linkMap, targetURL)
	return links, res.StatusCode, duration
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
	if linkHost == baseHost {
		return true
	}
	return strings.HasSuffix(linkHost, "."+baseHost)
}

// Crawl is the main recursive function that performs the web crawl.
func Crawl(targetURL *url.URL, startURL *url.URL, initialDepth, currentDepth int) {
	mu.Lock()
	if visited[targetURL.String()] || currentDepth < 0 {
		mu.Unlock()
		return
	}
	visited[targetURL.String()] = true
	allVisitedLinks[targetURL.String()] = true
	mu.Unlock()

	links, statusCode, duration := checkAndResolveURL(targetURL.String())

	mu.Lock()
	if statusCode >= 400 {
		totalErrors++
	} else if statusCode >= 300 {
		totalWarnings++
		redirects[targetURL.String()] = true
	} else if statusCode >= 200 && statusCode < 300 {
		totalCrawled++
		successfulPages[targetURL.String()] = duration
	}
	mu.Unlock()

	indent := strings.Repeat("  ", initialDepth-currentDepth)
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "Unknown Status"
	}
	statusColor := getColor(statusCode)

	if verbose {
		fmt.Printf("%s[+] Crawling: %s (Depth %d)\n", indent, targetURL, initialDepth-currentDepth)
		fmt.Printf("%s  ↳ [%s%d %s%s] Found %d links.\n", indent, statusColor, statusCode, statusText, ColorReset, len(links))
		if len(links) > 0 {
			fmt.Printf("%s  Links found: \n", indent)
			for _, link := range links {
				fmt.Printf("%s    - %s\n", indent, link)
			}
		}
	} else {
		// Non-verbose output with color applied to the status part.
		fmt.Printf(" [%s%d %s%s] (%dms) %s\n", statusColor, statusCode, statusText, ColorReset, duration.Milliseconds(), targetURL)
	}

	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		if !isSameDomain(parsedLink.Host, startURL.Host) {
			mu.Lock()
			if parsedLink.Scheme == "mailto" || parsedLink.Scheme == "https" || parsedLink.Scheme == "http" {
				externalLinks[link] = true
			}
			mu.Unlock()
		} else {
			mu.Lock()
			if !visited[link] {
				mu.Unlock()
				Crawl(parsedLink, startURL, initialDepth, currentDepth-1)
			} else {
				mu.Unlock()
			}
		}
	}
}

// getColor returns the appropriate ANSI color code for a given status code.
func getColor(statusCode int) string {
	if statusCode >= 400 {
		return ColorRed
	} else if statusCode >= 300 {
		return ColorYellow
	} else {
		return ColorGreen
	}
}

// printSummaryAndLinks displays the final summary and the list of visited links.
func printSummaryAndLinks(duration time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%s\n%s CRAWL SITEMAP REPORT\n%s\n\n", ColorGreen, strings.Repeat(" ", 18), ColorReset)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n[+] Crawled Pages (200 OK)")
	for url, dur := range successfulPages {
		fmt.Printf("  - %s (%dms)  %s%s\n", url, dur.Milliseconds(), ColorGreen, ColorReset)
	}

	fmt.Println("\n\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[~] Redirects (3xx)")
	if len(redirects) == 0 {
		fmt.Println("  No redirects found.")
	} else {
		for url := range redirects {
			fmt.Printf("  - %s\n", url)
		}
	}

	fmt.Println("\n\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[✗] Client & Server Errors (4xx/5xx)")
	if totalErrors == 0 {
		fmt.Println("  No errors found.")
	} else {
		// Placeholder for detailed error reporting if needed later.
	}

	fmt.Println("\n\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[!] External Links")
	if len(externalLinks) == 0 {
		fmt.Println("  No external links found.")
	} else {
		for link := range externalLinks {
			fmt.Printf("  - %s\n", link)
		}
	}

	fmt.Println("\n\n" + strings.Repeat("=", 60) + "\n")
	fmt.Println("Crawl Summary")
	fmt.Printf("  - Crawled %d pages in %s.\n", totalCrawled, duration.Round(time.Millisecond))
	fmt.Printf("  - %sSuccess:%s %d\n", ColorGreen, ColorReset, totalCrawled)
	fmt.Printf("  - %sWarnings:%s %d\n", ColorYellow, ColorReset, totalWarnings)
	fmt.Printf("  - %sErrors:%s %d\n", ColorRed, ColorReset, totalErrors)

	fmt.Println("\n\n" + strings.Repeat("=", 60) + "\n")
	fmt.Println("All Visited Links")
	for url := range allVisitedLinks {
		if _, ok := successfulPages[url]; ok {
			fmt.Printf("  - [%s✓%s] %s\n", ColorGreen, ColorReset, url)
		} else if _, ok := redirects[url]; ok {
			fmt.Printf("  - [%s~%s] %s\n", ColorYellow, ColorReset, url)
		} else {
			fmt.Printf("  - [%s✗%s] %s\n", ColorRed, ColorReset, url)
		}
	}
	fmt.Println("\n" + strings.Repeat("=", 60))
}

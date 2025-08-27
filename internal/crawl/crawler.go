package crawl

import (
	"CRAWLER/internal/parse"
	"CRAWLER/internal/util"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	visited         = make(map[string]bool)
	mu              sync.Mutex
	Verbose         bool // Now exported so main can set it
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
	linkMap := parse.ExtractLinks(body)
	links := parse.ResolveLinks(linkMap, targetURL)
	return links, res.StatusCode, duration
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
	statusColor := util.GetColor(statusCode)

	if Verbose {
		fmt.Printf("%s[+] Crawling: %s (Depth %d)\n", indent, targetURL, initialDepth-currentDepth)
		fmt.Printf("%s  ↳ [%s%d %s%s] Found %d links.\n", indent, statusColor, statusCode, statusText, util.ColorReset, len(links))
		if len(links) > 0 {
			fmt.Printf("%s  Links found: \n", indent)
			for _, link := range links {
				fmt.Printf("%s    - %s\n", indent, link)
			}
		}
	} else {
		fmt.Printf(" [%s%d %s%s] (%dms) %s\n", statusColor, statusCode, statusText, util.ColorReset, duration.Milliseconds(), targetURL)
	}

	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		if !util.IsSameDomain(parsedLink.Host, startURL.Host) {
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

// PrintSummaryAndLinks displays the final summary and the list of visited links.
func PrintSummaryAndLinks(duration time.Duration) {
	// ... (This function remains the same, just ensure it uses util.Color constants)
	// For brevity, I've omitted the body, but you should move the original function here.
	// Make sure to replace all Color constants with util.Color...
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%s\n%s CRAWL SITEMAP REPORT\n%s\n\n", util.ColorGreen, strings.Repeat(" ", 18), util.ColorReset)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n[+] Crawled Pages (200 OK)")
	for url, dur := range successfulPages {
		fmt.Printf("  - %s (%dms)  %s%s\n", url, dur.Milliseconds(), util.ColorGreen, util.ColorReset)
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
	fmt.Printf("  - %sSuccess:%s %d\n", util.ColorGreen, util.ColorReset, totalCrawled)
	fmt.Printf("  - %sWarnings:%s %d\n", util.ColorYellow, util.ColorReset, totalWarnings)
	fmt.Printf("  - %sErrors:%s %d\n", util.ColorRed, util.ColorReset, totalErrors)

	fmt.Println("\n\n" + strings.Repeat("=", 60) + "\n")
	fmt.Println("All Visited Links")
	for url := range allVisitedLinks {
		if _, ok := successfulPages[url]; ok {
			fmt.Printf("  - [%s✓%s] %s\n", util.ColorGreen, util.ColorReset, url)
		} else if _, ok := redirects[url]; ok {
			fmt.Printf("  - [%s~%s] %s\n", util.ColorYellow, util.ColorReset, url)
		} else {
			fmt.Printf("  - [%s✗%s] %s\n", util.ColorRed, util.ColorReset, url)
		}
	}
	fmt.Println("\n" + strings.Repeat("=", 60))
}

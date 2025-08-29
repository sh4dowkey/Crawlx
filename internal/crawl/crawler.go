package crawl

import (
	"CRAWLER/internal/parse"
	"CRAWLER/internal/util"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// CrawlError struct to hold categorized error information.
type CrawlError struct {
	Type    string // "Network" or "HTTP"
	Message string
}

var (
	visited         = make(map[string]bool)
	mu              sync.Mutex
	Verbose         bool
	totalCrawled    int
	totalWarnings   int
	totalErrors     int
	redirects       = make(map[string]bool)
	externalLinks   = make(map[string]bool)
	successfulPages = make(map[string]time.Duration)
	allVisitedLinks = make(map[string]bool)
	// The map now stores the new CrawlError struct.
	failedPages = make(map[string]CrawlError)
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// checkAndResolveURL returns a standard error if a network/parsing issue occurs.
func checkAndResolveURL(targetURL string) ([]string, int, time.Duration, error) {
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		return nil, 0, 0, fmt.Errorf("unsupported protocol scheme")
	}

	startTime := time.Now()
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Go-Web-Crawler/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer res.Body.Close()
	duration := time.Since(startTime)

	if res.StatusCode >= 300 && res.StatusCode < 400 {
		return nil, res.StatusCode, duration, nil
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return nil, res.StatusCode, duration, fmt.Errorf("failed to parse HTML: %w", err)
	}
	linkMap := parse.ExtractLinks(body)
	links := parse.ResolveLinks(linkMap, targetURL)
	return links, res.StatusCode, duration, nil
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

	links, statusCode, duration, err := checkAndResolveURL(targetURL.String())

	// Handle and categorize network or parsing errors.
	if err != nil {
		mu.Lock()
		failedPages[targetURL.String()] = CrawlError{Type: "Network", Message: err.Error()}
		totalErrors++
		mu.Unlock()
		fmt.Printf(" [%sFAIL%s] %s (%s)\n", util.ColorRed, util.ColorReset, targetURL.String(), err.Error())
		return
	}

	mu.Lock()
	// Handle and categorize HTTP errors.
	if statusCode >= 400 {
		totalErrors++
		failedPages[targetURL.String()] = CrawlError{
			Type:    "HTTP",
			Message: fmt.Sprintf("Status %d %s", statusCode, http.StatusText(statusCode)),
		}
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
			if !externalLinks[link] {
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

// PrintSummaryAndLinks now prints two distinct sections for errors.
func PrintSummaryAndLinks(duration time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%s\n%s CRAWL SITEMAP REPORT\n%s\n\n", util.ColorGreen, strings.Repeat(" ", 18), util.ColorReset)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n[+] Crawled Pages (200 OK)")
	if len(successfulPages) == 0 {
		fmt.Println("  No pages were crawled successfully.")
	} else {
		for url, dur := range successfulPages {
			fmt.Printf("  - %s (%dms)\n", url, dur.Milliseconds())
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[~] Redirects (3xx)")
	if len(redirects) == 0 {
		fmt.Println("  No redirects found.")
	} else {
		for url := range redirects {
			fmt.Printf("  - %s\n", url)
		}
	}

	// New, separated error sections.
	httpErrors := false
	networkErrors := false

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[✗] Client & Server Errors (4xx/5xx)")
	for url, errInfo := range failedPages {
		if errInfo.Type == "HTTP" {
			fmt.Printf("  - %s (%s)\n", url, errInfo.Message)
			httpErrors = true
		}
	}
	if !httpErrors {
		fmt.Println("  No server-side errors encountered.")
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[!] Network & Parsing Failures")
	for url, errInfo := range failedPages {
		if errInfo.Type == "Network" {
			fmt.Printf("  - %s (%s)\n", url, errInfo.Message)
			networkErrors = true
		}
	}
	if !networkErrors {
		fmt.Println("  No network or parsing failures encountered.")
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	fmt.Println("[!] External Links Found")
	if len(externalLinks) == 0 {
		fmt.Println("  No external links found.")
	} else {
		for link := range externalLinks {
			fmt.Printf("  - %s\n", link)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Println("Crawl Summary")
	fmt.Printf("  - Crawled %d pages in %s.\n", totalCrawled, duration.Round(time.Millisecond))
	fmt.Printf("  - %sSuccess:%s %d\n", util.ColorGreen, util.ColorReset, totalCrawled)
	fmt.Printf("  - %sWarnings:%s %d\n", util.ColorYellow, util.ColorReset, totalWarnings)
	fmt.Printf("  - %sErrors:%s %d\n", util.ColorRed, util.ColorReset, totalErrors)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
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

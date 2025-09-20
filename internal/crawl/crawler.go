package crawl

import (
	"CRAWLER/internal/parse"
	"CRAWLER/internal/util"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// CrawlTask holds the URL and its current depth for our job queue.
type CrawlTask struct {
	URL   *url.URL
	Depth int
}

// crawlResult is a struct for workers to send their findings.
type crawlResult struct {
	task        CrawlTask
	links       []string
	failedLinks []string
	err         error
	statusCode  int
	duration    time.Duration
}

// CrawlError holds categorized error information.
type CrawlError struct {
	Type    string
	Message string
}

// Global variables for the final report.
var (
	Verbose            bool
	mu                 sync.RWMutex
	totalCrawled       int
	totalErrors        int
	successfulPages    = make(map[string]time.Duration)
	allVisitedLinks    = make(map[string]bool)
	failedPages        = make(map[string]CrawlError)
	malformedLinks     = make(map[string]bool)
	externalLinks      = make(map[string]bool)
	initialDepthGlobal int
)

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

// Configuration constants
const (
	NumWorkers     = 10
	TaskBufferSize = 100
	MaxRetries     = 3
	RetryDelay     = time.Second
	RequestDelay   = 100 * time.Millisecond
)

// Crawl is the main setup function with proper concurrency handling.
func Crawl(startURL *url.URL, initialDepth int) {
	initialDepthGlobal = initialDepth

	// Create channels with appropriate buffer sizes
	tasks := make(chan CrawlTask, TaskBufferSize)
	results := make(chan crawlResult, TaskBufferSize)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start workers
	var workerWg sync.WaitGroup
	for i := 0; i < NumWorkers; i++ {
		workerWg.Add(1)
		go worker(ctx, &workerWg, tasks, results)
	}

	// Start dispatcher in a separate goroutine
	var dispatcherWg sync.WaitGroup
	dispatcherWg.Add(1)
	go dispatcher(ctx, &dispatcherWg, tasks, results, startURL)

	// Add the first task to start the crawl
	select {
	case tasks <- CrawlTask{URL: startURL, Depth: initialDepth}:
	case <-ctx.Done():
		return
	}

	// Wait for dispatcher to finish
	dispatcherWg.Wait()

	// Close tasks channel and wait for workers
	close(tasks)
	workerWg.Wait()
	close(results)

	// Drain any remaining results
	for res := range results {
		processResult(res)
	}
}

// worker processes crawl tasks with context for cancellation.
func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan CrawlTask, results chan<- crawlResult) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return // Channel closed
			}

			// Small delay to be polite to servers
			time.Sleep(RequestDelay)

			// Process the task with retries
			links, failedLinks, statusCode, duration, err := checkAndResolveURLWithRetry(task.URL.String())

			select {
			case results <- crawlResult{
				task:        task,
				links:       links,
				failedLinks: failedLinks,
				err:         err,
				statusCode:  statusCode,
				duration:    duration,
			}:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// dispatcher manages the crawl queue and coordinates work distribution.
func dispatcher(ctx context.Context, wg *sync.WaitGroup, tasks chan<- CrawlTask, results <-chan crawlResult, startURL *url.URL) {
	defer wg.Done()

	visited := make(map[string]bool)
	pendingTasks := 1 // We start with one task

	// Mark start URL as visited
	visited[startURL.String()] = true
	mu.Lock()
	allVisitedLinks[startURL.String()] = true
	mu.Unlock()

	for pendingTasks > 0 {
		select {
		case res := <-results:
			pendingTasks-- // One task completed

			// Process the result
			processResult(res)

			// Add new tasks if the crawl was successful and we are not at max depth
			if res.err == nil && res.statusCode < 400 && res.task.Depth > 0 {
				newTasksAdded := 0

				for _, linkStr := range res.links {
					parsedLink, err := url.Parse(linkStr)
					if err != nil {
						mu.Lock()
						malformedLinks[linkStr] = true
						mu.Unlock()
						continue
					}

					if util.IsSameDomain(parsedLink.Host, startURL.Host) {
						// Check if already visited (dispatcher is the only one who checks)
						if !visited[parsedLink.String()] {
							visited[parsedLink.String()] = true
							mu.Lock()
							allVisitedLinks[parsedLink.String()] = true
							mu.Unlock()

							newTask := CrawlTask{URL: parsedLink, Depth: res.task.Depth - 1}

							select {
							case tasks <- newTask:
								newTasksAdded++
							case <-ctx.Done():
								return
							default:
								// Channel full, skip this link to prevent blocking
								if Verbose {
									fmt.Printf("    [WARN] Queue full, skipping: %s\n", parsedLink.String())
								}
							}
						}
					} else {
						mu.Lock()
						externalLinks[linkStr] = true
						mu.Unlock()
					}
				}

				pendingTasks += newTasksAdded
			}

		case <-ctx.Done():
			return
		}
	}
}

// processResult updates the global state with clean output.
func processResult(res crawlResult) {
	mu.Lock()
	defer mu.Unlock()

	if res.err != nil {
		failedPages[res.task.URL.String()] = CrawlError{
			Type:    "Network Error",
			Message: fmt.Sprintf("Connection failed: %s", res.err.Error())}
		totalErrors++
		fmt.Printf(" [%sFAIL%s] %s (Network: %s)\n",
			util.ColorRed, util.ColorReset, res.task.URL.String(), res.err.Error())
		return
	}

	// Handle malformed links
	if len(res.failedLinks) > 0 {
		for _, f := range res.failedLinks {
			malformedLinks[f] = true
		}
	}

	statusText := http.StatusText(res.statusCode)
	if statusText == "" {
		statusText = "Unknown Status"
	}

	// Categorize results by status code ranges
	switch {
	case res.statusCode >= 200 && res.statusCode < 300:
		// 2xx Success
		totalCrawled++
		successfulPages[res.task.URL.String()] = res.duration
	case res.statusCode >= 300 && res.statusCode < 400:
		// 3xx Redirects - handle silently, don't add to report sections
		totalCrawled++
		successfulPages[res.task.URL.String()] = res.duration
	case res.statusCode >= 400 && res.statusCode < 500:
		// 4xx Client Errors
		totalErrors++
		failedPages[res.task.URL.String()] = CrawlError{
			Type:    "Client Error",
			Message: fmt.Sprintf("%d %s", res.statusCode, statusText)}
	case res.statusCode >= 500 && res.statusCode < 600:
		// 5xx Server Errors
		totalErrors++
		failedPages[res.task.URL.String()] = CrawlError{
			Type:    "Server Error",
			Message: fmt.Sprintf("%d %s", res.statusCode, statusText)}
	default:
		// Unknown status codes (rare but possible)
		totalErrors++
		failedPages[res.task.URL.String()] = CrawlError{
			Type:    "Unknown Status",
			Message: fmt.Sprintf("%d %s", res.statusCode, statusText)}
	}

	// Clean output formatting
	statusColor := util.GetColor(res.statusCode)
	currentDepthLevel := initialDepthGlobal - res.task.Depth
	indent := strings.Repeat("  ", currentDepthLevel)

	if Verbose {
		// Clean indexed format like the original
		fmt.Printf("%s[+] Crawling: %s (Depth %d)\n", indent, res.task.URL.String(), currentDepthLevel)
		fmt.Printf("%s  ↳ [%s%d %s%s] Found %d links (%dms)\n",
			indent, statusColor, res.statusCode, statusText, util.ColorReset,
			len(res.links), res.duration.Milliseconds())

		if len(res.links) > 0 && len(res.links) <= 10 {
			fmt.Printf("%s  Links found:\n", indent)
			for _, link := range res.links {
				fmt.Printf("%s    - %s\n", indent, link)
			}
		} else if len(res.links) > 10 {
			fmt.Printf("%s  Links found: (showing first 10 of %d)\n", indent, len(res.links))
			for i, link := range res.links[:10] {
				fmt.Printf("%s    - %s\n", indent, link)
				if i == 9 {
					fmt.Printf("%s    ... and %d more\n", indent, len(res.links)-10)
				}
			}
		}
	} else {
		// Simple clean output
		fmt.Printf(" [%s%d%s] %s (%dms)\n",
			statusColor, res.statusCode, util.ColorReset,
			res.task.URL.String(), res.duration.Milliseconds())
	}
}

// checkAndResolveURLWithRetry adds retry logic for better reliability.
func checkAndResolveURLWithRetry(targetURL string) ([]string, []string, int, time.Duration, error) {
	var lastErr error
	var totalDuration time.Duration

	for attempt := 0; attempt < MaxRetries; attempt++ {
		links, failedLinks, statusCode, duration, err := checkAndResolveURL(targetURL)
		totalDuration += duration

		if err == nil {
			return links, failedLinks, statusCode, totalDuration, nil
		}

		lastErr = err
		if attempt < MaxRetries-1 {
			time.Sleep(RetryDelay * time.Duration(attempt+1))
		}
	}

	return nil, nil, 0, totalDuration, fmt.Errorf("failed after %d attempts: %w", MaxRetries, lastErr)
}

// checkAndResolveURL fetches and parses a URL.
func checkAndResolveURL(targetURL string) ([]string, []string, int, time.Duration, error) {
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		return nil, nil, 0, 0, fmt.Errorf("unsupported protocol scheme")
	}

	startTime := time.Now()
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Go-Web-Crawler/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer res.Body.Close()
	duration := time.Since(startTime)

	if res.StatusCode >= 300 && res.StatusCode < 400 {
		return nil, nil, res.StatusCode, duration, nil
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return nil, nil, res.StatusCode, duration, fmt.Errorf("failed to parse HTML: %w", err)
	}

	linkMap := parse.ExtractLinks(body)
	links, failedLinks := parse.ResolveLinks(linkMap, targetURL)
	return links, failedLinks, res.StatusCode, duration, nil
}

// PrintSummaryAndLinks with clean, separated reporting.
func PrintSummaryAndLinks(duration time.Duration) {
	mu.RLock()
	defer mu.RUnlock()

	// Clear separation between crawl output and report
	fmt.Printf("\n")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("                    CRAWL REPORT\n")
	fmt.Println(strings.Repeat("=", 60))

	// Summary stats
	fmt.Printf("\nScan completed............ Crawled %d sites in %s seconds\n", totalCrawled, duration.Round(time.Millisecond))
	fmt.Printf("Total: %d | %sSuccess: %d%s | %sErrors: %d%s\n\n",
		len(allVisitedLinks),
		util.ColorGreen, totalCrawled, util.ColorReset,
		util.ColorRed, totalErrors, util.ColorReset)

	// Successful pages (2xx and 3xx - redirects handled transparently)
	if len(successfulPages) > 0 {
		fmt.Printf("Successful Pages (%d):\n", len(successfulPages))
		for pageURL, dur := range successfulPages {
			fmt.Printf("-  %s (%dms)\n", pageURL, dur.Milliseconds())
		}
		fmt.Printf("\n")
	}

	// All errors grouped by type
	clientErrors := make(map[string]CrawlError)
	serverErrors := make(map[string]CrawlError)
	networkErrors := make(map[string]CrawlError)
	unknownErrors := make(map[string]CrawlError)

	for pageURL, errInfo := range failedPages {
		switch errInfo.Type {
		case "Client Error":
			clientErrors[pageURL] = errInfo
		case "Server Error":
			serverErrors[pageURL] = errInfo
		case "Network Error":
			networkErrors[pageURL] = errInfo
		default:
			unknownErrors[pageURL] = errInfo
		}
	}

	// Client errors (4xx)
	if len(clientErrors) > 0 {
		fmt.Printf("Client Errors - 4xx (%d):\n", len(clientErrors))
		for pageURL, errInfo := range clientErrors {
			fmt.Printf("-  %s (%s)\n", pageURL, errInfo.Message)
		}
		fmt.Printf("\n")
	}

	// Server errors (5xx)
	if len(serverErrors) > 0 {
		fmt.Printf("Server Errors - 5xx (%d):\n", len(serverErrors))
		for pageURL, errInfo := range serverErrors {
			fmt.Printf("-  %s (%s)\n", pageURL, errInfo.Message)
		}
		fmt.Printf("\n")
	}

	// Network errors
	if len(networkErrors) > 0 {
		fmt.Printf("Network Errors (%d):\n", len(networkErrors))
		for pageURL, errInfo := range networkErrors {
			fmt.Printf("-  %s (%s)\n", pageURL, errInfo.Message)
		}
		fmt.Printf("\n")
	}

	// Unknown status codes (rare)
	if len(unknownErrors) > 0 {
		fmt.Printf("Unknown Status Codes (%d):\n", len(unknownErrors))
		for pageURL, errInfo := range unknownErrors {
			fmt.Printf("-  %s (%s)\n", pageURL, errInfo.Message)
		}
		fmt.Printf("\n")
	}

	// External links
	if len(externalLinks) > 0 {
		fmt.Printf("External Links (%d):\n", len(externalLinks))
		for link := range externalLinks {
			fmt.Printf("-  %s\n", link)
		}
		fmt.Printf("\n")
	}

	// Malformed links
	if len(malformedLinks) > 0 {
		fmt.Printf("Malformed Links (%d):\n", len(malformedLinks))
		for link := range malformedLinks {
			fmt.Printf("-  %s\n", link)
		}
		fmt.Printf("\n")
	}

	fmt.Println(strings.Repeat("=", 60))
}

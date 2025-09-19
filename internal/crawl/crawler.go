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
	Verbose         bool
	mu              sync.RWMutex
	totalCrawled    int
	totalWarnings   int
	totalErrors     int
	redirects       = make(map[string]bool)
	externalLinks   = make(map[string]bool)
	successfulPages = make(map[string]time.Duration)
	allVisitedLinks = make(map[string]bool)
	failedPages     = make(map[string]CrawlError)
	malformedLinks  = make(map[string]bool)
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
)

// Crawl is the main setup function with proper concurrency handling.
func Crawl(startURL *url.URL, initialDepth int) {
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
									fmt.Printf("[WARNING] Task queue full, skipping: %s\n", parsedLink.String())
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

// processResult updates the global state with proper locking.
func processResult(res crawlResult) {
	mu.Lock()
	defer mu.Unlock()

	if res.err != nil {
		failedPages[res.task.URL.String()] = CrawlError{Type: "Network", Message: res.err.Error()}
		totalErrors++
		fmt.Printf(" [%sFAIL%s] %s (%s)\n", util.ColorRed, util.ColorReset, res.task.URL.String(), res.err.Error())
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

	// Categorize results
	if res.statusCode >= 400 {
		totalErrors++
		failedPages[res.task.URL.String()] = CrawlError{Type: "HTTP", Message: fmt.Sprintf("Status %d %s", res.statusCode, statusText)}
	} else if res.statusCode >= 300 {
		totalWarnings++
		redirects[res.task.URL.String()] = true
	} else {
		totalCrawled++
		successfulPages[res.task.URL.String()] = res.duration
	}

	// Output formatting
	statusColor := util.GetColor(res.statusCode)
	if Verbose {
		fmt.Printf(" [%s%d%s] (Depth %d) %s (%dms)\n",
			statusColor, res.statusCode, util.ColorReset,
			res.task.Depth, res.task.URL.String(), res.duration.Milliseconds())
	} else {
		fmt.Printf(" [%s%d%s] (%dms) %s\n",
			statusColor, res.statusCode, util.ColorReset,
			res.duration.Milliseconds(), res.task.URL.String())
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
			time.Sleep(RetryDelay * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	return nil, nil, 0, totalDuration, fmt.Errorf("failed after %d attempts: %w", MaxRetries, lastErr)
}

// checkAndResolveURL fetches and parses a URL (unchanged from original).
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

// PrintSummaryAndLinks remains the same but with improved thread safety.
func PrintSummaryAndLinks(duration time.Duration) {
	mu.RLock()
	defer mu.RUnlock()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n%s CRAWL SITEMAP REPORT\n\n", strings.Repeat(" ", 18))
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
	fmt.Println("[?] Malformed Links Found (Skipped)")
	if len(malformedLinks) == 0 {
		fmt.Println("  No malformed links were found.")
	} else {
		for link := range malformedLinks {
			fmt.Printf("  - %s\n", link)
		}
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
	fmt.Println("All Links Found")
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

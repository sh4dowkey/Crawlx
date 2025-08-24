package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

// `main` is now the orchestrator, delegating tasks to other functions.
func main() {
	// 1. Parse and validate command-line flags.
	startURLStr, initialDepth := parseAndValidateFlags()

	// 2. Record the start time for the summary.
	startTime := time.Now()

	// 3. Parse the starting URL.
	startURL, err := url.Parse(startURLStr)
	if err != nil {
		log.Fatalf("Error parsing initial URL: %v", err)
	}

	fmt.Printf("[INFO] Starting crawl at: %s\n", time.Now().Format("3:04:05 PM MST"))

	// 4. Start the crawl, passing the startURL twice.
	Crawl(startURL, startURL, initialDepth, initialDepth)

	// 5. Calculate and print the final summary.
	duration := time.Since(startTime)
	printSummaryAndLinks(duration)
	os.Exit(0)
}

// parseAndValidateFlags handles all flag-related logic.
func parseAndValidateFlags() (string, int) {
	urlPtr := flag.String("url", "", "The URL where the crawler will start crawling")
	flag.StringVar(urlPtr, "u", "", "The URL where the crawler will start crawling")

	depthPtr := flag.Int("depth", 2, "The depth the crawler crawls the website ")
	flag.IntVar(depthPtr, "d", 2, "The depth the crawler crawls the website")

	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output with more details")

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

package main

import (
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

	// 4. Start the crawl, passing the startURL twice.
	Crawl(startURL, startURL, initialDepth, initialDepth)

	// 5. Calculate and print the final summary.
	duration := time.Since(startTime)
	printSummaryAndLinks(duration)
	os.Exit(0)
}

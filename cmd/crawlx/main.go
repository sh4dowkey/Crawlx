package main

import (
	"CRAWLER/internal/crawl"
	"CRAWLER/internal/util"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

// Define globals for our flags to make them accessible in customUsage.
var (
	urlStr  string
	depth   int
	verbose bool
)

func main() {
	// We call a new function to set up our flags.
	setupFlags()

	// Set the custom usage message AFTER setting up flags.
	flag.Usage = customUsage

	// Now, parse them.
	flag.Parse()

	// Validate the parsed flags.
	if err := validateFlags(); err != nil {
		fmt.Printf("\n%sError: %v%s\n", util.ColorRed, err, util.ColorReset)
		flag.Usage()
		os.Exit(1)
	}

	crawl.Verbose = verbose

	startTime := time.Now()

	startURL, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("Error parsing initial URL: %v", err)
	}

	fmt.Printf("[INFO] Starting crawl at: %s\n", time.Now().Format("3:04:05 PM MST"))
	crawl.Crawl(startURL, startURL, depth, depth)
	duration := time.Since(startTime)
	crawl.PrintSummaryAndLinks(duration)
	os.Exit(0)
}

// setupFlags defines all command-line flags.
func setupFlags() {
	flag.StringVar(&urlStr, "url", "", "The starting URL to crawl.")
	flag.StringVar(&urlStr, "u", "", "The starting URL to crawl (shorthand).")
	flag.IntVar(&depth, "depth", 2, "The maximum depth for recursive crawling.")
	flag.IntVar(&depth, "d", 2, "The maximum depth for recursive crawling (shorthand).")
	flag.BoolVar(&verbose, "verbose", false, "Enable detailed, verbose output.")
}

// validateFlags checks if the provided flags are valid.
func validateFlags() error {
	if urlStr == "" {
		return fmt.Errorf("the --url flag is required")
	}
	if depth < 0 {
		return fmt.Errorf("crawl depth cannot be a negative value")
	}
	return nil
}

// customUsage prints a professional, clean, and dynamic help message.
func customUsage() {

	fmt.Printf("\n%sCrawlX%s - A fast, concurrent web crawler built in Go.\n", util.ColorBold+util.ColorCyan, util.ColorReset)

	fmt.Printf("\n%sUSAGE:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  crawlx -u <STARTING_URL> [OPTIONS]\n")

	// Manually print flags for a clean, grouped layout.
	fmt.Printf("\n%sOPTIONS:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  -u, --url <string>      The starting URL to crawl. (Required)\n")
	fmt.Printf("  -d, --depth <int>       The maximum depth for recursive crawling. (Default: 2)\n")
	fmt.Printf("      --verbose           Enable detailed, verbose output. (Default: false)\n")

	fmt.Printf("\n%sEXAMPLES:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  %s# Crawl a website with a depth of 3\n", util.ColorWhite)
	fmt.Printf("  crawlx -u https://toscrape.com -d 3\n\n")
	fmt.Printf("  %s# Crawl with verbose output\n", util.ColorWhite)
	fmt.Printf("  crawlx --url https://example.com --verbose\n")
}

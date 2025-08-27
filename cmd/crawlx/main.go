package main

import (
	"CRAWLER/internal/crawl"
	"CRAWLER/internal/util"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Usage = customUsage

	startURLStr, initialDepth, verbose := parseAndValidateFlags()
	crawl.Verbose = verbose // Set the exported variable in the crawl package

	startTime := time.Now()

	startURL, err := url.Parse(startURLStr)
	if err != nil {
		log.Fatalf("Error parsing initial URL: %v", err)
	}

	fmt.Printf("[INFO] Starting crawl at: %s\n", time.Now().Format("3:04:05 PM MST"))

	crawl.Crawl(startURL, startURL, initialDepth, initialDepth)

	duration := time.Since(startTime)
	crawl.PrintSummaryAndLinks(duration)
	os.Exit(0)
}

func customUsage() {
	fmt.Printf("%s\n", strings.Repeat("━", 60))
	fmt.Printf("%sGO CRAWLER%s - A concurrent web crawler in Go\n", util.ColorBold+util.ColorCyan, util.ColorReset)
	fmt.Printf("%s\n", strings.Repeat("━", 60))

	fmt.Printf("\n%sUSAGE%s:\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  %sgo run ./cmd/crawlx [OPTIONS]%s\n", util.ColorWhite, util.ColorReset)

	fmt.Printf("\n%sOPTIONS%s:\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	flag.PrintDefaults()

	fmt.Printf("\n%sEXAMPLES%s:\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  %s# Crawl a website with default depth of 2\n", util.ColorWhite)
	fmt.Printf("  go run ./cmd/crawlx --url https://toscrape.com%s\n\n", util.ColorReset)
	fmt.Printf("  %s# Crawl with a depth of 5 and verbose output\n", util.ColorWhite)
	fmt.Printf("  go run ./cmd/crawlx -u https://toscrape.com -d 5 --verbose%s\n", util.ColorReset)

	fmt.Printf("%s\n", strings.Repeat("━", 60))
}

// parseAndValidateFlags handles all flag-related logic.
func parseAndValidateFlags() (string, int, bool) {
	urlPtr := flag.String("url", "", "The URL where the crawler will start crawling")
	flag.StringVar(urlPtr, "u", "", "The URL where the crawler will start crawling (shorthand)")

	depthPtr := flag.Int("depth", 2, "The depth the crawler crawls the website")
	flag.IntVar(depthPtr, "d", 2, "The depth the crawler crawls the website (shorthand)")

	verbosePtr := flag.Bool("verbose", false, "Enable verbose output with more details")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Error: The --url flag is required.")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	if *depthPtr < 0 {
		fmt.Println("Error: Crawl depth cannot be a negative value. Exiting...")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	return *urlPtr, *depthPtr, *verbosePtr
}

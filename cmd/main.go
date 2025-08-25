package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func customUsage() {
	fmt.Printf("%s\n", strings.Repeat("━", 60))
	fmt.Printf("%sGO CRAWLER%s - A concurrent web crawler in Go\n", ColorBold+ColorCyan, ColorReset)
	fmt.Printf("%s\n", strings.Repeat("━", 60))

	fmt.Printf("\n%sUSAGE%s:\n", ColorBold+ColorGreen, ColorReset)
	fmt.Printf("  %sgo run . [OPTIONS]%s\n", ColorWhite, ColorReset)

	fmt.Printf("\n%sOPTIONS%s:\n", ColorBold+ColorGreen, ColorReset)
	flag.PrintDefaults()

	fmt.Printf("\n%sEXAMPLES%s:\n", ColorBold+ColorGreen, ColorReset)
	fmt.Printf("  %s# Crawl a website with default depth of 2\n", ColorWhite)
	fmt.Printf("  go run . --url https://toscrape.com%s\n\n", ColorReset)
	fmt.Printf("  %s# Crawl with a depth of 5 and verbose output\n", ColorWhite)
	fmt.Printf("  go run . -u https://toscrape.com -d 5 --verbose%s\n", ColorReset)

	fmt.Printf("%s\n", strings.Repeat("━", 60))
}

func main() {

	flag.Usage = customUsage

	startURLStr, initialDepth := parseAndValidateFlags()

	startTime := time.Now()

	startURL, err := url.Parse(startURLStr)
	if err != nil {
		log.Fatalf("Error parsing initial URL: %v", err)
	}

	fmt.Printf("[INFO] Starting crawl at: %s\n", time.Now().Format("3:04:05 PM MST"))

	Crawl(startURL, startURL, initialDepth, initialDepth)

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

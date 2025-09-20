package main

import (
	"CRAWLER/internal/crawl"
	"CRAWLER/internal/util"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Define globals for flags.
var (
	urlStr  string
	depth   int
	verbose bool
	allowIP bool
)

func main() {
	setupFlags()
	flag.Usage = customUsage
	flag.Parse()

	// Validate the URL
	startURL, err := validateURL(urlStr, allowIP)
	if err != nil {
		fmt.Printf("%sError:%s %v\n", util.ColorRed, util.ColorReset, err)
		os.Exit(1)
	}

	// Set verbose mode before crawling
	crawl.Verbose = verbose

	// Simple startup message
	fmt.Printf("Crawling %s (depth: %d)\n", startURL.String(), depth)

	// Start the concurrent crawl
	startTime := time.Now()
	crawl.Crawl(startURL, depth)
	duration := time.Since(startTime)

	// Print summary
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
	flag.BoolVar(&verbose, "v", false, "Enable detailed, verbose output (shorthand).")

	flag.BoolVar(&allowIP, "allow-ip", false, "Allow crawling a host that is an IP address.")
	flag.BoolVar(&allowIP, "i", false, "Allow crawling a host that is an IP address (shorthand).")
}

// validateURL performs all pre-flight checks on the user-provided URL.
func validateURL(rawURL string, ipAllowed bool) (*url.URL, error) {
	if rawURL == "" {
		return nil, errors.New("the --url flag is required")
	}

	if strings.Count(rawURL, "://") > 1 {
		return nil, errors.New("malformed URL: multiple '://' sequences found")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, errors.New("invalid URL scheme, please use http:// or https://")
	}

	if parsedURL.Host == "" {
		return nil, errors.New("invalid URL, the host (domain name) is missing")
	}

	isIP, _ := regexp.MatchString(`^[0-9\.]+$`, parsedURL.Host)

	if isIP {
		if !ipAllowed {
			return nil, errors.New("crawling IP addresses is not allowed, use the --allow-ip flag to enable")
		}
	} else {
		if !strings.Contains(parsedURL.Host, ".") {
			return nil, errors.New("invalid domain name, host must contain a '.'")
		}
		parts := strings.Split(parsedURL.Host, ".")
		if len(parts) < 2 {
			return nil, errors.New("invalid domain name")
		}
		tld := parts[len(parts)-1]
		if _, err := strconv.Atoi(tld); err == nil {
			return nil, fmt.Errorf("invalid TLD '%s', cannot be a number", tld)
		}
	}

	return parsedURL, nil
}

// customUsage prints a clean help message.
func customUsage() {
	progName := filepath.Base(os.Args[0])

	fmt.Printf("\nCrawlX - A fast, concurrent web crawler built in Go.\n")

	fmt.Printf("\nUSAGE:\n")
	fmt.Printf("  %s -u <STARTING_URL> [OPTIONS]\n", progName)

	fmt.Printf("\nOPTIONS:\n")
	fmt.Printf("  -u, --url <string>      The starting URL to crawl (Required)\n")
	fmt.Printf("  -d, --depth <int>       Maximum crawling depth (Default: 2)\n")
	fmt.Printf("  -i, --allow-ip          Allow crawling IP addresses (Default: false)\n")
	fmt.Printf("  -v, --verbose           Show detailed crawl progress (Default: false)\n")

	fmt.Printf("\nEXAMPLES:\n")
	fmt.Printf("  %s -u https://example.com -d 3\n", progName)
	fmt.Printf("  %s -u https://toscrape.com -d 2 --verbose\n", progName)
	fmt.Printf("  %s -u http://127.0.0.1:8000 -d 1 --allow-ip\n", progName)
}

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

	// It returns the parsed URL object if validation succeeds.
	startURL, err := validateURL(urlStr, allowIP)
	if err != nil {
		fmt.Printf("\n%sValidation Error: %v%s\n", util.ColorRed, err, util.ColorReset)
		flag.Usage()
		os.Exit(1)
	}

	crawl.Verbose = verbose

	startTime := time.Now()
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
	flag.BoolVar(&verbose, "v", false, "Enable detailed, verbose output.")

	flag.BoolVar(&allowIP, "allow-ip", false, "Allow crawling a host that is an IP address.")
	flag.BoolVar(&allowIP, "i", false, "Allow crawling a host that is an IP address (shorthand).")
}

// validateURL performs all pre-flight checks on the user-provided URL.
func validateURL(rawURL string, ipAllowed bool) (*url.URL, error) {

	// Check 1: Ensure URL flag was provided.
	if rawURL == "" {
		return nil, errors.New("the --url flag is required")
	}

	// Check 2: Prevent malformed URLs like "https://https://..."
	if strings.Count(rawURL, "://") > 1 {
		return nil, errors.New("malformed URL: multiple '://' sequences found")
	}

	// Check 3: Parse the URL.
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	// Check 4: Validate the scheme.
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, errors.New("invalid URL scheme, please use http:// or https://")
	}

	// Check 5: Ensure host is not empty.
	if parsedURL.Host == "" {
		return nil, errors.New("invalid URL, the host (domain name) is missing")
	}

	// Check 6: Handle IP Addresses vs. Domain Names.
	isIP, _ := regexp.MatchString(`^[0-9\.]+$`, parsedURL.Host)

	if isIP {
		if !ipAllowed {
			return nil, errors.New("crawling IP addresses is not allowed, use the --allow-ip flag to enable")
		}
	} else {
		// Check 7: If it's a domain name, it must contain a dot.
		if !strings.Contains(parsedURL.Host, ".") {
			return nil, errors.New("invalid domain name, host must contain a '.'")
		}
		// Check 8 (Advanced): Ensure the Top-Level Domain (TLD) is not a number.
		parts := strings.Split(parsedURL.Host, ".")
		if len(parts) < 2 {
			return nil, errors.New("invalid domain name") // Should be caught by the check above, but good for safety.
		}
		tld := parts[len(parts)-1]
		if _, err := strconv.Atoi(tld); err == nil {
			return nil, fmt.Errorf("invalid TLD '%s', cannot be a number", tld)
		}
	}

	// If all checks pass, return the parsed URL and no error.
	return parsedURL, nil
}

// customUsage prints a help message.
func customUsage() {
	progName := filepath.Base(os.Args[0])

	fmt.Printf("\n%sCrawlX%s - A fast, concurrent web crawler built in Go.\n", util.ColorBold+util.ColorCyan, util.ColorReset)

	fmt.Printf("\n%sUSAGE:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  %s -u <STARTING_URL> [OPTIONS]\n", progName)

	// Manually print flags for a clean, grouped layout.
	fmt.Printf("\n%sOPTIONS:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  -u, --url <string>      The starting URL to crawl. (Required)\n")
	fmt.Printf("  -d, --depth <int>       The maximum depth for recursive crawling. (Default: 2)\n")
	fmt.Printf("  -i, --allow-ip          Allow crawling a host that is an IP address. (Default: false)\n")
	fmt.Printf("  -v  --verbose           Enable detailed, verbose output. (Default: false)\n")

	fmt.Printf("\n%sEXAMPLES:%s\n", util.ColorBold+util.ColorGreen, util.ColorReset)
	fmt.Printf("  %s# Crawl a website with a depth of 3\n", util.ColorWhite)
	fmt.Printf("  %s -u https://toscrape.com -d 3\n\n", progName)
	fmt.Printf("  %s# Crawl with verbose output\n", util.ColorWhite)
	fmt.Printf("  %s --url https://example.com --verbose\n", progName)
}

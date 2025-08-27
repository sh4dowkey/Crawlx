package util

import "strings"

// ANSI escape codes for coloring
const (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"
)

// IsSameDomain checks if the host of a link belongs to the base domain.
func IsSameDomain(linkHost, baseHost string) bool {
	if linkHost == baseHost {
		return true
	}
	return strings.HasSuffix(linkHost, "."+baseHost)
}

// GetColor returns the appropriate ANSI color code for a given status code.
func GetColor(statusCode int) string {
	if statusCode >= 400 {
		return ColorRed
	} else if statusCode >= 300 {
		return ColorYellow
	}
	return ColorGreen
}

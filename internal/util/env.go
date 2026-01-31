package util

import (
	"os"
	"strings"
)

// GetEnvironment returns the current environment: "development", "staging", or "production"
// Defaults to "development" in tests, "production" otherwise
func GetEnvironment() string {
	env := os.Getenv("MANDOR_ENV")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		// Default to development mode in tests, production otherwise
		if isTestMode() {
			env = "development"
		} else {
			env = "production"
		}
	}
	return strings.ToLower(env)
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return GetEnvironment() == "development" || GetEnvironment() == "dev"
}

// isTestMode detects if we're running in test mode
func isTestMode() bool {
	// Check if testing flag is set (go test sets this)
	return strings.Contains(os.Args[0], "test") || strings.Contains(os.Args[0], ".test")
}

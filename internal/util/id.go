package util

import (
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// GenerateNanoID generates a 12-character nanoid
func GenerateNanoID() (string, error) {
	return gonanoid.New(12)
}

// IsValidWorkspaceName validates workspace name
func IsValidWorkspaceName(name string) bool {
	// Allow only alphanumeric, hyphens, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	return matched
}

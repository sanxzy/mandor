package util

import (
	"crypto/rand"
	"regexp"
)

// GenerateID generates a 4-character alphanumeric ID
func GenerateID() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i := range bytes {
		bytes[i] = chars[int(bytes[i])%len(chars)]
	}
	return string(bytes), nil
}

// IsValidWorkspaceName validates workspace name
func IsValidWorkspaceName(name string) bool {
	// Allow only alphanumeric, hyphens, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	return matched
}

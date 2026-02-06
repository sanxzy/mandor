package util

import (
	"crypto/rand"
	"regexp"
	"strings"
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

// ToSlug converts a name to a slug format
// Converts to lowercase, replaces non-alphanumeric chars with hyphens,
// and collapses consecutive hyphens into single hyphen
func ToSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace non-alphanumeric characters with hyphens
	re := regexp.MustCompile("[^a-z0-9-]+")
	slug = re.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Collapse consecutive hyphens to single hyphen
	re = regexp.MustCompile("-+")
	slug = re.ReplaceAllString(slug, "-")

	return slug
}

// NextSequential generates a random 4-character base-62 ID with given prefix
// prefix is used as the ID prefix (e.g., "req", "scenario", "task")
// existingIDs parameter is kept for backward compatibility but not used
// Format: "prefix-XXXX" where XXXX is 4 random base-62 characters (0-9, a-z, A-Z)
// Example: NextSequential("req", []) might return "req-a7K2" or "req-mPx9"
// Base-62: 0-9 (0-9), a-z (10-35), A-Z (36-61)
// Capacity: 62^4 = 14,776,336 possible IDs per prefix
func NextSequential(prefix string, existingIDs []string) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const idLength = 4
	
	prefixLower := strings.ToLower(prefix)
	
	// Generate 4 random base-62 characters
	bytes := make([]byte, idLength)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to deterministic ID if random read fails
		return prefixLower + "-0000"
	}
	
	id := ""
	for _, b := range bytes {
		id += string(charset[int(b)%len(charset)])
	}
	
	return prefixLower + "-" + id
}

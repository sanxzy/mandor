package util

import (
	"os"
	"os/exec"
	"strings"
)

// GetGitUsername retrieves the git user.name from git config
func GetGitUsername() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// IsGitUserConfigured checks if git user is configured
func IsGitUserConfigured() bool {
	return GetGitUsername() != "unknown"
}

// GetGitUsernameWithWarning returns username and a warning if not configured
func GetGitUsernameWithWarning() (string, string) {
	username := GetGitUsername()
	if username == "unknown" {
		return username, "Warning: Git user not configured. Run 'git config user.name \"Your Name\"' to set your identity."
	}
	return username, ""
}

// GetCurrentDirectory returns the current working directory name
func GetCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extract the last part of the path (directory name)
	parts := strings.Split(cwd, string(os.PathSeparator))
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}
	return "", nil
}

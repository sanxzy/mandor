package ai

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		mandorDir := filepath.Join(dir, ".mandor")
		if _, err := os.Stat(mandorDir); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .mandor directory found in current directory or any parent directory")
		}
		dir = parent
	}
}

func FindProjectRootFrom(path string) (string, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	}

	for {
		mandorDir := filepath.Join(dir, ".mandor")
		if _, err := os.Stat(mandorDir); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .mandor directory found in path or any parent directory: %s", path)
		}
		dir = parent
	}
}

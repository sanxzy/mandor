package ai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	mandorDir := filepath.Join(tmpDir, ".mandor")
	os.MkdirAll(mandorDir, 0755)

	// Change to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test finding project root from current directory
	root, err := FindProjectRoot()
	if err != nil {
		t.Errorf("FindProjectRoot() error = %v, want nil", err)
	}

	// Verify a .mandor directory exists in root
	if _, err := os.Stat(filepath.Join(root, ".mandor")); os.IsNotExist(err) {
		t.Errorf("FindProjectRoot() returned path without .mandor directory: %s", root)
	}
}

func TestFindProjectRoot_NotFound(t *testing.T) {
	// Create a temporary directory without .mandor
	tmpDir := t.TempDir()

	// Change to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test that error is returned when .mandor not found
	_, err = FindProjectRoot()
	if err == nil {
		t.Error("FindProjectRoot() error = nil, want error")
	}
}

func TestFindProjectRoot_NestedDirectory(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	mandorDir := filepath.Join(tmpDir, ".mandor")
	os.MkdirAll(mandorDir, 0755)

	// Create a nested subdirectory
	nestedDir := filepath.Join(tmpDir, "nested", "deep", "dir")
	os.MkdirAll(nestedDir, 0755)

	// Change to the nested directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(nestedDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test finding project root from nested directory
	root, err := FindProjectRoot()
	if err != nil {
		t.Errorf("FindProjectRoot() error = %v, want nil", err)
	}

	// Verify a .mandor directory exists in root
	if _, err := os.Stat(filepath.Join(root, ".mandor")); os.IsNotExist(err) {
		t.Errorf("FindProjectRoot() returned path without .mandor directory: %s", root)
	}
}

func TestFindProjectRootFrom(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	mandorDir := filepath.Join(tmpDir, ".mandor")
	os.MkdirAll(mandorDir, 0755)

	// Test finding project root from specific path
	root, err := FindProjectRootFrom(tmpDir)
	if err != nil {
		t.Errorf("FindProjectRootFrom() error = %v, want nil", err)
	}

	if root != tmpDir {
		t.Errorf("FindProjectRootFrom() = %s, want %s", root, tmpDir)
	}
}

func TestFindProjectRootFrom_Nested(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	mandorDir := filepath.Join(tmpDir, ".mandor")
	os.MkdirAll(mandorDir, 0755)

	// Create a nested subdirectory
	nestedDir := filepath.Join(tmpDir, "nested", "deep", "dir")
	os.MkdirAll(nestedDir, 0755)

	// Test finding project root from nested path
	root, err := FindProjectRootFrom(nestedDir)
	if err != nil {
		t.Errorf("FindProjectRootFrom() error = %v, want nil", err)
	}

	if root != tmpDir {
		t.Errorf("FindProjectRootFrom() = %s, want %s", root, tmpDir)
	}
}

func TestFindProjectRootFrom_NotFound(t *testing.T) {
	// Create a temporary directory without .mandor
	tmpDir := t.TempDir()

	// Test that error is returned when .mandor not found
	_, err := FindProjectRootFrom(tmpDir)
	if err == nil {
		t.Error("FindProjectRootFrom() error = nil, want error")
	}
}

func TestFindProjectRootFrom_NonexistentPath(t *testing.T) {
	// Test with non-existent path
	nonexistentPath := "/nonexistent/path/that/should/not/exist"

	_, err := FindProjectRootFrom(nonexistentPath)
	if err == nil {
		t.Error("FindProjectRootFrom() error = nil, want error for non-existent path")
	}
}

func TestFindProjectRootFrom_RelativePath(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	mandorDir := filepath.Join(tmpDir, ".mandor")
	os.MkdirAll(mandorDir, 0755)

	// Create a nested subdirectory
	nestedDir := filepath.Join(tmpDir, "nested")
	os.MkdirAll(nestedDir, 0755)

	// Change to tmpDir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test with relative path
	root, err := FindProjectRootFrom("./nested")
	if err != nil {
		t.Errorf("FindProjectRootFrom() error = %v, want nil", err)
	}

	absPath, _ := filepath.Abs(".")
	if root != absPath {
		t.Errorf("FindProjectRootFrom() = %s, want %s", root, absPath)
	}
}

package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/cmd/workspace"
)

// TestInitCmdSuccess tests successful initialization
func TestInitCmdSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	cmdInit := workspace.NewInitCmd()
	err := cmdInit.RunE(cmdInit, []string{})

	if err != nil {
		t.Fatalf("Init command failed: %v", err)
	}

	// Verify .mandor directory exists
	mandorDir := filepath.Join(tmpDir, ".mandor")
	if _, err := os.Stat(mandorDir); os.IsNotExist(err) {
		t.Error(".mandor directory was not created")
	}

	// Verify workspace.json exists
	workspacePath := filepath.Join(mandorDir, "workspace.json")
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		t.Error("workspace.json was not created")
	}

	// Verify projects directory exists
	projectsDir := filepath.Join(mandorDir, "projects")
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		t.Error("projects directory was not created")
	}
}

// TestInitCmdWithCustomName tests init with custom workspace name
func TestInitCmdWithCustomName(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	cmdInit := workspace.NewInitCmd()
	cmdInit.Flags().Set("workspace-name", "my-workspace")
	err := cmdInit.RunE(cmdInit, []string{})

	if err != nil {
		t.Fatalf("Init command with custom name failed: %v", err)
	}

	// Verify workspace was created
	mandorDir := filepath.Join(tmpDir, ".mandor")
	if _, err := os.Stat(mandorDir); os.IsNotExist(err) {
		t.Error(".mandor directory was not created")
	}
}

// TestInitCmdAlreadyInitialized tests init fails when already initialized
func TestInitCmdAlreadyInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize first time
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Try to initialize again
	cmdInit2 := workspace.NewInitCmd()
	err := cmdInit2.RunE(cmdInit2, []string{})

	if err == nil {
		t.Error("Expected error when initializing already-initialized workspace")
	}
}

// TestInitCmdWorkspaceNameValidation tests invalid workspace names
func TestInitCmdWorkspaceNameValidation(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"valid-name", false},
		{"valid_name", false},
		{"validname123", false},
		{"123validname", false},
		{"invalid name!", true},
		{"invalid@name", true},
		{"invalid.name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			os.Chdir(testDir)

			cmdInit := workspace.NewInitCmd()
			cmdInit.Flags().Set("workspace-name", tt.name)
			err := cmdInit.RunE(cmdInit, []string{})

			if (err != nil) != tt.shouldErr {
				t.Errorf("workspace name %s: expected error=%v, got=%v", tt.name, tt.shouldErr, err != nil)
			}

			os.Chdir(tmpDir)
		})
	}
}

package feature_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"mandor/internal/cmd/feature"
)

func TestFeatureListCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewListCmd()
	cmd.SetArgs([]string{"--project", "testproject"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-initialized workspace")
	}
}

func TestFeatureListCmd_MissingProject(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewListCmd()
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing project")
	}
}

func TestFeatureListCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewListCmd()
	cmd.SetArgs([]string{"--project", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestFeatureListCmd_Success(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "draft")

	cmd := feature.NewListCmd()
	cmd.SetArgs([]string{"--project", "testproject"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}
}

func TestFeatureListCmd_JSON(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "draft")

	cmd := feature.NewListCmd()
	cmd.SetArgs([]string{"--project", "testproject", "--json"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\"features\"") {
		t.Errorf("Expected JSON output with 'features' key, got: %s", output)
	}
}

func TestNewFeatureListCmd(t *testing.T) {
	cmd := feature.NewListCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "list") {
		t.Errorf("Expected 'list' in use, got %q", cmd.Use)
	}
}

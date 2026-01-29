package feature_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"mandor/internal/cmd/feature"
)

func TestFeatureDetailCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-initialized workspace")
	}
}

func TestFeatureDetailCmd_MissingProject(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing project")
	}
}

func TestFeatureDetailCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestFeatureDetailCmd_FeatureNotFound(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-nonexistent", "--project", "testproject"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent feature")
	}
}

func TestFeatureDetailCmd_Success(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "active")

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject"})

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

func TestFeatureDetailCmd_JSON(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "active")

	cmd := feature.NewDetailCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--json"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\"id\"") {
		t.Errorf("Expected JSON output with 'id' key, got: %s", output)
	}
}

func TestNewFeatureDetailCmd(t *testing.T) {
	cmd := feature.NewDetailCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "detail") {
		t.Errorf("Expected 'detail' in use, got %q", cmd.Use)
	}
}

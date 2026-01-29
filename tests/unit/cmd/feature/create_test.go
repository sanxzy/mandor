package feature_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"mandor/internal/cmd/feature"
	"mandor/internal/domain"
)

func TestFeatureCreateCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--project", "testproject", "--goal", "Test goal"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-initialized workspace")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestFeatureCreateCmd_MissingProject(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--goal", "Test goal"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing project")
	}
}

func TestFeatureCreateCmd_MissingGoal(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--project", "testproject"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing goal")
	}
}

func TestFeatureCreateCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--project", "nonexistent", "--goal", "Test goal"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestFeatureCreateCmd_InvalidScope(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--project", "testproject", "--goal", "Test goal", "--scope", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid scope")
	}
}

func TestFeatureCreateCmd_InvalidPriority(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"Test Feature", "--project", "testproject", "--goal", "Test goal", "--priority", "P6"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

func TestFeatureCreateCmd_Success(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{"New Feature", "--project", "testproject", "--goal", "This is a test feature goal"})

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

	if !strings.Contains(output, "Feature created:") {
		t.Error("Expected 'Feature created:' in output")
	}
}

func TestFeatureCreateCmd_WithAllFlags(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{
		"Full Feature",
		"--project", "testproject",
		"--goal", "This is a test feature goal",
		"--scope", "frontend",
		"--priority", "P1",
	})

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

func TestFeatureCreateCmd_WithDepends(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "done")

	cmd := feature.NewCreateCmd()
	cmd.SetArgs([]string{
		"Dependent Feature",
		"--project", "testproject",
		"--goal", "This feature depends on another",
		"--depends", "testproject-feature-abc123",
	})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestNewFeatureCreateCmd(t *testing.T) {
	cmd := feature.NewCreateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "create") {
		t.Errorf("Expected 'create' in use, got %q", cmd.Use)
	}
}

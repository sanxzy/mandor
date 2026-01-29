package feature_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"mandor/internal/cmd/feature"
	"mandor/internal/domain"
)

func TestFeatureUpdateCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--name", "New Name"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-initialized workspace")
	}
}

func TestFeatureUpdateCmd_MissingProject(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--name", "New Name"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing project")
	}
}

func TestFeatureUpdateCmd_FeatureNotFound(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-nonexistent", "--project", "testproject", "--name", "New Name"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent feature")
	}
}

func TestFeatureUpdateCmd_InvalidPriority(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "draft")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--priority", "P10"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

func TestFeatureUpdateCmd_CancelledFeature(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "cancelled")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--name", "New Name"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for cancelled feature")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestFeatureUpdateCmd_Success(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "draft")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--name", "Updated Name"})

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

func TestFeatureUpdateCmd_Reopen(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "cancelled")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--reopen"})

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

func TestFeatureUpdateCmd_DryRun(t *testing.T) {
	tmpDir := setupTestFeatureWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForCmd(t, tmpDir, "testproject")
	writeTestFeatureForCmd(t, tmpDir, "testproject", "testproject-feature-abc123", "draft")

	cmd := feature.NewUpdateCmd()
	cmd.SetArgs([]string{"testproject-feature-abc123", "--project", "testproject", "--name", "New Name", "--dry-run"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[DRY RUN]") {
		t.Error("Expected '[DRY RUN]' in output")
	}
}

func TestNewFeatureUpdateCmd(t *testing.T) {
	cmd := feature.NewUpdateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "update") {
		t.Errorf("Expected 'update' in use, got %q", cmd.Use)
	}
}

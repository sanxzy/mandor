package project_test

import (
	"bytes"
	"os"
	"testing"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
)

func TestReopenCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewReopenCmd()
	cmd.SetArgs([]string{"test"})
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

func TestReopenCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewReopenCmd()
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestReopenCmd_NotDeleted(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "activeproj", domain.ProjectStatusActive)

	cmd := project.NewReopenCmd()
	cmd.SetArgs([]string{"activeproj"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-deleted project")
	}
}

func TestReopenCmd_Success(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewReopenCmd()
	cmd.SetArgs([]string{"deletedproj", "-y"})
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

func TestReopenCmd_Cancelled(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewReopenCmd()
	cmd.SetArgs([]string{"deletedproj"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error when cancelled, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}
}

func TestNewReopenCmd(t *testing.T) {
	cmd := project.NewReopenCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "reopen <id>" {
		t.Errorf("Expected use 'reopen <id>', got %q", cmd.Use)
	}
}

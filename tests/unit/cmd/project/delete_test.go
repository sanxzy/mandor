package project_test

import (
	"bytes"
	"os"
	"testing"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
)

func TestDeleteCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewDeleteCmd()
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

func TestDeleteCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestDeleteCmd_AlreadyDeleted(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"deletedproj"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for already deleted project")
	}
}

func TestDeleteCmd_SoftDelete(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"testproj", "-y"})
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

func TestDeleteCmd_DryRun(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"testproj", "--dry-run"})
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

func TestDeleteCmd_HardDelete(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"testproj", "--hard", "-y"})
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

func TestDeleteCmd_HardDeleteRequiresConfirmation(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"testproj", "--hard"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when confirmation doesn't match")
	}
}

func TestDeleteCmd_HardDeleteOnDeletedProject(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewDeleteCmd()
	cmd.SetArgs([]string{"deletedproj", "--hard", "-y"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error for hard delete on deleted project, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}
}

func TestNewDeleteCmd(t *testing.T) {
	cmd := project.NewDeleteCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "delete <id>" {
		t.Errorf("Expected use 'delete <id>', got %q", cmd.Use)
	}
}

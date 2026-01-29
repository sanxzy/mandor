package project_test

import (
	"bytes"
	"os"
	"testing"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
)

func TestDetailCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewDetailCmd()
	cmd.SetArgs([]string{"testproject"})
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

func TestDetailCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewDetailCmd()
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestDetailCmd_Success(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDetailCmd()
	cmd.SetArgs([]string{"testproj"})
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

func TestDetailCmd_JsonOutput(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewDetailCmd()
	cmd.SetArgs([]string{"testproj", "--json"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected JSON output, got empty")
	}
}

func TestDetailCmd_DeletedProject(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewDetailCmd()
	cmd.SetArgs([]string{"deletedproj"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error for deleted project, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}
}

func TestNewDetailCmd(t *testing.T) {
	cmd := project.NewDetailCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "detail <id>" {
		t.Errorf("Expected use 'detail <id>', got %q", cmd.Use)
	}
}

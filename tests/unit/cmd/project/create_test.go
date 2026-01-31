package project_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
)

func setupTestWorkspace(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mandorDir := filepath.Join(tmpDir, ".mandor")
	projectsDir := filepath.Join(mandorDir, "projects")

	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create .mandor directory: %v", err)
	}

	workspacePath := filepath.Join(mandorDir, "workspace.json")
	wsData := []byte("{\"version\":\"v1.0.0\"}")
	if err := os.WriteFile(workspacePath, wsData, 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write workspace.json: %v", err)
	}

	return tmpDir
}

func TestCreateCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"test", "--name", "Test", "--goal", string(make([]byte, 501))})

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

func TestCreateCmd_InvalidID(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"123invalid", "--name", "Test", "--goal", string(make([]byte, 501))})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}

func TestCreateCmd_MissingName(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"test", "--goal", string(make([]byte, 501)), "-y"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing name")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestCreateCmd_MissingGoal(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"test", "--name", "Test", "-y"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing goal")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestCreateCmd_GoalTooShort(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	// Force production mode for this test so validation fails
	oldEnv := os.Getenv("MANDOR_ENV")
	os.Setenv("MANDOR_ENV", "production")
	defer os.Setenv("MANDOR_ENV", oldEnv)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"test", "--name", "Test", "--goal", "short goal"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for goal too short")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestCreateCmd_ProjectExists(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	projectDir := filepath.Join(tmpDir, ".mandor", "projects", "existing")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"existing", "--name", "Test", "--goal", string(make([]byte, 501))})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for existing project")
	}
}

func TestCreateCmd_Success(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{"testproject", "--name", "Test Project", "--goal", string(make([]byte, 501))})

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

func TestCreateCmd_AllFlags(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewCreateCmd()
	cmd.SetArgs([]string{
		"fullproject",
		"--name", "Full Project",
		"--goal", string(make([]byte, 501)),
		"--task-dep", "cross_project_allowed",
		"--feature-dep", "same_project_only",
		"--issue-dep", "disabled",
		"--strict",
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

func TestNewCreateCmd(t *testing.T) {
	cmd := project.NewCreateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "create <id>" {
		t.Errorf("Expected use 'create <id>', got %q", cmd.Use)
	}
}

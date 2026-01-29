package project_test

import (
	"bytes"
	"os"
	"testing"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
)

func TestUpdateCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"test", "--name", "New Name"})
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

func TestUpdateCmd_ProjectNotFound(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"nonexistent", "--name", "New Name"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

func TestUpdateCmd_DeletedProject(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"deletedproj", "--name", "New Name"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for deleted project")
	}
}

func TestUpdateCmd_NoUpdates(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when no updates specified")
	}
}

func TestUpdateCmd_UpdateName(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj", "--name", "Updated Name"})
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

func TestUpdateCmd_UpdateGoal(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	newGoal := string(make([]byte, 501))
	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj", "--goal", newGoal})
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

func TestUpdateCmd_UpdateGoalTooShort(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj", "--goal", "short goal"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for goal too short")
	}
}

func TestUpdateCmd_UpdateStrict(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj", "--strict", "true"})
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

func TestUpdateCmd_InvalidStrictValue(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewUpdateCmd()
	cmd.SetArgs([]string{"testproj", "--strict", "invalid"})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid strict value")
	}
}

func TestNewUpdateCmd(t *testing.T) {
	cmd := project.NewUpdateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "update <id>" {
		t.Errorf("Expected use 'update <id>', got %q", cmd.Use)
	}
}

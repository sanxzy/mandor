package project_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"mandor/internal/cmd/project"
	"mandor/internal/domain"
	"mandor/internal/fs"
)

func writeTestProjectForCmd(t *testing.T, tmpDir, projectID string, status string) {
	t.Helper()

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}
	writer := fs.NewWriter(paths)

	project := &domain.Project{
		ID:        projectID,
		Name:      "Test Project",
		Goal:      string(make([]byte, 501)),
		Status:    status,
		Strict:    false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if err := writer.WriteProjectMetadata(projectID, project); err != nil {
		t.Fatalf("Failed to write project metadata: %v", err)
	}

	schema := domain.DefaultProjectSchema("", "", "")
	if err := writer.WriteProjectSchema(projectID, &schema); err != nil {
		t.Fatalf("Failed to write project schema: %v", err)
	}
}

func TestListCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewListCmd()
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

func TestListCmd_EmptyWorkspace(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := project.NewListCmd()
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

func TestListCmd_WithProjects(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "project1", domain.ProjectStatusActive)
	writeTestProjectForCmd(t, tmpDir, "project2", domain.ProjectStatusInitial)

	cmd := project.NewListCmd()
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

func TestListCmd_IncludeDeleted(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "active", domain.ProjectStatusActive)
	writeTestProjectForCmd(t, tmpDir, "deleted", domain.ProjectStatusDeleted)

	cmd := project.NewListCmd()
	cmd.SetArgs([]string{"--include-deleted"})
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

func TestListCmd_JsonOutput(t *testing.T) {
	tmpDir := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	writeTestProjectForCmd(t, tmpDir, "testproj", domain.ProjectStatusActive)

	cmd := project.NewListCmd()
	cmd.SetArgs([]string{"--json"})
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

func TestNewListCmd(t *testing.T) {
	cmd := project.NewListCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "list" {
		t.Errorf("Expected use 'list', got %q", cmd.Use)
	}
}

package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/cmd/workspace"
	"mandor/internal/service"
)

// TestStatusCmdEmpty tests status command with no projects
func TestStatusCmdEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Test status
	statusCmd := workspace.NewStatusCmd()
	err := statusCmd.RunE(statusCmd, []string{})

	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}
}

// TestStatusCmdJSON tests status command with JSON format
func TestStatusCmdJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Test status with --json flag
	statusCmd := workspace.NewStatusCmd()
	statusCmd.Flags().Set("json", "true")
	err := statusCmd.RunE(statusCmd, []string{})

	if err != nil {
		t.Fatalf("Status JSON command failed: %v", err)
	}
}

// TestStatusCmdSummary tests status command with summary format
func TestStatusCmdSummary(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Test status with --summary flag
	statusCmd := workspace.NewStatusCmd()
	statusCmd.Flags().Set("summary", "true")
	err := statusCmd.RunE(statusCmd, []string{})

	if err != nil {
		t.Fatalf("Status summary command failed: %v", err)
	}
}

// TestStatusCmdSpecificProject tests status command for specific project
func TestStatusCmdSpecificProject(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Create project directory
	projectID := "test-project"
	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	os.MkdirAll(projectDir, 0755)

	// Write minimal project metadata
	metadata := []byte(`{"id":"test-project","name":"Test","status":"initial","created_at":"2026-01-27T00:00:00Z","last_updated_at":"2026-01-27T00:00:00Z","created_by":"test","schema_version":"mandor.v1"}`)
	os.WriteFile(filepath.Join(projectDir, "project.jsonl"), metadata, 0644)

	// Test status with --project flag
	statusCmd := workspace.NewStatusCmd()
	statusCmd.Flags().Set("project", projectID)
	err := statusCmd.RunE(statusCmd, []string{})

	if err != nil {
		t.Fatalf("Status project command failed: %v", err)
	}
}

// TestStatusCmdNotInitialized tests status command without initialized workspace
func TestStatusCmdNotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Test status without initialization
	statusCmd := workspace.NewStatusCmd()
	err := statusCmd.RunE(statusCmd, []string{})

	if err == nil {
		t.Error("Expected error when workspace not initialized")
	}
}

// TestStatusAllFormats tests all output format combinations
func TestStatusAllFormats(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	formats := []struct {
		flag  string
		value string
	}{
		{"default", ""},
		{"json", "true"},
		{"summary", "true"},
	}

	for _, format := range formats {
		t.Run(format.flag, func(t *testing.T) {
			statusCmd := workspace.NewStatusCmd()
			if format.value != "" {
				statusCmd.Flags().Set(format.flag, format.value)
			}
			err := statusCmd.RunE(statusCmd, []string{})

			if err != nil {
				t.Errorf("Status with %s format failed: %v", format.flag, err)
			}
		})
	}
}

// TestStatusWorkspaceInfo tests status returns workspace information
func TestStatusWorkspaceInfo(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Get status via service
	statusSvc, _ := service.NewStatusService()
	status, _ := statusSvc.GetWorkspaceStatus("")

	if status.Workspace == nil {
		t.Error("Workspace info missing from status")
	}

	if status.Workspace.ID == "" {
		t.Error("Workspace ID missing")
	}

	if status.Workspace.Name == "" {
		t.Error("Workspace name missing")
	}

	if status.Workspace.CreatedBy == "" {
		t.Error("Workspace creator missing")
	}
}

// TestStatusEmptyProjects tests status with no projects
func TestStatusEmptyProjects(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Get status
	statusSvc, _ := service.NewStatusService()
	status, _ := statusSvc.GetWorkspaceStatus("")

	if len(status.Projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(status.Projects))
	}

	if status.Totals.Features != 0 {
		t.Errorf("Expected 0 features, got %d", status.Totals.Features)
	}

	if status.Totals.Tasks != 0 {
		t.Errorf("Expected 0 tasks, got %d", status.Totals.Tasks)
	}

	if status.Totals.Issues != 0 {
		t.Errorf("Expected 0 issues, got %d", status.Totals.Issues)
	}
}

// TestStatusDependencyTracking tests dependency tracking in status
func TestStatusDependencyTracking(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Get status
	statusSvc, _ := service.NewStatusService()
	status, _ := statusSvc.GetWorkspaceStatus("")

	// New workspace should have no circular dependencies
	if status.Dependencies.CircularDeps != 0 {
		t.Error("New workspace should have 0 circular dependencies")
	}

	// New workspace should have no cross-project dependencies
	if status.Dependencies.CrossProjectCount != 0 {
		t.Error("New workspace should have 0 cross-project dependencies")
	}
}

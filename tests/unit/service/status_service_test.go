package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/service"
)

// TestGetWorkspaceStatusEmpty tests status with no projects
func TestGetWorkspaceStatusEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	wsSvc, _ := service.NewWorkspaceService()
	wsSvc.InitWorkspace("")

	// Get status
	statusSvc, _ := service.NewStatusService()
	status, err := statusSvc.GetWorkspaceStatus("")

	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.Workspace == nil {
		t.Error("Workspace is nil")
	}

	if len(status.Projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(status.Projects))
	}

	if status.Totals.Features != 0 {
		t.Errorf("Expected 0 features, got %d", status.Totals.Features)
	}
}

// TestGetWorkspaceStatusNotInitialized tests status when workspace not initialized
func TestGetWorkspaceStatusNotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	statusSvc, _ := service.NewStatusService()
	_, err := statusSvc.GetWorkspaceStatus("")

	if err == nil {
		t.Error("Expected error when workspace not initialized")
	}
}

// TestGetWorkspaceStatusProjectNotFound tests status with invalid project
func TestGetWorkspaceStatusProjectNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	wsSvc, _ := service.NewWorkspaceService()
	wsSvc.InitWorkspace("")

	// Try to get status for non-existent project
	statusSvc, _ := service.NewStatusService()
	_, err := statusSvc.GetWorkspaceStatus("nonexistent")

	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

// TestGetProjectStatus tests detailed project status
func TestGetProjectStatus(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	wsSvc, _ := service.NewWorkspaceService()
	wsSvc.InitWorkspace("")

	// Create a project directory with metadata
	projectID := "test-project"
	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	os.MkdirAll(projectDir, 0755)

	// Write minimal project metadata
	metadata := []byte(`{"id":"test-project","name":"Test","status":"initial","created_at":"2026-01-27T00:00:00Z","last_updated_at":"2026-01-27T00:00:00Z","created_by":"test","schema_version":"mandor.v1"}`)
	os.WriteFile(filepath.Join(projectDir, "project.jsonl"), metadata, 0644)

	// Get project status
	statusSvc, _ := service.NewStatusService()
	status, err := statusSvc.GetProjectStatus(projectID)

	if err != nil {
		t.Fatalf("Failed to get project status: %v", err)
	}

	if status.ID != projectID {
		t.Errorf("Expected ID %s, got %s", projectID, status.ID)
	}

	if status.Stats.Features.Total != 0 {
		t.Errorf("Expected 0 features, got %d", status.Stats.Features.Total)
	}
}

// TestGetProjectStatusNotFound tests project status when project doesn't exist
func TestGetProjectStatusNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	wsSvc, _ := service.NewWorkspaceService()
	wsSvc.InitWorkspace("")

	// Try to get status for non-existent project
	statusSvc, _ := service.NewStatusService()
	_, err := statusSvc.GetProjectStatus("nonexistent")

	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

// TestStatusWithDependencies tests status returns dependency info
func TestStatusWithDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	wsSvc, _ := service.NewWorkspaceService()
	wsSvc.InitWorkspace("")

	// Get status
	statusSvc, _ := service.NewStatusService()
	status, _ := statusSvc.GetWorkspaceStatus("")

	if status.Dependencies.CircularDeps != 0 {
		t.Error("New workspace should have no circular dependencies")
	}

	if status.Dependencies.CrossProjectCount != 0 {
		t.Error("New workspace should have no cross-project dependencies")
	}
}

package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/service"
)

// TestInitWorkspace tests workspace initialization
func TestInitWorkspace(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, err := service.NewWorkspaceService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Initialize workspace
	ws, err := svc.InitWorkspace("")
	if err != nil {
		t.Fatalf("Failed to init workspace: %v", err)
	}

	if ws.Name == "" {
		t.Error("Workspace name is empty")
	}

	if ws.ID == "" {
		t.Error("Workspace ID is empty")
	}

	if ws.Version != "mandor.v1" {
		t.Errorf("Expected version mandor.v1, got %s", ws.Version)
	}

	// Check that .mandor directory exists
	mandorDir := filepath.Join(tmpDir, ".mandor")
	if _, err := os.Stat(mandorDir); os.IsNotExist(err) {
		t.Error(".mandor directory was not created")
	}

	// Check that workspace.json exists
	workspacePath := filepath.Join(mandorDir, "workspace.json")
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		t.Error("workspace.json was not created")
	}
}

// TestWorkspaceAlreadyExists tests that init fails when workspace exists
func TestWorkspaceAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	// Try to init again
	svc2, _ := service.NewWorkspaceService()
	_, err := svc2.InitWorkspace("")

	if err == nil {
		t.Error("Expected error when workspace already exists")
	}
}

// TestUpdateConfig tests configuration updates
func TestUpdateConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	// Update priority
	err := svc.UpdateWorkspaceConfig("default_priority", "P1")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify update
	value, _ := svc.GetConfigValue("default_priority")
	if value != "P1" {
		t.Errorf("Expected P1, got %v", value)
	}
}

// TestInvalidPriority tests that invalid priority is rejected
func TestInvalidPriority(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	err := svc.UpdateWorkspaceConfig("default_priority", "P9")
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

// TestGetWorkspace tests retrieving workspace
func TestGetWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	ws, err := svc.GetWorkspace()
	if err != nil {
		t.Fatalf("Failed to get workspace: %v", err)
	}

	if ws == nil {
		t.Error("Workspace is nil")
	}

	if ws.Config.DefaultPriority != "P3" {
		t.Errorf("Expected default priority P3, got %s", ws.Config.DefaultPriority)
	}

	if ws.Config.StrictMode != false {
		t.Error("Expected strict_mode to be false by default")
	}
}

// TestUpdateStrictMode tests strict mode configuration
func TestUpdateStrictMode(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	err := svc.UpdateWorkspaceConfig("strict_mode", true)
	if err != nil {
		t.Fatalf("Failed to update strict_mode: %v", err)
	}

	value, _ := svc.GetConfigValue("strict_mode")
	if value != true {
		t.Errorf("Expected true, got %v", value)
	}
}

// TestInvalidConfigKey tests invalid configuration key
func TestInvalidConfigKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	svc.InitWorkspace("")

	err := svc.UpdateWorkspaceConfig("invalid_key", "value")
	if err == nil {
		t.Error("Expected error for invalid config key")
	}
}

// TestWorkspaceIDUnique tests that workspace IDs are unique
func TestWorkspaceIDUnique(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)

	os.Chdir(tmpDir1)
	svc1, _ := service.NewWorkspaceService()
	ws1, _ := svc1.InitWorkspace("")

	os.Chdir(tmpDir2)
	svc2, _ := service.NewWorkspaceService()
	ws2, _ := svc2.InitWorkspace("")

	if ws1.ID == ws2.ID {
		t.Error("Workspace IDs should be unique")
	}
}

// TestWorkspaceTimestampFormat tests ISO8601 timestamp format
func TestWorkspaceTimestampFormat(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	svc, _ := service.NewWorkspaceService()
	ws, _ := svc.InitWorkspace("")

	// Check ISO8601 format (YYYY-MM-DDTHH:MM:SSZ)
	createdAtStr := ws.CreatedAt.Format("2006-01-02T15:04:05Z")
	if createdAtStr != ws.CreatedAt.UTC().Format("2006-01-02T15:04:05Z") {
		t.Error("Timestamp is not in UTC ISO8601 format")
	}
}

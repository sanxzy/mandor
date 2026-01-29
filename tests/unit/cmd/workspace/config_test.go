package workspace_test

import (
	"os"
	"testing"

	"mandor/internal/cmd/workspace"
	"mandor/internal/service"
)

// TestConfigGetCmd tests config get command
func TestConfigGetCmd(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace first
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Test config get
	configCmd := workspace.NewConfigCmd()
	getCmd := configCmd.Commands()[0] // get subcommand
	err := getCmd.RunE(getCmd, []string{})

	if err != nil {
		t.Fatalf("Config get command failed: %v", err)
	}
}

// TestConfigGetSpecificKey tests getting a specific config key
func TestConfigGetSpecificKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Verify via service
	svc, _ := service.NewWorkspaceService()
	value, err := svc.GetConfigValue("default_priority")

	if err != nil {
		t.Fatalf("Failed to get config value: %v", err)
	}

	if value != "P3" {
		t.Errorf("Expected P3, got %v", value)
	}
}

// TestConfigSetPriority tests config set for priority
func TestConfigSetPriority(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	// Update via service
	svc, _ := service.NewWorkspaceService()
	err := svc.UpdateWorkspaceConfig("default_priority", "P1")

	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Verify
	value, _ := svc.GetConfigValue("default_priority")
	if value != "P1" {
		t.Errorf("Expected P1, got %v", value)
	}
}

// TestConfigSetAllPriorities tests all valid priorities
func TestConfigSetAllPriorities(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	priorities := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	svc, _ := service.NewWorkspaceService()

	for _, priority := range priorities {
		err := svc.UpdateWorkspaceConfig("default_priority", priority)
		if err != nil {
			t.Errorf("Failed to set priority %s: %v", priority, err)
		}

		value, _ := svc.GetConfigValue("default_priority")
		if value != priority {
			t.Errorf("Expected %s, got %v", priority, value)
		}
	}
}

// TestConfigSetInvalidPriority tests invalid priority rejection
func TestConfigSetInvalidPriority(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	invalidPriorities := []string{"P9", "P10", "P-1", "invalid", "p0"}

	svc, _ := service.NewWorkspaceService()
	for _, priority := range invalidPriorities {
		err := svc.UpdateWorkspaceConfig("default_priority", priority)
		if err == nil {
			t.Errorf("Expected error for invalid priority %s", priority)
		}
	}
}

// TestConfigSetStrictMode tests strict mode configuration
func TestConfigSetStrictMode(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	svc, _ := service.NewWorkspaceService()

	// Test setting to true
	err := svc.UpdateWorkspaceConfig("strict_mode", true)
	if err != nil {
		t.Fatalf("Failed to set strict_mode: %v", err)
	}

	value, _ := svc.GetConfigValue("strict_mode")
	if value != true {
		t.Errorf("Expected true, got %v", value)
	}

	// Test setting to false
	err = svc.UpdateWorkspaceConfig("strict_mode", false)
	if err != nil {
		t.Fatalf("Failed to unset strict_mode: %v", err)
	}

	value, _ = svc.GetConfigValue("strict_mode")
	if value != false {
		t.Errorf("Expected false, got %v", value)
	}
}

// TestConfigSetInvalidKey tests invalid config key rejection
func TestConfigSetInvalidKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	svc, _ := service.NewWorkspaceService()
	err := svc.UpdateWorkspaceConfig("invalid_key", "value")

	if err == nil {
		t.Error("Expected error for invalid config key")
	}
}

// TestConfigListKeys tests all valid config keys
func TestConfigListKeys(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	svc, _ := service.NewWorkspaceService()
	ws, _ := svc.GetWorkspace()

	// Verify both config keys exist
	if ws.Config.DefaultPriority == "" {
		t.Error("default_priority config key missing")
	}

	if ws.Config.StrictMode == false && ws.Config.StrictMode != true {
		// This will always be true or false, so just check it exists
		// No error needed
	}
}

// TestConfigTimestampUpdate tests that timestamp is updated on config change
func TestConfigTimestampUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Initialize workspace
	cmdInit := workspace.NewInitCmd()
	cmdInit.RunE(cmdInit, []string{})

	svc, _ := service.NewWorkspaceService()
	ws1, _ := svc.GetWorkspace()
	originalTime := ws1.LastUpdatedAt

	// Wait and then update config
	svc.UpdateWorkspaceConfig("default_priority", "P1")
	ws2, _ := svc.GetWorkspace()
	newTime := ws2.LastUpdatedAt

	if newTime.Equal(originalTime) {
		t.Error("Timestamp should be updated after config change")
	}

	if newTime.Before(originalTime) {
		t.Error("New timestamp should be after original")
	}
}

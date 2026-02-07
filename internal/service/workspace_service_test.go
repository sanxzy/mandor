package service

import (
	"os"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

func setupTestWorkspaceService(t *testing.T) (*WorkspaceService, string) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	service := &WorkspaceService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

func TestNewWorkspaceService(t *testing.T) {
	service, err := NewWorkspaceService()
	if err != nil {
		t.Fatalf("NewWorkspaceService failed: %v", err)
	}

	if service == nil {
		t.Error("NewWorkspaceService returned nil")
	}
	if service.paths == nil {
		t.Error("Service paths should not be nil")
	}
	if service.reader == nil {
		t.Error("Service reader is nil")
	}
	if service.writer == nil {
		t.Error("Service writer is nil")
	}
}

func TestInitWorkspace_AlreadyInitialized(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "existing-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	_, err := service.InitWorkspace("new-name")
	if err == nil {
		t.Error("Expected error for already initialized workspace, got nil")
	}
}

func TestInitWorkspace_Success(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	workspace, err := service.InitWorkspace("test-workspace")
	if err != nil {
		t.Fatalf("InitWorkspace failed: %v", err)
	}

	if workspace.Name != "test-workspace" {
		t.Errorf("Workspace Name = %q, want %q", workspace.Name, "test-workspace")
	}
	if workspace.ID == "" {
		t.Error("Workspace ID should not be empty")
	}
	if workspace.Version != "mandor.v1" {
		t.Errorf("Workspace Version = %q, want %q", workspace.Version, "mandor.v1")
	}
}

func TestGetWorkspace_NotInitialized(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	_, err := service.GetWorkspace()
	if err == nil {
		t.Error("Expected error for uninitialized workspace, got nil")
	}
}

func TestGetWorkspace_Success(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	result, err := service.GetWorkspace()
	if err != nil {
		t.Fatalf("GetWorkspace failed: %v", err)
	}

	if result.Name != workspace.Name {
		t.Errorf("Workspace Name = %q, want %q", result.Name, workspace.Name)
	}
}

func TestUpdateWorkspaceConfig_DefaultPriority_Valid(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.UpdateWorkspaceConfig("default_priority", "P1")
	if err != nil {
		t.Errorf("UpdateWorkspaceConfig failed: %v", err)
	}

	updated, _ := service.GetWorkspace()
	if updated.Config.DefaultPriority != "P1" {
		t.Errorf("DefaultPriority = %q, want %q", updated.Config.DefaultPriority, "P1")
	}
}

func TestUpdateWorkspaceConfig_DefaultPriority_Invalid(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.UpdateWorkspaceConfig("default_priority", "P99")
	if err == nil {
		t.Error("Expected error for invalid priority, got nil")
	}
}

func TestUpdateWorkspaceConfig_StrictMode(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.UpdateWorkspaceConfig("strict_mode", true)
	if err != nil {
		t.Errorf("UpdateWorkspaceConfig failed: %v", err)
	}

	updated, _ := service.GetWorkspace()
	if !updated.Config.StrictMode {
		t.Error("StrictMode should be true")
	}
}

func TestUpdateWorkspaceConfig_GoalLengths(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.UpdateWorkspaceConfig("goal.lengths.feature", 400)
	if err != nil {
		t.Errorf("UpdateWorkspaceConfig failed: %v", err)
	}

	length := service.GetGoalLength("feature")
	if length != 400 {
		t.Errorf("Goal length = %d, want %d", length, 400)
	}
}

func TestUpdateWorkspaceConfig_UnknownKey(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.UpdateWorkspaceConfig("unknown_key", "value")
	if err == nil {
		t.Error("Expected error for unknown key, got nil")
	}
}

func TestGetConfigValue_DefaultPriority(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P2",
		},
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	value, err := service.GetConfigValue("default_priority")
	if err != nil {
		t.Fatalf("GetConfigValue failed: %v", err)
	}

	if value != "P2" {
		t.Errorf("Config value = %q, want %q", value, "P2")
	}
}

func TestGetConfigValue_UnknownKey(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	_, err := service.GetConfigValue("unknown_key")
	if err == nil {
		t.Error("Expected error for unknown key, got nil")
	}
}

func TestGetGoalLength(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config: domain.WorkspaceConfig{
			GoalLengths: domain.GoalLengths{
				Backlog: 500,
				Feature: 300,
				Task:    500,
				Issue:   200,
			},
		},
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	tests := []struct {
		entity   string
		expected int
	}{
		{"backlog", 500},
		{"feature", 300},
		{"task", 500},
		{"issue", 200},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.entity, func(t *testing.T) {
			length := service.GetGoalLength(tt.entity)
			if length != tt.expected {
				t.Errorf("GetGoalLength(%q) = %d, want %d", tt.entity, length, tt.expected)
			}
		})
	}
}

func TestSetGoalLength(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.SetGoalLength("feature", 450)
	if err != nil {
		t.Errorf("SetGoalLength failed: %v", err)
	}

	length := service.GetGoalLength("feature")
	if length != 450 {
		t.Errorf("Goal length = %d, want %d", length, 450)
	}
}

func TestSetGoalLength_InvalidEntity(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.SetGoalLength("unknown", 100)
	if err == nil {
		t.Error("Expected error for unknown entity, got nil")
	}
}

func TestResetGoalLength(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config: domain.WorkspaceConfig{
			GoalLengths: domain.GoalLengths{
				Feature: 999,
			},
		},
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.ResetGoalLength("feature")
	if err != nil {
		t.Errorf("ResetGoalLength failed: %v", err)
	}

	length := service.GetGoalLength("feature")
	if length != domain.FeatureGoalMinLength {
		t.Errorf("Goal length = %d, want %d (default)", length, domain.FeatureGoalMinLength)
	}
}

func TestResetGoalLength_InvalidEntity(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.ResetGoalLength("unknown")
	if err == nil {
		t.Error("Expected error for unknown entity, got nil")
	}
}

func TestResetAllGoalLengths(t *testing.T) {
	service, _ := setupTestWorkspaceService(t)

	os.MkdirAll(service.paths.WorkspaceRoot+"/.mandor", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config: domain.WorkspaceConfig{
			GoalLengths: domain.GoalLengths{
				Backlog: 999,
				Feature: 999,
				Task:    999,
				Issue:   999,
			},
		},
	}
	if err := service.writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	err := service.ResetAllGoalLengths()
	if err != nil {
		t.Errorf("ResetAllGoalLengths failed: %v", err)
	}

	defaults := domain.DefaultGoalLengths()

	tests := []struct {
		entity   string
		expected int
	}{
		{"backlog", defaults.Backlog},
		{"feature", defaults.Feature},
		{"task", defaults.Task},
		{"issue", defaults.Issue},
	}

	for _, tt := range tests {
		t.Run(tt.entity, func(t *testing.T) {
			length := service.GetGoalLength(tt.entity)
			if length != tt.expected {
				t.Errorf("GetGoalLength(%q) = %d, want %d (default)", tt.entity, length, tt.expected)
			}
		})
	}
}

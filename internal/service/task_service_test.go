package service

import (
	"os"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

func setupTestTaskService(t *testing.T) (*TaskService, string) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(tmpDir+"/.mandor/backlogs", 0755)
	os.MkdirAll(tmpDir+"/.mandor/features", 0755)
	os.MkdirAll(tmpDir+"/.mandor/specs", 0755)
	os.MkdirAll(tmpDir+"/.mandor/tasks", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			GoalLengths: domain.GoalLengths{
				Task: 500,
			},
		},
	}
	if err := writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	service := &TaskService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

func TestNewTaskServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewTaskServiceWithPaths(paths)

	if service == nil {
		t.Error("NewTaskServiceWithPaths returned nil")
	}
	if service.paths != paths {
		t.Error("Service paths not set correctly")
	}
	if service.reader == nil {
		t.Error("Service reader is nil")
	}
	if service.writer == nil {
		t.Error("Service writer is nil")
	}
}

func TestTaskWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestTaskService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

func TestParseTaskID_Valid(t *testing.T) {
	service := &TaskService{}

	backlogID, featureID, err := service.ParseTaskID("test-backlog-feature-test-feature-task-abc123")
	if err != nil {
		t.Fatalf("ParseTaskID failed: %v", err)
	}
	if backlogID != "test-backlog" {
		t.Errorf("backlogID = %q, want %q", backlogID, "test-backlog")
	}
	if featureID != "test-backlog-feature-test-feature" {
		t.Errorf("featureID = %q, want %q", featureID, "test-backlog-feature-test-feature")
	}
}

func TestParseTaskID_InvalidFormat(t *testing.T) {
	service := &TaskService{}

	tests := []struct {
		name   string
		taskID string
	}{
		{"Missing task separator", "test-backlog-feature-test-feature-abc123"},
		{"Missing feature separator", "test-backlog-task-abc123"},
		{"Empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ParseTaskID(tt.taskID)
			if err == nil {
				t.Errorf("Expected error for invalid task ID %q, got nil", tt.taskID)
			}
		})
	}
}

func TestExtractBacklogIDFromFeatureID_Valid(t *testing.T) {
	service := &TaskService{}

	backlogID, err := service.extractBacklogIDFromFeatureID("test-backlog-feature-test-feature")
	if err != nil {
		t.Fatalf("extractBacklogIDFromFeatureID failed: %v", err)
	}
	if backlogID != "test-backlog" {
		t.Errorf("backlogID = %q, want %q", backlogID, "test-backlog")
	}
}

func TestExtractBacklogIDFromFeatureID_Invalid(t *testing.T) {
	service := &TaskService{}

	backlogID, err := service.extractBacklogIDFromFeatureID("invalid")
	if err == nil {
		t.Error("Expected error for invalid feature ID, got nil")
	}
	if backlogID != "" {
		t.Errorf("backlogID should be empty on error, got %q", backlogID)
	}
}

func TestComparePriority(t *testing.T) {
	tests := []struct {
		p1       string
		p2       string
		wantSign int
	}{
		{"P0", "P1", -1},
		{"P1", "P0", 1},
		{"P3", "P3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.p1+"_vs_"+tt.p2, func(t *testing.T) {
			result := ComparePriority(tt.p1, tt.p2)
			if tt.wantSign == 0 {
				if result != 0 {
					t.Errorf("ComparePriority(%q, %q) = %d, want 0", tt.p1, tt.p2, result)
				}
			} else if (result < 0 && tt.wantSign > 0) || (result > 0 && tt.wantSign < 0) {
				t.Errorf("ComparePriority(%q, %q) = %d, want sign %d", tt.p1, tt.p2, result, tt.wantSign)
			}
		})
	}
}

func TestGetTaskGoalMinLength(t *testing.T) {
	service, _ := setupTestTaskService(t)

	minLen := service.getTaskGoalMinLength()
	if minLen != 500 {
		t.Errorf("getTaskGoalMinLength() = %d, want 500", minLen)
	}
}

func TestValidateStatusTransition_Valid(t *testing.T) {
	service := &TaskService{}

	tests := []struct {
		current string
		next    string
		valid   bool
	}{
		{domain.TaskStatusPending, domain.TaskStatusReady, true},
		{domain.TaskStatusPending, domain.TaskStatusInProgress, true},
		{domain.TaskStatusReady, domain.TaskStatusInProgress, true},
		{domain.TaskStatusInProgress, domain.TaskStatusDone, true},
		{domain.TaskStatusBlocked, domain.TaskStatusReady, true},
		{domain.TaskStatusReady, domain.TaskStatusDone, false},
		{domain.TaskStatusDone, domain.TaskStatusPending, false},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_to_"+tt.next, func(t *testing.T) {
			err := service.validateStatusTransition(tt.current, tt.next)
			if tt.valid && err != nil {
				t.Errorf("Expected valid transition from %s to %s, got error: %v", tt.current, tt.next, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid transition from %s to %s, got no error", tt.current, tt.next)
			}
		})
	}
}

func TestValidateNoCycle_NoCycle(t *testing.T) {
	service := &TaskService{}

	err := service.validateNoCycle("test-backlog", "", []string{"task-1"})
	if err != nil {
		t.Errorf("validateNoCycle failed: %v", err)
	}
}

func TestValidateNoDuplicateName_NoDuplicates(t *testing.T) {
	service, _ := setupTestTaskService(t)

	err := service.validateNoDuplicateName("test-backlog", "test-feature", "unique-name")
	if err != nil {
		t.Errorf("validateNoDuplicateName failed: %v", err)
	}
}

func TestValidateIAEScenariosExist_NoSpec(t *testing.T) {
	service, _ := setupTestTaskService(t)

	err := service.validateIAEScenariosExist("test-backlog", "nonexistent-spec", []string{"req-0001:scn-0001"})
	if err == nil {
		t.Error("Expected error for nonexistent spec, got nil")
	}
}

func TestCheckGatesBeforeInProgress_AllGatesMet(t *testing.T) {
	service := &TaskService{}

	gates := domain.ReadGates{
		IsReadBrief:        true,
		IsReadSpec:         true,
		IsReadSessionNotes: true,
	}

	err := service.checkGatesBeforeInProgress(gates)
	if err != nil {
		t.Errorf("checkGatesBeforeInProgress failed: %v", err)
	}
}

func TestCheckGatesBeforeInProgress_UnmetGates(t *testing.T) {
	service := &TaskService{}

	gates := domain.ReadGates{
		IsReadBrief:        false,
		IsReadSpec:         true,
		IsReadSessionNotes: true,
	}

	err := service.checkGatesBeforeInProgress(gates)
	if err == nil {
		t.Error("Expected error for unmet gates, got nil")
	}
}

func TestFindDependents_NoDependents(t *testing.T) {
	service, _ := setupTestTaskService(t)

	dependents, err := service.findDependents("test-backlog", "task-1")
	if err != nil {
		t.Errorf("findDependents failed: %v", err)
	}
	if len(dependents) != 0 {
		t.Errorf("dependents count = %d, want 0", len(dependents))
	}
}

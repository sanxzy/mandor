package domain

import (
	"testing"
	"time"
)

// ==========================
// Task Validation Tests
// ==========================

func TestValidateTaskID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid 12 chars", "auth-task-xyz", true},
		{"valid long ID", "auth-task-abc123def", true},
		{"empty", "", false},
		{"too short 11 chars", "auth-task-x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTaskID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateTaskID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestValidateTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"pending", TaskStatusPending, true},
		{"ready", TaskStatusReady, true},
		{"in_progress", TaskStatusInProgress, true},
		{"blocked", TaskStatusBlocked, true},
		{"done", TaskStatusDone, true},
		{"cancelled", TaskStatusCancelled, true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTaskStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidateTaskStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestValidateTaskGoalLength(t *testing.T) {
	tests := []struct {
		name     string
		goal     string
		minLen   int
		expected bool
	}{
		{"exactly min 500", string(make([]byte, 500)), 500, true},
		{"exceeds min 500", string(make([]byte, 501)), 500, true},
		{"too short", "short", 500, false},
		{"empty", "", 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTaskGoalLength(tt.goal, tt.minLen)
			if result != tt.expected {
				t.Errorf("ValidateTaskGoalLength(len=%d, minLen=%d) = %v, want %v",
					len(tt.goal), tt.minLen, result, tt.expected)
			}
		})
	}
}

// ==========================
// Task Status Constants Tests
// ==========================

func TestTaskStatusConstants(t *testing.T) {
	if TaskStatusPending != "pending" {
		t.Errorf("TaskStatusPending = %q, want %q", TaskStatusPending, "pending")
	}
	if TaskStatusReady != "ready" {
		t.Errorf("TaskStatusReady = %q, want %q", TaskStatusReady, "ready")
	}
	if TaskStatusInProgress != "in_progress" {
		t.Errorf("TaskStatusInProgress = %q, want %q", TaskStatusInProgress, "in_progress")
	}
	if TaskStatusBlocked != "blocked" {
		t.Errorf("TaskStatusBlocked = %q, want %q", TaskStatusBlocked, "blocked")
	}
	if TaskStatusDone != "done" {
		t.Errorf("TaskStatusDone = %q, want %q", TaskStatusDone, "done")
	}
	if TaskStatusCancelled != "cancelled" {
		t.Errorf("TaskStatusCancelled = %q, want %q", TaskStatusCancelled, "cancelled")
	}
}

func TestTaskGoalMinLength(t *testing.T) {
	if TaskGoalMinLength != 500 {
		t.Errorf("TaskGoalMinLength = %d, want 500", TaskGoalMinLength)
	}
}

// ==========================
// Task Struct Tests
// ==========================

func TestTaskStruct(t *testing.T) {
	_ = time.Now().UTC()
	task := Task{
		ID:                  "auth-task-abc123",
		FeatureID:           "auth-feature-abc",
		SpecID:              "test-cap-spec",
		BacklogID:           "auth",
		Name:                "Test Task",
		Goal:                string(make([]byte, 500)),
		Priority:            "P2",
		Status:              TaskStatusReady,
		DependsOn:           []string{"other-task"},
		IAEScenarios:        []string{"req-0001:scenario-0001"},
		ImplementationSteps: []string{"step1", "step2"},
		TestCases:           []string{"test1", "test2"},
		LibraryNeeds:        []string{"lib1"},
		ReadGates: ReadGates{
			IsReadBrief:        true,
			IsReadSpec:         true,
			IsReadSessionNotes: false,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if task.ID != "auth-task-abc123" {
		t.Errorf("Task.ID = %q, want %q", task.ID, "auth-task-abc123")
	}
	if task.Status != TaskStatusReady {
		t.Errorf("Task.Status = %q, want %q", task.Status, TaskStatusReady)
	}
	if len(task.ImplementationSteps) != 2 {
		t.Errorf("Task.ImplementationSteps count = %d, want 2", len(task.ImplementationSteps))
	}
}

func TestReadGates(t *testing.T) {
	tests := []struct {
		name    string
		gates   ReadGates
		isReady bool
	}{
		{"all gates true", ReadGates{true, true, true}, true},
		{"brief not read", ReadGates{false, true, true}, false},
		{"spec not read", ReadGates{true, false, true}, false},
		{"notes not read", ReadGates{true, true, false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isReady := tt.gates.IsReadBrief && tt.gates.IsReadSpec && tt.gates.IsReadSessionNotes
			if isReady != tt.isReady {
				t.Errorf("ReadGates ready state = %v, want %v", isReady, tt.isReady)
			}
		})
	}
}

func TestTaskCreateInput(t *testing.T) {
	input := TaskCreateInput{
		FeatureID:           "auth-feature-abc",
		SpecID:              "test-cap-spec",
		Name:                "New Task",
		Goal:                string(make([]byte, 500)),
		IAEScenarios:        []string{"req-0001:scenario-0001"},
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		LibraryNeeds:        []string{"lib1"},
		Priority:            "P3",
		DependsOn:           []string{},
	}

	if input.FeatureID != "auth-feature-abc" {
		t.Errorf("TaskCreateInput.FeatureID = %q, want %q", input.FeatureID, "auth-feature-abc")
	}
}

func TestTaskListInput(t *testing.T) {
	input := TaskListInput{
		FeatureID:      "auth-feature-abc",
		BacklogID:      "auth",
		Status:         TaskStatusReady,
		Priority:       "P2",
		IncludeDeleted: false,
		JSON:           true,
		Sort:           "created_at",
		Order:          "desc",
	}

	if input.FeatureID != "auth-feature-abc" {
		t.Errorf("TaskListInput.FeatureID = %q, want %q", input.FeatureID, "auth-feature-abc")
	}
}

func TestTaskDetailInput(t *testing.T) {
	input := TaskDetailInput{
		FeatureID:      "auth-feature-abc",
		TaskID:         "auth-task-abc123",
		JSON:           true,
		IncludeDeleted: false,
		Events:         true,
		Dependencies:   true,
		Timestamps:     true,
	}

	if input.TaskID != "auth-task-abc123" {
		t.Errorf("TaskDetailInput.TaskID = %q, want %q", input.TaskID, "auth-task-abc123")
	}
}

func TestTaskUpdateInput(t *testing.T) {
	name := "Updated Task"
	status := TaskStatusInProgress
	input := TaskUpdateInput{
		FeatureID: "auth-feature-abc",
		TaskID:    "auth-task-abc123",
		Name:      &name,
		Status:    &status,
		Reopen:    false,
		Cancel:    false,
		Force:     true,
		DryRun:    false,
	}

	if *input.Name != "Updated Task" {
		t.Errorf("TaskUpdateInput.Name = %q, want %q", *input.Name, "Updated Task")
	}
}

func TestTaskListItem(t *testing.T) {
	item := TaskListItem{
		ID:             "auth-task-abc123",
		Name:           "Test Task",
		Status:         TaskStatusReady,
		Priority:       "P2",
		FeatureID:      "auth-feature-abc",
		BacklogID:      "auth",
		DependsOnCount: 1,
		CreatedAt:      "2026-01-01T00:00:00Z",
		UpdatedAt:      "2026-01-02T00:00:00Z",
	}

	if item.DependsOnCount != 1 {
		t.Errorf("TaskListItem.DependsOnCount = %d, want 1", item.DependsOnCount)
	}
}

func TestTaskListOutput(t *testing.T) {
	output := TaskListOutput{
		Tasks: []TaskListItem{
			{ID: "t1"},
			{ID: "t2"},
		},
		Total:   2,
		Deleted: 0,
	}

	if len(output.Tasks) != 2 {
		t.Errorf("TaskListOutput.Tasks count = %d, want 2", len(output.Tasks))
	}
}

func TestTaskDetailOutput(t *testing.T) {
	output := TaskDetailOutput{
		ID:                  "auth-task-abc123",
		FeatureID:           "auth-feature-abc",
		BacklogID:           "auth",
		Name:                "Test Task",
		Goal:                "Test goal",
		Priority:            "P2",
		Status:              TaskStatusReady,
		DependsOn:           []string{},
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		LibraryNeeds:        []string{"lib1"},
		Events:              5,
		CreatedAt:           "2026-01-01T00:00:00Z",
		UpdatedAt:           "2026-01-02T00:00:00Z",
	}

	if output.Events != 5 {
		t.Errorf("TaskDetailOutput.Events = %d, want 5", output.Events)
	}
}

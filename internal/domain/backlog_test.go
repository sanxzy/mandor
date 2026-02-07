package domain

import (
	"testing"
	"time"
)

// ==========================
// Backlog Validation Tests
// ==========================

func TestValidateBacklogID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid simple", "auth", true},
		{"valid with numbers", "auth123", true},
		{"valid with hyphen", "my-backlog", true},
		{"valid with underscore", "my_backlog", true},
		{"valid mixed case", "Auth-Service_123", true},
		{"invalid starts with number", "123auth", false},
		{"invalid starts with hyphen", "-auth", false},
		{"invalid starts with underscore", "_auth", false},
		{"empty", "", false},
		{"invalid character @", "auth@123", false},
		{"invalid space", "auth 123", false},
		{"invalid dot", "auth.test", false},
		{"invalid slash", "auth/test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBacklogID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateBacklogID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestValidateGoalLength(t *testing.T) {
	tests := []struct {
		name     string
		goal     string
		minLen   int
		expected bool
	}{
		{"exactly min length", string(make([]byte, 500)), 500, true},
		{"exceeds min length", string(make([]byte, 501)), 500, true},
		{"well exceeds min", string(make([]byte, 1000)), 500, true},
		{"empty goal", "", 500, false},
		{"too short", "short", 500, false},
		{"300 chars with min 300", string(make([]byte, 300)), 300, true},
		{"299 chars with min 300", string(make([]byte, 299)), 300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateGoalLength(tt.goal, tt.minLen)
			if result != tt.expected {
				t.Errorf("ValidateGoalLength(len=%d, minLen=%d) = %v, want %v",
					len(tt.goal), tt.minLen, result, tt.expected)
			}
		})
	}
}

func TestValidateDependencyRule(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		expected bool
	}{
		{"valid same_backlog_only", "same_backlog_only", true},
		{"valid cross_backlog_allowed", "cross_backlog_allowed", true},
		{"valid disabled", "disabled", true},
		{"invalid empty", "", false},
		{"invalid random", "invalid_rule", false},
		{"invalid partial match", "same_backlog", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateDependencyRule(tt.rule)
			if result != tt.expected {
				t.Errorf("ValidateDependencyRule(%q) = %v, want %v", tt.rule, result, tt.expected)
			}
		})
	}
}

func TestValidateBooleanValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true lowercase", "true", true},
		{"true uppercase", "TRUE", true},
		{"false lowercase", "false", true},
		{"TRUE mixed case", "True", true},
		{"yes lowercase", "yes", true},
		{"no lowercase", "no", true},
		{"YES uppercase", "YES", true},
		{"NO uppercase", "NO", true},
		{"1", "1", true},
		{"0", "0", true},
		{"invalid maybe", "maybe", false},
		{"invalid empty", "", false},
		{"invalid truthy", "truthy", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBooleanValue(tt.value)
			if result != tt.expected {
				t.Errorf("ValidateBooleanValue(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestParseBooleanValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true lowercase", "true", true},
		{"TRUE uppercase", "TRUE", true},
		{"True mixed", "True", true},
		{"yes lowercase", "yes", true},
		{"YES uppercase", "YES", true},
		{"1", "1", true},
		{"false lowercase", "false", false},
		{"FALSE uppercase", "FALSE", false},
		{"no lowercase", "no", false},
		{"0", "0", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBooleanValue(tt.value)
			if result != tt.expected {
				t.Errorf("ParseBooleanValue(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestDefaultBacklogSchema(t *testing.T) {
	tests := []struct {
		name          string
		taskDep       string
		featureDep    string
		issueDep      string
		expectedTask  string
		expectedFeat  string
		expectedIssue string
	}{
		{"all empty defaults to same_backlog_only", "", "", "",
			"same_backlog_only", "cross_backlog_allowed", "same_backlog_only"},
		{"custom cross for task", "cross_backlog_allowed", "", "",
			"cross_backlog_allowed", "cross_backlog_allowed", "same_backlog_only"},
		{"custom disabled for issue", "", "", "disabled",
			"same_backlog_only", "cross_backlog_allowed", "disabled"},
		{"all custom", "cross_backlog_allowed", "same_backlog_only", "disabled",
			"cross_backlog_allowed", "same_backlog_only", "disabled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := DefaultBacklogSchema(tt.taskDep, tt.featureDep, tt.issueDep)

			if schema.Version != "mandor.v1" {
				t.Errorf("Version = %q, want %q", schema.Version, "mandor.v1")
			}
			if schema.Rules.Task.Dependency != tt.expectedTask {
				t.Errorf("Task.Dependency = %q, want %q", schema.Rules.Task.Dependency, tt.expectedTask)
			}
			if schema.Rules.Feature.Dependency != tt.expectedFeat {
				t.Errorf("Feature.Dependency = %q, want %q", schema.Rules.Feature.Dependency, tt.expectedFeat)
			}
			if schema.Rules.Issue.Dependency != tt.expectedIssue {
				t.Errorf("Issue.Dependency = %q, want %q", schema.Rules.Issue.Dependency, tt.expectedIssue)
			}
			if schema.Rules.Priority.Default != "P3" {
				t.Errorf("Priority.Default = %q, want %q", schema.Rules.Priority.Default, "P3")
			}
			if len(schema.Rules.Priority.Levels) != 6 {
				t.Errorf("Priority.Levels count = %d, want 6", len(schema.Rules.Priority.Levels))
			}
			if schema.Rules.Task.Cycle != CycleDisallowed {
				t.Errorf("Task.Cycle = %q, want %q", schema.Rules.Task.Cycle, CycleDisallowed)
			}
		})
	}
}

func TestBacklogStruct(t *testing.T) {
	now := time.Now().UTC()
	backlog := Backlog{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      string(make([]byte, 500)),
		Status:    BacklogStatusActive,
		Strict:    true,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if backlog.ID != "test-backlog" {
		t.Errorf("Backlog.ID = %q, want %q", backlog.ID, "test-backlog")
	}
	if backlog.Status != BacklogStatusActive {
		t.Errorf("Backlog.Status = %q, want %q", backlog.Status, BacklogStatusActive)
	}
}

func TestBacklogStatsStruct(t *testing.T) {
	stats := BacklogStats{
		Features: EntityStats{
			Total:        5,
			ByStatus:     map[string]int{"active": 3, "done": 2},
			AvgPriority:  "P2",
			BlockedCount: 1,
		},
		Tasks: EntityStats{
			Total:    10,
			ByStatus: map[string]int{"ready": 5, "in_progress": 5},
		},
		Issues: EntityStats{
			Total: 3,
		},
	}

	if stats.Features.Total != 5 {
		t.Errorf("Features.Total = %d, want 5", stats.Features.Total)
	}
	if stats.Tasks.Total != 10 {
		t.Errorf("Tasks.Total = %d, want 10", stats.Tasks.Total)
	}
	if stats.Features.BlockedCount != 1 {
		t.Errorf("Features.BlockedCount = %d, want 1", stats.Features.BlockedCount)
	}
}

// ==========================
// Backlog Status Constants Tests
// ==========================

func TestBacklogStatusConstants(t *testing.T) {
	if BacklogStatusInitial != "initial" {
		t.Errorf("BacklogStatusInitial = %q, want %q", BacklogStatusInitial, "initial")
	}
	if BacklogStatusActive != "active" {
		t.Errorf("BacklogStatusActive = %q, want %q", BacklogStatusActive, "active")
	}
	if BacklogStatusDone != "done" {
		t.Errorf("BacklogStatusDone = %q, want %q", BacklogStatusDone, "done")
	}
	if BacklogStatusDeleted != "deleted" {
		t.Errorf("BacklogStatusDeleted = %q, want %q", BacklogStatusDeleted, "deleted")
	}
}

func TestDependencyConstants(t *testing.T) {
	if DependencySameBacklogOnly != "same_backlog_only" {
		t.Errorf("DependencySameBacklogOnly = %q, want %q", DependencySameBacklogOnly, "same_backlog_only")
	}
	if DependencyCrossBacklogAllowed != "cross_backlog_allowed" {
		t.Errorf("DependencyCrossBacklogAllowed = %q, want %q", DependencyCrossBacklogAllowed, "cross_backlog_allowed")
	}
	if DependencyDisabled != "disabled" {
		t.Errorf("DependencyDisabled = %q, want %q", DependencyDisabled, "disabled")
	}
}

func TestCycleConstants(t *testing.T) {
	if CycleDisallowed != "disallowed" {
		t.Errorf("CycleDisallowed = %q, want %q", CycleDisallowed, "disallowed")
	}
	if CycleAllowed != "allowed" {
		t.Errorf("CycleAllowed = %q, want %q", CycleAllowed, "allowed")
	}
}

// ==========================
// Backlog Input/Output Tests
// ==========================

func TestBacklogCreateInput(t *testing.T) {
	input := BacklogCreateInput{
		ID:         "new-backlog",
		Name:       "New Backlog",
		Goal:       string(make([]byte, 500)),
		TaskDep:    "same_backlog_only",
		FeatureDep: "cross_backlog_allowed",
		IssueDep:   "same_backlog_only",
		Strict:     false,
	}

	if input.ID != "new-backlog" {
		t.Errorf("BacklogCreateInput.ID = %q, want %q", input.ID, "new-backlog")
	}
}

func TestBacklogUpdateInput(t *testing.T) {
	name := "Updated Name"
	strict := true
	input := BacklogUpdateInput{
		ID:     "test-backlog",
		Name:   &name,
		Strict: &strict,
	}

	if *input.Name != "Updated Name" {
		t.Errorf("BacklogUpdateInput.Name = %q, want %q", *input.Name, "Updated Name")
	}
	if *input.Strict != true {
		t.Errorf("BacklogUpdateInput.Strict = %v, want true", *input.Strict)
	}
}

func TestBacklogDeleteInput(t *testing.T) {
	input := BacklogDeleteInput{
		ID:   "test-backlog",
		Hard: true,
	}

	if input.ID != "test-backlog" {
		t.Errorf("BacklogDeleteInput.ID = %q, want %q", input.ID, "test-backlog")
	}
	if !input.Hard {
		t.Errorf("BacklogDeleteInput.Hard = false, want true")
	}
}

func TestBacklogReopenInput(t *testing.T) {
	input := BacklogReopenInput{
		ID:  "test-backlog",
		Yes: true,
	}

	if input.ID != "test-backlog" {
		t.Errorf("BacklogReopenInput.ID = %q, want %q", input.ID, "test-backlog")
	}
}

func TestBacklogListItem(t *testing.T) {
	item := BacklogListItem{
		ID:        "test-backlog",
		Name:      "Test",
		Goal:      "Test goal",
		Status:    BacklogStatusActive,
		Features:  5,
		Tasks:     10,
		Issues:    3,
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-02T00:00:00Z",
	}

	if item.Features != 5 {
		t.Errorf("BacklogListItem.Features = %d, want 5", item.Features)
	}
}

func TestBacklogListOutput(t *testing.T) {
	output := BacklogListOutput{
		Backlogs: []BacklogListItem{
			{ID: "backlog-1"},
			{ID: "backlog-2"},
		},
		Total:   2,
		Deleted: 0,
	}

	if len(output.Backlogs) != 2 {
		t.Errorf("BacklogListOutput.Backlogs count = %d, want 2", len(output.Backlogs))
	}
	if output.Total != 2 {
		t.Errorf("BacklogListOutput.Total = %d, want 2", output.Total)
	}
}

func TestBacklogDetailOutput(t *testing.T) {
	now := time.Now().UTC()
	output := BacklogDetailOutput{
		ID:        "test-backlog",
		Name:      "Test",
		Goal:      "Test goal",
		Status:    BacklogStatusActive,
		Strict:    false,
		Schema:    DefaultBacklogSchema("", "", ""),
		Stats:     BacklogStats{},
		Activity:  ActivityInfo{TotalEvents: 5},
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	if output.Activity.TotalEvents != 5 {
		t.Errorf("Activity.TotalEvents = %d, want 5", output.Activity.TotalEvents)
	}
}

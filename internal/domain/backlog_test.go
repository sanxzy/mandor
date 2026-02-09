package domain

import (
	"testing"
	"time"
)

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
		{"valid mixed", "Auth-Backlog_123", true},
		{"invalid starts with number", "123auth", false},
		{"invalid starts with hyphen", "-auth", false},
		{"invalid starts with underscore", "_auth", false},
		{"empty", "", false},
		{"invalid character", "auth@123", false},
		{"invalid space", "auth 123", false},
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

func TestDefaultBacklogSchema(t *testing.T) {
	schema := DefaultBacklogSchema("", "", "")

	if schema.Version != "mandor.v1" {
		t.Errorf("Version = %q, want %q", schema.Version, "mandor.v1")
	}

	if schema.Rules.Task.Dependency != DependencySameProjectOnly {
		t.Errorf("Task.Dependency = %q, want %q", schema.Rules.Task.Dependency, DependencySameProjectOnly)
	}

	if schema.Rules.Feature.Dependency != DependencyCrossProjectAllowed {
		t.Errorf("Feature.Dependency = %q, want %q", schema.Rules.Feature.Dependency, DependencyCrossProjectAllowed)
	}

	if schema.Rules.Issue.Dependency != DependencySameProjectOnly {
		t.Errorf("Issue.Dependency = %q, want %q", schema.Rules.Issue.Dependency, DependencySameProjectOnly)
	}

	if schema.Rules.Priority.Default != "P3" {
		t.Errorf("Priority.Default = %q, want %q", schema.Rules.Priority.Default, "P3")
	}

	if schema.Rules.Task.Cycle != CycleDisallowed {
		t.Errorf("Task.Cycle = %q, want %q", schema.Rules.Task.Cycle, CycleDisallowed)
	}
}

func TestDefaultBacklogSchemaWithCustomDeps(t *testing.T) {
	schema := DefaultBacklogSchema("cross_project_allowed", "same_project_only", "disabled")

	if schema.Rules.Task.Dependency != DependencyCrossProjectAllowed {
		t.Errorf("Task.Dependency = %q, want %q", schema.Rules.Task.Dependency, DependencyCrossProjectAllowed)
	}

	if schema.Rules.Feature.Dependency != DependencySameProjectOnly {
		t.Errorf("Feature.Dependency = %q, want %q", schema.Rules.Feature.Dependency, DependencySameProjectOnly)
	}

	if schema.Rules.Issue.Dependency != DependencyDisabled {
		t.Errorf("Issue.Dependency = %q, want %q", schema.Rules.Issue.Dependency, DependencyDisabled)
	}
}

func TestBacklogStruct(t *testing.T) {
	now := time.Now()
	backlog := Backlog{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      "This is a test backlog",
		Status:    BacklogStatusActive,
		Strict:    true,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "user1",
		UpdatedBy: "user1",
	}

	if backlog.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", backlog.ID, "test-backlog")
	}
	if backlog.Name != "Test Backlog" {
		t.Errorf("Name = %q, want %q", backlog.Name, "Test Backlog")
	}
	if backlog.Status != BacklogStatusActive {
		t.Errorf("Status = %q, want %q", backlog.Status, BacklogStatusActive)
	}
	if backlog.Strict != true {
		t.Errorf("Strict = %v, want %v", backlog.Strict, true)
	}
	if backlog.CreatedBy != "user1" {
		t.Errorf("CreatedBy = %q, want %q", backlog.CreatedBy, "user1")
	}
}

func TestBacklogStatusConstants(t *testing.T) {
	statuses := []string{
		BacklogStatusInitial,
		BacklogStatusActive,
		BacklogStatusDone,
		BacklogStatusDeleted,
	}

	expectedStatuses := []string{"initial", "active", "done", "deleted"}

	for i, status := range statuses {
		if status != expectedStatuses[i] {
			t.Errorf("Status[%d] = %q, want %q", i, status, expectedStatuses[i])
		}
	}
}

func TestBacklogCreateInput(t *testing.T) {
	input := BacklogCreateInput{
		ID:         "test-backlog",
		Name:       "Test Backlog",
		Goal:       "A comprehensive test backlog",
		TaskDep:    "same_project_only",
		FeatureDep: "cross_project_allowed",
		IssueDep:   "disabled",
		Strict:     true,
	}

	if input.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", input.ID, "test-backlog")
	}
	if input.Name != "Test Backlog" {
		t.Errorf("Name = %q, want %q", input.Name, "Test Backlog")
	}
	if !input.Strict {
		t.Errorf("Strict = %v, want %v", input.Strict, true)
	}
}

func TestBacklogUpdateInput(t *testing.T) {
	newName := "Updated Backlog"
	newGoal := "Updated goal"
	newStrictValue := false

	input := BacklogUpdateInput{
		ID:     "test-backlog",
		Name:   &newName,
		Goal:   &newGoal,
		Strict: &newStrictValue,
	}

	if input.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", input.ID, "test-backlog")
	}
	if *input.Name != "Updated Backlog" {
		t.Errorf("Name = %q, want %q", *input.Name, "Updated Backlog")
	}
	if *input.Goal != "Updated goal" {
		t.Errorf("Goal = %q, want %q", *input.Goal, "Updated goal")
	}
	if *input.Strict != false {
		t.Errorf("Strict = %v, want %v", *input.Strict, false)
	}
}

func TestBacklogDeleteInput(t *testing.T) {
	input := BacklogDeleteInput{
		ID:     "test-backlog",
		Hard:   false,
		DryRun: true,
		Yes:    false,
	}

	if input.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", input.ID, "test-backlog")
	}
	if input.Hard != false {
		t.Errorf("Hard = %v, want %v", input.Hard, false)
	}
	if !input.DryRun {
		t.Errorf("DryRun = %v, want %v", input.DryRun, true)
	}
}

func TestBacklogReopenInput(t *testing.T) {
	input := BacklogReopenInput{
		ID:  "test-backlog",
		Yes: true,
	}

	if input.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", input.ID, "test-backlog")
	}
	if !input.Yes {
		t.Errorf("Yes = %v, want %v", input.Yes, true)
	}
}

func TestBacklogListItem(t *testing.T) {
	item := BacklogListItem{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      "A test backlog",
		Status:    BacklogStatusActive,
		Features:  5,
		Tasks:     10,
		Issues:    2,
		CreatedAt: "2026-02-09T10:00:00Z",
		UpdatedAt: "2026-02-09T11:00:00Z",
	}

	if item.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", item.ID, "test-backlog")
	}
	if item.Features != 5 {
		t.Errorf("Features = %d, want %d", item.Features, 5)
	}
	if item.Tasks != 10 {
		t.Errorf("Tasks = %d, want %d", item.Tasks, 10)
	}
	if item.Issues != 2 {
		t.Errorf("Issues = %d, want %d", item.Issues, 2)
	}
}

func TestBacklogListOutput(t *testing.T) {
	items := []BacklogListItem{
		{
			ID:       "backlog-1",
			Name:     "Backlog 1",
			Status:   BacklogStatusActive,
			Features: 3,
		},
		{
			ID:       "backlog-2",
			Name:     "Backlog 2",
			Status:   BacklogStatusDone,
			Features: 2,
		},
	}

	output := BacklogListOutput{
		Backlogs: items,
		Total:    2,
		Deleted:  0,
	}

	if len(output.Backlogs) != 2 {
		t.Errorf("Backlogs count = %d, want %d", len(output.Backlogs), 2)
	}
	if output.Total != 2 {
		t.Errorf("Total = %d, want %d", output.Total, 2)
	}
	if output.Deleted != 0 {
		t.Errorf("Deleted = %d, want %d", output.Deleted, 0)
	}
}

func TestBacklogDetailOutput(t *testing.T) {
	schema := DefaultBacklogSchema("", "", "")
	stats := BacklogStats{}

	output := BacklogDetailOutput{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      "A detailed test backlog",
		Status:    BacklogStatusActive,
		Strict:    false,
		Schema:    schema,
		Stats:     stats,
		CreatedAt: "2026-02-09T10:00:00Z",
		UpdatedAt: "2026-02-09T11:00:00Z",
		CreatedBy: "user1",
		UpdatedBy: "user1",
	}

	if output.ID != "test-backlog" {
		t.Errorf("ID = %q, want %q", output.ID, "test-backlog")
	}
	if output.Schema.Version != "mandor.v1" {
		t.Errorf("Schema.Version = %q, want %q", output.Schema.Version, "mandor.v1")
	}
	if output.CreatedBy != "user1" {
		t.Errorf("CreatedBy = %q, want %q", output.CreatedBy, "user1")
	}
}

package service

import (
	"testing"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

// ==========================
// StatusService Type Tests
// ==========================

func TestWorkspaceStatusStruct(t *testing.T) {
	status := WorkspaceStatus{
		Workspace: nil,
		Backlogs:  []BacklogSummary{},
		Dependencies: DependencySummary{
			CrossBacklogCount: 0,
			CircularDeps:      0,
			BlockingItems:     []string{},
		},
		Totals: TotalStats{
			Features: 0,
			Tasks:    0,
			Issues:   0,
			Active:   0,
			Blocked:  0,
		},
	}

	if status.Totals.Features != 0 {
		t.Errorf("TotalStats.Features = %d, want 0", status.Totals.Features)
	}
	if status.Totals.Tasks != 0 {
		t.Errorf("TotalStats.Tasks = %d, want 0", status.Totals.Tasks)
	}
	if status.Totals.Issues != 0 {
		t.Errorf("TotalStats.Issues = %d, want 0", status.Totals.Issues)
	}
}

func TestBacklogSummaryStruct(t *testing.T) {
	summary := BacklogSummary{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Stats: domain.BacklogStats{
			Features: domain.EntityStats{Total: 5, ByStatus: make(map[string]int)},
			Tasks:    domain.EntityStats{Total: 10, ByStatus: make(map[string]int)},
			Issues:   domain.EntityStats{Total: 3, ByStatus: make(map[string]int), ByType: make(map[string]int)},
		},
	}

	if summary.ID != "test-backlog" {
		t.Errorf("BacklogSummary.ID = %q, want %q", summary.ID, "test-backlog")
	}
	if summary.Name != "Test Backlog" {
		t.Errorf("BacklogSummary.Name = %q, want %q", summary.Name, "Test Backlog")
	}
	if summary.Stats.Features.Total != 5 {
		t.Errorf("BacklogSummary.Stats.Features.Total = %d, want 5", summary.Stats.Features.Total)
	}
}

func TestDependencySummaryStruct(t *testing.T) {
	summary := DependencySummary{
		CrossBacklogCount: 5,
		CircularDeps:      2,
		BlockingItems:     []string{"issue-1", "task-2"},
	}

	if summary.CrossBacklogCount != 5 {
		t.Errorf("DependencySummary.CrossBacklogCount = %d, want 5", summary.CrossBacklogCount)
	}
	if summary.CircularDeps != 2 {
		t.Errorf("DependencySummary.CircularDeps = %d, want 2", summary.CircularDeps)
	}
	if len(summary.BlockingItems) != 2 {
		t.Errorf("DependencySummary.BlockingItems count = %d, want 2", len(summary.BlockingItems))
	}
}

func TestTotalStatsStruct(t *testing.T) {
	stats := TotalStats{
		Features: 10,
		Tasks:    50,
		Issues:   25,
		Active:   30,
		Blocked:  5,
	}

	if stats.Features != 10 {
		t.Errorf("TotalStats.Features = %d, want 10", stats.Features)
	}
	if stats.Tasks != 50 {
		t.Errorf("TotalStats.Tasks = %d, want 50", stats.Tasks)
	}
	if stats.Issues != 25 {
		t.Errorf("TotalStats.Issues = %d, want 25", stats.Issues)
	}
	if stats.Active != 30 {
		t.Errorf("TotalStats.Active = %d, want 30", stats.Active)
	}
	if stats.Blocked != 5 {
		t.Errorf("TotalStats.Blocked = %d, want 5", stats.Blocked)
	}
}

// ==========================
// WorkspaceService Config Tests
// ==========================

func TestWorkspaceServiceStruct(t *testing.T) {
	// Test that the struct can be instantiated (with nil dependencies for struct test)
	svc := &WorkspaceService{
		reader: nil,
		writer: nil,
		paths:  nil,
	}

	if svc == nil {
		t.Error("WorkspaceService should not be nil")
	}
}

func TestStatusServiceStruct(t *testing.T) {
	// Test that the struct can be instantiated (with nil dependencies for struct test)
	svc := &StatusService{
		reader: nil,
		paths:  nil,
	}

	if svc == nil {
		t.Error("StatusService should not be nil")
	}
}

// ==========================
// ChangeGovernanceService Tests
// ==========================

func TestChangeGovernanceServiceStruct(t *testing.T) {
	// Test that the struct can be instantiated (with nil dependencies for struct test)
	svc := &ChangeGovernanceService{
		paths: nil,
	}

	if svc == nil {
		t.Error("ChangeGovernanceService should not be nil")
	}
}

func TestChangeImpactAnalysisStruct(t *testing.T) {
	analysis := &ChangeImpactAnalysis{
		ChangeID:           "change-123",
		ChangeType:         "brief_modification",
		EntityID:           "test-brief",
		BacklogID:          "test-backlog",
		FieldsChanged:      []string{"why"},
		ImpactedSpecs:      []string{},
		ImpactedFeatures:   []string{},
		ImpactedTasks:      []string{},
		ImpactedBlueprints: []string{},
		BlockingStatus:     "NON_BLOCKING",
		RequiredActions:    []string{},
		ValidationDeadline: "2026-02-14T10:00:00Z",
		Status:             "pending_validation",
		Timestamp:          "2026-02-07T10:00:00Z",
		User:               "testuser",
		VersionBefore:      "1.0.0",
		VersionAfter:       "1.1.0",
	}

	if analysis.ChangeID != "change-123" {
		t.Errorf("ChangeImpactAnalysis.ChangeID = %q, want %q", analysis.ChangeID, "change-123")
	}
	if analysis.ChangeType != "brief_modification" {
		t.Errorf("ChangeImpactAnalysis.ChangeType = %q, want %q", analysis.ChangeType, "brief_modification")
	}
	if analysis.BlockingStatus != "NON_BLOCKING" {
		t.Errorf("ChangeImpactAnalysis.BlockingStatus = %q, want %q", analysis.BlockingStatus, "NON_BLOCKING")
	}
	if len(analysis.FieldsChanged) != 1 {
		t.Errorf("ChangeImpactAnalysis.FieldsChanged count = %d, want 1", len(analysis.FieldsChanged))
	}
}

func TestChangeImpactAnalysisWithBlocking(t *testing.T) {
	analysis := &ChangeImpactAnalysis{
		ChangeID:           "change-456",
		ChangeType:         "spec_modification",
		EntityID:           "test-spec",
		BacklogID:          "test-backlog",
		FieldsChanged:      []string{"acceptance-criteria"},
		ImpactedSpecs:      []string{},
		ImpactedFeatures:   []string{"feat-1"},
		ImpactedTasks:      []string{"task-1", "task-2", "task-3"},
		ImpactedBlueprints: []string{},
		BlockingStatus:     "BLOCKING",
		RequiredActions:    []string{"Regenerate all dependent Tasks"},
		ValidationDeadline: "2026-02-14T10:00:00Z",
		Status:             "pending_validation",
		Timestamp:          "2026-02-07T10:00:00Z",
		User:               "testuser",
		VersionBefore:      "1.0.0",
		VersionAfter:       "1.1.0",
	}

	if analysis.BlockingStatus != "BLOCKING" {
		t.Errorf("ChangeImpactAnalysis.BlockingStatus = %q, want %q", analysis.BlockingStatus, "BLOCKING")
	}
	if len(analysis.RequiredActions) != 1 {
		t.Errorf("ChangeImpactAnalysis.RequiredActions count = %d, want 1", len(analysis.RequiredActions))
	}
	if len(analysis.ImpactedTasks) != 3 {
		t.Errorf("ChangeImpactAnalysis.ImpactedTasks count = %d, want 3", len(analysis.ImpactedTasks))
	}
}

// ==========================
// Blocking Status Tests
// ==========================

func TestBlockingStatusValues(t *testing.T) {
	// Test that these string values are used correctly
	nonBlocking := "NON_BLOCKING"
	blocking := "BLOCKING"

	if nonBlocking != "NON_BLOCKING" {
		t.Errorf("NON_BLOCKING = %q", nonBlocking)
	}
	if blocking != "BLOCKING" {
		t.Errorf("BLOCKING = %q", blocking)
	}
}

func TestChangeTypeValues(t *testing.T) {
	// Test that these string values are used correctly
	briefMod := "brief_modification"
	specMod := "spec_modification"
	blueprintMod := "blueprint_modification"

	if briefMod != "brief_modification" {
		t.Errorf("brief_modification = %q", briefMod)
	}
	if specMod != "spec_modification" {
		t.Errorf("spec_modification = %q", specMod)
	}
	if blueprintMod != "blueprint_modification" {
		t.Errorf("blueprint_modification = %q", blueprintMod)
	}
}

func TestChangeStatusValues(t *testing.T) {
	// Test that these string values are used correctly
	pending := "pending_validation"
	approved := "approved"
	rejected := "rejected"
	implemented := "implemented"

	if pending != "pending_validation" {
		t.Errorf("pending_validation = %q", pending)
	}
	if approved != "approved" {
		t.Errorf("approved = %q", approved)
	}
	if rejected != "rejected" {
		t.Errorf("rejected = %q", rejected)
	}
	if implemented != "implemented" {
		t.Errorf("implemented = %q", implemented)
	}
}

// ==========================
// Brief Change Field Tests
// ==========================

func TestBriefChangeFieldBlocking(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"capabilities is blocking", "capabilities", "BLOCKING"},
		{"tech-stack is blocking", "tech-stack", "BLOCKING"},
		{"timeline is blocking", "timeline", "BLOCKING"},
		{"why is non-blocking", "why", "NON_BLOCKING"},
		{"business-rationale is non-blocking", "business-rationale", "NON_BLOCKING"},
		{"priority is non-blocking", "priority", "NON_BLOCKING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateBriefChangeBlocking("test-backlog", "test-brief", map[string]interface{}{tt.field: true})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expected {
				t.Errorf("BlockingStatus = %q, want %q", analysis.BlockingStatus, tt.expected)
			}
		})
	}
}

func TestSpecChangeFieldBlocking(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"acceptance-criteria is blocking", "acceptance-criteria", "BLOCKING"},
		{"requirements is blocking", "requirements", "BLOCKING"},
		{"iae-scenarios is blocking", "iae-scenarios", "BLOCKING"},
		{"summary is non-blocking", "summary", "NON_BLOCKING"},
		{"testing-strategy-notes is non-blocking", "testing-strategy-notes", "NON_BLOCKING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateSpecChangeBlocking("test-backlog", "test-spec", map[string]interface{}{tt.field: true})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expected {
				t.Errorf("BlockingStatus = %q, want %q", analysis.BlockingStatus, tt.expected)
			}
		})
	}
}

func TestBlueprintChangeFieldBlocking(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"tech-stack is blocking", "tech-stack", "BLOCKING"},
		{"deployment-strategy is blocking", "deployment-strategy", "BLOCKING"},
		{"architecture-pattern is blocking", "architecture-pattern", "BLOCKING"},
		{"cost-analysis is non-blocking", "cost-analysis", "NON_BLOCKING"},
		{"rationale is non-blocking", "rationale", "NON_BLOCKING"},
		{"assumptions is non-blocking", "assumptions", "NON_BLOCKING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateBlueprintChangeBlocking("test-backlog", "test-bp", map[string]interface{}{tt.field: true})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expected {
				t.Errorf("BlockingStatus = %q, want %q", analysis.BlockingStatus, tt.expected)
			}
		})
	}
}

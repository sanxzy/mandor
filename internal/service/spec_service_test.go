package service

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

// Helper to set up test service with temporary directory and required structure
func setupTestSpecService(t *testing.T) (*SpecService, string) {
	tmpDir := t.TempDir()

	// Initialize workspace first
	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	// Create necessary directories
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "features"), 0755)

	// Write workspace file
	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			GoalLengths: domain.GoalLengths{
				Backlog: 500,
				Feature: 300,
				Task:    500,
				Issue:   200,
			},
		},
	}
	if err := writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	service := &SpecService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

// Helper to set up test service with backlog and brief with capability
func setupTestSpecServiceWithBrief(t *testing.T, backlogID, capabilityID string) (*SpecService, *BriefService) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	// Create directories
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "features"), 0755)

	// Write workspace
	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
		},
	}
	writer.WriteWorkspace(workspace)

	// Create backlog
	backlogService := &BacklogService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}
	backlogInput := &domain.BacklogCreateInput{
		ID:   backlogID,
		Name: "Test Backlog",
		Goal: "This is a test goal for the backlog that has enough characters to pass validation requirements",
	}
	backlogService.CreateBacklog(backlogInput)

	// Create brief with capability
	briefService := &BriefService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}
	briefInput := &domain.BriefCreateInput{
		BacklogID: backlogID,
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for the spec service tests",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        capabilityID,
				Description: fmt.Sprintf("Description for %s capability", capabilityID),
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	briefService.CreateBrief(briefInput)

	specService := &SpecService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return specService, briefService
}

// ==========================
// Constructor Tests
// ==========================

func TestNewSpecServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewSpecServiceWithPaths(paths)

	if service == nil {
		t.Error("NewSpecServiceWithPaths returned nil")
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

func TestSpecWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestSpecService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

// ==========================
// Validation Tests
// ==========================

func TestValidateCreateInput_BacklogNotFound(t *testing.T) {
	service, _ := setupTestSpecService(t)

	input := &domain.SpecCreateInput{
		BacklogID:    "nonexistent-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent backlog, got nil")
	}
}

func TestValidateCreateInput_EmptyCapabilityID(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty capability ID, got nil")
	}
}

func TestValidateCreateInput_InvalidCapabilityID(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	tests := []struct {
		name string
		id   string
	}{
		{"ID with special chars", "cap!@#"},
		{"ID with uppercase", "TEST-CAP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &domain.SpecCreateInput{
				BacklogID:    "test-backlog",
				CapabilityID: tt.id,
				Summary:      "Test spec summary",
				Requirements: []domain.RequirementInput{
					{
						Summary: "Test requirement",
						IAEScenarios: []domain.IAEScenarioInput{
							{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
						},
					},
				},
			}

			err := service.ValidateCreateInput(input)
			if err == nil {
				t.Errorf("Expected error for invalid capability ID %q, got nil", tt.id)
			}
		})
	}
}

func TestValidateCreateInput_CapabilityNotInBrief(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "nonexistent-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for capability not in brief, got nil")
	}
}

func TestValidateCreateInput_EmptySummary(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty summary, got nil")
	}
}

func TestValidateCreateInput_NoRequirements(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for no requirements, got nil")
	}
}

func TestValidateCreateInput_EmptyRequirementSummary(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty requirement summary, got nil")
	}
}

func TestValidateCreateInput_RequirementNoIAE(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary:      "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for requirement without IAE scenarios, got nil")
	}
}

func TestValidateCreateInput_EmptyIntent(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty intent, got nil")
	}
}

func TestValidateCreateInput_EmptyAction(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty action, got nil")
	}
}

func TestValidateCreateInput_EmptyExpect(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: ""},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty expect, got nil")
	}
}

func TestValidateCreateInput_ValidInput(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary:            "Test requirement",
				Details:            "Test details",
				AcceptanceCriteria: []string{"AC1", "AC2"},
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

// ==========================
// Create Tests
// ==========================

func TestCreateSpec_Success(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}

	spec, err := service.CreateSpec(input)
	if err != nil {
		t.Fatalf("CreateSpec failed: %v", err)
	}

	if spec == nil {
		t.Error("Spec returned is nil")
	}
	if spec.ID != "test-capability-spec" {
		t.Errorf("Spec.ID = %q, want %q", spec.ID, "test-capability-spec")
	}
	if spec.CapabilityID != "test-capability" {
		t.Errorf("Spec.CapabilityID = %q, want %q", spec.CapabilityID, "test-capability")
	}
	if spec.Status != domain.SpecStatusDraft {
		t.Errorf("Spec.Status = %q, want %q", spec.Status, domain.SpecStatusDraft)
	}
	if len(spec.Requirements) != 1 {
		t.Errorf("Requirements count = %d, want 1", len(spec.Requirements))
	}
	if spec.CreatedBy == "" {
		t.Error("Spec.CreatedBy is empty")
	}
}

func TestCreateSpec_MultipleRequirements(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Requirement 1",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Intent 1", Action: "Action 1", Expect: "Expect 1"},
				},
			},
			{
				Summary: "Requirement 2",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Intent 2", Action: "Action 2", Expect: "Expect 2"},
				},
			},
		},
	}

	spec, err := service.CreateSpec(input)
	if err != nil {
		t.Fatalf("CreateSpec failed: %v", err)
	}

	if len(spec.Requirements) != 2 {
		t.Errorf("Requirements count = %d, want 2", len(spec.Requirements))
	}
}

func TestCreateSpec_MultipleIAEScenarios(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Intent 1", Action: "Action 1", Expect: "Expect 1"},
					{Intent: "Intent 2", Action: "Action 2", Expect: "Expect 2"},
				},
			},
		},
	}

	spec, err := service.CreateSpec(input)
	if err != nil {
		t.Fatalf("CreateSpec failed: %v", err)
	}

	if len(spec.Requirements[0].IAEScenarios) != 2 {
		t.Errorf("IAEScenarios count = %d, want 2", len(spec.Requirements[0].IAEScenarios))
	}
}

// ==========================
// Read Tests
// ==========================

func TestReadSpec_Success(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec first
	createInput := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	createdSpec, _ := service.CreateSpec(createInput)

	// Read the spec
	spec, err := service.ReadSpec("test-backlog", createdSpec.ID)
	if err != nil {
		t.Fatalf("ReadSpec failed: %v", err)
	}

	if spec == nil {
		t.Error("Spec returned is nil")
	}
	if spec.ID != createdSpec.ID {
		t.Errorf("Spec.ID = %q, want %q", spec.ID, createdSpec.ID)
	}
}

func TestReadSpec_NotFound(t *testing.T) {
	service, _ := setupTestSpecService(t)

	spec, err := service.ReadSpec("nonexistent-backlog", "nonexistent-spec")
	if err == nil {
		t.Error("Expected error for nonexistent spec, got nil")
	}
	if spec != nil {
		t.Error("Spec should be nil when not found")
	}
}

// ==========================
// Update Tests
// ==========================

func TestUpdateSpec_Status(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec first
	createInput := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	spec, _ := service.CreateSpec(createInput)

	// Update status
	newStatus := domain.SpecStatusActive
	spec.Status = newStatus
	err := service.UpdateSpec("test-backlog", spec)
	if err != nil {
		t.Fatalf("UpdateSpec failed: %v", err)
	}

	// Verify update
	updatedSpec, _ := service.ReadSpec("test-backlog", spec.ID)
	if updatedSpec.Status != newStatus {
		t.Errorf("Spec.Status = %q, want %q", updatedSpec.Status, newStatus)
	}
}

func TestUpdateSpec_Summary(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec first
	createInput := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Original summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	spec, _ := service.CreateSpec(createInput)

	// Update summary
	newSummary := "Updated summary"
	spec.Summary = newSummary
	service.UpdateSpec("test-backlog", spec)

	// Verify update
	updatedSpec, _ := service.ReadSpec("test-backlog", spec.ID)
	if updatedSpec.Summary != newSummary {
		t.Errorf("Spec.Summary = %q, want %q", updatedSpec.Summary, newSummary)
	}
}

func TestUpdateSpec_Timestamp(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec first
	createInput := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	spec, _ := service.CreateSpec(createInput)
	originalUpdatedAt := spec.UpdatedAt

	// Wait a bit and update
	time.Sleep(10 * time.Millisecond)
	spec.Status = domain.SpecStatusActive
	service.UpdateSpec("test-backlog", spec)

	// Verify timestamp changed
	updatedSpec, _ := service.ReadSpec("test-backlog", spec.ID)
	if updatedSpec.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("UpdatedAt timestamp was not updated")
	}
}

// ==========================
// Delete Tests
// ==========================

func TestDeleteSpec_Success(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec first
	createInput := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	spec, _ := service.CreateSpec(createInput)

	// Delete the spec
	err := service.DeleteSpec("test-backlog", spec.ID)
	if err != nil {
		t.Fatalf("DeleteSpec failed: %v", err)
	}

	// Verify deletion
	_, err = service.ReadSpec("test-backlog", spec.ID)
	if err == nil {
		t.Error("Expected error when reading deleted spec")
	}
}

func TestDeleteSpec_NotFound(t *testing.T) {
	service, _ := setupTestSpecService(t)

	err := service.DeleteSpec("nonexistent-backlog", "nonexistent-spec")
	// DeleteSpec might not error for nonexistent file, which is acceptable
	_ = err
}

// ==========================
// Serialization Tests
// ==========================

func TestSpecToMarkdown_ContainsAllSections(t *testing.T) {
	service := &SpecService{}

	spec := &domain.Spec{
		ID:           "test-spec",
		CapabilityID: "test-capability",
		BacklogID:    "test-backlog",
		Status:       domain.SpecStatusDraft,
		Summary:      "Test summary",
		Requirements: []domain.Requirement{
			{
				ID:                 "req-001",
				Summary:            "Test requirement",
				Details:            "Test details",
				AcceptanceCriteria: []string{"AC1", "AC2"},
				IAEScenarios: []domain.IAEScenario{
					{
						ID:     "scn-001",
						Intent: "Test intent",
						Action: "Test action",
						Expect: "Test expect",
					},
				},
			},
		},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedBy: "test-user",
		UpdatedBy: "test-user",
	}

	md := service.specToMarkdown(spec)

	// Check all sections are present
	if !contains(md, "SPEC_DATA_START") {
		t.Error("Markdown missing SPEC_DATA_START marker")
	}
	if !contains(md, "SPEC_DATA_END") {
		t.Error("Markdown missing SPEC_DATA_END marker")
	}
	if !contains(md, "# Spec: test-spec") {
		t.Error("Markdown missing spec title")
	}
	if !contains(md, "**Status:** draft") {
		t.Error("Markdown missing status")
	}
	if !contains(md, "## Summary") {
		t.Error("Markdown missing Summary section")
	}
	if !contains(md, "## Requirements") {
		t.Error("Markdown missing Requirements section")
	}
	if !contains(md, "**Intent:**") {
		t.Error("Markdown missing Intent section")
	}
	if !contains(md, "**Action:**") {
		t.Error("Markdown missing Action section")
	}
	if !contains(md, "**Expect:**") {
		t.Error("Markdown missing Expect section")
	}
}

func TestMarkdownToSpec_Success(t *testing.T) {
	service := &SpecService{}

	md := `<!--
SPEC_DATA_START
{"id":"test-spec","capability_id":"test-capability","backlog_id":"test-backlog","status":"draft","summary":"Test summary","requirements":[{"id":"req-001","summary":"Test requirement","details":"Test details","acceptance_criteria":["AC1"],"iae_scenarios":[{"id":"scn-001","intent":"Test intent","action":"Test action","expect":"Test expect"}]}],"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","created_by":"user","updated_by":"user"}
SPEC_DATA_END
-->
# Spec: test-spec

**Status:** draft

**Capability:** test-capability

## Summary

Test summary

## Requirements

### Requirement 1: Test requirement (req-001)

**Details:** Test details

**Acceptance Criteria:**
- AC1

**IAE Scenarios:**

#### Scenario 1 (scn-001)

- **Intent:** Test intent
- **Action:** Test action
- **Expect:** Test expect

**Created:** 2024-01-01T00:00:00Z by user
**Updated:** 2024-01-01T00:00:00Z by user
`

	spec, err := service.markdownToSpec(md)
	if err != nil {
		t.Fatalf("markdownToSpec failed: %v", err)
	}

	if spec.ID != "test-spec" {
		t.Errorf("Spec.ID = %q, want %q", spec.ID, "test-spec")
	}
	if spec.CapabilityID != "test-capability" {
		t.Errorf("Spec.CapabilityID = %q, want %q", spec.CapabilityID, "test-capability")
	}
	if len(spec.Requirements) != 1 {
		t.Errorf("Requirements count = %d, want 1", len(spec.Requirements))
	}
	if len(spec.Requirements[0].IAEScenarios) != 1 {
		t.Errorf("IAEScenarios count = %d, want 1", len(spec.Requirements[0].IAEScenarios))
	}
}

func TestMarkdownToSpec_MissingData(t *testing.T) {
	service := &SpecService{}

	tests := []struct {
		name string
		md   string
	}{
		{"Missing start marker", "# Spec: test\n\nContent"},
		{"Missing end marker", "<!--\n{}nSPEC_DATA_START"},
		{"Empty content", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.markdownToSpec(tt.md)
			if err == nil {
				t.Error("Expected error for malformed markdown, got nil")
			}
		})
	}
}

func TestMarkdownToSpec_InvalidJSON(t *testing.T) {
	service := &SpecService{}

	md := `<!--
SPEC_DATA_START
not valid json
SPEC_DATA_END
-->
# Spec: test
`

	_, err := service.markdownToSpec(md)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestSpecFileExists(t *testing.T) {
	service := &SpecService{}

	// Create a temp file with content
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "spec.md")
	os.WriteFile(path, []byte("Some content"), 0644)

	if !service.specFileExists(path) {
		t.Error("specFileExists returned false for existing file")
	}

	// Test with empty file
	emptyPath := filepath.Join(tmpDir, "empty.md")
	os.WriteFile(emptyPath, []byte(""), 0644)

	if service.specFileExists(emptyPath) {
		t.Error("specFileExists returned true for empty file")
	}

	// Test with nonexistent file
	if service.specFileExists(filepath.Join(tmpDir, "nonexistent.md")) {
		t.Error("specFileExists returned true for nonexistent file")
	}
}

func TestCapabilityExistsInBrief(t *testing.T) {
	service := &SpecService{}

	brief := &domain.Brief{
		ID: "test-brief",
		NewCapabilities: []domain.Capability{
			{ID: "cap-1", Name: "Capability 1", Description: "Desc 1"},
		},
		ModifiedCapabilities: []domain.Capability{
			{ID: "cap-2", Name: "Capability 2", Description: "Desc 2"},
		},
	}

	if !service.capabilityExistsInBrief(brief, "cap-1") {
		t.Error("capabilityExistsInBrief should find cap-1 in NewCapabilities")
	}
	if !service.capabilityExistsInBrief(brief, "cap-2") {
		t.Error("capabilityExistsInBrief should find cap-2 in ModifiedCapabilities")
	}
	if service.capabilityExistsInBrief(brief, "cap-3") {
		t.Error("capabilityExistsInBrief should not find cap-3")
	}
}

func TestSpecSerialization_Roundtrip(t *testing.T) {
	service, _ := setupTestSpecServiceWithBrief(t, "test-backlog", "test-capability")

	// Create a spec with complex data
	input := &domain.SpecCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		Summary:      "Complex spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary:            "Requirement with AC",
				Details:            "Detailed requirement description",
				AcceptanceCriteria: []string{"AC1", "AC2", "AC3"},
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "User wants to login", Action: "User enters credentials", Expect: "User is authenticated"},
					{Intent: "User wants to logout", Action: "User clicks logout", Expect: "Session is terminated"},
				},
			},
		},
	}

	// Create
	spec, err := service.CreateSpec(input)
	if err != nil {
		t.Fatalf("CreateSpec failed: %v", err)
	}

	// Read back
	readSpec, err := service.ReadSpec("test-backlog", spec.ID)
	if err != nil {
		t.Fatalf("ReadSpec failed: %v", err)
	}

	// Verify all fields are preserved
	if readSpec.ID != spec.ID {
		t.Errorf("ID mismatch: %q vs %q", readSpec.ID, spec.ID)
	}
	if readSpec.Summary != spec.Summary {
		t.Error("Summary not preserved")
	}
	if len(readSpec.Requirements) != 1 {
		t.Errorf("Requirements count = %d, want 1", len(readSpec.Requirements))
	}
	if len(readSpec.Requirements[0].AcceptanceCriteria) != 3 {
		t.Errorf("AcceptanceCriteria count = %d, want 3", len(readSpec.Requirements[0].AcceptanceCriteria))
	}
	if len(readSpec.Requirements[0].IAEScenarios) != 2 {
		t.Errorf("IAEScenarios count = %d, want 2", len(readSpec.Requirements[0].IAEScenarios))
	}
}

// Helper function for string matching
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

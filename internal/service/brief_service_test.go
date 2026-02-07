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

// Helper to set up test service with temporary directory
func setupTestBriefService(t *testing.T) (*BriefService, string) {
	tmpDir := t.TempDir()

	// Initialize workspace first
	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	// Create necessary directories (.mandor/backlogs, .mandor/briefs, etc.)
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

	service := &BriefService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

// Helper to set up test service with a backlog already created
func setupTestBriefServiceWithBacklog(t *testing.T, backlogID string) (*BriefService, *BacklogService) {
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

	// Create backlog service and create a backlog
	backlogService := &BacklogService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	backlogCreateInput := &domain.BacklogCreateInput{
		ID:   backlogID,
		Name: "Test Backlog",
		Goal: "This is a test goal for testing brief service with a proper backlog setup",
	}
	if err := backlogService.CreateBacklog(backlogCreateInput); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	briefService := &BriefService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return briefService, backlogService
}

// ==========================
// Constructor Tests
// ==========================

func TestNewBriefServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewBriefServiceWithPaths(paths)

	if service == nil {
		t.Error("NewBriefServiceWithPaths returned nil")
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

// ==========================
// Validation Tests
// ==========================

func TestBriefValidateCreateInput_NoBacklog(t *testing.T) {
	service, _ := setupTestBriefService(t)

	input := &domain.BriefCreateInput{
		BacklogID: "nonexistent-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent backlog, got nil")
	}
}

func TestBriefValidateCreateInput_EmptyName(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestBriefValidateCreateInput_InvalidName(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	tests := []struct {
		name  string
		input string
	}{
		{"Name with special chars", "Brief!@#$%"},
		{"Name with uppercase", "BRIEF_NAME"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &domain.BriefCreateInput{
				BacklogID: "test-backlog",
				Name:      tt.input,
				Why:       "This is a problem statement that is long enough to pass validation requirements",
				Capabilities: []domain.CapabilityInput{
					{
						Name:        "New Feature",
						Description: "This is a new feature description",
						Modified:    false,
					},
				},
				Impact: &domain.BriefImpact{},
			}

			err := service.ValidateCreateInput(input)
			if err == nil {
				t.Errorf("Expected error for invalid name %q, got nil", tt.input)
			}
		})
	}
}

func TestBriefValidateCreateInput_WhyTooShort(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "Too short",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for short Why section, got nil")
	}
}

func TestBriefValidateCreateInput_WhyTooLong(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a why string that is too long (> 5000 chars)
	longWhy := ""
	for i := 0; i < 501; i++ {
		longWhy += "This is a sentence. "
	}

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       longWhy,
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for long Why section, got nil")
	}
}

func TestBriefValidateCreateInput_NoCapabilities(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID:    "test-backlog",
		Name:         "Test Brief",
		Why:          "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{},
		Impact:       &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for no capabilities, got nil")
	}
}

func TestBriefValidateCreateInput_EmptyCapabilityName(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "",
				Description: "This is a description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty capability name, got nil")
	}
}

func TestBriefValidateCreateInput_InvalidCapabilityName(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "Feature!@#",
				Description: "This is a description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for invalid capability name, got nil")
	}
}

func TestBriefValidateCreateInput_EmptyCapabilityDescription(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty capability description, got nil")
	}
}

func TestBriefValidateCreateInput_BriefAlreadyExists(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create first brief
	firstInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "First Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for the first brief",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "Feature One",
				Description: "This is the first feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	service.CreateBrief(firstInput)

	// Try to create another brief for same backlog
	secondInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Second Brief",
		Why:       "This is another problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "Feature Two",
				Description: "This is the second feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(secondInput)
	if err == nil {
		t.Error("Expected error for duplicate brief, got nil")
	}
}

func TestBriefValidateCreateInput_ValidInput(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Valid Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for a valid brief test",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
			{
				Name:        "Modified Feature",
				Description: "This is a modified feature description",
				Modified:    true,
			},
		},
		Impact: &domain.BriefImpact{
			TechnicalStack:  []string{"Go", "PostgreSQL"},
			AffectedSystems: []string{"API", "Database"},
			Dependencies:    []string{"Service A", "Service B"},
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

func TestCreateBrief_Success(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for creating a brief",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	brief, err := service.CreateBrief(input)
	if err != nil {
		t.Fatalf("CreateBrief failed: %v", err)
	}

	if brief == nil {
		t.Error("Brief returned is nil")
	}
	if brief.ID == "" {
		t.Error("Brief ID is empty")
	}
	if brief.BacklogID != "test-backlog" {
		t.Errorf("Brief.BacklogID = %q, want test-backlog", brief.BacklogID)
	}
	if brief.Status != domain.BriefStatusDraft {
		t.Errorf("Brief.Status = %q, want draft", brief.Status)
	}
	if brief.CreatedBy == "" {
		t.Error("Brief.CreatedBy is empty")
	}
	if brief.UpdatedBy == "" {
		t.Error("Brief.UpdatedBy is empty")
	}
}

func TestCreateBrief_CapabilitiesSeparation(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for testing capability separation",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
			{
				Name:        "Modified Feature",
				Description: "This is a modified feature",
				Modified:    true,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	brief, err := service.CreateBrief(input)
	if err != nil {
		t.Fatalf("CreateBrief failed: %v", err)
	}

	if len(brief.NewCapabilities) != 1 {
		t.Errorf("NewCapabilities count = %d, want 1", len(brief.NewCapabilities))
	}
	if len(brief.ModifiedCapabilities) != 1 {
		t.Errorf("ModifiedCapabilities count = %d, want 1", len(brief.ModifiedCapabilities))
	}
	if len(brief.WhatChanges) != 2 {
		t.Errorf("WhatChanges count = %d, want 2", len(brief.WhatChanges))
	}
}

func TestCreateBrief_ImpactPreservation(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	impact := &domain.BriefImpact{
		TechnicalStack:  []string{"Go", "PostgreSQL"},
		AffectedSystems: []string{"API"},
		Dependencies:    []string{"Service A"},
	}

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for impact preservation testing",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: impact,
	}

	brief, err := service.CreateBrief(input)
	if err != nil {
		t.Fatalf("CreateBrief failed: %v", err)
	}

	if len(brief.Impact.TechnicalStack) != 2 {
		t.Errorf("TechnicalStack count = %d, want 2", len(brief.Impact.TechnicalStack))
	}
	if len(brief.Impact.AffectedSystems) != 1 {
		t.Errorf("AffectedSystems count = %d, want 1", len(brief.Impact.AffectedSystems))
	}
	if len(brief.Impact.Dependencies) != 1 {
		t.Errorf("Dependencies count = %d, want 1", len(brief.Impact.Dependencies))
	}
}

// ==========================
// Read Tests
// ==========================

func TestReadBrief_Success(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief first
	createInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for reading a brief",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	service.CreateBrief(createInput)

	// Read the brief
	brief, err := service.ReadBrief("test-backlog")
	if err != nil {
		t.Fatalf("ReadBrief failed: %v", err)
	}

	if brief == nil {
		t.Error("Brief returned is nil")
	}
	if brief.BacklogID != "test-backlog" {
		t.Errorf("Brief.BacklogID = %q, want test-backlog", brief.BacklogID)
	}
	if brief.Why != "This is a problem statement that is long enough to pass validation requirements for reading a brief" {
		t.Error("Brief.Why not preserved")
	}
}

func TestReadBrief_NotFound(t *testing.T) {
	service, _ := setupTestBriefService(t)

	brief, err := service.ReadBrief("nonexistent-backlog")
	if err == nil {
		t.Error("Expected error for nonexistent brief, got nil")
	}
	if brief != nil {
		t.Error("Brief should be nil when not found")
	}
}

// ==========================
// Update Tests
// ==========================

func TestUpdateBrief_Status(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief first
	createInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for updating brief status",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	brief, _ := service.CreateBrief(createInput)

	// Update status
	newStatus := domain.BriefStatusActive
	brief.Status = newStatus
	err := service.UpdateBrief("test-backlog", brief)
	if err != nil {
		t.Fatalf("UpdateBrief failed: %v", err)
	}

	// Verify update
	updatedBrief, _ := service.ReadBrief("test-backlog")
	if updatedBrief.Status != newStatus {
		t.Errorf("Brief.Status = %q, want %q", updatedBrief.Status, newStatus)
	}
}

func TestUpdateBrief_Metadata(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief first
	createInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for updating brief metadata",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	brief, _ := service.CreateBrief(createInput)

	// Update metadata
	newWhy := "This is an updated problem statement that is long enough to pass validation requirements for brief metadata"
	brief.Why = newWhy
	err := service.UpdateBrief("test-backlog", brief)
	if err != nil {
		t.Fatalf("UpdateBrief failed: %v", err)
	}

	// Verify update
	updatedBrief, _ := service.ReadBrief("test-backlog")
	if updatedBrief.Why != newWhy {
		t.Errorf("Brief.Why was not updated properly")
	}
}

func TestUpdateBrief_UpdatesTimestamp(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief first
	createInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for timestamp checking",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	brief, _ := service.CreateBrief(createInput)
	originalUpdatedAt := brief.UpdatedAt

	// Wait a bit and update
	time.Sleep(10 * time.Millisecond)
	brief.Status = domain.BriefStatusActive
	service.UpdateBrief("test-backlog", brief)

	// Verify timestamp changed
	updatedBrief, _ := service.ReadBrief("test-backlog")
	if updatedBrief.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("UpdatedAt timestamp was not updated")
	}
	if updatedBrief.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("UpdatedAt timestamp went backwards")
	}
}

// ==========================
// Delete Tests
// ==========================

func TestDeleteBrief_Success(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief first
	createInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for deleting a brief",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "New Feature",
				Description: "This is a new feature",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	service.CreateBrief(createInput)

	// Delete the brief
	err := service.DeleteBrief("test-backlog")
	if err != nil {
		t.Fatalf("DeleteBrief failed: %v", err)
	}

	// Verify deletion
	_, err = service.ReadBrief("test-backlog")
	if err == nil {
		t.Error("Expected error when reading deleted brief")
	}
}

func TestDeleteBrief_NonExistent(t *testing.T) {
	service, _ := setupTestBriefService(t)

	err := service.DeleteBrief("nonexistent-backlog")
	// DeleteBrief might not error for nonexistent file, which is acceptable
	// This test documents that behavior
	_ = err
}

// ==========================
// Workspace Initialization Tests
// ==========================

func TestBriefWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestBriefService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

// ==========================
// Serialization Tests (roundtrip)
// ==========================

func TestBriefSerialization_Roundtrip(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief with complex data
	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Complex Brief",
		Why:       "This is a detailed problem statement that explains why we need to make this change and what benefits it will bring to our system",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "API Enhancement",
				Description: "Enhance the REST API to support new endpoints",
				Modified:    false,
			},
			{
				Name:        "Database Optimization",
				Description: "Optimize database queries for better performance",
				Modified:    true,
			},
		},
		Impact: &domain.BriefImpact{
			TechnicalStack:  []string{"Go", "PostgreSQL", "Redis"},
			AffectedSystems: []string{"API", "Database", "Cache"},
			Dependencies:    []string{"Service A", "Service B", "External API"},
		},
	}

	// Create
	brief, err := service.CreateBrief(input)
	if err != nil {
		t.Fatalf("CreateBrief failed: %v", err)
	}

	// Read back
	readBrief, err := service.ReadBrief("test-backlog")
	if err != nil {
		t.Fatalf("ReadBrief failed: %v", err)
	}

	// Verify all fields are preserved
	if readBrief.ID != brief.ID {
		t.Errorf("ID mismatch: %q vs %q", readBrief.ID, brief.ID)
	}
	if readBrief.Why != brief.Why {
		t.Error("Why content not preserved")
	}
	if len(readBrief.NewCapabilities) != 1 {
		t.Errorf("NewCapabilities count = %d, want 1", len(readBrief.NewCapabilities))
	}
	if len(readBrief.ModifiedCapabilities) != 1 {
		t.Errorf("ModifiedCapabilities count = %d, want 1", len(readBrief.ModifiedCapabilities))
	}
	if len(readBrief.Impact.TechnicalStack) != 3 {
		t.Errorf("TechnicalStack count = %d, want 3", len(readBrief.Impact.TechnicalStack))
	}
	if len(readBrief.Impact.AffectedSystems) != 3 {
		t.Errorf("AffectedSystems count = %d, want 3", len(readBrief.Impact.AffectedSystems))
	}
	if len(readBrief.Impact.Dependencies) != 3 {
		t.Errorf("Dependencies count = %d, want 3", len(readBrief.Impact.Dependencies))
	}
}

// ==========================
// Edge Cases
// ==========================

func TestBriefValidateCreateInput_MinimalWhy(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a why string that is exactly 100 characters (minimum)
	minWhy := ""
	for i := 0; i < 10; i++ {
		minWhy += "Ten chars "
	}

	input := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Minimal Brief",
		Why:       minWhy,
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "Feature",
				Description: "Description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error for minimal valid Why, got: %v", err)
	}
}

func TestBriefCreateBrief_MultipleCapabilities(t *testing.T) {
	service, _ := setupTestBriefServiceWithBacklog(t, "test-backlog")

	// Create a brief with many capabilities
	capabilities := []domain.CapabilityInput{}
	for i := 1; i <= 5; i++ {
		numStr := fmt.Sprintf("%d", i)
		cap := domain.CapabilityInput{
			Name:        numStr + " Feature",
			Description: "This is feature " + numStr + " description",
			Modified:    i%2 == 0,
		}
		capabilities = append(capabilities, cap)
	}

	input := &domain.BriefCreateInput{
		BacklogID:    "test-backlog",
		Name:         "Multi Cap Brief",
		Why:          "This is a problem statement that is long enough to pass validation requirements for multiple capabilities",
		Capabilities: capabilities,
		Impact:       &domain.BriefImpact{},
	}

	brief, err := service.CreateBrief(input)
	if err != nil {
		t.Fatalf("CreateBrief failed: %v", err)
	}

	totalCaps := len(brief.NewCapabilities) + len(brief.ModifiedCapabilities)
	if totalCaps != 5 {
		t.Errorf("Total capabilities = %d, want 5", totalCaps)
	}
}

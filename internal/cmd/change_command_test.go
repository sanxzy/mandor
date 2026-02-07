package cmd

import (
	"strings"
	"testing"

	"mandor/internal/fs"
)

// TestChangeCommand_AnalyzePositiveCases tests successful analysis scenarios
func TestChangeCommand_AnalyzePositiveCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name           string
		request        *AnalyzeRequest
		expectedStatus string
	}{
		{
			name: "Brief blocking change",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"capabilities"},
			},
			expectedStatus: "BLOCKING",
		},
		{
			name: "Brief non-blocking change",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why"},
			},
			expectedStatus: "NON_BLOCKING",
		},
		{
			name: "Spec blocking change",
			request: &AnalyzeRequest{
				EntityType:   "spec",
				EntityID:     "test-spec",
				BacklogID:    "test-proj",
				FieldsChange: []string{"requirements"},
			},
			expectedStatus: "BLOCKING",
		},
		{
			name: "Spec non-blocking change",
			request: &AnalyzeRequest{
				EntityType:   "spec",
				EntityID:     "test-spec",
				BacklogID:    "test-proj",
				FieldsChange: []string{"summary"},
			},
			expectedStatus: "NON_BLOCKING",
		},
		{
			name: "Blueprint blocking change",
			request: &AnalyzeRequest{
				EntityType:   "blueprint",
				EntityID:     "test-bp",
				BacklogID:    "test-proj",
				FieldsChange: []string{"tech-stack"},
			},
			expectedStatus: "BLOCKING",
		},
		{
			name: "Blueprint non-blocking change",
			request: &AnalyzeRequest{
				EntityType:   "blueprint",
				EntityID:     "test-bp",
				BacklogID:    "test-proj",
				FieldsChange: []string{"cost-analysis"},
			},
			expectedStatus: "NON_BLOCKING",
		},
		{
			name: "Multiple fields with blocking field",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why", "capabilities", "priority"},
			},
			expectedStatus: "BLOCKING",
		},
		{
			name: "Multiple fields all non-blocking",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why", "business-rationale", "priority"},
			},
			expectedStatus: "NON_BLOCKING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Analyze(tt.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if response == nil {
				t.Fatal("response is nil")
			}

			if response.Analysis.BlockingStatus != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, response.Analysis.BlockingStatus)
			}

			if response.Analysis.ChangeID == "" {
				t.Error("change ID must not be empty")
			}

			if response.Analysis.EntityID != tt.request.EntityID {
				t.Error("entity ID mismatch")
			}
		})
	}
}

// TestChangeCommand_AnalyzeNegativeCases tests error scenarios
func TestChangeCommand_AnalyzeNegativeCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name        string
		request     *AnalyzeRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Missing entity type",
			request: &AnalyzeRequest{
				EntityType:   "",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why"},
			},
			expectError: true,
			errorMsg:    "entity type is required",
		},
		{
			name: "Invalid entity type",
			request: &AnalyzeRequest{
				EntityType:   "invalid",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why"},
			},
			expectError: true,
			errorMsg:    "invalid entity type",
		},
		{
			name: "Missing entity ID",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why"},
			},
			expectError: true,
			errorMsg:    "entity ID is required",
		},
		{
			name: "Missing backlog ID",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "",
				FieldsChange: []string{"why"},
			},
			expectError: true,
			errorMsg:    "backlog ID is required",
		},
		{
			name: "No fields specified",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{},
			},
			expectError: true,
			errorMsg:    "at least one field",
		},
		{
			name: "Only whitespace fields",
			request: &AnalyzeRequest{
				EntityType:   "brief",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"   ", "\t", "  "},
			},
			expectError: true,
			errorMsg:    "no fields",
		},
		{
			name: "Entity type case insensitivity - lowercase",
			request: &AnalyzeRequest{
				EntityType:   "BRIEF",
				EntityID:     "test-brief",
				BacklogID:    "test-proj",
				FieldsChange: []string{"why"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Analyze(tt.request)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if response != nil {
					t.Error("response should be nil on error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if response == nil {
					t.Fatal("response is nil")
				}
			}
		})
	}
}

// TestChangeCommand_ApprovePositiveCases tests successful approval scenarios
func TestChangeCommand_ApprovePositiveCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name    string
		request *ApproveRequest
	}{
		{
			name: "Approve with simple reason",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123456",
				Reason:    "All specs regenerated and validated",
				BacklogID: "test-proj",
			},
		},
		{
			name: "Approve with detailed reason",
			request: &ApproveRequest{
				ChangeID:  "spec-change-789012",
				Reason:    "Regenerated all dependent tasks with new acceptance criteria",
				BacklogID: "test-proj",
			},
		},
		{
			name: "Approve with long change ID",
			request: &ApproveRequest{
				ChangeID:  "blueprint-change-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Reason:    "Architecture reviewed and tasks updated accordingly",
				BacklogID: "test-proj",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Approve(tt.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if response == nil {
				t.Fatal("response is nil")
			}

			if !response.Approval.Approved {
				t.Error("approval status should be true")
			}

			if response.Approval.ChangeID != tt.request.ChangeID {
				t.Error("change ID mismatch in approval")
			}

			if response.Approval.ApprovalReason != tt.request.Reason {
				t.Error("reason mismatch in approval")
			}
		})
	}
}

// TestChangeCommand_ApproveNegativeCases tests approval error scenarios
func TestChangeCommand_ApproveNegativeCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name        string
		request     *ApproveRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Missing change ID",
			request: &ApproveRequest{
				ChangeID:  "",
				Reason:    "All specs regenerated",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "change ID is required",
		},
		{
			name: "Missing reason",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123",
				Reason:    "",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "approval reason is required",
		},
		{
			name: "Reason too short",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123",
				Reason:    "short",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "approval reason too short",
		},
		{
			name: "Reason with exactly 10 characters",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123",
				Reason:    "1234567890",
				BacklogID: "test-proj",
			},
			expectError: false,
		},
		{
			name: "Reason with 9 characters",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123",
				Reason:    "123456789",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "approval reason too short",
		},
		{
			name: "Change ID with only whitespace",
			request: &ApproveRequest{
				ChangeID:  "   ",
				Reason:    "All specs regenerated",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "change ID is required",
		},
		{
			name: "Reason with only whitespace",
			request: &ApproveRequest{
				ChangeID:  "brief-change-123",
				Reason:    "          ",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "approval reason is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Approve(tt.request)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if response != nil {
					t.Error("response should be nil on error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if response == nil {
					t.Fatal("response is nil")
				}
			}
		})
	}
}

// TestChangeCommand_RejectPositiveCases tests successful rejection scenarios
func TestChangeCommand_RejectPositiveCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name    string
		request *RejectRequest
	}{
		{
			name: "Reject with simple reason",
			request: &RejectRequest{
				ChangeID:  "brief-change-123456",
				Reason:    "Timeline impact too high",
				BacklogID: "test-proj",
			},
		},
		{
			name: "Reject with detailed reason",
			request: &RejectRequest{
				ChangeID:  "spec-change-789012",
				Reason:    "Major requirements change deferred to next sprint",
				BacklogID: "test-proj",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Reject(tt.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if response == nil {
				t.Fatal("response is nil")
			}

			if response.Status != "rejected" {
				t.Errorf("expected status 'rejected', got '%s'", response.Status)
			}

			if response.ChangeID != tt.request.ChangeID {
				t.Error("change ID mismatch in rejection")
			}
		})
	}
}

// TestChangeCommand_RejectNegativeCases tests rejection error scenarios
func TestChangeCommand_RejectNegativeCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name        string
		request     *RejectRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Missing change ID",
			request: &RejectRequest{
				ChangeID:  "",
				Reason:    "Timeline impact too high",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "change ID is required",
		},
		{
			name: "Missing reason",
			request: &RejectRequest{
				ChangeID:  "brief-change-123",
				Reason:    "",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "rejection reason is required",
		},
		{
			name: "Reason too short",
			request: &RejectRequest{
				ChangeID:  "brief-change-123",
				Reason:    "too short",
				BacklogID: "test-proj",
			},
			expectError: true,
			errorMsg:    "rejection reason too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Reject(tt.request)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if response != nil {
					t.Error("response should be nil on error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if response == nil {
					t.Fatal("response is nil")
				}
			}
		})
	}
}

// TestChangeCommand_ListChangesPositiveCases tests successful listing scenarios
func TestChangeCommand_ListChangesPositiveCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name    string
		request *ListChangesRequest
	}{
		{
			name: "List all changes",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "all",
				EntityType: "all",
			},
		},
		{
			name: "List pending changes for project",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "pending_validation",
				EntityType: "all",
			},
		},
		{
			name: "List approved brief changes",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "approved",
				EntityType: "brief",
			},
		},
		{
			name: "List rejected spec changes",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "rejected",
				EntityType: "spec",
			},
		},
		{
			name: "List blueprint changes",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "all",
				EntityType: "blueprint",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.ListChanges(tt.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if response == nil {
				t.Fatal("response is nil")
			}

			if response.Changes == nil {
				t.Error("changes list is nil")
			}
		})
	}
}

// TestChangeCommand_ListChangesNegativeCases tests list error scenarios
func TestChangeCommand_ListChangesNegativeCases(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	tests := []struct {
		name        string
		request     *ListChangesRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Missing backlog ID",
			request: &ListChangesRequest{
				BacklogID:  "",
				Status:     "all",
				EntityType: "all",
			},
			expectError: true,
			errorMsg:    "backlog ID is required",
		},
		{
			name: "Invalid status filter",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "invalid_status",
				EntityType: "all",
			},
			expectError: true,
			errorMsg:    "invalid status",
		},
		{
			name: "Invalid entity type filter",
			request: &ListChangesRequest{
				BacklogID:  "test-proj",
				Status:     "all",
				EntityType: "invalid_type",
			},
			expectError: true,
			errorMsg:    "invalid entity type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.ListChanges(tt.request)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if response != nil {
					t.Error("response should be nil on error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if response == nil {
					t.Fatal("response is nil")
				}
			}
		})
	}
}

// TestChangeCommand_PrintAnalysis tests formatting of analysis output
func TestChangeCommand_PrintAnalysis(t *testing.T) {
	paths, _ := fs.NewPaths()
	cmd := NewChangeCommand(paths)

	// Create a sample analysis
	analysis := &struct {
		ChangeID           string
		ChangeType         string
		EntityID           string
		BacklogID          string
		Status             string
		BlockingStatus     string
		Timestamp          string
		FieldsChanged      []string
		ImpactedSpecs      []string
		ImpactedFeatures   []string
		ImpactedTasks      []string
		ImpactedBlueprints []string
		RequiredActions    []string
		ValidationDeadline string
	}{
		ChangeID:           "brief-change-123",
		ChangeType:         "brief_modification",
		EntityID:           "test-brief",
		BacklogID:          "test-proj",
		Status:             "pending_validation",
		BlockingStatus:     "BLOCKING",
		Timestamp:          "2024-02-05T10:30:45Z",
		FieldsChanged:      []string{"capabilities"},
		ImpactedSpecs:      []string{"spec-1", "spec-2"},
		ImpactedFeatures:   []string{"feature-1"},
		ImpactedTasks:      []string{"task-1", "task-2"},
		ImpactedBlueprints: []string{"bp-1"},
		RequiredActions:    []string{"Regenerate all affected Specs"},
		ValidationDeadline: "2024-02-08T10:30:45Z",
	}

	// Create a ChangeImpactAnalysis from the sample
	analysisObj := &struct {
		ChangeID           string
		ChangeType         string
		EntityID           string
		BacklogID          string
		Status             string
		BlockingStatus     string
		Timestamp          string
		FieldsChanged      []string
		ImpactedSpecs      []string
		ImpactedFeatures   []string
		ImpactedTasks      []string
		ImpactedBlueprints []string
		RequiredActions    []string
		ValidationDeadline string
	}{
		ChangeID:           "brief-change-123",
		ChangeType:         "brief_modification",
		EntityID:           "test-brief",
		BacklogID:          "test-proj",
		Status:             "pending_validation",
		BlockingStatus:     "BLOCKING",
		Timestamp:          "2024-02-05T10:30:45Z",
		FieldsChanged:      []string{"capabilities"},
		ImpactedSpecs:      []string{"spec-1", "spec-2"},
		ImpactedFeatures:   []string{"feature-1"},
		ImpactedTasks:      []string{"task-1", "task-2"},
		ImpactedBlueprints: []string{"bp-1"},
		RequiredActions:    []string{"Regenerate all affected Specs"},
		ValidationDeadline: "2024-02-08T10:30:45Z",
	}

	_ = analysis
	_ = analysisObj

	// Note: PrintAnalysis requires ChangeImpactAnalysis type
	// This test validates structure is correct
	t.Run("Analyze output contains required fields", func(t *testing.T) {
		req := &AnalyzeRequest{
			EntityType:   "brief",
			EntityID:     "test-brief",
			BacklogID:    "test-proj",
			FieldsChange: []string{"capabilities"},
		}

		response, err := cmd.Analyze(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := cmd.PrintAnalysis(response.Analysis)

		// Verify output contains key sections
		requiredStrings := []string{
			"CHANGE IMPACT ANALYSIS",
			"Change ID:",
			"Blocking Status:",
			"BLOCKING",
			"Fields Changed:",
			"Impacted Entities:",
		}

		for _, s := range requiredStrings {
			if !strings.Contains(output, s) {
				t.Errorf("output missing required section: %s", s)
			}
		}
	})
}

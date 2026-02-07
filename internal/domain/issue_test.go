package domain

import (
	"testing"
	"time"
)

// ==========================
// Issue Validation Tests
// ==========================

func TestValidateIssueID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid 12 chars", "auth-issue-xy", true},
		{"valid long ID", "auth-issue-abc123def", true},
		{"empty", "", false},
		{"valid 12 chars", "auth-issue-x", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIssueID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateIssueID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestValidateIssueStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"open", IssueStatusOpen, true},
		{"ready", IssueStatusReady, true},
		{"in_progress", IssueStatusInProgress, true},
		{"blocked", IssueStatusBlocked, true},
		{"resolved", IssueStatusResolved, true},
		{"wontfix", IssueStatusWontFix, true},
		{"cancelled", IssueStatusCancelled, true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIssueStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidateIssueStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestValidateIssueType(t *testing.T) {
	tests := []struct {
		name      string
		issueType string
		expected  bool
	}{
		{"bug", IssueTypeBug, true},
		{"improvement", IssueTypeImprovement, true},
		{"debt", IssueTypeDebt, true},
		{"security", IssueTypeSecurity, true},
		{"performance", IssueTypePerformance, true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIssueType(tt.issueType)
			if result != tt.expected {
				t.Errorf("ValidateIssueType(%q) = %v, want %v", tt.issueType, result, tt.expected)
			}
		})
	}
}

func TestValidateIssueGoalLength(t *testing.T) {
	tests := []struct {
		name     string
		goal     string
		minLen   int
		expected bool
	}{
		{"exactly min 200", string(make([]byte, 200)), 200, true},
		{"exceeds min 200", string(make([]byte, 201)), 200, true},
		{"too short", "short", 200, false},
		{"empty", "", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIssueGoalLength(tt.goal, tt.minLen)
			if result != tt.expected {
				t.Errorf("ValidateIssueGoalLength(len=%d, minLen=%d) = %v, want %v",
					len(tt.goal), tt.minLen, result, tt.expected)
			}
		})
	}
}

func TestIsIssueTerminalStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"resolved is terminal", IssueStatusResolved, true},
		{"wontfix is terminal", IssueStatusWontFix, true},
		{"cancelled is terminal", IssueStatusCancelled, true},
		{"open is not terminal", IssueStatusOpen, false},
		{"ready is not terminal", IssueStatusReady, false},
		{"in_progress is not terminal", IssueStatusInProgress, false},
		{"blocked is not terminal", IssueStatusBlocked, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIssueTerminalStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsIssueTerminalStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

// ==========================
// Issue Constants Tests
// ==========================

func TestIssueStatusConstants(t *testing.T) {
	if IssueStatusOpen != "open" {
		t.Errorf("IssueStatusOpen = %q, want %q", IssueStatusOpen, "open")
	}
	if IssueStatusReady != "ready" {
		t.Errorf("IssueStatusReady = %q, want %q", IssueStatusReady, "ready")
	}
	if IssueStatusInProgress != "in_progress" {
		t.Errorf("IssueStatusInProgress = %q, want %q", IssueStatusInProgress, "in_progress")
	}
	if IssueStatusBlocked != "blocked" {
		t.Errorf("IssueStatusBlocked = %q, want %q", IssueStatusBlocked, "blocked")
	}
	if IssueStatusResolved != "resolved" {
		t.Errorf("IssueStatusResolved = %q, want %q", IssueStatusResolved, "resolved")
	}
	if IssueStatusWontFix != "wontfix" {
		t.Errorf("IssueStatusWontFix = %q, want %q", IssueStatusWontFix, "wontfix")
	}
	if IssueStatusCancelled != "cancelled" {
		t.Errorf("IssueStatusCancelled = %q, want %q", IssueStatusCancelled, "cancelled")
	}
}

func TestIssueTypeConstants(t *testing.T) {
	if IssueTypeBug != "bug" {
		t.Errorf("IssueTypeBug = %q, want %q", IssueTypeBug, "bug")
	}
	if IssueTypeImprovement != "improvement" {
		t.Errorf("IssueTypeImprovement = %q, want %q", IssueTypeImprovement, "improvement")
	}
	if IssueTypeDebt != "debt" {
		t.Errorf("IssueTypeDebt = %q, want %q", IssueTypeDebt, "debt")
	}
	if IssueTypeSecurity != "security" {
		t.Errorf("IssueTypeSecurity = %q, want %q", IssueTypeSecurity, "security")
	}
	if IssueTypePerformance != "performance" {
		t.Errorf("IssueTypePerformance = %q, want %q", IssueTypePerformance, "performance")
	}
}

func TestIssueGoalMinLength(t *testing.T) {
	if IssueGoalMinLength != 200 {
		t.Errorf("IssueGoalMinLength = %d, want 200", IssueGoalMinLength)
	}
}

// ==========================
// Issue Struct Tests
// ==========================

func TestIssueStruct(t *testing.T) {
	now := time.Now().UTC()
	issue := Issue{
		ID:                  "auth-issue-xyz123",
		BacklogID:           "auth",
		Name:                "Test Issue",
		Goal:                string(make([]byte, 200)),
		IssueType:           IssueTypeBug,
		Priority:            "P2",
		Status:              IssueStatusOpen,
		DependsOn:           []string{"other-issue"},
		AffectedFiles:       []string{"src/file.ts"},
		AffectedTests:       []string{"tests/file.test.ts"},
		ImplementationSteps: []string{"step1"},
		LibraryNeeds:        []string{"lib1"},
		CreatedAt:           now,
		LastUpdatedAt:       now,
		CreatedBy:           "testuser",
		LastUpdatedBy:       "testuser",
	}

	if issue.ID != "auth-issue-xyz123" {
		t.Errorf("Issue.ID = %q, want %q", issue.ID, "auth-issue-xyz123")
	}
	if issue.IssueType != IssueTypeBug {
		t.Errorf("Issue.IssueType = %q, want %q", issue.IssueType, IssueTypeBug)
	}
}

func TestIssueCreateInput(t *testing.T) {
	input := IssueCreateInput{
		BacklogID:           "auth",
		Name:                "New Issue",
		Goal:                string(make([]byte, 200)),
		IssueType:           IssueTypeSecurity,
		Priority:            "P1",
		DependsOn:           []string{},
		AffectedFiles:       []string{"src/auth.ts"},
		AffectedTests:       []string{"tests/auth.test.ts"},
		ImplementationSteps: []string{"fix vulnerability"},
		LibraryNeeds:        []string{},
	}

	if input.IssueType != IssueTypeSecurity {
		t.Errorf("IssueCreateInput.IssueType = %q, want %q", input.IssueType, IssueTypeSecurity)
	}
}

func TestIssueListInput(t *testing.T) {
	input := IssueListInput{
		BacklogID:      "auth",
		IssueType:      IssueTypeBug,
		Status:         IssueStatusOpen,
		Priority:       "P2",
		IncludeDeleted: false,
		JSON:           true,
		Sort:           "created_at",
		Order:          "desc",
	}

	if input.IssueType != IssueTypeBug {
		t.Errorf("IssueListInput.IssueType = %q, want %q", input.IssueType, IssueTypeBug)
	}
}

func TestIssueDetailInput(t *testing.T) {
	input := IssueDetailInput{
		BacklogID:      "auth",
		IssueID:        "auth-issue-xyz",
		JSON:           true,
		IncludeDeleted: false,
		Events:         true,
		Dependencies:   true,
		Timestamps:     true,
	}

	if input.IssueID != "auth-issue-xyz" {
		t.Errorf("IssueDetailInput.IssueID = %q, want %q", input.IssueID, "auth-issue-xyz")
	}
}

func TestIssueUpdateInput(t *testing.T) {
	status := IssueStatusResolved
	input := IssueUpdateInput{
		BacklogID: "auth",
		IssueID:   "auth-issue-xyz",
		Status:    &status,
		Resolve:   true,
		Reopen:    false,
		Cancel:    false,
		Force:     false,
		DryRun:    false,
	}

	if *input.Status != IssueStatusResolved {
		t.Errorf("IssueUpdateInput.Status = %q, want %q", *input.Status, IssueStatusResolved)
	}
}

func TestIssueListItem(t *testing.T) {
	item := IssueListItem{
		ID:                       "auth-issue-xyz",
		Name:                     "Test Issue",
		IssueType:                IssueTypeBug,
		Status:                   IssueStatusOpen,
		Priority:                 "P2",
		BacklogID:                "auth",
		DependsOnCount:           1,
		AffectedFilesCount:       2,
		AffectedTestsCount:       1,
		ImplementationStepsCount: 3,
		LibraryNeedsCount:        1,
		CreatedAt:                "2026-01-01T00:00:00Z",
		LastUpdatedAt:            "2026-01-02T00:00:00Z",
	}

	if item.DependsOnCount != 1 {
		t.Errorf("IssueListItem.DependsOnCount = %d, want 1", item.DependsOnCount)
	}
}

func TestIssueListOutput(t *testing.T) {
	output := IssueListOutput{
		Issues: []IssueListItem{
			{ID: "i1"},
			{ID: "i2"},
		},
		Total:   2,
		Deleted: 0,
	}

	if len(output.Issues) != 2 {
		t.Errorf("IssueListOutput.Issues count = %d, want 2", len(output.Issues))
	}
}

func TestIssueDetailOutput(t *testing.T) {
	output := IssueDetailOutput{
		ID:                  "auth-issue-xyz",
		BacklogID:           "auth",
		Name:                "Test Issue",
		Goal:                "Test goal",
		IssueType:           IssueTypeBug,
		Priority:            "P2",
		Status:              IssueStatusOpen,
		DependsOn:           []string{},
		AffectedFiles:       []string{"src/file.ts"},
		AffectedTests:       []string{"tests/file.test.ts"},
		ImplementationSteps: []string{"step1"},
		LibraryNeeds:        []string{"lib1"},
		Events:              5,
		CreatedAt:           "2026-01-01T00:00:00Z",
		LastUpdatedAt:       "2026-01-02T00:00:00Z",
	}

	if output.Events != 5 {
		t.Errorf("IssueDetailOutput.Events = %d, want 5", output.Events)
	}
}

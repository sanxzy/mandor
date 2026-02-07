package domain

import (
	"testing"
	"time"
)

// ==========================
// Feature Validation Tests
// ==========================

func TestValidateFeatureID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid 12 chars", "auth-feature", true},
		{"valid long ID", "auth-feature-abc123", true},
		{"valid with numbers", "abc123def456", true},
		{"empty", "", false},
		{"too short 11 chars", "auth-featur", false},
		{"exactly 12 chars", "ab-cd-ef-ghij", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFeatureID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateFeatureID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestValidateScope(t *testing.T) {
	tests := []struct {
		name     string
		scope    string
		expected bool
	}{
		{"frontend", "frontend", true},
		{"backend", "backend", true},
		{"fullstack", "fullstack", true},
		{"cli", "cli", true},
		{"desktop", "desktop", true},
		{"android", "android", true},
		{"flutter", "flutter", true},
		{"react-native", "react-native", true},
		{"ios", "ios", true},
		{"swift", "swift", true},
		{"empty scope", "", true},
		{"invalid scope", "invalid", false},
		{"partial match", "front", false},
		{"uppercase FRONTEND", "FRONTEND", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateScope(tt.scope)
			if result != tt.expected {
				t.Errorf("ValidateScope(%q) = %v, want %v", tt.scope, result, tt.expected)
			}
		})
	}
}

func TestValidateFeatureStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"draft", FeatureStatusDraft, true},
		{"active", FeatureStatusActive, true},
		{"done", FeatureStatusDone, true},
		{"blocked", FeatureStatusBlocked, true},
		{"cancelled", FeatureStatusCancelled, true},
		{"invalid status", "invalid", false},
		{"empty", "", false},
		{"partial match", "drafts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFeatureStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidateFeatureStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestValidateFeatureGoalLength(t *testing.T) {
	tests := []struct {
		name     string
		goal     string
		minLen   int
		expected bool
	}{
		{"exactly min 300", string(make([]byte, 300)), 300, true},
		{"exceeds min 300", string(make([]byte, 301)), 300, true},
		{"too short", "short", 300, false},
		{"empty", "", 300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateFeatureGoalLength(tt.goal, tt.minLen)
			if result != tt.expected {
				t.Errorf("ValidateFeatureGoalLength(len=%d, minLen=%d) = %v, want %v",
					len(tt.goal), tt.minLen, result, tt.expected)
			}
		})
	}
}

// ==========================
// Feature Status Constants Tests
// ==========================

func TestFeatureStatusConstants(t *testing.T) {
	if FeatureStatusDraft != "draft" {
		t.Errorf("FeatureStatusDraft = %q, want %q", FeatureStatusDraft, "draft")
	}
	if FeatureStatusActive != "active" {
		t.Errorf("FeatureStatusActive = %q, want %q", FeatureStatusActive, "active")
	}
	if FeatureStatusDone != "done" {
		t.Errorf("FeatureStatusDone = %q, want %q", FeatureStatusDone, "done")
	}
	if FeatureStatusBlocked != "blocked" {
		t.Errorf("FeatureStatusBlocked = %q, want %q", FeatureStatusBlocked, "blocked")
	}
	if FeatureStatusCancelled != "cancelled" {
		t.Errorf("FeatureStatusCancelled = %q, want %q", FeatureStatusCancelled, "cancelled")
	}
}

func TestFeatureGoalMinLength(t *testing.T) {
	if FeatureGoalMinLength != 300 {
		t.Errorf("FeatureGoalMinLength = %d, want 300", FeatureGoalMinLength)
	}
}

// ==========================
// Feature Struct Tests
// ==========================

func TestFeatureStruct(t *testing.T) {
	now := time.Now().UTC()
	feature := Feature{
		ID:           "auth-feature",
		CapabilityID: "test-cap",
		SpecID:       "test-cap-spec",
		BacklogID:    "auth",
		Name:         "Test Feature",
		Goal:         string(make([]byte, 300)),
		Scope:        "fullstack",
		Priority:     "P2",
		Status:       FeatureStatusActive,
		DependsOn:    []string{"other-feature"},
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    "testuser",
		UpdatedBy:    "testuser",
	}

	if feature.ID != "auth-feature" {
		t.Errorf("Feature.ID = %q, want %q", feature.ID, "auth-feature")
	}
	if feature.Status != FeatureStatusActive {
		t.Errorf("Feature.Status = %q, want %q", feature.Status, FeatureStatusActive)
	}
	if len(feature.DependsOn) != 1 {
		t.Errorf("Feature.DependsOn count = %d, want 1", len(feature.DependsOn))
	}
}

func TestFeatureCreateInput(t *testing.T) {
	input := FeatureCreateInput{
		BacklogID:    "auth",
		CapabilityID: "test-cap",
		SpecID:       "test-cap-spec",
		Name:         "New Feature",
		Goal:         string(make([]byte, 300)),
		Scope:        "frontend",
		Priority:     "P1",
		DependsOn:    []string{},
	}

	if input.BacklogID != "auth" {
		t.Errorf("FeatureCreateInput.BacklogID = %q, want %q", input.BacklogID, "auth")
	}
}

func TestFeatureListInput(t *testing.T) {
	input := FeatureListInput{
		BacklogID:      "auth",
		Scope:          "fullstack",
		IncludeDeleted: true,
		JSON:           true,
	}

	if input.BacklogID != "auth" {
		t.Errorf("FeatureListInput.BacklogID = %q, want %q", input.BacklogID, "auth")
	}
	if !input.IncludeDeleted {
		t.Error("FeatureListInput.IncludeDeleted = false, want true")
	}
}

func TestFeatureDetailInput(t *testing.T) {
	input := FeatureDetailInput{
		BacklogID:      "auth",
		FeatureID:      "auth-feature-abc",
		JSON:           true,
		IncludeDeleted: false,
	}

	if input.FeatureID != "auth-feature-abc" {
		t.Errorf("FeatureDetailInput.FeatureID = %q, want %q", input.FeatureID, "auth-feature-abc")
	}
}

func TestFeatureUpdateInput(t *testing.T) {
	name := "Updated Name"
	input := FeatureUpdateInput{
		BacklogID: "auth",
		FeatureID: "auth-feature-abc",
		Name:      &name,
		Status:    nil,
		Reopen:    false,
		Cancel:    true,
		DryRun:    true,
	}

	if *input.Name != "Updated Name" {
		t.Errorf("FeatureUpdateInput.Name = %q, want %q", *input.Name, "Updated Name")
	}
}

func TestFeatureDeleteInput(t *testing.T) {
	input := FeatureDeleteInput{
		BacklogID: "auth",
		FeatureID: "auth-feature-abc",
		Force:     true,
		Reason:    "No longer needed",
	}

	if input.FeatureID != "auth-feature-abc" {
		t.Errorf("FeatureDeleteInput.FeatureID = %q, want %q", input.FeatureID, "auth-feature-abc")
	}
}

func TestFeatureListItem(t *testing.T) {
	item := FeatureListItem{
		ID:        "auth-feature-abc",
		Name:      "Test",
		Goal:      "Test goal",
		Scope:     "fullstack",
		Priority:  "P2",
		Status:    FeatureStatusActive,
		DependsOn: 3,
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-02T00:00:00Z",
	}

	if item.DependsOn != 3 {
		t.Errorf("FeatureListItem.DependsOn = %d, want 3", item.DependsOn)
	}
}

func TestFeatureListOutput(t *testing.T) {
	output := FeatureListOutput{
		Features: []FeatureListItem{
			{ID: "f1"},
			{ID: "f2"},
		},
		Total:   2,
		Deleted: 0,
	}

	if len(output.Features) != 2 {
		t.Errorf("FeatureListOutput.Features count = %d, want 2", len(output.Features))
	}
}

func TestFeatureDetailOutput(t *testing.T) {
	output := FeatureDetailOutput{
		ID:           "auth-feature-abc",
		CapabilityID: "test-cap",
		SpecID:       "test-cap-spec",
		BacklogID:    "auth",
		Name:         "Test",
		Goal:         "Test goal",
		Scope:        "fullstack",
		Priority:     "P2",
		Status:       FeatureStatusActive,
		DependsOn:    []string{},
		Events:       5,
	}

	if output.Events != 5 {
		t.Errorf("FeatureDetailOutput.Events = %d, want 5", output.Events)
	}
}

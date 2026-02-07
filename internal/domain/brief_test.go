package domain

import (
	"testing"
	"time"
)

// ==========================
// Brief Validation Tests
// ==========================

func TestValidateBriefWhy(t *testing.T) {
	tests := []struct {
		name     string
		why      string
		expected bool
	}{
		{"exactly 100 chars", string(make([]byte, 100)), true},
		{"500 chars", string(make([]byte, 500)), true},
		{"5000 chars", string(make([]byte, 5000)), true},
		{"99 chars", string(make([]byte, 99)), false},
		{"5001 chars", string(make([]byte, 5001)), false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBriefWhy(tt.why)
			if result != tt.expected {
				t.Errorf("ValidateBriefWhy(len=%d) = %v, want %v", len(tt.why), result, tt.expected)
			}
		})
	}
}

func TestValidateCapabilityID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid simple", "test-capability", true},
		{"valid with numbers", "test-capability-123", true},
		{"valid all lowercase", "authservice", true},
		{"empty", "", false},
		{"invalid uppercase", "Test-Cap", false},
		{"invalid underscore", "test_capability", false},
		{"invalid @", "test@capability", false},
		{"invalid space", "test capability", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCapabilityID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateCapabilityID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

// ==========================
// Brief Status Constants Tests
// ==========================

func TestBriefStatusConstants(t *testing.T) {
	if BriefStatusDraft != "draft" {
		t.Errorf("BriefStatusDraft = %q, want %q", BriefStatusDraft, "draft")
	}
	if BriefStatusActive != "active" {
		t.Errorf("BriefStatusActive = %q, want %q", BriefStatusActive, "active")
	}
	if BriefStatusArchived != "archived" {
		t.Errorf("BriefStatusArchived = %q, want %q", BriefStatusArchived, "archived")
	}
}

// ==========================
// Brief Struct Tests
// ==========================

func TestBriefStruct(t *testing.T) {
	_ = time.Now().UTC()
	brief := Brief{
		ID:          "auth-brief",
		BacklogID:   "auth",
		Status:      BriefStatusActive,
		Why:         string(make([]byte, 500)),
		WhatChanges: []string{"New authentication flow"},
		Impact: BriefImpact{
			TechnicalStack:  []string{"Go", "PostgreSQL"},
			AffectedSystems: []string{"auth-service"},
			Dependencies:    []string{"external-idp"},
		},
		NewCapabilities: []Capability{
			{ID: "sso-login", Name: "SSO Login", Description: "Login via SSO"},
		},
		ModifiedCapabilities: []Capability{},
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
		CreatedBy:            "testuser",
		UpdatedBy:            "testuser",
	}

	if brief.ID != "auth-brief" {
		t.Errorf("Brief.ID = %q, want %q", brief.ID, "auth-brief")
	}
	if brief.Status != BriefStatusActive {
		t.Errorf("Brief.Status = %q, want %q", brief.Status, BriefStatusActive)
	}
	if len(brief.NewCapabilities) != 1 {
		t.Errorf("Brief.NewCapabilities count = %d, want 1", len(brief.NewCapabilities))
	}
}

func TestCapabilityStruct(t *testing.T) {
	tests := []struct {
		name       string
		capability Capability
	}{
		{"new capability", Capability{ID: "new-cap", Name: "New", Description: "desc"}},
		{"modified capability", Capability{ID: "existing-cap", Name: "Existing", Description: "updated"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.capability.ID == "" {
				t.Error("Capability ID should not be empty")
			}
			if tt.capability.Name == "" {
				t.Error("Capability Name should not be empty")
			}
		})
	}
}

func TestBriefImpactStruct(t *testing.T) {
	impact := BriefImpact{
		TechnicalStack:  []string{"Go", "React"},
		AffectedSystems: []string{"frontend", "backend"},
		Dependencies:    []string{"external-api"},
	}

	if len(impact.TechnicalStack) != 2 {
		t.Errorf("BriefImpact.TechnicalStack count = %d, want 2", len(impact.TechnicalStack))
	}
}

func TestBriefCreateInput(t *testing.T) {
	input := BriefCreateInput{
		BacklogID:   "auth",
		Name:        "Auth Brief",
		Why:         string(make([]byte, 500)),
		WhatChanges: []string{"New SSO login"},
		Impact:      nil,
		Capabilities: []CapabilityInput{
			{Name: "SSO Login", Description: "Login via SSO", Modified: false},
		},
	}

	if input.BacklogID != "auth" {
		t.Errorf("BriefCreateInput.BacklogID = %q, want %q", input.BacklogID, "auth")
	}
	if len(input.Capabilities) != 1 {
		t.Errorf("BriefCreateInput.Capabilities count = %d, want 1", len(input.Capabilities))
	}
}

func TestBriefUpdateInput(t *testing.T) {
	status := BriefStatusArchived
	input := BriefUpdateInput{
		ID:     "auth-brief",
		Status: &status,
	}

	if *input.Status != BriefStatusArchived {
		t.Errorf("BriefUpdateInput.Status = %q, want %q", *input.Status, BriefStatusArchived)
	}
}

// ==========================
// Spec Validation Tests
// ==========================

func TestValidateSpecStructure(t *testing.T) {
	tests := []struct {
		name     string
		spec     *Spec
		expected bool
	}{
		{"valid spec with requirements", &Spec{
			Requirements: []Requirement{
				{
					ID: "req-0001",
					IAEScenarios: []IAEScenario{
						{ID: "s1", Intent: "User wants to login", Action: "Click login", Expect: "Login page shown"},
					},
				},
			},
		}, true},
		{"nil spec", nil, false},
		{"empty requirements", &Spec{Requirements: []Requirement{}}, false},
		{"requirement without scenarios", &Spec{
			Requirements: []Requirement{
				{ID: "req-0001", IAEScenarios: []IAEScenario{}},
			},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSpecStructure(tt.spec)
			if result != tt.expected {
				t.Errorf("ValidateSpecStructure() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateRequirementStructure(t *testing.T) {
	tests := []struct {
		name     string
		req      *Requirement
		expected bool
	}{
		{"valid requirement", &Requirement{
			ID:           "req-0001",
			IAEScenarios: []IAEScenario{{ID: "s1", Intent: "i", Action: "a", Expect: "e"}},
		}, true},
		{"nil requirement", nil, false},
		{"empty scenarios", &Requirement{ID: "req-0001", IAEScenarios: []IAEScenario{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRequirementStructure(tt.req)
			if result != tt.expected {
				t.Errorf("ValidateRequirementStructure() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateIAEScenario(t *testing.T) {
	tests := []struct {
		name     string
		scenario *IAEScenario
		expected bool
	}{
		{"valid scenario", &IAEScenario{ID: "s1", Intent: "i", Action: "a", Expect: "e"}, true},
		{"nil scenario", nil, false},
		{"empty intent", &IAEScenario{ID: "s1", Intent: "", Action: "a", Expect: "e"}, false},
		{"empty action", &IAEScenario{ID: "s1", Intent: "i", Action: "", Expect: "e"}, false},
		{"empty expect", &IAEScenario{ID: "s1", Intent: "i", Action: "a", Expect: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIAEScenario(tt.scenario)
			if result != tt.expected {
				t.Errorf("ValidateIAEScenario() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ==========================
// Spec Status Constants Tests
// ==========================

func TestSpecStatusConstants(t *testing.T) {
	if SpecStatusDraft != "draft" {
		t.Errorf("SpecStatusDraft = %q, want %q", SpecStatusDraft, "draft")
	}
	if SpecStatusActive != "active" {
		t.Errorf("SpecStatusActive = %q, want %q", SpecStatusActive, "active")
	}
	if SpecStatusArchived != "archived" {
		t.Errorf("SpecStatusArchived = %q, want %q", SpecStatusArchived, "archived")
	}
}

// ==========================
// Spec Struct Tests
// ==========================

func TestSpecStruct(t *testing.T) {
	now := time.Now().UTC()
	spec := Spec{
		ID:           "test-cap-spec",
		CapabilityID: "test-cap",
		BacklogID:    "auth",
		Status:       SpecStatusActive,
		Summary:      "Test specification",
		Requirements: []Requirement{
			{
				ID:                 "req-0001",
				Summary:            "User authentication",
				Details:            "Users should be able to authenticate",
				AcceptanceCriteria: []string{"Can login with email", "Can login with SSO"},
				IAEScenarios: []IAEScenario{
					{ID: "s1", Intent: "Login", Action: "Click login", Expect: "Success"},
				},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if spec.ID != "test-cap-spec" {
		t.Errorf("Spec.ID = %q, want %q", spec.ID, "test-cap-spec")
	}
	if len(spec.Requirements) != 1 {
		t.Errorf("Spec.Requirements count = %d, want 1", len(spec.Requirements))
	}
}

func TestIAEScenarioStruct(t *testing.T) {
	scenario := IAEScenario{
		ID:     "s1234",
		Intent: "User wants to search",
		Action: "Enters query",
		Expect: "Results displayed",
	}

	if scenario.ID == "" {
		t.Error("IAEScenario.ID should not be empty")
	}
}

func TestRequirementStruct(t *testing.T) {
	req := Requirement{
		ID:                 "req-0001",
		Summary:            "Test requirement",
		Details:            "Detailed description",
		AcceptanceCriteria: []string{"Criteria 1", "Criteria 2"},
		IAEScenarios: []IAEScenario{
			{ID: "s1", Intent: "i", Action: "a", Expect: "e"},
		},
	}

	if len(req.AcceptanceCriteria) != 2 {
		t.Errorf("Requirement.AcceptanceCriteria count = %d, want 2", len(req.AcceptanceCriteria))
	}
}

func TestSpecCreateInput(t *testing.T) {
	input := SpecCreateInput{
		BacklogID:    "auth",
		CapabilityID: "test-cap",
		Summary:      "Test Spec",
		Requirements: []RequirementInput{
			{
				Summary:            "Req 1",
				Details:            "Details",
				AcceptanceCriteria: []string{"AC1"},
				IAEScenarios: []IAEScenarioInput{
					{Intent: "i", Action: "a", Expect: "e"},
				},
			},
		},
	}

	if input.BacklogID != "auth" {
		t.Errorf("SpecCreateInput.BacklogID = %q, want %q", input.BacklogID, "auth")
	}
}

func TestSpecUpdateInput(t *testing.T) {
	status := SpecStatusArchived
	input := SpecUpdateInput{
		ID:     "test-cap-spec",
		Status: &status,
	}

	if *input.Status != SpecStatusArchived {
		t.Errorf("SpecUpdateInput.Status = %q, want %q", *input.Status, SpecStatusArchived)
	}
}

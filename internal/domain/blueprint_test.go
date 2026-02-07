package domain

import (
	"testing"
	"time"
)

// ==========================
// Blueprint Validation Tests
// ==========================

func TestValidateBlueprintStructure(t *testing.T) {
	tests := []struct {
		name     string
		bp       *Blueprint
		expected bool
	}{
		{"valid blueprint", &Blueprint{
			ArchitectureDecisions: []ArchitectureDecision{
				{Title: "Use PostgreSQL", Decision: "Use PostgreSQL", Rationale: string(make([]byte, 50))},
			},
		}, true},
		{"nil blueprint", nil, false},
		{"empty decisions", &Blueprint{ArchitectureDecisions: []ArchitectureDecision{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBlueprintStructure(tt.bp)
			if result != tt.expected {
				t.Errorf("ValidateBlueprintStructure() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateArchitectureDecision(t *testing.T) {
	tests := []struct {
		name     string
		decision *ArchitectureDecision
		expected bool
	}{
		{"valid decision", &ArchitectureDecision{
			Title:     "Use PostgreSQL",
			Decision:  "Use PostgreSQL for data storage",
			Rationale: string(make([]byte, 50)),
		}, true},
		{"nil decision", nil, false},
		{"empty title", &ArchitectureDecision{Decision: "d", Rationale: string(make([]byte, 50))}, false},
		{"empty decision", &ArchitectureDecision{Title: "t", Rationale: string(make([]byte, 50))}, false},
		{"rationale too short", &ArchitectureDecision{Title: "t", Decision: "d", Rationale: string(make([]byte, 49))}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateArchitectureDecision(tt.decision)
			if result != tt.expected {
				t.Errorf("ValidateArchitectureDecision() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ==========================
// Blueprint Status Constants Tests
// ==========================

func TestBlueprintStatusConstants(t *testing.T) {
	if BlueprintStatusDraft != "draft" {
		t.Errorf("BlueprintStatusDraft = %q, want %q", BlueprintStatusDraft, "draft")
	}
	if BlueprintStatusActive != "active" {
		t.Errorf("BlueprintStatusActive = %q, want %q", BlueprintStatusActive, "active")
	}
	if BlueprintStatusArchived != "archived" {
		t.Errorf("BlueprintStatusArchived = %q, want %q", BlueprintStatusArchived, "archived")
	}
}

// ==========================
// Blueprint Struct Tests
// ==========================

func TestBlueprintStruct(t *testing.T) {
	now := time.Now().UTC()
	bp := Blueprint{
		ID:               "auth-blueprint",
		BriefID:          "auth-brief",
		BacklogID:        "auth",
		Status:           BlueprintStatusActive,
		Version:          "1.0.0",
		ProblemStatement: "Need secure authentication",
		Constraints:      []string{"Must use OAuth 2.0"},
		UserTypes:        []string{"End users", "Admins"},
		Goals: BlueprintGoals{
			InScope:  []string{"SSO login", "Token management"},
			OutScope: []string{"Legacy auth"},
		},
		ArchitectureDecisions: []ArchitectureDecision{
			{
				ID:                     "ADR-001",
				Title:                  "Use PostgreSQL",
				Decision:               "Use PostgreSQL for data storage",
				Rationale:              string(make([]byte, 50)),
				AlternativesConsidered: []string{"MySQL", "MongoDB"},
			},
		},
		DataModels: []DataModel{
			{
				Name:        "User",
				Description: "User account",
				Fields: []DataModelField{
					{Name: "ID", Type: "UUID", Required: true, Description: "Primary key"},
				},
			},
		},
		ImplementationStrategy: "Phase 1: Database schema\nPhase 2: API endpoints",
		Risks: []Risk{
			{ID: "R1", Description: "Performance", Mitigation: "Add caching"},
		},
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if bp.ID != "auth-blueprint" {
		t.Errorf("Blueprint.ID = %q, want %q", bp.ID, "auth-blueprint")
	}
	if len(bp.ArchitectureDecisions) != 1 {
		t.Errorf("Blueprint.ArchitectureDecisions count = %d, want 1", len(bp.ArchitectureDecisions))
	}
}

func TestArchitectureDecisionStruct(t *testing.T) {
	decision := ArchitectureDecision{
		ID:                     "ADR-001",
		Title:                  "Use PostgreSQL",
		Decision:               "Use PostgreSQL",
		Rationale:              string(make([]byte, 50)),
		AlternativesConsidered: []string{"MySQL", "MongoDB"},
	}

	if decision.ID != "ADR-001" {
		t.Errorf("ArchitectureDecision.ID = %q, want %q", decision.ID, "ADR-001")
	}
	if len(decision.AlternativesConsidered) != 2 {
		t.Errorf("ArchitectureDecision.AlternativesConsidered count = %d, want 2", len(decision.AlternativesConsidered))
	}
}

func TestDataModelStruct(t *testing.T) {
	model := DataModel{
		Name:        "User",
		Description: "User account model",
		Fields: []DataModelField{
			{Name: "ID", Type: "UUID", Required: true, Description: "Primary key"},
			{Name: "Email", Type: "string", Required: true, Description: "User email"},
		},
	}

	if len(model.Fields) != 2 {
		t.Errorf("DataModel.Fields count = %d, want 2", len(model.Fields))
	}
}

func TestDataModelFieldStruct(t *testing.T) {
	field := DataModelField{
		Name:        "ID",
		Type:        "UUID",
		Required:    true,
		Description: "Primary key",
	}

	if field.Required != true {
		t.Errorf("DataModelField.Required = false, want true")
	}
}

func TestRiskStruct(t *testing.T) {
	risk := Risk{
		ID:          "R1",
		Description: "Performance degradation",
		Mitigation:  "Add Redis caching layer",
	}

	if risk.ID != "R1" {
		t.Errorf("Risk.ID = %q, want %q", risk.ID, "R1")
	}
}

func TestBlueprintGoalsStruct(t *testing.T) {
	goals := BlueprintGoals{
		InScope:  []string{"Feature A", "Feature B"},
		OutScope: []string{"Legacy system"},
	}

	if len(goals.InScope) != 2 {
		t.Errorf("BlueprintGoals.InScope count = %d, want 2", len(goals.InScope))
	}
	if len(goals.OutScope) != 1 {
		t.Errorf("BlueprintGoals.OutScope count = %d, want 1", len(goals.OutScope))
	}
}

func TestBlueprintCreateInput(t *testing.T) {
	input := BlueprintCreateInput{
		BacklogID:        "auth",
		BriefID:          "auth-brief",
		ProblemStatement: "Need authentication",
		Constraints:      []string{"Must use OAuth"},
		UserTypes:        []string{"End users"},
		Goals: &BlueprintGoals{
			InScope:  []string{"SSO"},
			OutScope: []string{},
		},
		ArchitectureDecisions: []ArchitectureDecisionInput{
			{
				Title:                  "Use OAuth 2.0",
				Decision:               "Use OAuth 2.0 for auth",
				Rationale:              string(make([]byte, 50)),
				AlternativesConsidered: []string{"SAML", "Basic auth"},
			},
		},
		DataModels:             []DataModel{},
		ImplementationStrategy: "Phase 1 implementation",
		Risks:                  []RiskInput{},
	}

	if input.BacklogID != "auth" {
		t.Errorf("BlueprintCreateInput.BacklogID = %q, want %q", input.BacklogID, "auth")
	}
}

func TestArchitectureDecisionInput(t *testing.T) {
	input := ArchitectureDecisionInput{
		Title:                  "Use PostgreSQL",
		Decision:               "Decision text",
		Rationale:              string(make([]byte, 50)),
		AlternativesConsidered: []string{"MySQL"},
	}

	if len(input.AlternativesConsidered) != 1 {
		t.Errorf("ArchitectureDecisionInput.AlternativesConsidered count = %d, want 1", len(input.AlternativesConsidered))
	}
}

func TestRiskInput(t *testing.T) {
	input := RiskInput{
		Description: "Risk description",
		Mitigation:  "Mitigation strategy",
	}

	if input.Description != "Risk description" {
		t.Errorf("RiskInput.Description = %q, want %q", input.Description, "Risk description")
	}
}

func TestBlueprintUpdateInput(t *testing.T) {
	status := BlueprintStatusArchived
	input := BlueprintUpdateInput{
		ID:     "auth-blueprint",
		Status: &status,
	}

	if *input.Status != BlueprintStatusArchived {
		t.Errorf("BlueprintUpdateInput.Status = %q, want %q", *input.Status, BlueprintStatusArchived)
	}
}

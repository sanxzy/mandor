package domain

import "time"

const (
	SpecStatusDraft    = "draft"
	SpecStatusActive   = "active"
	SpecStatusArchived = "archived"
)

// IAEScenario represents an Intent-Action-Expect scenario
type IAEScenario struct {
	ID     string `json:"id"`     // scenario-XXXX (4-char base-62 random)
	Intent string `json:"intent"` // What the user intends to do
	Action string `json:"action"` // What action they take
	Expect string `json:"expect"` // What they expect to happen
}

// Requirement represents a requirement with IAE scenarios
type Requirement struct {
	ID                  string         `json:"id"`                    // req-XXXX (4-char base-62 random)
	Summary             string         `json:"summary"`               // Short requirement description
	Details             string         `json:"details"`               // Detailed requirement
	AcceptanceCriteria  []string       `json:"acceptance_criteria"`   // List of acceptance criteria
	IAEScenarios        []IAEScenario  `json:"iae_scenarios"`         // Intent-Action-Expect scenarios (min 1)
}

// Spec represents a specification for a capability
type Spec struct {
	ID            string         `json:"id"`             // spec-id ({capability-id}-spec)
	CapabilityID  string         `json:"capability_id"`  // Reference to Brief capability
	ProjectID     string         `json:"project_id"`     // Reference to project
	Status        string         `json:"status"`         // draft | active | archived
	Summary       string         `json:"summary"`        // Brief description of specification
	Requirements  []Requirement  `json:"requirements"`   // Minimum 1 requirement required
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedBy     string         `json:"created_by"`
	UpdatedBy     string         `json:"updated_by"`
}

// SpecCreateInput represents input for creating a Spec
type SpecCreateInput struct {
	ProjectID     string
	CapabilityID  string             // From Brief capability
	Summary       string             // Brief description
	Requirements  []RequirementInput  // Minimum 1
}

// RequirementInput represents input for a requirement
type RequirementInput struct {
	Summary            string           // Short requirement description
	Details            string           // Detailed requirement
	AcceptanceCriteria []string         // List of acceptance criteria
	IAEScenarios       []IAEScenarioInput // Minimum 1 per requirement
}

// IAEScenarioInput represents input for an IAE scenario
type IAEScenarioInput struct {
	Intent string
	Action string
	Expect string
}

// SpecUpdateInput represents input for updating a Spec
type SpecUpdateInput struct {
	ID           string
	Status       *string
	Summary      *string
	Requirements *[]Requirement
}

// ValidateSpecStructure validates that a Spec has minimum required structure
func ValidateSpecStructure(spec *Spec) bool {
	if spec == nil || len(spec.Requirements) == 0 {
		return false
	}
	for _, req := range spec.Requirements {
		if len(req.IAEScenarios) == 0 {
			return false
		}
	}
	return true
}

// ValidateRequirementStructure validates a requirement has at least 1 IAE scenario
func ValidateRequirementStructure(req *Requirement) bool {
	return req != nil && len(req.IAEScenarios) > 0
}

// ValidateIAEScenario validates an IAE scenario has all three fields
func ValidateIAEScenario(scenario *IAEScenario) bool {
	return scenario != nil && 
		len(scenario.Intent) > 0 && 
		len(scenario.Action) > 0 && 
		len(scenario.Expect) > 0
}

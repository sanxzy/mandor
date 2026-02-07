package domain

import "time"

const (
	BlueprintStatusDraft    = "draft"
	BlueprintStatusActive   = "active"
	BlueprintStatusArchived = "archived"
)

// ArchitectureDecision represents a key architectural decision
type ArchitectureDecision struct {
	ID                     string   `json:"id"`                      // Decision identifier
	Title                  string   `json:"title"`                   // Decision title
	Decision               string   `json:"decision"`                // What decision was made
	Rationale              string   `json:"rationale"`               // Why this decision was made (min 50 chars)
	AlternativesConsidered []string `json:"alternatives_considered"` // Options that were rejected
}

// DataModel represents a data model in the blueprint
type DataModel struct {
	Name        string           `json:"name"`        // Model name
	Description string           `json:"description"` // What this model represents
	Fields      []DataModelField `json:"fields"`      // Model fields
}

// DataModelField represents a field in a data model
type DataModelField struct {
	Name        string `json:"name"`        // Field name
	Type        string `json:"type"`        // Field type (UUID, string, date, etc.)
	Required    bool   `json:"required"`    // Whether field is required
	Description string `json:"description"` // Field description
}

// Risk represents a risk and its mitigation
type Risk struct {
	ID          string `json:"id"`          // Risk identifier
	Description string `json:"description"` // What could go wrong
	Mitigation  string `json:"mitigation"`  // How we will prevent or handle
}

// Blueprint represents the technical architecture document
type Blueprint struct {
	ID                     string                 `json:"id"`                      // blueprint-id ({backlog-id}-blueprint)
	BriefID                string                 `json:"brief_id"`                // Reference to Brief
	BacklogID              string                 `json:"backlog_id"`              // Reference to backlog
	Status                 string                 `json:"status"`                  // draft | active | archived
	Version                string                 `json:"version"`                 // Version number
	ProblemStatement       string                 `json:"problem_statement"`       // Context: problem being solved
	Constraints            []string               `json:"constraints"`             // Context: constraints
	UserTypes              []string               `json:"user_types"`              // Context: types of users
	Goals                  BlueprintGoals         `json:"goals"`                   // In scope and out of scope
	ArchitectureDecisions  []ArchitectureDecision `json:"architecture_decisions"`  // Min 1 decision required
	DataModels             []DataModel            `json:"data_models"`             // Data models
	ImplementationStrategy string                 `json:"implementation_strategy"` // High-level implementation approach
	Risks                  []Risk                 `json:"risks"`                   // Risks and mitigations
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
	CreatedBy              string                 `json:"created_by"`
	UpdatedBy              string                 `json:"updated_by"`
}

// BlueprintGoals represents in-scope and out-of-scope items
type BlueprintGoals struct {
	InScope  []string `json:"in_scope"`  // Items in scope
	OutScope []string `json:"out_scope"` // Items out of scope
}

// BlueprintCreateInput represents input for creating a Blueprint
type BlueprintCreateInput struct {
	BacklogID              string
	BriefID                string
	ProblemStatement       string
	Constraints            []string
	UserTypes              []string
	Goals                  *BlueprintGoals
	ArchitectureDecisions  []ArchitectureDecisionInput
	DataModels             []DataModel
	ImplementationStrategy string
	Risks                  []RiskInput
}

// ArchitectureDecisionInput represents input for an architecture decision
type ArchitectureDecisionInput struct {
	Title                  string
	Decision               string
	Rationale              string // Min 50 chars
	AlternativesConsidered []string
}

// RiskInput represents input for a risk
type RiskInput struct {
	Description string
	Mitigation  string
}

// BlueprintUpdateInput represents input for updating a Blueprint
type BlueprintUpdateInput struct {
	ID                     string
	Status                 *string
	ProblemStatement       *string
	Constraints            *[]string
	UserTypes              *[]string
	Goals                  *BlueprintGoals
	ArchitectureDecisions  *[]ArchitectureDecision
	DataModels             *[]DataModel
	ImplementationStrategy *string
	Risks                  *[]Risk
}

// ValidateBlueprintStructure validates minimum Blueprint structure
func ValidateBlueprintStructure(bp *Blueprint) bool {
	if bp == nil {
		return false
	}
	// Min 1 architecture decision required
	return len(bp.ArchitectureDecisions) >= 1
}

// ValidateArchitectureDecision validates a decision has rationale (min 50 chars)
func ValidateArchitectureDecision(decision *ArchitectureDecision) bool {
	return decision != nil &&
		len(decision.Title) > 0 &&
		len(decision.Decision) > 0 &&
		len(decision.Rationale) >= 50
}

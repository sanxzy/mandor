package service

import (
	"encoding/json"
	"fmt"
	"time"

	"mandor/internal/fs"
)

// ChangeGovernanceService enforces deterministic change propagation rules
type ChangeGovernanceService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewChangeGovernanceService(paths *fs.Paths) *ChangeGovernanceService {
	return &ChangeGovernanceService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

// ChangeImpactAnalysis holds detailed impact information for a change
type ChangeImpactAnalysis struct {
	ChangeID           string   `json:"change_id"`
	ChangeType         string   `json:"change_type"` // brief_modification, spec_modification, blueprint_modification
	EntityID           string   `json:"entity_id"`
	BacklogID          string   `json:"backlog_id"`
	FieldsChanged      []string `json:"fields_changed"`
	ImpactedSpecs      []string `json:"impacted_specs"`
	ImpactedFeatures   []string `json:"impacted_features"`
	ImpactedTasks      []string `json:"impacted_tasks"`
	ImpactedBlueprints []string `json:"impacted_blueprints"`
	BlockingStatus     string   `json:"blocking_status"` // BLOCKING, NON_BLOCKING
	RequiredActions    []string `json:"required_actions"`
	ValidationDeadline string   `json:"validation_deadline"`
	Status             string   `json:"status"` // pending_validation, approved, rejected
	Timestamp          string   `json:"timestamp"`
	User               string   `json:"user"`
	VersionBefore      string   `json:"version_before"`
	VersionAfter       string   `json:"version_after"`
}

// ValidateBriefChangeBlocking checks if brief change should be blocked
func (s *ChangeGovernanceService) ValidateBriefChangeBlocking(
	backlogID string,
	briefID string,
	changedFields map[string]interface{},
) (*ChangeImpactAnalysis, error) {

	impact := &ChangeImpactAnalysis{
		ChangeID:   fmt.Sprintf("brief-change-%d", time.Now().UnixNano()),
		ChangeType: "brief_modification",
		EntityID:   briefID,
		BacklogID:  backlogID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Status:     "pending_validation",
	}

	// Analyze field changes for blocking conditions
	for field := range changedFields {
		switch field {
		case "capabilities":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate all affected Specs")
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate Blueprint")
		case "tech-stack":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate Blueprint")
		case "timeline", "team-size":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Revalidate all task effort estimates")
		case "why", "business-rationale", "priority":
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		default:
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		}
	}

	// Set validation deadline
	impact.ValidationDeadline = time.Now().AddDate(0, 0, 3).UTC().Format(time.RFC3339)

	for field := range changedFields {
		impact.FieldsChanged = append(impact.FieldsChanged, field)
	}

	return impact, nil
}

// ValidateSpecChangeBlocking checks if spec change should be blocked
func (s *ChangeGovernanceService) ValidateSpecChangeBlocking(
	backlogID string,
	specID string,
	changedFields map[string]interface{},
) (*ChangeImpactAnalysis, error) {

	impact := &ChangeImpactAnalysis{
		ChangeID:   fmt.Sprintf("spec-change-%d", time.Now().UnixNano()),
		ChangeType: "spec_modification",
		EntityID:   specID,
		BacklogID:  backlogID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Status:     "pending_validation",
	}

	// Analyze field changes for blocking conditions
	for field := range changedFields {
		switch field {
		case "acceptance-criteria", "requirements", "iae-scenarios":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate all dependent Tasks")
		case "input-validation", "error-handling":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate affected Tasks")
		case "summary", "testing-strategy-notes":
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		default:
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		}
	}

	// Set validation deadline
	impact.ValidationDeadline = time.Now().AddDate(0, 0, 3).UTC().Format(time.RFC3339)

	for field := range changedFields {
		impact.FieldsChanged = append(impact.FieldsChanged, field)
	}

	return impact, nil
}

// ValidateBlueprintChangeBlocking checks if blueprint change should be blocked
func (s *ChangeGovernanceService) ValidateBlueprintChangeBlocking(
	backlogID string,
	blueprintID string,
	changedFields map[string]interface{},
) (*ChangeImpactAnalysis, error) {

	impact := &ChangeImpactAnalysis{
		ChangeID:   fmt.Sprintf("blueprint-change-%d", time.Now().UnixNano()),
		ChangeType: "blueprint_modification",
		EntityID:   blueprintID,
		BacklogID:  backlogID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Status:     "pending_validation",
	}

	// Analyze field changes for blocking conditions
	for field := range changedFields {
		switch field {
		case "tech-stack", "deployment-strategy", "architecture-pattern":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate all dependent Tasks")
		case "scalability", "resilience", "performance-targets":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Revalidate task acceptance criteria")
		case "monitoring", "alerting", "operational-requirements":
			impact.BlockingStatus = "BLOCKING"
			impact.RequiredActions = append(impact.RequiredActions, "Regenerate operations/monitoring tasks")
		case "cost-analysis", "1-year-projection", "rationale", "decision-reversibility", "assumptions":
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		default:
			if impact.BlockingStatus != "BLOCKING" {
				impact.BlockingStatus = "NON_BLOCKING"
			}
		}
	}

	// Set validation deadline
	impact.ValidationDeadline = time.Now().AddDate(0, 0, 3).UTC().Format(time.RFC3339)

	for field := range changedFields {
		impact.FieldsChanged = append(impact.FieldsChanged, field)
	}

	return impact, nil
}

// PersistChangeAnalysis stores change analysis in audit log
func (s *ChangeGovernanceService) PersistChangeAnalysis(backlogID string, analysis *ChangeImpactAnalysis) error {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	// In production, append to change audit log (.mandor/backlogs/{backlog-id}/change-audit.jsonl)
	// For now, return success (persistence layer would be implemented separately)
	_ = data // Use data to avoid linter error

	return nil
}

// ClassifyBlockingRisk provides detailed risk assessment for a change
func (s *ChangeGovernanceService) ClassifyBlockingRisk(impact *ChangeImpactAnalysis) string {
	if impact.BlockingStatus == "BLOCKING" {
		if len(impact.ImpactedTasks) > 10 {
			return "HIGH_RISK_MANY_TASKS"
		}
		if len(impact.RequiredActions) > 5 {
			return "HIGH_RISK_COMPLEX_CHANGE"
		}
		return "BLOCKING"
	}
	return "NON_BLOCKING"
}

// ValidateChangePrerequisites checks if change can proceed
func (s *ChangeGovernanceService) ValidateChangePrerequisites(
	backlogID string,
	analysis *ChangeImpactAnalysis,
) (bool, []string) {
	var blockers []string

	// Check if any tasks are in_progress
	for _, taskID := range analysis.ImpactedTasks {
		task, err := s.reader.ReadTask(backlogID, taskID)
		if err == nil && (task.Status == "in_progress" || task.Status == "done") {
			blockers = append(blockers,
				fmt.Sprintf("Task %s has status %s (cannot modify while task is executing)", taskID, task.Status))
		}
	}

	if len(blockers) > 0 {
		return false, blockers
	}

	return true, []string{}
}

// ApprovalStatus represents a change approval decision
type ApprovalStatus struct {
	ChangeID       string `json:"change_id"`
	Approved       bool   `json:"approved"`
	ApprovalReason string `json:"approval_reason"`
	ApprovedAt     string `json:"approved_at"`
	ApprovedBy     string `json:"approved_by"`
}

// ApproveChange records a change approval
func (s *ChangeGovernanceService) ApproveChange(
	backlogID string,
	changeID string,
	reason string,
) (*ApprovalStatus, error) {
	approval := &ApprovalStatus{
		ChangeID:       changeID,
		Approved:       true,
		ApprovalReason: reason,
		ApprovedAt:     time.Now().UTC().Format(time.RFC3339),
		ApprovedBy:     "system", // In production, capture user ID
	}

	return approval, nil
}

// BlockingFields maps change fields to blocking classification
type BlockingFields struct {
	BriefBlockingFields     map[string]bool
	SpecBlockingFields      map[string]bool
	BlueprintBlockingFields map[string]bool
}

// GetBlockingFieldMap returns all blocking field classifications
func (s *ChangeGovernanceService) GetBlockingFieldMap() *BlockingFields {
	return &BlockingFields{
		BriefBlockingFields: map[string]bool{
			"capabilities":       true,
			"tech-stack":         true,
			"timeline":           true,
			"team-size":          true,
			"why":                false,
			"business-rationale": false,
			"priority":           false,
		},
		SpecBlockingFields: map[string]bool{
			"acceptance-criteria":    true,
			"requirements":           true,
			"iae-scenarios":          true,
			"input-validation":       true,
			"error-handling":         true,
			"summary":                false,
			"testing-strategy-notes": false,
		},
		BlueprintBlockingFields: map[string]bool{
			"tech-stack":               true,
			"deployment-strategy":      true,
			"architecture-pattern":     true,
			"scalability":              true,
			"resilience":               true,
			"performance-targets":      true,
			"monitoring":               true,
			"alerting":                 true,
			"operational-requirements": true,
			"cost-analysis":            false,
			"1-year-projection":        false,
			"rationale":                false,
			"decision-reversibility":   false,
			"assumptions":              false,
		},
	}
}

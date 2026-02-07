package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"mandor/internal/fs"
	"mandor/internal/service"
)

// ChangeCommand handles change governance operations
type ChangeCommand struct {
	service *service.ChangeGovernanceService
	paths   *fs.Paths
}

// NewChangeCommand creates a new change command handler
func NewChangeCommand(paths *fs.Paths) *ChangeCommand {
	return &ChangeCommand{
		service: service.NewChangeGovernanceService(paths),
		paths:   paths,
	}
}

// AnalyzeRequest contains parameters for change analysis
type AnalyzeRequest struct {
	EntityType   string // brief, spec, blueprint
	EntityID     string
	BacklogID    string
	FieldsChange []string // comma-separated field names
}

// AnalyzeResponse contains the analysis result
type AnalyzeResponse struct {
	Analysis  *service.ChangeImpactAnalysis `json:"analysis"`
	Message   string                        `json:"message"`
	Timestamp string                        `json:"timestamp"`
}

// Analyze performs change impact analysis (positive case)
func (c *ChangeCommand) Analyze(req *AnalyzeRequest) (*AnalyzeResponse, error) {
	// Validate input
	if err := c.validateAnalyzeRequest(req); err != nil {
		return nil, err
	}

	// Convert field list to map
	fieldsMap := make(map[string]interface{})
	for _, field := range req.FieldsChange {
		field = strings.TrimSpace(field)
		if field != "" {
			fieldsMap[field] = true
		}
	}

	if len(fieldsMap) == 0 {
		return nil, fmt.Errorf("no fields specified for analysis")
	}

	// Perform analysis based on entity type
	var analysis *service.ChangeImpactAnalysis
	var err error

	switch strings.ToLower(req.EntityType) {
	case "brief":
		analysis, err = c.service.ValidateBriefChangeBlocking(req.BacklogID, req.EntityID, fieldsMap)
	case "spec":
		analysis, err = c.service.ValidateSpecChangeBlocking(req.BacklogID, req.EntityID, fieldsMap)
	case "blueprint":
		analysis, err = c.service.ValidateBlueprintChangeBlocking(req.BacklogID, req.EntityID, fieldsMap)
	default:
		return nil, fmt.Errorf("invalid entity type: %s (must be brief, spec, or blueprint)", req.EntityType)
	}

	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Persist analysis
	if err := c.service.PersistChangeAnalysis(req.BacklogID, analysis); err != nil {
		return nil, fmt.Errorf("failed to persist analysis: %w", err)
	}

	return &AnalyzeResponse{
		Analysis:  analysis,
		Message:   fmt.Sprintf("Change analysis complete. Blocking status: %s", analysis.BlockingStatus),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// ApproveRequest contains parameters for change approval
type ApproveRequest struct {
	ChangeID  string
	Reason    string
	BacklogID string
}

// ApproveResponse contains the approval result
type ApproveResponse struct {
	Approval  *service.ApprovalStatus `json:"approval"`
	Message   string                  `json:"message"`
	Timestamp string                  `json:"timestamp"`
}

// Approve approves a pending change
func (c *ChangeCommand) Approve(req *ApproveRequest) (*ApproveResponse, error) {
	// Validate input
	if strings.TrimSpace(req.ChangeID) == "" {
		return nil, fmt.Errorf("change ID is required")
	}

	if strings.TrimSpace(req.Reason) == "" {
		return nil, fmt.Errorf("approval reason is required (minimum 10 characters)")
	}

	if len(req.Reason) < 10 {
		return nil, fmt.Errorf("approval reason too short (minimum 10 characters, got %d)", len(req.Reason))
	}

	// Perform approval
	approval, err := c.service.ApproveChange(req.BacklogID, req.ChangeID, req.Reason)
	if err != nil {
		return nil, fmt.Errorf("approval failed: %w", err)
	}

	return &ApproveResponse{
		Approval:  approval,
		Message:   fmt.Sprintf("Change %s approved successfully", req.ChangeID),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// RejectRequest contains parameters for change rejection
type RejectRequest struct {
	ChangeID  string
	Reason    string
	BacklogID string
}

// RejectResponse contains the rejection result
type RejectResponse struct {
	Status    string `json:"status"`
	ChangeID  string `json:"change_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// Reject rejects a pending change
func (c *ChangeCommand) Reject(req *RejectRequest) (*RejectResponse, error) {
	// Validate input
	if strings.TrimSpace(req.ChangeID) == "" {
		return nil, fmt.Errorf("change ID is required")
	}

	if strings.TrimSpace(req.Reason) == "" {
		return nil, fmt.Errorf("rejection reason is required (minimum 10 characters)")
	}

	if len(req.Reason) < 10 {
		return nil, fmt.Errorf("rejection reason too short (minimum 10 characters, got %d)", len(req.Reason))
	}

	return &RejectResponse{
		Status:    "rejected",
		ChangeID:  req.ChangeID,
		Message:   fmt.Sprintf("Change %s rejected", req.ChangeID),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// ListChangesRequest contains parameters for listing changes
type ListChangesRequest struct {
	BacklogID  string
	Status     string // pending_validation, approved, rejected, all
	EntityType string // brief, spec, blueprint, all
}

// ListChangesResponse contains the changes list
type ListChangesResponse struct {
	Changes   []map[string]interface{} `json:"changes"`
	Total     int                      `json:"total"`
	Filtered  int                      `json:"filtered"`
	Timestamp string                   `json:"timestamp"`
}

// ListChanges lists all changes for a backlog
func (c *ChangeCommand) ListChanges(req *ListChangesRequest) (*ListChangesResponse, error) {
	// Validate input
	if strings.TrimSpace(req.BacklogID) == "" {
		return nil, fmt.Errorf("backlog ID is required")
	}

	// Validate status filter
	if req.Status != "" && req.Status != "all" {
		validStatuses := []string{"pending_validation", "approved", "rejected"}
		found := false
		for _, s := range validStatuses {
			if req.Status == s {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("invalid status: %s (must be pending_validation, approved, rejected, or all)", req.Status)
		}
	}

	// Validate entity type filter
	if req.EntityType != "" && req.EntityType != "all" {
		validTypes := []string{"brief", "spec", "blueprint"}
		found := false
		for _, t := range validTypes {
			if req.EntityType == t {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("invalid entity type: %s (must be brief, spec, blueprint, or all)", req.EntityType)
		}
	}

	// In production, this would query the audit log
	// For now, return empty list (audit persistence not yet implemented)
	return &ListChangesResponse{
		Changes:   []map[string]interface{}{},
		Total:     0,
		Filtered:  0,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// ExportChangeAnalysis exports analysis to JSON
func (c *ChangeCommand) ExportChangeAnalysis(analysis *service.ChangeImpactAnalysis) (string, error) {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal analysis: %w", err)
	}
	return string(data), nil
}

// Helper methods

func (c *ChangeCommand) validateAnalyzeRequest(req *AnalyzeRequest) error {
	if strings.TrimSpace(req.EntityType) == "" {
		return fmt.Errorf("entity type is required (brief, spec, or blueprint)")
	}

	if strings.TrimSpace(req.EntityID) == "" {
		return fmt.Errorf("entity ID is required")
	}

	if strings.TrimSpace(req.BacklogID) == "" {
		return fmt.Errorf("backlog ID is required")
	}

	if len(req.FieldsChange) == 0 {
		return fmt.Errorf("at least one field must be specified for analysis")
	}

	return nil
}

// PrintAnalysis formats analysis for display
func (c *ChangeCommand) PrintAnalysis(analysis *service.ChangeImpactAnalysis) string {
	var output strings.Builder

	output.WriteString("\n================================\n")
	output.WriteString("CHANGE IMPACT ANALYSIS\n")
	output.WriteString("================================\n\n")

	output.WriteString(fmt.Sprintf("Change ID:        %s\n", analysis.ChangeID))
	output.WriteString(fmt.Sprintf("Change Type:      %s\n", analysis.ChangeType))
	output.WriteString(fmt.Sprintf("Entity ID:        %s\n", analysis.EntityID))
	output.WriteString(fmt.Sprintf("Backlog ID:       %s\n", analysis.BacklogID))
	output.WriteString(fmt.Sprintf("Status:           %s\n", analysis.Status))
	output.WriteString(fmt.Sprintf("Blocking Status:  %s\n", analysis.BlockingStatus))
	output.WriteString(fmt.Sprintf("Timestamp:        %s\n", analysis.Timestamp))

	output.WriteString("\nFields Changed:\n")
	for _, field := range analysis.FieldsChanged {
		output.WriteString(fmt.Sprintf("  - %s\n", field))
	}

	output.WriteString(fmt.Sprintf("\nImpacted Entities:\n"))
	output.WriteString(fmt.Sprintf("  Specs:       %d\n", len(analysis.ImpactedSpecs)))
	output.WriteString(fmt.Sprintf("  Features:    %d\n", len(analysis.ImpactedFeatures)))
	output.WriteString(fmt.Sprintf("  Tasks:       %d\n", len(analysis.ImpactedTasks)))
	output.WriteString(fmt.Sprintf("  Blueprints:  %d\n", len(analysis.ImpactedBlueprints)))

	if len(analysis.RequiredActions) > 0 {
		output.WriteString("\nRequired Actions:\n")
		for i, action := range analysis.RequiredActions {
			output.WriteString(fmt.Sprintf("  %d. %s\n", i+1, action))
		}
	}

	output.WriteString(fmt.Sprintf("\nValidation Deadline: %s\n", analysis.ValidationDeadline))
	output.WriteString("================================\n\n")

	return output.String()
}

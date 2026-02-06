package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/util"
)

type BlueprintService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewBlueprintService() (*BlueprintService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &BlueprintService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewBlueprintServiceWithPaths(paths *fs.Paths) *BlueprintService {
	return &BlueprintService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *BlueprintService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *BlueprintService) ValidateCreateInput(input *domain.BlueprintCreateInput) error {
	if !s.reader.ProjectExists(input.ProjectID) {
		return domain.NewValidationError("Project not found: " + input.ProjectID)
	}

	if strings.TrimSpace(input.BriefID) == "" {
		return domain.NewValidationError("Brief ID is required.")
	}

	// Load and validate Brief exists
	brief, err := s.loadBrief(input.ProjectID)
	if err != nil {
		return domain.NewValidationError("Brief not found. Create a Brief first.")
	}

	// Verify all Brief capabilities have corresponding Specs
	if err := s.verifyAllSpecsExist(input.ProjectID, brief); err != nil {
		return err
	}

	// Verify each Spec is valid
	if err := s.verifyAllSpecsValid(input.ProjectID, brief); err != nil {
		return err
	}

	if strings.TrimSpace(input.ProblemStatement) == "" {
		return domain.NewValidationError("Problem statement is required.")
	}

	// Validate minimum 1 architecture decision
	if len(input.ArchitectureDecisions) == 0 {
		return domain.NewValidationError("Error: at least one architecture decision is required\nSolution: use --decisions 'title:rationale' with rationale min 50 chars")
	}

	// Validate each decision
	for i, decision := range input.ArchitectureDecisions {
		if strings.TrimSpace(decision.Title) == "" {
			return domain.NewValidationError(fmt.Sprintf("Architecture Decision %d: title cannot be empty.", i+1))
		}
		if strings.TrimSpace(decision.Decision) == "" {
			return domain.NewValidationError(fmt.Sprintf("Architecture Decision %d: decision statement cannot be empty.", i+1))
		}
		if len(decision.Rationale) < 50 {
			return domain.NewValidationError(fmt.Sprintf("Architecture Decision %d: rationale must be at least 50 characters (got %d).", i+1, len(decision.Rationale)))
		}
	}

	return nil
}

func (s *BlueprintService) CreateBlueprint(input *domain.BlueprintCreateInput) (*domain.Blueprint, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	// Generate blueprint ID as {project-id}-blueprint
	blueprintID := util.ToSlug(input.ProjectID) + "-blueprint"

	// Process architecture decisions
	decisions := []domain.ArchitectureDecision{}
	for _, decInput := range input.ArchitectureDecisions {
		decID := util.NextSequential("decision", []string{})
		decisions = append(decisions, domain.ArchitectureDecision{
			ID:                   decID,
			Title:                decInput.Title,
			Decision:             decInput.Decision,
			Rationale:            decInput.Rationale,
			AlternativesConsidered: decInput.AlternativesConsidered,
		})
	}

	// Process risks
	risks := []domain.Risk{}
	for _, riskInput := range input.Risks {
		riskID := util.NextSequential("risk", []string{})
		risks = append(risks, domain.Risk{
			ID:          riskID,
			Description: riskInput.Description,
			Mitigation:  riskInput.Mitigation,
		})
	}

	// Create Blueprint
	blueprint := &domain.Blueprint{
		ID:                    blueprintID,
		BriefID:               input.BriefID,
		ProjectID:             input.ProjectID,
		Status:                domain.BlueprintStatusDraft,
		Version:               "1.0",
		ProblemStatement:      input.ProblemStatement,
		Constraints:           input.Constraints,
		UserTypes:             input.UserTypes,
		Goals:                 *input.Goals,
		ArchitectureDecisions: decisions,
		DataModels:            input.DataModels,
		ImplementationStrategy: input.ImplementationStrategy,
		Risks:                 risks,
		CreatedAt:             now,
		UpdatedAt:             now,
		CreatedBy:             creator,
		UpdatedBy:             creator,
	}

	// Validate blueprint structure
	if !domain.ValidateBlueprintStructure(blueprint) {
		return nil, domain.NewValidationError("Error: Blueprint must have at least 1 architecture decision\nSolution: add --decisions with proper format")
	}

	// Save Blueprint
	if err := s.saveBlueprint(input.ProjectID, blueprint); err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (s *BlueprintService) ReadBlueprint(projectID string) (*domain.Blueprint, error) {
	blueprintPath := s.paths.BlueprintPath(projectID)

	if !s.blueprintFileExists(blueprintPath) {
		return nil, domain.NewValidationError("Blueprint not found for project: " + projectID)
	}

	blueprint, err := s.loadBlueprint(blueprintPath)
	if err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (s *BlueprintService) saveBlueprint(projectID string, blueprint *domain.Blueprint) error {
	blueprintPath := s.paths.BlueprintPath(projectID)

	// Build markdown representation
	content := s.blueprintToMarkdown(blueprint)

	// Write to markdown file
	if err := s.writer.WriteFile(blueprintPath, content); err != nil {
		return err
	}

	return nil
}

func (s *BlueprintService) loadBlueprint(blueprintPath string) (*domain.Blueprint, error) {
	content, err := s.reader.ReadFile(blueprintPath)
	if err != nil {
		return nil, err
	}

	blueprint, err := s.markdownToBlueprint(content)
	if err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (s *BlueprintService) blueprintFileExists(path string) bool {
	content, err := s.reader.ReadFile(path)
	return err == nil && len(strings.TrimSpace(content)) > 0
}

func (s *BlueprintService) blueprintToMarkdown(bp *domain.Blueprint) string {
	// Serialize Blueprint to JSON for storage
	bpJSON, _ := json.MarshalIndent(bp, "", "  ")

	sb := strings.Builder{}

	sb.WriteString("<!--\n")
	sb.WriteString("BLUEPRINT_DATA_START\n")
	sb.WriteString(string(bpJSON))
	sb.WriteString("\nBLUEPRINT_DATA_END\n")
	sb.WriteString("-->\n\n")

	sb.WriteString(fmt.Sprintf("# Blueprint: %s\n\n", bp.ID))
	sb.WriteString(fmt.Sprintf("**Version:** %s\n", bp.Version))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", bp.Status))

	sb.WriteString("## Problem Statement\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", bp.ProblemStatement))

	if len(bp.Constraints) > 0 {
		sb.WriteString("## Constraints\n\n")
		for _, constraint := range bp.Constraints {
			sb.WriteString(fmt.Sprintf("- %s\n", constraint))
		}
		sb.WriteString("\n")
	}

	if len(bp.UserTypes) > 0 {
		sb.WriteString("## User Types\n\n")
		for _, ut := range bp.UserTypes {
			sb.WriteString(fmt.Sprintf("- %s\n", ut))
		}
		sb.WriteString("\n")
	}

	if len(bp.Goals.InScope) > 0 || len(bp.Goals.OutScope) > 0 {
		sb.WriteString("## Goals\n\n")
		if len(bp.Goals.InScope) > 0 {
			sb.WriteString("### In Scope\n")
			for _, g := range bp.Goals.InScope {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}
		if len(bp.Goals.OutScope) > 0 {
			sb.WriteString("### Out of Scope\n")
			for _, g := range bp.Goals.OutScope {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("## Architecture Decisions\n\n")
	for i, dec := range bp.ArchitectureDecisions {
		sb.WriteString(fmt.Sprintf("### Decision %d: %s (%s)\n\n", i+1, dec.Title, dec.ID))
		sb.WriteString(fmt.Sprintf("**Decision:** %s\n\n", dec.Decision))
		sb.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", dec.Rationale))
		if len(dec.AlternativesConsidered) > 0 {
			sb.WriteString("**Alternatives Considered:**\n")
			for _, alt := range dec.AlternativesConsidered {
				sb.WriteString(fmt.Sprintf("- %s\n", alt))
			}
			sb.WriteString("\n")
		}
	}

	if len(bp.DataModels) > 0 {
		sb.WriteString("## Data Models\n\n")
		for i, dm := range bp.DataModels {
			sb.WriteString(fmt.Sprintf("### Model %d: %s\n\n", i+1, dm.Name))
			sb.WriteString(fmt.Sprintf("%s\n\n", dm.Description))
			sb.WriteString("**Fields:**\n")
			for _, field := range dm.Fields {
				required := "optional"
				if field.Required {
					required = "required"
				}
				sb.WriteString(fmt.Sprintf("- **%s** (%s, %s): %s\n", field.Name, field.Type, required, field.Description))
			}
			sb.WriteString("\n")
		}
	}

	if bp.ImplementationStrategy != "" {
		sb.WriteString("## Implementation Strategy\n\n")
		sb.WriteString(fmt.Sprintf("%s\n\n", bp.ImplementationStrategy))
	}

	if len(bp.Risks) > 0 {
		sb.WriteString("## Risks & Mitigations\n\n")
		for i, risk := range bp.Risks {
			sb.WriteString(fmt.Sprintf("### Risk %d (%s)\n\n", i+1, risk.ID))
			sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", risk.Description))
			sb.WriteString(fmt.Sprintf("**Mitigation:** %s\n\n", risk.Mitigation))
		}
	}

	sb.WriteString(fmt.Sprintf("**Created:** %s by %s\n", bp.CreatedAt.Format(time.RFC3339), bp.CreatedBy))
	sb.WriteString(fmt.Sprintf("**Updated:** %s by %s\n", bp.UpdatedAt.Format(time.RFC3339), bp.UpdatedBy))

	return sb.String()
}

func (s *BlueprintService) markdownToBlueprint(content string) (*domain.Blueprint, error) {
	// Extract JSON from markdown comment
	startMarker := "BLUEPRINT_DATA_START\n"
	endMarker := "\nBLUEPRINT_DATA_END"

	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return nil, domain.NewSystemError("Blueprint data not found in file", nil)
	}

	jsonData := content[startIdx+len(startMarker) : endIdx]

	var bp domain.Blueprint
	if err := json.Unmarshal([]byte(jsonData), &bp); err != nil {
		return nil, domain.NewSystemError("Failed to parse Blueprint data", err)
	}

	return &bp, nil
}

func (s *BlueprintService) loadBrief(projectID string) (*domain.Brief, error) {
	briefService := NewBriefServiceWithPaths(s.paths)
	return briefService.ReadBrief(projectID)
}

func (s *BlueprintService) verifyAllSpecsExist(projectID string, brief *domain.Brief) error {
	specService := NewSpecServiceWithPaths(s.paths)

	allCapabilities := append(brief.NewCapabilities, brief.ModifiedCapabilities...)
	for _, cap := range allCapabilities {
		specID := cap.ID + "-spec"
		_, err := specService.ReadSpec(projectID, specID)
		if err != nil {
			return domain.NewValidationError(fmt.Sprintf("Error: Capability '%s' does not have a valid Spec\nSolution: create a Spec for this capability first using 'mandor spec create'", cap.ID))
		}
	}
	return nil
}

func (s *BlueprintService) verifyAllSpecsValid(projectID string, brief *domain.Brief) error {
	specService := NewSpecServiceWithPaths(s.paths)

	allCapabilities := append(brief.NewCapabilities, brief.ModifiedCapabilities...)
	for _, cap := range allCapabilities {
		specID := cap.ID + "-spec"
		spec, err := specService.ReadSpec(projectID, specID)
		if err != nil {
			return err
		}

		// Verify spec has minimum required structure
		if !domain.ValidateSpecStructure(spec) {
			return domain.NewValidationError(fmt.Sprintf("Error: Spec '%s' is invalid\nSolution: ensure it has at least 1 requirement with 1 IAE scenario", specID))
		}
	}
	return nil
}

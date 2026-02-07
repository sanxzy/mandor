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

type SpecService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewSpecService() (*SpecService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &SpecService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewSpecServiceWithPaths(paths *fs.Paths) *SpecService {
	return &SpecService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *SpecService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *SpecService) ValidateCreateInput(input *domain.SpecCreateInput) error {
	if !s.reader.BacklogExists(input.BacklogID) {
		return domain.NewValidationError("Backlog not found: " + input.BacklogID)
	}

	if strings.TrimSpace(input.CapabilityID) == "" {
		return domain.NewValidationError("Capability ID is required.")
	}

	// Validate capability ID format
	if !domain.ValidateCapabilityID(input.CapabilityID) {
		return domain.NewValidationError(fmt.Sprintf("Error: Invalid capability ID '%s'\nSolution: use alphanumeric characters and hyphens only", input.CapabilityID))
	}

	// TODO: Verify capability exists in Brief
	brief, err := s.loadBriefForBacklog(input.BacklogID)
	if err != nil {
		return domain.NewValidationError("Brief not found for backlog. Create a Brief first with this capability.")
	}

	if !s.capabilityExistsInBrief(brief, input.CapabilityID) {
		return domain.NewValidationError(fmt.Sprintf("Capability '%s' not found in Brief. Create it in the Brief first.", input.CapabilityID))
	}

	if strings.TrimSpace(input.Summary) == "" {
		return domain.NewValidationError("Spec summary is required.")
	}

	// Validate minimum 1 requirement
	if len(input.Requirements) == 0 {
		return domain.NewValidationError("Error: at least one requirement is required\nSolution: provide --requirements with proper format")
	}

	// Validate each requirement
	for i, req := range input.Requirements {
		if strings.TrimSpace(req.Summary) == "" {
			return domain.NewValidationError(fmt.Sprintf("Requirement %d: summary cannot be empty.", i+1))
		}

		if len(req.IAEScenarios) == 0 {
			return domain.NewValidationError(fmt.Sprintf("Requirement %d: must have at least 1 IAE scenario.", i+1))
		}

		for j, scenario := range req.IAEScenarios {
			if strings.TrimSpace(scenario.Intent) == "" {
				return domain.NewValidationError(fmt.Sprintf("Requirement %d, Scenario %d: Intent cannot be empty.", i+1, j+1))
			}
			if strings.TrimSpace(scenario.Action) == "" {
				return domain.NewValidationError(fmt.Sprintf("Requirement %d, Scenario %d: Action cannot be empty.", i+1, j+1))
			}
			if strings.TrimSpace(scenario.Expect) == "" {
				return domain.NewValidationError(fmt.Sprintf("Requirement %d, Scenario %d: Expect cannot be empty.", i+1, j+1))
			}
		}
	}

	return nil
}

func (s *SpecService) CreateSpec(input *domain.SpecCreateInput) (*domain.Spec, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	// Generate spec ID as {capability-id}-spec
	specID := input.CapabilityID + "-spec"

	// Process requirements and generate IDs
	requirements := []domain.Requirement{}
	for _, reqInput := range input.Requirements {
		reqID := util.NextSequential("req", []string{})

		scenarios := []domain.IAEScenario{}
		for _, scnInput := range reqInput.IAEScenarios {
			scnID := util.NextSequential("scenario", []string{})
			scenarios = append(scenarios, domain.IAEScenario{
				ID:     scnID,
				Intent: scnInput.Intent,
				Action: scnInput.Action,
				Expect: scnInput.Expect,
			})
		}

		requirements = append(requirements, domain.Requirement{
			ID:                 reqID,
			Summary:            reqInput.Summary,
			Details:            reqInput.Details,
			AcceptanceCriteria: reqInput.AcceptanceCriteria,
			IAEScenarios:       scenarios,
		})
	}

	// Create Spec
	spec := &domain.Spec{
		ID:           specID,
		CapabilityID: input.CapabilityID,
		BacklogID:    input.BacklogID,
		Status:       domain.SpecStatusDraft,
		Summary:      input.Summary,
		Requirements: requirements,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    creator,
		UpdatedBy:    creator,
	}

	// Validate spec structure
	if !domain.ValidateSpecStructure(spec) {
		return nil, domain.NewValidationError("Error: Spec validation failed: minimum 1 requirement with 1 IAE scenario required\nSolution: ensure each requirement has Intent, Action, Expect fields")
	}

	// Save Spec
	if err := s.saveSpec(input.BacklogID, spec); err != nil {
		return nil, err
	}

	return spec, nil
}

func (s *SpecService) ReadSpec(backlogID, specID string) (*domain.Spec, error) {
	specPath := s.paths.SpecPath(backlogID, specID)

	if !s.specFileExists(specPath) {
		return nil, domain.NewValidationError("Spec not found: " + specID)
	}

	spec, err := s.loadSpec(specPath)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func (s *SpecService) saveSpec(backlogID string, spec *domain.Spec) error {
	specPath := s.paths.SpecPath(backlogID, spec.ID)

	// Build markdown representation
	content := s.specToMarkdown(spec)

	// Write to markdown file
	if err := s.writer.WriteFile(specPath, content); err != nil {
		return err
	}

	return nil
}

func (s *SpecService) loadSpec(specPath string) (*domain.Spec, error) {
	content, err := s.reader.ReadFile(specPath)
	if err != nil {
		return nil, err
	}

	spec, err := s.markdownToSpec(content)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func (s *SpecService) specFileExists(path string) bool {
	content, err := s.reader.ReadFile(path)
	return err == nil && len(strings.TrimSpace(content)) > 0
}

func (s *SpecService) specToMarkdown(spec *domain.Spec) string {
	// Serialize Spec to JSON for storage
	specJSON, _ := json.MarshalIndent(spec, "", "  ")

	sb := strings.Builder{}

	sb.WriteString("<!--\n")
	sb.WriteString("SPEC_DATA_START\n")
	sb.WriteString(string(specJSON))
	sb.WriteString("\nSPEC_DATA_END\n")
	sb.WriteString("-->\n\n")

	sb.WriteString(fmt.Sprintf("# Spec: %s\n\n", spec.ID))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", spec.Status))
	sb.WriteString(fmt.Sprintf("**Capability:** %s\n\n", spec.CapabilityID))

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", spec.Summary))

	sb.WriteString("## Requirements\n\n")
	for i, req := range spec.Requirements {
		sb.WriteString(fmt.Sprintf("### Requirement %d: %s (%s)\n\n", i+1, req.Summary, req.ID))
		if req.Details != "" {
			sb.WriteString(fmt.Sprintf("**Details:** %s\n\n", req.Details))
		}

		if len(req.AcceptanceCriteria) > 0 {
			sb.WriteString("**Acceptance Criteria:**\n")
			for _, ac := range req.AcceptanceCriteria {
				sb.WriteString(fmt.Sprintf("- %s\n", ac))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("**IAE Scenarios:**\n\n")
		for j, iae := range req.IAEScenarios {
			sb.WriteString(fmt.Sprintf("#### Scenario %d (%s)\n\n", j+1, iae.ID))
			sb.WriteString(fmt.Sprintf("- **Intent:** %s\n", iae.Intent))
			sb.WriteString(fmt.Sprintf("- **Action:** %s\n", iae.Action))
			sb.WriteString(fmt.Sprintf("- **Expect:** %s\n\n", iae.Expect))
		}
	}

	sb.WriteString(fmt.Sprintf("**Created:** %s by %s\n", spec.CreatedAt.Format(time.RFC3339), spec.CreatedBy))
	sb.WriteString(fmt.Sprintf("**Updated:** %s by %s\n", spec.UpdatedAt.Format(time.RFC3339), spec.UpdatedBy))

	return sb.String()
}

func (s *SpecService) markdownToSpec(content string) (*domain.Spec, error) {
	// Extract JSON from markdown comment
	startMarker := "SPEC_DATA_START\n"
	endMarker := "\nSPEC_DATA_END"

	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return nil, domain.NewSystemError("Spec data not found in file", nil)
	}

	jsonData := content[startIdx+len(startMarker) : endIdx]

	var spec domain.Spec
	if err := json.Unmarshal([]byte(jsonData), &spec); err != nil {
		return nil, domain.NewSystemError("Failed to parse Spec data", err)
	}

	return &spec, nil
}

func (s *SpecService) loadBriefForBacklog(backlogID string) (*domain.Brief, error) {
	briefService := NewBriefServiceWithPaths(s.paths)
	return briefService.ReadBrief(backlogID)
}

func (s *SpecService) capabilityExistsInBrief(brief *domain.Brief, capabilityID string) bool {
	for _, cap := range brief.NewCapabilities {
		if cap.ID == capabilityID {
			return true
		}
	}
	for _, cap := range brief.ModifiedCapabilities {
		if cap.ID == capabilityID {
			return true
		}
	}
	return false
}

func (s *SpecService) UpdateSpec(backlogID string, spec *domain.Spec) error {
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	spec.UpdatedAt = now
	spec.UpdatedBy = updater

	return s.saveSpec(backlogID, spec)
}

func (s *SpecService) DeleteSpec(backlogID, specID string) error {
	specPath := s.paths.SpecPath(backlogID, specID)
	return s.writer.DeleteFile(specPath)
}

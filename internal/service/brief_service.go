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

type BriefService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewBriefService() (*BriefService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &BriefService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewBriefServiceWithPaths(paths *fs.Paths) *BriefService {
	return &BriefService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *BriefService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *BriefService) ValidateCreateInput(input *domain.BriefCreateInput) error {
	if !s.reader.BacklogExists(input.BacklogID) {
		return domain.NewValidationError("Backlog not found: " + input.BacklogID)
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.NewValidationError("Brief name is required.")
	}

	// Validate brief ID format (alphanumeric and hyphens only)
	briefID := util.ToSlug(input.Name)
	if !domain.ValidateCapabilityID(briefID) {
		return domain.NewValidationError(fmt.Sprintf("Error: Invalid Brief name '%s'\nSolution: use alphanumeric characters and hyphens only", input.Name))
	}

	// Check if brief already exists for this backlog
	briefPath := s.paths.BriefPath(input.BacklogID)
	if s.briefFileExists(briefPath) {
		return domain.NewValidationError("A Brief already exists for this backlog. Update or delete the existing Brief first.")
	}

	// Validate Why section
	if !domain.ValidateBriefWhy(input.Why) {
		return domain.NewValidationError(fmt.Sprintf("Error: Why section must be 100-5000 characters (got %d)\nSolution: expand or condense your problem statement", len(input.Why)))
	}

	// Validate at least one capability
	if len(input.Capabilities) == 0 {
		return domain.NewValidationError("Error: at least one capability is required\nSolution: provide at least one capability")
	}

	// Validate each capability
	for _, cap := range input.Capabilities {
		if strings.TrimSpace(cap.Name) == "" {
			return domain.NewValidationError("Capability name cannot be empty.")
		}
		capID := util.ToSlug(cap.Name)
		if !domain.ValidateCapabilityID(capID) {
			return domain.NewValidationError(fmt.Sprintf("Error: Invalid capability name '%s'\nSolution: use alphanumeric characters and hyphens only", cap.Name))
		}
		if strings.TrimSpace(cap.Description) == "" {
			return domain.NewValidationError("Capability description cannot be empty.")
		}
	}

	return nil
}

func (s *BriefService) CreateBrief(input *domain.BriefCreateInput) (*domain.Brief, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	briefID := util.ToSlug(input.Name)

	// Process capabilities
	newCapabilities := []domain.Capability{}
	modifiedCapabilities := []domain.Capability{}

	for _, cap := range input.Capabilities {
		capID := util.ToSlug(cap.Name)
		capability := domain.Capability{
			ID:          capID,
			Name:        cap.Name,
			Description: cap.Description,
		}

		if cap.Modified {
			modifiedCapabilities = append(modifiedCapabilities, capability)
		} else {
			newCapabilities = append(newCapabilities, capability)
		}
	}

	// Create Brief
	brief := &domain.Brief{
		ID:                   briefID,
		BacklogID:            input.BacklogID,
		Status:               domain.BriefStatusDraft,
		Why:                  input.Why,
		Impact:               *input.Impact,
		NewCapabilities:      newCapabilities,
		ModifiedCapabilities: modifiedCapabilities,
		CreatedAt:            now,
		UpdatedAt:            now,
		CreatedBy:            creator,
		UpdatedBy:            creator,
	}

	// Build What Changes list
	for _, cap := range append(newCapabilities, modifiedCapabilities...) {
		brief.WhatChanges = append(brief.WhatChanges, cap.Description)
	}

	// Save Brief
	if err := s.saveBrief(input.BacklogID, brief); err != nil {
		return nil, err
	}

	return brief, nil
}

func (s *BriefService) ReadBrief(backlogID string) (*domain.Brief, error) {
	briefPath := s.paths.BriefPath(backlogID)

	if !s.briefFileExists(briefPath) {
		return nil, domain.NewValidationError("Brief not found for backlog: " + backlogID)
	}

	brief, err := s.loadBrief(briefPath)
	if err != nil {
		return nil, err
	}

	return brief, nil
}

func (s *BriefService) saveBrief(backlogID string, brief *domain.Brief) error {
	briefPath := s.paths.BriefPath(backlogID)

	// Create specs directory if it doesn't exist
	specsDir := s.paths.SpecsDirPath(backlogID)
	if err := s.writer.CreateDirectory(specsDir); err != nil {
		return err
	}

	// Build markdown with embedded JSON
	content := s.briefToMarkdown(brief)

	// Write to markdown file
	if err := s.writer.WriteFile(briefPath, content); err != nil {
		return err
	}

	return nil
}

func (s *BriefService) loadBrief(briefPath string) (*domain.Brief, error) {
	content, err := s.reader.ReadFile(briefPath)
	if err != nil {
		return nil, err
	}

	brief, err := s.markdownToBrief(content)
	if err != nil {
		return nil, err
	}

	return brief, nil
}

func (s *BriefService) briefFileExists(path string) bool {
	content, err := s.reader.ReadFile(path)
	return err == nil && len(strings.TrimSpace(content)) > 0
}

func (s *BriefService) briefToMarkdown(brief *domain.Brief) string {
	// Serialize Brief to JSON for storage
	briefJSON, _ := json.MarshalIndent(brief, "", "  ")

	// Build markdown representation with embedded JSON
	sb := strings.Builder{}

	sb.WriteString("<!--\n")
	sb.WriteString("BRIEF_DATA_START\n")
	sb.WriteString(string(briefJSON))
	sb.WriteString("\nBRIEF_DATA_END\n")
	sb.WriteString("-->\n\n")

	sb.WriteString(fmt.Sprintf("# Brief: %s\n\n", brief.ID))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", brief.Status))

	sb.WriteString("## Why\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", brief.Why))

	sb.WriteString("## What Changes\n\n")
	for _, change := range brief.WhatChanges {
		sb.WriteString(fmt.Sprintf("- %s\n", change))
	}
	sb.WriteString("\n")

	if len(brief.NewCapabilities) > 0 {
		sb.WriteString("## New Capabilities\n\n")
		for _, cap := range brief.NewCapabilities {
			sb.WriteString(fmt.Sprintf("### %s (%s)\n\n", cap.Name, cap.ID))
			sb.WriteString(fmt.Sprintf("%s\n\n", cap.Description))
		}
	}

	if len(brief.ModifiedCapabilities) > 0 {
		sb.WriteString("## Modified Capabilities\n\n")
		for _, cap := range brief.ModifiedCapabilities {
			sb.WriteString(fmt.Sprintf("### %s (%s)\n\n", cap.Name, cap.ID))
			sb.WriteString(fmt.Sprintf("%s\n\n", cap.Description))
		}
	}

	if len(brief.Impact.TechnicalStack) > 0 {
		sb.WriteString("## Technical Stack\n\n")
		for _, tech := range brief.Impact.TechnicalStack {
			sb.WriteString(fmt.Sprintf("- %s\n", tech))
		}
		sb.WriteString("\n")
	}

	if len(brief.Impact.AffectedSystems) > 0 {
		sb.WriteString("## Affected Systems\n\n")
		for _, sys := range brief.Impact.AffectedSystems {
			sb.WriteString(fmt.Sprintf("- %s\n", sys))
		}
		sb.WriteString("\n")
	}

	if len(brief.Impact.Dependencies) > 0 {
		sb.WriteString("## Dependencies\n\n")
		for _, dep := range brief.Impact.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s\n", dep))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("**Created:** %s by %s\n", brief.CreatedAt.Format(time.RFC3339), brief.CreatedBy))
	sb.WriteString(fmt.Sprintf("**Updated:** %s by %s\n", brief.UpdatedAt.Format(time.RFC3339), brief.UpdatedBy))

	return sb.String()
}

func (s *BriefService) markdownToBrief(content string) (*domain.Brief, error) {
	// Extract JSON from markdown comment
	startMarker := "BRIEF_DATA_START\n"
	endMarker := "\nBRIEF_DATA_END"

	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return nil, domain.NewSystemError("Brief data not found in file", nil)
	}

	jsonData := content[startIdx+len(startMarker) : endIdx]

	var brief domain.Brief
	if err := json.Unmarshal([]byte(jsonData), &brief); err != nil {
		return nil, domain.NewSystemError("Failed to parse Brief data", err)
	}

	return &brief, nil
}

func (s *BriefService) UpdateBrief(backlogID string, brief *domain.Brief) error {
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	brief.UpdatedAt = now
	brief.UpdatedBy = updater

	return s.saveBrief(backlogID, brief)
}

func (s *BriefService) DeleteBrief(backlogID string) error {
	briefPath := s.paths.BriefPath(backlogID)
	return s.writer.DeleteFile(briefPath)
}

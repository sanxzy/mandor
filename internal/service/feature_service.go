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

type FeatureService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewFeatureService() (*FeatureService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &FeatureService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewFeatureServiceWithPaths(paths *fs.Paths) *FeatureService {
	return &FeatureService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *FeatureService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *FeatureService) ValidateCreateInput(input *domain.FeatureCreateInput) error {
	if !s.reader.ProjectExists(input.ProjectID) {
		return domain.NewValidationError("Project not found: " + input.ProjectID)
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.NewValidationError("Feature name is required.")
	}

	if strings.TrimSpace(input.Goal) == "" {
		return domain.NewValidationError("Feature goal is required (--goal).")
	}

	if input.Scope != "" && !domain.ValidateScope(input.Scope) {
		return domain.NewValidationError("Invalid scope. Valid options: frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift")
	}

	if !domain.ValidatePriority(input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Priority == "" {
		input.Priority = "P3"
	}

	if err := s.validateDependencies(input.ProjectID, "", input.DependsOn); err != nil {
		return err
	}

	return nil
}

func (s *FeatureService) validateDependencies(projectID, selfID string, dependsOn []string) error {
	for _, depID := range dependsOn {
		if depID == selfID {
			return domain.NewValidationError("Self-dependency detected. Entity cannot depend on itself.")
		}

		dep, err := s.reader.ReadFeature(projectID, depID)
		if err != nil {
			if _, ok := err.(*domain.MandorError); ok {
				return domain.NewValidationError("Dependency not found: " + depID)
			}
			return err
		}

		if dep.Status == domain.FeatureStatusCancelled {
			return domain.NewValidationError("Dependency is cancelled: " + depID)
		}
	}

	if err := s.validateNoCycle(projectID, selfID, dependsOn); err != nil {
		return err
	}

	return nil
}

func (s *FeatureService) validateNoCycle(projectID, selfID string, dependsOn []string) error {
	visited := make(map[string]bool)
	var dfs func(featureID string) bool

	dfs = func(featureID string) bool {
		if featureID == selfID {
			return true
		}
		if visited[featureID] {
			return false
		}
		visited[featureID] = true

		f, err := s.reader.ReadFeature(projectID, featureID)
		if err != nil {
			return false
		}

		for _, dep := range f.DependsOn {
			if dfs(dep) {
				return true
			}
		}
		return false
	}

	for _, depID := range dependsOn {
		visited = make(map[string]bool)
		if dfs(depID) {
			return domain.NewValidationError("Circular dependency detected.")
		}
	}

	return nil
}

func (s *FeatureService) CreateFeature(input *domain.FeatureCreateInput) (*domain.Feature, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	nanoid, err := util.GenerateID()
	if err != nil {
		return nil, domain.NewSystemError("Failed to generate feature ID", err)
	}

	featureID := input.ProjectID + "-feature-" + nanoid

	feature := &domain.Feature{
		ID:        featureID,
		ProjectID: input.ProjectID,
		Name:      input.Name,
		Goal:      input.Goal,
		Scope:     input.Scope,
		Priority:  input.Priority,
		Status:    domain.FeatureStatusDraft,
		DependsOn: input.DependsOn,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: creator,
		UpdatedBy: creator,
	}

	if len(input.DependsOn) > 0 {
		allDone, err := s.checkDependenciesDone(input.ProjectID, input.DependsOn)
		if err != nil {
			return nil, err
		}
		if !allDone {
			feature.Status = domain.FeatureStatusBlocked
		}
	}

	if err := s.writer.WriteFeature(input.ProjectID, feature); err != nil {
		return nil, err
	}

	event := &domain.FeatureEvent{
		Layer: "feature",
		Type:  "created",
		ID:    featureID,
		By:    creator,
		Ts:    now,
	}
	if err := s.writer.AppendFeatureEvent(input.ProjectID, event); err != nil {
		return nil, err
	}

	return feature, nil
}

func (s *FeatureService) checkDependenciesDone(projectID string, dependsOn []string) (bool, error) {
	for _, depID := range dependsOn {
		dep, err := s.reader.ReadFeature(projectID, depID)
		if err != nil {
			return false, domain.NewValidationError("Dependency not found: " + depID)
		}
		if dep.Status != domain.FeatureStatusDone && dep.Status != domain.FeatureStatusCancelled {
			return false, nil
		}
	}
	return true, nil
}

func (s *FeatureService) ListFeatures(input *domain.FeatureListInput) (*domain.FeatureListOutput, error) {
	if !s.reader.ProjectExists(input.ProjectID) {
		return nil, domain.NewValidationError("Project not found: " + input.ProjectID)
	}

	var features []domain.FeatureListItem
	deletedCount := 0

	err := s.reader.ReadNDJSON(s.paths.ProjectFeaturesPath(input.ProjectID), func(raw []byte) error {
		var f domain.Feature
		if err := json.Unmarshal(raw, &f); err != nil {
			return err
		}
		if !input.IncludeDeleted && f.Status == domain.FeatureStatusCancelled {
			return nil
		}

		item := domain.FeatureListItem{
			ID:        f.ID,
			Name:      f.Name,
			Goal:      f.Goal,
			Scope:     f.Scope,
			Priority:  f.Priority,
			Status:    f.Status,
			DependsOn: len(f.DependsOn),
			CreatedAt: f.CreatedAt.Format(time.RFC3339),
			UpdatedAt: f.UpdatedAt.Format(time.RFC3339),
		}
		features = append(features, item)

		if f.Status == domain.FeatureStatusCancelled {
			deletedCount++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &domain.FeatureListOutput{
		Features: features,
		Total:    len(features),
		Deleted:  deletedCount,
	}, nil
}

func (s *FeatureService) GetFeatureDetail(input *domain.FeatureDetailInput) (*domain.FeatureDetailOutput, error) {
	feature, err := s.reader.ReadFeature(input.ProjectID, input.FeatureID)
	if err != nil {
		return nil, err
	}

	if !input.IncludeDeleted && feature.Status == domain.FeatureStatusCancelled {
		return nil, domain.NewValidationError("Feature not found: " + input.FeatureID)
	}

	events, _ := s.reader.CountEventLines(input.ProjectID)

	return &domain.FeatureDetailOutput{
		ID:        feature.ID,
		ProjectID: feature.ProjectID,
		Name:      feature.Name,
		Goal:      feature.Goal,
		Scope:     feature.Scope,
		Priority:  feature.Priority,
		Status:    feature.Status,
		DependsOn: feature.DependsOn,
		Reason:    feature.Reason,
		Events:    events,
		CreatedAt: feature.CreatedAt.Format(time.RFC3339),
		UpdatedAt: feature.UpdatedAt.Format(time.RFC3339),
		CreatedBy: feature.CreatedBy,
		UpdatedBy: feature.UpdatedBy,
	}, nil
}

func (s *FeatureService) ValidateUpdateInput(input *domain.FeatureUpdateInput) error {
	feature, err := s.reader.ReadFeature(input.ProjectID, input.FeatureID)
	if err != nil {
		return err
	}

	if feature.Status == domain.FeatureStatusCancelled && !input.Reopen && !input.Cancel {
		return domain.NewValidationError("Feature is cancelled. Use --reopen to reopen, or --cancel to confirm cancellation.")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return domain.NewValidationError("Feature name cannot be empty.")
	}

	if input.Goal != nil && strings.TrimSpace(*input.Goal) == "" {
		return domain.NewValidationError("Feature goal cannot be empty.")
	}

	if input.Priority != nil && !domain.ValidatePriority(*input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Status != nil && !domain.ValidateFeatureStatus(*input.Status) {
		return domain.NewValidationError("Invalid status. Valid options: draft, active, done, blocked, cancelled")
	}

	if input.Scope != nil && !domain.ValidateScope(*input.Scope) {
		return domain.NewValidationError("Invalid scope. Valid options: frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift")
	}

	if input.DependsOn != nil {
		if err := s.validateDependencies(input.ProjectID, input.FeatureID, *input.DependsOn); err != nil {
			return err
		}
	}

	return nil
}

func (s *FeatureService) UpdateFeature(input *domain.FeatureUpdateInput) ([]string, error) {
	feature, err := s.reader.ReadFeature(input.ProjectID, input.FeatureID)
	if err != nil {
		return nil, err
	}

	if input.DryRun {
		return []string{"[DRY RUN] Would update feature: " + input.FeatureID}, nil
	}

	var changes []string
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	if input.Reopen {
		if feature.Status != domain.FeatureStatusCancelled {
			return nil, domain.NewValidationError("Feature is not cancelled. Nothing to reopen.")
		}
		feature.Status = domain.FeatureStatusDraft
		feature.Reason = ""
		changes = append(changes, "status", "reason")
	}

	if input.Cancel {
		if feature.Status == domain.FeatureStatusCancelled {
			return nil, domain.NewValidationError("Feature is already cancelled.")
		}

		dependents, err := s.findDependents(input.ProjectID, input.FeatureID)
		if err != nil {
			return nil, err
		}
		if len(dependents) > 0 && !input.Force {
			return nil, domain.NewValidationError("Feature has " + fmt.Sprintf("%d", len(dependents)) + " dependent(s). Use --force to cancel anyway.")
		}

		if input.Reason == nil || *input.Reason == "" {
			return nil, domain.NewValidationError("Cancellation reason is required (--reason).")
		}

		feature.Status = domain.FeatureStatusCancelled
		feature.Reason = *input.Reason
		changes = append(changes, "status", "reason")
	}

	if input.Name != nil && *input.Name != feature.Name {
		feature.Name = *input.Name
		changes = append(changes, "name")
	}

	if input.Goal != nil && *input.Goal != feature.Goal {
		feature.Goal = *input.Goal
		changes = append(changes, "goal")
	}

	if input.Scope != nil && *input.Scope != feature.Scope {
		feature.Scope = *input.Scope
		changes = append(changes, "scope")
	}

	if input.Priority != nil && *input.Priority != feature.Priority {
		feature.Priority = *input.Priority
		changes = append(changes, "priority")
	}

	if input.DependsOn != nil {
		feature.DependsOn = *input.DependsOn
		changes = append(changes, "depends_on")
	}

	feature.UpdatedAt = now
	feature.UpdatedBy = updater

	if err := s.writer.ReplaceFeature(input.ProjectID, feature); err != nil {
		return nil, err
	}

	event := &domain.FeatureEvent{
		Layer:   "feature",
		Type:    "updated",
		ID:      input.FeatureID,
		By:      updater,
		Ts:      now,
		Changes: changes,
	}
	if err := s.writer.AppendFeatureEvent(input.ProjectID, event); err != nil {
		return nil, err
	}

	return changes, nil
}

func (s *FeatureService) findDependents(projectID, featureID string) ([]string, error) {
	var dependents []string
	err := s.reader.ReadNDJSON(s.paths.ProjectFeaturesPath(projectID), func(raw []byte) error {
		var f domain.Feature
		if err := json.Unmarshal(raw, &f); err != nil {
			return err
		}
		for _, dep := range f.DependsOn {
			if dep == featureID {
				dependents = append(dependents, f.ID)
			}
		}
		return nil
	})
	return dependents, err
}

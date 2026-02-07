package service

import (
	"fmt"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/util"
)

type BacklogService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewBacklogService() (*BacklogService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &BacklogService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewBacklogServiceWithPaths(paths *fs.Paths) *BacklogService {
	return &BacklogService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *BacklogService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *BacklogService) ValidateCreateInput(input *domain.BacklogCreateInput) error {
	if !domain.ValidateBacklogID(input.ID) {
		return domain.NewValidationError("Invalid backlog ID. Must start with letter, contain only alphanumeric, hyphens, underscores.")
	}

	if s.reader.BacklogExists(input.ID) {
		return domain.NewValidationError("Backlog already exists: " + input.ID)
	}

	if !s.writer.BacklogsDirWritable() {
		return domain.NewPermissionError("Permission denied. Cannot create backlog directory.")
	}

	return nil
}

func (s *BacklogService) CreateBacklog(input *domain.BacklogCreateInput) error {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	backlog := &domain.Backlog{
		ID:        input.ID,
		Name:      input.Name,
		Goal:      input.Goal,
		Status:    domain.BacklogStatusInitial,
		Strict:    input.Strict,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: creator,
		UpdatedBy: creator,
	}

	if err := s.writer.CreateBacklogDir(input.ID); err != nil {
		return err
	}

	if err := s.writer.WriteBacklogMetadata(input.ID, backlog); err != nil {
		return err
	}

	schema := domain.DefaultBacklogSchema(input.TaskDep, input.FeatureDep, input.IssueDep)
	if err := s.writer.WriteBacklogSchema(input.ID, &schema); err != nil {
		return err
	}

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.BacklogEvent{
	// 	Layer: "backlog",
	// 	Type:  "created",
	// 	ID:    input.ID,
	// 	By:    creator,
	// 	Ts:    now,
	// }
	// if err := s.writer.AppendBacklogEvent(input.ID, event); err != nil {
	// 	return err
	// }

	return nil
}

func (s *BacklogService) ListBacklogs(includeDeleted, includeGoal bool) (*domain.BacklogListOutput, error) {
	backlogIDs, err := s.reader.ListBacklogs(includeDeleted)
	if err != nil {
		return nil, err
	}

	var backlogs []domain.BacklogListItem
	deletedCount := 0

	for _, id := range backlogIDs {
		backlog, err := s.reader.ReadBacklogMetadata(id)
		if err != nil {
			continue
		}

		if !includeDeleted && backlog.Status == domain.BacklogStatusDeleted {
			continue
		}

		features, _ := s.reader.CountEntityLines(s.paths.BacklogFeaturesPath(id))
		tasks, _ := s.reader.CountEntityLines(s.paths.BacklogTasksPath(id))
		issues, _ := s.reader.CountEntityLines(s.paths.BacklogIssuesPath(id))

		item := domain.BacklogListItem{
			ID:        id,
			Name:      backlog.Name,
			Goal:      backlog.Goal,
			Status:    backlog.Status,
			Features:  features,
			Tasks:     tasks,
			Issues:    issues,
			CreatedAt: backlog.CreatedAt.Format(time.RFC3339),
			UpdatedAt: backlog.UpdatedAt.Format(time.RFC3339),
		}

		if !includeGoal {
			item.Goal = ""
		}

		backlogs = append(backlogs, item)

		if backlog.Status == domain.BacklogStatusDeleted {
			deletedCount++
		}
	}

	return &domain.BacklogListOutput{
		Backlogs: backlogs,
		Total:    len(backlogs),
		Deleted:  deletedCount,
	}, nil
}

func (s *BacklogService) GetBacklogDetail(backlogID string) (*domain.BacklogDetailOutput, error) {
	backlog, err := s.reader.ReadBacklogMetadata(backlogID)
	if err != nil {
		return nil, err
	}

	schema, err := s.reader.ReadBacklogSchema(backlogID)
	if err != nil {
		return nil, err
	}

	features, _ := s.reader.CountEntityLines(s.paths.BacklogFeaturesPath(backlogID))
	tasks, _ := s.reader.CountEntityLines(s.paths.BacklogTasksPath(backlogID))
	issues, _ := s.reader.CountEntityLines(s.paths.BacklogIssuesPath(backlogID))

	return &domain.BacklogDetailOutput{
		ID:     backlogID,
		Name:   backlog.Name,
		Goal:   backlog.Goal,
		Status: backlog.Status,
		Strict: backlog.Strict,
		Schema: *schema,
		Stats: domain.BacklogStats{
			Features: domain.EntityStats{Total: features},
			Tasks:    domain.EntityStats{Total: tasks},
			Issues:   domain.EntityStats{Total: issues},
		},
		Activity: domain.ActivityInfo{
			TotalEvents:  0,
			LastActivity: "",
		},
		CreatedAt: backlog.CreatedAt.Format(time.RFC3339),
		UpdatedAt: backlog.UpdatedAt.Format(time.RFC3339),
		CreatedBy: backlog.CreatedBy,
		UpdatedBy: backlog.UpdatedBy,
	}, nil
}

func (s *BacklogService) ValidateUpdateInput(input *domain.BacklogUpdateInput) error {
	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return err
	}

	if backlog.Status == domain.BacklogStatusDeleted {
		return domain.NewValidationError("Cannot update deleted backlog: " + input.ID)
	}

	if !s.writer.CheckBacklogWritable(input.ID) {
		return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.BacklogMetadataPath(input.ID))
	}

	return nil
}

func (s *BacklogService) UpdateBacklog(input *domain.BacklogUpdateInput) ([]string, error) {
	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return nil, err
	}

	var changes []string
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	if input.Name != nil {
		if *input.Name == "" {
			return nil, domain.NewValidationError("Backlog name cannot be empty.")
		}
		backlog.Name = *input.Name
		changes = append(changes, "name")
	}

	if input.Goal != nil {
		if *input.Goal == "" {
			return nil, domain.NewValidationError("Backlog goal cannot be empty.")
		}
		minLen := s.GetBacklogGoalMinLength()
		if !domain.ValidateGoalLength(*input.Goal, minLen) {
			return nil, domain.NewValidationError(fmt.Sprintf("Backlog goal must be at least %d characters.", minLen))
		}
		backlog.Goal = *input.Goal
		changes = append(changes, "goal")
	}

	if input.Strict != nil {
		backlog.Strict = *input.Strict
		changes = append(changes, "strict")
	}

	backlog.UpdatedAt = now
	backlog.UpdatedBy = updater

	if err := s.writer.WriteBacklogMetadata(input.ID, backlog); err != nil {
		return nil, err
	}

	schemaChanged := false
	if input.TaskDep != nil || input.FeatureDep != nil || input.IssueDep != nil {
		schema, err := s.reader.ReadBacklogSchema(input.ID)
		if err != nil {
			return nil, err
		}

		if input.TaskDep != nil {
			if !domain.ValidateDependencyRule(*input.TaskDep) {
				return nil, domain.NewValidationError("Invalid value for --task-dep. Valid options: same_backlog_only, cross_backlog_allowed, disabled")
			}
			schema.Rules.Task.Dependency = *input.TaskDep
			changes = append(changes, "task_dep")
			schemaChanged = true
		}

		if input.FeatureDep != nil {
			if !domain.ValidateDependencyRule(*input.FeatureDep) {
				return nil, domain.NewValidationError("Invalid value for --feature-dep. Valid options: same_backlog_only, cross_backlog_allowed, disabled")
			}
			schema.Rules.Feature.Dependency = *input.FeatureDep
			changes = append(changes, "feature_dep")
			schemaChanged = true
		}

		if input.IssueDep != nil {
			if !domain.ValidateDependencyRule(*input.IssueDep) {
				return nil, domain.NewValidationError("Invalid value for --issue-dep. Valid options: same_backlog_only, cross_backlog_allowed, disabled")
			}
			schema.Rules.Issue.Dependency = *input.IssueDep
			changes = append(changes, "issue_dep")
			schemaChanged = true
		}

		if schemaChanged {
			if err := s.writer.WriteBacklogSchema(input.ID, schema); err != nil {
				return nil, err
			}
		}
	}

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.BacklogEvent{
	// 	Layer:   "backlog",
	// 	Type:    "updated",
	// 	ID:      input.ID,
	// 	By:      updater,
	// 	Ts:      now,
	// 	Changes: changes,
	// }
	// if err := s.writer.AppendBacklogEvent(input.ID, event); err != nil {
	// 	return nil, err
	// }

	return changes, nil
}

func (s *BacklogService) ValidateDeleteInput(input *domain.BacklogDeleteInput) error {
	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return err
	}

	if backlog.Status == domain.BacklogStatusDeleted && !input.Hard {
		return domain.NewValidationError("Backlog is already deleted: " + input.ID + ". Use --hard to permanently remove.")
	}

	if !input.Hard && !input.DryRun {
		if !s.writer.CheckBacklogWritable(input.ID) {
			return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.BacklogMetadataPath(input.ID))
		}
	}

	if input.Hard && !input.DryRun {
		if !s.writer.CheckBacklogWritable(input.ID) {
			return domain.NewPermissionError("Permission denied. Cannot delete " + s.paths.BacklogDirPath(input.ID))
		}
	}

	return nil
}

func (s *BacklogService) DeleteBacklog(input *domain.BacklogDeleteInput) (string, error) {
	if input.DryRun {
		if input.Hard {
			return "[DRY RUN] Would hard delete backlog: " + input.ID, nil
		}
		return "[DRY RUN] Would soft delete backlog: " + input.ID, nil
	}

	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return "", err
	}

	if input.Hard {
		if err := s.writer.DeleteBacklogDir(input.ID); err != nil {
			return "", err
		}
		return "Backlog permanently deleted: " + input.ID, nil
	}

	updater := util.GetGitUsername()
	now := time.Now().UTC()

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.BacklogEvent{
	// 	Layer: "backlog",
	// 	Type:  "deleted",
	// 	ID:    input.ID,
	// 	By:    updater,
	// 	Ts:    now,
	// }
	// if err := s.writer.AppendBacklogEvent(input.ID, event); err != nil {
	// 	return "", err
	// }

	backlog.Status = domain.BacklogStatusDeleted
	backlog.UpdatedAt = now
	backlog.UpdatedBy = updater

	if err := s.writer.WriteBacklogMetadata(input.ID, backlog); err != nil {
		return "", err
	}

	return "Backlog deleted: " + input.ID, nil
}

func (s *BacklogService) ValidateReopenInput(input *domain.BacklogReopenInput) error {
	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return err
	}

	if backlog.Status != domain.BacklogStatusDeleted {
		return domain.NewValidationError("Backlog is not deleted: " + input.ID + ". Nothing to reopen.")
	}

	if !s.writer.CheckBacklogWritable(input.ID) {
		return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.BacklogMetadataPath(input.ID))
	}

	return nil
}

func (s *BacklogService) ReopenBacklog(input *domain.BacklogReopenInput) (string, error) {
	backlog, err := s.reader.ReadBacklogMetadata(input.ID)
	if err != nil {
		return "", err
	}

	updater := util.GetGitUsername()
	now := time.Now().UTC()

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.BacklogEvent{
	// 	Layer: "backlog",
	// 	Type:  "reopened",
	// 	ID:    input.ID,
	// 	By:    updater,
	// 	Ts:    now,
	// }
	// if err := s.writer.AppendBacklogEvent(input.ID, event); err != nil {
	// 	return "", err
	// }

	backlog.Status = domain.BacklogStatusInitial
	backlog.UpdatedAt = now
	backlog.UpdatedBy = updater

	if err := s.writer.WriteBacklogMetadata(input.ID, backlog); err != nil {
		return "", err
	}

	return "Backlog reopened: " + input.ID, nil
}

func (s *BacklogService) GetBacklog(backlogID string) (*domain.Backlog, error) {
	return s.reader.ReadBacklogMetadata(backlogID)
}

func (s *BacklogService) GetBacklogGoalMinLength() int {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return domain.GoalMinLength
	}
	if ws.Config.GoalLengths.Backlog > 0 {
		return ws.Config.GoalLengths.Backlog
	}
	return domain.GoalMinLength
}

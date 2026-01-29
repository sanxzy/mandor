package service

import (
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/util"
)

type ProjectService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewProjectService() (*ProjectService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &ProjectService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewProjectServiceWithPaths(paths *fs.Paths) *ProjectService {
	return &ProjectService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *ProjectService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *ProjectService) ValidateCreateInput(input *domain.ProjectCreateInput) error {
	if !domain.ValidateProjectID(input.ID) {
		return domain.NewValidationError("Invalid project ID. Must start with letter, contain only alphanumeric, hyphens, underscores.")
	}

	if s.reader.ProjectExists(input.ID) {
		return domain.NewValidationError("Project already exists: " + input.ID)
	}

	if !s.writer.ProjectsDirWritable() {
		return domain.NewPermissionError("Permission denied. Cannot create project directory.")
	}

	return nil
}

func (s *ProjectService) CreateProject(input *domain.ProjectCreateInput) error {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	project := &domain.Project{
		ID:        input.ID,
		Name:      input.Name,
		Goal:      input.Goal,
		Status:    domain.ProjectStatusInitial,
		Strict:    input.Strict,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: creator,
		UpdatedBy: creator,
	}

	if err := s.writer.CreateProjectDir(input.ID); err != nil {
		return err
	}

	if err := s.writer.WriteProjectMetadata(input.ID, project); err != nil {
		return err
	}

	schema := domain.DefaultProjectSchema(input.TaskDep, input.FeatureDep, input.IssueDep)
	if err := s.writer.WriteProjectSchema(input.ID, &schema); err != nil {
		return err
	}

	event := &domain.ProjectEvent{
		Layer: "project",
		Type:  "created",
		ID:    input.ID,
		By:    creator,
		Ts:    now,
	}
	if err := s.writer.AppendProjectEvent(input.ID, event); err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) ListProjects(includeDeleted, includeGoal bool) (*domain.ProjectListOutput, error) {
	projectIDs, err := s.reader.ListProjects(includeDeleted)
	if err != nil {
		return nil, err
	}

	var projects []domain.ProjectListItem
	deletedCount := 0

	for _, id := range projectIDs {
		project, err := s.reader.ReadProjectMetadata(id)
		if err != nil {
			continue
		}

		if !includeDeleted && project.Status == domain.ProjectStatusDeleted {
			continue
		}

		features, _ := s.reader.CountEntityLines(s.paths.ProjectFeaturesPath(id))
		tasks, _ := s.reader.CountEntityLines(s.paths.ProjectTasksPath(id))
		issues, _ := s.reader.CountEntityLines(s.paths.ProjectIssuesPath(id))

		item := domain.ProjectListItem{
			ID:        id,
			Name:      project.Name,
			Goal:      project.Goal,
			Status:    project.Status,
			Features:  features,
			Tasks:     tasks,
			Issues:    issues,
			CreatedAt: project.CreatedAt.Format(time.RFC3339),
			UpdatedAt: project.UpdatedAt.Format(time.RFC3339),
		}

		if !includeGoal {
			item.Goal = ""
		}

		projects = append(projects, item)

		if project.Status == domain.ProjectStatusDeleted {
			deletedCount++
		}
	}

	return &domain.ProjectListOutput{
		Projects: projects,
		Total:    len(projects),
		Deleted:  deletedCount,
	}, nil
}

func (s *ProjectService) GetProjectDetail(projectID string) (*domain.ProjectDetailOutput, error) {
	project, err := s.reader.ReadProjectMetadata(projectID)
	if err != nil {
		return nil, err
	}

	schema, err := s.reader.ReadProjectSchema(projectID)
	if err != nil {
		return nil, err
	}

	features, _ := s.reader.CountEntityLines(s.paths.ProjectFeaturesPath(projectID))
	tasks, _ := s.reader.CountEntityLines(s.paths.ProjectTasksPath(projectID))
	issues, _ := s.reader.CountEntityLines(s.paths.ProjectIssuesPath(projectID))
	events, _ := s.reader.CountEventLines(projectID)

	lastActivity := ""
	if events > 0 {
		lastActivity = project.UpdatedAt.Format(time.RFC3339)
	}

	return &domain.ProjectDetailOutput{
		ID:     projectID,
		Name:   project.Name,
		Goal:   project.Goal,
		Status: project.Status,
		Strict: project.Strict,
		Schema: *schema,
		Stats: domain.ProjectStats{
			Features: domain.EntityStats{Total: features},
			Tasks:    domain.EntityStats{Total: tasks},
			Issues:   domain.EntityStats{Total: issues},
		},
		Activity: domain.ActivityInfo{
			TotalEvents:  events,
			LastActivity: lastActivity,
		},
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
		UpdatedAt: project.UpdatedAt.Format(time.RFC3339),
		CreatedBy: project.CreatedBy,
		UpdatedBy: project.UpdatedBy,
	}, nil
}

func (s *ProjectService) ValidateUpdateInput(input *domain.ProjectUpdateInput) error {
	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return err
	}

	if project.Status == domain.ProjectStatusDeleted {
		return domain.NewValidationError("Cannot update deleted project: " + input.ID)
	}

	if !s.writer.CheckProjectWritable(input.ID) {
		return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.ProjectMetadataPath(input.ID))
	}

	return nil
}

func (s *ProjectService) UpdateProject(input *domain.ProjectUpdateInput) ([]string, error) {
	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return nil, err
	}

	var changes []string
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	if input.Name != nil {
		if *input.Name == "" {
			return nil, domain.NewValidationError("Project name cannot be empty.")
		}
		project.Name = *input.Name
		changes = append(changes, "name")
	}

	if input.Goal != nil {
		if *input.Goal == "" {
			return nil, domain.NewValidationError("Project goal cannot be empty.")
		}
		if !domain.ValidateGoalLength(*input.Goal) {
			return nil, domain.NewValidationError("Project goal must be at least 500 characters.")
		}
		project.Goal = *input.Goal
		changes = append(changes, "goal")
	}

	if input.Strict != nil {
		project.Strict = *input.Strict
		changes = append(changes, "strict")
	}

	project.UpdatedAt = now
	project.UpdatedBy = updater

	if err := s.writer.WriteProjectMetadata(input.ID, project); err != nil {
		return nil, err
	}

	schemaChanged := false
	if input.TaskDep != nil || input.FeatureDep != nil || input.IssueDep != nil {
		schema, err := s.reader.ReadProjectSchema(input.ID)
		if err != nil {
			return nil, err
		}

		if input.TaskDep != nil {
			if !domain.ValidateDependencyRule(*input.TaskDep) {
				return nil, domain.NewValidationError("Invalid value for --task-dep. Valid options: same_project_only, cross_project_allowed, disabled")
			}
			schema.Rules.Task.Dependency = *input.TaskDep
			changes = append(changes, "task_dep")
			schemaChanged = true
		}

		if input.FeatureDep != nil {
			if !domain.ValidateDependencyRule(*input.FeatureDep) {
				return nil, domain.NewValidationError("Invalid value for --feature-dep. Valid options: same_project_only, cross_project_allowed, disabled")
			}
			schema.Rules.Feature.Dependency = *input.FeatureDep
			changes = append(changes, "feature_dep")
			schemaChanged = true
		}

		if input.IssueDep != nil {
			if !domain.ValidateDependencyRule(*input.IssueDep) {
				return nil, domain.NewValidationError("Invalid value for --issue-dep. Valid options: same_project_only, cross_project_allowed, disabled")
			}
			schema.Rules.Issue.Dependency = *input.IssueDep
			changes = append(changes, "issue_dep")
			schemaChanged = true
		}

		if schemaChanged {
			if err := s.writer.WriteProjectSchema(input.ID, schema); err != nil {
				return nil, err
			}
		}
	}

	event := &domain.ProjectEvent{
		Layer:   "project",
		Type:    "updated",
		ID:      input.ID,
		By:      updater,
		Ts:      now,
		Changes: changes,
	}
	if err := s.writer.AppendProjectEvent(input.ID, event); err != nil {
		return nil, err
	}

	return changes, nil
}

func (s *ProjectService) ValidateDeleteInput(input *domain.ProjectDeleteInput) error {
	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return err
	}

	if project.Status == domain.ProjectStatusDeleted && !input.Hard {
		return domain.NewValidationError("Project is already deleted: " + input.ID + ". Use --hard to permanently remove.")
	}

	if !input.Hard && !input.DryRun {
		if !s.writer.CheckProjectWritable(input.ID) {
			return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.ProjectMetadataPath(input.ID))
		}
	}

	if input.Hard && !input.DryRun {
		if !s.writer.CheckProjectWritable(input.ID) {
			return domain.NewPermissionError("Permission denied. Cannot delete " + s.paths.ProjectDirPath(input.ID))
		}
	}

	return nil
}

func (s *ProjectService) DeleteProject(input *domain.ProjectDeleteInput) (string, error) {
	if input.DryRun {
		if input.Hard {
			return "[DRY RUN] Would hard delete project: " + input.ID, nil
		}
		return "[DRY RUN] Would soft delete project: " + input.ID, nil
	}

	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return "", err
	}

	if input.Hard {
		if err := s.writer.DeleteProjectDir(input.ID); err != nil {
			return "", err
		}
		return "Project permanently deleted: " + input.ID, nil
	}

	updater := util.GetGitUsername()
	now := time.Now().UTC()

	event := &domain.ProjectEvent{
		Layer: "project",
		Type:  "deleted",
		ID:    input.ID,
		By:    updater,
		Ts:    now,
	}
	if err := s.writer.AppendProjectEvent(input.ID, event); err != nil {
		return "", err
	}

	project.Status = domain.ProjectStatusDeleted
	project.UpdatedAt = now
	project.UpdatedBy = updater

	if err := s.writer.WriteProjectMetadata(input.ID, project); err != nil {
		return "", err
	}

	return "Project deleted: " + input.ID, nil
}

func (s *ProjectService) ValidateReopenInput(input *domain.ProjectReopenInput) error {
	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return err
	}

	if project.Status != domain.ProjectStatusDeleted {
		return domain.NewValidationError("Project is not deleted: " + input.ID + ". Nothing to reopen.")
	}

	if !s.writer.CheckProjectWritable(input.ID) {
		return domain.NewPermissionError("Permission denied. Cannot write to " + s.paths.ProjectMetadataPath(input.ID))
	}

	return nil
}

func (s *ProjectService) ReopenProject(input *domain.ProjectReopenInput) (string, error) {
	project, err := s.reader.ReadProjectMetadata(input.ID)
	if err != nil {
		return "", err
	}

	updater := util.GetGitUsername()
	now := time.Now().UTC()

	event := &domain.ProjectEvent{
		Layer: "project",
		Type:  "reopened",
		ID:    input.ID,
		By:    updater,
		Ts:    now,
	}
	if err := s.writer.AppendProjectEvent(input.ID, event); err != nil {
		return "", err
	}

	project.Status = domain.ProjectStatusInitial
	project.UpdatedAt = now
	project.UpdatedBy = updater

	if err := s.writer.WriteProjectMetadata(input.ID, project); err != nil {
		return "", err
	}

	return "Project reopened: " + input.ID, nil
}

func (s *ProjectService) GetProject(projectID string) (*domain.Project, error) {
	return s.reader.ReadProjectMetadata(projectID)
}

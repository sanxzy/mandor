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

type IssueService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewIssueService() (*IssueService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &IssueService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewIssueServiceWithPaths(paths *fs.Paths) *IssueService {
	return &IssueService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *IssueService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *IssueService) ValidateCreateInput(input *domain.IssueCreateInput) error {
	if !s.reader.ProjectExists(input.ProjectID) {
		return domain.NewValidationError("Project not found: " + input.ProjectID)
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.NewValidationError("Issue name is required (--name).")
	}

	if strings.TrimSpace(input.Goal) == "" {
		return domain.NewValidationError("Issue goal is required (--goal).")
	}

	if !domain.ValidateIssueGoalLength(input.Goal) {
		minLen := domain.IssueGoalMinLength
		if util.IsDevelopment() {
			minLen = domain.IssueGoalMinLengthDevelopment
		}
		return domain.NewValidationError(fmt.Sprintf("Issue goal must be at least %d characters. Current length: %d characters.", minLen, len(input.Goal)))
	}

	if strings.TrimSpace(input.IssueType) == "" {
		return domain.NewValidationError("Issue type is required (--type).")
	}

	if !domain.ValidateIssueType(input.IssueType) {
		return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
	}

	if len(input.AffectedFiles) == 0 {
		return domain.NewValidationError("Affected files are required (--affected-files).")
	}

	if len(input.AffectedTests) == 0 {
		return domain.NewValidationError("Affected tests are required (--affected-tests).")
	}

	if len(input.ImplementationSteps) == 0 {
		return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
	}

	if !domain.ValidatePriority(input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Priority == "" {
		input.Priority = "P2"
	}

	if err := s.validateDependencies(input.ProjectID, "", input.DependsOn); err != nil {
		return err
	}

	return nil
}

func (s *IssueService) validateDependencies(projectID, selfID string, dependsOn []string) error {
	// Read schema to check if cross-project dependencies are allowed
	schema, err := s.reader.ReadProjectSchema(projectID)
	if err != nil {
		return domain.NewSystemError("Cannot read project schema", err)
	}
	allowCrossProject := schema.Rules.Issue.Dependency != "same_project_only" && schema.Rules.Issue.Dependency != "disabled"

	for _, depID := range dependsOn {
		if depID == selfID {
			return domain.NewValidationError("Self-dependency detected. Issue cannot depend on itself.")
		}

		// Parse issue ID to get project
		depProjectID := extractProjectIDFromIssueID(depID)
		if depProjectID == "" {
			return domain.NewValidationError("Invalid issue ID format: " + depID)
		}

		// Check if cross-project and not allowed
		if depProjectID != projectID && !allowCrossProject {
			return domain.NewValidationError(fmt.Sprintf("Cross-project dependency detected: %s -> %s. Cross-project dependencies are disabled.", selfID, depID))
		}

		dep, err := s.reader.ReadIssue(depProjectID, depID)
		if err != nil {
			if _, ok := err.(*domain.MandorError); ok {
				return domain.NewValidationError("Dependency not found: " + depID)
			}
			return err
		}

		if dep.Status == domain.IssueStatusCancelled {
			return domain.NewValidationError("Dependency is cancelled: " + depID)
		}
	}

	if err := s.validateNoCycle(projectID, selfID, dependsOn); err != nil {
		return err
	}

	return nil
}

func (s *IssueService) validateNoCycle(projectID, selfID string, dependsOn []string) error {
	visited := make(map[string]bool)
	var dfs func(issueID string) bool

	dfs = func(issueID string) bool {
		if issueID == selfID {
			return true
		}
		if visited[issueID] {
			return false
		}
		visited[issueID] = true

		i, err := s.reader.ReadIssue(projectID, issueID)
		if err != nil {
			return false
		}

		for _, dep := range i.DependsOn {
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

func (s *IssueService) CreateIssue(input *domain.IssueCreateInput) (*domain.Issue, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	nanoid, err := util.GenerateID()
	if err != nil {
		return nil, domain.NewSystemError("Failed to generate issue ID", err)
	}

	issueID := input.ProjectID + "-issue-" + nanoid

	issue := &domain.Issue{
		ID:                  issueID,
		ProjectID:           input.ProjectID,
		Name:                input.Name,
		Goal:                input.Goal,
		IssueType:           input.IssueType,
		Priority:            input.Priority,
		Status:              domain.IssueStatusOpen,
		DependsOn:           input.DependsOn,
		AffectedFiles:       input.AffectedFiles,
		AffectedTests:       input.AffectedTests,
		ImplementationSteps: input.ImplementationSteps,
		LibraryNeeds:        input.LibraryNeeds,
		CreatedAt:           now,
		LastUpdatedAt:       now,
		CreatedBy:           creator,
		LastUpdatedBy:       creator,
	}

	if len(input.DependsOn) > 0 {
		allResolved, err := s.checkDependenciesResolved(input.ProjectID, input.DependsOn)
		if err != nil {
			return nil, err
		}
		if allResolved {
			issue.Status = domain.IssueStatusReady
		} else {
			issue.Status = domain.IssueStatusBlocked
		}
	} else {
		issue.Status = domain.IssueStatusReady
	}

	if err := s.writer.WriteIssue(input.ProjectID, issue); err != nil {
		return nil, err
	}

	event := &domain.IssueEvent{
		Layer: "issue",
		Type:  "created",
		ID:    issueID,
		By:    creator,
		Ts:    now,
	}
	if err := s.writer.AppendIssueEvent(input.ProjectID, event); err != nil {
		return nil, err
	}

	if issue.Status == domain.IssueStatusReady {
		readyEvent := &domain.IssueEvent{
			Layer: "issue",
			Type:  "ready",
			ID:    issueID,
			By:    "system",
			Ts:    now,
		}
		if err := s.writer.AppendIssueEvent(input.ProjectID, readyEvent); err != nil {
			return nil, err
		}
	}

	if issue.Status == domain.IssueStatusBlocked {
		blockedEvent := &domain.IssueEvent{
			Layer: "issue",
			Type:  "blocked",
			ID:    issueID,
			By:    "system",
			Ts:    now,
		}
		if err := s.writer.AppendIssueEvent(input.ProjectID, blockedEvent); err != nil {
			return nil, err
		}
	}

	return issue, nil
}

func (s *IssueService) checkDependenciesResolved(projectID string, dependsOn []string) (bool, error) {
	for _, depID := range dependsOn {
		dep, err := s.reader.ReadIssue(projectID, depID)
		if err != nil {
			return false, domain.NewValidationError("Dependency not found: " + depID)
		}
		if dep.Status != domain.IssueStatusResolved && dep.Status != domain.IssueStatusWontFix {
			return false, nil
		}
	}
	return true, nil
}

func (s *IssueService) ListIssues(input *domain.IssueListInput) (*domain.IssueListOutput, error) {
	if !s.reader.ProjectExists(input.ProjectID) {
		return nil, domain.NewValidationError("Project not found: " + input.ProjectID)
	}

	var issues []domain.IssueListItem
	deletedCount := 0

	err := s.reader.ReadNDJSON(s.paths.ProjectIssuesPath(input.ProjectID), func(raw []byte) error {
		var i domain.Issue
		if err := json.Unmarshal(raw, &i); err != nil {
			return err
		}
		if !input.IncludeDeleted && i.Status == domain.IssueStatusCancelled {
			return nil
		}

		if input.IssueType != "" && i.IssueType != input.IssueType {
			return nil
		}

		if input.Status != "" && i.Status != input.Status {
			return nil
		}

		if input.Priority != "" && i.Priority != input.Priority {
			return nil
		}

		item := domain.IssueListItem{
			ID:                       i.ID,
			Name:                     i.Name,
			IssueType:                i.IssueType,
			Status:                   i.Status,
			Priority:                 i.Priority,
			ProjectID:                i.ProjectID,
			DependsOnCount:           len(i.DependsOn),
			AffectedFilesCount:       len(i.AffectedFiles),
			AffectedTestsCount:       len(i.AffectedTests),
			ImplementationStepsCount: len(i.ImplementationSteps),
			LibraryNeedsCount:        len(i.LibraryNeeds),
			CreatedAt:                i.CreatedAt.Format(time.RFC3339),
			LastUpdatedAt:            i.LastUpdatedAt.Format(time.RFC3339),
		}
		issues = append(issues, item)

		if i.Status == domain.IssueStatusCancelled {
			deletedCount++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &domain.IssueListOutput{
		Issues:  issues,
		Total:   len(issues),
		Deleted: deletedCount,
	}, nil
}

func (s *IssueService) GetIssueDetail(input *domain.IssueDetailInput) (*domain.IssueDetailOutput, error) {
	issue, err := s.reader.ReadIssue(input.ProjectID, input.IssueID)
	if err != nil {
		return nil, err
	}

	if !input.IncludeDeleted && issue.Status == domain.IssueStatusCancelled {
		return nil, domain.NewValidationError("Issue not found: " + input.IssueID)
	}

	events, _ := s.reader.CountEventLines(input.ProjectID)

	return &domain.IssueDetailOutput{
		ID:                  issue.ID,
		ProjectID:           issue.ProjectID,
		Name:                issue.Name,
		Goal:                issue.Goal,
		IssueType:           issue.IssueType,
		Priority:            issue.Priority,
		Status:              issue.Status,
		DependsOn:           issue.DependsOn,
		Reason:              issue.Reason,
		AffectedFiles:       issue.AffectedFiles,
		AffectedTests:       issue.AffectedTests,
		ImplementationSteps: issue.ImplementationSteps,
		LibraryNeeds:        issue.LibraryNeeds,
		Events:              events,
		CreatedAt:           issue.CreatedAt.Format(time.RFC3339),
		LastUpdatedAt:       issue.LastUpdatedAt.Format(time.RFC3339),
		CreatedBy:           issue.CreatedBy,
		LastUpdatedBy:       issue.LastUpdatedBy,
	}, nil
}

func (s *IssueService) ValidateUpdateInput(input *domain.IssueUpdateInput) error {
	issue, err := s.reader.ReadIssue(input.ProjectID, input.IssueID)
	if err != nil {
		return err
	}

	if issue.Status == domain.IssueStatusCancelled && !input.Reopen && !input.Cancel {
		return domain.NewValidationError("Issue is cancelled. Use --reopen to reopen, or --cancel to confirm cancellation.")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return domain.NewValidationError("Issue name cannot be empty.")
	}

	if input.Goal != nil && strings.TrimSpace(*input.Goal) == "" {
		return domain.NewValidationError("Issue goal cannot be empty.")
	}

	if input.IssueType != nil && !domain.ValidateIssueType(*input.IssueType) {
		return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
	}

	if input.Priority != nil && !domain.ValidatePriority(*input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Status != nil && !domain.ValidateIssueStatus(*input.Status) {
		return domain.NewValidationError("Invalid status. Valid options: open, ready, in_progress, blocked, resolved, wontfix, cancelled")
	}

	if input.DependsOn != nil {
		if err := s.validateDependencies(input.ProjectID, input.IssueID, *input.DependsOn); err != nil {
			return err
		}
	}

	return nil
}

func (s *IssueService) UpdateIssue(input *domain.IssueUpdateInput) ([]string, error) {
	issue, err := s.reader.ReadIssue(input.ProjectID, input.IssueID)
	if err != nil {
		return nil, err
	}

	if input.DryRun {
		return []string{"[DRY RUN] Would update issue: " + input.IssueID}, nil
	}

	var changes []string
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	if input.Reopen {
		if !domain.IsIssueTerminalStatus(issue.Status) {
			return nil, domain.NewValidationError("Issue is not in terminal state. Only resolved, wontfix, or cancelled issues can be reopened.")
		}
		issue.Status = domain.IssueStatusOpen
		issue.Reason = ""
		changes = append(changes, "status", "reason")
	}

	if input.Cancel {
		if issue.Status == domain.IssueStatusCancelled {
			return nil, domain.NewValidationError("Issue is already cancelled.")
		}

		if input.Reason == nil || *input.Reason == "" {
			return nil, domain.NewValidationError("Cancellation reason is required (--reason).")
		}

		issue.Status = domain.IssueStatusCancelled
		issue.Reason = *input.Reason
		changes = append(changes, "status", "reason")
	}

	if input.Name != nil && *input.Name != issue.Name {
		issue.Name = *input.Name
		changes = append(changes, "name")
	}

	if input.Goal != nil && *input.Goal != issue.Goal {
		issue.Goal = *input.Goal
		changes = append(changes, "goal")
	}

	if input.IssueType != nil && *input.IssueType != issue.IssueType {
		issue.IssueType = *input.IssueType
		changes = append(changes, "issue_type")
	}

	if input.Priority != nil && *input.Priority != issue.Priority {
		issue.Priority = *input.Priority
		changes = append(changes, "priority")
	}

	if input.DependsOn != nil {
		issue.DependsOn = *input.DependsOn
		changes = append(changes, "depends_on")
	}

	if input.AffectedFiles != nil {
		issue.AffectedFiles = *input.AffectedFiles
		changes = append(changes, "affected_files")
	}

	if input.AffectedTests != nil {
		issue.AffectedTests = *input.AffectedTests
		changes = append(changes, "affected_tests")
	}

	if input.ImplementationSteps != nil {
		issue.ImplementationSteps = *input.ImplementationSteps
		changes = append(changes, "implementation_steps")
	}

	if input.LibraryNeeds != nil {
		issue.LibraryNeeds = *input.LibraryNeeds
		changes = append(changes, "library_needs")
	}

	if input.Start {
		if issue.Status != domain.IssueStatusOpen && issue.Status != domain.IssueStatusReady {
			return nil, domain.NewValidationError("Issue is not in startable state (open or ready).")
		}
		issue.Status = domain.IssueStatusInProgress
		changes = append(changes, "status")
	}

	if input.Resolve {
		if issue.Status == domain.IssueStatusResolved || issue.Status == domain.IssueStatusWontFix || issue.Status == domain.IssueStatusCancelled {
			return nil, domain.NewValidationError("Issue is already resolved, wontfix, or cancelled.")
		}
		issue.Status = domain.IssueStatusResolved
		changes = append(changes, "status")
	}

	if input.WontFix {
		if issue.Status == domain.IssueStatusResolved || issue.Status == domain.IssueStatusWontFix || issue.Status == domain.IssueStatusCancelled {
			return nil, domain.NewValidationError("Issue is already resolved, wontfix, or cancelled.")
		}
		if input.Reason == nil || *input.Reason == "" {
			return nil, domain.NewValidationError("Wontfix reason is required (--reason).")
		}
		issue.Status = domain.IssueStatusWontFix
		issue.Reason = *input.Reason
		changes = append(changes, "status", "reason")
	}

	if input.Status != nil && *input.Status != issue.Status {
		if err := s.validateStatusTransition(issue.Status, *input.Status); err != nil {
			return nil, err
		}
		issue.Status = *input.Status
		changes = append(changes, "status")
	}

	issue.LastUpdatedAt = now
	issue.LastUpdatedBy = updater

	if err := s.writer.ReplaceIssue(input.ProjectID, issue); err != nil {
		return nil, err
	}

	if input.Resolve || input.WontFix {
		unblocked, err := s.unblockDependents(input.ProjectID, issue.ID)
		if err != nil {
			return nil, err
		}
		if unblocked {
			changes = append(changes, "dependent_unblocked")
		}
	}

	event := &domain.IssueEvent{
		Layer:   "issue",
		Type:    "updated",
		ID:      input.IssueID,
		By:      updater,
		Ts:      now,
		Changes: changes,
	}
	if err := s.writer.AppendIssueEvent(input.ProjectID, event); err != nil {
		return nil, err
	}

	return changes, nil
}

func (s *IssueService) validateStatusTransition(current, next string) error {
	validTransitions := map[string][]string{
		domain.IssueStatusOpen:       {domain.IssueStatusReady, domain.IssueStatusInProgress, domain.IssueStatusBlocked, domain.IssueStatusResolved, domain.IssueStatusWontFix, domain.IssueStatusCancelled},
		domain.IssueStatusReady:      {domain.IssueStatusInProgress, domain.IssueStatusBlocked, domain.IssueStatusResolved, domain.IssueStatusWontFix, domain.IssueStatusCancelled},
		domain.IssueStatusInProgress: {domain.IssueStatusBlocked, domain.IssueStatusResolved, domain.IssueStatusWontFix, domain.IssueStatusCancelled},
		domain.IssueStatusBlocked:    {domain.IssueStatusReady, domain.IssueStatusResolved, domain.IssueStatusWontFix, domain.IssueStatusCancelled},
	}

	allowed, ok := validTransitions[current]
	if !ok {
		return domain.NewValidationError(fmt.Sprintf("Cannot transition from %s", current))
	}

	for _, allowedStatus := range allowed {
		if next == allowedStatus {
			return nil
		}
	}

	return domain.NewValidationError(fmt.Sprintf("Invalid status transition from %s to %s", current, next))
}

func (s *IssueService) FindDependents(projectID, issueID string) ([]string, error) {
	var dependents []string
	err := s.reader.ReadNDJSON(s.paths.ProjectIssuesPath(projectID), func(raw []byte) error {
		var i domain.Issue
		if err := json.Unmarshal(raw, &i); err != nil {
			return err
		}
		for _, dep := range i.DependsOn {
			if dep == issueID {
				dependents = append(dependents, i.ID)
			}
		}
		return nil
	})
	return dependents, err
}

func (s *IssueService) GetWorkspace() (*domain.Workspace, error) {
	return s.reader.ReadWorkspace()
}

func (s *IssueService) ProjectExists(projectID string) bool {
	return s.reader.ProjectExists(projectID)
}

func (s *IssueService) ReadDependency(projectID, issueID string) (*domain.Issue, error) {
	return s.reader.ReadIssue(projectID, issueID)
}

func (s *IssueService) GetIssueEvents(projectID, issueID string) ([]domain.IssueEvent, error) {
	var events []domain.IssueEvent
	err := s.reader.ReadNDJSON(s.paths.ProjectEventsPath(projectID), func(raw []byte) error {
		var event domain.IssueEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			return err
		}
		if event.Layer == "issue" && event.ID == issueID {
			events = append(events, event)
		}
		return nil
	})
	return events, err
}

func (s *IssueService) unblockDependents(projectID, resolvedIssueID string) (bool, error) {
	unblockedAny := false
	now := time.Now().UTC()

	err := s.reader.ReadNDJSON(s.paths.ProjectIssuesPath(projectID), func(raw []byte) error {
		var issue domain.Issue
		if err := json.Unmarshal(raw, &issue); err != nil {
			return err
		}

		if issue.Status != domain.IssueStatusBlocked {
			return nil
		}

		hasResolved := false
		allResolved := true
		for _, depID := range issue.DependsOn {
			if depID == resolvedIssueID {
				hasResolved = true
			}
			// Parse the dependency ID to get the project it belongs to
			depProjectID := extractProjectIDFromIssueID(depID)
			if depProjectID == "" {
				return domain.NewValidationError("Invalid issue ID format: " + depID)
			}
			dep, err := s.reader.ReadIssue(depProjectID, depID)
			if err != nil {
				return err
			}
			if dep.Status != domain.IssueStatusResolved && dep.Status != domain.IssueStatusWontFix {
				allResolved = false
			}
		}

		if hasResolved && allResolved {
			issue.Status = domain.IssueStatusReady
			issue.LastUpdatedAt = now
			if err := s.writer.ReplaceIssue(projectID, &issue); err != nil {
				return err
			}

			event := &domain.IssueEvent{
				Layer: "issue",
				Type:  "ready",
				ID:    issue.ID,
				By:    "system",
				Ts:    now,
			}
			if err := s.writer.AppendIssueEvent(projectID, event); err != nil {
				return err
			}
			unblockedAny = true
		}

		return nil
	})

	return unblockedAny, err
}

// extractProjectIDFromIssueID extracts the project ID from an issue ID
// Issue ID format: <project>-issue-<nanoid>
func extractProjectIDFromIssueID(issueID string) string {
	parts := strings.Split(issueID, "-issue-")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

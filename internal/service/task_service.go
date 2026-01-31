package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/util"
)

type TaskService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewTaskService() (*TaskService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &TaskService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func NewTaskServiceWithPaths(paths *fs.Paths) *TaskService {
	return &TaskService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}
}

func (s *TaskService) WorkspaceInitialized() bool {
	return s.reader.WorkspaceExists()
}

func (s *TaskService) ParseTaskID(taskID string) (projectID, featureID string, err error) {
	// Task ID format: <project>-feature-<feature>-task-<nanoid>
	// Find the last occurrence of "-task-" to handle project IDs with hyphens
	taskSeparator := "-task-"
	taskIdx := strings.LastIndex(taskID, taskSeparator)
	if taskIdx == -1 {
		return "", "", domain.NewValidationError(fmt.Sprintf("Invalid task ID format: %s", taskID))
	}

	featureIDStr := taskID[:taskIdx]

	// Parse feature ID: <project>-feature-<feature>
	featureSeparator := "-feature-"
	featureIdx := strings.Index(featureIDStr, featureSeparator)
	if featureIdx == -1 {
		return "", "", domain.NewValidationError(fmt.Sprintf("Invalid task ID format: %s", taskID))
	}

	projectID = featureIDStr[:featureIdx]
	featureID = featureIDStr
	return projectID, featureID, nil
}

func (s *TaskService) extractProjectIDFromFeatureID(featureID string) (string, error) {
	parts := strings.Split(featureID, "-feature-")
	if len(parts) != 2 {
		return "", domain.NewValidationError(fmt.Sprintf("Invalid feature ID format: %s", featureID))
	}
	return parts[0], nil
}

func (s *TaskService) ValidateCreateInput(input *domain.TaskCreateInput) error {
	if input.FeatureID == "" {
		return domain.NewValidationError("Feature ID is required (--feature).")
	}

	projectID, err := s.extractProjectIDFromFeatureID(input.FeatureID)
	if err != nil {
		return domain.NewValidationError("Invalid feature ID format.")
	}

	if !s.reader.ProjectExists(projectID) {
		return domain.NewValidationError("Project not found: " + projectID)
	}

	feature, err := s.reader.ReadFeature(projectID, input.FeatureID)
	if err != nil {
		return domain.NewValidationError("Feature not found: " + input.FeatureID)
	}

	if feature.Status == domain.FeatureStatusCancelled {
		return domain.NewValidationError("Cannot create task for cancelled feature.")
	}
	if feature.Status == domain.FeatureStatusDone {
		return domain.NewValidationError("Cannot create task for completed feature.")
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.NewValidationError("Task name is required.")
	}

	if strings.TrimSpace(input.Goal) == "" {
		return domain.NewValidationError("Task goal is required (--goal).")
	}

	if !domain.ValidateTaskGoalLength(input.Goal) {
		minLen := domain.TaskGoalMinLength
		if util.IsDevelopment() {
			minLen = domain.TaskGoalMinLengthDevelopment
		}
		return domain.NewValidationError(fmt.Sprintf("Task goal must be at least %d characters. Current length: %d characters.", minLen, len(input.Goal)))
	}

	if len(input.ImplementationSteps) == 0 {
		return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
	}

	if len(input.TestCases) == 0 {
		return domain.NewValidationError("Test cases are required (--test-cases).")
	}

	if len(input.DerivableFiles) == 0 {
		return domain.NewValidationError("Derivable files are required (--derivable-files).")
	}

	if len(input.LibraryNeeds) == 0 {
		return domain.NewValidationError("Library needs are required (--library-needs).")
	}

	if !domain.ValidatePriority(input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Priority == "" {
		input.Priority = "P3"
	}

	if err := s.validateDependencies(projectID, "", input.DependsOn); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) validateDependencies(projectID, selfID string, dependsOn []string) error {
	// Read schema to check if cross-project dependencies are allowed
	schema, err := s.reader.ReadProjectSchema(projectID)
	if err != nil {
		return domain.NewSystemError("Cannot read project schema", err)
	}
	allowCrossProject := schema.Rules.Task.Dependency != "same_project_only" && schema.Rules.Task.Dependency != "disabled"

	for _, depID := range dependsOn {
		if depID == selfID {
			return domain.NewValidationError("Self-dependency detected. Task cannot depend on itself.")
		}

		depProjectID, _, err := s.ParseTaskID(depID)
		if err != nil {
			return domain.NewValidationError("Invalid dependency ID format: " + depID)
		}

		if depProjectID != projectID && !allowCrossProject {
			return domain.NewValidationError(fmt.Sprintf("Cross-project dependency detected: %s -> %s. Cross-project dependencies are disabled.", selfID, depID))
		}

		dep, err := s.reader.ReadTask(depProjectID, depID)
		if err != nil {
			if _, ok := err.(*domain.MandorError); ok {
				return domain.NewValidationError("Dependency not found: " + depID)
			}
			return err
		}

		if dep.Status == domain.TaskStatusCancelled || dep.Status == domain.TaskStatusDone {
			return domain.NewValidationError(fmt.Sprintf("Dependency is not actionable: %s (status: %s)", depID, dep.Status))
		}
	}

	if err := s.validateNoCycle(projectID, selfID, dependsOn); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) validateNoCycle(projectID, selfID string, dependsOn []string) error {
	visited := make(map[string]bool)
	var dfs func(taskID string) bool

	dfs = func(taskID string) bool {
		if taskID == selfID {
			return true
		}
		if visited[taskID] {
			return false
		}
		visited[taskID] = true

		// Extract project ID from task ID for cross-project dependencies
		depProjectID, _, err := s.ParseTaskID(taskID)
		if err != nil {
			return false
		}

		t, err := s.reader.ReadTask(depProjectID, taskID)
		if err != nil {
			return false
		}

		for _, dep := range t.DependsOn {
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

func (s *TaskService) CreateTask(input *domain.TaskCreateInput) (*domain.Task, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	projectID, err := s.extractProjectIDFromFeatureID(input.FeatureID)
	if err != nil {
		return nil, domain.NewValidationError("Invalid feature ID format.")
	}

	nanoid, err := util.GenerateID()
	if err != nil {
		return nil, domain.NewSystemError("Failed to generate task ID", err)
	}

	taskID := input.FeatureID + "-task-" + nanoid

	task := &domain.Task{
		ID:                  taskID,
		FeatureID:           input.FeatureID,
		ProjectID:           projectID,
		Name:                input.Name,
		Goal:                input.Goal,
		Priority:            input.Priority,
		Status:              domain.TaskStatusReady,
		DependsOn:           input.DependsOn,
		ImplementationSteps: input.ImplementationSteps,
		TestCases:           input.TestCases,
		DerivableFiles:      input.DerivableFiles,
		LibraryNeeds:        input.LibraryNeeds,
		CreatedAt:           now,
		UpdatedAt:           now,
		CreatedBy:           creator,
		UpdatedBy:           creator,
	}

	if len(input.DependsOn) > 0 {
		allDone, err := s.checkDependenciesDone(projectID, input.DependsOn)
		if err != nil {
			return nil, err
		}
		if allDone {
			task.Status = domain.TaskStatusReady
		} else {
			task.Status = domain.TaskStatusBlocked
		}
	}

	if err := s.writer.WriteTask(projectID, task); err != nil {
		return nil, err
	}

	event := &domain.TaskEvent{
		Layer: "task",
		Type:  "created",
		ID:    taskID,
		By:    creator,
		Ts:    now,
	}
	if err := s.writer.AppendTaskEvent(projectID, event); err != nil {
		return nil, err
	}

	if task.Status == domain.TaskStatusReady && len(input.DependsOn) == 0 {
		readyEvent := &domain.TaskEvent{
			Layer: "task",
			Type:  "ready",
			ID:    taskID,
			By:    "system",
			Ts:    now,
		}
		if err := s.writer.AppendTaskEvent(projectID, readyEvent); err != nil {
			return nil, err
		}
	}

	if task.Status == domain.TaskStatusBlocked {
		blockedEvent := &domain.TaskEvent{
			Layer: "task",
			Type:  "blocked",
			ID:    taskID,
			By:    "system",
			Ts:    now,
		}
		if err := s.writer.AppendTaskEvent(projectID, blockedEvent); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *TaskService) checkDependenciesDone(projectID string, dependsOn []string) (bool, error) {
	for _, depID := range dependsOn {
		// Extract project ID from task ID for cross-project dependencies
		depProjectID, _, err := s.ParseTaskID(depID)
		if err != nil {
			return false, domain.NewValidationError("Invalid dependency ID format: " + depID)
		}

		dep, err := s.reader.ReadTask(depProjectID, depID)
		if err != nil {
			return false, domain.NewValidationError("Dependency not found: " + depID)
		}
		if dep.Status != domain.TaskStatusDone && dep.Status != domain.TaskStatusCancelled {
			return false, nil
		}
	}
	return true, nil
}

func (s *TaskService) ListTasks(input *domain.TaskListInput) (*domain.TaskListOutput, error) {
	var tasks []domain.TaskListItem
	deletedCount := 0

	projects, err := s.reader.ListProjects(false)
	if err != nil {
		return nil, err
	}

	for _, projectID := range projects {
		if input.ProjectID != "" && projectID != input.ProjectID {
			continue
		}

		err := s.reader.ReadNDJSON(s.paths.ProjectTasksPath(projectID), func(raw []byte) error {
			var t domain.Task
			if err := json.Unmarshal(raw, &t); err != nil {
				return err
			}

			if input.FeatureID != "" && t.FeatureID != input.FeatureID {
				return nil
			}

			if input.Status != "" && t.Status != input.Status {
				return nil
			}

			if input.Priority != "" && t.Priority != input.Priority {
				return nil
			}

			if !input.IncludeDeleted && t.Status == domain.TaskStatusCancelled {
				deletedCount++
				return nil
			}

			item := domain.TaskListItem{
				ID:             t.ID,
				Name:           t.Name,
				Status:         t.Status,
				Priority:       t.Priority,
				FeatureID:      t.FeatureID,
				ProjectID:      t.ProjectID,
				DependsOnCount: len(t.DependsOn),
				CreatedAt:      t.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      t.UpdatedAt.Format(time.RFC3339),
			}
			tasks = append(tasks, item)

			if t.Status == domain.TaskStatusCancelled {
				deletedCount++
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	sortBy := input.Sort
	if sortBy == "" {
		sortBy = "priority"
	}
	orderDesc := input.Order != "asc"

	sort.Slice(tasks, func(i, j int) bool {
		switch sortBy {
		case "priority":
			if orderDesc {
				return comparePriority(tasks[i].Priority, tasks[j].Priority) > 0
			}
			return comparePriority(tasks[i].Priority, tasks[j].Priority) < 0
		case "created_at":
			if orderDesc {
				return tasks[i].CreatedAt > tasks[j].CreatedAt
			}
			return tasks[i].CreatedAt < tasks[j].CreatedAt
		case "name":
			if orderDesc {
				return tasks[i].Name > tasks[j].Name
			}
			return tasks[i].Name < tasks[j].Name
		default:
			return tasks[i].ID < tasks[j].ID
		}
	})

	return &domain.TaskListOutput{
		Tasks:   tasks,
		Total:   len(tasks),
		Deleted: deletedCount,
	}, nil
}

func comparePriority(p1, p2 string) int {
	levels := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	for i, level := range levels {
		if p1 == level {
			return i
		}
	}
	return 3
}

func (s *TaskService) GetTaskDetail(input *domain.TaskDetailInput) (*domain.TaskDetailOutput, error) {
	projectID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return nil, err
	}

	task, err := s.reader.ReadTask(projectID, input.TaskID)
	if err != nil {
		return nil, err
	}

	if !input.IncludeDeleted && task.Status == domain.TaskStatusCancelled {
		return nil, domain.NewValidationError("Task not found: " + input.TaskID)
	}

	events, _ := s.reader.CountEventLines(projectID)

	return &domain.TaskDetailOutput{
		ID:                  task.ID,
		FeatureID:           task.FeatureID,
		ProjectID:           task.ProjectID,
		Name:                task.Name,
		Goal:                task.Goal,
		Priority:            task.Priority,
		Status:              task.Status,
		DependsOn:           task.DependsOn,
		Reason:              task.Reason,
		ImplementationSteps: task.ImplementationSteps,
		TestCases:           task.TestCases,
		DerivableFiles:      task.DerivableFiles,
		LibraryNeeds:        task.LibraryNeeds,
		Events:              events,
		CreatedAt:           task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           task.UpdatedAt.Format(time.RFC3339),
		CreatedBy:           task.CreatedBy,
		UpdatedBy:           task.UpdatedBy,
	}, nil
}

func (s *TaskService) ValidateUpdateInput(input *domain.TaskUpdateInput) error {
	projectID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return err
	}

	task, err := s.reader.ReadTask(projectID, input.TaskID)
	if err != nil {
		return err
	}

	if task.Status == domain.TaskStatusDone {
		return domain.NewValidationError("Cannot modify done task.")
	}

	if task.Status == domain.TaskStatusCancelled && !input.Reopen && !input.Cancel {
		return domain.NewValidationError("Task is cancelled. Use --reopen to reopen, or --cancel to confirm cancellation.")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return domain.NewValidationError("Task name cannot be empty.")
	}

	if input.Priority != nil && !domain.ValidatePriority(*input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if input.Status != nil && !domain.ValidateTaskStatus(*input.Status) {
		return domain.NewValidationError("Invalid status. Valid options: pending, ready, in_progress, blocked, done, cancelled")
	}

	if input.DependsOn != nil {
		if err := s.validateDependencies(projectID, input.TaskID, *input.DependsOn); err != nil {
			return err
		}
	}

	if input.DependsAdd != nil {
		allDeps := append(task.DependsOn, *input.DependsAdd...)
		if err := s.validateDependencies(projectID, input.TaskID, allDeps); err != nil {
			return err
		}
	}

	if input.DependsRemove != nil {
		depSet := make(map[string]bool)
		for _, dep := range task.DependsOn {
			depSet[dep] = true
		}
		for _, remove := range *input.DependsRemove {
			delete(depSet, remove)
		}
		var remaining []string
		for dep := range depSet {
			remaining = append(remaining, dep)
		}
		if err := s.validateDependencies(projectID, input.TaskID, remaining); err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) UpdateTask(input *domain.TaskUpdateInput) ([]string, error) {
	projectID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return nil, err
	}

	task, err := s.reader.ReadTask(projectID, input.TaskID)
	if err != nil {
		return nil, err
	}

	if input.DryRun {
		return []string{"[DRY RUN] Would update task: " + input.TaskID}, nil
	}

	var changes []string
	updater := util.GetGitUsername()
	now := time.Now().UTC()

	if input.Reopen {
		if task.Status != domain.TaskStatusCancelled {
			return nil, domain.NewValidationError("Task is not cancelled. Nothing to reopen.")
		}
		task.Status = domain.TaskStatusPending
		task.Reason = ""
		changes = append(changes, "status", "reason")
	}

	if input.Cancel {
		if task.Status == domain.TaskStatusCancelled {
			return nil, domain.NewValidationError("Task is already cancelled.")
		}

		dependents, err := s.findDependents(projectID, input.TaskID)
		if err != nil {
			return nil, err
		}
		if len(dependents) > 0 && !input.Force {
			return nil, domain.NewValidationError("Task has " + fmt.Sprintf("%d", len(dependents)) + " dependent(s). Use --force to cancel anyway.")
		}

		if input.Reason == nil || *input.Reason == "" {
			return nil, domain.NewValidationError("Cancellation reason is required (--reason).")
		}

		task.Status = domain.TaskStatusCancelled
		task.Reason = *input.Reason
		changes = append(changes, "status", "reason")
	}

	if input.Name != nil && *input.Name != task.Name {
		task.Name = *input.Name
		changes = append(changes, "name")
	}

	if input.Goal != nil && *input.Goal != task.Goal {
		task.Goal = *input.Goal
		changes = append(changes, "goal")
	}

	if input.Priority != nil && *input.Priority != task.Priority {
		task.Priority = *input.Priority
		changes = append(changes, "priority")
	}

	if input.ImplementationSteps != nil {
		task.ImplementationSteps = *input.ImplementationSteps
		changes = append(changes, "implementation_steps")
	}

	if input.TestCases != nil {
		task.TestCases = *input.TestCases
		changes = append(changes, "test_cases")
	}

	if input.DerivableFiles != nil {
		task.DerivableFiles = *input.DerivableFiles
		changes = append(changes, "derivable_files")
	}

	if input.LibraryNeeds != nil {
		task.LibraryNeeds = *input.LibraryNeeds
		changes = append(changes, "library_needs")
	}

	if input.DependsOn != nil {
		task.DependsOn = *input.DependsOn
		changes = append(changes, "depends_on")
	}

	if input.DependsAdd != nil {
		task.DependsOn = append(task.DependsOn, *input.DependsAdd...)
		changes = append(changes, "depends_on")
	}

	if input.DependsRemove != nil {
		depSet := make(map[string]bool)
		for _, dep := range task.DependsOn {
			depSet[dep] = true
		}
		for _, remove := range *input.DependsRemove {
			delete(depSet, remove)
		}
		var remaining []string
		for dep := range depSet {
			remaining = append(remaining, dep)
		}
		task.DependsOn = remaining
		changes = append(changes, "depends_on")
	}

	if input.Status != nil && *input.Status != task.Status {
		if err := s.validateStatusTransition(task.Status, *input.Status); err != nil {
			return nil, err
		}
		task.Status = *input.Status
		changes = append(changes, "status")
	}

	task.UpdatedAt = now
	task.UpdatedBy = updater

	if err := s.writer.ReplaceTask(projectID, task); err != nil {
		return nil, err
	}

	if input.Status != nil && *input.Status == domain.TaskStatusDone {
		unblocked, err := s.unblockDependents(projectID, input.TaskID)
		if err != nil {
			return nil, err
		}
		if unblocked {
			changes = append(changes, "dependent_unblocked")
		}
	}

	event := &domain.TaskEvent{
		Layer:   "task",
		Type:    "updated",
		ID:      input.TaskID,
		By:      updater,
		Ts:      now,
		Changes: changes,
	}
	if err := s.writer.AppendTaskEvent(projectID, event); err != nil {
		return nil, err
	}

	return changes, nil
}

func (s *TaskService) validateStatusTransition(current, next string) error {
	validTransitions := map[string][]string{
		domain.TaskStatusPending:    {domain.TaskStatusReady, domain.TaskStatusInProgress, domain.TaskStatusCancelled},
		domain.TaskStatusReady:      {domain.TaskStatusInProgress, domain.TaskStatusCancelled},
		domain.TaskStatusInProgress: {domain.TaskStatusDone, domain.TaskStatusBlocked, domain.TaskStatusCancelled},
		domain.TaskStatusBlocked:    {domain.TaskStatusReady, domain.TaskStatusCancelled},
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

func (s *TaskService) findDependents(projectID, taskID string) ([]string, error) {
	var dependents []string
	err := s.reader.ReadNDJSON(s.paths.ProjectTasksPath(projectID), func(raw []byte) error {
		var t domain.Task
		if err := json.Unmarshal(raw, &t); err != nil {
			return err
		}
		for _, dep := range t.DependsOn {
			if dep == taskID {
				dependents = append(dependents, t.ID)
			}
		}
		return nil
	})
	return dependents, err
}

func (s *TaskService) unblockDependents(projectID, doneTaskID string) (bool, error) {
	unblockedAny := false
	now := time.Now().UTC()

	// First handle same-project dependencies
	var allTasks []*domain.Task
	err := s.reader.ReadNDJSON(s.paths.ProjectTasksPath(projectID), func(raw []byte) error {
		var task domain.Task
		if err := json.Unmarshal(raw, &task); err != nil {
			return err
		}
		allTasks = append(allTasks, &task)
		return nil
	})
	if err != nil {
		return false, err
	}

	// Track which tasks need to be written back
	tasksToWrite := make(map[string]*domain.Task)
	eventsToAppend := []*domain.TaskEvent{}

	// Process all tasks to find those that should unblock
	for _, task := range allTasks {
		if task.Status != domain.TaskStatusBlocked {
			continue
		}

		hasDone := false
		allDone := true
		for _, depID := range task.DependsOn {
			if depID == doneTaskID {
				hasDone = true
			}
			// Parse the dependency ID to get the project it belongs to
			depProjectID, _, err := s.ParseTaskID(depID)
			if err != nil {
				return false, err
			}
			dep, err := s.reader.ReadTask(depProjectID, depID)
			if err != nil {
				return false, err
			}
			if dep.Status != domain.TaskStatusDone && dep.Status != domain.TaskStatusCancelled {
				allDone = false
			}
		}

		if hasDone && allDone {
			task.Status = domain.TaskStatusReady
			task.UpdatedAt = now
			tasksToWrite[task.ID] = task
			unblockedAny = true

			event := &domain.TaskEvent{
				Layer: "task",
				Type:  "ready",
				ID:    task.ID,
				By:    "system",
				Ts:    now,
			}
			eventsToAppend = append(eventsToAppend, event)
		}
	}

	// If we have same-project updates, write them
	if unblockedAny {
		if err := s.writer.ReplaceTasks(projectID, allTasks, tasksToWrite); err != nil {
			return false, err
		}
		for _, event := range eventsToAppend {
			if err := s.writer.AppendTaskEvent(projectID, event); err != nil {
				return false, err
			}
		}
	}

	// Now handle cross-project dependencies: find all projects and check for tasks that depend on doneTaskID
	projects, err := s.reader.ListProjects(false)
	if err != nil {
		return unblockedAny, err
	}

	for _, otherProjectID := range projects {
		if otherProjectID == projectID {
			continue // Already handled
		}

		var otherProjectTasks []*domain.Task
		err := s.reader.ReadNDJSON(s.paths.ProjectTasksPath(otherProjectID), func(raw []byte) error {
			var task domain.Task
			if err := json.Unmarshal(raw, &task); err != nil {
				return err
			}
			otherProjectTasks = append(otherProjectTasks, &task)
			return nil
		})
		if err != nil {
			continue // Skip if project tasks can't be read
		}

		otherTasksToWrite := make(map[string]*domain.Task)
		otherEventsToAppend := []*domain.TaskEvent{}

		// Process all tasks in other project
		for _, task := range otherProjectTasks {
			if task.Status != domain.TaskStatusBlocked {
				continue
			}

			hasDone := false
			allDone := true
			for _, depID := range task.DependsOn {
				if depID == doneTaskID {
					hasDone = true
				}
				depProjectID, _, err := s.ParseTaskID(depID)
				if err != nil {
					continue
				}
				dep, err := s.reader.ReadTask(depProjectID, depID)
				if err != nil {
					allDone = false
					continue
				}
				if dep.Status != domain.TaskStatusDone && dep.Status != domain.TaskStatusCancelled {
					allDone = false
				}
			}

			if hasDone && allDone {
				task.Status = domain.TaskStatusReady
				task.UpdatedAt = now
				otherTasksToWrite[task.ID] = task
				unblockedAny = true

				event := &domain.TaskEvent{
					Layer: "task",
					Type:  "ready",
					ID:    task.ID,
					By:    "system",
					Ts:    now,
				}
				otherEventsToAppend = append(otherEventsToAppend, event)
			}
		}

		// Write updates for other project if any
		if len(otherTasksToWrite) > 0 {
			if err := s.writer.ReplaceTasks(otherProjectID, otherProjectTasks, otherTasksToWrite); err != nil {
				continue // Skip error for other project
			}
			for _, event := range otherEventsToAppend {
				if err := s.writer.AppendTaskEvent(otherProjectID, event); err != nil {
					continue
				}
			}
		}
	}

	return unblockedAny, nil
}

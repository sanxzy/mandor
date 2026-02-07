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

func (s *TaskService) ParseTaskID(taskID string) (backlogID, featureID string, err error) {
	// Task ID format: <backlog>-feature-<feature>-task-<nanoid>
	// Find the last occurrence of "-task-" to handle backlog IDs with hyphens
	taskSeparator := "-task-"
	taskIdx := strings.LastIndex(taskID, taskSeparator)
	if taskIdx == -1 {
		return "", "", domain.NewValidationError(fmt.Sprintf("Invalid task ID format: %s", taskID))
	}

	featureIDStr := taskID[:taskIdx]

	// Parse feature ID: <backlog>-feature-<feature>
	featureSeparator := "-feature-"
	featureIdx := strings.Index(featureIDStr, featureSeparator)
	if featureIdx == -1 {
		return "", "", domain.NewValidationError(fmt.Sprintf("Invalid task ID format: %s", taskID))
	}

	backlogID = featureIDStr[:featureIdx]
	featureID = featureIDStr
	return backlogID, featureID, nil
}

func (s *TaskService) extractBacklogIDFromFeatureID(featureID string) (string, error) {
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

	backlogID, err := s.extractBacklogIDFromFeatureID(input.FeatureID)
	if err != nil {
		return domain.NewValidationError("Invalid feature ID format.")
	}

	if !s.reader.BacklogExists(backlogID) {
		return domain.NewValidationError("Backlog not found: " + backlogID)
	}

	feature, err := s.reader.ReadFeature(backlogID, input.FeatureID)
	if err != nil {
		return domain.NewValidationError("Feature not found: " + input.FeatureID)
	}

	if feature.Status == domain.FeatureStatusCancelled {
		return domain.NewValidationError("Cannot create task for cancelled feature.")
	}
	if feature.Status == domain.FeatureStatusDone {
		return domain.NewValidationError("Cannot create task for completed feature.")
	}

	// Validate SpecID matches Feature's SpecID
	if input.SpecID == "" {
		return domain.NewValidationError("Spec ID is required (--spec-id).")
	}
	if input.SpecID != feature.SpecID {
		return domain.NewValidationError(fmt.Sprintf("Spec ID mismatch. Task Spec ID (%s) must match Feature's Spec ID (%s).", input.SpecID, feature.SpecID))
	}

	// Validate IAE scenarios exist in the Spec
	if len(input.IAEScenarios) == 0 {
		return domain.NewValidationError("IAE scenarios are required (--iae-scenarios).")
	}

	if err := s.validateIAEScenariosExist(backlogID, input.SpecID, input.IAEScenarios); err != nil {
		return err
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.NewValidationError("Task name is required.")
	}

	if err := s.validateNoDuplicateName(backlogID, input.FeatureID, input.Name); err != nil {
		return err
	}

	if strings.TrimSpace(input.Goal) == "" {
		return domain.NewValidationError("Task goal is required (--goal).")
	}

	minLen := s.getTaskGoalMinLength()
	if !domain.ValidateTaskGoalLength(input.Goal, minLen) {
		return domain.NewValidationError(fmt.Sprintf("Task goal must be at least %d characters. Current length: %d characters.", minLen, len(input.Goal)))
	}

	if len(input.ImplementationSteps) == 0 {
		return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
	}

	if len(input.TestCases) == 0 {
		return domain.NewValidationError("Test cases are required (--test-cases).")
	}

	// Apply default priority if not specified
	if input.Priority == "" {
		// Use default priority from workspace config
		ws, err := s.reader.ReadWorkspace()
		if err == nil && ws.Config.DefaultPriority != "" {
			input.Priority = ws.Config.DefaultPriority
		} else {
			input.Priority = "P3"
		}
	}

	if !domain.ValidatePriority(input.Priority) {
		return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
	}

	if err := s.validateDependencies(backlogID, "", input.DependsOn); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) validateDependencies(backlogID, selfID string, dependsOn []string) error {
	// Read schema to check if cross-backlog dependencies are allowed
	schema, err := s.reader.ReadBacklogSchema(backlogID)
	if err != nil {
		return domain.NewSystemError("Cannot read backlog schema", err)
	}
	allowCrossBacklog := schema.Rules.Task.Dependency != "same_backlog_only" && schema.Rules.Task.Dependency != "disabled"

	for _, depID := range dependsOn {
		if depID == selfID {
			return domain.NewValidationError("Self-dependency detected. Task cannot depend on itself.")
		}

		depBacklogID, _, err := s.ParseTaskID(depID)
		if err != nil {
			return domain.NewValidationError("Invalid dependency ID format: " + depID)
		}

		if depBacklogID != backlogID && !allowCrossBacklog {
			return domain.NewValidationError(fmt.Sprintf("Cross-backlog dependency detected: %s -> %s. Cross-backlog dependencies are disabled.", selfID, depID))
		}

		dep, err := s.reader.ReadTask(depBacklogID, depID)
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

	if err := s.validateNoCycle(backlogID, selfID, dependsOn); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) validateNoCycle(backlogID, selfID string, dependsOn []string) error {
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

		// Extract backlog ID from task ID for cross-backlog dependencies
		depBacklogID, _, err := s.ParseTaskID(taskID)
		if err != nil {
			return false
		}

		t, err := s.reader.ReadTask(depBacklogID, taskID)
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

func (s *TaskService) validateNoDuplicateName(backlogID, featureID, name string) error {
	var tasks []domain.Task
	err := s.reader.ReadNDJSON(s.paths.BacklogTasksPath(backlogID), func(raw []byte) error {
		var t domain.Task
		if err := json.Unmarshal(raw, &t); err != nil {
			return err
		}
		tasks = append(tasks, t)
		return nil
	})
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if t.FeatureID == featureID && t.Name == name && t.Status != domain.TaskStatusCancelled {
			return domain.NewValidationError("A task with this name already exists in the feature: " + name)
		}
	}
	return nil
}

func (s *TaskService) validateIAEScenariosExist(backlogID, specID string, iaeScenarios []string) error {
	// Load the spec to validate scenarios exist
	specService := NewSpecServiceWithPaths(s.paths)
	spec, err := specService.ReadSpec(backlogID, specID)
	if err != nil {
		return domain.NewValidationError(fmt.Sprintf("Spec not found: %s", specID))
	}

	// Build a map of available requirement:scenario combinations
	validScenarios := make(map[string]bool)
	for _, req := range spec.Requirements {
		for _, scenario := range req.IAEScenarios {
			key := fmt.Sprintf("%s:%s", req.ID, scenario.ID)
			validScenarios[key] = true
		}
	}

	// Validate all referenced scenarios exist
	for _, iae := range iaeScenarios {
		if !validScenarios[iae] {
			return domain.NewValidationError(fmt.Sprintf("IAE scenario not found in spec: %s. Valid format: req-XXXX:scenario-YYYY", iae))
		}
	}

	return nil
}

func (s *TaskService) CreateTask(input *domain.TaskCreateInput) (*domain.Task, error) {
	creator := util.GetGitUsername()
	now := time.Now().UTC()

	backlogID, err := s.extractBacklogIDFromFeatureID(input.FeatureID)
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
		SpecID:              input.SpecID,
		BacklogID:           backlogID,
		Name:                input.Name,
		Goal:                input.Goal,
		Priority:            input.Priority,
		Status:              domain.TaskStatusReady,
		DependsOn:           input.DependsOn,
		IAEScenarios:        input.IAEScenarios,
		ImplementationSteps: input.ImplementationSteps,
		TestCases:           input.TestCases,
		LibraryNeeds:        input.LibraryNeeds,
		ReadGates: domain.ReadGates{
			IsReadBrief:        false,
			IsReadSpec:         false,
			IsReadSessionNotes: false,
		},
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: creator,
		UpdatedBy: creator,
	}

	if len(input.DependsOn) > 0 {
		allDone, err := s.checkDependenciesDone(backlogID, input.DependsOn)
		if err != nil {
			return nil, err
		}
		if allDone {
			task.Status = domain.TaskStatusReady
		} else {
			task.Status = domain.TaskStatusBlocked
		}
	}

	if err := s.writer.WriteTask(backlogID, task); err != nil {
		return nil, err
	}

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.TaskEvent{
	// 	Layer: "task",
	// 	Type:  "created",
	// 	ID:    taskID,
	// 	By:    creator,
	// 	Ts:    now,
	// }
	// 	if err := s.writer.AppendTaskEvent(backlogID, event); err != nil {
	// 		return nil, err
	// 	}
	// }

	// if task.Status == domain.TaskStatusReady && len(input.DependsOn) == 0 {
	// 	readyEvent := &domain.TaskEvent{
	// 		Layer: "task",
	// 		Type:  "ready",
	// 		ID:    taskID,
	// 		By:    "system",
	// 		Ts:    now,
	// 	}
	// 	if err := s.writer.AppendTaskEvent(backlogID, readyEvent); err != nil {
	// 		return nil, err
	// 	}
	// }

	// if task.Status == domain.TaskStatusBlocked {
	// 	blockedEvent := &domain.TaskEvent{
	// 		Layer: "task",
	// 		Type:  "blocked",
	// 		ID:    taskID,
	// 		By:    "system",
	// 		Ts:    now,
	// 	}
	// 	if err := s.writer.AppendTaskEvent(backlogID, blockedEvent); err != nil {
	// 		return nil, err
	// 	}
	// }

	return task, nil
}

func (s *TaskService) checkDependenciesDone(backlogID string, dependsOn []string) (bool, error) {
	for _, depID := range dependsOn {
		// Extract backlog ID from task ID for cross-backlog dependencies
		depBacklogID, _, err := s.ParseTaskID(depID)
		if err != nil {
			return false, domain.NewValidationError("Invalid dependency ID format: " + depID)
		}

		dep, err := s.reader.ReadTask(depBacklogID, depID)
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

	backlogs, err := s.reader.ListBacklogs(false)
	if err != nil {
		return nil, err
	}

	for _, backlogID := range backlogs {
		if input.BacklogID != "" && backlogID != input.BacklogID {
			continue
		}

		err := s.reader.ReadNDJSON(s.paths.BacklogTasksPath(backlogID), func(raw []byte) error {
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
				BacklogID:      t.BacklogID,
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
				return ComparePriority(tasks[i].Priority, tasks[j].Priority) > 0
			}
			return ComparePriority(tasks[i].Priority, tasks[j].Priority) < 0
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

// ComparePriority compares two priority levels and returns negative if p1 > p2, positive if p1 < p2, 0 if equal
func ComparePriority(p1, p2 string) int {
	levels := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	for i, level := range levels {
		if p1 == level {
			p1Index := i
			for j, l := range levels {
				if p2 == l {
					return p1Index - j
				}
			}
			return -1
		}
	}
	return 1
}

func (s *TaskService) GetTaskDetail(input *domain.TaskDetailInput) (*domain.TaskDetailOutput, error) {
	backlogID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return nil, err
	}

	task, err := s.reader.ReadTask(backlogID, input.TaskID)
	if err != nil {
		return nil, err
	}

	if !input.IncludeDeleted && task.Status == domain.TaskStatusCancelled {
		return nil, domain.NewValidationError("Task not found: " + input.TaskID)
	}

	return &domain.TaskDetailOutput{
		ID:                  task.ID,
		FeatureID:           task.FeatureID,
		BacklogID:           task.BacklogID,
		Name:                task.Name,
		Goal:                task.Goal,
		Priority:            task.Priority,
		Status:              task.Status,
		DependsOn:           task.DependsOn,
		Reason:              task.Reason,
		ImplementationSteps: task.ImplementationSteps,
		TestCases:           task.TestCases,
		LibraryNeeds:        task.LibraryNeeds,
		Events:              0,
		CreatedAt:           task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           task.UpdatedAt.Format(time.RFC3339),
		CreatedBy:           task.CreatedBy,
		UpdatedBy:           task.UpdatedBy,
	}, nil
}

func (s *TaskService) ValidateUpdateInput(input *domain.TaskUpdateInput) error {
	backlogID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return err
	}

	task, err := s.reader.ReadTask(backlogID, input.TaskID)
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
		if err := s.validateDependencies(backlogID, input.TaskID, *input.DependsOn); err != nil {
			return err
		}
	}

	if input.DependsAdd != nil {
		allDeps := append(task.DependsOn, *input.DependsAdd...)
		if err := s.validateDependencies(backlogID, input.TaskID, allDeps); err != nil {
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
		if err := s.validateDependencies(backlogID, input.TaskID, remaining); err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) UpdateTask(input *domain.TaskUpdateInput) ([]string, error) {
	backlogID, _, err := s.ParseTaskID(input.TaskID)
	if err != nil {
		return nil, err
	}

	task, err := s.reader.ReadTask(backlogID, input.TaskID)
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
		task.Status = domain.TaskStatusReady
		task.Reason = ""
		changes = append(changes, "status", "reason")
	}

	if input.Cancel {
		if task.Status == domain.TaskStatusCancelled {
			return nil, domain.NewValidationError("Task is already cancelled.")
		}

		dependents, err := s.findDependents(backlogID, input.TaskID)
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
		minLen := s.getTaskGoalMinLength()
		if !domain.ValidateTaskGoalLength(*input.Goal, minLen) {
			return nil, domain.NewValidationError(fmt.Sprintf("Task goal must be at least %d characters. Current length: %d characters.", minLen, len(*input.Goal)))
		}
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
		// Enforce gates before transitioning to in_progress
		if *input.Status == domain.TaskStatusInProgress {
			if err := s.checkGatesBeforeInProgress(task.ReadGates); err != nil {
				return nil, err
			}
		}
		task.Status = *input.Status
		changes = append(changes, "status")
	}

	task.UpdatedAt = now
	task.UpdatedBy = updater

	if err := s.writer.ReplaceTask(backlogID, task); err != nil {
		return nil, err
	}

	if input.Status != nil && *input.Status == domain.TaskStatusDone {
		unblocked, err := s.unblockDependents(backlogID, input.TaskID)
		if err != nil {
			return nil, err
		}
		if unblocked {
			changes = append(changes, "dependent_unblocked")
		}
	}

	// TODO: Event system being removed - Phase 1 commented out
	// event := &domain.TaskEvent{
	// 	Layer:   "task",
	// 	Type:    "updated",
	// 	ID:      input.TaskID,
	// 	By:      updater,
	// 	Ts:      now,
	// 	Changes: changes,
	// }
	// if err := s.writer.AppendTaskEvent(backlogID, event); err != nil {
	// 	return nil, err
	// }

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

func (s *TaskService) checkGatesBeforeInProgress(gates domain.ReadGates) error {
	unmetGates := []string{}

	if !gates.IsReadBrief {
		unmetGates = append(unmetGates, "is-read-brief")
	}
	if !gates.IsReadSpec {
		unmetGates = append(unmetGates, "is-read-spec")
	}
	if !gates.IsReadSessionNotes {
		unmetGates = append(unmetGates, "is-read-session-notes")
	}

	if len(unmetGates) > 0 {
		solution := "Set all gates before transitioning to in_progress:\n"
		for _, gate := range unmetGates {
			solution += fmt.Sprintf("  mandor task set-gate <task-id> --%s\n", gate)
		}
		return domain.NewValidationError(fmt.Sprintf("Error: Cannot transition to in_progress - %d unmet gates: %s\nSolution: %s",
			len(unmetGates),
			strings.Join(unmetGates, ", "),
			strings.TrimSuffix(solution, "\n")))
	}

	return nil
}

func (s *TaskService) findDependents(backlogID, taskID string) ([]string, error) {
	var dependents []string
	err := s.reader.ReadNDJSON(s.paths.BacklogTasksPath(backlogID), func(raw []byte) error {
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

func (s *TaskService) unblockDependents(backlogID, doneTaskID string) (bool, error) {
	unblockedAny := false
	now := time.Now().UTC()

	// First handle same-backlog dependencies
	var allTasks []*domain.Task
	err := s.reader.ReadNDJSON(s.paths.BacklogTasksPath(backlogID), func(raw []byte) error {
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
			// Parse the dependency ID to get the backlog it belongs to
			depBacklogID, _, err := s.ParseTaskID(depID)
			if err != nil {
				return false, err
			}
			dep, err := s.reader.ReadTask(depBacklogID, depID)
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
		}
	}

	// If we have same-backlog updates, write them
	if unblockedAny {
		if err := s.writer.ReplaceTasks(backlogID, allTasks, tasksToWrite); err != nil {
			return false, err
		}
	}

	// Now handle cross-backlog dependencies: find all backlogs and check for tasks that depend on doneTaskID
	backlogs, err := s.reader.ListBacklogs(false)
	if err != nil {
		return unblockedAny, err
	}

	for _, otherBacklogID := range backlogs {
		if otherBacklogID == backlogID {
			continue // Already handled
		}

		var otherBacklogTasks []*domain.Task
		err := s.reader.ReadNDJSON(s.paths.BacklogTasksPath(otherBacklogID), func(raw []byte) error {
			var task domain.Task
			if err := json.Unmarshal(raw, &task); err != nil {
				return err
			}
			otherBacklogTasks = append(otherBacklogTasks, &task)
			return nil
		})
		if err != nil {
			continue // Skip if backlog tasks can't be read
		}

		otherTasksToWrite := make(map[string]*domain.Task)

		// Process all tasks in other backlog
		for _, task := range otherBacklogTasks {
			if task.Status != domain.TaskStatusBlocked {
				continue
			}

			hasDone := false
			allDone := true
			for _, depID := range task.DependsOn {
				if depID == doneTaskID {
					hasDone = true
				}
				depBacklogID, _, err := s.ParseTaskID(depID)
				if err != nil {
					continue
				}
				dep, err := s.reader.ReadTask(depBacklogID, depID)
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
			}
		}

		// Write updates for other backlog if any
		if len(otherTasksToWrite) > 0 {
			if err := s.writer.ReplaceTasks(otherBacklogID, otherBacklogTasks, otherTasksToWrite); err != nil {
				continue // Skip error for other backlog
			}
		}
	}

	return unblockedAny, nil
}

func (s *TaskService) getTaskGoalMinLength() int {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return domain.TaskGoalMinLength
	}
	if ws.Config.GoalLengths.Task > 0 {
		return ws.Config.GoalLengths.Task
	}
	return domain.TaskGoalMinLength
}

// ReadTask reads a task by ID from filesystem
func (s *TaskService) ReadTask(backlogID, featureID, taskID string) (*domain.Task, error) {
	return s.reader.ReadTask(backlogID, taskID)
}

// SaveTask saves a task object directly to filesystem (used by gates commands)
func (s *TaskService) SaveTask(backlogID, featureID string, task *domain.Task) error {
	task.UpdatedAt = time.Now().UTC()
	task.UpdatedBy = util.GetGitUsername()
	return s.writer.ReplaceTask(backlogID, task)
}

package track

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

// handleWorkspace handles workspace-level tracking
func handleWorkspace(cmd *cobra.Command, _ string) error {
	backlogSvc, err := service.NewBacklogService()
	if err != nil {
		return err
	}

	if !backlogSvc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Get all backlogs
	output, err := backlogSvc.ListBacklogs(false, true)
	if err != nil {
		return err
	}

	response := &TrackResponse{
		Scope:    "workspace",
		Backlogs: []BacklogTrackItem{},
		Summary:  SummaryStats{ByStatus: make(map[string]int)},
	}

	// Build backlog items
	for _, backlog := range output.Backlogs {
		backlogItem := BacklogTrackItem{
			ID:        backlog.ID,
			Name:      backlog.Name,
			Status:    backlog.Status,
			CreatedAt: backlog.CreatedAt,
			UpdatedAt: backlog.UpdatedAt,
			Features:  backlog.Features,
			Tasks:     backlog.Tasks,
			Issues:    backlog.Issues,
		}

		if globalFlags.Verbose {
			backlogItem.Description = backlog.Goal
		}

		response.Backlogs = append(response.Backlogs, backlogItem)
		response.Summary.Total++
	}

	return outputResponse(cmd, response)
}

// handleBacklog handles backlog-level tracking
func handleBacklog(cmd *cobra.Command, backlogID string) error {
	backlogSvc, err := service.NewBacklogService()
	if err != nil {
		return err
	}

	if !backlogSvc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Get backlog
	backlog, err := backlogSvc.GetBacklog(backlogID)
	if err != nil {
		return err
	}

	// Get issues
	issueSvc, err := service.NewIssueService()
	if err != nil {
		return err
	}

	input := &domain.IssueListInput{
		BacklogID: backlogID,
	}

	issuesOutput, err := issueSvc.ListIssues(input)
	if err != nil {
		return err
	}

	response := &TrackResponse{
		Scope:   "backlog",
		ID:      backlog.ID,
		Name:    backlog.Name,
		Issues:  []IssueTrackItem{},
		Summary: SummaryStats{ByStatus: make(map[string]int)},
	}

	// Build issue items
	statusMap := make(map[string]int)

	for _, issue := range issuesOutput.Issues {
		issueItem := IssueTrackItem{
			ID:        issue.ID,
			Title:     issue.Name,
			Type:      issue.IssueType,
			Status:    issue.Status,
			Priority:  issue.Priority,
			BacklogID: issue.BacklogID,
			CreatedAt: issue.CreatedAt,
			UpdatedAt: issue.LastUpdatedAt,
		}

		response.Issues = append(response.Issues, issueItem)
		statusMap[issue.Status]++
		response.Summary.Total++
	}

	response.Summary.ByStatus = statusMap

	// Build recommendations
	if response.Summary.Total > 0 {
		response.RecommendedNextCommands = []string{
			fmt.Sprintf("mandor track issue <issue_id>    # View specific issue details"),
			fmt.Sprintf("mandor track backlog %s --verbose  # See blockers and relationships", backlogID),
		}
	}

	return outputResponse(cmd, response)
}

// handleFeature handles feature-level tracking
func handleFeature(cmd *cobra.Command, featureID string) error {
	featureSvc, err := service.NewFeatureService()
	if err != nil {
		return err
	}

	if !featureSvc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Parse feature ID to get backlog
	parts := strings.Split(featureID, "-feature-")
	if len(parts) != 2 {
		return domain.NewValidationError(fmt.Sprintf("Invalid feature ID format: %s", featureID))
	}
	backlogID := parts[0]

	// Get feature detail
	input := &domain.FeatureDetailInput{
		BacklogID: backlogID,
		FeatureID: featureID,
	}

	feature, err := featureSvc.GetFeatureDetail(input)
	if err != nil {
		return err
	}

	// Get tasks
	taskSvc, err := service.NewTaskService()
	if err != nil {
		return err
	}

	taskListInput := &domain.TaskListInput{
		FeatureID: featureID,
		BacklogID: backlogID,
	}

	tasksOutput, err := taskSvc.ListTasks(taskListInput)
	if err != nil {
		return err
	}

	response := &TrackResponse{
		Scope:   "feature",
		ID:      feature.ID,
		Name:    feature.Name,
		Tasks:   []TaskTrackItem{},
		Summary: SummaryStats{ByStatus: make(map[string]int)},
	}

	// Build task items
	statusMap := make(map[string]int)
	completedCount := 0

	for _, task := range tasksOutput.Tasks {
		taskItem := TaskTrackItem{
			ID:        task.ID,
			Name:      task.Name,
			Status:    task.Status,
			Priority:  task.Priority,
			FeatureID: task.FeatureID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}

		response.Tasks = append(response.Tasks, taskItem)
		statusMap[task.Status]++
		response.Summary.Total++

		if task.Status == domain.TaskStatusDone {
			completedCount++
		}
	}

	response.Summary.ByStatus = statusMap
	if response.Summary.Total > 0 {
		response.Summary.CompletionPercent = (completedCount * 100) / response.Summary.Total
	}

	// Build recommendations
	if response.Summary.Total > 0 {
		response.RecommendedNextCommands = []string{
			fmt.Sprintf("mandor track task <task_id>      # View specific task details"),
			fmt.Sprintf("mandor track feature %s --verbose # See task blockers", featureID),
			fmt.Sprintf("mandor track backlog %s           # View parent backlog", backlogID),
		}
	}

	return outputResponse(cmd, response)
}

// handleTask handles task-level tracking
func handleTask(cmd *cobra.Command, taskID string) error {
	taskSvc, err := service.NewTaskService()
	if err != nil {
		return err
	}

	if !taskSvc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Parse task ID
	backlogID, featureID, err := taskSvc.ParseTaskID(taskID)
	if err != nil {
		return err
	}

	// Get task detail
	input := &domain.TaskDetailInput{
		FeatureID: featureID,
		TaskID:    taskID,
	}

	task, err := taskSvc.GetTaskDetail(input)
	if err != nil {
		return err
	}

	response := &TrackResponse{
		Scope: "task",
		ID:    task.ID,
		Name:  task.Name,
		Tasks: []TaskTrackItem{
			{
				ID:                  task.ID,
				Name:                task.Name,
				Status:              task.Status,
				Priority:            task.Priority,
				FeatureID:           task.FeatureID,
				CreatedAt:           task.CreatedAt,
				UpdatedAt:           task.UpdatedAt,
				Goal:                task.Goal,
				ImplementationSteps: task.ImplementationSteps,
				TestCases:           task.TestCases,
				BlockedBy:           task.DependsOn,
			},
		},
		Summary: SummaryStats{
			Total:    1,
			ByStatus: map[string]int{task.Status: 1},
		},
	}

	// Build recommendations
	response.RecommendedNextCommands = []string{
		fmt.Sprintf("mandor track feature %s          # View all tasks in feature", featureID),
		fmt.Sprintf("mandor track backlog %s          # View parent backlog", backlogID),
	}

	if len(task.DependsOn) > 0 {
		response.RecommendedNextCommands = append(response.RecommendedNextCommands,
			fmt.Sprintf("mandor track task %s             # Check blocking task", task.DependsOn[0]))
	}

	return outputResponse(cmd, response)
}

// handleIssue handles issue-level tracking
func handleIssue(cmd *cobra.Command, issueID string) error {
	issueSvc, err := service.NewIssueService()
	if err != nil {
		return err
	}

	if !issueSvc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Parse issue ID to get backlog
	parts := strings.Split(issueID, "-issue-")
	if len(parts) != 2 {
		return domain.NewValidationError(fmt.Sprintf("Invalid issue ID format: %s", issueID))
	}
	backlogID := parts[0]

	// Get issue detail
	input := &domain.IssueDetailInput{
		BacklogID: backlogID,
		IssueID:   issueID,
	}

	issue, err := issueSvc.GetIssueDetail(input)
	if err != nil {
		return err
	}

	response := &TrackResponse{
		Scope: "issue",
		ID:    issue.ID,
		Name:  issue.Name,
		Issues: []IssueTrackItem{
			{
				ID:          issue.ID,
				Title:       issue.Name,
				Type:        issue.IssueType,
				Status:      issue.Status,
				Priority:    issue.Priority,
				BacklogID:   issue.BacklogID,
				CreatedAt:   issue.CreatedAt,
				UpdatedAt:   issue.LastUpdatedAt,
				Description: issue.Goal,
				BlockedBy:   issue.DependsOn,
			},
		},
		Summary: SummaryStats{
			Total:    1,
			ByStatus: map[string]int{issue.Status: 1},
		},
	}

	// Build recommendations
	response.RecommendedNextCommands = []string{
		fmt.Sprintf("mandor track backlog %s           # View parent backlog", backlogID),
		fmt.Sprintf("mandor track issue %s --verbose  # See all relationships", issueID),
	}

	if len(issue.DependsOn) > 0 {
		response.RecommendedNextCommands = append(response.RecommendedNextCommands,
			fmt.Sprintf("mandor track issue %s            # Check blocking issue", issue.DependsOn[0]))
	}

	return outputResponse(cmd, response)
}

// handleAutoScope attempts to auto-resolve scope from ID
func handleAutoScope(cmd *cobra.Command, id string) error {
	// Check ID prefix conventions
	if strings.HasPrefix(id, "task-") || strings.HasPrefix(id, "t-") {
		return handleTask(cmd, id)
	}
	if strings.HasPrefix(id, "issue-") || strings.HasPrefix(id, "i-") {
		return handleIssue(cmd, id)
	}
	if strings.HasPrefix(id, "feature-") || strings.HasPrefix(id, "f-") {
		return handleFeature(cmd, id)
	}
	if strings.HasPrefix(id, "backlog-") || strings.HasPrefix(id, "b-") {
		return handleBacklog(cmd, id)
	}

	// Try to resolve from data store
	taskSvc, _ := service.NewTaskService()
	if taskSvc != nil && strings.Contains(id, "-task-") {
		_, _, err := taskSvc.ParseTaskID(id)
		if err == nil {
			return handleTask(cmd, id)
		}
	}

	issueSvc, _ := service.NewIssueService()
	if issueSvc != nil && strings.Contains(id, "-issue-") {
		return handleIssue(cmd, id)
	}

	featureSvc, _ := service.NewFeatureService()
	if featureSvc != nil && strings.Contains(id, "-feature-") {
		return handleFeature(cmd, id)
	}

	backlogSvc, _ := service.NewBacklogService()
	if backlogSvc != nil {
		_, err := backlogSvc.GetBacklog(id)
		if err == nil {
			return handleBacklog(cmd, id)
		}
	}

	return domain.NewValidationError(fmt.Sprintf("ID '%s' not found. Checked: tasks, issues, features, backlogs", id))
}

package track

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// outputResponse outputs the response in the requested format
func outputResponse(cmd *cobra.Command, response *TrackResponse) error {
	// Validate flags
	if err := response.ParseFlags(); err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	// Handle different output formats
	if globalFlags.JSON {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(response)
	}

	if globalFlags.CSV {
		return outputCSV(out, response)
	}

	if globalFlags.Tree {
		return outputTree(out, response)
	}

	if globalFlags.Graph {
		return outputGraph(out, response)
	}

	// Default: table output
	return outputTable(out, response)
}

// outputTable outputs human-readable table format
func outputTable(out io.Writer, response *TrackResponse) error {
	switch response.Scope {
	case "workspace":
		return outputWorkspaceTable(out, response)
	case "backlog":
		return outputBacklogTable(out, response)
	case "feature":
		return outputFeatureTable(out, response)
	case "task":
		return outputTaskTable(out, response)
	case "issue":
		return outputIssueTable(out, response)
	}

	return nil
}

func outputWorkspaceTable(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Workspace Overview\n")
	fmt.Fprintf(out, "%-20s %-30s %-10s %10s %10s %10s\n", "Backlog ID", "Backlog", "Status", "Features", "Tasks", "Issues")
	fmt.Fprintf(out, strings.Repeat("-", 100)+"\n")

	for _, backlog := range response.Backlogs {
		fmt.Fprintf(out, "%-20s %-30s %-10s %10d %10d %10d\n",
			truncate(backlog.ID, 20), truncate(backlog.Name, 30), backlog.Status, backlog.Features, backlog.Tasks, backlog.Issues)
	}

	fmt.Fprintf(out, "\n")
	outputSummaryStats(out, response.Summary)
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputBacklogTable(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Backlog: %s\n", response.Name)
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "%-40s %-15s %-8s %s\n", "Issue ID", "Type", "Status", "Title")
	fmt.Fprintf(out, strings.Repeat("-", 80)+"\n")

	issues := response.Issues
	if globalFlags.GroupBy == "status" {
		issues = groupIssuesByStatus(issues)
	} else if globalFlags.GroupBy == "priority" {
		issues = groupIssuesByPriority(issues)
	}

	for _, issue := range issues {
		fmt.Fprintf(out, "%-40s %-15s %-8s %s\n",
			truncate(issue.ID, 40), issue.Type, issue.Status, truncate(issue.Title, 30))
	}

	fmt.Fprintf(out, "\n")
	outputSummaryStats(out, response.Summary)
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputFeatureTable(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Feature: %s\n", response.Name)
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "%-44s %-10s %-8s %s\n", "Task ID", "Status", "Priority", "Name")
	fmt.Fprintf(out, strings.Repeat("-", 80)+"\n")

	tasks := response.Tasks
	if globalFlags.GroupBy == "status" {
		tasks = groupTasksByStatus(tasks)
	} else if globalFlags.GroupBy == "priority" {
		tasks = groupTasksByPriority(tasks)
	}

	for _, task := range tasks {
		fmt.Fprintf(out, "%-44s %-10s %-8s %s\n",
			truncate(task.ID, 44), task.Status, task.Priority, truncate(task.Name, 20))
	}

	fmt.Fprintf(out, "\n")
	outputSummaryStats(out, response.Summary)
	if response.Summary.CompletionPercent > 0 {
		fmt.Fprintf(out, "Completion: %d%%\n", response.Summary.CompletionPercent)
	}
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputTaskTable(out io.Writer, response *TrackResponse) error {
	if len(response.Tasks) == 0 {
		return nil
	}

	task := response.Tasks[0]
	fmt.Fprintf(out, "Task: %s\n", task.Name)
	fmt.Fprintf(out, "ID:       %s\n", task.ID)
	fmt.Fprintf(out, "Status:   %s\n", task.Status)
	fmt.Fprintf(out, "Priority: %s\n", task.Priority)
	fmt.Fprintf(out, "Feature:  %s\n", task.FeatureID)

	if globalFlags.Verbose {
		fmt.Fprintf(out, "Goal:     %s\n", task.Goal)
		if len(task.ImplementationSteps) > 0 {
			fmt.Fprintf(out, "Steps:    %v\n", task.ImplementationSteps)
		}
		if len(task.BlockedBy) > 0 {
			fmt.Fprintf(out, "BlockedBy: %v\n", task.BlockedBy)
		}
		if len(task.Blocks) > 0 {
			fmt.Fprintf(out, "Blocks:   %v\n", task.Blocks)
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputIssueTable(out io.Writer, response *TrackResponse) error {
	if len(response.Issues) == 0 {
		return nil
	}

	issue := response.Issues[0]
	fmt.Fprintf(out, "Issue: %s\n", issue.Title)
	fmt.Fprintf(out, "ID:       %s\n", issue.ID)
	fmt.Fprintf(out, "Type:     %s\n", issue.Type)
	fmt.Fprintf(out, "Status:   %s\n", issue.Status)
	fmt.Fprintf(out, "Priority: %s\n", issue.Priority)
	fmt.Fprintf(out, "Backlog: %s\n", issue.BacklogID)

	if globalFlags.Verbose {
		fmt.Fprintf(out, "Description: %s\n", issue.Description)
		if len(issue.BlockedBy) > 0 {
			fmt.Fprintf(out, "BlockedBy: %v\n", issue.BlockedBy)
		}
		if len(issue.Blocks) > 0 {
			fmt.Fprintf(out, "Blocks:   %v\n", issue.Blocks)
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputSummaryStats(out io.Writer, stats SummaryStats) {
	fmt.Fprintf(out, "Total: %d\n", stats.Total)
	if len(stats.ByStatus) > 0 {
		fmt.Fprintf(out, "By Status: ")
		statuses := make([]string, 0)
		for status := range stats.ByStatus {
			statuses = append(statuses, status)
		}
		sort.Strings(statuses)
		for i, status := range statuses {
			if i > 0 {
				fmt.Fprintf(out, ", ")
			}
			fmt.Fprintf(out, "%s: %d", status, stats.ByStatus[status])
		}
		fmt.Fprintf(out, "\n")
	}
	if stats.CompletionPercent > 0 {
		fmt.Fprintf(out, "Completion: %d%%\n", stats.CompletionPercent)
	}
}

func outputRecommendedCommands(out io.Writer, response *TrackResponse) {
	if len(response.RecommendedNextCommands) > 0 {
		fmt.Fprintf(out, "Recommended next commands:\n")
		for _, cmd := range response.RecommendedNextCommands {
			fmt.Fprintf(out, "  %s\n", cmd)
		}
	}
}

// outputCSV outputs CSV format
func outputCSV(out io.Writer, response *TrackResponse) error {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	switch response.Scope {
	case "task":
		if len(response.Tasks) > 0 {
			writer.Write([]string{"ID", "Name", "Status", "Priority", "Feature"})
			for _, task := range response.Tasks {
				writer.Write([]string{task.ID, task.Name, task.Status, task.Priority, task.FeatureID})
			}
		}
	case "issue":
		if len(response.Issues) > 0 {
			writer.Write([]string{"ID", "Title", "Type", "Status", "Priority"})
			for _, issue := range response.Issues {
				writer.Write([]string{issue.ID, issue.Title, issue.Type, issue.Status, issue.Priority})
			}
		}
	case "feature":
		if len(response.Tasks) > 0 {
			writer.Write([]string{"ID", "Name", "Status", "Priority", "Feature"})
			for _, task := range response.Tasks {
				writer.Write([]string{task.ID, task.Name, task.Status, task.Priority, task.FeatureID})
			}
		}
	case "backlog":
		if len(response.Issues) > 0 {
			writer.Write([]string{"ID", "Title", "Type", "Status", "Priority"})
			for _, issue := range response.Issues {
				writer.Write([]string{issue.ID, issue.Title, issue.Type, issue.Status, issue.Priority})
			}
		}
	}

	return nil
}

// outputTree outputs tree visualization
func outputTree(out io.Writer, response *TrackResponse) error {
	switch response.Scope {
	case "workspace":
		return outputWorkspaceTree(out, response)
	case "backlog":
		return outputBacklogTree(out, response)
	case "feature":
		return outputFeatureTree(out, response)
	case "task":
		return outputTaskTree(out, response)
	case "issue":
		return outputIssueTree(out, response)
	}
	return nil
}

// outputWorkspaceTree shows workspace -> backlogs hierarchy
func outputWorkspaceTree(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Workspace Overview\n")
	for i, backlog := range response.Backlogs {
		isLast := i == len(response.Backlogs)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}
		fmt.Fprintf(out, "%s%s (%s) [%d features, %d tasks, %d issues]\n",
			prefix, backlog.Name, backlog.Status, backlog.Features, backlog.Tasks, backlog.Issues)
	}
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputBacklogTree shows backlog -> issues with blocking relationships
func outputBacklogTree(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Backlog: %s\n", response.Name)
	for i, issue := range response.Issues {
		isLast := i == len(response.Issues)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}
		fmt.Fprintf(out, "%s[%s] %s (%s) - %s\n",
			prefix, issue.Type, issue.Title, issue.Status, issue.Priority)

		// Show blocking relationships if verbose
		if globalFlags.Verbose {
			if len(issue.BlockedBy) > 0 {
				for j, blocker := range issue.BlockedBy {
					isLastBlocker := j == len(issue.BlockedBy)-1
					subPrefix := "    "
					if isLast {
						subPrefix = "    "
					}
					blockMarker := "├── "
					if isLastBlocker {
						blockMarker = "└── "
					}
					fmt.Fprintf(out, "%sblocked by: %s%s\n", subPrefix, blockMarker, blocker)
				}
			}
		}
	}
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputFeatureTree shows feature -> tasks with blocking relationships
func outputFeatureTree(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Feature: %s\n", response.Name)
	if response.Summary.Total > 0 {
		completion := response.Summary.CompletionPercent
		fmt.Fprintf(out, "  (%d/%d complete, %d%%)\n", response.Summary.ByStatus["done"], response.Summary.Total, completion)
	}

	for i, task := range response.Tasks {
		isLast := i == len(response.Tasks)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}
		fmt.Fprintf(out, "%s[%s] %s (%s) - %s\n",
			prefix, task.Status, task.Name, task.Priority, task.ID)

		// Show blocking relationships if verbose
		if globalFlags.Verbose {
			subPrefix := "    "
			if isLast {
				subPrefix = "    "
			}

			if len(task.BlockedBy) > 0 {
				for j, blocker := range task.BlockedBy {
					isLastBlocker := j == len(task.BlockedBy)-1
					blockMarker := "├── "
					if isLastBlocker {
						blockMarker = "└── "
					}
					fmt.Fprintf(out, "%sblocked by: %s%s\n", subPrefix, blockMarker, blocker)
				}
			}

			if len(task.Blocks) > 0 {
				for j, blocked := range task.Blocks {
					isLastBlocked := j == len(task.Blocks)-1
					blockMarker := "├── "
					if isLastBlocked {
						blockMarker = "└── "
					}
					fmt.Fprintf(out, "%sblocks: %s%s\n", subPrefix, blockMarker, blocked)
				}
			}
		}
	}
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputTaskTree shows task with full dependency chain
func outputTaskTree(out io.Writer, response *TrackResponse) error {
	if len(response.Tasks) == 0 {
		return nil
	}

	task := response.Tasks[0]
	fmt.Fprintf(out, "Task: %s\n", task.Name)
	fmt.Fprintf(out, "├── ID:       %s\n", task.ID)
	fmt.Fprintf(out, "├── Status:   %s\n", task.Status)
	fmt.Fprintf(out, "├── Priority: %s\n", task.Priority)
	fmt.Fprintf(out, "└── Feature:  %s\n", task.FeatureID)

	if globalFlags.Verbose {
		if task.Goal != "" {
			fmt.Fprintf(out, "\nGoal:\n  %s\n", task.Goal)
		}

		if len(task.ImplementationSteps) > 0 {
			fmt.Fprintf(out, "\nImplementation Steps:\n")
			for i, step := range task.ImplementationSteps {
				isLast := i == len(task.ImplementationSteps)-1
				prefix := "├── "
				if isLast {
					prefix = "└── "
				}
				fmt.Fprintf(out, "  %s%s\n", prefix, step)
			}
		}

		if len(task.BlockedBy) > 0 {
			fmt.Fprintf(out, "\nBlocked By:\n")
			for i, blocker := range task.BlockedBy {
				isLast := i == len(task.BlockedBy)-1
				prefix := "├── "
				if isLast {
					prefix = "└── "
				}
				fmt.Fprintf(out, "  %s%s\n", prefix, blocker)
			}
		}

		if len(task.Blocks) > 0 {
			fmt.Fprintf(out, "\nBlocks:\n")
			for i, blocked := range task.Blocks {
				isLast := i == len(task.Blocks)-1
				prefix := "├── "
				if isLast {
					prefix = "└── "
				}
				fmt.Fprintf(out, "  %s%s\n", prefix, blocked)
			}
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputIssueTree shows issue with full dependency chain
func outputIssueTree(out io.Writer, response *TrackResponse) error {
	if len(response.Issues) == 0 {
		return nil
	}

	issue := response.Issues[0]
	fmt.Fprintf(out, "Issue: %s\n", issue.Title)
	fmt.Fprintf(out, "├── ID:       %s\n", issue.ID)
	fmt.Fprintf(out, "├── Type:     %s\n", issue.Type)
	fmt.Fprintf(out, "├── Status:   %s\n", issue.Status)
	fmt.Fprintf(out, "└── Priority: %s\n", issue.Priority)

	if globalFlags.Verbose {
		if issue.Description != "" {
			fmt.Fprintf(out, "\nDescription:\n  %s\n", issue.Description)
		}

		if len(issue.BlockedBy) > 0 {
			fmt.Fprintf(out, "\nBlocked By:\n")
			for i, blocker := range issue.BlockedBy {
				isLast := i == len(issue.BlockedBy)-1
				prefix := "├── "
				if isLast {
					prefix = "└── "
				}
				fmt.Fprintf(out, "  %s%s\n", prefix, blocker)
			}
		}

		if len(issue.Blocks) > 0 {
			fmt.Fprintf(out, "\nBlocks:\n")
			for i, blocked := range issue.Blocks {
				isLast := i == len(issue.Blocks)-1
				prefix := "├── "
				if isLast {
					prefix = "└── "
				}
				fmt.Fprintf(out, "  %s%s\n", prefix, blocked)
			}
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputGraph outputs ASCII graph visualization
func outputGraph(out io.Writer, response *TrackResponse) error {
	switch response.Scope {
	case "workspace":
		return outputWorkspaceGraph(out, response)
	case "backlog":
		return outputBacklogGraph(out, response)
	case "feature":
		return outputFeatureGraph(out, response)
	case "task":
		return outputTaskGraph(out, response)
	case "issue":
		return outputIssueGraph(out, response)
	}
	return nil
}

// outputWorkspaceGraph shows backlog status summary
func outputWorkspaceGraph(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Workspace Overview - Backlog Status\n\n")

	// Create a simple status distribution graph
	statusCounts := make(map[string]int)
	for _, backlog := range response.Backlogs {
		statusCounts[backlog.Status]++
	}

	statusOrder := []string{"initial", "active", "release", "archived"}
	for _, status := range statusOrder {
		if count, ok := statusCounts[status]; ok && count > 0 {
			bar := strings.Repeat("█", count)
			fmt.Fprintf(out, "%-15s %s (%d)\n", status, bar, count)
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputBacklogGraph shows issue blocking graph
func outputBacklogGraph(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Backlog: %s - Issue Blocking Graph\n\n", response.Name)

	// Build a dependency map
	blockedByMap := make(map[string][]string)
	blocksMap := make(map[string][]string)

	for _, issue := range response.Issues {
		if len(issue.BlockedBy) > 0 {
			blockedByMap[issue.ID] = issue.BlockedBy
		}
		if len(issue.Blocks) > 0 {
			blocksMap[issue.ID] = issue.Blocks
		}
	}

	// Show issues with blocking relationships
	for _, issue := range response.Issues {
		if len(issue.BlockedBy) > 0 || len(issue.Blocks) > 0 {
			fmt.Fprintf(out, "[%s] %s\n", issue.Status, truncate(issue.Title, 40))

			for _, blocker := range issue.BlockedBy {
				fmt.Fprintf(out, "    ↑ blocked by: %s\n", blocker)
			}
			for _, blocked := range issue.Blocks {
				fmt.Fprintf(out, "    ↓ blocks: %s\n", blocked)
			}
			fmt.Fprintf(out, "\n")
		}
	}

	if len(blockedByMap) == 0 && len(blocksMap) == 0 {
		fmt.Fprintf(out, "No blocking relationships found\n\n")
	}

	outputRecommendedCommands(out, response)
	return nil
}

// outputFeatureGraph shows task dependency graph
func outputFeatureGraph(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Feature: %s - Task Dependency Graph\n\n", response.Name)

	// Build dependency chains
	type Node struct {
		id    string
		name  string
		level int
	}

	// Find root tasks (not blocked)
	rootTasks := []TaskTrackItem{}
	for _, task := range response.Tasks {
		if len(task.BlockedBy) == 0 {
			rootTasks = append(rootTasks, task)
		}
	}

	// Show chains starting from roots
	for _, root := range rootTasks {
		fmt.Fprintf(out, "[%s] %s\n", root.Status, truncate(root.Name, 40))
		showTaskChain(out, root, response.Tasks, 1)
	}

	// Show blocked tasks not in chains
	processedTasks := make(map[string]bool)
	for _, root := range rootTasks {
		markTaskChain(&processedTasks, root, response.Tasks)
	}

	unlinked := []TaskTrackItem{}
	for _, task := range response.Tasks {
		if !processedTasks[task.ID] && len(task.BlockedBy) > 0 {
			unlinked = append(unlinked, task)
		}
	}

	if len(unlinked) > 0 {
		fmt.Fprintf(out, "\nIndependent chains:\n")
		for _, task := range unlinked {
			if len(task.BlockedBy) > 0 {
				fmt.Fprintf(out, "[%s] %s\n", task.Status, truncate(task.Name, 40))
				for _, blocker := range task.BlockedBy {
					fmt.Fprintf(out, "    ← %s\n", blocker)
				}
			}
		}
	}

	if len(rootTasks) == 0 && len(unlinked) == 0 {
		fmt.Fprintf(out, "No dependencies found\n")
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// showTaskChain recursively shows task dependency chain
func showTaskChain(out io.Writer, task TaskTrackItem, allTasks []TaskTrackItem, depth int) {
	// Find tasks blocked by this one
	for _, other := range allTasks {
		for _, blocker := range other.BlockedBy {
			if blocker == task.ID {
				indent := strings.Repeat("    ", depth)
				fmt.Fprintf(out, "%s↓ [%s] %s\n", indent, other.Status, truncate(other.Name, 35))
				showTaskChain(out, other, allTasks, depth+1)
			}
		}
	}
}

// markTaskChain marks all tasks in a chain as processed
func markTaskChain(processed *map[string]bool, task TaskTrackItem, allTasks []TaskTrackItem) {
	(*processed)[task.ID] = true
	for _, other := range allTasks {
		for _, blocker := range other.BlockedBy {
			if blocker == task.ID && !(*processed)[other.ID] {
				markTaskChain(processed, other, allTasks)
			}
		}
	}
}

// outputTaskGraph shows task with blocking context
func outputTaskGraph(out io.Writer, response *TrackResponse) error {
	if len(response.Tasks) == 0 {
		return nil
	}

	task := response.Tasks[0]
	fmt.Fprintf(out, "Task: %s\n\n", task.Name)

	// Show blocking relationships
	if len(task.BlockedBy) > 0 {
		fmt.Fprintf(out, "Blocked By (dependencies):\n")
		for _, blocker := range task.BlockedBy {
			fmt.Fprintf(out, "    %s\n", blocker)
			fmt.Fprintf(out, "         ↓\n")
		}
	}

	fmt.Fprintf(out, "    [%s] %s\n", task.Status, task.ID)

	if len(task.Blocks) > 0 {
		fmt.Fprintf(out, "         ↓\n")
		fmt.Fprintf(out, "Blocks (dependents):\n")
		for _, blocked := range task.Blocks {
			fmt.Fprintf(out, "    %s\n", blocked)
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// outputIssueGraph shows issue with blocking context
func outputIssueGraph(out io.Writer, response *TrackResponse) error {
	if len(response.Issues) == 0 {
		return nil
	}

	issue := response.Issues[0]
	fmt.Fprintf(out, "Issue: %s\n\n", issue.Title)

	// Show blocking relationships
	if len(issue.BlockedBy) > 0 {
		fmt.Fprintf(out, "Blocked By (dependencies):\n")
		for _, blocker := range issue.BlockedBy {
			fmt.Fprintf(out, "    %s\n", blocker)
			fmt.Fprintf(out, "         ↓\n")
		}
	}

	fmt.Fprintf(out, "    [%s] %s\n", issue.Status, issue.ID)

	if len(issue.Blocks) > 0 {
		fmt.Fprintf(out, "         ↓\n")
		fmt.Fprintf(out, "Blocks (dependents):\n")
		for _, blocked := range issue.Blocks {
			fmt.Fprintf(out, "    %s\n", blocked)
		}
	}

	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)
	return nil
}

// Helper functions

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func groupTasksByStatus(tasks []TaskTrackItem) []TaskTrackItem {
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Status != tasks[j].Status {
			return statusOrder(tasks[i].Status) < statusOrder(tasks[j].Status)
		}
		return tasks[i].Name < tasks[j].Name
	})
	return tasks
}

func groupTasksByPriority(tasks []TaskTrackItem) []TaskTrackItem {
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority != tasks[j].Priority {
			return priorityOrder(tasks[i].Priority) < priorityOrder(tasks[j].Priority)
		}
		return tasks[i].Name < tasks[j].Name
	})
	return tasks
}

func groupIssuesByStatus(issues []IssueTrackItem) []IssueTrackItem {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Status != issues[j].Status {
			return statusOrder(issues[i].Status) < statusOrder(issues[j].Status)
		}
		return issues[i].Title < issues[j].Title
	})
	return issues
}

func groupIssuesByPriority(issues []IssueTrackItem) []IssueTrackItem {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Priority != issues[j].Priority {
			return priorityOrder(issues[i].Priority) < priorityOrder(issues[j].Priority)
		}
		return issues[i].Title < issues[j].Title
	})
	return issues
}

func statusOrder(status string) int {
	switch status {
	case "ready":
		return 0
	case "in_progress":
		return 1
	case "blocked":
		return 2
	case "done":
		return 3
	case "cancelled":
		return 4
	default:
		return 5
	}
}

func priorityOrder(priority string) int {
	switch priority {
	case "P0":
		return 0
	case "P1":
		return 1
	case "P2":
		return 2
	case "P3":
		return 3
	case "P4":
		return 4
	case "P5":
		return 5
	default:
		return 6
	}
}

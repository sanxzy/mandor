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
	if globalFlags.Summary {
		return outputTableSummary(out, response)
	}

	switch response.Scope {
	case "workspace":
		return outputWorkspaceTable(out, response)
	case "project":
		return outputProjectTable(out, response)
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
	fmt.Fprintf(out, "%-30s %-10s %10s %10s %10s\n", "Project", "Status", "Features", "Tasks", "Issues")
	fmt.Fprintf(out, strings.Repeat("-", 70)+"\n")

	for _, proj := range response.Projects {
		fmt.Fprintf(out, "%-30s %-10s %10d %10d %10d\n",
			truncate(proj.Name, 30), proj.Status, proj.Features, proj.Tasks, proj.Issues)
	}

	fmt.Fprintf(out, "\n")
	outputSummaryStats(out, response.Summary)
	fmt.Fprintf(out, "\n")
	outputRecommendedCommands(out, response)

	return nil
}

func outputProjectTable(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Project: %s\n", response.Name)
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
	fmt.Fprintf(out, "Project:  %s\n", issue.ProjectID)

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

func outputTableSummary(out io.Writer, response *TrackResponse) error {
	fmt.Fprintf(out, "Summary for %s\n", response.Scope)
	outputSummaryStats(out, response.Summary)
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
	case "project":
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
	// Tree not fully implemented for brevity
	return nil
}

// outputGraph outputs ASCII graph visualization
func outputGraph(out io.Writer, response *TrackResponse) error {
	// Graph not fully implemented for brevity
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

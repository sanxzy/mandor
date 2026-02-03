package task

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

func NewSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "summary <feature_id>",
		Short:      "Display task summary for a feature",
		Deprecated: "Use 'mandor track feature <id>' instead. Use '--summary' flag for aggregated counts only.",
		Long:       "Display a summary of tasks in a feature grouped by status with priority information.\n\nDEPRECATED: Use 'mandor track feature <id>' instead, optionally with '--summary' flag for count-only output.",
		Args:       cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			featureID := args[0]

			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			// Group tasks by status
			tasksByStatus := make(map[string][]domain.TaskListItem)
			statusOrder := []string{
				domain.TaskStatusReady,
				domain.TaskStatusInProgress,
				domain.TaskStatusBlocked,
				domain.TaskStatusPending,
				domain.TaskStatusDone,
				domain.TaskStatusCancelled,
			}

			input := &domain.TaskListInput{
				FeatureID: featureID,
			}

			output, err := svc.ListTasks(input)
			if err != nil {
				return err
			}

			if output == nil || output.Total == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No tasks in feature %s.\n", featureID)
				return nil
			}

			// Group by status
			for _, task := range output.Tasks {
				tasksByStatus[task.Status] = append(tasksByStatus[task.Status], task)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Tasks in %s:\n\n", featureID)

			// Print by status
			for _, status := range statusOrder {
				tasks := tasksByStatus[status]
				if len(tasks) == 0 {
					continue
				}

				statusLabel := getTaskStatusLabel(status)
				fmt.Fprintf(out, "%s (%d tasks)\n", statusLabel, len(tasks))
				fmt.Fprintf(out, "| # | Task ID | Name | Priority | Status |\n")
				fmt.Fprintf(out, "|---|---------|------|----------|--------|\n")

				sort.Slice(tasks, func(i, j int) bool {
					return tasks[i].CreatedAt < tasks[j].CreatedAt
				})

				for idx, task := range tasks {
					name := task.Name
					if len(name) > 40 {
						name = name[:37] + "..."
					}
					fmt.Fprintf(out, "| %d | %s | %s | %s | %s |\n",
						idx+1,
						task.ID,
						name,
						task.Priority,
						task.Status)
				}

				fmt.Fprintln(out)
			}

			return nil
		},
	}

	return cmd
}

func getTaskStatusLabel(status string) string {
	labels := map[string]string{
		domain.TaskStatusPending:    "Pending",
		domain.TaskStatusReady:      "Ready",
		domain.TaskStatusInProgress: "In Progress",
		domain.TaskStatusBlocked:    "Blocked",
		domain.TaskStatusDone:       "Done",
		domain.TaskStatusCancelled:  "Cancelled",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return "Unknown"
}

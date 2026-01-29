package task

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	listFeatureID      string
	listProjectID      string
	listStatus         string
	listPriority       string
	listJSON           bool
	listIncludeDeleted bool
	listSort           string
	listOrder          string
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [--feature <id>] [--project <id>] [--status <status>] [--priority <priority>] [--json] [--include-deleted]",
		Short: "List tasks",
		Long:  "List all tasks in the workspace or filter by feature/project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if listStatus != "" {
				if !domain.ValidateTaskStatus(listStatus) {
					return domain.NewValidationError(fmt.Sprintf("Invalid status: '%s'. Valid values: pending, ready, in_progress, blocked, done, cancelled", listStatus))
				}
			}

			if listPriority != "" {
				if !domain.ValidatePriority(listPriority) {
					return domain.NewValidationError(fmt.Sprintf("Invalid priority: '%s'. Valid values: P0, P1, P2, P3, P4, P5", listPriority))
				}
			}

			input := &domain.TaskListInput{
				FeatureID:      listFeatureID,
				ProjectID:      listProjectID,
				Status:         listStatus,
				Priority:       listPriority,
				IncludeDeleted: listIncludeDeleted,
				JSON:           listJSON,
				Sort:           listSort,
				Order:          listOrder,
			}

			output, err := svc.ListTasks(input)
			if err != nil {
				return err
			}

			if listJSON {
				out := cmd.OutOrStdout()
				encoder := json.NewEncoder(out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(output)
			}

			out := cmd.OutOrStdout()

			if listFeatureID != "" {
				fmt.Fprintf(out, "Tasks in %s:\n", listFeatureID)
			} else if listProjectID != "" {
				fmt.Fprintf(out, "Tasks in project %s:\n", listProjectID)
			} else {
				fmt.Fprintf(out, "All tasks:\n")
			}

			fmt.Fprintf(out, "%-44s %-10s %-12s %-6s %s\n", "ID", "Status", "Priority", "Feature", "Name")
			fmt.Fprintln(out, strings.Repeat("-", 100))

			for _, t := range output.Tasks {
				name := t.Name
				if len(name) > 30 {
					name = name[:27] + "..."
				}
				featureShort := t.FeatureID
				if len(featureShort) > 6 {
					featureShort = featureShort[:6] + "..."
				}
				fmt.Fprintf(out, "%-44s %-10s %-12s %-6s %s\n", t.ID, t.Status, t.Priority, featureShort, name)
			}

			fmt.Fprintf(out, "\nTotal: %d", output.Total)
			if output.Deleted > 0 {
				fmt.Fprintf(out, " (%d deleted)", output.Deleted)
			}
			fmt.Fprintln(out)

			return nil
		},
	}

	cmd.Flags().StringVarP(&listFeatureID, "feature", "f", "", "Filter by feature ID")
	cmd.Flags().StringVarP(&listProjectID, "project", "p", "", "Filter by project ID")
	cmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (pending, ready, in_progress, blocked, done, cancelled)")
	cmd.Flags().StringVar(&listPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&listIncludeDeleted, "include-deleted", false, "Include deleted tasks")
	cmd.Flags().StringVar(&listSort, "sort", "priority", "Sort field: priority, created_at, name")
	cmd.Flags().StringVar(&listOrder, "order", "desc", "Sort order: asc, desc")

	return cmd
}

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
	blockedProjectID string
	blockedFeatureID string
	blockedPriority  string
	blockedJSON      bool
)

func NewBlockedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocked [--project <id>] [--feature <id>] [--priority <priority>] [--json]",
		Short: "List blocked tasks",
		Long:  "List all tasks with status='blocked' that are waiting for dependencies to complete.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if blockedPriority != "" {
				if !domain.ValidatePriority(blockedPriority) {
					return domain.NewValidationError(fmt.Sprintf("Invalid priority: '%s'. Valid values: P0, P1, P2, P3, P4, P5", blockedPriority))
				}
			}

			input := &domain.TaskListInput{
				FeatureID:      blockedFeatureID,
				ProjectID:      blockedProjectID,
				Status:         domain.TaskStatusBlocked,
				Priority:       blockedPriority,
				IncludeDeleted: false,
				JSON:           blockedJSON,
				Sort:           "priority",
				Order:          "asc",
			}

			output, err := svc.ListTasks(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if blockedJSON {
				result := map[string]interface{}{
					"tasks": output.Tasks,
					"total": output.Total,
				}
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Fprintln(out, string(jsonBytes))
				return nil
			}

			if output.Total == 0 {
				fmt.Fprintln(out, "No blocked tasks found.")
				if blockedFeatureID != "" {
					fmt.Fprintf(out, "\nUnblock tasks: mandor task update <id> --depends <id,...>\n")
				} else if blockedProjectID != "" {
					fmt.Fprintf(out, "\nCreate a feature: mandor feature create <name> --project %s\n", blockedProjectID)
				}
				return nil
			}

			if blockedFeatureID != "" {
				fmt.Fprintf(out, "Blocked tasks in %s:\n", blockedFeatureID)
			} else if blockedProjectID != "" {
				fmt.Fprintf(out, "Blocked tasks in project %s:\n", blockedProjectID)
			} else {
				fmt.Fprintln(out, "Blocked tasks (waiting for dependencies):")
			}

			fmt.Fprintf(out, "%-44s %-12s %-6s %s\n", "ID", "Priority", "Feature", "Name")
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
				fmt.Fprintf(out, "%-44s %-12s %-6s %s\n", t.ID, t.Priority, featureShort, name)
			}

			fmt.Fprintf(out, "\nTotal: %d\n", output.Total)

			return nil
		},
	}

	cmd.Flags().StringVarP(&blockedProjectID, "project", "p", "", "Filter by project ID")
	cmd.Flags().StringVarP(&blockedFeatureID, "feature", "f", "", "Filter by feature ID")
	cmd.Flags().StringVar(&blockedPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&blockedJSON, "json", false, "Output as JSON")

	return cmd
}

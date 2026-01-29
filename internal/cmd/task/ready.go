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
	readyProjectID string
	readyFeatureID string
	readyPriority  string
	readyJSON      bool
)

func NewReadyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ready [--project <id>] [--feature <id>] [--priority <priority>] [--json]",
		Short: "List ready tasks",
		Long:  "List all tasks with status='ready' that are available to work on (no blocking dependencies).",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if readyPriority != "" {
				if !domain.ValidatePriority(readyPriority) {
					return domain.NewValidationError(fmt.Sprintf("Invalid priority: '%s'. Valid values: P0, P1, P2, P3, P4, P5", readyPriority))
				}
			}

			input := &domain.TaskListInput{
				FeatureID:      readyFeatureID,
				ProjectID:      readyProjectID,
				Status:         domain.TaskStatusReady,
				Priority:       readyPriority,
				IncludeDeleted: false,
				JSON:           readyJSON,
				Sort:           "priority",
				Order:          "asc",
			}

			output, err := svc.ListTasks(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if readyJSON {
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
				fmt.Fprintln(out, "No ready tasks found.")
				if readyFeatureID != "" {
					fmt.Fprintf(out, "\nCreate a task: mandor task create <name> --feature %s\n", readyFeatureID)
				} else if readyProjectID != "" {
					fmt.Fprintf(out, "\nCreate a feature: mandor feature create <name> --project %s\n", readyProjectID)
				}
				return nil
			}

			if readyFeatureID != "" {
				fmt.Fprintf(out, "Ready tasks in %s:\n", readyFeatureID)
			} else if readyProjectID != "" {
				fmt.Fprintf(out, "Ready tasks in project %s:\n", readyProjectID)
			} else {
				fmt.Fprintln(out, "Ready tasks (available to work on):")
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

	cmd.Flags().StringVarP(&readyProjectID, "project", "p", "", "Filter by project ID")
	cmd.Flags().StringVarP(&readyFeatureID, "feature", "f", "", "Filter by feature ID")
	cmd.Flags().StringVar(&readyPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&readyJSON, "json", false, "Output as JSON")

	return cmd
}

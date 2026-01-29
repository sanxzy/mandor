package project

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var jsonDetailOutput bool

func NewDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail <id>",
		Short: "Show project details",
		Long:  "Show detailed information about a specific project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			detail, err := svc.GetProjectDetail(args[0])
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if jsonDetailOutput {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(detail)
			}

			deletedLabel := ""
			if detail.Status == domain.ProjectStatusDeleted {
				deletedLabel = " [DELETED]"
			}

			fmt.Fprintln(out, "╔════════════════════════════════════════════════════════════╗")
			fmt.Fprintf(out, "║ PROJECT: %s%s", detail.ID, deletedLabel)
			for i := 0; i < 57-len(detail.ID)-len(deletedLabel); i++ {
				fmt.Fprint(out, " ")
			}
			fmt.Fprintln(out, "║")
			fmt.Fprintln(out, "╚════════════════════════════════════════════════════════════╝")
			fmt.Fprintln(out)
			fmt.Fprintf(out, "ID:          %s\n", detail.ID)
			fmt.Fprintf(out, "Name:        %s\n", detail.Name)
			fmt.Fprintf(out, "Goal:        %s\n", detail.Goal)
			fmt.Fprintf(out, "Status:      %s\n", detail.Status)
			fmt.Fprintf(out, "Strict:      %t\n", detail.Strict)
			fmt.Fprintln(out)
			fmt.Fprintf(out, "Created:     %s\n", detail.CreatedAt)
			fmt.Fprintf(out, "Updated:     %s\n", detail.UpdatedAt)
			fmt.Fprintf(out, "Created by:  %s\n", detail.CreatedBy)
			fmt.Fprintf(out, "Updated by:  %s\n", detail.UpdatedBy)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "SCHEMA")
			fmt.Fprintln(out, "══════")
			fmt.Fprintf(out, "Version:     %s\n", detail.Schema.Version)
			fmt.Fprintln(out, "Dependency Rules:")
			fmt.Fprintf(out, "  - Task:    %s\n", detail.Schema.Rules.Task.Dependency)
			fmt.Fprintf(out, "  - Feature: %s\n", detail.Schema.Rules.Feature.Dependency)
			fmt.Fprintf(out, "  - Issue:   %s\n", detail.Schema.Rules.Issue.Dependency)
			fmt.Fprintf(out, "Priority:    %s (default: %s)\n", joinLevels(detail.Schema.Rules.Priority.Levels), detail.Schema.Rules.Priority.Default)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "STATISTICS")
			fmt.Fprintln(out, "══════════")
			fmt.Fprintf(out, "Features:    %d total\n", detail.Stats.Features.Total)
			fmt.Fprintf(out, "Tasks:       %d total\n", detail.Stats.Tasks.Total)
			fmt.Fprintf(out, "Issues:      %d total\n", detail.Stats.Issues.Total)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "ACTIVITY")
			fmt.Fprintln(out, "══════")
			fmt.Fprintf(out, "Total events: %d\n", detail.Activity.TotalEvents)
			if detail.Activity.LastActivity != "" {
				fmt.Fprintf(out, "Last activity: %s\n", detail.Activity.LastActivity)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonDetailOutput, "json", false, "Output in JSON format")

	return cmd
}

func joinLevels(levels []string) string {
	result := ""
	for i, l := range levels {
		if i > 0 {
			result += ", "
		}
		result += l
	}
	return result
}

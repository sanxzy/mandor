package project

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	includeDeleted bool
	includeGoal    bool
	jsonOutput     bool
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Long:  "List all projects in the workspace.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			output, err := svc.ListProjects(includeDeleted, includeGoal)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if jsonOutput {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(output)
			}

			if len(output.Projects) == 0 {
				fmt.Fprintln(out, "No projects in workspace.")
				fmt.Fprintln(out, "Run 'mandor project create <id> --name \"<name>\" --goal \"<goal>\"' to create your first project.")
				return nil
			}

			deletedStr := ""
			if includeDeleted && output.Deleted > 0 {
				deletedStr = fmt.Sprintf(", %d deleted", output.Deleted)
			}
			fmt.Fprintf(out, "Projects (%d total%s)\n", output.Total, deletedStr)
			fmt.Fprintln(out, "═══════════════════════════════════════════")
			fmt.Fprintln(out)

			for i, p := range output.Projects {
				num := i + 1
				deletedLabel := ""
				if p.Status == domain.ProjectStatusDeleted {
					deletedLabel = " [DELETED]"
				}
				fmt.Fprintf(out, "[%d] %s%s\n", num, p.ID, deletedLabel)
				fmt.Fprintf(out, "    Name:        %s\n", p.Name)
				if includeGoal && p.Goal != "" {
					goalDisplay := p.Goal
					if len(goalDisplay) > 50 {
						goalDisplay = goalDisplay[:47] + "..."
					}
					fmt.Fprintf(out, "    Goal:        %s\n", goalDisplay)
				}
				fmt.Fprintf(out, "    Status:      %s\n", p.Status)
				fmt.Fprintf(out, "    Features:    %d\n", p.Features)
				fmt.Fprintf(out, "    Tasks:       %d\n", p.Tasks)
				fmt.Fprintf(out, "    Issues:      %d\n", p.Issues)
				fmt.Fprintf(out, "    Created:     %s\n", p.CreatedAt)
				fmt.Fprintf(out, "    Updated:     %s\n", p.UpdatedAt)
				fmt.Fprintln(out)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&includeDeleted, "include-deleted", false, "Include deleted projects")
	cmd.Flags().BoolVar(&includeGoal, "include-goal", false, "Show goal for each project")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

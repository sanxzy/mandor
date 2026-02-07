package backlog

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	updateName       string
	updateGoal       string
	updateTaskDep    string
	updateFeatureDep string
	updateIssueDep   string
	updateStrict     string
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update project metadata",
		Long:  "Update metadata for an existing project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewBacklogService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.BacklogUpdateInput{
				ID:         args[0],
				Name:       nil,
				Goal:       nil,
				TaskDep:    nil,
				FeatureDep: nil,
				IssueDep:   nil,
				Strict:     nil,
			}

			if updateName != "" {
				input.Name = &updateName
			}
			if updateGoal != "" {
				input.Goal = &updateGoal
			}
			if updateTaskDep != "" {
				input.TaskDep = &updateTaskDep
			}
			if updateFeatureDep != "" {
				input.FeatureDep = &updateFeatureDep
			}
			if updateIssueDep != "" {
				input.IssueDep = &updateIssueDep
			}
			if updateStrict != "" {
				if !domain.ValidateBooleanValue(updateStrict) {
					return domain.NewValidationError("Invalid value for --strict. Use: true, false, yes, no, 1, or 0.")
				}
				val := domain.ParseBooleanValue(updateStrict)
				input.Strict = &val
			}

			if input.Name == nil && input.Goal == nil && input.TaskDep == nil && input.FeatureDep == nil && input.IssueDep == nil && input.Strict == nil {
				return domain.NewValidationError("No updates specified. Use --name, --goal, --task-dep, --feature-dep, --issue-dep, or --strict.")
			}

			if err := svc.ValidateUpdateInput(input); err != nil {
				return err
			}

			changes, err := svc.UpdateBacklog(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if len(changes) == 0 {
				fmt.Fprintf(out, "✓ No changes to project: %s\n", args[0])
				return nil
			}

			fmt.Fprintf(out, "✓ Project updated: %s\n", args[0])
			fmt.Fprintln(out, "  Changes:")
			for _, change := range changes {
				switch change {
				case "name":
					fmt.Fprintf(out, "    - name set to: %s\n", updateName)
				case "goal":
					goalDisplay := updateGoal
					if len(goalDisplay) > 50 {
						goalDisplay = goalDisplay[:47] + "..."
					}
					fmt.Fprintf(out, "    - goal updated (%d chars)\n", len(updateGoal))
				case "strict":
					fmt.Fprintf(out, "    - strict: %s\n", updateStrict)
				case "task_dep":
					fmt.Fprintf(out, "    - task_dep: %s\n", updateTaskDep)
				case "feature_dep":
					fmt.Fprintf(out, "    - feature_dep: %s\n", updateFeatureDep)
				case "issue_dep":
					fmt.Fprintf(out, "    - issue_dep: %s\n", updateIssueDep)
				}
			}

			// Show full project state
			fmt.Fprintln(out)
			detail, err := svc.GetBacklogDetail(args[0])
			if err == nil {
				fmt.Fprintf(out, "BACKLOG: %s", detail.ID)
				if detail.Status == domain.BacklogStatusDeleted {
					fmt.Fprint(out, " [DELETED]")
				}
				fmt.Fprintln(out)
				fmt.Fprintln(out, "------------------------------------------------------------")

				fmt.Fprintf(out, "  Name:   %s\n", detail.Name)
				fmt.Fprintf(out, "  Status: %s\n", detail.Status)
				fmt.Fprintf(out, "  Strict: %t\n", detail.Strict)

				if detail.Goal != "" {
					goalPreview := detail.Goal
					if len(goalPreview) > 200 {
						goalPreview = goalPreview[:197] + "..."
					}
					fmt.Fprintf(out, "\n  Goal:\n    %s\n", goalPreview)
				}

				fmt.Fprintln(out)
				fmt.Fprintln(out, "  Schema Rules:")
				fmt.Fprintf(out, "    Task:     %s\n", detail.Schema.Rules.Task.Dependency)
				fmt.Fprintf(out, "    Feature:  %s\n", detail.Schema.Rules.Feature.Dependency)
				fmt.Fprintf(out, "    Issue:    %s\n", detail.Schema.Rules.Issue.Dependency)

				fmt.Fprintln(out)
				fmt.Fprintln(out, "  Stats:")
				fmt.Fprintf(out, "    Features: %d\n", detail.Stats.Features.Total)
				fmt.Fprintf(out, "    Tasks:    %d\n", detail.Stats.Tasks.Total)
				fmt.Fprintf(out, "    Issues:   %d\n", detail.Stats.Issues.Total)

				fmt.Fprintf(out, "\n  Created: %s by %s\n", detail.CreatedAt, detail.CreatedBy)
				fmt.Fprintf(out, "  Updated: %s by %s\n", detail.UpdatedAt, detail.UpdatedBy)

				// AI-friendly recommendations
				fmt.Fprintln(out)
				fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
				fmt.Fprintln(out, "  RECOMMENDATIONS:")
				fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")

				switch detail.Status {
				case domain.BacklogStatusInitial:
					fmt.Fprintf(out, "  • Activate backlog? Run: mandor backlog update %s --strict true\n", args[0])
					fmt.Fprintf(out, "  • Create features: mandor feature create --backlog %s --name \"Feature Name\"\n", args[0])
				case domain.BacklogStatusActive:
					fmt.Fprintf(out, "  • View all features: mandor feature list --backlog %s\n", args[0])
					fmt.Fprintf(out, "  • View all tasks: mandor task list --backlog %s\n", args[0])
					fmt.Fprintf(out, "  • View all issues: mandor issue list --backlog %s\n", args[0])
				case domain.BacklogStatusDone:
					fmt.Fprintf(out, "  • Backlog complete. Archive with: mandor backlog delete %s\n", args[0])
				}

				fmt.Fprintf(out, "  • View full details: mandor backlog detail %s\n", args[0])
				fmt.Fprintf(out, "  • List all backlogs: mandor backlog list\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&updateName, "name", "n", "", "New project name")
	cmd.Flags().StringVarP(&updateGoal, "goal", "g", "", "New project goal (min 500 chars)")
	cmd.Flags().StringVar(&updateTaskDep, "task-dep", "", "Update task dependency rule (same_project_only, cross_project_allowed, disabled)")
	cmd.Flags().StringVar(&updateFeatureDep, "feature-dep", "", "Update feature dependency rule (same_project_only, cross_project_allowed, disabled)")
	cmd.Flags().StringVar(&updateIssueDep, "issue-dep", "", "Update issue dependency rule (same_project_only, cross_project_allowed, disabled)")
	cmd.Flags().StringVar(&updateStrict, "strict", "", "Toggle strict mode (true/false/yes/no/1/0)")

	return cmd
}

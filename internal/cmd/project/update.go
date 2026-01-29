package project

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
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.ProjectUpdateInput{
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

			changes, err := svc.UpdateProject(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if len(changes) == 0 {
				fmt.Fprintf(out, "✓ No changes to project: %s\n", args[0])
				return nil
			}

			project, err := svc.GetProject(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "✓ Project updated: %s\n", args[0])
			fmt.Fprintln(out, "  Changes:")
			for _, change := range changes {
				switch change {
				case "name":
					fmt.Fprintf(out, "    - name: %s\n", project.Name)
				case "goal":
					goalDisplay := project.Goal
					if len(goalDisplay) > 50 {
						goalDisplay = goalDisplay[:47] + "..."
					}
					fmt.Fprintf(out, "    - goal: %s\n", goalDisplay)
				case "strict":
					fmt.Fprintf(out, "    - strict: %t\n", project.Strict)
				case "task_dep":
					fmt.Fprintf(out, "    - task_dep: %s\n", updateTaskDep)
				case "feature_dep":
					fmt.Fprintf(out, "    - feature_dep: %s\n", updateFeatureDep)
				case "issue_dep":
					fmt.Fprintf(out, "    - issue_dep: %s\n", updateIssueDep)
				}
			}
			fmt.Fprintf(out, "  Updated: %s\n", project.UpdatedAt.Format("2006-01-02T15:04:05Z"))

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

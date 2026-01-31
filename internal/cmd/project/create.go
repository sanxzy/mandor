package project

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	name       string
	goal       string
	taskDep    string
	featureDep string
	issueDep   string
	strict     bool
	yesFlag    bool
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a new project",
		Long:  "Create a new project in the workspace with the specified ID, name, and goal.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.ProjectCreateInput{
				ID:         args[0],
				Name:       name,
				Goal:       goal,
				TaskDep:    taskDep,
				FeatureDep: featureDep,
				IssueDep:   issueDep,
				Strict:     strict,
			}

			if !yesFlag && name == "" {
				out := cmd.OutOrStdout()
				in := bufio.NewReader(os.Stdin)
				fmt.Fprint(out, "Project name: ")
				line, _ := in.ReadString('\n')
				name = line[:len(line)-1]
				if len(name) == 0 {
					return domain.NewValidationError("Project name is required.")
				}
			}

			if name == "" {
				return domain.NewValidationError("Project name is required.")
			}
			if goal == "" {
				return domain.NewValidationError("Project goal is required.")
			}

			if !domain.ValidateGoalLength(goal) {
				minLen := domain.GoalMinLength
				if util.IsDevelopment() {
					minLen = domain.GoalMinLengthDevelopment
				}
				return domain.NewValidationError(fmt.Sprintf("Project goal must be at least %d characters. Current length: %d characters.", minLen, len(goal)))
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			if err := svc.CreateProject(input); err != nil {
				return err
			}

			project, err := svc.GetProject(input.ID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Project created: %s\n", input.ID)
			fmt.Fprintf(out, "  Name:        %s\n", project.Name)
			fmt.Fprintf(out, "  Goal:        %s\n", project.Goal)
			fmt.Fprintf(out, "  Task Dep:    %s\n", taskDep)
			fmt.Fprintf(out, "  Feature Dep: %s\n", featureDep)
			fmt.Fprintf(out, "  Issue Dep:   %s\n", issueDep)
			fmt.Fprintf(out, "  Strict:      %t\n", strict)
			fmt.Fprintf(out, "  Location:    .mandor/projects/%s/\n", input.ID)
			fmt.Fprintf(out, "  Created:     %s\n", project.CreatedAt.Format("2006-01-02T15:04:05Z"))
			fmt.Fprintf(out, "  Creator:     %s\n", project.CreatedBy)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Next steps:")
			fmt.Fprintln(out, "  1. Add features: mandor feature create \"Feature Name\" --project "+input.ID)
			fmt.Fprintln(out, "  2. Add tasks: mandor task create \"Task Name\" --project "+input.ID)
			fmt.Fprintln(out, "  3. Add issues: mandor issue create \"Issue Name\" --project "+input.ID)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Project display name")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Project goal/objectives (required, min 500 characters)")
	cmd.Flags().StringVar(&taskDep, "task-dep", "same_project_only", "same_project_only,Task dependency rule ( cross_project_allowed, disabled)")
	cmd.Flags().StringVar(&featureDep, "feature-dep", "cross_project_allowed", "Feature dependency rule (same_project_only, cross_project_allowed, disabled)")
	cmd.Flags().StringVar(&issueDep, "issue-dep", "same_project_only", "Issue dependency rule (same_project_only, cross_project_allowed, disabled)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Enforce strict dependency rules")
	cmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Non-interactive mode")

	return cmd
}

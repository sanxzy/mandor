package backlog

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
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
		Short: "Create a new backlog",
		Long:  "Create a new backlog in the workspace with the specified ID, name, and goal.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewBacklogService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.BacklogCreateInput{
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
				fmt.Fprint(out, "Backlog name: ")
				line, _ := in.ReadString('\n')
				name = line[:len(line)-1]
				if len(name) == 0 {
					return domain.NewValidationError("Backlog name is required.")
				}
			}

			if name == "" {
				return domain.NewValidationError("Backlog name is required.")
			}
			if goal == "" {
				return domain.NewValidationError("Backlog goal is required.")
			}

			minLen := svc.GetBacklogGoalMinLength()
			if !domain.ValidateGoalLength(goal, minLen) {
				return domain.NewValidationError(fmt.Sprintf("Backlog goal must be at least %d characters. Current length: %d characters.", minLen, len(goal)))
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			if err := svc.CreateBacklog(input); err != nil {
				return err
			}

			backlog, err := svc.GetBacklog(input.ID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Backlog created: %s\n", input.ID)
			fmt.Fprintf(out, "  Name:        %s\n", backlog.Name)
			fmt.Fprintf(out, "  Goal:        %s\n", backlog.Goal)
			fmt.Fprintf(out, "  Task Dep:    %s\n", taskDep)
			fmt.Fprintf(out, "  Feature Dep: %s\n", featureDep)
			fmt.Fprintf(out, "  Issue Dep:   %s\n", issueDep)
			fmt.Fprintf(out, "  Strict:      %t\n", strict)
			fmt.Fprintf(out, "  Location:    .mandor/backlogs/%s/\n", input.ID)
			fmt.Fprintf(out, "  Created:     %s\n", backlog.CreatedAt.Format("2006-01-02T15:04:05Z"))
			fmt.Fprintf(out, "  Creator:     %s\n", backlog.CreatedBy)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Next steps:")
			fmt.Fprintln(out, "  1. Add features: mandor feature create \"Feature Name\" --backlog "+input.ID)
			fmt.Fprintln(out, "  2. Add tasks: mandor task create \"Task Name\" --backlog "+input.ID)
			fmt.Fprintln(out, "  3. Add issues: mandor issue create \"Issue Name\" --backlog "+input.ID)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Backlog display name")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Backlog goal/objectives (required, min 500 characters)")
	cmd.Flags().StringVar(&taskDep, "task-dep", "same_backlog_only", "Task dependency rule (same_backlog_only, cross_backlog_allowed, disabled)")
	cmd.Flags().StringVar(&featureDep, "feature-dep", "cross_backlog_allowed", "Feature dependency rule (same_backlog_only, cross_backlog_allowed, disabled)")
	cmd.Flags().StringVar(&issueDep, "issue-dep", "same_backlog_only", "Issue dependency rule (same_backlog_only, cross_backlog_allowed, disabled)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Enforce strict dependency rules")
	cmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Non-interactive mode")

	return cmd
}

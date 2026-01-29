package task

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	createFeatureID string
	createGoal      string
	createImplSteps string
	createTestCases string
	createDerivable string
	createLibraries string
	createPriority  string
	createDependsOn string
	createYes       bool
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name> --feature <id> --goal <text> --implementation-steps <steps> --test-cases <cases> --derivable-files <files> --library-needs <libs> [--priority <P0-P5>] [--depends-on <ids>] [-y]",
		Short: "Create a new task",
		Long:  "Create a new task in the specified feature with the given details.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if createFeatureID == "" {
				return domain.NewValidationError("Feature ID is required (--feature).")
			}

			if createGoal == "" {
				return domain.NewValidationError("Task goal is required (--goal).")
			}

			if len(createGoal) > 500 {
				return domain.NewValidationError("Goal must be 500 characters or less.")
			}

			implSteps := splitByComma(createImplSteps)
			if len(implSteps) == 0 || (len(implSteps) == 1 && implSteps[0] == "") {
				return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
			}

			testCases := splitByComma(createTestCases)
			if len(testCases) == 0 || (len(testCases) == 1 && testCases[0] == "") {
				return domain.NewValidationError("Test cases are required (--test-cases).")
			}

			derivableFiles := splitByComma(createDerivable)
			if len(derivableFiles) == 0 || (len(derivableFiles) == 1 && derivableFiles[0] == "") {
				return domain.NewValidationError("Derivable files are required (--derivable-files).")
			}

			libraries := splitByComma(createLibraries)
			if len(libraries) == 0 || (len(libraries) == 1 && libraries[0] == "") {
				return domain.NewValidationError("Library needs are required (--library-needs).")
			}

			var dependsOnList []string
			if createDependsOn != "" {
				dependsOnList = splitByComma(createDependsOn)
			}

			input := &domain.TaskCreateInput{
				FeatureID:           createFeatureID,
				Name:                args[0],
				Goal:                createGoal,
				ImplementationSteps: implSteps,
				TestCases:           testCases,
				DerivableFiles:      derivableFiles,
				LibraryNeeds:        libraries,
				Priority:            createPriority,
				DependsOn:           dependsOnList,
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			task, err := svc.CreateTask(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Task created: %s\n", task.ID)
			fmt.Fprintf(out, "  Name:               %s\n", task.Name)
			fmt.Fprintf(out, "  Feature:            %s\n", task.FeatureID)
			fmt.Fprintf(out, "  Priority:           %s\n", task.Priority)
			fmt.Fprintf(out, "  Status:             %s\n", task.Status)
			fmt.Fprintf(out, "  Goal:               %s\n", truncate(task.Goal, 50))
			fmt.Fprintf(out, "  Implementation Steps: %d\n", len(task.ImplementationSteps))
			fmt.Fprintf(out, "  Test Cases:         %d\n", len(task.TestCases))
			fmt.Fprintf(out, "  Derivable Files:    %d\n", len(task.DerivableFiles))
			fmt.Fprintf(out, "  Library Needs:      %d\n", len(task.LibraryNeeds))
			if len(task.DependsOn) > 0 {
				fmt.Fprintf(out, "  Depends on:         %d task(s)\n", len(task.DependsOn))
			}

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&createFeatureID, "feature", "f", "", "Feature ID (required)")
	cmd.Flags().StringVarP(&createGoal, "goal", "g", "", "Task goal (required, max 500 chars)")
	cmd.Flags().StringVar(&createImplSteps, "implementation-steps", "", "Implementation steps (comma-separated, required)")
	cmd.Flags().StringVar(&createTestCases, "test-cases", "", "Test cases (comma-separated, required)")
	cmd.Flags().StringVar(&createDerivable, "derivable-files", "", "Derivable files (comma-separated, required)")
	cmd.Flags().StringVar(&createLibraries, "library-needs", "", "Required libraries (comma-separated, required)")
	cmd.Flags().StringVar(&createPriority, "priority", "P3", "Priority (P0-P5)")
	cmd.Flags().StringVar(&createDependsOn, "depends-on", "", "Comma-separated task IDs this task depends on")
	cmd.Flags().BoolVarP(&createYes, "yes", "y", false, "Skip confirmation prompts")

	return cmd
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

package task

import (
	"fmt"
	"strings"

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
	createLibraries string
	createPriority  string
	createDependsOn string
	createYes       bool
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <feature_id> <name> --goal <text> --implementation-steps <steps> --test-cases <cases> --library-needs <libs> [--priority <P0-P5>] [--depends-on <ids>] [-y]",
		Short: "Create a new task",
		Long:  "Create a new task in the specified feature with the given details.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			createFeatureID = args[0]
			taskName := args[1]

			if createGoal == "" {
				return domain.NewValidationError("Task goal is required (--goal).")
			}

			implSteps := splitByPipe(createImplSteps)
			if len(implSteps) == 0 || (len(implSteps) == 1 && implSteps[0] == "") {
				return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
			}

			testCases := splitByPipe(createTestCases)
			if len(testCases) == 0 || (len(testCases) == 1 && testCases[0] == "") {
				return domain.NewValidationError("Test cases are required (--test-cases).")
			}

			if createLibraries == "" {
				return domain.NewValidationError("Library needs are required (--library-needs).")
			}

			libraries := splitByPipe(createLibraries)

			var dependsOnList []string
			if createDependsOn != "" {
				dependsOnList = splitByPipe(createDependsOn)
			}

			input := &domain.TaskCreateInput{
				FeatureID:           createFeatureID,
				Name:                taskName,
				Goal:                createGoal,
				ImplementationSteps: implSteps,
				TestCases:           testCases,
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

	cmd.Flags().StringVarP(&createGoal, "goal", "g", "", "Task goal (required, min 500 chars)")
	cmd.Flags().StringVar(&createImplSteps, "implementation-steps", "", "Implementation steps (pipe-separated, required)")
	cmd.Flags().StringVar(&createTestCases, "test-cases", "", "Test cases (pipe-separated, required)")
	cmd.Flags().StringVar(&createLibraries, "library-needs", "", "Required libraries (pipe-separated, required). Use \"none\" if no external libraries are needed.")
	cmd.Flags().StringVar(&createPriority, "priority", "", "Priority (P0-P5, default from config)")
	cmd.Flags().StringVar(&createDependsOn, "depends-on", "", "Pipe-separated task IDs this task depends on")
	cmd.Flags().BoolVarP(&createYes, "yes", "y", false, "Skip confirmation prompts")

	return cmd
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func splitByPipe(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
			result = append(result, strings.TrimSpace(s[start:i]))
			start = i + 1
		}
	}
	result = append(result, strings.TrimSpace(s[start:]))
	return result
}

package task

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)



func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <feature_id> <name> --spec-id <spec-id> --iae-scenarios <req-XXXX:scenario-YYYY>|<...> --goal <text> --implementation-steps <steps> --test-cases <cases> [--library-needs <libs>] [--priority <P0-P5>] [--depends-on <ids>] [-y]",
		Short: "Create a new task",
		Long:  "Create a new task with IAE scenarios. Spec ID must match Feature's spec_id. IAE scenarios are pipe-separated req-XXXX:scenario-YYYY references.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			featureID := args[0]
			taskName := args[1]

			// Get flags
			specID, _ := cmd.Flags().GetString("spec-id")
			if specID == "" {
				return domain.NewValidationError("Spec ID is required (--spec-id).")
			}

			iaeStr, _ := cmd.Flags().GetString("iae-scenarios")
			if iaeStr == "" {
				return domain.NewValidationError("IAE scenarios are required (--iae-scenarios, format: req-XXXX:scenario-YYYY or req-X:s1|req-X:s2).")
			}

			goalText, _ := cmd.Flags().GetString("goal")
			if goalText == "" {
				return domain.NewValidationError("Task goal is required (--goal).")
			}

			implStr, _ := cmd.Flags().GetString("implementation-steps")
			testStr, _ := cmd.Flags().GetString("test-cases")
			libStr, _ := cmd.Flags().GetString("library-needs")
			priorityStr, _ := cmd.Flags().GetString("priority")
			dependsOnStr, _ := cmd.Flags().GetString("depends-on")

			// Parse IAE scenarios (pipe-separated)
			iaeScenarios := splitByPipe(iaeStr)
			if len(iaeScenarios) == 0 || (len(iaeScenarios) == 1 && iaeScenarios[0] == "") {
				return domain.NewValidationError("IAE scenarios are required (--iae-scenarios).")
			}

			implSteps := splitByPipe(implStr)
			if len(implSteps) == 0 || (len(implSteps) == 1 && implSteps[0] == "") {
				return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
			}

			testCases := splitByPipe(testStr)
			if len(testCases) == 0 || (len(testCases) == 1 && testCases[0] == "") {
				return domain.NewValidationError("Test cases are required (--test-cases).")
			}

			libraries := splitByPipe(libStr)

			var dependsOnList []string
			if dependsOnStr != "" {
				dependsOnList = splitByPipe(dependsOnStr)
			}

			input := &domain.TaskCreateInput{
				FeatureID:           featureID,
				SpecID:              specID,
				Name:                taskName,
				Goal:                goalText,
				IAEScenarios:        iaeScenarios,
				ImplementationSteps: implSteps,
				TestCases:           testCases,
				LibraryNeeds:        libraries,
				Priority:            priorityStr,
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
			fmt.Fprintf(out, "✓ Task created: %s\n", task.ID)
			fmt.Fprintf(out, "  Name:                 %s\n", task.Name)
			fmt.Fprintf(out, "  Feature:              %s\n", task.FeatureID)
			fmt.Fprintf(out, "  Spec ID:              %s\n", task.SpecID)
			fmt.Fprintf(out, "  Priority:             %s\n", task.Priority)
			fmt.Fprintf(out, "  Status:               %s\n", task.Status)
			fmt.Fprintf(out, "  Goal:                 %s\n", truncate(task.Goal, 50))
			fmt.Fprintf(out, "  IAE Scenarios:        %d\n", len(task.IAEScenarios))
			fmt.Fprintf(out, "  Implementation Steps: %d\n", len(task.ImplementationSteps))
			fmt.Fprintf(out, "  Test Cases:           %d\n", len(task.TestCases))
			fmt.Fprintf(out, "  Library Needs:        %d\n", len(task.LibraryNeeds))
			fmt.Fprintf(out, "  Read Gates:           brief=%v, spec=%v, notes=%v\n", task.ReadGates.IsReadBrief, task.ReadGates.IsReadSpec, task.ReadGates.IsReadSessionNotes)
			if len(task.DependsOn) > 0 {
				fmt.Fprintf(out, "  Depends on:           %d task(s)\n", len(task.DependsOn))
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

	cmd.Flags().String("spec-id", "", "Spec ID (required, must match Feature's spec_id)")
	cmd.Flags().String("iae-scenarios", "", "IAE scenarios (required, format: req-XXXX:scenario-YYYY or req-X:s1|req-X:s2)")
	cmd.Flags().StringP("goal", "g", "", "Task goal (required, min 500 chars)")
	cmd.Flags().String("implementation-steps", "", "Implementation steps (pipe-separated, required)")
	cmd.Flags().String("test-cases", "", "Test cases (pipe-separated, required)")
	cmd.Flags().String("library-needs", "", "Required libraries (pipe-separated, optional). Use \"none\" if no external libraries are needed.")
	cmd.Flags().String("priority", "", "Priority (P0-P5, default from config)")
	cmd.Flags().String("depends-on", "", "Pipe-separated task IDs this task depends on")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	cmd.MarkFlagRequired("spec-id")
	cmd.MarkFlagRequired("iae-scenarios")
	cmd.MarkFlagRequired("goal")
	cmd.MarkFlagRequired("implementation-steps")
	cmd.MarkFlagRequired("test-cases")

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

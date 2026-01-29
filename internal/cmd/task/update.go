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
	updateTaskID        string
	updateName          string
	updateGoal          string
	updatePriority      string
	updateImplSteps     string
	updateTestCases     string
	updateDerivable     string
	updateLibraries     string
	updateStatus        string
	updateReason        string
	updateDependsOn     string
	updateDependsAdd    string
	updateDependsRemove string
	updateReopen        bool
	updateCancel        bool
	updateForce         bool
	updateDryRun        bool
	updateYes           bool
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <task_id> [--name] [--priority] [--goal] [--implementation-steps] [--test-cases] [--derivable-files] [--library-needs] [--status <ready|in_progress|done>] [--cancel --reason] [--reopen] [--depends <ids>]",
		Short: "Update a task",
		Long:  "Update task properties, change status, cancel, or reopen.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			taskID := args[0]

			var namePtr, goalPtr, priorityPtr, statusPtr, reasonPtr *string
			var implStepsPtr, testCasesPtr, derivablePtr, librariesPtr *[]string
			var dependsOnPtr, dependsAddPtr, dependsRemovePtr *[]string

			if updateName != "" {
				namePtr = &updateName
			}
			if updateGoal != "" {
				goalPtr = &updateGoal
			}
			if updatePriority != "" {
				if !domain.ValidatePriority(updatePriority) {
					return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
				}
				priorityPtr = &updatePriority
			}
			if updateStatus != "" {
				if !domain.ValidateTaskStatus(updateStatus) {
					return domain.NewValidationError("Invalid status. Valid options: pending, ready, in_progress, blocked, done, cancelled")
				}
				statusPtr = &updateStatus
			}
			if updateReason != "" {
				reasonPtr = &updateReason
			}

			if updateImplSteps != "" {
				steps := splitByComma(updateImplSteps)
				implStepsPtr = &steps
			}
			if updateTestCases != "" {
				cases := splitByComma(updateTestCases)
				testCasesPtr = &cases
			}
			if updateDerivable != "" {
				files := splitByComma(updateDerivable)
				derivablePtr = &files
			}
			if updateLibraries != "" {
				libs := splitByComma(updateLibraries)
				librariesPtr = &libs
			}

			if updateDependsOn != "" {
				deps := splitByComma(updateDependsOn)
				dependsOnPtr = &deps
			}
			if updateDependsAdd != "" {
				deps := splitByComma(updateDependsAdd)
				dependsAddPtr = &deps
			}
			if updateDependsRemove != "" {
				deps := splitByComma(updateDependsRemove)
				dependsRemovePtr = &deps
			}

			input := &domain.TaskUpdateInput{
				TaskID:              taskID,
				Name:                namePtr,
				Goal:                goalPtr,
				Priority:            priorityPtr,
				ImplementationSteps: implStepsPtr,
				TestCases:           testCasesPtr,
				DerivableFiles:      derivablePtr,
				LibraryNeeds:        librariesPtr,
				Status:              statusPtr,
				Reason:              reasonPtr,
				DependsOn:           dependsOnPtr,
				DependsAdd:          dependsAddPtr,
				DependsRemove:       dependsRemovePtr,
				Reopen:              updateReopen,
				Cancel:              updateCancel,
				Force:               updateForce,
				DryRun:              updateDryRun,
			}

			if err := svc.ValidateUpdateInput(input); err != nil {
				return err
			}

			changes, err := svc.UpdateTask(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if updateDryRun {
				fmt.Fprintln(out, "[DRY RUN] Changes:")
			} else {
				fmt.Fprintf(out, "Task updated: %s\n", taskID)
			}
			for _, change := range changes {
				fmt.Fprintf(out, "  - %s\n", change)
			}

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" && !updateDryRun {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&updateName, "name", "", "New task name")
	cmd.Flags().StringVar(&updateGoal, "goal", "", "New task goal")
	cmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (P0-P5)")
	cmd.Flags().StringVar(&updateImplSteps, "implementation-steps", "", "Update implementation steps (comma-separated)")
	cmd.Flags().StringVar(&updateTestCases, "test-cases", "", "Update test cases (comma-separated)")
	cmd.Flags().StringVar(&updateDerivable, "derivable-files", "", "Update derivable files (comma-separated)")
	cmd.Flags().StringVar(&updateLibraries, "library-needs", "", "Update library needs (comma-separated)")
	cmd.Flags().StringVar(&updateStatus, "status", "", "New status (ready, in_progress, done)")
	cmd.Flags().StringVar(&updateReason, "reason", "", "Cancellation reason (required with --cancel)")
	cmd.Flags().StringVar(&updateDependsOn, "depends", "", "Set all dependencies (comma-separated)")
	cmd.Flags().StringVar(&updateDependsAdd, "depends-add", "", "Add dependencies (comma-separated)")
	cmd.Flags().StringVar(&updateDependsRemove, "depends-remove", "", "Remove dependencies (comma-separated)")
	cmd.Flags().BoolVar(&updateReopen, "reopen", false, "Reopen a cancelled task")
	cmd.Flags().BoolVar(&updateCancel, "cancel", false, "Cancel the task")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Force operation (e.g., cancel with dependents)")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
	cmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Skip confirmation")

	return cmd
}

func splitByComma(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, strings.TrimSpace(s[start:i]))
			start = i + 1
		}
	}
	result = append(result, strings.TrimSpace(s[start:]))
	return result
}

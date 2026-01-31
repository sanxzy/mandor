package task

import (
	"fmt"

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
				steps := splitByPipe(updateImplSteps)
				implStepsPtr = &steps
			}
			if updateTestCases != "" {
				cases := splitByPipe(updateTestCases)
				testCasesPtr = &cases
			}
			if updateDerivable != "" {
				files := splitByPipe(updateDerivable)
				derivablePtr = &files
			}
			if updateLibraries != "" {
				libs := splitByPipe(updateLibraries)
				librariesPtr = &libs
			}

			if updateDependsOn != "" {
				deps := splitByPipe(updateDependsOn)
				dependsOnPtr = &deps
			}
			if updateDependsAdd != "" {
				deps := splitByPipe(updateDependsAdd)
				dependsAddPtr = &deps
			}
			if updateDependsRemove != "" {
				deps := splitByPipe(updateDependsRemove)
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
				for _, change := range changes {
					fmt.Fprintf(out, "  - %s\n", change)
				}
				return nil
			}

			fmt.Fprintf(out, "Task updated: %s\n", taskID)
			for _, change := range changes {
				fmt.Fprintf(out, "  - %s\n", change)
			}

			// Show task detail unless status is "done"
			if statusPtr == nil || *statusPtr != domain.TaskStatusDone {
				fmt.Fprintln(out)
				detailInput := &domain.TaskDetailInput{
					TaskID:         taskID,
					JSON:           false,
					IncludeDeleted: false,
					Events:         false,
					Dependencies:   false,
					Timestamps:     false,
				}

				detailOutput, err := svc.GetTaskDetail(detailInput)
				if err == nil {
					fmt.Fprintf(out, "Task: %s\n", detailOutput.ID)
					fmt.Fprintf(out, "  Name:               %s\n", detailOutput.Name)
					fmt.Fprintf(out, "  Feature:            %s\n", detailOutput.FeatureID)
					fmt.Fprintf(out, "  Project:            %s\n", detailOutput.ProjectID)
					fmt.Fprintf(out, "  Status:             %s\n", detailOutput.Status)
					fmt.Fprintf(out, "  Priority:           %s\n", detailOutput.Priority)
					fmt.Fprintf(out, "  Goal:               %s\n", detailOutput.Goal)
					fmt.Fprintf(out, "  Implementation Steps (%d):\n", len(detailOutput.ImplementationSteps))
					for i, step := range detailOutput.ImplementationSteps {
						fmt.Fprintf(out, "    %d. %s\n", i+1, step)
					}
					fmt.Fprintf(out, "  Test Cases (%d):\n", len(detailOutput.TestCases))
					for i, tc := range detailOutput.TestCases {
						fmt.Fprintf(out, "    %d. %s\n", i+1, tc)
					}
					fmt.Fprintf(out, "  Derivable Files (%d):\n", len(detailOutput.DerivableFiles))
					for _, f := range detailOutput.DerivableFiles {
						fmt.Fprintf(out, "    - %s\n", f)
					}
					fmt.Fprintf(out, "  Library Needs (%d):\n", len(detailOutput.LibraryNeeds))
					for _, lib := range detailOutput.LibraryNeeds {
						fmt.Fprintf(out, "    - %s\n", lib)
					}
					if len(detailOutput.DependsOn) > 0 {
						fmt.Fprintf(out, "  Depends On (%d):\n", len(detailOutput.DependsOn))
						for _, dep := range detailOutput.DependsOn {
							fmt.Fprintf(out, "    - %s\n", dep)
						}
					}
					fmt.Fprintf(out, "  Created:   %s\n", detailOutput.CreatedAt)
					fmt.Fprintf(out, "  Updated:   %s\n", detailOutput.UpdatedAt)
					fmt.Fprintf(out, "  CreatedBy: %s\n", detailOutput.CreatedBy)
					fmt.Fprintf(out, "  UpdatedBy: %s\n", detailOutput.UpdatedBy)
					fmt.Fprintf(out, "  Events:    %d\n", detailOutput.Events)
				}
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

	cmd.Flags().StringVar(&updateName, "name", "", "New task name")
	cmd.Flags().StringVar(&updateGoal, "goal", "", "New task goal")
	cmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (P0-P5)")
	cmd.Flags().StringVar(&updateImplSteps, "implementation-steps", "", "Update implementation steps (pipe-separated)")
	cmd.Flags().StringVar(&updateTestCases, "test-cases", "", "Update test cases (pipe-separated)")
	cmd.Flags().StringVar(&updateDerivable, "derivable-files", "", "Update derivable files (pipe-separated)")
	cmd.Flags().StringVar(&updateLibraries, "library-needs", "", "Update library needs (pipe-separated)")
	cmd.Flags().StringVar(&updateStatus, "status", "", "New status (ready, in_progress, done)")
	cmd.Flags().StringVar(&updateReason, "reason", "", "Cancellation reason (required with --cancel)")
	cmd.Flags().StringVar(&updateDependsOn, "depends", "", "Set all dependencies (pipe-separated)")
	cmd.Flags().StringVar(&updateDependsAdd, "depends-add", "", "Add dependencies (pipe-separated)")
	cmd.Flags().StringVar(&updateDependsRemove, "depends-remove", "", "Remove dependencies (pipe-separated)")
	cmd.Flags().BoolVar(&updateReopen, "reopen", false, "Reopen a cancelled task")
	cmd.Flags().BoolVar(&updateCancel, "cancel", false, "Cancel the task")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Force operation (e.g., cancel with dependents)")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
	cmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Skip confirmation")

	return cmd
}

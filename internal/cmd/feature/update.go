package feature

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	updateProjectID string
	updateName      string
	updateGoal      string
	updateScope     string
	updatePriority  string
	updateStatus    string
	updateReason    string
	updateDependsOn string
	updateReopen    bool
	updateCancel    bool
	updateForce     bool
	updateDryRun    bool
	updateYes       bool
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <feature_id> [--project <id>] [--name] [--goal] [--scope] [--priority] [--status] [--cancel --reason] [--reopen] [--depends]",
		Short: "Update a feature",
		Long:  "Update feature properties, change status, cancel, or reopen.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := updateProjectID
			if projectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			featureID := args[0]

			var dependsOnList *[]string
			if updateDependsOn != "" {
				list := splitDependsOn(updateDependsOn)
				dependsOnList = &list
			}

			var namePtr, goalPtr, scopePtr, priorityPtr, statusPtr, reasonPtr *string
			if updateName != "" {
				namePtr = &updateName
			}
			if updateGoal != "" {
				goalPtr = &updateGoal
			}
			if updateScope != "" {
				scopePtr = &updateScope
			}
			if updatePriority != "" {
				priorityPtr = &updatePriority
			}
			if updateStatus != "" {
				statusPtr = &updateStatus
			}
			if updateReason != "" {
				reasonPtr = &updateReason
			}

			input := &domain.FeatureUpdateInput{
				ProjectID: projectID,
				FeatureID: featureID,
				Name:      namePtr,
				Goal:      goalPtr,
				Scope:     scopePtr,
				Priority:  priorityPtr,
				Status:    statusPtr,
				Reason:    reasonPtr,
				DependsOn: dependsOnList,
				Reopen:    updateReopen,
				Cancel:    updateCancel,
				Force:     updateForce,
				DryRun:    updateDryRun,
			}

			if err := svc.ValidateUpdateInput(input); err != nil {
				return err
			}

			changes, err := svc.UpdateFeature(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if updateDryRun {
				fmt.Fprintln(out, "[DRY RUN] Changes:")
			} else {
				fmt.Fprintln(out, "Feature updated:", featureID)
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

	cmd.Flags().StringVarP(&updateProjectID, "project", "p", "", "Project ID (required)")
	cmd.Flags().StringVar(&updateName, "name", "", "New feature name")
	cmd.Flags().StringVar(&updateGoal, "goal", "", "New feature goal")
	cmd.Flags().StringVar(&updateScope, "scope", "", "New scope (frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift)")
	cmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (P0-P5)")
	cmd.Flags().StringVar(&updateStatus, "status", "", "New status (draft, active, done, blocked, cancelled)")
	cmd.Flags().StringVar(&updateReason, "reason", "", "Cancellation reason (required with --cancel)")
	cmd.Flags().StringVar(&updateDependsOn, "depends", "", "Comma-separated feature IDs this feature depends on")
	cmd.Flags().BoolVar(&updateReopen, "reopen", false, "Reopen a cancelled feature")
	cmd.Flags().BoolVar(&updateCancel, "cancel", false, "Cancel the feature")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Force operation (e.g., cancel with dependents)")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
	cmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Skip confirmation")

	return cmd
}

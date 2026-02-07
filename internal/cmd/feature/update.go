package feature

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	updateBacklogID string
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
		Use:   "update <feature_id> [--backlog <id>] [--name] [--goal] [--scope] [--priority] [--status] [--cancel --reason] [--reopen] [--depends]",
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

			backlogID := updateBacklogID
			if backlogID == "" {
				return domain.NewValidationError("Backlog ID is required (--backlog).")
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
				BacklogID: backlogID,
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
				for _, change := range changes {
					fmt.Fprintf(out, "  - %s\n", change)
				}
				return nil
			}

			fmt.Fprintf(out, "✓ Feature updated: %s\n", featureID)
			if len(changes) > 0 {
				fmt.Fprintln(out, "  Changes:")
				for _, change := range changes {
					fmt.Fprintf(out, "    - %s\n", change)
				}
			}

			// Show full feature state
			fmt.Fprintln(out)
			detailInput := &domain.FeatureDetailInput{
				BacklogID:      backlogID,
				FeatureID:      featureID,
				JSON:           false,
				IncludeDeleted: false,
			}

			detail, err := svc.GetFeatureDetail(detailInput)
			if err == nil {
				fmt.Fprintf(out, "FEATURE: %s", detail.ID)
				if detail.Status == domain.FeatureStatusCancelled {
					fmt.Fprint(out, " [CANCELLED]")
				}
				fmt.Fprintln(out)
				fmt.Fprintln(out, "------------------------------------------------------------")

				fmt.Fprintf(out, "  Name:         %s\n", detail.Name)
				fmt.Fprintf(out, "  Status:       %s\n", detail.Status)
				fmt.Fprintf(out, "  Priority:     %s\n", detail.Priority)
				if detail.Scope != "" {
					fmt.Fprintf(out, "  Scope:        %s\n", detail.Scope)
				}
				fmt.Fprintf(out, "  Capability:   %s\n", detail.CapabilityID)
				fmt.Fprintf(out, "  Spec:         %s\n", detail.SpecID)

				if detail.Goal != "" {
					goalPreview := detail.Goal
					if len(goalPreview) > 200 {
						goalPreview = goalPreview[:197] + "..."
					}
					fmt.Fprintf(out, "\n  Goal:\n    %s\n", goalPreview)
				}

				if len(detail.DependsOn) > 0 {
					fmt.Fprintf(out, "\n  Dependencies (%d):\n", len(detail.DependsOn))
					for _, depID := range detail.DependsOn {
						fmt.Fprintf(out, "    - %s\n", depID)
					}
				}

				if detail.Reason != "" && detail.Status == domain.FeatureStatusCancelled {
					fmt.Fprintf(out, "\n  Cancellation Reason: %s\n", detail.Reason)
				}

				fmt.Fprintf(out, "\n  Created: %s by %s\n", detail.CreatedAt, detail.CreatedBy)
				fmt.Fprintf(out, "  Updated: %s by %s\n", detail.UpdatedAt, detail.UpdatedBy)

				// AI-friendly recommendations
				fmt.Fprintln(out)
				fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
				fmt.Fprintln(out, "  RECOMMENDATIONS:")
				fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")

				switch detail.Status {
				case domain.FeatureStatusDraft:
					fmt.Fprintf(out, "  • Ready to start? Run: mandor feature update %s --status active\n", featureID)
					fmt.Fprintf(out, "  • Add dependencies? Run: mandor feature update %s --depends <feature-id>\n", featureID)
				case domain.FeatureStatusActive:
					fmt.Fprintf(out, "  • Mark complete? Run: mandor feature update %s --status done\n", featureID)
					fmt.Fprintf(out, "  • Blocked? Run: mandor feature update %s --status blocked\n", featureID)
				case domain.FeatureStatusBlocked:
					fmt.Fprintf(out, "  • Unblocked? Resume with: mandor feature update %s --status active\n", featureID)
				case domain.FeatureStatusDone:
					fmt.Fprintf(out, "  • Feature complete. Check for dependent features that may now be unblocked.\n")
					fmt.Fprintf(out, "  • View dependents: mandor feature list --backlog %s\n", backlogID)
				}

				fmt.Fprintf(out, "  • View full details: mandor feature detail %s --backlog %s\n", featureID, backlogID)
				fmt.Fprintf(out, "  • List all features: mandor feature list --backlog %s\n", backlogID)
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

	cmd.Flags().StringVarP(&updateBacklogID, "backlog", "b", "", "Backlog ID (required)")
	cmd.Flags().StringVar(&updateName, "name", "", "New feature name")
	cmd.Flags().StringVar(&updateGoal, "goal", "", "New feature goal")
	cmd.Flags().StringVar(&updateScope, "scope", "", "New scope (frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift)")
	cmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (P0-P5)")
	cmd.Flags().StringVar(&updateStatus, "status", "", "New status (draft, active, done, blocked, cancelled)")
	cmd.Flags().StringVar(&updateReason, "reason", "", "Cancellation reason (required with --cancel)")
	cmd.Flags().StringVar(&updateDependsOn, "depends", "", "Pipe-separated feature IDs this feature depends on")
	cmd.Flags().BoolVar(&updateReopen, "reopen", false, "Reopen a cancelled feature")
	cmd.Flags().BoolVar(&updateCancel, "cancel", false, "Cancel the feature")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Force operation (e.g., cancel with dependents)")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
	cmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Skip confirmation")

	return cmd
}

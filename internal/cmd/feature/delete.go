package feature

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	deleteProjectID string
	deleteForce     bool
	deleteReason    string
	deleteYes       bool
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <feature_id> [--project <id>] [--reason <text>]",
		Short: "Delete a feature",
		Long:  "Delete a feature from a project. Requires reason unless --force is used.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := deleteProjectID
			if projectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			featureID := args[0]

			// Confirm deletion unless --yes is provided
			if !deleteYes && !deleteForce {
				fmt.Fprintf(cmd.OutOrStdout(), "Delete feature %s? This action cannot be undone. Type 'yes' to confirm: ", featureID)
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" {
					fmt.Fprintln(cmd.OutOrStdout(), "Deletion cancelled.")
					return nil
				}
			}

			input := &domain.FeatureDeleteInput{
				ProjectID: projectID,
				FeatureID: featureID,
				Force:     deleteForce,
				Reason:    deleteReason,
			}

			if err := svc.ValidateDeleteInput(input); err != nil {
				return err
			}

			if err := svc.DeleteFeature(input); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Feature deleted: %s\n", featureID)

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&deleteProjectID, "project", "p", "", "Project ID (required)")
	cmd.Flags().StringVar(&deleteReason, "reason", "", "Reason for deletion")
	cmd.Flags().BoolVar(&deleteForce, "force", false, "Force deletion without reason")
	cmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")

	return cmd
}

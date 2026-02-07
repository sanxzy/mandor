package spec

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <spec-id>",
	Short: "Delete a Spec",
	Long:  "Archive or permanently delete a Spec",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var (
	deleteHard bool
	deleteYes  bool
)

func init() {
	deleteCmd.Flags().BoolVar(&deleteHard, "hard", false, "Permanently delete instead of archiving")
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
}

func runDelete(cmd *cobra.Command, args []string) error {
	specID := args[0]

	var backlogID string
	for p := cmd; p != nil; p = p.Parent() {
		if val, err := p.Flags().GetString("backlog"); err == nil && val != "" {
			backlogID = val
			break
		}
	}
	if backlogID == "" {
		return domain.NewValidationError("Backlog ID is required (--backlog).")
	}

	svc, err := service.NewSpecService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Load spec
	spec, err := svc.ReadSpec(backlogID, specID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	if !deleteYes {
		return domain.NewValidationError("Confirmation required. Use --yes to confirm deletion.")
	}

	if deleteHard {
		if err := svc.DeleteSpec(backlogID, specID); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Spec permanently deleted: %s\n", specID)
	} else {
		// Soft delete - mark as archived
		spec.Status = domain.SpecStatusArchived
		if err := svc.UpdateSpec(backlogID, spec); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Spec archived: %s\n", specID)
	}

	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}

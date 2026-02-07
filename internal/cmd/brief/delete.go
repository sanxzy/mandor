package brief

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <brief-id>",
	Short: "Delete a Brief",
	Long:  "Archive or permanently delete a Brief",
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

	svc, err := service.NewBriefService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Load brief
	brief, err := svc.ReadBrief(backlogID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	if !deleteYes {
		return domain.NewValidationError("Confirmation required. Use --yes to confirm deletion.")
	}

	if deleteHard {
		if err := svc.DeleteBrief(backlogID); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Brief permanently deleted: %s\n", backlogID)
	} else {
		// Soft delete - mark as archived
		brief.Status = domain.BriefStatusArchived
		if err := svc.UpdateBrief(backlogID, brief); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Brief archived: %s\n", backlogID)
	}

	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}

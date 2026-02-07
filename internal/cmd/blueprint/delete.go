package blueprint

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <blueprint-id>",
	Short: "Delete a Blueprint",
	Long:  "Archive or permanently delete a Blueprint",
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
	blueprintID := args[0]

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

	svc, err := service.NewBlueprintService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Load blueprint
	blueprint, err := svc.ReadBlueprint(backlogID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	if !deleteYes {
		return domain.NewValidationError("Confirmation required. Use --yes to confirm deletion.")
	}

	if deleteHard {
		if err := svc.DeleteBlueprint(backlogID); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Blueprint permanently deleted: %s\n", blueprintID)
	} else {
		// Soft delete - mark as archived
		blueprint.Status = domain.BlueprintStatusArchived
		if err := svc.UpdateBlueprint(backlogID, blueprint); err != nil {
			return err
		}
		fmt.Fprintf(out, "✓ Blueprint archived: %s\n", blueprintID)
	}

	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}

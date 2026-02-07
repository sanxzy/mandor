package blueprint

import (
	"github.com/spf13/cobra"
)

var (
	backlogID string
)

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Manage Blueprints (technical architecture)",
	Long:  "Create, read, update, and manage Blueprints - technical architecture documents linking Briefs and Specs",
}

func init() {
	blueprintCmd.PersistentFlags().StringVarP(&backlogID, "backlog", "b", "", "Backlog ID (required)")

	blueprintCmd.AddCommand(GetCreateCmd())
	blueprintCmd.AddCommand(GetDetailCmd())
	blueprintCmd.AddCommand(GetUpdateCmd())
	blueprintCmd.AddCommand(GetDeleteCmd())
	blueprintCmd.AddCommand(GetValidateCmd())
}

func GetBlueprintCmd() *cobra.Command {
	return blueprintCmd
}

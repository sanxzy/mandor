package brief

import (
	"github.com/spf13/cobra"
)

var (
	backlogID string
)

var briefCmd = &cobra.Command{
	Use:   "brief",
	Short: "Manage Briefs (backlog intent documents)",
	Long:  "Create, read, update, and manage Briefs - the root specification documents for backlogs",
}

func init() {
	briefCmd.PersistentFlags().StringVarP(&backlogID, "backlog", "b", "", "Backlog ID (required)")

	briefCmd.AddCommand(GetCreateCmd())
	briefCmd.AddCommand(GetReadCmd())
	briefCmd.AddCommand(GetUpdateCmd())
	briefCmd.AddCommand(GetDeleteCmd())
	briefCmd.AddCommand(GetValidateCmd())
}

func GetBriefCmd() *cobra.Command {
	return briefCmd
}

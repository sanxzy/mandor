package brief

import (
	"github.com/spf13/cobra"
)

var (
	projectID string
)

var briefCmd = &cobra.Command{
	Use:   "brief",
	Short: "Manage Briefs (project intent documents)",
	Long:  "Create, read, update, and manage Briefs - the root specification documents for projects",
}

func init() {
	briefCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "Project ID (required)")
	
	briefCmd.AddCommand(GetCreateCmd())
	briefCmd.AddCommand(GetReadCmd())
	briefCmd.AddCommand(GetUpdateCmd())
	briefCmd.AddCommand(GetDeleteCmd())
	briefCmd.AddCommand(GetValidateCmd())
}

func GetBriefCmd() *cobra.Command {
	return briefCmd
}

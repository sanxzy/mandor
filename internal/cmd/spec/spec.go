package spec

import (
	"github.com/spf13/cobra"
)

var (
	projectID string
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Manage Specs (capability specifications)",
	Long:  "Create, read, update, and manage Specs - detailed specifications for Brief capabilities",
}

func init() {
	specCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "Project ID (required)")
	
	specCmd.AddCommand(GetCreateCmd())
	specCmd.AddCommand(GetDetailCmd())
	specCmd.AddCommand(GetUpdateCmd())
	specCmd.AddCommand(GetDeleteCmd())
	specCmd.AddCommand(GetValidateCmd())
}

func GetSpecCmd() *cobra.Command {
	return specCmd
}

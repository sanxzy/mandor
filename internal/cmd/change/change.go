package change

import (
	"github.com/spf13/cobra"
	"mandor/internal/fs"
)

// NewChangeCmd creates the change command
func NewChangeCmd() *cobra.Command {
	changeCmd := &cobra.Command{
		Use:   "change",
		Short: "Manage change governance and impact analysis",
		Long:  "Commands for analyzing, approving, and tracking changes across Briefs, Specs, and Blueprints",
	}

	paths, _ := fs.NewPaths()

	changeCmd.AddCommand(NewAnalyzeCmd(paths))
	changeCmd.AddCommand(NewApproveCmd(paths))
	changeCmd.AddCommand(NewRejectCmd(paths))
	changeCmd.AddCommand(NewListCmd(paths))

	return changeCmd
}

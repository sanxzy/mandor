package issue

import (
	"github.com/spf13/cobra"
)

func NewIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue commands",
		Long:  "Commands for managing issues within projects.",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewReadyCmd())
	cmd.AddCommand(NewBlockedCmd())

	return cmd
}

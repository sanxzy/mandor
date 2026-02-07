package backlog

import (
	"github.com/spf13/cobra"
)

func NewBacklogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlog",
		Short: "Backlog management commands",
		Long:  "Commands for managing backlogs in the workspace.",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewReopenCmd())

	return cmd
}

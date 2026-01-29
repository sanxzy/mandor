package task

import (
	"github.com/spf13/cobra"
)

func NewTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task commands",
		Long:  "Commands for managing tasks within features.",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewReadyCmd())
	cmd.AddCommand(NewBlockedCmd())

	return cmd
}

package project

import (
	"github.com/spf13/cobra"
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project management commands",
		Long:  "Commands for managing projects in the workspace.",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewReopenCmd())

	return cmd
}

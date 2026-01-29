package feature

import (
	"github.com/spf13/cobra"
)

func NewFeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature",
		Short: "Feature commands",
		Long:  "Commands for managing features within projects.",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())

	return cmd
}

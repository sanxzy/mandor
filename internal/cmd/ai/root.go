package ai

import (
	"github.com/spf13/cobra"
)

func NewAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-assisted documentation commands",
	}

	cmd.AddCommand(NewClaudeCmd())
	cmd.AddCommand(NewAgentsCmd())

	return cmd
}

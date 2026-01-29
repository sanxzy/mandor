package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/ai"
	"mandor/internal/domain"
)

func NewAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Generate AGENTS.md for multi-agent coordination",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := ai.FindProjectRoot()
			if err != nil {
				return domain.NewValidationError("failed to find project root: " + err.Error())
			}

			now := time.Now().UTC().Format("2006-01-02")
			projectName := filepath.Base(projectRoot)

			content := generateAIDoc(projectName, now)

			targetPath := filepath.Join(projectRoot, "AGENTS.md")

			if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
				return domain.NewSystemError("failed to write file", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Generated: %s\n", targetPath)
			return nil
		},
	}

	return cmd
}

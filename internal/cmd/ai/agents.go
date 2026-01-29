package ai

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
)

func NewAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Generate AGENTS.md for multi-agent coordination",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return domain.NewValidationError("failed to get current directory: " + err.Error())
			}

			now := time.Now().UTC().Format("2006-01-02")
			projectName := filepath.Base(cwd)

			content := generateAIDoc(projectName, now)

			targetPath := filepath.Join(cwd, "AGENTS.md")

			if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
				return domain.NewSystemError("failed to write file", err)
			}

			cmd.OutOrStdout().Write([]byte("Generated: " + targetPath + "\n"))
			return nil
		},
	}

	return cmd
}

func generateAIDoc(projectName, createdAt string) string {
	return "# Project Task Management\n\n" +
		"This project uses **Mandor CLI** for event-based task management.\n" +
		"All tasks, features, and issues must be tracked using Mandor.\n\n" +
		"---\n\n" +
		"## Required Tool\n\n" +
		"- Mandor CLI is mandatory for all development work.\n" +
		"- Task management outside of Mandor is not allowed.\n\n" +
		"---\n\n" +
		"## Getting Started\n\n" +
		"To view available commands and usage instructions, run:\n\n" +
		"mandor populate\n\n" +
		"---\n\n" +
		"## Critical Rules\n\n" +
		"- Always use Mandor for task creation, tracking, and updates.\n" +
		"- All development must be tied to Mandor-managed tasks.\n"
}

package ai

import (
	"os"
	"path/filepath"
	"time"

	"mandor/internal/domain"

	"github.com/spf13/cobra"
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
	return "# Mandor Essential Commands\n\n" +
		"Use these three commands to manage your work:\n\n" +
		"---\n\n" +
		"## 1. mandor populate\n\n" +
		"View all available commands and usage instructions.\n\n" +
		"```bash\n" +
		"mandor populate\n" +
		"```\n\n" +
		"Shows complete reference guide with examples for creating projects, features, tasks, " +
		"issues, and managing dependencies.\n\n" +
		"---\n\n" +
		"## 2. mandor track\n\n" +
		"Check status of workspace, projects, features, and tasks.\n\n" +
		"```bash\n" +
		"# View entire workspace\n" +
		"mandor track\n\n" +
		"# View specific project\n" +
		"mandor track project <project-id>\n\n" +
		"# View feature with tasks\n" +
		"mandor track feature <feature-id>\n\n" +
		"# View single task\n" +
		"mandor track task <task-id>\n" +
		"```\n\n" +
		"Shows what's ready, blocked, in progress, or done. Use before starting work.\n\n" +
		"---\n\n" +
		"## 3. mandor session note\n\n" +
		"Record and read session progress.\n\n" +
		"```bash\n" +
		"# Log what you completed\n" +
		"mandor session note \"Completed v0.4.4 release and testing\"\n\n" +
		"# Read last 50 notes\n" +
		"mandor session note --read\n\n" +
		"# Read more notes\n" +
		"mandor session note --read --offset 100\n" +
		"```\n\n" +
		"End each session with a note. Start next session by reading notes to resume work.\n\n" +
		"---\n\n" +
		"## Quick Workflow\n\n" +
		"1. `mandor track` - See what's ready\n" +
		"2. `mandor populate` - Learn how to create/update work\n" +
		"3. `mandor session note \"done\"` - Log progress before ending session\n" +
		"4. `mandor session note --read` - Check progress when resuming\n"
}

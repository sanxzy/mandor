package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/ai"
	"mandor/internal/ai/templates"
	"mandor/internal/domain"
)

var (
	claudeGoal     string
	claudeReplace  bool
	claudeTemplate string
)

func NewClaudeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claude",
		Short: "Generate CLAUDE.md for the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := ai.FindProjectRoot()
			if err != nil {
				return domain.NewValidationError("failed to find project root: " + err.Error())
			}

			targetPath := filepath.Join(projectRoot, "CLAUDE.md")

			if !claudeReplace {
				if _, err := os.Stat(targetPath); err == nil {
					return domain.NewValidationError("CLAUDE.md already exists. Use --replace to overwrite.")
				}
			}

			data := templates.ClaudeTemplateData{
				ProjectName:   filepath.Base(projectRoot),
				Goal:          claudeGoal,
				MandorVersion: "0.0.1",
				CreatedAt:     time.Now().UTC().Format("2006-01-02"),
			}

			content := templates.GenerateClaudeMD(data)

			tmpFile, err := os.CreateTemp("", "claude-md-*.tmp")
			if err != nil {
				return domain.NewSystemError("failed to create temp file", err)
			}
			tmpPath := tmpFile.Name()
			defer os.Remove(tmpPath)

			if _, err := tmpFile.WriteString(content); err != nil {
				tmpFile.Close()
				return domain.NewSystemError("failed to write temp file", err)
			}
			tmpFile.Close()

			if err := os.Rename(tmpPath, targetPath); err != nil {
				return domain.NewSystemError("failed to rename temp file", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Generated: %s\n", targetPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&claudeGoal, "goal", "g", "", "Project goal for documentation")
	cmd.Flags().BoolVar(&claudeReplace, "replace", false, "Replace existing CLAUDE.md")
	cmd.Flags().StringVar(&claudeTemplate, "template", "", "Custom template file")

	return cmd
}

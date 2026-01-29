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
	agentsRole     string
	agentsProtocol string
	agentsReplace  bool
	agentCount     int = 3
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

			targetPath := filepath.Join(projectRoot, "AGENTS.md")

			if !agentsReplace {
				if _, err := os.Stat(targetPath); err == nil {
					return domain.NewValidationError("AGENTS.md already exists. Use --replace to overwrite.")
				}
			}

			data := templates.AgentsTemplateData{
				ProjectName:   filepath.Base(projectRoot),
				Goal:          agentsRole,
				MandorVersion: "0.0.1",
				AgentCount:    agentCount,
				Protocol:      agentsProtocol,
				CreatedAt:     time.Now().UTC().Format("2006-01-02"),
			}

			content := templates.GenerateAgentsMD(data)

			tmpFile, err := os.CreateTemp("", "agents-md-*.tmp")
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

	cmd.Flags().StringVarP(&agentsRole, "role", "r", "", "Primary agent role or project goal")
	cmd.Flags().IntVar(&agentCount, "agents", 3, "Number of agents")
	cmd.Flags().StringVar(&agentsProtocol, "protocol", "sequential", "Coordination protocol")
	cmd.Flags().BoolVar(&agentsReplace, "replace", false, "Replace existing AGENTS.md")

	return cmd
}

package workspace

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"mandor/internal/ai"
	"mandor/internal/service"
	"mandor/internal/util"
)

func NewInitCmd() *cobra.Command {
	var (
		workspaceName string
		strict        bool
		aiAgent       string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Mandor workspace",
		Long: `Initialize a new Mandor workspace in the current directory.

Creates a .mandor/ directory with workspace metadata and backlog storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			ws, err := svc.InitWorkspace(workspaceName)
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			var agentType ai.AgentType
			if aiAgent != "" {
				agentType = mapAgentFlagToType(aiAgent)
			} else if !isInteractive() {
				agentType = ai.AgentClaude
			} else {
				agent, err := selectAIAgent()
				if err != nil {
					return err
				}
				agentType = mapAgentSelectionToType(agent)
			}

			generator := ai.NewGenerator()
			if err := generator.Write(agentType, cwd); err != nil {
				return err
			}

			fmt.Printf("Workspace initialized: %s\n", ws.Name)
			fmt.Printf("  Location: .mandor/\n")
			fmt.Printf("  ID: %s\n", ws.ID)
			fmt.Printf("  Creator: %s\n", ws.CreatedBy)
			fmt.Printf("  Created: %s\n", ws.CreatedAt.Format("2006-01-02T15:04:05Z"))

			if agentType == ai.AgentClaude {
				fmt.Printf("\nGenerated: CLAUDE.md\n")
			} else if agentType == ai.AgentGeneral {
				fmt.Printf("\nGenerated: AGENTS.md\n")
			}

			username, warning := util.GetGitUsernameWithWarning()
			if username == "unknown" {
				fmt.Printf("\nWarning: Git user not configured. Events will show 'unknown' as creator.\n")
				fmt.Printf("  Run: git config user.name \"Your Name\"\n")
			}

			fmt.Printf("\nNext steps:\n")
			fmt.Printf("  1. Create a backlog: mandor backlog create <backlog_id> --name \"Backlog Name\"\n")
			fmt.Printf("  2. View workspace: mandor track\n")
			fmt.Printf("  3. Check config: mandor config get\n")

			_ = warning
			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceName, "workspace-name", "", "", "Custom workspace name (default: current directory)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Enforce strict dependency rules (deprecated)")
	cmd.Flags().StringVar(&aiAgent, "ai-agent", "", "AI agent type: claude, general, none")

	return cmd
}

func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func selectAIAgent() (string, error) {
	var agent string
	prompt := &survey.Select{
		Message: "Which AI agent will you use?",
		Options: []string{
			"Claude Code (Anthropic)",
			"General AI Agent (OpenAI, Gemini, etc.)",
			"None (skip AI documentation)",
		},
		Default: "Claude Code (Anthropic)",
	}
	err := survey.AskOne(prompt, &agent)
	return agent, err
}

func mapAgentSelectionToType(agent string) ai.AgentType {
	switch agent {
	case "Claude Code (Anthropic)":
		return ai.AgentClaude
	case "General AI Agent (OpenAI, Gemini, etc.)":
		return ai.AgentGeneral
	default:
		return ai.AgentNone
	}
}

func mapAgentFlagToType(flag string) ai.AgentType {
	switch flag {
	case "claude":
		return ai.AgentClaude
	case "general":
		return ai.AgentGeneral
	case "none":
		return ai.AgentNone
	default:
		return ai.AgentNone
	}
}

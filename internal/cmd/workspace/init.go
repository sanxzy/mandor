package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		skillsDir     string
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

			if agentType != ai.AgentNone {
				if err := injectSkills(agentType, cwd, skillsDir); err != nil {
					fmt.Printf("Warning: Failed to inject skills: %v\n", err)
				}
			}

			fmt.Printf("Workspace initialized: %s\n", ws.Name)
			fmt.Printf("  Location: .mandor/\n")
			fmt.Printf("  ID: %s\n", ws.ID)
			fmt.Printf("  Creator: %s\n", ws.CreatedBy)
			fmt.Printf("  Created: %s\n", ws.CreatedAt.Format("2006-01-02T15:04:05Z"))

			locations := ai.GetAgentLocations(agentType)
			fmt.Printf("\nGenerated files:\n")
			if agentType == ai.AgentClaude {
				fmt.Printf("  - CLAUDE.md\n")
			} else if agentType == ai.AgentGeneral {
				fmt.Printf("  - AGENTS.md\n")
			}
			for _, loc := range locations {
				fmt.Printf("  - %s\n", loc)
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
	cmd.Flags().StringVar(&aiAgent, "ai-agent", "", "AI agent type: claude, opencode, copilot, factory-droid, windsurf, antigravity, cline, roocode, gemini-cli, amazon-q, qoder, none")
	cmd.Flags().StringVar(&skillsDir, "skills-dir", "", "Path to skills directory (default: downloads from GitHub)")

	return cmd
}

func injectSkills(agentType ai.AgentType, workspacePath, skillsDir string) error {
	locations := ai.GetAgentLocations(agentType)
	mapper := ai.NewSkillMapper()

	tempDir, err := os.MkdirTemp("", "mandor-skills-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download skills from GitHub
	if err := mapper.DownloadAllSkillsFromGitHub("sanxzy", "mandor", "main", "skill-templates", tempDir); err != nil {
		return fmt.Errorf("failed to download skills: %w", err)
	}

	// Map skill names to directory names and file names
	skillFiles := map[string]string{
		"planner":   "planner.md",
		"specs":     "specs.md",
		"blueprint": "blueprint.md",
	}

	for skillName, filename := range skillFiles {
		skillPath := filepath.Join(tempDir, filename)
		if !fileExists(skillPath) {
			continue
		}

		skill, err := mapper.ParseSkillTemplate(skillPath)
		if err != nil {
			continue
		}

		for _, location := range locations {
			// Create skill directory: <location>/mandor-<skill>/
			skillDirName := fmt.Sprintf("mandor-%s", skillName)
			targetDir := filepath.Join(workspacePath, location, skillDirName)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				continue
			}

			// Write SKILL.md with frontmatter
			targetPath := filepath.Join(targetDir, "SKILL.md")
			skillContent := fmt.Sprintf(`---
name: "%s"
description: %s
category: %s
tags: [%s]
---

%s
`,
				skill.Name,
				skill.Description,
				skill.Category,
				strings.Join(skill.Tags, ", "),
				skill.Content,
			)
			if err := os.WriteFile(targetPath, []byte(skillContent), 0644); err != nil {
				continue
			}
		}
	}

	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getExtension(format string) string {
	switch format {
	case "claude-skill", "claude-subagent", "opencode-skill", "amazon-q-agent", "roocode-mcp":
		return "json"
	default:
		return "md"
	}
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
			"OpenCode",
			"GitHub Copilot",
			"Factory Droid",
			"Windsurf",
			"Google Antigravity",
			"Cline",
			"RooCode",
			"Gemini CLI",
			"Amazon Q Developer",
			"Qoder",
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
	case "OpenCode":
		return ai.AgentOpenCode
	case "GitHub Copilot":
		return ai.AgentCopilot
	case "Factory Droid":
		return ai.AgentFactoryDroid
	case "Windsurf":
		return ai.AgentWindsurf
	case "Google Antigravity":
		return ai.AgentAntigravity
	case "Cline":
		return ai.AgentCline
	case "RooCode":
		return ai.AgentRooCode
	case "Gemini CLI":
		return ai.AgentGeminiCLI
	case "Amazon Q Developer":
		return ai.AgentAmazonQ
	case "Qoder":
		return ai.AgentQoder
	default:
		return ai.AgentNone
	}
}

func mapAgentFlagToType(flag string) ai.AgentType {
	switch flag {
	case "claude":
		return ai.AgentClaude
	case "opencode":
		return ai.AgentOpenCode
	case "copilot":
		return ai.AgentCopilot
	case "factory-droid":
		return ai.AgentFactoryDroid
	case "windsurf":
		return ai.AgentWindsurf
	case "antigravity":
		return ai.AgentAntigravity
	case "cline":
		return ai.AgentCline
	case "roocode":
		return ai.AgentRooCode
	case "gemini-cli":
		return ai.AgentGeminiCLI
	case "amazon-q":
		return ai.AgentAmazonQ
	case "qoder":
		return ai.AgentQoder
	case "none":
		return ai.AgentNone
	default:
		return ai.AgentNone
	}
}

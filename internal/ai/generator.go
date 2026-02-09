package ai

import (
	"os"
	"path/filepath"
	"time"
)

type AgentType string

const (
	AgentClaude       AgentType = "claude"
	AgentOpenCode     AgentType = "opencode"
	AgentCopilot      AgentType = "copilot"
	AgentFactoryDroid AgentType = "factory-droid"
	AgentWindsurf     AgentType = "windsurf"
	AgentAntigravity  AgentType = "antigravity"
	AgentCline        AgentType = "cline"
	AgentRooCode      AgentType = "roocode"
	AgentGeminiCLI    AgentType = "gemini-cli"
	AgentAmazonQ      AgentType = "amazon-q"
	AgentQoder        AgentType = "qoder"
	AgentGeneral      AgentType = "general"
	AgentNone         AgentType = "none"
)

func GetAgentLocations(agent AgentType) []string {
	switch agent {
	case AgentClaude:
		return []string{".claude/skills/"}
	case AgentOpenCode:
		return []string{".opencode/skills/"}
	case AgentCopilot:
		return []string{".github/skills/"}
	case AgentFactoryDroid:
		return []string{".factory/droids/"}
	case AgentWindsurf:
		return []string{".windsurf/workflows/"}
	case AgentAntigravity:
		return []string{".agent/skills/"}
	case AgentCline:
		return []string{".cline/skills/"}
	case AgentRooCode:
		return []string{".roo/skills/"}
	case AgentGeminiCLI:
		return []string{".gemini/skills/"}
	case AgentAmazonQ:
		return []string{".amazon-q/"}
	case AgentQoder:
		return []string{".qoder/skills/"}
	default:
		return []string{".skills/"}
	}
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(agent AgentType, projectName, createdAt string) string {
	switch agent {
	case AgentClaude:
		return generateClaudeTemplate(projectName, createdAt)
	case AgentGeneral:
		return generateAgentsTemplate(projectName, createdAt)
	default:
		return ""
	}
}

func (g *Generator) Write(agent AgentType, workspacePath string) error {
	if agent == AgentNone {
		return nil
	}

	content := g.Generate(agent, filepath.Base(workspacePath), time.Now().UTC().Format("2006-01-02"))

	var targetPath string
	if agent == AgentClaude {
		targetPath = filepath.Join(workspacePath, "CLAUDE.md")
	} else {
		targetPath = filepath.Join(workspacePath, "AGENTS.md")
	}

	return os.WriteFile(targetPath, []byte(content), 0644)
}

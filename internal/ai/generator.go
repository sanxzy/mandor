package ai

import (
	"os"
	"path/filepath"
	"time"
)

type AgentType string

const (
	AgentClaude  AgentType = "claude"
	AgentGeneral AgentType = "general"
	AgentNone    AgentType = "none"
)

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

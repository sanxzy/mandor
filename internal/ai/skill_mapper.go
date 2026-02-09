package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	Name          string   `yaml:"name" json:"name"`
	Description   string   `yaml:"description" json:"description"`
	Category      string   `yaml:"category,omitempty" json:"category,omitempty"`
	Tags          []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	License       string   `yaml:"license,omitempty" json:"license,omitempty"`
	Compatibility string   `yaml:"compatibility,omitempty" json:"compatibility,omitempty"`
	Content       string   `yaml:"-" json:"-"`
}

type SkillMapper struct{}

func NewSkillMapper() *SkillMapper {
	return &SkillMapper{}
}

func (m *SkillMapper) ParseSkillTemplate(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	content := string(data)

	startMarker := "---\n"
	endMarker := "\n---\n"

	startIdx := strings.Index(content, startMarker)
	if startIdx != 0 {
		return nil, fmt.Errorf("invalid skill format: missing opening ---")
	}

	endIdx := strings.Index(content[startIdx+len(startMarker):], endMarker)
	if endIdx == -1 {
		return nil, fmt.Errorf("invalid skill format: missing closing ---")
	}

	frontmatter := content[len(startMarker) : len(startMarker)+endIdx]
	skillContent := strings.TrimSpace(content[len(startMarker)+endIdx+len(endMarker):])

	var skill Skill
	if err := yaml.Unmarshal([]byte(frontmatter), &skill); err != nil {
		return nil, fmt.Errorf("parse YAML frontmatter: %w", err)
	}

	skill.Content = skillContent
	return &skill, nil
}

func (m *SkillMapper) ToClaudeCode(skill *Skill) string {
	claudeSkill := map[string]interface{}{
		"name":        skill.Name,
		"description": skill.Description,
		"category":    skill.Category,
		"tags":        skill.Tags,
		"steps":       m.extractSteps(skill.Content),
		"questions":   m.extractQuestions(skill.Content),
		"guardrails":  m.extractGuardrails(skill.Content),
	}

	output, _ := json.MarshalIndent(claudeSkill, "", "  ")
	return fmt.Sprintf("---\n%s\n---\n%s", string(output), skill.Content)
}

func (m *SkillMapper) ToClaudeSubagent(skill *Skill) string {
	subagent := map[string]interface{}{
		"name":        skill.Name,
		"description": skill.Description,
		"model":       "claude-sonnet-4-5",
		"tools":       []string{"read", "write", "edit", "grep", "glob", "bash"},
		"system":      skill.Content,
	}

	output, _ := json.MarshalIndent(subagent, "", "  ")
	return string(output)
}

func (m *SkillMapper) ToOpenCode(skill *Skill) string {
	opencodeSkill := map[string]interface{}{
		"name":         skill.Name,
		"description":  skill.Description,
		"category":     skill.Category,
		"tags":         skill.Tags,
		"permissions":  map[string]string{"read": "allow", "write": "ask", "execute": "ask"},
		"questions":    m.extractQuestions(skill.Content),
		"instructions": skill.Content,
	}

	output, _ := json.MarshalIndent(opencodeSkill, "", "  ")
	return fmt.Sprintf("---\n%s\n---\n%s", string(output), skill.Content)
}

func (m *SkillMapper) ToCopilotAgent(skill *Skill) string {
	copilotAgent := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
tools:
  - read
  - write
  - edit
  - grep
  - glob
  - bash
mcpServers: []
---

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Content,
	)

	return copilotAgent
}

func (m *SkillMapper) ToFactoryDroid(skill *Skill) string {
	droid := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
model: sonnet-4
reasoningEffort: high
toolCategories:
  - read-only
  - edit
  - execute
  - web
---

# %s

%s

## Workflow

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Name,
		skill.Description,
		m.formatStepsForDroid(skill.Content),
	)

	return droid
}

func (m *SkillMapper) ToWindsurfWorkflow(skill *Skill) string {
	workflow := fmt.Sprintf(`---
name: %s
description: %s
category: %s
trigger: slash /%s
steps:
%s
---

# %s Workflow

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.ReplaceAll(skill.Name, " ", "-"),
		m.formatWorkflowSteps(skill.Content),
		skill.Name,
		skill.Description,
	)

	return workflow
}

func (m *SkillMapper) ToAntigravitySkill(skill *Skill) string {
	antigravity := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
version: 1.0.0
rules:
  - type: execution
    mode: sequential
  - type: validation
    requireConfirmation: true
---

# %s

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Name,
		skill.Description,
	)

	return antigravity
}

func (m *SkillMapper) ToClineSkill(skill *Skill) string {
	cline := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
loading: progressive
---

# %s

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Name,
		skill.Content,
	)

	return cline
}

func (m *SkillMapper) ToRooCodeMCP(skill *Skill) string {
	mcpConfig := map[string]interface{}{
		"name":        skill.Name,
		"description": skill.Description,
		"transport":   "stdio",
		"command":     fmt.Sprintf("mandor-skill-%s", skill.Name),
		"args":        []string{"--skill", skill.Name},
		"env":         map[string]string{},
	}

	output, _ := json.MarshalIndent(mcpConfig, "", "  ")
	return fmt.Sprintf("# RooCode MCP Configuration for %s\n\n%s\n", skill.Name, string(output))
}

func (m *SkillMapper) ToGeminiCLISkill(skill *Skill) string {
	gemini := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
activation: activate_skill
version: 1.0.0
---

# %s

%s

## Usage

Use `+"`"+`activate_skill("%s")`+"`"+` to activate this skill.
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Name,
		skill.Description,
		skill.Name,
	)

	return gemini
}

func (m *SkillMapper) ToAmazonQAgent(skill *Skill) string {
	qAgent := map[string]interface{}{
		"name":        skill.Name,
		"description": skill.Description,
		"type":        "agent",
		"version":     "1.0",
		"context": map[string]interface{}{
			"hooks": []string{"pre-execution", "post-execution"},
		},
		"tools":        []string{"ReadFile", "WriteFile", "EditFile", "Glob", "Grep", "Bash"},
		"mcpServers":   []string{},
		"permissions":  map[string]string{"read": "allow", "write": "ask"},
		"instructions": skill.Content,
	}

	output, _ := json.MarshalIndent(qAgent, "", "  ")
	return fmt.Sprintf("# Amazon Q Agent: %s\n\n%s\n", skill.Name, string(output))
}

func (m *SkillMapper) ToQoderSkill(skill *Skill) string {
	qoder := fmt.Sprintf(`---
name: %s
description: %s
category: %s
tags: %s
priority: high
scope: project
---

# %s

%s
`,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.Name,
		skill.Content,
	)

	return qoder
}

func (m *SkillMapper) ToGenericMarkdown(skill *Skill, agent AgentType) string {
	agentNote := fmt.Sprintf("<!-- MANDOR_SKILL: %s | AGENT: %s | DATE: %s -->",
		skill.Name, agent, time.Now().Format("2006-01-02"))

	return fmt.Sprintf(`%s

---
# %s

## Description

%s

## Metadata
- **Category**: %s
- **Tags**: %s
- **License**: %s
- **Compatibility**: %s

---

## Skill Content

%s

---
*Generated by Mandor Skill Mapper for %s*
`,
		agentNote,
		skill.Name,
		skill.Description,
		skill.Category,
		strings.Join(skill.Tags, ", "),
		skill.License,
		skill.Compatibility,
		skill.Content,
		agent,
	)
}

func (m *SkillMapper) ToPromptTemplate(skill *Skill) string {
	return fmt.Sprintf(`## Skill: %s

### Description
%s

### Instructions
%s

---
*Mandor Skill - Generated for AI Agent*
`,
		skill.Name,
		skill.Description,
		skill.Content,
	)
}

func (m *SkillMapper) ToOpenSpec(skill *Skill) string {
	return fmt.Sprintf(`---
name: %s
description: %s
license: MIT
compatibility: agentskills.io/v1
---

# %s

%s

---
*Generated by Mandor Skill Mapper*
`,
		skill.Name,
		skill.Description,
		skill.Name,
		skill.Content,
	)
}

func (m *SkillMapper) Map(skill *Skill, agent AgentType, format string) (string, error) {
	switch format {
	case "claude-skill":
		return m.ToClaudeCode(skill), nil
	case "claude-subagent":
		return m.ToClaudeSubagent(skill), nil
	case "opencode-skill":
		return m.ToOpenCode(skill), nil
	case "copilot-agent":
		return m.ToCopilotAgent(skill), nil
	case "factory-droid":
		return m.ToFactoryDroid(skill), nil
	case "windsurf-workflow":
		return m.ToWindsurfWorkflow(skill), nil
	case "antigravity-skill":
		return m.ToAntigravitySkill(skill), nil
	case "cline-skill":
		return m.ToClineSkill(skill), nil
	case "roocode-mcp":
		return m.ToRooCodeMCP(skill), nil
	case "gemini-cli-skill":
		return m.ToGeminiCLISkill(skill), nil
	case "amazon-q-agent":
		return m.ToAmazonQAgent(skill), nil
	case "qoder-skill":
		return m.ToQoderSkill(skill), nil
	case "markdown":
		return m.ToGenericMarkdown(skill, agent), nil
	case "prompt":
		return m.ToPromptTemplate(skill), nil
	case "open-spec":
		return m.ToOpenSpec(skill), nil
	default:
		return m.ToGenericMarkdown(skill, agent), nil
	}
}

func (m *SkillMapper) MapAndWrite(skill *Skill, agent AgentType, format, outputPath string) error {
	mapped, err := m.Map(skill, agent, format)
	if err != nil {
		return fmt.Errorf("map skill: %w", err)
	}

	filename := fmt.Sprintf("%s.%s", skill.Name, getExtension(format))
	targetPath := filepath.Join(outputPath, filename)

	return os.WriteFile(targetPath, []byte(mapped), 0644)
}

func (m *SkillMapper) extractSteps(content string) []string {
	var steps []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "1. **") || strings.HasPrefix(line, "2. **") ||
			strings.HasPrefix(line, "3. **") || strings.HasPrefix(line, "4. **") ||
			strings.HasPrefix(line, "5. **") || strings.HasPrefix(line, "6. **") ||
			strings.HasPrefix(line, "7. **") || strings.HasPrefix(line, "8. **") ||
			strings.HasPrefix(line, "9. **") || strings.HasPrefix(line, "10. **") {
			step := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "0123456789. **"), "**"))
			steps = append(steps, step)
		}
	}
	return steps
}

func (m *SkillMapper) extractQuestions(content string) []string {
	var questions []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Ask:") || strings.Contains(line, "ask:") ||
			strings.Contains(line, "Ask user") || strings.Contains(line, "Confirm") {
			q := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			q = strings.TrimSpace(strings.TrimPrefix(q, "Ask: "))
			q = strings.TrimSpace(strings.TrimPrefix(q, "ask: "))
			if q != "" && len(q) > 5 {
				questions = append(questions, q)
			}
		}
	}
	return questions
}

func (m *SkillMapper) extractGuardrails(content string) []string {
	var guardrails []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") && (strings.Contains(line, "Don't") ||
			strings.Contains(line, "ERROR") ||
			strings.Contains(line, "MUST") ||
			strings.Contains(line, "refuse") ||
			strings.Contains(line, "Do NOT") ||
			strings.Contains(line, "Prerequisite")) {
			guardrails = append(guardrails, strings.TrimPrefix(line, "- "))
		}
	}
	return guardrails
}

func (m *SkillMapper) formatStepsForDroid(content string) string {
	steps := m.extractSteps(content)
	var result string
	for i, step := range steps {
		result += fmt.Sprintf("  - step %d: %s\n", i+1, step)
	}
	return result
}

func (m *SkillMapper) formatQuestionsForDroid(content string) string {
	questions := m.extractQuestions(content)
	var result string
	for _, q := range questions {
		result += fmt.Sprintf("- %s\n", q)
	}
	return result
}

func (m *SkillMapper) formatGuardrailsForDroid(content string) string {
	guardrails := m.extractGuardrails(content)
	var result string
	for _, g := range guardrails {
		result += fmt.Sprintf("- %s\n", g)
	}
	return result
}

func (m *SkillMapper) formatWorkflowSteps(content string) string {
	steps := m.extractSteps(content)
	var result string
	for i, step := range steps {
		result += fmt.Sprintf("  - name: Step %d\n    description: %s\n", i+1, step)
	}
	return result
}

func (m *SkillMapper) formatQuestionsForWorkflow(content string) string {
	questions := m.extractQuestions(content)
	var result string
	for _, q := range questions {
		result += fmt.Sprintf("- %s\n", q)
	}
	return result
}

func (m *SkillMapper) formatQuestionsList(content string) string {
	questions := m.extractQuestions(content)
	var result string
	for _, q := range questions {
		result += fmt.Sprintf("- [ ] %s\n", q)
	}
	return result
}

func (m *SkillMapper) formatGuardrailsList(content string) string {
	guardrails := m.extractGuardrails(content)
	var result string
	for _, g := range guardrails {
		result += fmt.Sprintf("- **%s**\n", g)
	}
	return result
}

func (m *SkillMapper) formatGeminiSteps(content string) string {
	steps := m.extractSteps(content)
	var result string
	for i, step := range steps {
		result += fmt.Sprintf("%d. %s\n", i+1, step)
	}
	return result
}

func (m *SkillMapper) formatStepsForQoder(content string) string {
	steps := m.extractSteps(content)
	var result string
	for _, step := range steps {
		result += fmt.Sprintf("→ %s\n", step)
	}
	return result
}

const (
	DefaultSkillRepoOwner = "sanxzy"
	DefaultSkillRepoName  = "mandor"
	DefaultSkillBranch    = "main"
	DefaultSkillPath      = "skill-templates"
)

func (m *SkillMapper) DownloadSkillFromGitHub(skillName, owner, repo, branch, githubPath, saveDir string) error {
	if owner == "" {
		owner = DefaultSkillRepoOwner
	}
	if repo == "" {
		repo = DefaultSkillRepoName
	}
	if branch == "" {
		branch = DefaultSkillBranch
	}
	if githubPath == "" {
		githubPath = DefaultSkillPath
	}

	skillFileName := skillName
	if skillName == "planner" {
		skillFileName = "planner.md"
	} else if skillName == "specs" || skillName == "spec" {
		skillFileName = "specs.md"
	} else if skillName == "blueprint" {
		skillFileName = "blueprint.md"
	} else {
		skillFileName = skillName + ".md"
	}

	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s/%s", owner, repo, branch, githubPath, skillFileName)

	resp, err := http.Get(rawURL)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", skillName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: HTTP %d", skillName, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", skillName, err)
	}

	targetPath := filepath.Join(saveDir, skillFileName)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(targetPath, body, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", skillName, err)
	}

	return nil
}

func (m *SkillMapper) DownloadAllSkillsFromGitHub(owner, repo, branch, githubPath, saveDir string) error {
	skills := []string{"planner", "specs", "blueprint"}

	for _, skill := range skills {
		if err := m.DownloadSkillFromGitHub(skill, owner, repo, branch, githubPath, saveDir); err != nil {
			return err
		}
	}

	return nil
}

func getExtension(format string) string {
	switch format {
	case "claude-skill", "claude-subagent", "opencode-skill", "amazon-q-agent", "roocode-mcp":
		return "json"
	case "copilot-agent", "factory-droid", "windsurf-workflow", "antigravity-skill", "cline-skill", "gemini-cli-skill", "qoder-skill", "open-spec":
		return "md"
	case "markdown", "prompt":
		return "md"
	default:
		return "md"
	}
}

func GetFormatsForAgent(agent AgentType) []string {
	switch agent {
	case AgentClaude:
		return []string{"markdown"}
	case AgentOpenCode:
		return []string{"markdown"}
	case AgentCopilot:
		return []string{"markdown"}
	case AgentFactoryDroid:
		return []string{"markdown"}
	case AgentWindsurf:
		return []string{"markdown"}
	case AgentAntigravity:
		return []string{"markdown"}
	case AgentCline:
		return []string{"markdown"}
	case AgentRooCode:
		return []string{"markdown"}
	case AgentGeminiCLI:
		return []string{"markdown"}
	case AgentAmazonQ:
		return []string{"markdown"}
	case AgentQoder:
		return []string{"markdown"}
	default:
		return []string{"markdown"}
	}
}

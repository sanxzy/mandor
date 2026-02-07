# Plan: Remove `mandor ai` Command and Integrate into `mandor init`

## Overview

Remove the standalone `mandor ai` command and make AI documentation generation automatic during workspace initialization. The `mandor init` command will prompt users to select their AI agent, then inject the appropriate documentation file (`CLAUDE.md` for Claude Code, `AGENTS.md` for other AI agents).

---

## Current State Analysis

### Files to Remove
| File | Purpose |
|------|---------|
| `internal/cmd/ai/root.go` | Root AI command wrapper |
| `internal/cmd/ai/claude.go` | Claude-specific documentation generator |
| `internal/cmd/ai/agents.go` | General AI agents documentation generator |
| `internal/cmd/ai/` directory | Entire ai command package |

### Root Command Change
**File**: `internal/cmd/root.go`
- Remove import: `"mandor/internal/cmd/ai"`
- Remove line: `rootCmd.AddCommand(ai.NewAICmd())`

### Current Functionality
- `mandor ai claude` → Generates `CLAUDE.md`
- `mandor ai agents` → Generates `AGENTS.md`
- Both use similar template generation logic

---

## Proposed Changes

### 1. Update `mandor init` Command

**File**: `internal/cmd/workspace/init.go`

#### New Flags
```go
var (
    workspaceName string
    skipConfirm   bool
    strict        bool
    aiAgent       string  // NEW: --ai-agent flag (claude, general, none)
    yes           bool    // NEW: -y flag for skipping prompts
)
```

#### Behavior
1. If `--ai-agent` is provided, skip prompt and use specified agent
2. If no flag and `--yes` is NOT set, show interactive prompt
3. Generate appropriate file based on selection:
   - `claude` → Generate `CLAUDE.md`
   - `general` or other → Generate `AGENTS.md`
   - `none` → Skip AI documentation generation

#### Interactive Prompt Design
```bash
? Which AI agent will you use? (Use arrow keys)
  ▸ Claude Code (Anthropic)
    General AI Agent (OpenAI, Gemini, etc.)
    None (skip AI documentation)
```

---

### 2. Create AI Documentation Generator Module

**New File**: `internal/ai/generator.go`

```go
package ai

type AgentType string

const (
    AgentClaude  AgentType = "claude"
    AgentGeneral AgentType = "general"
    AgentNone   AgentType = "none"
)

type Generator struct{}

func (g *Generator) Generate(agent AgentType, projectName, createdAt string) string
func (g *Generator) Write(agent AgentType, workspacePath string) error
```

---

### 3. AI Documentation Templates

**New File**: `internal/ai/templates.go`

Separate templates for each agent type:

```go
var claudeTemplate = `... CLAUDE.md content ...`
var agentsTemplate = `... AGENTS.md content ...`
```

---

### 4. Move Existing Templates

Move content from `internal/cmd/ai/agents.go` to `internal/ai/templates.go`:
- `generateAIDoc()` → Refactored into template system

---

## Implementation Steps

### Phase 1: Create AI Generator Module

#### Step 1.1: Create `internal/ai/generator.go`
```go
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
    AgentNone   AgentType = "none"
)

type Generator struct{}

func NewGenerator() *Generator {
    return &Generator{}
}

func (g *Generator) Generate(agent AgentType, projectName, createdAt string) string {
    switch agent {
    case AgentClaude:
        return g.generateClaudeTemplate(projectName, createdAt)
    case AgentGeneral:
        return g.generateAgentsTemplate(projectName, createdAt)
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

// Template methods...
```

#### Step 1.2: Create `internal/ai/templates.go`
- Move `generateAIDoc()` logic from `internal/cmd/ai/agents.go`
- Create separate templates for Claude vs General agents

---

### Phase 2: Update `mandor init` Command

#### Step 2.1: Modify `internal/cmd/workspace/init.go`

**Add imports**:
```go
import (
    "github.com/AlecAivazis/survey/v2"
    "mandor/internal/ai"
)
```

**Update command struct**:
```go
var (
    workspaceName string
    skipConfirm   bool
    strict        bool
    aiAgent       string  // NEW: --ai-agent flag
    yes           bool    // NEW: -y flag
)
```

**Add AI agent selection prompt**:
```go
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

func mapAgentToType(agent string) ai.AgentType {
    switch agent {
    case "Claude Code (Anthropic)":
        return ai.AgentClaude
    case "General AI Agent (OpenAI, Gemini, etc.)":
        return ai.AgentGeneral
    default:
        return ai.AgentNone
    }
}
```

**Update RunE**:
```go
RunE: func(cmd *cobra.Command, args []string) error {
    // ... existing workspace initialization ...

    // AI Documentation Generation
    var agentType ai.AgentType
    if aiAgent != "" {
        // Flag provided
        agentType = mapFlagToAgentType(aiAgent)
    } else if !yes {
        // Interactive prompt
        agent, err := selectAIAgent()
        if err != nil {
            return err
        }
        agentType = mapAgentToType(agent)
    } else {
        // Default to Claude
        agentType = ai.AgentClaude
    }

    // Generate AI documentation
    generator := ai.NewGenerator()
    if err := generator.Write(agentType, cwd); err != nil {
        return err
    }

    // Report generated file
    if agentType == ai.AgentClaude {
        fmt.Printf("✓ Generated: CLAUDE.md\n")
    } else if agentType == ai.AgentGeneral {
        fmt.Printf("✓ Generated: AGENTS.md\n")
    }

    return nil
},
```

#### Step 2.2: Add Flags
```go
cmd.Flags().StringVar(&aiAgent, "ai-agent", "", "AI agent type: claude, general, none")
cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip all prompts (use defaults)")
```

---

### Phase 3: Remove `mandor ai` Command

#### Step 3.1: Remove AI Command Files
```bash
rm -rf internal/cmd/ai/
```

#### Step 3.2: Update Root Command
**File**: `internal/cmd/root.go`
- Remove import: `"mandor/internal/cmd/ai"`
- Remove line: `rootCmd.AddCommand(ai.NewAICmd())`

---

### Phase 4: Update Dependencies

**File**: `go.mod`
```diff
require github.com/spf13/cobra v1.10.2
+require github.com/AlecAivazis/survey/v2 v2.3.7
```

```bash
go get github.com/AlecAivazis/survey/v2@v2.3.7
go mod tidy
```

---

### Phase 5: Update Documentation

#### Step 5.1: Update `docs/RELEASE.md`
- Remove references to `mandor ai` command
- Add AI agent selection to init flow

#### Step 5.2: Update `README.md`
- Update command reference
- Add AI agent selection documentation

#### Step 5.3: Update `CHANGELOG.md`
```markdown
## [0.7.0] - TBD

### Changed
- Removed `mandor ai` command - AI documentation is now auto-generated during `mandor init`
- `mandor init` prompts for AI agent selection (Claude Code or General)
- AI documentation file generated automatically: CLAUDE.md or AGENTS.md

### Added
- `--ai-agent` flag for `mandor init` (claude, general, none)
- `-y` flag for non-interactive initialization
```

---

## File Changes Summary

### Files to Create
| File | Purpose |
|------|---------|
| `internal/ai/generator.go` | AI documentation generator |
| `internal/ai/templates.go` | Claude and AGENTS templates |

### Files to Modify
| File | Changes |
|------|---------|
| `internal/cmd/workspace/init.go` | Add AI agent prompt, flags, file generation |
| `internal/cmd/root.go` | Remove ai command import/addition |
| `go.mod` | Add survey dependency |
| `docs/RELEASE.md` | Update documentation |
| `README.md` | Update command reference |
| `CHANGELOG.md` | Document breaking change |

### Files to Delete
| File | Purpose |
|------|---------|
| `internal/cmd/ai/root.go` | Removed AI command root |
| `internal/cmd/ai/claude.go` | Merged into generator |
| `internal/cmd/ai/agents.go` | Merged into generator |

---

## Testing Strategy

### Unit Tests
1. **Agent type mapping**: Test `mapFlagToAgentType()` and `mapAgentToType()`
2. **Template generation**: Verify correct template used per agent
3. **File generation**: Test `Write()` method with mock filesystem
4. **Flag parsing**: Test `--ai-agent` and `-y` flags

### Integration Tests
1. **`mandor init` with `--ai-agent claude`**: Verify `CLAUDE.md` created
2. **`mandor init` with `--ai-agent general`**: Verify `AGENTS.md` created
3. **`mandor init` with `--ai-agent none`**: Verify no file created
4. **Interactive mode**: Verify prompt appears without `-y` flag
5. **Non-interactive mode**: Verify `-y` skips prompt and uses default

### Backward Compatibility
- Verify `mandor ai` is no longer registered
- Verify `mandor init` still works without AI flags (backward compatible)

---

## Usage Examples

### Interactive Mode (Default)
```bash
$ mandor init
✓ Workspace initialized: my-project
  Location: .mandor/
  ID: abc123
  Creator: developer

? Which AI agent will you use? › - Use arrow keys
   ▸ Claude Code (Anthropic)
     General AI Agent (OpenAI.)
     None (, Gemini, etcskip AI documentation)

✓ Generated: CLAUDE.md

Next steps:
  1. Create a project: mandor project create <project_id> --name "Project Name"
```

### Non-Interactive with Flag
```bash
$ mandor init --ai-agent claude -y
✓ Workspace initialized: my-project
✓ Generated: CLAUDE.md
```

### Skip AI Documentation
```bash
$ mandor init --ai-agent none
✓ Workspace initialized: my-project
  (No AI documentation generated)
```

### Invalid Agent Type
```bash
$ mandor init --ai-agent invalid
Error: invalid AI agent type. Valid options: claude, general, none
```

---

## Migration Guide for Users

**Before (v0.6.x)**:
```bash
mandor init
mandor ai claude          # Generate CLAUDE.md
```

**After (v0.7.0)**:
```bash
mandor init              # Prompts for AI agent
# OR
mandor init --ai-agent claude -y  # Auto-generate CLAUDE.md
```

---

## Rollback Plan

If issues arise, revert to commit before this change:
```bash
git revert <commit-hash>
go mod tidy
```

---

## Estimated Effort

| Task | Effort |
|------|--------|
| Create AI generator module | 1-2 hours |
| Update mandor init | 2-3 hours |
| Remove mandor ai command | 30 minutes |
| Update documentation | 1 hour |
| Testing | 2 hours |
| **Total** | **6-8 hours** |

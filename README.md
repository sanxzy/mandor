# Mandor - Event-Based Task Manager CLI for AI Agent Workflows

<p align="center">
  <strong>Compact context | Deterministic JSONL output | Dependency enforcement</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#commands">Commands</a>
</p>

---

## Overview

Mandor is a CLI tool for managing tasks, features, and issues in AI agent workflows. It provides:

- **Compact Context**: Structured JSON output instead of verbose descriptions
- **Event-Based**: All changes logged in `events.jsonl` for audit trail
- **Dependency Tracking**: Automatic status based on dependencies
- **Schema-Driven**: Configurable rules per project
- **Cross-Platform**: Works on macOS, Linux, and Windows

---

## Installation

### From Source

```bash
git clone https://github.com/budisantoso/mandor.git
cd mandor
go build -o build/mandor ./cmd/mandor
sudo mv build/mandor /usr/local/bin/
```

### From NPM

```bash
npm install -g @mandor/cli
npx @mandor/cli init "My Project"
```

### Verify

```bash
mandor --version
mandor --help
```

---

## Quick Start

```bash
# 1. Initialize workspace
mandor init "My Project"

# 2. Create project
mandor project create api --name "API Service" \
  --goal "Implement comprehensive REST API service."

# 3. Create feature
mandor feature create "User Authentication" --project api \
  --goal "Implement user login and registration."

# 4. Create task
mandor task create "Password Hashing" \
  --feature api-feature-xxx \
  --goal "Implement bcrypt password hashing" \
  --implementation-steps "Install bcrypt,Create utility,Write tests" \
  --test-cases "Hash validation,Password comparison" \
  --derivable-files "src/utils/password.ts" \
  --library-needs "bcrypt"

# 5. Check status
mandor status
```

---

## Commands

### Workspace

| Command | Description |
|---------|-------------|
| `mandor init <name>` | Initialize workspace |
| `mandor status` | Show workspace status |
| `mandor config get <key>` | Get config value |
| `mandor config set <key> <value>` | Set config value |

### Project

| Command | Description |
|---------|-------------|
| `mandor project create <id> --name <name> --goal <goal>` | Create project |
| `mandor project list` | List projects |
| `mandor project detail <id>` | Show project details |
| `mandor project update <id>` | Update project |
| `mandor project delete <id>` | Delete project |
| `mandor project reopen <id>` | Reopen deleted project |

### Feature

| Command | Description |
|---------|-------------|
| `mandor feature create <name> --project <id> --goal <goal>` | Create feature |
| `mandor feature list [--project <id>]` | List features |
| `mandor feature detail <id>` | Show feature details |
| `mandor feature update <id>` | Update feature |

**Status Flow:** `draft → active → done` (blocked if dependencies not met)

### Task

| Command | Description |
|---------|-------------|
| `mandor task create <name> --feature <id> --goal <goal> --implementation-steps <steps> --test-cases <cases> --derivable-files <files> --library-needs <libs>` | Create task |
| `mandor task list [--project <id>] [--feature <id>]` | List tasks |
| `mandor task detail <id>` | Show task details |
| `mandor task update <id>` | Update task |
| `mandor task ready [--project <id>]` | List ready tasks |
| `mandor task blocked [--project <id>]` | List blocked tasks |

**Status Flow:** `pending → ready → in_progress → done` (blocked if deps not met)

### Issue

| Command | Description |
|---------|-------------|
| `mandor issue create <name> --project <id> --type <type> --goal <goal>` | Create issue |
| `mandor issue list [--project <id>]` | List issues |
| `mandor issue detail <id>` | Show issue details |
| `mandor issue update <id>` | Update issue |
| `mandor issue ready [--project <id>]` | List ready issues |
| `mandor issue blocked [--project <id>]` | List blocked issues |

**Types:** `bug`, `improvement`, `debt`, `security`, `performance`  
**Status Flow:** `open → ready → in_progress → resolved` (or `wontfix`/`cancelled`)

### Utility

| Command | Description |
|---------|-------------|
| `mandor completion [bash\|zsh\|fish]` | Generate shell completions |
| `mandor populate` | Show full command reference |

---

## Dependency Management

### ID Format

| Entity | Format | Example |
|--------|--------|---------|
| Project | `<id>` | `api` |
| Feature | `<project>-feature-<nanoid>` | `api-feature-abc123` |
| Task | `<feature_id>-task-<nanoid>` | `api-feature-abc-task-xyz789` |
| Issue | `<project>-issue-<nanoid>` | `api-issue-abc123` |

### Status Rules

- **Feature**: No deps → `draft`, all done → `active`, else `blocked`
- **Task**: No deps → `ready`, all done → `ready`, else `pending`
- **Issue**: No deps → `ready`, all resolved → `ready`, else `open`

### Blocking

Cannot cancel entities that others depend on.

---

## Configuration

### Priority Levels

| Priority | Description |
|----------|-------------|
| P0 | Critical - Must do |
| P1 | High - Important |
| P2 | Medium - Should do |
| P3 | Normal - Default |
| P4 | Low - Nice to have |
| P5 | Minimal - Can defer |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | System error |
| 2 | Validation error |
| 3 | Permission error |

---

## Support

- [Development Guide](docs/DEVELOPMENT_GUIDE.md)
- [Context Rot Research](docs/background/context_rot.md)
- Issues: https://github.com/sanxzy/mandor/issues

---

**Built with ❤️ for AI Agent Workflows**

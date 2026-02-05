# Mandor - Deterministic Task Manager CLI for AI Agent Workflows

<p align="center">
  <img src="logo.png" alt="Mandor Logo" width="600">
</p>

<p align="center">
  <strong>Stop writing markdown plans. Start shipping features with deterministic task tracking.</strong>
</p>

<p align="center">
  <strong>Dependency-aware | Structured storage | CLI-native | Built for AI agents</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#why-mandor">Why Mandor</a> •
  <a href="#core-concepts">Core Concepts</a> •
  <a href="#commands">Commands</a> •
  <a href="#examples">Examples</a>
</p>

---

## Why Mandor

**No More Markdown Plan Files**

Traditional workflows scatter task state across markdown files, spreadsheets, and Slack messages. Dependencies are manual, status is fiction, and progress is invisible until code review.

Mandor brings **deterministic task management** to AI agent workflows:

- **Single Source of Truth**: All state in structured JSONL files—queryable, reproducible, auditable
- **Automatic Dependency Resolution**: Mark tasks done → dependents auto-transition to ready
- **Schema-Driven**: Enforce implementation steps, test cases, and library needs upfront
- **CLI-Native**: Works in terminal, scripts, and CI/CD pipelines
- **Dependency Tracking**: Full support for same-project and cross-project dependencies

## Overview

Mandor is a CLI tool for managing tasks, features, and issues in AI agent workflows:

- **Structured Storage**: All data in JSONL format with full audit trail
- **Real-Time Status**: Query tasks/issues by status (ready, blocked, in_progress)
- **Dependency Tracking**: Automatic status transitions when dependencies complete
- **Cross-Platform**: Go binary for macOS, Linux, Windows (arm64 & x64)

---


## Core Concepts

### Entity Types

| Type | Purpose | Status Values |
|------|---------|---------------|
| **Workspace** | Top-level container for all projects | (single instance per directory) |
| **Project** | Container for features and issues | (active/deleted) |
| **Feature** | Logical grouping of related tasks | draft, active, done, blocked, cancelled |
| **Task** | Work items within a feature | pending, ready, in_progress, done, blocked, cancelled |
| **Issue** | Problems, bugs, or improvement requests | open, ready, in_progress, resolved, wontfix, blocked, cancelled |

### Dependency Types

- **Task Dependencies**: One task can depend on multiple tasks
- **Feature Dependencies**: Features can depend on other features
- **Issue Dependencies**: Issues can depend on other issues

### Status Transitions

**Tasks:**
```
pending → {ready, in_progress, cancelled}
ready → {in_progress, cancelled}
in_progress → {done, blocked, cancelled}
blocked → {ready, cancelled}
done → (terminal)
cancelled → (terminal)
```

**Features:**
```
draft → {active, blocked, cancelled}
active → {done, blocked, cancelled}
blocked → {draft, active, done, cancelled}
done → {cancelled}
cancelled → {draft}
```

**Issues:**
```
open → {ready, in_progress, blocked, resolved, wontfix, cancelled}
ready → {in_progress, blocked, resolved, wontfix, cancelled}
in_progress → {blocked, resolved, wontfix, cancelled}
blocked → {ready, resolved, wontfix, cancelled}
resolved → (terminal, can reopen to any status)
wontfix → (terminal, can reopen to any status)
cancelled → (terminal, can reopen to any status)
```

---


## Installation

### Install with curl

```bash
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh
mandor --help
```

### Install from npm

```bash
npm install -g @mandors/cli
mandor --help
```

---


## Quick Start

### 1. Initialize Workspace

```bash
mandor init "My Project"
```

### 2. Create a Project

```bash
mandor project create api --name "API Development" \
  --goal "Build REST API service with authentication and endpoints"
```

### 3. Create a Feature

```bash
mandor feature create "Authentication" --project api \
  --goal "Implement JWT-based authentication with login and refresh flows for secure API access" \
  --scope backend
```

### 4. Create Tasks with Dependencies

```bash
# Create first task (no dependencies)
mandor task create api-feature-xxxx "JWT Parser" \
  --goal "Parse and validate JWT tokens in incoming requests with expiry and signature verification" \
  --implementation-steps "Setup crypto library|Add token validation|Handle expiry|Return errors" \
  --test-cases "Valid token accepted|Expired token rejected|Invalid signature rejected" \
  --library-needs "golang-jwt" \
  --priority P1

# Create dependent task (depends on JWT Parser)
mandor task create api-feature-xxxx "Login Endpoint" \
  --goal "Accept user credentials and return JWT token with refresh token flow" \
  --implementation-steps "Setup endpoint|Validate credentials|Generate JWT|Return tokens" \
  --test-cases "Valid creds return token|Invalid creds rejected|Tokens properly formatted" \
  --depends-on api-feature-xxxx-task-xxxx \
  --priority P1
```

### 5. View Task Progress

```bash
# See all tasks in feature with visualization
mandor track feature api-feature-xxxx

# Get task details
mandor task detail <task-id>
```

### 6. Mark Tasks Complete

```bash
# Get task ID from track output
mandor task update <task-id> --status in_progress
mandor task update <task-id> --status done

# Dependent tasks auto-transition to "ready"
mandor track feature api-feature-xxxx  # Now shows "Login Endpoint" as ready
```

---


## Commands Reference

### Workspace Commands

```bash
# Initialize a new workspace
mandor init [--workspace-name <name>] [-y]

# View workspace and project status
mandor status [--project <id>] [--summary] [--json]

# Manage configuration
mandor config get <key>
mandor config set <key> <value>
mandor config list
mandor config reset <key>

# Display all commands and best practices
mandor populate

# Show version
mandor version

# Generate shell completions
mandor completion [bash|zsh|fish]

# AI-assisted documentation
mandor ai --help
```

### Track Commands

```bash
# Track workspace status
mandor track

# Track project status
mandor track project <project-id>

# Track feature with tasks
mandor track feature <feature-id> [--verbose]

# Track specific task
mandor track task <task-id>

# Track issue
mandor track issue <issue-id>
```

### Session Commands

```bash
# Add a progress note (AI agents use this to track work)
mandor session note "Completed v0.4.4 release and testing"

# Read recent session notes (last 50 entries by default)
mandor session note --read

# Read more notes with offset
mandor session note --read --offset 100
```

### Project Commands

```bash
# Create a project
mandor project create <id> --name <name> --goal <goal> [OPTIONS]

# Show project details
mandor project detail <project-id>

# Update project
mandor project update <project-id> [--name <name>] [--goal <goal>] [--status <status>]
```

### Feature Commands

```bash
# Create a feature
mandor feature create <name> --project <id> --goal <goal> [--scope <scope>] [--priority <priority>]

# List features
mandor feature list --project <id>

# Show feature details
mandor feature detail <feature-id> --project <id>

# Update feature
mandor feature update <feature-id> --project <id> [--name <text>] [--goal <goal>] [--scope <scope>] [--priority <priority>] [--status <status>] [--depends <ids>] [--cancel --reason <text>] [--reopen] [--dry-run]
```

### Task Commands

```bash
# Create a task
mandor task create <feature_id> <name> --goal <goal> \
  --implementation-steps <steps> --test-cases <cases> \
  --derivable-files <files> --library-needs <libs> \
  [--priority <priority>] [--depends-on <ids>]

# Show task details
mandor task detail <task-id>

# Update task
mandor task update <task-id> [--name <text>] [--goal <goal>] [--priority <priority>] \
  [--status <status>] [--depends-add <ids>] [--depends-remove <ids>] [--cancel --reason <text>] [--dry-run]
```

### Issue Commands

```bash
# Create an issue
mandor issue create <name> --project <id> --type <type> --goal <goal> \
  --affected-files <files> --affected-tests <tests> \
  --implementation-steps <steps> [--priority <priority>] [--depends-on <ids>] [--library-needs <libs>]

# Show issue details
mandor issue detail <issue-id> --project <id>

# Update issue
mandor issue update <issue-id> [--name <text>] [--goal <goal>] [--priority <priority>] \
  [--type <type>] [--status <status>] [--start] [--resolve] [--wontfix] [--reason <text>] [--cancel --reason <text>] [--dry-run]
```

### AI Commands

```bash
# AI-assisted documentation generation
mandor ai agents
mandor ai claude
```

---


## Common Workflows

### Replace This (Markdown Plan Files)

```markdown
# PLAN.md
## Phase 1: Authentication
- [ ] JWT parser (depends on cryptography)
- [ ] Login endpoint (depends on JWT parser)
- [ ] Refresh token (depends on JWT parser)

Status: Last updated 3 days ago (probably stale!)
```

### With This (Mandor)

```bash
# Create structured plan
mandor feature create "Authentication" --project api \
  --goal "Implement JWT and login endpoints" \
  --scope backend

# Create tasks with explicit dependencies
mandor task create api-feature-xxx "JWT Parser" \
  --goal "Validate JWT tokens..." \
  --implementation-steps "Step 1|Step 2" \
  --test-cases "Test invalid tokens|Test expired" \
  --derivable-files "jwt.go|jwt_test.go" \
  --library-needs "jsonwebtoken" \
  --priority P1

mandor task create api-feature-xxx "Login Endpoint" \
  --goal "Accept credentials and return JWT..." \
  --depends-on api-feature-xxx-task-001 \
  --priority P1

# Real-time progress queries
mandor track feature api-feature-xxx                 # See all tasks and status
mandor track task api-feature-xxx-task-001          # See specific task details
```

**Benefits:**
- No file sync required
- Dependencies auto-validated
- Blocking tasks auto-detected
- Structured JSONL storage
- Queryable via CLI or JSON
- Works in CI/CD pipelines

### Dependency Management

```bash
# View all projects and their status
mandor status

# Check a specific project
mandor status --project api

# View feature dependencies and progress
mandor track project api

# Create tasks with dependencies
mandor task create api-feature-xxx "Task" \
  --goal "..." \
  --implementation-steps "..." \
  --test-cases "..." \
  --derivable-files "..." \
  --library-needs "..." \
  --depends-on api-feature-xxx-task-001

# See all feature tasks with status
mandor track feature api-feature-xxx

# Mark as done (auto-unblocks dependents)
mandor task update api-feature-xxx-task-001 --status done

# Verify dependents auto-transitioned to ready
mandor track feature api-feature-xxx
```

### Issue Tracking

```bash
# Create a bug issue
mandor issue create "Memory leak in auth handler" \
  --project api \
  --type bug \
  --priority P0 \
  --goal "Goroutine not cleaned up in token refresh handler..." \
  --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
  --affected-tests "src/handlers/auth_test.go" \
  --implementation-steps "Identify|Fix|Add tests|Verify" \
  --library-needs "none"

# View issue details
mandor issue detail api-issue-abc123

# Start working on an issue
mandor issue update api-issue-abc123 --start

# Mark as resolved
mandor issue update api-issue-abc123 --resolve

# Mark as won't fix with reason
mandor issue update api-issue-abc123 --wontfix --reason "Working as intended"

# See project issues with track
mandor track project api
```

### Configuration

```bash
# Set default priority
mandor config set default_priority P2

# Enable strict mode
mandor config set strict_mode true

# View all configuration
mandor config list

# Get specific value
mandor config get default_priority

# Reset to default
mandor config reset default_priority
```

---


## Best Practices

### 1. Use Meaningful IDs

Project and feature IDs should be:
- Short but descriptive
- Lowercase with hyphens
- Consistent naming convention

```bash
# Good
mandor project create user-auth
mandor feature create jwt-tokens

# Avoid
mandor project create p1
mandor feature create f123
```

### 2. Write Clear Goals (min char requirement enforced)

Goals should include:
- What is being built/fixed
- Why it matters
- Technical requirements
- Acceptance criteria

```bash
# Good
--goal "Implement JWT-based authentication with login and refresh flows for secure API access"

# Avoid
--goal "Add authentication"
```

### 3. Use Scopes for Features

Organize by scope:
- `frontend`, `backend`, `fullstack`
- `cli`, `desktop`, `android`, `flutter`, `react-native`, `ios`, `swift`

```bash
mandor feature create "Login UI" --project api --scope frontend
mandor feature create "Login API" --project api --scope backend
```

### 4. Keep Dependencies Shallow

Deep dependency chains (>5 levels) are hard to manage. Consider breaking into smaller features.

```bash
# Good: tasks depend on other tasks in same feature
mandor task create api-feature-xxxx "Task B" \
  --goal "..." \
  --implementation-steps "..." \
  --test-cases "..." \
  --depends-on api-feature-xxxx-task-xxxx

# Consider splitting if: task chains exceed 5 levels
```

### 5. Use Issues for Bugs, Tasks for Features

- **Tasks**: Feature work, implementation, refactoring
- **Issues**: Bugs, improvements, technical debt, security, performance

```bash
# Feature work
mandor task create api-feature-xxxx "Add OAuth2" \
  --goal "..." \
  --implementation-steps "..." \
  --test-cases "..." \
  --library-needs "oauth2-lib"

# Bug fix
mandor issue create "Fix auth timeout" --project api --type bug
```

### 6. Document Cancellation Reasons

Always provide clear reasons when cancelling:

```bash
mandor task update api-feature-xxxx-task-xxxx --cancel --reason "Superseded by feature X"
mandor feature update api-feature-xxxx --project api --cancel --reason "Sticking with JWT, OAuth2 adds too much complexity"
```

### 7. Use Pipe Separators For Lists

For flags accepting multiple values, use pipe separators:

```bash
# Implementation steps
--implementation-steps "Step 1|Step 2|Step 3"

# Test cases
--test-cases "Case 1|Case 2|Case 3"

# Dependencies
--depends-on task-1|task-2|task-3
```

### 8. Use --dry-run Before Major Changes

Before making significant updates, preview with `--dry-run`:

```bash
mandor task update api-feature-xxxx-task-xxxx --status done --dry-run
mandor feature update api-feature-xxxx --project api --cancel --reason "..." --dry-run
```

### 9. Dependency Auto-Resolution

- Mark task done → dependents auto-transition to ready
- Mark issue resolved → dependents auto-transition to ready
- Manual block → must manually unblock

### 10. Configure Early, Rarely Change

Configure workspace defaults at the start:

```bash
mandor init "Project Name"
mandor config set default_priority P2
mandor config set strict_mode true
```

### 11. Review Status Regularly

```bash
# Workspace overview
mandor status

# Project summary
mandor status --project api

# See feature progress
mandor track project api

# See feature tasks
mandor track feature api-feature-xxxx

# See task details
mandor track task api-feature-xxxx-task-xxxx
```

---


## Troubleshooting

### "Command not found"

Ensure mandor is in your PATH:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### "Project not found"

Check the project ID and ensure you're in the correct workspace:

```bash
mandor status
```

### "Entity not found"

Verify the entity ID exists:

```bash
mandor track feature <feature-id>
mandor track project <project-id>
```

### "Cross-project dependency detected"

The project doesn't allow cross-project dependencies:

```bash
# Check project config
mandor project detail <project-id>

# Create new project with cross-project enabled
mandor project create <id> --name "..." --goal "..." --task-dep cross_project_allowed
```

### "Invalid status transition"

The transition isn't allowed by the state machine:

```bash
# Tasks: pending → ready → in_progress → done
# Features: draft → active → done
# Issues: open → ready → in_progress → resolved
```

### "Cannot create task for cancelled feature"

Reopen the feature first:

```bash
mandor feature update <feature-id> --project <project-id> --reopen
```

---


## Configuration Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `default_priority` | string | P3 | Default priority for new entities (P0-P5) |
| `strict_mode` | boolean | false | Enable strict dependency validation |
| `goal.lengths.project` | integer | 500 | Min chars for project goal |
| `goal.lengths.feature` | integer | 300 | Min chars for feature goal |
| `goal.lengths.task` | integer | 500 | Min chars for task goal |
| `goal.lengths.issue` | integer | 200 | Min chars for issue goal |

---


## File Structure

```
.mandor/
├── workspace.json          # Workspace metadata
├── config.json             # Workspace configuration
├── session-notes.jsonl     # AI agent session progress notes (NDJSON)
└── projects/
    └── <project-id>/
        ├── project.json    # Project metadata
        ├── features.jsonl  # Feature records
        ├── tasks.jsonl     # Task records
        └── issues.jsonl    # Issue records
```

**Session Notes Format (session-notes.jsonl):**
```json
{"timestamp":"2026-02-04T12:45:00Z","note":"Completed v0.4.4 release and testing"}
{"timestamp":"2026-02-04T14:20:00Z","note":"Started performance optimization - blocked on benchmarks"}
```

---


## Support

- Issues: https://github.com/sanxzy/mandor/issues
- Documentation: `/docs` directory
- Repository: https://github.com/sanxzy/mandor

---


**Built for AI Agent Workflows**

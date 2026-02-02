# Mandor - Event-Based Task Manager CLI for AI Agent Workflows

<p align="center">
  <img src="logo.png" alt="Mandor Logo" width="600">
</p>

<p align="center">
  <strong>Stop writing markdown plans. Start shipping features with deterministic task tracking.</strong>
</p>

<p align="center">
  <strong>Event-sourced | Dependency-aware | CLI-native | Built for AI agents</strong>
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

- **Single Source of Truth**: All state in `events.jsonl`—queryable, reproducible, auditable
- **Automatic Dependency Resolution**: Mark tasks done → dependents auto-transition to ready
- **Schema-Driven**: Enforce implementation steps, test cases, and library needs upfront
- **CLI-Native**: Works in terminal, scripts, and CI/CD pipelines
- **Event-Sourced**: Full audit trail of every status change

## Overview

Mandor is a CLI tool for managing tasks, features, and issues in AI agent workflows:

- **Event-Based Architecture**: All changes logged in `events.jsonl` with immutable timestamps
- **JSONL Format**: Deterministic, append-only storage for reproducibility
- **Dependency Tracking**: Automatic status transitions when dependencies complete
- **Cross-Platform**: Go binary for macOS, Linux, Windows (arm64 & x64)

---

## Core Concepts

### Entity Hierarchy

```
Workspace
  └── Projects
        └── Features
              └── Tasks
        └── Issues
```

### Entity Types

| Type | Purpose | Status Values |
|------|---------|---------------|
| **Task** | Work items within a feature | pending, ready, in_progress, done, blocked, cancelled |
| **Feature** | Logical grouping of related tasks | draft, active, done, blocked, cancelled |
| **Issue** | Problems, bugs, or improvement requests | open, ready, in_progress, resolved, wontfix, blocked, cancelled |

### Dependency Types

- **Task Dependencies**: One task can depend on multiple tasks
- **Feature Dependencies**: Features can depend on other features
- **Issue Dependencies**: Issues can depend on other issues

### Status Transitions

**Tasks:**
```
pending → ready → in_progress → done
pending → ready → blocked → ready
ready → cancelled
cancelled → ready (reopen)
```

**Features:**
```
draft → active → done
draft → blocked (dependency not done)
draft → cancelled
cancelled → draft (reopen)
```

**Issues:**
```
open → ready → in_progress → resolved
open → ready → in_progress → wontfix
open → ready → blocked → ready
resolved → ready (reopen)
wontfix → ready (reopen)
open → cancelled
cancelled → open (reopen)
```

---

## Installation

### Build from Source

```bash
git clone https://github.com/sanxzy/mandor.git
cd mandor
go build -o ./mandor ./cmd/mandor
```

### Use from Binaries

```bash
./mandor --help
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
mandor task create "JWT Parser" --feature api-feature-xxx \
  --goal "Parse and validate JWT tokens in incoming requests with expiry and signature verification" \
  --implementation-steps "Setup crypto library|Add token validation|Handle expiry|Return errors" \
  --test-cases "Valid token accepted|Expired token rejected|Invalid signature rejected" \
  --derivable-files "jwt_validator.go|jwt_test.go" \
  --library-needs "golang-jwt" \
  --priority P1

# Create dependent task (depends on JWT Parser)
mandor task create "Login Endpoint" --feature api-feature-xxx \
  --goal "Accept user credentials and return JWT token with refresh token flow" \
  --implementation-steps "Setup endpoint|Validate credentials|Generate JWT|Return tokens" \
  --test-cases "Valid creds return token|Invalid creds rejected|Tokens properly formatted" \
  --derivable-files "login_handler.go|login_test.go" \
  --library-needs "none" \
  --depends-on api-task-xxx-001 \
  --priority P1
```

### 5. View Task Progress

```bash
# See all tasks in feature
mandor task list --feature api-feature-xxx

# See tasks ready to work on
mandor task ready --feature api-feature-xxx

# See blocked/waiting tasks
mandor task blocked --feature api-feature-xxx
```

### 6. Mark Tasks Complete

```bash
# Get task ID from list
mandor task update <task-id> --status in_progress
mandor task update <task-id> --status done

# Dependent tasks auto-transition to "ready"
mandor task ready --feature api-feature-xxx  # Now shows "Login Endpoint" as ready
```

---

## Commands Reference

### Workspace Commands

```bash
# Initialize a new workspace
mandor init [--workspace-name <name>] [-y]

# View workspace and project status
mandor status [--project <id>] [--summary] [--json]

# Display feature summary grouped by priority
mandor summary [--project <id>]

# Manage configuration
mandor config get <key>
mandor config set <key> <value>
mandor config list
mandor config reset <key>

# Display all commands and best practices
mandor populate [--markdown|--json]

# Show version
mandor version

# Generate shell completions
mandor completion [bash|zsh|fish]
```

### Project Commands

```bash
# Create a project
mandor project create <id> --name <name> --goal <goal> [OPTIONS]

# List projects
mandor project list [--json]

# Show project details
mandor project detail <project-id>

# Update project
mandor project update <project-id> --name <name> [--goal <goal>]

# Delete project (soft delete)
mandor project delete <project-id> [--hard]

# Reopen a soft-deleted project
mandor project reopen <project-id>
```

### Feature Commands

```bash
# Create a feature
mandor feature create <name> --project <id> --goal <goal> [--scope <scope>] [--priority <priority>] [--depends <ids>]

# List features
mandor feature list --project <id> [--json] [--include-deleted]

# Show feature details
mandor feature detail <feature-id> --project <id>

# Update feature
mandor feature update <feature-id> --project <id> [--name <text>] [--goal <goal>] [--scope <scope>] [--priority <priority>] [--status <status>] [--depends <ids>] [--cancel --reason <text>] [--reopen] [--dry-run]

# List features by status filter
mandor feature list --project <id> [--status <status>]
```

### Task Commands

```bash
# Create a task
mandor task create <name> --feature <id> --goal <goal> \
  --implementation-steps <steps> --test-cases <cases> \
  --derivable-files <files> --library-needs <libs> \
  [--priority <priority>] [--depends-on <ids>] [-y]

# List tasks in a feature
mandor task list --feature <id> [--status <status>] [--priority <priority>] [--json] [--include-deleted] [--sort <field>] [--order <asc|desc>]

# Show ready tasks
mandor task ready --feature <id> [--priority <priority>] [--json]

# Show blocked tasks
mandor task blocked --feature <id> [--priority <priority>] [--json]

# Show task details
mandor task detail <task-id>

# Update task
mandor task update <task-id> [--name <text>] [--goal <goal>] [--priority <priority>] \
  [--implementation-steps <steps>] [--test-cases <cases>] [--derivable-files <files>] \
  [--library-needs <libs>] [--status <status>] [--depends <ids>] [--depends-add <ids>] \
  [--depends-remove <ids>] [--cancel --reason <text>] [--reopen] [--dry-run] [--force]
```

### Issue Commands

```bash
# Create an issue
mandor issue create <name> --project <id> --type <type> --goal <goal> \
  --affected-files <files> --affected-tests <tests> \
  --implementation-steps <steps> \
  [--priority <priority>] [--depends-on <ids>] [--library-needs <libs>] [-y]

# List issues in project
mandor issue list [--project <id>] [--type <type>] [--status <status>] [--priority <priority>] [--json] [--sort <field>] [--order <asc|desc>] [--verbose]

# Show ready issues
mandor issue ready [--project <id>] [--type <type>] [--priority <priority>] [--json]

# Show blocked issues
mandor issue blocked [--project <id>] [--type <type>] [--priority <priority>] [--json]

# Show issue details
mandor issue detail <issue-id> [--project <id>]

# Update issue
mandor issue update <issue-id> [--project <id>] [--name <text>] [--goal <goal>] [--type <type>] [--priority <priority>] \
  [--status <status>] [--reason <text>] [--depends-on <ids>] [--depends-add <ids>] [--depends-remove <ids>] \
  [--affected-files <files>] [--affected-tests <tests>] [--implementation-steps <steps>] [--library-needs <libs>] \
  [--start] [--resolve] [--wontfix] [--reopen] [--cancel] [--force] [--dry-run]
```

### AI Commands

```bash
# Generate AGENTS.md for multi-agent coordination
mandor ai agents

# Generate CLAUDE.md for the project
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
mandor task create "JWT Parser" --feature auth-feature-id \
  --goal "Validate JWT tokens..." \
  --implementation-steps "Step 1|Step 2" \
  --test-cases "Test invalid tokens|Test expired" \
  --derivable-files "jwt.go|jwt_test.go" \
  --library-needs "jsonwebtoken" \
  --priority P1

mandor task create "Login Endpoint" --feature auth-feature-id \
  --goal "Accept credentials and return JWT..." \
  --depends-on jwt-parser-task-id \
  --priority P1

# Real-time progress queries
mandor task ready --feature auth-feature-id           # See what's available now
mandor task blocked --feature auth-feature-id         # See what's waiting
```

**Benefits:**
- No file sync required
- Dependencies auto-validated
- Blocking tasks auto-detected
- Reproducible state (`events.jsonl`)
- Queryable via CLI or JSON
- Works in CI/CD pipelines

### Dependency Management

```bash
# View all projects and their status
mandor status

# Check a specific project
mandor status --project api

# View feature dependencies
mandor feature list --project api

# Create tasks with dependencies
mandor task create "Step 2" --feature feature-id \
  --goal "..." \
  --implementation-steps "..." \
  --test-cases "..." \
  --derivable-files "..." \
  --library-needs "..." \
  --depends-on task-id-1|task-id-2

# See what's blocking progress
mandor task blocked --feature feature-id

# Mark as done (auto-unblocks dependents)
mandor task update task-id --status done

# Dependents auto-transition to ready
mandor task ready --feature feature-id
```

### Issue Tracking

```bash
# Create a bug issue
mandor issue create "Fix memory leak in auth" \
  --project api \
  --type bug \
  --priority P0 \
  --goal "Goroutine leak in token refresh handler..." \
  --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
  --affected-tests "src/handlers/auth_test.go" \
  --implementation-steps "Identify leak|Add cleanup|Test|Verify"

# List open issues
mandor issue list --project api --status open

# Filter by type and priority
mandor issue list --project api --type bug --priority P0

# See ready issues
mandor issue ready --project api --type bug

# Start working
mandor issue update issue-id --start

# Mark as resolved
mandor issue update issue-id --resolve

# Mark as won't fix with reason
mandor issue update issue-id --wontfix --reason "Working as intended"
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

### 2. Write Clear Goals

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
mandor task create "Task B" --feature feature-id --depends-on task-a-id

# Consider splitting if: task chains exceed 5 levels
```

### 5. Use Issues for Bugs, Tasks for Features

- **Tasks**: Feature work, implementation, refactoring
- **Issues**: Bugs, improvements, technical debt, security, performance

```bash
# Feature work
mandor task create "Add OAuth2" --feature api-auth

# Bug fix
mandor issue create "Fix auth timeout" --project api --type bug
```

### 6. Document Cancellation Reasons

Always provide clear reasons when cancelling:

```bash
mandor task update task-id --cancel --reason "Superseded by feature X"
mandor feature update feature-id --project api --cancel --reason "Sticking with JWT, OAuth2 adds too much complexity"
```

### 7. Use Pipe-Separated Lists

For flags accepting multiple values, use pipe separators:

```bash
# Implementation steps
--implementation-steps "Step 1|Step 2|Step 3"

# Test cases
--test-cases "Case 1|Case 2|Case 3"

# Dependencies
--depends-on task-1|task-2|task-3
```

### 8. Use --dry-run to Preview Changes

Before making significant updates, preview with `--dry-run`:

```bash
mandor task update task-id --status done --cancel --dry-run
mandor feature update feature-id --project api --cancel --reason "..." --dry-run
```

### 9. Set Configuration Early

Configure workspace defaults at the start:

```bash
mandor init "Project Name"
mandor config set default_priority P2
mandor config set strict_mode true
```

### 10. Review Status Regularly

```bash
# Workspace overview
mandor status

# Project summary
mandor status --project api

# Feature priorities
mandor summary --project api

# See ready work
mandor task ready --feature feature-id
mandor issue ready --project api

# See blockers
mandor task blocked --feature feature-id
mandor issue blocked --project api
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
mandor project list
```

### "Entity not found"

Verify the entity ID exists:

```bash
mandor task list --feature <feature-id>
mandor issue list --project <project-id>
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
mandor feature update <id> --project <pid> --reopen
```

---

## Configuration Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `default_priority` | string | P3 | Default priority for new entities (P0-P5) |
| `strict_mode` | boolean | false | Enable strict dependency validation |

---

## File Structure

```
.mandor/
├── workspace.json          # Workspace metadata
├── config.json             # Workspace configuration
└── projects/
    └── <project-id>/
        ├── project.json    # Project metadata
        ├── features.jsonl  # Feature event log
        ├── tasks.jsonl     # Task event log
        ├── issues.jsonl    # Issue event log
        └── events.jsonl    # All events
```

---

## Support

- Issues: https://github.com/sanxzy/mandor/issues
- Documentation: `/docs` directory
- Repository: https://github.com/sanxzy/mandor

---

**Built for AI Agent Workflows**

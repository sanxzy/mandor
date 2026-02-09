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
| **Workspace** | Top-level container for all backlogs | (single instance per directory) |
| **Backlog** | Container for features and issues | initial, active, done, deleted |
| **Feature** | Logical grouping of related tasks | draft, active, done, blocked, cancelled |
| **Task** | Work items within a feature | pending, ready, in_progress, done, blocked, cancelled |
| **Issue** | Problems, bugs, or improvement requests | open, ready, in_progress, resolved, wontfix, blocked, cancelled |

> **Note**: Backlogs replace Projects. The system automatically migrates existing projects to backlogs on first use.

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

This initializes your workspace and automatically migrates any existing projects to backlogs.

### 2. Create a Backlog

```bash
mandor backlog create api --name "API Development" \
  --goal "Build REST API service with authentication and endpoints. This backlog covers the complete implementation of secure API endpoints including authentication, authorization, and comprehensive testing to ensure production-ready quality."
```

### 3. Create a Feature

```bash
mandor feature create auth-feature --backlog api \
  --name "JWT Authentication" \
  --goal "Implement JWT-based authentication with login and refresh flows for secure API access" \
  --scope backend
```

### 4. Create Tasks with Dependencies

```bash
# Create first task (no dependencies)
mandor task create jwt-parser auth-feature \
  --backlog api \
  --goal "Parse and validate JWT tokens in incoming requests with expiry and signature verification including error handling" \
  --implementation-steps "Setup crypto library|Add token validation|Handle expiry|Return errors" \
  --test-cases "Valid token accepted|Expired token rejected|Invalid signature rejected" \
  --library-needs "golang-jwt" \
  --priority P1

# Create dependent task (depends on JWT Parser)
mandor task create login-endpoint auth-feature \
  --backlog api \
  --goal "Accept user credentials and return JWT token with refresh token flow" \
  --implementation-steps "Setup endpoint|Validate credentials|Generate JWT|Return tokens" \
  --test-cases "Valid creds return token|Invalid creds rejected|Tokens properly formatted" \
  --depends-on jwt-parser \
  --priority P1
```

### 5. View Task Progress

```bash
# See all tasks in feature with visualization
mandor track feature auth-feature

# Get task details
mandor task detail jwt-parser
```

### 6. Mark Tasks Complete

```bash
# Get task ID from track output
mandor task update jwt-parser --status in_progress
mandor task update jwt-parser --status done

# Dependent tasks auto-transition to "ready"
mandor track feature auth-feature  # Now shows "login-endpoint" as ready
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

# Track backlog status
mandor track backlog <backlog-id>

# Track feature with tasks
mandor track feature <feature-id> [--verbose]

# Track specific task
mandor track task <task-id>

# Track issue
mandor track issue <issue-id>
```

### Backlog Commands

```bash
# Create a new backlog
mandor backlog create <backlog-id> --name "<name>" --goal "<goal>"

# Show backlog details
mandor backlog detail <backlog-id>

# Update backlog metadata
mandor backlog update <backlog-id> [--name "<name>"] [--goal "<goal>"]

# Delete a backlog (soft-delete)
mandor backlog delete <backlog-id> [-y]

# Reopen a deleted backlog
mandor backlog reopen <backlog-id> [-y]

# List all backlogs
mandor backlog list [--deleted]
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
mandor task create api-feature-xxxx "JWT Parser" \
  --goal "Validate JWT tokens..." \
  --implementation-steps "Step 1|Step 2" \
  --test-cases "Test invalid tokens|Test expired" \
  --library-needs "jsonwebtoken" \
  --priority P1

mandor task create api-feature-xxxx "Login Endpoint" \
  --goal "Accept credentials and return JWT..." \
  --depends-on api-feature-xxxx-task-xxxx \
  --priority P1

# Real-time progress queries
mandor track feature api-feature-xxxx                 # See all tasks and status
mandor track task api-feature-xxxx-task-xxxx         # See specific task details
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
mandor task create api-feature-xxxx "Task" \
  --goal "..." \
  --implementation-steps "..." \
  --test-cases "..." \
  --library-needs "..." \
  --depends-on api-feature-xxxx-task-xxxx

# See all feature tasks with status
mandor track feature api-feature-xxxx

# Mark as done (auto-unblocks dependents)
mandor task update api-feature-xxxx-task-xxxx --status done

# Verify dependents auto-transitioned to ready
mandor track feature api-feature-xxxx
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

Minimum goal lengths:
- **Backlog**: 500 characters
- **Feature**: 300 characters
- **Task**: 500 characters
- **Issue**: 200 characters

```bash
# Good
--goal "Implement JWT-based authentication with login and refresh flows for secure API access with proper error handling and comprehensive test coverage for all edge cases."

# Avoid
--goal "Add authentication"
```

### 3. Use Scopes for Features

Organize features by scope to group related work:
- `frontend`, `backend`, `fullstack`
- `cli`, `desktop`, `android`, `flutter`, `react-native`, `ios`, `swift`

```bash
mandor feature create login-ui --backlog api --scope frontend \
  --goal "Create responsive login user interface with form validation and error handling for web and mobile platforms"

mandor feature create login-api --backlog api --scope backend \
  --goal "Implement login API endpoint with JWT token generation and refresh token management for secure user authentication"
```

### 4. Keep Dependencies Shallow

Deep dependency chains (>5 levels) are hard to manage. Consider breaking into smaller features or backlogs.

```bash
# Good: tasks depend on other tasks in same feature/backlog
mandor task create task-b auth-feature \
  --backlog api \
  --goal "Implement token refresh logic with automatic expiration handling and validation" \
  --implementation-steps "Parse refresh token|Validate expiration|Generate new token|Return response" \
  --test-cases "Valid refresh accepted|Expired refresh rejected|Token properly formatted" \
  --depends-on jwt-parser

# Consider splitting if: task chains exceed 5 levels
# Move related features to a new backlog if single backlog becomes too large
```

### 5. Use Issues for Bugs, Tasks for Features

- **Tasks**: Feature work, implementation, refactoring within a feature
- **Issues**: Bugs, improvements, technical debt, security, performance within a backlog

```bash
# Feature work - task within a feature
mandor task create oauth-integration auth-feature \
  --backlog api \
  --goal "Integrate OAuth2 authentication provider with existing JWT system for third-party integrations" \
  --implementation-steps "Setup OAuth2 client|Add provider endpoints|Implement token exchange|Update user model" \
  --test-cases "OAuth flow complete|Token exchange succeeds|User created|Session established" \
  --library-needs "oauth2-lib"

# Bug fix - issue within a backlog
mandor issue create auth-timeout --backlog api --type bug \
  --goal "Fix authentication timeout issue where tokens expire prematurely causing user sessions to terminate unexpectedly"
```

### 6. Document Cancellation Reasons

Always provide clear reasons when cancelling:

```bash
mandor task update jwt-parser --status cancelled --reason "Superseded by feature X"
mandor feature update oauth-feature --backlog api --status cancelled --reason "Sticking with JWT, OAuth2 adds too much complexity"
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

### 8. Use -y Flag for Automatic Confirmation

For non-interactive automation, use the `-y` or `--yes` flag:

```bash
# Create without prompt
mandor backlog create api --name "API" --goal "..." --yes

# Delete without prompt
mandor backlog delete api --yes

# Reopen without prompt
mandor backlog reopen api --yes
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

# Backlog summary
mandor status --backlog api

# See backlog progress
mandor track backlog api

# See feature tasks
mandor track feature auth-feature

# See task details
mandor track task jwt-parser
```

---


## Troubleshooting

### "Command not found"

Ensure mandor is in your PATH:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### "Backlog not found"

Check the backlog ID and ensure you're in the correct workspace:

```bash
mandor status
mandor backlog list
```

### "Entity not found"

Verify the entity ID exists:

```bash
mandor track feature <feature-id>
mandor track backlog <backlog-id>
mandor track task <task-id>
```

### "Cross-backlog dependency detected"

If you need dependencies across backlogs, plan your backlog structure accordingly:

```bash
# Check backlog details
mandor backlog detail <backlog-id>

# Dependencies are typically scoped within a backlog or feature
# Consider merging related work into the same backlog if dependencies are complex
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
mandor feature update <feature-id> --status active
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
├── backlogs/               # Backlog containers (replaces projects/)
│   └── <backlog-id>/
│       ├── backlog.json    # Backlog metadata
│       ├── schema.json     # Backlog dependency rules
│       ├── features.jsonl  # Feature records
│       ├── tasks.jsonl     # Task records
│       └── issues.jsonl    # Issue records
└── .projects_backup_*      # Automatic backup of old projects/ directory
```

**Note**: The `projects/` directory is automatically migrated to `backlogs/` on first use. Backups are retained for safety.

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

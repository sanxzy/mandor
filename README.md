# Mandor - Event-Based Task Manager CLI for AI Agent Workflows

<p align="center">
  <strong>Deterministic JSONL output | Streaming-native architecture | Schema-driven task management</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#core-concepts">Core Concepts</a> •
  <a href="#commands">Commands</a> •
  <a href="#examples">Examples</a>
</p>

---

## Overview

Mandor is a CLI tool for managing tasks, features, and issues in AI agent workflows:

- **Event-Based Architecture**: All changes logged in `events.jsonl`
- **JSONL Format**: Deterministic, append-only storage
- **Dependency Tracking**: Automatic status based on dependencies
- **Cross-Platform**: Go binary for macOS, Linux, Windows

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
| **Task** | Work items within a feature | ready, in_progress, done, blocked, cancelled |
| **Feature** | Logical grouping of related tasks | draft, active, done, blocked, cancelled |
| **Issue** | Problems, bugs, or improvement requests | ready, in_progress, resolved, wontfix, blocked |

### Dependency Types

- **Task Dependencies**: One task can depend on multiple tasks
- **Feature Dependencies**: Features can depend on other features
- **Issue Dependencies**: Issues can depend on other issues

### Status Transitions

**Tasks:**
```
ready → in_progress → done
ready → blocked (dependency not done)
blocked → ready (dependencies resolved)
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
ready → in_progress → resolved
ready → in_progress → wontfix
ready → blocked (dependency not done)
blocked → ready (dependencies resolved)
resolved → ready (reopen)
wontfix → ready (reopen)
```

---

## Installation

### Build from Source

```bash
git clone https://github.com/budisantoso/mandor.git
cd mandor
go build -o ./binaries/mandor ./cmd/mandor
```

### Use from Binaries

```bash
./binaries/mandor --help
```

---

## Quick Start

### 1. Initialize Workspace

```bash
export MANDOR_ENV=development
./binaries/mandor init --workspace-name "My Project"
```

### 2. Create a Project

```bash
./binaries/mandor project create api --name "API Development" \
  --goal "Build REST API endpoints"
```

### 3. Create a Feature

```bash
./binaries/mandor feature create "Auth Feature" --project api \
  --goal "Implement authentication system" \
  --scope backend
```

### 4. Create Tasks

```bash
./binaries/mandor task create "JWT Parser" \
  --feature api-feature-xxx \
  --goal "Parse and validate JWT tokens" \
  --implementation-steps "none" \
  --test-cases "none" \
  --derivable-files "none" \
  --library-needs "none"
```

---

## Commands Reference

### Workspace Commands

```bash
# Initialize a new workspace
mandor init --workspace-name "Name"

# View workspace status
mandor status

# Manage configuration
mandor config get <key>
mandor config set <key> <value>
mandor config reset <key>
```

### Project Commands

```bash
# Create a project
mandor project create <id> --name "Name" --goal "Description"

# List projects
mandor project list

# Show project details
mandor project detail <project-id>

# Delete a project
mandor project delete <project-id>
```

### Feature Commands

```bash
# Create a feature
mandor feature create "Name" --project <id> --goal "Description" [--scope <value>]

# List features
mandor feature list --project <id> [--scope <value>]

# Show feature details
mandor feature detail <feature-id> --project <id> [--include-deleted]

# Update feature
mandor feature update <id> --project <id> [--status <value>] [--cancel --reason] [--reopen]

# Delete feature (soft delete)
mandor feature update <id> --project <id> --cancel --reason "Reason"
```

### Task Commands

```bash
# Create a task
mandor task create "Name" --feature <id> -g "Goal" \
  --implementation-steps "step1|step2" \
  --test-cases "test1|test2" \
  --derivable-files "file1|file2" \
  --library-needs "lib1|lib2" \
  [--priority <P0-P5>] \
  [--depends-on <task-id>]

# List tasks
mandor task list --feature <id> [--status <value>] [--priority <value>]

# Show task details
mandor task detail <task-id> [--include-deleted] [--events] [--dependencies]

# Update task
mandor task update <id> [--status <value>] \
  [--cancel --reason] [--reopen] \
  [--depends <ids>] [--depends-add <ids>] [--depends-remove <ids>]

# Complete a task (in_progress → done)
mandor task update <id> --status in_progress
mandor task update <id> --status done
```

### Issue Commands

```bash
# Create an issue
mandor issue create "Name" --project <id> -t <type> -g "Goal" \
  --affected-files "file1" \
  --affected-tests "test1" \
  --implementation-steps "step1" \
  [--priority <P0-P5>] \
  [--depends-on <issue-id>]

# List issues
mandor issue list --project <id> [--status <value>] [--type <value>]

# Show issue details
mandor issue detail <issue-id> [--include-deleted]

# Update issue
mandor issue update <id> [--status <value>] [--start] [--resolve] \
  [--wontfix --reason] [--reopen] \
  [--depends-on <ids>]
```

---

## Configuration

### Available Config Keys

```bash
# Default priority for new entities
mandor config set default_priority P3

# Strict mode for dependency rules
mandor config set strict_mode true

# View current config
mandor config get default_priority
mandor config get strict_mode

# Reset to defaults
mandor config reset default_priority
```

### Priority Values

| Priority | Use Case |
|----------|----------|
| P0 | Critical / Security |
| P1 | High / Blocker |
| P2 | Medium-High |
| P3 | Medium (default) |
| P4 | Medium-Low |
| P5 | Low / Nice to have |

### Scope Values

Valid scope values for features: `frontend`, `backend`, `fullstack`, `cli`, `desktop`, `android`, `flutter`, `react-native`, `ios`, `swift`

### Issue Types

Valid issue types: `bug`, `improvement`, `debt`, `security`, `performance`

---

## Dependency Management

### Creating Dependencies

```bash
# Task depends on another task
mandor task create "Task B" --feature f1 \
  --depends-on task-a-id

# Issue depends on another issue
mandor issue create "Issue B" --project p1 \
  --depends-on issue-a-id

# Feature depends on another feature
mandor feature create "Feature B" --project p1 \
  --depends-on feature-a-id
```

### Managing Dependencies

```bash
# Add dependencies
mandor task update <task-id> --depends-add "id1|id2"

# Replace all dependencies
mandor task update <task-id> --depends "id1|id2"

# Remove dependencies
mandor task update <task-id> --depends-remove "id1"
```

### Dependency Behaviors

- **Auto-blocking**: Entities start `blocked` if dependencies aren't `done`
- **Auto-unblocking**: When a dependency becomes `done`, dependents automatically transition to `ready`
- **Cross-project**: Dependencies can span projects if enabled in project config

### Cross-Project Dependencies

Projects can be configured to allow or disallow cross-project dependencies:

```bash
# Allow cross-project task dependencies
mandor project create p1 --task-dep cross_project_allowed

# Restrict to same-project only (default)
mandor project create p2 --task-dep same_project_only
```

### Circular Dependency Prevention

Mandor automatically prevents:
- Self-dependencies (A → A)
- Two-node cycles (A → B → A)
- N-node cycles (A → B → C → A)
- Cross-project cycles

---

## Status Management

### Checking Status

```bash
# List blocked tasks
mandor task list --feature <id> --status blocked

# List ready tasks
mandor task list --feature <id> --status ready

# List all tasks
mandor task list --feature <id>
```

### Transitioning Status

```bash
# Task workflow
mandor task update <id> --status in_progress
mandor task update <id> --status done

# Feature workflow
mandor feature update <id> --project <pid> --status active
mandor feature update <id> --project <pid> --status done

# Issue workflow
mandor issue update <id> --start
mandor issue update <id> --resolve
# or
mandor issue update <id> --wontfix --reason "Reason"
```

### Cancel and Reopen

```bash
# Cancel (soft delete)
mandor task update <id> --cancel --reason "Why cancelled"
mandor feature update <id> --project <pid> --cancel --reason "Why cancelled"
mandor issue update <id> --wontfix --reason "Why wontfix"

# Reopen
mandor task update <id> --reopen
mandor feature update <id> --project <pid> --reopen
mandor issue update <id> --reopen
```

**Note**: Cancelled entities can be viewed with `--include-deleted` flag.

---

## Filtering and Querying

### Filter by Status

```bash
mandor task list --feature <id> --status ready
mandor task list --feature <id> --status blocked
mandor task list --feature <id> --status done
```

### Filter by Priority

```bash
mandor task list --feature <id> --priority P0
mandor task list --feature <id> --priority P2
```

### Filter by Scope (Features)

```bash
mandor feature list --project <id> --scope backend
mandor feature list --project <id> --scope frontend
```

### Filter by Type (Issues)

```bash
mandor issue list --project <id> --type bug
mandor issue list --project <id> --type security
```

---

## Event System

### Event Log Location

All events are stored in `.mandor/events.jsonl` in your workspace.

### Event Types

| Event | Description |
|-------|-------------|
| entity_created | New entity created |
| status_changed | Status transitioned |
| dependency_added | Dependency added |
| dependency_removed | Dependency removed |
| dependent_unblocked | Dependent entity became ready |
| entity_reopened | Cancelled entity reopened |
| entity_cancelled | Entity cancelled |

### Viewing Events

```bash
# View events for a specific entity
mandor task detail <id> --events

# Events are also shown in default detail view
mandor task detail <id>
```

---

## JSONL Format

Mandor uses JSONL (JSON Lines) for event storage:

```json
{"timestamp":"2026-02-01T10:00:00Z","type":"entity_created","entity":"task","id":"task-abc","name":"JWT Parser"}
{"timestamp":"2026-02-01T10:01:00Z","type":"status_changed","entity":"task","id":"task-abc","from":"ready","to":"in_progress"}
{"timestamp":"2026-02-01T10:02:00Z","type":"status_changed","entity":"task","id":"task-abc","from":"in_progress","to":"done"}
{"timestamp":"2026-02-01T10:02:00Z","type":"dependent_unblocked","entity":"task","id":"task-xyz","dependency":"task-abc"}
```

### Parsing Events

```bash
# View all events
cat .mandor/events.jsonl

# Filter for specific entity
grep "task-abc" .mandor/events.jsonl

# Count events by type
grep '"type":"status_changed"' .mandor/events.jsonl | wc -l
```

---

## Examples

### Complete Feature Workflow

```bash
# Setup
mandor init --workspace-name "API Project"
mandor project create api --name "API" --goal "Build REST API"
mandor feature create "Auth" --project api --goal "Authentication" --scope backend

# Create tasks with dependencies
mandor task create "JWT Parser" --feature auth-xxx \
  -g "Parse JWT tokens" --implementation-steps "none" \
  --test-cases "none" --derivable-files "none" --library-needs "none"

mandor task create "JWT Validator" --feature auth-xxx \
  -g "Validate JWT tokens" --implementation-steps "none" \
  --test-cases "none" --derivable-files "none" --library-needs "none" \
  --depends-on jwt-parser-task-id

# Execute workflow
mandor task update jwt-parser-id --status in_progress
mandor task update jwt-parser-id --status done

# Validator automatically becomes ready
mandor task update jwt-validator-id --status in_progress
mandor task update jwt-validator-id --status done

# Mark feature done
mandor feature update auth-xxx --project api --status active
mandor feature update auth-xxx --project api --status done
```

### Multi-Project Dependencies

```bash
# Create projects with cross-project enabled
mandor project create core --name "Core" --goal "Core library"
mandor project create api --name "API" --goal "API layer" --task-dep cross_project_allowed

# Create task in core
mandor feature create lib --project core --goal "Core library"
mandor task create "DB Connection" --feature lib-xxx \
  -g "Database connection" --implementation-steps "none" \
  --test-cases "none" --derivable-files "none" --library-needs "none"

# Create task in api depending on core task
mandor feature create endpoints --project api --goal "API endpoints"
mandor task create "User Endpoint" --feature endpoints-xxx \
  -g "User API endpoint" --implementation-steps "none" \
  --test-cases "none" --derivable-files "none" --library-needs "none" \
  --depends-on db-connection-task-id

# Complete core task, api task becomes ready
mandor task update db-connection-id --status in_progress
mandor task update db-connection-id --status done
```

### Issue Tracking with Dependencies

```bash
# Create blocking issues
mandor issue create "Data Validation Bug" --project api -t bug \
  -g "Fix data validation" --affected-files "validate.js" \
  --affected-tests "validate_test.js" --implementation-steps "step1"

mandor issue create "Cache Issue" --project api -t bug \
  -g "Fix cache issue" --affected-files "cache.js" \
  --affected-tests "cache_test.js" --implementation-steps "step1"

# Create dependent issue
mandor issue create "Database Crash" --project api -t bug \
  -g "Fix crash" --affected-files "db.js" \
  --affected-tests "db_test.js" --implementation-steps "step1" \
  --depends-on data-validation-id|cache-issue-id

# Resolve blocking issues
mandor issue update data-validation-id --resolve
mandor issue update cache-issue-id --resolve

# Dependent issue automatically becomes ready
mandor issue update database-crash-id --start
mandor issue update database-crash-id --resolve
```

### Cancel and Reopen Workflow

```bash
# Create feature and tasks
mandor feature create "Experiment" --project api --goal "Try something"
mandor task create "Try X" --feature experiment-xxx -g "Test X" \
  --implementation-steps "none" --test-cases "none" \
  --derivable-files "none" --library-needs "none"

# Cancel the feature
mandor feature update experiment-xxx --project api --cancel --reason "Not needed"

# Cannot create new tasks in cancelled feature
mandor task create "Try Y" --feature experiment-xxx -g "Test Y" \
  --implementation-steps "none" --test-cases "none" \
  --derivable-files "none" --library-needs "none"
# Error: "Cannot create task for cancelled feature"

# Reopen the feature
mandor feature update experiment-xxx --project api --reopen

# Can create tasks again
mandor task create "Try Y" --feature experiment-xxx -g "Test Y" \
  --implementation-steps "none" --test-cases "none" \
  --derivable-files "none" --library-needs "none"
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

### 2. Set Default Priority

Configure default priority in workspace:

```bash
mandor config set default_priority P3
```

### 3. Use Scopes for Features

Assign scopes to help filter and organize:

```bash
mandor feature create "Login UI" --project api --scope frontend
mandor feature create "Login API" --project api --scope backend
```

### 4. Keep Dependencies Shallow

Deep dependency chains (>5 levels) can be hard to manage. Consider breaking into smaller features.

### 5. Use Issues for Bugs, Tasks for Work

- **Tasks**: Work to be done (implementations, refactoring)
- **Issues**: Problems to be fixed or improvements to be made

### 6. Document Cancellation Reasons

Always provide clear reasons when cancelling:

```bash
mandor task update <id> --cancel --reason "Superseded by feature X"
```

### 7. Review Blocked Tasks Regularly

```bash
mandor task list --feature <id> --status blocked
```

### 8. Use Configuration for Consistency

Set up workspace configuration early:

```bash
mandor config set default_priority P2
mandor config set strict_mode true
```

---

## Troubleshooting

### "Command not found"

Ensure mandor is in your PATH:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### "Project not found"

Check the project ID and ensure you're in the correct workspace.

```bash
mandor project list
```

### "Entity not found"

Verify the entity ID exists:

```bash
mandor task list --feature <feature-id>
```

### "Cross-project dependency detected"

The project doesn't allow cross-project dependencies:

```bash
# Check project config
mandor project detail <project-id>

# Create new project with cross-project enabled
mandor project create <id> --task-dep cross_project_allowed
```

### "Invalid status transition"

The transition isn't allowed by the state machine:

```bash
# Tasks: ready -> in_progress -> done
# Features: draft -> active -> done
```

### "Cannot create task for cancelled feature"

Reopen the feature first:

```bash
mandor feature update <id> --project <pid> --reopen
```

---

## File Structure

```
mandor/
├── cmd/
│   └── mandor/           # CLI entry point
├── internal/
│   ├── cmd/              # Command implementations
│   ├── domain/           # Domain models
│   ├── service/          # Business logic
│   ├── repository/       # Data access
│   └── events/           # Event handling
├── .mandor/              # Workspace (created by init)
│   ├── workspace.json    # Workspace config
│   └── projects/         # Projects directory
│       └── <project>/
│           ├── project.json
│           ├── features/
│           ├── tasks/
│           ├── issues/
│           └── events.jsonl
└── binaries/
    └── mandor            # Built binary
```

---

## Support

- Issues: https://github.com/budisantoso/mandor/issues
- Documentation: `/docs` directory

---

**Built for AI Agent Workflows**

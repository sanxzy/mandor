# Mandor - Event-Based Task Manager CLI for AI Agent Workflows

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
mandor task create api-feature-xxx "JWT Parser" \
  --goal "Parse and validate JWT tokens in incoming requests with expiry and signature verification" \
  --implementation-steps "Setup crypto library|Add token validation|Handle expiry|Return errors" \
  --test-cases "Valid token accepted|Expired token rejected|Invalid signature rejected" \
  --derivable-files "jwt_validator.go|jwt_test.go" \
  --library-needs "golang-jwt" \
  --priority P1

# Create dependent task (depends on JWT Parser)
mandor task create api-feature-xxx "Login Endpoint" \
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
mandor task list api-feature-xxx

# See tasks ready to work on
mandor task ready api-feature-xxx

# See blocked/waiting tasks
mandor task blocked api-feature-xxx

# See summary grouped by status
mandor task summary api-feature-xxx
```

### 6. Mark Tasks Complete

```bash
# Get task ID from list
mandor task update <task-id> --status in_progress
mandor task update <task-id> --status done

# Dependent tasks auto-transition to "ready"
mandor task ready api-feature-xxx  # Now shows "Login Endpoint" as ready
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
# Create a task (positional arguments: feature_id, name)
mandor task create <feature_id> "<name>" \
  --goal "Task description (min 500 chars)" \
  --implementation-steps "step1|step2|step3" \
  --test-cases "test1|test2|test3" \
  --derivable-files "file1.go|file2.go" \
  --library-needs "lib1|lib2" \
  [--priority <P0-P5>] \
  [--depends-on <task-id>]

# List tasks in a feature (positional argument: feature_id)
mandor task list <feature_id> [--status <value>] [--priority <value>] [--json]

# Show ready tasks (available to work on)
mandor task ready <feature_id> [--priority <P0-P5>] [--json]

# Show blocked tasks (waiting on dependencies)
mandor task blocked <feature_id> [--priority <P0-P5>] [--json]

# Show task summary (grouped by status)
mandor task summary <feature_id>

# Show task details
mandor task detail <task-id> [--include-deleted] [--events] [--dependencies]

# Update task status
mandor task update <id> --status <value>

# Transition task to work on it
mandor task update <id> --status in_progress

# Mark task complete (auto-unblocks dependents)
mandor task update <id> --status done

# Block task manually (external dependency)
mandor task update <id> --status blocked --reason "Waiting on API response"

# Cancel task
mandor task update <id> --cancel --reason "Superseded by feature X"

# Reopen cancelled task
mandor task update <id> --reopen
```

### Issue Commands

```bash
# Create an issue
mandor issue create <project_id> "<name>" \
  --type <type> \
  --goal "Issue description (min 200 chars)" \
  --affected-files "file1|file2" \
  --affected-tests "test1|test2" \
  --implementation-steps "step1|step2" \
  [--priority <P0-P5>] \
  [--depends-on <issue-id>]

# List issues in a project
mandor issue list <project_id> [--status <value>] [--type <value>] [--json]

# Show ready issues (available to fix)
mandor issue ready <project_id> [--type <type>] [--priority <P0-P5>] [--json]

# Show blocked issues (waiting on dependencies)
mandor issue blocked <project_id> [--type <type>] [--priority <P0-P5>] [--json]

# Show issue summary (grouped by status)
mandor issue summary <project_id>

# Show issue details
mandor issue detail <issue-id> [--include-deleted]

# Update issue status
mandor issue update <id> --status <value>

# Start working on an issue
mandor issue update <id> --start

# Mark issue resolved
mandor issue update <id> --resolve

# Mark issue as won't fix
mandor issue update <id> --wontfix --reason "Working as intended"

# Reopen a resolved/wontfix issue
mandor issue update <id> --reopen

# Block issue manually
mandor issue update <id> --status blocked --reason "Waiting on infrastructure"

# Cancel issue (duplicate/no longer relevant)
mandor issue update <id> --cancel --reason "Duplicate of issue #123"
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
# Initialize and create project
export MANDOR_ENV=development
mandor init "API Project"
mandor project create api --name "REST API" --goal "Build production REST API with authentication"

# Create feature (structured epic)
mandor feature create "Authentication" --project api \
  --goal "Implement JWT-based auth with login, logout, and token refresh flows" \
  --scope backend

# Create first task (no dependencies)
mandor task create api-feature-xxx "JWT Parser" \
  --goal "Parse and validate JWT tokens with signature verification and expiry checks" \
  --implementation-steps "Import crypto library|Implement token parsing|Add signature validation|Handle expiry" \
  --test-cases "Valid token accepted|Invalid signature rejected|Expired token rejected" \
  --derivable-files "jwt_parser.go|jwt_parser_test.go" \
  --library-needs "golang-jwt" \
  --priority P1

# Create dependent task
mandor task create api-feature-xxx "JWT Validator" \
  --goal "Create middleware to validate JWT tokens in incoming requests" \
  --implementation-steps "Create middleware|Validate token|Check expiry|Return errors" \
  --test-cases "Valid requests pass|Invalid requests rejected|Error responses correct" \
  --derivable-files "auth_middleware.go|auth_middleware_test.go" \
  --library-needs "none" \
  --depends-on <jwt-parser-id> \
  --priority P1

# Create another dependent task
mandor task create api-feature-xxx "Login Endpoint" \
  --goal "Accept credentials and return JWT token pair with refresh token" \
  --implementation-steps "Setup endpoint|Validate credentials|Generate JWT|Store refresh token|Return tokens" \
  --test-cases "Valid creds return tokens|Invalid creds rejected|Tokens formatted correctly" \
  --derivable-files "login_handler.go|login_handler_test.go" \
  --library-needs "none" \
  --depends-on <jwt-parser-id> \
  --priority P1

# View progress and blocking status
mandor task list api-feature-xxx              # See all tasks
mandor task ready api-feature-xxx             # Only JWT Parser is ready
mandor task blocked api-feature-xxx           # Validator and Login waiting on Parser
mandor task summary api-feature-xxx           # Grouped summary

# Execute workflow
mandor task update <jwt-parser-id> --status in_progress
mandor task update <jwt-parser-id> --status done

# Now Validator and Login auto-transition to ready
mandor task ready api-feature-xxx             # Both now ready to start

# Work on dependent tasks
mandor task update <jwt-validator-id> --status in_progress
mandor task update <jwt-validator-id> --status done

mandor task update <login-endpoint-id> --status in_progress
mandor task update <login-endpoint-id> --status done

# Mark feature complete
mandor feature update api-feature-xxx --project api --status active
mandor feature update api-feature-xxx --project api --status done
```

### Multi-Project Dependencies

```bash
# Create projects with cross-project dependencies enabled
mandor project create core --name "Core Library" \
  --goal "Shared database and utility libraries"
mandor project create api --name "API Service" \
  --goal "REST API service" --task-dep cross_project_allowed

# Create shared library feature and task
mandor feature create "Database Layer" --project core \
  --goal "Connection pool and query builder for database access" \
  --scope backend

mandor task create core-feature-xxx "Database Connection Pool" \
  --goal "Implement connection pool with health checks and auto-reconnect" \
  --implementation-steps "Create pool|Setup health check|Auto-reconnect|Connection limits" \
  --test-cases "Pool creates connections|Health check works|Reconnect on failure" \
  --derivable-files "db_pool.go|db_pool_test.go" \
  --library-needs "pgx|pgxpool" \
  --priority P0

# Create API feature depending on core library task
mandor feature create "User Endpoints" --project api \
  --goal "User CRUD endpoints backed by database" \
  --scope backend

mandor task create api-feature-xxx "User API Handler" \
  --goal "Create REST endpoints for user CRUD operations" \
  --implementation-steps "Setup handler|Implement GET|Implement POST|Implement DELETE" \
  --test-cases "GET returns user|POST creates user|DELETE removes user" \
  --derivable-files "user_handler.go|user_handler_test.go" \
  --library-needs "gin|none" \
  --depends-on <db-pool-task-id> \
  --priority P0

# Check progress
mandor task list core-feature-xxx          # DB Pool ready to start
mandor task blocked api-feature-xxx        # User Handler blocked on DB Pool

# Complete core task, dependents auto-unblock
mandor task update <db-pool-task-id> --status in_progress
mandor task update <db-pool-task-id> --status done

# API task now ready
mandor task ready api-feature-xxx          # User Handler now ready
```

### Issue Tracking with Dependencies

```bash
# Create blocker issues
mandor issue create api "Data Validation Bug" \
  --type bug \
  --goal "Fix data validation allowing invalid emails to pass through" \
  --affected-files "validate.go|email_validator.go" \
  --affected-tests "validate_test.go|email_validator_test.go" \
  --implementation-steps "Add email regex|Add domain check|Add test cases" \
  --priority P0

mandor issue create api "Cache Consistency Issue" \
  --type bug \
  --goal "Fix race condition in cache invalidation causing stale data" \
  --affected-files "cache.go|invalidation.go" \
  --affected-tests "cache_test.go" \
  --implementation-steps "Add mutex|Refactor invalidation|Add concurrency tests" \
  --priority P0

# Create dependent issue (depends on blockers)
mandor issue create api "High Memory Usage in Production" \
  --type bug \
  --goal "Fix memory leak causing OOM errors in production due to cache and validation issues" \
  --affected-files "main.go|memory.go" \
  --affected-tests "memory_test.go" \
  --implementation-steps "Profile memory|Identify leaks|Fix validation and cache" \
  --depends-on <data-validation-bug-id>|<cache-issue-id> \
  --priority P0

# Check status
mandor issue ready api          # Both blockers ready
mandor issue blocked api        # Memory issue blocked

# Resolve blocking issues (order doesn't matter)
mandor issue update <data-validation-bug-id> --start
mandor issue update <data-validation-bug-id> --resolve

mandor issue update <cache-issue-id> --start
mandor issue update <cache-issue-id> --resolve

# Check progress
mandor issue ready api          # Memory issue now ready
mandor issue summary api        # See all status groups

# Fix dependent issue
mandor issue update <memory-issue-id> --start
mandor issue update <memory-issue-id> --resolve
```

### Cancel and Reopen Workflow

```bash
# Create experimental feature and tasks
mandor feature create "OAuth2 Investigation" --project api \
  --goal "Research OAuth2 integration options" \
  --scope backend

mandor task create api-feature-xxx "Compare OAuth2 Libraries" \
  --goal "Research and compare auth0, okta, and open-source options" \
  --implementation-steps "Create comparison spreadsheet|Evaluate pros/cons|Estimate effort" \
  --test-cases "Evaluation complete|Team consensus reached" \
  --derivable-files "oauth2_comparison.md" \
  --library-needs "none" \
  --priority P2

mandor task create api-feature-xxx "Proof of Concept" \
  --goal "Build minimal OAuth2 integration demo" \
  --implementation-steps "Setup oauth lib|Create login flow|Test flow" \
  --test-cases "Login works|Token refresh works" \
  --derivable-files "poc_oauth.go|poc_oauth_test.go" \
  --library-needs "none" \
  --priority P2

# List tasks
mandor task list api-feature-xxx             # Show feature tasks

# Change requirements, decide not to pursue OAuth2
mandor feature update api-feature-xxx --project api \
  --cancel --reason "Sticking with JWT, OAuth2 adds too much complexity"

# Try to create new task (fails)
mandor task create api-feature-xxx "Integration Tests" \
  --goal "Add integration tests" \
  --implementation-steps "none" \
  --test-cases "none" \
  --derivable-files "none" \
  --library-needs "none"
# Error: "Cannot create task for cancelled feature"

# Change mind - reopen the feature to continue research
mandor feature update api-feature-xxx --project api --reopen

# Can now create tasks again
mandor task list api-feature-xxx             # Previous tasks still exist
mandor task create api-feature-xxx "Hybrid Approach" \
  --goal "Combine JWT with optional OAuth2 for third-party apps" \
  --implementation-steps "Design hybrid flow|Implement dual auth|Test both flows" \
  --test-cases "JWT still works|OAuth2 works|Both interoperable" \
  --derivable-files "hybrid_auth.go" \
  --library-needs "oauth2lib" \
  --priority P3
```

---

## Best Practices

### Mandor vs. Markdown Plan Files

| Problem | Markdown Plans | Mandor |
|---------|---|---|
| Single source of truth | ❌ Scattered across multiple files | ✓ Centralized `events.jsonl` |
| Dependency tracking | ❌ Manual, often wrong | ✓ Automatic status transitions |
| Progress visibility | ❌ Requires manual updates | ✓ Real-time status queries |
| Audit trail | ❌ Git history only | ✓ Immutable event log |
| Blocking detection | ❌ Must review files | ✓ `mandor task blocked <id>` |
| Schema validation | ❌ Free-form text | ✓ Enforced structure |
| Automation | ❌ Parse text with regex | ✓ JSON queryable for scripts |

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

### 9. Stop Writing Markdown Plan Files

Replace this workflow:

```markdown
# PLAN.md
## Phase 1: Authentication
- [ ] JWT parser (depends on cryptography)
- [ ] Login endpoint (depends on JWT parser)
- ...
# Status: Last updated 3 days ago
```

With this:

```bash
# Create structured plan
mandor feature create "Authentication" --project api \
  --goal "Implement JWT and login endpoints" \
  --scope backend

# Create tasks with explicit dependencies
mandor task create "JWT Parser" --feature auth-xxx \
  -g "Validate JWT tokens..." \
  --implementation-steps "Step 1|Step 2" \
  --test-cases "Test invalid tokens|Test expired" \
  --library-needs "jsonwebtoken" \
  --priority P1

mandor task create "Login Endpoint" --feature auth-xxx \
  -g "Accept credentials and return JWT..." \
  --depends-on jwt-parser-id \
  --priority P1

# Real-time progress queries
mandor task ready auth-xxx           # See what's available now
mandor task blocked auth-xxx         # See what's waiting
mandor task summary auth-xxx         # See grouped status
```

Benefits:
- No file sync required
- Dependencies auto-validated
- Blocking tasks auto-detected
- Reproducible state (`events.jsonl`)
- Queryable via CLI or JSON
- Works in CI/CD pipelines

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

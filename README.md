# Mandor - Deterministic Task Manager for AI Agent Workflows

<p align="center">
  <img src="logo.png" alt="Mandor Logo" width="600">
</p>

<p align="center">
  <strong>Stop writing markdown plans that go stale.</strong>
</p>

<p align="center">
  <strong>Deterministic task management with structured briefs, specifications, and blueprints.</strong>
</p>

<p align="center">
  Gate-enforced | Dependency-aware | Structured storage | CLI-native
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#why-mandor">Why Mandor</a> •
  <a href="#workflow">Workflow</a> •
  <a href="#commands">Commands</a> •
  <a href="#best-practices">Best Practices</a>
</p>

---

## Why Mandor

Traditional workflows scatter task state across markdown files, spreadsheets, and Slack. Dependencies are manual, status is fiction, and context is lost between sessions.

Mandor brings **deterministic task management** to AI agent workflows:

- **Single Source of Truth**: Brief → Spec → Blueprint → Feature → Task pipeline
- **Gate-Enforced Progression**: All three gates (brief read, spec read, session notes) required before starting work
- **Automatic Dependency Resolution**: Mark task done → dependents auto-transition to ready
- **Cross-Feature Dependencies**: Task A (Feature 1) can depend on Task B (Feature 2) with cascading unblocking
- **Schema-Driven**: Enforce implementation steps, test cases, requirements upfront
- **CLI-Native**: Works in terminal, scripts, and CI/CD pipelines

---

## Alternatives

Mandor is one approach to deterministic task management. Other tools exist in this space:

- **Beads**: Git-backed issue tracker focusing on rapid iteration with minimal overhead
  - GitHub: https://github.com/steveyegge/beads

- **OpenSpec**: Fluid, iterative specs with AI tool integration (Claude Code, Cursor, Copilot)
  - GitHub: https://github.com/Fission-AI/OpenSpec

---

## Installation

### Build from Source

```bash
cd Mandor
go build -o mandor ./cmd/mandor
./mandor --help
```

### Install to PATH

```bash
go build -o ~/.local/bin/mandor ./cmd/mandor
export PATH="$HOME/.local/bin:$PATH"
mandor --help
```

---

## Quick Start

### 1. Initialize Workspace

```bash
mandor init -y
```

### 2. Create Project

```bash
mandor project create api --name "API Service" \
  --goal "REST API with JWT authentication, user management, and data endpoints"
```

### 3. Create Brief (Project Intent & Capabilities)

```bash
mandor brief create -p api \
  --name "JWT Authentication" \
  --why "Need secure, stateless authentication for API endpoints" \
  --capabilities "jwt-auth:JWT login flow|refresh:Token refresh endpoint" \
  --tech-stack "golang,jwt,redis"
```

Creates: `jwt-auth` (brief ID)

### 4. Create Spec (Requirements with IAE Scenarios)

```bash
mandor spec create -p api \
  --capability jwt-auth \
  --summary "JWT authentication specification" \
  --requirements "Setup:Enable login:Import JWT lib|Validate:Verify token:Parse JWT|Refresh:Update token:Issue new JWT"
```

Creates: `jwt-auth-spec` (spec ID with auto-generated requirement IDs: req-0001, req-0002, etc.)

### 5. Create Blueprint (Technical Architecture)

```bash
mandor blueprint create -p api \
  --brief jwt-auth \
  --problem "Secure API authentication without server-side session state" \
  --decisions "Use JWT for stateless auth|HTTP-only cookies for token storage|Refresh tokens on expiry" \
  --goals-in-scope "Authentication,Authorization,Token refresh"
```

Creates: `api-blueprint` (blueprint ID)

### 6. Create Feature (from Spec)

```bash
mandor feature create "JWT Authentication" -p api \
  --capability jwt-auth \
  --spec-id jwt-auth-spec \
  --scope backend \
  -g "Implement JWT-based authentication with login flow, token validation, and refresh mechanism"
```

Creates: `jwt-tokens` (feature ID, one-to-one mapping to spec-id)

### 7. Create Tasks (with IAE Scenarios)

```bash
# Task 1: Setup JWT Library
mandor task create jwt-tokens "Setup JWT Library" \
  --spec-id jwt-auth-spec \
  --iae-scenarios "req-0001:scenario-0001|req-0001:scenario-0002" \
  -g "Configure golang-jwt library with crypto setup and validation" \
  --implementation-steps "Import library|Configure settings|Add crypto|Add validation" \
  --test-cases "JWT creates|JWT validates|JWT signature verifies" \
  --library-needs "golang-jwt"
```

Creates: `jwt-tokens-task-a7K2` (status=ready, all gates=false)

```bash
# Task 2: Token Validation (depends on Task 1)
mandor task create jwt-tokens "Token Validation Middleware" \
  --spec-id jwt-auth-spec \
  --iae-scenarios "req-0002:scenario-0001" \
  --depends-on jwt-tokens-task-a7K2 \
  -g "Implement middleware to validate JWT tokens on all protected endpoints" \
  --implementation-steps "Create middleware|Extract token|Validate|Return errors" \
  --test-cases "Valid passes|Expired rejects|Invalid rejects"
```

Creates: `jwt-tokens-task-b3M8` (status=blocked, waiting for task-a7K2)

### 8. Track Progress

```bash
# See feature with all tasks
mandor track feature jwt-tokens
```

Output shows:
- Task 1: status=ready, gates=false
- Task 2: status=blocked (waiting for Task 1)

### 9. Set Gates & Start Work

```bash
# Read brief, spec, and session notes first
# Then set gates for Task 1:
mandor task set-gate jwt-tokens-task-a7K2 --is-read-brief
mandor task set-gate jwt-tokens-task-a7K2 --is-read-spec
mandor task set-gate jwt-tokens-task-a7K2 --is-read-session-notes

# Transition to in_progress (requires all gates=true)
mandor task update jwt-tokens-task-a7K2 --status in_progress

# Work on implementation...

# Mark complete
mandor task update jwt-tokens-task-a7K2 --status done

# Task 2 automatically transitions: blocked → ready
mandor track feature jwt-tokens  # Now shows Task 2 as ready
```

### 10. Repeat for Task 2

```bash
# Set gates for Task 2
mandor task set-gate jwt-tokens-task-b3M8 --is-read-brief
mandor task set-gate jwt-tokens-task-b3M8 --is-read-spec
mandor task set-gate jwt-tokens-task-b3M8 --is-read-session-notes

# Start work
mandor task update jwt-tokens-task-b3M8 --status in_progress
```

---

## Workflow

### Complete Pipeline

```
1. Brief (Intent & Capabilities)
   ↓
2. Spec (Requirements & IAE Scenarios)
   ↓
3. Blueprint (Architecture Decisions)
   ↓
4. Feature (Spec Mapping)
   ↓
5. Task (Work Items with IAE References)
```

### Gate Enforcement

Every task has three read gates:
- **IsReadBrief**: Have you read the Brief?
- **IsReadSpec**: Have you read the Spec?
- **IsReadSessionNotes**: Have you read session notes?

**All three gates MUST be true before ready → in_progress transition.**

Gates are NOT required for:
- ready → cancelled
- in_progress → done
- blocked → ready (auto-transition)

### Dependency Auto-Resolution

Tasks with `--depends-on`:
- Auto-assigned `status=blocked` on creation
- Stay blocked until ALL dependencies are `done`
- Auto-transition `blocked → ready` when dependencies complete
- Works cross-feature (Task A in Feature 1 depends on Task B in Feature 2)
- Cascading: Task A done → unblocks Task B → Task B done → unblocks Task C

---

## Commands

### Essential Three Commands

#### 1. mandor populate
View all available commands and usage instructions.

```bash
mandor populate                 # Full reference
mandor populate | grep "brief"  # Find relevant sections
```

#### 2. mandor track
Check status of workspace, projects, features, tasks, issues.

```bash
mandor track                        # Workspace overview
mandor track project <project-id>   # Project features
mandor track feature <feature-id>   # Feature with tasks
mandor track task <task-id>         # Task with gate status
mandor track task <task-id> --json  # Machine-readable
```

#### 3. mandor session note
Record and read session progress (for AI agents).

```bash
mandor session note "Completed Task 1 and dependencies"
mandor session note --read           # Show last 50 notes
mandor session note --read --offset 100  # Show more notes
```

### Full Command Reference

**For complete command documentation, run:**

```bash
mandor populate     # Shows all commands with examples
mandor -h          # Shows available commands
mandor <cmd> -h    # Shows specific command help
```

### Quick Command List

| Category | Commands |
|----------|----------|
| **Workspace** | `init`, `config` |
| **Project** | `create`, `detail`, `update`, `delete`, `reopen` |
| **Brief** | `create`, `read`, `update`, `delete`, `validate` |
| **Spec** | `create`, `detail`, `update`, `delete`, `validate` |
| **Blueprint** | `create`, `detail`, `update`, `delete`, `validate` |
| **Feature** | `create`, `detail`, `update`, `delete` |
| **Task** | `create`, `detail`, `set-gate`, `read-gates`, `update` |
| **Issue** | `create`, `detail`, `update` |
| **Track** | `track` (workspace\|project\|feature\|task\|issue) |
| **Session** | `session note` |

---

## Best Practices

### 1. Follow the Complete Workflow

Brief → Spec → Blueprint → Feature → Task

Don't skip steps. Each phase builds context for the next.

### 2. Write Comprehensive Briefs

Include:
- Clear problem statement and motivation
- Well-defined capabilities with descriptions
- Technical stack that will be used
- Affected systems and dependencies

```bash
mandor brief create -p api \
  --name "Authentication" \
  --why "Secure stateless API authentication" \
  --capabilities "login:User login|tokens:Token management" \
  --tech-stack "golang,jwt,postgres"
```

### 3. Detailed Specs with IAE Scenarios

Create specs with Intent-Action-Expectation scenarios:

```bash
mandor spec create -p api \
  --capability login \
  --summary "Login specification" \
  --requirements "Login:User enters creds:Check password|Token:Generate JWT:Return token|Refresh:Extend session:Issue new token"
```

Specs become reference documents for tasks and gates.

### 4. Link Tasks to Spec Requirements

Tasks reference specific requirement-scenario pairs:

```bash
mandor task create feature-id "Implement Login" \
  --spec-id spec-id \
  --iae-scenarios "req-0001:scenario-0001|req-0002:scenario-0001"
```

Traces implementation back to requirements.

### 5. Gate Enforcement Discipline

Always follow the gate workflow:

1. Read Brief thoroughly
2. Read Spec with requirements
3. Read session notes from previous work
4. Set all three gates
5. Transition to in_progress
6. Implement
7. Mark done (auto-unblocks dependents)

```bash
# Check task status
mandor track task <task-id>

# Set gates
mandor task set-gate <task-id> --is-read-brief
mandor task set-gate <task-id> --is-read-spec
mandor task set-gate <task-id> --is-read-session-notes

# Start work
mandor task update <task-id> --status in_progress
```

### 6. Dependency Management

Create dependent tasks with clear relationships:

```bash
# Task A (no dependencies)
mandor task create feature-id "Task A" \
  --spec-id spec-id \
  --iae-scenarios "req-0001:scenario-0001" \
  -g "..." \
  --implementation-steps "s1|s2" \
  --test-cases "t1|t2"

# Task B (depends on Task A)
mandor task create feature-id "Task B" \
  --spec-id spec-id \
  --iae-scenarios "req-0002:scenario-0001" \
  --depends-on <task-a-id> \
  -g "..." \
  --implementation-steps "s1|s2" \
  --test-cases "t1|t2"
```

Task B auto-blocks until Task A is done.

### 7. Configuration for Your Team

Set defaults early, rarely change:

```bash
mandor config set default_priority P2
mandor config set strict_mode true
mandor config set goal.lengths.task 500
```

### 8. Document Status Changes

Always explain why you're cancelling or changing status:

```bash
mandor task update <task-id> --status cancelled \
  --reason "Superseded by feature X"
```

### 9. Use Pipe Separators for Lists

```bash
--implementation-steps "Step 1|Step 2|Step 3"
--test-cases "Test 1|Test 2|Test 3"
--iae-scenarios "req-0001:scenario-0001|req-0002:scenario-0001"
--depends-on "task-1|task-2|task-3"
```

### 10. Track Regularly Before Starting

Always check status first:

```bash
mandor track feature <feature-id>   # See all task states
mandor track task <task-id>         # Check gate status and dependencies
```

### 11. AI Agent Session Management

Log your progress between sessions:

```bash
# End of session
mandor session note "Completed task setup and testing, next: validation middleware"

# Start of session
mandor session note --read  # See what was done
```

---

## Status Transitions

### Task Lifecycle

```
New (with depends-on)          New (no depends)
    ↓                              ↓
blocked                          ready
    ↓                              ↓
    └──────→ ready ←──────┐        │
             ↓  ↑          │       │
    in_progress │          │       │
             ↓  │          │       │
           done ├──cancelled│──────┤
                │          │       │
                └──────────┴───────┘
```

**Rules:**
- New task without `--depends-on`: status=ready, all gates=false
- New task with `--depends-on`: status=blocked
- blocked → ready: auto-transition when dependencies done
- ready → in_progress: requires all three gates=true
- ready → cancelled: allowed without gates
- in_progress → done: allowed without gates
- done: immutable, no transitions out
- Error messages show which gates are unmet

---

## File Structure

```
.mandor/
├── workspace.json          # Workspace metadata
├── config.json             # Configuration
├── session-notes.jsonl     # AI session progress (NDJSON)
└── projects/
    └── <project-id>/
        ├── project.json
        ├── briefs/
        │   └── <brief-id>.md          # Brief document
        ├── specs/
        │   └── <spec-id>.md           # Spec with requirements
        ├── blueprints.jsonl           # Blueprint records
        ├── features.jsonl             # Feature records
        ├── tasks.jsonl                # Task records
        └── issues.jsonl               # Issue records
```

---

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `default_priority` | string | P3 | Default priority for new entities (P0-P5) |
| `strict_mode` | boolean | false | Enable strict validation |
| `goal.lengths.project` | integer | 500 | Min chars for project goal |
| `goal.lengths.feature` | integer | 300 | Min chars for feature goal |
| `goal.lengths.task` | integer | 500 | Min chars for task goal |
| `goal.lengths.issue` | integer | 200 | Min chars for issue goal |

```bash
mandor config list                                  # Show all
mandor config get default_priority                  # Get one
mandor config set default_priority P2               # Set one
mandor config reset default_priority                # Reset to default
```

---

## Troubleshooting

### Gate Transition Error

**Error:** "Cannot transition to in_progress: gates not set"

**Solution:** Set all three gates before transitioning

```bash
mandor track task <task-id>  # Check which gates are false

mandor task set-gate <task-id> --is-read-brief
mandor task set-gate <task-id> --is-read-spec
mandor task set-gate <task-id> --is-read-session-notes

mandor task update <task-id> --status in_progress
```

### Task Blocked by Dependencies

**Error:** "Task is blocked by dependencies"

**Solution:** Complete all blocking tasks first

```bash
mandor track task <task-id>  # See which tasks are blocking

# Complete each blocking task
mandor task update <blocking-task-id> --status done

# Dependent task auto-transitions to ready
mandor track task <task-id>
```

### Feature Not Found

**Error:** "Feature not found" or "Entity not found"

**Solution:** Verify ID and check workspace

```bash
mandor track project <project-id>  # List all features
mandor track feature <feature-id>  # Check if exists
```

---

## Development

### Build

```bash
cd Mandor
go build -o mandor ./cmd/mandor
```

### Run Tests

```bash
go test ./...
```

### View All Commands

```bash
./mandor --help
./mandor populate  # Full reference with examples
```

---

## Key Features

✓ **Structured Workflow**: Brief → Spec → Blueprint → Feature → Task
✓ **Gate Enforcement**: Three read gates required before starting work
✓ **Dependency Tracking**: Auto-blocking and auto-unblocking with cascading
✓ **Cross-Feature Dependencies**: Tasks can depend across feature boundaries
✓ **Session Tracking**: Session notes for AI agent progress
✓ **CLI-Native**: Terminal, scripts, CI/CD pipelines
✓ **Deterministic**: Single source of truth in JSONL files
✓ **Auditable**: Full change history

---

## Support

- **Documentation**: Run `mandor populate` for complete command reference
- **Help**: `mandor <command> --help` for any command
- **Repository**: https://github.com/budisantoso/mandor

---

**Built for AI Agent Workflows**

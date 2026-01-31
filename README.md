# Mandor - Event-Based Task Manager CLI for AI Agent Workflows

<p align="center">
<img src="logo.png" alt="Mandor CLI">
</p>

<p align="center">
  <strong>Deterministic JSONL output | Streaming-native architecture | Schema-driven task management</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
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

## Background: Why Mandor Was Built

Research on **Context Rot** reveals a critical challenge for AI agents: LLM performance degrades significantly as input token count increases.

### The Problem

AI agents working on long tasks accumulate conversation history, task notes, and context. Research shows:

| Factor | Impact |
|--------|--------|
| Input Length | Performance drops 10-40% as tokens increase |
| Irrelevant Content | Causes 15-30% error rate |
| Task Complexity | Reasoning degrades faster than retrieval |

Even simple retrieval tasks show degradation at scale. Benchmarks like "Needle in a Haystack" (NIAH) show near-perfect scores, but they test simple keyword matching - not real-world reasoning.

### Why Structured Task Management Helps

Instead of stuffing everything into the context window:

```bash
# Instead of: "Remember the 15 tasks from our conversation..."

# Use Mandor to externalize state:
mandor task list --project api --status pending
# Returns compact JSON for parsing

mandor task detail auth-feature-abc-task-xyz789
# Exact state, no ambiguity
```

Mandor provides:
- **Compact Context**: Replace verbose descriptions with structured JSON
- **Deterministic Output**: JSONL format is reliable to parse
- **Complete Audit Trail**: Event log shows what changed and when
- **Dependency Enforcement**: Auto-blocking prevents invalid states

---

## Installation

Mandor can be installed via curl or npm.

### Option 1: curl (Recommended for macOS/Linux)

The fastest way to install Mandor on macOS or Linux:

```bash
# Install latest stable version (to ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh -s -- --prefix /usr/local/bin

# Install latest prerelease (beta versions)
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh -s -- --prerelease

# Install specific version
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh -s -- --version v0.0.16
```

**Default install location:** `$HOME/.local/bin/mandor`

**Add to PATH:**
```bash
# For bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc

# For zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

**Verify installation:**
```bash
mandor --help
```

### Option 2: NPM (Cross-platform)

Install via npm for macOS, Linux, or Windows:

```bash
# Install globally
npm install -g @mandor/cli

# Verify installation
mandor --help
```

**Or use npx to run without installing:**
```bash
npx @mandor/cli init "My Project"
```

**Programmatic usage:**
```javascript
const mandor = require('@mandor/cli');

const cli = new mandor.Mandor({ json: true, cwd: '/project/path' });
await cli.init('My Project');
await cli.projectCreate('api', { name: 'API Service' });
const tasks = await cli.taskList({ project: 'api', status: 'pending' });
```

### Platform Support

| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| curl | ✅ arm64, x64 | ✅ arm64, x64 | ❌ |
| npm | ✅ arm64, x64 | ✅ arm64, x64 | ✅ arm64, x64 |

### Troubleshooting

**curl: command not found**
Install curl: `brew install curl` (macOS) or `sudo apt install curl` (Linux)

**mandor: command not found**
Ensure the install directory is in your PATH (see above)

**Permission denied**
```bash
# Fix permissions for ~/.local/bin
chmod +x ~/.local/bin/mandor
```

**NPM permission errors**
```bash
# Use npx (no install)
npx @mandor/cli init "My Project"

# Or fix npm global permissions
mkdir -p ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

---

## Quick Start

```bash
# 1. Initialize workspace
mandor init "My Project"

# 2. Create project
mandor project create api --name "API Service" --goal "Implement REST API"

# 3. Create feature
mandor feature create "User Auth" --project api --goal "Implement login/logout"

# 4. Get feature ID and create task
FEATURE_ID=$(mandor feature list --project api --json | jq -r '.[0].id')
mandor task create "Password Hashing" \
  --feature $FEATURE_ID \
  --goal "Implement bcrypt hashing" \
  --implementation-steps "Install bcrypt|Create utility|Write tests" \
  --test-cases "Hash validation|Password comparison" \
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
| `mandor config get/set/list` | Manage configuration |

### Project

| Command | Description |
|---------|-------------|
| `mandor project create <id> --name --goal` | Create project |
| `mandor project list` | List projects |
| `mandor project detail <id>` | Show project details |
| `mandor project update <id>` | Update metadata |
| `mandor project delete <id>` | Delete project |

### Feature

| Command | Description |
|---------|-------------|
| `mandor feature create <name> --project --goal` | Create feature |
| `mandor feature list [--project <id>]` | List features |
| `mandor feature detail <id>` | Show feature details |
| `mandor feature update <id>` | Update/cancel/reopen |

**Status flow:** `draft` → `active` → `done` (or `blocked` → `cancelled`)

### Task

| Command | Description |
|---------|-------------|
| `mandor task create <name> --feature --goal --implementation-steps --test-cases --derivable-files --library-needs` | Create task |
| `mandor task list [--feature <id>] [--project <id>] [--status <status>]` | List tasks |
| `mandor task detail <id>` | Show task details |
| `mandor task update <id>` | Update task |
| `mandor task ready [--project <id>] [--priority <P0-P5>]` | List ready tasks |
| `mandor task blocked [--project <id>]` | List blocked tasks |

**Status flow:** `pending` → `ready` → `in_progress` → `done` (or `blocked` → `cancelled`)

**Note on `--library-needs`:** This flag is required. Provide comma-separated library names (e.g., `"bcrypt,lodash"`), or use `"none"` if the task requires no new external libraries.

### Issue

| Command | Description |
|---------|-------------|
| `mandor issue create <name> --project --type --goal --affected-files --affected-tests --implementation-steps` | Create issue |
| `mandor issue list [--project <id>] [--type <type>] [--status <status>]` | List issues |
| `mandor issue detail <id>` | Show issue details |
| `mandor issue update <id>` | Update/resolve/wontfix/cancel |
| `mandor issue ready [--project <id>]` | List ready issues |
| `mandor issue blocked [--project <id>]` | List blocked issues |

**Issue types:** `bug`, `improvement`, `debt`, `security`, `performance`
**Status flow:** `open` → `ready` → `in_progress` → `resolved` (or `wontfix`/`blocked` → `cancelled`)

### Utility

| Command | Description |
|---------|-------------|
| `mandor populate [--markdown\|--json]` | Full CLI reference |
| `mandor completion [bash\|zsh\|fish]` | Shell completion |

### AI Documentation

| Command | Description |
|---------|-------------|
| `mandor ai claude` | Generate CLAUDE.md for the project |
| `mandor ai agents` | Generate AGENTS.md for multi-agent coordination |

Generate AI assistant documentation files:

```bash
# Generate CLAUDE.md for Claude Code
mandor ai claude

# Generate AGENTS.md for multi-agent coordination
mandor ai agents
```

---

## Entity Types

| Entity | File | Description |
|--------|------|-------------|
| Workspace | `.mandor/workspace.json` | Root container |
| Project | `.mandor/projects/<id>/project.jsonl` | Feature/task/issue grouping |
| Feature | `.mandor/projects/<id>/features.jsonl` | High-level functionality |
| Task | `.mandor/projects/<id>/tasks.jsonl` | Work item implementing feature |
| Issue | `.mandor/projects/<id>/issues.jsonl` | Bug/improvement/debt |
| Events | `.mandor/projects/<id>/events.jsonl` | Append-only audit trail |

### ID Format

| Entity | Format | Example |
|--------|--------|---------|
| Project | `<id>` | `api` |
| Feature | `<project>-feature-<nanoid>` | `api-feature-abc123` |
| Task | `<feature_id>-task-<nanoid>` | `api-feature-abc-task-xyz789` |
| Issue | `<project>-issue-<nanoid>` | `api-issue-abc123` |

---

## File Structure

```
.mandor/
├── workspace.json          # Workspace metadata
└── projects/
    └── <project_id>/
        ├── project.jsonl      # Project metadata
        ├── schema.json        # Project rules
        ├── features.jsonl     # Feature state
        ├── tasks.jsonl        # Task state
        ├── issues.jsonl       # Issue state
        └── events.jsonl       # Append-only audit trail
```

---

## Dependency Management

### Status Based on Dependencies

- **Feature**: No deps → `draft`, all done → `active`, otherwise blocked
- **Task**: No deps → `ready`, all done → `ready`, otherwise pending
- **Issue**: No deps → `ready`, all resolved → `ready`, otherwise open

### Blocking

Cannot cancel entities that other entities depend on. Use `--force` to override.

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

### Scope Options (Features)

`frontend`, `backend`, `fullstack`, `cli`, `desktop`, `mobile`

---

## Examples

### Complete Workflow

```bash
mandor init "My Project"
mandor project create api --name "API Service" --goal "Implement API"

# Create features
mandor feature create "User Auth" --project api --goal "Login/logout/registration"

# Get feature ID for creating tasks
AUTH_FEATURE_ID=$(mandor feature list --project api --json | jq -r '.[] | select(.name == "User Auth") | .id')

# Create tasks under the User Auth feature
mandor task create "Password Hashing" \
  --feature $AUTH_FEATURE_ID \
  --goal "Implement bcrypt hashing" \
  --implementation-steps "Install bcrypt|Create utility|Write tests" \
  --test-cases "Hash validation|Password comparison" \
  --derivable-files "src/utils/password.ts" \
  --library-needs "bcrypt"

mandor task create "JWT Token Management" \
  --feature $AUTH_FEATURE_ID \
  --goal "Implement JWT token generation and validation" \
  --implementation-steps "Review JWT spec|Implement token generation|Add validation middleware" \
  --test-cases "Token generation works|Token validation works|Expired tokens rejected" \
  --derivable-files "src/utils/jwt.ts|src/middleware/auth.ts" \
  --library-needs "jsonwebtoken"

# Task with no new external libraries
mandor task create "Update Login Endpoint" \
  --feature $AUTH_FEATURE_ID \
  --goal "Refactor existing login endpoint" \
  --implementation-steps "Review current endpoint|Refactor logic|Update tests" \
  --test-cases "Endpoint returns correct status|Authentication works|Errors handled" \
  --derivable-files "src/handlers/auth.ts" \
  --library-needs "none"

mandor issue create "Fix security vulnerability" --project api \
  --type security --goal "Fix JWT signing vulnerability" \
  --affected-files "src/utils/jwt.ts" \
  --affected-tests "src/utils/jwt.test.ts" \
  --implementation-steps "Review JWT library|Update to secure version|Verify signature"

mandor status
```

### Issue Lifecycle

```bash
mandor issue create "Security Fix" --project api \
  --type security --goal "Fix vulnerability"

mandor issue update api-issue-xxx --status in_progress
mandor issue update api-issue-xxx --resolve  # or --wontfix
mandor issue update api-issue-xxx --reopen   # if needed
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | System error (I/O, internal) |
| 2 | Validation error (not found, invalid input) |
| 3 | Permission error |

---

## Support

- Issues: https://github.com/budisantoso/mandor/issues
- Documentation: `/docs` directory

---

**Built for AI Agent Workflows**

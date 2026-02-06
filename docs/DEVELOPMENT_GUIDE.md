# Mandor Development Guide

This guide covers setting up the development environment, running tests, building, and understanding the codebase architecture.

## Development Environment

### Prerequisites

- **Go 1.21+** - Download from [golang.org](https://golang.org/dl/)
- **Node.js 16+** - For npm package management and building
- **Git** - For version control

### Setup

```bash
# Clone the repository
git clone https://github.com/sanxzy/mandor.git
cd mandor

# Download dependencies
go mod download

# Verify installation
go version
```

## Running Tests

### All Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v
```

### Unit Tests

```bash
# Unit tests only
go test ./tests/unit/... -v

# Service layer tests
go test ./tests/unit/service/... -v

# Command layer tests
go test ./tests/unit/cmd/... -v

# Domain validation tests
go test ./tests/unit/domain/... -v

# File I/O tests
go test ./tests/unit/fs/... -v
```

### Test Coverage

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Building

### Current Platform

```bash
# Build binary to build/mandor
go build -o build/mandor ./cmd/mandor

# Verify build
./build/mandor --version
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o build/mandor-linux-amd64 ./cmd/mandor

# macOS
GOOS=darwin GOARCH=amd64 go build -o build/mandor-darwin-amd64 ./cmd/mandor
GOOS=darwin GOARCH=arm64 go build -o build/mandor-darwin-arm64 ./cmd/mandor

# Windows
GOOS=windows GOARCH=amd64 go build -o build/mandor-windows-amd64.exe ./cmd/mandor
```

### Release Build with NPM

```bash
# From root directory
npm run build

# Binaries built to binaries/ directory with platform-specific subdirs:
# - binaries/darwin-amd64/mandor (Intel Mac)
# - binaries/darwin-arm64/mandor (Apple Silicon)
# - binaries/linux-amd64/mandor (Linux x86_64)
# - binaries/linux-arm64/mandor (Linux ARM64)
# - binaries/windows-amd64/mandor.exe (Windows x86_64)
# - binaries/windows-arm64/mandor.exe (Windows ARM64)

# Compressed archives also created:
# - binaries/darwin-amd64.tar.gz
# - binaries/linux-amd64.tar.gz
# etc.
```

## NPM Package Build Commands

The NPM package (`@mandors/cli`) wraps the Go binary for cross-platform distribution.

```bash
# Build supported platforms (attempts all 6, only succeeds on compatible systems)
npm run build

# Build specific platforms (from package.json scripts)
npm run build:darwin:x64    # Darwin x86_64 (Intel Macs) - GOOS=darwin GOARCH=amd64
npm run build:darwin:arm64  # Darwin ARM64 (Apple Silicon) - GOOS=darwin GOARCH=arm64
npm run build:linux:x64     # Linux x86_64 - GOOS=linux GOARCH=amd64
npm run build:linux:arm64   # Linux ARM64 - GOOS=linux GOARCH=arm64
npm run build:win32:x64     # Windows x86_64 - GOOS=windows GOARCH=amd64
npm run build:win32:arm64   # Windows ARM64 - GOOS=windows GOARCH=arm64
```

Binaries are output to `binaries/` directory with platform-specific subdirectories. Compressed tar.gz archives are automatically created for distribution and GitHub releases.

### Package Structure

```
npm/
├── bin/
│   └── mandor              # CLI wrapper script
├── lib/
│   ├── index.js            # Package entry point
│   ├── api.js              # Programmatic Node.js API
│   ├── config.js           # Configuration management
│   ├── download.js         # Binary download logic
│   ├── install.js          # Post-install hook
│   └── resolve.js          # Version resolution
└── scripts/
    └── build.js            # Cross-platform build script
```

## Code Style

### Pre-commit Hooks

The project uses pre-commit hooks for automated code quality checks.

```bash
# Install pre-commit tool
brew install pre-commit  # macOS
pip install pre-commit   # or via pip

# Install hooks in this repo
cd Mandor
pre-commit install

# Run on all files (auto before commit)
pre-commit run --all-files

# Run on staged files only
pre-commit run
```

#### Configured Hooks

| Hook | Description | Excluded Paths |
|------|-------------|----------------|
| `go-fmt` | Formats Go code | None |
| `go-vet` | Static analysis | `tests/` |
| `go-mod-tidy` | Tidies Go modules | None |
| `go-build` | Builds Go packages | `tests/` |
| `go-unit-tests` | Runs unit tests | None |
| `eslint` | Lints JavaScript | `npm/lib/` |

#### Troubleshooting

**Hook fails with "no Go files" error**
- This happens when hooks run on test directories
- `go-vet` and `go-build` exclude `tests/` directory
- If error persists, run `pre-commit clean` then `pre-commit install`

**Hooks not running on commit**
- Verify hooks are installed: `pre-commit hooks`
- Check hook configuration in `.pre-commit-config.yaml`
- Run manually: `pre-commit run --all-files`

### Formatting

```bash
# Format all Go code
go fmt ./...

# Show what would be formatted
go fmt -n ./...
```

### Linting

```bash
# Install golangci-lint (if not installed)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run linter
golangci-lint run ./...
```

### Vet

```bash
# Static analysis
go vet ./...
```

## Project Structure

```
mandor/
├── cmd/
│   └── mandor/
│       └── main.go                    # CLI entry point
├── internal/
│   ├── ai/                            # AI integration (Claude agents)
│   │   └── project.go
│   ├── cmd/                           # Command handlers (Cobra)
│   │   ├── root.go                    # Root command
│   │   ├── completion.go              # Shell completion
│   │   ├── populate.go                # Populate reference data
│   │   ├── track.go                   # Track workspace/projects/features/tasks
│   │   ├── version.go                 # Version command
│   │   ├── brief/                     # Brief commands
│   │   │   └── create.go
│   │   ├── blueprint/                 # Blueprint commands
│   │   │   └── create.go
│   │   ├── feature/                   # Feature commands
│   │   │   ├── create.go
│   │   │   └── update.go
│   │   ├── issue/                     # Issue commands
│   │   │   ├── create.go
│   │   │   ├── ready.go
│   │   │   └── blocked.go
│   │   ├── project/                   # Project commands
│   │   │   └── create.go
│   │   ├── session/                   # Session commands
│   │   │   └── note.go
│   │   ├── spec/                      # Spec commands
│   │   │   └── create.go
│   │   ├── task/                      # Task commands
│   │   │   ├── create.go
│   │   │   ├── update.go
│   │   │   └── set_gate.go
│   │   └── workspace/                 # Workspace commands
│   │       └── init.go
│   ├── domain/                        # Domain entities
│   │   ├── brief.go
│   │   ├── blueprint.go
│   │   ├── feature.go
│   │   ├── issue.go
│   │   ├── project.go
│   │   ├── spec.go
│   │   ├── task.go
│   │   ├── workspace.go
│   │   └── validation.go
│   ├── fs/                            # File I/O (NDJSON reading/writing)
│   │   ├── io.go
│   │   ├── paths.go
│   │   └── errors.go
│   ├── service/                       # Business logic layer
│   │   ├── brief_service.go
│   │   ├── blueprint_service.go
│   │   ├── feature_service.go
│   │   ├── issue_service.go
│   │   ├── project_service.go
│   │   ├── spec_service.go
│   │   ├── task_service.go
│   │   ├── workspace_service.go
│   │   └── dependency_resolver.go
│   └── util/                          # Utilities
│       ├── idgen.go                   # ID generation (nanoid)
│       ├── priority.go                # Priority parsing/comparison
│       └── errors.go
├── tests/
│   ├── unit/                          # Unit tests
│   │   ├── cmd/                       # Command tests
│   │   ├── domain/                    # Domain entity tests
│   │   ├── fs/                        # File I/O tests
│   │   └── service/                   # Service layer tests
│   └── cli/                           # CLI integration tests
├── npm/                               # NPM package wrapper
│   ├── bin/                           # CLI executable wrapper
│   ├── lib/                           # Node.js package code
│   └── scripts/                       # Build scripts
├── docs/                              # Documentation
│   ├── DEVELOPMENT_GUIDE.md           # This file
│   ├── RELEASE.md                     # Release process
│   └── background/                    # Additional documentation
└── scripts/                           # Utility scripts
```

### Key Components

**Command Layer (internal/cmd/)**
- Routes CLI arguments to service layer via Cobra
- Handles input validation and flag parsing
- Formats and outputs results to stdout

**Service Layer (internal/service/)**
- Implements business logic and validation rules
- Manages task dependencies and status transitions
- Enforces gate requirements (brief read, spec read, session notes)
- Handles cross-project dependencies

**Domain Layer (internal/domain/)**
- Defines core entities: Workspace, Project, Brief, Spec, Blueprint, Feature, Task, Issue
- Specifies validation rules and constraints
- Entity lifecycle: creation, updates, status transitions

**File I/O Layer (internal/fs/)**
- Reads/writes NDJSON files for persistent storage
- Streaming parser for large workspaces (no memory bloat)
- Path management for .mandor/ directory structure

**Utility Packages (internal/util/)**
- Nanoid ID generation for entity IDs
- Priority parsing (P0-P5)
- Custom error types for validation and system errors

## Core Architecture

### Data Model

Entities are stored in `.mandor/projects/<project-id>/` using NDJSON format:

```
.mandor/
├── workspace.json            # Workspace metadata
├── config.json               # Configuration
├── session-notes.jsonl       # AI session progress (NDJSON)
└── projects/
    └── <project-id>/
        ├── project.json      # Project metadata
        ├── briefs/
        │   └── <brief-id>.md # Brief document (markdown)
        ├── specs/
        │   └── <spec-id>.md  # Spec document (markdown)
        ├── blueprints.jsonl  # Blueprint records (NDJSON)
        ├── features.jsonl    # Feature records (NDJSON)
        ├── tasks.jsonl       # Task records (NDJSON)
        └── issues.jsonl      # Issue records (NDJSON)
```

### Workflow Pipeline

```
Brief (Intent & Capabilities)
  ↓
Spec (Requirements with IAE Scenarios)
  ↓
Blueprint (Architecture Decisions)
  ↓
Feature (Maps to Spec)
  ↓
Task (Implementation Work)
```

### Task Status Lifecycle

- **ready**: Awaiting work (all gates must be set before transition)
- **blocked**: Waiting for dependencies to complete (auto-transition to ready when dependencies done)
- **in_progress**: Active work
- **done**: Completed (immutable, no further transitions)
- **cancelled**: Abandoned (soft-delete with full data preservation)

### Dependency Resolution

- Tasks can depend on other tasks (same or different projects)
- Multiple dependencies use pipe-separated format: `--depends-on "id1|id2|id3"`
- ALL dependencies must be done before dependent task can progress
- Auto-transition from blocked → ready when ALL dependencies complete
- Cascading transitions: Task A done → unblocks Task B → auto-transitions B to ready

### Gate Enforcement

Every task has three read gates that must all be true before ready → in_progress:
- **IsReadBrief**: Have you read the Brief?
- **IsReadSpec**: Have you read the Spec?
- **IsReadSessionNotes**: Have you read session notes?

## Adding New Features

### 1. Add Domain Entity

Create `internal/domain/entity.go`:

```go
type Entity struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"` // ready, in_progress, done
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validation
func (e *Entity) Validate() error {
	if e.Name == "" {
		return NewValidationError("Name required")
	}
	return nil
}
```

### 2. Add Service Methods

Create `internal/service/entity_service.go`:

```go
type EntityService struct {
	reader *fs.Reader
	writer *fs.Writer
}

func (s *EntityService) CreateEntity(ctx context.Context, input EntityInput) (*Entity, error) {
	// Validation
	// Business logic
	// File I/O
	return entity, nil
}
```

### 3. Add File I/O Methods

Extend `internal/fs/io.go`:

```go
func (r *Reader) ReadEntity(projectID, entityID string) (*Entity, error) {
	var entity *Entity
	err := r.ReadNDJSON(r.paths.ProjectEntitiesPath(projectID), func(raw []byte) error {
		var e Entity
		if err := json.Unmarshal(raw, &e); err != nil {
			return err
		}
		if e.ID == entityID {
			entity = &e
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, NewValidationError("Entity not found: " + entityID)
	}
	return entity, nil
}

func (w *Writer) WriteEntity(projectID string, entity *Entity) error {
	return w.AppendNDJSON(w.paths.ProjectEntitiesPath(projectID), entity)
}
```

### 4. Add Command Handler

Create `internal/cmd/entity/create.go`:

```go
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create entity",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := service.NewEntityService(workspace)
		entity, err := svc.CreateEntity(ctx, input)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", entity.ID)
		return nil
	},
}
```

### 5. Add Tests

Create `tests/unit/service/entity_service_test.go`:

```go
func TestEntityCreate(t *testing.T) {
	tmpDir := t.TempDir()
	svc := setupTestService(t, tmpDir)
	
	entity, err := svc.CreateEntity(ctx, validInput)
	if err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}
	if entity.ID == "" {
		t.Error("Entity ID not generated")
	}
}
```

## Testing Strategy

### Unit Tests

Location: `tests/unit/`

Run all unit tests:
```bash
go test ./tests/unit/... -v
```

Run specific test layer:
```bash
go test ./tests/unit/service/... -v    # Service layer tests
go test ./tests/unit/cmd/... -v        # Command handler tests
go test ./tests/unit/domain/... -v     # Domain validation tests
go test ./tests/unit/fs/... -v         # File I/O tests
```

### CLI Integration Tests

Location: `tests/cli/`

Tests full end-to-end workflows with actual CLI commands.

### Test Coverage

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Helpers

```go
// tests/unit/helpers_test.go - Common setup for tests
func setupTestWorkspace(t *testing.T) (string, *WorkspaceService) {
	tmpDir := t.TempDir()
	
	// Initialize workspace
	ws := NewWorkspaceService(tmpDir)
	err := ws.Init("Test Workspace")
	if err != nil {
		t.Fatalf("Failed to initialize workspace: %v", err)
	}
	
	return tmpDir, ws
}
```

## Debugging

### Enable Verbose Output

```bash
# Build with debug flags
go build -gcflags="all=-N -l" -o build/mandor-debug ./cmd/mandor

# Use delve for debugging
dlv debug ./cmd/mandor -- args...
```

### Logging

The CLI writes errors to stderr. For debugging:

```bash
# Capture full output
mandor command 2>&1 | tee debug.log
```

## Performance Considerations

### Large Workspaces

- **Batch operations**: Use track command with filters
- **NDJSON reading**: Uses streaming (no loading entire file into memory)
- **Dependency resolution**: O(n) depth-first search for cycle detection
- **Status transitions**: O(n) to find affected tasks when dependency completes

### Memory Usage

- Streaming NDJSON parser (no large memory allocation)
- Efficient line-by-line processing for large files
- Temp file writes for atomic updates prevent memory bloat

## Contributing

### Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/my-feature`
3. **Implement** changes with tests
4. **Run** `go fmt ./...` before committing
5. **Ensure** all tests pass: `go test ./...`
6. **Submit** pull request

### Commit Messages

```
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
Examples:
  - feat(task): add priority validation
  - fix(service): resolve circular dependency bug
  - docs(readme): update installation instructions
```

### Code Review Guidelines

- All code must have tests
- Follow existing patterns and conventions
- Document complex logic (especially dependency resolution)
- Keep PRs focused and small
- Ensure gate enforcement is preserved in task transitions

## Release Process

### Version Bumping

1. Update version in `package.json`
2. Update `CHANGELOG.md` with new features/fixes
3. Build binaries: `npm run build`
4. Create GitHub release: `gh release create vX.Y.Z -t "Title" -F CHANGELOG.md`
5. Upload binaries: `gh release upload vX.Y.Z binaries/*.tar.gz`
6. Commit changes: `git add package.json CHANGELOG.md && git commit -m "vX.Y.Z: description"`
7. Push to origin: `git push`
8. Publish to npm: `npm publish --access public`

### Release Checklist

```bash
# 1. Update version and changelog
# (edit package.json and CHANGELOG.md)

# 2. Build all platforms
npm run build

# 3. Create git tag and GitHub release
gh release create vX.Y.Z -t "vX.Y.Z - Description" -F CHANGELOG.md

# 4. Upload binaries (darwin-arm64.tar.gz, linux-arm64.tar.gz, etc.)
gh release upload vX.Y.Z binaries/*.tar.gz

# 5. Commit and push
git add package.json CHANGELOG.md
git commit -m "vX.Y.Z: Production release with [description]"
git push

# 6. Publish to npm
npm publish --access public

# 7. Verify
npm info @mandors/cli | grep -A5 dist-tags
```

## Troubleshooting

### Common Issues

**"Workspace not initialized"**
```bash
# Initialize workspace first
mandor init "My Project"
```

**"Project not found"**
```bash
# Check project exists
mandor track
```

**"Permission denied"**
```bash
# Check directory permissions
ls -la .mandor/
```

**"Gates not set" error on task transition**
```bash
# Check which gates are false
mandor track task <task-id>

# Set all three gates
mandor task set-gate <task-id> --is-read-brief
mandor task set-gate <task-id> --is-read-spec
mandor task set-gate <task-id> --is-read-session-notes
```

### Getting Help

- Check existing issues: https://github.com/sanxzy/mandor/issues
- Review documentation: `/docs` directory
- Run with `--help` for command options
- See [README.md](../README.md) for command reference
- Run `mandor populate` for complete usage guide

---

## Additional Resources

| Resource | Description |
|----------|-------------|
| [README.md](../README.md) | Complete user guide with workflow examples |
| [CHANGELOG.md](../CHANGELOG.md) | Version history and release notes |
| [RELEASE.md](RELEASE.md) | Detailed release process documentation |
| [AGENTS.md](../AGENTS.md) | Essential commands for AI agents |

---

**Last Updated**: February 2026

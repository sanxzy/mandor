# Mandor Development Guide

This guide covers setting up the development environment, running tests, building, and understanding the codebase architecture.

## Development Environment

### Prerequisites

- **Go 1.21+** - Download from [golang.org](https://golang.org/dl/)
- **Git** - For version control
- **jq** - For JSON processing in scripts (optional but recommended)

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

# Windows
GOOS=windows GOARCH=amd64 go build -o build/mandor-windows-amd64.exe ./cmd/mandor
```

### Release Build

```bash
# Build all platforms
./scripts/build.sh

# Or manually
GOOS=linux GOARCH=amd64 go build -o release/mandor-linux-x64 ./cmd/mandor
GOOS=darwin GOARCH=amd64 go build -o release/mandor-darwin-x64 ./cmd/mandor
GOOS=windows GOARCH=amd64 go build -o release/mandor-windows-x64.exe ./cmd/mandor
```

## NPM Package Build Commands

The NPM package (`@mandors/cli`) wraps the Go binary for cross-platform distribution.

```bash
cd npm

# Build supported platforms (attempts all 6, skips unsupported)
npm run build

# Build specific platforms
npm run build:darwin:x64    # Darwin x64 (Intel Macs)
npm run build:darwin:arm64  # Darwin arm64 (Apple Silicon)
npm run build:linux:x64     # Linux x64
npm run build:linux:arm64   # Linux arm64
npm run build:win32:x64     # Windows x64
npm run build:win32:arm64   # Windows arm64
```

Binaries are output to `npm/binaries/` directory as tar.gz archives for npm package distribution.

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

### Programmatic Usage

```javascript
const mandor = require('@mandors/cli');

const cli = new mandor.Mandor({ json: true, cwd: '/project/path' });
await cli.init('My Project');
await cli.projectCreate('api', { name: 'API Service' });
const tasks = await cli.taskList({ project: 'api', status: 'pending' });
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
│       └── main.go              # CLI entry point
├── internal/
│   ├── cmd/                      # Command handlers (Cobra)
│   │   ├── root.go               # Root command
│   │   ├── workspace/            # Workspace commands
│   │   │   ├── init.go
│   │   │   ├── status.go
│   │   │   └── config.go
│   │   ├── project/              # Project commands
│   │   │   ├── project.go
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── detail.go
│   │   │   ├── update.go
│   │   │   ├── delete.go
│   │   │   └── reopen.go
│   │   ├── feature/              # Feature commands
│   │   │   ├── feature.go
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── detail.go
│   │   │   └── update.go
│   │   ├── task/                 # Task commands
│   │   │   ├── task.go
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── detail.go
│   │   │   └── update.go
│   │   └── issue/                # Issue commands
│   │       ├── issue.go
│   │       ├── create.go
│   │       ├── list.go
│   │       ├── detail.go
│   │       └── update.go
│   ├── service/                  # Business logic layer
│   │   ├── workspace_service.go
│   │   ├── status_service.go
│   │   ├── project_service.go
│   │   ├── feature_service.go
│   │   ├── task_service.go
│   │   └── issue_service.go
│   ├── domain/                   # Data models & validation
│   │   ├── workspace.go
│   │   ├── project.go
│   │   ├── feature.go
│   │   ├── task.go
│   │   └── issue.go
│   ├── fs/                       # Filesystem I/O
│   │   ├── paths.go
│   │   └── io.go
│   └── util/                     # Utilities
│       ├── id.go
│       └── git.go
├── tests/
│   └── unit/                     # Unit tests
│       ├── cmd/
│       │   ├── workspace/
│       │   ├── project/
│       │   ├── feature/
│       │   ├── task/
│       │   └── issue/
│       └── service/
│           ├── workspace_service_test.go
│           ├── project_service_test.go
│           ├── feature_service_test.go
│           ├── task_service_test.go
│           └── issue_service_test.go
├── docs/                         # Documentation
│   ├── prd.md
│   ├── rules/
│   │   ├── dependency-rules.md
│   │   ├── status-type-reference.md
│   │   └── event-type-reference.md
│   ├── plans/
│   │   └── commands/
│   └── test/
│       ├── integration_test.md
│       └── integration_task_test.md
├── scripts/
│   └── build.sh
├── build/                        # Build output
├── IMPL
│   ├── IMPLEMENT_SUMMARY.md
│   └── AGENTS.md
├── README.md
├── DEVELOPMENT_GUIDE.md
├── go.mod
└── go.sum
```

## Architecture

### Layer Diagram

```
┌─────────────────┐
│   CLI (Cobra)   │  Command handlers - Parse flags, call services
├─────────────────┤
│    Service      │  Business logic - Validation, status transitions
├─────────────────┤
│     Domain      │  Types - Structs, constants, validation functions
├─────────────────┤
│   Filesystem    │  JSONL I/O - Read/write NDJSON files
├─────────────────┤
│     Util        │  ID generation, git integration
└─────────────────┘
```

### Key Design Decisions

1. **NDJSON Format**: Append-only events, replace for current state
   - `events.jsonl`: Append-only audit trail
   - `features.jsonl`, `tasks.jsonl`, `issues.jsonl`: Current state (replace)

2. **No Delete**: Soft delete via status, preserve audit trail
   - Cancelled features/tasks/issues remain in files
   - Filtered out by default in list commands

3. **Atomic Writes**: Write to temp file, then rename
   - Prevents data corruption on interruption
   - Uses `os.Rename` for atomic replacement

4. **Event Sourcing**: Current state computed from events
   - Every change emits an event
   - `events.jsonl` is the source of truth

5. **DFS Cycle Detection**: Linear time complexity
   - Validates no circular dependencies
   - Uses depth-first search algorithm

## Adding New Commands

### 1. Create Command File

Create `internal/cmd/<entity>/<command>.go`:

```go
package entity

import (
	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	flag1 string
	flag2 bool
)

func NewCommandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command <args>",
		Short: "Short description",
		Long:  "Long description",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewEntityService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized")
			}

			// Business logic here

			return nil
		},
	}

	cmd.Flags().StringVarP(&flag1, "flag", "f", "", "Description")
	cmd.Flags().BoolVar(&flag2, "flag2", false, "Description")

	return cmd
}
```

### 2. Register Command

Add to `internal/cmd/<entity>/<entity>.go`:

```go
func NewEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity",
		Short: "Entity commands",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDetailCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewCommandCmd())  // Add here

	return cmd
}
```

Add to `internal/cmd/root.go`:

```go
import "mandor/internal/cmd/entity"

// In NewRootCmd():
rootCmd.AddCommand(entity.NewEntityCmd())
```

### 3. Add Tests

Create `tests/unit/cmd/entity/command_test.go`:

```go
package entity_test

import (
	"testing"
	"mandor/internal/cmd/entity"
)

func TestNewCommandCmd(t *testing.T) {
	cmd := entity.NewCommandCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
}
```

## Adding New Entity Types

### 1. Define Domain Types

Create `internal/domain/<entity>.go`:

```go
package domain

import "time"

const (
	EntityStatusActive = "active"
	// ... more statuses
)

type Entity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	// ... more fields
}

// Input types for service methods
type EntityCreateInput struct {
	Name   string
	Status string
	// ... more fields
}

// Output types for display
type EntityOutput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
```

### 2. Implement Service

Create `internal/service/<entity>_service.go`:

```go
package service

type EntityService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

func NewEntityService() (*EntityService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &EntityService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

func (s *EntityService) CreateEntity(input *domain.EntityCreateInput) (*domain.Entity, error) {
	// Validation
	if input.Name == "" {
		return nil, domain.NewValidationError("Name is required")
	}

	// Business logic
	// File I/O

	return entity, nil
}

// List, Detail, Update methods...
}
```

### 3. Add Filesystem Methods

Extend `internal/fs/io.go`:

```go
func (r *Reader) ReadEntity(projectID, entityID string) (*domain.Entity, error) {
	var entity *domain.Entity
	err := r.ReadNDJSON(r.paths.ProjectEntitiesPath(projectID), func(raw []byte) error {
		var e domain.Entity
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
		return nil, domain.NewValidationError("Entity not found: " + entityID)
	}
	return entity, nil
}

func (w *Writer) WriteEntity(projectID string, entity *domain.Entity) error {
	return w.AppendNDJSON(w.paths.ProjectEntitiesPath(projectID), entity)
}
```

## Testing Strategy

### Unit Tests

Location: `tests/unit/`

```go
func TestEntityCreate(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	// ...

	// Test validation
	err := svc.ValidateCreateInput(invalidInput)
	if err == nil {
		t.Error("Expected validation error")
	}

	// Test success case
	entity, err := svc.CreateEntity(validInput)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
```

### Integration Tests

See `docs/test/integration_test.md` for comprehensive integration test scenarios.

### Test Helpers

```go
// tests/unit/service/helpers_test.go
func setupTestService(t *testing.T) (*Service, string) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create minimal workspace structure
	// ...

	return svc, tmpDir
}
```

## Debugging

### Enable Verbose Output

```bash
# Build with debug flags
go build -gcflags="all=-N -l" -o build/mandor-debug ./cmd/mandor

# Use delve or gdb
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

- **Batch operations**: Use list commands with filters
- **JSONL reading**: Uses streaming (no loading entire file)
- **Cycle detection**: O(n) DFS complexity

### Memory Usage

- Streaming JSONL parser (no large memory allocation)
- Efficient NDJSON line-by-line processing
- Temp file writes prevent memory bloat

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
- Document complex logic
- Keep PRs focused and small

## Release Process

### Version Bumping

1. Update version in code (if applicable)
2. Update CHANGELOG.md
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`

### Building Release

```bash
# Clean build directory
rm -rf build/*

# Build all platforms
./scripts/build.sh

# Sign binaries (if applicable)
# ...

# Create release notes
# ...
```

## Troubleshooting

### Common Issues

**"Workspace not initialized"**
```bash
# Initialize first
mandor init "My Project"
```

**"Project not found"**
```bash
# Check project exists
mandor project list
```

**"Permission denied"**
```bash
# Check directory permissions
ls -la .mandor/
```

### Getting Help

- Check existing issues: https://github.com/budisantoso/mandor/issues
- Review documentation: `/docs` directory
- Run with `--help` for command options

---

## Additional Resources

| Resource | Description |
|----------|-------------|
| [PRD](docs/prd.md) | Product requirements document |
| [Dependency Rules](docs/rules/dependency-rules.md) | Dependency validation rules |
| [Status Reference](docs/rules/status-type-reference.md) | Status transitions |
| [Event Reference](docs/rules/event-type-reference.md) | Event types |
| [IMPLEMENT_SUMMARY.md](IMPLEMENT_SUMMARY.md) | Implementation progress |
| [AGENTS.md](AGENTS.md) | Agent workflow guide |

---

**Happy Coding! 🚀**

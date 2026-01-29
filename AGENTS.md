# AGENTS.md - Mandor CLI Development Guide

## Project Overview
Mandor is an event-based task manager CLI for AI agent workflows. Go 1.21+, Cobra CLI framework, JSONL output format.

## Build & Test Commands
```bash
go mod download              # Install dependencies
go test ./...                # Run all tests
go test ./tests/unit/... -v  # Run all unit tests with verbose output
go test ./tests/unit/... -run TestName  # Run single test by name
go build -o binaries/mandor ./cmd/mandor  # Build CLI
go run ./cmd/mandor <command>  # Run CLI directly
go fmt ./...                 # Format code
```

## NPM Package (@mandor/cli)

### Build Commands
```bash
npm run build               # Build supported platforms (attempts all 6)
npm run build:darwin:x64    # Build Darwin x64
npm run build:darwin:arm64  # Build Darwin arm64
npm run build:linux:x64     # Build Linux x64
npm run build:linux:arm64   # Build Linux arm64
npm run build:win32:x64     # Build Windows x64
npm run build:win32:arm64   # Build Windows arm64
```

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

### NPM Installation
```bash
npm install @mandor/cli      # Install via npm
npx mandor init "My Project" # Run directly
```

### Programmatic Usage
```javascript
const mandor = require('@mandor/cli');

const cli = new mandor.Mandor({ json: true, cwd: '/project/path' });
await cli.init('My Project');
await cli.projectCreate('api', { name: 'API Service' });
const tasks = await cli.taskList({ project: 'api', status: 'pending' });
```

## Code Style Guidelines

### Go Conventions
- `gofmt` formatting required before commit
- Imports: stdlib first, then external (blank line between groups)
- `PascalCase` for types/interfaces, `camelCase` for functions/variables
- Package-level constants in `UPPER_SNAKE_CASE`
- Receiver methods: `(r *ReceiverType)`, `(s *Service)`, `(w *Writer)`

### Error Handling (CRITICAL)
- **Exit codes**: 0=success, 1=system error, 2=validation error, 3=permission error
- Always use typed `domain.MandorError` for domain errors
- Service layer raises errors, CLI layer formats output
- Never let panics escape to user
```go
return domain.NewValidationError("user-friendly message")   // exit 2
return domain.NewSystemError("operation failed", err)       // exit 1
return domain.NewPermissionError("cannot write to file")    // exit 3
```

### File I/O (CRITICAL)
- **Atomic writes**: Write to temp file, rename to target
- **NDJSON format**: One JSON object per line, events.jsonl append-only
- Use `filepath.Join()`, never string concatenation
- Test write permissions before creating directories

### Configuration & State
- Timestamps: UTC ISO8601 format (`time.Now().UTC()`)
- Status derived from events, never stored as field
- events.jsonl is append-only: never edit, only append

### Validation
- Validate all user input at service layer (not CLI)
- Workspace name: alphanumeric, hyphens, underscores only
- Priority: P0-P5 only (use `domain.ValidatePriority()`)

### Testing
- Unit tests in `*_test.go` files in same package
- Use `t.TempDir()` for filesystem tests
- Test both success and error paths

## Command Pattern
```go
func NewFooCmd() *cobra.Command {
  cmd := &cobra.Command{
    Use: "foo",
    RunE: func(cmd *cobra.Command, args []string) error {
      svc, err := service.NewFooService()
      if err != nil { return err }
      return nil
    },
  }
  cmd.Flags().StringVarP(&flag, "flag", "f", "", "Description")
  return cmd
}
```

## Directory Structure
```
cmd/mandor/main.go           # CLI entry point
internal/cmd/                # Command handlers (Cobra)
internal/service/            # Business logic layer
internal/domain/             # Data models & validation
internal/fs/                 # Filesystem I/O
internal/util/               # Utilities
tests/unit/                  # Unit tests (mirrors internal structure)
npm/                         # NPM package (cross-platform binary distribution)
```

## Important Reminders
- All PRDs in `/docs/`
- Exit codes: use `domain.ExitCode` constants
- Commands added to `root.go`'s `rootCmd.AddCommand()`
- Tests MUST pass before considering work complete
- Never include `master_docs/` in git commits
- **Pre-commit hooks**: Configured in `.pre-commit-config.yaml`. Hooks exclude `tests/` directory for non-test file checks

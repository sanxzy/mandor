# Changelog

All notable changes to this project will be documented in this file.

## [0.4.4] - 2026-02-04

### Added
- Comprehensive test coverage: 61 test scenarios across 16 test suites with 100% pass rate
- Validated dependency resolution (linear, fan-in, fan-out, diamond, deep chains)
- Validated cross-project task dependencies (3-project chains)
- Validated real-world workflows (feature releases, bug fixes, security/debt tracking)
- Validated status transitions with dependent unblocking cascades
- Validated data persistence and soft-delete recovery
- Validated error handling and circular dependency prevention

### Tested Features
- Task/feature/issue creation with full metadata (steps, test cases, files, libraries)
- Multiple dependencies using pipe-separated format: --depends-on "id1|id2|id3"
- ALL-must-be-done semantics for multiple dependencies
- Cross-project task dependencies when enabled (cross_project_allowed flag)
- Feature status transitions: draft→active→done (requires intermediate states)
- Done state immutability (cannot modify completed tasks)
- Soft-delete (cancel) with full data preservation and recovery
- Self-dependency and circular dependency detection (exit code 2)
- Valid scopes: frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift
- Valid issue types: bug, improvement, debt, security, performance
- Valid priorities: P0-P5

## [0.4.1] - 2026-02-03

### Changed
- Rewritten `mandor populate` command with comprehensive, accurate reference
- Updated help output to reflect current commands (replaced deprecated references)
- Removed event-sourcing language, replaced with "structured storage"
- Added track command documentation with all output formats
- Improved status transition diagrams
- Added quick workflows section

### Fixed
- Populate command no longer references deprecated task ready/blocked/list commands
- Populate command no longer references deprecated issue ready/blocked/list commands
- Populate command no longer references deprecated summary command
- Removed event-sourcing terminology from populate output

## [0.4.0] - 2026-02-03

### Changed
- Event removal complete: Phases 1-4 tested and verified
- Updated README: removed deprecated commands (task list/ready/blocked, issue list/ready/blocked, summary)
- Replaced event-sourcing terminology with "structured storage"
- All status transitions accurately documented from service layer
- Installation now uses GitHub curl script instead of manual build
- Entity Types table now includes Workspace and Project

### Added
- `mandor track` command: unified workspace/project/feature/task/issue tracking
- Full support for same-project and cross-project dependencies

### Removed
- Deprecated commands: task list, task ready, task blocked, issue list, issue ready, issue blocked, summary
- Event infrastructure (events.jsonl, event append methods, event structs)
- Build from source instructions (replaced with curl installation)

## [0.3.11] - 2026-02-02

### Changed
- Comprehensive README.md with all commands, workflows, and best practices
- Enhanced `mandor populate` output with complete command reference covering all 12 command groups
- Added detailed status transitions, dependency rules, and common workflows documentation

## [0.3.10] - 2026-02-02

### Changed
- Production release with stable v0.3.x feature set

## [0.3.9] - 2026-02-02

### Changed
- Enhanced README with "Why Mandor" section emphasizing markdownless task management
- Updated comprehensive tutorials with current CLI command signatures
- Improved examples showing real-world authentication and database workflows
- Added table comparing Mandor vs. Markdown plan files
- Enhanced `mandor populate` output with "No Markdown Plans" messaging at the top
- Added clear value proposition: deterministic state, automatic dependencies, real-time queries

### Added
- "Stop Writing Markdown Plans" section in `mandor populate` command output
- New "No More Markdown Plans" best practices section in populate output
- Complete workflow examples for multi-project dependencies
- Practical issue tracking examples with blocking dependencies
- Cancel/reopen workflow examples with real scenarios

## [0.3.8] - 2026-02-02

### Added
- New `mandor task summary <feature_id>` command: displays tasks grouped by status in markdown table format
- New `mandor issue summary <project_id>` command: displays issues grouped by status in markdown table format
- Both summary commands show item counts per status group and sort by creation date within groups

### Changed
- Task commands refactored to use positional arguments for feature_id: `mandor task list <feature_id>`, `mandor task ready <feature_id>`, `mandor task blocked <feature_id>`
- Task create command signature changed to: `mandor task create <feature_id> <name> [flags]`
- All task commands now follow consistent positional argument pattern

## [0.3.5] - 2026-02-02

### Changed
- Task goal validation: removed maximum character limit, now only enforces minimum 500 character requirement
- Updated task create command flag description from "max 500 chars" to "min 500 chars"
- Issue goal validation enforces minimum character requirement only (200 chars)

### Added
- Proper npm package.json for @mandors/cli distribution

## [0.0.25] - Previous Release

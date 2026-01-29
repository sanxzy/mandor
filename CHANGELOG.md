# Changelog

All notable changes to Mandor will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial CLI commands: `init`, `task`, `project`, `feature`, `issue`, `populate`, `status`, `config`, `completion`
- Event-based task management with NDJSON storage
- Workspace and project management
- Shell completion support (bash, zsh, fish)
- Configuration management with JSON and environment variable overrides

### Changed

- Improved error handling with typed domain errors (exit codes 1-3)
- Refactored service layer for better separation of concerns
- Updated file I/O to use atomic writes

## [0.0.14] - 2026-01-29

### Added

- Bundled cross-platform binaries in npm package
- NPM package `@mandor/cli` for easy installation
- Post-install hook for automatic binary extraction
- Binary resolution system for cross-platform support
- Cross-platform build script supporting darwin/linux/win32 x64/arm64

### Changed

- Package structure reorganized for npm distribution
- Updated npm scripts for build and publish workflow
- Refined .npmignore to include binaries directory

### Fixed

- Tarball extraction path alignment between install.js and resolve.js
- CLI wrapper script path resolution
- Binary caching for offline installation

## [0.0.13] - 2026-01-29 [YANKED]

### Changed

- npm package configuration updates
- Build script improvements

### Fixed

- Binary path resolution issues
- Wrapper script path corrections

## [0.0.12] - 2026-01-29 [YANKED]

### Changed

- npm package structure
- Build process refinements

## [0.0.11] - 2026-01-29 [YANKED]

### Fixed

- Tarball filename format (os-arch naming)
- Binary extraction process

## [0.0.10] - 2026-01-29 [YANKED]

### Added

- Manual tarball extraction for bundled binaries

## [0.0.9] - 2026-01-29 [YANKED]

### Fixed

- Direct file path handling for binary resolution

## [0.0.8] - 2026-01-29 [YANKED]

### Added

- Support for both directory and direct file binary structures

## [0.0.7] - 2026-01-29 [YANKED]

### Added

- Debug output for troubleshooting npm installation

## [0.0.6] - 2026-01-29 [YANKED]

### Added

- Enhanced debug logging for binary resolution

## [0.0.5] - 2026-01-29 [YANKED]

### Fixed

- Binary path format consistency (version-platform-arch/mandor)

## [0.0.4] - 2026-01-29 [YANKED]

### Fixed

- install.js extraction logic for bundled binaries

## [0.0.3] - 2026-01-29 [YANKED]

### Changed

- Updated .npmignore to include binaries in package

## [0.0.2] - 2026-01-29 [YANKED]

### Added

- Initial bundled binary support in install.js

## [0.0.1] - 2026-01-29 [YANKED]

### Added

- Initial npm package structure
- Basic install.js postinstall hook

## [0.0.0] - 2026-01-27

### Added

- Initial project setup
- Go project structure with Cobra CLI framework
- Internal package organization (cmd, service, domain, fs, util)
- Basic README and documentation
- Git repository initialization

<!--
Template for future releases:

## [X.X.X] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes to existing features

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements
-->

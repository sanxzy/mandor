# Release Process

This document outlines the complete release process for Mandor, including building cross-platform binaries, creating GitHub releases, and publishing to npm.

## Prerequisites

### Required Tools

```bash
# Go (for building binaries)
go version  # Must be 1.21+

# Node.js & npm (for npm package)
node --version  # Must be 16+
npm --version

# Git
git --version
```

### NPM Access

Ensure you have publish access to `@mandor/cli` on npmjs.com:

```bash
npm whoami  # Should show your username
npm access ls-collaborators @mandor/cli  # Verify publish access
```

### GitHub Access

Ensure you have push access to the repository and permissions to create releases.

## Version Management

Mandor uses semantic versioning: `MAJOR.MINOR.PATCH`

| Version Type | Description | Example |
|-------------|-------------|---------|
| PATCH | Bug fixes | 0.0.1 → 0.0.2 |
| MINOR | New features (backward compatible) | 0.0.2 → 0.1.0 |
| MAJOR | Breaking changes | 0.1.0 → 1.0.0 |
| PRERELEASE | Beta releases | 0.1.0 → 0.1.1-beta.0 |

### Version Bump Commands

```bash
# Patch release (bug fixes)
npm run version:patch

# Minor release (new features)
npm run version:minor

# Major release (breaking changes)
npm run version:major

# Prerelease (beta)
npm run version:beta
```

Each command will:
1. Bump version in `package.json`
2. Create a git commit
3. Create a git tag
4. Push to origin with tags

## Building Binaries

### Local Build (Current Platform Only)

```bash
npm run build
```

This builds for your current platform only:
- **macOS ARM64**: `binaries/darwin-arm64/`
- **macOS x64**: Requires macOS x64 machine or CI
- **Linux ARM64**: `binaries/linux-arm64/`
- **Linux x64**: Requires Linux machine or CI
- **Windows**: Requires Windows machine or CI

### Cross-Platform CI Build

For complete cross-platform binaries, use GitHub Actions:

1. Push your changes to `feature/npm` branch
2. GitHub Actions will build all platforms automatically
3. Download artifacts from the Actions tab

Alternatively, build manually on each platform:

```bash
# macOS ARM64 (native)
npm run build:darwin:arm64

# macOS x64 (if on macOS x64 machine)
npm run build:darwin:x64

# Linux (via Docker or Linux machine)
docker run --rm -v $PWD:/app -w /app golang:1.21 sh -c "npm run build:linux:x64 && npm run build:linux:arm64"

# Windows (via PowerShell or Windows machine)
npm run build:win32:x64
npm run build:win32:arm64
```

### Verify Builds

Check that all binaries were created:

```bash
ls -la npm/binaries/
# Expected output:
# darwin-arm64.tar.gz
# darwin-arm64/mandor
# linux-arm64.tar.gz
# linux-arm64/mandor
# ... (other platforms if built)
```

## GitHub Release Process

### Step 1: Prepare Release Notes

Draft release notes in `CHANGELOG.md` or GitHub's release notes editor.

Include:
- Summary of changes since last release
- New features
- Bug fixes
- Known issues
- Upgrade instructions (if breaking changes)

### Step 2: Create GitHub Release

**Option A: Via GitHub Web UI**

1. Navigate to: https://github.com/sanxzy/mandor/releases
2. Click "Draft a new release"
3. Select the tag version (created by `npm version` command)
4. Set target branch: `feature/npm`
5. Fill in release title and description
6. Upload binary artifacts:
   - `npm/binaries/darwin-arm64.tar.gz`
   - `npm/binaries/linux-arm64.tar.gz`
   - (other platforms if built)
7. Check "Set as the latest release" (unless prerelease)
8. Click "Publish release"

**Option B: Via GitHub CLI**

```bash
# Create release
gh release create v0.0.14 \
  --title "Mandor v0.0.14" \
  --notes "Release notes here" \
  npm/binaries/darwin-arm64.tar.gz \
  npm/binaries/linux-arm64.tar.gz

# Or from file
gh release create v0.0.14 \
  --notes-file RELEASE_NOTES.md \
  npm/binaries/*.tar.gz
```

### Step 3: Verify Release

After publishing, verify:
1. Release page is accessible
2. Downloads work correctly
3. Binary checksums match local files

## NPM Package Publishing

### Step 1: Build Binaries

```bash
npm run build
```

This runs `preversion` and `prepublishOnly` hooks which build binaries.

### Step 2: Publish to Beta (Testing)

```bash
npm run publish:beta
```

This publishes to npm with the `beta` tag.

**Test the beta release:**

```bash
cd /tmp
npm install @mandors/cli@beta
./node_modules/.bin/mandor --help
```

### Step 3: Publish to Latest (Production)

After testing the beta release:

```bash
npm run publish:latest
```

Or manually:

```bash
npm publish --access public
```

### Verify NPM Package

```bash
# Check package info
npm view @mandors/cli

# Install and test
npm install @mandor/cli
npx mandor init test-project
```

## Complete Release Checklist

### Before Release

- [ ] All tests pass: `go test ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] Build succeeds: `npm run build`
- [ ] No uncommitted changes: `git status`
- [ ] Release notes drafted

### During Release

- [ ] Version bumped: `npm run version:patch|minor|major|beta`
- [ ] Git tag pushed: Verify with `git tag -l`
- [ ] GitHub release created with binary assets
- [ ] NPM beta published: `npm run publish:beta`
- [ ] Beta tested: `npm install @mandor/cli@beta`

### After Release

- [ ] NPM latest published: `npm run publish:latest`
- [ ] Installation verified: `npm install @mandor/cli`
- [ ] CLI functional: `npx mandor --help`
- [ ] Version shown correctly: `npx mandor config get version`
- [ ] Release announced (if applicable)

## Troubleshooting

### Build Fails on Some Platforms

Cross-platform builds require different environments:

```bash
# Linux binaries from Linux/Docker
docker run --rm -v $PWD:/app -w /app golang:1.21 sh -c "npm run build:linux:x64 && npm run build:linux:arm64"

# macOS binaries require macOS
# Windows binaries require Windows
```

### NPM Publish Fails with 403

- Verify npm access: `npm access ls-collaborators @mandor/cli`
- Check package.json name matches registry
- Ensure 2FA is disabled or configured correctly

### GitHub Release Upload Fails

- Files must be under 2GB each
- Use GitHub CLI or web UI for uploads
- Check file permissions

### Binary Fails to Execute After Install

```bash
# Check cache location
ls -la ~/.mandor/bin/

# Verify binary
file ~/.mandor/bin/latest-darwin-arm64/mandor

# Should output: Mach-O 64-bit executable
```

## Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `NPM_TOKEN` | NPM registry authentication | Required for CI publishing |
| `GITHUB_TOKEN` | GitHub API authentication | Required for GitHub Actions |

## CI/CD Setup

For automated releases, create GitHub Actions workflows:

### Build Workflow (on push to feature/npm)

```yaml
name: Build
on:
  push:
    branches: [feature/npm]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: macos-latest
            platform: darwin
            arch: x64
          - os: macos-latest
            platform: darwin
            arch: arm64
          - os: ubuntu-latest
            platform: linux
            arch: x64
          - os: ubuntu-latest
            platform: linux
            arch: arm64
          - os: windows-latest
            platform: win32
            arch: x64
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Build
        run: npm run build:${{ matrix.platform }}:${{ matrix.arch }}
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: ${{ matrix.platform }}-${{ matrix.arch }}
          path: binaries/${{ matrix.platform }}-${{ matrix.arch }}.tar.gz
```

### Release Workflow (on tag push)

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Download all artifacts
        uses: actions/download-artifact@v4
        with:
          path: binaries
      - name: Create Release
        run: |
          gh release create ${{ github.ref_name }} \
            --title "Mandor ${{ github.ref_name }}" \
            binaries/**/*.tar.gz
```

## Rollback Procedure

If a release has critical issues:

### NPM Rollback

```bash
# Deprecate the bad version
npm deprecate @mandor/cli@0.0.14 "Bug in v0.0.14, use v0.0.15"

# Unpublish (within 72 hours, only if necessary)
npm unpublish @mandor/cli@0.0.14
```

### GitHub Rollback

1. Edit the release description to mark as broken
2. Create a new release with the fix
3. Update any links pointing to the broken release

## Security Considerations

- Never commit NPM tokens to git
- Use GitHub Secrets for CI/CD
- Verify checksums before publishing
- Review permissions before granting access

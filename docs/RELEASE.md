# Release Process

Complete release workflow for Mandor v0.4.2+, from version bump through npm publishing.

## Prerequisites

- Go 1.21+: `go version`
- Node.js 16+: `node --version`
- GitHub CLI: `gh --version` (must be authenticated)
- npm access to `@mandors/cli` package
- Git push access to `sanxzy/mandor` repository

## Release Workflow

### 1. Update Version

Edit `Mandor/package.json` and increment version:

```json
{
  "name": "@mandors/cli",
  "version": "0.4.2"
}
```

### 2. Build Binaries for Production

From `Mandor/` folder:

```bash
npm run build
```

This builds cross-platform binaries:
- `binaries/darwin-arm64.tar.gz` (macOS ARM64)
- `binaries/linux-arm64.tar.gz` (Linux ARM64)
- Other platform binaries (if your system supports them)

**Output structure:**
```
binaries/
├── darwin-arm64/
│   └── mandor           # Binary
├── darwin-arm64.tar.gz  # Archive
├── linux-arm64/
│   └── mandor           # Binary
└── linux-arm64.tar.gz   # Archive
```

### 3. Create GitHub Release

From `Mandor/` folder:

```bash
gh release create v0.4.2 \
  --title "v0.4.2" \
  --notes "Release notes from CHANGELOG.md" \
  --repo sanxzy/mandor
```

**Notes format from CHANGELOG.md:**
```
Stable release with comprehensive documentation

## [0.4.2] - 2026-02-03

### Changed
- Updated feature X
- Fixed bug Y

### Fixed
- Issue resolved
```

### 4. Upload Binaries to Release

```bash
gh release upload v0.4.2 \
  binaries/darwin-arm64.tar.gz \
  binaries/linux-arm64.tar.gz \
  --repo sanxzy/mandor
```

### 5. Commit Version Change

```bash
git add package.json
git commit -m "v0.4.2: Release description here"
```

Pre-commit hooks will validate formatting and tests.

### 6. Publish to NPM

```bash
npm publish --access public
```

Publishes `@mandors/cli@0.4.2` to npm registry.

**Output:**
```
npm notice Publishing to https://registry.npmjs.org/ with tag latest
+ @mandors/cli@0.4.2
```

### 7. Push to Remote

```bash
git push origin main
```

Verify: `git status` shows "up to date with origin/main"

## Complete Release Checklist

**Before Release**
- [ ] Code changes committed
- [ ] CHANGELOG.md updated with release notes
- [ ] No uncommitted changes: `git status`
- [ ] Binaries build successfully: `npm run build`

**During Release**
- [ ] Version updated in `package.json`
- [ ] `npm run build` completes without errors
- [ ] GitHub release created: `gh release create v0.4.X`
- [ ] Binaries uploaded: `gh release upload v0.4.X binaries/*.tar.gz`
- [ ] Changes committed: `git commit -m "v0.4.X: ..."`
- [ ] NPM published: `npm publish --access public`
- [ ] Pushed to main: `git push origin main`

**After Release**
- [ ] Verify on npm: `npm view @mandors/cli`
- [ ] Test install: `npm install @mandors/cli@0.4.2`
- [ ] GitHub release accessible: https://github.com/sanxzy/mandor/releases/tag/v0.4.2

## Verifying Release

### Check GitHub Release

```bash
gh release view v0.4.2 --repo sanxzy/mandor
```

### Check NPM Package

```bash
npm view @mandors/cli@0.4.2
npm view @mandors/cli  # Latest version
```

### Test Installation

```bash
cd /tmp
npm install @mandors/cli@0.4.2
npx mandor populate
```

## Troubleshooting

### Build Fails

```bash
# Clean and rebuild
rm -rf binaries/
npm run build
```

### GitHub Release Upload Fails

- Files must exist: `ls -la binaries/*.tar.gz`
- GitHub CLI authenticated: `gh auth status`
- Repository access: `gh repo view sanxzy/mandor`

### NPM Publish Fails

```bash
# Verify authentication
npm whoami

# Check package name matches registry
npm view @mandors/cli

# If 403 error: verify publish access
npm access ls-collaborators @mandors/cli
```

### Binary Doesn't Run After Install

```bash
# Extract and test binary
tar -xzf binaries/darwin-arm64.tar.gz
./binaries/darwin-arm64/mandor --version
```

## Version Strategy

Mandor uses semantic versioning:

- **PATCH** (0.4.1 → 0.4.2): Bug fixes, minor improvements
- **MINOR** (0.4.0 → 0.5.0): New features, backward compatible
- **MAJOR** (0.x.x → 1.0.0): Breaking changes, major refactors

## Release History

Releases are tracked in:
- `Mandor/CHANGELOG.md` - Development changelog
- `Mandor/package.json` - NPM package version
- GitHub Releases - https://github.com/sanxzy/mandor/releases
- NPM Registry - https://npmjs.com/package/@mandors/cli

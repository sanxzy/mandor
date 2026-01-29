# Issue: Install Script Binary Extraction Path Mismatch

## Summary

The `scripts/install.sh` script fails to properly install the Mandor binary because the tarball structure in GitHub releases does not match what the install script expects.

## Symptoms

When running the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh
```

The installation appears to succeed, but the binary is not found:

```
Extracting...
chmod: /Users/budisantoso/.local/bin/mandor: No such file or directory
```

The binary is actually extracted to a subdirectory:

```bash
$ ls -la ~/.local/bin/
drwxr-xr-x darwin-arm64/
```

## Root Cause

### Install Script Expectation

The `scripts/install.sh` (line 94-95) expects the tarball to contain the `mandor` binary at the root level:

```bash
tar -xzf "$TARFILE" -C "$INSTALL_DIR"
chmod 755 "${INSTALL_DIR}/mandor"
```

### Build Script Behavior

The `npm/scripts/build.js` creates tarballs using:

```bash
tar -czf "${archivePath}" -C "${sourceDir}" ${binaryName}
```

Where `sourceDir` is `binaries/darwin-arm64/` and `binaryName` is `mandor`.

This produces a tarball containing:
```
mandor         # Binary at root level ✓
```

### Historical Issue

Previously, the build script used:

```bash
tar -czf "${archivePath}" -C "${sourceDir}" .
```

This produced a tarball containing:
```
darwin-arm64/
└── mandor     # Binary inside subdirectory ✗
```

### The Problem

1. When the GitHub release was created, the build script was producing tarballs with the subdirectory structure
2. The install script was updated to expect binaries at the root level
3. But the old tarballs in the release still contained the subdirectory structure
4. Even after fixing the build script, the old assets in the release need to be deleted and re-uploaded

## Resolution Steps

### For Users (Workaround)

1. Delete the partially extracted files:
   ```bash
   rm -rf ~/.local/bin/darwin-arm64
   ```

2. Manually extract the binary:
   ```bash
   cd ~/.local/bin
   curl -fsSL https://github.com/sanxzy/mandor/releases/download/v0.1.7/darwin-arm64.tar.gz | tar -xzf -
   chmod 755 mandor
   ```

3. Or use npm instead:
   ```bash
   npm install -g @mandor/cli
   ```

### For Maintainers (Fix)

1. Delete the old assets from the GitHub release:
   ```bash
   gh release delete-asset v0.1.7 darwin-arm64.tar.gz --repo sanxzy/mandor
   gh release delete-asset v0.1.7 linux-arm64.tar.gz --repo sanxzy/mandor
   ```

2. Upload the corrected tarballs:
   ```bash
   gh release upload v0.1.7 binaries/darwin-arm64.tar.gz binaries/linux-arm64.tar.gz --repo sanxzy/mandor
   ```

3. Verify the tarball structure:
   ```bash
   curl -sL "https://github.com/sanxzy/mandor/releases/download/v0.1.7/darwin-arm64.tar.gz" | tar -tzf -
   # Should output: mandor (not darwin-arm64/mandor)
   ```

## Prevention

### Build Script Fix

The build script was fixed to create tarballs with the correct structure:

```javascript
// npm/scripts/build.js - createArchive function
execSync(`tar -czf "${archivePath}" -C "${sourceDir}" ${binaryName}`, {
  stdio: 'pipe',
  shell: true
});
```

### Release Checklist

When creating a new release:

1. Update version in `package.json`
2. Build with version ldflags:
   ```bash
   go build -ldflags "-X mandor/internal/cmd.version=X.X.X" -o binaries/mandor ./cmd/mandor
   ```
3. Run npm build:
   ```bash
   npm run build
   ```
4. Verify tarball structure:
   ```bash
   tar -tzf binaries/darwin-arm64.tar.gz
   # Should output: mandor (not darwin-arm64/mandor)
   ```
5. Delete old assets from GitHub release (if overwriting)
6. Upload new tarballs
7. Test installation

## Files Involved

| File | Role |
|------|------|
| `scripts/install.sh` | Extracts binary from tarball |
| `npm/scripts/build.js` | Creates distribution tarballs |
| `internal/cmd/version.go` | Version command implementation |
| `package.json` | Version source of truth |

## Related Commands

```bash
# Build with version
go build -ldflags "-X mandor/internal/cmd.version=$(node -p "require('./package.json').version")" -o binaries/mandor ./cmd/mandor

# Build npm packages
npm run build

# Create release
gh release create vX.X.X --title "vX.X.X" --notes CHANGELOG.md --repo sanxzy/mandor

# Upload assets
gh release upload vX.X.X binaries/darwin-arm64.tar.gz binaries/linux-arm64.tar.gz --repo sanxzy/mandor

# Publish npm
npm publish --access public
```

## Timeline

- **2026-01-29**: Issue identified during v0.1.6/v0.1.7 release process
- **Cause**: Build script and install script expectations were misaligned
- **Fix**: Updated build script to create tarballs with binary at root level
- **Resolution**: Deleted and re-uploaded release assets with correct structure

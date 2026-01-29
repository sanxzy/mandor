/**
 * @fileoverview Cross-platform build script for Mandor binaries
 * @description Compiles Go binaries for all supported platforms and creates distribution archives
 * @version 0.0.1
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

/**
 * @typedef {Object} PlatformConfig
 * @property {string} os - Operating system (darwin, linux, win32)
 * @property {string} arch - Architecture (x64, arm64)
 */

/** @type {PlatformConfig[]} All supported platform configurations */
const PLATFORMS = [
  { os: 'darwin', arch: 'x64' },
  { os: 'darwin', arch: 'arm64' },
  { os: 'linux', arch: 'x64' },
  { os: 'linux', arch: 'arm64' },
  { os: 'win32', arch: 'x64' },
  { os: 'win32', arch: 'arm64' }
];

/**
 * Builds the Mandor binary for a specific platform
 * @param {PlatformConfig} platform - Platform configuration
 * @param {string} sourceDir - Source directory containing Go code
 * @returns {string|null} Path to the compiled binary, or null if unsupported
 * @example
 * const binaryPath = buildForPlatform({ os: 'darwin', arch: 'arm64' }, './cmd/mandor');
 */
function buildForPlatform(platform, sourceDir) {
  const { os, arch } = platform;
  const outputDir = path.join(__dirname, '..', 'binaries', `${os}-${arch}`);
  const outputPath = path.join(outputDir, os === 'win32' ? 'mandor.exe' : 'mandor');

  fs.mkdirSync(outputDir, { recursive: true });

  console.log(`Building for ${os}-${arch}...`);

  try {
    execSync(`GOOS=${os} GOARCH=${arch} go build -o "${outputPath}" ${sourceDir}`, {
      stdio: 'pipe',
      shell: process.platform === 'win32'
    });
    return outputPath;
  } catch (error) {
    // Check stderr for unsupported pair message
    const stderr = error.stderr ? error.stderr.toString() : '';
    if (stderr.includes('unsupported GOOS/GOARCH pair')) {
      console.log(`  ⊘ Not supported on this system`);
      return null;
    }
    console.error(`  ✗ Build failed:`, error.message);
    throw error;
  }
}

/**
 * Creates a distribution archive for a platform
 * @param {PlatformConfig} platform - Platform configuration
 * @returns {string|null} Path to the archive file, or null if binary doesn't exist
 * @example
 * const archivePath = createArchive({ os: 'linux', arch: 'x64' });
 */
function createArchive(platform) {
  const { os, arch } = platform;
  const sourceDir = path.join(__dirname, '..', 'binaries', `${os}-${arch}`);

  // Check if binary exists
  const binaryExists = fs.existsSync(path.join(sourceDir, os === 'win32' ? 'mandor.exe' : 'mandor'));
  if (!binaryExists) {
    return null;
  }

  const archivePath = path.join(__dirname, '..', 'binaries', `${os}-${arch}.tar.gz`);

  console.log(`Creating archive for ${os}-${arch}...`);

  execSync(`tar -czf "${archivePath}" -C "${sourceDir}" .`, {
    stdio: 'inherit',
    shell: true
  });

  return archivePath;
}

/**
 * Main build function - attempts to build all platforms
 * @returns {Object[]} Build results for each successfully built platform
 * @example
 * const results = mainBuild();
 * console.log(`Built ${results.length} platforms`);
 */
function mainBuild() {
  const sourceDir = path.join(__dirname, '..', '..', 'cmd', 'mandor');
  const results = [];
  const unsupported = [];

  console.log(`Running on: ${os.platform()}/${os.arch()}`);
  console.log('Building cross-platform binaries...\n');

  // Build all platforms
  for (const platform of PLATFORMS) {
    const binaryPath = buildForPlatform(platform, sourceDir);
    if (binaryPath) {
      results.push({
        platform: platform.os,
        arch: platform.arch,
        binaryPath,
        status: 'built'
      });
      console.log(`  ✓ Built successfully\n`);
    } else {
      unsupported.push({
        platform: platform.os,
        arch: platform.arch,
        status: 'unsupported'
      });
    }
  }

  // Create archives for successfully built platforms
  const archiveResults = [];
  for (const platform of PLATFORMS) {
    const archivePath = createArchive(platform);
    if (archivePath) {
      const stats = fs.statSync(archivePath);
      archiveResults.push({
        platform: platform.os,
        arch: platform.arch,
        archivePath,
        archiveSize: stats.size
      });

      const existing = results.find(r => r.platform === platform.os && r.arch === platform.arch);
      if (existing) {
        existing.archivePath = archivePath;
        existing.archiveSize = stats.size;
      }
    }
  }

  // Summary
  console.log('─'.repeat(50));
  console.log(`Build complete!`);
  console.log(`  Built: ${archiveResults.length} platforms`);
  console.log(`  Unsupported: ${unsupported.length} platforms`);

  if (unsupported.length > 0) {
    console.log('\nUnsupported combinations (need different build environment):');
    unsupported.forEach(u => {
      console.log(`  - ${u.platform}/${u.arch}`);
    });
  }

  return archiveResults;
}

// Run if executed directly
if (require.main === module) {
  const results = mainBuild();
  if (results.length > 0) {
    console.log('\nBuilt platforms:');
    console.table(results.map(r => ({
      Platform: `${r.platform}/${r.arch}`,
      'Archive Size': `${(r.archiveSize / 1024).toFixed(1)} KB`
    })));

    console.log('Archives location: npm/binaries/');
  }
}

module.exports = {
  PLATFORMS,
  buildForPlatform,
  createArchive,
  mainBuild
};

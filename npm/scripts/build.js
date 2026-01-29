/**
 * @fileoverview Cross-platform build script for Mandor binaries
 * @description Compiles Go binaries for all supported platforms and creates distribution archives
 * @version 0.0.1
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

/**
 * @typedef {Object} PlatformConfig
 * @property {string} os - Operating system (darwin, linux, win32)
 * @property {string} arch - Architecture (x64, arm64)
 */

/** @type {PlatformConfig[]} Supported platform configurations */
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
 * @returns {string} Path to the compiled binary
 * @throws {Error} If build fails
 * @example
 * const binaryPath = buildForPlatform({ os: 'darwin', arch: 'x64' }, './cmd/mandor');
 */
function buildForPlatform(platform, sourceDir) {
  const { os, arch } = platform;
  const outputDir = path.join(__dirname, '..', 'binaries', `${os}-${arch}`);
  const outputPath = path.join(outputDir, os === 'win32' ? 'mandor.exe' : 'mandor');

  fs.mkdirSync(outputDir, { recursive: true });

  console.log(`Building for ${os}-${arch}...`);

  try {
    execSync(`GOOS=${os} GOARCH=${arch} go build -o "${outputPath}" ${sourceDir}`, {
      stdio: 'inherit',
      shell: process.platform === 'win32'
    });
    return outputPath;
  } catch (error) {
    console.error(`Failed to build for ${os}-${arch}:`, error.message);
    throw error;
  }
}

/**
 * Creates a distribution archive for a platform
 * @param {PlatformConfig} platform - Platform configuration
 * @returns {string} Path to the archive file
 * @example
 * const archivePath = createArchive({ os: 'darwin', arch: 'x64' });
 */
function createArchive(platform) {
  const { os, arch } = platform;
  const sourceDir = path.join(__dirname, '..', 'binaries', `${os}-${arch}`);
  const archivePath = path.join(__dirname, '..', 'binaries', `${os}-${arch}.tar.gz`);

  console.log(`Creating archive for ${os}-${arch}...`);

  execSync(`tar -czf "${archivePath}" -C "${sourceDir}" .`, {
    stdio: 'inherit',
    shell: true
  });

  return archivePath;
}

/**
 * Gets the binary filename for a platform
 * @param {PlatformConfig} platform - Platform configuration
 * @returns {string} Binary filename
 * @example
 * const filename = getBinaryName({ os: 'linux', arch: 'x64' });
 * // Returns: 'mandor-linux-x64'
 */
function getBinaryName(platform) {
  const { os, arch } = platform;
  return `mandor-${os}-${arch}`;
}

/**
 * Main build function - builds all platforms and creates archives
 * @returns {Object[]} Build results for each platform
 * @throws {Error} If any build fails
 * @example
 * const results = mainBuild();
 * console.log(`Built ${results.length} platforms`);
 */
function mainBuild() {
  const sourceDir = path.join(__dirname, '..', '..', 'cmd', 'mandor');
  const results = [];

  console.log('Starting cross-platform build...\n');

  for (const platform of PLATFORMS) {
    try {
      const binaryPath = buildForPlatform(platform, sourceDir);
      const archivePath = createArchive(platform);
      const stats = fs.statSync(archivePath);

      results.push({
        platform: platform.os,
        arch: platform.arch,
        binaryPath,
        archivePath,
        archiveSize: stats.size
      });

      console.log(`✓ ${platform.os}-${platform.arch}: ${(stats.size / 1024).toFixed(1)} KB\n`);
    } catch (error) {
      console.error(`✗ ${platform.os}-${platform.arch}: Build failed\n`);
      throw error;
    }
  }

  return results;
}

// Run if executed directly
if (require.main === module) {
  mainBuild().then(results => {
    console.log('Build complete!');
    console.table(results.map(r => ({
      Platform: `${r.platform}/${r.arch}`,
      'Archive Size': `${(r.archiveSize / 1024).toFixed(1)} KB`
    })));
  }).catch(error => {
    console.error('Build failed:', error);
    process.exit(1);
  });
}

module.exports = {
  PLATFORMS,
  buildForPlatform,
  createArchive,
  getBinaryName,
  mainBuild
};

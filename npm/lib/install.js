/**
 * @fileoverview Post-install hook for Mandor CLI
 * @description Handles binary download and caching during npm install
 * @version 0.0.1
 */

const { downloadBinary, getCurrentPlatform } = require('./download');
const fs = require('fs');
const path = require('path');
const os = require('os');

/** @type {string} Cache directory for binaries */
const CACHE_DIR = path.join(__dirname, '..', '.cache');

/** @type {string} Bundled binaries directory */
const BUNDLE_DIR = path.join(__dirname, '..', 'binaries');

/**
 * Installs the Mandor binary for the current platform
 * @async
 * @param {Object} [options] - Installation options
 * @param {string} [options.version] - Version to install (default: 'latest')
 * @returns {Promise<string>} Path to the installed binary
 * @throws {Error} If download fails
 * @example
 * // Called automatically by npm postinstall
 * await install();
 */
async function install(options = {}) {
  const version = options.version || 'latest';
  const { platform, arch } = getCurrentPlatform();

  console.log(`Installing Mandor ${version} for ${platform}-${arch}...`);

  let binaryPath;

  // Try bundled binary first
  const bundledPath = useBundledBinary(platform, arch);
  if (bundledPath) {
    console.log(`✓ Using bundled binary`);
    binaryPath = bundledPath;
  } else {
    // Download from GitHub releases
    binaryPath = await downloadBinary(version, platform, arch);
  }

  console.log(`✓ Mandor installed: ${binaryPath}`);

  return binaryPath;
}

/**
 * Uses bundled binary if available for current platform
 * @param {string} platform - Target platform
 * @param {string} arch - Target architecture
 * @returns {string|null} Path to binary or null if not bundled
 */
function useBundledBinary(platform, arch) {
  const filename = `mandor-${platform}-${arch}`;
  
  // Try both: directory structure and direct binary file
  const bundledBinaryDir = path.join(BUNDLE_DIR, filename, 'mandor');
  const bundledBinaryFile = path.join(BUNDLE_DIR, filename);

  console.log(`DEBUG: Looking for binary for ${platform}-${arch}`);
  console.log(`DEBUG: BUNDLE_DIR: ${BUNDLE_DIR}`);
  console.log(`DEBUG: Files in BUNDLE_DIR: ${fs.readdirSync(BUNDLE_DIR).join(', ')}`);

  let bundledBinary = null;
  let foundMethod = '';

  if (fs.existsSync(bundledBinaryDir)) {
    bundledBinary = bundledBinaryDir;
    foundMethod = 'directory structure';
  } else if (fs.existsSync(bundledBinaryFile) && fs.statSync(bundledBinaryFile).isFile()) {
    bundledBinary = bundledBinaryFile;
    foundMethod = 'direct file';
  }

  if (!bundledBinary) {
    console.log(`DEBUG: No bundled binary found`);
    return null;
  }

  console.log(`DEBUG: Found bundled binary via ${foundMethod}: ${bundledBinary}`);

  const cacheDir = path.join(os.homedir(), '.mandor', 'bin');
  if (!fs.existsSync(cacheDir)) {
    fs.mkdirSync(cacheDir, { recursive: true });
  }

  const dest = path.join(cacheDir, filename);

  try {
    fs.copyFileSync(bundledBinary, dest);
    fs.chmodSync(dest, '755');
    return dest;
  } catch (e) {
    console.log(`DEBUG: Failed to copy: ${e.message}`);
    return null;
  }
}

/**
 * Cleans up old binary caches
 * @returns {number} Number of files removed
 * @example
 * const removed = cleanupCache();
 * console.log(`Removed ${removed} old binary files`);
 */
function cleanupCache() {
  if (!fs.existsSync(CACHE_DIR)) return 0;

  const files = fs.readdirSync(CACHE_DIR);
  let removed = 0;

  for (const file of files) {
    const filePath = path.join(CACHE_DIR, file);
    const stats = fs.statSync(filePath);

    // Remove files older than 30 days
    const thirtyDaysAgo = Date.now() - (30 * 24 * 60 * 60 * 1000);
    if (stats.mtimeMs < thirtyDaysAgo) {
      fs.unlinkSync(filePath);
      removed++;
    }
  }

  return removed;
}

/**
 * Gets the installed binary version
 * @returns {string|null} Version string or null if not installed
 * @example
 * const version = getInstalledVersion();
 * if (version) { console.log(`Using Mandor ${version}`); }
 */
function getInstalledVersion() {
  const versionPath = path.join(CACHE_DIR, 'version.txt');
  if (fs.existsSync(versionPath)) {
    return fs.readFileSync(versionPath, 'utf-8').trim();
  }
  return null;
}

// Run install on postinstall
if (require.main === module || process.env.npm_lifecycle_event === 'postinstall') {
  install().catch(error => {
    console.error('Failed to install Mandor:', error.message);
    process.exit(1);
  });
}

module.exports = {
  install,
  cleanupCache,
  getInstalledVersion
};

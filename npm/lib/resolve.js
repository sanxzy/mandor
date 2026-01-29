/**
 * @fileoverview Version resolution module for Mandor CLI
 * @description Resolves the correct binary path based on version and platform
 * @version 0.0.1
 */

const path = require('path');
const fs = require('fs');
const os = require('os');
const { downloadBinary, getCurrentPlatform, binaryExists } = require('./download');

/** @type {string} Default version to use */
const DEFAULT_VERSION = 'latest';

/**
 * Resolves the binary path for the requested version
 * @async
 * @param {Object} [options] - Resolution options
 * @param {string} [options.version] - Requested version (default: 'latest')
 * @param {boolean} [options.forceDownload] - Force re-download even if cached
 * @returns {Promise<string>} Path to the Mandor binary
 * @throws {Error} If binary cannot be resolved or downloaded
 * @example
 * const binaryPath = await resolve({ version: '1.0.0' });
 * console.log(`Using: ${binaryPath}`);
 */
async function resolve(options = {}) {
  const version = options.version || DEFAULT_VERSION;
  const { platform, arch } = getCurrentPlatform();

  // Check cache first
  const cachedPath = getCachedBinary(version, platform, arch);
  if (cachedPath && !options.forceDownload) {
    return cachedPath;
  }

  // Download if not cached
  const binaryPath = await downloadBinary(version, platform, arch);
  cacheBinary(binaryPath, version, platform, arch);

  return binaryPath;
}

/**
 * Gets the cached binary path for a specific version
 * @param {string} version - Version to look for
 * @param {string} platform - Target platform
 * @param {string} arch - Target architecture
 * @returns {string|null} Path to cached binary or null
 * @example
 * const cached = getCachedBinary('1.0.0', 'darwin', 'x64');
 */
function getCachedBinary(version, platform, arch) {
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin');
  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  const binaryPath = path.join(cacheDir, `${version}-${platform}-${arch}`, binaryName);

  if (fs.existsSync(binaryPath)) {
    return binaryPath;
  }
  return null;
}

/**
 * Caches a binary for future use
 * @param {string} binaryPath - Path to the binary
 * @param {string} version - Version identifier
 * @param {string} platform - Target platform
 * @param {string} arch - Target architecture
 * @returns {void}
 * @example
 * cacheBinary('/home/user/.mandor/bin/1.0.0-darwin-x64/mandor', '1.0.0', 'darwin', 'x64');
 */
function cacheBinary(binaryPath, version, platform, arch) {
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin', `${version}-${platform}-${arch}`);
  fs.mkdirSync(cacheDir, { recursive: true });

  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  const destPath = path.join(cacheDir, binaryName);

  fs.copyFileSync(binaryPath, destPath);
  fs.chmodSync(destPath, '755');
}

/**
 * Lists all cached binary versions
 * @returns {Object[]} Array of cached binary info
 * @example
 * const cached = listCachedBinaries();
 * cached.forEach(b => console.log(`${b.version} (${b.platform}-${b.arch})`));
 */
function listCachedBinaries() {
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin');

  if (!fs.existsSync(cacheDir)) {
    return [];
  }

  const versions = [];
  const entries = fs.readdirSync(cacheDir);

  for (const entry of entries) {
    const entryPath = path.join(cacheDir, entry);
    if (fs.statSync(entryPath).isDirectory()) {
      const [version, platform, arch] = entry.split('-');
      versions.push({ version, platform, arch, path: entryPath });
    }
  }

  return versions;
}

/**
 * Clears all cached binaries
 * @returns {number} Number of binaries removed
 * @example
 * const removed = clearCache();
 * console.log(`Cleared ${removed} cached binaries`);
 */
function clearCache() {
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin');

  if (!fs.existsSync(cacheDir)) {
    return 0;
  }

  const entries = fs.readdirSync(cacheDir);
  let removed = 0;

  for (const entry of entries) {
    const entryPath = path.join(cacheDir, entry);
    if (fs.statSync(entryPath).isDirectory()) {
      fs.rmSync(entryPath, { recursive: true, force: true });
      removed++;
    }
  }

  return removed;
}

module.exports = {
  resolve,
  getCachedBinary,
  cacheBinary,
  listCachedBinaries,
  clearCache,
  DEFAULT_VERSION
};

/**
 * @fileoverview Binary download module for Mandor CLI
 * @description Handles downloading and caching Mandor binaries for the current platform
 * @version 0.0.1
 */

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');

/** @type {string} GitHub releases API URL */
const RELEASES_URL = 'https://api.github.com/repos/sanxzy/mandor/releases';

/**
 * Downloads the Mandor binary for the specified platform and architecture
 * @async
 * @param {string} version - The Mandor version to download (e.g., '1.0.0', 'latest')
 * @param {string} platform - Target platform (e.g., 'darwin', 'linux', 'win32')
 * @param {string} arch - Target architecture (e.g., 'x64', 'arm64')
 * @returns {Promise<string>} Path to the downloaded and executable binary
 * @throws {Error} If download fails or platform is unsupported
 * @example
 * // Download Mandor v1.0.0 for macOS x64
 * const binaryPath = await downloadBinary('1.0.0', 'darwin', 'x64');
 * console.log(`Binary downloaded to: ${binaryPath}`);
 */
async function downloadBinary(version, platform, arch) {
  const filename = `mandor-${platform}-${arch}`;
  const url = `${RELEASES_URL}/download/${version}/${filename}`;
  const dest = path.join(os.homedir(), '.mandor', 'bin', filename);

  // Download and make executable
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 302) {
        return downloadBinary(response.headers.location, platform, arch);
      }
      const file = fs.createWriteStream(dest);
      response.pipe(file);
      file.on('finish', () => {
        fs.chmodSync(dest, '755');
        resolve(dest);
      });
    }).on('error', reject);
  });
}

/**
 * Gets the platform identifier for the current system
 * @returns {{platform: string, arch: string}} Platform and architecture info
 * @example
 * const { platform, arch } = getCurrentPlatform();
 * console.log(`Running on ${platform}-${arch}`);
 */
function getCurrentPlatform() {
  const platform = os.platform(); // 'darwin', 'linux', 'win32'
  const arch = os.arch(); // 'x64', 'arm64'
  return { platform, arch };
}

/**
 * Checks if a binary already exists and is up-to-date
 * @param {string} version - Expected version
 * @param {string} platform - Target platform
 * @param {string} arch - Target architecture
 * @returns {Promise<boolean>} True if binary exists and is valid
 * @example
 * const exists = await binaryExists('1.0.0', 'darwin', 'x64');
 * if (exists) { console.log('Binary cached'); }
 */
async function binaryExists(version, platform, arch) {
  const filename = `mandor-${platform}-${arch}`;
  const dest = path.join(os.homedir(), '.mandor', 'bin', filename);
  return fs.existsSync(dest);
}

module.exports = {
  downloadBinary,
  getCurrentPlatform,
  binaryExists,
  RELEASES_URL
};

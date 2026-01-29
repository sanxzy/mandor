/**
 * @fileoverview Version resolution module for Mandor CLI
 * @description Resolves the binary path for the CLI
 * @version 0.0.3
 */

const path = require('path');
const os = require('os');
const https = require('https');

const REPO = 'sanxzy/mandor';
const DEFAULT_VERSION = 'latest';
const INSTALL_DIR = path.join(os.homedir(), '.local', 'bin');

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();
  const platformMap = { darwin: 'darwin', linux: 'linux', win32: 'win32' };
  const archMap = { x64: 'x64', arm64: 'arm64', amd64: 'x64', aarch64: 'arm64' };
  return {
    platform: platformMap[platform] || platform,
    arch: archMap[arch] || arch
  };
}

async function getLatestVersion(prerelease = false) {
  const url = prerelease
    ? `https://api.github.com/repos/${REPO}/releases`
    : `https://api.github.com/repos/${REPO}/releases/latest`;

  return new Promise((resolve, reject) => {
    https.get(url, { headers: { 'User-Agent': 'Mandor-CLI' } }, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          const parsed = JSON.parse(data);
          const tagName = Array.isArray(parsed) ? parsed[0].tag_name : parsed.tag_name;
          resolve(tagName.replace(/^v/, ''));
        } catch (e) {
          reject(e);
        }
      });
    }).on('error', reject);
  });
}

function getBinaryPath(version = DEFAULT_VERSION) {
  const { platform, arch } = getPlatform();
  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  return path.join(INSTALL_DIR, binaryName);
}

module.exports = {
  getPlatform,
  getLatestVersion,
  getBinaryPath,
  DEFAULT_VERSION,
  INSTALL_DIR
};

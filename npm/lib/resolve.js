/**
 * @fileoverview Version resolution module for Mandor CLI
 * @description Resolves the correct binary path based on version and platform
 * @version 0.0.2
 */

const path = require('path');
const fs = require('fs');
const os = require('os');
const { downloadBinary, getCurrentPlatform } = require('./download');

const REPO = 'sanxzy/mandor';
const GITHUB_API = 'https://api.github.com';
const DEFAULT_VERSION = 'latest';

function getCachedBinary(version, platform, arch) {
  const osArch = `${platform}-${arch}`;
  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  const binaryPath = path.join(os.homedir(), '.mandor', 'bin', version, osArch, binaryName);

  if (fs.existsSync(binaryPath)) {
    return binaryPath;
  }
  return null;
}

async function getLatestVersion(prerelease = false) {
  const url = prerelease
    ? `${GITHUB_API}/repos/${REPO}/releases`
    : `${GITHUB_API}/repos/${REPO}/releases/latest`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to fetch releases: ${response.statusText}`);
  }
  const data = await response.json();
  const tagName = Array.isArray(data) ? data[0].tag_name : data.tag_name;
  return tagName.replace(/^v/, '');
}

function cacheBinary(binaryPath, version, platform, arch) {
  const osArch = `${platform}-${arch}`;
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin', version, osArch);
  fs.mkdirSync(cacheDir, { recursive: true });

  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  const destPath = path.join(cacheDir, binaryName);

  fs.copyFileSync(binaryPath, destPath);
  fs.chmodSync(destPath, '755');
}

async function resolve(options = {}) {
  const version = options.version || DEFAULT_VERSION;
  const { platform, arch } = getCurrentPlatform();
  const prerelease = options.prerelease || false;

  let resolveVersion = version;
  if (version === 'latest') {
    resolveVersion = await getLatestVersion(prerelease);
  }

  const cachedPath = getCachedBinary(resolveVersion, platform, arch);
  if (cachedPath && !options.forceDownload) {
    return cachedPath;
  }

  const osArch = `${platform}-${arch}`;
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${resolveVersion}/${osArch}.tar.gz`;
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'mandor-download-'));
  const tarball = path.join(tempDir, `${osArch}.tar.gz`);

  console.log(`Downloading Mandor ${resolveVersion}...`);

  const response = await fetch(downloadUrl);
  if (!response.ok) {
    fs.rmSync(tempDir, { recursive: true });
    throw new Error(`Download failed: ${response.statusText}`);
  }

  const file = fs.createWriteStream(tarball);
  await new Promise((resolve, reject) => {
    response.body.pipe(file);
    file.on('finish', resolve);
    file.on('error', reject);
  });

  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin', resolveVersion, osArch);
  fs.mkdirSync(cacheDir, { recursive: true });

  const { execSync } = require('child_process');
  execSync(`tar -xzf "${tarball}" -C "${cacheDir}"`, { stdio: 'pipe' });
  fs.chmodSync(path.join(cacheDir, binaryName), '755');

  fs.rmSync(tempDir, { recursive: true });

  return path.join(cacheDir, binaryName);
}

function listCachedBinaries() {
  const cacheDir = path.join(os.homedir(), '.mandor', 'bin');

  if (!fs.existsSync(cacheDir)) {
    return [];
  }

  const versions = [];
  const entries = fs.readdirSync(cacheDir);

  for (const entry of entries) {
    const versionPath = path.join(cacheDir, entry);
    if (fs.statSync(versionPath).isDirectory()) {
      const subEntries = fs.readdirSync(versionPath);
      for (const subEntry of subEntries) {
        const subPath = path.join(versionPath, subEntry);
        if (fs.statSync(subPath).isDirectory()) {
          const parts = subEntry.split('-');
          const arch = parts.pop();
          const platform = parts.join('-');
          versions.push({ version: entry, platform, arch, path: subPath });
        }
      }
    }
  }

  return versions;
}

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
  DEFAULT_VERSION,
  getLatestVersion
};

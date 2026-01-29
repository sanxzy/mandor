/**
 * @fileoverview Post-install hook for Mandor CLI
 * @description Downloads binary from GitHub releases during npm install
 * @version 0.0.3
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const REPO = 'sanxzy/mandor';
const GITHUB_API = 'https://api.github.com';
const CACHE_DIR = path.join(os.homedir(), '.mandor', 'bin');

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

async function install(options = {}) {
  const { platform, arch } = getPlatform();
  const version = options.version || 'latest';
  const prerelease = options.prerelease || false;
  const osArch = `${platform}-${arch}`;
  const assetName = `${osArch}.tar.gz`;

  console.log('Mandor Installer');
  console.log('================');
  console.log(`OS: ${osArch}`);

  let installVersion = version;
  if (version === 'latest') {
    console.log(`Fetching latest ${prerelease ? 'prerelease' : 'release'}...`);
    installVersion = await getLatestVersion(prerelease);
  }

  console.log(`Version: ${installVersion}`);
  console.log('');

  const cachePath = path.join(CACHE_DIR, installVersion, osArch);
  const binaryPath = path.join(cachePath, platform === 'win32' ? 'mandor.exe' : 'mandor');

  if (fs.existsSync(binaryPath)) {
    console.log(`Using cached binary: ${binaryPath}`);
    return binaryPath;
  }

  console.log('Downloading from GitHub releases...');
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${installVersion}/${assetName}`;
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'mandor-install-'));
  const tarball = path.join(tempDir, assetName);

  const response = await fetch(downloadUrl);
  if (!response.ok) {
    fs.rmSync(tempDir, { recursive: true });
    throw new Error(`Download failed: ${response.statusText} (${downloadUrl})`);
  }

  const file = fs.createWriteStream(tarball);
  await new Promise((resolve, reject) => {
    response.body.pipe(file);
    file.on('finish', resolve);
    file.on('error', reject);
  });

  if (!fs.existsSync(cachePath)) {
    fs.mkdirSync(cachePath, { recursive: true });
  }

  execSync(`tar -xzf "${tarball}" -C "${cachePath}"`, { stdio: 'inherit' });
  fs.chmodSync(binaryPath, '755');

  fs.rmSync(tempDir, { recursive: true });

  console.log(`Installed: ${binaryPath}`);
  return binaryPath;
}

if (require.main === module || process.env.npm_lifecycle_event === 'postinstall') {
  const prerelease = process.argv.includes('--prerelease') || process.argv.includes('-p');
  install({ prerelease }).catch(error => {
    console.error('Failed to install Mandor:', error.message);
    process.exit(1);
  });
}

module.exports = { install, getLatestVersion, getPlatform };

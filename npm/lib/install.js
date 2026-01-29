/**
 * @fileoverview Post-install hook for Mandor CLI
 * @description Downloads binary from GitHub releases during npm install
 * @version 0.0.4
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');
const https = require('https');

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

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        return downloadFile(res.headers.location, dest).then(resolve).catch(reject);
      }
      res.pipe(file);
      file.on('finish', () => {
        file.close(resolve);
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
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

  await downloadFile(downloadUrl, tarball);

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

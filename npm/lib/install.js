/**
 * @fileoverview Post-install hook for Mandor CLI
 * @description Downloads and extracts binary from GitHub releases
 * @version 0.0.5
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');
const https = require('https');

const REPO = 'sanxzy/mandor';
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
  const binaryName = platform === 'win32' ? 'mandor.exe' : 'mandor';

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

  const binaryPath = path.join(INSTALL_DIR, binaryName);

  if (fs.existsSync(binaryPath)) {
    console.log(`Already installed: ${binaryPath}`);
    return binaryPath;
  }

  console.log('Downloading from GitHub releases...');
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${installVersion}/${osArch}.tar.gz`;
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'mandor-'));
  const tarball = path.join(tempDir, `${osArch}.tar.gz`);

  await downloadFile(downloadUrl, tarball);

  fs.mkdirSync(INSTALL_DIR, { recursive: true });
  execSync(`tar -xzf "${tarball}" -C "${tempDir}"`, { stdio: 'inherit' });

  const extractedBinary = path.join(tempDir, binaryName);
  fs.copyFileSync(extractedBinary, binaryPath);
  fs.chmodSync(binaryPath, '755');

  fs.rmSync(tempDir, { recursive: true });

  console.log(`Installed: ${binaryPath}`);
  console.log('');
  console.log('Add to PATH:');
  console.log(`  export PATH="${INSTALL_DIR}:$PATH"`);
  return binaryPath;
}

if (require.main === module || process.env.npm_lifecycle_event === 'postinstall') {
  install().catch(error => {
    console.error('Failed to install Mandor:', error.message);
    process.exit(1);
  });
}

module.exports = { install, getLatestVersion, getPlatform };

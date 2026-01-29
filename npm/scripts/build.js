/**
 * @fileoverview Cross-platform build script for Mandor binaries
 * @description Compiles Go binaries for all supported platforms and creates distribution archives
 * @version 0.0.2
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const ROOT_DIR = path.join(__dirname, '..', '..');
const BINARIES_DIR = path.join(ROOT_DIR, 'binaries');

const PLATFORMS = [
  { os: 'darwin', arch: 'x64' },
  { os: 'darwin', arch: 'arm64' },
  { os: 'linux', arch: 'x64' },
  { os: 'linux', arch: 'arm64' },
  { os: 'win32', arch: 'x64' },
  { os: 'win32', arch: 'arm64' }
];

function buildForPlatform(platform) {
  const { os, arch } = platform;
  const outputDir = path.join(BINARIES_DIR, `${os}-${arch}`);
  const binaryName = os === 'win32' ? 'mandor.exe' : 'mandor';
  const outputPath = path.join(outputDir, binaryName);

  fs.mkdirSync(outputDir, { recursive: true });

  console.log(`Building for ${os}-${arch}...`);

  const pkg = JSON.parse(fs.readFileSync(path.join(ROOT_DIR, 'package.json'), 'utf8'));
  const version = pkg.version;

  try {
    execSync(`GOOS=${os} GOARCH=${arch} go build -ldflags "-X mandor/internal/cmd.version=${version}" -o "${outputPath}" ./cmd/mandor`, {
      stdio: 'pipe',
      shell: process.platform === 'win32'
    });
    return outputPath;
  } catch (error) {
    const stderr = error.stderr ? error.stderr.toString() : '';
    if (stderr.includes('unsupported GOOS/GOARCH pair')) {
      console.log(`  Not supported on this system`);
      return null;
    }
    console.error(`  Build failed:`, error.message);
    throw error;
  }
}

function createArchive(platform) {
  const { os, arch } = platform;
  const sourceDir = path.join(BINARIES_DIR, `${os}-${arch}`);
  const binaryName = os === 'win32' ? 'mandor.exe' : 'mandor';

  if (!fs.existsSync(path.join(sourceDir, binaryName))) {
    return null;
  }

  const archivePath = path.join(BINARIES_DIR, `${os}-${arch}.tar.gz`);

  try {
    execSync(`tar -czf "${archivePath}" -C "${sourceDir}" ${binaryName}`, {
      stdio: 'pipe',
      shell: true
    });
    return archivePath;
  } catch (error) {
    console.error(`  Failed to create archive:`, error.message);
    return null;
  }
}

function uploadToGithubReleases(results) {
  let version;
  try {
    version = execSync('git describe --tags --abbrev=0', { encoding: 'utf8' }).trim();
  } catch {
    const pkg = require(path.join(ROOT_DIR, 'package.json'));
    version = pkg.version;
  }

  const archives = results.filter(r => r.archivePath && fs.existsSync(r.archivePath));

  if (archives.length === 0) {
    console.log('No archives to upload');
    return;
  }

  console.log(`\nUploading ${archives.length} binaries to GitHub release ${version}...`);

  for (const archive of archives) {
    const assetPath = archive.archivePath;
    const assetName = `${path.basename(path.dirname(assetPath))}.tar.gz`;

    try {
      console.log(`  Uploading ${assetName}...`);
      execSync(`gh release upload "${version}" "${assetPath}" --repo "${process.env.GITHUB_REPOSITORY || 'sanxzy/mandor'}"`, {
        stdio: 'pipe',
        shell: true
      });
      console.log(`  Uploaded ${assetName}`);
    } catch (error) {
      console.error(`  Failed to upload ${assetName}:`, error.message);
    }
  }
}

function cleanBuildDirs() {
  if (fs.existsSync(BINARIES_DIR)) {
    for (const entry of fs.readdirSync(BINARIES_DIR)) {
      if (entry.endsWith('.tar.gz')) {
        fs.unlinkSync(path.join(BINARIES_DIR, entry));
      }
    }
  }
}

function mainBuild() {
  console.log(`Running on: ${os.platform()}/${os.arch()}`);
  console.log('Building cross-platform binaries...\n');

  cleanBuildDirs();

  const results = [];
  const unsupported = [];

  for (const platform of PLATFORMS) {
    const binaryPath = buildForPlatform(platform);
    if (binaryPath) {
      results.push({ platform: platform.os, arch: platform.arch, binaryPath, status: 'built' });
      console.log(`  Built successfully\n`);
    } else {
      unsupported.push({ platform: platform.os, arch: platform.arch, status: 'unsupported' });
    }
  }

  for (const platform of PLATFORMS) {
    const archivePath = createArchive(platform);
    if (archivePath) {
      const stats = fs.statSync(archivePath);
      const existing = results.find(r => r.platform === platform.os && r.arch === platform.arch);
      if (existing) {
        existing.archivePath = archivePath;
        existing.archiveSize = stats.size;
      }
    }
  }

  console.log('─'.repeat(50));
  console.log(`Build complete!`);
  console.log(`  Built: ${results.length} platforms`);
  console.log(`  Unsupported: ${unsupported.length} platforms`);

  if (unsupported.length > 0) {
    console.log('\nUnsupported combinations (need different build environment):');
    unsupported.forEach(u => {
      console.log(`  - ${u.platform}/${u.arch}`);
    });
  }

  console.log(`\nArchives location: ${BINARIES_DIR}/`);

  uploadToGithubReleases(results);

  return results;
}

if (require.main === module) {
  const results = mainBuild();
  if (results.length > 0) {
    console.log('\nBuilt platforms:');
    console.table(results.map(r => ({
      Platform: `${r.platform}/${r.arch}`,
      'Archive Size': `${(r.archiveSize / 1024).toFixed(1)} KB`
    })));
  }
}

module.exports = { PLATFORMS, buildForPlatform, createArchive, mainBuild };

#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

// GitHub release settings
const REPO = 'hugompham/frontforge';
const VERSION = require('../package.json').version;

// Determine platform and architecture
function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  const platformMap = {
    darwin: {
      x64: 'darwin-amd64',
      arm64: 'darwin-arm64'
    },
    linux: {
      x64: 'linux-amd64'
    },
    win32: {
      x64: 'windows-amd64'
    }
  };

  if (!platformMap[platform]) {
    console.error(`Unsupported platform: ${platform}`);
    process.exit(1);
  }

  if (!platformMap[platform][arch]) {
    console.error(`Unsupported architecture: ${arch} on ${platform}`);
    process.exit(1);
  }

  return platformMap[platform][arch];
}

// Get binary name
function getBinaryName() {
  const platformStr = getPlatform();
  const ext = process.platform === 'win32' ? '.exe' : '';
  return `frontforge-${platformStr}${ext}`;
}

// Download binary from GitHub releases
function downloadBinary(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);

    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        return downloadBinary(response.headers.location, dest)
          .then(resolve)
          .catch(reject);
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode}`));
        return;
      }

      response.pipe(file);

      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

// Main installation function
async function install() {
  console.log('📦 Installing frontforge...');

  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, '..', 'bin');
  const binaryPath = path.join(binDir, 'frontforge' + (process.platform === 'win32' ? '.exe' : ''));

  // Create bin directory if it doesn't exist
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  try {
    // Try to download from GitHub releases
    const downloadUrl = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;

    console.log(`⬇️  Downloading binary from GitHub releases...`);
    console.log(`   Platform: ${getPlatform()}`);

    await downloadBinary(downloadUrl, binaryPath);

    // Make binary executable (Unix-like systems)
    if (process.platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }

    console.log('Installation complete!');
    console.log('\nUsage:');
    console.log('   npx create-frontend-app');
    console.log('   or');
    console.log('   create-frontend-app');
  } catch (error) {
    console.error('Installation failed:', error.message);
    console.log('\nAlternative installation:');
    console.log('   You can download the binary directly from:');
    console.log(`   https://github.com/${REPO}/releases/tag/v${VERSION}`);
    process.exit(1);
  }
}

// For local development: copy from ../bin/
function installLocal() {
  console.log('🔧 Local installation mode...');

  const platformStr = getPlatform();
  const ext = process.platform === 'win32' ? '.exe' : '';
  const sourceBinary = path.join(__dirname, '..', '..', 'bin', `frontforge-${platformStr}${ext}`);
  const destBinary = path.join(__dirname, '..', 'bin', `frontforge${ext}`);

  if (!fs.existsSync(sourceBinary)) {
    console.error(`❌ Local binary not found: ${sourceBinary}`);
    console.log('   Please run "make build-all" first');
    process.exit(1);
  }

  const binDir = path.join(__dirname, '..', 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  fs.copyFileSync(sourceBinary, destBinary);

  if (process.platform !== 'win32') {
    fs.chmodSync(destBinary, 0o755);
  }

  console.log('Local installation complete!');
}

// Check if we're in development mode (local installation)
const isLocal = fs.existsSync(path.join(__dirname, '..', '..', 'bin'));

if (isLocal) {
  installLocal();
} else {
  install();
}

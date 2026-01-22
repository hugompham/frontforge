#!/usr/bin/env node

const { spawn } = require('child_process');
const { existsSync } = require('fs');
const path = require('path');

// Determine binary name based on platform
const binaryName = 'frontforge' + (process.platform === 'win32' ? '.exe' : '');
const binaryPath = path.join(__dirname, binaryName);

// Check if binary exists before trying to execute
if (!existsSync(binaryPath)) {
  console.error('Error: frontforge binary not found.');
  console.error('');
  console.error('The binary should have been downloaded during installation.');
  console.error('Please try reinstalling:');
  console.error('');
  console.error('  npm uninstall -g frontforge');
  console.error('  npm install -g frontforge');
  console.error('');
  console.error('If the issue persists, you can download the binary manually from:');
  console.error('  https://github.com/hugompham/frontforge/releases');
  process.exit(1);
}

// Execute the Go binary with all arguments
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  shell: false
});

child.on('exit', (code) => {
  process.exit(code || 0);
});

child.on('error', (err) => {
  if (err.code === 'ENOENT') {
    console.error('Error: Binary not found at:', binaryPath);
    console.error('');
    console.error('Please try reinstalling frontforge:');
    console.error('  npm uninstall -g frontforge && npm install -g frontforge');
  } else if (err.code === 'EACCES') {
    console.error('Error: Permission denied when trying to execute binary.');
    console.error('');
    console.error('Try fixing permissions:');
    console.error('  chmod +x', binaryPath);
  } else {
    console.error('Error: Failed to start binary:', err.message);
    console.error('Code:', err.code);
  }
  process.exit(1);
});

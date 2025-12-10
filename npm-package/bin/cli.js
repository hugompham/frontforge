#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

// Determine binary name based on platform
const binaryName = 'frontforge' + (process.platform === 'win32' ? '.exe' : '');
const binaryPath = path.join(__dirname, binaryName);

// Execute the Go binary with all arguments
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  shell: false
});

child.on('exit', (code) => {
  process.exit(code);
});

child.on('error', (err) => {
  console.error('Failed to start binary:', err.message);
  process.exit(1);
});

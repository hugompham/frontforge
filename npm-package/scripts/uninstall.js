#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const binDir = path.join(__dirname, '..', 'bin');
const binaryPath = path.join(binDir, 'frontforge' + (process.platform === 'win32' ? '.exe' : ''));

console.log('Uninstalling frontforge...');

if (fs.existsSync(binaryPath)) {
  fs.unlinkSync(binaryPath);
  console.log('Binary removed');
}

if (fs.existsSync(binDir) && fs.readdirSync(binDir).length === 0) {
  fs.rmdirSync(binDir);
}

console.log('Uninstall complete');

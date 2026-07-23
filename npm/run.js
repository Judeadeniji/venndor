#!/usr/bin/env node

const { spawnSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const binName = os.platform() === 'win32' ? 'venndor.exe' : 'venndor';
const binPath = path.join(__dirname, '..', 'bin', binName);

if (!fs.existsSync(binPath)) {
  console.error(`Error: venndor binary not found at ${binPath}.`);
  console.error(`You might need to reinstall the package or run 'node npm/install.js'.`);
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), {
  stdio: 'inherit'
});

if (result.error) {
  console.error(`Failed to execute venndor: ${result.error.message}`);
  process.exit(1);
}

process.exit(result.status);

#!/usr/bin/env node
const { spawnSync } = require('child_process');
const os = require('os');
const path = require('path');

const platform = os.platform();
const arch = os.arch();

const mapPlatform = { win32: 'win32', darwin: 'darwin', linux: 'linux' };
const mapArch = { x64: 'x64', arm64: 'arm64' };

const osStr = mapPlatform[platform];
const archStr = mapArch[arch];

if (!osStr || !archStr) {
  console.error(`Unsupported platform/architecture: ${platform}-${arch}`);
  process.exit(1);
}

const pkgName = `@judeadeniji/venndor-${osStr}-${archStr}`;
let binPath;

try {
  const pkgDir = path.dirname(require.resolve(`${pkgName}/package.json`));
  const binName = platform === 'win32' ? 'venndor.exe' : 'venndor';
  binPath = path.join(pkgDir, binName);
} catch (e) {
  console.error(`Error: Could not find native binary package ${pkgName}.`);
  console.error('Make sure you have not disabled optional dependencies.');
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
if (result.error) {
  console.error(`Failed to execute venndor: ${result.error.message}`);
  process.exit(1);
}
process.exit(result.status);

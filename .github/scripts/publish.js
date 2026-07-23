const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const version = process.env.VERSION;
if (!version) {
  console.error('VERSION env var is required');
  process.exit(1);
}

const tag = version.includes('-') ? 'next' : 'latest';

const platforms = [
  { os: 'linux', cpu: 'x64', dirPattern: 'linux_amd64' },
  { os: 'linux', cpu: 'arm64', dirPattern: 'linux_arm64' },
  { os: 'darwin', cpu: 'x64', dirPattern: 'darwin_amd64' },
  { os: 'darwin', cpu: 'arm64', dirPattern: 'darwin_arm64' },
  { os: 'win32', cpu: 'x64', dirPattern: 'windows_amd64' }
];

const distDir = path.join(__dirname, '../../dist');
const optionalDependencies = {};

for (const p of platforms) {
  const pkgName = `venndor-${p.os}-${p.cpu}`;
  optionalDependencies[pkgName] = version;
  
  const pkgDir = path.join(__dirname, '../../npm-packages', pkgName);
  fs.mkdirSync(pkgDir, { recursive: true });

  const binName = p.os === 'win32' ? 'venndor.exe' : 'venndor';
  
  // Find the exact directory in dist/
  const distDirs = fs.readdirSync(distDir);
  const targetDir = distDirs.find(d => d.includes(p.dirPattern) && !d.includes('.tar.gz') && !d.includes('.zip'));
  
  if (!targetDir) {
    console.warn(`Warning: Could not find build directory for ${p.dirPattern}`);
    continue;
  }

  const srcBin = path.join(distDir, targetDir, binName);
  
  // Create package.json
  const pkgJson = {
    name: pkgName,
    version: version,
    description: `venndor native binary for ${p.os} ${p.cpu}`,
    os: [p.os],
    cpu: [p.cpu],
    repository: "https://github.com/judeadeniji/venndor.git",
    author: "Jude Adeniji",
    license: "MIT"
  };

  fs.writeFileSync(path.join(pkgDir, 'package.json'), JSON.stringify(pkgJson, null, 2));
  
  if (fs.existsSync(srcBin)) {
    fs.copyFileSync(srcBin, path.join(pkgDir, binName));
    if (p.os !== 'win32') {
      fs.chmodSync(path.join(pkgDir, binName), 0o755);
    }
  } else {
    console.warn(`Warning: Binary not found at ${srcBin}`);
    continue;
  }

  console.log(`Publishing ${pkgName}...`);
  execSync(`npm publish --access public --tag ${tag}`, { cwd: pkgDir, stdio: 'inherit' });
}

// Update main package.json
const mainPkgPath = path.join(__dirname, '../../package.json');
const mainPkg = JSON.parse(fs.readFileSync(mainPkgPath, 'utf8'));
mainPkg.version = version;
mainPkg.optionalDependencies = optionalDependencies;
fs.writeFileSync(mainPkgPath, JSON.stringify(mainPkg, null, 2));

console.log(`Publishing main package...`);
execSync(`npm publish --access public --tag ${tag}`, { cwd: path.join(__dirname, '../..'), stdio: 'inherit' });

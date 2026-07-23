const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { execSync } = require('child_process');
const packageJson = require('../package.json');

const version = packageJson.version;
const platform = os.platform();
const arch = os.arch();

const mapPlatform = {
  win32: 'windows',
  darwin: 'darwin',
  linux: 'linux'
};

const mapArch = {
  x64: 'x86_64',
  arm64: 'arm64'
};

const osStr = mapPlatform[platform] || platform;
const archStr = mapArch[arch] || arch;

// Following goreleaser default name_template for archives
const filename = `venndor_${osStr}_${archStr}.tar.gz`;
const url = `https://github.com/judeadeniji/venndor/releases/download/v${version}/${filename}`;

const binDir = path.join(__dirname, '..', 'bin');
const binName = platform === 'win32' ? 'venndor.exe' : 'venndor';
const binPath = path.join(binDir, binName);

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 301 || response.statusCode === 302) {
        return download(response.headers.location, dest).then(resolve).catch(reject);
      }
      if (response.statusCode !== 200) {
        return reject(new Error(`Failed to download: ${response.statusCode}`));
      }
      const file = fs.createWriteStream(dest);
      response.pipe(file);
      file.on('finish', () => {
        file.close(resolve);
      });
    }).on('error', reject);
  });
}

async function install() {
  console.log(`Downloading venndor v${version} for ${osStr}-${archStr}...`);
  const tmpTarball = path.join(__dirname, filename);
  try {
    await download(url, tmpTarball);
    console.log('Extracting...');
    // using built in tar command since we are on mostly unix systems. For Windows, tar is available in modern Windows 10+.
    execSync(`tar -xzf "${tmpTarball}" -C "${binDir}"`);
    fs.unlinkSync(tmpTarball);
    
    // Ensure binary is executable
    if (os.platform() !== 'win32') {
      fs.chmodSync(binPath, 0o755);
    }
    console.log('Installed successfully!');
  } catch (err) {
    console.error(`Failed to install venndor binary:`, err.message);
    process.exit(1);
  }
}

// In local dev scenarios or CI, don't try to download binary from GitHub since it may not be released yet
if (process.env.SKIP_BINARY_DOWNLOAD) {
  console.log("Skipping binary download");
  process.exit(0);
}

install();

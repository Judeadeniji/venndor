# venndor 📦

> Vendor npm packages directly into your repo as owned source, while retaining upstream tracking, updates, and custom patching natively.

**venndor** is a pragmatic CLI tool designed to cleanly inline Node dependencies directly into a `vendor/` directory in your repository. It acts exactly like standard workspace management but gives you total ownership over the raw source code of the dependencies you pull.

With `venndor`, you can make local hacks or emergency fixes to npm packages seamlessly via `.patch` generation and reapply them automatically when upgrading the package!

## Features

- **Intelligent Package Manager Detection**: Automatically intercepts `npm`, `yarn`, `pnpm`, or `bun` via lockfiles or `corepack`.
- **Zero Configuration**: Non-destructive integration with your `package.json` workspaces. Modifies files via surgically precise AST-like string modifications.
- **Diff & Patch Generation**: Generates standard unified diff patches and stores them automatically in `patches/`.
- **Smart Upgrades**: Upgrading vendored packages intelligently fetches new versions while flawlessly attempting to re-apply any local modifications via `patch -p2`.

## Installation

Install globally using your favorite package manager:

```bash
npm install -g venndor
```

Or via Go:
```bash
go install github.com/judeadeniji/venndor/cmd/vendor@latest
```

## Usage

**`venndor init`**
Bootstraps the project explicitly, setting up `vendor.yaml`, `vendor-lock.json`, and configuring `package.json` workspaces. (You don't *need* to run this, `add` will do it automatically!)

**`venndor add <pkg>[@version]`**
Downloads the package from the registry and extracts it precisely to `vendor/<pkg>`. It injects a clean `#vendor/<pkg>` alias into the `package.json` `imports` map.
```bash
venndor add is-even
venndor add lodash@4.17.21
```

**`venndor diff <pkg>`**
Generates a `.patch` file for any edits made to the vendored source code inside the `vendor/` directory, comparing it against the untouched original source fetched from `.vendor-cache/`.
```bash
venndor diff is-even
```

**`venndor status`**
Lists all vendored packages, showing whether they have local `[PATCHED]` modifications, and checking the registry for any available upstream versions.
```bash
venndor status --check-updates
```

**`venndor update [pkg]`**
Upgrades the specified package (or all packages if run without arguments) to the latest version, fetching the new tarball and re-applying your generated `.patch` files on top.
```bash
venndor update is-even
venndor update
```

**`venndor sync`**
Re-applies the workspace maps to `package.json` and syncs `node_modules` across the board using your natively detected package manager.

**`venndor remove <pkg>`**
Safely purges the package entirely out of the tracking manifest, deletes the folder, strips out the import configurations, removes leftover `.patch` overrides, and cleans `node_modules`.
```bash
venndor remove is-even
```

## Architecture Notes

- Raw `.tgz` npm registry tarballs are cached directly into `node_modules/.vendor-cache` to facilitate true, pristine offline diffing operations.
- Built via standard library Go + Cobra CLI for incredible performance.
- Binaries are compiled strictly via GitHub Actions + GoReleaser.

## License

MIT

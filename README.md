# venndor 📦

[![npm version](https://img.shields.io/npm/v/%40judeadeniji%2Fvenndor.svg)](https://www.npmjs.com/package/@judeadeniji/venndor)
[![Release](https://img.shields.io/github/v/release/judeadeniji/venndor?include_prereleases)](https://github.com/Judeadeniji/venndor/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> Pull an npm package's actual code into your repo, keep it up to date, and keep your own edits when you upgrade it.

`venndor` copies a dependency's real source code into a `vendor/` folder in your project, instead of leaving it hidden away in `node_modules`. You get to read it, change it, and review it like any other file in your repo, without giving up the ability to pull in upstream updates later.

It doesn't reinvent anything: `venndor` just fetches the package and lets your normal package manager (npm, yarn, pnpm, or bun) install it like usual. It only adds the wiring to make that work smoothly.

## Why

Sometimes you need to fix or tweak a dependency right now: a bug that can't wait for the maintainer, a small behavior change, or you just want that code to be visible and reviewable instead of buried in `node_modules`.

The usual ways to do this each cost you something:

- **Forking the package**: works, but now you're maintaining a whole separate copy that needs to stay in sync
- **`patch-package`**: patches the installed files automatically, but the code itself is still hidden in `node_modules`, so it never shows up in a code review
- **git submodule or subtree**: fine if you're pulling from source, but what's published on npm often doesn't match the raw git repo (built files, trimmed folders, etc.)

`venndor` downloads the exact package you'd get from `npm install`, puts it in your repo as normal files, and keeps track of any changes you make so they can be reapplied automatically the next time you upgrade.

## Install

```bash
npm install -g @judeadeniji/venndor
```

```bash
go install github.com/judeadeniji/venndor/cmd/vendor@latest
```

## Quickstart

```bash
venndor add is-even        # download it, add it to your project, done
```

```js
import isEven from "#vendor/is-even";
```

That's it. `is-even`'s code now lives in `vendor/is-even` in your repo, and that import just works, no extra config needed.

Now say you edit that code directly:

```bash
venndor diff is-even        # save your edit as a patch file
venndor update is-even      # later: get the newest version, reapply your edit on top
```

## Commands

| Command | What it does |
|---|---|
| `venndor add <pkg>[@version]` | Downloads the package and adds it to `vendor/<pkg>`, so you can import it right away |
| `venndor init` | Sets up the project for vendoring. You don't need to run this yourself, `add` does it automatically the first time |
| `venndor diff <pkg>` | Saves your local edits as a patch file, so they aren't lost |
| `venndor status [--check-updates]` | Shows which packages are vendored and which have local edits. Add `--check-updates` to also check npm for newer versions |
| `venndor update [pkg]` | Downloads the newest version (or one you specify) and reapplies your saved edits on top. If an edit can't be reapplied cleanly, it tells you instead of guessing |
| `venndor sync` | Fixes up your project setup and reinstalls, useful after cloning the repo fresh or editing `package.json` by hand |
| `venndor remove <pkg>` | Deletes the package and everything venndor added for it, then cleans up |

```bash
venndor add lodash@4.17.21
venndor update            # update everything that has a saved edit
venndor remove is-even
```

## How it works

- **It figures out your package manager on its own** (npm, yarn, pnpm, or bun), so you don't have to tell it.
- **It edits `package.json` carefully**, keeping the formatting and field order you already had, instead of rewriting the whole file.
- **A clean copy of each package is cached locally**, so your saved edits can always be compared against the original, untouched version.
- **Your edits are saved as plain patch files**, the same kind of file `git` and other tools already use. If a saved edit can't be applied to a new version (because the upstream code changed in the same spot), venndor won't silently overwrite your change. It tells you so you can fix it by hand.
- Built in Go, with prebuilt downloads for every major OS so installing it doesn't require compiling anything yourself.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT

# vendor-cli — Design Doc

A CLI tool to vendor npm packages directly into a repo as owned source, while still
tracking upstream versions, supporting local patches, and staying out of the way of
normal dependency management until explicitly opted into.

**Language/runtime:** Go
**Distribution:** npm, via prebuilt per-platform binaries (`optionalDependencies` + `os`/`cpu`, no postinstall script)

---

## Problem

Sometimes you want a dependency to live *inside* your app rather than as an external
package — so you can own it, modify it, and still pull upstream updates later. Not a
one-time copy-paste, not a full package manager reimplementation.

---

## Core Design Decisions

- **Source of truth: npm tarball, not git source.** npm builds/transpiles output can
  differ meaningfully from git source. Vendoring the tarball guarantees you get exactly
  what a real `npm install` would give you, and makes version pinning trivial (registry
  handles exact version resolution — no unreliable npm-name → git-repo mapping needed).
- **No recursive dependency vendoring.** Only the target package's source is vendored.
  Its own `dependencies` stay declared in its `package.json` and resolve normally via
  the package manager. Vendoring the full transitive tree would mean reimplementing
  npm's resolver — explicitly out of scope.
- **Workspaces, not raw folders.** Each vendored package gets its own `package.json`
  in `vendor/<pkg>` and is registered as a workspace member. This lets the existing
  package manager (npm/yarn/pnpm) do all resolution/hoisting/dependency installation
  for the vendored package's own deps — no custom resolver logic needed.
- **Scoped packages mirror npm structure.** `@org/pkg` → `vendor/@org/pkg`. No
  flattening scheme; workspaces and `imports` both handle scoped paths fine natively.
- **Patch-based update reconciliation**, not git merge. Since there's no ongoing git
  history relationship with upstream (tarball-based, not git clone), local
  modifications are tracked as a diff against a baseline snapshot, regenerated and
  reapplied against each new version on update.
- **Patch application relies on the parent repo's git**, not a nested git repo per
  vendored package. The vendored files are already tracked as normal files in the
  parent repo's working tree — `git apply` operates from repo root. No `.git` init
  needed inside `vendor/<pkg>`. Requires git installed on PATH; tool should fail
  clearly if git is missing (not silently degrade).
- **Explicit flags over interactive prompts by default**, to keep the tool CI-friendly:
  `--yes`/`-y` skips confirmation prompts (e.g. uncommitted-changes warning before update).
- **Update-checking is opt-in, not default.** `vendor status` reads local manifest state
  only (fast, offline). `vendor status --check-updates` additionally hits the npm
  registry per package to compare against latest — kept as an explicit flag since it's
  a network call per package and could be slow with many vendored packages.

---

## Keeping vendored code out of `node_modules` by default

Workspace membership is what triggers auto-symlinking into `node_modules`. Two paths:

1. **Not a workspace member (default-safe idea, later superseded):** package sits in
   `vendor/<pkg>` as a plain folder, imported via the `imports` field subpath alias
   (e.g. `#vendor/some-lib`) — pure Node.js resolution, no `node_modules` involvement,
   no bundler alias config needed (Node 4.7+/TypeScript/Vite all defer to Node's own
   resolution for `imports`).
2. **Workspace member (the path ultimately chosen):** registered in `workspaces`
   (npm/yarn) or `pnpm-workspace.yaml` (pnpm), so the package manager symlinks it into
   `node_modules` and installs its own `dependencies` normally.

**Decision: always workspace-register (steps 4–5 of the `vendor add` flow are never
skipped)**, even for zero-dependency packages — consistency over marginal savings, and
avoids conditional detection logic. This means vendored packages *do* land in
`node_modules` (as a symlink) and in the lockfile — necessary so their own transitive
dependencies resolve correctly. The "kept out of dependencies" framing refers to the
root `dependencies`/`devDependencies` list, not `node_modules` presence.

The `imports` field alias is still used on top of this, so code reads as
`import { x } from "#vendor/some-lib"` regardless of the underlying resolution
mechanism.

---

## `vendor add <pkg>[@version]` — full flow

1. **Resolve + fetch** — hit npm registry for the package (latest or pinned version),
   get tarball URL, download, unpack into `vendor/<pkg>/` (mirroring npm's structure
   for scoped packages: `vendor/@org/pkg`)
2. **Snapshot baseline** — copy of unpacked files stashed (e.g. `.vendor-cache/<pkg>-<version>/`)
   as the "before" state for future patch diffing. Plain file copy, not a git repo.
3. **Detect package manager** — check root `package.json` `packageManager` field
   (corepack-style) or lockfile presence (`package-lock.json` / `pnpm-lock.yaml` /
   `yarn.lock`). If multiple or none detected, prompt the user.
4. **Register as workspace** — add `vendor/<pkg>` to root `workspaces` array
   (npm/yarn) or `pnpm-workspace.yaml` packages list. Always done, never conditional.
5. **Run install** — trigger the detected package manager's install command so it
   symlinks the vendored package into `node_modules` and resolves its own dependencies.
6. **Write import map entry** — add to root `package.json` `imports`:
   ```json
   "imports": {
     "#vendor/fast-decode-uri": "./vendor/fast-decode-uri/index.js"
   }
   ```
7. **Record in manifest** — write to `vendor.yaml` and `vendor-lock.json` (see schema below).

---

## CLI Commands

| Command | Purpose |
|---|---|
| `vendor add <pkg>[@version]` | Full add flow (steps above) |
| `vendor remove <pkg>` | Delete `vendor/<pkg>`, unregister workspace entry, remove import map entry, drop manifest entries, re-run install to clean up the `node_modules` symlink |
| `vendor update <pkg>[@version]` (no arg = all) | Fetch new version, diff/reapply local patch if present, flag conflicts, update baseline + version |
| `vendor diff <pkg>` | Show local modifications vs. last snapshot baseline |
| `vendor status` | List vendored packages + patched state, offline/local-only by default |
| `vendor status --check-updates` | Same as above, plus registry check for newer versions available |
| `vendor sync` / `vendor install` | Re-apply workspace registration + import map + install; useful after a fresh clone if `package.json` was hand-edited |
| `--yes` / `-y` flag | Skips interactive confirmation prompts (e.g. uncommitted-changes warning), for CI/non-interactive use |

---

## `vendor update` flow (with patch reconciliation)

1. **Determine target version** — explicit `@version`, or latest from registry if omitted
2. **Fetch new tarball** to a temp location (not directly into `vendor/<pkg>` yet)
3. **Check `patched` flag** in `vendor.yaml`:
   - If `false`: simply replace `vendor/<pkg>` with new tarball contents, update
     baseline hash + version, done.
   - If `true`: continue to patch reconciliation.
4. **Pre-update safety check** — run `git status --porcelain -- vendor/<pkg>`. If
   there are uncommitted changes, warn the user:
   > ⚠ `vendor/<pkg>` has uncommitted changes. `vendor update` will treat all current
   > differences from the baseline as your patch. If you're mid-edit, commit or stash
   > first, then re-run.
   Prompt to continue (`y/N`), or skip prompt entirely with `--yes`.
5. **Generate patch** — diff current `vendor/<pkg>` contents against the stored
   baseline snapshot (`.vendor-cache/<pkg>-<old-version>/`), producing a `.patch` file
   representing local modifications only (independent of the version bump itself).
6. **Apply patch to new version** — `git apply` the generated patch against the
   freshly fetched new-version files (parent repo's git, no per-package `.git` needed).
7. **Success:** patch applies cleanly → write patched result into `vendor/<pkg>`,
   update baseline to the new unpatched snapshot, bump version in both manifest files.
8. **Conflict:** patch doesn't apply cleanly (upstream changed the same lines) →
   don't silently overwrite. Write the new tarball as-is into `vendor/<pkg>`, save
   the failed patch (e.g. via `git apply --reject`, producing `.rej` files for the
   parts that failed), back up the old patched version (e.g.
   `.vendor-backup/<pkg>-<old-version>/`), and print a clear message telling the user
   to manually reconcile.

**Git requirement:** `vendor update` on a patched package requires git installed and
the tool running inside a git repo. `vendor add` and `vendor diff` don't need git.
Tool should check for `.git` existence upfront and fail with a clear error message
rather than a cryptic `git apply` failure if git/repo context is missing.

---

## Manifest Schema

**`vendor.yaml`** — human-facing, reviewed in PRs, hand-editable:

```yaml
version: 1

packages:
  fast-decode-uri:
    version: "1.0.3"           # pinned semver
    source: npm                # npm | git (future)
    path: vendor/fast-decode-uri
    import: "#vendor/fast-decode-uri"
    patched: false             # true once local mods exist
    notes: ""                  # optional freeform, e.g. why it's vendored
```

**`vendor-lock.json`** — machine-owned, not meant for hand-editing:

```json
{
  "lockfileVersion": 1,
  "packages": {
    "fast-decode-uri": {
      "version": "1.0.3",
      "resolved": "https://registry.npmjs.org/fast-decode-uri/-/fast-decode-uri-1.0.3.tgz",
      "integrity": "sha512-...",
      "baselineHash": "a1b2c3...",
      "vendoredAt": "2026-07-21T10:00:00Z",
      "patchFile": null
    }
  }
}
```

Notes on the split:
- `patched` lives in the human file (visible in code review), but diff mechanics
  (`patchFile` path, `baselineHash`) live in the lock since they're derived/generated.
- `import` is stored per-package rather than assumed, in case of a custom alias
  instead of the default `#vendor/<name>` pattern.
- `path` is stored explicitly (even though predictable) to handle scoped package
  naming (`@org/pkg`) cleanly.

---

## Go Implementation Notes

- **npm registry calls:** plain `net/http` + `encoding/json`
- **Tarball extraction:** `archive/tar` + `compress/gzip`, both standard library —
  no external deps needed for `.tgz` unpacking
- **YAML parsing** (`vendor.yaml`): `gopkg.in/yaml.v3`
- **Git shelling:** `os/exec` calling the user's installed `git` binary — matches the
  decision to rely on the parent repo's git rather than vendoring a git library or
  requiring nested repos
- **`package.json` manipulation:** needs to preserve field order/formatting reasonably
  well since it's a file humans also hand-edit. Go's `encoding/json` doesn't preserve
  key order on unmarshal/remarshal — consider a surgical-edit approach (e.g. something
  like `github.com/tidwall/sjson`) instead of full unmarshal/remarshal.
- **Cross-compilation:** trivial via `GOOS`/`GOARCH` env vars, no cross-compiler
  toolchain needed — straightforward CI build matrix for all target platforms.

---

## Distribution via npm (no postinstall scripts)

Same pattern used by esbuild, swc, turbo:

- Publish separate tiny platform packages, e.g. `vendor-cli-darwin-arm64`,
  `vendor-cli-linux-x64`, `vendor-cli-win32-x64` — each containing just the prebuilt
  Go binary for that target.
- Each platform package declares `"os": [...]` and `"cpu": [...]` in its
  `package.json` — npm uses these to skip installing packages that don't match the
  current machine.
- The main package (`vendor-cli`) lists all platform packages under
  `optionalDependencies`; npm only actually downloads the one matching the user's
  platform — no network fetch at install time beyond normal npm resolution.
- Main package's `bin` field points to a small JS shim (or directly to the binary
  path) that execs the matching platform binary.
- CI cross-compiles the Go binary per target (`GOOS`/`GOARCH`), then publishes each
  platform package alongside the main package on release.

No postinstall script, no download-on-install step — npm's own platform-matching
resolution (`os`/`cpu` fields) handles everything.

---

## Open / Deferred Items

- `vendor init` (first-time setup: PM detection, initial `vendor.yaml`/lock creation,
  `imports` scaffold) — discussed as a next step, not yet detailed.
- Actual Go project structure/package layout — not yet sketched.
- GitHub Actions build matrix specifics — not yet sketched.
- Git-source vendoring as an opt-in fallback (vs. npm tarball default) — mentioned
  as a possibility for packages where source matters more than build output, not
  fully designed.

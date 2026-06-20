# bigmark — VS Code / VSCodium extension

Insert ASCII section-banner comments (minimap landmarks) into the active editor
by shelling out to the `bigmark` CLI. Works in VS Code, VSCodium, and other
Code OSS builds (it uses only the stable extension API).

## Prerequisites

- The `bigmark` binary built and on your `PATH` (or point `bigmark.binaryPath`
  at it). From the repo root: `make build` produces `./bigmark`.
- `figlet` installed (bigmark depends on it). On Debian/Ubuntu it lands at
  `/usr/bin/figlet`, which bigmark finds automatically. If your editor can't
  find it, set `bigmark.figletPath` — the extension passes it through as
  `BIGMARK_FIGLET`, because the editor's environment often has a stripped
  `PATH`.

## Flatpak VSCodium (figlet bridge)

If you run VSCodium (or VS Code) as a **Flatpak**, the editor is sandboxed and
**cannot see the host's `/usr`** — so even though `figlet` is installed on your
system, the extension's child process can't reach `/usr/bin/figlet`, and bigmark
fails with "required dependency 'figlet' was not found." The Flatpak only
exposes your home directory (`filesystems=host` covers `/home`, not `/usr`).

`figlet` links only against libc (which the Flatpak runtime has), so the clean
fix is a **self-contained copy under your home directory** that the sandbox can
run directly — no `flatpak-spawn`, no host round-trip:

```sh
# 1. copy the figlet binary + fonts somewhere the sandbox can see
mkdir -p ~/.local/lib/figlet/fonts
cp "$(readlink -f /usr/bin/figlet)" ~/.local/lib/figlet/figlet
cp -r /usr/share/figlet/* ~/.local/lib/figlet/fonts/

# 2. a one-line shim that supplies the font dir (figlet's compiled-in default
#    path /usr/share/figlet doesn't exist in the sandbox)
cat > ~/.local/bin/figlet-host <<'EOF'
#!/bin/sh
exec ~/.local/lib/figlet/figlet -d ~/.local/lib/figlet/fonts "$@"
EOF
chmod +x ~/.local/bin/figlet-host
```

Then point the extension at it in VSCodium settings:

```jsonc
"bigmark.figletPath": "/home/<you>/.local/bin/figlet-host"
```

Also set `bigmark.binaryPath` to the **absolute** path of the bigmark binary —
a bare `bigmark` won't resolve in the sandbox's `PATH`.

> The sandbox boundary forces *some* bridge here; this is the most self-contained
> form. Trade-off: the copied figlet won't pick up a system `apt upgrade figlet`
> (figlet is stable, so this rarely matters).

## Install

From this directory, build it:

```sh
npm install
npm run compile
```

Then pick one:

### Option A — Extension Development Host (no CLI needed)

Open this folder in VSCodium and press **F5**. This launches a second window
with the extension loaded. Best for trying it out; nothing is installed
permanently.

### Option B — Package a `.vsix` and install it

```sh
npx --yes @vscode/vsce package        # produces bigmark-0.1.0.vsix
```

Install via the UI — **Extensions panel → ⋯ menu → "Install from VSIX…"** — or,
if you have the CLI on PATH, `codium --install-extension bigmark-0.1.0.vsix`
(it's `codium`, not `code`, on VSCodium).

> Note: VSCodium uses the [Open VSX](https://open-vsx.org) registry, not
> Microsoft's marketplace. Local `.vsix` install works regardless. To publish to
> Open VSX, use [`ovsx`](https://github.com/eclipse/openvsx/tree/master/cli)
> rather than `vsce publish`.

## Usage

- Run **bigmark: Insert banner** from the Command Palette, or press
  **`Ctrl+K Ctrl+B`** with the editor focused.
- If you have text selected, it pre-fills the label and the banner replaces the
  selection; otherwise the banner is inserted at the cursor.

The comment style is chosen automatically from the active file's language:

- **Non-canvas modes** (tier1/2/3, rotated) use a per-line prefix via `-c`
  (`// ` for C-like languages, `# ` for Python/shell/YAML, `-- ` for SQL/Lua,
  etc.). Override with `bigmark.commentPrefix`.
- **Canvas mode** wraps the art in a block comment via `--lang`
  (`js/ts/php/css/html/python/ruby`), falling back to `js`.

## Settings

| Setting | Default | Description |
| --- | --- | --- |
| `bigmark.binaryPath` | `bigmark` | Path to the bigmark executable. |
| `bigmark.mode` | `tier1` | `tier1` \| `tier2` \| `tier3` \| `rotated` \| `canvas`. |
| `bigmark.width` | `80` | Max line width (`-w`). |
| `bigmark.commentPrefix` | `""` | Override the auto-derived `-c` prefix (non-canvas modes). |
| `bigmark.figletPath` | `""` | Path to figlet, passed as `BIGMARK_FIGLET`. |

## Not in this version

- The optional second "note" argument (one-step prompt only).
- `--rand` / `--frand` / `--seed`, `--brush`, `--vstretch`, `--solid`.
- Marketplace / Open VSX publishing (local `.vsix` / dev host only).

# bigmark (Go) — project notes for Claude

Go port of the PHP `bigmark`. Single static binary that shells out to `figlet`.

## Layout

```
cmd/bigmark/main.go     flag parsing, seed/random derivation, dispatch
cmd/bigmark/help.go     -h help text and --brushes ramp
internal/render/        the five renderers + figlet helper
  figlet.go             figlet exec wrapper, padding/maxLen helpers, die()
  render.go             Tier1/Tier2/Tier3/Rotated/Canvas
docs/DESIGN.md          design rationale (shared with the PHP repo)
docs/bigmark.php.reference   the original PHP, kept for diffing during the port
```

Module path: `github.com/GravisTechGregB/bigmark`.

## Parity with the PHP version

All five modes produce **byte-for-byte identical** output to the PHP tool for a
given fixed seed and inputs. This is the regression check — if you change a
renderer, diff against `docs/bigmark.php.reference`:

```sh
go build -o bigmark ./cmd/bigmark
diff <(./bigmark --canvas --seed 100 "CSS") <(php docs/bigmark.php.reference --canvas --seed 100 "CSS")
```

The ONE intentional divergence: `--rand`/`--frand` with no explicit `--seed`
pick different fonts/brushes than PHP, because Go's `math/rand` != PHP's Mersenne
Twister. Within Go a seed always reproduces; that's what matters.

## Key design decisions (same as PHP)

- **Reproducibility via seed.** Random modes print a `--seed N` replay command.
  `--seed` does NOT force a mode; mode comes from `--canvas` presence.
- **Bold-only `--rand` font pool** (`banner`/`block`); thin fonts carve to mush.
  `--frand` uses the full pool and works in any mode without enabling canvas.
- **Exact widths.** Frame/canvas rows built to a fixed display-column count.
- **`-r` forces `banner`** — only stock font whose strokes survive 90° rotation.

## Gotchas

- Width math assumes single-column brushes. figlet output is ASCII (1 byte/char)
  so `len()` is fine there, but a multi-rune `--brush` would break alignment.
- The help text lives in `help.go` as a raw string — edit it when adding flags.
- Hand-rolled arg parser (not `flag`) to support `--rand`/`-rand` dual-dash and
  positional-args-after-flags, matching the PHP behavior.

## Build / test

```sh
make build     # -> ./bigmark
make test      # smoke test all modes
make fmt vet   # format + static checks
```

# bigmark

Generate ASCII section banners sized for code comments, designed to be
recognizable **landmarks in the editor minimap**. Navigate a long file by the
silhouette of its sections instead of scrolling and reading.

The minimap is the intended viewing distance: up close a banner is just a block
of comment characters; zoomed out it reads as a shape, a word, or a bold bar you
can jump to.

> Go port of the original PHP tool. Single static binary; requires only
> [`figlet`](http://www.figlet.org/) on your PATH for the letterforms.

## Install

Download a binary from [Releases](../../releases), or build from source:

```sh
go install github.com/GravisTechGregB/bigmark/cmd/bigmark@latest
# or
git clone … && cd bigmark && go build -o bigmark ./cmd/bigmark
```

Then ensure `figlet` is installed:

```sh
sudo apt install figlet     # Debian/Ubuntu
brew install figlet         # macOS
```

## Quick start

```sh
bigmark "parse" "raw bytes -> tokens"   # tier 1: full-width framed figlet
bigmark -2 "validate" "reject bad"      # tier 2: one-line box
bigmark -3 "normalize whitespace"       # tier 3: tick
bigmark -r "parse"                      # rotated 90: reads top-to-bottom
bigmark --canvas "parse"                # carve the word into a filled field
bigmark --canvas --rand "parse"         # random style + an exact-replay line
bigmark -h                              # full help
```

## The five modes

| Mode | Flag | Use |
|---|---|---|
| Framed figlet | (default) | Major section landmark, full-width frame |
| Box | `-2` | A section divider, one line tall |
| Tick | `-3` | A subsection marker, minimal |
| Rotated | `-r` | Tall vertical landmark; word runs down the file |
| Canvas | `--canvas` | Word carved as a shape into a filled minimap field |

## Reproducible randomness

`--rand` (canvas) and `--frand` (font, any mode) pick a random **seed**, derive
all settings from it, and print the exact command to regenerate the identical
output. Replaying uses `--seed N`, never `--rand`, so a render you like is frozen
forever.

```sh
$ bigmark --canvas --rand "PARSE"
// bigmark --canvas --seed 4821 --lang js -f banner --brush '█' --vstretch 2 -w 80 'PARSE'
/* ...art... */
```

> Seeds are reproducible within this Go build, but **not** across the Go and PHP
> versions (different RNGs). Pick one and stay with it.

## Language-aware comment syntax (canvas)

`--lang` wraps the canvas in the right comment delimiters: `js` `ts` `php` `css`
(`/* */`), `html` (`<!-- -->`), `python`/`py` (`"""..."""`), `ruby`/`rb`
(`=begin`/`=end`).

## Notes & limitations

- Canvas art only reads in the minimap when VS Code's
  `"editor.minimap.renderCharacters": true` is set.
- Only **bold** fonts (`banner`, `block`) carve legibly; thin/diagonal fonts come
  out as mush. `--rand` rolls only the bold ones; force any font with `-f`.
- `-r` uses the `banner` font (the only stock font that stays legible rotated).

See [`docs/DESIGN.md`](docs/DESIGN.md) for the full design rationale.

## License

MIT — see [LICENSE](LICENSE).

# bigmark design notes

## The idea

A name (or section banner) does two different jobs depending on how you read it:

- **Explanation** — read line by line, up close. The name carries meaning.
- **Wayfinding** — read at a glance, while jumping around. The shape carries
  position.

`bigmark` is about the second job. The editor minimap is a zoomed-out map of the
file; a banner comment, designed right, becomes a landmark on that map. You stop
counting separators and start recognizing terrain.

## Why the minimap is the unit of resolution

At minimap zoom each source line renders as roughly one pixel row, and the
columns are squished too. So:

- **Big bold shapes read; fine detail vanishes.** Work in blocks, like a tapestry
  meant to be seen from across the room.
- **Silhouette beats text.** A full-width `+===+` rule is two solid horizontal
  bars; a rotated word is a vertical streak; carved negative space is a dense
  block with holes. Each is a distinct silhouette.
- **Distinct beats descriptive.** Sections must differ from *each other* at a
  blur — different lengths, different first letters, different shapes.

## The tiers (visual weight = importance)

1. **Framed figlet** — major movements. Full-width frame = strongest landmark.
2. **Box** — sections. A single bold bar.
3. **Tick** — subsections. A lightweight anchor.

## Rotated mode

Turning figlet 90° makes the word run down the file's long axis — a tall
vertical landmark. Critical constraint: **only fonts with orthogonal strokes
survive rotation.** Diagonal-stroke fonts (`standard`, `slant`) turn to mush
when rotated because `/` and `\` point the wrong way. `banner` (pure `#`,
no diagonals) is the only stock font that stays legible, so `-r` forces it.

## Canvas mode (carving)

The whole 80×N comment block is a canvas. Render the word big, then either:

- **carve** — fill the field, punch the word out as negative space, or
- **solid** — empty field, fill the word.

Negative space often reads better in a minimap because a dense block catches the
eye and the holes form the letters.

Two hard constraints:

1. **`renderCharacters` must be on.** VS Code's minimap can render tiny real
   glyphs (art shows) or solid color blocks per token (art vanishes — the whole
   comment becomes one flat rectangle). Canvas mode requires the former.
2. **Brush must be an area-filler, not a line.** `█` (and `▓▒░` for shading,
   `#@MW` as ASCII fills) read as ink when tiny. Diagonal/line chars (`/\|_`)
   do not fill area — useless for carving.

### The brush-density ramp

ASCII artists pick characters by ink weight, not meaning:

```
ASCII 10-step :  .:-=+*#%@
Unicode shade :  ░▒▓█
half blocks   :  ▀▄▌▐     (sub-line detail; beat the vertical squish)
fat ASCII fill:  # @ M W & 8 0
```

### Carve-safe fonts

Only **bold** fonts carve into readable letters. Tested the full stock set on
real words (not isolated glyphs — that's misleading):

- **Good:** `banner`, `block`.
- **Mush when carved:** `standard`, `big`, `slant`, `small`, `shadow`, `script`,
  and the `sm*` variants — thin/diagonal strokes leave noisy holes.
- **Collapse:** `digital` (becomes a tiny box), `mini`/`bubble` (too thin/round).

So `--rand` rolls only `banner`/`block`. `--frand` rolls the full pool for when
you explicitly want variety and will eyeball the result.

## Reproducible randomness

Random output you can't reproduce is useless. So `--rand`/`--frand` pick a
**seed** and derive every setting deterministically from it via `mt_srand`. The
printed replay command pins `--seed N` (never `--rand`), freezing the render.
`--seed` deliberately does not select a mode, so a tier-1 replay line stays
tier-1 and a canvas one stays canvas.

## Vertical pre-stretch

The minimap squishes vertically. `--vstretch N` repeats each canvas row N times
so the art ends up proportioned (not flat) once the minimap squashes it — the
same trick low-res sprite artists used to pre-distort for non-square pixels.

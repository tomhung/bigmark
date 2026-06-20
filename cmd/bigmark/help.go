package main

import "fmt"

func printHelp() {
	fmt.Print(`bigmark — generate ASCII section banners sized for code comments, designed to
be recognizable landmarks in the editor minimap.

THREE TIERS (visual weight = section importance):
  (default / -1)  full-width framed figlet  -> major movements
  -2              one-line box               -> sections
  -3              lightweight tick           -> subsections

USAGE
  bigmark "parse" "raw bytes -> tokens"        # tier 1, big framed figlet
  bigmark -2 "validate" "reject bad tokens"    # tier 2, boxed
  bigmark -3 "normalize whitespace"            # tier 3, tick
  bigmark -r "parse"                           # ROTATED 90: reads top-to-bottom
  bigmark --canvas "parse"                     # CARVE word into a filled field
  bigmark --canvas --rand "parse"              # random style + an exact-replay line
  bigmark --frand "parse"                      # random font, any mode

OPTIONS
  -r          rotated mode: figlet turned 90deg so the word runs DOWN the file
              (a tall vertical landmark, centered). Forces the 'banner' font,
              the only stock font that stays legible rotated. -ccw reverses.
  -ccw        rotated counter-clockwise (with -r)
  --canvas    canvas mode: render the word as a SHAPE in a filled field so it
              reads in the minimap. By default the word is CARVED as negative
              space; --solid fills the word instead.
  --solid     canvas mode: fill the word instead of carving it
  --brush C   fill character for canvas mode (default █ full block). See --brushes
  --brushes   print the brush-density ramp and exit
  --vstretch N  repeat each canvas row N times to counter the minimap's vertical
              squish (default 2; use 1 to disable)
  --rand      canvas: randomize brush, font (bold-only), vstretch, carve/solid
              from a random seed. Prints an exact-replay command above the art.
  --frand     randomize the FONT only, from the FULL pool. Works in ANY mode and
              does NOT enable canvas by itself.
  --seed N    replay a specific render. Same seed -> identical output.
  --lang L    canvas comment syntax: js|ts|php|css (/* */), html (<!-- -->),
              python|py (""" """), ruby|rb (=begin/=end). Default js.
  -w N        max line width (default 80)
  -c "PFX "   comment prefix (default "# "), e.g. -c "// "
  -f FONT     figlet font. Any font in tier 1; only bold fonts (banner, block)
              carve legibly in canvas mode (thin fonts come out as mush).
  --header    print a one-line definition of the convention and exit
  -h          this help

Requires: figlet on PATH.
`)
}

func printBrushes() {
	fmt.Print(`Brush density ramp (lightest -> darkest), pick a --brush:

  ASCII 10-step :  .:-=+*#%@
  Unicode shade :  ░▒▓█   (light, medium, dark, full block)
  half blocks   :  ▀▄▌▐   (top, bottom, left, right - sub-line detail)
  fat ASCII fill:  # @ M W & 8 0

Minimap tip: solid area-fill reads best. Default brush is the full block █.
Diagonal/line chars (/ \ | _) do NOT fill area - avoid for canvas.
`)
}

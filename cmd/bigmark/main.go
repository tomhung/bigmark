// bigmark — generate ASCII section banners sized for code comments, designed to
// be recognizable landmarks in the editor minimap. Go port of the original PHP
// tool; shells out to figlet for letterforms.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/GravisTechGregB/bigmark/internal/config"
	"github.com/GravisTechGregB/bigmark/internal/render"
)

const fullBlock = "█" // █

func main() {
	args := os.Args[1:]

	// defaults
	tier := 1
	width := 80
	prefix := "# "
	font := ""   // tier-1 font ("" = auto)
	rotate := "" // "", "cw", "ccw"
	canvas := false
	carve := true
	cfont := "banner" // canvas mask font
	brush := fullBlock
	vstretch := 2
	lang := "js"
	var seed *int
	doRand := false
	doFRand := false

	// track explicitly-set knobs so random modes don't override them
	setBrush, setVStretch, setCarve, setCFont := false, false, false, false

	// ---- config file: overrides built-in defaults, overridden by flags ----
	// (precedence: built-in defaults -> config -> command-line flags)
	if cfg, err := config.Load(); err == nil {
		if cfg.Width != nil {
			width = *cfg.Width
		}
		if cfg.Prefix != nil {
			prefix = *cfg.Prefix
		}
		if cfg.Brush != nil {
			brush = *cfg.Brush
		}
		if cfg.Lang != nil {
			lang = *cfg.Lang
		}
		if cfg.Font != nil {
			font = *cfg.Font
			cfont = *cfg.Font
		}
		if cfg.VStretch != nil {
			vstretch = *cfg.VStretch
		}
		if cfg.Carve != nil {
			carve = *cfg.Carve
		}
	}

	var pos []string
	next := func(i *int) string {
		*i++
		if *i < len(args) {
			return args[*i]
		}
		return ""
	}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			printHelp()
			return
		case a == "--header":
			fmt.Println(prefix + "bigmark: full-width figlet banner comments used as minimap landmarks.")
			fmt.Println(prefix + "         major sections get one; navigate by silhouette.")
			return
		case a == "--config":
			// print the config path being read and whether it exists
			p := config.Path()
			if _, err := os.Stat(p); err == nil {
				fmt.Printf("%s (exists)\n", p)
			} else {
				fmt.Printf("%s (not found)\n", p)
			}
			return
		case a == "-1":
			tier = 1
		case a == "-2":
			tier = 2
		case a == "-3":
			tier = 3
		case a == "-r":
			if rotate == "" {
				rotate = "cw"
			}
		case a == "-ccw":
			rotate = "ccw"
		case a == "--canvas":
			canvas = true
		case a == "--solid":
			canvas = true
			carve = false
			setCarve = true
		case a == "--brush":
			brush = next(&i)
			if brush == "" {
				brush = fullBlock
			}
			setBrush = true
		case a == "--brushes":
			printBrushes()
			return
		case a == "--vstretch":
			v, _ := strconv.Atoi(next(&i))
			if v < 1 {
				v = 1
			}
			vstretch = v
			setVStretch = true
		case a == "--lang":
			lang = strings.ToLower(next(&i))
		case a == "--seed":
			v, _ := strconv.Atoi(next(&i))
			seed = &v
		case a == "--rand" || a == "-rand":
			doRand = true
			canvas = true
		case a == "--frand" || a == "-frand":
			doFRand = true
		case a == "-w":
			width, _ = strconv.Atoi(next(&i))
			if width == 0 {
				width = 80
			}
		case a == "-c":
			prefix = next(&i)
		case a == "-f":
			font = next(&i)
			if font == "" {
				font = "banner"
			}
			cfont = font
			setCFont = true
		default:
			if strings.HasPrefix(a, "-") {
				fmt.Fprintf(os.Stderr, "unknown option %q (try: bigmark -h)\n", a)
				os.Exit(2)
			}
			pos = append(pos, a)
		}
	}

	label := "SECTION"
	if len(pos) > 0 {
		label = strings.ToUpper(pos[0])
	}
	note := ""
	if len(pos) > 1 {
		note = pos[1]
	}

	// ---- seed / random derivation (deterministic from seed) ----
	if (doRand || doFRand) && seed == nil {
		s := 1000 + rand.Intn(9000) // human-friendly 4-digit seed
		seed = &s
	}
	if seed != nil {
		rng := rand.New(rand.NewSource(int64(*seed)))
		brushPool := []string{fullBlock, "▓", "▒", "#", "@", "M", "W", "8"}
		boldPool := []string{"banner", "block"}
		wildPool := []string{"banner", "block", "standard", "big", "small", "slant",
			"smslant", "shadow", "smshadow", "lean", "script", "smscript"}
		tier1Pool := []string{"standard", "big", "small", "slant", "shadow", "banner",
			"block", "lean", "script", "smslant", "smshadow", "smscript"}

		if doRand && !setBrush {
			brush = brushPool[rng.Intn(len(brushPool))]
		}
		if doRand && !setVStretch {
			vstretch = 1 + rng.Intn(3)
		}
		if doRand && !setCarve {
			carve = rng.Intn(2) == 1
		}
		if !setCFont {
			if doFRand {
				cfont = wildPool[rng.Intn(len(wildPool))]
			} else if doRand {
				cfont = boldPool[rng.Intn(len(boldPool))]
			}
		}
		if doFRand && font == "" {
			font = tier1Pool[rng.Intn(len(tier1Pool))]
		}
	}

	// every render mode needs figlet — verify it up front with a clear message
	// (help/header/config/brushes already returned above and don't need it).
	render.CheckFiglet()

	// ---- dispatch ----
	switch {
	case canvas:
		render.Canvas(render.CanvasOpts{
			Word: label, Max: width, Brush: brush, Carve: carve,
			VStretch: vstretch, Lang: lang, Seed: seed, Font: cfont,
		})
	case rotate != "":
		render.Rotated(label, prefix, rotate, width)
	case tier == 3:
		render.Tier3(label, note, prefix)
	case tier == 2:
		render.Tier2(label, note, prefix, width)
	default:
		// tier 1: emit a reproduce line if a random flag picked the font
		if (doFRand || doRand) && seed != nil && font != "" {
			repro := fmt.Sprintf("%sbigmark --seed %d -f %s -w %d %s", prefix, *seed, font, width, shellQuote(label))
			if note != "" {
				repro += " " + shellQuote(note)
			}
			fmt.Println(repro)
		}
		render.Tier1(label, note, prefix, width, font)
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// Package render holds the five bigmark output modes: tier1 (framed figlet),
// tier2 (box), tier3 (tick), rotated (90deg figlet), and canvas (carve/fill).
package render

import (
	"fmt"
	"strings"
)

// Tier3 emits a lightweight subsection tick.
func Tier3(label, note, prefix string) {
	text := strings.ToLower(strings.TrimSpace(label + " " + note))
	fmt.Printf("%s--> %s\n", prefix, text)
}

// Tier2 emits a one-line ASCII box padded to max width.
func Tier2(label, note, prefix string, max int) {
	text := label
	if note != "" {
		text = label + " - " + note
	}
	inner := max - len(prefix) - 4 // "| " + " |"
	if inner < len(text) {
		text = text[:inner]
	}
	barw := max - len(prefix) - 2
	bar := prefix + "+" + strings.Repeat("=", barw) + "+"
	mid := prefix + "| " + padRight(text, inner) + " |"
	fmt.Printf("%s\n%s\n%s\n", bar, mid, bar)
}

// Tier1 renders a figlet word centered inside a full-width frame, with balanced
// top/bottom padding and an optional centered subtitle. font=="" auto-selects.
func Tier1(word, note, prefix string, max int, font string) {
	plen := len(prefix)
	inner := max - plen - 4
	if inner < 10 {
		die("width %d too small for a frame", max)
	}

	fonts := []string{"standard", "small", "mini"}
	if font != "" {
		fonts = []string{font}
	}
	var art []string
	var lastErr error
	for _, f := range fonts {
		a, ok, err := figletFit(word, f, inner)
		if ok {
			art = a
			break
		}
		if err != nil {
			lastErr = err // figlet broke, not just overflow
		}
	}
	if art == nil {
		if lastErr != nil {
			die("could not render with figlet: %v", lastErr)
		}
		die("%q won't fit in width %d even at the smallest font (try a shorter word or wider -w).", word, max)
	}

	barw := max - plen - 2
	bar := prefix + "+" + strings.Repeat("=", barw) + "+"

	emit := func(line string) {
		pad := inner - len(line)
		if pad < 0 {
			pad = 0
		}
		l := pad / 2
		r := pad - l
		fmt.Printf("%s| %s%s%s |\n", prefix, strings.Repeat(" ", l), line, strings.Repeat(" ", r))
	}

	fmt.Println(bar)
	emit("") // top padding (balances the frame)
	for _, line := range art {
		emit(line)
	}
	if note != "" {
		emit("")
		if len(note) > inner {
			note = note[:inner-1] + "."
		}
		emit(note)
	} else {
		emit("") // bottom padding when there's no note
	}
	fmt.Println(bar)
}

// Rotated renders the figlet 'banner' font turned 90deg so the word reads
// top-to-bottom, centered within the width. dir is "cw" or "ccw".
func Rotated(word, prefix, dir string, max int) {
	rows, err := figlet(word, "banner", 0)
	if err != nil {
		die("figlet 'banner' font not available")
	}
	w := maxLen(rows)
	for i := range rows {
		rows[i] = padRight(rows[i], w)
	}
	h := len(rows)
	grid := make([][]byte, h)
	for i, r := range rows {
		grid[i] = []byte(r)
	}

	var out []string
	if dir == "cw" {
		// clockwise: column x, bottom-to-top -> reads top-to-bottom
		for x := 0; x < w; x++ {
			var sb strings.Builder
			for y := h - 1; y >= 0; y-- {
				sb.WriteByte(grid[y][x])
			}
			out = append(out, sb.String())
		}
	} else {
		// counter-clockwise
		for x := w - 1; x >= 0; x-- {
			var sb strings.Builder
			for y := 0; y < h; y++ {
				sb.WriteByte(grid[y][x])
			}
			out = append(out, sb.String())
		}
	}

	// center the rotated block: every line is h chars wide; pad all lines by
	// the same left margin so the column stays a block, centered in the usable
	// width (max minus the comment prefix).
	avail := max - len(prefix)
	left := (avail - h) / 2
	if left < 0 {
		left = 0
	}
	pad := strings.Repeat(" ", left)
	for _, line := range out {
		fmt.Printf("%s%s\n", prefix, strings.TrimRight(pad+line, " "))
	}
}

// CanvasOpts bundles the many canvas parameters.
type CanvasOpts struct {
	Word     string
	Max      int
	Brush    string
	Carve    bool
	VStretch int
	Lang     string
	Seed     *int // nil = no reproduce line
	Font     string
}

var canvasDelims = map[string][2]string{
	"js": {"/*", "*/"}, "ts": {"/*", "*/"}, "php": {"/*", "*/"}, "css": {"/*", "*/"},
	"html":   {"<!--", "-->"},
	"python": {`"""`, `"""`}, "py": {`"""`, `"""`},
	"ruby": {"=begin", "=end"}, "rb": {"=begin", "=end"},
}

var canvasLineComment = map[string]string{
	"js": "// ", "ts": "// ", "php": "// ", "python": "# ", "py": "# ",
	"ruby": "# ", "rb": "# ", // css/html have no line comment
}

// Canvas carves (or fills) a word as a shape into a filled field, emitted as a
// single block comment so no per-line prefix pollutes the picture.
func Canvas(o CanvasOpts) {
	W := o.Max - 2 // +2 border brushes -> every row is exactly Max display cols
	if W < 10 {
		die("width too small for canvas")
	}

	delim, ok := canvasDelims[o.Lang]
	if !ok {
		die("unknown --lang %q (js|ts|php|css|html|python|ruby)", o.Lang)
	}
	open, close := delim[0], delim[1]

	// reproduce line: the exact command that regenerates this art.
	if o.Seed != nil {
		cmd := fmt.Sprintf("bigmark --canvas --seed %d --lang %s -f %s --brush %s --vstretch %d",
			*o.Seed, o.Lang, o.Font, shellQuote(o.Brush), o.VStretch)
		if !o.Carve {
			cmd += " --solid"
		}
		cmd += fmt.Sprintf(" -w %d %s", o.Max, shellQuote(o.Word))
		if lc, ok := canvasLineComment[o.Lang]; ok {
			fmt.Println(lc + cmd)
		} else {
			fmt.Printf("%s %s %s\n", open, cmd, close)
		}
	}

	rows, err := figlet(o.Word, o.Font, W)
	if err != nil {
		die("figlet font %q not available", o.Font)
	}
	gw := maxLen(rows)
	if gw > W { // too wide: trim mask to fit
		for i := range rows {
			rows[i] = padRight(rows[i], gw)[:W]
		}
		gw = W
	}
	padL := (W - gw) / 2
	if padL < 0 {
		padL = 0
	}

	// pad mask to full canvas width, centered, with blank margin top/bottom
	blank := strings.Repeat(" ", W)
	mask := []string{blank, blank}
	for _, r := range rows {
		r = padRight(r, gw)
		mask = append(mask, strings.Repeat(" ", padL)+r+strings.Repeat(" ", W-padL-gw))
	}
	mask = append(mask, blank, blank)

	// emit: ink = brush, hole = space. carve inverts (word becomes negative
	// space); solid fills the word. Delimiters sit on their own lines.
	fmt.Println(open)
	fmt.Println(strings.Repeat(o.Brush, W+2)) // top brush rule
	for _, line := range mask {
		var sb strings.Builder
		for x := 0; x < W; x++ {
			ink := x < len(line) && line[x] != ' '
			paint := ink
			if o.Carve {
				paint = !ink
			}
			if paint {
				sb.WriteString(o.Brush)
			} else {
				sb.WriteByte(' ')
			}
		}
		row := o.Brush + sb.String() + o.Brush
		for v := 0; v < o.VStretch; v++ {
			fmt.Println(row)
		}
	}
	fmt.Println(strings.Repeat(o.Brush, W+2)) // bottom brush rule
	fmt.Println(close)
}

// shellQuote wraps s in single quotes for the reproduce command (matching PHP's
// escapeshellarg for the common case).
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

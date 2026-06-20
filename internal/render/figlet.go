package render

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// figlet runs the figlet binary with the given font and width cap, returning
// the rendered rows with leading/trailing blank lines stripped. width<=0 means
// no -w cap.
func figlet(word, font string, width int) ([]string, error) {
	args := []string{"-f", font}
	if width > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", width))
	}
	args = append(args, word)
	out, err := exec.Command("figlet", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("figlet font %q not available: %w", font, err)
	}
	raw := strings.TrimRight(string(out), "\n")
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("figlet produced no output for %q", word)
	}
	lines := strings.Split(raw, "\n")
	// drop blank leading/trailing lines
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("figlet produced only blank lines for %q", word)
	}
	return lines, nil
}

// figletFit renders word and returns the rows only if the widest line fits cap.
func figletFit(word, font string, cap int) ([]string, bool) {
	lines, err := figlet(word, font, cap)
	if err != nil {
		return nil, false
	}
	if maxLen(lines) > cap {
		return nil, false
	}
	return lines, true
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func maxLen(lines []string) int {
	m := 0
	for _, l := range lines {
		if len(l) > m {
			m = len(l)
		}
	}
	return m
}

// padRight pads s to length n with spaces (assumes single-byte chars, which is
// true for figlet ASCII output).
func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

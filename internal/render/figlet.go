package render

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// commonFigletPaths are checked when figlet isn't on PATH — covers the case
// where the environment (e.g. VS Code's integrated terminal) has a stripped
// PATH that's missing /usr/bin. Override with $BIGMARK_FIGLET.
var commonFigletPaths = []string{
	"/usr/bin/figlet",
	"/usr/local/bin/figlet",
	"/opt/homebrew/bin/figlet", // macOS arm64 homebrew
	"/bin/figlet",
}

// figletPath resolves the figlet executable: $BIGMARK_FIGLET first, then PATH,
// then known install locations. Returns "" if none is runnable.
func figletPath() string {
	if p := os.Getenv("BIGMARK_FIGLET"); p != "" {
		return p
	}
	if p, err := exec.LookPath("figlet"); err == nil {
		return p
	}
	for _, p := range commonFigletPaths {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

// CheckFiglet verifies the figlet binary is present and runnable, exiting with
// a clear, actionable message if not. Call once at startup so a missing/broken
// figlet fails loudly instead of surfacing as a confusing "won't fit" later.
func CheckFiglet() {
	path := figletPath()
	if path == "" {
		fmt.Fprintln(os.Stderr, "bigmark: required dependency 'figlet' was not found.")
		fmt.Fprintln(os.Stderr, "  Looked on PATH and in: "+strings.Join(commonFigletPaths, ", "))
		fmt.Fprintln(os.Stderr, "  install it:  sudo apt install figlet   (Debian/Ubuntu)")
		fmt.Fprintln(os.Stderr, "               brew install figlet       (macOS)")
		fmt.Fprintln(os.Stderr, "  or point bigmark at it:  export BIGMARK_FIGLET=/path/to/figlet")
		os.Exit(127) // 127 = command not found, conventional
	}
	// present but does it actually run? (broken install, bad perms, etc.)
	if out, err := exec.Command(path, "test").Output(); err != nil || len(out) == 0 {
		fmt.Fprintf(os.Stderr, "bigmark: 'figlet' was found at %s but failed to run.\n", path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error: %v\n", err)
		}
		fmt.Fprintln(os.Stderr, "  try running 'figlet test' yourself to see the problem.")
		os.Exit(127)
	}
}

// figlet runs the figlet binary with the given font and width cap, returning
// the rendered rows with leading/trailing blank lines stripped. width<=0 means
// no -w cap.
func figlet(word, font string, width int) ([]string, error) {
	bin := figletPath()
	if bin == "" {
		return nil, fmt.Errorf("figlet not found (install it, or set BIGMARK_FIGLET)")
	}
	args := []string{"-f", font}
	if width > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", width))
	}
	args = append(args, word)
	out, err := exec.Command(bin, args...).Output()
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
// On failure it also returns the error (nil when the only problem is overflow),
// so callers can distinguish "figlet is broken" from "word is too wide".
func figletFit(word, font string, cap int) ([]string, bool, error) {
	lines, err := figlet(word, font, cap)
	if err != nil {
		return nil, false, err
	}
	if maxLen(lines) > cap {
		return nil, false, nil // genuinely too wide, not an error
	}
	return lines, true, nil
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

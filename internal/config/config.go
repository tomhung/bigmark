// Package config loads default settings from a simple key=value file. Values
// set here override built-in defaults but are themselves overridden by any
// command-line flag the user passes (defaults -> config -> flags).
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds the user-settable defaults. A nil-ish zero value means "not set
// in the file", so the caller only applies keys that were present.
type Config struct {
	Width    *int
	Prefix   *string
	Brush    *string
	Lang     *string
	Font     *string
	VStretch *int
	Carve    *bool // true=carve (default), false=solid
}

// Path returns the config file path that will be read, honoring:
//
//	$BIGMARK_CONFIG                       (explicit override)
//	$XDG_CONFIG_HOME/bigmark/config
//	$HOME/.config/bigmark/config          (fallback)
//
// It does not check existence; use Load for that.
func Path() string {
	if p := os.Getenv("BIGMARK_CONFIG"); p != "" {
		return p
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "bigmark", "config")
}

// Load reads and parses the config file. A missing file is not an error: it
// returns an empty Config and nil. Parse errors on individual lines are
// ignored (best-effort), so a typo never blocks the tool.
func Load() (*Config, error) {
	path := Path()
	if path == "" {
		return &Config{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return &Config{}, err
	}
	defer f.Close()

	c := &Config{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := stripComment(sc.Text())
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(strings.ToLower(key))
		// NOTE: do not TrimSpace the value's leading/trailing on prefix, where a
		// trailing space is meaningful (e.g. prefix="// "). We trim only the
		// outer whitespace that the file format adds, then keep the rest.
		val = strings.Trim(val, "\r\n")
		val = strings.TrimLeft(val, " \t")
		// allow optional quoting to preserve trailing spaces: prefix="// "
		val = unquote(val)

		switch key {
		case "width":
			if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
				c.Width = &n
			}
		case "prefix":
			v := val
			c.Prefix = &v
		case "brush":
			v := strings.TrimSpace(val)
			c.Brush = &v
		case "lang":
			v := strings.ToLower(strings.TrimSpace(val))
			c.Lang = &v
		case "font":
			v := strings.TrimSpace(val)
			c.Font = &v
		case "vstretch":
			if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
				if n < 1 {
					n = 1
				}
				c.VStretch = &n
			}
		case "carve":
			v := strings.ToLower(strings.TrimSpace(val))
			b := v == "true" || v == "1" || v == "yes" || v == "carve"
			c.Carve = &b
		}
	}
	return c, sc.Err()
}

// stripComment removes a trailing `#` comment, but only when the `#` begins a
// line or is preceded by whitespace. This lets `brush=#` keep its `#` value
// while `width=80  # inline note` still strips the note. A `#` inside a quoted
// value is also preserved.
func stripComment(line string) string {
	inQuote := false
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			inQuote = !inQuote
		case '#':
			if inQuote {
				continue
			}
			// comment only if at start or after whitespace
			if i == 0 || line[i-1] == ' ' || line[i-1] == '\t' {
				return line[:i]
			}
		}
	}
	return line
}

// unquote removes a single matching pair of surrounding double quotes, so a
// value like "// " keeps its trailing space.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

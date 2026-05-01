package cmd

import (
	"github.com/target-ops/gitswitch/internal/style"
)

// Color prefixes / reset suffix used inline in command output.
//
// These are package-level vars (not consts) so the cobra root command
// can flip them off at runtime when --no-color is set, and so we
// honour the auto-detected NO_COLOR / non-TTY state from the style
// package on first import.
//
// Why we keep prefix-style strings instead of going full lipgloss:
// the cmd/ files concatenate these into f-strings and printf calls
// dozens of times. Migrating all 80+ sites to lipgloss.Render() at
// once would balloon the diff and risk subtle output regressions.
// Letting the prefix/suffix shape stay while routing through `style`
// gets us NO_COLOR support and per-flag override with one small file
// instead of a sweeping rewrite.
var (
	green  string
	yellow string
	red    string
	blue   string
	dim    string
	bold   string
	reset  string
)

func init() { applyColorState() }

// disableColors zeros all color variables — used when --no-color or
// NO_COLOR or non-TTY stdout is detected.
func disableColors() {
	green, yellow, red, blue, dim, bold, reset = "", "", "", "", "", "", ""
}

// enableColors sets the standard ANSI prefixes/suffix.
func enableColors() {
	green = "\033[32m"
	yellow = "\033[33m"
	red = "\033[31m"
	blue = "\033[34m"
	dim = "\033[2m"
	bold = "\033[1m"
	reset = "\033[0m"
}

// applyColorState reads the current style.IsEnabled() and updates
// the prefix vars accordingly. Called once at init; called again by
// root.go after parsing --no-color so the flag wins over auto-detect.
func applyColorState() {
	if style.IsEnabled() {
		enableColors()
	} else {
		disableColors()
	}
}

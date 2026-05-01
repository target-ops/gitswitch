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

// enableColors sets truecolor (24-bit) ANSI escapes for our palette.
// Hex → r;g;b mappings match internal/style:
//
//	#22c55e (emerald-500) → 34;197;94
//	#ef4444 (red-500)     → 239;68;68
//	#f59e0b (amber-500)   → 245;158;11
//	#3b82f6 (blue-500)    → 59;130;246
//
// Terminals that don't support truecolor downsample to 256 / 16
// colour automatically — the punchy palette degrades gracefully.
// Was previously bare ANSI-16 (\033[32m etc.) which on modern dark
// themes reads as muted dim grey rather than the punchy semantic
// signal we want.
func enableColors() {
	green = "\033[38;2;34;197;94m"
	yellow = "\033[38;2;245;158;11m"
	red = "\033[38;2;239;68;68m"
	blue = "\033[38;2;59;130;246m"
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

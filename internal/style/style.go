// Package style is the gitswitch CLI's single source of truth for
// colours, icons, and text decoration. Every command in cmd/ should
// render output through these helpers — no raw "\033[...]" escapes,
// no fmt.Sprintf-with-color-codes inline.
//
// Why centralised:
//   - One place to flip on --no-color / honour NO_COLOR.
//   - One place to keep contrast accessible (we don't pick colours
//     that fail on a light terminal).
//   - One place to evolve when we add more visual elements (bordered
//     panels, tables, progress markers).
package style

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// Disable strips all colour and decoration from style output —
// equivalent to running with --no-color or NO_COLOR set.
//
// Honours the de facto NO_COLOR convention (https://no-color.org)
// and auto-detects when stdout is not a tty (piped to less, grep,
// redirected to a file).
//
// SetEnabled() / IsEnabled() let the cobra root command flip this
// based on a --no-color flag at runtime.
var enabled = autoDetect()

func autoDetect() bool {
	// Belt + braces: NO_COLOR wins regardless of FORCE_COLOR (privacy
	// trumps marketing).
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// FORCE_COLOR is the de facto override for "I know stdout isn't a
	// tty but I still want colour" (recorders, CI logs, asciinema).
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return false
	}
	return true
}

// SetEnabled forces colour on or off, overriding the auto-detected
// default. Call from cmd/root.go after parsing --no-color.
func SetEnabled(b bool) { enabled = b }

// IsEnabled reports whether style helpers will emit colour.
func IsEnabled() bool { return enabled }

// ----- semantic styles -----------------------------------------------------
//
// Use these helpers in command output. Resist reaching for raw
// lipgloss in cmd/ files — every visual choice should land here so
// the look stays coherent.

var (
	// Foundational colours. We use ANSI 16 for portability — every
	// terminal renders these the same, no truecolor surprises.
	colourSuccess = lipgloss.Color("2")  // green
	colourError   = lipgloss.Color("1")  // red
	colourWarn    = lipgloss.Color("3")  // yellow
	colourAccent  = lipgloss.Color("4")  // blue
	colourMuted   = lipgloss.Color("8")  // bright black (gray)

	successStyle = lipgloss.NewStyle().Foreground(colourSuccess).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(colourError).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(colourWarn).Bold(true)
	accentStyle  = lipgloss.NewStyle().Foreground(colourAccent)
	mutedStyle   = lipgloss.NewStyle().Foreground(colourMuted)
	boldStyle    = lipgloss.NewStyle().Bold(true)
	dimStyle     = lipgloss.NewStyle().Faint(true)
)

// Render functions — call these from cmd/.

func Success(s string) string { return render(successStyle, s) }
func Error(s string) string   { return render(errorStyle, s) }
func Warn(s string) string    { return render(warnStyle, s) }
func Accent(s string) string  { return render(accentStyle, s) }
func Muted(s string) string   { return render(mutedStyle, s) }
func Dim(s string) string     { return render(dimStyle, s) }
func Bold(s string) string    { return render(boldStyle, s) }

// Path renders an OS path with subtle styling. Different from Muted
// so we can later distinguish "file path" from "secondary text" if
// we ever want to.
func Path(s string) string { return render(mutedStyle, s) }

// Code renders an inline command/code fragment (e.g., `gitswitch use`)
// with light accent so it stands out without shouting.
func Code(s string) string { return render(accentStyle, s) }

// ----- icons ---------------------------------------------------------------

const (
	IconOK    = "✓"
	IconBad   = "✗"
	IconWarn  = "•"
	IconArrow = "→"
)

// CheckMark, Cross, Bullet — semantic wrappers so cmd/ stays readable.
func CheckMark() string { return Success(IconOK) }
func Cross() string     { return Error(IconBad) }
func Bullet() string    { return Warn(IconWarn) }

// ----- multi-line helpers --------------------------------------------------

// Heading renders a bold header, intended for the first line of a
// command's output ("✓ active identity: work", "✗ identity drift").
func Heading(icon, text string) string {
	if !enabled {
		return icon + " " + text
	}
	return icon + " " + boldStyle.Render(text)
}

// KV renders a "  label  value" row with the label in muted style
// and the value plain. The label width is padded so a series of KV
// rows align; pass labelWidth = 0 for "use the longest label you've
// rendered so far" (state-free, so callers must compute themselves).
func KV(label, value string, labelWidth int) string {
	pad := labelWidth - len(label)
	if pad < 0 {
		pad = 0
	}
	return fmt.Sprintf("  %s%s   %s",
		Muted(label+":"),
		spaces(pad),
		value,
	)
}

// Hint is a dim "→ next step" line, used to suggest the obvious
// next command after a successful operation.
func Hint(text string) string {
	return Muted("  " + IconArrow + " " + text)
}

// Divider renders a thin horizontal rule of the given width.
func Divider(width int) string {
	if width <= 0 {
		width = 50
	}
	r := make([]byte, width)
	for i := range r {
		r[i] = '-'
	}
	return Muted(string(r))
}

// ----- panels --------------------------------------------------------------

// BoxStyle picks which palette to wrap a panel in.
type BoxStyle int

const (
	BoxNeutral BoxStyle = iota
	BoxSuccess
	BoxError
	BoxWarn
)

// Box renders `content` inside a thin rounded border. In no-color mode
// it returns the content as-is — borders without colour just look
// like noise on plain output.
func Box(content string, kind BoxStyle) string {
	if !enabled {
		return content
	}
	colour := colourMuted
	switch kind {
	case BoxSuccess:
		colour = colourSuccess
	case BoxError:
		colour = colourError
	case BoxWarn:
		colour = colourWarn
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colour).
		Padding(0, 1).
		Render(content)
}

func render(st lipgloss.Style, s string) string {
	if !enabled {
		return s
	}
	return st.Render(s)
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	r := make([]byte, n)
	for i := range r {
		r[i] = ' '
	}
	return string(r)
}

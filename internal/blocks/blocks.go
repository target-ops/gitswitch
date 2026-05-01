// Package blocks does idempotent edits of plain-text config files
// (~/.gitconfig, ~/.ssh/config, …) using sentinel comment markers.
//
// We mark each gitswitch-managed block with a pair of comments:
//
//	# >>> gitswitch:<name>
//	...
//	# <<< gitswitch:<name>
//
// On every Upsert we strip the previous block (if any) and append the
// new one, preserving everything else in the file. This lets us
// re-run `gitswitch use` safely without duplicating includes or
// growing the file each time.
package blocks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Marker is the comment sentinel style — typically "#" for
// gitconfig/sshconfig but extracted so future formats (e.g. ";"-comment
// configs) can reuse the same machinery.
type Marker struct {
	CommentPrefix string // "#" for gitconfig + sshconfig
	OpenWord      string // "gitswitch:" — concatenated with `name`
}

// DefaultMarker is the canonical "# >>> gitswitch:<name>" marker.
var DefaultMarker = Marker{CommentPrefix: "#", OpenWord: "gitswitch:"}

// Upsert reads `path`, removes any existing block named `name`, appends
// the supplied `body` wrapped in fresh sentinel comments, and writes
// the result back atomically with `mode` perms. Creates parent dirs
// and the file itself when missing. The file's other contents are
// preserved exactly (including non-managed blocks and free-form text).
func Upsert(path, name, body string, mode os.FileMode) error {
	return UpsertWithMarker(path, name, body, mode, DefaultMarker)
}

// UpsertWithMarker is the same as Upsert with a custom comment style.
func UpsertWithMarker(path, name, body string, mode os.FileMode, m Marker) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	stripped := stripBlock(string(existing), name, m)
	updated := appendBlock(stripped, name, body, m)

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op if rename succeeded

	if _, err := tmp.WriteString(updated); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// Remove strips the named block from the file (no-op if not present).
// Useful for an eventual `gitswitch use --unbind`.
func Remove(path, name string, mode os.FileMode) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	stripped := stripBlock(string(existing), name, DefaultMarker)
	if stripped == string(existing) {
		return nil
	}
	return os.WriteFile(path, []byte(stripped), mode)
}

func stripBlock(content, name string, m Marker) string {
	pattern := blockPattern(name, m)
	return pattern.ReplaceAllString(content, "")
}

func appendBlock(content, name, body string, m Marker) string {
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if content != "" && !strings.HasSuffix(content, "\n\n") {
		content += "\n"
	}
	open := fmt.Sprintf("%s >>> %s%s", m.CommentPrefix, m.OpenWord, name)
	closeLine := fmt.Sprintf("%s <<< %s%s", m.CommentPrefix, m.OpenWord, name)
	body = strings.TrimRight(body, "\n")
	return content + open + "\n" + body + "\n" + closeLine + "\n"
}

// blockPattern matches a full sentinel-wrapped block, including its
// trailing newline if present. The (?s) flag lets `.` cross newlines.
func blockPattern(name string, m Marker) *regexp.Regexp {
	cp := regexp.QuoteMeta(m.CommentPrefix)
	ow := regexp.QuoteMeta(m.OpenWord + name)
	return regexp.MustCompile(
		`(?s)` +
			`[ \t]*` + cp + `[ \t]*>>>[ \t]*` + ow + `[ \t]*\n` + // open line
			`.*?` + // body, non-greedy
			`[ \t]*` + cp + `[ \t]*<<<[ \t]*` + ow + `[ \t]*\n?`, // close line
	)
}

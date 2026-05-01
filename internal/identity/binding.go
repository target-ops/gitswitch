package identity

import (
	"path/filepath"
	"strings"
)

// FindBindingForDir returns the most-specific binding whose Directory
// is an ancestor of (or equal to) `dir`. Nil when no binding matches.
//
// "Most specific" means longest prefix — so a binding for /work/team-a
// wins over a binding for /work when the cwd is /work/team-a/repo.
// This lets a parent directory hold a default identity while
// subdirectories override.
func (c *Config) FindBindingForDir(dir string) *Binding {
	dir = filepath.Clean(dir)
	var best *Binding
	for i := range c.Bindings {
		b := &c.Bindings[i]
		bdir := filepath.Clean(b.Directory)
		if !isAncestorOrEqual(bdir, dir) {
			continue
		}
		if best == nil || len(filepath.Clean(best.Directory)) < len(bdir) {
			best = b
		}
	}
	return best
}

// isAncestorOrEqual reports whether `parent` is the same path as or
// an ancestor of `child`. We compare cleaned paths and require the
// child to start with parent + a separator (so /work doesn't match
// /workspace).
func isAncestorOrEqual(parent, child string) bool {
	if parent == child {
		return true
	}
	sep := string(filepath.Separator)
	if !strings.HasSuffix(parent, sep) {
		parent += sep
	}
	return strings.HasPrefix(child+sep, parent)
}

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/target-ops/gitswitch/internal/identity"
)

// TestSyncBindingBlocksMultipleDirsOneIdentity is the regression test for
// the bug where two directories bound to the same identity clobbered each
// other's includeIf block. After sync, ~/.gitconfig must hold one block
// per bound directory, the legacy bare-name block must be gone, and other
// identities must be untouched.
func TestSyncBindingBlocksMultipleDirsOneIdentity(t *testing.T) {
	gc := filepath.Join(t.TempDir(), ".gitconfig")
	// A legacy single-block-per-identity entry, as older gitswitch wrote.
	legacy := "# >>> gitswitch:private\n" +
		"[includeIf \"gitdir:/old/dir/\"]\n    path = /x/private.gitconfig\n" +
		"# <<< gitswitch:private\n"
	if err := os.WriteFile(gc, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := &identity.Config{Bindings: []identity.Binding{
		{Directory: "/work/a", Identity: "private"},
		{Directory: "/work/b", Identity: "private"},
		{Directory: "/work/c", Identity: "other"},
	}}
	if err := syncBindingBlocks(cfg, "private", gc); err != nil {
		t.Fatal(err)
	}

	out, _ := os.ReadFile(gc)
	s := string(out)

	if !strings.Contains(s, "gitdir:/work/a/") || !strings.Contains(s, "gitdir:/work/b/") {
		t.Errorf("both private dirs should have includeIf blocks:\n%s", s)
	}
	if n := strings.Count(s, ">>> gitswitch:private:"); n != 2 {
		t.Errorf("expected 2 per-binding blocks, got %d:\n%s", n, s)
	}
	// Legacy bare block (open line "gitswitch:private" with no :hash) gone.
	if strings.Contains(s, ">>> gitswitch:private\n") {
		t.Errorf("legacy bare block should be removed:\n%s", s)
	}
	if strings.Contains(s, "gitdir:/old/dir/") {
		t.Errorf("legacy block contents should be gone:\n%s", s)
	}
	// Unrelated identity must not be added.
	if strings.Contains(s, "gitdir:/work/c/") {
		t.Errorf("should not touch the 'other' identity:\n%s", s)
	}

	// Idempotent: a second sync yields identical content.
	if err := syncBindingBlocks(cfg, "private", gc); err != nil {
		t.Fatal(err)
	}
	out2, _ := os.ReadFile(gc)
	if string(out2) != s {
		t.Errorf("syncBindingBlocks is not idempotent:\nfirst:\n%s\nsecond:\n%s", s, out2)
	}
}

func TestBindingBlockNameDistinctAndStable(t *testing.T) {
	a := bindingBlockName("private", "/work/a")
	b := bindingBlockName("private", "/work/b")
	if a == b {
		t.Errorf("distinct dirs should yield distinct block names, both = %q", a)
	}
	if !strings.HasPrefix(a, "private:") {
		t.Errorf("block name should be prefixed by identity name: %q", a)
	}
	// Stable across equivalent paths (trailing slash is cleaned away).
	if a != bindingBlockName("private", "/work/a/") {
		t.Errorf("block name should be stable across equivalent paths")
	}
}

package blocks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpsertIdempotentAcrossManyBlocks guards the appendBlock fix: adding
// several blocks and then re-adding them must not keep leaking blank lines
// into the file (the bug that made multi-directory bindings grow
// ~/.gitconfig on every `gitswitch use`).
func TestUpsertIdempotentAcrossManyBlocks(t *testing.T) {
	f := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(f, []byte("[user]\n    name = base\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	write := func() {
		for _, n := range []string{"a", "b", "c"} {
			if err := Upsert(f, n, "[x]\n    k = "+n+"\n", 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}

	write()
	first, _ := os.ReadFile(f)
	for i := 0; i < 3; i++ {
		write()
	}
	again, _ := os.ReadFile(f)

	if string(first) != string(again) {
		t.Errorf("Upsert not idempotent.\nfirst:\n%q\nafter re-runs:\n%q", first, again)
	}
	// User content preserved; no run of 3+ newlines (i.e. >1 blank line).
	if !strings.Contains(string(again), "name = base") {
		t.Errorf("user content lost:\n%s", again)
	}
	if strings.Contains(string(again), "\n\n\n") {
		t.Errorf("blank lines leaked:\n%q", again)
	}
}

// TestRemoveOneOfManyKeepsOthers verifies removing a block leaves the
// others intact and doesn't strand blank lines.
func TestRemoveOneOfManyKeepsOthers(t *testing.T) {
	f := filepath.Join(t.TempDir(), "gitconfig")
	_ = os.WriteFile(f, []byte(""), 0o600)
	for _, n := range []string{"a", "b", "c"} {
		if err := Upsert(f, n, "body-"+n+"\n", 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := Remove(f, "b", 0o600); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(f)
	s := string(out)
	if strings.Contains(s, "gitswitch:b") || strings.Contains(s, "body-b") {
		t.Errorf("block b should be gone:\n%s", s)
	}
	if !strings.Contains(s, "body-a") || !strings.Contains(s, "body-c") {
		t.Errorf("blocks a and c should remain:\n%s", s)
	}
	if strings.Contains(s, "\n\n\n") {
		t.Errorf("blank lines leaked after remove:\n%q", s)
	}
}

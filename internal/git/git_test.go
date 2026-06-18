package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEffectiveIdentityRespectsLocalOverGlobal is the regression test for
// the `doctor` bug: a directory binding (modeled here as a local repo
// override, which is what gitswitch's includeIf block effectively produces)
// must win over the global config. GlobalIdentity must still report the
// global values so doctor can show both.
func TestEffectiveIdentityRespectsLocalOverGlobal(t *testing.T) {
	// Isolated config so we never touch the developer's real ~/.gitconfig.
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "global.gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	// Global identity.
	run("", "config", "--global", "user.name", "Global Name")
	run("", "config", "--global", "user.email", "global@example.com")

	// A repo with a local override.
	repo := t.TempDir()
	run(repo, "init", "-q")
	run(repo, "config", "user.name", "Local Name")
	run(repo, "config", "user.email", "local@example.com")

	// Run our functions from inside the repo.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}

	name, email, _ := EffectiveIdentity()
	if name != "Local Name" || email != "local@example.com" {
		t.Errorf("EffectiveIdentity() = %q <%s>, want \"Local Name\" <local@example.com>", name, email)
	}

	gName, gEmail, _ := GlobalIdentity()
	if gName != "Global Name" || gEmail != "global@example.com" {
		t.Errorf("GlobalIdentity() = %q <%s>, want \"Global Name\" <global@example.com>", gName, gEmail)
	}
}

// TestLocalOverridesAndUnset covers the shadow-detection primitives: a
// repo with identity keys in its local .git/config is reported by
// LocalOverrides, and UnsetLocal clears them so nothing is left to
// shadow a gitswitch binding.
func TestLocalOverridesAndUnset(t *testing.T) {
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "global.gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)

	repo := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "local@example.com")
	run("config", "core.sshCommand", "ssh -i /tmp/key")

	got := LocalOverrides(repo)
	want := map[string]bool{"user.email": true, "core.sshCommand": true}
	if len(got) != len(want) {
		t.Fatalf("LocalOverrides() = %v, want keys %v", got, want)
	}
	for _, k := range got {
		if !want[k] {
			t.Errorf("LocalOverrides() returned unexpected key %q", k)
		}
	}

	// A directory with no local overrides reports none.
	clean := t.TempDir()
	runIn(t, clean, "init", "-q")
	if ov := LocalOverrides(clean); len(ov) != 0 {
		t.Errorf("LocalOverrides(clean repo) = %v, want empty", ov)
	}

	// UnsetLocal clears them; a second unset is a no-op (not an error).
	for _, k := range []string{"user.email", "core.sshCommand"} {
		if err := UnsetLocal(repo, k); err != nil {
			t.Fatalf("UnsetLocal(%q): %v", k, err)
		}
		if err := UnsetLocal(repo, k); err != nil {
			t.Errorf("UnsetLocal(%q) second call should be a no-op, got %v", k, err)
		}
	}
	if ov := LocalOverrides(repo); len(ov) != 0 {
		t.Errorf("after UnsetLocal, LocalOverrides() = %v, want empty", ov)
	}
}

func runIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

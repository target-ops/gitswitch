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

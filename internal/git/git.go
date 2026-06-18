// Package git wraps the small subset of `git config` we need.
package git

import (
	"errors"
	"os/exec"
	"strings"
)

// GlobalIdentity returns the global user.name and user.email from git config.
// Either string may be empty if the corresponding config is unset.
func GlobalIdentity() (name, email string, err error) {
	name, _ = configGet("--global", "user.name")
	email, _ = configGet("--global", "user.email")
	return name, email, nil
}

// GlobalSigning returns the global signing key path and signing format
// (e.g. "ssh", "openpgp"). Either string may be empty if unset.
func GlobalSigning() (signingKey, format string) {
	signingKey, _ = configGet("--global", "user.signingkey")
	format, _ = configGet("--global", "gpg.format")
	return signingKey, format
}

// GlobalGet returns a single value from the global git config.
// Empty string + nil error when the key is unset.
func GlobalGet(key string) (string, error) {
	return configGet("--global", key)
}

// configGet runs `git config <args...>` and returns the trimmed stdout.
// Returns an empty string and a wrapped error when git exits non-zero
// (e.g., the key isn't set).
func configGet(args ...string) (string, error) {
	full := append([]string{"config"}, args...)
	out, err := exec.Command("git", full...).Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return "", nil // unset config — treat as empty, not fatal
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// EffectiveEmail returns `git config user.email` from the current
// working directory — meaning the value git would actually write into
// the next commit's author header. Inside a repo this respects local
// repo config and includeIf chains; outside a repo it falls back to
// the global value.
func EffectiveEmail() (string, error) {
	return configGet("user.email")
}

// EffectiveIdentity returns the user.name and user.email git would
// actually use for the next commit in the current working directory.
// Like EffectiveEmail, this respects local repo config and includeIf
// chains (which is how gitswitch binds a directory to an identity);
// outside a repo it falls back to the global values. Either string may
// be empty if the corresponding config is unset.
func EffectiveIdentity() (name, email string, err error) {
	name, _ = configGet("user.name")
	email, _ = configGet("user.email")
	return name, email, nil
}

// IdentityKeys are the git config keys gitswitch manages through its
// per-identity gitconfig (included via includeIf). When any of these is
// also set in a repo's LOCAL .git/config, the local value wins — git's
// precedence is local > global/includeIf — and silently defeats a
// gitswitch binding. These are the keys worth checking for shadowing.
var IdentityKeys = []string{"user.name", "user.email", "user.signingkey", "core.sshCommand"}

// LocalOverrides returns the subset of IdentityKeys that are set in the
// LOCAL .git/config of the repo at dir. Pass "" for the current working
// directory. The result is empty when dir isn't a git repo or nothing
// identity-related is set locally — so a non-empty result means a
// gitswitch binding for that location is being shadowed.
func LocalOverrides(dir string) []string {
	var set []string
	for _, k := range IdentityKeys {
		if localGet(dir, k) != "" {
			set = append(set, k)
		}
	}
	return set
}

// localGet reads a single key from dir's local .git/config only (never
// global/includeIf). Empty string when unset or dir isn't a repo.
func localGet(dir, key string) string {
	args := []string{}
	if dir != "" {
		args = append(args, "-C", dir)
	}
	args = append(args, "config", "--local", "--get", key)
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "" // unset, or not inside a repo — both are "no override"
	}
	return strings.TrimSpace(string(out))
}

// UnsetLocal removes key from dir's local .git/config. Pass "" for the
// current working directory. Treated as success when the key wasn't set
// in the first place — that's what the caller wants.
func UnsetLocal(dir, key string) error {
	args := []string{}
	if dir != "" {
		args = append(args, "-C", dir)
	}
	args = append(args, "config", "--local", "--unset", key)
	err := exec.Command("git", args...).Run()
	if err == nil {
		return nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.ExitCode() == 5 {
		// `git config --unset` exits 5 when the key didn't exist.
		return nil
	}
	return err
}

// SetGlobal writes a value into the user's global git config.
func SetGlobal(key, value string) error {
	return exec.Command("git", "config", "--global", key, value).Run()
}

// UnsetGlobal clears a key from the global git config. Treated as
// success when the key wasn't set in the first place — that's what
// the caller almost always wants.
func UnsetGlobal(key string) error {
	cmd := exec.Command("git", "config", "--global", "--unset", key)
	err := cmd.Run()
	if err == nil {
		return nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.ExitCode() == 5 {
		// `git config --unset` exits 5 when the key didn't exist.
		return nil
	}
	return err
}

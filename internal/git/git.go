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

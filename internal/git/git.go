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

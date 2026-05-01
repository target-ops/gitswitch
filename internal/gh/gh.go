// Package gh shells out to the GitHub CLI for the small set of
// operations gitswitch needs. We deliberately don't link against any
// GitHub API client — `gh` already handles auth, 2FA, scopes, and rate
// limits.
package gh

import (
	"errors"
	"os/exec"
	"strings"
)

// ErrNotInstalled is returned when `gh` is not on PATH. Callers should
// downgrade gh-related checks to a friendly hint rather than a hard fail.
var ErrNotInstalled = errors.New("gh CLI not installed")

// IsInstalled reports whether `gh` is on PATH.
func IsInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// ActiveLogin returns the currently active gh user's login (e.g. "OfirHaim1")
// by calling `gh api user --jq .login`. Empty string + nil error when not
// logged in or `gh` is missing.
func ActiveLogin() (string, error) {
	if !IsInstalled() {
		return "", ErrNotInstalled
	}
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return "", nil // not authenticated → empty, not fatal
	}
	return strings.TrimSpace(string(out)), nil
}

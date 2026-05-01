// Package hookscript generates the small shell wrapper that git
// invokes on every commit. The wrapper does nothing but exec
// `gitswitch guard check` — keeping the real logic in Go where it
// can be tested, type-checked, and changed without re-installing the
// hook on every user's machine.
package hookscript

import "fmt"

// Version identifies the hook script schema. Bump when the script
// content changes in a way that older installs need to be migrated
// (e.g., new flags, new env vars). Used in the marker comment so
// `gitswitch guard install` can recognise its own files.
const Version = 1

// MarkerComment uniquely identifies a gitswitch-installed hook. We
// look for it before overwriting the file — that's how we know the
// existing hook is ours and safe to replace, vs. a hand-written one
// that should be preserved.
const MarkerComment = "# gitswitch-hook-version:"

// PreCommit returns the bytes of the pre-commit hook script.
//
// We ship it as a plain shell wrapper rather than embedding logic
// here, for three reasons:
//   - The hook lives outside the gitswitch binary's reach to be
//     overwritten on upgrade — a thin wrapper means every install
//     of gitswitch picks up the latest check logic automatically.
//   - The wrapper is small enough to audit at a glance.
//   - If gitswitch is ever uninstalled, the wrapper fails fast with
//     a clear "command not found" instead of silently letting bad
//     commits through.
func PreCommit() []byte {
	return []byte(fmt.Sprintf(`#!/usr/bin/env bash
%s %d
#
# This hook is installed by `+"`gitswitch guard install`"+`. It refuses
# git commits where the active identity is wrong for the current
# directory. To override for a single commit, use:
#     git commit --no-verify
# To remove the hook entirely:
#     gitswitch guard uninstall
#
set -e
exec gitswitch guard check "$@"
`, MarkerComment, Version))
}

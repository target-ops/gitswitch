package cmd

import (
	"fmt"
	"hash/fnv"
	"path/filepath"

	"github.com/target-ops/gitswitch/internal/blocks"
	"github.com/target-ops/gitswitch/internal/identity"
)

// bindingBlockName is the sentinel name for the includeIf block of a
// single (identity, directory) binding. Keying per binding — rather than
// per identity — is what lets several directories share one identity
// without their includeIf blocks clobbering each other in ~/.gitconfig.
// Older gitswitch keyed blocks by the bare identity name, so a second
// `use <id> <dir>` overwrote the first directory's block. The directory
// is folded into a short stable hash so the sentinel stays comment- and
// regex-safe and bounded in length.
func bindingBlockName(name, dir string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(filepath.Clean(dir)))
	return fmt.Sprintf("%s:%08x", name, h.Sum32())
}

// syncBindingBlocks rewrites ~/.gitconfig (at path) so there is exactly
// one includeIf block per directory currently bound to `name` in cfg. It
// first removes the legacy single-block-per-identity form (keyed by the
// bare identity name) written by older gitswitch versions, migrating the
// user forward on the next `use`/`rename`. Idempotent: re-running with
// the same cfg yields the same file.
//
// Note: this only adds/refreshes blocks for directories still bound to
// `name`; callers that drop a binding must remove that binding's block
// (bindingBlockName) themselves before calling this.
func syncBindingBlocks(cfg *identity.Config, name, path string) error {
	// First remove the legacy bare-name block and every current
	// per-binding block, THEN re-append them. Stripping all up front
	// (rather than upserting each in place) keeps the result order-stable
	// and idempotent — repeated in-place upserts would otherwise reshuffle
	// blocks and leak blank lines into ~/.gitconfig.
	//
	// Removing "gitswitch:<name>" only matches the legacy block: the regex
	// anchors on the exact name, so it never touches "gitswitch:<name>:<hash>".
	if err := blocks.Remove(path, name, 0o600); err != nil {
		return err
	}
	for _, b := range cfg.Bindings {
		if b.Identity != name {
			continue
		}
		if err := blocks.Remove(path, bindingBlockName(name, b.Directory), 0o600); err != nil {
			return err
		}
	}
	for _, b := range cfg.Bindings {
		if b.Identity != name {
			continue
		}
		body := buildIncludeIfBlock(b.Directory, identity.GitconfigPath(name))
		if err := blocks.Upsert(path, bindingBlockName(name, b.Directory), body, 0o600); err != nil {
			return err
		}
	}
	return nil
}

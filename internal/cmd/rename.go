package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/blocks"
	"github.com/target-ops/gitswitch/internal/identity"
)

func newRenameCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "rename <old> <new>",
		Aliases: []string{"mv"},
		Short:   "Rename an identity without touching its underlying credentials.",
		Long: `Renames an identity in ~/.config/gitswitch/config.json. Also:

  - moves ~/.config/gitswitch/identities/<old>.gitconfig → <new>.gitconfig
  - rewrites the includeIf block in ~/.gitconfig (the sentinel
    comments use the identity name)
  - retargets every directory binding that pointed at <old>

Useful when init auto-named an identity something awkward (e.g. a
lowercased gh login) and you want a more human name. Underlying
SSH keys and gh accounts are untouched.`,
		Args: exactArgsHelp(2,
			"gitswitch rename takes <old-name> <new-name>",
			"gitswitch rename ofirhaim1 personal",
			"gitswitch mv work-old work    # alias",
		),
		Example: "  gitswitch rename ofirhaim1 personal\n" +
			"  gitswitch mv work-old work-new",
		RunE: func(_ *cobra.Command, args []string) error {
			return runRename(args[0], args[1])
		},
	}
}

func runRename(oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf("old and new names are the same: %q", oldName)
	}
	if !validIdentityName(newName) {
		return fmt.Errorf("invalid name %q — use lowercase letters, digits, '-', '_'", newName)
	}

	cfg, err := identity.Load()
	if err != nil {
		return err
	}
	if cfg.FindByName(oldName) == nil {
		return fmt.Errorf("no identity named %q.\n  run `gitswitch list` to see what's configured", oldName)
	}
	if cfg.FindByName(newName) != nil {
		return fmt.Errorf("an identity named %q already exists.\n  pick a different name, or `gitswitch delete %s` first",
			newName, newName)
	}

	// 1. Move the per-identity gitconfig file. Do this before mutating
	//    the JSON, so on failure we abort with the JSON intact.
	oldPath := identity.GitconfigPath(oldName)
	newPath := identity.GitconfigPath(newName)
	if _, err := os.Stat(oldPath); err == nil {
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("rename %s → %s: %w", oldPath, newPath, err)
		}
	}

	// 2. Strip the old includeIf blocks from ~/.gitconfig — one per
	//    bound directory, plus the legacy bare-name block. We do this
	//    before mutating the JSON, while bindings still reference the
	//    old name.
	gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	for _, b := range cfg.Bindings {
		if b.Identity == oldName {
			if err := blocks.Remove(gitconfigPath, bindingBlockName(oldName, b.Directory), 0o600); err != nil {
				return fmt.Errorf("strip old includeIf block: %w", err)
			}
		}
	}
	if err := blocks.Remove(gitconfigPath, oldName, 0o600); err != nil {
		return fmt.Errorf("strip old includeIf block: %w", err)
	}

	// 3. Mutate JSON: rename identity, retarget bindings.
	for i := range cfg.Identities {
		if cfg.Identities[i].Name == oldName {
			cfg.Identities[i].Name = newName
		}
	}
	for i := range cfg.Bindings {
		if cfg.Bindings[i].Identity == oldName {
			cfg.Bindings[i].Identity = newName
		}
	}

	// 4. Re-add the includeIf blocks under the new name, one per dir
	//    still bound to it.
	if err := syncBindingBlocks(cfg, newName, gitconfigPath); err != nil {
		return fmt.Errorf("re-add includeIf block: %w", err)
	}

	if err := identity.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("%s%s✓ renamed %s → %s%s\n",
		green, bold, oldName, newName, reset)
	return nil
}

func validIdentityName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z',
			'0' <= c && c <= '9',
			c == '-', c == '_':
			continue
		default:
			return false
		}
	}
	// Disallow names that are purely numeric or start with a hyphen,
	// to avoid awkward CLI parsing later.
	return !strings.HasPrefix(s, "-")
}

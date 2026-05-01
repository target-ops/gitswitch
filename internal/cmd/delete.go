package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/blocks"
	"github.com/target-ops/gitswitch/internal/identity"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <name>",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove an identity (and its bindings) from gitswitch.",
		Long: `Removes the named identity from ~/.config/gitswitch/config.json
along with any directory bindings that point at it. Strips the
sentinel-marked includeIf block from ~/.gitconfig and deletes the
per-identity gitconfig file at
~/.config/gitswitch/identities/<name>.gitconfig.

Does NOT touch your SSH keys, GPG keys, or gh CLI authentication —
those aren't gitswitch's to delete. Use this when you've added a
duplicate by accident or want to retire an identity from gitswitch
without nuking the underlying credentials.`,
		Args: exactArgsHelp(1,
			"gitswitch delete requires an identity name",
			"gitswitch delete work        (prompts before removing)",
			"gitswitch delete work -y     (skips the confirmation)",
		),
		Example: "  gitswitch delete work\n" +
			"  gitswitch rm personal -y\n" +
			"  gitswitch list             # see what's configured first",
		RunE: func(cmd *cobra.Command, args []string) error {
			yes, _ := cmd.Flags().GetBool("yes")
			return runDelete(args[0], yes)
		},
	}
	cmd.Flags().BoolP("yes", "y", false, "skip the confirmation prompt")
	return cmd
}

func runDelete(name string, assumeYes bool) error {
	cfg, err := identity.Load()
	if err != nil {
		return err
	}

	id := cfg.FindByName(name)
	if id == nil {
		return fmt.Errorf("no identity named %q.\n  run `gitswitch list` to see what's configured", name)
	}

	// Find any bindings that reference this identity — they go too.
	var refBindings []identity.Binding
	for _, b := range cfg.Bindings {
		if b.Identity == name {
			refBindings = append(refBindings, b)
		}
	}

	// Plan, before any prompt.
	fmt.Println(bold + "About to remove:" + reset)
	fmt.Println()
	fmt.Printf("  %sidentity:%s     %s%s%s   %s\n",
		dim, reset, green, name, reset, summarize(*id))
	if len(refBindings) > 0 {
		fmt.Printf("  %sbindings:%s\n", dim, reset)
		for _, b := range refBindings {
			fmt.Printf("    %s%s%s\n", dim, shortPath(b.Directory), reset)
		}
	}
	fmt.Println()
	fmt.Printf("  %swill not touch:%s\n", dim, reset)
	if id.SSHKey != "" {
		fmt.Printf("    %sSSH key %s%s\n", dim, id.SSHKey, reset)
	}
	if id.GHAccount != "" {
		fmt.Printf("    %sgh CLI auth (%s — manage with `gh auth logout`)%s\n",
			dim, id.GHAccount, reset)
	}
	fmt.Println()

	if !assumeYes {
		var ok bool
		err := huh.NewConfirm().
			Title("Remove identity " + name + "?").
			Affirmative("Yes, remove").
			Negative("Cancel").
			Value(&ok).
			Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				fmt.Println(yellow + "aborted." + reset)
				return nil
			}
			return err
		}
		if !ok {
			fmt.Println(yellow + "cancelled — nothing removed." + reset)
			return nil
		}
	}

	// 1. strip the includeIf block from ~/.gitconfig (no-op if not present)
	gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	if err := blocks.Remove(gitconfigPath, name, 0o600); err != nil {
		return fmt.Errorf("strip includeIf block from ~/.gitconfig: %w", err)
	}

	// 2. delete the per-identity gitconfig file
	perID := identity.GitconfigPath(name)
	if err := os.Remove(perID); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", perID, err)
	}

	// 3. drop the identity from the JSON store
	newIDs := make([]identity.Identity, 0, len(cfg.Identities))
	for _, x := range cfg.Identities {
		if x.Name != name {
			newIDs = append(newIDs, x)
		}
	}
	cfg.Identities = newIDs

	// 4. drop bindings pointing at this identity
	newBindings := make([]identity.Binding, 0, len(cfg.Bindings))
	for _, b := range cfg.Bindings {
		if b.Identity != name {
			newBindings = append(newBindings, b)
		}
	}
	cfg.Bindings = newBindings

	if err := identity.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("%s%s✓ removed %s%s\n", green, bold, name, reset)
	return nil
}

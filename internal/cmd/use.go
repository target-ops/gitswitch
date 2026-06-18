package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/blocks"
	"github.com/target-ops/gitswitch/internal/git"
	"github.com/target-ops/gitswitch/internal/identity"
)

func newUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <identity> [<directory>]",
		Short: "Bind an identity to a directory.",
		Long: `Configures git so that any repo inside <directory> automatically
uses <identity>'s name, email, signing key, and SSH key. Once bound,
no manual switching is needed — every "cd" into the directory tree
is the switch.

Concretely, "use" does three things:

  1. Writes ~/.config/gitswitch/identities/<identity>.gitconfig with
     [user], [commit], [gpg], and [core.sshCommand] sections.

  2. Adds an idempotent includeIf block to ~/.gitconfig that points
     at the per-identity file when the cwd is inside <directory>.

  3. Records the binding in ~/.config/gitswitch/config.json so
     "gitswitch why" can explain the link later.

Re-running "use" with the same identity + directory is a no-op (the
sentinel-marked block is replaced in place). Run with --unbind to
remove the binding.`,
		Args: rangeArgsHelp(1, 2,
			"gitswitch use takes <identity> [<directory>]",
			"gitswitch use work ~/work      # bind work identity to ~/work",
			"gitswitch use personal ~/code",
			"gitswitch use work             # refresh per-identity config (no binding)",
		),
		Example: "  gitswitch use work ~/work\n" +
			"  gitswitch use personal ~/code\n" +
			"  gitswitch use work ~/work --unbind   # reverse the binding",
		RunE: func(cmd *cobra.Command, args []string) error {
			unbind, _ := cmd.Flags().GetBool("unbind")
			return runUse(args, unbind)
		},
	}
	cmd.Flags().Bool("unbind", false, "remove the binding instead of creating it")
	return cmd
}

func runUse(args []string, unbind bool) error {
	name := args[0]
	cfg, err := identity.Load()
	if err != nil {
		return err
	}
	id := cfg.FindByName(name)
	if id == nil {
		return fmt.Errorf("no identity named %q. Run `gitswitch list` to see what's configured, or `gitswitch init` to detect identities", name)
	}

	// 1. Always (re)write the per-identity gitconfig — easier to keep
	// it consistent with the stored Identity than to diff.
	if err := identity.WriteGitconfig(*id); err != nil {
		return fmt.Errorf("write per-identity gitconfig: %w", err)
	}

	if len(args) == 1 {
		// Just refresh the per-identity file, no directory binding.
		fmt.Printf("%s✓%s refreshed %s\n",
			green, reset, identity.GitconfigPath(name))
		fmt.Printf("%s  (no directory specified — pass <directory> to bind)%s\n", dim, reset)
		return nil
	}

	dir, err := resolveDir(args[1])
	if err != nil {
		return err
	}

	if unbind {
		if err := removeBinding(cfg, name, dir); err != nil {
			return err
		}
		if err := identity.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("%s✓%s unbound %s from %s\n", green, reset, name, dir)
		return nil
	}

	// 2. Update ~/.gitconfig with an idempotent includeIf block.
	gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	blockBody := buildIncludeIfBlock(dir, identity.GitconfigPath(name))
	if err := blocks.Upsert(gitconfigPath, name, blockBody, 0o600); err != nil {
		return fmt.Errorf("update ~/.gitconfig: %w", err)
	}

	// 3. Record the binding.
	addBinding(cfg, name, dir)
	if err := identity.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("%s%s✓ bound %s%s%s%s → %s%s%s\n",
		green, bold, reset, bold, name, reset, dim, dir, reset)
	fmt.Println()
	fmt.Printf("  %sper-identity config:%s %s\n", dim, reset, identity.GitconfigPath(name))
	fmt.Printf("  %sincludeIf in:%s        %s\n", dim, reset, gitconfigPath)

	// 4. If the bound directory is itself a repo with identity keys set
	// in its local .git/config, those win over the includeIf we just
	// wrote (git precedence: local > global) and silently defeat the
	// binding. Offer to clear them so the binding actually takes effect.
	authoritative, err := clearShadowingOverrides(dir)
	if err != nil {
		return err
	}

	fmt.Println()
	if authoritative {
		fmt.Printf("%scd into %s and your git identity becomes %s automatically.%s\n",
			dim, dir, name, reset)
	} else {
		fmt.Printf("%sclear the local overrides above, then cd into %s switches to %s automatically.%s\n",
			dim, dir, name, reset)
	}
	fmt.Println()
	fmt.Printf("Verify it worked:  %scd %s && gitswitch doctor%s\n", dim, dir, reset)
	return nil
}

// clearShadowingOverrides checks dir's local .git/config for identity
// keys that would shadow the binding we just wrote, and — with the
// user's confirmation — removes them. Returns whether the binding is now
// authoritative: true when there were no overrides or they were removed,
// false when conflicting local config remains.
func clearShadowingOverrides(dir string) (bool, error) {
	overrides := git.LocalOverrides(dir)
	if len(overrides) == 0 {
		return true, nil
	}

	fmt.Println()
	fmt.Printf("  %s⚠ local .git/config overrides this binding:%s\n", yellow, reset)
	fmt.Printf("      %s%s%s\n", bold, strings.Join(overrides, ", "), reset)
	fmt.Printf("  %sthese win over the includeIf (git precedence: local > global),%s\n", dim, reset)
	fmt.Printf("  %sso the binding won't take effect until they're removed.%s\n", dim, reset)
	fmt.Println()

	var remove bool
	if err := huh.NewConfirm().
		Title("Remove these local overrides so the binding takes effect?").
		Affirmative("Yes, remove").
		Negative("No, keep them").
		Value(&remove).
		Run(); err != nil {
		// User aborted, or there's no TTY to prompt on (scripts/CI). The
		// binding is already written — don't fail the command; leave the
		// overrides and fall through to the manual instructions below.
		remove = false
	}

	if !remove {
		fmt.Printf("  %s• kept local overrides — they still shadow this binding. Clear them with:%s\n", dim, reset)
		for _, k := range overrides {
			fmt.Printf("      %sgit -C %s config --local --unset %s%s\n", dim, dir, k, reset)
		}
		return false, nil
	}
	for _, k := range overrides {
		if err := git.UnsetLocal(dir, k); err != nil {
			return false, fmt.Errorf("unset local %s: %w", k, err)
		}
	}
	fmt.Printf("  %s✓ removed local overrides — binding is now authoritative%s\n", green, reset)
	return true, nil
}

// buildIncludeIfBlock renders the body of the sentinel-wrapped block
// we drop into ~/.gitconfig. Trailing slash on gitdir is mandatory —
// without it the pattern matches the literal directory name only and
// nothing inside it. This is the single most common gotcha in every
// "managing multiple Git identities" tutorial.
func buildIncludeIfBlock(dir, includePath string) string {
	dir = strings.TrimRight(dir, "/") + "/"
	return fmt.Sprintf("[includeIf \"gitdir:%s\"]\n    path = %s\n",
		dir, includePath)
}

// resolveDir expands ~/ and ./, validates the directory exists, and
// returns an absolute path. Pre-flight check so users don't bind to a
// path with a typo and only discover the breakage at commit time.
func resolveDir(in string) (string, error) {
	if strings.HasPrefix(in, "~/") {
		in = filepath.Join(os.Getenv("HOME"), in[2:])
	}
	abs, err := filepath.Abs(in)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("directory does not exist: %s", abs)
		}
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", abs)
	}
	return abs, nil
}

func addBinding(c *identity.Config, name, dir string) {
	for i, b := range c.Bindings {
		if b.Directory == dir {
			c.Bindings[i].Identity = name
			return
		}
	}
	c.Bindings = append(c.Bindings, identity.Binding{Directory: dir, Identity: name})
}

func removeBinding(c *identity.Config, name, dir string) error {
	for i, b := range c.Bindings {
		if b.Directory == dir && b.Identity == name {
			c.Bindings = append(c.Bindings[:i], c.Bindings[i+1:]...)
			gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
			return blocks.Remove(gitconfigPath, name, 0o600)
		}
	}
	return fmt.Errorf("%s is not bound to %s", name, dir)
}

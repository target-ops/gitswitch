package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/git"
	"github.com/target-ops/gitswitch/internal/hookscript"
	"github.com/target-ops/gitswitch/internal/identity"
)

// HooksDir is the directory we install our hook into. We use our own
// directory (rather than ~/.git-hooks or the user's existing one) so
// `guard install` can be reasoned about and reversed cleanly.
func HooksDir() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "gitswitch", "hooks")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "gitswitch", "hooks")
}

func preCommitPath() string {
	return filepath.Join(HooksDir(), "pre-commit")
}

func newGuardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Install the pre-commit hook that refuses wrong-author commits.",
		Long: `Pre-commit hook that refuses commits when the active identity
doesn't match the directory you're committing from. The killer
feature: makes "I committed as the wrong person" structurally
impossible while gitswitch is in charge.

  gitswitch guard install     install the hook globally
  gitswitch guard uninstall   remove the hook
  gitswitch guard status      show whether the hook is active
  gitswitch guard check       (called by the hook itself; not for humans)
`,
	}
	cmd.AddCommand(newGuardInstallCmd())
	cmd.AddCommand(newGuardUninstallCmd())
	cmd.AddCommand(newGuardStatusCmd())
	cmd.AddCommand(newGuardCheckCmd())
	return cmd
}

func newGuardInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the global pre-commit hook.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			force, _ := cmd.Flags().GetBool("force")
			return runGuardInstall(force)
		},
	}
	cmd.Flags().Bool("force", false, "overwrite an existing core.hooksPath set by another tool")
	return cmd
}

func newGuardUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the pre-commit hook and unset core.hooksPath if we set it.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGuardUninstall()
		},
	}
}

func newGuardStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether the guard hook is active.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGuardStatus()
		},
	}
}

func newGuardCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check",
		Short:  "Run the guard check (called by the installed hook, not humans).",
		Hidden: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return runGuardCheck(args)
		},
	}
}

func runGuardInstall(force bool) error {
	if err := os.MkdirAll(HooksDir(), 0o700); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	// Refuse to overwrite a non-gitswitch hook unless --force is set.
	hookPath := preCommitPath()
	if existing, err := os.ReadFile(hookPath); err == nil {
		if !bytes.Contains(existing, []byte(hookscript.MarkerComment)) && !force {
			return fmt.Errorf(
				"a non-gitswitch pre-commit hook already exists at %s.\n"+
					"  re-run with --force to overwrite, or remove that file by hand first",
				hookPath,
			)
		}
	}

	if err := os.WriteFile(hookPath, hookscript.PreCommit(), 0o755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	// Set git's global core.hooksPath if we can do so cleanly.
	current, _ := git.GlobalGet("core.hooksPath")
	hooksDir := HooksDir()
	switch {
	case current == "":
		if err := git.SetGlobal("core.hooksPath", hooksDir); err != nil {
			return fmt.Errorf("set core.hooksPath: %w", err)
		}
	case sameDir(current, hooksDir):
		// Already pointed at us.
	case force:
		if err := git.SetGlobal("core.hooksPath", hooksDir); err != nil {
			return fmt.Errorf("set core.hooksPath: %w", err)
		}
	default:
		return fmt.Errorf(
			"git config --global core.hooksPath is already set to %s.\n"+
				"  re-run with --force to switch to gitswitch's hooks dir\n"+
				"  (note: your existing hooks at that path will stop firing)\n"+
				"  or copy %s into %s to keep both",
			current, hookPath, current,
		)
	}

	fmt.Printf("%s%s✓ guard installed%s\n", green, bold, reset)
	fmt.Println()
	fmt.Printf("  %shook:%s         %s\n", dim, reset, hookPath)
	fmt.Printf("  %score.hooksPath:%s %s\n", dim, reset, hooksDir)
	fmt.Println()
	fmt.Println(dim + "Every `git commit` now checks identity vs directory binding." + reset)
	fmt.Println(dim + "Override once: " + reset + "git commit --no-verify")
	fmt.Println(dim + "Remove later:  " + reset + "gitswitch guard uninstall")
	return nil
}

func runGuardUninstall() error {
	hookPath := preCommitPath()
	if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove hook: %w", err)
	}

	current, _ := git.GlobalGet("core.hooksPath")
	if sameDir(current, HooksDir()) {
		if err := git.UnsetGlobal("core.hooksPath"); err != nil {
			return fmt.Errorf("unset core.hooksPath: %w", err)
		}
	}

	fmt.Printf("%s✓ guard removed%s\n", green, reset)
	if current != "" && !sameDir(current, HooksDir()) {
		fmt.Printf("%s  (your core.hooksPath at %s was left in place)%s\n",
			dim, current, reset)
	}
	return nil
}

func runGuardStatus() error {
	hookPath := preCommitPath()
	hookOK := false
	if data, err := os.ReadFile(hookPath); err == nil {
		hookOK = bytes.Contains(data, []byte(hookscript.MarkerComment))
	}
	hooksPath, _ := git.GlobalGet("core.hooksPath")
	pathOK := sameDir(hooksPath, HooksDir())

	row("hook script", hookPath, hookOK)
	if hooksPath == "" {
		row("core.hooksPath", "(unset)", false)
	} else {
		row("core.hooksPath", hooksPath, pathOK)
	}

	fmt.Println()
	switch {
	case hookOK && pathOK:
		fmt.Println(green + bold + "✓ guard is active" + reset)
	case hookOK && !pathOK:
		fmt.Println(yellow + bold + "• hook installed but not wired into git" + reset)
		fmt.Println(dim + "  fix: gitswitch guard install --force" + reset)
	case !hookOK && pathOK:
		fmt.Println(yellow + bold + "• core.hooksPath set but hook script missing" + reset)
		fmt.Println(dim + "  fix: gitswitch guard install" + reset)
	default:
		fmt.Println(yellow + bold + "• guard is not installed" + reset)
		fmt.Println(dim + "  fix: gitswitch guard install" + reset)
	}
	return nil
}

// runGuardCheck is the workhorse fired on every commit. Stay fast and
// fail safely: any internal error exits 0 (don't block legitimate
// commits when gitswitch itself is broken).
func runGuardCheck(_ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	cfg, err := identity.Load()
	if err != nil || cfg == nil {
		return nil
	}

	binding := cfg.FindBindingForDir(cwd)
	if binding == nil {
		// Not a directory we manage — don't interfere.
		return nil
	}
	id := cfg.FindByName(binding.Identity)
	if id == nil {
		// Orphan binding (identity was removed); nothing useful to enforce.
		return nil
	}

	effective, _ := git.EffectiveEmail()
	if strings.EqualFold(effective, id.Email) {
		return nil
	}

	// Mismatch. This is the disaster-prevention moment.
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%s%s✗ gitswitch guard: blocked commit%s\n",
		red, bold, reset)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "  %sin directory:%s   %s\n", dim, reset, cwd)
	fmt.Fprintf(os.Stderr, "  %sexpected:%s       %s   (bound identity: %s)\n",
		dim, reset, id.Email, id.Name)
	fmt.Fprintf(os.Stderr, "  %sgot:%s            %s\n", dim, reset, effective)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "  %sfix:%s            gitswitch use %s %s\n",
		dim, reset, id.Name, binding.Directory)
	fmt.Fprintf(os.Stderr, "                  (or: git commit --no-verify to override this once)\n")
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
	return nil
}

// sameDir compares two paths semantically — resolves both through
// filepath.Clean. Not symlink-aware; that's intentional, the user
// shouldn't be expected to keep two semantically identical hook dirs
// hooked up via symlink.
func sameDir(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

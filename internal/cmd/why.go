package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/git"
	"github.com/target-ops/gitswitch/internal/identity"
)

func newWhyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "why",
		Short: "Explain the active identity for the current directory.",
		Long: `Explains, in plain English, why your git identity is what it is
right now: which binding matched, which per-identity gitconfig was
included, what user.email resolves to, and whether the layers agree.

The honest counterweight to any "automatic" tool — magic you can't
inspect is just a bug waiting to happen.`,
		Example: "  cd ~/work && gitswitch why\n" +
			"  cd ~/personal && gitswitch why",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runWhy()
		},
	}
}

func runWhy() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := identity.Load()
	if err != nil {
		return err
	}

	binding := cfg.FindBindingForDir(cwd)
	effective, _ := git.EffectiveEmail()

	// Case 1: no binding for this dir — global identity wins.
	if binding == nil {
		fmt.Println(bold + "no gitswitch binding active for this directory" + reset)
		fmt.Println()
		fmt.Printf("  %scwd:%s              %s\n", dim, reset, cwd)
		gitName, gitEmail, _ := git.GlobalIdentity()
		fmt.Printf("  %sglobal identity:%s  %s\n", dim, reset,
			renderIdentity(gitName, gitEmail))
		fmt.Println()
		fmt.Println(dim + "Bind this directory to an identity:" + reset)
		fmt.Printf("  %sgitswitch use <name> %s%s\n", dim, cwd, reset)
		if len(cfg.Identities) == 0 {
			fmt.Println()
			fmt.Println(dim + "(no identities configured yet — run `gitswitch init` first)" + reset)
		}
		return nil
	}

	id := cfg.FindByName(binding.Identity)
	if id == nil {
		// Orphan binding — the identity it pointed at was removed.
		fmt.Println(yellow + bold + "• binding points at a missing identity" + reset)
		fmt.Println()
		fmt.Printf("  %sbound directory:%s  %s\n", dim, reset, binding.Directory)
		fmt.Printf("  %sbound identity:%s   %s   %s(removed?)%s\n",
			dim, reset, binding.Identity, dim, reset)
		fmt.Println()
		fmt.Println(dim + "Re-bind or remove the stale binding:" + reset)
		fmt.Printf("  %sgitswitch use <name> %s%s\n", dim, binding.Directory, reset)
		fmt.Printf("  %sgitswitch use %s %s --unbind%s\n",
			dim, binding.Identity, binding.Directory, reset)
		return nil
	}

	// Case 2 & 3 — render the resolved chain.
	matches := strings.EqualFold(effective, id.Email)
	// A repo's local .git/config wins over the includeIf (git precedence:
	// local > global). When identity keys are set locally, the binding is
	// being shadowed — even if the effective email happens to match, it's
	// a coincidence, not the binding doing its job.
	overrides := git.LocalOverrides("")
	shadowed := len(overrides) > 0
	switch {
	case !matches:
		fmt.Println(red + bold + "✗ identity drift" + reset)
	case shadowed:
		fmt.Println(yellow + bold + "⚠ binding shadowed by local .git/config" + reset)
	default:
		fmt.Println(green + bold + "✓ active identity: " + id.Name + reset)
	}
	fmt.Println()
	fmt.Printf("  %scwd:%s               %s\n", dim, reset, cwd)
	fmt.Printf("  %sbound directory:%s   %s\n", dim, reset, binding.Directory)

	switch {
	case !matches:
		fmt.Printf("  %suser.email:%s        %s   %s\n",
			dim, reset, effective, red+"← MISMATCH"+reset)
		fmt.Printf("  %sexpected:%s          %s\n", dim, reset, id.Email)
	case shadowed:
		fmt.Printf("  %suser.email:%s        %s   %s\n",
			dim, reset, effective, yellow+"(from local .git/config, not the binding)"+reset)
	default:
		fmt.Printf("  %suser.email:%s        %s   %s\n",
			dim, reset, effective, green+"(matches binding ✓)"+reset)
	}

	if id.GitName != "" {
		fmt.Printf("  %suser.name:%s         %s\n", dim, reset, id.GitName)
	}
	if id.SSHKey != "" {
		fmt.Printf("  %sssh key:%s           %s\n", dim, reset, shortPath(id.SSHKey))
	}
	if id.SigningKey != "" {
		fmt.Printf("  %ssigning key:%s       %s\n", dim, reset, shortPath(id.SigningKey))
	}
	if id.GHAccount != "" {
		fmt.Printf("  %sgh account:%s        %s\n", dim, reset, id.GHAccount)
	}

	fmt.Println()
	fmt.Printf("  %sresolved by:%s       %s\n", dim, reset,
		fmt.Sprintf(`includeIf "gitdir:%s/" in ~/.gitconfig`,
			strings.TrimRight(binding.Directory, "/")))
	fmt.Printf("  %sper-identity file:%s %s\n", dim, reset,
		shortPath(identity.GitconfigPath(id.Name)))

	if shadowed {
		fmt.Println()
		fmt.Println(yellow + bold + "⚠ this repo's local .git/config overrides the binding" + reset)
		fmt.Printf("  %sset locally, so they win over the includeIf above:%s\n", dim, reset)
		fmt.Printf("      %s%s%s\n", bold, strings.Join(overrides, ", "), reset)
		fmt.Println()
		fmt.Printf("  %sgitswitch isn't really in control here. Take over with:%s\n", dim, reset)
		fmt.Printf("  %sgitswitch use %s %s%s\n", dim, id.Name, binding.Directory, reset)
		fmt.Printf("  %s(it will offer to remove the local overrides)%s\n", dim, reset)
	}

	if !matches {
		fmt.Println()
		fmt.Println(red + bold + "fix:" + reset)
		fmt.Printf("  %sgitswitch use %s %s%s\n",
			dim, id.Name, binding.Directory, reset)
		fmt.Printf("  %s(re-writes the includeIf block; idempotent)%s\n", dim, reset)
		os.Exit(1)
	}

	return nil
}

// renderIdentity formats "name <email>" tolerating either being empty.
func renderIdentity(name, email string) string {
	switch {
	case name != "" && email != "":
		return fmt.Sprintf("%s <%s>", name, email)
	case email != "":
		return email
	case name != "":
		return name
	default:
		return "(unset)"
	}
}

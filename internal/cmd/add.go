package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/identity"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <name>",
		Aliases: []string{"new"},
		Short:   "Manually add an identity (when init didn't catch it).",
		Long: `Adds a new identity to ~/.config/gitswitch/config.json. For
identities init can already detect (existing git config + ssh key +
gh login on this machine), prefer running "gitswitch init". Use
"add" when you want to register an identity whose key/account isn't
on this machine yet, or when init merged two identities you wanted
kept separate.

Required: <name> + --email. Everything else is optional.

Run with no flags (just the name) and you'll get an interactive
prompt for each field.`,
		Args: exactArgsHelp(1,
			"gitswitch add requires an identity name",
			"gitswitch add work --email you@company.com",
			"gitswitch add personal              (interactive prompts for the rest)",
		),
		Example: "  gitswitch add work --email you@company.com --gh you-work\n" +
			"  gitswitch add personal     # interactive: prompts for email, key, gh, etc.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0])
		},
	}
	cmd.Flags().String("email", "", "git user.email for this identity (required)")
	cmd.Flags().String("git-name", "", "git user.name for this identity (defaults to <name>)")
	cmd.Flags().String("ssh-key", "", "path to the SSH private key (e.g. ~/.ssh/id_ed25519_work)")
	cmd.Flags().String("signing-key", "", "path to the SSH public key for commit signing (e.g. ~/.ssh/id_ed25519_work.pub)")
	cmd.Flags().String("gh", "", "gh CLI account login (e.g. octocat)")
	cmd.Flags().String("vendor", "github", "vendor — github / gitlab / bitbucket")
	cmd.Flags().BoolP("yes", "y", false, "skip the interactive prompts; require flags")
	return cmd
}

func runAdd(cmd *cobra.Command, name string) error {
	if !validIdentityName(name) {
		return fmt.Errorf("invalid name %q — use lowercase letters, digits, '-', '_'", name)
	}

	cfg, err := identity.Load()
	if err != nil {
		return err
	}
	if cfg.FindByName(name) != nil {
		return fmt.Errorf("an identity named %q already exists.\n  use `gitswitch rename` or `gitswitch delete` first", name)
	}

	email, _ := cmd.Flags().GetString("email")
	gitName, _ := cmd.Flags().GetString("git-name")
	sshKey, _ := cmd.Flags().GetString("ssh-key")
	signingKey, _ := cmd.Flags().GetString("signing-key")
	ghAccount, _ := cmd.Flags().GetString("gh")
	vendor, _ := cmd.Flags().GetString("vendor")
	assumeYes, _ := cmd.Flags().GetBool("yes")

	// Interactive prompts for missing fields when not in --yes mode.
	if !assumeYes {
		if email == "" {
			if err := huh.NewInput().
				Title("Email for " + name).
				Description("Goes into git config user.email").
				Value(&email).
				Validate(validateEmail).
				Run(); err != nil {
				return abortOnUserCancel(err)
			}
		}
		if gitName == "" {
			gitName = name
			if err := huh.NewInput().
				Title("Display name").
				Description("Goes into git config user.name (default: identity name)").
				Value(&gitName).
				Run(); err != nil {
				return abortOnUserCancel(err)
			}
		}
		if sshKey == "" {
			if err := huh.NewInput().
				Title("SSH private key path").
				Description("Optional. Leave empty if you'll add it later.").
				Value(&sshKey).
				Run(); err != nil {
				return abortOnUserCancel(err)
			}
		}
		if vendor == "github" && ghAccount == "" {
			if err := huh.NewInput().
				Title("gh CLI account login").
				Description("Optional. Used by `gitswitch switch` to flip `gh auth`.").
				Value(&ghAccount).
				Run(); err != nil {
				return abortOnUserCancel(err)
			}
		}
	}

	if email == "" {
		return errors.New("--email is required (or run interactively without --yes)")
	}
	if err := validateEmail(email); err != nil {
		return err
	}

	// Derive signing key from SSH key when both are unset and the
	// .pub file is sitting next to the private key.
	if signingKey == "" && sshKey != "" {
		pub := expandHome(sshKey) + ".pub"
		if _, err := os.Stat(pub); err == nil {
			signingKey = sshKey + ".pub"
		}
	}

	id := identity.Identity{
		Name:       name,
		Email:      email,
		GitName:    coalesce(gitName, name),
		SSHKey:     expandHome(sshKey),
		SigningKey: expandHome(signingKey),
		GHAccount:  ghAccount,
		Vendor:     vendor,
	}
	cfg.Upsert(id)
	if err := identity.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("%s%s✓ added %s%s   %s\n",
		green, bold, name, reset, summarize(id))
	fmt.Println()
	fmt.Printf("%sNext: bind it to a directory.%s\n", dim, reset)
	fmt.Printf("  %sgitswitch use %s ~/some/dir%s\n", dim, name, reset)
	return nil
}

func validateEmail(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
		return errors.New("doesn't look like an email address")
	}
	return nil
}

func abortOnUserCancel(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		fmt.Println(yellow + "aborted." + reset)
		return nil
	}
	return err
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// expandHome turns "~/foo" into "$HOME/foo". Path values that aren't
// home-relative (or are empty) come back unchanged.
func expandHome(p string) string {
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "~/") {
		return strings.Replace(p, "~", os.Getenv("HOME"), 1)
	}
	return p
}

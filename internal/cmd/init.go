package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/discover"
	"github.com/target-ops/gitswitch/internal/identity"
	"github.com/target-ops/gitswitch/internal/legacy"
)

func newInitCommand() *cobra.Command {
	var assumeYes bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Auto-detect identities on this machine and set up gitswitch.",
		Long: `Reads your existing git config, ~/.ssh/config, public keys, and gh
auth state. Surfaces every identity it can find. Asks you to name each
one. Writes the gitswitch config — but does not yet apply anything to
~/.gitconfig or ~/.ssh/config; that happens when you run "gitswitch use".

This command is read-only with respect to the rest of your system. The
only file it writes is ~/.config/gitswitch/config.json.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInit(assumeYes)
		},
	}
	cmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "skip confirmation prompts; accept proposed names")
	return cmd
}

func runInit(assumeYes bool) error {
	// 0. surface legacy v0.2.x config if present — many users coming
	//    from brew will have one and miss it otherwise.
	if legacy.Exists() {
		fmt.Println(yellow + "Found a v0.2.x config at " + legacy.LegacyPath() + reset)
		fmt.Println(dim + "  to import it instead of (or in addition to) auto-detection:" + reset)
		fmt.Println(dim + "    gitswitch migrate" + reset)
		fmt.Println()
	}

	// 1. discover
	found := discover.Scan()
	if len(found) == 0 {
		fmt.Println(yellow + "No identities detected." + reset)
		fmt.Println("  • Set git globals first: " + dim + "git config --global user.email you@example.com" + reset)
		fmt.Println("  • Or generate a key:    " + dim + "ssh-keygen -t ed25519 -C \"you@example.com\"" + reset)
		fmt.Println("  • Then re-run:          " + dim + "gitswitch init" + reset)
		return nil
	}

	// 2. show what was found, before any prompts
	fmt.Println(bold + "Found " + countWord(len(found), "identity", "identities") + " on this machine:" + reset)
	fmt.Println()
	for i, d := range found {
		printDetected(i, d)
	}
	fmt.Println()

	// 3. load existing config; we'll merge into it
	cfg, err := identity.Load()
	if err != nil {
		return fmt.Errorf("load existing config: %w", err)
	}

	// 4. for each detected identity, propose a name and (optionally) ask the user
	var proposed []identity.Identity
	usedNames := map[string]bool{}
	for _, ex := range cfg.Identities {
		usedNames[ex.Name] = true
	}

	for _, d := range found {
		name := proposeName(d, usedNames)
		if !assumeYes {
			if err := promptName(&name, d); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					fmt.Println(yellow + "aborted." + reset)
					return nil
				}
				return err
			}
		}
		usedNames[name] = true
		proposed = append(proposed, identity.Identity{
			Name:       name,
			Email:      d.Email,
			GitName:    d.GitName,
			SSHKey:     d.SSHKey,
			SigningKey: d.SigningKey,
			GHAccount:  d.GHAccount,
			Vendor:     d.Vendor,
		})
	}

	// 5. final confirmation
	if !assumeYes {
		fmt.Println()
		fmt.Println(bold + "About to write:" + reset)
		for _, p := range proposed {
			fmt.Printf("  %s%s%s  %s\n", green, p.Name, reset, summarize(p))
		}
		fmt.Println()

		var ok bool
		err := huh.NewConfirm().
			Title("Write " + countWord(len(proposed), "identity", "identities") + " to " + identity.Path() + "?").
			Affirmative("Yes, write").
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
			fmt.Println(yellow + "cancelled — nothing written." + reset)
			return nil
		}
	}

	// 6. merge and save
	for _, p := range proposed {
		cfg.Upsert(p)
	}
	if err := identity.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Println()
	fmt.Println(green + bold + "✓ wrote " + identity.Path() + reset)
	fmt.Println()
	fmt.Println(dim + "Next: bind an identity to a directory." + reset)
	for _, p := range proposed {
		fmt.Printf("  %sgitswitch use %s ~/some/dir%s\n", dim, p.Name, reset)
	}
	return nil
}

// printDetected renders one Detected row for the pre-prompt summary.
func printDetected(i int, d discover.Detected) {
	header := fmt.Sprintf("  %d. %s", i+1, summarizeDetected(d))
	fmt.Println(bold + header + reset)
	if len(d.Sources) > 0 {
		fmt.Printf("     %sfound in: %s%s\n", dim, strings.Join(d.Sources, ", "), reset)
	}
}

// summarizeDetected renders one identity in a compact "email · gh · vendor" form.
func summarizeDetected(d discover.Detected) string {
	parts := []string{}
	if d.Email != "" {
		parts = append(parts, d.Email)
	}
	if d.GHAccount != "" {
		parts = append(parts, "gh: "+d.GHAccount)
	}
	if d.SSHKey != "" {
		parts = append(parts, "key: "+shortPath(d.SSHKey))
	}
	if len(parts) == 0 {
		return "(unidentified)"
	}
	return strings.Join(parts, " · ")
}

// summarize is the same idea but for a saved Identity.
func summarize(id identity.Identity) string {
	parts := []string{id.Email}
	if id.GHAccount != "" {
		parts = append(parts, "gh: "+id.GHAccount)
	}
	if id.SSHKey != "" {
		parts = append(parts, "key: "+shortPath(id.SSHKey))
	}
	return strings.Join(parts, " · ")
}

// genericVendorNames are too vague to make good identity names. We
// fall through them when proposing.
var genericVendorNames = map[string]bool{
	"github": true, "gitlab": true, "bitbucket": true,
}

// proposeName picks a sensible default name for a detected identity,
// avoiding collisions with names already in usedNames. Prefers gh
// account and email local-part over generic vendor names.
func proposeName(d discover.Detected, usedNames map[string]bool) string {
	candidates := []string{
		d.GHAccount,
		localPart(d.Email),
	}
	for _, c := range candidates {
		c = sanitizeName(c)
		if c == "" || genericVendorNames[c] {
			continue
		}
		if !usedNames[c] {
			return c
		}
	}
	// Last resort: append a numeric suffix.
	base := sanitizeName(localPart(d.Email))
	if base == "" || genericVendorNames[base] {
		base = "identity"
	}
	for i := 2; i < 100; i++ {
		n := fmt.Sprintf("%s%d", base, i)
		if !usedNames[n] {
			return n
		}
	}
	return base
}

func promptName(name *string, d discover.Detected) error {
	return huh.NewInput().
		Title("Name this identity").
		Description(summarizeDetected(d)).
		Value(name).
		Validate(func(s string) error {
			s = strings.TrimSpace(s)
			if s == "" {
				return errors.New("name cannot be empty")
			}
			if strings.ContainsAny(s, " \t/\\") {
				return errors.New("no spaces or slashes")
			}
			return nil
		}).
		Run()
}

// sanitizeName lowercases and strips characters that don't belong in an
// identity name (we use it as a directory-ish key in config).
func sanitizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z', '0' <= c && c <= '9', c == '-', c == '_':
			out = append(out, c)
		}
	}
	return string(out)
}

func localPart(email string) string {
	if i := strings.IndexByte(email, '@'); i > 0 {
		return email[:i]
	}
	return email
}

func shortPath(p string) string {
	home := homeDir()
	if home != "" && strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}

func homeDir() string { return getenv("HOME") }
func getenv(k string) string {
	// Indirection so tests can stub later without depending on `os`.
	return osEnv(k)
}

func countWord(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

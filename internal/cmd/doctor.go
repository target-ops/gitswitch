package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/gh"
	"github.com/target-ops/gitswitch/internal/git"
	"github.com/target-ops/gitswitch/internal/identity"
	"github.com/target-ops/gitswitch/internal/ssh"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Verify your git/ssh/gh identities all agree.",
		Long: `Reports who you are, end-to-end, across the layers that have
opinions about it: git config, ssh, and the GitHub CLI. Prints one
green line if everything agrees, or a red diff if they disagree.

This is the command you should run any time you suspect "wait, did
my last commit go out as the right person?"`,
		Example: "  gitswitch doctor",
		RunE:    runDoctor,
	}
}

func runDoctor(_ *cobra.Command, _ []string) error {
	// 1. git config — report the identity git will ACTUALLY commit as in
	// this directory (respects local repo config and the includeIf blocks
	// gitswitch writes), not just the global config. Surface the global
	// value separately when a directory binding overrides it, so the
	// switch is visible rather than silently hidden.
	gitName, gitEmail, _ := git.EffectiveIdentity()
	row("git (here)", renderIdentity(gitName, gitEmail), gitName != "" && gitEmail != "")

	globalName, globalEmail, _ := git.GlobalIdentity()
	if globalName != gitName || globalEmail != gitEmail {
		row("git (global)", renderIdentity(globalName, globalEmail), globalName != "" && globalEmail != "")
	}

	// 2. ssh ~/.ssh/config Host blocks
	blocks, err := ssh.ParseConfig()
	if err != nil {
		row("ssh-config", fmt.Sprintf("error reading ~/.ssh/config: %v", err), false)
	} else if len(blocks) == 0 {
		row("ssh-config", "no ~/.ssh/config", false)
	} else {
		var hosts []string
		for _, b := range blocks {
			hosts = append(hosts, b.Host)
		}
		row("ssh-config", strings.Join(hosts, ", "), true)
	}

	// 3. ssh -T against each git-hosting host we found
	var sshLogins []string
	for _, b := range blocks {
		if !looksLikeGitHost(b.HostName, b.Host) {
			continue
		}
		// Git-hosting providers all expect the literal user "git".
		welcome, _ := ssh.TestAuth("git@"+b.Host, 6*time.Second)
		switch {
		case welcome == "":
			row("ssh-auth ("+b.Host+")", "no greeting (auth may have failed)", false)
		default:
			row("ssh-auth ("+b.Host+")", welcome, true)
			sshLogins = append(sshLogins, parseLoginFromGreeting(welcome))
		}
	}

	// 4. gh active login
	if !gh.IsInstalled() {
		row("gh", "not installed (skipping)", false)
	} else {
		login, _ := gh.ActiveLogin()
		if login == "" {
			row("gh", "not authenticated (run `gh auth login`)", false)
		} else {
			row("gh", login, true)
		}
	}

	// 5. Is a gitswitch binding being silently shadowed by this repo's
	// local .git/config? Local config wins over the includeIf (git
	// precedence: local > global), so the binding can be defeated without
	// any sign. Detect it so the summary below can't lie with "all agree".
	var shadowOverrides []string
	var shadowIdentity string
	if cwd, err := os.Getwd(); err == nil {
		if cfg, err := identity.Load(); err == nil {
			if b := cfg.FindBindingForDir(cwd); b != nil {
				if ov := git.LocalOverrides(""); len(ov) > 0 {
					shadowOverrides = ov
					shadowIdentity = b.Identity
				}
			}
		}
	}

	// 6. cross-check + summary: shadowing first (the root cause), then
	// gh-vs-ssh, then the all-clear.
	ghLogin, _ := gh.ActiveLogin()
	mismatch := false
	if ghLogin != "" {
		for _, s := range sshLogins {
			if s != "" && s != ghLogin {
				mismatch = true
			}
		}
	}
	hasCrossCheck := ghLogin != "" && len(sshLogins) > 0

	switch {
	case len(shadowOverrides) > 0:
		fmt.Println()
		fmt.Println(yellow + bold + "⚠ binding shadowed by local .git/config" + reset)
		fmt.Printf("  %q is bound to this directory, but these keys in this repo's\n", shadowIdentity)
		fmt.Printf("  .git/config override it: %s%s%s\n", bold, strings.Join(shadowOverrides, ", "), reset)
		fmt.Printf("  gitswitch isn't in control here. Take over with: %sgitswitch use %s .%s\n",
			dim, shadowIdentity, reset)
	case mismatch:
		fmt.Println()
		fmt.Println(red + bold + "✗ identity mismatch" + reset)
		fmt.Printf("  gh says you are:  %s\n", ghLogin)
		fmt.Printf("  ssh says you are: %s\n", strings.Join(sshLogins, ", "))
	case hasCrossCheck:
		fmt.Println()
		fmt.Println(green + bold + "✓ all layers agree" + reset)
	}
	return nil
}

func row(label, value string, ok bool) {
	mark := green + "✓" + reset
	if !ok {
		mark = yellow + "•" + reset
	}
	fmt.Printf("  %s %s%-22s%s %s\n", mark, dim, label, reset, value)
}

func looksLikeGitHost(hostname, alias string) bool {
	for _, s := range []string{hostname, alias} {
		s = strings.ToLower(s)
		if strings.Contains(s, "github.com") || strings.Contains(s, "gitlab.com") ||
			strings.Contains(s, "bitbucket.org") {
			return true
		}
	}
	return false
}

// parseLoginFromGreeting extracts "octocat" from "Hi octocat! You've ..."
func parseLoginFromGreeting(g string) string {
	g = strings.TrimSpace(g)
	for _, prefix := range []string{"Hi ", "Welcome ", "Hello "} {
		if strings.HasPrefix(g, prefix) {
			rest := strings.TrimPrefix(g, prefix)
			// "octocat!" → "octocat"
			if idx := strings.IndexAny(rest, "!,. "); idx > 0 {
				return rest[:idx]
			}
			return rest
		}
	}
	return ""
}

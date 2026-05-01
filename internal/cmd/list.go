package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/identity"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Show all configured identities and bindings.",
		Long: `Lists every identity gitswitch knows about plus every directory
binding currently active. Reads ~/.config/gitswitch/config.json.

This is what you should run before "gitswitch use" if you've
forgotten the names of your identities.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runList()
		},
	}
}

func runList() error {
	cfg, err := identity.Load()
	if err != nil {
		return err
	}

	if len(cfg.Identities) == 0 {
		fmt.Println(yellow + "No identities configured." + reset)
		fmt.Println(dim + "  Run `gitswitch init` to detect identities on this machine." + reset)
		return nil
	}

	// Identities table — align names so emails line up.
	maxName := len("name")
	for _, id := range cfg.Identities {
		if len(id.Name) > maxName {
			maxName = len(id.Name)
		}
	}

	fmt.Println(bold + "Identities" + reset)
	fmt.Println()
	for _, id := range cfg.Identities {
		fmt.Printf("  %s%-*s%s   %s\n",
			green, maxName, id.Name, reset, summarize(id))
	}

	fmt.Println()
	if len(cfg.Bindings) > 0 {
		fmt.Println(bold + "Bindings" + reset)
		fmt.Println()
		bindings := append([]identity.Binding(nil), cfg.Bindings...)
		// Sort by directory for stable output across runs.
		sort.Slice(bindings, func(i, j int) bool {
			return bindings[i].Directory < bindings[j].Directory
		})
		// Align directories so the arrows line up.
		maxDir := 0
		for _, b := range bindings {
			s := shortPath(b.Directory)
			if len(s) > maxDir {
				maxDir = len(s)
			}
		}
		for _, b := range bindings {
			fmt.Printf("  %s%-*s%s  %s→%s  %s%s%s\n",
				dim, maxDir, shortPath(b.Directory), reset,
				dim, reset,
				green, b.Identity, reset)
		}
	} else {
		fmt.Println(dim + "(no directory bindings yet — `gitswitch use <name> <dir>` to add one)" + reset)
	}
	return nil
}

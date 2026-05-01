package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/identity"
	"github.com/target-ops/gitswitch/internal/legacy"
)

func newMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Import identities from a v0.2.x ~/.config.ini.",
		Long: `Reads ~/.config.ini left over by the Python v0.2.x release and
imports each identity into the new JSON config at
~/.config/gitswitch/config.json. Idempotent: re-running upserts (by
identity name) rather than duplicating.

The legacy file is left in place after migration. Once you've
verified the new config works, you can delete it by hand:

  rm ~/.config.ini
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			yes, _ := cmd.Flags().GetBool("yes")
			return runMigrate(yes)
		},
	}
	cmd.Flags().BoolP("yes", "y", false, "skip the confirmation prompt")
	return cmd
}

func runMigrate(assumeYes bool) error {
	if !legacy.Exists() {
		fmt.Println(yellow + "no legacy ~/.config.ini found — nothing to migrate." + reset)
		return nil
	}

	imported, err := legacy.Parse()
	if err != nil {
		return fmt.Errorf("parse %s: %w", legacy.LegacyPath(), err)
	}
	if len(imported) == 0 {
		fmt.Println(yellow + "~/.config.ini is empty or malformed — nothing to import." + reset)
		return nil
	}

	cfg, err := identity.Load()
	if err != nil {
		return err
	}

	fmt.Println(bold + "Found " + countWord(len(imported), "identity", "identities") + " in " + legacy.LegacyPath() + ":" + reset)
	fmt.Println()
	for i, id := range imported {
		fmt.Printf("  %d. %s%s%s   %s%s%s\n",
			i+1, bold, id.Name, reset, dim, summarize(id), reset)
	}
	fmt.Println()

	if !assumeYes {
		var ok bool
		err := huh.NewConfirm().
			Title("Import them into " + identity.Path() + "?").
			Affirmative("Yes, import").
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
			fmt.Println(yellow + "cancelled — nothing imported." + reset)
			return nil
		}
	}

	for _, id := range imported {
		cfg.Upsert(id)
	}
	if err := identity.Save(cfg); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(green + bold + "✓ imported " + countWord(len(imported), "identity", "identities") + reset)
	fmt.Println()
	fmt.Println(dim + "Next: bind each one to a directory." + reset)
	for _, id := range imported {
		fmt.Printf("  %sgitswitch use %s ~/some/dir%s\n", dim, id.Name, reset)
	}
	fmt.Println()
	fmt.Println(dim + "When ready, you can remove the legacy file: rm " + legacy.LegacyPath() + reset)
	return nil
}

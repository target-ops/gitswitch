package cmd

import (
	"github.com/spf13/cobra"

	"github.com/target-ops/gitswitch/internal/style"
)

// NewRootCommand wires up the cobra command tree.
func NewRootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "gitswitch",
		Short: "Stop committing as the wrong person.",
		// Cobra auto-generates a `completion` subcommand that emits
		// shell completion scripts. Useful, but in `gitswitch --help`
		// it adds noise next to our five real commands. Hide it from
		// the listing — power users still find it via
		// `gitswitch completion --help`.
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		Long: `gitswitch — manage multiple Git identities (SSH, gh CLI, signing) ` +
			`per directory.

Bind a directory to an identity and gitswitch keeps git, ssh, gh, and
commit-signing in lockstep — no more accidentally committing personal
work to a company repo.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		// PersistentPreRun fires before any subcommand. Use it to honour
		// --no-color before the subcommand renders anything.
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			noColor, _ := cmd.Flags().GetBool("no-color")
			if noColor {
				style.SetEnabled(false)
				applyColorState()
			}
		},
	}

	root.PersistentFlags().Bool("no-color", false,
		"disable colour and decoration in output (also honours $NO_COLOR)")

	root.AddCommand(newDoctorCommand())
	root.AddCommand(newInitCommand())
	root.AddCommand(newUseCommand())
	root.AddCommand(newGuardCommand())
	root.AddCommand(newWhyCommand())
	root.AddCommand(newListCommand())
	root.AddCommand(newAddCommand())
	root.AddCommand(newDeleteCommand())
	root.AddCommand(newRenameCommand())
	return root
}

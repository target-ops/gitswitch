package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand wires up the cobra command tree.
func NewRootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "gitswitch",
		Short: "Stop committing as the wrong person.",
		Long: `gitswitch — manage multiple Git identities (SSH, gh CLI, signing) ` +
			`per directory.

Bind a directory to an identity and gitswitch keeps git, ssh, gh, and
commit-signing in lockstep — no more accidentally committing personal
work to a company repo.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newDoctorCommand())
	root.AddCommand(newInitCommand())
	return root
}

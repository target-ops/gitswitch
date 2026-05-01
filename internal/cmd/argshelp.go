package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// argError formats a friendly "you used this wrong" error for argument
// validation failures. The default cobra message ("accepts 1 arg(s),
// received 0") is technically correct but tells the user nothing
// useful — they have to leave the error and run --help to recover.
//
// Usage: pass a one-line description of what's wrong, then a list of
// worked examples. The output is:
//
//   <usage>
//
//   examples:
//     <ex 1>
//     <ex 2>
//
//   run "<cmdpath> --help" for the full reference.
func argError(cmd *cobra.Command, usage string, examples ...string) error {
	var sb strings.Builder
	sb.WriteString(usage)
	if len(examples) > 0 {
		sb.WriteString("\n\nexamples:")
		for _, ex := range examples {
			sb.WriteString("\n  ")
			sb.WriteString(ex)
		}
	}
	if cmd != nil {
		sb.WriteString("\n\nrun \"")
		sb.WriteString(cmd.CommandPath())
		sb.WriteString(" --help\" for the full reference.")
	}
	return errors.New(sb.String())
}

// exactArgsHelp is cobra.ExactArgs(n) with a useful error message on
// miscount. Pass the same usage + example shape as argError.
func exactArgsHelp(n int, usage string, examples ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == n {
			return nil
		}
		return argError(cmd, usage, examples...)
	}
}

// rangeArgsHelp is cobra.RangeArgs(min, max) with a useful error message.
func rangeArgsHelp(minArgs, maxArgs int, usage string, examples ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= minArgs && len(args) <= maxArgs {
			return nil
		}
		return argError(cmd, usage, examples...)
	}
}

package require

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NoArgs returns an error if any args are included.
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf(
			"%q accepts no arguments\n\nUsage:  %s",
			cmd.CommandPath(),
			cmd.UseLine(),
		)
	}
	return nil
}

// ExactArgs returns error if there are no exactly n args.
func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return errors.Errorf(
				"%q requires %d %s\n\nUsage: %s",
				cmd.CommandPath(),
				n,
				"arguments",
				cmd.UseLine(),
			)
		}
		return nil
	}
}

// MinimumNArgs returns error if there's not at least N args.
func MinimumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return errors.Errorf(
				"%q requires at least %d %s\n\nUsage:  %s",
				cmd.CommandPath(),
				n,
				"arguments",
				cmd.UseLine(),
			)
		}
		return nil
	}
}

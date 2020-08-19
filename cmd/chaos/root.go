package main

import (
	"io"

	"github.com/girikuncoro/chaos/pkg/action"
	"github.com/spf13/cobra"
)

func newRootCmd(actionConfig *action.Configuration, out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "chaos",
		Short:        "The Chaos management tool for Kubernetes",
		Long:         "The Chaos management tool for Kubernetes",
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	settings.AddFlags(flags)

	// Add subcommands
	cmd.AddCommand(
		newInitCmd(actionConfig, out),
		newRepoCmd(out),
		newInstallCmd(actionConfig, out),
	)

	return cmd, nil
}

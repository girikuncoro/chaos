package main

import (
	"io"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/action"
	"github.com/spf13/cobra"
)

const installDesc = `
This command installs an experiment chart into the target cluster.
This doesn't run chaos experiment yet.
`

func newInstallCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInstall(cfg)

	cmd := &cobra.Command{
		Use:   "install [CHART]",
		Short: "install an experiment chart",
		Long:  installDesc,
		Args:  require.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := runInstall(args, client, out)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func runInstall(args []string, client *action.Install, out io.Writer) error {
	debug("Installing experiment chart")

	// TODO: Currently single repo contains single experiment chart.
	// When we have repo that contains many charts, this needs FIX.
	chart := args[0]

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return err
	}

	debug("CHART PATH: %s\n", cp)
	return client.Run(cp)
}

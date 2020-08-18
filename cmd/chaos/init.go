package main

import (
	"io"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/action"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newInitCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInit(cfg)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize litmus chaos operator",
		Long:  "This command installs Litmus Chaos operator to a target cluster.",
		Args:  require.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			err := runInit(args, client, out)
			if err != nil {
				return err
			}
			return nil
		},
	}

	addInitFlags(cmd, cmd.Flags(), client)
	return cmd
}

func addInitFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.Init) {
	f.StringVarP(&client.Image, "litmus-image", "i", "", "override litmus image")
}

func runInit(args []string, client *action.Init, out io.Writer) error {
	debug("Initializing Litmus on Kubernetes cluster")
	return client.Run()
}

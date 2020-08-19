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
	f.StringVarP(&client.Opts.ImageSpec, "litmus-image", "i", "", "override litmus image")
	f.IntVar(&client.Opts.Replicas, "replicas", 1, "Amount of chaos-operator instances to run on the cluster")
}

func runInit(args []string, client *action.Init, out io.Writer) error {
	debug("Initializing Litmus on Kubernetes cluster")
	return client.Run()
}

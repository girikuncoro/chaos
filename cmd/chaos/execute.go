package main

import (
	"io"
	"time"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/action"
	"github.com/girikuncoro/chaos/pkg/cli/values"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const executeDesc = `
This command runs chaos on specified Kubernetes resource object.
The run argument must be a chart reference.
`

func newExecuteCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewExecute(cfg)
	valueOpts := &values.Options{}
	cmd := &cobra.Command{
		Use:     "execute [NAME] [RESOURCE_KIND]/[RESOURCE_NAME] [EXPERIMENT]",
		Aliases: []string{"exec"},
		Short:   "execute chaos experiment on specified object using installed chart",
		Long:    executeDesc,
		Args:    require.MinimumNArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			err := runExecute(args, client, valueOpts, out)
			if err != nil {
				return err
			}
			return nil
		},
	}

	addExecuteFlags(cmd, cmd.Flags(), client, valueOpts)
	return cmd
}

func addExecuteFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.Execute, valueOpts *values.Options) {
	f.BoolVar(&client.Wait, "wait", false, "if set, will wait until chaos experiments have been completely executed. It will wait for as long as --timeout")
	f.DurationVar(&client.Timeout, "timeout", 300*time.Second, "time to wait for waiting experiment completion")
	addValueOptionsFlags(f, valueOpts)
}

func runExecute(args []string, client *action.Execute, valueOpts *values.Options, out io.Writer) error {
	client.Namespace = settings.Namespace()

	name, exp, err := client.NameAndChart(args)
	if err != nil {
		return err
	}
	client.TestName = name
	client.ExperimentName = exp
	return client.Run()
}

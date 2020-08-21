package main

import (
	"fmt"
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
	f.StringVar(&client.ChaosTestOptions.ServiceAccount, "service-account", "", "service account override name to be used to perform the chaos test")
	f.StringVar(&client.ChaosTestOptions.AnnotationCheck, "annotation-check", "true", "flag to control annotation checks on applications as prerequisites for chaos")
	f.StringVar(&client.ChaosTestOptions.EngineState, "engine-state", "active", "flag to control the state of the chaosengine")
	f.StringVar(&client.ChaosTestOptions.AuxiliaryAppInfo, "auxiliary-app-info", "", "flag to specify one or more app namespace-label pairs whose health is also monitored as part of the chaos experiment")
	f.StringVar(&client.ChaosTestOptions.AppLabel, "app-label", "", "flag to specify unique label of application under test")
	f.StringVar(&client.ChaosTestOptions.AppKind, "app-kind", "", "flag to specify resource kind of application under test")
	f.BoolVar(&client.ChaosTestOptions.Monitoring, "monitoring", false, "flag to enable collection of simple chaos metrics")
	f.StringVar(&client.ChaosTestOptions.JobCleanupPolicy, "job-cleanup-policy", "delete", "flag to control cleanup of chaos experiment job post execution of chaos")
	f.IntVar(&client.ChaosTestOptions.ChaosInterval, "chaos-interval", 10, "flag to control interval between chaos tests")
	f.IntVar(&client.ChaosTestOptions.ChaosDuration, "chaos-duration", 20, "flag to control duration of the chaos tests")
	addValueOptionsFlags(f, valueOpts)
}

func runExecute(args []string, client *action.Execute, valueOpts *values.Options, out io.Writer) error {
	client.Namespace = settings.Namespace()

	name, exp, dep, err := client.ExperimentAndDeployment(args)
	if err != nil {
		return err
	}
	client.TestName = name
	client.Experiment = exp
	client.ChaosTestOptions.AppKind = dep.TypeMeta.Kind
	client.ChaosTestOptions.AppLabel = buildAppLabel(dep.Labels)
	return client.Run(dep)
}

// buildAppLabel constructs label with 'key=value' format from sets of object labels.
// If there are multiple labels given, the first set is chosen.
func buildAppLabel(labels map[string]string) string {
	for k, v := range labels {
		return fmt.Sprintf("%s=%s", k, v)
	}
	return ""
}

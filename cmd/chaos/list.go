package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/girikuncoro/chaos/pkg/action"
	"github.com/girikuncoro/chaos/pkg/chaostest"
	"github.com/girikuncoro/chaos/pkg/cli/output"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var listHelp = `
This command lists all of the executed chaos tests for a specified namespace (uses current namespace context if namespace not specified).
`

func newListCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewList(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list tests",
		Long:    listHelp,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}

			if client.Short {
				for _, res := range results {
					fmt.Fprintln(out, res.Name)
				}
				return nil
			}

			return outfmt.Write(out, newChaosTestListWriter(results))
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&client.Short, "short", "q", false, "output short (quiet) listing format")
	bindOutputFlag(cmd, &outfmt)
	return cmd
}

type chaosTestElement struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Updated     string `json:"updated"`
	Status      string `json:"status"`
	Experiments string `json:"experiments"`
}

type chaosTestListWriter struct {
	chaosTests []chaosTestElement
}

func newChaosTestListWriter(chaosTests []*chaostest.ChaosTest) *chaosTestListWriter {
	elements := make([]chaosTestElement, 0, len(chaosTests))
	for _, ct := range chaosTests {
		element := chaosTestElement{
			Name:        ct.Name,
			Namespace:   ct.Namespace,
			Updated:     "TODO",
			Status:      "TODO",
			Experiments: buildExperiments(ct.Experiments),
		}
		elements = append(elements, element)
	}
	return &chaosTestListWriter{elements}
}

func (t *chaosTestListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "UPDATED", "STATUS", "EXPERIMENTS")
	for _, t := range t.chaosTests {
		table.AddRow(t.Name, t.Namespace, t.Updated, t.Status, t.Experiments)
	}
	return output.EncodeTable(out, table)
}

func buildExperiments(results []*chaostest.ExperimentResult) string {
	s := make([]string, len(results))
	for i, r := range results {
		s[i] = r.Experiment + "=" + r.Result
	}
	return strings.Join(s, ",")
}

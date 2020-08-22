package main

import (
	"io"

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
			return outfmt.Write(out, newChaosTestListWriter(results))
		},
	}

	bindOutputFlag(cmd, &outfmt)
	return cmd
}

type chaosTestElement struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Updated   string `json:"updated"`
	Status    string `json:"status"`
	Chart     string `json:"chart"`
}

type chaosTestListWriter struct {
	chaosTests []chaosTestElement
}

func newChaosTestListWriter(chaosTests []*chaostest.ChaosTest) *chaosTestListWriter {
	_ = make([]chaosTestElement, 0, len(chaosTests))
	return nil
}

func (t *chaosTestListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "UPDATED", "STATUS", "CHART")
	for _, t := range t.chaosTests {
		table.AddRow(t.Name, t.Namespace, t.Updated, t.Status, t.Chart)
	}
	return output.EncodeTable(out, table)
}

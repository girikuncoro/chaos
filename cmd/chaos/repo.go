package main

import (
	"io"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/spf13/cobra"
)

var repoChaos = `
This command consists of multiple subcommands to interact with experiment chart repositories.

It can be used to add and list experiment chart repositories.
`

func newRepoCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo add|list [ARGS]",
		Short: "add and list experiment chart repositories",
		Long:  repoChaos,
		Args:  require.NoArgs,
	}

	cmd.AddCommand(newRepoAddCmd(out))

	return cmd
}

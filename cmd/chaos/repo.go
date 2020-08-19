package main

import (
	"io"
	"os"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/pkg/errors"
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
	cmd.AddCommand(newRepoListCmd(out))

	return cmd
}

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

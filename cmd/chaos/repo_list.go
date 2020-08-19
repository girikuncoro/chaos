package main

import (
	"io"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/cli/output"
	"github.com/girikuncoro/chaos/pkg/repo"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newRepoListCmd(out io.Writer) *cobra.Command {
	var outfmt output.Format
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list chart repositories",
		Args:    require.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := repo.LoadFile(settings.RepositoryConfig)
			if isNotExist(err) || len(f.Repositories) == 0 {
				return errors.New("no repositories to show")
			}
			return outfmt.Write(out, &repoListWriter{f.Repositories})
		},
	}

	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type repoListWriter struct {
	repos []*repo.Entry
}

func (r *repoListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "URL")
	for _, re := range r.repos {
		table.AddRow(re.Name, re.URL)
	}
	return output.EncodeTable(out, table)
}

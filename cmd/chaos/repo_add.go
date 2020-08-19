package main

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/repo"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type repoAddOptions struct {
	name     string
	url      string
	repoFile string
}

func newRepoAddCmd(out io.Writer) *cobra.Command {
	o := &repoAddOptions{}

	cmd := &cobra.Command{
		Use:   "add [NAME] [URL]",
		Short: "add an experiment chart repository",
		Args:  require.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			o.url = args[1]
			o.repoFile = settings.RepositoryConfig
			return o.run(out)
		},
	}

	return cmd
}

func (o *repoAddOptions) run(out io.Writer) error {
	// Ensure file directory exists
	err := os.MkdirAll(filepath.Dir(o.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// File lock to ensure config atomicity
	fileLock := flock.New(strings.Replace(o.repoFile, filepath.Ext(o.repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(o.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(o.name) {
		return errors.Errorf("repository name (%s) already exists, please specify different name", o.name)
	}

	c := repo.Entry{
		Name: o.name,
		URL:  o.url,
	}

	_, err = repo.NewChartRepository(&c)
	if err != nil {
		return err
	}

	// TODO: download index file, update, and write

	return nil
}

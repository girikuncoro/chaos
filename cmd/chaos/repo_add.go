package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/girikuncoro/chaos/cmd/chaos/require"
	"github.com/girikuncoro/chaos/pkg/getter"
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

	e := repo.Entry{
		Name: o.name,
		URL:  o.url,
	}

	client, _ := getter.NewHTTPGetter()
	r, err := repo.NewChartRepository(&e, client)
	if err != nil {
		return err
	}

	if e.ExperimentFile, err = r.DownloadExperimentFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid experiment chart repository or cannot be reached", o.url)
	}

	f.Update(&e)

	if err := f.WriteFile(o.repoFile, 0644); err != nil {
		return err
	}
	fmt.Fprintf(out, "%q has been added to your repositories\n", o.name)
	return nil
}

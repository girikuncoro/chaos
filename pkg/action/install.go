package action

import (
	"bytes"
	"io/ioutil"

	"github.com/girikuncoro/chaos/pkg/cli"
	"github.com/girikuncoro/chaos/pkg/repo"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Install performs experiment chart installation in given namespace.
type Install struct {
	cfg *Configuration

	Namespace string
}

// NewInstall creates a new Install object with given configuration.
func NewInstall(cfg *Configuration) *Install {
	return &Install{
		cfg: cfg,
	}
}

// LocateChart looks for a chart directory in known places.
// Currently it excepts charts to be present in cache directory.
func (i *Install) LocateChart(name string, settings *cli.EnvSettings) (string, error) {
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil || len(f.Repositories) == 0 {
		return "", errors.Wrap(err, "no repositories exist, need to add repo first")
	}

	e := f.Get(name)
	if e == nil {
		return "", errors.Errorf("entry %s is not found", name)
	}
	return e.ExperimentFile, nil
}

// Run executes the install operation.
func (i *Install) Run(chartPath string) error {
	// Check reachability of cluster
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(chartPath)
	if err != nil {
		return errors.Wrapf(err, "could not load the experiment file (%s) from cache", chartPath)
	}

	res, err := i.cfg.KubeClient.Build(bytes.NewBuffer(b), true)
	if err != nil {
		return errors.Wrap(err, "unable to build kubernetes objects from experiment file")
	}

	if _, err := i.cfg.KubeClient.Create(res); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

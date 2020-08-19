package action

import (
	"bytes"

	"github.com/girikuncoro/chaos/cmd/chaos/installer"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Init performs an init operation.
type Init struct {
	cfg *Configuration

	Opts      installer.Options
	Namespace string
	Image     string
}

// NewInit creates a new Init object with given configuration.
func NewInit(cfg *Configuration) *Init {
	return &Init{
		cfg: cfg,
	}
}

// Run executes the init operation.
func (i *Init) Run() error {
	// Check reachability of cluster
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return err
	}

	var manifests []string
	var err error
	if manifests, err = installer.LitmusManifests(&i.Opts); err != nil {
		return err
	}

	for _, manifest := range manifests {
		res, err := i.cfg.KubeClient.Build(bytes.NewBufferString(manifest), true)
		if err != nil {
			return errors.Wrap(err, "unable to build kubernetes objects from manifest")
		}

		if _, err := i.cfg.KubeClient.Create(res); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

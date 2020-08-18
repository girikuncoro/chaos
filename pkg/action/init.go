package action

import (
	"fmt"

	"github.com/girikuncoro/chaos/cmd/chaos/installer"
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

	fmt.Println(manifests)

	// TODO: init local config

	// TODO: Create Litmus namespace
	//
	// ns := &v1.Namespace{
	// 	TypeMeta: metav1.TypeMeta{
	// 		APIVersion: "v1",
	// 		Kind:       "Namespace",
	// 	},
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "litmus",
	// 		Labels: map[string]string{
	// 			"name": "litmus",
	// 		},
	// 	},
	// }

	// TODO: build and create the Kubernetes resource
	//
	// resourceList, err := i.cfg.KubeClient.Build(bytes.NewBuffer(buf), true)
	// if err != nil {
	// 	return err
	// }
	// if _, err := i.cfg.KubeClient.Create(resourceList); err != nil && !apierrors.IsAlreadyExists(err) {
	// 	return err
	// }

	return nil
}

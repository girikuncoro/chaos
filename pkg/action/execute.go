package action

import (
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
)

// Execute performs an exec operation.
type Execute struct {
	cfg *Configuration
	ChartPathOptions

	Namespace      string
	ExperimentName string
	Wait           bool
	Timeout        time.Duration
}

// NewExecute creates a new Execute object with given configuration.
func NewExecute(cfg *Configuration) *Execute {
	return &Execute{
		cfg: cfg,
	}
}

// DeploymentAndChart returns deployment resource and chart that should be used.
func (e *Execute) DeploymentAndChart(args []string) (*appsv1.Deployment, error) {
	kind, name, _ := args[0], args[1], args[2]
	if kind != "deployment" && kind != "deploy" {
		return nil, errors.New("currently only supports executing experiment on Deployment resource kind")
	}

	// TODO: Refactor to pass runtime.Object and inteface with generic Get instead
	dep, err := e.cfg.KubeClient.GetDeployment(name, e.Namespace)
	if err != nil {
		return nil, errors.Errorf("deployment %s in namespace %s doesn't exist", name, e.Namespace)
	}

	return dep, nil
}

// Run executes the execute operation.
func (e *Execute) Run() error {
	return nil
}

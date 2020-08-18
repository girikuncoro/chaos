package action

import (
	"github.com/pkg/errors"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/girikuncoro/chaos/pkg/kube"
)

var (
	// errMissingExperiment indicates that an experiment was not provided.
	errMissingExperiment = errors.New("no experiment provided")
)

// Configuration injects the dependencies that all actions share.
type Configuration struct {
	// RESTClientGetter is an interface that loads Kubernetes clients.
	RESTClientGetter RESTClientGetter

	// KubeClient is a Kubernetes API client.
	KubeClient kube.Interface

	Log func(string, ...interface{})
}

// RESTClientGetter gets the rest client
type RESTClientGetter interface {
	ToRESTConfig() (*rest.Config, error)
}

// DebugLog sets the logger that writes debug strings.
type DebugLog func(format string, v ...interface{})

// KubernetesClientSet creates a new kubernetes ClientSet based on the configuration.
func (c *Configuration) KubernetesClientSet() (kubernetes.Interface, error) {
	conf, err := c.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate config for kubernetes client")
	}
	return kubernetes.NewForConfig(conf)
}

// Init initializes action configuration
func (c *Configuration) Init(getter genericclioptions.RESTClientGetter, namespace string, log DebugLog) error {
	kc := kube.New(getter)
	kc.Log = log
	c.RESTClientGetter = getter
	c.KubeClient = kc
	c.Log = log

	return nil
}

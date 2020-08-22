package action

import (
	"github.com/pkg/errors"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/girikuncoro/chaos/pkg/kube"
	"github.com/girikuncoro/chaos/pkg/storage"
	"github.com/girikuncoro/chaos/pkg/storage/driver"
	litmuschaos "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned"
)

var (
	// errMissingExperiment indicates that an experiment was not provided.
	errMissingExperiment = errors.New("no experiment provided")
)

// Configuration injects the dependencies that all actions share.
type Configuration struct {
	// RESTClientGetter is an interface that loads Kubernetes clients.
	RESTClientGetter RESTClientGetter

	// ChaosTests stores records of chaos tests.
	ChaosTests *storage.Storage

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

// LitmusChaosClientSet creates a new litmus chaos ClientSet.
func (c *Configuration) LitmusChaosClientSet() (*litmuschaos.Clientset, error) {
	conf, err := c.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate config for litmus chaos client")
	}
	return litmuschaos.NewForConfig(conf)
}

// Init initializes action configuration
func (c *Configuration) Init(getter genericclioptions.RESTClientGetter, namespace string, log DebugLog) error {
	kc := kube.New(getter)
	kc.Log = log

	// TODO: Streamline kube client and litmus client
	lazyClient := &lazyClient{
		namespace: namespace,
		clientFn:  c.LitmusChaosClientSet,
	}

	// TODO: Pass this through init and add switch cases
	d := driver.NewLitmusCRD(newChaosEngineClient(lazyClient))
	d.Log = log
	store := storage.Init(d)

	c.RESTClientGetter = getter
	c.KubeClient = kc
	c.ChaosTests = store
	c.Log = log

	return nil
}

package kube

import (
	"sync"

	"github.com/pkg/errors"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/deprecated/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// Client represents a client that communicates with the Kubernetes API.
type Client struct {
	Factory   Factory
	Log       func(string, ...interface{})
	Namespace string
}

var addToScheme sync.Once

// New creates a new Client.
func New(getter genericclioptions.RESTClientGetter) *Client {
	if getter == nil {
		getter = genericclioptions.NewConfigFlags(true)
	}
	// Add CRDs to the scheme, missing by default.
	addToScheme.Do(func() {
		if err := apiextv1.AddToScheme(scheme.Scheme); err != nil {
			panic(err)
		}
		if err := apiextv1beta1.AddToScheme(scheme.Scheme); err != nil {
			panic(err)
		}
	})
	return &Client{
		Factory: cmdutil.NewFactory(getter),
		Log:     nopLogger,
	}
}

var nopLogger = func(_ string, _ ...interface{}) {}

// IsReachable tests connectivity to the cluster
func (c *Client) IsReachable() error {
	client, err := c.Factory.KubernetesClientSet()
	if err == genericclioptions.ErrEmptyConfig {
		return errors.New("Kubernetes cluster unreachable")
	}
	if err != nil {
		return errors.Wrap(err, "Kubernetes cluster unreachable")
	}
	if _, err := client.ServerVersion(); err != nil {
		return errors.Wrap(err, "Kubernetes cluster unreachable")
	}
	return nil
}

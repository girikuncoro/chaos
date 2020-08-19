package kube

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/deprecated/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// ErrNoObjectsVisited indicates that during a visit operation, no matching objects were found.
var ErrNoObjectsVisited = errors.New("no objects visited")

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

// Create creates Kubernetes resources specified in the resource list.
func (c *Client) Create(resources ResourceList) (*Result, error) {
	c.Log("creating %d resource(s)", len(resources))
	if err := perform(resources, createResource); err != nil {
		return nil, err
	}
	return &Result{Created: resources}, nil
}

func (c *Client) namespace() string {
	if c.Namespace != "" {
		return c.Namespace
	}
	if ns, _, err := c.Factory.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return v1.NamespaceDefault
}

// newBuilder returns a new resource builder for structured api objects.
func (c *Client) newBuilder() *resource.Builder {
	return c.Factory.NewBuilder().
		ContinueOnError().
		NamespaceParam(c.namespace()).
		DefaultNamespace().
		Flatten()
}

// Build validates for Kubernetes objects and returns unstructured infos.
func (c *Client) Build(reader io.Reader, validate bool) (ResourceList, error) {
	schema, err := c.Factory.Validator(validate)
	if err != nil {
		return nil, err
	}
	result, err := c.newBuilder().
		Unstructured().
		Schema(schema).
		Stream(reader, "").
		Do().Infos()
	return result, err
}

// GetDeployment fetches deployment object from given name
func (c *Client) GetDeployment(name, namespace string) (*appsv1.Deployment, error) {
	client, err := c.Factory.KubernetesClientSet()
	if err != nil {
		return nil, errors.Wrap(err, "Kubernetes cluster unreachable")
	}

	dep, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "deployment %s is not found", name)
	}
	return dep, nil
}

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

func perform(infos ResourceList, fn func(*resource.Info) error) error {
	if len(infos) == 0 {
		return ErrNoObjectsVisited
	}

	errs := make(chan error)
	go batchPerform(infos, fn, errs)

	for range infos {
		err := <-errs
		if err != nil {
			return err
		}
	}
	return nil
}

func batchPerform(infos ResourceList, fn func(*resource.Info) error, errs chan<- error) {
	var kind string
	var wg sync.WaitGroup
	for _, info := range infos {
		currentKind := info.Object.GetObjectKind().GroupVersionKind().Kind
		if kind != currentKind {
			wg.Wait()
			kind = currentKind
		}
		wg.Add(1)
		go func(i *resource.Info) {
			errs <- fn(i)
			wg.Done()
		}(info)
	}
}

func createResource(info *resource.Info) error {
	obj, err := resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
	if err != nil {
		return err
	}
	return info.Refresh(obj, true)
}

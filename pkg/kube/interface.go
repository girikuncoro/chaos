package kube

import (
	"io"

	appsv1 "k8s.io/api/apps/v1"
)

// Interface represents a client that communicates with Kubernetes API.
//
// KubernetesClient must be concurrency safe.
type Interface interface {
	// Create creates one or more resources.
	Create(resources ResourceList) (*Result, error)
	// GetDeployment fetches deployment object from given name
	GetDeployment(name, namespace string) (*appsv1.Deployment, error)
	// Build creates a resource list from a Reader
	Build(reader io.Reader, validate bool) (ResourceList, error)
	// IsReachable checks whether the client is able to connect to the cluster
	IsReachable() error
}

var _ Interface = (*Client)(nil)

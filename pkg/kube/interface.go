package kube

// Interface represents a client that communicates with Kubernetes API.
//
// KubernetesClient must be concurrency safe.
type Interface interface {
	// IsReachable checks whether the client is able to connect to the cluster
	IsReachable() error
}

var _ Interface = (*Client)(nil)

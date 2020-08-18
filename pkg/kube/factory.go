package kube

import (
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/validation"
)

// Factory provides abstractions that allow kubectl command to be extended
// across multiple types of resources and different API sets.
type Factory interface {
	// ToRawKubeConfigLoader returns kubeconfig loader as-is
	ToRawKubeConfigLoader() clientcmd.ClientConfig
	// KubernetesClientSet gives back an external clientset
	KubernetesClientSet() (*kubernetes.Clientset, error)
	// NewBuilder implements common patterns for CLI interactions with generic resources.
	NewBuilder() *resource.Builder
	// Validator returns schema that can validate objects.
	Validator(validate bool) (validation.Schema, error)
}

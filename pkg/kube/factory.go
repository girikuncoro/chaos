package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Factory provides abstractions that allow kubectl command to be extended
// across multiple types of resources and different API sets.
type Factory interface {
	// ToRawKubeConfigLoader returns kubeconfig loader as-is
	ToRawKubeConfigLoader() clientcmd.ClientConfig
	// KubernetesClientSet gives back an external clientset
	KubernetesClientSet() (*kubernetes.Clientset, error)
}

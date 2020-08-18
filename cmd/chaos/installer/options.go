package installer

const (
	defaultImage = "litmuschaos/chaos-operator"
)

// Options control how to install Litmus into a cluster, upgrade, and uninstall Litmus from a cluster.
type Options struct {
	// Namespace is the Kubernetes namespace to use to deploy Litmus.
	Namespace string
	// ImageSpec identifies the image Litmus will use when deployed.
	ImageSpec string
}

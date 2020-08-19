package installer

const (
	defaultImage     = "litmuschaos/chaos-operator"
	defaultNamespace = "litmus"
)

// Options control how to install Litmus into a cluster, upgrade, and uninstall Litmus from a cluster.
type Options struct {
	// Namespace is the Kubernetes namespace to use to deploy Litmus.
	Namespace string
	// ImageSpec identifies the image Litmus will use when deployed.
	ImageSpec string
	// Replicas identify number of chaos-operator instances to run on the cluster
	Replicas int
}

func (opts *Options) getReplicas() *int32 {
	replicas := int32(1)
	if opts.Replicas > 1 {
		replicas = int32(opts.Replicas)
	}
	return &replicas
}

func (opts *Options) getImage() string {
	if opts.ImageSpec == "" {
		return defaultImage + ":latest"
	}
	return opts.ImageSpec
}

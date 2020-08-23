package installer

const (
	defaultOperatorImage = "litmuschaos/chaos-operator"
	defaultRunnerImage   = "litmuschaos/chaos-runner"
	defaultNamespace     = "litmus"
)

// Options control how to install Litmus into a cluster, upgrade, and uninstall Litmus from a cluster.
type Options struct {
	// Namespace is the Kubernetes namespace to use to deploy Litmus.
	Namespace string
	// ImageSpec identifies the image Litmus operator will use when deployed.
	OperatorImageSpec string
	// RunnerImageSpec identifies the image Litmus runner will use when deployed.
	RunnerImageSpec string
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

func (opts *Options) getOperatorImage() string {
	if opts.OperatorImageSpec == "" {
		return defaultOperatorImage + ":latest"
	}
	return opts.OperatorImageSpec
}

func (opts *Options) getRunnerImage() string {
	if opts.RunnerImageSpec == "" {
		return defaultRunnerImage + ":latest"
	}
	return opts.RunnerImageSpec
}

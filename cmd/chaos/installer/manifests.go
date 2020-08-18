package installer

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

const (
	deploymentName     = "chaos-operator-ce"
	serviceAccountName = "litmus"
)

// LitmusManifests gets the Deployment, ServiceAccount manifests
func LitmusManifests(opts *Options) ([]string, error) {
	if opts.Namespace == "" {
		opts.Namespace = defaultNamespace
	}

	ns := Namespace(opts.Namespace)
	sa := ServiceAccount(opts.Namespace)
	dep, err := Deployment(opts)
	if err != nil {
		return []string{}, err
	}

	objs := []runtime.Object{ns, sa, dep}

	manifests := make([]string, len(objs))
	for i, obj := range objs {
		o, err := yaml.Marshal(obj)
		if err != nil {
			return []string{}, err
		}
		manifests[i] = string(o)
	}
	return manifests, err
}

// Namespace gets a namespace object that can be used to generate
// manifest as string.
func Namespace(name string) *v1.Namespace {
	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
	}
	return ns
}

// ServiceAccount gets a service account object that can be used to generate
// manifest as string.
func ServiceAccount(namespace string) *v1.ServiceAccount {
	sa := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
			Labels: map[string]string{
				"name": serviceAccountName,
			},
		},
	}
	return sa
}

// Deployment gets a deployment object that can be used to generate
// manifest as string.
func Deployment(opts *Options) (*appsv1.Deployment, error) {
	dep, err := generateDeployment(opts)
	if err != nil {
		return nil, err
	}
	dep.TypeMeta = metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	return dep, nil
}

func generateDeployment(opts *Options) (*appsv1.Deployment, error) {
	labels := generateLabels(map[string]string{"name": "chaos-operator"})
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      deploymentName,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "chaos-operator",
							Image: "litmuschaos/chaos-operator:1.7.0",
						},
					},
				},
			},
		},
	}
	return d, nil
}

func generateLabels(labels map[string]string) map[string]string {
	labels["app"] = "litmus"
	return labels
}

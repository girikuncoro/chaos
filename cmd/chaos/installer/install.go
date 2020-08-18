package installer

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

// Install uses Kubernetes client to install Litmus chaos operator.
func Install() error {
	// TODO: create deployment, crd of Litmus
	return nil
}

// LitmusManifests gets the Deployment, ServiceAccount manifests
func LitmusManifests(opts *Options) ([]string, error) {
	dep, err := Deployment(opts)
	if err != nil {
		return []string{}, err
	}
	objs := []runtime.Object{dep}
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

func Deployment(opts *Options) (*appsv1.Deployment, error) {
	dep, err := generateDeployment(opts)
	if err != nil {
		return nil, err
	}
	dep.TypeMeta = metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	}
	return dep, nil
}

func generateDeployment(opts *Options) (*appsv1.Deployment, error) {
	labels := generateLabels(map[string]string{"name": "chaos-operator"})
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      "litmus",
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

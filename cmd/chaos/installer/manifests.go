package installer

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

const (
	deploymentName         = "chaos-operator"
	serviceAccountName     = "litmus"
	clusterRoleName        = "litmus"
	clusterRoleBindingName = "litmus"
)

// LitmusManifests gets the Deployment, ServiceAccount manifests
func LitmusManifests(opts *Options) ([]string, error) {
	if opts.Namespace == "" {
		opts.Namespace = defaultNamespace
	}

	ns := Namespace(opts.Namespace)
	sa := ServiceAccount(opts.Namespace)
	cr := ClusterRole()
	crb := ClusterRoleBinding(opts.Namespace)
	dep := Deployment(opts)

	objs := []runtime.Object{ns, sa, cr, crb, dep}

	manifests := make([]string, len(objs))
	for i, obj := range objs {
		o, err := yaml.Marshal(obj)
		if err != nil {
			return []string{}, err
		}
		manifests[i] = string(o)
	}
	manifests = append(manifests, chaosEngineCRD, chaosExperimentCRD, chaosResultCRD)
	return manifests, nil
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

// ClusterRole gets a cluster role object that can be used to generate
// manifest as string.
func ClusterRole() *rbacv1beta1.ClusterRole {
	cr := generateClusterRole()
	cr.TypeMeta = metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1beta1",
		Kind:       "ClusterRole",
	}
	return cr
}

func generateClusterRole() *rbacv1beta1.ClusterRole {
	labels := generateLabels(map[string]string{"name": clusterRoleName})
	cr := &rbacv1beta1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   clusterRoleName,
			Labels: labels,
		},
		Rules: []rbacv1beta1.PolicyRule{
			{
				APIGroups: []string{"", "apps", "batch", "litmuschaos.io", "apps.openshift.io"},
				Resources: []string{"pods", "jobs", "deployments", "replicationcontrollers", "daemonsets", "replicasets", "statefulsets", "deploymentconfigs", "events", "configmaps", "services", "secrets", "chaosengines", "chaosexperiments", "chaosresults"},
				Verbs:     []string{"get", "create", "update", "patch", "delete", "list", "watch", "deletecollection"},
			},
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Resources: []string{"validatingwebhookconfigurations"},
				Verbs:     []string{"get", "create", "list", "delete", "update"},
			},
		},
	}
	return cr
}

// ClusterRoleBinding gets a cluster role binding object that can be used to generate
// manifest as string.
func ClusterRoleBinding(serviceAccountNamespace string) *rbacv1beta1.ClusterRoleBinding {
	crb := &rbacv1beta1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
			Labels: map[string]string{
				"name": clusterRoleBindingName,
			},
		},
		RoleRef: rbacv1beta1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: serviceAccountNamespace,
			},
		},
	}
	return crb
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
func Deployment(opts *Options) *appsv1.Deployment {
	dep := generateDeployment(opts)
	dep.TypeMeta = metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	return dep
}

func generateDeployment(opts *Options) *appsv1.Deployment {
	labels := generateLabels(map[string]string{"name": "chaos-operator"})
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      deploymentName,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: opts.getReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []v1.Container{
						{
							Name:            "chaos-operator",
							Image:           opts.getImage(),
							Command:         []string{"chaos-operator"},
							ImagePullPolicy: v1.PullAlways,
							Env: []v1.EnvVar{
								{Name: "CHAOS_RUNNER_IMAGE", Value: opts.getImage()},
								{Name: "WATCH_NAMESPACE", Value: ""},
								{
									Name: "POD_NAME",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{Name: "OPERATOR_NAME", Value: "chaos-operator"},
							},
						},
					},
				},
			},
		},
	}
	return d
}

func generateLabels(labels map[string]string) map[string]string {
	labels["app"] = "litmus"
	return labels
}

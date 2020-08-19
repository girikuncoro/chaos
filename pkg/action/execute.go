package action

import (
	"bytes"
	"strings"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

var ErrInvalidArgFormat = errors.New("should provide argument with format 'resource_kind/resource_name'")

// Execute performs an exec operation.
type Execute struct {
	cfg *Configuration
	ChartPathOptions

	Namespace      string
	TestName       string
	ExperimentName string
	Wait           bool
	Timeout        time.Duration
}

// NewExecute creates a new Execute object with given configuration.
func NewExecute(cfg *Configuration) *Execute {
	return &Execute{
		cfg: cfg,
	}
}

// NameAndChart returns name and experirment that should be used. This returns error when deployment name is invalid.
func (e *Execute) NameAndChart(args []string) (string, string, error) {
	name := args[0]
	exp := args[2]

	resourceKind, resourceName, err := resourceKindAndName(args[1])
	if err != nil {
		return name, exp, err
	}

	if resourceKind != "deployment" && resourceKind != "deploy" {
		return name, exp, errors.New("currently only supports executing experiment on Deployment resource kind")
	}

	// TODO: Refactor to pass runtime.Object and inteface with generic Get instead
	_, err = e.cfg.KubeClient.GetDeployment(resourceName, e.Namespace)
	if err != nil {
		return name, exp, errors.Errorf("deployment %s in namespace %s doesn't exist", resourceName, e.Namespace)
	}

	// TODO: Validate chaos experiment exists in the target cluster

	return name, exp, nil
}

// Run executes the execute operation.
func (e *Execute) Run() error {
	sa := ServiceAccount(e.ExperimentName+"-sa", e.Namespace)
	r := Role(e.ExperimentName+"-sa", e.Namespace)
	rb := RoleBinding(e.ExperimentName+"-sa", e.Namespace)

	objs := []runtime.Object{sa, r, rb}

	for _, obj := range objs {
		o, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}

		res, err := e.cfg.KubeClient.Build(bytes.NewBuffer(o), true)
		if err != nil {
			return errors.Wrap(err, "unable to build kubernetes objects from manifest")
		}

		if _, err := e.cfg.KubeClient.Create(res); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	// Step 4: Annotate deployment

	// Step 5: Chaos Engine

	// TODO: Override experiment values

	return nil
}

// resourceKindAndName extracts resource kind and resource name
// from input with format 'resource_kind/resource_name'
func resourceKindAndName(param string) (string, string, error) {
	if !strings.Contains(param, "/") {
		return "", "", ErrInvalidArgFormat
	}

	s := strings.Split(param, "/")
	if len(s) > 2 {
		return "", "", ErrInvalidArgFormat
	}

	kind, name := s[0], s[1]
	return kind, name, nil
}

// ServiceAccount gets a service account object that can be used to generate
// manifest as string.
func ServiceAccount(name, namespace string) *v1.ServiceAccount {
	sa := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"name": name,
			},
		},
	}
	return sa
}

// Role gets a role object that can be used to generate
// manifest as string.
func Role(name, namespace string) *rbacv1beta1.Role {
	r := &rbacv1beta1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"name": name},
		},
		Rules: []rbacv1beta1.PolicyRule{
			{
				APIGroups: []string{"", "litmuschaos.io", "batch", "apps"},
				Resources: []string{"pods", "jobs", "pods/exec", "pods/log", "events", "chaosengines", "chaosexperiments", "chaosresults"},
				Verbs:     []string{"create", "list", "get", "patch", "update", "delete"},
			},
		},
	}
	return r
}

// RoleBinding gets a role binding object that can be used to generate
// manifest as string.
func RoleBinding(name, namespace string) *rbacv1beta1.RoleBinding {
	rb := &rbacv1beta1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"name": name},
		},
		RoleRef: rbacv1beta1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     name,
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      name,
				Namespace: namespace,
			},
		},
	}
	return rb
}

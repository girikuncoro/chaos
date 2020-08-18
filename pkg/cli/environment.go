package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	namespace string
	config    *genericclioptions.ConfigFlags

	// KubeConfig is the path to the kubeconfig file.
	KubeConfig string
	// KubeContext is the name of the kubeconfig context.
	KubeContext string
	// Kubernetes API Server Endpoint for authentication.
	KubeAPIServer string
	// Debug indicates whether or not Chaos is running in Debug mode.
	Debug bool
}

func New() *EnvSettings {
	env := &EnvSettings{
		namespace:     os.Getenv("CHAOS_NAMESPACE"),
		KubeContext:   os.Getenv("CHAOS_KUBECONTEXT"),
		KubeAPIServer: os.Getenv("CHAOS_KUBEAPISERVER"),
	}
	env.Debug, _ = strconv.ParseBool(os.Getenv("CHAOS_DEBUG"))

	// bind to kubernetes config flags
	env.config = &genericclioptions.ConfigFlags{
		Namespace:  &env.namespace,
		Context:    &env.KubeContext,
		APIServer:  &env.KubeAPIServer,
		KubeConfig: &env.KubeConfig,
	}
	return env
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.namespace, "namespace", "n", s.namespace, "namespace scope for this request")
	fs.StringVar(&s.KubeConfig, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.KubeAPIServer, "kube-apiserver", s.KubeAPIServer, "the address and the port for Kubernetes API server")
	fs.BoolVar(&s.Debug, "debug", s.Debug, "enable verbose output")
}

func (s *EnvSettings) EnvVars() map[string]string {
	envvars := map[string]string{
		"CHAOS_BIN":       os.Args[0],
		"CHAOS_DEBUG":     fmt.Sprint(s.Debug),
		"CHAOS_NAMESPACE": s.Namespace(),
	}
	if s.KubeConfig != "" {
		envvars["KUBECONFIG"] = s.KubeConfig
	}
	return envvars
}

// Namespace gets the namespace from configuration.
func (s *EnvSettings) Namespace() string {
	if ns, _, err := s.config.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

// RESTClientGetter gets the kubeconfig from EnvSettings.
func (s *EnvSettings) RESTClientGetter() genericclioptions.RESTClientGetter {
	return s.config
}

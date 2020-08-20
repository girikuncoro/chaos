module github.com/girikuncoro/chaos

go 1.13

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.6.0

require (
	github.com/Sirupsen/logrus v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/gofrs/flock v0.7.1
	github.com/gosuri/uitable v0.0.4
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/cli-runtime v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.18.8
	sigs.k8s.io/yaml v1.2.0
)

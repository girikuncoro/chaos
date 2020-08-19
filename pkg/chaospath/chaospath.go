package chaospath

import (
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

// ConfigPath returns the default path of local config directory used by chaos
func ConfigPath(elem ...string) string {
	return filepath.Join(homedir.HomeDir(), ".chaos", filepath.Join(elem...))
}

// CachePath returns the default path of local chaos directory used by chaos
func CachePath(elem ...string) string {
	return filepath.Join(homedir.HomeDir(), ".chaos", filepath.Join(elem...))
}

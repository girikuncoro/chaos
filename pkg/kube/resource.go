package kube

import "k8s.io/cli-runtime/pkg/resource"

// ResourceList provides convenience methods for comparing collections of Infos.
type ResourceList []*resource.Info

// Append adds an Info to the Result.
func (r *ResourceList) Append(val *resource.Info) {
	*r = append(*r, val)
}

// Get returns the Info from the result that matches the name and kind.
func (r ResourceList) Get(info *resource.Info) *resource.Info {
	for _, i := range r {
		if isMatchingInfo(i, info) {
			return i
		}
	}
	return nil
}

// isMatchingInfo returns true if infos match on Name and GroupVersionKind.
func isMatchingInfo(a, b *resource.Info) bool {
	return a.Name == b.Name && a.Namespace == b.Namespace && a.Mapping.GroupVersionKind.Kind == b.Mapping.GroupVersionKind.Kind
}

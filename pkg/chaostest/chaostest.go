package chaostest

// ChaosTest describes an execution of an experiment chart
// on kubernetes resource object.
type ChaosTest struct {
	Name string `json:"name,omitempty"`
	Info *Info  `json:"info,omitempty"`
	// TODO: Make chart an object
	Chart     string `json:"chart,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// SetStatus is a helper for setting status on a ChaosTest.
func (t *ChaosTest) SetStatus(status Status, msg string) {
	t.Info.Status = status
	t.Info.Description = msg
}

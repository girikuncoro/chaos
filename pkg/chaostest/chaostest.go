package chaostest

// ChaosTest describes an execution of an experiment chart
// on kubernetes resource object.
type ChaosTest struct {
	Name        string              `json:"name,omitempty"`
	Info        *Info               `json:"info,omitempty"`
	Experiments []*ExperimentResult `json:"experiments,omitempty"`
	Namespace   string              `json:"namespace,omitempty"`
}

// ExperimentResult describes the result of executed chaos test experiment.
type ExperimentResult struct {
	Experiment string `json:"experiment"`
	Result     string `json:"result"`
	Phase      string `json:"phase"`
}

// SetStatus is a helper for setting status on a ChaosTest.
func (t *ChaosTest) SetStatus(status Status, msg string) {
	t.Info.Status = status
	t.Info.Description = msg
}

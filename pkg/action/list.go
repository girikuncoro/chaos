package action

import (
	"fmt"
	"strings"

	"github.com/girikuncoro/chaos/pkg/chaostest"
)

// ListStates represents zero or more status code that a list item may have set
type ListStates uint

const (
	// ListExecuted filters on status "executed"
	ListExecuted ListStates = 1 << iota
	// ListPendingExecute filters on status "pending" (execution in progress)
	ListPendingExecute
	// ListPass filters on status "pass" (executed experiment results pass)
	ListPass
	// ListFail filters on status "fail" (executed experiment results failure)
	ListFail
	// ListUnknown filters on an unknown status
	ListUnknown
)

// FromName takes state name and returns ListStates representation.
func (s ListStates) FromName(str string) ListStates {
	switch str {
	case "executed":
		return ListExecuted
	case "pending-execute":
		return ListPendingExecute
	case "pass":
		return ListPass
	case "fail":
		return ListFail
	}
	return ListUnknown
}

// List is the action for listing chaos tests.
type List struct {
	cfg *Configuration

	AllNamespaces bool
	Short         bool
	StateMask     ListStates
}

// NewList constructs a new List object
func NewList(cfg *Configuration) *List {
	return &List{
		StateMask: ListExecuted,
		cfg:       cfg,
	}
}

// Run executes the list command, returning a set of matches.
func (l *List) Run() ([]*chaostest.ChaosTest, error) {
	if err := l.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}

	// TODO: Implement filter
	results, err := l.cfg.ChaosTests.List()
	if err != nil {
		return nil, err
	}

	// Conclude the status of chaos test
	for _, res := range results {
		for _, exp := range res.Experiments {
			if exp.Phase == "" {
				res.SetStatus(chaostest.StatusPending, fmt.Sprintf("experiment is starting"))
			}
			if strings.ToLower(exp.Phase) == chaostest.StatusRunning.String() {
				res.SetStatus(chaostest.StatusRunning, fmt.Sprintf("experiment %q is currently running", exp.Experiment))
			}
			if strings.ToLower(exp.Phase) == chaostest.StatusCompleted.String() {
				res.SetStatus(chaostest.StatusCompleted, "all experiments have been completed")
			}
		}
	}
	return results, nil
}

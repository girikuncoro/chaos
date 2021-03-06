package driver

import (
	"github.com/girikuncoro/chaos/pkg/chaostest"
	chaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/typed/litmuschaos/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ Driver = (*LitmusCRDs)(nil)

// LitmusCRDsDriverName is the string name of the driver.
const LitmusCRDsDriverName = "LitmusCRD"

type LitmusCRDs struct {
	chaosEngineImpl chaosv1alpha1.ChaosEngineInterface
	chaosResultImpl chaosv1alpha1.ChaosResultInterface
	Log             func(string, ...interface{})
}

func NewLitmusCRD(ceImpl chaosv1alpha1.ChaosEngineInterface, crImpl chaosv1alpha1.ChaosResultInterface) *LitmusCRDs {
	return &LitmusCRDs{
		chaosEngineImpl: ceImpl,
		chaosResultImpl: crImpl,
		Log:             func(_ string, _ ...interface{}) {},
	}
}

// Name returns the name of driver.
func (l *LitmusCRDs) Name() string {
	return LitmusCRDsDriverName
}

// List fetches all chaostests.
func (l *LitmusCRDs) List() ([]*chaostest.ChaosTest, error) {
	opts := metav1.ListOptions{}
	chaosEngineList, err := l.chaosEngineImpl.List(opts)
	if err != nil {
		l.Log("list: failed to list: %s", err)
		return nil, err
	}

	var results []*chaostest.ChaosTest

	for _, item := range chaosEngineList.Items {
		var experimentResults []*chaostest.ExperimentResult
		for _, exp := range item.Spec.Experiments {
			rname := item.Name + "-" + exp.Name
			er := &chaostest.ExperimentResult{
				Experiment: exp.Name,
			}

			opts := metav1.GetOptions{}
			res, err := l.chaosResultImpl.Get(rname, opts)
			if err != nil && !apierrors.IsNotFound(err) {
				l.Log("get: failed to get %s: %s", rname, err)
				return nil, err
			}
			if res != nil {
				er.Result = res.Status.ExperimentStatus.Verdict
				er.Phase = res.Status.ExperimentStatus.Phase
			}
			experimentResults = append(experimentResults, er)
		}

		ct := &chaostest.ChaosTest{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Info:        &chaostest.Info{},
			Experiments: experimentResults,
		}
		results = append(results, ct)
	}
	return results, nil
}

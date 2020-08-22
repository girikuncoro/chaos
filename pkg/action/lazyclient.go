package action

import (
	"sync"

	v1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	litmuschaos "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned"
	chaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/typed/litmuschaos/v1alpha1"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// lazyClient is a workaround because litmus chaos CRD
// is a separate clientset outside of kubernetes/client-go,
// need to find a way to streamline this.
type lazyClient struct {
	initClient sync.Once
	client     litmuschaos.Interface
	clientErr  error

	clientFn  func() (*litmuschaos.Clientset, error)
	namespace string
}

func (s *lazyClient) init() error {
	s.initClient.Do(func() {
		s.client, s.clientErr = s.clientFn()
	})
	return s.clientErr
}

// chaosEngineClient implements chaosv1alpha1.ChaosEngineInterface
type chaosEngineClient struct{ *lazyClient }

var _ chaosv1alpha1.ChaosEngineInterface

func newChaosEngineClient(lc *lazyClient) *chaosEngineClient {
	return &chaosEngineClient{lazyClient: lc}
}

func (ce *chaosEngineClient) Create(chaosEngine *v1alpha1.ChaosEngine) (*v1alpha1.ChaosEngine, error) {
	if err := ce.init(); err != nil {
		return nil, err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).Create(chaosEngine)
}

func (ce *chaosEngineClient) Update(chaosEngine *v1alpha1.ChaosEngine) (*v1alpha1.ChaosEngine, error) {
	if err := ce.init(); err != nil {
		return nil, err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).Update(chaosEngine)
}

func (ce *chaosEngineClient) UpdateStatus(chaosEngine *v1alpha1.ChaosEngine) (*v1alpha1.ChaosEngine, error) {
	if err := ce.init(); err != nil {
		return nil, err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).UpdateStatus(chaosEngine)
}

func (ce *chaosEngineClient) Delete(name string, opts *v1.DeleteOptions) error {
	if err := ce.init(); err != nil {
		return err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).Delete(name, opts)
}

func (ce *chaosEngineClient) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return errors.New("function is not implemented yet")
}

func (ce *chaosEngineClient) Get(name string, opts v1.GetOptions) (*v1alpha1.ChaosEngine, error) {
	if err := ce.init(); err != nil {
		return nil, err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).Get(name, opts)
}

func (ce *chaosEngineClient) List(opts v1.ListOptions) (*v1alpha1.ChaosEngineList, error) {
	if err := ce.init(); err != nil {
		return nil, err
	}
	return ce.client.LitmuschaosV1alpha1().ChaosEngines(ce.namespace).List(opts)
}

func (ce *chaosEngineClient) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("function is not implemented yet")
}

func (ce *chaosEngineClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ChaosEngine, err error) {
	return nil, errors.New("function is not implemented yet")
}

// chaosResultClient implements chaosv1alpha1.ChaosResultInterface
type chaosResultClient struct{ *lazyClient }

var _ chaosv1alpha1.ChaosResultInterface

func newChaosResultClient(lc *lazyClient) *chaosResultClient {
	return &chaosResultClient{lazyClient: lc}
}

func (cr *chaosResultClient) Create(chaosResult *v1alpha1.ChaosResult) (*v1alpha1.ChaosResult, error) {
	return nil, errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) Update(chaosResult *v1alpha1.ChaosResult) (*v1alpha1.ChaosResult, error) {
	return nil, errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) UpdateStatus(chaosResult *v1alpha1.ChaosResult) (*v1alpha1.ChaosResult, error) {
	return nil, errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) Delete(name string, opts *v1.DeleteOptions) error {
	return errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) Get(name string, opts v1.GetOptions) (*v1alpha1.ChaosResult, error) {
	if err := cr.init(); err != nil {
		return nil, err
	}
	return cr.client.LitmuschaosV1alpha1().ChaosResults(cr.namespace).Get(name, opts)
}

func (cr *chaosResultClient) List(opts v1.ListOptions) (*v1alpha1.ChaosResultList, error) {
	return nil, errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("function is not implemented yet")
}

func (cr *chaosResultClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ChaosResult, err error) {
	return nil, errors.New("function is not implemented yet")
}

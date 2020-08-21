package driver

import (
	"github.com/girikuncoro/chaos/pkg/chaostest"
	"github.com/pkg/errors"
)

var (
	// ErrChaosNotFound indicates that a chaostest is not found
	ErrChaosTestNotFound = errors.New("chaostest: not found")
	// ErrChaosTestExists indicates that a release already exists.
	ErrChaosTestExists = errors.New("chaostest: already exists")
)

// Queryor is the interface that wraps the List method.
type Queryor interface {
	List() ([]*chaostest.ChaosTest, error)
}

type Driver interface {
	Queryor
	Name() string
}

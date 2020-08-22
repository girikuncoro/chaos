package storage

import (
	"github.com/girikuncoro/chaos/pkg/chaostest"
	"github.com/girikuncoro/chaos/pkg/storage/driver"
)

// Storage represents a storage engine for a ChaosTest.
type Storage struct {
	driver.Driver
	Log func(string, ...interface{})
}

// ListChaosTests returns all chaos tests from storage.
func (s *Storage) ListChaosTests() ([]*chaostest.ChaosTest, error) {
	s.Log("listing all chaos tests in storage")
	return s.Driver.List()
}

// Init initializes a new storage backend with the driver d.
func Init(d driver.Driver) *Storage {
	return &Storage{
		Driver: d,
		Log:    func(_ string, _ ...interface{}) {},
	}
}

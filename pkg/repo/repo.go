package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"sigs.k8s.io/yaml"
)

// File represents the repositories.yaml file.
type File struct {
	Generated    time.Time `json:"generated"`
	Repositories []*Entry  `json:"repositories"`
}

// NewFile generates an empty repositories file.
func NewFile() *File {
	return &File{
		Generated:    time.Now(),
		Repositories: []*Entry{},
	}
}

// Add adds one or more repo entries to repo file.
func (r *File) Add(re ...*Entry) {
	r.Repositories = append(r.Repositories, re...)
}

// Has returns true if given name exists in repository.
func (r *File) Has(name string) bool {
	entry := r.Get(name)
	return entry != nil
}

// Get returns an entry with given name if exists, returns nil otherwise.
func (r *File) Get(name string) *Entry {
	for _, entry := range r.Repositories {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

// WriteFile writes repositories file to give path.
func (r *File) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}

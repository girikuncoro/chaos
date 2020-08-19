package repo

import (
	"net/url"

	"github.com/pkg/errors"
)

// Entry represents a collection of parameters for experiment chart repository
type Entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ChartRepository represents a chart repository
type ChartRepository struct {
	Config *Entry
}

func NewChartRepository(cfg *Entry) (*ChartRepository, error) {
	_, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, errors.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	// TODO: Init client to get the experiment chart from repo url

	return &ChartRepository{
		Config: cfg,
	}, nil
}

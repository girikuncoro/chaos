package repo

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/girikuncoro/chaos/pkg/chaospath"
	"github.com/girikuncoro/chaos/pkg/getter"
	"github.com/pkg/errors"
)

// Entry represents a collection of parameters for experiment chart repository
type Entry struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	ExperimentFile string `json:"experimentFile"`
}

// ChartRepository represents a chart repository
type ChartRepository struct {
	Config         *Entry
	ExperimentFile string
	Client         getter.Getter
	CachePath      string
}

// NewChartRepository constructs ChartRepository
func NewChartRepository(cfg *Entry, client getter.Getter) (*ChartRepository, error) {
	_, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, errors.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	// TODO: Init client to get the experiment chart from repo url
	return &ChartRepository{
		Config:    cfg,
		Client:    client,
		CachePath: chaospath.CachePath("repository"),
	}, nil
}

// DownloadExperimentFile fetches the experiments from a repository.
// Currently the url expects experiments.yaml, should be pass like
// e.g. https://hub.litmuschaos.io/api/chaos/1.7.0?file=charts/generic
//
// TODO(giri): Figure out how to index the experiments
func (r *ChartRepository) DownloadExperimentFile() (string, error) {
	parsedURL, err := url.Parse(r.Config.URL)
	if err != nil {
		return "", err
	}
	parsedURL.RawQuery = path.Join(parsedURL.RawQuery, "experiments.yaml")

	experimentURL := parsedURL.String()
	resp, err := r.Client.Get(experimentURL,
		getter.WithURL(r.Config.URL),
	)
	if err != nil {
		return "", err
	}

	exp, err := ioutil.ReadAll(resp)
	if err != nil {
		return "", err
	}

	expFile, err := loadExperiment(exp)
	if err != nil {
		return "", err
	}

	// TODO: Create experiment list in cache directory

	// Create the experiment file in the cache directory
	fname := filepath.Join(r.CachePath, r.Config.Name+".yaml")
	os.MkdirAll(filepath.Dir(fname), 0755)
	return fname, ioutil.WriteFile(fname, expFile, 0644)
}

// loadExperiment loads an experiment file and does minimal sanity check.
//
// TODO: currently only checks it's empty or not.
func loadExperiment(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, errors.New("no valid experiment file found")
	}
	return data, nil
}

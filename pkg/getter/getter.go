package getter

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type options struct {
	url     string
	timeout time.Duration
}

// Option allows specifcying configurations when performing
// Get operations.
type Option func(*options)

// WithURL informs the server name will be used for fetching objects.
func WithURL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

// WithTimeout sets the timeout for requests
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// Getter is an interface to support GET to the specified URL.
type Getter interface {
	// Get file content by url string
	Get(url string, options ...Option) (*bytes.Buffer, error)
}

// HTTPGetter is the HTTP(/S) backend handler
type HTTPGetter struct {
	opts options
}

// Get performs a Get from repo.Getter and returns the body.
func (g *HTTPGetter) Get(href string, options ...Option) (*bytes.Buffer, error) {
	for _, opt := range options {
		opt(&g.opts)
	}
	return g.get(href)
}

func (g *HTTPGetter) get(href string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return buf, err
	}

	client := g.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return buf, err
	}
	if resp.StatusCode != 200 {
		return buf, errors.Errorf("failed to fetch %s : %s", href, resp.Status)
	}

	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	return buf, err
}

// NewHTTPGetter constructs http/https client as a Getter
func NewHTTPGetter(options ...Option) (Getter, error) {
	var client HTTPGetter

	for _, opt := range options {
		opt(&client.opts)
	}

	return &client, nil
}

func (g *HTTPGetter) httpClient() *http.Client {
	// TODO(giri): Shouldn't hardcode to non-TLS
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: g.opts.timeout,
	}
	return client
}

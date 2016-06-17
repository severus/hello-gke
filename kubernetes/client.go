package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/severusio/hello-gke/kubernetes/api"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	// APIEndpoint defines the base path for Kubernetes API resources.
	APIEndpoint = "/api/v1"
	defaultPod  = "/namespaces/default/pods"
)

// Config represents a Kubernetes API client configuration.
type Config struct {
	BaseURL  string
	Username string
	Password string
}

// Client is a client for the Kubernetes master.
type Client struct {
	endpointURL string
	httpClient  *http.Client
}

func httpClient(cfg *Config) *http.Client {
	return &http.Client{
		Transport: NewAuthTransport(cfg),
	}
}

// NewClient returns a new Kubernetes client.
func NewClient(cfg *Config) (*Client, error) {
	validURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %v", cfg.BaseURL, err)
	}

	return &Client{
		endpointURL: strings.TrimSuffix(validURL.String(), "/") + APIEndpoint,
		httpClient:  httpClient(cfg),
	}, nil
}

// GetPods returns all pods in the cluster, regardless of status.
func (c *Client) GetPods(ctx context.Context) ([]api.Pod, error) {
	getURL := c.endpointURL + defaultPod

	// Make request to Kubernetes API
	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: GET %q : %v", getURL, err)
	}
	resp, err := ctxhttp.Do(ctx, c.httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: GET %q: %v", getURL, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read request body for GET %q: %v", getURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error %d GET %q: %q: %v", resp.StatusCode, getURL, string(body), err)
	}

	var podList api.PodList
	if err := json.Unmarshal(body, &podList); err != nil {
		return nil, fmt.Errorf("failed to decode list of pod resources: %v", err)
	}
	return podList.Items, nil
}

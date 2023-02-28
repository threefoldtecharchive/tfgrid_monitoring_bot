package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// GridProxyClient struct
type GridProxyClient struct {
	endpoint string
}

// NewProxyBus generates a new proxy bus
func NewGridProxyClient(endpoint string) (*GridProxyClient, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	return &GridProxyClient{
		endpoint,
	}, nil
}

func (r GridProxyClient) Ping() error {
	pingURL := fmt.Sprintf("%s/ping", r.endpoint)

	resp, err := http.Get(pingURL)
	if err != nil {
		return errors.Wrapf(err, "failed to send GET request to %s", pingURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("endpoint %s returned %d status code", pingURL, resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for endpoint %s", pingURL)
	}
	if !strings.Contains(string(bodyBytes), "pong") {
		return fmt.Errorf("proxy server response does not have pong: '%s'", string(bodyBytes))
	}
	return nil
}

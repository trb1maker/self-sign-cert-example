package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/trb1maker/self-sign-cert-example/internal/http/creds"
	"github.com/trb1maker/self-sign-cert-example/internal/http/dto"
)

var errFailedAddCert = errors.New("failed to add cert to pool")

func NewClient(port uint16, certPath string) (*Client, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %w", err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(cert) {
		return nil, errFailedAddCert
	}

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool}, //nolint:gosec
			},
		},
		port: port,
	}, nil
}

type Client struct {
	client *http.Client
	port   uint16
}

func (c *Client) Time(ctx context.Context) (time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://localhost:%d/now", c.port), nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(creds.User, creds.Password)

	resp, err := c.client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get time: %w", err)
	}
	defer resp.Body.Close()

	ts := new(dto.TimeResponse) //nolint:varnamelen

	if err := json.NewDecoder(resp.Body).Decode(ts); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return ts.Time(), nil
}

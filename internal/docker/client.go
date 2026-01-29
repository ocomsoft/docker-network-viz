// Package docker provides a wrapper around the Docker SDK client for fetching
// network and container information. It handles client initialization with
// proper error handling.
package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

// Client wraps the Docker SDK client with additional functionality
// for network topology visualization. It provides methods for fetching
// networks, containers, and building network-to-container mappings.
type Client struct {
	// cli is the underlying Docker SDK client.
	cli client.APIClient
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithDockerClient sets a custom Docker API client.
// This is primarily used for testing with mock clients.
func WithDockerClient(cli client.APIClient) ClientOption {
	return func(c *Client) {
		c.cli = cli
	}
}

// NewClient creates a new Docker client wrapper with the given options.
// It initializes the Docker SDK client using environment configuration
// and API version negotiation.
//
// The client can be configured with the following options:
//   - WithDockerClient: Inject a custom Docker API client (useful for testing)
//
// Returns an error if the Docker client cannot be initialized.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{}

	// Apply options first to allow WithDockerClient to be used
	for _, opt := range opts {
		opt(c)
	}

	// If no custom client was provided, create one from environment
	if c.cli == nil {
		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create Docker client: %w", err)
		}

		c.cli = cli
	}

	return c, nil
}

// Ping checks if the Docker daemon is accessible.
// This is useful for verifying connectivity before performing operations.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping Docker daemon: %w", err)
	}

	return nil
}

// Close closes the underlying Docker client connection.
// It should be called when the client is no longer needed.
func (c *Client) Close() error {
	return c.cli.Close()
}

// APIClient returns the underlying Docker API client.
// This can be used for operations not covered by the wrapper.
func (c *Client) APIClient() client.APIClient {
	return c.cli
}

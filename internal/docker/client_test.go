// Package docker provides tests for the Docker client wrapper.
package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// mockAPIClient is a mock implementation of the Docker API client for testing.
type mockAPIClient struct {
	client.APIClient

	// Mock function implementations
	pingFunc             func(ctx context.Context) (types.Ping, error)
	closeFunc            func() error
	networkListFunc      func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error)
	networkInspectFunc   func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error)
	containerListFunc    func(ctx context.Context, opts container.ListOptions) ([]types.Container, error)
	containerInspectFunc func(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

// Ping implements the Ping method of the Docker API client.
func (m *mockAPIClient) Ping(ctx context.Context) (types.Ping, error) {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return types.Ping{}, nil
}

// Close implements the Close method of the Docker API client.
func (m *mockAPIClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

// NetworkList implements the NetworkList method of the Docker API client.
func (m *mockAPIClient) NetworkList(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
	if m.networkListFunc != nil {
		return m.networkListFunc(ctx, opts)
	}
	return nil, nil
}

// NetworkInspect implements the NetworkInspect method of the Docker API client.
func (m *mockAPIClient) NetworkInspect(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
	if m.networkInspectFunc != nil {
		return m.networkInspectFunc(ctx, networkID, opts)
	}
	return network.Inspect{}, nil
}

// ContainerList implements the ContainerList method of the Docker API client.
func (m *mockAPIClient) ContainerList(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
	if m.containerListFunc != nil {
		return m.containerListFunc(ctx, opts)
	}
	return nil, nil
}

// ContainerInspect implements the ContainerInspect method of the Docker API client.
func (m *mockAPIClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	if m.containerInspectFunc != nil {
		return m.containerInspectFunc(ctx, containerID)
	}
	return types.ContainerJSON{}, nil
}

// TestNewClient_WithMockClient tests client creation with a mock Docker client.
func TestNewClient_WithMockClient(t *testing.T) {
	mock := &mockAPIClient{}

	c, err := NewClient(
		WithDockerClient(mock),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c == nil {
		t.Fatal("expected client, got nil")
	}

	if c.APIClient() != mock {
		t.Error("expected mock client to be set")
	}
}

// TestClient_Ping_Success tests successful ping to Docker daemon.
func TestClient_Ping_Success(t *testing.T) {
	mock := &mockAPIClient{
		pingFunc: func(ctx context.Context) (types.Ping, error) {
			return types.Ping{APIVersion: "1.41"}, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = c.Ping(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestClient_Ping_Failure tests failed ping to Docker daemon.
func TestClient_Ping_Failure(t *testing.T) {
	expectedErr := errors.New("docker daemon not accessible")
	mock := &mockAPIClient{
		pingFunc: func(ctx context.Context) (types.Ping, error) {
			return types.Ping{}, expectedErr
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = c.Ping(context.Background())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// TestClient_Close tests closing the Docker client.
func TestClient_Close(t *testing.T) {
	closeCalled := false
	mock := &mockAPIClient{
		closeFunc: func() error {
			closeCalled = true
			return nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = c.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !closeCalled {
		t.Error("expected Close to be called on mock")
	}
}

// TestClient_APIClient tests the APIClient getter.
func TestClient_APIClient(t *testing.T) {
	mock := &mockAPIClient{}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if c.APIClient() != mock {
		t.Error("expected APIClient to return the mock")
	}
}

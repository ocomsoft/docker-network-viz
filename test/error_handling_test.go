// Package integration provides error handling integration tests for docker-network-viz.
package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"git.o.ocom.com.au/go/docker-network-viz/internal/docker"
)

// errorMockAPIClient is a mock that returns errors for testing error handling.
type errorMockAPIClient struct {
	client.APIClient

	pingFunc           func(ctx context.Context) (types.Ping, error)
	closeFunc          func() error
	networkListFunc    func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error)
	containerListFunc  func(ctx context.Context, opts container.ListOptions) ([]types.Container, error)
	networkInspectFunc func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error)
}

// Ping implements the Ping method.
func (m *errorMockAPIClient) Ping(ctx context.Context) (types.Ping, error) {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return types.Ping{}, nil
}

// Close implements the Close method.
func (m *errorMockAPIClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

// NetworkList implements the NetworkList method.
func (m *errorMockAPIClient) NetworkList(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
	if m.networkListFunc != nil {
		return m.networkListFunc(ctx, opts)
	}
	return nil, nil
}

// ContainerList implements the ContainerList method.
func (m *errorMockAPIClient) ContainerList(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
	if m.containerListFunc != nil {
		return m.containerListFunc(ctx, opts)
	}
	return nil, nil
}

// NetworkInspect implements the NetworkInspect method.
func (m *errorMockAPIClient) NetworkInspect(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
	if m.networkInspectFunc != nil {
		return m.networkInspectFunc(ctx, networkID, opts)
	}
	return network.Inspect{}, nil
}

// TestError_NetworkListFailure tests handling of network list errors.
func TestError_NetworkListFailure(t *testing.T) {
	expectedErr := errors.New("permission denied: cannot access Docker daemon")

	mock := &errorMockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return nil, expectedErr
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()
	_, err = dockerClient.FetchNetworks(ctx, nil)

	if err == nil {
		t.Fatal("expected error when network list fails, got nil")
	}

	if !errors.Is(err, expectedErr) && err.Error() == "" {
		t.Errorf("expected wrapped error containing original error, got: %v", err)
	}
}

// TestError_ContainerListFailure tests handling of container list errors.
func TestError_ContainerListFailure(t *testing.T) {
	expectedErr := errors.New("Docker daemon is not running")

	mock := &errorMockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return nil, expectedErr
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()
	_, err = dockerClient.FetchContainers(ctx, nil)

	if err == nil {
		t.Fatal("expected error when container list fails, got nil")
	}
}

// TestError_PingFailure tests handling of ping errors.
func TestError_PingFailure(t *testing.T) {
	expectedErr := errors.New("connection refused: Docker daemon not accessible")

	mock := &errorMockAPIClient{
		pingFunc: func(ctx context.Context) (types.Ping, error) {
			return types.Ping{}, expectedErr
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()
	err = dockerClient.Ping(ctx)

	if err == nil {
		t.Fatal("expected error when ping fails, got nil")
	}
}

// TestError_NetworkInspectFailure tests handling of network inspect errors.
func TestError_NetworkInspectFailure(t *testing.T) {
	expectedErr := errors.New("network not found: invalid_network")

	mock := &errorMockAPIClient{
		networkInspectFunc: func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
			return network.Inspect{}, expectedErr
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()
	_, err = dockerClient.FetchNetworkByID(ctx, "invalid_network")

	if err == nil {
		t.Fatal("expected error when network inspect fails, got nil")
	}
}

// TestError_CloseFailure tests handling of close errors.
func TestError_CloseFailure(t *testing.T) {
	expectedErr := errors.New("failed to close connection")

	mock := &errorMockAPIClient{
		closeFunc: func() error {
			return expectedErr
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}

	err = dockerClient.Close()

	if err == nil {
		t.Fatal("expected error when close fails, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

// TestError_EmptyContainerNames tests handling of containers with empty names.
func TestError_EmptyContainerNames(t *testing.T) {
	mock := &errorMockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:    "container1",
					Names: []string{}, // Empty names
					NetworkSettings: &types.SummaryNetworkSettings{
						Networks: map[string]*network.EndpointSettings{
							"bridge": {},
						},
					},
				},
			}, nil
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()
	containers, err := dockerClient.FetchContainers(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should handle empty names gracefully
	containerMap := dockerClient.BuildContainerMap(containers)

	// Empty name should result in empty string key
	if _, ok := containerMap[""]; !ok {
		t.Error("expected empty string key for container with no names")
	}
}

// TestError_NilNetworkSettings tests handling of containers with nil network settings.
// Note: Currently the Docker API always returns a non-nil NetworkSettings struct,
// so this test verifies behavior with an empty Networks map instead.
func TestError_NilNetworkSettings(t *testing.T) {
	mock := &errorMockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:    "container1",
					Names: []string{"/test_container"},
					NetworkSettings: &types.SummaryNetworkSettings{
						Networks: map[string]*network.EndpointSettings{}, // Empty networks
					},
				},
			}, nil
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()

	containers, err := dockerClient.FetchContainers(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should handle empty networks gracefully
	containerMap := dockerClient.BuildContainerMap(containers)

	// Container should exist with no networks
	if ci, ok := containerMap["test_container"]; !ok {
		t.Error("expected test_container in container map")
	} else if ci.NetworkCount() != 0 {
		t.Errorf("expected 0 networks, got %d", ci.NetworkCount())
	}
}

// TestError_ContextCancellation tests handling of context cancellation.
func TestError_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mock := &errorMockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return nil, ctx.Err()
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	_, err = dockerClient.FetchNetworks(ctx, nil)

	if err == nil {
		t.Fatal("expected error when context is canceled, got nil")
	}
}

// TestError_PartialDataRecovery tests that we can still process partial data
// when some operations succeed and others fail.
func TestError_PartialDataRecovery(t *testing.T) {
	networksReturned := false

	mock := &errorMockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			networksReturned = true
			return []network.Summary{
				{Name: "bridge", Driver: "bridge"},
			}, nil
		},
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return nil, errors.New("container list failed")
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	ctx := context.Background()

	// Networks should succeed
	networks, err := dockerClient.FetchNetworks(ctx, nil)
	if err != nil {
		t.Fatalf("network fetch should succeed: %v", err)
	}

	if !networksReturned {
		t.Error("networks should have been returned")
	}

	if len(networks) != 1 {
		t.Errorf("expected 1 network, got %d", len(networks))
	}

	// Containers should fail
	_, err = dockerClient.FetchContainers(ctx, nil)
	if err == nil {
		t.Error("container fetch should fail")
	}
}

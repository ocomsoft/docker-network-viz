// Package docker provides tests for the Docker network wrapper.
package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types/network"
)

// TestClient_FetchNetworks_Success tests successful network listing.
func TestClient_FetchNetworks_Success(t *testing.T) {
	expectedNetworks := []network.Summary{
		{ID: "net1", Name: "bridge", Driver: "bridge"},
		{ID: "net2", Name: "host", Driver: "host"},
		{ID: "net3", Name: "custom_net", Driver: "bridge"},
	}

	mock := &mockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return expectedNetworks, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	networks, err := c.FetchNetworks(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that networks are sorted alphabetically by name
	if len(networks) != 3 {
		t.Fatalf("expected 3 networks, got %d", len(networks))
	}

	// Verify sorting order
	if networks[0].Name != "bridge" {
		t.Errorf("expected first network to be 'bridge', got '%s'", networks[0].Name)
	}
	if networks[1].Name != "custom_net" {
		t.Errorf("expected second network to be 'custom_net', got '%s'", networks[1].Name)
	}
	if networks[2].Name != "host" {
		t.Errorf("expected third network to be 'host', got '%s'", networks[2].Name)
	}
}

// TestClient_FetchNetworks_Empty tests fetching an empty network list.
func TestClient_FetchNetworks_Empty(t *testing.T) {
	mock := &mockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return []network.Summary{}, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	networks, err := c.FetchNetworks(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(networks) != 0 {
		t.Errorf("expected 0 networks, got %d", len(networks))
	}
}

// TestClient_FetchNetworks_Error tests error handling when listing networks fails.
func TestClient_FetchNetworks_Error(t *testing.T) {
	expectedErr := errors.New("failed to list networks")
	mock := &mockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return nil, expectedErr
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	networks, err := c.FetchNetworks(context.Background(), nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if networks != nil {
		t.Error("expected nil networks on error")
	}
}

// TestClient_FetchNetworkByID_Success tests successful network inspection by ID.
func TestClient_FetchNetworkByID_Success(t *testing.T) {
	expectedNet := network.Inspect{
		Name:   "test_network",
		ID:     "abc123",
		Driver: "bridge",
		Containers: map[string]network.EndpointResource{
			"container1": {Name: "web"},
			"container2": {Name: "db"},
		},
	}

	mock := &mockAPIClient{
		networkInspectFunc: func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
			if networkID != "abc123" {
				t.Errorf("unexpected network ID: %s", networkID)
			}
			return expectedNet, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	net, err := c.FetchNetworkByID(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if net.Name != "test_network" {
		t.Errorf("expected network name 'test_network', got '%s'", net.Name)
	}

	if net.Driver != "bridge" {
		t.Errorf("expected driver 'bridge', got '%s'", net.Driver)
	}

	if len(net.Containers) != 2 {
		t.Errorf("expected 2 containers, got %d", len(net.Containers))
	}
}

// TestClient_FetchNetworkByID_Error tests error handling when network inspection fails.
func TestClient_FetchNetworkByID_Error(t *testing.T) {
	mock := &mockAPIClient{
		networkInspectFunc: func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
			return network.Inspect{}, errors.New("network not found")
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = c.FetchNetworkByID(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// TestClient_FetchNetworkByName_Success tests successful network inspection by name.
func TestClient_FetchNetworkByName_Success(t *testing.T) {
	expectedNet := network.Inspect{
		Name:   "my_network",
		ID:     "def456",
		Driver: "overlay",
	}

	mock := &mockAPIClient{
		networkInspectFunc: func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
			if networkID != "my_network" {
				t.Errorf("unexpected network name: %s", networkID)
			}
			return expectedNet, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	net, err := c.FetchNetworkByName(context.Background(), "my_network")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if net.Name != "my_network" {
		t.Errorf("expected network name 'my_network', got '%s'", net.Name)
	}
}

// TestClient_FetchNetworkByName_Error tests error handling when network inspection by name fails.
func TestClient_FetchNetworkByName_Error(t *testing.T) {
	mock := &mockAPIClient{
		networkInspectFunc: func(ctx context.Context, networkID string, opts network.InspectOptions) (network.Inspect, error) {
			return network.Inspect{}, errors.New("network not found")
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = c.FetchNetworkByName(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// TestConvertToNetworkInfo tests conversion of Docker network summary to internal model.
func TestConvertToNetworkInfo(t *testing.T) {
	summary := network.Summary{
		Name:   "test_net",
		Driver: "bridge",
	}

	info := ConvertToNetworkInfo(summary)

	if info.Name != "test_net" {
		t.Errorf("expected name 'test_net', got '%s'", info.Name)
	}

	if info.Driver != "bridge" {
		t.Errorf("expected driver 'bridge', got '%s'", info.Driver)
	}
}

// TestConvertNetworksToNetworkInfos tests bulk conversion of network summaries.
func TestConvertNetworksToNetworkInfos(t *testing.T) {
	summaries := []network.Summary{
		{Name: "net1", Driver: "bridge"},
		{Name: "net2", Driver: "host"},
		{Name: "net3", Driver: "overlay"},
	}

	infos := ConvertNetworksToNetworkInfos(summaries)

	if len(infos) != 3 {
		t.Fatalf("expected 3 infos, got %d", len(infos))
	}

	for i, info := range infos {
		if info.Name != summaries[i].Name {
			t.Errorf("info[%d]: expected name '%s', got '%s'", i, summaries[i].Name, info.Name)
		}
		if info.Driver != summaries[i].Driver {
			t.Errorf("info[%d]: expected driver '%s', got '%s'", i, summaries[i].Driver, info.Driver)
		}
	}
}

// TestConvertNetworksToNetworkInfos_Empty tests conversion of empty network list.
func TestConvertNetworksToNetworkInfos_Empty(t *testing.T) {
	summaries := []network.Summary{}

	infos := ConvertNetworksToNetworkInfos(summaries)

	if len(infos) != 0 {
		t.Errorf("expected 0 infos, got %d", len(infos))
	}
}

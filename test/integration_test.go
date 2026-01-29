// Package integration provides integration tests for docker-network-viz.
// These tests verify end-to-end functionality using mock Docker data.
package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"git.o.ocom.com.au/go/docker-network-viz/internal/docker"
	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
	"git.o.ocom.com.au/go/docker-network-viz/internal/output"
)

// mockAPIClient is a mock implementation of the Docker API client for testing.
type mockAPIClient struct {
	client.APIClient

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

// createMockNetworks creates mock network data for testing.
func createMockNetworks() []network.Summary {
	return []network.Summary{
		{Name: "bridge", Driver: "bridge", ID: "net1"},
		{Name: "frontend_net", Driver: "bridge", ID: "net2"},
		{Name: "backend_net", Driver: "bridge", ID: "net3"},
	}
}

// createMockContainers creates mock container data for testing.
func createMockContainers() []types.Container {
	return []types.Container{
		{
			ID:    "container1",
			Names: []string{"/web_app"},
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"frontend_net": {
						Aliases: []string{"web", "web.local"},
					},
				},
			},
		},
		{
			ID:    "container2",
			Names: []string{"/api"},
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"frontend_net": {
						Aliases: []string{"api"},
					},
					"backend_net": {
						Aliases: []string{"api-backend"},
					},
				},
			},
		},
		{
			ID:    "container3",
			Names: []string{"/postgres"},
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"backend_net": {
						Aliases: []string{"db", "postgres.local"},
					},
				},
			},
		},
		{
			ID:    "container4",
			Names: []string{"/redis"},
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"backend_net": {
						Aliases: []string{"cache"},
					},
				},
			},
		},
	}
}

// TestIntegration_FetchAndVisualize tests the full flow of fetching Docker data
// and generating visualization output.
func TestIntegration_FetchAndVisualize(t *testing.T) {
	ctx := context.Background()

	mockNetworks := createMockNetworks()
	mockContainers := createMockContainers()

	mock := &mockAPIClient{
		networkListFunc: func(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
			return mockNetworks, nil
		},
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return mockContainers, nil
		},
	}

	// Create client with mock
	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	// Fetch networks
	networks, err := dockerClient.FetchNetworks(ctx, nil)
	if err != nil {
		t.Fatalf("failed to fetch networks: %v", err)
	}

	if len(networks) != 3 {
		t.Errorf("expected 3 networks, got %d", len(networks))
	}

	// Fetch containers
	containers, err := dockerClient.FetchContainers(ctx, &docker.ContainerListOptions{All: true})
	if err != nil {
		t.Fatalf("failed to fetch containers: %v", err)
	}

	if len(containers) != 4 {
		t.Errorf("expected 4 containers, got %d", len(containers))
	}

	// Build mappings
	networkToContainers := dockerClient.BuildNetworkToContainersMap(containers)

	// Verify network mappings
	if len(networkToContainers["frontend_net"]) != 2 {
		t.Errorf("expected 2 containers on frontend_net, got %d", len(networkToContainers["frontend_net"]))
	}

	if len(networkToContainers["backend_net"]) != 3 {
		t.Errorf("expected 3 containers on backend_net, got %d", len(networkToContainers["backend_net"]))
	}
}

// TestIntegration_NetworkTreeOutput tests that network tree output is formatted correctly.
func TestIntegration_NetworkTreeOutput(t *testing.T) {
	netInfo := models.NewNetworkInfo("frontend_net", "bridge")

	containers := []models.ContainerInfo{
		{
			Name:     "api",
			Aliases:  []string{"api"},
			Networks: []string{"frontend_net", "backend_net"},
		},
		{
			Name:     "web_app",
			Aliases:  []string{"web", "web.local"},
			Networks: []string{"frontend_net"},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Verify the output contains expected elements
	expectedElements := []string{
		"Network:",
		"frontend_net",
		"bridge",
		"api",
		"web_app",
		"alias:",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, result)
		}
	}

	// Verify tree structure symbols are present
	treeSymbols := []string{"├──", "└──", "│"}
	foundAnySymbol := false
	for _, symbol := range treeSymbols {
		if strings.Contains(result, symbol) {
			foundAnySymbol = true
			break
		}
	}
	if !foundAnySymbol {
		t.Errorf("expected output to contain tree symbols, got:\n%s", result)
	}
}

// TestIntegration_ContainerTreeOutput tests that container tree output is formatted correctly.
func TestIntegration_ContainerTreeOutput(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "api",
		Aliases:  []string{"api-service"},
		Networks: []string{"frontend_net", "backend_net"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"frontend_net": {
			{Name: "api", Networks: []string{"frontend_net", "backend_net"}},
			{Name: "web_app", Networks: []string{"frontend_net"}},
		},
		"backend_net": {
			{Name: "api", Networks: []string{"frontend_net", "backend_net"}},
			{Name: "postgres", Networks: []string{"backend_net"}},
			{Name: "redis", Networks: []string{"backend_net"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	// Verify the output contains expected elements
	expectedElements := []string{
		"Container:",
		"api",
		"Network:",
		"frontend_net",
		"backend_net",
		"connects to:",
		"web_app",
		"postgres",
		"redis",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, result)
		}
	}
}

// TestIntegration_ReachabilityCalculation tests that container reachability is calculated correctly.
func TestIntegration_ReachabilityCalculation(t *testing.T) {
	networkToContainers := map[string][]models.ContainerInfo{
		"frontend_net": {
			{Name: "api", Networks: []string{"frontend_net", "backend_net"}},
			{Name: "nginx", Networks: []string{"frontend_net"}},
			{Name: "web_app", Networks: []string{"frontend_net"}},
		},
		"backend_net": {
			{Name: "api", Networks: []string{"frontend_net", "backend_net"}},
			{Name: "postgres", Networks: []string{"backend_net"}},
			{Name: "redis", Networks: []string{"backend_net"}},
		},
	}

	// Test reachability from api on frontend_net
	reachable := output.ReachableContainers("api", "frontend_net", networkToContainers)

	if len(reachable) != 2 {
		t.Errorf("expected 2 reachable containers from api on frontend_net, got %d", len(reachable))
	}

	// Should contain nginx and web_app but not api itself
	expectedReachable := map[string]bool{"nginx": true, "web_app": true}
	for _, name := range reachable {
		if !expectedReachable[name] {
			t.Errorf("unexpected container %q in reachable list", name)
		}
	}

	// Test reachability from api on backend_net
	reachable = output.ReachableContainers("api", "backend_net", networkToContainers)

	if len(reachable) != 2 {
		t.Errorf("expected 2 reachable containers from api on backend_net, got %d", len(reachable))
	}

	// Should be sorted alphabetically
	if reachable[0] != "postgres" || reachable[1] != "redis" {
		t.Errorf("expected sorted reachable list [postgres, redis], got %v", reachable)
	}
}

// TestIntegration_EmptyNetwork tests handling of networks with no containers.
func TestIntegration_EmptyNetwork(t *testing.T) {
	netInfo := models.NewNetworkInfo("isolated_net", "bridge")
	containers := []models.ContainerInfo{}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	if !strings.Contains(result, "no containers") {
		t.Errorf("expected output to indicate no containers, got:\n%s", result)
	}
}

// TestIntegration_ContainerWithNoReachability tests container with no reachable peers.
func TestIntegration_ContainerWithNoReachability(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "isolated",
		Aliases:  []string{},
		Networks: []string{"private_net"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"private_net": {
			{Name: "isolated", Networks: []string{"private_net"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	if !strings.Contains(result, "(none)") {
		t.Errorf("expected output to indicate no reachable containers, got:\n%s", result)
	}
}

// TestIntegration_ContainerMapBuilding tests the container map building process.
func TestIntegration_ContainerMapBuilding(t *testing.T) {
	ctx := context.Background()
	mockContainers := createMockContainers()

	mock := &mockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return mockContainers, nil
		},
	}

	dockerClient, err := docker.NewClient(docker.WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	containers, err := dockerClient.FetchContainers(ctx, nil)
	if err != nil {
		t.Fatalf("failed to fetch containers: %v", err)
	}

	containerMap := dockerClient.BuildContainerMap(containers)

	// Verify container map contents
	if _, ok := containerMap["web_app"]; !ok {
		t.Error("expected web_app in container map")
	}

	if _, ok := containerMap["api"]; !ok {
		t.Error("expected api in container map")
	}

	// Verify api container has correct networks
	apiContainer := containerMap["api"]
	if !apiContainer.HasNetwork("frontend_net") {
		t.Error("expected api to have frontend_net")
	}
	if !apiContainer.HasNetwork("backend_net") {
		t.Error("expected api to have backend_net")
	}
}

// TestIntegration_SortedOutput tests that output is consistently sorted.
func TestIntegration_SortedOutput(t *testing.T) {
	netInfo := models.NewNetworkInfo("test_net", "bridge")

	// Containers in reverse alphabetical order
	containers := []models.ContainerInfo{
		{Name: "zebra", Aliases: []string{"z"}, Networks: []string{"test_net"}},
		{Name: "apple", Aliases: []string{"a"}, Networks: []string{"test_net"}},
		{Name: "mango", Aliases: []string{"m"}, Networks: []string{"test_net"}},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Find positions of container names in output
	applePos := strings.Index(result, "apple")
	mangoPos := strings.Index(result, "mango")
	zebraPos := strings.Index(result, "zebra")

	if applePos == -1 || mangoPos == -1 || zebraPos == -1 {
		t.Fatalf("expected all container names in output, got:\n%s", result)
	}

	// Verify alphabetical order
	if !(applePos < mangoPos && mangoPos < zebraPos) {
		t.Errorf("expected containers in alphabetical order, got:\n%s", result)
	}
}

// Package docker provides tests for the Docker container wrapper.
package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// createTestContainer creates a test container with proper types.
func createTestContainer(name string, networks map[string][]string) types.Container {
	netSettings := &types.SummaryNetworkSettings{
		Networks: make(map[string]*network.EndpointSettings),
	}

	for netName, aliases := range networks {
		netSettings.Networks[netName] = &network.EndpointSettings{
			Aliases: aliases,
		}
	}

	return types.Container{
		ID:              "id_" + name,
		Names:           []string{"/" + name},
		NetworkSettings: netSettings,
	}
}

// TestClient_FetchContainers_Success tests successful container listing.
func TestClient_FetchContainers_Success(t *testing.T) {
	expectedContainers := []types.Container{
		createTestContainer("web", map[string][]string{"bridge": {"web.local"}}),
		createTestContainer("api", map[string][]string{"bridge": {"api.local"}}),
		createTestContainer("db", map[string][]string{"backend": {"postgres"}}),
	}

	mock := &mockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			if !opts.All {
				t.Error("expected All option to be true")
			}
			return expectedContainers, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containers, err := c.FetchContainers(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that containers are sorted alphabetically by name
	if len(containers) != 3 {
		t.Fatalf("expected 3 containers, got %d", len(containers))
	}

	// Verify sorting order (api, db, web)
	names := []string{}
	for _, cont := range containers {
		names = append(names, sanitizeContainerName(cont.Names))
	}

	if names[0] != "api" {
		t.Errorf("expected first container to be 'api', got '%s'", names[0])
	}
	if names[1] != "db" {
		t.Errorf("expected second container to be 'db', got '%s'", names[1])
	}
	if names[2] != "web" {
		t.Errorf("expected third container to be 'web', got '%s'", names[2])
	}
}

// TestClient_FetchContainers_WithOptions tests container listing with custom options.
func TestClient_FetchContainers_WithOptions(t *testing.T) {
	mock := &mockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			if opts.All {
				t.Error("expected All option to be false")
			}
			return []types.Container{}, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opts := &ContainerListOptions{All: false}
	_, err = c.FetchContainers(context.Background(), opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestClient_FetchContainers_Empty tests fetching an empty container list.
func TestClient_FetchContainers_Empty(t *testing.T) {
	mock := &mockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return []types.Container{}, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containers, err := c.FetchContainers(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(containers) != 0 {
		t.Errorf("expected 0 containers, got %d", len(containers))
	}
}

// TestClient_FetchContainers_Error tests error handling when listing containers fails.
func TestClient_FetchContainers_Error(t *testing.T) {
	expectedErr := errors.New("failed to list containers")
	mock := &mockAPIClient{
		containerListFunc: func(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
			return nil, expectedErr
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containers, err := c.FetchContainers(context.Background(), nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if containers != nil {
		t.Error("expected nil containers on error")
	}
}

// TestClient_FetchContainerByID_Success tests successful container inspection by ID.
func TestClient_FetchContainerByID_Success(t *testing.T) {
	expectedContainer := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   "abc123",
			Name: "/test_container",
			State: &types.ContainerState{
				Status: "running",
			},
		},
	}

	mock := &mockAPIClient{
		containerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			if containerID != "abc123" {
				t.Errorf("unexpected container ID: %s", containerID)
			}
			return expectedContainer, nil
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containerJSON, err := c.FetchContainerByID(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if containerJSON.Name != "/test_container" {
		t.Errorf("expected container name '/test_container', got '%s'", containerJSON.Name)
	}

	if containerJSON.State.Status != "running" {
		t.Errorf("expected state 'running', got '%s'", containerJSON.State.Status)
	}
}

// TestClient_FetchContainerByID_Error tests error handling when container inspection fails.
func TestClient_FetchContainerByID_Error(t *testing.T) {
	mock := &mockAPIClient{
		containerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			return types.ContainerJSON{}, errors.New("container not found")
		},
	}

	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = c.FetchContainerByID(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// TestClient_BuildContainerMap tests building a container map.
func TestClient_BuildContainerMap(t *testing.T) {
	containers := []types.Container{
		createTestContainer("web", map[string][]string{
			"frontend": {"web.local", "www"},
			"backend":  {"web-internal"},
		}),
		createTestContainer("api", map[string][]string{
			"backend": {"api.local"},
		}),
	}

	mock := &mockAPIClient{}
	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containerMap := c.BuildContainerMap(containers)

	if len(containerMap) != 2 {
		t.Fatalf("expected 2 containers in map, got %d", len(containerMap))
	}

	// Check web container
	web, ok := containerMap["web"]
	if !ok {
		t.Fatal("expected 'web' in container map")
	}

	if web.Name != "web" {
		t.Errorf("expected name 'web', got '%s'", web.Name)
	}

	if !web.HasNetwork("frontend") {
		t.Error("expected web to have network 'frontend'")
	}

	if !web.HasNetwork("backend") {
		t.Error("expected web to have network 'backend'")
	}

	// Check api container
	api, ok := containerMap["api"]
	if !ok {
		t.Fatal("expected 'api' in container map")
	}

	if !api.HasNetwork("backend") {
		t.Error("expected api to have network 'backend'")
	}
}

// TestClient_BuildContainerMap_Empty tests building a container map from empty list.
func TestClient_BuildContainerMap_Empty(t *testing.T) {
	mock := &mockAPIClient{}
	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	containerMap := c.BuildContainerMap([]types.Container{})

	if len(containerMap) != 0 {
		t.Errorf("expected 0 containers in map, got %d", len(containerMap))
	}
}

// TestClient_BuildNetworkToContainersMap tests building network-to-containers mapping.
func TestClient_BuildNetworkToContainersMap(t *testing.T) {
	containers := []types.Container{
		createTestContainer("web", map[string][]string{
			"frontend": {"web.local"},
			"backend":  {},
		}),
		createTestContainer("api", map[string][]string{
			"backend": {"api.local"},
		}),
		createTestContainer("db", map[string][]string{
			"backend": {"postgres"},
		}),
	}

	mock := &mockAPIClient{}
	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	netMap := c.BuildNetworkToContainersMap(containers)

	// Check frontend network
	frontend, ok := netMap["frontend"]
	if !ok {
		t.Fatal("expected 'frontend' in network map")
	}

	if len(frontend) != 1 {
		t.Errorf("expected 1 container on frontend, got %d", len(frontend))
	}

	if frontend[0].Name != "web" {
		t.Errorf("expected 'web' on frontend, got '%s'", frontend[0].Name)
	}

	// Check backend network
	backend, ok := netMap["backend"]
	if !ok {
		t.Fatal("expected 'backend' in network map")
	}

	if len(backend) != 3 {
		t.Errorf("expected 3 containers on backend, got %d", len(backend))
	}

	// Verify containers are sorted alphabetically
	if backend[0].Name != "api" {
		t.Errorf("expected first backend container to be 'api', got '%s'", backend[0].Name)
	}
	if backend[1].Name != "db" {
		t.Errorf("expected second backend container to be 'db', got '%s'", backend[1].Name)
	}
	if backend[2].Name != "web" {
		t.Errorf("expected third backend container to be 'web', got '%s'", backend[2].Name)
	}
}

// TestClient_BuildNetworkToContainersMap_Empty tests building network map from empty list.
func TestClient_BuildNetworkToContainersMap_Empty(t *testing.T) {
	mock := &mockAPIClient{}
	c, err := NewClient(WithDockerClient(mock))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	netMap := c.BuildNetworkToContainersMap([]types.Container{})

	if len(netMap) != 0 {
		t.Errorf("expected 0 networks in map, got %d", len(netMap))
	}
}

// TestConvertToContainerInfo tests conversion of Docker container to internal model.
func TestConvertToContainerInfo(t *testing.T) {
	cont := createTestContainer("web", map[string][]string{
		"frontend": {"web.local", "www"},
		"backend":  {"web-internal"},
	})

	info := ConvertToContainerInfo(cont)

	if info.Name != "web" {
		t.Errorf("expected name 'web', got '%s'", info.Name)
	}

	if !info.HasNetwork("frontend") {
		t.Error("expected network 'frontend'")
	}

	if !info.HasNetwork("backend") {
		t.Error("expected network 'backend'")
	}

	if !info.HasAlias("web.local") {
		t.Error("expected alias 'web.local'")
	}

	if !info.HasAlias("www") {
		t.Error("expected alias 'www'")
	}

	if !info.HasAlias("web-internal") {
		t.Error("expected alias 'web-internal'")
	}
}

// TestConvertContainersToContainerInfos tests bulk conversion of containers.
func TestConvertContainersToContainerInfos(t *testing.T) {
	containers := []types.Container{
		createTestContainer("web", map[string][]string{"bridge": {}}),
		createTestContainer("api", map[string][]string{"bridge": {}}),
		createTestContainer("db", map[string][]string{"backend": {}}),
	}

	infos := ConvertContainersToContainerInfos(containers)

	if len(infos) != 3 {
		t.Fatalf("expected 3 infos, got %d", len(infos))
	}

	expectedNames := []string{"web", "api", "db"}
	for i, info := range infos {
		if info.Name != expectedNames[i] {
			t.Errorf("info[%d]: expected name '%s', got '%s'", i, expectedNames[i], info.Name)
		}
	}
}

// TestConvertContainersToContainerInfos_Empty tests conversion of empty container list.
func TestConvertContainersToContainerInfos_Empty(t *testing.T) {
	containers := []types.Container{}

	infos := ConvertContainersToContainerInfos(containers)

	if len(infos) != 0 {
		t.Errorf("expected 0 infos, got %d", len(infos))
	}
}

// TestSanitizeContainerName tests the sanitizeContainerName helper function.
func TestSanitizeContainerName(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"/web"}, "web"},
		{[]string{"/my_container"}, "my_container"},
		{[]string{"no_slash"}, "no_slash"},
		{[]string{}, ""},
		{[]string{"/container1", "/container2"}, "container1"},
	}

	for _, tc := range testCases {
		result := sanitizeContainerName(tc.input)
		if result != tc.expected {
			t.Errorf("sanitizeContainerName(%v) = '%s', expected '%s'", tc.input, result, tc.expected)
		}
	}
}

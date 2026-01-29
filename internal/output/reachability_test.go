package output

import (
	"testing"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

func TestReachableContainers_ReturnsOtherContainersOnSameNetwork(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"bridge": {
			{Name: "web", Networks: []string{"bridge"}},
			{Name: "api", Networks: []string{"bridge"}},
			{Name: "db", Networks: []string{"bridge"}},
		},
	}

	result := ReachableContainers("web", "bridge", netMap)

	expected := []string{"api", "db"}
	if len(result) != len(expected) {
		t.Errorf("expected %d containers, got %d", len(expected), len(result))
	}

	for i, name := range expected {
		if result[i] != name {
			t.Errorf("expected result[%d] = %q, got %q", i, name, result[i])
		}
	}
}

func TestReachableContainers_ExcludesSelfFromResults(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"backend": {
			{Name: "api", Networks: []string{"backend"}},
			{Name: "worker", Networks: []string{"backend"}},
		},
	}

	result := ReachableContainers("api", "backend", netMap)

	for _, name := range result {
		if name == "api" {
			t.Errorf("self container 'api' should not be in result")
		}
	}

	if len(result) != 1 || result[0] != "worker" {
		t.Errorf("expected [worker], got %v", result)
	}
}

func TestReachableContainers_ReturnsEmptySliceWhenNoOtherContainers(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"isolated": {
			{Name: "lonely", Networks: []string{"isolated"}},
		},
	}

	result := ReachableContainers("lonely", "isolated", netMap)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestReachableContainers_ReturnsEmptySliceWhenNetworkNotFound(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"existing": {
			{Name: "container1", Networks: []string{"existing"}},
		},
	}

	result := ReachableContainers("container1", "nonexistent", netMap)

	if len(result) != 0 {
		t.Errorf("expected empty slice for nonexistent network, got %v", result)
	}
}

func TestReachableContainers_ResultsAreSortedAlphabetically(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"network": {
			{Name: "self", Networks: []string{"network"}},
			{Name: "zebra", Networks: []string{"network"}},
			{Name: "apple", Networks: []string{"network"}},
			{Name: "mango", Networks: []string{"network"}},
		},
	}

	result := ReachableContainers("self", "network", netMap)

	expected := []string{"apple", "mango", "zebra"}
	if len(result) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(result))
	}

	for i, name := range expected {
		if result[i] != name {
			t.Errorf("expected sorted result[%d] = %q, got %q", i, name, result[i])
		}
	}
}

func TestReachableContainers_EmptyNetMap(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{}

	result := ReachableContainers("any", "any", netMap)

	if len(result) != 0 {
		t.Errorf("expected empty slice for empty netMap, got %v", result)
	}
}

func TestReachableContainers_NilNetMap(t *testing.T) {
	result := ReachableContainers("any", "any", nil)

	if len(result) != 0 {
		t.Errorf("expected empty slice for nil netMap, got %v", result)
	}
}

func TestReachableContainers_MultipleNetworks(t *testing.T) {
	netMap := map[string][]models.ContainerInfo{
		"frontend": {
			{Name: "nginx", Networks: []string{"frontend"}},
			{Name: "api", Networks: []string{"frontend", "backend"}},
		},
		"backend": {
			{Name: "api", Networks: []string{"frontend", "backend"}},
			{Name: "db", Networks: []string{"backend"}},
			{Name: "cache", Networks: []string{"backend"}},
		},
	}

	// Test frontend reachability from api
	frontendResult := ReachableContainers("api", "frontend", netMap)
	if len(frontendResult) != 1 || frontendResult[0] != "nginx" {
		t.Errorf("expected [nginx] for frontend, got %v", frontendResult)
	}

	// Test backend reachability from api
	backendResult := ReachableContainers("api", "backend", netMap)
	expected := []string{"cache", "db"}
	if len(backendResult) != len(expected) {
		t.Errorf("expected %d results for backend, got %d", len(expected), len(backendResult))
	}
	for i, name := range expected {
		if backendResult[i] != name {
			t.Errorf("expected backend result[%d] = %q, got %q", i, name, backendResult[i])
		}
	}
}

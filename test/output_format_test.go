// Package integration provides output format integration tests for docker-network-viz.
package integration

import (
	"bytes"
	"strings"
	"testing"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
	"git.o.ocom.com.au/go/docker-network-viz/internal/output"
)

// TestOutputFormat_NetworkTreeStructure tests that network tree output has correct structure.
func TestOutputFormat_NetworkTreeStructure(t *testing.T) {
	netInfo := models.NewNetworkInfo("test_network", "bridge")

	containers := []models.ContainerInfo{
		{
			Name:     "container_a",
			Aliases:  []string{"alias1", "alias2"},
			Networks: []string{"test_network"},
		},
		{
			Name:     "container_b",
			Aliases:  []string{"alias3"},
			Networks: []string{"test_network"},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")

	// First line should be network header
	if !strings.HasPrefix(lines[0], "Network:") {
		t.Errorf("first line should start with 'Network:', got: %s", lines[0])
	}

	// Should contain network name and driver
	if !strings.Contains(lines[0], "test_network") {
		t.Errorf("first line should contain network name, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "bridge") {
		t.Errorf("first line should contain driver, got: %s", lines[0])
	}

	// Should contain both containers
	foundContainerA := false
	foundContainerB := false
	for _, line := range lines {
		if strings.Contains(line, "container_a") {
			foundContainerA = true
		}
		if strings.Contains(line, "container_b") {
			foundContainerB = true
		}
	}

	if !foundContainerA {
		t.Error("output should contain container_a")
	}
	if !foundContainerB {
		t.Error("output should contain container_b")
	}
}

// TestOutputFormat_ContainerTreeStructure tests that container tree output has correct structure.
func TestOutputFormat_ContainerTreeStructure(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "test_container",
		Aliases:  []string{"alias1"},
		Networks: []string{"network_a", "network_b"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"network_a": {
			{Name: "test_container", Networks: []string{"network_a", "network_b"}},
			{Name: "peer_a", Networks: []string{"network_a"}},
		},
		"network_b": {
			{Name: "test_container", Networks: []string{"network_a", "network_b"}},
			{Name: "peer_b", Networks: []string{"network_b"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")

	// First line should be container header
	if !strings.HasPrefix(lines[0], "Container:") {
		t.Errorf("first line should start with 'Container:', got: %s", lines[0])
	}

	// Should contain container name
	if !strings.Contains(lines[0], "test_container") {
		t.Errorf("first line should contain container name, got: %s", lines[0])
	}

	// Should contain both networks
	foundNetworkA := false
	foundNetworkB := false
	for _, line := range lines {
		if strings.Contains(line, "network_a") {
			foundNetworkA = true
		}
		if strings.Contains(line, "network_b") {
			foundNetworkB = true
		}
	}

	if !foundNetworkA {
		t.Error("output should contain network_a")
	}
	if !foundNetworkB {
		t.Error("output should contain network_b")
	}
}

// TestOutputFormat_TreeSymbols tests that correct tree symbols are used.
func TestOutputFormat_TreeSymbols(t *testing.T) {
	netInfo := models.NewNetworkInfo("test", "bridge")

	containers := []models.ContainerInfo{
		{Name: "first", Aliases: []string{"a"}, Networks: []string{"test"}},
		{Name: "middle", Aliases: []string{"b"}, Networks: []string{"test"}},
		{Name: "last", Aliases: []string{"c"}, Networks: []string{"test"}},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Should have branch symbols for non-last items
	if !strings.Contains(result, output.TreeBranch) {
		t.Error("output should contain tree branch symbol")
	}

	// Should have end symbol for last item
	if !strings.Contains(result, output.TreeEnd) {
		t.Error("output should contain tree end symbol")
	}
}

// TestOutputFormat_AliasDisplay tests that aliases are displayed correctly.
func TestOutputFormat_AliasDisplay(t *testing.T) {
	netInfo := models.NewNetworkInfo("test", "bridge")

	containers := []models.ContainerInfo{
		{
			Name:     "container",
			Aliases:  []string{"alias_one", "alias_two", "alias_three"},
			Networks: []string{"test"},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Should display all aliases
	for _, alias := range []string{"alias_one", "alias_two", "alias_three"} {
		if !strings.Contains(result, alias) {
			t.Errorf("output should contain alias %q", alias)
		}
	}

	// Should have "alias:" labels
	aliasCount := strings.Count(result, "alias:")
	if aliasCount != 3 {
		t.Errorf("expected 3 alias labels, got %d", aliasCount)
	}
}

// TestOutputFormat_ConnectsToDisplay tests that "connects to:" section is formatted correctly.
func TestOutputFormat_ConnectsToDisplay(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "main",
		Aliases:  []string{},
		Networks: []string{"shared_net"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"shared_net": {
			{Name: "main", Networks: []string{"shared_net"}},
			{Name: "peer1", Networks: []string{"shared_net"}},
			{Name: "peer2", Networks: []string{"shared_net"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	// Should contain "connects to:" label
	if !strings.Contains(result, "connects to:") {
		t.Error("output should contain 'connects to:' label")
	}

	// Should contain peer containers
	if !strings.Contains(result, "peer1") {
		t.Error("output should contain peer1")
	}
	if !strings.Contains(result, "peer2") {
		t.Error("output should contain peer2")
	}

	// Should NOT contain self
	// Count occurrences of "main" - should appear only once (in header)
	mainCount := strings.Count(result, "main")
	if mainCount != 1 {
		t.Errorf("expected container name 'main' to appear once (in header), got %d occurrences", mainCount)
	}
}

// TestOutputFormat_NoContainersMessage tests the message for empty networks.
func TestOutputFormat_NoContainersMessage(t *testing.T) {
	netInfo := models.NewNetworkInfo("empty_net", "bridge")
	containers := []models.ContainerInfo{}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	if !strings.Contains(result, "no containers") {
		t.Errorf("expected 'no containers' message, got:\n%s", result)
	}
}

// TestOutputFormat_NoReachableContainersMessage tests the message for isolated containers.
func TestOutputFormat_NoReachableContainersMessage(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "lonely",
		Aliases:  []string{},
		Networks: []string{"isolated_net"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"isolated_net": {
			{Name: "lonely", Networks: []string{"isolated_net"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	if !strings.Contains(result, "(none)") {
		t.Errorf("expected '(none)' message for isolated container, got:\n%s", result)
	}
}

// TestOutputFormat_MultipleNetworksPerContainer tests display of multi-homed containers.
func TestOutputFormat_MultipleNetworksPerContainer(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "multihomed",
		Aliases:  []string{},
		Networks: []string{"frontend", "backend", "management"},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"frontend": {
			{Name: "multihomed", Networks: []string{"frontend", "backend", "management"}},
			{Name: "web", Networks: []string{"frontend"}},
		},
		"backend": {
			{Name: "multihomed", Networks: []string{"frontend", "backend", "management"}},
			{Name: "db", Networks: []string{"backend"}},
		},
		"management": {
			{Name: "multihomed", Networks: []string{"frontend", "backend", "management"}},
			{Name: "monitor", Networks: []string{"management"}},
		},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	// Should contain all three networks
	for _, net := range []string{"frontend", "backend", "management"} {
		if !strings.Contains(result, net) {
			t.Errorf("output should contain network %q", net)
		}
	}

	// Should contain reachable containers from each network
	for _, peer := range []string{"web", "db", "monitor"} {
		if !strings.Contains(result, peer) {
			t.Errorf("output should contain reachable container %q", peer)
		}
	}
}

// TestOutputFormat_SortedAliases tests that aliases are sorted alphabetically.
func TestOutputFormat_SortedAliases(t *testing.T) {
	netInfo := models.NewNetworkInfo("test", "bridge")

	containers := []models.ContainerInfo{
		{
			Name:     "container",
			Aliases:  []string{"zebra", "apple", "mango"}, // Unsorted
			Networks: []string{"test"},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Find positions of aliases
	applePos := strings.Index(result, "apple")
	mangoPos := strings.Index(result, "mango")
	zebraPos := strings.Index(result, "zebra")

	if applePos == -1 || mangoPos == -1 || zebraPos == -1 {
		t.Fatalf("expected all aliases in output, got:\n%s", result)
	}

	// Verify alphabetical order
	if !(applePos < mangoPos && mangoPos < zebraPos) {
		t.Errorf("expected aliases in alphabetical order, got:\n%s", result)
	}
}

// TestOutputFormat_SortedNetworksInContainerTree tests that networks are sorted.
func TestOutputFormat_SortedNetworksInContainerTree(t *testing.T) {
	containerInfo := &models.ContainerInfo{
		Name:     "container",
		Aliases:  []string{},
		Networks: []string{"zebra_net", "alpha_net", "middle_net"}, // Unsorted
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"zebra_net":  {{Name: "container", Networks: containerInfo.Networks}},
		"alpha_net":  {{Name: "container", Networks: containerInfo.Networks}},
		"middle_net": {{Name: "container", Networks: containerInfo.Networks}},
	}

	var buf bytes.Buffer
	output.PrintContainerTree(&buf, containerInfo, networkToContainers)

	result := buf.String()

	// Find positions of networks
	alphaPos := strings.Index(result, "alpha_net")
	middlePos := strings.Index(result, "middle_net")
	zebraPos := strings.Index(result, "zebra_net")

	if alphaPos == -1 || middlePos == -1 || zebraPos == -1 {
		t.Fatalf("expected all networks in output, got:\n%s", result)
	}

	// Verify alphabetical order
	if !(alphaPos < middlePos && middlePos < zebraPos) {
		t.Errorf("expected networks in alphabetical order, got:\n%s", result)
	}
}

// TestOutputFormat_LongNames tests handling of long container/network names.
func TestOutputFormat_LongNames(t *testing.T) {
	longName := "this_is_a_very_long_container_name_that_might_cause_formatting_issues"
	longNetworkName := "extremely_long_network_name_for_testing_purposes_and_formatting"

	netInfo := models.NewNetworkInfo(longNetworkName, "bridge")
	containers := []models.ContainerInfo{
		{
			Name:     longName,
			Aliases:  []string{"short"},
			Networks: []string{longNetworkName},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Should contain the full long names
	if !strings.Contains(result, longName) {
		t.Errorf("output should contain long container name")
	}
	if !strings.Contains(result, longNetworkName) {
		t.Errorf("output should contain long network name")
	}
}

// TestOutputFormat_SpecialCharactersInNames tests handling of special characters.
func TestOutputFormat_SpecialCharactersInNames(t *testing.T) {
	// Docker allows underscores, hyphens, and periods in names
	specialName := "container-with_special.name"

	netInfo := models.NewNetworkInfo("test-network_name.v2", "bridge")
	containers := []models.ContainerInfo{
		{
			Name:     specialName,
			Aliases:  []string{"alias-with_periods.v1"},
			Networks: []string{"test-network_name.v2"},
		},
	}

	var buf bytes.Buffer
	output.PrintNetworkTree(&buf, *netInfo, containers)

	result := buf.String()

	// Should handle special characters correctly
	if !strings.Contains(result, specialName) {
		t.Errorf("output should contain container name with special characters")
	}
	if !strings.Contains(result, "alias-with_periods.v1") {
		t.Errorf("output should contain alias with special characters")
	}
}

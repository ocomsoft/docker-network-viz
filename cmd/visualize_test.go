package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/network"
	"github.com/spf13/viper"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// TestVisualizeCommandExists verifies that the visualize command is properly defined.
func TestVisualizeCommandExists(t *testing.T) {
	if visualizeCmd == nil {
		t.Fatal("visualize command should not be nil")
	}

	if visualizeCmd.Use != "visualize" {
		t.Errorf("visualize command Use should be 'visualize', got %q", visualizeCmd.Use)
	}

	if visualizeCmd.Short == "" {
		t.Error("visualize command Short should not be empty")
	}
}

// TestVisualizeCommandHasFlags verifies that the visualize command has required flags.
func TestVisualizeCommandHasFlags(t *testing.T) {
	// Check for only-network flag
	onlyNetworkFlag := visualizeCmd.Flags().Lookup("only-network")
	if onlyNetworkFlag == nil {
		t.Error("visualize command should have an only-network flag")
	}

	// Check for container flag
	containerFlagObj := visualizeCmd.Flags().Lookup("container")
	if containerFlagObj == nil {
		t.Error("visualize command should have a container flag")
	}

	// Check for no-aliases flag
	noAliasesFlag := visualizeCmd.Flags().Lookup("no-aliases")
	if noAliasesFlag == nil {
		t.Error("visualize command should have a no-aliases flag")
	}
}

// TestPrintVisualizationNetworkTree verifies network tree output.
func TestPrintVisualizationNetworkTree(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	networks := []network.Summary{
		{Name: "bridge", Driver: "bridge"},
		{Name: "test_net", Driver: "bridge"},
	}

	containerMap := map[string]*models.ContainerInfo{
		"web": {
			Name:     "web",
			Aliases:  []string{"www"},
			Networks: []string{"bridge"},
		},
		"db": {
			Name:     "db",
			Aliases:  []string{"database"},
			Networks: []string{"bridge"},
		},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"bridge": {
			{Name: "web", Aliases: []string{"www"}, Networks: []string{"bridge"}},
			{Name: "db", Aliases: []string{"database"}, Networks: []string{"bridge"}},
		},
		"test_net": {},
	}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error: %v", err)
	}

	output := buf.String()

	// Verify network section header
	if !strings.Contains(output, "=== Networks ===") {
		t.Error("output should contain '=== Networks ===' header")
	}

	// Verify network names are present
	if !strings.Contains(output, "Network: bridge (bridge)") {
		t.Error("output should contain 'Network: bridge (bridge)'")
	}

	if !strings.Contains(output, "Network: test_net (bridge)") {
		t.Error("output should contain 'Network: test_net (bridge)'")
	}

	// Verify container reachability section header
	if !strings.Contains(output, "=== Containers (Reachability) ===") {
		t.Error("output should contain '=== Containers (Reachability) ===' header")
	}

	// Verify containers are listed
	if !strings.Contains(output, "Container: web") {
		t.Error("output should contain 'Container: web'")
	}

	if !strings.Contains(output, "Container: db") {
		t.Error("output should contain 'Container: db'")
	}
}

// TestPrintVisualizationWithOnlyNetworkFilter verifies network filtering.
func TestPrintVisualizationWithOnlyNetworkFilter(t *testing.T) {
	// Reset viper for this test
	viper.Reset()
	viper.Set("only-network", "bridge")

	networks := []network.Summary{
		{Name: "bridge", Driver: "bridge"},
		{Name: "other_net", Driver: "bridge"},
	}

	containerMap := map[string]*models.ContainerInfo{
		"web": {
			Name:     "web",
			Aliases:  []string{},
			Networks: []string{"bridge"},
		},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"bridge":    {{Name: "web", Aliases: []string{}, Networks: []string{"bridge"}}},
		"other_net": {},
	}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error: %v", err)
	}

	output := buf.String()

	// Should contain bridge network
	if !strings.Contains(output, "Network: bridge (bridge)") {
		t.Error("output should contain 'Network: bridge (bridge)'")
	}

	// Should NOT contain other_net
	if strings.Contains(output, "Network: other_net") {
		t.Error("output should not contain 'Network: other_net' when filter is set")
	}
}

// TestPrintVisualizationWithContainerFilter verifies container filtering.
func TestPrintVisualizationWithContainerFilter(t *testing.T) {
	// Reset viper for this test
	viper.Reset()
	viper.Set("container", "web")

	networks := []network.Summary{
		{Name: "bridge", Driver: "bridge"},
	}

	containerMap := map[string]*models.ContainerInfo{
		"web": {
			Name:     "web",
			Aliases:  []string{},
			Networks: []string{"bridge"},
		},
		"db": {
			Name:     "db",
			Aliases:  []string{},
			Networks: []string{"bridge"},
		},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"bridge": {
			{Name: "web", Aliases: []string{}, Networks: []string{"bridge"}},
			{Name: "db", Aliases: []string{}, Networks: []string{"bridge"}},
		},
	}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error: %v", err)
	}

	output := buf.String()

	// Should contain web container
	if !strings.Contains(output, "Container: web") {
		t.Error("output should contain 'Container: web'")
	}

	// Should NOT contain db container in the Container section
	// (it may appear in the network tree, but not as a top-level container)
	lines := strings.Split(output, "\n")
	dbContainerFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Container: db") {
			dbContainerFound = true
			break
		}
	}

	if dbContainerFound {
		t.Error("output should not contain 'Container: db' as a header when filter is set")
	}
}

// TestPrintVisualizationWithNoAliases verifies alias hiding.
func TestPrintVisualizationWithNoAliases(t *testing.T) {
	// Reset viper for this test
	viper.Reset()
	viper.Set("no-aliases", true)

	networks := []network.Summary{
		{Name: "bridge", Driver: "bridge"},
	}

	containerMap := map[string]*models.ContainerInfo{
		"web": {
			Name:     "web",
			Aliases:  []string{"www", "webapp"},
			Networks: []string{"bridge"},
		},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"bridge": {
			{Name: "web", Aliases: []string{"www", "webapp"}, Networks: []string{"bridge"}},
		},
	}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error: %v", err)
	}

	output := buf.String()

	// Aliases should not appear
	if strings.Contains(output, "alias: www") {
		t.Error("output should not contain aliases when no-aliases flag is set")
	}

	if strings.Contains(output, "alias: webapp") {
		t.Error("output should not contain aliases when no-aliases flag is set")
	}
}

// TestRemoveAliasesFromContainers verifies the alias removal function.
func TestRemoveAliasesFromContainers(t *testing.T) {
	containers := []models.ContainerInfo{
		{Name: "web", Aliases: []string{"www", "webapp"}, Networks: []string{"bridge"}},
		{Name: "db", Aliases: []string{"database"}, Networks: []string{"bridge"}},
	}

	result := removeAliasesFromContainers(containers)

	if len(result) != 2 {
		t.Errorf("expected 2 containers, got %d", len(result))
	}

	// Check that aliases are removed
	for _, c := range result {
		if len(c.Aliases) != 0 {
			t.Errorf("container %s should have no aliases, got %v", c.Name, c.Aliases)
		}
	}

	// Check that names and networks are preserved
	if result[0].Name != "web" {
		t.Errorf("expected first container name to be 'web', got %q", result[0].Name)
	}

	if len(result[0].Networks) != 1 || result[0].Networks[0] != "bridge" {
		t.Errorf("expected first container networks to be ['bridge'], got %v", result[0].Networks)
	}

	// Check that original containers are not modified
	if len(containers[0].Aliases) != 2 {
		t.Error("original container aliases should not be modified")
	}
}

// TestPrintVisualizationEmptyNetworks verifies behavior with no networks.
func TestPrintVisualizationEmptyNetworks(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	networks := []network.Summary{}
	containerMap := map[string]*models.ContainerInfo{}
	networkToContainers := map[string][]models.ContainerInfo{}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error with empty data: %v", err)
	}

	output := buf.String()

	// Should still have section headers
	if !strings.Contains(output, "=== Networks ===") {
		t.Error("output should contain '=== Networks ===' header even with no networks")
	}

	if !strings.Contains(output, "=== Containers (Reachability) ===") {
		t.Error("output should contain '=== Containers (Reachability) ===' header even with no containers")
	}
}

// TestVisualizeCommandHelp verifies that help works for visualize command.
func TestVisualizeCommandHelp(t *testing.T) {
	// Get help output through UsageString which is more reliable in tests
	output := visualizeCmd.UsageString()

	if !strings.Contains(output, "visualize") {
		t.Error("help output should contain 'visualize'")
	}

	if !strings.Contains(output, "--only-network") {
		t.Error("help output should contain '--only-network'")
	}

	if !strings.Contains(output, "--container") {
		t.Error("help output should contain '--container'")
	}

	if !strings.Contains(output, "--no-aliases") {
		t.Error("help output should contain '--no-aliases'")
	}
}

// TestPrintVisualizationMultipleNetworksPerContainer tests containers on multiple networks.
func TestPrintVisualizationMultipleNetworksPerContainer(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	networks := []network.Summary{
		{Name: "backend", Driver: "bridge"},
		{Name: "frontend", Driver: "bridge"},
	}

	containerMap := map[string]*models.ContainerInfo{
		"api": {
			Name:     "api",
			Aliases:  []string{},
			Networks: []string{"backend", "frontend"},
		},
		"db": {
			Name:     "db",
			Aliases:  []string{},
			Networks: []string{"backend"},
		},
		"web": {
			Name:     "web",
			Aliases:  []string{},
			Networks: []string{"frontend"},
		},
	}

	networkToContainers := map[string][]models.ContainerInfo{
		"backend": {
			{Name: "api", Aliases: []string{}, Networks: []string{"backend", "frontend"}},
			{Name: "db", Aliases: []string{}, Networks: []string{"backend"}},
		},
		"frontend": {
			{Name: "api", Aliases: []string{}, Networks: []string{"backend", "frontend"}},
			{Name: "web", Aliases: []string{}, Networks: []string{"frontend"}},
		},
	}

	buf := new(bytes.Buffer)
	err := printVisualization(buf, networks, containerMap, networkToContainers)

	if err != nil {
		t.Errorf("printVisualization should not return error: %v", err)
	}

	output := buf.String()

	// API container should show connectivity to both networks
	if !strings.Contains(output, "Container: api") {
		t.Error("output should contain 'Container: api'")
	}

	// API should show it can reach db through backend
	// API should show it can reach web through frontend
	if !strings.Contains(output, "Network: backend") {
		t.Error("output should show api connected to backend network")
	}

	if !strings.Contains(output, "Network: frontend") {
		t.Error("output should show api connected to frontend network")
	}
}

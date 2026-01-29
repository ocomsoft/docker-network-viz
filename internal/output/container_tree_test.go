package output

import (
	"bytes"
	"strings"
	"testing"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

func TestPrintContainerTree_SingleNetworkNoReachability(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "isolated",
		Aliases:  []string{},
		Networks: []string{"solo_network"},
	}
	netMap := map[string][]models.ContainerInfo{
		"solo_network": {
			{Name: "isolated", Networks: []string{"solo_network"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Check header
	if !strings.Contains(output, "Container: isolated") {
		t.Errorf("missing container header:\n%s", output)
	}

	// Check network
	if !strings.Contains(output, "Network: solo_network") {
		t.Errorf("missing network name:\n%s", output)
	}

	// Check connects to with none
	if !strings.Contains(output, "connects to:") {
		t.Errorf("missing 'connects to:' label:\n%s", output)
	}

	if !strings.Contains(output, "(none)") {
		t.Errorf("expected '(none)' for isolated container:\n%s", output)
	}
}

func TestPrintContainerTree_SingleNetworkWithReachability(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "api",
		Aliases:  []string{},
		Networks: []string{"backend"},
	}
	netMap := map[string][]models.ContainerInfo{
		"backend": {
			{Name: "api", Networks: []string{"backend"}},
			{Name: "db", Networks: []string{"backend"}},
			{Name: "cache", Networks: []string{"backend"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Check reachable containers (alphabetically sorted)
	if !strings.Contains(output, "cache") {
		t.Errorf("missing reachable container 'cache':\n%s", output)
	}

	if !strings.Contains(output, "db") {
		t.Errorf("missing reachable container 'db':\n%s", output)
	}

	// Self should not appear
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "\u2514\u2500\u2500 api") || strings.Contains(line, "\u251c\u2500\u2500 api") {
			// Skip the Container: header line
			if !strings.HasPrefix(line, "Container:") {
				t.Errorf("self container 'api' should not appear in reachability list:\n%s", output)
			}
		}
	}
}

func TestPrintContainerTree_MultipleNetworks(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "api",
		Aliases:  []string{},
		Networks: []string{"frontend", "backend"},
	}
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

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Check both networks are present (sorted alphabetically)
	if !strings.Contains(output, "Network: backend") {
		t.Errorf("missing backend network:\n%s", output)
	}

	if !strings.Contains(output, "Network: frontend") {
		t.Errorf("missing frontend network:\n%s", output)
	}

	// Check reachable containers per network
	if !strings.Contains(output, "nginx") {
		t.Errorf("missing reachable container 'nginx' from frontend:\n%s", output)
	}

	if !strings.Contains(output, "db") {
		t.Errorf("missing reachable container 'db' from backend:\n%s", output)
	}

	if !strings.Contains(output, "cache") {
		t.Errorf("missing reachable container 'cache' from backend:\n%s", output)
	}
}

func TestPrintContainerTree_NetworksSortedAlphabetically(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "service",
		Aliases:  []string{},
		Networks: []string{"zebra_net", "alpha_net", "beta_net"},
	}
	netMap := map[string][]models.ContainerInfo{
		"zebra_net": {{Name: "service", Networks: []string{"zebra_net", "alpha_net", "beta_net"}}},
		"alpha_net": {{Name: "service", Networks: []string{"zebra_net", "alpha_net", "beta_net"}}},
		"beta_net":  {{Name: "service", Networks: []string{"zebra_net", "alpha_net", "beta_net"}}},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Find network lines
	alphaIdx := strings.Index(output, "Network: alpha_net")
	betaIdx := strings.Index(output, "Network: beta_net")
	zebraIdx := strings.Index(output, "Network: zebra_net")

	if alphaIdx == -1 || betaIdx == -1 || zebraIdx == -1 {
		t.Fatalf("missing networks in output:\n%s", output)
	}

	if !(alphaIdx < betaIdx && betaIdx < zebraIdx) {
		t.Errorf("networks not in alphabetical order: alpha=%d, beta=%d, zebra=%d\n%s",
			alphaIdx, betaIdx, zebraIdx, output)
	}
}

func TestPrintContainerTree_ReachableContainersSortedAlphabetically(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "api",
		Aliases:  []string{},
		Networks: []string{"network"},
	}
	netMap := map[string][]models.ContainerInfo{
		"network": {
			{Name: "api", Networks: []string{"network"}},
			{Name: "zebra", Networks: []string{"network"}},
			{Name: "alpha", Networks: []string{"network"}},
			{Name: "beta", Networks: []string{"network"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	alphaIdx := strings.Index(output, "alpha")
	betaIdx := strings.Index(output, "beta")
	zebraIdx := strings.Index(output, "zebra")

	if alphaIdx == -1 || betaIdx == -1 || zebraIdx == -1 {
		t.Fatalf("missing reachable containers:\n%s", output)
	}

	if !(alphaIdx < betaIdx && betaIdx < zebraIdx) {
		t.Errorf("reachable containers not in alphabetical order:\n%s", output)
	}
}

func TestPrintContainerTree_TreePrefixesCorrect(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "service",
		Aliases:  []string{},
		Networks: []string{"net1", "net2"},
	}
	netMap := map[string][]models.ContainerInfo{
		"net1": {
			{Name: "service", Networks: []string{"net1", "net2"}},
			{Name: "other1", Networks: []string{"net1"}},
		},
		"net2": {
			{Name: "service", Networks: []string{"net1", "net2"}},
			{Name: "other2", Networks: []string{"net2"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// First network should have branch prefix
	foundFirstNet := false
	foundLastNet := false
	for _, line := range lines {
		if strings.Contains(line, "Network: net1") {
			if !strings.Contains(line, "\u251c\u2500\u2500") {
				t.Errorf("first network should have branch prefix:\n%s", line)
			}
			foundFirstNet = true
		}
		if strings.Contains(line, "Network: net2") {
			if !strings.Contains(line, "\u2514\u2500\u2500") {
				t.Errorf("last network should have end prefix:\n%s", line)
			}
			foundLastNet = true
		}
	}

	if !foundFirstNet || !foundLastNet {
		t.Errorf("did not find expected network lines:\n%s", output)
	}
}

func TestPrintContainerTree_EmptyNetworks(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "orphan",
		Aliases:  []string{},
		Networks: []string{},
	}
	netMap := map[string][]models.ContainerInfo{}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Should just have the container header
	expected := "Container: orphan\n"
	if output != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestPrintContainerTree_DoesNotModifyOriginalNetworks(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "service",
		Aliases:  []string{},
		Networks: []string{"zebra", "alpha"},
	}
	netMap := map[string][]models.ContainerInfo{
		"zebra": {{Name: "service", Networks: []string{"zebra", "alpha"}}},
		"alpha": {{Name: "service", Networks: []string{"zebra", "alpha"}}},
	}

	// Keep original order
	originalOrder := []string{c.Networks[0], c.Networks[1]}

	PrintContainerTree(&buf, c, netMap)

	// Verify original slice is not modified
	if c.Networks[0] != originalOrder[0] || c.Networks[1] != originalOrder[1] {
		t.Errorf("original networks slice was modified: expected %v, got %v",
			originalOrder, c.Networks)
	}
}

func TestPrintContainerTree_SingleReachableContainer(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "web",
		Aliases:  []string{},
		Networks: []string{"frontend"},
	}
	netMap := map[string][]models.ContainerInfo{
		"frontend": {
			{Name: "web", Networks: []string{"frontend"}},
			{Name: "nginx", Networks: []string{"frontend"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()

	// Single reachable container should have end prefix
	if !strings.Contains(output, "\u2514\u2500\u2500 nginx") {
		t.Errorf("single reachable container should have end prefix:\n%s", output)
	}
}

func TestPrintContainerTree_MultipleReachableContainersPrefixes(t *testing.T) {
	var buf bytes.Buffer
	c := &models.ContainerInfo{
		Name:     "api",
		Aliases:  []string{},
		Networks: []string{"backend"},
	}
	netMap := map[string][]models.ContainerInfo{
		"backend": {
			{Name: "api", Networks: []string{"backend"}},
			{Name: "cache", Networks: []string{"backend"}},
			{Name: "db", Networks: []string{"backend"}},
			{Name: "worker", Networks: []string{"backend"}},
		},
	}

	PrintContainerTree(&buf, c, netMap)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Find lines with reachable containers (cache, db, worker - sorted)
	var reachableLines []string
	for _, line := range lines {
		if strings.Contains(line, "cache") || strings.Contains(line, "db") || strings.Contains(line, "worker") {
			reachableLines = append(reachableLines, line)
		}
	}

	if len(reachableLines) != 3 {
		t.Fatalf("expected 3 reachable container lines, got %d:\n%s", len(reachableLines), output)
	}

	// First two should have branch prefix, last should have end prefix
	if !strings.Contains(reachableLines[0], "\u251c\u2500\u2500") {
		t.Errorf("first reachable should have branch prefix:\n%s", reachableLines[0])
	}

	if !strings.Contains(reachableLines[1], "\u251c\u2500\u2500") {
		t.Errorf("second reachable should have branch prefix:\n%s", reachableLines[1])
	}

	if !strings.Contains(reachableLines[2], "\u2514\u2500\u2500") {
		t.Errorf("last reachable should have end prefix:\n%s", reachableLines[2])
	}
}

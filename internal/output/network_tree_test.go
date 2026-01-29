package output

import (
	"bytes"
	"strings"
	"testing"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

func TestPrintNetworkTree_EmptyContainers(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "bridge", Driver: "bridge"}

	PrintNetworkTree(&buf, net, []models.ContainerInfo{})

	output := buf.String()
	expected := "Network: bridge (bridge)\n\u2514\u2500\u2500 (no containers)\n"

	if output != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestPrintNetworkTree_SingleContainer(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "backend", Driver: "bridge"}
	containers := []models.ContainerInfo{
		{Name: "api", Aliases: []string{}, Networks: []string{"backend"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()

	// Check header
	if !strings.Contains(output, "Network: backend (bridge)") {
		t.Errorf("missing network header in output:\n%s", output)
	}

	// Check container with last-item prefix
	if !strings.Contains(output, "\u2514\u2500\u2500 api") {
		t.Errorf("expected container with last-item prefix, got:\n%s", output)
	}
}

func TestPrintNetworkTree_MultipleContainers(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "frontend", Driver: "overlay"}
	containers := []models.ContainerInfo{
		{Name: "nginx", Aliases: []string{}, Networks: []string{"frontend"}},
		{Name: "web", Aliases: []string{}, Networks: []string{"frontend"}},
		{Name: "api", Aliases: []string{}, Networks: []string{"frontend"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// Check header
	if lines[0] != "Network: frontend (overlay)" {
		t.Errorf("expected header 'Network: frontend (overlay)', got %q", lines[0])
	}

	// Containers should be sorted: api, nginx, web
	// First two should have branch prefix, last should have end prefix
	if !strings.Contains(lines[1], "\u251c\u2500\u2500 api") {
		t.Errorf("expected api with branch prefix on line 1, got %q", lines[1])
	}

	if !strings.Contains(lines[2], "\u251c\u2500\u2500 nginx") {
		t.Errorf("expected nginx with branch prefix on line 2, got %q", lines[2])
	}

	if !strings.Contains(lines[3], "\u2514\u2500\u2500 web") {
		t.Errorf("expected web with end prefix on line 3, got %q", lines[3])
	}
}

func TestPrintNetworkTree_ContainerWithAliases(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "bridge", Driver: "bridge"}
	containers := []models.ContainerInfo{
		{Name: "web_app", Aliases: []string{"web", "web.local"}, Networks: []string{"bridge"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()

	// Check container name
	if !strings.Contains(output, "web_app") {
		t.Errorf("missing container name in output:\n%s", output)
	}

	// Check aliases are present (sorted)
	if !strings.Contains(output, "alias: web") {
		t.Errorf("missing alias 'web' in output:\n%s", output)
	}

	if !strings.Contains(output, "alias: web.local") {
		t.Errorf("missing alias 'web.local' in output:\n%s", output)
	}
}

func TestPrintNetworkTree_MultipleContainersWithAliases(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "services", Driver: "bridge"}
	containers := []models.ContainerInfo{
		{Name: "redis", Aliases: []string{"cache", "redis-server"}, Networks: []string{"services"}},
		{Name: "postgres", Aliases: []string{"db"}, Networks: []string{"services"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()

	// Postgres comes before redis alphabetically
	// Check postgres appears with branch prefix (not last)
	if !strings.Contains(output, "\u251c\u2500\u2500 postgres") {
		t.Errorf("expected postgres with branch prefix:\n%s", output)
	}

	// Check redis appears with end prefix (last)
	if !strings.Contains(output, "\u2514\u2500\u2500 redis") {
		t.Errorf("expected redis with end prefix:\n%s", output)
	}

	// Check aliases are present
	if !strings.Contains(output, "alias: db") {
		t.Errorf("missing postgres alias 'db':\n%s", output)
	}

	if !strings.Contains(output, "alias: cache") {
		t.Errorf("missing redis alias 'cache':\n%s", output)
	}
}

func TestPrintNetworkTree_SortsContainersByName(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "test", Driver: "bridge"}
	// Input in unsorted order
	containers := []models.ContainerInfo{
		{Name: "zebra", Aliases: []string{}, Networks: []string{"test"}},
		{Name: "apple", Aliases: []string{}, Networks: []string{"test"}},
		{Name: "mango", Aliases: []string{}, Networks: []string{"test"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// Skip header line, check container order
	if !strings.Contains(lines[1], "apple") {
		t.Errorf("expected 'apple' first, got %q", lines[1])
	}

	if !strings.Contains(lines[2], "mango") {
		t.Errorf("expected 'mango' second, got %q", lines[2])
	}

	if !strings.Contains(lines[3], "zebra") {
		t.Errorf("expected 'zebra' third, got %q", lines[3])
	}
}

func TestPrintNetworkTree_SortsAliasesByName(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "test", Driver: "bridge"}
	containers := []models.ContainerInfo{
		{Name: "service", Aliases: []string{"zulu", "alpha", "bravo"}, Networks: []string{"test"}},
	}

	PrintNetworkTree(&buf, net, containers)

	output := buf.String()

	// Find alias lines
	lines := strings.Split(output, "\n")
	var aliasLines []string
	for _, line := range lines {
		if strings.Contains(line, "alias:") {
			aliasLines = append(aliasLines, line)
		}
	}

	if len(aliasLines) != 3 {
		t.Errorf("expected 3 alias lines, got %d", len(aliasLines))
	}

	// Check sorted order
	if !strings.Contains(aliasLines[0], "alpha") {
		t.Errorf("expected alpha first, got %q", aliasLines[0])
	}

	if !strings.Contains(aliasLines[1], "bravo") {
		t.Errorf("expected bravo second, got %q", aliasLines[1])
	}

	if !strings.Contains(aliasLines[2], "zulu") {
		t.Errorf("expected zulu third, got %q", aliasLines[2])
	}
}

func TestPrintNetworkTree_DoesNotModifyOriginalContainers(t *testing.T) {
	var buf bytes.Buffer
	net := models.NetworkInfo{Name: "test", Driver: "bridge"}
	original := []models.ContainerInfo{
		{Name: "zebra", Aliases: []string{}, Networks: []string{"test"}},
		{Name: "apple", Aliases: []string{}, Networks: []string{"test"}},
	}

	// Keep a copy of original order
	originalOrder := []string{original[0].Name, original[1].Name}

	PrintNetworkTree(&buf, net, original)

	// Verify original slice is not modified
	if original[0].Name != originalOrder[0] || original[1].Name != originalOrder[1] {
		t.Errorf("original slice was modified: expected %v, got [%s, %s]",
			originalOrder, original[0].Name, original[1].Name)
	}
}

func TestPrintNetworkTree_DifferentDriverTypes(t *testing.T) {
	testCases := []struct {
		driver   string
		expected string
	}{
		{"bridge", "Network: test (bridge)"},
		{"overlay", "Network: test (overlay)"},
		{"host", "Network: test (host)"},
		{"macvlan", "Network: test (macvlan)"},
		{"none", "Network: test (none)"},
	}

	for _, tc := range testCases {
		t.Run(tc.driver, func(t *testing.T) {
			var buf bytes.Buffer
			net := models.NetworkInfo{Name: "test", Driver: tc.driver}

			PrintNetworkTree(&buf, net, []models.ContainerInfo{})

			if !strings.HasPrefix(buf.String(), tc.expected) {
				t.Errorf("expected output to start with %q, got %q", tc.expected, buf.String())
			}
		})
	}
}

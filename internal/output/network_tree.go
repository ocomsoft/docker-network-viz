// Package output provides tree-style formatters for Docker network topology visualization.
package output

import (
	"fmt"
	"io"
	"sort"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// PrintNetworkTree prints a tree-style representation of a Docker network
// and its connected containers to the provided writer.
//
// The output format shows the network name and driver, followed by a tree
// of containers connected to that network. Each container's aliases are
// shown as nested items beneath the container name.
//
// Example output:
//
//	Network: bridge (bridge)
//	├── web_app
//	│   ├── alias: web
//	│   └── alias: web.local
//	├── redis
//	│   └── alias: redis
//	└── postgres
//	    └── alias: db
//
// Parameters:
//   - w: The io.Writer to write the output to
//   - net: The NetworkInfo containing the network name and driver
//   - containers: Slice of ContainerInfo for containers connected to this network
func PrintNetworkTree(w io.Writer, net models.NetworkInfo, containers []models.ContainerInfo) {
	cw := NewColorWriter(w)

	fmt.Fprintf(w, "%s %s (%s)\n",
		cw.Label("Network:"),
		cw.Network(net.Name),
		net.Driver)

	if len(containers) == 0 {
		fmt.Fprintf(w, "%s (no containers)\n", cw.Tree(TreeEnd))
		return
	}

	// Sort containers by name for consistent output
	sortedContainers := make([]models.ContainerInfo, len(containers))
	copy(sortedContainers, containers)
	sort.Slice(sortedContainers, func(i, j int) bool {
		return sortedContainers[i].Name < sortedContainers[j].Name
	})

	for i, c := range sortedContainers {
		prefix := TreeBranch
		indent := TreeVertical
		if i == len(sortedContainers)-1 {
			prefix = TreeEnd
			indent = TreeSpace
		}

		fmt.Fprintf(w, "%s %s\n", cw.Tree(prefix), cw.Container(c.Name))

		// Sort aliases for consistent output
		sortedAliases := c.SortedAliases()
		for j, a := range sortedAliases {
			aliasPrefix := TreeBranch
			if j == len(sortedAliases)-1 {
				aliasPrefix = TreeEnd
			}
			fmt.Fprintf(w, "%s%s %s %s\n",
				cw.Tree(indent),
				cw.Tree(aliasPrefix),
				cw.Label("alias:"),
				cw.Alias(a))
		}
	}
}

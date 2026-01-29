// Package output provides tree-style formatters for Docker network topology visualization.
package output

import (
	"fmt"
	"io"
	"sort"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// PrintContainerTree prints a tree-style representation of a container's
// network connectivity and reachability to other containers.
//
// The output shows the container name, followed by each network it belongs to,
// and under each network, the list of other containers that can be reached
// through that network.
//
// Example output:
//
//	Container: api
//	├── Network: frontend_net
//	│   └── connects to:
//	│       └── nginx
//	└── Network: backend_net
//	    └── connects to:
//	        ├── postgres
//	        └── redis
//
// Parameters:
//   - w: The io.Writer to write the output to
//   - c: Pointer to the ContainerInfo for the container being displayed
//   - netMap: Map of network names to slices of ContainerInfo for containers on each network
func PrintContainerTree(w io.Writer, c *models.ContainerInfo, netMap map[string][]models.ContainerInfo) {
	cw := NewColorWriter(w)

	fmt.Fprintf(w, "%s %s\n", cw.Label("Container:"), cw.Container(c.Name))

	// Sort networks for consistent output
	sortedNetworks := make([]string, len(c.Networks))
	copy(sortedNetworks, c.Networks)
	sort.Strings(sortedNetworks)

	for i, net := range sortedNetworks {
		prefix := TreeBranch
		indent := TreeVertical
		if i == len(sortedNetworks)-1 {
			prefix = TreeEnd
			indent = TreeSpace
		}

		fmt.Fprintf(w, "%s %s %s\n", cw.Tree(prefix), cw.Label("Network:"), cw.Network(net))
		fmt.Fprintf(w, "%s%s %s\n", cw.Tree(indent), cw.Tree(TreeEnd), cw.Label("connects to:"))

		others := ReachableContainers(c.Name, net, netMap)
		if len(others) == 0 {
			fmt.Fprintf(w, "%s    %s (none)\n", cw.Tree(indent), cw.Tree(TreeEnd))
			continue
		}

		for j, o := range others {
			op := TreeBranch
			if j == len(others)-1 {
				op = TreeEnd
			}
			fmt.Fprintf(w, "%s    %s %s\n", cw.Tree(indent), cw.Tree(op), cw.Container(o))
		}
	}
}

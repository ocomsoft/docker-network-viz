// Package output provides tree-style formatters for Docker network topology visualization.
// It includes functions to print network trees, container reachability trees,
// and calculate container connectivity across networks.
package output

import (
	"sort"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// ReachableContainers returns a sorted list of container names that can be reached
// from a container on the specified network. It excludes the source container itself
// from the results.
//
// Parameters:
//   - self: The name of the source container (will be excluded from results)
//   - network: The network name to check for reachable containers
//   - netMap: A map of network names to slices of ContainerInfo for containers on that network
//
// Returns a sorted slice of container names that share the same network as the source
// container, excluding the source container itself. Returns an empty slice if no other
// containers are found on the network.
func ReachableContainers(self, network string, netMap map[string][]models.ContainerInfo) []string {
	var result []string
	for _, c := range netMap[network] {
		if c.Name != self {
			result = append(result, c.Name)
		}
	}
	sort.Strings(result)
	return result
}

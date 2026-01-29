// Package models provides data structures for representing Docker network topology.
// These models are used to store and manipulate information about containers,
// networks, and their relationships for visualization purposes.
package models

import (
	"sort"
)

// ContainerInfo represents a Docker container's network-related information.
// It stores the container's name, network aliases, and the networks it belongs to.
// This struct is used for building network topology views and determining
// container reachability across networks.
type ContainerInfo struct {
	// Name is the container's name without the leading slash.
	// Example: "web_app" not "/web_app"
	Name string

	// Aliases are the network-scoped aliases assigned to this container.
	// Aliases allow containers to be discovered by alternative names within a network.
	Aliases []string

	// Networks contains the names of all networks this container is connected to.
	// A container can be connected to multiple networks simultaneously.
	Networks []string
}

// NewContainerInfo creates a new ContainerInfo with the given name.
// The Aliases and Networks slices are initialized as empty slices.
func NewContainerInfo(name string) *ContainerInfo {
	return &ContainerInfo{
		Name:     name,
		Aliases:  []string{},
		Networks: []string{},
	}
}

// AddAlias adds a network alias to the container if it doesn't already exist.
// Returns true if the alias was added, false if it already existed.
func (c *ContainerInfo) AddAlias(alias string) bool {
	for _, existing := range c.Aliases {
		if existing == alias {
			return false
		}
	}
	c.Aliases = append(c.Aliases, alias)
	return true
}

// AddNetwork adds a network name to the container if it doesn't already exist.
// Returns true if the network was added, false if it already existed.
func (c *ContainerInfo) AddNetwork(network string) bool {
	for _, existing := range c.Networks {
		if existing == network {
			return false
		}
	}
	c.Networks = append(c.Networks, network)
	return true
}

// HasNetwork checks if the container is connected to the specified network.
func (c *ContainerInfo) HasNetwork(network string) bool {
	for _, n := range c.Networks {
		if n == network {
			return true
		}
	}
	return false
}

// HasAlias checks if the container has the specified alias.
func (c *ContainerInfo) HasAlias(alias string) bool {
	for _, a := range c.Aliases {
		if a == alias {
			return true
		}
	}
	return false
}

// SortedNetworks returns a copy of the Networks slice sorted alphabetically.
// This is useful for consistent output when displaying network information.
func (c *ContainerInfo) SortedNetworks() []string {
	sorted := make([]string, len(c.Networks))
	copy(sorted, c.Networks)
	sort.Strings(sorted)
	return sorted
}

// SortedAliases returns a copy of the Aliases slice sorted alphabetically.
// This is useful for consistent output when displaying alias information.
func (c *ContainerInfo) SortedAliases() []string {
	sorted := make([]string, len(c.Aliases))
	copy(sorted, c.Aliases)
	sort.Strings(sorted)
	return sorted
}

// NetworkCount returns the number of networks this container is connected to.
func (c *ContainerInfo) NetworkCount() int {
	return len(c.Networks)
}

// AliasCount returns the number of aliases this container has.
func (c *ContainerInfo) AliasCount() int {
	return len(c.Aliases)
}

// Clone creates a deep copy of the ContainerInfo.
// This is useful when you need to modify container information
// without affecting the original.
func (c *ContainerInfo) Clone() *ContainerInfo {
	aliases := make([]string, len(c.Aliases))
	copy(aliases, c.Aliases)

	networks := make([]string, len(c.Networks))
	copy(networks, c.Networks)

	return &ContainerInfo{
		Name:     c.Name,
		Aliases:  aliases,
		Networks: networks,
	}
}

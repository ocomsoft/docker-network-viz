// Package models provides data structures for docker-network-viz.
package models

// NetworkInfo represents a Docker network's basic information.
// It stores the network's name and driver type for visualization purposes.
// This struct is used to decouple the output package from Docker API types.
type NetworkInfo struct {
	// Name is the network's name.
	// Example: "bridge", "frontend_net", "backend_net"
	Name string

	// Driver is the network driver type.
	// Common values: "bridge", "host", "overlay", "macvlan", "none"
	Driver string
}

// NewNetworkInfo creates a new NetworkInfo with the given name and driver.
func NewNetworkInfo(name, driver string) *NetworkInfo {
	return &NetworkInfo{
		Name:   name,
		Driver: driver,
	}
}

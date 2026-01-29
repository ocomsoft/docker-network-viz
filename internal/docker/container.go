// Package docker provides Docker client wrapper functionality.
package docker

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// ContainerListOptions provides options for filtering container lists.
type ContainerListOptions struct {
	// All includes stopped containers when set to true.
	// When false, only running containers are returned.
	All bool

	// Filters is a map of filter names to filter values.
	// Supported filters include: ancestor, before, expose, exited,
	// health, id, isolation, is-task, label, name, network, publish,
	// since, status, volume
	Filters map[string][]string
}

// FetchContainers retrieves all Docker containers from the daemon.
// It returns a slice of types.Container sorted alphabetically by name.
//
// The options parameter can be used to filter containers and control
// whether stopped containers are included.
func (c *Client) FetchContainers(ctx context.Context, opts *ContainerListOptions) ([]types.Container, error) {
	listOpts := container.ListOptions{
		All: true, // Default to all containers
	}

	if opts != nil {
		listOpts.All = opts.All
	}

	containers, err := c.cli.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %w", err)
	}

	// Sort containers alphabetically by name for consistent output
	sort.Slice(containers, func(i, j int) bool {
		nameI := sanitizeContainerName(containers[i].Names)
		nameJ := sanitizeContainerName(containers[j].Names)
		return nameI < nameJ
	})

	return containers, nil
}

// FetchContainerByID retrieves a specific Docker container by its ID.
// It returns detailed JSON information about the container.
func (c *Client) FetchContainerByID(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	containerJSON, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to inspect Docker container %s: %w", containerID, err)
	}

	return containerJSON, nil
}

// BuildContainerMap creates a map of container names to ContainerInfo structs.
// This provides quick lookup of container information by name.
func (c *Client) BuildContainerMap(containers []types.Container) map[string]*models.ContainerInfo {
	containerMap := make(map[string]*models.ContainerInfo, len(containers))

	for _, cont := range containers {
		name := sanitizeContainerName(cont.Names)
		ci := models.NewContainerInfo(name)

		// Add all networks and their aliases
		for netName, netSettings := range cont.NetworkSettings.Networks {
			ci.AddNetwork(netName)

			// Add aliases from network settings
			if netSettings != nil {
				for _, alias := range netSettings.Aliases {
					ci.AddAlias(alias)
				}
			}
		}

		containerMap[name] = ci
	}

	return containerMap
}

// BuildNetworkToContainersMap creates a mapping from network names to the
// containers connected to each network. This is essential for determining
// container reachability within networks.
//
// The returned map has network names as keys and slices of ContainerInfo
// as values. Each ContainerInfo contains the container's name, aliases,
// and the networks it belongs to.
func (c *Client) BuildNetworkToContainersMap(containers []types.Container) map[string][]models.ContainerInfo {
	// First build the container map to get complete ContainerInfo objects
	containerMap := c.BuildContainerMap(containers)

	// Now build the network to containers mapping
	networkToContainers := make(map[string][]models.ContainerInfo)

	for _, cont := range containers {
		name := sanitizeContainerName(cont.Names)
		ci := containerMap[name]

		for netName := range cont.NetworkSettings.Networks {
			// Dereference the pointer to store a copy in the map
			networkToContainers[netName] = append(networkToContainers[netName], *ci)
		}
	}

	// Sort containers within each network for consistent output
	for netName, containerList := range networkToContainers {
		sort.Slice(containerList, func(i, j int) bool {
			return containerList[i].Name < containerList[j].Name
		})
		networkToContainers[netName] = containerList
	}

	return networkToContainers
}

// ConvertToContainerInfo converts a Docker types.Container to our internal
// ContainerInfo model. This decouples the output package from Docker API types.
func ConvertToContainerInfo(cont types.Container) *models.ContainerInfo {
	name := sanitizeContainerName(cont.Names)
	ci := models.NewContainerInfo(name)

	for netName, netSettings := range cont.NetworkSettings.Networks {
		ci.AddNetwork(netName)

		if netSettings != nil {
			for _, alias := range netSettings.Aliases {
				ci.AddAlias(alias)
			}
		}
	}

	return ci
}

// ConvertContainersToContainerInfos converts a slice of Docker types.Container
// to a slice of internal ContainerInfo models.
func ConvertContainersToContainerInfos(containers []types.Container) []*models.ContainerInfo {
	result := make([]*models.ContainerInfo, len(containers))
	for i, cont := range containers {
		result[i] = ConvertToContainerInfo(cont)
	}
	return result
}

// sanitizeContainerName removes the leading slash from container names.
// Docker container names are stored with a leading "/" which we remove
// for cleaner display and consistency.
func sanitizeContainerName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

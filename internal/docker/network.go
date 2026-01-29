// Package docker provides Docker client wrapper functionality.
package docker

import (
	"context"
	"fmt"
	"sort"

	"github.com/docker/docker/api/types/network"

	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

// NetworkListOptions provides options for filtering network lists.
type NetworkListOptions struct {
	// Filters is a map of filter names to filter values.
	// Supported filters include: driver, id, label, name, scope, type
	Filters map[string][]string
}

// FetchNetworks retrieves all Docker networks from the daemon.
// It returns a slice of network.Summary sorted alphabetically by name.
//
// The options parameter can be used to filter networks by various criteria.
// Pass nil or an empty NetworkListOptions for no filtering.
func (c *Client) FetchNetworks(ctx context.Context, opts *NetworkListOptions) ([]network.Summary, error) {
	listOpts := network.ListOptions{}
	if opts != nil && opts.Filters != nil {
		// Convert our filters to the Docker SDK filter format
		// The Docker SDK expects filters.Args which we build from our map
		for _, driver := range opts.Filters["driver"] {
			listOpts.Filters.Add("driver", driver)
		}
	}

	networks, err := c.cli.NetworkList(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker networks: %w", err)
	}

	// Sort networks alphabetically by name for consistent output
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})

	return networks, nil
}

// FetchNetworkByID retrieves a specific Docker network by its ID.
// It returns detailed information about the network including connected containers.
func (c *Client) FetchNetworkByID(ctx context.Context, networkID string) (network.Inspect, error) {
	net, err := c.cli.NetworkInspect(ctx, networkID, network.InspectOptions{})
	if err != nil {
		return network.Inspect{}, fmt.Errorf("failed to inspect Docker network %s: %w", networkID, err)
	}

	return net, nil
}

// FetchNetworkByName retrieves a specific Docker network by its name.
// It first lists networks filtered by name, then returns the matching network.
func (c *Client) FetchNetworkByName(ctx context.Context, name string) (network.Inspect, error) {
	net, err := c.cli.NetworkInspect(ctx, name, network.InspectOptions{})
	if err != nil {
		return network.Inspect{}, fmt.Errorf("failed to inspect Docker network %s: %w", name, err)
	}

	return net, nil
}

// ConvertToNetworkInfo converts a Docker network.Summary to our internal NetworkInfo model.
// This decouples the output package from Docker API types.
func ConvertToNetworkInfo(net network.Summary) *models.NetworkInfo {
	return models.NewNetworkInfo(net.Name, net.Driver)
}

// ConvertNetworksToNetworkInfos converts a slice of Docker network summaries
// to a slice of internal NetworkInfo models.
func ConvertNetworksToNetworkInfos(networks []network.Summary) []*models.NetworkInfo {
	result := make([]*models.NetworkInfo, len(networks))
	for i, net := range networks {
		result[i] = ConvertToNetworkInfo(net)
	}
	return result
}

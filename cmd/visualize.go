// Package cmd provides the CLI commands for the docker-network-viz tool.
// This file contains the visualize command which displays Docker network topology.
package cmd

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/docker/docker/api/types/network"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.o.ocom.com.au/go/docker-network-viz/internal/docker"
	"git.o.ocom.com.au/go/docker-network-viz/internal/models"
	"git.o.ocom.com.au/go/docker-network-viz/internal/output"
)

var (
	// onlyNetwork filters output to show only the specified network.
	onlyNetwork string

	// containerFilter filters output to show only the specified container.
	containerFilter string

	// noAliases disables the display of container aliases.
	noAliases bool

	// visualizeCmd represents the visualize command.
	visualizeCmd = &cobra.Command{
		Use:   "visualize",
		Short: "Display Docker network topology",
		Long: `Visualize Docker network topology in a tree-style format.

This command displays two views:
1. Network tree: Shows each network with its connected containers and aliases
2. Container reachability: Shows each container with the networks it belongs to
   and other containers it can reach through those networks

Examples:
  # Show all networks and containers
  docker-network-viz visualize

  # Show only a specific network
  docker-network-viz visualize --only-network bridge

  # Show only a specific container's connectivity
  docker-network-viz visualize --container web_app

  # Hide container aliases
  docker-network-viz visualize --no-aliases`,
		RunE: runVisualize,
	}
)

func init() {
	// Add visualize command to root
	rootCmd.AddCommand(visualizeCmd)

	// Local flags for visualize command
	visualizeCmd.Flags().StringVar(&onlyNetwork, "only-network", "",
		"show only the specified network")
	visualizeCmd.Flags().StringVar(&containerFilter, "container", "",
		"show only the specified container's connectivity")
	visualizeCmd.Flags().BoolVar(&noAliases, "no-aliases", false,
		"hide container aliases in the output")

	// Bind flags to viper
	_ = viper.BindPFlag("only-network", visualizeCmd.Flags().Lookup("only-network"))
	_ = viper.BindPFlag("container", visualizeCmd.Flags().Lookup("container"))
	_ = viper.BindPFlag("no-aliases", visualizeCmd.Flags().Lookup("no-aliases"))
}

// runVisualize executes the visualize command logic.
// It fetches Docker networks and containers, then prints the network topology
// in a tree-style format.
func runVisualize(cmd *cobra.Command, _ []string) error {
	ctx := context.Background()

	// Initialize Docker client
	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer func() {
		_ = client.Close()
	}()

	// Fetch networks
	networks, err := client.FetchNetworks(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch networks: %w", err)
	}

	// Fetch containers
	containers, err := client.FetchContainers(ctx, &docker.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to fetch containers: %w", err)
	}

	// Build mappings
	containerMap := client.BuildContainerMap(containers)
	networkToContainers := client.BuildNetworkToContainersMap(containers)

	// Get output writer
	writer := cmd.OutOrStdout()

	// Apply filters and print output
	return printVisualization(writer, networks, containerMap, networkToContainers)
}

// printVisualization handles the actual output of the network topology.
// It respects the command flags for filtering and formatting.
func printVisualization(
	w io.Writer,
	networks []network.Summary,
	containerMap map[string]*models.ContainerInfo,
	networkToContainers map[string][]models.ContainerInfo,
) error {
	onlyNetworkFlag := viper.GetString("only-network")
	containerFlag := viper.GetString("container")
	noAliasesFlag := viper.GetBool("no-aliases")

	// Print network tree section
	fmt.Fprintln(w, "=== Networks ===")

	for _, net := range networks {
		// Filter by network name if specified
		if onlyNetworkFlag != "" && net.Name != onlyNetworkFlag {
			continue
		}

		netInfo := models.NewNetworkInfo(net.Name, net.Driver)
		netContainers := networkToContainers[net.Name]

		// Apply alias filtering if needed
		if noAliasesFlag {
			netContainers = removeAliasesFromContainers(netContainers)
		}

		output.PrintNetworkTree(w, *netInfo, netContainers)
		fmt.Fprintln(w)
	}

	// Print container reachability section
	fmt.Fprintln(w, "=== Containers (Reachability) ===")

	// Sort container names for consistent output
	containerNames := make([]string, 0, len(containerMap))
	for name := range containerMap {
		containerNames = append(containerNames, name)
	}
	sort.Strings(containerNames)

	for _, name := range containerNames {
		// Filter by container name if specified
		if containerFlag != "" && name != containerFlag {
			continue
		}

		container := containerMap[name]
		output.PrintContainerTree(w, container, networkToContainers)
		fmt.Fprintln(w)
	}

	return nil
}

// removeAliasesFromContainers creates a copy of the container list with aliases removed.
// This is used when the --no-aliases flag is set.
func removeAliasesFromContainers(containers []models.ContainerInfo) []models.ContainerInfo {
	result := make([]models.ContainerInfo, len(containers))
	for i, c := range containers {
		result[i] = models.ContainerInfo{
			Name:     c.Name,
			Aliases:  []string{}, // Empty aliases
			Networks: c.Networks,
		}
	}
	return result
}

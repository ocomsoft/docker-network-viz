# Docker Client Package

The `docker` package provides a wrapper around the Docker SDK client for fetching network and container information. It handles client initialization with proper error handling.

## Overview

This package contains three main components:

1. **client.go** - Docker client wrapper with initialization and lifecycle management
2. **network.go** - Network-related operations (list, inspect, convert)
3. **container.go** - Container-related operations (list, inspect, mapping functions)

## Usage

### Creating a Client

```go
import (
    "git.o.ocom.com.au/go/docker-network-viz/internal/docker"
)

// Create client with default options (reads from environment)
client, err := docker.NewClient()
if err != nil {
    return fmt.Errorf("failed to create Docker client: %w", err)
}
defer client.Close()

// Create client with a mock for testing
mockClient := &MockAPIClient{}
client, err := docker.NewClient(
    docker.WithDockerClient(mockClient),
)
```

### Fetching Networks

```go
ctx := context.Background()

// Fetch all networks
networks, err := client.FetchNetworks(ctx, nil)
if err != nil {
    return fmt.Errorf("failed to fetch networks: %w", err)
}

for _, net := range networks {
    fmt.Printf("Network: %s (%s)\n", net.Name, net.Driver)
}

// Fetch specific network by ID or name
network, err := client.FetchNetworkByID(ctx, "network_id")
network, err := client.FetchNetworkByName(ctx, "bridge")
```

### Fetching Containers

```go
ctx := context.Background()

// Fetch all containers (including stopped)
containers, err := client.FetchContainers(ctx, nil)
if err != nil {
    return fmt.Errorf("failed to fetch containers: %w", err)
}

// Fetch only running containers
opts := &docker.ContainerListOptions{All: false}
containers, err := client.FetchContainers(ctx, opts)

// Fetch specific container
containerJSON, err := client.FetchContainerByID(ctx, "container_id")
```

### Building Container Maps

```go
// Build a map of container name -> ContainerInfo
containerMap := client.BuildContainerMap(containers)

// Build network-to-containers mapping (essential for reachability)
networkToContainers := client.BuildNetworkToContainersMap(containers)

// Access containers on a specific network
for _, cont := range networkToContainers["bridge"] {
    fmt.Printf("Container: %s\n", cont.Name)
}
```

### Converting to Internal Models

```go
// Convert Docker types to internal models
networkInfo := docker.ConvertToNetworkInfo(networkSummary)
networkInfos := docker.ConvertNetworksToNetworkInfos(networks)

containerInfo := docker.ConvertToContainerInfo(container)
containerInfos := docker.ConvertContainersToContainerInfos(containers)
```

## Client Options

| Option | Description |
|--------|-------------|
| `WithDockerClient(client.APIClient)` | Injects a custom Docker API client (for testing) |

## Types

### Client

The main wrapper struct that provides all Docker operations.

```go
type Client struct {
    cli client.APIClient
}
```

### ContainerListOptions

Options for filtering container lists.

```go
type ContainerListOptions struct {
    All     bool              // Include stopped containers
    Filters map[string][]string // Filter by various criteria
}
```

### NetworkListOptions

Options for filtering network lists.

```go
type NetworkListOptions struct {
    Filters map[string][]string // Filter by driver, id, label, name, scope, type
}
```

## Methods

### Client Methods

| Method | Description |
|--------|-------------|
| `NewClient(opts ...ClientOption)` | Creates a new Docker client wrapper |
| `Ping(ctx)` | Checks if Docker daemon is accessible |
| `Close()` | Closes the client connection |
| `APIClient()` | Returns the underlying Docker API client |

### Network Methods

| Method | Description |
|--------|-------------|
| `FetchNetworks(ctx, opts)` | Lists all Docker networks |
| `FetchNetworkByID(ctx, id)` | Gets network details by ID |
| `FetchNetworkByName(ctx, name)` | Gets network details by name |
| `ConvertToNetworkInfo(net)` | Converts Docker network to internal model |
| `ConvertNetworksToNetworkInfos(nets)` | Bulk converts networks |

### Container Methods

| Method | Description |
|--------|-------------|
| `FetchContainers(ctx, opts)` | Lists all Docker containers |
| `FetchContainerByID(ctx, id)` | Gets container details by ID |
| `BuildContainerMap(containers)` | Creates name -> ContainerInfo map |
| `BuildNetworkToContainersMap(containers)` | Creates network -> containers mapping |
| `ConvertToContainerInfo(cont)` | Converts Docker container to internal model |
| `ConvertContainersToContainerInfos(conts)` | Bulk converts containers |

## Testing

The package includes comprehensive unit tests with mocked Docker responses:

```bash
go test -v ./internal/docker/...
```

All tests use a `mockAPIClient` struct that implements the `client.APIClient` interface, allowing for isolated testing without requiring a running Docker daemon.

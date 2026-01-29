# Output Package

The `output` package provides tree-style formatters for Docker network topology visualization. It includes functions to print network trees, container reachability trees, calculate container connectivity across networks, and support for colored terminal output.

## Overview

This package is responsible for rendering Docker network and container information in a human-readable tree format. It produces output that is:

- SSH-friendly (ASCII tree characters)
- Sorted alphabetically for consistency
- Color-coded for better readability (when terminal supports it)

## Files

| File | Description |
|------|-------------|
| `color.go` | Color support utilities and ColorWriter |
| `container_tree.go` | Container reachability tree formatter |
| `network_tree.go` | Network tree formatter |
| `reachability.go` | Container reachability calculations |
| `tree_symbols.go` | Tree drawing symbol constants |

## Color Support

The package provides automatic color detection and colored output through the `ColorWriter` type.

### Color Scheme

| Element | Color | Method |
|---------|-------|--------|
| Network names | Cyan (Bold) | `Network()` |
| Container names | Green | `Container()` |
| Aliases | Yellow | `Alias()` |
| Labels | Magenta | `Label()` |
| Tree characters | Blue | `Tree()` |

### ColorWriter

```go
// Create a new ColorWriter
cw := output.NewColorWriter(os.Stdout)

// Check if color is enabled
if cw.IsEnabled() {
    fmt.Println("Colors enabled")
}

// Use color methods
fmt.Println(cw.Network("bridge"))      // Cyan bold
fmt.Println(cw.Container("web_app"))   // Green
fmt.Println(cw.Alias("api.local"))     // Yellow
fmt.Println(cw.Label("Network:"))      // Magenta
fmt.Println(cw.Tree("├──"))            // Blue
```

Color is automatically disabled when:
- Output is not a terminal (piped or redirected)
- The `--no-color` flag is set (via Viper configuration)
- The `NO_COLOR` environment variable is set

## Functions

### PrintNetworkTree

Prints a tree-style representation of a Docker network and its connected containers.

```go
func PrintNetworkTree(w io.Writer, net models.NetworkInfo, containers []models.ContainerInfo)
```

**Parameters:**
- `w` - The io.Writer to write output to (e.g., os.Stdout)
- `net` - NetworkInfo containing network name and driver
- `containers` - Slice of ContainerInfo for containers on this network

**Example Output:**
```
Network: bridge (bridge)
├── web_app
│   ├── alias: web
│   └── alias: web.local
├── redis
│   └── alias: redis
└── postgres
    └── alias: db
```

### PrintContainerTree

Prints a tree-style representation of a container's network connectivity and reachability.

```go
func PrintContainerTree(w io.Writer, c *models.ContainerInfo, netMap map[string][]models.ContainerInfo)
```

**Parameters:**
- `w` - The io.Writer to write output to
- `c` - Pointer to ContainerInfo for the container being displayed
- `netMap` - Map of network names to containers on each network

**Example Output:**
```
Container: api
├── Network: frontend_net
│   └── connects to:
│       └── nginx
└── Network: backend_net
    └── connects to:
        ├── postgres
        └── redis
```

### ReachableContainers

Returns a sorted list of container names reachable from a container on a specific network.

```go
func ReachableContainers(self, network string, netMap map[string][]models.ContainerInfo) []string
```

**Parameters:**
- `self` - Name of the source container (excluded from results)
- `network` - Network name to check for reachable containers
- `netMap` - Map of network names to containers

**Returns:**
- Sorted slice of container names that share the network with the source container

## Tree Symbols

The package uses Unicode box-drawing characters for tree formatting, defined as constants:

| Constant | Value | Description |
|----------|-------|-------------|
| `TreeBranch` | `├──` | Branch (not last item) |
| `TreeEnd` | `└──` | End (last item) |
| `TreeVertical` | `│   ` | Vertical line continuation |
| `TreeSpace` | `    ` | Empty space for indentation |

## Usage Example

```go
package main

import (
    "os"

    "git.o.ocom.com.au/go/docker-network-viz/internal/models"
    "git.o.ocom.com.au/go/docker-network-viz/internal/output"
)

func main() {
    // Create network info
    net := models.NetworkInfo{Name: "backend", Driver: "bridge"}

    // Create container info
    containers := []models.ContainerInfo{
        {Name: "api", Aliases: []string{"api.local"}, Networks: []string{"backend"}},
        {Name: "db", Aliases: []string{"postgres"}, Networks: []string{"backend"}},
    }

    // Print network tree (with colors if terminal supports it)
    output.PrintNetworkTree(os.Stdout, net, containers)

    // Build network map for reachability
    netMap := map[string][]models.ContainerInfo{
        "backend": containers,
    }

    // Print container tree
    output.PrintContainerTree(os.Stdout, &containers[0], netMap)
}
```

## Sorting Behavior

All output is sorted alphabetically for consistent, predictable results:

- Networks are sorted alphabetically by name
- Containers are sorted alphabetically by name
- Aliases are sorted alphabetically
- Reachable containers are sorted alphabetically

This ensures the same input always produces the same output, making it suitable for:
- Automated testing
- Diff comparisons
- Documentation examples

## Testing

Run tests with:

```bash
go test -v ./internal/output/...
```

The test suite covers:
- Empty input handling
- Single and multiple items
- Sorting correctness
- Tree prefix correctness
- Immutability of input data
- Color writer functionality
- Terminal detection

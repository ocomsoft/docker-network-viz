# ContainerInfo Model

The `ContainerInfo` struct represents a Docker container's network-related information. It is the core data structure used throughout the docker-network-viz application for building network topology views and determining container reachability.

## Package

```go
import "git.o.ocom.com.au/go/docker-network-viz/internal/models"
```

## Struct Definition

```go
type ContainerInfo struct {
    Name     string
    Aliases  []string
    Networks []string
}
```

## Fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | The container's name without the leading slash (e.g., "web_app" not "/web_app") |
| `Aliases` | `[]string` | Network-scoped aliases assigned to the container for discovery |
| `Networks` | `[]string` | Names of all networks this container is connected to |

## Constructor

### NewContainerInfo

Creates a new `ContainerInfo` with the given name. Aliases and Networks are initialized as empty slices.

```go
func NewContainerInfo(name string) *ContainerInfo
```

**Example:**
```go
container := models.NewContainerInfo("web_app")
```

## Methods

### AddAlias

Adds a network alias to the container if it does not already exist.

```go
func (c *ContainerInfo) AddAlias(alias string) bool
```

**Returns:** `true` if the alias was added, `false` if it already existed.

**Example:**
```go
container := models.NewContainerInfo("web_app")
added := container.AddAlias("web")      // returns true
duplicate := container.AddAlias("web")  // returns false
```

### AddNetwork

Adds a network name to the container if it does not already exist.

```go
func (c *ContainerInfo) AddNetwork(network string) bool
```

**Returns:** `true` if the network was added, `false` if it already existed.

**Example:**
```go
container := models.NewContainerInfo("web_app")
container.AddNetwork("bridge")
container.AddNetwork("frontend")
```

### HasNetwork

Checks if the container is connected to the specified network.

```go
func (c *ContainerInfo) HasNetwork(network string) bool
```

**Example:**
```go
if container.HasNetwork("bridge") {
    // Container is on the bridge network
}
```

### HasAlias

Checks if the container has the specified alias.

```go
func (c *ContainerInfo) HasAlias(alias string) bool
```

**Example:**
```go
if container.HasAlias("web") {
    // Container has the "web" alias
}
```

### SortedNetworks

Returns a copy of the Networks slice sorted alphabetically. The original slice is not modified.

```go
func (c *ContainerInfo) SortedNetworks() []string
```

**Example:**
```go
networks := container.SortedNetworks()
// Returns: ["alpha", "beta", "zebra"] even if original order was different
```

### SortedAliases

Returns a copy of the Aliases slice sorted alphabetically. The original slice is not modified.

```go
func (c *ContainerInfo) SortedAliases() []string
```

### NetworkCount

Returns the number of networks this container is connected to.

```go
func (c *ContainerInfo) NetworkCount() int
```

### AliasCount

Returns the number of aliases this container has.

```go
func (c *ContainerInfo) AliasCount() int
```

### Clone

Creates a deep copy of the ContainerInfo. Useful when you need to modify container information without affecting the original.

```go
func (c *ContainerInfo) Clone() *ContainerInfo
```

**Example:**
```go
original := models.NewContainerInfo("web")
original.AddNetwork("bridge")

clone := original.Clone()
clone.Name = "modified"  // Does not affect original
```

## Usage Example

```go
package main

import (
    "fmt"
    "git.o.ocom.com.au/go/docker-network-viz/internal/models"
)

func main() {
    // Create a new container
    container := models.NewContainerInfo("api_server")

    // Add networks
    container.AddNetwork("frontend")
    container.AddNetwork("backend")

    // Add aliases
    container.AddAlias("api")
    container.AddAlias("api.local")

    // Check connectivity
    if container.HasNetwork("backend") {
        fmt.Printf("%s can reach backend services\n", container.Name)
    }

    // Display sorted networks for consistent output
    for _, net := range container.SortedNetworks() {
        fmt.Printf("  Network: %s\n", net)
    }
}
```

## Testing

The ContainerInfo model has comprehensive unit tests covering:

- Constructor functionality
- Adding aliases (new and duplicates)
- Adding networks (new and duplicates)
- Checking for network membership
- Checking for alias existence
- Sorted retrieval of networks and aliases
- Count methods
- Deep cloning
- Direct field access

Run tests with:
```bash
go test -v ./internal/models/...
```

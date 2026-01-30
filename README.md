# docker-network-viz

A command-line tool for visualizing Docker network topology in a tree-style format.

## Overview

docker-network-viz connects to the Docker daemon and displays:

- **Network-centric view**: Lists all networks with their connected containers and aliases
- **Container-centric view**: Shows each container's network memberships and reachability to other containers

The tool provides colored output for better readability when run in a terminal and can be installed as a Docker CLI plugin.

## Installation

### Using go install

```bash
go install git.o.ocom.com.au/go/docker-network-viz@latest
```

### Build from source

```bash
git clone git.o.ocom.com.au/go/docker-network-viz.git
cd docker-network-viz
make build
```

### Install as Docker Plugin

To use this tool as a Docker CLI plugin (e.g., `docker network-viz`), build the binary and place it in one of these directories:

1. **`~/.docker/cli-plugins/`** - recommended for user-specific installation
2. **`/usr/local/lib/docker/cli-plugins/`** - system-wide installation
3. **`/usr/lib/docker/cli-plugins/`** - alternative system-wide location

Using the Makefile:

```bash
# Install to user's Docker CLI plugins directory
make install

# Install system-wide (requires sudo)
make install-system
```

Or manually:

```bash
# Build the binary
make build

# Make it executable (already done by make build)
chmod +x docker-network-viz

# Install to Docker CLI plugins directory (user-specific)
mkdir -p ~/.docker/cli-plugins/
cp docker-network-viz ~/.docker/cli-plugins/

# Verify installation
docker network-viz
```

### Using Docker

Build and run using Docker (see [DOCKER.md](DOCKER.md) for detailed Docker usage):

```bash
# Build the Docker image
./docker-build.sh

# Run with Docker
./docker-run.sh

# Or manually
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock docker-network-viz:latest

# With flags
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock docker-network-viz:latest --no-color
```

**Note**: The Docker socket must be mounted to access the Docker daemon. If you encounter permission errors, see the [Docker troubleshooting guide](DOCKER.md#troubleshooting).

### Quick Demo

To see the tool in action with sample containers:

```bash
./docker-demo.sh
```

This will create demo networks and containers to visualize.

## Usage

The tool can be run directly without any subcommands to display the full network visualization:

```bash
# Display all networks and containers (default behavior)
docker-network-viz

# Or as a Docker plugin
docker network-viz
```

### Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to configuration file | `$HOME/.docker-network-viz.yaml` |
| `--no-color` | Disable colored output | `false` |
| `--only-network` | Show only the specified network | (all networks) |
| `--container` | Show only the specified container's connectivity | (all containers) |
| `--no-aliases` | Hide container aliases in the output | `false` |

### Examples

```bash
# Show all networks and containers with colored output
docker-network-viz

# Disable colored output (useful for piping or logging)
docker-network-viz --no-color

# Show only a specific network
docker-network-viz --only-network bridge

# Show only a specific container's connectivity
docker-network-viz --container web_app

# Hide container aliases for cleaner output
docker-network-viz --no-aliases

# Combine multiple flags
docker-network-viz --only-network frontend_net --no-aliases

# Use the explicit visualize subcommand
docker-network-viz visualize --only-network backend
```

### Environment Variables

Flags can also be set via environment variables with the `DNV_` prefix:

| Variable | Equivalent Flag |
|----------|-----------------|
| `DNV_NO_COLOR` | `--no-color` |
| `DNV_ONLY_NETWORK` | `--only-network` |
| `DNV_CONTAINER` | `--container` |
| `DNV_NO_ALIASES` | `--no-aliases` |

Example:

```bash
# Set environment variable to always disable color
export DNV_NO_COLOR=true
docker-network-viz
```

### Configuration File

Create a configuration file at `$HOME/.docker-network-viz.yaml` or `./.docker-network-viz.yaml`:

```yaml
no-color: false
only-network: ""
container: ""
no-aliases: false
```

## Output Format

### Colored Output

When running in a terminal, the output uses colors for better readability:

| Element | Color |
|---------|-------|
| Network names | **Cyan (Bold)** |
| Container names | Green |
| Aliases | Yellow |
| Labels (Network:, Container:, alias:, connects to:) | Magenta |
| Tree characters | Blue |

Color is automatically disabled when output is piped or redirected, or when the `--no-color` flag is set.

### Network Tree

The first section shows each network with its connected containers and aliases:

```
=== Networks ===
Network: bridge (bridge)
├── web_app
│   ├── alias: web
│   └── alias: web.local
├── redis
│   └── alias: redis
└── postgres
    └── alias: db

Network: frontend_net (bridge)
├── nginx
└── web_app
```

### Container Reachability Tree

The second section shows each container with the networks it belongs to and which containers it can reach through those networks:

```
=== Containers (Reachability) ===
Container: api
├── Network: frontend_net
│   └── connects to:
│       └── nginx
└── Network: backend_net
    └── connects to:
        ├── postgres
        └── redis

Container: nginx
└── Network: frontend_net
    └── connects to:
        └── api
```

This helps you quickly understand:
- Which containers can communicate with each other
- Through which networks the communication happens
- Whether a container is accidentally exposed on multiple networks

## Project Structure

```
docker-network-viz/
├── cmd/
│   └── docker-network-viz/    # CLI entry point
│       ├── main.go            # Main entry point
│       ├── root.go            # Root command with global flags
│       └── visualize.go       # Visualize command implementation
├── internal/
│   ├── docker/                # Docker client wrapper
│   │   ├── client.go          # Client initialization
│   │   ├── container.go       # Container operations
│   │   └── network.go         # Network operations
│   ├── models/                # Data structures
│   │   ├── container.go       # ContainerInfo model
│   │   └── network.go         # NetworkInfo model
│   └── output/                # Output formatters
│       ├── color.go           # Color support utilities
│       ├── container_tree.go  # Container tree formatter
│       ├── network_tree.go    # Network tree formatter
│       ├── reachability.go    # Reachability calculations
│       └── tree_symbols.go    # Tree drawing symbols
├── test/                      # Integration tests
├── Makefile                   # Build automation
├── go.mod
├── go.sum
└── README.md
```

## Development

### Prerequisites

- Go 1.24 or later
- Docker (for running the tool)
- goimports: `go install golang.org/x/tools/cmd/goimports@latest`
- golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make cover

# Run tests for a specific package
go test -v ./internal/docker/...
```

### Linting

```bash
# Run golangci-lint
make lint

# Run all quality checks (format, lint, test)
make check
```

### Building

```bash
# Build the binary
make build

# Clean build artifacts
make clean

# Clean and rebuild
make clean build
```

## Dependencies

- `github.com/docker/docker/client` - Docker API client
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/fatih/color` - Terminal color output
- `github.com/rs/zerolog` - Structured logging

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

Copyright Ocom. All rights reserved.

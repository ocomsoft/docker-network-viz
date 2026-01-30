# docker-network-viz CLI

This directory contains the Cobra CLI implementation for the docker-network-viz tool.

## Overview

The CLI is structured using the Cobra command library with Viper for configuration management. Each command is in its own source file following Ocom Go standards.

## File Structure

| File | Description |
|------|-------------|
| `main.go` | Entry point that executes the root command |
| `root.go` | Root command definition with global flags and Viper integration |
| `visualize.go` | The visualization command that displays network topology |

## Commands

### Root Command (Default Behavior)

The root command can be run without any subcommand to display the full network visualization. This is the primary way to use the tool.

```bash
# Run without subcommand (recommended)
docker-network-viz

# Or as Docker plugin
docker network-viz
```

**Global Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `$HOME/.docker-network-viz.yaml` |
| `--no-color` | Disable colored output | `false` |

**Visualization Flags (available on root and visualize commands):**

| Flag | Description | Default |
|------|-------------|---------|
| `--only-network` | Show only the specified network | (all networks) |
| `--container` | Show only the specified container's connectivity | (all containers) |
| `--no-aliases` | Hide container aliases in the output | `false` |

### Visualize Subcommand

The `visualize` command is also available as an explicit subcommand with the same functionality:

```bash
docker-network-viz visualize [flags]
```

This is equivalent to running without a subcommand.

## Usage Examples

```bash
# Show all networks and containers (default behavior)
docker-network-viz

# Disable colored output (useful for piping to files or other tools)
docker-network-viz --no-color

# Show only a specific network
docker-network-viz --only-network bridge

# Show only a specific container's connectivity
docker-network-viz --container web_app

# Hide container aliases for cleaner output
docker-network-viz --no-aliases

# Combine multiple flags
docker-network-viz --only-network frontend_net --no-aliases --no-color

# Using the explicit visualize subcommand
docker-network-viz visualize --only-network bridge
```

## Configuration

The CLI supports configuration via:

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefix: `DNV_`)
3. **Configuration file** (lowest priority)

### Environment Variables

Environment variables use the `DNV_` prefix with underscores replacing dashes:

| Variable | Equivalent Flag |
|----------|-----------------|
| `DNV_NO_COLOR` | `--no-color` |
| `DNV_ONLY_NETWORK` | `--only-network` |
| `DNV_CONTAINER` | `--container` |
| `DNV_NO_ALIASES` | `--no-aliases` |

Example:

```bash
# Always disable colors
export DNV_NO_COLOR=true
docker-network-viz

# Filter to specific network via environment
DNV_ONLY_NETWORK=bridge docker-network-viz
```

### Configuration File

The configuration file is searched for in the following locations:

1. Path specified by `--config` flag
2. `$HOME/.docker-network-viz.yaml`
3. `./.docker-network-viz.yaml`

**Example configuration file:**

```yaml
no-color: false
only-network: ""
container: ""
no-aliases: false
```

## Output Format

The visualization produces two sections:

### Network Tree Section

Shows each network with its connected containers and aliases:

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
```

### Container Reachability Section

Shows each container with networks it belongs to and other reachable containers:

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
```

## Colored Output

When running in a terminal, the output uses ANSI colors for better readability:

| Element | Color |
|---------|-------|
| Network names | Cyan (Bold) |
| Container names | Green |
| Aliases | Yellow |
| Labels | Magenta |
| Tree characters | Blue |

Color is automatically disabled when:
- Output is piped or redirected to a file
- The `--no-color` flag is set
- The `DNV_NO_COLOR` environment variable is set to `true`
- The `NO_COLOR` environment variable is set (standard convention)

## Testing

Run the CLI tests with:

```bash
go test -v ./cmd/...
```

The tests verify:
- Root command functionality and default behavior
- Flag parsing and validation
- Viper configuration binding
- Output formatting

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/fatih/color` - Terminal color output
- `github.com/rs/zerolog` - Structured logging

# Contributing to docker-network-viz

Thank you for your interest in contributing to docker-network-viz. This document provides guidelines and standards for contributing to the project.

## Development Environment Setup

### Prerequisites

1. **Go 1.24 or later**
   ```bash
   # Check your Go version
   go version
   ```

2. **Docker** - Required for running the tool and integration tests
   ```bash
   # Verify Docker is running
   docker ps
   ```

3. **Development tools**
   ```bash
   # Install goimports
   go install golang.org/x/tools/cmd/goimports@latest

   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Getting the Code

```bash
git clone git.o.ocom.com.au/go/docker-network-viz.git
cd docker-network-viz
make deps
```

## Code Standards

### File Organization

Following Ocom Go standards, each component should be in its own source file:

- **One Handler per file**
- **One Model per file**
- **One Cobra Command per file**
- **One Controller per file**

Smaller source files are preferred over larger ones.

### Code Quality

All code must:

1. **Pass golangci-lint with zero issues**
   ```bash
   make lint
   ```

2. **Be formatted with goimports**
   ```bash
   make fmt
   ```

3. **Have comprehensive unit tests**
   ```bash
   make test
   ```

### DRY Principle

Never duplicate code. Before writing new functionality:

1. Check if similar functionality already exists
2. Consider creating a generic function if code could be reused
3. Use the existing helper functions and utilities

## Development Workflow

### Before Making Changes

1. Create a feature branch
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Ensure tests pass on clean main
   ```bash
   make check
   ```

### Making Changes

1. **Write tests first** (TDD encouraged)

2. **Implement functionality**

3. **Run quality checks**
   ```bash
   make check
   ```

4. **Ensure pristine test output**
   - Tests should not produce warnings
   - If errors are expected, capture and test them explicitly

### Code Documentation

- **Complex functions**: Require inline godoc comments
- **Simple getters/setters**: Do not require documentation
- **Public APIs**: Must have godoc comments describing parameters and return values
- **Package documentation**: Each package should have a package comment

Example of good function documentation:

```go
// PrintNetworkTree prints a tree-style representation of a Docker network
// and its connected containers to the provided writer.
//
// The output format shows the network name and driver, followed by a tree
// of containers connected to that network. Each container's aliases are
// shown as nested items beneath the container name.
//
// Parameters:
//   - w: The io.Writer to write the output to
//   - net: The NetworkInfo containing the network name and driver
//   - containers: Slice of ContainerInfo for containers connected to this network
func PrintNetworkTree(w io.Writer, net models.NetworkInfo, containers []models.ContainerInfo) {
    // ...
}
```

## Testing Requirements

### Unit Tests

Every piece of functionality must have comprehensive unit tests:

```bash
# Run all unit tests
make test

# Run tests for a specific package
go test -v ./internal/docker/...

# Run tests with coverage
make cover
```

### Test Output

- **Test output must be pristine** - no warnings, no errors, no noise
- If tests intentionally produce errors, capture and verify them
- Use table-driven tests for better coverage

### Integration Tests

Integration tests are located in the `test/` directory and require a running Docker daemon:

```bash
# Run integration tests (requires Docker)
go test -v ./test/...
```

## Commit Guidelines

### Commit Messages

Write clear, descriptive commit messages:

```
Add --only-network flag for filtering output

- Implement flag parsing in root.go and visualize.go
- Add filtering logic in printVisualization function
- Include unit tests for the new functionality
```

### What to Commit

- Do not commit generated files (binaries, coverage reports)
- Do not commit IDE-specific files
- See `.gitignore` for excluded files

### Before Committing

Always run:

```bash
make check
```

This runs formatting, linting, and all tests.

## Pull Request Process

1. **Ensure all tests pass**
   ```bash
   make check
   ```

2. **Update documentation** if you changed functionality

3. **Create a pull request** with:
   - Clear description of changes
   - Reference to any related issues
   - Screenshots or examples if UI/output changed

4. **Address review feedback** promptly

## Project Structure

```
docker-network-viz/
├── cmd/docker-network-viz/    # CLI commands (one file per command)
│   ├── main.go                # Entry point
│   ├── root.go                # Root command with global flags
│   └── visualize.go           # Visualize command
├── internal/
│   ├── docker/                # Docker client wrapper
│   │   ├── client.go          # Client initialization
│   │   ├── container.go       # Container operations
│   │   └── network.go         # Network operations
│   ├── models/                # Data structures
│   │   ├── container.go       # ContainerInfo model
│   │   └── network.go         # NetworkInfo model
│   └── output/                # Output formatters
│       ├── color.go           # Color support
│       ├── container_tree.go  # Container tree formatter
│       ├── network_tree.go    # Network tree formatter
│       ├── reachability.go    # Reachability calculations
│       └── tree_symbols.go    # Tree drawing symbols
├── test/                      # Integration tests
├── Makefile                   # Build automation
├── .golangci.yml              # Linter configuration
└── README.md                  # Project documentation
```

## Technology Stack

This project uses the following libraries (per Ocom standards):

| Purpose | Library |
|---------|---------|
| CLI Framework | github.com/spf13/cobra |
| Configuration | github.com/spf13/viper |
| Docker Client | github.com/docker/docker/client |
| Color Output | github.com/fatih/color |
| Logging | github.com/rs/zerolog |

Do not add new dependencies without discussion.

## Common Tasks

### Adding a New Command

1. Create a new file in `cmd/docker-network-viz/` (e.g., `newcmd.go`)
2. Define the command with Cobra
3. Register flags with Viper
4. Add the command to root in `init()`
5. Write tests in `newcmd_test.go`
6. Update documentation

### Adding a New Model

1. Create a new file in `internal/models/` (e.g., `newmodel.go`)
2. Add godoc comments
3. Implement any helper methods
4. Write tests in `newmodel_test.go`
5. Create markdown documentation if complex

### Adding Output Formatters

1. Create a new file in `internal/output/`
2. Follow the existing patterns for color support
3. Ensure sorted, consistent output
4. Write comprehensive tests

## Troubleshooting

### Common Issues

**golangci-lint fails:**
```bash
# Update golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run with verbose output
golangci-lint run -v ./...
```

**Tests fail with Docker errors:**
```bash
# Ensure Docker is running
docker ps

# Check Docker socket permissions
ls -la /var/run/docker.sock
```

**Imports not sorted:**
```bash
# Run goimports on all files
make fmt
```

## Getting Help

If you have questions or need help:

1. Check existing documentation in the repository
2. Look at existing code for patterns
3. Ask questions in pull request comments

## Code of Conduct

- Be respectful and constructive
- Focus on the code, not the person
- Welcome newcomers and help them learn

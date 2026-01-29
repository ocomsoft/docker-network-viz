# Makefile Documentation

This document describes the Makefile targets available for building, testing, and maintaining the docker-network-viz project.

## Overview

The Makefile provides standardized commands for common development tasks following Ocom Go standards. It integrates with golangci-lint for code quality and supports installation as a Docker CLI plugin.

## Prerequisites

Before using the Makefile, ensure you have the following tools installed:

- **Go**: Version 1.24 or later
- **goimports**: `go install golang.org/x/tools/cmd/goimports@latest`
- **golangci-lint**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Targets

### help

Display all available Makefile targets with descriptions.

```bash
make help
```

### build

Compile the binary to `docker-network-viz` in the project root directory.

```bash
make build
```

The binary is built with optimizations (`-s -w` flags) to reduce binary size.

### test

Run all tests with race detection and coverage reporting.

```bash
make test
```

This runs `go test -v -race -cover ./...` to execute all tests across all packages.

### lint

Run golangci-lint and goimports to check code quality.

```bash
make lint
```

This target first formats code with goimports, then runs golangci-lint with the configuration from `.golangci.yml`.

### fmt

Format code with goimports.

```bash
make fmt
```

### install

Install the binary to the user's Docker CLI plugins directory (`~/.docker/cli-plugins/`).

```bash
make install
```

After installation, you can use:
```bash
docker network-viz
```

### install-system

Install the binary to the system-wide Docker CLI plugins directory (`/usr/local/lib/docker/cli-plugins/`). Requires sudo.

```bash
make install-system
```

### clean

Remove build artifacts including the compiled binary and coverage files.

```bash
make clean
```

### tidy

Tidy and verify Go modules.

```bash
make tidy
```

### deps

Download all Go dependencies.

```bash
make deps
```

### cover

Run tests with coverage and generate an HTML coverage report.

```bash
make cover
```

Opens `coverage.html` with detailed coverage information.

### check

Run all quality checks: format, lint, and test.

```bash
make check
```

### all

Build and test the project (combines lint, test, and build).

```bash
make all
```

## Common Workflows

### Development

```bash
# Format, lint, and test
make check

# Build and install locally
make install
```

### Continuous Integration

```bash
# Run all quality checks
make lint test

# Build for release
make build
```

### Clean Build

```bash
make clean build
```

## Configuration

The Makefile uses the following configuration files:

- `.golangci.yml`: Configuration for golangci-lint
- `go.mod`: Go module definition

## Environment Variables

The Makefile respects standard Go environment variables:

- `GOCACHE`: Go build cache directory
- `GOPATH`: Go workspace path

## File Locations

| File | Description |
|------|-------------|
| `docker-network-viz` | Compiled binary (project root) |
| `~/.docker/cli-plugins/docker-network-viz` | User installation location |
| `/usr/local/lib/docker/cli-plugins/docker-network-viz` | System installation location |
| `coverage.out` | Test coverage data |
| `coverage.html` | HTML coverage report |

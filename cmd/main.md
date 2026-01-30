# main.go

## Overview

The `main.go` file serves as the entry point for the docker-network-viz CLI tool. This tool provides a tree-style visualization of Docker network topology, showing the relationships between networks, containers, and their connectivity.

## Purpose

- Initialize logging with zerolog
- Configure and execute the Cobra root command
- Serve as the primary executable entry point

## Dependencies

- `github.com/rs/zerolog` - Structured logging
- `github.com/spf13/cobra` - CLI framework (to be implemented)

## Usage

Build and run the application:

```bash
go build -o docker-network-viz .
./docker-network-viz
```

## Future Enhancements

The main function will be expanded to:

1. Initialize the Cobra root command
2. Configure global flags (--verbose, --json, etc.)
3. Execute the command tree
4. Handle graceful shutdown

## Related Files

- `../../internal/docker/` - Docker client wrapper
- `../../internal/models/` - Data structures
- `../../internal/output/` - Output formatters

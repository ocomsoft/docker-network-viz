# Integration Tests

This directory contains comprehensive integration tests for docker-network-viz.

## Test Files

### integration_test.go

End-to-end integration tests that verify the complete flow of fetching Docker data and generating visualization output.

**Tests included:**

| Test Name | Description |
|-----------|-------------|
| `TestIntegration_FetchAndVisualize` | Tests full flow of fetching networks/containers and building mappings |
| `TestIntegration_NetworkTreeOutput` | Verifies network tree output format contains expected elements |
| `TestIntegration_ContainerTreeOutput` | Verifies container tree output format contains expected elements |
| `TestIntegration_ReachabilityCalculation` | Tests container reachability calculation across networks |
| `TestIntegration_EmptyNetwork` | Tests handling of networks with no containers |
| `TestIntegration_ContainerWithNoReachability` | Tests container with no reachable peers |
| `TestIntegration_ContainerMapBuilding` | Tests container map building process |
| `TestIntegration_SortedOutput` | Verifies output is consistently sorted alphabetically |

### cli_test.go

CLI command execution tests that verify command-line interface behavior.

**Tests included:**

| Test Name | Description |
|-----------|-------------|
| `TestCLI_HelpCommand` | Tests `--help` flag output contains expected elements |
| `TestCLI_VisualizeHelpCommand` | Tests `visualize --help` output |
| `TestCLI_UnknownFlag` | Verifies unknown flags produce errors |
| `TestCLI_InvalidConfigFile` | Tests graceful handling of non-existent config files |
| `TestCLI_NoColorFlag` | Tests `--no-color` flag acceptance |
| `TestCLI_VisualizeFlags` | Tests all visualize command flags |
| `TestCLI_RootFlagsPassthrough` | Tests root command accepts visualize flags |

### error_handling_test.go

Error handling tests that verify proper error responses for various failure scenarios.

**Tests included:**

| Test Name | Description |
|-----------|-------------|
| `TestError_NetworkListFailure` | Tests handling of network list errors |
| `TestError_ContainerListFailure` | Tests handling of container list errors |
| `TestError_PingFailure` | Tests handling of Docker daemon ping errors |
| `TestError_NetworkInspectFailure` | Tests handling of network inspect errors |
| `TestError_CloseFailure` | Tests handling of client close errors |
| `TestError_EmptyContainerNames` | Tests handling of containers with empty names |
| `TestError_NilNetworkSettings` | Tests handling of empty network settings |
| `TestError_ContextCancellation` | Tests handling of context cancellation |
| `TestError_PartialDataRecovery` | Tests partial data processing when some operations fail |

### output_format_test.go

Output format tests that verify the structure and formatting of generated output.

**Tests included:**

| Test Name | Description |
|-----------|-------------|
| `TestOutputFormat_NetworkTreeStructure` | Verifies network tree output structure |
| `TestOutputFormat_ContainerTreeStructure` | Verifies container tree output structure |
| `TestOutputFormat_TreeSymbols` | Tests correct tree symbols are used |
| `TestOutputFormat_AliasDisplay` | Tests alias display formatting |
| `TestOutputFormat_ConnectsToDisplay` | Tests "connects to:" section formatting |
| `TestOutputFormat_NoContainersMessage` | Tests empty network message |
| `TestOutputFormat_NoReachableContainersMessage` | Tests isolated container message |
| `TestOutputFormat_MultipleNetworksPerContainer` | Tests multi-homed container display |
| `TestOutputFormat_SortedAliases` | Verifies aliases are sorted alphabetically |
| `TestOutputFormat_SortedNetworksInContainerTree` | Verifies networks are sorted |
| `TestOutputFormat_LongNames` | Tests handling of long container/network names |
| `TestOutputFormat_SpecialCharactersInNames` | Tests special characters in names |

## Running Tests

Run all integration tests:

```bash
go test -v ./test/...
```

Run specific test file:

```bash
go test -v ./test/integration_test.go
go test -v ./test/cli_test.go
go test -v ./test/error_handling_test.go
go test -v ./test/output_format_test.go
```

Run specific test:

```bash
go test -v ./test/... -run TestIntegration_FetchAndVisualize
```

## Mock Data

Tests use mock Docker API clients to simulate various scenarios without requiring a running Docker daemon. The mock clients are defined in each test file and implement the `client.APIClient` interface.

### Mock Container Data

The integration tests use a predefined set of mock containers:

- `web_app` - Frontend container on `frontend_net`
- `api` - Multi-homed container on `frontend_net` and `backend_net`
- `postgres` - Database container on `backend_net`
- `redis` - Cache container on `backend_net`

### Mock Network Data

The integration tests use three mock networks:

- `bridge` - Default bridge network
- `frontend_net` - Frontend tier network
- `backend_net` - Backend tier network

## Test Coverage

These integration tests cover:

1. **End-to-end flow** - Complete data fetching and visualization pipeline
2. **CLI execution** - Command-line interface behavior and flag handling
3. **Output format** - Correct tree structure, symbols, and formatting
4. **Error handling** - Graceful handling of various failure scenarios
5. **Edge cases** - Empty data, special characters, long names, isolated containers

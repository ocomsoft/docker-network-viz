# Docker Usage Guide

This document describes how to build and run docker-network-viz using Docker.

## Overview

The docker-network-viz tool is packaged as a Docker container using a multi-stage build for minimal image size and security.

## Multi-Stage Dockerfile

The Dockerfile uses two stages:

1. **Builder Stage**: Uses `golang:1.22-alpine` to compile the Go binary with static linking
2. **Runtime Stage**: Uses minimal `alpine:latest` with only the compiled binary and CA certificates

### Image Features

- Minimal size (~15MB)
- Non-root user (`appuser`)
- Statically linked binary
- CA certificates for HTTPS support

## Building the Image

### Using the Build Script

```bash
./docker-build.sh
```

The script will:
- Build the Docker image as `docker-network-viz:latest`
- Display image details
- Show usage instructions

### Using Make

```bash
make docker-build
```

### Custom Image Name/Tag

```bash
IMAGE_NAME=myorg/network-viz IMAGE_TAG=v1.0.0 ./docker-build.sh
```

### Manual Build

```bash
docker build -t docker-network-viz:latest .
```

## Running the Container

### Using the Run Script

```bash
# Run with default settings
./docker-run.sh

# Pass flags to docker-network-viz (after --)
./docker-run.sh -- --no-color
./docker-run.sh -- --only-network bridge
./docker-run.sh -- --container web_app --no-aliases
```

### Using Make

```bash
make docker-run
```

### Manual Run

```bash
# Basic usage
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest

# With flags
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest --no-color --only-network bridge
```

## Important Notes

### Docker Socket Mounting

The container **must** have access to the Docker socket to inspect networks and containers:

```bash
-v /var/run/docker.sock:/var/run/docker.sock
```

This allows the container to communicate with the Docker daemon on the host.

### Security Considerations

1. **Socket Access**: Mounting the Docker socket gives the container full access to Docker. Only run this in trusted environments.

2. **Non-Root User**: The container runs as a non-root user (`appuser`, UID 1000) for security.

3. **Read-Only**: The tool only reads Docker data; it doesn't modify containers or networks.

### Permissions

The user running the container must have permission to access `/var/run/docker.sock`. This typically means:

- Being in the `docker` group on the host
- Running as root (not recommended)

## Examples

### Inspect a Specific Network

```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest --only-network my-custom-network
```

### Show Only a Specific Container

```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest --container web_app
```

### Disable Colors (for CI/CD)

```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest --no-color
```

### Hide Aliases

```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker-network-viz:latest --no-aliases
```

## Integration with Docker Compose

You can add docker-network-viz to a docker-compose.yml file:

```yaml
version: '3.8'

services:
  network-viz:
    image: docker-network-viz:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    command: ["--no-color"]
```

Run it:

```bash
docker-compose run --rm network-viz
```

## Troubleshooting

### Permission Denied Accessing Docker Socket

**Error**: `permission denied while trying to connect to the Docker daemon socket`

**Solutions**:

1. Add your user to the docker group:
   ```bash
   sudo usermod -aG docker $USER
   # Log out and back in for changes to take effect
   ```

2. Ensure the Docker socket has correct permissions:
   ```bash
   sudo chmod 666 /var/run/docker.sock  # Not recommended for production
   ```

3. Run with sudo (not recommended):
   ```bash
   sudo docker run --rm -v /var/run/docker.sock:/var/run/docker.sock docker-network-viz:latest
   ```

**Note**: In some environments (like Codespaces or Docker-in-Docker), you may need to use the host's Docker socket or run the binary directly instead of in a container.

### Image Not Found

**Error**: `Unable to find image 'docker-network-viz:latest' locally`

**Solution**: Build the image first:

```bash
./docker-build.sh
```

### Cannot Connect to Docker Daemon

**Error**: `Cannot connect to the Docker daemon at unix:///var/run/docker.sock`

**Solution**: Ensure Docker is running and the socket is mounted correctly.

## Development

### Building for Multiple Architectures

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t docker-network-viz:latest .
```

### Testing the Image

```bash
make docker-test
```

This will build the image and run `--help` to verify it works.

## CI/CD Usage

Example GitHub Actions workflow:

```yaml
- name: Build Docker Image
  run: ./docker-build.sh

- name: Test Docker Image
  run: |
    docker run --rm \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker-network-viz:latest --no-color > output.txt

- name: Upload Output
  uses: actions/upload-artifact@v3
  with:
    name: network-topology
    path: output.txt
```

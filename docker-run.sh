#!/bin/bash
# Run script for docker-network-viz Docker container

set -e

# Configuration
IMAGE_NAME="${IMAGE_NAME:-docker-network-viz}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
FULL_IMAGE_NAME="${IMAGE_NAME}:${IMAGE_TAG}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if image exists
if ! docker images "${IMAGE_NAME}" | grep -q "${IMAGE_TAG}"; then
    echo -e "${YELLOW}Image ${FULL_IMAGE_NAME} not found. Building...${NC}"
    ./docker-build.sh
fi

echo -e "${GREEN}Running ${FULL_IMAGE_NAME}${NC}"

# Parse command line arguments
# Everything after -- is passed to docker-network-viz
DOCKER_ARGS=""
APP_ARGS=""
PARSE_APP_ARGS=false

for arg in "$@"; do
    if [ "$arg" = "--" ]; then
        PARSE_APP_ARGS=true
        continue
    fi

    if [ "$PARSE_APP_ARGS" = true ]; then
        APP_ARGS="$APP_ARGS $arg"
    else
        DOCKER_ARGS="$DOCKER_ARGS $arg"
    fi
done

# Run the container
# Mount Docker socket to allow the container to inspect Docker
docker run --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    ${DOCKER_ARGS} \
    "${FULL_IMAGE_NAME}" \
    ${APP_ARGS}

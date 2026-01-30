#!/bin/bash
# Build script for docker-network-viz Docker image

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

echo -e "${GREEN}Building Docker image: ${FULL_IMAGE_NAME}${NC}"

# Build the image
docker build -t "${FULL_IMAGE_NAME}" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Successfully built ${FULL_IMAGE_NAME}${NC}"
    echo ""
    echo "Image details:"
    docker images "${IMAGE_NAME}" | grep "${IMAGE_TAG}"
    echo ""
    echo -e "${YELLOW}To run the container:${NC}"
    echo "  ./docker-run.sh"
    echo ""
    echo -e "${YELLOW}Or manually:${NC}"
    echo "  docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ${FULL_IMAGE_NAME}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

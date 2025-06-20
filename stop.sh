#!/bin/bash

# Stop Podman containers for CompareFlow

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Stopping CompareFlow Database...${NC}"

# Stop and remove PostgreSQL container
if podman ps -a | grep -q compareflow-postgres; then
    echo -e "${YELLOW}Stopping PostgreSQL container...${NC}"
    podman stop compareflow-postgres 2>/dev/null || true
    podman rm compareflow-postgres 2>/dev/null || true
    echo -e "${GREEN}PostgreSQL stopped and removed${NC}"
else
    echo -e "${YELLOW}PostgreSQL container not found${NC}"
fi

echo -e "${GREEN}Done!${NC}"
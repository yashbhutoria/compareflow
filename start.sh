#!/bin/bash

# Podman startup script for CompareFlow

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting CompareFlow Database with Podman${NC}"

# Check if podman is available
if ! command -v podman &> /dev/null; then
    echo -e "${RED}Error: Podman is not installed${NC}"
    echo "Please install Podman: https://podman.io/getting-started/installation"
    exit 1
fi

# Stop existing container if running
echo -e "${YELLOW}Cleaning up existing containers...${NC}"
podman stop compareflow-postgres 2>/dev/null || true
podman rm compareflow-postgres 2>/dev/null || true

# Start PostgreSQL with Podman
echo -e "${GREEN}Starting PostgreSQL container...${NC}"
podman run -d \
    --name compareflow-postgres \
    -e POSTGRES_USER=compareflow \
    -e POSTGRES_PASSWORD=compareflow123 \
    -e POSTGRES_DB=compareflow \
    -p 5432:5432 \
    postgres:15

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
for i in {1..30}; do
    if podman exec compareflow-postgres pg_isready -U compareflow &> /dev/null; then
        echo -e "${GREEN}PostgreSQL is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}PostgreSQL failed to start after 30 seconds${NC}"
        echo "Check container logs:"
        echo "  podman logs compareflow-postgres"
        exit 1
    fi
    echo -n "."
    sleep 1
done
echo

# Create admin user (optional)
echo -e "${YELLOW}Creating admin user...${NC}"
podman exec compareflow-postgres psql -U compareflow -d compareflow -c "
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
" 2>/dev/null || true

# Display connection info
echo -e "${GREEN}Database is ready!${NC}"
echo -e "${YELLOW}Connection details:${NC}"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: compareflow"
echo "  Username: compareflow"
echo "  Password: compareflow123"
echo ""
echo -e "${YELLOW}Connection string:${NC}"
echo "  postgresql://compareflow:compareflow123@localhost:5432/compareflow?sslmode=disable"
echo ""
echo -e "${GREEN}Starting CompareFlow application...${NC}"
echo "Default login: admin / admin123"
echo ""
go run cmd/compareflow/main.go

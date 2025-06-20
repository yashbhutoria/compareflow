#!/bin/bash

# Setup sample data for CompareFlow

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up sample data for CompareFlow${NC}"

# Check if PostgreSQL is running
if ! podman ps | grep -q compareflow-postgres; then
    echo -e "${RED}Error: PostgreSQL is not running${NC}"
    echo "Please run ./start.sh first"
    exit 1
fi

# Wait a moment for database to be ready
sleep 2

# Create sample tables
echo -e "${YELLOW}Creating sample tables...${NC}"
podman exec -i compareflow-postgres psql -U compareflow -d compareflow < scripts/create_sample_tables.sql

# Add sample connections and validations
echo -e "${YELLOW}Adding sample connections and validations...${NC}"
podman exec -i compareflow-postgres psql -U compareflow -d compareflow < scripts/add_sample_data.sql

# Add more validations
echo -e "${YELLOW}Adding additional validations...${NC}"
podman exec -i compareflow-postgres psql -U compareflow -d compareflow < scripts/add_more_validations.sql

echo -e "${GREEN}Sample data setup complete!${NC}"
echo ""
echo -e "${YELLOW}What was created:${NC}"
echo "  - 2 PostgreSQL connections (both pointing to the same database for demo)"
echo "  - Sample tables: products, orders, customers (and their _dw counterparts)"
echo "  - 8 validations to test different scenarios"
echo "  - One validation will fail (Orders) to show error handling"
echo ""
echo -e "${GREEN}You can now:${NC}"
echo "  1. Login to the web UI at http://localhost:8080"
echo "  2. Use credentials: admin / admin123"
echo "  3. Go to Connections to see the sample connections"
echo "  4. Go to Validations to see and run the sample validations"
echo "  5. The 'Orders Table Row Count' validation will fail (intentionally)"
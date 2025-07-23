#!/bin/bash

echo "=== Wallet Tracker Improvements Test ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Build the project
echo "1. Building project..."
if go build -o wallet-tracker cmd/wallet-tracker/main.go; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Check if services are running
echo -e "\n2. Checking services..."
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓ Services are running${NC}"
else
    echo -e "${RED}✗ Services not running. Starting them...${NC}"
    docker-compose up -d
    sleep 10
fi

# Test basic functionality
echo -e "\n3. Testing basic wallet tracking..."
export WALLET_TRACKER_APP_LOG_LEVEL=info
export WALLET_TRACKER_APP_LOG_FORMAT=text

if timeout 30 ./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa 2>&1 | grep -q "Error\|error"; then
    echo -e "${RED}✗ Errors detected during tracking${NC}"
else
    echo -e "${GREEN}✓ Basic tracking works${NC}"
fi

# Test debug logging
echo -e "\n4. Testing debug logging..."
export WALLET_TRACKER_APP_LOG_LEVEL=debug
if ./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa 2>&1 | head -20 | grep -q "DEBUG\|debug"; then
    echo -e "${GREEN}✓ Debug logging works${NC}"
else
    echo -e "${RED}✗ Debug logging not working${NC}"
fi

# Test Redis connection
echo -e "\n5. Testing Redis cache..."
if docker exec wallet-tracker-redis redis-cli ping | grep -q "PONG"; then
    echo -e "${GREEN}✓ Redis is accessible${NC}"
else
    echo -e "${RED}✗ Redis not accessible${NC}"
fi

# Test Neo4j connection
echo -e "\n6. Testing Neo4j connection..."
if curl -s http://localhost:7474 | grep -q "neo4j"; then
    echo -e "${GREEN}✓ Neo4j is accessible${NC}"
else
    echo -e "${RED}✗ Neo4j not accessible${NC}"
fi

echo -e "\n=== Test Summary ==="
echo "Check the output above for any red ✗ marks."
echo "If everything is green ✓, the improvements are working!"
echo
echo "Try these commands manually:"
echo "  ./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
echo "  ./wallet-tracker tracker websocket --all"
echo "  ./wallet-tracker redis get --exchanges binance --limit 3"

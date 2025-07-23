#!/bin/bash

echo "🚀 Building Wallet Tracker V2 (Simplified)..."

# Clean up problematic directories
echo "Cleaning up..."
rm -rf data/redis/appendonlydir 2>/dev/null || true
find . -name "*.go" -type f -exec chmod 644 {} \;

# Update dependencies
echo "Updating dependencies..."
go get github.com/jedib0t/go-pretty/v6/table
go get github.com/jedib0t/go-pretty/v6/text
go get github.com/spf13/cobra
go mod tidy -e

# Build with verbose output
echo "Building..."
go build -v -o wallet-tracker cmd/wallet-tracker/main.go 2>&1

if [ -f "./wallet-tracker" ]; then
    echo "✅ Build successful!"
    echo ""
    echo "Available commands:"
    echo ""
    ./wallet-tracker --help
    echo ""
    echo "Try the enhanced tracker:"
    echo "./wallet-tracker tracker trackv2 --help"
else
    echo "❌ Build failed. Trying alternative approach..."
    
    # Try building without the problematic imports
    echo "Building core components only..."
    cd cli/command/tracker
    go build -o ../../../wallet-tracker-core track_v2.go multichain.go webui.go
    cd ../../..
    
    if [ -f "./wallet-tracker-core" ]; then
        echo "✅ Core components built successfully!"
        echo "Run: ./wallet-tracker-core"
    fi
fi

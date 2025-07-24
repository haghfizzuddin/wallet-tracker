#!/bin/bash

echo "🔨 Building Real-time Monitor..."
echo ""

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Build the realtime monitor
go build -o realtime-monitor realtime_monitor.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📦 Binary created: realtime-monitor"
    echo ""
    echo "Usage:"
    echo "  ./realtime-monitor 0x0179eEd08227F3de8e7B2B50c91bb2E34DE5c659"
    echo ""
    echo "The monitor will:"
    echo "  - Check for new transactions every 30 seconds"
    echo "  - Alert on high-risk activities"
    echo "  - Update risk scores in real-time"
    echo "  - Track transaction patterns"
    echo ""
    echo "Press Ctrl+C to stop monitoring"
else
    echo "❌ Build failed!"
    echo ""
    echo "Common issues:"
    echo "  1. Check if all dependencies are installed"
    echo "  2. Run: go mod init && go mod tidy"
    echo "  3. Ensure config file exists"
    exit 1
fi

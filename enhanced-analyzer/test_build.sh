#!/bin/bash

echo "🔧 Testing build after fixes..."
echo ""

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Build the behavioral analyzer
echo "Building behavioral analyzer..."
go build -o behavioral-analyzer advanced_behavioral_analyzer.go

if [ $? -eq 0 ]; then
    echo "✅ Behavioral analyzer built successfully!"
else
    echo "❌ Behavioral analyzer build failed!"
    exit 1
fi

# Build the realtime monitor
echo ""
echo "Building realtime monitor..."
go build -o realtime-monitor realtime_monitor.go

if [ $? -eq 0 ]; then
    echo "✅ Realtime monitor built successfully!"
else
    echo "❌ Realtime monitor build failed!"
    exit 1
fi

echo ""
echo "🎉 All builds successful!"
echo ""
echo "Testing with an address..."
./behavioral-analyzer 0x0179eEd08227F3de8e7B2B50c91bb2E34DE5c659

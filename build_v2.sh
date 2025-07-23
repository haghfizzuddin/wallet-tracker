#!/bin/bash

echo "🚀 Building Wallet Tracker with V2 API Support"
echo "=============================================="
echo ""
echo "✨ NEW: One API key now works for 50+ chains!"
echo ""

# Build the new version
go build -o tracker tracker_v2_api.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 What's New in V2:"
    echo "  • Single API key for all EVM chains"
    echo "  • Support for 50+ blockchains"
    echo "  • Unified endpoint for all chains"
    echo ""
    echo "📊 Supported Networks:"
    echo "  • Ethereum (ETH)"
    echo "  • Binance Smart Chain (BSC)"
    echo "  • Polygon (MATIC)"
    echo "  • Arbitrum (ARB)"
    echo "  • Optimism (OP)"
    echo "  • Base (BASE)"
    echo "  • Avalanche (AVAX)"
    echo "  • Fantom (FTM)"
    echo "  • Blast (BLAST)"
    echo "  • Scroll (SCROLL)"
    echo "  • Bitcoin (BTC) - No API key needed"
    echo ""
    echo "🔑 Setup:"
    echo "  ./tracker config"
    echo ""
    echo "📝 Examples:"
    echo "  ./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
    echo "  ./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network ARB"
    echo "  ./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network BASE"
    echo ""
    chmod +x tracker
else
    echo "❌ Build failed!"
fi

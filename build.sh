#!/bin/bash

# Universal Wallet Tracker Build Script

echo "🚀 Building Universal Wallet Tracker..."

# Build the tracker
go build -o tracker universal_tracker.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📝 Quick Start Guide:"
    echo "─────────────────────────────────────────────────"
    echo ""
    echo "1️⃣  Add API Keys (optional but recommended):"
    echo "   Edit universal_tracker.go and add your keys at the top:"
    echo "   - EtherscanAPIKey: \"YOUR_KEY\""
    echo "   - BscscanAPIKey: \"YOUR_KEY\""
    echo ""
    echo "2️⃣  Basic Usage:"
    echo "   ./tracker <wallet_address>"
    echo ""
    echo "3️⃣  Examples:"
    echo "   ./tracker 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa         # Bitcoin"
    echo "   ./tracker 0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8  # Ethereum"
    echo "   ./tracker 0x123...abc --network BSC                   # Force BSC"
    echo "   ./tracker 0x123...abc --flow                          # Show flow"
    echo ""
    echo "4️⃣  Get Help:"
    echo "   ./tracker --help"
    echo "   ./tracker setup"
    echo ""
    chmod +x tracker
else
    echo "❌ Build failed!"
fi

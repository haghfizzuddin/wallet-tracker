#!/bin/bash

# Universal Wallet Tracker Build Script

echo "🚀 Building Universal Wallet Tracker (Secure Version)..."

# Build the tracker
go build -o tracker universal_tracker_secure.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📝 Quick Start Guide:"
    echo "─────────────────────────────────────────────────"
    echo ""
    echo "1️⃣  Configure API Keys (choose one method):"
    echo ""
    echo "   Method A - Environment Variables (Recommended):"
    echo "   export ETHERSCAN_API_KEY=your_key_here"
    echo "   export BSCSCAN_API_KEY=your_key_here"
    echo "   export POLYGON_API_KEY=your_key_here"
    echo ""
    echo "   Method B - Config File:"
    echo "   cp tracker-config.json.example tracker-config.json"
    echo "   # Edit tracker-config.json with your keys"
    echo ""
    echo "   Method C - Interactive Setup:"
    echo "   ./tracker config"
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
    echo "4️⃣  Get API Keys:"
    echo "   • Etherscan: https://etherscan.io/apis"
    echo "   • BscScan: https://bscscan.com/apis"
    echo "   • PolygonScan: https://polygonscan.com/apis"
    echo ""
    echo "🔐 Security: Your API keys are never stored in the code!"
    echo ""
    chmod +x tracker
else
    echo "❌ Build failed!"
fi

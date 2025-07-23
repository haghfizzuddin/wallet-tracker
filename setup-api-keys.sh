#!/bin/bash

echo "🔐 Wallet Tracker API Key Setup Helper"
echo "====================================="
echo ""
echo "This script will help you set up API keys for Ethereum and BSC tracking."
echo ""

# Check if config file exists
if [ -f "tracker-config.json" ]; then
    echo "⚠️  Config file already exists. Backing up to tracker-config.json.backup"
    cp tracker-config.json tracker-config.json.backup
fi

echo "📝 Please enter your API keys (or press Enter to skip):"
echo ""

# Read API keys
read -p "Etherscan API Key: " ETHERSCAN_KEY
read -p "BscScan API Key: " BSCSCAN_KEY
read -p "PolygonScan API Key: " POLYGON_KEY

# Create config file
cat > tracker-config.json << EOF
{
  "etherscan_api_key": "${ETHERSCAN_KEY:-}",
  "bscscan_api_key": "${BSCSCAN_KEY:-}",
  "polygon_api_key": "${POLYGON_KEY:-}"
}
EOF

echo ""
echo "✅ Config file created: tracker-config.json"
echo ""

# Show next steps
if [ -z "$ETHERSCAN_KEY" ] && [ -z "$BSCSCAN_KEY" ] && [ -z "$POLYGON_KEY" ]; then
    echo "⚠️  No API keys were entered. To get free API keys:"
    echo ""
    echo "1. Etherscan: https://etherscan.io/apis"
    echo "2. BscScan: https://bscscan.com/apis"
    echo "3. PolygonScan: https://polygonscan.com/apis"
    echo ""
    echo "Then run this script again or edit tracker-config.json directly."
else
    echo "🎉 You can now track wallets on:"
    [ -n "$ETHERSCAN_KEY" ] && echo "   ✓ Ethereum"
    [ -n "$BSCSCAN_KEY" ] && echo "   ✓ Binance Smart Chain"
    [ -n "$POLYGON_KEY" ] && echo "   ✓ Polygon"
    echo ""
    echo "Try: ./tracker 0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8"
fi

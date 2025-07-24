#!/bin/bash

echo "🔍 Testing Advanced Behavioral Analyzer with Real Addresses"
echo "=========================================================="

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Build first
echo "Building analyzer..."
go build -o behavioral-analyzer advanced_behavioral_analyzer.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "✅ Build successful! Starting tests..."
echo ""

# Test addresses
declare -A test_addresses=(
    # Known hack addresses
    ["0x098b716b8aaf21512996dc57eb0615e2383e2f96"]="Ronin Bridge Hacker"
    ["0xb66cd966670d962c227b3eaba30a872dbfb995db"]="Euler Finance Hacker"
    
    # Tornado Cash addresses (mixers)
    ["0x910cbd523d972eb0a6f4cae4618ad62622b39dbf"]="Tornado Cash 1 ETH"
    
    # Exchange addresses (should be low risk)
    ["0x28c6c06298d514db089934071355e5743bf21d60"]="Binance Hot Wallet"
    
    # Random address for testing
    ["0x9263e7871a6c9487ce985adfe1f65e66fab1ec81"]="Unknown Address"
)

# Test each address
for address in "${!test_addresses[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: ${test_addresses[$address]}"
    echo "Address: $address"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    ./behavioral-analyzer "$address"
    
    echo ""
    echo "Press Enter to continue to next address..."
    read
done

echo ""
echo "🏁 All tests completed!"

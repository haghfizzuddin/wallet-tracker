#!/bin/bash

echo "🔄 Testing the fixed analyzer with real addresses..."
echo ""

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Rebuild
echo "📦 Building analyzer..."
go build -o behavioral-analyzer advanced_behavioral_analyzer.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"
echo ""

# Test 1: Ronin Bridge Hacker (address that caused the error)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 1: Ronin Bridge Hacker"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./behavioral-analyzer 0x098b716b8aaf21512996dc57eb0615e2383e2f96

echo ""
echo "Press Enter to continue to next test..."
read

# Test 2: A random address
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 2: Random Address"  
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./behavioral-analyzer 0x9263e7871a6c9487ce985adfe1f65e66fab1ec81

echo ""
echo "Press Enter to continue to next test..."
read

# Test 3: Binance Hot Wallet (should be low risk)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 3: Binance Hot Wallet (Exchange)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./behavioral-analyzer 0x28c6c06298d514db089934071355e5743bf21d60

echo ""
echo "🏁 All tests completed!"

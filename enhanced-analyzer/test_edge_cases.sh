#!/bin/bash

echo "Testing edge cases..."

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Build first
go build -o behavioral-analyzer advanced_behavioral_analyzer.go 2>/dev/null

# Test with a brand new address (likely no transactions)
echo "Test: Brand new address with no transactions"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./behavioral-analyzer 0x0000000000000000000000000000000000000001

echo ""
echo "Test: Invalid address format"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./behavioral-analyzer invalid-address 2>&1 || echo "Handled invalid address gracefully"

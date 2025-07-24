#!/bin/bash

echo "Building Advanced Behavioral Analyzer..."

cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer

# Build the analyzer
go build -o behavioral-analyzer advanced_behavioral_analyzer.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "📦 Binary created: behavioral-analyzer"
    echo ""
    echo "Usage examples:"
    echo "  ./behavioral-analyzer 0x9263e7871a6c9487ce985adfe1f65e66fab1ec81"
    echo "  ./behavioral-analyzer 0x098b716b8aaf21512996dc57eb0615e2383e2f96"
    echo ""
    echo "Make sure you have:"
    echo "  1. enhanced-analyzer-config.json with your Etherscan API key"
    echo "  2. known_addresses.json (optional, will work without it)"
else
    echo "❌ Build failed!"
    echo ""
    echo "Common issues:"
    echo "  1. Missing Go dependencies - run: go mod init && go mod tidy"
    echo "  2. Missing config file - create enhanced-analyzer-config.json"
    echo "  3. Syntax errors - check the error messages above"
    exit 1
fi

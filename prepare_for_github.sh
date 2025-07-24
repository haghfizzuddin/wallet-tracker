#!/bin/bash

echo "🧹 Preparing wallet-tracker for GitHub..."
echo ""

cd /home/haghfizzuddin/repo/wallet-tracker

# Make scripts executable
chmod +x enhanced-analyzer/build_behavioral.sh
chmod +x enhanced-analyzer/test_behavioral.sh
chmod +x enhanced-analyzer/test_fixed_analyzer.sh
chmod +x enhanced-analyzer/test_edge_cases.sh

# Clean up binaries and temp files
echo "Cleaning binaries and temporary files..."
find . -type f -name "analyzer" -delete 2>/dev/null
find . -type f -name "analyzer-*" -delete 2>/dev/null
find . -type f -name "*-analyzer" -delete 2>/dev/null
find . -type f -name "behavioral-analyzer" -delete 2>/dev/null
find . -type f -name "tracker" -delete 2>/dev/null
find . -type f -name "security-analyzer*" -delete 2>/dev/null
find . -type f -name "*.log" -delete 2>/dev/null
find . -type f -name "*.tmp" -delete 2>/dev/null
find . -type f -name "hack_analysis_*.json" -delete 2>/dev/null

# Remove old/duplicate Go files
echo "Removing duplicate analyzer files..."
rm -f enhanced-analyzer/real_analyzer*.go
rm -f enhanced-analyzer/analyzer_v2.go
rm -f enhanced-analyzer/simple_analyzer.go
rm -f enhanced-analyzer/mock_analyzer*.go
rm -f enhanced-analyzer/advanced_risk_analyzer.go
rm -f enhanced-analyzer/calibrated_analyzer.go
rm -f enhanced-analyzer/dynamic_analyzer.go

# Remove root level duplicate files
echo "Removing root level duplicates..."
rm -f multi_chain_tracker.go
rm -f realtime_tracker*.go
rm -f standalone_v2.go
rm -f tracker_v2_api.go
rm -f universal_tracker*.go
rm -f btc_transactions_fix.go
rm -f analyze
rm -f example.txt
rm -f setup-api-keys.sh
rm -f define_schema.sh
rm -f integrated_tracker.sh

# Remove old build scripts
echo "Removing old build scripts..."
rm -f build_enhanced_analyzer*.sh
rm -f build_secure.sh
rm -f build_simple.sh
rm -f build_v2.sh
rm -f build.sh

# Remove old test scripts
echo "Removing old test scripts..."
rm -f test_improvements.sh
rm -f test_mock_analyzer.sh
rm -f test_real_hacks.sh
rm -f test_simple_mock.sh
rm -f enhanced-analyzer/test_fix.sh

# Remove old plan files (keeping updated ones)
rm -f ENHANCED_ANALYZER_SUMMARY.md
rm -f PHASE1_IMPROVEMENTS.md
rm -f PHASE2_PLAN.md

# Final cleanup in enhanced-analyzer
cd enhanced-analyzer
rm -f README.md
mv README_BEHAVIORAL.md README.md 2>/dev/null || true
rm -f etherscan_tagged_addresses.json  # Duplicate of known_addresses.json

echo ""
echo "✅ Repository cleaned and ready for GitHub!"
echo ""
echo "📁 Final structure:"
echo "   Core files:"
echo "   - enhanced-analyzer/advanced_behavioral_analyzer.go"
echo "   - enhanced-analyzer/realtime_monitor.go"
echo "   - enhanced-analyzer/gas_pattern_analyzer.go"
echo "   - enhanced-analyzer/final_analyzer.go (legacy reference)"
echo ""
echo "   Configuration:"
echo "   - enhanced-analyzer/enhanced-analyzer-config.json.example"
echo "   - enhanced-analyzer/known_addresses.json"
echo ""
echo "   Documentation:"
echo "   - README.md (main)"
echo "   - enhanced-analyzer/README.md (detailed)"
echo "   - enhanced-analyzer/GAS_PATTERNS_EXPLAINED.md"
echo "   - PHASE3_PLAN.md"
echo "   - PHASE4_PLAN.md"
echo "   - CONTRIBUTING.md"
echo ""
echo "   Scripts:"
echo "   - enhanced-analyzer/build_behavioral.sh"
echo "   - enhanced-analyzer/test_behavioral.sh"
echo "   - enhanced-analyzer/test_fixed_analyzer.sh"
echo ""
echo "📝 Next steps:"
echo "   1. Review .gitignore"
echo "   2. Add enhanced-analyzer-config.json to .gitignore if not already"
echo "   3. git add ."
echo "   4. git commit -m 'Major refactor: Advanced behavioral analyzer with real-time monitoring'"
echo "   5. git push origin master"

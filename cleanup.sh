#!/bin/bash

echo "🧹 Cleaning up wallet-tracker repository..."
echo ""

cd /home/haghfizzuddin/repo/wallet-tracker

# Remove binary files
echo "Removing binary files..."
rm -f analyzer behavioral-analyzer tracker security-analyzer security-analyzer-real
rm -f enhanced-analyzer/analyzer* enhanced-analyzer/behavioral-analyzer
rm -f enhanced-analyzer/calibrated-analyzer enhanced-analyzer/dynamic-analyzer
rm -f enhanced-analyzer/enhanced-analyzer* enhanced-analyzer/final-analyzer
rm -f enhanced-analyzer/mock_analyzer enhanced-analyzer/mock_analyzer_simple
rm -f enhanced-analyzer/risk-analyzer

# Remove temporary and test output files
echo "Removing temporary files..."
rm -f enhanced-analyzer/hack_analysis_*.json

# Remove old/duplicate files
echo "Removing duplicate analyzer files..."
rm -f enhanced-analyzer/real_analyzer.go
rm -f enhanced-analyzer/real_analyzer_fixed.go
rm -f enhanced-analyzer/real_analyzer_v2.go
rm -f enhanced-analyzer/analyzer_v2.go
rm -f enhanced-analyzer/simple_analyzer.go
rm -f enhanced-analyzer/mock_analyzer.go
rm -f enhanced-analyzer/mock_analyzer_simple.go
rm -f enhanced-analyzer/advanced_risk_analyzer.go
rm -f enhanced-analyzer/calibrated_analyzer.go
rm -f enhanced-analyzer/dynamic_analyzer.go

# Remove duplicate build scripts
echo "Removing old build scripts..."
rm -f build_enhanced_analyzer.sh
rm -f build_enhanced_analyzer_v2.sh
rm -f build_secure.sh
rm -f build_simple.sh
rm -f build_v2.sh

# Remove old test scripts
echo "Removing old test scripts..."
rm -f test_improvements.sh
rm -f test_mock_analyzer.sh
rm -f test_real_hacks.sh
rm -f test_simple_mock.sh
rm -f enhanced-analyzer/test_fix.sh

# Remove duplicate tracker files
echo "Removing duplicate tracker files..."
rm -f multi_chain_tracker.go
rm -f realtime_tracker.go
rm -f realtime_tracker_fixed.go
rm -f standalone_v2.go
rm -f tracker_v2_api.go
rm -f universal_tracker.go
rm -f universal_tracker_secure.go
rm -f btc_transactions_fix.go

# Remove misc files
echo "Removing misc files..."
rm -f analyze
rm -f example.txt
rm -f integrated_tracker.sh
rm -f define_schema.sh
rm -f setup-api-keys.sh

# Keep only the final, working analyzer
echo "Organizing enhanced-analyzer directory..."
cd enhanced-analyzer
rm -f README.md  # We'll use README_BEHAVIORAL.md as the main one
mv README_BEHAVIORAL.md README.md 2>/dev/null || true

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "Remaining structure:"
echo "- Main behavioral analyzer: enhanced-analyzer/advanced_behavioral_analyzer.go"
echo "- Real-time monitor: enhanced-analyzer/realtime_monitor.go"
echo "- Gas pattern analyzer: enhanced-analyzer/gas_pattern_analyzer.go"
echo "- Final analyzer (legacy): enhanced-analyzer/final_analyzer.go"
echo "- Documentation: enhanced-analyzer/README.md, GAS_PATTERNS_EXPLAINED.md"
echo "- Build scripts: enhanced-analyzer/build_behavioral.sh"
echo "- Test scripts: enhanced-analyzer/test_behavioral.sh, test_fixed_analyzer.sh"

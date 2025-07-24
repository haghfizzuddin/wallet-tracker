# Advanced Behavioral Blockchain Analyzer

A sophisticated blockchain security analyzer that combines behavioral pattern analysis, statistical methods, and real-time data integration to detect suspicious activities without relying solely on hardcoded addresses.

## Features

### Phase 1 + 2 Combined Implementation

1. **Behavioral Pattern Analysis**
   - Transaction velocity monitoring
   - Value concentration patterns
   - Time-based anomaly detection
   - Gas price manipulation detection
   - Circular transaction patterns
   - New address behavior analysis

2. **Statistical Methods**
   - Benford's Law analysis for transaction values
   - Entropy calculations for address randomness
   - Clustering coefficient for network analysis
   - Temporal anomaly detection using standard deviation
   - Z-score analysis for outlier detection

3. **Real-time Data Integration**
   - Live Etherscan API integration
   - Transaction history analysis
   - Real-time label checking
   - Failed transaction detection
   - Contract creation monitoring

4. **Heuristic Rules**
   - Newly created address detection
   - High-value transaction monitoring
   - Mixer interaction detection
   - Suspicious method call identification
   - Automated behavior detection

## Installation

1. Navigate to the enhanced-analyzer directory:
```bash
cd /home/haghfizzuddin/repo/wallet-tracker/enhanced-analyzer
```

2. Build the analyzers:
```bash
# Build the behavioral analyzer
chmod +x build_behavioral.sh
./build_behavioral.sh

# Build the real-time monitor
go build -o realtime-monitor realtime_monitor.go
```

## Usage

### 1. One-time Address Analysis

Analyze any Ethereum address for suspicious patterns:

```bash
./behavioral-analyzer 0x9263e7871a6c9487ce985adfe1f65e66fab1ec81
```

Example output:
```
🔍 Analyzing address: 0x9263e7871a6c9487ce985adfe1f65e66fab1ec81
================================================================================
📊 Fetching transaction history...
   Found 125 transactions
🏷️  Checking real-time labels...
🧠 Analyzing behavioral patterns...
📈 Performing statistical analysis...
⚡ Checking real-time risk indicators...

📊 ANALYSIS RESULTS
================================================================================
🎯 Risk Score: 0.75/1.00 (Confidence: 85.0%)

📈 Statistical Analysis:
   • Benford's Law Score: 0.68
   • Velocity Score: 0.82
   • Entropy Score: 0.45
   • Clustering Score: 0.71
   • Temporal Anomaly: 0.55

🚩 Behavioral Patterns Detected:
   • Rapid outgoing transfers detected: 45.2 ETH in recent transactions [Severity: 0.90]
   • Detected 25 transactions in one hour (threshold: 20) [Severity: 0.85]
   • Called suspicious method: flashLoan() - Flash Loan [Severity: 0.80]
```

### 2. Real-time Monitoring

Monitor an address continuously for new suspicious activities:

```bash
./realtime-monitor 0x9263e7871a6c9487ce985adfe1f65e66fab1ec81
```

The monitor will:
- Check for new transactions every 30 seconds
- Alert on high-risk activities
- Update risk scores in real-time
- Track alert history

### 3. Batch Testing

Test multiple addresses at once:

```bash
chmod +x test_behavioral.sh
./test_behavioral.sh
```

## Configuration

### API Keys

Create `enhanced-analyzer-config.json`:
```json
{
  "etherscan_api_key": "YOUR_ETHERSCAN_API_KEY",
  "infura_url": "YOUR_INFURA_URL"  // Optional
}
```

### Risk Thresholds

Default thresholds (can be modified in code):
- High Value: 10 ETH
- Velocity: 20 transactions/hour
- Gas Anomaly: 3x average
- New Address Age: 60 minutes
- Benford Deviation: 15%

### Known Addresses (Optional)

The analyzer works without `known_addresses.json`, but you can enhance detection by adding known addresses:

```json
{
  "exchanges": {
    "0x...": "Exchange Name"
  },
  "mixers": {
    "0x...": "Mixer Name"
  },
  "hackers": {
    "0x...": {
      "name": "Hack Name",
      "amount_stolen": "Amount",
      "date": "YYYY-MM-DD",
      "hack_type": "Type"
    }
  }
}
```

## Risk Score Interpretation

- **0.0 - 0.2**: Minimal Risk ✅
- **0.2 - 0.4**: Low Risk ℹ️
- **0.4 - 0.6**: Medium Risk ⚠️
- **0.6 - 0.8**: High Risk ⚡
- **0.8 - 1.0**: Critical Risk 🚨

## Key Improvements Over Hardcoded System

1. **Dynamic Detection**: Identifies suspicious behavior patterns regardless of whether the address is in a database
2. **Statistical Analysis**: Uses mathematical models to detect anomalies
3. **Real-time Updates**: Fetches latest data from Etherscan
4. **Behavioral Patterns**: Recognizes attack patterns like rapid drainage, mixer sequences, flash loans
5. **Confidence Scoring**: Provides confidence levels based on data availability
6. **Comprehensive Analysis**: Combines multiple detection methods for accuracy

## Detection Methods

### Behavioral Patterns
- **High Velocity**: Many transactions in short time
- **Value Concentration**: Large fund movements
- **New Address Activity**: New addresses with high-value transactions
- **Gas Anomalies**: Abnormal gas prices (front-running indicators)
- **Circular Patterns**: Money laundering indicators
- **Time Patterns**: Bot/automated behavior detection

### Statistical Methods
- **Benford's Law**: Natural distribution of transaction values
- **Entropy Analysis**: Randomness in address interactions
- **Clustering Coefficient**: Network topology analysis
- **Temporal Anomalies**: Time gap irregularities
- **Variance Analysis**: Consistency in behavior

### Real-time Indicators
- Failed transaction ratios
- Contract creation patterns
- Known address interactions
- Suspicious method calls
- Exchange/mixer interactions

## Future Enhancements (Phase 3-4)

- Machine Learning models for pattern recognition
- Graph neural networks for transaction flow analysis
- Cross-chain analysis capabilities
- Integration with threat intelligence feeds
- Community reporting system
- Advanced visualization dashboard

## Troubleshooting

1. **No Etherscan API Key**: Add your key to `enhanced-analyzer-config.json`
2. **Build Errors**: Ensure Go is installed and dependencies are met
3. **No Transactions Found**: Check if the address exists and has transactions
4. **Rate Limits**: Etherscan has rate limits; use multiple API keys if needed

## Examples

### Analyze a Known Hacker Address
```bash
./behavioral-analyzer 0x098b716b8aaf21512996dc57eb0615e2383e2f96
```

### Monitor an Exchange Hot Wallet
```bash
./realtime-monitor 0x28c6c06298d514db089934071355e5743bf21d60
```

### Check a Suspicious Address
```bash
./behavioral-analyzer 0xSUSPICIOUS_ADDRESS_HERE
```

## Contributing

To add new detection methods:
1. Add pattern detection in `analyzeBehavioralPatterns()`
2. Add statistical analysis in `performStatisticalAnalysis()`
3. Update risk calculation in `calculateFinalRiskScore()`
4. Add new suspicious patterns to `suspiciousPatterns` map

## License

MIT License - See LICENSE file for details

# Wallet Tracker - Advanced Blockchain Security Analyzer

A sophisticated blockchain security analysis tool that uses behavioral patterns, statistical analysis, and real-time data to detect suspicious activities on Ethereum and other EVM-compatible chains.

## 🚀 Features

### Core Capabilities
- **Behavioral Pattern Analysis**: Detects suspicious patterns without relying on hardcoded addresses
- **Statistical Analysis**: Uses Benford's Law, entropy calculations, and clustering algorithms
- **Real-time Monitoring**: Continuous monitoring with alert system
- **Gas Anomaly Detection**: Identifies front-running, MEV attacks, and exploit patterns
- **Multi-source Integration**: Etherscan API, real-time labels, and transaction history analysis

### Detection Methods
1. **Behavioral Patterns**
   - Transaction velocity monitoring
   - Value concentration analysis
   - Circular transaction detection
   - New address behavior tracking
   - Time-based anomaly detection

2. **Statistical Methods**
   - Benford's Law for natural distribution
   - Entropy analysis for randomness
   - Clustering coefficient calculation
   - Temporal anomaly detection
   - Z-score analysis for outliers

3. **Gas Pattern Analysis**
   - Front-running detection
   - Sandwich attack identification
   - Gas war pattern recognition
   - Exploit execution patterns
   - Censorship evasion detection

## 📋 Prerequisites

- Go 1.19 or higher
- Etherscan API key (free from [Etherscan.io](https://etherscan.io/apis))
- Git

## 🛠️ Installation & Setup

### 1. Clone the repository
```bash
git clone https://github.com/haghfizzuddin/wallet-tracker.git
cd wallet-tracker
```

### 2. Set up Etherscan API key
```bash
cd enhanced-analyzer

# Create config file from example
cp enhanced-analyzer-config.json.example enhanced-analyzer-config.json

# Edit the file and add your Etherscan API key
nano enhanced-analyzer-config.json
# or use any text editor to add your key:
# {
#   "etherscan_api_key": "YOUR_ETHERSCAN_API_KEY_HERE"
# }
```

### 3. Build the analyzers

#### Option A: Using build scripts (Recommended)
```bash
# Make build scripts executable
chmod +x build_behavioral.sh build_realtime_monitor.sh

# Build the behavioral analyzer
./build_behavioral.sh

# Build the real-time monitor
./build_realtime_monitor.sh
```

#### Option B: Manual build
```bash
# Build behavioral analyzer
go build -o behavioral-analyzer advanced_behavioral_analyzer.go

# Build real-time monitor
go build -o realtime-monitor realtime_monitor.go
```

## 🔍 Usage

### Basic Address Analysis
Analyze any Ethereum address for suspicious patterns:

```bash
# Analyze a specific address
./behavioral-analyzer 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49

# Example outputs:
# - Known exchange address (low risk)
./behavioral-analyzer 0x28c6c06298d514db089934071355e5743bf21d60

# - MEV bot address (medium-high risk)
./behavioral-analyzer 0x633dCF31bb890b26279C9a0480754DC09E27c01E

# - Suspicious address (high risk)
./behavioral-analyzer 0x098b716b8aaf21512996dc57eb0615e2383e2f96
```

### Real-time Monitoring
Monitor an address continuously for new transactions:

```bash
# Start monitoring (checks every 30 seconds)
./realtime-monitor 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49

# Press Ctrl+C to stop and see summary
```

### Batch Testing
Test multiple addresses using the test script:

```bash
chmod +x test_behavioral.sh
./test_behavioral.sh
```

## 📊 Understanding the Output

### Risk Score Levels
- **0.0 - 0.2**: Minimal Risk ✅ - Safe to interact
- **0.2 - 0.4**: Low Risk ℹ️ - Generally safe with standard precautions
- **0.4 - 0.6**: Medium Risk ⚠️ - Exercise caution, investigate further
- **0.6 - 0.8**: High Risk ⚡ - Avoid interaction unless necessary
- **0.8 - 1.0**: Critical Risk 🚨 - Do not interact, likely malicious

### Example Output Explained
```
🔍 Analyzing address: 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49
================================================================================
📊 Fetching transaction history...
   Found 523 transactions            # Total transactions analyzed
🧠 Analyzing behavioral patterns...  # Pattern detection in progress
📈 Performing statistical analysis... # Mathematical analysis

📊 ANALYSIS RESULTS
================================================================================
🎯 Risk Score: 0.72/1.00 (Confidence: 89.2%)
   │                    │                │
   │                    │                └── How confident the analysis is
   │                    └──────────────────── Maximum risk is 1.00
   └──────────────────────────────────────── Current risk level

📈 Statistical Analysis:
   • Benford's Law Score: 0.65  # Natural distribution check (higher = suspicious)
   • Velocity Score: 0.78       # Transaction speed (higher = faster/suspicious)
   • Entropy Score: 0.42        # Randomness measure (lower = more predictable)
   • Clustering Score: 0.81     # Network connectivity (higher = more connected)
   • Temporal Anomaly: 0.58     # Time pattern irregularities

🚩 Behavioral Patterns Detected:
   • Rapid outgoing transfers: 125.3 ETH [Severity: 0.90]
     └── Large amount leaving quickly (possible hack/drain)
   
   • High transaction velocity: 45 tx/hour [Severity: 0.85]
     └── Too many transactions per hour (bot/automated activity)
   
   • Potential front-running: 0xabc...def used 2500 Gwei [Severity: 0.85]
     └── Extremely high gas price to get priority (MEV activity)
```

## 🏗️ Project Structure

```
wallet-tracker/
├── enhanced-analyzer/                    # Main analyzer directory
│   ├── advanced_behavioral_analyzer.go   # Core behavioral analysis
│   ├── realtime_monitor.go              # Real-time monitoring
│   ├── gas_pattern_analyzer.go          # Gas anomaly detection
│   ├── known_addresses.json             # Optional known addresses
│   ├── enhanced-analyzer-config.json    # Your API configuration
│   ├── build_behavioral.sh              # Build script for analyzer
│   ├── build_realtime_monitor.sh        # Build script for monitor
│   └── test_behavioral.sh               # Test script
├── README.md                            # This file
├── CONTRIBUTING.md                      # Contribution guidelines
├── LICENSE                              # MIT License
└── .gitignore                          # Git ignore rules
```

## 🔧 Configuration

### API Configuration (`enhanced-analyzer-config.json`)
```json
{
  "etherscan_api_key": "YOUR_ETHERSCAN_API_KEY",
  "infura_url": "https://mainnet.infura.io/v3/YOUR_PROJECT_ID" // Optional
}
```

### Risk Thresholds (in code)
You can adjust detection sensitivity by modifying thresholds in `advanced_behavioral_analyzer.go`:

```go
RiskThresholds{
    HighValueThreshold:    10.0,  // ETH - Transactions above this are flagged
    VelocityThreshold:     20,    // tx/hour - More than this is suspicious
    GasAnomalyMultiplier:  3.0,   // 3x average gas = anomaly
    NewAddressAgeMinutes:  60,    // New addresses younger than this are flagged
    BenfordDeviationLimit: 0.15,  // 15% deviation from Benford's Law
}
```

## 🐛 Troubleshooting

### Common Issues

1. **"No such file or directory" when running analyzer**
   ```bash
   # Make sure you built it first:
   ./build_behavioral.sh
   ```

2. **"Failed to load config" error**
   ```bash
   # Create config file with your API key:
   cp enhanced-analyzer-config.json.example enhanced-analyzer-config.json
   # Edit and add your Etherscan API key
   ```

3. **Build fails with Go errors**
   ```bash
   # Initialize Go modules:
   go mod init
   go mod tidy
   ```

4. **"No transaction history found"**
   - Check if the address is valid
   - Ensure your Etherscan API key is correct
   - The address might be new with no transactions

5. **Rate limit errors**
   - Free Etherscan API has rate limits
   - Wait a few seconds between requests
   - Consider getting a paid API key for heavy usage

## 📚 Understanding Detection Methods

### Why High Gas Prices Matter
High gas prices often indicate:
- **Front-running**: Paying high gas to execute before a victim
- **MEV attacks**: Sandwich attacks on DEX trades
- **Exploit execution**: Ensuring malicious transactions succeed
- **Competition**: Gas wars during NFT mints or token launches

### Behavioral Patterns We Detect
- **Rapid drainage**: Multiple large withdrawals quickly
- **Mixer sequences**: Interaction with privacy protocols
- **Circular transfers**: Money laundering patterns
- **New address activity**: Fresh addresses with large transactions
- **Bot patterns**: Automated trading or attack behavior

## 🚦 Next Steps

1. **For basic usage**: Run the analyzer on any address you want to check
2. **For investigations**: Use real-time monitor during active incidents
3. **For development**: Check CONTRIBUTING.md to add new detection methods
4. **For enterprise**: See PHASE3_PLAN.md and PHASE4_PLAN.md for roadmap

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on adding new features.

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

## ⚠️ Disclaimer

This tool provides risk analysis based on behavioral patterns and should not be the sole basis for security decisions. Always perform additional due diligence.

## 📞 Support

- **Issues**: Use GitHub Issues for bug reports
- **Discussions**: Use GitHub Discussions for questions
- **Security**: For security issues, please email directly

---

Built with ❤️ by the blockchain security community

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
- Etherscan API key
- Git

## 🛠️ Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/wallet-tracker.git
cd wallet-tracker
```

2. Set up configuration:
```bash
cd enhanced-analyzer
cp enhanced-analyzer-config.json.example enhanced-analyzer-config.json
# Edit the file and add your Etherscan API key
```

3. Build the analyzers:
```bash
chmod +x build_behavioral.sh
./build_behavioral.sh

# For real-time monitoring
go build -o realtime-monitor realtime_monitor.go
```

## 🔍 Usage

### Basic Address Analysis
Analyze any Ethereum address for suspicious patterns:

```bash
./behavioral-analyzer 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49
```

### Real-time Monitoring
Monitor an address continuously:

```bash
./realtime-monitor 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49
```

### Batch Testing
Test multiple addresses:

```bash
./test_behavioral.sh
```

## 📊 Output Example

```
🔍 Analyzing address: 0x742d35Cc6634C0532925a3b844Bc9e7595f8fA49
================================================================================
📊 Fetching transaction history...
   Found 523 transactions
🧠 Analyzing behavioral patterns...
📈 Performing statistical analysis...

📊 ANALYSIS RESULTS
================================================================================
🎯 Risk Score: 0.72/1.00 (Confidence: 89.2%)

📈 Statistical Analysis:
   • Benford's Law Score: 0.65
   • Velocity Score: 0.78
   • Entropy Score: 0.42
   • Clustering Score: 0.81
   • Temporal Anomaly: 0.58

🚩 Behavioral Patterns Detected:
   • Rapid outgoing transfers detected: 125.3 ETH [Severity: 0.90]
   • High transaction velocity: 45 tx/hour [Severity: 0.85]
   • Potential front-running: 0xabc...def used 2500 Gwei [Severity: 0.85]

💡 Recommendations:
   ⚡ HIGH RISK: Exercise extreme caution with this address.
   🔍 Perform additional due diligence before any interaction.
```

## 🏗️ Architecture

```
wallet-tracker/
├── enhanced-analyzer/          # Main analyzer implementation
│   ├── advanced_behavioral_analyzer.go  # Core behavioral analysis
│   ├── realtime_monitor.go             # Real-time monitoring
│   ├── gas_pattern_analyzer.go         # Gas anomaly detection
│   ├── final_analyzer.go               # Legacy analyzer
│   ├── known_addresses.json            # Known malicious addresses
│   └── README.md                       # Detailed documentation
├── cli/                        # CLI command structure
├── cmd/                        # Application entry points
├── domain/                     # Domain models
├── pkg/                        # Shared packages
└── docker-compose.yml          # Docker configuration
```

## 🔧 Configuration

### Risk Thresholds
Edit thresholds in `advanced_behavioral_analyzer.go`:

```go
RiskThresholds{
    HighValueThreshold:    10.0,  // ETH
    VelocityThreshold:     20,    // tx/hour
    GasAnomalyMultiplier:  3.0,   // 3x average
    NewAddressAgeMinutes:  60,    // minutes
    BenfordDeviationLimit: 0.15,  // 15% deviation
}
```

### API Configuration
Create `enhanced-analyzer-config.json`:

```json
{
  "etherscan_api_key": "YOUR_API_KEY_HERE",
  "infura_url": "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
}
```

## 📈 Risk Score Interpretation

- **0.0 - 0.2**: Minimal Risk ✅
- **0.2 - 0.4**: Low Risk ℹ️
- **0.4 - 0.6**: Medium Risk ⚠️
- **0.6 - 0.8**: High Risk ⚡
- **0.8 - 1.0**: Critical Risk 🚨

## 🚦 Roadmap

### ✅ Phase 1 & 2 (Completed)
- Behavioral pattern analysis
- Statistical methods
- Real-time data integration
- Gas anomaly detection

### 📋 Phase 3 (Planned)
- Machine learning models
- Graph neural networks
- Advanced feature engineering
- Automated retraining

### 📋 Phase 4 (Future)
- Cross-chain analysis
- External intelligence integration
- Smart contract vulnerability detection
- Enterprise API

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Ethereum community for blockchain data
- Etherscan for API access
- Go community for excellent libraries

## ⚠️ Disclaimer

This tool is for educational and research purposes only. Always perform your own due diligence before interacting with any blockchain address.

## 📞 Support

- Create an issue for bug reports
- Join our Discord for discussions
- Check the [documentation](enhanced-analyzer/README.md) for detailed usage

---

**Note**: This tool does not guarantee 100% accuracy in detecting malicious addresses. It provides risk indicators based on behavioral patterns and should be used as part of a comprehensive security approach.

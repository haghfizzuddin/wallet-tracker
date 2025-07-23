# 🚀 Universal Wallet Tracker

A powerful multi-chain blockchain wallet tracker supporting 50+ networks with a single API key!

![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)
![Chains Supported](https://img.shields.io/badge/chains-50%2B-green.svg)

## ✨ Features

- 🌐 **Multi-Chain Support**: Track wallets across Bitcoin, Ethereum, BSC, Polygon, Arbitrum, and 50+ other chains
- 🔑 **Single API Key**: Use Etherscan's V2 API - one key for all EVM chains!
- 📊 **Rich Transaction Details**: View amounts, fees, timestamps, and USD values
- 💸 **Fund Flow Visualization**: ASCII diagrams showing money flow
- 🎨 **Beautiful Terminal UI**: Colored tables and formatted output
- 🔒 **Secure**: API keys stored separately from code
- ⚡ **Real-Time Data**: Live blockchain data with current prices

## 📋 Supported Networks

### No API Key Required
- Bitcoin (BTC)

### With Etherscan API Key (One key for all!)
- Ethereum (ETH)
- Binance Smart Chain (BSC)
- Polygon (MATIC)
- Arbitrum (ARB)
- Optimism (OP)
- Base (BASE)
- Avalanche (AVAX)
- Fantom (FTM)
- Blast (BLAST)
- Scroll (SCROLL)
- And 40+ more chains!

## 🚀 Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/haghfizzuddin/wallet-tracker.git
cd wallet-tracker
chmod +x build_v2.sh
./build_v2.sh
```

### 2. Configure API Key
Get your free API key from [Etherscan](https://etherscan.io/apis) (works for all chains!)

```bash
# Method 1: Interactive setup
./tracker config

# Method 2: Environment variable
export ETHERSCAN_API_KEY=your_key_here

# Method 3: Config file
cp tracker-config.json.example tracker-config.json
# Edit the file with your key
```

### 3. Track Wallets!
```bash
# Bitcoin (no API key needed)
./tracker 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

# Ethereum
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Arbitrum
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network ARB

# Show fund flow diagram
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --flow
```

## 📸 Screenshots

### Transaction Tracking
```
╔══════════════════════════════════════════════════════════════════╗
║         UNIVERSAL BLOCKCHAIN TRACKER v5.0 (V2 API)              ║
╚══════════════════════════════════════════════════════════════════╝

📊 Tracking: 0xd8dA...6045
🌐 Network: Ethereum (Chain ID: 1)
🕒 Time: 2025-01-23 12:45:00
──────────────────────────────────────────────────────────────────────

 #  TYPE  TIME   FROM → TO         AMOUNT        USD VALUE  STATUS
 1  IN    12:30  0x123...→[TRACKED] 1.5000 ETH   $3,525.00  ✅
 2  OUT   11:45  [TRACKED]→0x456... 0.5000 ETH   $1,175.00  ✅
```

### Fund Flow Visualization
```
💸 Fund Flow Visualization
──────────────────────────────────────────────────────────────────────
📥 INFLOWS:
  0x123...def         1.500000 ETH

  [0xd8dA...6045]

📤 OUTFLOWS:
  0x456...789         0.500000 ETH
```

## 🛠️ Development

### Prerequisites
- Go 1.19 or higher
- Git

### Building from Source
```bash
# Standard build
go build -o tracker tracker_v2_api.go

# Build with all features
make build
```

### Project Structure
```
wallet-tracker/
├── tracker_v2_api.go          # Main V2 API implementation
├── universal_tracker_secure.go # Secure config management
├── build_v2.sh               # Build script
├── tracker-config.json.example # Config template
└── README.md                 # This file
```

## 🔧 Configuration

### API Keys
- Get free API key: https://etherscan.io/apis
- Works for ALL supported chains (V2 API)
- Multiple storage options (env, file, interactive)

### Config File Format
```json
{
  "etherscan_api_key": "YOUR_API_KEY_HERE"
}
```

### Environment Variables
```bash
export ETHERSCAN_API_KEY=your_key_here
```

## 📚 Usage Examples

### Basic Tracking
```bash
# Auto-detect network
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Specify network
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network BSC

# Limit transactions
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --limit 20

# Show fund flow
./tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --flow
```

### Advanced Usage
```bash
# Track on multiple chains
for chain in ETH BSC MATIC ARB OP; do
  echo "Checking $chain..."
  ./tracker 0xYourAddress --network $chain
done

# Export to file (redirect output)
./tracker 0xAddress > wallet_report.txt
```

## 🚦 Development Phases

### ✅ Phase 1: Core Functionality
- Basic wallet tracking
- Bitcoin support
- Terminal UI improvements

### ✅ Phase 2: Multi-Chain Support
- Etherscan V2 API integration
- 50+ chain support
- Unified API key management
- Enhanced UI with colors and tables

### 🚧 Phase 3: Advanced Features (Planned)
- Web dashboard
- REST API server
- Real-time monitoring
- Risk scoring
- Machine learning analytics
- Cross-chain portfolio tracking

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### How to Contribute
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🐛 Troubleshooting

### Common Issues

**"API key issue" error**
- Make sure you have configured your API key correctly
- Check that your key has not exceeded rate limits

**"No transactions found"**
- The wallet might be new or have no transactions
- Try a known active wallet for testing

**Network auto-detection issues**
- Use `--network` flag to specify the chain explicitly

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Etherscan](https://etherscan.io) for the amazing V2 API
- [CoinGecko](https://coingecko.com) for price data
- The Go community for excellent libraries
- All contributors and users

## 📞 Support

- 🐛 Create an [issue](https://github.com/haghfizzuddin/wallet-tracker/issues) for bugs
- 💬 Check [discussions](https://github.com/haghfizzuddin/wallet-tracker/discussions) for Q&A
- ⭐ Star the project if you find it useful!
- 🍴 Fork and contribute!

## 🔗 Links

- [Etherscan API Documentation](https://docs.etherscan.io)
- [Supported Chain IDs](https://chainlist.org)
- [Go Documentation](https://pkg.go.dev)

---

<p align="center">
Made with ❤️ by <a href="https://github.com/haghfizzuddin">haghfizzuddin</a>
</p>

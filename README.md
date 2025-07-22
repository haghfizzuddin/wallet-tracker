# Wallet Tracker CLI

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

## 🚀 Overview

Wallet Tracker CLI is a powerful command-line tool for tracking cryptocurrency wallet transactions across blockchain networks. It features real-time transaction monitoring, exchange detection, and visual analytics through Neo4j graph database integration.

### Key Features
- 🔍 **Real-time wallet tracking** - Monitor blockchain transactions as they happen
- 🏦 **Exchange detection** - Identify transactions to/from major exchanges
- 📊 **Visual analytics** - Neo4j integration with NeoDash for graph visualization
- 🔄 **WebSocket support** - Stream live transaction data
- 💾 **Redis caching** - Improved performance with intelligent caching
- 🛡️ **Robust error handling** - Retry mechanisms and graceful failure recovery

## 📋 Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Track Wallet](#track-wallet)
  - [WebSocket Streaming](#websocket-streaming)
  - [Exchange Detection](#exchange-detection)
  - [Visual Analytics](#visual-analytics)
- [Development](#development)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- Go 1.19 or higher
- Docker & Docker Compose (for Neo4j and Redis)
- Git

## Installation

### Option 1: Install from source

```bash
# Clone the repository
git clone https://github.com/haghfizzuddin/wallet-tracker.git
cd wallet-tracker

# Install dependencies
go mod download

# Build the binary
go build -o wallet-tracker cmd/wallet-tracker/main.go

# Make it executable
chmod +x wallet-tracker
```

### Option 2: Using go install

```bash
go install github.com/haghfizzuddin/wallet-tracker/cmd/wallet-tracker@latest
```

### Option 3: Download pre-built binary

Check the [releases page](https://github.com/haghfizzuddin/wallet-tracker/releases) for pre-built binaries for your platform.

## Configuration

### Quick Start with Docker Compose

```bash
# Start Neo4j and Redis
docker-compose up -d

# Copy environment example
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### Configuration File (config.yaml)

Create a `config.yaml` for advanced configuration:

```yaml
app:
  log_level: info      # debug, info, warn, error
  log_format: json     # json or text

database:
  uri: neo4j://localhost:7687
  username: neo4j
  password: your_password_here

api:
  rate_limit: 10       # requests per second
  max_retries: 3
  retry_delay: 1s

redis:
  host: localhost
  port: 6379
  ttl: 1h             # cache time-to-live
```

### Environment Variables

All configuration can be overridden with environment variables:

```bash
export WALLET_TRACKER_DATABASE_URI=neo4j://localhost:7687
export WALLET_TRACKER_DATABASE_USERNAME=neo4j
export WALLET_TRACKER_DATABASE_PASSWORD=your_password
export WALLET_TRACKER_APP_LOG_LEVEL=debug
```

## Usage

### Track Wallet

Basic wallet tracking:

```bash
./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

![CLI Screenshot](img/2022-06-23_19-39.png)

### Track with Network Specification

```bash
./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa --network BTC
```

### WebSocket Streaming

Stream all transactions in real-time:

```bash
./wallet-tracker tracker websocket --all
```

![WebSocket Streaming](img/2022-06-23_20-05.png)

### Exchange Detection

Detect transactions involving known exchange wallets:

```bash
./wallet-tracker tracker track --wallet 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa --detect-exchanges
```

### Visual Analytics

Start NeoDash for graph visualization:

```bash
# Start NeoDash server
./wallet-tracker neodash start
```

Then open http://localhost:5005 in your browser.

![NeoDash Overview](img/2022-06-23_20-01.png)

#### View specific transaction graphs:

![Transaction Graph](img/2022-06-23_20-02.png)

#### Explore the entire graph database:

![Full Database](img/2022-06-23_20-03.png)

### Get Exchange Wallets

Query known exchange wallets from Redis:

```bash
./wallet-tracker redis get --exchanges binance --limit 5
```

## Development

### Project Structure

```
wallet-tracker/
├── cmd/wallet-tracker/    # Main application entry point
├── cli/command/          # CLI commands implementation
├── domain/               # Domain models and interfaces
├── pkg/                  # Shared packages
│   ├── cache/           # Caching layer
│   ├── config/          # Configuration management
│   ├── errors/          # Error handling
│   ├── logger/          # Structured logging
│   ├── progress/        # Progress indicators
│   └── retry/           # Retry mechanisms
├── neodash/             # NeoDash dashboard configuration
└── docker-compose.yml   # Local development setup
```

### Building from Source

```bash
# Regular build
go build -o wallet-tracker cmd/wallet-tracker/main.go

# Build with version information
go build -ldflags "-X main.version=1.0.0" -o wallet-tracker cmd/wallet-tracker/main.go

# Cross-compilation examples
GOOS=linux GOARCH=amd64 go build -o wallet-tracker-linux-amd64 cmd/wallet-tracker/main.go
GOOS=darwin GOARCH=amd64 go build -o wallet-tracker-darwin-amd64 cmd/wallet-tracker/main.go
GOOS=windows GOARCH=amd64 go build -o wallet-tracker-windows-amd64.exe cmd/wallet-tracker/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/retry
```

## Architecture

The wallet tracker uses a modular architecture:

1. **CLI Layer**: Command-line interface using Cobra
2. **Domain Layer**: Business logic and models
3. **Infrastructure Layer**: External service integrations (blockchain APIs, Neo4j, Redis)
4. **Package Layer**: Shared utilities and cross-cutting concerns

### Technology Stack

- **Language**: Go
- **Database**: Neo4j (graph database)
- **Cache**: Redis
- **CLI Framework**: Cobra
- **Visualization**: NeoDash
- **Container**: Docker

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Update documentation
- Follow Go best practices
- Use conventional commits

## Roadmap

- [ ] Multi-blockchain support (Ethereum, BSC, Polygon)
- [ ] Advanced pattern detection algorithms
- [ ] REST API wrapper
- [ ] Real-time alerting system
- [ ] Machine learning integration
- [ ] Distributed processing support

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Original concept inspired by blockchain analysis needs
- Built with love by the open-source community
- Special thanks to all contributors

## Support

- 📧 Create an [issue](https://github.com/haghfizzuddin/wallet-tracker/issues) for bug reports
- 💬 Join our [discussions](https://github.com/haghfizzuddin/wallet-tracker/discussions) for questions
- 📖 Check the [wiki](https://github.com/haghfizzuddin/wallet-tracker/wiki) for detailed documentation

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/haghfizzuddin">haghfizzuddin</a>
</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/haghfizzuddin/wallet-tracker.svg?style=for-the-badge
[contributors-url]: https://github.com/haghfizzuddin/wallet-tracker/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/haghfizzuddin/wallet-tracker.svg?style=for-the-badge
[forks-url]: https://github.com/haghfizzuddin/wallet-tracker/network/members
[stars-shield]: https://img.shields.io/github/stars/haghfizzuddin/wallet-tracker.svg?style=for-the-badge
[stars-url]: https://github.com/haghfizzuddin/wallet-tracker/stargazers
[issues-shield]: https://img.shields.io/github/issues/haghfizzuddin/wallet-tracker.svg?style=for-the-badge
[issues-url]: https://github.com/haghfizzuddin/wallet-tracker/issues
[license-shield]: https://img.shields.io/github/license/haghfizzuddin/wallet-tracker.svg?style=for-the-badge
[license-url]: https://github.com/haghfizzuddin/wallet-tracker/blob/main/LICENSE

# Contributing to Wallet Tracker

We love your input! We want to make contributing to Wallet Tracker as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## We Develop with Github
We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## We Use [Github Flow](https://guides.github.com/introduction/flow/index.html)
Pull requests are the best way to propose changes to the codebase. We actively welcome your pull requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issues](https://github.com/yourusername/wallet-tracker/issues)
We use GitHub issues to track public bugs. Report a bug by [opening a new issue](); it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/wallet-tracker.git
   cd wallet-tracker
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up configuration**
   ```bash
   cd enhanced-analyzer
   cp enhanced-analyzer-config.json.example enhanced-analyzer-config.json
   # Add your Etherscan API key
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## Code Style

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused

## Adding New Features

### Detection Methods
To add a new detection method:

1. Add your detection function to `advanced_behavioral_analyzer.go`:
   ```go
   func (ba *BehavioralAnalyzer) detectNewPattern(txs []Transaction) []BehavioralFlag {
       // Your detection logic
   }
   ```

2. Call it from `analyzeBehavioralPatterns()`:
   ```go
   newFlags := ba.detectNewPattern(txs)
   flags = append(flags, newFlags...)
   ```

3. Update risk calculation if needed in `calculateFinalRiskScore()`

### Statistical Analysis
To add new statistical methods:

1. Add calculation to `performStatisticalAnalysis()`
2. Update `StatisticalScores` struct if needed
3. Document the mathematical approach

## Testing

- Write unit tests for new functions
- Test with real blockchain data
- Include edge cases
- Document test scenarios

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update API documentation if applicable
- Include examples

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the version numbers in any examples files to the new version
3. The PR will be merged once you have the sign-off of at least one maintainer

## Community

- Be respectful and inclusive
- Help others when you can
- Ask questions if you're unsure
- Share your knowledge

## Recognition

Contributors will be recognized in:
- The README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to Wallet Tracker! 🚀

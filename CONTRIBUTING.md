# Contributing to Wallet Tracker

First off, thank you for considering contributing to Wallet Tracker! It's people like you that make Wallet Tracker such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include logs and error messages**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain which behavior you expected**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go style guide
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

### Go Styleguide

* Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
* Use `gofmt` to format your code
* Use `golint` to lint your code
* Write meaningful variable names
* Comment exported functions and types

### Testing

* Write unit tests for new functionality
* Ensure all tests pass before submitting PR
* Aim for high test coverage
* Use table-driven tests where appropriate

## Project Structure

```
wallet-tracker/
├── cmd/              # Command line applications
├── cli/              # CLI implementation
├── domain/           # Business logic
├── pkg/              # Shared packages
└── tests/            # Integration tests
```

## Getting Started

1. Set up your development environment:
   ```bash
   git clone https://github.com/haghfizzuddin/wallet-tracker.git
   cd wallet-tracker
   go mod download
   ```

2. Run tests:
   ```bash
   go test ./...
   ```

3. Build the project:
   ```bash
   go build -o wallet-tracker cmd/wallet-tracker/main.go
   ```

## Questions?

Feel free to open an issue with your question or reach out to the maintainers directly.

Thank you for contributing! 🎉

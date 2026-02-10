# Contributing to golog

Thank you for considering contributing to golog! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue on GitHub with:
- A clear title and description
- Steps to reproduce the bug
- Expected behavior
- Actual behavior
- Go version and operating system

### Suggesting Enhancements

We welcome suggestions for new features or improvements. Please create an issue with:
- A clear title and description
- The motivation for the enhancement
- Examples of how it would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write tests** for any new functionality
3. **Ensure all tests pass** by running `make test`
4. **Format your code** by running `make fmt`
5. **Check for issues** by running `make vet`
6. **Update documentation** if needed
7. **Write clear commit messages**
8. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.25.0 or later
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/muleiwu/golog.git
cd golog

# Install dependencies
make install

# Run tests
make test

# Format code
make fmt

# Run all checks
make check
```

## Code Style

- Follow standard Go formatting (`gofmt`)
- Write clear, descriptive variable names
- Add comments for exported functions and types
- Keep functions focused and small
- Use meaningful commit messages

## Testing

- Write tests for all new functionality
- Ensure existing tests pass
- Aim for high code coverage
- Include both unit tests and benchmarks where appropriate

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

## Documentation

- Update README.md for significant changes
- Add godoc comments for all exported functions, types, and methods
- Update CHANGELOG.md with your changes
- Include examples for new features

## Commit Messages

Write clear and meaningful commit messages:

```
Add feature: brief description

Detailed explanation of what changed and why.
Include any relevant context or links to issues.
```

## Code Review Process

1. All submissions require review
2. We may suggest changes or improvements
3. Once approved, a maintainer will merge your PR

## License

By contributing to golog, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue if you have any questions about contributing!

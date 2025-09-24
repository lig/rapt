# Contributing to Rapt

Thank you for your interest in contributing to Rapt! This document provides guidelines and information for contributors.

## Code of Conduct

This project and everyone participating in it is governed by our commitment to providing a welcoming and inspiring community for all. By participating, you agree to uphold this code of conduct.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured to access your cluster

### Development Setup

1. **Fork the repository** on Codeberg
2. **Clone your fork**:
   ```bash
   git clone https://codeberg.org/YOUR_USERNAME/rapt.git
   cd rapt
   ```

3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://codeberg.org/lig/rapt.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Build the project**:
   ```bash
   go build -o rapt main.go
   ```

6. **Run tests**:
   ```bash
   go test ./...
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards below

3. **Test your changes**:
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

4. **Commit your changes** with a clear commit message:
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** on Codeberg

### Coding Standards

#### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use `gofmt -s` for simplified formatting
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

#### Commit Message Format

We use conventional commits format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

Examples:
```
feat: add support for custom tool arguments
fix: resolve CRD creation error in init command
docs: update README with installation instructions
```

#### Testing

- Write tests for new functionality
- Ensure all tests pass before submitting
- Aim for good test coverage
- Use table-driven tests where appropriate

#### Documentation

- Update documentation for new features
- Add examples for new commands
- Update README.md if needed
- Document any breaking changes

## Project Structure

```
rapt/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and global flags
│   ├── init.go         # Initialize CRD command
│   ├── add.go          # Add tool command
│   ├── purge.go        # Purge command
│   └── version.go      # Version command
├── internal/           # Internal packages
│   ├── app/rapt/       # Application logic
│   └── k8s/            # Kubernetes client and CRD handling
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── LICENSE              # MIT License
├── README.md           # Project documentation
└── CONTRIBUTING.md     # This file
```

## Areas for Contribution

### High Priority

- **Tool execution**: Implement the actual job execution functionality
- **Tool management**: Add commands to list, update, and delete tools
- **Error handling**: Improve error messages and handling
- **Testing**: Add comprehensive test coverage
- **Documentation**: Improve documentation and examples

### Medium Priority

- **Configuration**: Add configuration file support
- **Logging**: Implement structured logging
- **Metrics**: Add basic metrics collection
- **Validation**: Improve input validation
- **Performance**: Optimize Kubernetes API calls

### Low Priority

- **Web UI**: Consider a web interface for tool management
- **Plugin system**: Allow for extensible tool types
- **Integration**: Integrate with popular CI/CD systems

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Environment information**:
   - OS and version
   - Go version
   - Kubernetes version
   - Rapt version

2. **Steps to reproduce**:
   - Clear, numbered steps
   - Expected behavior
   - Actual behavior

3. **Additional context**:
   - Error messages
   - Logs (if applicable)
   - Screenshots (if applicable)

### Feature Requests

When requesting features, please include:

1. **Use case**: Why is this feature needed?
2. **Proposed solution**: How should it work?
3. **Alternatives**: What other approaches have you considered?
4. **Additional context**: Any other relevant information

## Review Process

1. **Automated checks**: All PRs must pass CI checks
2. **Code review**: At least one maintainer must approve
3. **Testing**: Changes must be tested
4. **Documentation**: Documentation must be updated if needed

## Release Process

Releases are managed by maintainers:

1. **Version bump**: Update version in code
2. **Changelog**: Update CHANGELOG.md
3. **Tag**: Create a git tag
4. **Build**: CI automatically builds and releases
5. **Announcement**: Announce the release

## Getting Help

- **Issues**: [Report bugs or request features](https://codeberg.org/lig/rapt/issues)
- **Discussions**: [Join the community](https://codeberg.org/lig/rapt/discussions)
- **Email**: Contact maintainers directly if needed

## License

By contributing to Rapt, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Rapt! 🚀
# Rapt

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Rapt (from pronunciation of "wrapped") is a CLI utility written in Go that orchestrates running predefined jobs (commands or tools) in a Kubernetes environment.

## Overview

Rapt enables you to define, manage, and execute jobs in Kubernetes clusters using a simple command-line interface. It is designed to streamline workflows by providing a consistent and extensible way to run and manage containerized tasks.

## Features

- **Kubernetes Integration**: Native support for Kubernetes clusters
- **Custom Resource Definitions**: Uses CRDs to define and manage tools
- **Simple CLI Interface**: Easy-to-use command-line interface
- **Extensible**: Define custom tools with flexible configuration
- **Namespace Support**: Work with specific Kubernetes namespaces

## Installation

### From Source

```bash
git clone https://codeberg.org/lig/rapt.git
cd rapt
go build -o rapt main.go
sudo mv rapt /usr/local/bin/
```

### Using Go Install

```bash
go install codeberg.org/lig/rapt@main
```

## Quick Start

1. **Initialize Rapt in your cluster**:
   ```bash
   rapt init
   ```

2. **Add a tool definition**:
   ```bash
   rapt add my-tool --image alpine:latest --command "echo hello world"
   ```

3. **Run a tool**:
   ```bash
   rapt run my-tool
   ```

4. **Clean up**:
   ```bash
   rapt purge
   ```

## Commands

### `rapt init`
Install the Rapt CRD in your Kubernetes cluster. This command sets up the necessary CustomResourceDefinition so that Rapt can manage and orchestrate predefined jobs in your cluster.

```bash
rapt init [--namespace <namespace>]
```

### `rapt add`
Add a new tool definition to your cluster. Tools are defined using Kubernetes Custom Resources.

```bash
rapt add <tool-name> [flags]
```

**Flags:**
- `--image`: Container image to run
- `--command`: Command to execute (overrides ENTRYPOINT)
- `--env`: Environment variables (format: KEY=VALUE)
- `--help-text`: Help message for the tool
- `--arg`: Tool arguments (format: name:description:required:default)

### `rapt purge`
Remove the Rapt CRD and all associated resources from your Kubernetes cluster.

```bash
rapt purge [--namespace <namespace>]
```

⚠️ **Warning**: This operation is destructive and cannot be undone.

## Tool Definition Schema

Tools are defined using Kubernetes Custom Resources with the following schema:

```yaml
apiVersion: rapt.dev/v1alpha1
kind: Tool
metadata:
  name: my-tool
spec:
  help: "Description of what this tool does"
  arguments:
    - name: "input"
      description: "Input file path"
      required: true
    - name: "output"
      description: "Output file path"
      required: false
      default: "output.txt"
  jobTemplate:
    image: "alpine:latest"
    command: ["sh", "-c", "echo 'Hello from my-tool'"]
    env:
      - name: "ENV_VAR"
        value: "example-value"
```

## Examples

### Simple Echo Tool
```bash
rapt add echo-tool \
  --image alpine:latest \
  --command "echo" \
  --arg "message:Message to echo:true:" \
  --help-text "Echo a message"
```

### Database Migration Tool
```bash
rapt add db-migrate \
  --image postgres:15 \
  --command "psql" \
  --env "PGHOST=db-service" \
  --env "PGUSER=migrator" \
  --arg "database:Database name:true:" \
  --arg "script:Migration script path:true:" \
  --help-text "Run database migrations"
```

## Development

### Prerequisites
- Go 1.24 or later
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured to access your cluster

### Building from Source
```bash
git clone https://codeberg.org/lig/rapt.git
cd rapt
go mod download
go build -o rapt main.go
```

### Running Tests
```bash
go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Support

- **Issues**: [Report bugs or request features](https://codeberg.org/lig/rapt/issues)
- **Discussions**: [Join the community](https://codeberg.org/lig/rapt/discussions)

## Roadmap

For detailed development plans and future features, see our [ROADMAP.md](ROADMAP.md).

**Recent Achievements:**
- ✅ Complete MVP functionality with all essential commands
- ✅ Real-time log streaming in `rapt run`
- ✅ Comprehensive job log management with `rapt logs`
- ✅ Enhanced job naming with readable timestamps

---

**Author**: Serge Matveenko <lig@countzero.co>
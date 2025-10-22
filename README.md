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
go install codeberg.org/lig/rapt@latest
```

## Quick Start

1. **Initialize Rapt in your cluster**:
   ```bash
   rapt init
   ```

2. **Add a tool definition**:
   ```bash
   rapt add echo-tool --image alpine:latest --command "echo Hello from Rapt!"
   ```

3. **List available tools**:
   ```bash
   rapt list
   ```

4. **Run the tool**:
   ```bash
   rapt run echo-tool
   ```

5. **View tool details**:
   ```bash
   rapt describe echo-tool
   ```

6. **Clean up**:
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
- `-i, --image`: (Required) Container image to run
- `-c, --command`: Command to execute (overrides ENTRYPOINT). Specify as a single string.
- `-e, --env`: Environment variables in the form NAME=VALUE. Can be specified multiple times.
- `--dry-run`: Print the Tool CR YAML without applying it to the cluster

**Examples:**
```bash
# Simple tool
rapt add lstool --image alpine:latest --command "ls -la"

# Tool with environment variables
rapt add echo --image busybox -e FOO=bar -e BAZ=qux --command "echo \$FOO \$BAZ"

# Preview without creating
rapt add my-tool --image alpine:latest --command "whoami" --dry-run
```

### `rapt run`
Execute a tool by creating a Kubernetes Job from the tool definition.

```bash
rapt run <tool-name> [flags]
```

**Flags:**
- `-a, --arg`: Tool argument in the form key=value. Can be specified multiple times.
- `-e, --env`: Environment variable in the form key=value. Can be specified multiple times.
- `-m, --mount`: Mount local file into container in the form local-path:container-path. Can be specified multiple times.
- `-w, --wait`: Wait for the job to complete before exiting
- `-f, --follow`: Follow job logs in real-time (default behavior)
- `-t, --timeout`: Timeout in seconds when waiting for job completion (default: 300)

**Examples:**
```bash
# Run a simple tool
rapt run echo-tool

# Run with environment variables
rapt run my-tool --env DEBUG=true

# Run with file mounts
rapt run script-runner --mount ./script.sh:/app/script.sh --mount ./config.yaml:/etc/config.yaml

# Run and wait for completion
rapt run data-processor --wait --timeout 600
```

**Note**: By default, logs are streamed in real-time, making it feel like running a local command.

### `rapt list`
List all available tools in the cluster.

```bash
rapt list [flags]
```

**Flags:**
- `-o, --output`: Output format: table, json, yaml (default: table)
- `-A, --all-namespaces`: List tools from all namespaces

### `rapt describe`
Show detailed information about a specific tool.

```bash
rapt describe <tool-name> [flags]
```

**Flags:**
- `-o, --output`: Output format: table, json, yaml (default: table)

### `rapt logs`
View logs from previous job executions.

```bash
rapt logs <tool-name> [job-name]
```

### `rapt status`
Check the status of jobs created from a tool.

```bash
rapt status <tool-name>
```

### `rapt delete`
Delete a tool definition from the cluster.

```bash
rapt delete <tool-name>
```

### `rapt purge`
Remove the Rapt CRD and all associated resources from your Kubernetes cluster.

```bash
rapt purge [--namespace <namespace>]
```

⚠️ **Warning**: This operation is destructive and cannot be undone.

## Tool Definition Schema

**Note**: The schema below is for direct Kubernetes resource creation using `kubectl` or for understanding the underlying CRD structure. When using `rapt add`, you don't need to write YAML manually - the tool creates it for you. However, `rapt add` currently has limited support for defining tool arguments and help text. For complex tools with arguments, you may need to create the YAML file directly and apply it with `kubectl apply -f tool.yaml`.

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

### Example 1: Simple Echo Tool

**Step 1: Add the tool**
```bash
rapt add echo-tool \
  --image alpine:latest \
  --command "echo Hello from echo-tool"
```

**Step 2: Run the tool**
```bash
rapt run echo-tool
```

**Step 3: View tool details**
```bash
rapt describe echo-tool
```

### Example 2: Tool with Environment Variables

**Step 1: Add the tool**
```bash
rapt add greeter \
  --image alpine:latest \
  --command "sh -c 'echo Hello \$NAME, welcome to \$PLACE'" \
  --env "NAME=User" \
  --env "PLACE=Kubernetes"
```

**Step 2: Run with default environment**
```bash
rapt run greeter
```

**Step 3: Run with custom environment variables**
```bash
rapt run greeter --env NAME=Alice --env PLACE=Production
```

### Example 3: File Processing Tool

**Step 1: Add the tool**
```bash
rapt add file-processor \
  --image alpine:latest \
  --command "sh -c 'cat /data/input.txt && echo Processed > /data/output.txt'"
```

**Step 2: Run with mounted files**
```bash
echo "Sample content" > input.txt
rapt run file-processor \
  --mount ./input.txt:/data/input.txt \
  --mount ./output.txt:/data/output.txt
```

**Step 3: Check the output**
```bash
cat output.txt
```

### Example 4: Database Tool

**Step 1: Add the tool**
```bash
rapt add db-query \
  --image postgres:15 \
  --command "psql" \
  --env "PGHOST=db-service.default.svc.cluster.local" \
  --env "PGUSER=readonly" \
  --env "PGDATABASE=mydb"
```

**Step 2: Run the tool**
```bash
rapt run db-query --env PGPASSWORD=secret
```

For more examples, see the [examples/](examples/) directory.

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
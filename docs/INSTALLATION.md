# Installation Guide

This guide covers different ways to install Rapt on your system.

## Prerequisites

- Go 1.24 or later (for building from source)
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured to access your cluster

## Installation Methods

### 1. Pre-built Binaries

Download the latest release from the [releases page](https://codeberg.org/lig/rapt/releases).

#### Linux
```bash
# Download and extract
curl -L https://codeberg.org/lig/rapt/releases/download/v1.0.0/rapt-v1.0.0-linux-amd64.tar.gz | tar -xz

# Make executable and move to PATH
chmod +x rapt
sudo mv rapt /usr/local/bin/
```

#### macOS
```bash
# Download and extract
curl -L https://codeberg.org/lig/rapt/releases/download/v1.0.0/rapt-v1.0.0-darwin-amd64.tar.gz | tar -xz

# Make executable and move to PATH
chmod +x rapt
sudo mv rapt /usr/local/bin/
```

#### Windows
```powershell
# Download and extract
Invoke-WebRequest -Uri "https://codeberg.org/lig/rapt/releases/download/v1.0.0/rapt-v1.0.0-windows-amd64.zip" -OutFile "rapt.zip"
Expand-Archive -Path "rapt.zip" -DestinationPath "."
Move-Item "rapt.exe" "C:\Windows\System32\"
```

### 2. Go Install

```bash
go install codeberg.org/lig/rapt@main
```

### 3. Build from Source

```bash
# Clone the repository
git clone https://codeberg.org/lig/rapt.git
cd rapt

# Build
go build -o rapt main.go

# Install
sudo mv rapt /usr/local/bin/
```

### 4. Package Managers

#### Homebrew (macOS)
```bash
brew install lig/rapt/rapt
```

#### Scoop (Windows)
```powershell
scoop bucket add lig https://codeberg.org/lig/scoop-bucket.git
scoop install rapt
```

#### Arch Linux (AUR)
```bash
yay -S rapt
```

## Verification

After installation, verify that Rapt is working correctly:

```bash
# Check version
rapt version

# Check help
rapt --help
```

## Initial Setup

1. **Initialize Rapt in your cluster**:
   ```bash
   rapt init
   ```

2. **Verify CRD installation**:
   ```bash
   kubectl get crd tools.rapt.dev
   ```

3. **Test with a simple tool**:
   ```bash
   rapt add echo-tool --image alpine:latest --command "echo" --arg "message:Message to echo:true:" --help-text "Echo a message"
   ```

## Uninstallation

### Remove Rapt from Cluster
```bash
rapt purge
```

### Remove Binary
```bash
# Remove from PATH
sudo rm /usr/local/bin/rapt

# Or if installed via package manager
# brew uninstall rapt
# scoop uninstall rapt
```

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Make sure the binary is executable
chmod +x rapt

# Check PATH
echo $PATH
which rapt
```

#### Kubernetes Connection Issues
```bash
# Test kubectl connection
kubectl cluster-info

# Check current context
kubectl config current-context
```

#### CRD Installation Fails
```bash
# Check cluster permissions
kubectl auth can-i create crd

# Check if CRD already exists
kubectl get crd tools.rapt.dev
```

### Getting Help

- Check the [troubleshooting guide](TROUBLESHOOTING.md)
- Report issues on [Codeberg](https://codeberg.org/lig/rapt/issues)
- Join discussions on [Codeberg](https://codeberg.org/lig/rapt/discussions)

## Development Installation

For development purposes:

```bash
# Clone repository
git clone https://codeberg.org/lig/rapt.git
cd rapt

# Install dependencies
go mod download

# Run tests
go test ./...

# Build with development version
go build -ldflags="-X 'codeberg.org/lig/rapt/cmd.Version=dev' -X 'codeberg.org/lig/rapt/cmd.Commit=$(git rev-parse HEAD)'" -o rapt main.go
```

## Upgrading

### From Previous Versions

1. **Backup existing tools** (if any):
   ```bash
   kubectl get tools -o yaml > tools-backup.yaml
   ```

2. **Download new version**:
   ```bash
   # Using go install
   go install codeberg.org/lig/rapt@main
   
   # Or download binary
   curl -L https://codeberg.org/lig/rapt/releases/download/v1.1.0/rapt-v1.1.0-linux-amd64.tar.gz | tar -xz
   ```

3. **Verify upgrade**:
   ```bash
   rapt version
   ```

4. **Test functionality**:
   ```bash
   rapt --help
   kubectl get crd tools.rapt.dev
   ```

### Breaking Changes

Check the [CHANGELOG](CHANGELOG.md) for breaking changes between versions. Some upgrades may require:

- Recreating tools with updated schemas
- Updating Kubernetes cluster permissions
- Migrating configuration files

## Security Considerations

### Cluster Permissions

Rapt requires the following permissions:
- Create/read/update/delete CustomResourceDefinitions
- Create/read/update/delete Custom Resources
- Create/read/update/delete Jobs
- Read ConfigMaps and Secrets (for environment variables)

### Network Security

- Rapt communicates with the Kubernetes API server
- Ensure proper network policies if using network segmentation
- Use TLS for all communications

### Container Security

- Tools run in Kubernetes Jobs with the specified container images
- Ensure container images are from trusted sources
- Use non-root users in containers when possible
- Apply appropriate security contexts and policies
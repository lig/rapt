# Rapt Custom Resource Definition (CRD) Schema

This document describes the schema for Rapt's Custom Resource Definition (CRD) used to define tools in Kubernetes.

## Overview

Rapt uses Kubernetes Custom Resource Definitions to define and manage tools. Each tool is represented as a `Tool` resource with the API version `rapt.dev/v1alpha1`.

## Tool Resource Schema

### Metadata

Standard Kubernetes metadata fields:

```yaml
metadata:
  name: string          # Required: Name of the tool
  namespace: string     # Optional: Kubernetes namespace (defaults to "default")
  labels: map[string]string    # Optional: Labels for organization
  annotations: map[string]string  # Optional: Annotations for metadata
```

### Spec

The `spec` section defines the tool's behavior and configuration:

```yaml
spec:
  help: string                    # Optional: Help text for the tool
  arguments: []Argument          # Optional: List of tool arguments
  jobTemplate: JobTemplate       # Required: Kubernetes Job template
```

#### Argument Schema

Each argument defines a parameter that the tool accepts:

```yaml
name: string          # Required: Argument name
description: string   # Optional: Human-readable description
required: boolean     # Optional: Whether argument is required (default: false)
default: string       # Optional: Default value for optional arguments
```

#### JobTemplate Schema

The `jobTemplate` defines how the tool's job will be executed:

```yaml
image: string         # Required: Container image to run
command: []string     # Optional: Command to execute (overrides ENTRYPOINT)
args: []string        # Optional: Arguments to pass to the command
env: []EnvVar         # Optional: Environment variables
```

#### EnvVar Schema

Environment variables for the job:

```yaml
name: string          # Required: Environment variable name
value: string         # Optional: Environment variable value
valueFrom: ValueFrom  # Optional: Source for the value (e.g., secret, configmap)
```

## Complete Example

```yaml
apiVersion: rapt.dev/v1alpha1
kind: Tool
metadata:
  name: database-backup
  namespace: production
  labels:
    app: database
    environment: production
  annotations:
    rapt.dev/version: "1.0.0"
    rapt.dev/author: "ops-team"
spec:
  help: "Create a backup of the production database"
  arguments:
    - name: "database"
      description: "Name of the database to backup"
      required: true
    - name: "compression"
      description: "Compression algorithm to use (gzip, bzip2, none)"
      required: false
      default: "gzip"
    - name: "retention-days"
      description: "Number of days to retain the backup"
      required: false
      default: "30"
  jobTemplate:
    image: "postgres:15-alpine"
    command: ["pg_dump"]
    args: ["--host=db-service", "--port=5432", "--username=backup-user"]
    env:
      - name: "PGPASSWORD"
        valueFrom:
          secretKeyRef:
            name: db-backup-secret
            key: password
      - name: "BACKUP_DIR"
          value: "/backups"
      - name: "COMPRESSION"
          value: "gzip"
```

## Field Descriptions

### Required Fields

- `metadata.name`: Unique name for the tool within the namespace
- `spec.jobTemplate.image`: Container image to use for the job

### Optional Fields

#### spec.help
- **Type**: `string`
- **Description**: Help text displayed when users run `rapt help <tool-name>`
- **Example**: `"Create a backup of the production database"`

#### spec.arguments
- **Type**: `[]Argument`
- **Description**: List of arguments that the tool accepts
- **Validation**: Each argument must have a unique name

#### spec.jobTemplate.command
- **Type**: `[]string`
- **Description**: Command to execute in the container (overrides ENTRYPOINT)
- **Example**: `["pg_dump", "--verbose"]`

#### spec.jobTemplate.args
- **Type**: `[]string`
- **Description**: Arguments to pass to the command
- **Example**: `["--host=localhost", "--port=5432"]`

#### spec.jobTemplate.env
- **Type**: `[]EnvVar`
- **Description**: Environment variables to set in the container
- **Note**: Can reference Kubernetes secrets and configmaps

## Validation Rules

### Tool Name
- Must be a valid Kubernetes resource name
- Must be unique within the namespace
- Cannot be changed after creation

### Arguments
- Argument names must be unique within a tool
- Required arguments cannot have default values
- Optional arguments should have default values

### Job Template
- Image must be a valid container image reference
- Command and args are mutually exclusive with container's ENTRYPOINT
- Environment variables must have unique names

## Best Practices

### Naming
- Use descriptive, kebab-case names for tools
- Include environment or purpose in the name when appropriate
- Avoid generic names like "tool" or "job"

### Arguments
- Provide clear, descriptive help text
- Use appropriate defaults for optional arguments
- Group related arguments logically
- Validate argument values when possible

### Images
- Use specific image tags instead of `latest`
- Prefer minimal base images (alpine, distroless)
- Ensure images are from trusted sources
- Document image requirements and dependencies

### Environment Variables
- Use environment variables for configuration
- Reference Kubernetes secrets for sensitive data
- Provide sensible defaults
- Document required environment variables

### Security
- Run containers as non-root users when possible
- Use read-only root filesystems when appropriate
- Limit container capabilities
- Use network policies to restrict access

## Migration and Versioning

### Versioning
- Use semantic versioning for tool versions
- Document breaking changes in tool updates
- Maintain backward compatibility when possible

### Migration
- Tools can be updated by modifying the CRD
- Breaking changes may require recreating the tool
- Test tool updates in non-production environments

## Troubleshooting

### Common Issues

1. **Invalid Image**: Ensure the image exists and is accessible
2. **Missing Arguments**: Check that all required arguments are provided
3. **Environment Variables**: Verify that referenced secrets/configmaps exist
4. **Permissions**: Ensure the service account has necessary permissions

### Debugging
- Use `kubectl describe tool <tool-name>` to see tool details
- Check job logs with `kubectl logs job/<job-name>`
- Verify CRD installation with `kubectl get crd tools.rapt.dev`

## Future Enhancements

Planned additions to the schema:
- Resource limits and requests
- Node selection and affinity
- Volume mounts and persistent storage
- Init containers
- Sidecar containers
- Health checks and probes
- Retry policies and backoff strategies
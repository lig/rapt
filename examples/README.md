# Rapt Tool Examples

This directory contains example tool definitions that demonstrate how to use Rapt with various types of workloads.

## Examples

### echo-tool.yaml
A simple tool that echoes messages to stdout. Demonstrates:
- Basic tool definition
- Required and optional arguments
- Environment variables
- Simple Alpine-based container

**Usage:**
```bash
rapt add echo-tool --image alpine:latest --command "echo" --arg "message:Message to echo:true:" --help-text "Echo a message"
```

### database-migrate.yaml
A database migration tool using Flyway. Demonstrates:
- Complex tool with multiple arguments
- Database connectivity
- Environment-specific configuration
- Professional tool integration

**Usage:**
```bash
rapt add db-migrate \
  --image flyway/flyway:latest \
  --command "flyway migrate" \
  --env "FLYWAY_URL=jdbc:postgresql://db-service:5432" \
  --env "FLYWAY_USER=migrator" \
  --arg "database:Database name:true:" \
  --arg "script:Migration script path:true:" \
  --help-text "Run database migrations using Flyway"
```

### file-processor.yaml
A file processing tool with multiple operations. Demonstrates:
- Multiple operation modes
- File path handling
- Algorithm selection
- Complex argument structure

**Usage:**
```bash
rapt add file-processor \
  --image alpine:latest \
  --command "sh -c 'echo Processing file with operation: $OPERATION'" \
  --arg "operation:Operation to perform (compress, encrypt, convert):true:" \
  --arg "input-file:Path to input file:true:" \
  --arg "output-file:Path to output file:false:output" \
  --help-text "Process files with various operations"
```

## Creating Your Own Tools

When creating your own tools, consider:

1. **Image Selection**: Choose appropriate base images
   - `alpine:latest` for lightweight tools
   - `ubuntu:latest` for tools requiring more packages
   - Specific tool images (e.g., `node:18`, `python:3.11`)

2. **Arguments**: Define clear, descriptive arguments
   - Use descriptive names
   - Provide helpful descriptions
   - Set appropriate defaults
   - Mark required arguments clearly

3. **Environment Variables**: Use environment variables for:
   - Configuration
   - Secrets (consider using Kubernetes secrets)
   - Runtime parameters

4. **Commands**: Design commands that are:
   - Idempotent when possible
   - Well-documented
   - Handle errors gracefully

## Best Practices

### Security
- Use non-root users when possible
- Avoid hardcoded secrets
- Use Kubernetes secrets for sensitive data
- Validate input parameters

### Performance
- Use appropriate resource limits
- Consider parallel execution
- Optimize container images
- Use efficient base images

### Reliability
- Handle errors gracefully
- Provide meaningful error messages
- Use health checks when appropriate
- Implement retry logic for transient failures

### Maintainability
- Keep tools focused and single-purpose
- Document tool behavior clearly
- Version your tools appropriately
- Test tools thoroughly

## Contributing Examples

If you have a useful tool example, please:

1. Create a new YAML file with a descriptive name
2. Include comprehensive comments
3. Add usage instructions
4. Test the example thoroughly
5. Submit a pull request

We welcome examples for:
- Development tools (linters, formatters, test runners)
- Data processing (ETL, analytics, reporting)
- Infrastructure (monitoring, backup, deployment)
- Security (scanners, auditors, compliance)
- And more!
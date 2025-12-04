---
description: Build the issue2md project using make build
allowed-tools:
  - Bash(make:build)
---

# Build the issue2md Project

Execute `make build` to compile the issue2md-cli and issue2md-web binaries.

If the build fails, analyze the error output and provide troubleshooting steps based on the specific failure type.

## Steps:

1. First, run `make build` to compile both CLI and web applications
2. If build succeeds, confirm the binaries were created in the bin/ directory
3. If build fails, analyze the error output and determine the root cause:
   - Go compilation errors (syntax, type errors, missing imports)
   - Missing dependencies
   - Permission issues
   - Environment problems

## Build Success Output Expectations:
- Both issue2md-cli and issue2md-web binaries in bin/ directory
- No error messages in the build output
- Build completion message with binary paths

## Common Error Analysis:
- **Go compile errors**: Check for syntax issues, missing imports, type mismatches
- **Module issues**: Run `go mod tidy` and check go.mod file
- **Permission errors**: Verify write permissions to bin/ directory
- **Go version issues**: Ensure Go 1.21+ is installed

Provide specific remediation steps based on the actual error encountered.
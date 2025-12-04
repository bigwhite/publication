---
description: Build Docker image for the issue2md project using make docker-build with optional tag parameter
allowed-tools:
  - Bash(docker:*)
  - Bash(make:docker-build)
parameters:
  - name: tag
    description: Docker image tag (defaults to 'latest' if not provided)
    required: false
    default: "latest"
---

# Build Docker Image for issue2md Project

Build a Docker image for the issue2md project using `make docker-build`.

## Usage:
- `/docker-build` - Build image with tag 'latest' (default)
- `/docker-build v1.0.0` - Build image with tag 'v1.0.0'
- `/docker-build my-custom-tag` - Build image with tag 'my-custom-tag'

## Steps:

1. Check if Docker is running and accessible
2. Verify the Dockerfile exists in the project root
3. Set the DOCKER_TAG environment variable based on the provided parameter
4. Execute `make docker-build` with the specified tag
5. Verify the image was built successfully
6. Display image information

## Expected Results:
- Docker image `issue2md:{tag}` should be created
- Build should complete without errors
- Image size and creation time should be displayed

## Error Handling:
If the build fails, analyze and provide solutions for:
- Docker daemon not running
- Missing Dockerfile
- Build context errors
- Permission issues
- Insufficient disk space
- Network connectivity issues for dependency downloads

## Environment Variables:
The command will set `DOCKER_TAG={provided_tag}` before running make docker-build.

## Examples:
```bash
/docker-build                    # Uses tag 'latest'
/docker-build v1.2.3            # Uses tag 'v1.2.3'
/docker-build production-build  # Uses tag 'production-build'
```
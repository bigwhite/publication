#!/bin/sh

# docker-entrypoint.sh for issue2md
# Supports both CLI and web service modes

set -e

# Function to show usage
show_usage() {
    cat << EOF
issue2md Docker Entrypoint

Usage:
  docker run [options] issue2md [command] [args]

Commands:
  web              Start web service (default)
  cli [args]       Run CLI with provided arguments
  help             Show this help message

Examples:
  # Start web service
  docker run -p 8080:8080 issue2md

  # Run CLI command
  docker run issue2md cli facebook/react 12345 --output=issue.md

  # Run CLI with help
  docker run issue2md cli --help

Environment Variables:
  PORT             Web service port (default: 8080)
  GITHUB_TOKEN     GitHub API token
EOF
}

# Main execution logic
case "${1:-web}" in
    "web")
        echo "Starting issue2md web service on port ${PORT:-8080}..."
        exec ./issue2mdweb
        ;;
    "cli")
        shift  # Remove 'cli' from arguments
        echo "Running issue2md CLI with args: $*"
        exec ./issue2md "$@"
        ;;
    "help"|"-h"|"--help")
        show_usage
        exit 0
        ;;
    *)
        # If no command specified, default to web service
        if [ $# -eq 0 ]; then
            echo "No command specified, starting web service..."
            exec ./issue2mdweb
        else
            # Otherwise, treat as CLI command
            echo "Running issue2md CLI with args: $*"
            exec ./issue2md "$@"
        fi
        ;;
esac
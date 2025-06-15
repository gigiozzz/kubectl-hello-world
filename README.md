# kubectl-hello-world

A kubectl plugin demonstrating best practices for building kubectl extensions using Go, Cobra, and Kubernetes cli-runtime.

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)
[![Kubernetes](https://img.shields.io/badge/kubernetes-1.28+-blue.svg)](https://kubernetes.io)

## Overview

This project provides a single kubectl plugin with multiple subcommands showcasing kubectl plugin development practices:

- **`kubectl hello-world greetings`** - Basic hello world functionality with cluster information
- **`kubectl hello-world list-pods`** - Advanced pod listing with multiple output formats

The plugin demonstrates proper project organization, shared utilities, and integration with kubectl's library and configuration system following Kubernetes community standards.

## Features

- ✅ **kubectl Integration** - Seamless integration with kubectl configuration (contexts, namespaces, kubeconfig)
- ✅ **Multiple Output Formats** - Support for JSON, YAML, table, wide, and name formats
- ✅ **Cobra CLI Framework** - Professional command-line interface with help and flag parsing
- ✅ **Modular Project Structure** - Modular command organization with shared utilities
- ✅ **Shared Configuration** - Common kubectl flags and client setup across commands
- ✅ **Error Handling** - Proper error handling and user-friendly messages
- ✅ **Logging** - Structured logging using klog (same as kubectl)
- ✅ **Best Practices** - Follows kubectl plugin development patterns

## Installation

### Prerequisites

- Go 1.23 or later
- kubectl installed and configured
- Access to a Kubernetes cluster
- Docker (for containerized builds)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/gigiozzz/kubectl-hello-world.git
cd kubectl-hello-world

# Build using Makefile (recommended)
make build

# Or build manually
go mod tidy
go build -o bin/kubectl-hello-world .

# Install the plugin
make install
# Or manually: sudo cp bin/kubectl-hello-world /usr/local/bin/

# Verify installation
kubectl plugin list
```

### Using Makefile (Recommended)

The project includes a comprehensive Makefile for easy development:

```bash
# Show all available targets
make help

# Development workflow
make deps          # Download dependencies
make fmt           # Format code
make lint          # Run linters
make test          # Run tests
make build         # Build binary
make install       # Install to /usr/local/bin

# All-in-one check
make check         # Run fmt, lint, vet, and test

# Build for all platforms
make build-all     # Creates binaries in dist/

# Docker workflow
make docker-build  # Build Docker image
make docker-run    # Run in container
make docker-push   # Push to registry

# Testing
make test-install      # Test installation
make integration-test  # Run integration tests

# Cleanup
make clean         # Clean build artifacts
make clean-all     # Clean everything including Docker
```

### Docker Installation

```bash
# Pull the image
docker pull docker.io/gigiozzz/kubectl-hello-world:latest

# Run with your kubeconfig
docker run --rm -it \
  -v ~/.kube:/root/.kube:ro \
  docker.io/gigiozzz/kubectl-hello-world:latest greetings --name "Docker User"

# Create an alias for easier usage
alias kubectl-hello-world='docker run --rm -it -v ~/.kube:/root/.kube:ro docker.io/gigiozzz/kubectl-hello-world:latest'
kubectl-hello-world list-pods
```

## Usage

### Root Command

```bash
# Show help and available subcommands
kubectl hello-world --help
```

### Greetings Subcommand

Basic hello world functionality with cluster information:

```bash
# Simple greeting
kubectl hello-world greetings
# Output: Hello, World! 👋

# Greet someone specific
kubectl hello-world greetings Alice
kubectl hello-world greetings --name Gigi

# Use with different contexts/namespaces
kubectl hello-world greetings --context=production --namespace=default Bob

# Verbose output (shows all namespaces)
kubectl hello-world greetings --verbose Alice
```

### List-Pods Subcommand

Advanced pod listing with multiple output formats:

```bash
# Default table format
kubectl hello-world list-pods

# Wide format (like kubectl get pods -o wide)
kubectl hello-world list-pods -o wide

# JSON output
kubectl hello-world list-pods -o json

# YAML output
kubectl hello-world list-pods -o yaml

# Names only
kubectl hello-world list-pods -o name

# Different namespace
kubectl hello-world list-pods --namespace=kube-system

# Different context
kubectl hello-world list-pods --context=production
```

## Architecture

### Shared Utilities (`cmd/util.go`)

The project uses a `CommonOptions` struct that provides shared functionality across all commands:

```go
type CommonOptions struct {
    ConfigFlags *genericclioptions.ConfigFlags  // kubectl config integration
    Clientset   kubernetes.Interface            // Kubernetes client
    Namespace   string                          // Current namespace
    Context     string                          // Current context
    Kubeconfig  string                         // Kubeconfig file path
    ClusterName string                         // Current cluster name
}
```

**Key shared methods:**
- `Complete()` - Initialize Kubernetes client and configuration
- `Validate()` - Perform common validation checks
- `AddConfigFlags()` - Add standard kubectl flags to commands
- `PrintContextInfo()` - Display current context information
- `GetNamespaceCount()` - Get cluster namespace count
- `ListNamespaces()` - List all namespace names

### Command Structure

Each command follows the same pattern:

```go
type CommandOptions struct {
    *CommonOptions              // Embed shared options
    // Command-specific fields
}

func (o *CommandOptions) Complete(cmd *cobra.Command, args []string) error {
    // 1. Complete common options
    // 2. Handle command-specific arguments and flags
}

func (o *CommandOptions) Validate() error {
    // 1. Validate common options
    // 2. Perform command-specific validation
}

func (o *CommandOptions) Run() error {
    // Execute command logic
}
```


**Happy kubectl plugin development! 🚀**

*This project serves as a comprehensive example of professional kubectl plugin development with proper architecture, shared utilities, and best practices that scale.*
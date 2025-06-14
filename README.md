# kubectl-hello-world

A kubectl plugin demonstrating best practices for building kubectl extensions using Go, Cobra, and Kubernetes cli-runtime.

## Features

- ✅ **kubectl Integration** - Seamless integration with kubectl configuration (contexts, namespaces, kubeconfig)
- ✅ **Multiple Output Formats** - Support for JSON, YAML, table, wide, and name formats
- ✅ **Cobra CLI Framework** - Professional command-line interface with help and flag parsing
- ✅ **Error Handling** - Proper error handling and user-friendly messages
- ✅ **Logging** - Structured logging using klog (same as kubectl)
- ✅ **Best Practices** - Follows kubectl plugin development patterns

## Installation

### Prerequisites

- Go 1.21 or later
- kubectl installed and configured
- Access to a Kubernetes cluster

### Build from Source

```bash
# Clone the repository
git clone https://github.com/gigiozzz/kubectl-hello-world.git
cd kubectl-hello-world

# Build the hello plugin
go mod tidy
go build -o kubectl-hello .
sudo mv kubectl-hello /usr/local/bin/

# Verify installation
kubectl plugin list
```

## Usage

```bash

# Greet someone specific
kubectl hello Alice
# Output: Hello, Alice! 👋

# Use with different contexts/namespaces
kubectl hello --context=production --namespace=default Bob
kubectl hello --kubeconfig=/path/to/config Alice

# Verbose output
kubectl hello --verbose Alice

# Help
kubectl hello --help
```

**Happy kubectl plugin development! 🚀**
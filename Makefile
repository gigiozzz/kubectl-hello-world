# Makefile for Go Kubectl plugin

# Variables
BINARY_NAME := kubectl-hello-world
DOCKER_REGISTRY := docker.io/gigiozzz
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(BINARY_NAME)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_VERSION := 1.23
LDFLAGS := -s -w \
	-X github.com/gigiozzz/kubectl-hello-world/cmd.Version=$(VERSION) \
	-X github.com/gigiozzz/kubectl-hello-world/cmd.CommitHash=$(COMMIT_HASH) \
	-X github.com/gigiozzz/kubectl-hello-world/cmd.BuildDate=$(BUILD_DATE)

# Default target
.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...
	goimports -w .

.PHONY: lint
lint: ## Run linters
	@echo "🔍 Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: vet
vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

.PHONY: check
check: fmt lint vet test ## Run all checks (format, lint, vet, test)

##@ Build

.PHONY: build
build: deps ## Build the binary
	@echo "🔨 Building $(BINARY_NAME) $(VERSION)..."
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) .

.PHONY: build-all
build-all: deps ## Build for all platforms
	@echo "🔨 Building for all platforms..."
	@mkdir -p dist
	# Linux amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 .
	# Linux arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 .
	# macOS amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	# macOS arm64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	# Windows amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Built binaries in dist/"

.PHONY: install
install: build ## Install binary to /usr/local/bin
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"

.PHONY: uninstall
uninstall: ## Uninstall binary from /usr/local/bin
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled successfully"

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image $(DOCKER_IMAGE):$(VERSION)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.

.PHONY: docker-push
docker-push: docker-build ## Push Docker image to registry
	@echo "🚀 Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run --rm -it \
		-v ~/.kube:/root/.kube:ro \
		$(DOCKER_IMAGE):$(VERSION) greetings --name "Docker User"

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "🧹 Cleaning Docker images..."
	docker rmi $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest 2>/dev/null || true

##@ Testing

.PHONY: test-install
test-install: install ## Test installation
	@echo "🧪 Testing installation..."
	kubectl plugin list | grep $(BINARY_NAME) || (echo "❌ Plugin not found" && exit 1)
	$(BINARY_NAME) --help
	@echo "✅ Installation test passed"

.PHONY: test-commands
test-commands: ## Test all commands
	@echo "🧪 Testing commands..."
	@echo "Testing greetings command..."
	$(BINARY_NAME) greetings --name "Test User"
	@echo "Testing list-pods command..."
	$(BINARY_NAME) list-pods --help
	@echo "✅ Command tests passed"

.PHONY: integration-test
integration-test: build ## Run integration tests (requires kubectl access)
	@echo "🧪 Running integration tests..."
	@echo "Testing greetings with verbose output..."
	./bin/$(BINARY_NAME) greetings --verbose --name "Integration Test"
	@echo "Testing list-pods..."
	./bin/$(BINARY_NAME) list-pods || echo "No pods found (expected in some environments)"
	@echo "Testing list-pods with different formats..."
	./bin/$(BINARY_NAME) list-pods -o name || echo "No pods found"
	@echo "✅ Integration tests completed"

##@ Release

.PHONY: release-dry-run
release-dry-run: build-all ## Dry run release process
	@echo "🎯 Release dry run for version $(VERSION)..."
	@echo "Would create release with:"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT_HASH)"
	@echo "  Date:    $(BUILD_DATE)"
	@echo "  Binaries:"
	@ls -la dist/

.PHONY: release
release: check build-all docker-build ## Create a release
	@echo "🎯 Creating release $(VERSION)..."
	@if [ "$(VERSION)" = "dev" ]; then \
		echo "❌ Cannot release with version 'dev'. Please set VERSION or create a git tag."; \
		exit 1; \
	fi
	@echo "✅ Release $(VERSION) ready"
	@echo "📦 Binaries built in dist/"
	@echo "🐳 Docker images built"
	@echo "Next steps:"
	@echo "  1. Push to registry: make docker-push"
	@echo "  2. Create GitHub release with binaries in dist/"

##@ Cleanup

.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean -cache

.PHONY: clean-all
clean-all: clean docker-clean ## Clean everything including Docker images

##@ Info

.PHONY: info
info: ## Show build info
	@echo "📋 Build Information:"
	@echo "  Binary:     $(BINARY_NAME)"
	@echo "  Version:    $(VERSION)"
	@echo "  Commit:     $(COMMIT_HASH)"
	@echo "  Build Date: $(BUILD_DATE)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Docker:     $(DOCKER_IMAGE):$(VERSION)"

.PHONY: version
version: ## Show version
	@echo $(VERSION)
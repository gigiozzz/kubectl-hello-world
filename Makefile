# Makefile for Go Kubectl plugin

# Variables
APP_NAME := kubectl-hello-world
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")
REGISTRY := docker.io
USERNAME := gigiozzz
IMAGE_NAME := $(REGISTRY)/$(USERNAME)/$(APP_NAME)
TAGGED_IMAGE := $(IMAGE_NAME):$(VERSION)
LATEST_IMAGE := $(IMAGE_NAME):latest

# Docker build arguments
DOCKER_BUILD_ARGS := --platform=linux/amd64,linux/arm64

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build build-dev build-prod push push-latest run run-dev stop clean test docker-login check-deps

# Default target
help: ## Show this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Prerequisites check
check-deps: ## Check if required tools are installed
	@echo "$(GREEN)Checking dependencies...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed$(NC)"; exit 1; }
# @command -v git >/dev/null 2>&1 || { echo "$(RED)Git is required but not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All dependencies found$(NC)"

# Initialize go module if not exists
init: ## Initialize Go module
#	@if [ ! -f go.mod ]; then \
#		echo "$(YELLOW)Initializing Go module...$(NC)"; \
#		go mod init $(APP_NAME); \
#	fi
	@go mod tidy

# Local development
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

run-dev: ## Run application locally
	@echo "$(GREEN)Starting development plugin...$(NC)"
	go run main.go

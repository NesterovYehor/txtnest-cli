# Makefile for txtnest-cli project

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
BUILD_DIR := build
CLI_BINARY := $(BUILD_DIR)/txtnest-cli
SSH_BINARY := $(BUILD_DIR)/txtnest-ssh

# Default target: show help
all: help

# Build the CLI tool
build-cli:
	$(GOBUILD) -o $(CLI_BINARY) ./cmd/cli/main.go

# Build the SSH tool
build-ssh:
	$(GOBUILD) -o $(SSH_BINARY) ./cmd/ssh/main.go

# Run the CLI tool directly
run-cli:
	$(GORUN) ./cmd/cli/main.go

# Run the SSH tool directly
run-ssh:
	$(GORUN) ./cmd/ssh/main.go

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Show available make targets
help:
	@echo "Makefile for txtnest-cli project"
	@echo "Usage:"
	@echo "  make build-cli     Build the CLI tool"
	@echo "  make build-ssh     Build the SSH tool"
	@echo "  make run-cli       Run the CLI tool"
	@echo "  make run-ssh       Run the SSH tool"
	@echo "  make clean         Remove all build artifacts"


# Makefile for koto

# Version information
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Build flags
LDFLAGS := -ldflags "\
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.date=$(DATE)'"

# Directories
BIN_DIR := bin
CMD_DIR := cmd/koto

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building koto $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/koto ./$(CMD_DIR)
	@echo "Build complete: $(BIN_DIR)/koto"

# Build for release (with GoReleaser)
.PHONY: release
release:
	goreleaser release --clean

# Build snapshot (without publishing)
.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --clean

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -rf dist
	rm -f coverage.out coverage.html

# Install the binary
.PHONY: install
install: build
	go install $(LDFLAGS) ./$(CMD_DIR)

# Run the application
.PHONY: run
run: build
	./$(BIN_DIR)/koto

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run linters
.PHONY: lint
lint:
	golangci-lint run

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install the binary"
	@echo "  make run            - Build and run the application"
	@echo "  make release        - Build release with GoReleaser"
	@echo "  make snapshot       - Build snapshot with GoReleaser"
	@echo "  make version        - Show version information"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linters"
	@echo "  make help           - Show this help message"

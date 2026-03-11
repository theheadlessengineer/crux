.PHONY: help build test test-integration lint fmt vet clean hooks coverage dev

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Compile binary to ./bin/crux
	@echo "Building crux..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o ./bin/crux ./cmd/crux
	@echo "✓ Build complete: bin/crux"

test: ## Run all tests with race detector
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests passed"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...
	@echo "✓ Integration tests passed"

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run
	@echo "✓ Linting passed"

fmt: ## Format all Go files
	@echo "Formatting code..."
	@gofmt -w .
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet passed"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out
	@echo "✓ Clean complete"

hooks: ## Install pre-commit hooks
	@echo "Installing pre-commit hooks..."
	@mkdir -p .git/hooks
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✓ Pre-commit hooks installed"

coverage: test ## Run tests and open coverage report
	@go tool cover -html=coverage.out

dev: fmt lint test ## Run format, lint, and test (development workflow)
	@echo "✅ Development checks passed"

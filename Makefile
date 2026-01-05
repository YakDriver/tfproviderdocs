.PHONY: build test clean install fmt lint modern modern-check deps test-coverage ci copyright-check tidy vet goreleaser-check help

default: build

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build tool
	@go build ./...

install: ## Install tool
	@go install .

test: ## Run tests
	@go test ./...

test-coverage: ## Run tests with coverage
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

vet: ## Run go vet
	@go vet ./...

fmt: ## Format code
	@go fmt ./...

tidy: ## Tidy go.mod
	@go mod tidy

deps: ## Download dependencies
	@go mod download

lint: ## Run golangci-lint
	@golangci-lint run

modern-check: ## Check for modern Go code
	@echo "make: Checking for modern Go code..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...

modern: ## Fix modern Go code issues
	@echo "make: Fixing checks for modern Go code..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

copyright-check: ## Check copyright headers
	@echo "make: Checking copyright headers..."
	@go run github.com/YakDriver/copyplop@latest check

goreleaser-check: ## Check goreleaser configuration
	@goreleaser check

ci: tidy build test vet modern-check copyright-check ## Run all CI checks locally

clean: ## Clean build artifacts
	@rm -rf bin/ coverage.out coverage.html

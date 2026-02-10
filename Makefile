.PHONY: test build clean fmt vet lint install examples help

# Variables
BINARY_NAME=golog
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.out

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

fmt: ## Format code
	$(GOFMT) ./...

vet: ## Run go vet
	$(GOVET) ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

tidy: ## Tidy go modules
	$(GO) mod tidy

install: ## Install dependencies
	$(GO) mod download

examples: ## Build all examples
	@echo "Building examples..."
	@$(GOBUILD) -o bin/basic examples/basic/main.go
	@$(GOBUILD) -o bin/webserver examples/webserver/main.go
	@$(GOBUILD) -o bin/custom examples/custom/main.go
	@echo "Examples built in bin/ directory"

run-basic: ## Run basic example
	$(GO) run examples/basic/main.go

run-webserver: ## Run webserver example
	$(GO) run examples/webserver/main.go

run-custom: ## Run custom example
	$(GO) run examples/custom/main.go

clean: ## Clean build artifacts
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Cleaned build artifacts"

check: fmt vet test ## Run all checks (fmt, vet, test)

ci: check ## Run CI pipeline
	@echo "CI pipeline completed successfully"

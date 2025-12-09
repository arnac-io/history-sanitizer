.PHONY: build test clean install run help

# Binary name
BINARY_NAME=history-sanitizer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: mod-download ## Build the binary (downloads deps if needed)
	$(GOBUILD) -o $(BINARY_NAME) -v

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

deps: mod-download mod-tidy ## Download and tidy dependencies (development)

mod-download: ## Download Go module dependencies
	@$(GOMOD) download

mod-tidy: ## Clean up go.mod and go.sum (run manually after adding/removing deps)
	$(GOMOD) tidy

mod-verify: ## Verify dependencies have expected content (good for CI)
	$(GOMOD) verify

mod-update: ## Update all dependencies to latest versions
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

mod-vendor: ## Vendor dependencies (copy deps to vendor/)
	$(GOMOD) vendor

mod-check: ## Check if go.mod needs tidying (non-destructive)
	@$(GOMOD) tidy -v 2>&1 | grep -q "unused" && echo "⚠️  Run 'make mod-tidy' to clean up go.mod" || echo "✓ go.mod is clean"

install: build ## Install the binary
	mv $(BINARY_NAME) $(GOPATH)/bin/

run: build ## Build and run with default settings
	./$(BINARY_NAME) --help

# Cross compilation
build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux -v

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v

build-mac: ## Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-macos-amd64 -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-macos-arm64 -v

build-all: build-linux build-windows build-mac ## Build for all platforms

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

check: fmt vet test ## Run all checks (format, vet, test)

ci: mod-verify check build ## Run all CI checks (for continuous integration)


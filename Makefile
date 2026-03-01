.PHONY: build run test lint fmt vet clean install snapshot check tidy deps help

# Default target
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# --- Build ---

build: ## Build the binary
	cd projectsCLI && go build -o projects ./cmd/projects

run: build ## Build and run (usage: make run ARGS="ls")
	./projectsCLI/projects $(ARGS)

install: build ## Install to /usr/local/bin
	cp projectsCLI/projects /usr/local/bin/projects

# --- Quality ---

test: ## Run tests
	cd projectsCLI && go test ./...

lint: ## Run golangci-lint
	cd projectsCLI && golangci-lint run ./...

fmt: ## Format code
	cd projectsCLI && gofmt -w .

vet: ## Run go vet
	cd projectsCLI && go vet ./...

check: fmt vet test ## Run all checks (fmt + vet + test)

# --- Release ---

snapshot: ## Build snapshot release (no publish)
	goreleaser release --snapshot --clean

# --- Clean ---

clean: ## Remove build artifacts
	rm -f projectsCLI/projects
	rm -rf dist/

# --- Dependencies ---

tidy: ## Run go mod tidy
	cd projectsCLI && go mod tidy

deps: ## Download dependencies
	cd projectsCLI && go mod download

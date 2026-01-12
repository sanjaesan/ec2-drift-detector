.PHONY: build test coverage run clean help install

BINARY_NAME=drift-detector
CMD_PATH=./cmd/drift-detector
COVERAGE_FILE=coverage.out

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod verify
	go mod tidy

build: ## Build the application
	go build -o $(BINARY_NAME) $(CMD_PATH)

test: ## Run tests
	go test -v ./...

coverage: ## Run tests with coverage
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report: coverage.html"
	@go tool cover -func=$(COVERAGE_FILE) | grep total

test-race: ## Run tests with race detection
	go test -race ./...

run-mock: ## Run with mock data
	go run $(CMD_PATH) \
		--instances=i-1234567890abcdef0,i-0987654321fedcba0 \
		--terraform-state=testdata/terraform.tfstate \
		--mock

run-mock-concurrent: ## Run with mock data (concurrent)
	go run $(CMD_PATH) \
		--instances=i-1234567890abcdef0,i-0987654321fedcba0 \
		--terraform-state=testdata/terraform.tfstate \
		--mock \
		--concurrent

run-mock-json: ## Run with JSON output
	go run $(CMD_PATH) \
		--instances=i-1234567890abcdef0,i-0987654321fedcba0 \
		--terraform-state=testdata/terraform.tfstate \
		--mock \
		--format=json

run: ## Run with real AWS (set INSTANCES variable)
	@if [ -z "$(INSTANCES)" ]; then \
		echo "Usage: make run INSTANCES=i-xxx,i-yyy"; \
		exit 1; \
	fi
	go run $(CMD_PATH) --instances=$(INSTANCES) --terraform-state=testdata/terraform.tfstate

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Install golangci-lint first" && exit 1)
	golangci-lint run ./...

all: fmt vet test build ## Run all checks and build
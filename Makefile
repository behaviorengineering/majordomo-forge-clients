.DEFAULT_GOAL := help

.PHONY: help build test tidy fmt vet lint smoke hooks-install

GOWORK ?= off

help: ## List available make targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

build: ## Build gitvalet CLI into bin/
	GOWORK=$(GOWORK) go build -o bin/gitvalet ./cmd/gitvalet

test: ## Run unit tests with race detector
	GOWORK=$(GOWORK) go test -race -count=1 ./...

tidy: ## Run go mod tidy
	GOWORK=$(GOWORK) go mod tidy

fmt: ## Format Go sources with gofmt
	gofmt -w .

vet: ## Run go vet
	GOWORK=$(GOWORK) go vet ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

hooks-install: ## Use .githooks/pre-commit (gofmt on staged .go files)
	chmod +x .githooks/pre-commit
	git config core.hooksPath .githooks

smoke: build ## Smoke-test CLI agent guide, help, version, usage
	./bin/gitvalet
	./bin/gitvalet help
	./bin/gitvalet -h
	./bin/gitvalet version
	./bin/gitvalet nope || test $$? = 2
	./bin/gitvalet sync --dry-run

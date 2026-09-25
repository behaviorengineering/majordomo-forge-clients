.PHONY: build test tidy fmt vet lint

GOWORK ?= off

build:
	GOWORK=$(GOWORK) go build -o bin/majordomo-forge ./cmd/majordomo-forge

test:
	GOWORK=$(GOWORK) go test -race -count=1 ./...

tidy:
	GOWORK=$(GOWORK) go mod tidy

fmt:
	gofmt -w .

vet:
	GOWORK=$(GOWORK) go vet ./...

lint:
	golangci-lint run ./...

smoke:
	./bin/majordomo-forge
	./bin/majordomo-forge help
	./bin/majordomo-forge version
	./bin/majordomo-forge nope || test $$? = 2
	./bin/majordomo-forge sync --dry-run

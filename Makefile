.PHONY: build test tidy fmt vet lint

GOWORK ?= off

build:
	GOWORK=$(GOWORK) go build -o bin/gitvalet ./cmd/gitvalet

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
	./bin/gitvalet
	./bin/gitvalet help
	./bin/gitvalet version
	./bin/gitvalet nope || test $$? = 2
	./bin/gitvalet sync --dry-run

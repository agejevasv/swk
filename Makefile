BINARY := swk
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/agejevasv/swk/cmd.Version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test ./...

test-verbose:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

lint:
	go vet ./...
	staticcheck ./...

clean:
	rm -f $(BINARY) coverage.out

install:
	go install $(LDFLAGS) .

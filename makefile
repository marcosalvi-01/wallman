BIN=wallman
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
EXT?=
PREFIX ?= /usr/local

LDFLAGS=-s -w -X wallman/cmd.Version=$(VERSION)
BUILDFLAGS=-trimpath

.DEFAULT_GOAL := build

build: generate
	CGO_ENABLED=0 go build $(BUILDFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) .

# Build for specific platform (used by CI)
build-release:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILDFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN)-$(GOOS)-$(GOARCH)$(EXT) .

clean:
	rm -f $(BIN) $(BIN)-* || true
	rm log.log || true

generate:
	go generate ./...

fmt:
	go tool gofumpt -w .

lint: fmt
	golangci-lint run ./...

vet:
	go vet ./...
	go tool sqlc vet ./...

test: vet
	go test ./...

.PHONY: build build-release clean generate fmt lint vet test


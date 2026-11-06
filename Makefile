BINARY_NAME=go-git-swap
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

.PHONY: build build-all test clean install lint vuln coverage

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/go-git-swap

build-all:
	$(foreach p,$(PLATFORMS),GOOS=$(word 1,$(subst /, ,$(p))) GOARCH=$(word 2,$(subst /, ,$(p))) CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-$(subst /,-,$(p)) ./cmd/go-git-swap;)

test:
	go test -race ./...

clean:
	rm -rf bin/ coverage.out

install:
	go install $(LDFLAGS) ./cmd/go-git-swap

lint:
	golangci-lint run

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.DEFAULT_GOAL := build

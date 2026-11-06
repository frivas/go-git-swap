BINARY_NAME=go-git-swap
VERSION=$(shell git describe --tags --always --dirty)
# -trimpath keeps local paths (e.g. /Users/<you>) out of the binaries.
LDFLAGS=-trimpath -ldflags "-s -w -X main.Version=$(VERSION)"
PLATFORMS=linux/amd64 linux/arm64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: build build-all checksums test clean install lint vuln coverage

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/go-git-swap

build-all:
	$(foreach p,$(PLATFORMS),GOOS=$(word 1,$(subst /, ,$(p))) GOARCH=$(word 2,$(subst /, ,$(p))) CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-$(subst /,-,$(p))$(if $(findstring windows,$(p)),.exe) ./cmd/go-git-swap;)

checksums: build-all
	cd bin && shasum -a 256 $(BINARY_NAME)-* > SHA256SUMS

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

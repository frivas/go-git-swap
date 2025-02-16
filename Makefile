BINARY_NAME=go-git-swap
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.PHONY: build build-all test clean install lint coverage

# Build for the current platform
build:
    go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/go-git-swap

# Build for all supported platforms
build-all:
    GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd/go-git-swap
    GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 ./cmd/go-git-swap
    GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/go-git-swap
    GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd/go-git-swap

# Run tests
test:
    go test -v ./...

# Clean build artifacts
clean:
    rm -rf bin/

# Install locally
install: build
    mv bin/${BINARY_NAME} ${GOPATH}/bin/${BINARY_NAME}

# Run linter
lint:
    golangci-lint run

# Generate test coverage
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

.DEFAULT_GOAL := build
```

3. Let's create a LICENSE file (MIT License):

<augment_code_snippet path="LICENSE" mode="EDIT">
```text
MIT License

Copyright (c) 2024 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

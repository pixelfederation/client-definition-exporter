 .PHONY: build install clean test integration modverify modtidy

VERSION=`egrep -o '[0-9]+\.[0-9a-z.\-]+' version.go`
GIT_SHA=`git rev-parse --short HEAD || echo`

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

build:
	@echo "Building definition_exporter..."
	@mkdir -p bin
	@go build -ldflags "-X cmd.GitSHA=${GIT_SHA}" -o bin/definition_exporter .

install:
	@echo "Installing definition_exporter..."
	@install -c bin/definition_exporter /usr/local/bin/definition_exporter

clean:
	@rm -f bin/*

test:
	@echo "Running tests..."
	#@go test `go list ./... | grep -v vendor/`

fmt:
	@gofmt -l -w -s $(SRC)

modtidy:
	@go mod tidy

modverify:
	@go mod verify

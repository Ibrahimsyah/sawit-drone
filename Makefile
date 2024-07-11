.PHONY: clean all init

all: build/main

build/main:
	@echo "Building..."
	go build -o $@ *.go

init: clean
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

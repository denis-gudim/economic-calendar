.PHONY: all default lint build test

all: lint build test

default: lint test

lint:
	golangci-lint run ./...

build: build-api build-loader

build-api:
	go get -d -v ./...
	go build -o build/bin/api ./cmd/api/...

build-loader:
	go get -d -v ./...
	go build -o build/bin/loader ./cmd/api/...

test:
	go test -v -cover ./...
# Makefile 'go'
# ref: https://golangdocs.com/makefiles-golang

BINARY_NAME=todoClient

.DEFAULT_GOAL := run

OS_NAME := $(shell uname -s | tr A-Z a-z)

.PHONY: all
all: build test

.PHONY: run
run: build
	./bin/${BINARY_NAME}-${OS_NAME}

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-${OS_NAME} *.go


.PHONY: clean
clean:
	go clean
	rm bin/${BINARY_NAME}-${OS_NAME}

.PHONY: vet
vet:
	go vet

.PHONY: test
test:
	go test -v ./...
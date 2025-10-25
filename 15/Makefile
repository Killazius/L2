.PHONY: test lint vet build clear
MAIN_FILE ?= main.go
BIN_NAME ?= shell
BUILD_DIR ?= ./bin
GO_BUILD_FLAGS ?=
GO_TEST_FLAGS ?= -v -race -parallel 5 -shuffle=on -coverprofile=./cover.out -covermode=atomic


all: build

build:
	mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BIN_NAME) cmd/$(MAIN_FILE)

test:
	go test $(GO_TEST_FLAGS) ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clear:
	rm -rf $(BUILD_DIR) ./cover.out

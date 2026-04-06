.PHONY: build run test lint clean

APP_NAME=firefly-iii-go
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

.PHONY: build build-frontend build-backend run dev test lint clean

APP_NAME=quillow
BUILD_DIR=bin

build: build-frontend build-backend

build-frontend:
	cd web && npm install && npm run build
	rm -rf cmd/server/frontend
	cp -r web/dist cmd/server/frontend

build-backend:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

dev:
	cd web && npm run dev &
	go run ./cmd/server

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html web/dist web/node_modules

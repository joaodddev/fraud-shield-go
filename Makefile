.PHONY: run build tidy lint

APP_NAME=fraud-shield-go
BUILD_DIR=bin

run:
	@go run ./cmd/server/main.go

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/main.go

tidy:
	@go mod tidy

lint:
	@golangci-lint run ./...
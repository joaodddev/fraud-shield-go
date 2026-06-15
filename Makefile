.PHONY: run build tidy lint infra infra-down logs-kafka

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

infra:
	@docker-compose up -d

infra-down:
	@docker-compose down

logs-kafka:
	@docker logs -f fraud-shield-kafka
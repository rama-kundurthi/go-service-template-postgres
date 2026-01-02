.PHONY: run test gen up down build

BINARY_NAME ?= server
BIN_DIR ?= bin

run:
	@PORT=8080 go run ./cmd/server

test:
	@go test ./...

gen:
	@sqlc generate

build:
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

up:
	@docker compose up --build

down:
	@docker compose down -v

.PHONY: build run watch test clean fmt lint pre-commit

BINARY_NAME=scheduler
OUTPUT_DIR=bin
MAIN_FILE=./cmd/app/main.go

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

dev:
	air

watch:
	docker compose --env-file ./.env.development -f ./docker-compose.development.yml -p scheduler-test up --watch

run:
	make build
	./$(OUTPUT_DIR)/scheduler

test:
	go test ./test/...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(OUTPUT_DIR)

pre-commit:
	pre-commit install --hook-type commit-msg

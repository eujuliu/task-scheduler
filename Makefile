.PHONY: build run watch test coverage clean fmt lint check pre-commit

BINARY_NAME=scheduler
OUTPUT_DIR=bin
MAIN_FILE=./cmd/app/main.go
COVERAGE_DIR=coverage

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
	go test ./... -count=1

coverage:
	make clean
	mkdir $(COVERAGE_DIR)
	go test -coverprofile=coverage/cover.out ./...
	go tool cover -html=coverage/cover.out -o coverage/cover.html

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...

check:
	make fmt
	make lint

clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(COVERAGE_DIR)

pre-commit:
	pre-commit install --hook-type commit-msg

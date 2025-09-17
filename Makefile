.PHONY: build run watch watch_down test test_with_cache coverage clean fmt lint check pre-commit test_stripe

BINARY_NAME=scheduler
OUTPUT_DIR=bin
MAIN_FILE=./cmd/app/main.go
COVERAGE_DIR=coverage

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

dev:
	air

watch:
	cat .env* > .env.docker; \
	docker compose --env-file ./.env.docker -f ./docker-compose.development.yml -p scheduler-test --profile scheduler_dev watch --prune; \
	make watch_down

watch_down:
	docker stop $$(docker ps -a -q --filter "label=com.docker.compose.project=scheduler-test"); \
	docker system prune -a --filter "label=com.docker.compose.project=scheduler-test" -f; \
	rm -f .env.docker

run:
	make build
	./$(OUTPUT_DIR)/scheduler

test:
	go test -tags=unit -count=1 -short -v ./...

test_with_cache:
	go test -tags=unit -short -v ./...

coverage:
	make clean; \
	mkdir $(COVERAGE_DIR); \
	go test -coverprofile=coverage/cover.out ./...; \
	go tool cover -html=coverage/cover.out -o coverage/cover.html

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...

check:
	make fmt
	make lint

clean:
	rm -rf $(OUTPUT_DIR); \
	rm -rf $(COVERAGE_DIR)

pre-commit:
	pre-commit install

test_stripe:
	# npm i -g live-server
	live-server --port=5500 --host="localhost" --watch=./pkg/stripe --entry-file=./pkg/stripe/index.html

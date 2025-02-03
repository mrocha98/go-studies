all: build test

build:
	@echo "Building..."
	@go build -o ./bin/api ./cmd/api

run:
	@go run cmd/api/main.go

docker-run-dev:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f docker-compose.dev.yml up -d --build; \
	fi

docker-down-dev:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f docker-compose.dev.yml down; \
	fi

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v


clean:
	@echo "Cleaning..."
	@rm -f main

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run-dev docker-down-dev itest

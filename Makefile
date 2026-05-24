.PHONY: help dev build run test clean swagger-gen admin-cli

help:
	@echo "Available commands:"
	@echo "  make dev          - Run development server with live reload"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the built application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make swagger-gen  - Generate Swagger documentation"
	@echo "  make admin-cli    - Promote a user to admin (EMAIL=... or PHONE=...)"

dev:
	@echo "Starting development server with live reload..."
	go run github.com/air-verse/air@latest

build:
	@echo "Building application..."
	go build -o tmp/main ./cmd/api/main.go

run: build
	@echo "Running application..."
	./tmp/main

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf tmp/
	go clean

swagger-gen:
	@echo "Generating Swagger documentation..."
	go run github.com/swaggo/swag/cmd/swag@latest init --generalInfo cmd/api/main.go --output docs

admin-cli:
	@echo "Promoting a user to admin..."
	@if [ -n "$(EMAIL)" ]; then \
		go run ./cmd/admin --email "$(EMAIL)"; \
	elif [ -n "$(PHONE)" ]; then \
		go run ./cmd/admin --phone "$(PHONE)"; \
	elif [ -n "$(ARGS)" ]; then \
		go run ./cmd/admin $(ARGS); \
	else \
		echo "Provide EMAIL=... or PHONE=... (or ARGS=\"--phone ...\")"; \
		exit 1; \
	fi

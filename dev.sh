#!/bin/bash

# Development server startup script with hot reload

echo "Starting development server with hot reload..."
go run github.com/swaggo/swag/cmd/swag@latest init --generalInfo cmd/api/main.go --output docs
go run ./cmd/migrate
go run github.com/air-verse/air@latest


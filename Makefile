.PHONY: test coverage coverage-html build run help

test:
	go test -v ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

help:
	@echo "Available commands:"
	@echo "  make test          - Run unit tests"
	@echo "  make coverage      - Run tests with coverage report"
	@echo "  make coverage-html - Run tests and open coverage in browser"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"

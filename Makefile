# Variables
TF_DIR := terraform/environments/dev

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

localstack-up:
	docker compose -f docker-compose.localStack.yml up -d

localstack-down:
	docker compose -f docker-compose.localStack.yml down

# Infrastructure Commands
infra-init:
	cd $(TF_DIR) && terraform init

infra-up:
	cd $(TF_DIR) && terraform apply -auto-approve

infra-down:
	cd $(TF_DIR) && terraform destroy -auto-approve

# The "Full Setup" Command
# Starts docker, waits 2 seconds for LocalStack to breathe, then applies TF
up:
	docker compose -f docker-compose.localStack.yml up -d
	@echo "Waiting for LocalStack to be ready..."
	@sleep 2
	cd $(TF_DIR) && terraform init && terraform apply -auto-approve

down:
	cd $(TF_DIR) && terraform destroy -auto-approve
	docker compose -f docker-compose.localStack.yml down

check-infra:
	@echo "--- DynamoDB Tables ---"
	@aws --endpoint-url=http://localhost:4566 dynamodb list-tables
	@echo "\n--- Event Buses ---"
	@aws --endpoint-url=http://localhost:4566 events list-event-buses
	@echo "\n--- SQS Queues ---"
	@aws --endpoint-url=http://localhost:4566 sqs list-queues

help:
	@echo "Available commands:"
	@echo "  make test          - Run unit tests"
	@echo "  make coverage      - Run tests with coverage report"
	@echo "  make coverage-html - Run tests and open coverage in browser"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make localstack-up   - Start LocalStack services"
	@echo "  make localstack-down - Stop LocalStack services"

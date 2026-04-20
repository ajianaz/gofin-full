.PHONY: build run test lint docker-up docker-down migrate-up migrate-down test-integration-infra test-integration-teardown

COMPOSE_DIR := ../deployments/docker

# Build the server binary
build:
	go build -o bin/server ./cmd/server

# Run the server locally
run: build
	./bin/server

# Run all tests
test:
	go test ./... -v -count=1

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

# Run only unit tests
test-unit:
	go test ./tests/unit/... -v -count=1

# Run only integration tests (requires running containers)
test-integration:
	AUTH_PROVIDER=disabled \
	AUTH_JWT_SECRET=test-jwt-secret-for-integration-tests-32ch \
	DB_HOST=localhost DB_PORT=5433 DB_DATABASE=gofin_test DB_USERNAME=gofin_test DB_PASSWORD=gofin_test \
	REDIS_HOST=localhost REDIS_PORT=6380 \
	go test ./tests/integration/... -v -count=1

# Start integration test infrastructure (PostgreSQL + Redis on non-default ports)
test-integration-infra:
	docker compose -f $(COMPOSE_DIR)/docker-compose.test.yml up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3

# Stop integration test infrastructure
test-integration-teardown:
	docker compose -f $(COMPOSE_DIR)/docker-compose.test.yml down -v

# Lint
lint:
	golangci-lint run ./...

# Docker Compose
docker-up:
	docker compose -f $(COMPOSE_DIR)/docker-compose.yml up -d

docker-down:
	docker compose -f $(COMPOSE_DIR)/docker-compose.yml down

docker-logs:
	docker compose -f $(COMPOSE_DIR)/docker-compose.yml logs -f

# Database migrations
migrate-up:
	goose -dir migrations/postgres postgres "$(DB_DSN)" up

migrate-down:
	goose -dir migrations/postgres postgres "$(DB_DSN)" down

# Tidy dependencies
tidy:
	go mod tidy

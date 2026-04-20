.PHONY: api-build api-run api-test api-test-unit api-test-integration api-lint api-tidy \
        web-dev web-build web-install \
        docker-up docker-down docker-selfhost docker-selfhost-down docker-test \
        help

COMPOSE_DIR := deployments/docker

# ─── API (Go backend) ───────────────────────────────────────

api-build:
	cd api && $(MAKE) build

api-run:
	cd api && $(MAKE) run

api-test:
	cd api && $(MAKE) test

api-test-unit:
	cd api && $(MAKE) test-unit

api-test-integration:
	cd api && $(MAKE) test-integration

api-lint:
	cd api && $(MAKE) lint

api-tidy:
	cd api && $(MAKE) tidy

# ─── Web (SvelteKit frontend) ───────────────────────────────

web-install:
	cd web && bun install

web-dev:
	cd web && bun run dev

web-build:
	cd web && bun run build

web-lint:
	cd web && bun run lint

# ─── Docker ─────────────────────────────────────────────────

docker-up:
	docker compose -f $(COMPOSE_DIR)/docker-compose.yml up -d

docker-down:
	docker compose -f $(COMPOSE_DIR)/docker-compose.yml down

docker-selfhost:
	docker compose -f $(COMPOSE_DIR)/docker-compose.selfhost.yml up -d

docker-selfhost-down:
	docker compose -f $(COMPOSE_DIR)/docker-compose.selfhost.yml down

docker-test:
	docker compose -f $(COMPOSE_DIR)/docker-compose.test.yml up --build --abort-on-container-exit

# ─── Help ───────────────────────────────────────────────────

help:
	@echo "Gofin Full — Monorepo Makefile"
	@echo ""
	@echo "API (backend):"
	@echo "  api-build            Build Go binary"
	@echo "  api-run              Build and run server locally"
	@echo "  api-test             Run all tests"
	@echo "  api-test-unit        Run unit tests only"
	@echo "  api-test-integration Run integration tests (needs infra)"
	@echo "  api-lint             Lint Go code"
	@echo "  api-tidy             Tidy Go modules"
	@echo ""
	@echo "Web (frontend):"
	@echo "  web-install          Install Node dependencies"
	@echo "  web-dev              Start SvelteKit dev server"
	@echo "  web-build            Build for production"
	@echo "  web-lint             Lint Svelte/TS code"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up            Start dev stack (Postgres + Redis + Keycloak)"
	@echo "  docker-down          Stop dev stack"
	@echo "  docker-selfhost      Start self-hosted stack (Caddy + API + Web)"
	@echo "  docker-selfhost-down Stop self-hosted stack"
	@echo "  docker-test          Run full test suite in containers"

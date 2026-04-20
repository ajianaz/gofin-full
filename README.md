# Gofin Full

Personal finance tracker monorepo — Go backend + SvelteKit frontend.

## Structure

```
gofin-full/
├── api/                # Go backend (Fiber)
├── web/                # SvelteKit frontend
├── mobile/             # Future Flutter app (placeholder)
├── deployments/docker/ # Docker configs
├── docs/               # OpenAPI spec, runbook, research
├── scripts/            # Utility scripts
├── Makefile            # Root orchestrator
└── .env.example        # Environment variables
```

## Quick Start

### Backend

```bash
cd api
cp .env.example .env
go build -o bin/server ./cmd/server
./bin/server
```

### Frontend

```bash
cd web
npm install
npm run dev
```

### Docker (self-hosted)

```bash
cp .env.example .env
# Edit .env — set AUTH_JWT_SECRET, DB_PASSWORD, DOMAIN
docker compose -f deployments/docker/docker-compose.selfhost.yml up -d
```

## Makefile Targets

```bash
make help            # List all targets
make api-build       # Build backend
make web-dev         # Start frontend dev server
make docker-selfhost # Start self-hosted stack
```

## API

Runs on `:8080`. See `docs/openapi.yaml` for full API spec.

- Health: `GET /health`
- API docs: `GET /api/v1/docs`
- OpenAPI: `GET /api/v1/openapi.json`

## License

Private repository.

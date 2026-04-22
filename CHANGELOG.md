# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.1.0] - 2026-04-22

### Added
- Go API backend with Fiber v2 — 125+ endpoints covering personal finance domain
- SvelteKit 5 frontend with Tailwind CSS v4, shadcn-svelte, dark mode, i18n (id/en)
- Double-entry bookkeeping with split transactions
- Hierarchical RBAC with 21 group-level roles and 3 wallet-level roles
- UUID v7 primary keys for all entities
- Real-time notifications via Server-Sent Events (SSE)
- JWT authentication with refresh tokens
- OAuth2 login (Google, GitHub)
- Optional Keycloak (OIDC) integration
- Budget tracking, bill management, piggy banks, recurring transactions
- Rules engine, currency management, exchange rates
- CSV/OFX export and reconciliation
- Analytics: spending by category, period, net worth
- Audit trail and webhook support
- Prometheus metrics endpoint
- Self-hosted Docker Compose deployment with Caddy (HTTPS), PostgreSQL 17, Redis 7
- Automated daily database backup with retention
- Auto-seed admin user on first startup
- E2E tests (Playwright) and unit/integration tests

### Documentation
- Comprehensive README (monorepo overview, self-host quick start)
- API and Web READMEs
- Architecture overview, production runbook
- OpenAPI 3.0 specification
- GitHub issue and PR templates

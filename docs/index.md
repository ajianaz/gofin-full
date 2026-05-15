---
layout: home

hero:
  name: "Gofin"
  text: "Self-Hosted Finance Tracker"
  tagline: "Track wallets, transactions, budgets, and bills. Double-entry bookkeeping with hierarchical RBAC. Open source, Docker-ready."
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: GitHub
      link: https://github.com/ajianaz/gofin-full

features:
  - title: 🏦 Double-Entry Bookkeeping
    details: Every transaction has a source and destination wallet. Full auditability with split transactions, tags, attachments, and notes.
  - title: 🛡️ Hierarchical RBAC
    details: 21 group-level roles + 3 wallet-level roles (owner/editor/viewer). Fine-grained permissions from read-only to full owner access.
  - title: 📊 Analytics & Reports
    details: Spending by category and period, net worth tracking, budget analysis. Visual charts powered by Chart.js.
  - title: 🔄 Real-time Notifications
    details: Server-Sent Events (SSE) for instant updates. No polling needed — live notification panel in the UI.
  - title: 🌍 i18n & Dark Mode
    details: Indonesian and English. Dark mode with system preference detection. Responsive design for all screen sizes.
  - title: 🐳 Docker Self-Host
    details: Single `docker compose` command. Caddy auto-provisions HTTPS. Automated daily database backups included.
---

## Quick Start

```bash
git clone https://github.com/ajianaz/gofin-full.git
cd gofin-full
cp .env.example .env
# Edit .env — set AUTH_JWT_SECRET, DOMAIN, ADMIN_EMAIL
make docker-selfhost
```

Open `https://your-domain` — that's it. Caddy handles HTTPS automatically.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Fiber v2, PostgreSQL 17, Redis 7 |
| Frontend | SvelteKit 5, Svelte 5, Tailwind CSS 4, shadcn-svelte |
| Auth | JWT (golang-jwt/v5), OAuth2 (Google, GitHub), Keycloak OIDC |
| Deployment | Docker Compose, Caddy (HTTPS) |

## Stats

- **135+ API endpoints** across 29 resource groups
- **43 web pages** with full CRUD
- **21 group-level roles** + 3 wallet-level roles
- **CSV/OFX export**, rules engine, recurring transactions, bill tracking

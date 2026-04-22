# Gofin Web

SvelteKit 5 frontend for [Gofin](../README.md) — a self-hosted personal finance tracker.

## Features

- **40+ pages** covering the full personal finance domain
- **i18n** — Indonesian (default) + English
- **Dark mode** with system preference detection
- **Responsive** — mobile-first layout
- **Accessible** — keyboard navigation, ARIA labels

## Tech Stack

| Technology | Purpose |
|-----------|---------|
| SvelteKit 5 | Framework (file-based routing, SSR) |
| Svelte 5 | UI with runes ($state, $derived, $effect) |
| Tailwind CSS 4 | Utility-first styling |
| shadcn-svelte | UI component library |
| Lucide Svelte | Icon library |
| TanStack Table | Headless table for data grids |
| Playwright | E2E testing |

## Development

```bash
# Install dependencies
bun install

# Start dev server (proxies /api to localhost:8080)
bun run dev

# Type check
bun run check
```

Open http://localhost:5173. The API must be running on port 8080 (use `make docker-dev` from the monorepo root).

## Project Structure

```
web/
├── src/
│   ├── routes/
│   │   ├── (auth)/         # login, register, forgot-password, reset-password
│   │   └── (app)/          # dashboard, wallets, transactions, budgets, settings, etc.
│   ├── lib/
│   │   ├── components/     # UI components (shadcn-svelte based)
│   │   ├── services/       # API client functions
│   │   ├── stores/         # Svelte 5 stores (auth, theme, i18n)
│   │   ├── i18n/           # Translation files (id/, en/)
│   │   └── types/          # TypeScript type definitions
│   └── app.html            # HTML template
├── static/                 # Static assets
├── tests/e2e/              # Playwright E2E tests
├── playwright.config.ts
├── svelte.config.js
├── tailwind.config.ts
└── vite.config.ts
```

## Route Groups

| Group | Path | Auth | Pages |
|-------|------|------|-------|
| `(auth)` | `/login`, `/register`, `/forgot-password`, `/reset-password` | Public | Authentication flows |
| `(app)` | `/dashboard`, `/wallets`, `/transactions`, etc. | Protected | 36 authenticated pages |

## i18n

Translations live in `src/lib/i18n/`. Default language is Indonesian (`id`), with English (`en`) support. The language store detects browser preference and falls back to Indonesian.

## Build

```bash
bun run build
```

Output goes to `.svelte-kit/build/` using the Node adapter.

## E2E Tests

```bash
# Run all E2E tests (requires running API)
bunx playwright test

# Run with UI
bunx playwright test --ui

# Run specific file
bunx playwright test tests/e2e/auth.spec.ts
```

Tests use the API proxy (Vite) to register test users via `/api/v1/auth/register` before navigating protected pages.

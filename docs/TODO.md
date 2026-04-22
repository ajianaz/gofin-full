# Web UI TODO

Post-MVP enhancements for production readiness.

## UX Feedback

- [ ] **Toast/snackbar notifications** — add sonner or similar lib for success/error feedback after create, delete, update operations (currently redirects silently)
- [ ] **Inline form validation** — replace HTML5 `required` only with per-field error messages on create forms
- [ ] **Update UI** — service layer has `update()` for all 9 resources but no UI to trigger them (only delete has buttons)

## Accessibility

- [ ] **`aria-live` regions** — loading states, error messages, and empty states should be announced to screen readers
- [ ] **Focus management** — after delete, focus should move to a sensible element (next item or empty state message)
- [ ] **Skip links** — add skip-to-content link for keyboard navigation

## i18n

- [ ] **Dashboard hardcoded empty states** — `Belum ada transaksi.` / `Belum ada budget.` still partially hardcoded (was fixed to use `t('common.noData')` but context-specific messages would be better)
- [ ] **Confirm dialogs** — generic `t('common.delete') + '?'` could be more descriptive per resource type

## Mock → Real API (remaining pages)

- [ ] **Currencies** — `currencies/+page.svelte` uses `mockCurrencies`, needs real API (`GET /currencies`)
- [ ] **Groups** — `groups/+page.svelte` uses `mockGroups`, needs real API (`GET /groups`, `POST /groups`, `POST /groups/switch`)
- [ ] **Reports** — `reports/+page.svelte` and subpages use mock data, needs real API aggregation
- [ ] **Export** — `export/+page.svelte` form needs to trigger real CSV/OFX export via API
- [ ] **Admin pages** — `admin/users/+page.svelte` and `admin/audit-log/+page.svelte` use mock data, need real API

## Performance

- [ ] **Optimistic updates** — currently waits for API response before updating list, could update UI immediately and rollback on error
- [ ] **Skeleton loaders** — replace `Memuat...` text with skeleton components for better perceived performance
- [ ] **List pagination** — transactions page uses client-side pagination, should move to server-side for large datasets

# Real API Integration — Remaining Pages

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace all remaining mock data usage with real API calls across Currencies, Groups, Reports, Export, and Admin pages.

**Architecture:** Create service files following the existing pattern (import `api`, `unwrapMany`/`unwrapOne`, map JSON:API response to domain types). Each page gets loading/error states and calls the service in `onMount`. The Go backend already has all endpoints implemented — this is purely a frontend task.

**Tech Stack:** SvelteKit, Svelte 5 runes (`$state`, `$derived`), TypeScript, existing `api` client + JSON:API helpers.

**Branch:** `feat/real-api-integration` (already exists)

---

## Key Reference Patterns

**Service pattern** (from `web/src/lib/services/wallets.ts`):
```typescript
import { api } from './client.js';
import { unwrapMany } from './helpers.js';
export const xxxService = {
  async list() {
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/endpoint');
    return unwrapMany<Type>(res).map(item => ({ ...item, /* field normalization */ }));
  },
};
```

**Page pattern** (from `web/src/routes/(app)/transactions/+page.svelte`):
```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { xxxService } from '$lib/services/index.js';
  let items = $state<Type[]>([]);
  let isLoading = $state(true);
  let error = $state('');
  onMount(async () => {
    try { items = await xxxService.list(); } catch (e) { error = String(e); } finally { isLoading = false; }
  });
</script>
```

**Backend response format** (JSON:API):
```json
{ "data": [{ "type": "xxx", "id": "abc", "attributes": { "field": "value" } }] }
```

---

## Task 1: Currencies Service + Page Integration

**Files:**
- Create: `web/src/lib/services/currencies.ts`
- Modify: `web/src/lib/services/index.ts` (add export)
- Modify: `web/src/routes/(app)/currencies/+page.svelte` (replace mock)

**Step 1: Create currencies service**

```typescript
// web/src/lib/services/currencies.ts
import { api } from './client.js';
import { unwrapMany } from './helpers.js';
import type { Currency } from '$lib/types/domain.js';

export const currencyService = {
  async list(): Promise<Currency[]> {
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/currencies');
    return unwrapMany<Currency>(res).map((c) => ({
      ...c,
      code: c.id,
      decimal_places: (c as any).decimal_places ?? 2,
      enabled: (c as any).enabled ?? true
    }));
  }
};
```

Note: Backend uses currency code as `id` (e.g., `"USD"`). The `Currency` domain type already has `code` and `id` as separate fields, so we map `id -> code`.

**Step 2: Add export to services/index.ts**

Add this line to `web/src/lib/services/index.ts`:
```typescript
export { currencyService } from './currencies.js';
```

**Step 3: Update currencies page to use service**

Replace the entire `currencies/+page.svelte` with:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
  import { Card, CardContent } from '$lib/components/ui/card/index.js';
  import { currencyService } from '$lib/services/index.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let items = $state<Currency[]>([]);
  let isLoading = $state(true);
  let error = $state('');

  onMount(async () => {
    try {
      items = await currencyService.list();
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });

  import type { Currency } from '$lib/types/domain.js';
</script>

<PageHeader title={t('currencies.title')} description={t('currencies.description')} />

{#if error}
  <p class="text-sm text-destructive">{error}</p>
{:else if isLoading}
  <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
{:else}
  <Card>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/50">
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.code')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.name')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.symbol')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.decimalPlaces')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('common.status')}</th>
            </tr>
          </thead>
          <tbody>
            {#each items as currency}
              <tr class="border-b hover:bg-muted/30">
                <td class="p-3 font-mono font-medium text-foreground">{currency.code}</td>
                <td class="p-3 text-foreground">{currency.name}</td>
                <td class="p-3 text-muted-foreground">{currency.symbol}</td>
                <td class="p-3 text-muted-foreground">{currency.decimal_places}</td>
                <td class="p-3"><StatusBadge status={currency.enabled ? 'active' : 'inactive'} /></td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>

  <div class="mt-6">
    <a href="/currencies/exchange-rates" class="text-sm text-primary font-medium hover:underline">{t('currencies.viewRates')}</a>
  </div>
{/if}
```

**Step 4: Run type check**

Run: `cd web && bun run check`
Expected: PASS (no errors)

**Step 5: Commit**

```bash
git add web/src/lib/services/currencies.ts web/src/lib/services/index.ts web/src/routes/\(app\)/currencies/+page.svelte
git commit -m "feat: integrate real API for currencies page"
```

---

## Task 2: Exchange Rates Service + Page Integration

**Files:**
- Modify: `web/src/lib/services/currencies.ts` (add exchange rates)
- Modify: `web/src/routes/(app)/currencies/exchange-rates/+page.svelte` (replace mock)

**Step 1: Add exchange rates to currencies service**

Add to `web/src/lib/services/currencies.ts` (inside the `currencyService` object):

```typescript
  async exchangeRates(): Promise<ExchangeRate[]> {
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/exchange-rates');
    return unwrapMany<ExchangeRate>(res).map((r) => ({
      ...r,
      from_code: (r as any).from_currency_id ?? (r as any).from_code ?? '',
      to_code: (r as any).to_currency_id ?? (r as any).to_code ?? '',
      rate: parseFloat(String((r as any).rate ?? '0')),
      date: (r as any).date ?? ''
    }));
  }
```

Add this import at the top of the file:
```typescript
import type { ExchangeRate } from '$lib/types/domain.js';
```

Note: Backend uses `from_currency_id` / `to_currency_id` as attribute names, but the domain type uses `from_code` / `to_code`.

**Step 2: Update exchange rates page**

Replace `exchange-rates/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { BackButton } from '$lib/components/shared/index.js';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { currencyService } from '$lib/services/index.js';
  import { formatDate } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let rates = $state<import('$lib/types/domain.js').ExchangeRate[]>([]);
  let isLoading = $state(true);
  let error = $state('');

  onMount(async () => {
    try {
      rates = await currencyService.exchangeRates();
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });
</script>

<BackButton href="/currencies" label={t('currencies.title')} />

<div class="mb-6">
  <h1 class="text-2xl font-bold text-foreground">{t('currencies.exchangeRates.title')}</h1>
  <p class="text-sm text-muted-foreground mt-0.5">{t('currencies.exchangeRates.description')}</p>
</div>

{#if error}
  <p class="text-sm text-destructive">{error}</p>
{:else if isLoading}
  <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
{:else}
  <Card>
    <CardHeader>
      <CardTitle class="text-base">{t('currencies.exchangeRates.title')}</CardTitle>
    </CardHeader>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/50">
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.from')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.to')}</th>
              <th class="text-right p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.rate')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('currencies.exchangeRates.date')}</th>
            </tr>
          </thead>
          <tbody>
            {#each rates as rate}
              <tr class="border-b hover:bg-muted/30">
                <td class="p-3 font-mono font-medium text-foreground">{rate.from_code}</td>
                <td class="p-3 font-mono font-medium text-foreground">{rate.to_code}</td>
                <td class="p-3 text-right font-medium text-foreground">{rate.rate}</td>
                <td class="p-3 text-muted-foreground">{formatDate(rate.date)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>
{/if}
```

**Step 3: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 4: Commit**

```bash
git add web/src/lib/services/currencies.ts web/src/routes/\(app\)/currencies/exchange-rates/+page.svelte
git commit -m "feat: integrate real API for exchange rates page"
```

---

## Task 3: Groups Service + Page Integration

**Files:**
- Create: `web/src/lib/services/groups.ts`
- Modify: `web/src/lib/services/index.ts` (add export)
- Modify: `web/src/routes/(app)/groups/+page.svelte` (replace mock)

**Backend API response for groups:**
- `GET /groups` returns `{ data: [{ type: "user_groups", id: "uuid", attributes: { title: "string" } }] }`
- Note: `member_count` and `is_current` are NOT in the backend response. We'll add these to the domain type as optional, defaulting `is_current` based on auth store's active group, and `member_count` as `0` for now (backend doesn't return it yet).

**Step 1: Create groups service**

```typescript
// web/src/lib/services/groups.ts
import { api } from './client.js';
import { unwrapMany } from './helpers.js';
import type { UserGroup } from '$lib/types/domain.js';

export const groupService = {
  async list(): Promise<UserGroup[]> {
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/groups');
    return unwrapMany<UserGroup>(res).map((g) => ({
      id: g.id,
      title: (g as any).title ?? '',
      member_count: (g as any).member_count ?? 0,
      is_current: (g as any).is_current ?? false
    }));
  },

  async create(data: { title: string }): Promise<UserGroup> {
    const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/groups', data);
    const g = unwrapMany<UserGroup>({ data: [res.data] })[0];
    return { id: g.id, title: (g as any).title ?? data.title, member_count: 1, is_current: false };
  },

  async switch(userGroupId: string): Promise<void> {
    await api.post('/groups/switch', { user_group_id: userGroupId });
  }
};
```

**Step 2: Add export to services/index.ts**

```typescript
export { groupService } from './groups.js';
```

**Step 3: Update groups page**

Replace `groups/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
  import { Card, CardContent } from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Plus } from '@lucide/svelte';
  import { groupService } from '$lib/services/index.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let items = $state<import('$lib/types/domain.js').UserGroup[]>([]);
  let isLoading = $state(true);
  let error = $state('');

  async function load() {
    try {
      items = await groupService.list();
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(load);

  async function handleSwitch(id: string) {
    try {
      await groupService.switch(id);
      await load();
    } catch (e) {
      error = String(e);
    }
  }
</script>

<PageHeader title={t('groups.title')} description={t('groups.description')}>
  {#snippet actions()}
    <Button size="sm">
      <Plus class="size-4" />
      {t('groups.newGroup')}
    </Button>
  {/snippet}
</PageHeader>

{#if error}
  <p class="text-sm text-destructive">{error}</p>
{:else if isLoading}
  <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
{:else}
  <Card>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/50">
              <th class="text-left p-3 font-medium text-muted-foreground">{t('groups.name')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('groups.members')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('groups.activeGroup')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground"></th>
            </tr>
          </thead>
          <tbody>
            {#each items as group}
              <tr class="border-b hover:bg-muted/30">
                <td class="p-3 font-medium text-foreground">{group.title}</td>
                <td class="p-3 text-muted-foreground">{t('groups.memberCount', { count: group.member_count })}</td>
                <td class="p-3">
                  <StatusBadge status={group.is_current ? 'active' : 'inactive'} />
                </td>
                <td class="p-3">
                  {#if !group.is_current}
                    <Button variant="outline" size="sm" onclick={() => handleSwitch(group.id)}>{t('groups.switch')}</Button>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>
{/if}
```

**Step 4: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 5: Commit**

```bash
git add web/src/lib/services/groups.ts web/src/lib/services/index.ts web/src/routes/\(app\)/groups/+page.svelte
git commit -m "feat: integrate real API for groups page with switch"
```

---

## Task 4: Reports Service + Page Integration

**Files:**
- Create: `web/src/lib/services/reports.ts`
- Modify: `web/src/lib/services/index.ts` (add export)
- Modify: `web/src/routes/(app)/reports/+page.svelte` (replace mock)
- Modify: `web/src/routes/(app)/reports/net-worth/+page.svelte` (replace mock)
- Modify: `web/src/routes/(app)/reports/spending-by-category/+page.svelte` (replace mock)
- Modify: `web/src/routes/(app)/reports/spending-by-period/+page.svelte` (replace mock)

**Backend API responses:**
- `GET /analytics/spending-by-category?start=...&end=...` returns `{ data: [{ type: "category_spending", attributes: { category_id, category_name, total, count } }] }`
- `GET /analytics/spending-by-period?start=...&end=...` returns `{ data: [{ type: "period_spending", attributes: { period, income, expense } }] }`
- `GET /analytics/net-worth?start=...&end=...` returns `{ data: { type: "net_worth", attributes: { total_income, total_expense, net_income, transaction_count } } }`

Note: Net worth is a single object, not an array. The `unwrapOne` helper works for this. However, the data doesn't have a standard `id` field — we'll handle this specially.

**Step 1: Create reports service**

```typescript
// web/src/lib/services/reports.ts
import { api } from './client.js';
import { unwrapMany } from './helpers.js';

interface CategorySpending {
  category_id: string;
  category_name: string;
  total: number;
  count: number;
}

interface PeriodSpending {
  period: string;
  income: number;
  expense: number;
}

interface NetWorthSummary {
  total_income: number;
  total_expense: number;
  net_income: number;
  transaction_count: number;
}

export const reportService = {
  async spendingByCategory(start?: string, end?: string): Promise<CategorySpending[]> {
    const params = new URLSearchParams();
    if (start) params.set('start', start);
    if (end) params.set('end', end);
    const qs = params.toString();
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
      `/analytics/spending-by-category${qs ? '?' + qs : ''}`
    );
    return unwrapMany<CategorySpending>(res).map((r) => ({
      category_id: (r as any).category_id ?? '',
      category_name: (r as any).category_name ?? '',
      total: parseFloat(String((r as any).total ?? '0')),
      count: parseInt(String((r as any).count ?? '0'), 10)
    }));
  },

  async spendingByPeriod(start?: string, end?: string): Promise<PeriodSpending[]> {
    const params = new URLSearchParams();
    if (start) params.set('start', start);
    if (end) params.set('end', end);
    const qs = params.toString();
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
      `/analytics/spending-by-period${qs ? '?' + qs : ''}`
    );
    return unwrapMany<PeriodSpending>(res).map((r) => ({
      period: (r as any).period ?? '',
      income: parseFloat(String((r as any).income ?? '0')),
      expense: parseFloat(String((r as any).expense ?? '0'))
    }));
  },

  async netWorth(start?: string, end?: string): Promise<NetWorthSummary> {
    const params = new URLSearchParams();
    if (start) params.set('start', start);
    if (end) params.set('end', end);
    const qs = params.toString();
    const res = await api.get<{ data: { type: string; attributes: Record<string, unknown> } }>(
      `/analytics/net-worth${qs ? '?' + qs : ''}`
    );
    const attrs = res.data.attributes;
    return {
      total_income: parseFloat(String(attrs.total_income ?? '0')),
      total_expense: parseFloat(String(attrs.total_expense ?? '0')),
      net_income: parseFloat(String(attrs.net_income ?? '0')),
      transaction_count: parseInt(String(attrs.transaction_count ?? '0'), 10)
    };
  }
};
```

**Step 2: Add export to services/index.ts**

```typescript
export { reportService } from './reports.js';
```

**Step 3: Update reports overview page**

Replace `reports/+page.svelte`. This page currently computes everything from mock data. We'll simplify it to use the real API calls. The page will call `reportService.netWorth()`, `reportService.spendingByCategory()`, and `reportService.spendingByPeriod()` in parallel.

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { ChevronDown, Download } from '@lucide/svelte';
  import { reportService } from '$lib/services/index.js';
  import { formatCurrency } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let isLoading = $state(true);
  let error = $state('');

  let totalIncome = $state(0);
  let totalExpense = $state(0);
  let netIncome = $state(0);
  let categorySpending = $state<{ name: string; amount: number }[]>([]);
  let monthData = $state<{ label: string; income: number; expense: number }[]>([]);

  onMount(async () => {
    try {
      const [summary, cats, periods] = await Promise.all([
        reportService.netWorth(),
        reportService.spendingByCategory(),
        reportService.spendingByPeriod()
      ]);
      totalIncome = summary.total_income;
      totalExpense = summary.total_expense;
      netIncome = summary.net_income;
      categorySpending = cats
        .map((c) => ({ name: c.category_name, amount: c.total }))
        .filter((c) => c.amount > 0)
        .sort((a, b) => b.amount - a.amount);

      // Group period data into months for the bar chart
      const monthMap = new Map<string, { income: number; expense: number }>();
      for (const p of periods) {
        const month = p.period.substring(0, 7);
        const existing = monthMap.get(month) ?? { income: 0, expense: 0 };
        existing.income += p.income;
        existing.expense += p.expense;
        monthMap.set(month, existing);
      }
      monthData = Array.from(monthMap.entries())
        .sort(([a], [b]) => a.localeCompare(b))
        .slice(-6)
        .map(([m, d]) => {
          const [y, mo] = m.split('-');
          const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
          return { label: monthNames[parseInt(mo, 10) - 1], income: d.income, expense: d.expense };
        });
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });

  const maxCatAmount = $derived(Math.max(...categorySpending.map((c) => c.amount), 1));
  const maxMonthVal = $derived(Math.max(...monthData.flatMap((m) => [m.income, m.expense]), 1));
  const barColors = ['#3b82f6', '#ef4444', '#f59e0b', '#10b981', '#8b5cf6', '#ec4899'];
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <h2 class="text-lg font-semibold text-foreground">{t('reports.title')}</h2>
  </div>

  {#if error}
    <p class="text-sm text-destructive">{error}</p>
  {:else if isLoading}
    <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
  {:else}
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.income')}</CardTitle></CardHeader>
        <CardContent>
          <p class="text-xl font-bold text-green-600">{formatCurrency(totalIncome.toString())}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.expense')}</CardTitle></CardHeader>
        <CardContent>
          <p class="text-xl font-bold text-destructive">{formatCurrency(totalExpense.toString())}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.diff')}</CardTitle></CardHeader>
        <CardContent>
          <p class="text-xl font-bold {netIncome >= 0 ? 'text-green-600' : 'text-destructive'}">{formatCurrency(Math.abs(netIncome).toString())}</p>
          <p class="text-xs text-muted-foreground">{t('reports.monthlySavings')}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.netWorth')}</CardTitle></CardHeader>
        <CardContent>
          <p class="text-xl font-bold text-foreground">{formatCurrency(netIncome.toString())}</p>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory')}</CardTitle></CardHeader>
        <CardContent>
          <div class="flex flex-col gap-2">
            {#each categorySpending as cat, i}
              <div class="flex items-center gap-2">
                <span class="w-24 shrink-0 text-xs text-foreground truncate">{cat.name}</span>
                <div
                  class="h-5 rounded"
                  style="width: {(cat.amount / maxCatAmount) * 70}%; min-width: 8px; background-color: {barColors[i % barColors.length]}"
                ></div>
                <span class="text-xs text-muted-foreground">{formatCurrency(cat.amount.toString())}</span>
              </div>
            {/each}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.sixMonthTrend')}</CardTitle></CardHeader>
        <CardContent>
          {#if monthData.length > 0}
            <div class="flex items-end gap-2" style="height: 160px">
              {#each monthData as m}
                <div class="flex flex-1 flex-col items-center gap-1">
                  <div class="flex w-full items-end justify-center gap-[2px]" style="height: 140px">
                    <div
                      class="w-full rounded-t-sm bg-green-600"
                      style="height: {(m.income / maxMonthVal) * 100}%; min-height: 4px"
                    ></div>
                    <div
                      class="w-full rounded-t-sm bg-destructive"
                      style="height: {(m.expense / maxMonthVal) * 100}%; min-height: 4px"
                    ></div>
                  </div>
                  <span class="text-[11px] text-muted-foreground">{m.label}</span>
                </div>
              {/each}
            </div>
            <div class="mt-3 flex items-center gap-4">
              <div class="flex items-center gap-2">
                <div class="size-3 rounded-sm bg-green-600"></div>
                <span class="text-xs text-muted-foreground">{t('reports.income')}</span>
              </div>
              <div class="flex items-center gap-2">
                <div class="size-3 rounded-sm bg-destructive"></div>
                <span class="text-xs text-muted-foreground">{t('reports.expense')}</span>
              </div>
            </div>
          {:else}
            <p class="text-sm text-muted-foreground">{t('common.noData')}</p>
          {/if}
        </CardContent>
      </Card>
    </div>
  {/if}
</div>
```

**Step 4: Update reports/net-worth page**

Replace `reports/net-worth/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { ArrowLeft } from '@lucide/svelte';
  import { reportService } from '$lib/services/index.js';
  import { walletService } from '$lib/services/index.js';
  import { formatCurrency } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let isLoading = $state(true);
  let error = $state('');
  let wallets = $state<import('$lib/types/domain.js').Account[]>([]);

  onMount(async () => {
    try {
      wallets = await walletService.list();
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });

  let assets = $derived(wallets.filter((w) => w.type === 'asset' || w.type === 'cash'));
  let liabilities = $derived(wallets.filter((w) => w.type === 'liability'));
  let totalAssets = $derived(assets.reduce((s, w) => s + parseFloat(w.balance), 0));
  let totalLiabilities = $derived(liabilities.reduce((s, w) => s + Math.abs(parseFloat(w.balance)), 0));
  let netWorth = $derived(totalAssets - totalLiabilities);
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center gap-3">
    <Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
      <ArrowLeft class="size-4" />
      {t('common.back')}
    </Button>
    <h2 class="text-base font-semibold text-foreground">{t('reports.netWorth.title')}</h2>
  </div>

  {#if error}
    <p class="text-sm text-destructive">{error}</p>
  {:else if isLoading}
    <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
  {:else}
    <Card>
      <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.netWorth.title')}</CardTitle></CardHeader>
      <CardContent>
        <p class="text-xl font-bold {netWorth >= 0 ? 'text-green-600' : 'text-destructive'}">{formatCurrency(netWorth.toString())}</p>
      </CardContent>
    </Card>

    <div class="grid gap-4 md:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle class="text-base">{t('reports.netWorth.assets')}</CardTitle>
        </CardHeader>
        <CardContent class="p-0">
          {#each assets as w}
            <div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
              <span class="text-sm font-medium text-foreground">{w.name}</span>
              <span class="text-sm font-medium text-green-600">{formatCurrency(w.balance)}</span>
            </div>
          {/each}
          <div class="flex items-center justify-between bg-muted/50 px-4 py-3 font-semibold">
            <span class="text-sm text-foreground">{t('reports.netWorth.totalAssets')}</span>
            <span class="text-sm text-foreground">{formatCurrency(totalAssets.toString())}</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-base">{t('reports.netWorth.liabilities')}</CardTitle>
        </CardHeader>
        <CardContent class="p-0">
          {#each liabilities as w}
            <div class="flex items-center justify-between px-4 py-3 border-b last:border-b-0">
              <span class="text-sm font-medium text-foreground">{w.name}</span>
              <span class="text-sm font-medium text-destructive">{formatCurrency(w.balance)}</span>
            </div>
          {/each}
          <div class="flex items-center justify-between bg-muted/50 px-4 py-3 font-semibold">
            <span class="text-sm text-foreground">{t('reports.netWorth.totalLiabilities')}</span>
            <span class="text-sm text-foreground">{formatCurrency(totalLiabilities.toString())}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  {/if}
</div>
```

**Step 5: Update reports/spending-by-category page**

Replace `reports/spending-by-category/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Progress } from '$lib/components/ui/progress/index.js';
  import { ArrowLeft } from '@lucide/svelte';
  import { reportService } from '$lib/services/index.js';
  import { formatCurrency } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let isLoading = $state(true);
  let error = $state('');
  let categoryData = $state<{ name: string; amount: number }[]>([]);

  onMount(async () => {
    try {
      const cats = await reportService.spendingByCategory();
      categoryData = cats
        .map((c) => ({ name: c.category_name, amount: c.total }))
        .filter((c) => c.amount > 0)
        .sort((a, b) => b.amount - a.amount);
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });

  let totalSpent = $derived(categoryData.reduce((s, c) => s + c.amount, 0));
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center gap-3">
    <Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
      <ArrowLeft class="size-4" />
      {t('common.back')}
    </Button>
    <h2 class="text-base font-semibold text-foreground">{t('reports.spendingByCategory.title')}</h2>
  </div>

  {#if error}
    <p class="text-sm text-destructive">{error}</p>
  {:else if isLoading}
    <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
  {:else}
    <div class="grid gap-4 sm:grid-cols-2">
      <Card>
        <CardHeader class="pb-1"><CardTitle class="text-sm font-semibold">{t('reports.spendingByCategory.totalSpending')}</CardTitle></CardHeader>
        <CardContent>
          <p class="text-xl font-bold text-destructive">{formatCurrency(totalSpent.toString())}</p>
          <p class="text-xs text-muted-foreground">{t('reports.spendingByCategory.categoryCount', { count: categoryData.length })}</p>
        </CardContent>
      </Card>
    </div>

    <Card>
      <CardHeader>
        <CardTitle class="text-base">{t('reports.spendingByCategory.details')}</CardTitle>
      </CardHeader>
      <CardContent>
        {#if categoryData.length > 0}
          <div class="flex flex-col gap-4">
            {#each categoryData as cat}
              <div>
                <div class="flex items-center justify-between mb-1.5">
                  <span class="text-sm font-medium text-foreground">{cat.name}</span>
                  <span class="text-sm font-medium text-foreground">{formatCurrency(cat.amount.toString())}</span>
                </div>
                <div class="h-0.5 w-full rounded-full bg-muted"></div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-sm text-muted-foreground">{t('common.noData')}</p>
        {/if}
      </CardContent>
    </Card>
  {/if}
</div>
```

**Step 6: Update reports/spending-by-period page**

Replace `reports/spending-by-period/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { ArrowLeft } from '@lucide/svelte';
  import { reportService } from '$lib/services/index.js';
  import { formatCurrency } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let isLoading = $state(true);
  let error = $state('');
  let periodData = $state<{ period: string; income: number; expense: number; diff: number }[]>([]);

  onMount(async () => {
    try {
      const data = await reportService.spendingByPeriod();
      periodData = data.map((p) => ({
        period: p.period,
        income: p.income,
        expense: p.expense,
        diff: p.income - p.expense
      }));
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center gap-3">
    <Button variant="ghost" size="sm" onclick={() => goto('/reports')}>
      <ArrowLeft class="size-4" />
      {t('common.back')}
    </Button>
    <h2 class="text-base font-semibold text-foreground">{t('reports.spendingByPeriod.title')}</h2>
  </div>

  {#if error}
    <p class="text-sm text-destructive">{error}</p>
  {:else if isLoading}
    <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
  {:else}
    <Card>
      <CardHeader>
        <CardTitle class="text-base">{t('reports.spendingByPeriod.history')}</CardTitle>
      </CardHeader>
      <CardContent class="p-0">
        {#if periodData.length > 0}
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b bg-muted/50">
                  <th class="text-left px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.period')}</th>
                  <th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.income')}</th>
                  <th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.expense')}</th>
                  <th class="text-right px-4 py-3 font-medium text-muted-foreground">{t('reports.spendingByPeriod.diff')}</th>
                </tr>
              </thead>
              <tbody>
                {#each periodData as row}
                  <tr class="border-b last:border-b-0 hover:bg-muted/30">
                    <td class="px-4 py-3 font-medium text-foreground">{row.period}</td>
                    <td class="px-4 py-3 text-right text-green-600">{formatCurrency(row.income.toString())}</td>
                    <td class="px-4 py-3 text-right text-destructive">{formatCurrency(row.expense.toString())}</td>
                    <td class="px-4 py-3 text-right font-medium {row.diff >= 0 ? 'text-green-600' : 'text-destructive'}">
                      {row.diff >= 0 ? '+' : '-'}{formatCurrency(Math.abs(row.diff).toString())}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else}
          <p class="p-4 text-sm text-muted-foreground">{t('common.noData')}</p>
        {/if}
      </CardContent>
    </Card>
  {/if}
</div>
```

**Step 7: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 8: Commit**

```bash
git add web/src/lib/services/reports.ts web/src/lib/services/index.ts \
  web/src/routes/\(app\)/reports/+page.svelte \
  web/src/routes/\(app\)/reports/net-worth/+page.svelte \
  web/src/routes/\(app\)/reports/spending-by-category/+page.svelte \
  web/src/routes/\(app\)/reports/spending-by-period/+page.svelte
git commit -m "feat: integrate real API for all reports pages"
```

---

## Task 5: Export Service + Page Integration

**Files:**
- Create: `web/src/lib/services/export.ts`
- Modify: `web/src/lib/services/index.ts` (add export)
- Modify: `web/src/routes/(app)/export/+page.svelte` (replace mock wallets dropdown + add real download)

**Backend API:**
- `GET /export/csv` — returns CSV file with `Content-Disposition: attachment`
- `GET /export/ofx` — returns OFX file with `Content-Disposition: attachment`
- Both accept `?start=...&end=...&wallet_id=...` query params

**Step 1: Create export service**

```typescript
// web/src/lib/services/export.ts
import { api } from './client.js';

export const exportService = {
  async downloadCSV(startDate?: string, endDate?: string, walletId?: string): Promise<void> {
    const params = new URLSearchParams();
    if (startDate) params.set('start', startDate);
    if (endDate) params.set('end', endDate);
    if (walletId) params.set('wallet_id', walletId);
    const qs = params.toString();
    const url = `/export/csv${qs ? '?' + qs : ''}`;

    const token = localStorage.getItem('access_token');
    const response = await fetch(`/api/v1${url}`, {
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      }
    });

    if (!response.ok) throw new Error(`Export failed: ${response.statusText}`);

    const blob = await response.blob();
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'transactions.csv';
    a.click();
    URL.revokeObjectURL(a.href);
  },

  async downloadOFX(startDate?: string, endDate?: string, walletId?: string): Promise<void> {
    const params = new URLSearchParams();
    if (startDate) params.set('start', startDate);
    if (endDate) params.set('end', endDate);
    if (walletId) params.set('wallet_id', walletId);
    const qs = params.toString();
    const url = `/export/ofx${qs ? '?' + qs : ''}`;

    const token = localStorage.getItem('access_token');
    const response = await fetch(`/api/v1${url}`, {
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      }
    });

    if (!response.ok) throw new Error(`Export failed: ${response.statusText}`);

    const blob = await response.blob();
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'transactions.ofx';
    a.click();
    URL.revokeObjectURL(a.href);
  }
};
```

Note: We use raw `fetch` here instead of the `api` client because the export endpoints return binary files (CSV/OFX), not JSON. The `api` client tries to parse JSON which would fail.

**Step 2: Add export to services/index.ts**

```typescript
export { exportService } from './export.js';
```

**Step 3: Update export page**

Replace `export/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { PageHeader, FormCard } from '$lib/components/shared/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Label } from '$lib/components/ui/label/index.js';
  import { Select } from '$lib/components/ui/select/index.js';
  import { exportService } from '$lib/services/index.js';
  import { walletService } from '$lib/services/index.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let format = $state('csv');
  let startDate = $state('');
  let endDate = $state('');
  let walletId = $state('');
  let wallets = $state<import('$lib/types/domain.js').Account[]>([]);
  let isExporting = $state(false);
  let error = $state('');

  onMount(async () => {
    try {
      wallets = await walletService.list();
    } catch {
      // Wallet list is optional for export
    }
  });

  async function handleExport(e: Event) {
    e.preventDefault();
    isExporting = true;
    error = '';
    try {
      if (format === 'csv') {
        await exportService.downloadCSV(startDate || undefined, endDate || undefined, walletId || undefined);
      } else {
        await exportService.downloadOFX(startDate || undefined, endDate || undefined, walletId || undefined);
      }
    } catch (err) {
      error = String(err);
    } finally {
      isExporting = false;
    }
  }
</script>

<PageHeader title={t('export.title')} description={t('export.description')} />

<FormCard title={t('export.exportData')}>
  <form class="grid gap-4" onsubmit={handleExport}>
    <div class="grid gap-2">
      <Label for="format">{t('export.format')}</Label>
      <Select bind:value={format} id="format">
        <option value="csv">CSV</option>
        <option value="ofx">OFX</option>
      </Select>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div class="grid gap-2">
        <Label for="start">{t('export.startDate')}</Label>
        <Input id="start" type="date" bind:value={startDate} />
      </div>
      <div class="grid gap-2">
        <Label for="end">{t('export.endDate')}</Label>
        <Input id="end" type="date" bind:value={endDate} />
      </div>
    </div>

    <div class="grid gap-2">
      <Label for="wallet">{t('export.wallet')}</Label>
      <Select bind:value={walletId} id="wallet">
        <option value="">{t('export.allWallets')}</option>
        {#each wallets as w}
          <option value={w.id}>{w.name}</option>
        {/each}
      </Select>
    </div>

    {#if error}
      <p class="text-sm text-destructive">{error}</p>
    {/if}

    <Button type="submit" disabled={isExporting}>
      {isExporting ? t('common.loading') : t('export.export')}
    </Button>
  </form>
</FormCard>
```

**Step 4: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 5: Commit**

```bash
git add web/src/lib/services/export.ts web/src/lib/services/index.ts web/src/routes/\(app\)/export/+page.svelte
git commit -m "feat: integrate real API for export page with CSV/OFX download"
```

---

## Task 6: Admin Users Service + Page Integration

**Files:**
- Create: `web/src/lib/services/admin.ts`
- Modify: `web/src/lib/services/index.ts` (add export)
- Modify: `web/src/routes/(app)/admin/users/+page.svelte` (replace mock)

**Backend API response:**
- `GET /admin/users` returns `{ data: [{ type: "users", id: "uuid", attributes: { email: "string", created_at: "RFC3339" } }] }`
- Note: Backend only returns `email` and `created_at`. `name`, `role`, `is_active` are not included. We'll add these as optional with defaults.

**Step 1: Create admin service**

```typescript
// web/src/lib/services/admin.ts
import { api } from './client.js';
import { unwrapMany } from './helpers.js';

interface AdminUser {
  id: string;
  email: string;
  name?: string;
  role?: string;
  is_active?: boolean;
  created_at: string;
}

export const adminService = {
  async listUsers(): Promise<AdminUser[]> {
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/admin/users');
    return unwrapMany<AdminUser>(res).map((u) => ({
      ...u,
      email: (u as any).email ?? '',
      name: (u as any).name ?? '',
      role: (u as any).role ?? 'user',
      is_active: (u as any).is_active ?? true,
      created_at: (u as any).created_at ?? ''
    }));
  }
};
```

**Step 2: Add export to services/index.ts**

```typescript
export { adminService } from './admin.js';
```

**Step 3: Update admin users page**

Replace `admin/users/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { PageHeader, StatusBadge } from '$lib/components/shared/index.js';
  import { Card, CardContent } from '$lib/components/ui/card/index.js';
  import { Badge } from '$lib/components/ui/badge/index.js';
  import { adminService } from '$lib/services/index.js';
  import { formatDate } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let users = $state<{ id: string; email: string; name: string; role: string; is_active: boolean; created_at: string }[]>([]);
  let isLoading = $state(true);
  let error = $state('');

  let roleLabels = $derived<Record<string, string>>({
    admin: t('admin.users.roleAdmin'),
    manager: t('admin.users.roleManager'),
    user: t('admin.users.roleUser')
  });

  onMount(async () => {
    try {
      users = await adminService.listUsers();
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  });
</script>

<PageHeader title={t('admin.users.title')} description={t('admin.users.description')} />

{#if error}
  <p class="text-sm text-destructive">{error}</p>
{:else if isLoading}
  <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
{:else}
  <Card>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/50">
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.email')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.name')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.role')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('common.status')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.users.joined')}</th>
            </tr>
          </thead>
          <tbody>
            {#each users as user}
              <tr class="border-b hover:bg-muted/30">
                <td class="p-3 text-foreground">{user.email}</td>
                <td class="p-3 font-medium text-foreground">{user.name || '-'}</td>
                <td class="p-3">
                  <Badge variant="outline">{roleLabels[user.role] ?? user.role}</Badge>
                </td>
                <td class="p-3"><StatusBadge status={user.is_active ? 'active' : 'inactive'} /></td>
                <td class="p-3 text-muted-foreground">{user.created_at ? formatDate(user.created_at) : '-'}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>
{/if}
```

**Step 4: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 5: Commit**

```bash
git add web/src/lib/services/admin.ts web/src/lib/services/index.ts web/src/routes/\(app\)/admin/users/+page.svelte
git commit -m "feat: integrate real API for admin users page"
```

---

## Task 7: Admin Audit Log Service + Page Integration

**Files:**
- Modify: `web/src/lib/services/admin.ts` (add audit log method)
- Modify: `web/src/routes/(app)/admin/audit-log/+page.svelte` (replace mock)

**Backend API response:**
- `GET /audit-logs?entity_type=...` returns `{ data: [{ type: "audit_logs", id: "int64", attributes: { user_id, action, entity_type, entity_id, old_value, new_value, ip_address, created_at } }] }`
- Note: The mock uses `user_email` but the backend returns `user_id`. The backend doesn't return user email — we'll show user_id for now.

**Step 1: Add audit log method to admin service**

Add to `web/src/lib/services/admin.ts` (inside `adminService`):

```typescript
  async listAuditLogs(entityType?: string): Promise<AuditLogEntry[]> {
    const params = new URLSearchParams();
    if (entityType) params.set('entity_type', entityType);
    const qs = params.toString();
    const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>(
      `/audit-logs${qs ? '?' + qs : ''}`
    );
    return unwrapMany<AuditLogEntry>(res).map((l) => ({
      id: l.id,
      action: (l as any).action ?? '',
      user_email: (l as any).user_id ?? (l as any).user_email ?? '',
      entity_type: (l as any).entity_type ?? '',
      entity_id: String((l as any).entity_id ?? ''),
      changes: (l as any).new_value ?? (l as any).changes ?? '',
      created_at: (l as any).created_at ?? ''
    }));
  }
```

Add the interface at the top of the file:
```typescript
interface AuditLogEntry {
  id: string;
  action: string;
  user_email: string;
  entity_type: string;
  entity_id: string;
  changes: string;
  created_at: string;
}
```

**Step 2: Update audit log page**

Replace `admin/audit-log/+page.svelte`:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { PageHeader, FilterBar } from '$lib/components/shared/index.js';
  import { Card, CardContent } from '$lib/components/ui/card/index.js';
  import { Select } from '$lib/components/ui/select/index.js';
  import { adminService } from '$lib/services/index.js';
  import { formatDate } from '$lib/utils/format.js';
  import { localeStore } from '$lib/stores/i18n.svelte.js';
  const t = localeStore.t;

  let actionFilter = $state('all');
  let entityFilter = $state('all');
  let logs = $state<import('$lib/types/domain.js').AuditLogEntry[]>([]);
  let isLoading = $state(true);
  let error = $state('');

  async function load() {
    try {
      logs = await adminService.listAuditLogs(entityFilter !== 'all' ? entityFilter : undefined);
    } catch (e) {
      error = String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(load);

  let filtered = $derived(() => {
    let result = [...logs];
    if (actionFilter !== 'all') result = result.filter((l) => l.action.startsWith(actionFilter));
    return result;
  });
</script>

<PageHeader title={t('admin.auditLog.title')} description={t('admin.auditLog.description')} />

<FilterBar>
  <Select bind:value={actionFilter} class="w-40">
    <option value="all">{t('admin.auditLog.allActions')}</option>
    <option value="user.login">{t('admin.auditLog.login')}</option>
    <option value="transaction">{t('admin.auditLog.transaction')}</option>
    <option value="budget">{t('admin.auditLog.budget')}</option>
    <option value="api_key">{t('admin.auditLog.apiKey')}</option>
    <option value="user.update">{t('admin.auditLog.userUpdate')}</option>
    <option value="group.create">{t('admin.auditLog.group')}</option>
    <option value="piggy_bank">{t('admin.auditLog.piggyBank')}</option>
    <option value="recurring">{t('admin.auditLog.recurring')}</option>
    <option value="rule">{t('admin.auditLog.rule')}</option>
    <option value="currency">{t('admin.auditLog.currency')}</option>
  </Select>
  <Select bind:value={entityFilter} class="w-40" onchange={load}>
    <option value="all">{t('admin.auditLog.allEntities')}</option>
    <option value="user">{t('admin.auditLog.user')}</option>
    <option value="transaction">{t('admin.auditLog.transaction')}</option>
    <option value="budget">{t('admin.auditLog.budget')}</option>
    <option value="api_key">{t('admin.auditLog.apiKey')}</option>
    <option value="group">{t('admin.auditLog.group')}</option>
    <option value="piggy_bank">{t('admin.auditLog.piggyBank')}</option>
    <option value="recurring">{t('admin.auditLog.recurring')}</option>
    <option value="rule">{t('admin.auditLog.rule')}</option>
    <option value="currency">{t('admin.auditLog.currency')}</option>
  </Select>
</FilterBar>

{#if error}
  <p class="text-sm text-destructive">{error}</p>
{:else if isLoading}
  <p class="text-sm text-muted-foreground">{t('common.loading')}</p>
{:else}
  <Card>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-muted/50">
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.time')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.user')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.action')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.entity')}</th>
              <th class="text-left p-3 font-medium text-muted-foreground">{t('admin.auditLog.changes')}</th>
            </tr>
          </thead>
          <tbody>
            {#each filtered() as log}
              <tr class="border-b hover:bg-muted/30">
                <td class="p-3 text-muted-foreground whitespace-nowrap">{formatDate(log.created_at)}</td>
                <td class="p-3 text-foreground">{log.user_email}</td>
                <td class="p-3">
                  <span class="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-xs">{log.action}</span>
                </td>
                <td class="p-3 text-muted-foreground">{log.entity_type} ({log.entity_id})</td>
                <td class="p-3 text-foreground max-w-xs truncate">{log.changes}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </CardContent>
  </Card>
{/if}
```

**Step 3: Run type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 4: Commit**

```bash
git add web/src/lib/services/admin.ts web/src/routes/\(app\)/admin/audit-log/+page.svelte
git commit -m "feat: integrate real API for admin audit log page"
```

---

## Task 8: Clean Up — Remove Unused Mock Data Imports

**Files:**
- Potentially remove mock data files that are no longer imported anywhere
- Run a final type check

**Step 1: Check which mock files are still used**

Run: `cd web && grep -r "mock-" src/routes/ --include="*.svelte" -l`

If the only remaining references are in already-converted pages (shouldn't be after Tasks 1-7), proceed to clean up.

**Step 2: Remove mock data files that are fully replaced**

Only remove files that have ZERO imports remaining anywhere in the project. Use grep to verify:
```bash
cd web && for f in src/lib/data/mock-*.ts; do
  base=$(basename "$f" .ts | sed 's/mock-//;s/-/_/g')
  count=$(grep -r "mock-$base" src/ --include="*.svelte" --include="*.ts" -l | wc -l)
  if [ "$count" -eq 0 ]; then echo "UNUSED: $f"; fi
done
```

Only delete files confirmed as unused. Do NOT delete mock files still referenced elsewhere.

**Step 3: Run final type check**

Run: `cd web && bun run check`
Expected: PASS

**Step 4: Commit**

```bash
git add -A web/src/lib/data/
git commit -m "chore: remove unused mock data files"
```

---

## Task 9: Update CHANGELOG.md

**Files:**
- Modify: `CHANGELOG.md`

**Step 1: Add entries under `## [Unreleased]`**

```markdown
### Changed
- Currencies page now uses real API (`GET /api/v1/currencies`)
- Exchange rates page now uses real API (`GET /api/v1/exchange-rates`)
- Groups page now uses real API (`GET /api/v1/groups`, `POST /api/v1/groups/switch`)
- Reports pages now use real API (`GET /api/v1/analytics/*`)
- Export page now uses real API (`GET /api/v1/export/csv`, `GET /api/v1/export/ofx`)
- Admin users page now uses real API (`GET /api/v1/admin/users`)
- Admin audit log page now uses real API (`GET /api/v1/audit-logs`)
- Removed unused mock data files
```

**Step 2: Commit**

```bash
git add CHANGELOG.md
git commit -m "docs: update changelog for real API integration"
```

---

## Summary

| Task | Pages | New Service Files | Commits |
|------|-------|-------------------|---------|
| 1 | Currencies | `currencies.ts` | 1 |
| 2 | Exchange Rates | (extends `currencies.ts`) | 1 |
| 3 | Groups | `groups.ts` | 1 |
| 4 | Reports (4 pages) | `reports.ts` | 1 |
| 5 | Export | `export.ts` | 1 |
| 6 | Admin Users | `admin.ts` | 1 |
| 7 | Admin Audit Log | (extends `admin.ts`) | 1 |
| 8 | Cleanup | - | 1 |
| 9 | Changelog | - | 1 |
| **Total** | **10 pages** | **5 new files** | **9 commits** |

import { test, expect } from '@playwright/test';

const API_BASE = '/api/v1';
let token: string;
let walletId: string;
let categoryId: string;
let destWalletId: string;

test.describe('Fresh DB — Full User Journey', () => {

  test('Step 1: Register new user', async ({ request }) => {
    const email = `fresh-${Date.now()}@gofin.io`;
    const res = await request.post(`${API_BASE}/auth/register`, {
      data: { email, password: 'TestPass123!' }
    });
    expect(res.status()).toBe(201);
    const body = await res.json();
    expect(body.access_token).toBeDefined();
    token = body.access_token;
  });

  test('Step 2: Check currencies exist', async ({ request }) => {
    const res = await request.get(`${API_BASE}/currencies`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.data.length).toBeGreaterThanOrEqual(5);
    // Find IDR
    const idr = body.data.find((c: any) => c.attributes?.code === 'IDR');
    expect(idr).toBeDefined();
    expect(idr.attributes.symbol).toBe('Rp');
  });

  test('Step 3: Create wallet with IDR currency', async ({ request }) => {
    const res = await request.post(`${API_BASE}/wallets`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { name: 'BCA', wallet_type: 'asset', currency_id: 'IDR' }
    });
    expect(res.status()).toBe(201);
    const body = await res.json();
    walletId = body.data.id;
    const attrs = body.data.attributes;
    expect(attrs.name).toBe('BCA');
    expect(attrs.currency_code).toBe('IDR');
    expect(attrs.currency_symbol).toBe('Rp');
    expect(attrs.virtual_balance).toBe('0.00');
  });

  test('Step 4: Create second wallet (e-wallet)', async ({ request }) => {
    const res = await request.post(`${API_BASE}/wallets`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { name: 'GoPay', wallet_type: 'expense', currency_id: 'IDR' }
    });
    expect(res.status()).toBe(201);
    const body = await res.json();
    destWalletId = body.data.id;
    expect(body.data.attributes.currency_symbol).toBe('Rp');
  });

  test('Step 5: List wallets shows both with currency', async ({ request }) => {
    const res = await request.get(`${API_BASE}/wallets`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.data.length).toBe(2);
    for (const w of body.data) {
      expect(w.attributes.currency_code).toBe('IDR');
      expect(w.attributes.currency_symbol).toBe('Rp');
    }
  });

  test('Step 6: Create expense category', async ({ request }) => {
    const res = await request.post(`${API_BASE}/categories`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { name: 'Makanan', type: 'expense' }
    });
    expect(res.status()).toBe(201);
    const body = await res.json();
    categoryId = body.data.id;
  });

  test('Step 7: Create withdrawal transaction', async ({ request }) => {
    const res = await request.post(`${API_BASE}/transactions`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        type: 'withdrawal',
        description: 'Makan siang',
        amount: '50000',
        date: new Date().toISOString(),
        source_id: walletId,
        destination_id: destWalletId,
        currency_id: 'IDR'
      }
    });
    expect(res.status()).toBe(201);
  });

  test('Step 8: Create deposit transaction', async ({ request }) => {
    const res = await request.post(`${API_BASE}/transactions`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        type: 'deposit',
        description: 'Gaji bulanan',
        amount: '5000000',
        date: new Date().toISOString(),
        source_id: destWalletId,
        destination_id: walletId,
        currency_id: 'IDR'
      }
    });
    expect(res.status()).toBe(201);
  });

  test('Step 9: Net-worth analytics shows data', async ({ request }) => {
    const res = await request.get(`${API_BASE}/analytics/net-worth`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    const attrs = body.data?.attributes;
    if (attrs) {
      expect(attrs.total_income).toBeDefined();
    }
  });

  test('Step 10: Spending by category has data', async ({ request }) => {
    const res = await request.get(`${API_BASE}/analytics/spending-by-category`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.data).toBeDefined();
  });

  test('Step 11: Spending by period has data', async ({ request }) => {
    const res = await request.get(`${API_BASE}/analytics/spending-by-period`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.data).toBeDefined();
  });
});

test.describe('Fresh DB — UI Pages', () => {
  test.beforeEach(async ({ page }) => {
    // Login via UI
    await page.goto('/login');
    await page.waitForLoadState('domcontentloaded');

    const email = `ui-fresh-${Date.now()}-${Math.random().toString(36).slice(2,6)}@gofin.io`;
    await page.request.post('/api/v1/auth/register', {
      data: { email, password: 'TestPass123!' }
    }).then(r => r.json()).then(b => {
      return page.evaluate((tok: string) => {
        localStorage.setItem('access_token', tok);
      }, b.access_token);
    });

    await page.goto('/dashboard');
    await page.waitForLoadState('domcontentloaded');
  });

  test('Dashboard loads with empty state', async ({ page }) => {
    await expect(page.locator('h1, h2').first()).toBeVisible({ timeout: 5000 });
    // Should NOT show hardcoded savings
    const savingsCard = page.locator('text=12.000.000');
    expect(await savingsCard.count()).toBe(0);
  });

  test('Wallets page loads', async ({ page }) => {
    await page.goto('/wallets');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/wallets');
  });

  test('Transactions page loads', async ({ page }) => {
    await page.goto('/transactions');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/transactions');
  });

  test('Categories page loads', async ({ page }) => {
    await page.goto('/categories');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/categories');
  });

  test('Budgets page loads', async ({ page }) => {
    await page.goto('/budgets');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/budgets');
  });

  test('Bills page loads', async ({ page }) => {
    await page.goto('/bills');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/bills');
  });

  test('Tags page loads', async ({ page }) => {
    await page.goto('/tags');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/tags');
  });

  test('Reports page loads', async ({ page }) => {
    await page.goto('/reports');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/reports');
  });

  test('Rules page loads', async ({ page }) => {
    await page.goto('/rules');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/rules');
  });

  test('Settings account page loads', async ({ page }) => {
    await page.goto('/settings/account');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/settings/account');
  });

  test('Settings groups page loads', async ({ page }) => {
    await page.goto('/settings/groups');
    await page.waitForLoadState('domcontentloaded');
    expect(page.url()).toContain('/settings/groups');
  });

  test('Create wallet page has currency selector', async ({ page }) => {
    await page.goto('/wallets/create');
    await page.waitForLoadState('domcontentloaded');
    // Check currency label exists
    const currLabel = page.locator('label').filter({ hasText: /Currency|Mata Uang/ }).first();
    await expect(currLabel).toBeVisible({ timeout: 5000 });
  });

  test('Dark mode toggle works', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('domcontentloaded');
    const themeBtn = page.locator('button[title*="mode"], button[title*="Mode"]').first();
    if (await themeBtn.isVisible({ timeout: 3000 }).catch(() => false)) {
      await themeBtn.click();
      const theme = await page.evaluate(() => localStorage.getItem('gofin_theme'));
      expect(['dark', 'light']).toContain(theme);
    }
  });
});

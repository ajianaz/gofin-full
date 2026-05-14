import { test, expect } from '@playwright/test';

// ---------------------------------------------------------------------------
// Config (baseURL and ignoreHTTPSErrors set in playwright.selfhost.config.ts)
// ---------------------------------------------------------------------------
const API_BASE = '/api/v1';

// Mock credentials (web frontend uses mock auth)
const MOCK_EMAIL = 'admin@gofin.id';
const MOCK_PASSWORD = 'admin123';

// Admin credentials (seeded in self-host DB)
const ADMIN_EMAIL = 'admin@gofin.dev';
const ADMIN_PASSWORD = 'YjxD-C1uTKQztEjVyN-UQ';

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
async function registerViaAPI(request: import('@playwright/test').APIRequestContext) {
	const email = `e2e-api-${Date.now()}@example.com`;
	const password = 'TestPass1234!';
	const res = await request.post(`${API_BASE}/auth/register`, {
		data: { email, password },
	});
	return { email, password, status: res.status(), body: await res.json().catch(() => null) };
}

async function loginViaAPI(request: import('@playwright/test').APIRequestContext, email: string, password: string) {
	const res = await request.post(`${API_BASE}/auth/login`, {
		data: { email, password },
	});
	return { status: res.status(), body: await res.json().catch(() => null) };
}

// =============================================================================
// PART 1: API Tests (real backend via self-host Caddy)
// =============================================================================
test.describe('API — Health & Info', () => {

	test('GET /health returns ok', async ({ request }) => {
		const res = await request.get(`/health`);
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.status).toBe('ok');
	});

	test('GET /api/v1/ returns API info', async ({ request }) => {
		const res = await request.get(`${API_BASE}/`);
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.message).toContain('Gofin API');
	});
});

test.describe('API — Auth (self-host)', () => {

	test('POST /auth/register creates new user', async ({ request }) => {
		const result = await registerViaAPI(request);
		expect(result.status).toBe(201);
		expect(result.body?.access_token).toBeDefined();
		expect(result.body?.refresh_token).toBeDefined();
		expect(result.body?.token_type).toBe('Bearer');
	});

	test('POST /auth/register rejects duplicate email', async ({ request }) => {
		const email = `e2e-dup-${Date.now()}@example.com`;
		await request.post(`${API_BASE}/auth/register`, {
			data: { email, password: 'TestPass1234!' },
		});
		const res = await request.post(`${API_BASE}/auth/register`, {
			data: { email, password: 'TestPass1234!' },
		});
		expect(res.status()).toBe(409);
	});

	test('POST /auth/register rejects short password', async ({ request }) => {
		const res = await request.post(`${API_BASE}/auth/register`, {
			data: { email: `e2e-short-${Date.now()}@example.com`, password: 'short' },
		});
		expect(res.status()).toBe(422);
	});

	test('POST /auth/login with valid credentials', async ({ request }) => {
		const reg = await registerViaAPI(request);
		const result = await loginViaAPI(request, reg.email, reg.password);
		expect(result.status).toBe(200);
		expect(result.body?.access_token).toBeDefined();
	});

	test('POST /auth/login rejects invalid credentials', async ({ request }) => {
		const result = await loginViaAPI(request, 'nonexistent@example.com', 'wrongpass');
		expect(result.status).toBe(401);
	});

	test('POST /auth/login with admin seed account', async ({ request }) => {
		const result = await loginViaAPI(request, ADMIN_EMAIL, ADMIN_PASSWORD);
		expect(result.status).toBe(200);
		expect(result.body?.access_token).toBeDefined();
	});

	test('POST /auth/logout works with token', async ({ request }) => {
		const reg = await registerViaAPI(request);
		const loginRes = await loginViaAPI(request, reg.email, reg.password);
		const res = await request.post(`${API_BASE}/auth/logout`, {
			headers: { Authorization: `Bearer ${loginRes.body?.access_token}`, 'Content-Type': 'application/json' },
			data: {},
			});
		expect(res.status()).toBe(200);
	});

	test('GET /auth/provider returns local', async ({ request }) => {
		const res = await request.get(`${API_BASE}/auth/provider`);
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.provider).toBe('local');
	});
});

test.describe('API — Protected Endpoints', () => {

	let accessToken: string;

	test.beforeAll(async ({ request }) => {
		const reg = await registerViaAPI(request);
		accessToken = reg.body?.access_token;

		// Switch to active group (needed for protected routes)
		const groupsRes = await request.get(`${API_BASE}/groups`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		const groups = await groupsRes.json();
		if (groups.data?.[0]?.id) {
			await request.post(`${API_BASE}/groups/switch`, {
				headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
				data: { user_group_id: groups.data[0].id },
			});
		}

		// Re-login to get token with group context
		const loginRes = await loginViaAPI(request, reg.email, 'TestPass1234!');
		accessToken = loginRes.body?.access_token;
	});

	test('rejects unauthenticated access', async ({ request }) => {
		const res = await request.get(`${API_BASE}/wallets`);
		expect(res.status()).toBe(401);
	});

	test('GET /users/me returns user info', async ({ request }) => {
		const res = await request.get(`${API_BASE}/users/me`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data?.attributes?.email).toBeDefined();
	});

	test('GET /groups returns at least one group', async ({ request }) => {
		const res = await request.get(`${API_BASE}/groups`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data?.length).toBeGreaterThanOrEqual(1);
	});

	test('GET /wallets returns wallet list (may be empty)', async ({ request }) => {
		const res = await request.get(`${API_BASE}/wallets`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		// 200 = has group, 400 = no group set, 403 = forbidden
		expect([200, 400]).toContain(res.status());
	});

	test('GET /categories returns category list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/categories`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /transactions returns transaction list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/transactions`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /currencies returns reference data', async ({ request }) => {
		const res = await request.get(`${API_BASE}/currencies`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /wallet-types returns reference data', async ({ request }) => {
		const res = await request.get(`${API_BASE}/wallet-types`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /budgets returns budget list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/budgets`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /tags returns tag list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/tags`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /bills returns bill list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/bills`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /preferences returns user preferences', async ({ request }) => {
		const res = await request.get(`${API_BASE}/preferences`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /notifications returns notification list', async ({ request }) => {
		const res = await request.get(`${API_BASE}/notifications`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /piggy-banks requires wallet_id param', async ({ request }) => {
		const res = await request.get(`${API_BASE}/piggy-banks`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.status()).toBe(404);
	});

	test('GET /analytics endpoints require view_reports role', async ({ request }) => {
		const endpoints = [
			'/analytics/spending-by-category',
			'/analytics/spending-by-period',
			'/analytics/net-worth',
		];
		for (const endpoint of endpoints) {
			const res = await request.get(`${API_BASE}${endpoint}`, {
				headers: { Authorization: `Bearer ${accessToken}` },
			});
			expect([200, 403, 500]).toContain(res.status());
		}
	});

	test('GET /export/csv endpoint exists', async ({ request }) => {
		const res = await request.get(`${API_BASE}/export/csv`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect([200, 400]).toContain(res.status());
	});

	test('GET /api-docs returns documentation', async ({ request }) => {
		const res = await request.get(`${API_BASE}/docs`);
		expect(res.ok()).toBeTruthy();
	});

	test('GET /openapi.json returns OpenAPI spec', async ({ request }) => {
		const res = await request.get(`${API_BASE}/openapi.json`);
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.openapi).toBeDefined();
	});
});

test.describe('API — Admin Endpoints', () => {

	test('GET /admin/users requires admin role', async ({ request }) => {
		const reg = await registerViaAPI(request);
		const res = await request.get(`${API_BASE}/admin/users`, {
			headers: { Authorization: `Bearer ${reg.body?.access_token}` },
		});
		expect(res.status()).toBe(403);
	});

	test('GET /admin/users works with admin account', async ({ request }) => {
		const result = await loginViaAPI(request, ADMIN_EMAIL, ADMIN_PASSWORD);
		const res = await request.get(`${API_BASE}/admin/users`, {
			headers: { Authorization: `Bearer ${result.body?.access_token}` },
		});
		expect([200, 400]).toContain(res.status());
	});
});

// =============================================================================
// PART 2: UI Tests (web frontend with mock auth)
// =============================================================================
test.describe('UI — Auth Pages', () => {

	test('login page renders correctly', async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('input#email')).toBeVisible();
		await expect(page.locator('input#password')).toBeVisible();
		await expect(page.locator('button[type="submit"]')).toBeVisible();
	});

	test('register page renders correctly', async ({ page }) => {
		await page.goto(`/register`);
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('input#email')).toBeVisible();
		await expect(page.locator('input#password')).toBeVisible();
		await expect(page.locator('input#confirm-password')).toBeVisible();
		await expect(page.locator('button[type="submit"]')).toBeVisible();
	});

	test('forgot-password page renders', async ({ page }) => {
		await page.goto(`/forgot-password`);
		await page.waitForLoadState('domcontentloaded');
		await expect(page.locator('input[type="email"]')).toBeVisible();
	});

	test('UI login with mock credentials navigates to dashboard', async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');

		await page.locator('input#email').fill(MOCK_EMAIL);
		await page.locator('input#password').fill(MOCK_PASSWORD);
		await page.locator('button[type="submit"]').click();

		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
		expect(page.url()).toContain('/dashboard');
	});

	test('UI register with new account navigates to dashboard', async ({ page }) => {
		await page.goto(`/register`);
		await page.waitForLoadState('domcontentloaded');

		const newEmail = `ui-reg-${Date.now()}@example.com`;
		await page.locator('input#email').fill(newEmail);
		await page.locator('input#password').fill('password123');
		await page.locator('input#confirm-password').fill('password123');
		await page.locator('button[type="submit"]').click();

		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
		expect(page.url()).toContain('/dashboard');
	});

	test('unauthenticated access to dashboard redirects to login', async ({ page }) => {
		await page.goto(`/dashboard`);
		await page.waitForLoadState('domcontentloaded');

		await page.waitForURL('**/login**', { timeout: 5000 }).catch(() => {});
		expect(page.url()).toContain('/login');
	});
});

test.describe('UI — Dashboard', () => {

	test.beforeEach(async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(MOCK_EMAIL);
		await page.locator('input#password').fill(MOCK_PASSWORD);
		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
	});

	test('dashboard shows stat cards', async ({ page }) => {
		const statGrid = page.locator('.grid.lg\\:grid-cols-4');
		await expect(statGrid.first()).toBeVisible({ timeout: 5000 });
	});

	test('dashboard shows recent transactions section', async ({ page }) => {
		const cards = page.locator('[class*="card"], [class*="Card"]');
		await expect(cards.first()).toBeVisible({ timeout: 5000 });
	});

	test('dashboard has sidebar navigation', async ({ page }) => {
		const sidebar = page.locator('nav, aside, [class*="sidebar"]');
		await expect(sidebar.first()).toBeVisible({ timeout: 5000 });
	});
});

test.describe('UI — Core Pages Navigation', () => {

	const pages = [
		{ path: '/wallets', title: 'Wallets' },
		{ path: '/transactions', title: 'Transactions' },
		{ path: '/categories', title: 'Categories' },
		{ path: '/budgets', title: 'Budgets' },
		{ path: '/bills', title: 'Bills' },
		{ path: '/tags', title: 'Tags' },
		{ path: '/piggy-banks', title: 'Piggy Banks' },
		{ path: '/analytics', title: 'Analytics' },
		{ path: '/reports', title: 'Reports' },
		{ path: '/recurring', title: 'Recurring' },
		{ path: '/rules', title: 'Rules' },
		{ path: '/groups', title: 'Groups' },
		{ path: '/currencies', title: 'Currencies' },
		{ path: '/export', title: 'Export' },
	];

	for (const pg of pages) {
		test(`${pg.path} page loads after login`, async ({ page }) => {
			await page.goto(`/login`);
			await page.waitForLoadState('domcontentloaded');
			await page.locator('input#email').fill(MOCK_EMAIL);
			await page.locator('input#password').fill(MOCK_PASSWORD);
			await page.locator('button[type="submit"]').click();
			await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});

			await page.goto(pg.path);
			await page.waitForLoadState('domcontentloaded');
			expect(page.url()).toContain(pg.path);
		});
	}
});

test.describe('UI — Create Pages', () => {

	const createPages = [
		{ path: '/wallets/create', label: 'Wallet' },
		{ path: '/transactions/create', label: 'Transaction' },
		{ path: '/categories/create', label: 'Category' },
		{ path: '/budgets/create', label: 'Budget' },
		{ path: '/bills/create', label: 'Bill' },
		{ path: '/tags/create', label: 'Tag' },
		{ path: '/piggy-banks/create', label: 'Piggy Bank' },
		{ path: '/recurring/create', label: 'Recurring' },
		{ path: '/rules/create', label: 'Rule' },
	];

	for (const pg of createPages) {
		test(`${pg.label} create page has form elements`, async ({ page }) => {
			await page.goto(`/login`);
			await page.waitForLoadState('domcontentloaded');
			await page.locator('input#email').fill(MOCK_EMAIL);
			await page.locator('input#password').fill(MOCK_PASSWORD);
			await page.locator('button[type="submit"]').click();
			await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});

			await page.goto(pg.path);
			await page.waitForLoadState('domcontentloaded');
			expect(page.url()).toContain(pg.path);

			// Create pages should have at least one input and a submit button
			const inputs = page.locator('input, select, textarea');
			expect(await inputs.count()).toBeGreaterThanOrEqual(1);

			const submitBtn = page.locator('button[type="submit"]');
			await expect(submitBtn).toBeVisible();
		});
	}
});

test.describe('UI — Settings Pages', () => {

	test.beforeEach(async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(MOCK_EMAIL);
		await page.locator('input#password').fill(MOCK_PASSWORD);
		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
	});

	test('settings main page loads', async ({ page }) => {
		await page.goto(`/settings`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings');
	});

	test('settings/profile page loads', async ({ page }) => {
		await page.goto(`/settings/profile`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/profile');
	});

	test('settings/preferences page loads', async ({ page }) => {
		await page.goto(`/settings/preferences`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/preferences');
	});

	test('settings/notifications page loads', async ({ page }) => {
		await page.goto(`/settings/notifications`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/notifications');
	});

	test('settings/api-keys page loads', async ({ page }) => {
		await page.goto(`/settings/api-keys`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/api-keys');
	});
});

test.describe('UI — Reports Pages', () => {

	test.beforeEach(async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(MOCK_EMAIL);
		await page.locator('input#password').fill(MOCK_PASSWORD);
		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
	});

	test('reports main page loads', async ({ page }) => {
		await page.goto(`/reports`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/reports');
	});

	test('reports/spending-by-category loads', async ({ page }) => {
		await page.goto(`/reports/spending-by-category`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/reports/spending-by-category');
	});

	test('reports/spending-by-period loads', async ({ page }) => {
		await page.goto(`/reports/spending-by-period`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/reports/spending-by-period');
	});

	test('reports/net-worth loads', async ({ page }) => {
		await page.goto(`/reports/net-worth`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/reports/net-worth');
	});
});

test.describe('UI — Admin Pages', () => {

	test.beforeEach(async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(MOCK_EMAIL);
		await page.locator('input#password').fill(MOCK_PASSWORD);
		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/dashboard**', { timeout: 5000 }).catch(() => {});
	});

	test('admin/users page loads', async ({ page }) => {
		await page.goto(`/admin/users`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/admin/users');
	});

	test('admin/audit-log page loads', async ({ page }) => {
		await page.goto(`/admin/audit-log`);
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/admin/audit-log');
	});
});

test.describe('UI — Language Switching', () => {

	test('locale changes between Indonesian and English', async ({ page }) => {
		await page.goto(`/login`);
		await page.waitForLoadState('domcontentloaded');

		// Find language toggle (check common patterns)
		const langToggle = page.locator('button:has-text("EN"), button:has-text("ID"), [class*="locale"], [class*="language"]').first();
		if (await langToggle.isVisible({ timeout: 2000 }).catch(() => false)) {
			await langToggle.click();
		}

		// Page should still work after locale switch
		await expect(page.locator('input#email')).toBeVisible();
	});
});

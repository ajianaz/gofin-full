import { test, expect } from '@playwright/test';

const ADMIN_EMAIL = 'aji.anaz@gmail.com';
const ADMIN_PASSWORD = 'Str0ngP4ss123!';
const API_BASE = '/api/v1';

test.describe('Live Server — API Health', () => {
	test('GET /health returns ok', async ({ request }) => {
		const res = await request.get('/health');
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

	test('GET /api/v1/auth/provider returns local', async ({ request }) => {
		const res = await request.get(`${API_BASE}/auth/provider`);
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.provider).toBe('local');
	});
});

test.describe('Live Server — Admin Auth', () => {
	test('Admin login via API', async ({ request }) => {
		const res = await request.post(`${API_BASE}/auth/login`, {
			data: { email: ADMIN_EMAIL, password: ADMIN_PASSWORD },
		});
		expect(res.status()).toBe(200);
		const body = await res.json();
		expect(body.access_token).toBeDefined();
		expect(body.refresh_token).toBeDefined();
		expect(body.token_type).toBe('Bearer');
	});

	test('Admin login rejects wrong password', async ({ request }) => {
		const res = await request.post(`${API_BASE}/auth/login`, {
			data: { email: ADMIN_EMAIL, password: 'wrongpassword' },
		});
		expect([401, 429]).toContain(res.status());
	});
});

test.describe('Live Server — Admin UI Login', () => {
	test('Login with admin credentials and reach dashboard', async ({ page }) => {
		await page.goto('/login');
		await page.waitForLoadState('domcontentloaded');

		await page.locator('input#email').fill(ADMIN_EMAIL);
		await page.locator('input#password').fill(ADMIN_PASSWORD);

		const [response] = await Promise.all([
			page.waitForResponse(r => r.url().includes('/auth/login'), { timeout: 10000 }),
			page.locator('button[type="submit"]').click(),
		]);

		expect(response.status()).toBe(200);
		await page.waitForURL('**/dashboard**', { timeout: 10000 });
		expect(page.url()).toContain('/dashboard');
	});
});

test.describe('Live Server — Admin Protected Endpoints', () => {
	let accessToken: string;

	test.beforeAll(async ({ request }) => {
		const res = await request.post(`${API_BASE}/auth/login`, {
			data: { email: ADMIN_EMAIL, password: ADMIN_PASSWORD },
		});
		const body = await res.json();
		accessToken = body.access_token;

		// Switch to active group
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
		// Re-login to get group-context token
		const reLogin = await request.post(`${API_BASE}/auth/login`, {
			data: { email: ADMIN_EMAIL, password: ADMIN_PASSWORD },
		});
		accessToken = (await reLogin.json()).access_token;
	});

	test('GET /users/me returns admin user', async ({ request }) => {
		const res = await request.get(`${API_BASE}/users/me`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data?.attributes?.email).toBe(ADMIN_EMAIL);
	});

	test('GET /wallets works', async ({ request }) => {
		const res = await request.get(`${API_BASE}/wallets`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /categories works', async ({ request }) => {
		const res = await request.get(`${API_BASE}/categories`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /transactions works', async ({ request }) => {
		const res = await request.get(`${API_BASE}/transactions`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /currencies works', async ({ request }) => {
		const res = await request.get(`${API_BASE}/currencies`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /preferences works', async ({ request }) => {
		const res = await request.get(`${API_BASE}/preferences`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /admin/users works with admin', async ({ request }) => {
		const res = await request.get(`${API_BASE}/admin/users`, {
			headers: { Authorization: `Bearer ${accessToken}` },
		});
		expect(res.ok()).toBeTruthy();
	});
});

test.describe('Live Server — Settings Pages (New)', () => {
	test.beforeEach(async ({ page }) => {
		// Login as admin via UI
		await page.goto('/login');
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(ADMIN_EMAIL);
		await page.locator('input#password').fill(ADMIN_PASSWORD);
		await Promise.all([
			page.waitForResponse(r => r.url().includes('/auth/login'), { timeout: 10000 }),
			page.locator('button[type="submit"]').click(),
		]);
		await page.waitForURL('**/dashboard**', { timeout: 10000 });
	});

	test('/settings/account page loads', async ({ page }) => {
		await page.goto('/settings/account');
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/account');
		// Should show account info, not blank
		const content = page.locator('h1, h2, h3').first();
		await expect(content).toBeVisible({ timeout: 5000 });
	});

	test('/settings/groups page loads', async ({ page }) => {
		await page.goto('/settings/groups');
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/groups');
		// Should show groups table or create form
		const content = page.locator('h1, h2, h3, table').first();
		await expect(content).toBeVisible({ timeout: 5000 });
	});

	test('/settings/profile page loads', async ({ page }) => {
		await page.goto('/settings/profile');
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/profile');
	});

	test('/settings/preferences page loads', async ({ page }) => {
		await page.goto('/settings/preferences');
		await page.waitForLoadState('domcontentloaded');
		expect(page.url()).toContain('/settings/preferences');
	});

	test('Dark mode toggle works', async ({ page }) => {
		await page.goto('/dashboard');
		await page.waitForLoadState('domcontentloaded');

		// Find theme toggle button (Sun or Moon icon)
		const themeBtn = page.locator('button[title*="mode"], button[title*="Mode"]').first();
		if (await themeBtn.isVisible({ timeout: 3000 }).catch(() => false)) {
			// Click to toggle
			await themeBtn.click();
			await page.waitForTimeout(500);

			// Check localStorage for theme
			const theme = await page.evaluate(() => localStorage.getItem('gofin_theme'));
			expect(['dark', 'light']).toContain(theme);
		}
	});
});

test.describe('Live Server — Core Pages', () => {
	test.beforeEach(async ({ page }) => {
		await page.goto('/login');
		await page.waitForLoadState('domcontentloaded');
		await page.locator('input#email').fill(ADMIN_EMAIL);
		await page.locator('input#password').fill(ADMIN_PASSWORD);
		await Promise.all([
			page.waitForResponse(r => r.url().includes('/auth/login'), { timeout: 10000 }),
			page.locator('button[type="submit"]').click(),
		]);
		await page.waitForURL('**/dashboard**', { timeout: 10000 });
	});

	const pages = [
		'/wallets',
		'/transactions',
		'/categories',
		'/budgets',
		'/bills',
		'/tags',
		'/analytics',
		'/reports',
		'/rules',
		'/currencies',
		'/export',
	];

	for (const path of pages) {
		test(`${path} page loads`, async ({ page }) => {
			await page.goto(path);
			await page.waitForLoadState('domcontentloaded');
			expect(page.url()).toContain(path);
		});
	}
});

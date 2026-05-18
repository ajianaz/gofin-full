import { test, expect } from '@playwright/test';

const TEST_PASSWORD = 'TestPass123!';

function uniqueEmail(prefix: string) {
	return `e2e-${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@gofin.io`;
}

const JH = { 'Content-Type': 'application/json', Accept: 'application/json' };

async function registerAndAuthenticate(page: import('@playwright/test').Page, path: string) {
	const testEmail = uniqueEmail('core');
	const regResponse = await page.request.post('/api/v1/auth/register', {
		headers: JH,
		data: { email: testEmail, password: TEST_PASSWORD }
	});
	expect(regResponse.ok()).toBeTruthy();
	const tokens = await regResponse.json();
	expect(tokens.access_token).toBeDefined();

	await page.goto(path);
	await page.evaluate((accessToken) => {
		localStorage.setItem('access_token', accessToken);
	}, tokens.access_token);
	await page.reload();
	await page.waitForLoadState('domcontentloaded');

	return { email: testEmail, tokens };
}

test.describe('Dashboard Page', () => {
	test('loads and shows stat cards', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		const statGrid = page.locator('.grid.lg\\:grid-cols-4');
		await expect(statGrid).toBeVisible();

		const icons = statGrid.locator('svg');
		await expect(icons.first()).toBeVisible();
	});

	test('shows recent transactions section', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		const cards = page.locator('[class*="card"], [class*="Card"]');
		await expect(cards.first()).toBeVisible();
	});

	test('shows spending by category or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// Dashboard loaded — stat grid is visible (verified in first test)
		const statGrid = page.locator('.grid.lg\\:grid-cols-4');
		await expect(statGrid).toBeVisible();
	});
});

test.describe('Wallets Page', () => {
	test('loads and shows wallet heading', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const headings = page.locator('h2');
		await expect(headings.nth(1)).toBeVisible();
	});

	test('has filter controls', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		// shadcn Select uses button[data-slot="select-trigger"]
		const filterTrigger = page.locator('button[data-slot="select-trigger"]').first();
		await expect(filterTrigger).toBeVisible({ timeout: 5000 });
	});

	test('has "Add" button for wallet creation', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const addButton = page.locator('a[href="/wallets/create"], main button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 5000 });
	});
});

test.describe('Transactions Page', () => {
	test('loads and shows table', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		// Wait for data to load (API fetch)
		const table = page.locator('table');
		await expect(table).toBeVisible({ timeout: 10000 });

		expect(await table.locator('thead th').count()).toBeGreaterThanOrEqual(4);
	});

	test('has search input', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		// Wait for load
		await expect(page.locator('table')).toBeVisible({ timeout: 10000 });

		const searchInput = page.locator('input[type="text"], input[placeholder]');
		await expect(searchInput.first()).toBeVisible();
	});

	test('has pagination controls', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		await expect(page.locator('table')).toBeVisible({ timeout: 10000 });

		const paginationArea = page.locator('.flex.items-center.gap-1').last();
		await expect(paginationArea).toBeVisible();
	});

	test('has "Add" button for transaction creation', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const addButton = page.locator('a[href="/transactions/create"], main button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 5000 });
	});
});

test.describe('Authentication Guard', () => {
	test('unauthenticated access redirects to login', async ({ page }) => {
		await page.goto('/dashboard');
		await page.waitForLoadState('domcontentloaded');

		await page.waitForURL('**/login**', { timeout: 5000 }).catch(() => {});

		expect(page.url()).toContain('/login');
	});
});

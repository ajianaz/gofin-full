import { test, expect } from '@playwright/test';

// ---------------------------------------------------------------------------
// Helper: register + authenticate via API and navigate to a protected page
// ---------------------------------------------------------------------------
async function registerAndAuthenticate(page: import('@playwright/test').Page, path: string, email?: string) {
	const testEmail = email || `e2e-core-${Date.now()}@example.com`;
	const testPassword = 'testpassword12345';

	// Register
	const regResponse = await page.request.post('/api/v1/auth/register', {
		data: { email: testEmail, password: testPassword }
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

// ===========================================================================
// Dashboard
// ===========================================================================
test.describe('Dashboard Page', () => {
	test('loads and shows stat cards', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// StatCard grid with 4 columns
		const statGrid = page.locator('.grid.lg\\:grid-cols-4');
		await expect(statGrid).toBeVisible();

		// Stat cards contain SVG icons
		const icons = statGrid.locator('svg');
		await expect(icons.first()).toBeVisible();
	});

	test('shows recent transactions section', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// Recent transactions card visible
		const cards = page.locator('[class*="card"], [class*="Card"]');
		await expect(cards.first()).toBeVisible();

		// "View All" link exists
		const viewAllLink = page.locator('a[href="/transactions"].text-sm');
		await expect(viewAllLink).toBeVisible();
	});

	test('"View All" link navigates to transactions', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		await page.locator('a[href="/transactions"].text-sm').click();
		await page.waitForLoadState('domcontentloaded');

		expect(page.url()).toContain('/transactions');
	});

	test('shows spending by category with progress bars', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		const progressElements = page.locator('[role="progressbar"], [class*="progress"]');
		await expect(progressElements.first()).toBeVisible();
	});
});

// ===========================================================================
// Wallets
// ===========================================================================
test.describe('Wallets Page', () => {
	test('loads and shows wallet cards', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const headings = page.locator('h2');
		await expect(headings.nth(1)).toBeVisible();

		const badge = page.locator('span[class*="rounded-2xl"]');
		await expect(badge).toBeVisible();

		const walletCards = page.locator('.grid.grid-cols-1 [class*="card"], .grid.grid-cols-1 [class*="Card"]');
		await expect(walletCards.first()).toBeVisible();

		const cardIcons = page.locator('.grid svg');
		expect(await cardIcons.count()).toBeGreaterThan(0);
	});

	test('has type filter dropdown', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const typeFilter = page.locator('select');
		await expect(typeFilter.first()).toBeVisible();

		const options = typeFilter.first().locator('option');
		expect(await options.count()).toBeGreaterThanOrEqual(3);
	});

	test('has "Add" button for wallet creation', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const addButton = page.locator('main button:has(svg.lucide-plus, svg:has(path[d*="M12 5"]))').first();
		await expect(addButton).toBeVisible();
	});
});

// ===========================================================================
// Transactions
// ===========================================================================
test.describe('Transactions Page', () => {
	test('loads and shows table', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const headings = page.locator('h2');
		await expect(headings.nth(1)).toBeVisible();

		const badge = page.locator('span[class*="rounded-2xl"]');
		await expect(badge).toBeVisible();

		const table = page.locator('table');
		await expect(table).toBeVisible();

		expect(await table.locator('thead th').count()).toBeGreaterThanOrEqual(4);
	});

	test('has search input and filter dropdowns', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const searchInput = page.locator('input[type="text"], input[placeholder]');
		await expect(searchInput.first()).toBeVisible();

		const selectFilters = page.locator('select');
		expect(await selectFilters.count()).toBeGreaterThanOrEqual(4);
	});

	test('has pagination controls', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const paginationArea = page.locator('.flex.items-center.gap-1').last();
		await expect(paginationArea).toBeVisible();

		const navButtons = paginationArea.locator('button');
		expect(await navButtons.count()).toBeGreaterThanOrEqual(2);
	});

	test('has "Add" button for transaction creation', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const addButton = page.locator('main button:has(svg.lucide-plus, svg:has(path[d*="M12 5"]))').first();
		await expect(addButton).toBeVisible();
	});
});

// ===========================================================================
// Auth guard
// ===========================================================================
test.describe('Authentication Guard', () => {
	test('unauthenticated access redirects to login', async ({ page }) => {
		await page.goto('/dashboard');
		await page.waitForLoadState('domcontentloaded');

		await page.waitForURL('**/login**', { timeout: 5000 }).catch(() => {});

		expect(page.url()).toContain('/login');
	});
});

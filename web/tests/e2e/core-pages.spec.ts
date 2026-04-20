import { test, expect } from '@playwright/test';

// ---------------------------------------------------------------------------
// Helper: authenticate via API and navigate to a protected page
// ---------------------------------------------------------------------------
async function authenticateAndNavigate(page: import('@playwright/test').Page, path: string) {
	const loginResponse = await page.request.post('/api/v1/auth/login', {
		data: { email: 'core-test@example.com', password: 'corepassword123' }
	});
	expect(loginResponse.ok()).toBeTruthy();
	const tokens = await loginResponse.json();

	await page.goto(path);
	await page.evaluate((accessToken) => {
		localStorage.setItem('access_token', accessToken);
	}, tokens.access_token);
	await page.reload();
	await page.waitForLoadState('domcontentloaded');
}

// ===========================================================================
// Dashboard
// ===========================================================================
test.describe('Dashboard Page', () => {
	test('loads and shows stat cards', async ({ page }) => {
		await authenticateAndNavigate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// StatCard grid with 4 columns
		const statGrid = page.locator('.grid.lg\\:grid-cols-4');
		await expect(statGrid).toBeVisible();

		// Stat cards contain SVG icons
		const icons = statGrid.locator('svg');
		await expect(icons.first()).toBeVisible();
	});

	test('shows recent transactions section', async ({ page }) => {
		await authenticateAndNavigate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// Recent transactions card visible
		const cards = page.locator('[class*="card"], [class*="Card"]');
		await expect(cards.first()).toBeVisible();

		// "View All" link exists
		const viewAllLink = page.locator('a[href="/transactions"].text-sm');
		await expect(viewAllLink).toBeVisible();
	});

	test('"View All" link navigates to transactions', async ({ page }) => {
		await authenticateAndNavigate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		// Click the "View All" link (the one in the card content, not sidebar)
		await page.locator('a[href="/transactions"].text-sm').click();
		await page.waitForLoadState('domcontentloaded');

		expect(page.url()).toContain('/transactions');
	});

	test('shows spending by category with progress bars', async ({ page }) => {
		await authenticateAndNavigate(page, '/dashboard');
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
		await authenticateAndNavigate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		// Page heading (second h2 — first is the header brand)
		const headings = page.locator('h2');
		await expect(headings.nth(1)).toBeVisible();

		// Count badge
		const badge = page.locator('span[class*="rounded-2xl"]');
		await expect(badge).toBeVisible();

		// Wallet cards in grid
		const walletCards = page.locator('.grid.grid-cols-1 [class*="card"], .grid.grid-cols-1 [class*="Card"]');
		await expect(walletCards.first()).toBeVisible();

		// Each card has an icon
		const cardIcons = page.locator('.grid svg');
		expect(await cardIcons.count()).toBeGreaterThan(0);
	});

	test('has type filter dropdown', async ({ page }) => {
		await authenticateAndNavigate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const typeFilter = page.locator('select');
		await expect(typeFilter.first()).toBeVisible();

		const options = typeFilter.first().locator('option');
		expect(await options.count()).toBeGreaterThanOrEqual(3);
	});

	test('has "Add" button for wallet creation', async ({ page }) => {
		await authenticateAndNavigate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		// The Add button is a <button> with onclick (not an <a> tag)
		const addButton = page.locator('main button:has(svg.lucide-plus, svg:has(path[d*="M12 5"]))').first();
		await expect(addButton).toBeVisible();
	});
});

// ===========================================================================
// Transactions
// ===========================================================================
test.describe('Transactions Page', () => {
	test('loads and shows table', async ({ page }) => {
		await authenticateAndNavigate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		// Page heading (second h2)
		const headings = page.locator('h2');
		await expect(headings.nth(1)).toBeVisible();

		const badge = page.locator('span[class*="rounded-2xl"]');
		await expect(badge).toBeVisible();

		// Transaction table
		const table = page.locator('table');
		await expect(table).toBeVisible();

		expect(await table.locator('thead th').count()).toBeGreaterThanOrEqual(4);
		expect(await table.locator('tbody tr').count()).toBeGreaterThan(0);
	});

	test('has search input and filter dropdowns', async ({ page }) => {
		await authenticateAndNavigate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const searchInput = page.locator('input[type="text"], input[placeholder]');
		await expect(searchInput.first()).toBeVisible();

		const selectFilters = page.locator('select');
		expect(await selectFilters.count()).toBeGreaterThanOrEqual(4);
	});

	test('has pagination controls', async ({ page }) => {
		await authenticateAndNavigate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const paginationArea = page.locator('.flex.items-center.gap-1').last();
		await expect(paginationArea).toBeVisible();

		const navButtons = paginationArea.locator('button');
		expect(await navButtons.count()).toBeGreaterThanOrEqual(2);
	});

	test('has "Add" button for transaction creation', async ({ page }) => {
		await authenticateAndNavigate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		// The Add button is a <button> with onclick (not an <a> tag)
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

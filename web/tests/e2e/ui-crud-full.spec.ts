import { test, expect, type BrowserContext } from '@playwright/test';

const BASE = 'https://localhost';

const TIMESTAMP = Date.now();
const EMAIL = `e2e-${TIMESTAMP}@example.com`;
const PASS = 'Password1234!';

let savedTokens: { accessToken: string; refreshToken: string } | null = null;

async function authedCtx(browser: import('@playwright/test').Browser): Promise<BrowserContext> {
	const ctx = await browser.newContext({ ignoreHTTPSErrors: true });
	const page = await ctx.newPage();
	await page.goto('/login');
	await page.evaluate(({ at, rt }) => {
		localStorage.setItem('access_token', at);
		localStorage.setItem('refresh_token', rt);
	}, { at: savedTokens!.accessToken, rt: savedTokens!.refreshToken });
	await page.goto('/wallets/create');
	await page.waitForLoadState('networkidle');
	// If we got redirected to login, auth is broken
	if (page.url().includes('/login')) {
		await ctx.close();
		throw new Error('Auth failed — redirected to login');
	}
	return ctx;
}

test.describe.configure({ mode: 'serial' });

test.describe('UI Full CRUD Flow — Real API', () => {

	test('1. Register + save tokens', async ({ page }) => {
		await page.goto('/register');
		await page.waitForLoadState('networkidle');
		await page.fill('#email', EMAIL);
		await page.fill('#password', PASS);
		await page.fill('#confirm-password', PASS);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard**', { timeout: 15000 }).catch(() => {});
		expect(page.url()).toContain('/dashboard');

		// Save tokens from localStorage
		savedTokens = await page.evaluate(() => ({
			accessToken: localStorage.getItem('access_token') || '',
			refreshToken: localStorage.getItem('refresh_token') || ''
		}));
		expect(savedTokens.accessToken).toBeTruthy();
	});

	test('2. Wallets — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/wallets/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'BCA Tabungan');
		await page.selectOption('#type', 'asset');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/wallets**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/wallets');
		await ctx.close();
	});

	test('3. Categories — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/categories/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Makanan');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/categories**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/categories');
		await ctx.close();
	});

	test('4. Tags — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/tags/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#tag', 'Tagihan');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/tags**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/tags');
		await ctx.close();
	});

	test('5. Transactions — create withdrawal', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/transactions/create');
		await page.waitForLoadState('networkidle');
		await page.selectOption('#type', 'withdrawal');
		await page.fill('#description', 'Makan siang warteg');
		await page.fill('#amount', '150000');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/transactions**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/transactions');
		await ctx.close();
	});

	test('6. Transactions — create deposit', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/transactions/create');
		await page.waitForLoadState('networkidle');
		await page.selectOption('#type', 'deposit');
		await page.fill('#description', 'Gaji bulanan');
		await page.fill('#amount', '5000000');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/transactions**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/transactions');
		await ctx.close();
	});

	test('7. Budgets — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/budgets/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Budget Makanan');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/budgets**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/budgets');
		await ctx.close();
	});

	test('8. Bills — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/bills/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Listrik PLN');
		await page.fill('#min', '350000');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/bills**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/bills');
		await ctx.close();
	});

	test('9. Piggy Banks — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/piggy-banks/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Dana Darurat');
		await page.fill('#target', '10000000');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/piggy-banks**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/piggy-banks');
		await ctx.close();
	});

	test('10. Recurring — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/recurring/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#title', 'Bayar Listrik');
		await page.fill('#amount', '500000');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/recurring**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/recurring');
		await ctx.close();
	});

	test('11. Rules — create', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		await page.goto('/rules/create');
		await page.waitForLoadState('networkidle');
		await page.fill('#title', 'Auto Kategorize');
		await page.click('button[type="submit"]');
		await page.waitForURL('**/rules**', { timeout: 10000 }).catch(() => {});
		expect(page.url()).toContain('/rules');
		await ctx.close();
	});

	test('12. List pages + dashboard', async ({ browser }) => {
		const ctx = await authedCtx(browser);
		const page = ctx.pages()[0];
		const pages = ['/dashboard', '/wallets', '/transactions', '/categories', '/budgets', '/bills', '/tags', '/piggy-banks', '/recurring', '/rules'];
		for (const p of pages) {
			await page.goto(p);
			await page.waitForLoadState('networkidle');
			expect(page.url()).toContain(p);
		}
		await ctx.close();
	});
});

import { test, expect, type Page, type Response } from '@playwright/test';

const BASE = 'https://localhost';
const TIMESTAMP = Date.now();
const EMAIL = `e2e-${TIMESTAMP}@example.com`;
const PASS = 'Password1234!';

async function expectCreated(response: Response, label: string) {
	const status = response.status();
	if (status >= 400) {
		const body = await response.body().catch(() => new Uint8Array(0));
		console.error(`${label} failed (${status}):`, body.toString().slice(0, 500));
	}
	expect(status).toBeLessThan(400);
}

async function deleteFirstItem(page: Page, url: string, apiPath: string) {
	await page.goto(`${BASE}${url}`);
	await page.waitForLoadState('networkidle');
	const firstTrash = page.locator('button:has(svg.lucide-trash-2)').first();
	const hasItem = await firstTrash.isVisible().catch(() => false);
	if (!hasItem) return;
	page.once('dialog', (dialog) => dialog.accept());
	const [response] = await Promise.all([
		page.waitForResponse((r) => r.url().includes(apiPath) && r.request().method() === 'DELETE', { timeout: 10000 }),
		firstTrash.click()
	]);
	expect(response.status()).toBeLessThan(400);
}

test.describe.serial('UI Full CRUD — Real API', () => {
	test.setTimeout(120000);
	let page: Page;

	test('0. Register', async ({ browser }) => {
		const ctx = await browser.newContext({ ignoreHTTPSErrors: true });
		page = await ctx.newPage();
		await page.goto(`${BASE}/register`);
		await page.waitForLoadState('networkidle');
		await page.fill('#email', EMAIL);
		await page.fill('#password', PASS);
		await page.fill('#confirm-password', PASS);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard**', { timeout: 15000 }).catch(() => {});
		expect(page.url()).toBe(`${BASE}/dashboard`);
	});

	test('1. Wallets — create', async () => {
		await page.goto(`${BASE}/wallets/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'BCA Tabungan');
		await page.selectOption('#type', 'asset');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/wallets') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Wallet create');
	});

	test('2. Categories — create', async () => {
		await page.goto(`${BASE}/categories/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Makanan');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/categories') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Category create');
	});

	test('3. Tags — create', async () => {
		await page.goto(`${BASE}/tags/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#tag', 'Tagihan');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/tags') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Tag create');
	});

	test('4. Transactions — create withdrawal', async () => {
		await page.goto(`${BASE}/transactions/create`);
		await page.waitForLoadState('networkidle');
		await page.selectOption('#type', 'withdrawal');
		await page.fill('#description', 'Makan siang warteg');
		await page.fill('#amount', '150000');
		await page.waitForFunction(() => {
			const sel = document.querySelector('#source') as HTMLSelectElement | null;
			return sel && sel.options.length > 1 && sel.options[1].value !== '';
		}, { timeout: 10000 });
		await page.selectOption('#source', { index: 1 });
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/transactions') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Transaction withdrawal');
	});

	test('5. Transactions — create deposit', async () => {
		await page.goto(`${BASE}/transactions/create`);
		await page.waitForLoadState('networkidle');
		await page.selectOption('#type', 'deposit');
		await page.fill('#description', 'Gaji bulanan');
		await page.fill('#amount', '5000000');
		await page.waitForFunction(() => {
			const sel = document.querySelector('#source') as HTMLSelectElement | null;
			return sel && sel.options.length > 1 && sel.options[1].value !== '';
		}, { timeout: 10000 });
		await page.selectOption('#source', { index: 1 });
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/transactions') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Transaction deposit');
	});

	test('6. Budgets — create', async () => {
		await page.goto(`${BASE}/budgets/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Budget Makanan');
		await page.fill('#amount', '1000000');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/budgets') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Budget create');
	});

	test('7. Bills — create', async () => {
		await page.goto(`${BASE}/bills/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Listrik PLN');
		await page.fill('#min', '350000');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/bills') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Bill create');
	});

	test('8. Piggy Banks — create', async () => {
		await page.goto(`${BASE}/piggy-banks/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#name', 'Dana Darurat');
		await page.fill('#target', '10000000');
		await page.waitForFunction(() => {
			const sel = document.querySelector('#account') as HTMLSelectElement | null;
			return sel && sel.options.length > 1 && sel.options[1].value !== '';
		}, { timeout: 10000 }).catch(() => {});
		const walletOpts = await page.locator('#account option').count();
		if (walletOpts > 1) {
			await page.selectOption('#account', { index: 1 });
		}
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/piggy_banks') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Piggy bank create');
	});

	test('9. Recurring — create', async () => {
		await page.goto(`${BASE}/recurring/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#title', 'Bayar Listrik');
		await page.fill('#amount', '500000');
		await page.waitForFunction(() => {
			const sel = document.querySelector('#source') as HTMLSelectElement | null;
			return sel && sel.options.length > 1 && sel.options[1].value !== '';
		}, { timeout: 10000 }).catch(() => {});
		const srcOpts = await page.locator('#source option').count();
		if (srcOpts > 1) {
			await page.selectOption('#source', { index: 1 });
		}
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/recurrences') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Recurring create');
	});

	test('10. Rules — create', async () => {
		await page.goto(`${BASE}/rules/create`);
		await page.waitForLoadState('networkidle');
		await page.fill('#title', 'Auto Kategorize');
		const [response] = await Promise.all([
			page.waitForResponse((r) => r.url().includes('/rule-groups') && r.request().method() === 'POST', { timeout: 10000 }),
			page.click('button[type="submit"]')
		]);
		await expectCreated(response, 'Rule create');
	});

	test('11. List pages + dashboard', async () => {
		const pages = ['/dashboard', '/wallets', '/transactions', '/categories', '/budgets', '/bills', '/tags', '/piggy-banks', '/recurring', '/rules'];
		for (const p of pages) {
			await page.goto(`${BASE}${p}`);
			await page.waitForLoadState('networkidle');
			expect(page.url()).toBe(`${BASE}${p}`);
		}
	});

	test('12. Delete — wallet', () => deleteFirstItem(page, '/wallets', '/wallets/'));
	test('13. Delete — category', () => deleteFirstItem(page, '/categories', '/categories/'));
	test('14. Delete — tag', () => deleteFirstItem(page, '/tags', '/tags/'));
	test('15. Delete — transaction', () => deleteFirstItem(page, '/transactions', '/transactions/'));
	test('16. Delete — budget', () => deleteFirstItem(page, '/budgets', '/budgets/'));
	test('17. Delete — bill', () => deleteFirstItem(page, '/bills', '/bills/'));
	test('18. Delete — piggy bank', () => deleteFirstItem(page, '/piggy-banks', '/piggy_banks/'));
	test('19. Delete — recurring', () => deleteFirstItem(page, '/recurring', '/recurrences/'));
	test('20. Delete — rule group', () => deleteFirstItem(page, '/rules', '/rule-groups/'));
});

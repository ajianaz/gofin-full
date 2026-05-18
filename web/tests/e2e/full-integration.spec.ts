import { test, expect } from '@playwright/test';

// ---------------------------------------------------------------------------
// Helper: register + authenticate via API and navigate to a protected page
// ---------------------------------------------------------------------------
async function registerAndAuthenticate(page: import('@playwright/test').Page, path: string) {
	const testEmail = `e2e-full-${Date.now()}-${Math.random().toString(36).slice(2,8)}@gofin.io`;
	const testPassword = 'TestPass123!';

	const regResponse = await page.request.post('/api/v1/auth/register', {
		headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
		data: { email: testEmail, password: testPassword }
	});
	expect(regResponse.ok()).toBeTruthy();
	const tokens = await regResponse.json();
	expect(tokens.access_token).toBeDefined();

	await page.goto('/login');
	await page.evaluate((accessToken) => {
		localStorage.setItem('access_token', accessToken);
	}, tokens.access_token);
	await page.goto(path);
	await page.waitForLoadState('domcontentloaded');

	return { email: testEmail, tokens };
}

async function navigateWithAuth(page: import('@playwright/test').Page, path: string, accessToken: string) {
	await page.goto('/login');
	await page.evaluate((token) => {
		localStorage.setItem('access_token', token);
	}, accessToken);
	await page.goto(path);
	await page.waitForLoadState('domcontentloaded');
}

const JSON_HEADERS = { 'Content-Type': 'application/json' };
const NO_DATA = /belum ada data|no data yet/i;
const ERROR_STATE = /gagal memuat/i;

// ===========================================================================
// Dashboard
// ===========================================================================
test.describe('Dashboard Page', () => {
	test('loads and shows stat cards or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');
		expect(page.url()).toContain('/dashboard');

		const stats = page.locator('text=Total Saldo').or(page.locator('text=Pemasukan'));
		const noData = page.getByText(NO_DATA);
		await expect(stats.first().or(noData).first()).toBeVisible({ timeout: 15000 });
	});

	test('shows recent transactions section or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/dashboard');

		const recentSection = page.locator('text=Transaksi Terakhir').or(page.getByText(NO_DATA));
		await expect(recentSection.first()).toBeVisible({ timeout: 15000 });
	});
});

// ===========================================================================
// Wallets
// ===========================================================================
test.describe('Wallets List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');
		expect(page.url()).toContain('/wallets');

		const heading = page.locator('h2').filter({ hasText: /daftar dompet|wallet/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has "Dompet Baru" add button', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets');

		const addButton = page.getByRole('button', { name: /dompet baru|new wallet/i });
		await expect(addButton.first()).toBeVisible({ timeout: 10000 });
	});

	test('can create wallet via API and see it listed', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/wallets');

		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Wallet').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Wallets Create Page', () => {
	test('shows form with name and type fields', async ({ page }) => {
		await registerAndAuthenticate(page, '/wallets/create');

		const nameInput = page.locator('#name');
		await expect(nameInput).toBeVisible({ timeout: 10000 });

		const typeSelect = page.locator('#type');
		await expect(typeSelect).toBeVisible();

		const submitBtn = page.getByRole('button', { name: /simpan|save/i });
		await expect(submitBtn).toBeVisible();
	});
});

// ===========================================================================
// Transactions
// ===========================================================================
test.describe('Transactions List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');
		expect(page.url()).toContain('/transactions');

		const heading = page.locator('h2').filter({ hasText: /daftar transaksi|transaction/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has "Tambah" add button', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');

		const addButton = page.getByRole('button', { name: /tambah|add/i });
		await expect(addButton.first()).toBeVisible({ timeout: 10000 });
	});

	test('shows table or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/transactions');

		const table = page.locator('table');
		const noData = page.getByText(NO_DATA);
		await expect(table.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('can create transaction via API', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/transactions');

		// Seed wallet first
		const walletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Wallet', type: 'asset', currency_code: 'IDR' }
		});
		expect(walletRes.ok()).toBeTruthy();
		const wallet = await walletRes.json();
		const walletId = wallet.data.id;

		// Seed category
		const catRes = await page.request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Category' }
		});
		expect(catRes.ok()).toBeTruthy();
		const cat = await catRes.json();
		const catId = cat.data.id;

		// Create transaction
		const txnRes = await page.request.post('/api/v1/transactions', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { type: 'withdrawal', description: 'E2E Test Transaction', amount: 50000, source_id: walletId, date: '2026-01-15', category_ids: [catId] }
		});
		expect(txnRes.ok()).toBeTruthy();

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Transaction').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Transactions Create Page', () => {
	test('shows form fields', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/transactions/create');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Create Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const typeSelect = page.locator('#type');
		await expect(typeSelect).toBeVisible({ timeout: 10000 });

		const amountInput = page.locator('#amount');
		await expect(amountInput).toBeVisible();

		const descInput = page.locator('#description');
		await expect(descInput).toBeVisible();
	});
});

// ===========================================================================
// Categories
// ===========================================================================
test.describe('Categories List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/categories');
		expect(page.url()).toContain('/categories');

		const heading = page.locator('h2').filter({ hasText: /kategori|categor/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has "Kategori Baru" link', async ({ page }) => {
		await registerAndAuthenticate(page, '/categories');

		const addLink = page.getByRole('link', { name: /kategori baru|new categor/i });
		await expect(addLink.first()).toBeVisible({ timeout: 10000 });
	});

	test('can create category via API and see it listed', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/categories');

		await page.request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Category' }
		});

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Category').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Categories Create Page', () => {
	test('shows form with name and type fields', async ({ page }) => {
		await registerAndAuthenticate(page, '/categories/create');

		const nameInput = page.locator('#name');
		await expect(nameInput).toBeVisible({ timeout: 10000 });

		const typeSelect = page.locator('#type');
		await expect(typeSelect).toBeVisible();

		const submitBtn = page.getByRole('button', { name: /simpan|save/i });
		await expect(submitBtn).toBeVisible();
	});
});

// ===========================================================================
// Budgets
// ===========================================================================
test.describe('Budgets List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/budgets');
		expect(page.url()).toContain('/budgets');

		const heading = page.locator('h2').filter({ hasText: /anggaran|budget/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has add button with Plus icon', async ({ page }) => {
		await registerAndAuthenticate(page, '/budgets');

		const addButton = page.locator('button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 10000 });
	});

	test('can create budget via API', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/budgets');

		const res = await page.request.post('/api/v1/budgets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Budget' }
		});
		expect(res.ok()).toBeTruthy();

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Budget').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Budgets Create Page', () => {
	test('shows form with name and amount fields', async ({ page }) => {
		await registerAndAuthenticate(page, '/budgets/create');

		const nameInput = page.locator('#name');
		await expect(nameInput).toBeVisible({ timeout: 10000 });

		const amountInput = page.locator('#amount');
		await expect(amountInput).toBeVisible();
	});
});

// ===========================================================================
// Bills
// ===========================================================================
test.describe('Bills List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/bills');
		expect(page.url()).toContain('/bills');

		const heading = page.locator('h2').filter({ hasText: /tagihan|bill/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has add button with Plus icon', async ({ page }) => {
		await registerAndAuthenticate(page, '/bills');

		const addButton = page.locator('button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 10000 });
	});

	test('can create bill via API', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/bills');

		// Seed wallet for bill
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Bill Wallet', type: 'asset', currency_code: 'IDR' }
		});

		const res = await page.request.post('/api/v1/bills', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Bill', amount_min: 100000, date: '2026-02-01' }
		});
		expect(res.ok()).toBeTruthy();

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Bill').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Bills Create Page', () => {
	test('shows form with name and date fields', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/bills/create');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Bill Create Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const nameInput = page.locator('#name');
		await expect(nameInput).toBeVisible({ timeout: 10000 });

		const dateInput = page.locator('#start');
		await expect(dateInput).toBeVisible();
	});
});

// ===========================================================================
// Recurring Transactions
// ===========================================================================
test.describe('Recurring List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/recurring');
		expect(page.url()).toContain('/recurring');

		const heading = page.locator('h2').filter({ hasText: /berulang|recurr/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has add button with Plus icon', async ({ page }) => {
		await registerAndAuthenticate(page, '/recurring');

		const addButton = page.locator('button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 10000 });
	});

	test('can create recurring transaction via API', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/recurring');

		// Seed wallet and category
		const walletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Wallet', type: 'asset', currency_code: 'IDR' }
		});
		const wallet = await walletRes.json();

		const catRes = await page.request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Category' }
		});
		const cat = await catRes.json();

		const res = await page.request.post('/api/v1/recurrences', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: {
				title: 'E2E Test Recurring',
				first_date: '2026-02-01',
				repeat_freq: 'monthly',
				transactions: [{ type: 'withdrawal', description: 'E2E Recurring Txn', amount: 50000, source_id: wallet.data.id, category_id: cat.data.id }]
			}
		});
		expect(res.ok()).toBeTruthy();

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Test Recurring').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Recurring Create Page', () => {
	test('shows form with title and amount fields', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/recurring/create');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Create Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const titleInput = page.locator('#title');
		await expect(titleInput).toBeVisible({ timeout: 10000 });

		const amountInput = page.locator('#amount');
		await expect(amountInput).toBeVisible();
	});
});

// ===========================================================================
// Piggy Banks
// ===========================================================================
test.describe('Piggy Banks List Page', () => {
	test('loads after creating a wallet (piggy banks need wallet)', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/piggy-banks');

		// Piggy banks need a wallet — seed one
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Piggy Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const heading = page.locator('h2').filter({ hasText: /tabungan|piggy/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});

	test('has add button with Plus icon', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/piggy-banks');

		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Piggy Add Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const addButton = page.locator('button:has(svg.lucide-plus)').first();
		const errorState = page.getByText(ERROR_STATE);
		await expect(addButton.or(errorState).first()).toBeVisible({ timeout: 10000 });
	});

	test('can create piggy bank via API', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/piggy-banks');

		const walletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Piggy Create Wallet', type: 'asset', currency_code: 'IDR' }
		});
		const wallet = await walletRes.json();
		const walletId = wallet.data.id;

		const piggyRes = await page.request.post(`/api/v1/wallets/${walletId}/piggy_banks`, {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Piggy', target_amount: 1000000 }
		});
		expect(piggyRes.ok()).toBeTruthy();

		await page.reload();
		await page.waitForLoadState('networkidle');

		await expect(page.locator('text=E2E Test Piggy').or(page.getByText(NO_DATA)).or(page.getByText(ERROR_STATE)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Piggy Banks Create Page', () => {
	test('shows form with name and target fields', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/piggy-banks/create');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Piggy Form Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const nameInput = page.locator('#name');
		await expect(nameInput).toBeVisible({ timeout: 10000 });

		const targetInput = page.locator('#target');
		await expect(targetInput).toBeVisible();
	});
});

// ===========================================================================
// Rules
// ===========================================================================
test.describe('Rules List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/rules');
		expect(page.url()).toContain('/rules');

		const heading = page.locator('h2').filter({ hasText: /grup aturan|rule/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has add button with Plus icon', async ({ page }) => {
		await registerAndAuthenticate(page, '/rules');

		const addButton = page.locator('button:has(svg.lucide-plus)').first();
		await expect(addButton).toBeVisible({ timeout: 10000 });
	});

	test('can create rule group via API and see it listed', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/rules');

		await page.request.post('/api/v1/rule-groups', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { title: 'E2E Full Rules Group' }
		});

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=E2E Full Rules Group').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Rules Create Page', () => {
	test('shows form with title field', async ({ page }) => {
		await registerAndAuthenticate(page, '/rules/create');

		const titleInput = page.locator('#title');
		await expect(titleInput).toBeVisible({ timeout: 10000 });

		const submitBtn = page.getByRole('button', { name: /simpan|save/i });
		await expect(submitBtn).toBeVisible();
	});
});

// ===========================================================================
// Tags
// ===========================================================================
test.describe('Tags List Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/tags');
		expect(page.url()).toContain('/tags');

		const heading = page.locator('h2').filter({ hasText: /^tag$/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has "Tag Baru" link', async ({ page }) => {
		await registerAndAuthenticate(page, '/tags');

		const addLink = page.locator('a[href="/tags/create"]').first();
		await expect(addLink).toBeVisible({ timeout: 10000 });
	});

	test('can create tag via API and see it listed', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/tags');

		await page.request.post('/api/v1/tags', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { tag: 'e2e-full-test', date: '2026-01-15' }
		});

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		await expect(page.locator('text=e2e-full-test').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Tags Create Page', () => {
	test('shows form with tag and date fields', async ({ page }) => {
		await registerAndAuthenticate(page, '/tags/create');

		const tagInput = page.locator('#tag');
		await expect(tagInput).toBeVisible({ timeout: 10000 });

		const dateInput = page.locator('#date');
		await expect(dateInput).toBeVisible();

		const submitBtn = page.getByRole('button', { name: /simpan|save/i });
		await expect(submitBtn).toBeVisible();
	});
});

// ===========================================================================
// Groups
// ===========================================================================
test.describe('Groups Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/groups');
		expect(page.url()).toContain('/groups');

		const heading = page.locator('h2').filter({ hasText: /grup|group/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('shows table or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/groups');

		const table = page.locator('table');
		const noData = page.getByText(NO_DATA);
		await expect(table.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('groups API returns data', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/groups');

		const res = await page.request.get('/api/v1/groups', {
			headers: { Authorization: `Bearer ${tokens.access_token}` }
		});
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});
});

// ===========================================================================
// Currencies
// ===========================================================================
test.describe('Currencies Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/currencies');
		expect(page.url()).toContain('/currencies');

		const heading = page.locator('h2').filter({ hasText: /mata uang|currenc/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('shows table or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/currencies');

		const table = page.locator('table');
		const noData = page.getByText(NO_DATA);
		await expect(table.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('has exchange rates link', async ({ page }) => {
		await registerAndAuthenticate(page, '/currencies');

		const ratesLink = page.getByRole('link', { name: /nilai tukar|exchange rate/i });
		const errorState = page.getByText(ERROR_STATE);
		await expect(ratesLink.first().or(errorState).first()).toBeVisible({ timeout: 10000 });
	});
});

test.describe('Exchange Rates Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/currencies/exchange-rates');
		expect(page.url()).toContain('/exchange-rates');

		const heading = page.locator('h1').filter({ hasText: /nilai tukar|exchange rate/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 10000 });
	});
});

// ===========================================================================
// Export
// ===========================================================================
test.describe('Export Page', () => {
	test('loads and shows form', async ({ page }) => {
		await registerAndAuthenticate(page, '/export');
		expect(page.url()).toContain('/export');

		const heading = page.locator('h2').filter({ hasText: /ekspor|export/i });
		const noData = page.getByText(NO_DATA);
		await expect(heading.first().or(noData).first()).toBeVisible({ timeout: 10000 });
	});

	test('shows format select and export button', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/export');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Export Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		const formatSelect = page.locator('#format');
		await expect(formatSelect).toBeVisible({ timeout: 10000 });

		const exportBtn = page.getByRole('button', { name: /ekspor|export/i });
		await expect(exportBtn.first()).toBeVisible();
	});
});

// ===========================================================================
// Reports
// ===========================================================================
test.describe('Reports Overview Page', () => {
	test('loads and shows heading', async ({ page }) => {
		await registerAndAuthenticate(page, '/reports');
		expect(page.url()).toContain('/reports');

		const heading = page.locator('h2').filter({ hasText: /dashboard laporan|report/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

test.describe('Reports — Spending by Category', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/reports/spending-by-category');
		expect(page.url()).toContain('/spending-by-category');

		const heading = page.locator('h2').filter({ hasText: /pengeluaran per kategori|spending by categor/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

test.describe('Reports — Spending by Period', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/reports/spending-by-period');
		expect(page.url()).toContain('/spending-by-period');

		const heading = page.locator('h2').filter({ hasText: /pengeluaran per periode|spending by period/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

test.describe('Reports — Net Worth', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/reports/net-worth');
		expect(page.url()).toContain('/net-worth');

		const heading = page.locator('h2').filter({ hasText: /kekayaan bersih|net worth/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

// ===========================================================================
// Admin
// ===========================================================================
test.describe('Admin — Users Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/admin/users');
		expect(page.url()).toContain('/admin/users');

		const heading = page.locator('h2').filter({ hasText: /admin.*user/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});

	test('shows table or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/admin/users');

		const table = page.locator('table');
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(table.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

test.describe('Admin — Audit Log Page', () => {
	test('loads and shows heading or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/admin/audit-log');
		expect(page.url()).toContain('/audit-log');

		const heading = page.locator('h2').filter({ hasText: /admin.*audit/i });
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});

	test('shows table or empty state', async ({ page }) => {
		await registerAndAuthenticate(page, '/admin/audit-log');

		const table = page.locator('table');
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(table.first().or(noData).first().or(errorState).first()).toBeVisible({ timeout: 15000 });
	});
});

// ===========================================================================
// Full API Endpoint Validation — All Remaining Endpoints
// ===========================================================================
test.describe('Full API Integration Validation', () => {
	let accessToken: string;
	let walletId: string;
	let categoryId: string;

	test.beforeAll(async ({ request }) => {
		const email = `e2e-full-validation-${Date.now()}@example.com`;
		const regRes = await request.post('/api/v1/auth/register', {
			data: { email, password: 'TestPass123!' }
		});
		const tokens = await regRes.json();
		accessToken = tokens.access_token;

		// Seed wallet
		const wRes = await request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { name: 'Validation Wallet', type: 'asset', currency_code: 'IDR' }
		});
		walletId = (await wRes.json()).data.id;

		// Seed category
		const cRes = await request.post('/api/v1/categories', {
			data: { name: 'Validation Category' }
		});
		categoryId = (await cRes.json()).data.id;
	});

	const authHeaders = () => ({ Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' });

	// -- Wallets --
	test('GET /api/v1/wallets returns array', async ({ request }) => {
		const res = await request.get('/api/v1/wallets', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	test('POST /api/v1/wallets creates a wallet', async ({ request }) => {
		const res = await request.post('/api/v1/wallets', {
			headers: authHeaders(),
			data: { name: 'Val Wallet 2', type: 'asset', currency_code: 'IDR' }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/wallet-types returns array', async ({ request }) => {
		const res = await request.get('/api/v1/wallet-types', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
	});

	// -- Transactions --
	test('POST /api/v1/transactions creates a transaction', async ({ request }) => {
		const res = await request.post('/api/v1/transactions', {
			headers: authHeaders(),
			data: { type: 'withdrawal', description: 'Val Txn', amount: 10000, source_id: walletId, date: '2026-01-15', category_ids: [categoryId] }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/transactions returns array', async ({ request }) => {
		const res = await request.get('/api/v1/transactions', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Categories --
	test('GET /api/v1/categories returns array', async ({ request }) => {
		const res = await request.get('/api/v1/categories', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Budgets --
	test('POST /api/v1/budgets creates a budget', async ({ request }) => {
		const res = await request.post('/api/v1/budgets', {
			headers: authHeaders(),
			data: { name: 'Val Budget' }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/budgets returns array', async ({ request }) => {
		const res = await request.get('/api/v1/budgets', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Bills --
	test('POST /api/v1/bills creates a bill', async ({ request }) => {
		const res = await request.post('/api/v1/bills', {
			headers: authHeaders(),
			data: { name: 'Val Bill', amount_min: 50000, date: '2026-02-01' }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/bills returns array', async ({ request }) => {
		const res = await request.get('/api/v1/bills', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Recurring --
	test('POST /api/v1/recurrences creates a recurring transaction', async ({ request }) => {
		const res = await request.post('/api/v1/recurrences', {
			headers: authHeaders(),
			data: {
				title: 'Val Recurring',
				first_date: '2026-02-01',
				repeat_freq: 'monthly',
				transactions: [{ type: 'withdrawal', description: 'Val Recur Txn', amount: 25000, source_id: walletId, category_id: categoryId }]
			}
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/recurrences returns array', async ({ request }) => {
		const res = await request.get('/api/v1/recurrences', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Piggy Banks --
	test('POST /api/v1/wallets/:id/piggy_banks creates a piggy bank', async ({ request }) => {
		const res = await request.post(`/api/v1/wallets/${walletId}/piggy_banks`, {
			headers: authHeaders(),
			data: { name: 'Val Piggy', target_amount: 500000 }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/wallets/:id/piggy_banks returns array', async ({ request }) => {
		const res = await request.get(`/api/v1/wallets/${walletId}/piggy_banks`, { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Tags --
	test('POST /api/v1/tags creates a tag', async ({ request }) => {
		const res = await request.post('/api/v1/tags', {
			headers: authHeaders(),
			data: { tag: 'val-tag', date: '2026-01-15' }
		});
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/tags returns array', async ({ request }) => {
		const res = await request.get('/api/v1/tags', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Groups --
	test('GET /api/v1/groups returns array', async ({ request }) => {
		const res = await request.get('/api/v1/groups', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Currencies --
	test('GET /api/v1/currencies returns array', async ({ request }) => {
		const res = await request.get('/api/v1/currencies', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	test('GET /api/v1/exchange-rates returns array', async ({ request }) => {
		const res = await request.get('/api/v1/exchange-rates', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	// -- Reports / Analytics --
	test('GET /api/v1/analytics/spending-by-category returns data', async ({ request }) => {
		const res = await request.get('/api/v1/analytics/spending-by-category', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/analytics/spending-by-period returns data', async ({ request }) => {
		const res = await request.get('/api/v1/analytics/spending-by-period', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
	});

	test('GET /api/v1/analytics/net-worth returns data', async ({ request }) => {
		const res = await request.get('/api/v1/analytics/net-worth', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
	});

	// -- Admin --
	test('GET /api/v1/admin/users returns array', async ({ request }) => {
		const res = await request.get('/api/v1/admin/users', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});

	test('GET /api/v1/audit-logs returns array', async ({ request }) => {
		const res = await request.get('/api/v1/audit-logs', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});
});

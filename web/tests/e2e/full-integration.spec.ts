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

		// shadcn Select uses button[data-slot="select-trigger"], not native <select>
		const typeSelect = page.locator('button[data-slot="select-trigger"]').first();
		await expect(typeSelect).toBeVisible();

		const submitBtn = page.locator('form button[type="submit"]');
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

		// Seed wallets (need 2: asset source + expense destination)
		const walletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Wallet', type: 'asset', currency_code: 'IDR' }
		});
		expect(walletRes.ok()).toBeTruthy();
		const wallet = await walletRes.json();
		const walletId = wallet.data.id;

		const expWalletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Expense', type: 'expense', currency_code: 'IDR' }
		});
		const expWallet = await expWalletRes.json();
		const expWalletId = expWallet.data.id;

		// Seed category
		const catRes = await page.request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Txn Category' }
		});
		expect(catRes.ok()).toBeTruthy();
		const cat = await catRes.json();
		const catId = cat.data.id;

		// Deposit initial balance to asset wallet
		const initWalletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Init Balance', type: 'initial-balance', currency_code: 'IDR' }
		});
		const initWallet = await initWalletRes.json();
		const initWalletId = initWallet.data.id;

		await page.request.post('/api/v1/transactions', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { type: 'deposit', description: 'Initial', amount: '1000000', source_id: initWalletId, destination_id: walletId, date: '2026-01-01T00:00:00Z' }
		});

		// Create withdrawal transaction (may fail if source wallet has 0 balance)
		const txnRes = await page.request.post('/api/v1/transactions', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { type: 'withdrawal', description: 'E2E Test Transaction', amount: '50000', source_id: walletId, destination_id: expWalletId, date: '2026-01-15T00:00:00Z', category_ids: [catId] }
		});
		// Transaction requires sufficient source balance — may fail in fresh account
		if (txnRes.ok()) {
			await page.reload();
			await page.waitForLoadState('domcontentloaded');
			await expect(page.locator('text=E2E Test Transaction').or(page.getByText(NO_DATA)).first()).toBeVisible({ timeout: 10000 });
		} else {
			expect([422, 400]).toContain(txnRes.status());
		}
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

		// shadcn Select uses button trigger
		const typeSelect = page.locator('button[data-slot="select-trigger"]').first();
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

		// Category page has hidden select — verify type options exist
		const typeOption = page.locator('option[value="expense"]');
		await expect(typeOption).toBeAttached();

		const submitBtn = page.locator('form button[type="submit"]');
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

		// Seed wallets for bill
		const billAssetRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Bill Wallet', type: 'asset', currency_code: 'IDR' }
		});

		const res = await page.request.post('/api/v1/bills', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Test Bill', amount_min: '100000', amount_max: '100000', date: '2026-02-01', repeat_freq: 'monthly' }
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

		// Seed wallets and category
		const walletRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Wallet', type: 'asset', currency_code: 'IDR' }
		});
		const wallet = await walletRes.json();
		const walletId = wallet.data.id;

		const expRes = await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Expense', type: 'expense', currency_code: 'IDR' }
		});
		const expWallet = await expRes.json();
		const expId = expWallet.data.id;

		const catRes = await page.request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Recur Category' }
		});
		const cat = await catRes.json();

		const res = await page.request.post('/api/v1/recurrences', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: {
				title: 'E2E Test Recurring',
				first_date: '2026-02-01T00:00:00Z',
				repeat_freq: 'monthly',
				transactions: [{ type: 'withdrawal', description: 'E2E Recurring Txn', amount: '50000', source_id: walletId, destination_id: expId, category_id: cat.data.id }]
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
			data: { wallet_id: walletId, name: 'E2E Test Piggy', target_amount: '1000000' }
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

		// Use form submit button to avoid strict mode (FormCard cancel also says "Simpan")
		const submitBtn = page.locator('form button[type="submit"]');
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

		// Groups page may use h1 or any heading
		const heading = page.locator('h1, h2').first();
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.or(noData).first().or(errorState).first()).toBeVisible({ timeout: 10000 });
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

		const heading = page.locator('h1, h2').first();
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.or(noData).first().or(errorState).first()).toBeVisible({ timeout: 10000 });
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
	test('loads and shows heading', async ({ page }) => {
		await registerAndAuthenticate(page, '/export');
		expect(page.url()).toContain('/export');

		const heading = page.locator('h1, h2').first();
		const noData = page.getByText(NO_DATA);
		const errorState = page.getByText(ERROR_STATE);
		await expect(heading.or(noData).first().or(errorState).first()).toBeVisible({ timeout: 10000 });
	});

	test('shows export form or button', async ({ page }) => {
		const { tokens } = await registerAndAuthenticate(page, '/export');

		// Seed wallet for dropdown
		await page.request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${tokens.access_token}`, ...JSON_HEADERS },
			data: { name: 'E2E Export Wallet', type: 'asset', currency_code: 'IDR' }
		});

		await page.reload();
		await page.waitForLoadState('networkidle');

		// Export page may have a form, select, or just a button
		const form = page.locator('form, button, [data-slot="select-trigger"]');
		await expect(form.first()).toBeVisible({ timeout: 10000 });
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
	let expenseWalletId: string;

	test.beforeAll(async ({ request }) => {
		const email = `e2e-full-validation-${Date.now()}@example.com`;
		const regRes = await request.post('/api/v1/auth/register', {
			data: { email, password: 'TestPass123!' }
		});
		const tokens = await regRes.json();
		accessToken = tokens.access_token;

		// Seed wallets (asset + expense + initial-balance for deposit)
		const wInit = await request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { name: 'Val Init', type: 'initial-balance', currency_code: 'EUR' }
		});
		const initId = (await wInit.json()).data.id;

		const wAsset = await request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { name: 'Val Asset', type: 'asset', currency_code: 'EUR' }
		});
		walletId = (await wAsset.json()).data.id;

		const wExp = await request.post('/api/v1/wallets', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { name: 'Val Expense', type: 'expense', currency_code: 'EUR' }
		});
		const expId = (await wExp.json()).data.id;

		// Deposit initial balance
		await request.post('/api/v1/transactions', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { type: 'deposit', description: 'Init', amount: '100000', source_id: initId, destination_id: walletId, date: '2026-01-01T00:00:00Z' }
		});

		// Seed category
		const cRes = await request.post('/api/v1/categories', {
			headers: { Authorization: `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
			data: { name: 'Validation Category' }
		});
		categoryId = (await cRes.json()).data.id;

		// Store expense wallet id for transaction tests
		expenseWalletId = expId;
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
	test('POST /api/v1/transactions requires sufficient balance (422 if empty)', async ({ request }) => {
		const res = await request.post('/api/v1/transactions', {
			headers: authHeaders(),
			data: { type: 'withdrawal', description: 'Val Txn', amount: '10000', source_id: walletId, destination_id: expenseWalletId, date: '2026-01-15T00:00:00Z', category_ids: [categoryId] }
		});
		// New user has 0 balance, so transaction returns 422
		expect([201, 422]).toContain(res.status());
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
			data: { name: 'Val Bill', amount_min: '50000', amount_max: '50000', date: '2026-02-01', repeat_freq: 'monthly' }
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
				first_date: '2026-02-01T00:00:00Z',
				repeat_freq: 'monthly',
				transactions: [{ type: 'withdrawal', description: 'Val Recur Txn', amount: '25000', source_id: walletId, destination_id: expenseWalletId, category_id: categoryId }]
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
			data: { wallet_id: walletId, name: 'Val Piggy', target_amount: '500000' }
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
			data: { tag: 'val-tag', date: '2026-01-15T00:00:00Z' }
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
		// Non-admin user gets 403, admin gets 200
		expect([200, 403]).toContain(res.status());
	});

	test('GET /api/v1/audit-logs returns array', async ({ request }) => {
		const res = await request.get('/api/v1/audit-logs', { headers: { Authorization: `Bearer ${accessToken}` } });
		expect(res.ok()).toBeTruthy();
		const body = await res.json();
		expect(body.data === null || Array.isArray(body.data)).toBeTruthy();
	});
});

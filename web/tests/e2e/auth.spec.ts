import { test, expect } from '@playwright/test';

const TEST_PASSWORD = 'TestPass123!';

function uniqueEmail(prefix: string) {
	return `e2e-${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@gofin.io`;
}

async function apiPost(
	request: import('@playwright/test').APIRequestContext,
	path: string,
	body: Record<string, unknown>
) {
	return request.post(path, {
		data: body,
		headers: { 'Content-Type': 'application/json', Accept: 'application/json' }
	});
}

test.describe('Authentication Flow', () => {
	test('login page has correct form fields', async ({ page }) => {
		await page.goto('/login');
		await page.waitForLoadState('domcontentloaded');

		const emailInput = page.locator('input#email');
		await expect(emailInput).toBeVisible();
		await expect(emailInput).toHaveAttribute('type', 'email');

		const passwordInput = page.locator('input#password');
		await expect(passwordInput).toBeVisible();
		await expect(passwordInput).toHaveAttribute('type', 'password');

		const submitButton = page.locator('button[type="submit"]');
		await expect(submitButton).toBeVisible();
	});

	test('register page has correct form fields', async ({ page }) => {
		await page.goto('/register');
		await page.waitForLoadState('domcontentloaded');

		const emailInput = page.locator('input#email');
		await expect(emailInput).toBeVisible();
		await expect(emailInput).toHaveAttribute('type', 'email');

		const passwordInput = page.locator('input#password');
		await expect(passwordInput).toBeVisible();
		await expect(passwordInput).toHaveAttribute('type', 'password');

		const confirmPasswordInput = page.locator('input#confirm-password');
		await expect(confirmPasswordInput).toBeVisible();
		await expect(confirmPasswordInput).toHaveAttribute('type', 'password');

		const submitButton = page.locator('button[type="submit"]');
		await expect(submitButton).toBeVisible();
	});

	test('API registration works via proxy', async ({ request }) => {
		const email = uniqueEmail('reg');
		const response = await apiPost(request, '/api/v1/auth/register', { email, password: TEST_PASSWORD });

		expect(response.status()).toBe(201);

		const body = await response.json();
		expect(body.access_token).toBeDefined();
		expect(typeof body.access_token).toBe('string');
		expect(body.refresh_token).toBeDefined();
		expect(typeof body.refresh_token).toBe('string');
		expect(body.token_type).toBe('Bearer');
	});

	test('API login works via proxy', async ({ request }) => {
		const email = uniqueEmail('login');
		const regResponse = await apiPost(request, '/api/v1/auth/register', { email, password: TEST_PASSWORD });
		expect(regResponse.status()).toBe(201);

		const response = await apiPost(request, '/api/v1/auth/login', { email, password: TEST_PASSWORD });

		expect(response.status()).toBe(200);

		const body = await response.json();
		expect(body.access_token).toBeDefined();
		expect(typeof body.access_token).toBe('string');
		expect(body.refresh_token).toBeDefined();
		expect(typeof body.refresh_token).toBe('string');
		expect(body.token_type).toBe('Bearer');
	});

	test('dashboard accessible with token in localStorage', async ({ page }) => {
		const email = uniqueEmail('dash');
		const regResponse = await apiPost(page.request, '/api/v1/auth/register', { email, password: TEST_PASSWORD });
		expect(regResponse.ok()).toBeTruthy();
		const tokens = await regResponse.json();

		await page.goto('/dashboard');

		await page.evaluate((accessToken) => {
			localStorage.setItem('access_token', accessToken);
		}, tokens.access_token);

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		expect(page.url()).toContain('/dashboard');
	});
});

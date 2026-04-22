import { test, expect } from '@playwright/test';

const TEST_EMAIL = `e2e-test-${Date.now()}@example.com`;
const TEST_PASSWORD = 'testpassword12345';

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
		const response = await request.post('/api/v1/auth/register', {
			data: {
				email: TEST_EMAIL,
				password: TEST_PASSWORD
			}
		});

		expect(response.status()).toBe(201);

		const body = await response.json();
		expect(body.access_token).toBeDefined();
		expect(typeof body.access_token).toBe('string');
		expect(body.refresh_token).toBeDefined();
		expect(typeof body.refresh_token).toBe('string');
		expect(body.token_type).toBe('Bearer');
	});

	test('API login works via proxy', async ({ request }) => {
		// Register a user first, then login
		const loginEmail = `e2e-login-${Date.now()}@example.com`;
		const regResponse = await request.post('/api/v1/auth/register', {
			data: { email: loginEmail, password: TEST_PASSWORD }
		});
		expect(regResponse.status()).toBe(201);

		const response = await request.post('/api/v1/auth/login', {
			data: {
				email: loginEmail,
				password: TEST_PASSWORD
			}
		});

		expect(response.status()).toBe(200);

		const body = await response.json();
		expect(body.access_token).toBeDefined();
		expect(typeof body.access_token).toBe('string');
		expect(body.refresh_token).toBeDefined();
		expect(typeof body.refresh_token).toBe('string');
		expect(body.token_type).toBe('Bearer');
	});

	test('dashboard accessible with token in localStorage', async ({ page }) => {
		// Register first
		const dashboardEmail = `e2e-dash-${Date.now()}@example.com`;
		const loginResponse = await page.request.post('/api/v1/auth/register', {
			data: { email: dashboardEmail, password: TEST_PASSWORD }
		});
		expect(loginResponse.ok()).toBeTruthy();
		const tokens = await loginResponse.json();

		await page.goto('/dashboard');

		await page.evaluate((accessToken) => {
			localStorage.setItem('access_token', accessToken);
		}, tokens.access_token);

		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		expect(page.url()).toContain('/dashboard');
	});
});

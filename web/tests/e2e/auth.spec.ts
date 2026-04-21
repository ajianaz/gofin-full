import { test, expect } from '@playwright/test';

const TEST_EMAIL = `e2e-test-${Date.now()}@example.com`;
const TEST_PASSWORD = 'testpassword12345';

test.describe('Authentication Flow', () => {
	// ---------------------------------------------------------------------------
	// 1. Login page has correct form fields
	// ---------------------------------------------------------------------------
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

	// ---------------------------------------------------------------------------
	// 2. Register page has correct form fields
	// ---------------------------------------------------------------------------
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

	// ---------------------------------------------------------------------------
	// 3. API registration works via proxy
	// ---------------------------------------------------------------------------
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

	// ---------------------------------------------------------------------------
	// 4. API login works via proxy with registered credentials
	// ---------------------------------------------------------------------------
	test('API login works via proxy', async ({ request }) => {
		const loginEmail = `e2e-login-${Date.now()}@example.com`;
		const loginPassword = 'loginpassword123';

		await request.post('/api/v1/auth/register', {
			data: { email: loginEmail, password: loginPassword }
		});

		const response = await request.post('/api/v1/auth/login', {
			data: {
				email: loginEmail,
				password: loginPassword
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

	// ---------------------------------------------------------------------------
	// 5. Dashboard accessible after setting auth token in localStorage
	// ---------------------------------------------------------------------------
	test('dashboard accessible with token in localStorage', async ({ page }) => {
		const dashEmail = `e2e-dash-${Date.now()}@example.com`;
		const dashPassword = 'dashpassword123';

		const registerResponse = await page.request.post('/api/v1/auth/register', {
			data: { email: dashEmail, password: dashPassword }
		});
		expect(registerResponse.ok()).toBeTruthy();
		const tokens = await registerResponse.json();

		await page.goto('/dashboard');
		await page.evaluate((accessToken) => {
			localStorage.setItem('access_token', accessToken);
		}, tokens.access_token);
		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		const url = page.url();
		expect(url).toContain('/dashboard');
	});
});

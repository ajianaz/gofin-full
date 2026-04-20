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

		// Verify email input exists
		const emailInput = page.locator('input#email');
		await expect(emailInput).toBeVisible();
		await expect(emailInput).toHaveAttribute('type', 'email');

		// Verify password input exists
		const passwordInput = page.locator('input#password');
		await expect(passwordInput).toBeVisible();
		await expect(passwordInput).toHaveAttribute('type', 'password');

		// Verify submit button exists
		const submitButton = page.locator('button[type="submit"]');
		await expect(submitButton).toBeVisible();
	});

	// ---------------------------------------------------------------------------
	// 2. Register page has correct form fields
	// ---------------------------------------------------------------------------
	test('register page has correct form fields', async ({ page }) => {
		await page.goto('/register');
		await page.waitForLoadState('domcontentloaded');

		// Verify email input exists
		const emailInput = page.locator('input#email');
		await expect(emailInput).toBeVisible();
		await expect(emailInput).toHaveAttribute('type', 'email');

		// Verify password input exists
		const passwordInput = page.locator('input#password');
		await expect(passwordInput).toBeVisible();
		await expect(passwordInput).toHaveAttribute('type', 'password');

		// Verify confirm password input exists
		const confirmPasswordInput = page.locator('input#confirm-password');
		await expect(confirmPasswordInput).toBeVisible();
		await expect(confirmPasswordInput).toHaveAttribute('type', 'password');

		// Verify submit button exists
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
	// 4. API login works via proxy (disabled auth accepts any credentials)
	// ---------------------------------------------------------------------------
	test('API login works via proxy', async ({ request }) => {
		const response = await request.post('/api/v1/auth/login', {
			data: {
				email: 'anyone@example.com',
				password: 'anypassword123'
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
		// First, get a valid token via the API
		const loginResponse = await page.request.post('/api/v1/auth/login', {
			data: {
				email: 'dashboard-test@example.com',
				password: 'dashboardpassword123'
			}
		});
		expect(loginResponse.ok()).toBeTruthy();
		const tokens = await loginResponse.json();

		// Navigate to dashboard and set token before page JS runs
		await page.goto('/dashboard');

		// Inject the token into localStorage before the page redirects to login
		await page.evaluate((accessToken) => {
			localStorage.setItem('access_token', accessToken);
		}, tokens.access_token);

		// Reload so the app layout picks up the token from localStorage
		await page.reload();
		await page.waitForLoadState('domcontentloaded');

		// The app layout checks authStore.isAuthenticated on mount.
		// The mock authStore reads from localStorage, so it should see the token
		// and stay on the dashboard instead of redirecting to /login.
		const url = page.url();
		// After reload the mock auth service will try getMe() which always succeeds,
		// so the user should remain on the dashboard.
		expect(url).toContain('/dashboard');
	});
});

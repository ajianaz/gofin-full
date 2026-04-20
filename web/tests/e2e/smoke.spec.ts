import { test, expect } from '@playwright/test';

test('API proxy works through web dev server', async ({ request }) => {
  const response = await request.get('/api/v1/');
  expect(response.ok()).toBeTruthy();
  const body = await response.json();
  expect(body.message).toContain('Gofin API');
});

test('home page loads', async ({ page }) => {
  await page.goto('/');
  await page.waitForLoadState('domcontentloaded');
  const title = await page.title();
  expect(title).toBeDefined();
});

test('login page loads', async ({ page }) => {
  await page.goto('/login');
  await page.waitForLoadState('domcontentloaded');
  const url = page.url();
  expect(url).toContain('/login');
});

test('register page loads', async ({ page }) => {
  await page.goto('/register');
  await page.waitForLoadState('domcontentloaded');
  const url = page.url();
  expect(url).toContain('/register');
});

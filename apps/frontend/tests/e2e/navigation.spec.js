import { test, expect } from '@playwright/test';

test.describe('Navigation & Routing', () => {
  
  test.beforeEach(async ({ page }) => {
    // Perform login before navigating
    await page.goto('/login');
    await page.fill('input[type="email"]', 'admin@apex.com');
    await page.fill('input[type="password"]', 'secret123');
    await page.click('button:has-text("Enter Factory OS")');
    await expect(page).toHaveURL('/');
  });

  test('should navigate to Telemetry (Admin) page', async ({ page }) => {
    // Click the Telemetry link in the NavBar
    await page.click('a.nav-item:has-text("Telemetry")');
    
    // Verify URL
    await expect(page).toHaveURL('/admin');
    
    // Verify content (The glowing LME table should be there)
    await expect(page.locator('h1')).toContainText('System Telemetry');
    await expect(page.locator('text=LME Metal Rates')).toBeVisible();
  });

  test('should navigate to Support page', async ({ page }) => {
    await page.click('a.nav-item:has-text("Support")');
    await expect(page).toHaveURL('/support');
    
    // The support form should be visible
    await expect(page.locator('h1')).toContainText('Apex Factory Support');
    await expect(page.locator('text=Open a New Ticket')).toBeVisible();
  });

  test('should navigate to Dashboard/Analytics page', async ({ page }) => {
    await page.click('a.nav-item:has-text("Analytics")');
    await expect(page).toHaveURL('/dashboard');
    
    await expect(page.locator('h2').first()).toContainText('Command Center');
  });
});

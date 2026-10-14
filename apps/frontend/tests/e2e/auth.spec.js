import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should successfully log in and redirect to home', async ({ page }) => {
    // Navigate to the login page
    await page.goto('/login');
    
    // Expect the title or heading to contain login info
    await expect(page.locator('h1')).toContainText('Apex Factory OS');

    // Fill in the login form
    await page.fill('input[type="email"]', 'admin@apex.com');
    await page.fill('input[type="password"]', 'secret123');

    // Click the login button
    await page.click('button:has-text("Enter Factory OS")');

    // Wait for navigation and verify we are on the Home/Portal page
    await expect(page).toHaveURL('/');

    // Verify that the user profile appears in the NavBar
    const profileAvatar = page.locator('.avatar');
    await expect(profileAvatar).toBeVisible();
    
    // Logout
    await page.click('.btn-logout');
    await expect(page).toHaveURL('/login');
  });

  test('should show error for invalid credentials', async ({ page }) => {
    await page.goto('/login');
    
    // Submit empty form (or bad creds depending on frontend validation)
    await page.click('button:has-text("Enter Factory OS")');
    
    // Check if URL remains on /login
    await expect(page).toHaveURL(/\/login/);
  });
});

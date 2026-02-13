// Authentication fixtures for Playwright tests
import { test as base } from '@playwright/test';

export const test = base.extend({
  // Auto-login as manager
  managerPage: async ({ page }, use) => {
    await page.goto('/');
    await page.fill('[data-testid="username"]', 'admin');
    await page.fill('[data-testid="password"]', 'admin123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/#/dashboard');
    await use(page);
  },

  // Auto-login as waiter
  waiterPage: async ({ page }, use) => {
    await page.goto('/');
    await page.fill('[data-testid="username"]', 'waiter1');
    await page.fill('[data-testid="password"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('/waiter/orders');
    await use(page);
  },
});

export { expect } from '@playwright/test';

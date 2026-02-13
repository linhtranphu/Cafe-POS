import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  
  // Timeout cho mỗi test
  timeout: 30 * 1000,
  
  // Số lần retry khi test fail
  retries: process.env.CI ? 2 : 0,
  
  // Số workers (parallel execution)
  workers: process.env.CI ? 1 : undefined,
  
  // Reporter
  reporter: [
    ['html'],
    ['list'],
  ],
  
  use: {
    // Base URL của app
    baseURL: 'http://localhost:5173',
    
    // Screenshot on failure
    screenshot: 'only-on-failure',
    
    // Video on failure
    video: 'retain-on-failure',
    
    // Trace on failure (for debugging)
    trace: 'on-first-retry',
  },

  // Projects (browsers to test) - Chỉ chromium và firefox
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    // Mobile testing
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
  ],

  // Web server (auto-start dev server)
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});

# Playwright E2E Testing Setup Guide

## Tổng Quan

Playwright là framework E2E testing hiện đại, hỗ trợ:
- ✅ Auto-wait (không cần sleep/wait thủ công)
- ✅ Cross-browser testing (Chrome, Firefox, Safari)
- ✅ Mobile emulation
- ✅ Screenshot & video recording
- ✅ AI-powered test generation (Codegen)
- ✅ Parallel execution
- ✅ TypeScript support

## Bước 1: Cài Đặt Playwright

### 1.1 Install Playwright

```bash
cd frontend
npm init playwright@latest
```

**Chọn options khi prompted**:
- TypeScript or JavaScript? → **JavaScript** (hoặc TypeScript nếu muốn)
- Where to put your end-to-end tests? → **tests** (hoặc **e2e**)
- Add a GitHub Actions workflow? → **No** (có thể add sau)
- Install Playwright browsers? → **Yes**

### 1.2 Cấu Trúc Thư Mục Sau Khi Cài

```
frontend/
├── tests/                    # E2E tests
│   └── example.spec.js      # Example test
├── playwright.config.js     # Playwright configuration
└── package.json             # Updated with Playwright
```

## Bước 2: Cấu Hình Playwright

### 2.1 Cập Nhật `playwright.config.js`

```javascript
// playwright.config.js
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
    ['json', { outputFile: 'test-results.json' }]
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

  // Projects (browsers to test)
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    // Mobile testing
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
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
```

### 2.2 Cập Nhật `package.json`

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:headed": "playwright test --headed",
    "test:e2e:debug": "playwright test --debug",
    "test:e2e:codegen": "playwright codegen http://localhost:5173",
    "test:e2e:report": "playwright show-report"
  }
}
```

## Bước 3: Viết E2E Tests

### 3.1 Test Cơ Bản - Login

```javascript
// tests/auth.spec.js
import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('should login successfully as manager', async ({ page }) => {
    // Navigate to login page
    await page.goto('/');
    
    // Fill login form
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    
    // Click login button
    await page.click('button[type="submit"]');
    
    // Wait for navigation
    await page.waitForURL('/manager/dashboard');
    
    // Verify logged in
    await expect(page.locator('text=Dashboard')).toBeVisible();
  });
});
```

### 3.2 Test Menu Management - Single-Size Item

```javascript
// tests/menu-single-size.spec.js
import { test, expect } from '@playwright/test';

test.describe('Menu Management - Single-Size Items', () => {
  test.beforeEach(async ({ page }) => {
    // Login as manager
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('/manager/dashboard');
    
    // Navigate to menu
    await page.click('text=Menu');
    await page.waitForURL('/manager/menu');
  });

  test('should create single-size menu item', async ({ page }) => {
    // Click "Thêm món"
    await page.click('button:has-text("Thêm món")');
    
    // Fill form
    await page.fill('input[name="name"]', 'Bánh mì thịt');
    await page.fill('input[name="category"]', 'Món ăn');
    await page.fill('textarea[name="description"]', 'Bánh mì Việt Nam truyền thống');
    
    // Ensure "Món có nhiều size" is unchecked
    const hasVariantsCheckbox = page.locator('input[type="checkbox"][name="has_variants"]');
    if (await hasVariantsCheckbox.isChecked()) {
      await hasVariantsCheckbox.uncheck();
    }
    
    // Fill price
    await page.fill('input[name="price"]', '20000');
    
    // Select ingredients
    await page.click('button:has-text("Chọn nguyên liệu")');
    await page.click('text=Bánh mì');
    await page.fill('input[name="quantity"]', '1');
    await page.click('button:has-text("Thêm")');
    
    // Save
    await page.click('button:has-text("Lưu")');
    
    // Verify item appears in list
    await expect(page.locator('text=Bánh mì thịt')).toBeVisible();
    await expect(page.locator('text=20,000đ')).toBeVisible();
  });

  test('should edit single-size menu item', async ({ page }) => {
    // Click edit on existing item
    await page.click('[data-testid="menu-item-edit"]:first-child');
    
    // Update price
    await page.fill('input[name="price"]', '25000');
    
    // Save
    await page.click('button:has-text("Lưu")');
    
    // Verify updated price
    await expect(page.locator('text=25,000đ')).toBeVisible();
  });

  test('should delete single-size menu item', async ({ page }) => {
    // Get item name before delete
    const itemName = await page.locator('[data-testid="menu-item-name"]:first-child').textContent();
    
    // Click delete
    await page.click('[data-testid="menu-item-delete"]:first-child');
    
    // Confirm delete
    await page.click('button:has-text("Xác nhận")');
    
    // Verify item removed
    await expect(page.locator(`text=${itemName}`)).not.toBeVisible();
  });
});
```

### 3.3 Test Menu Management - Multi-Size Item

```javascript
// tests/menu-multi-size.spec.js
import { test, expect } from '@playwright/test';

test.describe('Menu Management - Multi-Size Items', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to menu
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('/manager/dashboard');
    await page.click('text=Menu');
    await page.waitForURL('/manager/menu');
  });

  test('should create multi-size menu item', async ({ page }) => {
    // Click "Thêm món"
    await page.click('button:has-text("Thêm món")');
    
    // Fill basic info
    await page.fill('input[name="name"]', 'Cà phê sữa đá');
    await page.fill('input[name="category"]', 'Cà phê');
    await page.fill('textarea[name="description"]', 'Cà phê phin truyền thống');
    
    // Check "Món có nhiều size"
    await page.check('input[type="checkbox"][name="has_variants"]');
    
    // Add Size M
    await page.click('button:has-text("Thêm size")');
    await page.fill('input[name="variants[0].id"]', 'M');
    await page.fill('input[name="variants[0].name"]', 'Size M');
    await page.fill('input[name="variants[0].price"]', '25000');
    await page.check('input[name="variants[0].is_default"]');
    
    // Select ingredients for Size M
    await page.click('button[data-variant="0"]:has-text("Chọn nguyên liệu")');
    await page.click('text=Cà phê');
    await page.fill('input[name="quantity"]', '20');
    await page.click('button:has-text("Thêm")');
    
    // Add Size L
    await page.click('button:has-text("Thêm size")');
    await page.fill('input[name="variants[1].id"]', 'L');
    await page.fill('input[name="variants[1].name"]', 'Size L');
    await page.fill('input[name="variants[1].price"]', '30000');
    
    // Select ingredients for Size L
    await page.click('button[data-variant="1"]:has-text("Chọn nguyên liệu")');
    await page.click('text=Cà phê');
    await page.fill('input[name="quantity"]', '30');
    await page.click('button:has-text("Thêm")');
    
    // Save
    await page.click('button:has-text("Lưu")');
    
    // Verify item appears with variants
    await expect(page.locator('text=Cà phê sữa đá')).toBeVisible();
    await expect(page.locator('text=Size M - 25,000đ')).toBeVisible();
    await expect(page.locator('text=Size L - 30,000đ')).toBeVisible();
  });

  test('should toggle single-size to multi-size', async ({ page }) => {
    // Create single-size item first
    await page.click('button:has-text("Thêm món")');
    await page.fill('input[name="name"]', 'Test Item');
    await page.fill('input[name="category"]', 'Test');
    await page.fill('input[name="price"]', '20000');
    await page.click('button:has-text("Lưu")');
    
    // Edit item
    await page.click('[data-testid="menu-item-edit"]:has-text("Test Item")');
    
    // Toggle to multi-size
    await page.check('input[type="checkbox"][name="has_variants"]');
    
    // Verify price field disappears
    await expect(page.locator('input[name="price"]')).not.toBeVisible();
    
    // Add variant
    await page.click('button:has-text("Thêm size")');
    await page.fill('input[name="variants[0].id"]', 'M');
    await page.fill('input[name="variants[0].name"]', 'Size M');
    await page.fill('input[name="variants[0].price"]', '20000');
    await page.check('input[name="variants[0].is_default"]');
    
    // Save
    await page.click('button:has-text("Lưu")');
    
    // Verify converted to multi-size
    await expect(page.locator('text=Size M - 20,000đ')).toBeVisible();
  });
});
```

### 3.4 Test Order Flow with Variants

```javascript
// tests/order-with-variants.spec.js
import { test, expect } from '@playwright/test';

test.describe('Order Flow with Variants', () => {
  test.beforeEach(async ({ page }) => {
    // Login as waiter
    await page.goto('/');
    await page.fill('input[name="username"]', 'waiter1');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('/waiter/orders');
  });

  test('should order single-size item', async ({ page }) => {
    // Start new order
    await page.click('button:has-text("Đơn mới")');
    
    // Tap single-size item
    await page.click('[data-testid="menu-item"]:has-text("Bánh mì thịt")');
    
    // Verify item added to order
    await expect(page.locator('[data-testid="order-item"]:has-text("Bánh mì thịt")')).toBeVisible();
    await expect(page.locator('text=20,000đ')).toBeVisible();
    
    // Complete order
    await page.click('button:has-text("Hoàn tất")');
    
    // Verify order created
    await expect(page.locator('text=Đơn hàng đã tạo')).toBeVisible();
  });

  test('should order multi-size item with variant selection', async ({ page }) => {
    // Start new order
    await page.click('button:has-text("Đơn mới")');
    
    // Tap multi-size item
    await page.click('[data-testid="menu-item"]:has-text("Cà phê sữa đá")');
    
    // Verify variant options appear
    await expect(page.locator('text=Chọn size')).toBeVisible();
    await expect(page.locator('text=Size M - 25,000đ')).toBeVisible();
    await expect(page.locator('text=Size L - 30,000đ')).toBeVisible();
    
    // Select Size L
    await page.click('button:has-text("Size L")');
    
    // Verify item added with variant
    await expect(page.locator('text=Cà phê sữa đá (Size L)')).toBeVisible();
    await expect(page.locator('text=30,000đ')).toBeVisible();
    
    // Complete order
    await page.click('button:has-text("Hoàn tất")');
    
    // Verify order created
    await expect(page.locator('text=Đơn hàng đã tạo')).toBeVisible();
  });

  test('should order mixed items (single + multi-size)', async ({ page }) => {
    // Start new order
    await page.click('button:has-text("Đơn mới")');
    
    // Add single-size item
    await page.click('[data-testid="menu-item"]:has-text("Bánh mì thịt")');
    
    // Add multi-size item (Size M)
    await page.click('[data-testid="menu-item"]:has-text("Cà phê sữa đá")');
    await page.click('button:has-text("Size M")');
    
    // Add another multi-size item (Size L)
    await page.click('[data-testid="menu-item"]:has-text("Cà phê sữa đá")');
    await page.click('button:has-text("Size L")');
    
    // Verify all items in order
    await expect(page.locator('text=Bánh mì thịt')).toBeVisible();
    await expect(page.locator('text=Cà phê sữa đá (Size M)')).toBeVisible();
    await expect(page.locator('text=Cà phê sữa đá (Size L)')).toBeVisible();
    
    // Verify total
    const total = await page.locator('[data-testid="order-total"]').textContent();
    expect(total).toContain('75,000đ'); // 20k + 25k + 30k
    
    // Complete order
    await page.click('button:has-text("Hoàn tất")');
  });
});
```

## Bước 4: AI-Powered Test Generation với Codegen

### 4.1 Sử Dụng Playwright Codegen

Playwright Codegen tự động generate test code khi bạn interact với app.

```bash
# Start codegen
npm run test:e2e:codegen

# Hoặc với specific URL
npx playwright codegen http://localhost:5173
```

**Cách sử dụng**:
1. Browser window sẽ mở
2. Interact với app như bình thường (click, type, navigate)
3. Playwright sẽ tự động generate code trong Inspector window
4. Copy code vào test file

**Ví dụ**: Codegen sẽ generate code như này khi bạn login:

```javascript
await page.goto('http://localhost:5173/');
await page.getByLabel('Username').click();
await page.getByLabel('Username').fill('admin');
await page.getByLabel('Password').click();
await page.getByLabel('Password').fill('password123');
await page.getByRole('button', { name: 'Đăng nhập' }).click();
```

### 4.2 Record Test với Specific Actions

```bash
# Record test và save vào file
npx playwright codegen --target javascript -o tests/recorded-test.spec.js http://localhost:5173
```

### 4.3 Generate Tests từ Existing App

```bash
# Generate tests cho specific page
npx playwright codegen http://localhost:5173/manager/menu
```

## Bước 5: Chạy Tests

### 5.1 Chạy Tất Cả Tests

```bash
# Run all tests
npm run test:e2e

# Run with UI mode (interactive)
npm run test:e2e:ui

# Run in headed mode (see browser)
npm run test:e2e:headed

# Run specific test file
npx playwright test tests/menu-single-size.spec.js

# Run tests matching pattern
npx playwright test --grep "should create"
```

### 5.2 Debug Tests

```bash
# Debug mode (step through tests)
npm run test:e2e:debug

# Debug specific test
npx playwright test tests/menu-single-size.spec.js --debug
```

### 5.3 View Test Report

```bash
# Show HTML report
npm run test:e2e:report
```

## Bước 6: Best Practices

### 6.1 Use Data Test IDs

Thêm `data-testid` vào components để dễ select:

```vue
<!-- MenuView.vue -->
<template>
  <div>
    <button data-testid="add-menu-item">Thêm món</button>
    
    <div 
      v-for="item in items" 
      :key="item.id"
      data-testid="menu-item"
    >
      <span data-testid="menu-item-name">{{ item.name }}</span>
      <button data-testid="menu-item-edit">Edit</button>
      <button data-testid="menu-item-delete">Delete</button>
    </div>
  </div>
</template>
```

Trong test:

```javascript
await page.click('[data-testid="add-menu-item"]');
await expect(page.locator('[data-testid="menu-item-name"]').first()).toBeVisible();
```

### 6.2 Create Page Object Models

```javascript
// tests/pages/MenuPage.js
export class MenuPage {
  constructor(page) {
    this.page = page;
    this.addButton = page.locator('[data-testid="add-menu-item"]');
    this.nameInput = page.locator('input[name="name"]');
    this.priceInput = page.locator('input[name="price"]');
    this.saveButton = page.locator('button:has-text("Lưu")');
  }

  async goto() {
    await this.page.goto('/manager/menu');
  }

  async createSingleSizeItem(name, price) {
    await this.addButton.click();
    await this.nameInput.fill(name);
    await this.priceInput.fill(price);
    await this.saveButton.click();
  }
}
```

Sử dụng trong test:

```javascript
import { MenuPage } from './pages/MenuPage';

test('should create item using page object', async ({ page }) => {
  const menuPage = new MenuPage(page);
  await menuPage.goto();
  await menuPage.createSingleSizeItem('Test Item', '20000');
});
```

### 6.3 Use Fixtures for Common Setup

```javascript
// tests/fixtures.js
import { test as base } from '@playwright/test';

export const test = base.extend({
  // Auto-login fixture
  authenticatedPage: async ({ page }, use) => {
    await page.goto('/');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('/manager/dashboard');
    await use(page);
  },
});

export { expect } from '@playwright/test';
```

Sử dụng:

```javascript
import { test, expect } from './fixtures';

test('should access menu', async ({ authenticatedPage }) => {
  // Already logged in!
  await authenticatedPage.click('text=Menu');
  await expect(authenticatedPage).toHaveURL('/manager/menu');
});
```

## Bước 7: CI/CD Integration

### 7.1 GitHub Actions

```yaml
# .github/workflows/e2e-tests.yml
name: E2E Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      
      - name: Install Playwright Browsers
        run: |
          cd frontend
          npx playwright install --with-deps
      
      - name: Run E2E tests
        run: |
          cd frontend
          npm run test:e2e
      
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

## Bước 8: Advanced Features

### 8.1 Visual Regression Testing

```javascript
test('should match screenshot', async ({ page }) => {
  await page.goto('/manager/menu');
  await expect(page).toHaveScreenshot('menu-page.png');
});
```

### 8.2 Network Mocking

```javascript
test('should handle API errors', async ({ page }) => {
  // Mock API response
  await page.route('**/api/menu', route => {
    route.fulfill({
      status: 500,
      body: JSON.stringify({ error: 'Server error' })
    });
  });
  
  await page.goto('/manager/menu');
  await expect(page.locator('text=Error loading menu')).toBeVisible();
});
```

### 8.3 Mobile Testing

```javascript
test('should work on mobile', async ({ page }) => {
  // Emulate mobile device
  await page.setViewportSize({ width: 375, height: 667 });
  
  await page.goto('/manager/menu');
  
  // Test mobile-specific UI
  await expect(page.locator('[data-testid="mobile-menu"]')).toBeVisible();
});
```

## Tổng Kết

### Commands Chính

```bash
# Install
npm init playwright@latest

# Generate tests (AI-powered)
npm run test:e2e:codegen

# Run tests
npm run test:e2e              # Headless
npm run test:e2e:ui           # UI mode
npm run test:e2e:headed       # See browser
npm run test:e2e:debug        # Debug mode

# View report
npm run test:e2e:report
```

### Workflow Đề Xuất

1. **Development**: Dùng Codegen để generate test nhanh
2. **Refinement**: Refactor generated code, add assertions
3. **Organization**: Tạo Page Objects cho reusability
4. **CI/CD**: Integrate vào pipeline
5. **Maintenance**: Update tests khi UI changes

### Lợi Ích

- ✅ AI-powered test generation (Codegen)
- ✅ Auto-wait, không cần sleep
- ✅ Cross-browser testing
- ✅ Mobile emulation
- ✅ Screenshot & video on failure
- ✅ Parallel execution
- ✅ Great debugging tools

**Playwright là lựa chọn tốt nhất cho E2E testing hiện đại!** 🚀

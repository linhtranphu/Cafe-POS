# Playwright E2E Testing - Quick Start

## Cài Đặt Nhanh (5 phút)

### Bước 1: Install Playwright

```bash
cd frontend
npm init playwright@latest
```

Chọn:
- JavaScript ✅
- Test folder: `tests` ✅
- GitHub Actions: No ❌
- Install browsers: Yes ✅

### Bước 2: Cấu Hình

File `playwright.config.js` đã được tạo tự động. Chỉ cần update `baseURL`:

```javascript
use: {
  baseURL: 'http://localhost:5173',
}
```

### Bước 3: Viết Test Đầu Tiên

Tạo file `tests/login.spec.js`:

```javascript
import { test, expect } from '@playwright/test';

test('should login', async ({ page }) => {
  await page.goto('/');
  await page.fill('input[name="username"]', 'admin');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/manager/dashboard');
});
```

### Bước 4: Chạy Test

```bash
npm run test:e2e
```

## AI-Powered Test Generation

### Generate Test Tự Động

```bash
npm run test:e2e:codegen
```

Hoặc:

```bash
npx playwright codegen http://localhost:5173
```

**Cách dùng**:
1. Browser mở → Interact với app
2. Code tự động generate trong Inspector
3. Copy code vào test file
4. Done! ✅

## Commands Hữu Ích

```bash
# Run tests
npm run test:e2e                # Headless mode
npm run test:e2e:ui             # UI mode (interactive)
npm run test:e2e:headed         # See browser
npm run test:e2e:debug          # Debug mode

# Generate tests (AI)
npm run test:e2e:codegen        # Auto-generate tests

# View report
npm run test:e2e:report         # HTML report
```

## Example: Test Menu với Variants

```javascript
test('create multi-size item', async ({ page }) => {
  // Login
  await page.goto('/');
  await page.fill('input[name="username"]', 'admin');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  // Navigate to menu
  await page.click('text=Menu');
  
  // Create item
  await page.click('button:has-text("Thêm món")');
  await page.fill('input[name="name"]', 'Cà phê');
  await page.check('input[name="has_variants"]');
  
  // Add variant
  await page.click('button:has-text("Thêm size")');
  await page.fill('input[name="variants[0].id"]', 'M');
  await page.fill('input[name="variants[0].name"]', 'Size M');
  await page.fill('input[name="variants[0].price"]', '25000');
  
  // Save
  await page.click('button:has-text("Lưu")');
  
  // Verify
  await expect(page.locator('text=Cà phê')).toBeVisible();
  await expect(page.locator('text=Size M')).toBeVisible();
});
```

## Tips

### 1. Dùng Data Test IDs

```vue
<button data-testid="add-menu-item">Thêm món</button>
```

```javascript
await page.click('[data-testid="add-menu-item"]');
```

### 2. Auto-wait (Không cần sleep!)

```javascript
// ❌ Không cần
await page.click('button');
await page.waitForTimeout(1000);

// ✅ Playwright tự động wait
await page.click('button');
await expect(page.locator('text=Success')).toBeVisible();
```

### 3. Debug với UI Mode

```bash
npm run test:e2e:ui
```

- Time travel debugging
- Watch mode
- Step through tests
- See screenshots

## Tài Liệu Chi Tiết

Xem `frontend/PLAYWRIGHT_SETUP_GUIDE.md` để biết thêm:
- Page Object Models
- Fixtures
- Network mocking
- Visual regression testing
- CI/CD integration

## Tổng Kết

**3 bước để bắt đầu**:
1. `npm init playwright@latest` - Install
2. `npm run test:e2e:codegen` - Generate tests (AI)
3. `npm run test:e2e` - Run tests

**Đơn giản vậy thôi!** 🚀

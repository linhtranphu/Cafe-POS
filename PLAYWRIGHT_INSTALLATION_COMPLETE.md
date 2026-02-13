# Playwright Installation Complete ✅

## Đã Cài Đặt Thành Công

### Packages Installed
- ✅ @playwright/test v1.58.2
- ✅ Chromium browser
- ✅ Firefox browser
- ✅ FFmpeg (for video recording)

### Files Created
- ✅ `frontend/playwright.config.js` - Configuration
- ✅ `frontend/tests/menu-variants.spec.js` - Example test
- ✅ `frontend/package.json` - Updated with test scripts

## Commands Sẵn Sàng Sử Dụng

### 1. Generate Tests với AI (Codegen) 🤖

```bash
cd frontend
npm run test:e2e:codegen
```

**Cách dùng**:
1. Browser sẽ mở tại http://localhost:5173
2. Interact với app (click, type, navigate)
3. Code tự động generate trong Inspector window
4. Copy code vào test file

**Lưu ý**: Dev server phải đang chạy (`npm run dev`)

### 2. Run Tests

```bash
cd frontend

# Run all tests (headless)
npm run test:e2e

# Run with UI mode (interactive, recommended)
npm run test:e2e:ui

# Run with browser visible
npm run test:e2e:headed

# Debug mode (step through tests)
npm run test:e2e:debug
```

### 3. View Test Report

```bash
cd frontend
npm run test:e2e:report
```

## Quick Start Workflow

### Bước 1: Start Dev Server

```bash
cd frontend
npm run dev
```

Để terminal này chạy.

### Bước 2: Generate Test với Codegen (Terminal mới)

```bash
cd frontend
npm run test:e2e:codegen
```

1. Browser mở → Login vào app
2. Navigate đến Menu page
3. Click "Thêm món"
4. Fill form
5. Click "Lưu"
6. Code tự động generate!

### Bước 3: Copy Code vào Test File

Copy generated code từ Inspector vào `frontend/tests/menu-variants.spec.js`

### Bước 4: Run Test

```bash
npm run test:e2e:ui
```

## Example Test Structure

```javascript
import { test, expect } from '@playwright/test';

test('create menu item', async ({ page }) => {
  // Navigate
  await page.goto('/');
  
  // Login
  await page.fill('input[name="username"]', 'admin');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  // Create item
  await page.click('text=Menu');
  await page.click('button:has-text("Thêm món")');
  await page.fill('input[name="name"]', 'Test Item');
  await page.fill('input[name="price"]', '20000');
  await page.click('button:has-text("Lưu")');
  
  // Verify
  await expect(page.locator('text=Test Item')).toBeVisible();
});
```

## Tips

### 1. Dùng Codegen để Generate Tests Nhanh

Thay vì viết test thủ công, dùng Codegen:
```bash
npm run test:e2e:codegen
```

### 2. Dùng UI Mode để Debug

```bash
npm run test:e2e:ui
```

- Time travel debugging
- Watch mode
- Step through tests
- See screenshots

### 3. Add Data Test IDs vào Components

```vue
<button data-testid="add-menu-item">Thêm món</button>
```

```javascript
await page.click('[data-testid="add-menu-item"]');
```

### 4. Auto-wait (Không cần sleep!)

Playwright tự động wait cho elements:

```javascript
// ❌ Không cần
await page.click('button');
await page.waitForTimeout(1000);

// ✅ Playwright tự động wait
await page.click('button');
await expect(page.locator('text=Success')).toBeVisible();
```

## Troubleshooting

### Dev server không chạy

Nếu test fail vì không connect được:

1. Start dev server manually:
   ```bash
   cd frontend
   npm run dev
   ```

2. Trong terminal khác, run tests:
   ```bash
   cd frontend
   npm run test:e2e
   ```

### Browser không mở trong Codegen

Ensure dev server đang chạy:
```bash
# Terminal 1
cd frontend
npm run dev

# Terminal 2
cd frontend
npm run test:e2e:codegen
```

## Next Steps

1. **Generate tests cho existing features**:
   ```bash
   npm run test:e2e:codegen
   ```
   - Login flow
   - Menu management
   - Order creation

2. **Implement menu size variants feature** (theo spec)

3. **Generate tests cho menu variants**:
   - Create single-size item
   - Create multi-size item
   - Order with variants

4. **Run tests trong CI/CD** (optional)

## Documentation

- **Quick Start**: `PLAYWRIGHT_QUICK_START.md`
- **Full Guide**: `frontend/PLAYWRIGHT_SETUP_GUIDE.md`
- **Example Tests**: `frontend/tests/menu-variants.spec.js`

## Summary

✅ Playwright installed successfully
✅ Config file created
✅ Test scripts added to package.json
✅ Example test created
✅ Ready to generate tests with AI (Codegen)

**Bắt đầu ngay**:
```bash
cd frontend
npm run test:e2e:codegen
```

🚀 Happy Testing!

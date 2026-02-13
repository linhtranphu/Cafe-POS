# Playwright E2E Tests - Integrated into Tasks

## Tổng Quan

Đã tích hợp Playwright E2E testing vào `.kiro/specs/menu-size-variants/tasks.md` với chi tiết cụ thể cho từng test case.

## Files Đã Tạo

### 1. Test Fixtures
- ✅ `frontend/tests/fixtures/auth.js` - Auto-login fixtures
  - `managerPage` - Auto-login as manager
  - `waiterPage` - Auto-login as waiter

### 2. Test Helpers
- ✅ `frontend/tests/helpers/selectors.js` - Common selectors
  - Auth selectors
  - Menu selectors
  - Order selectors
  - Variant selectors

### 3. Test Files
- ✅ `frontend/tests/menu-single-size.spec.js` - Single-size item tests
  - Create single-size item
  - Edit single-size item
  - Delete single-size item

- ✅ `frontend/tests/menu-multi-size.spec.js` - Multi-size item tests
  - Create multi-size item with variants
  - Toggle single-size to multi-size

- ✅ `frontend/tests/menu-variants.spec.js` - Example tests (template)

## Tasks.md Updates

### Section 6.5: E2E Testing with Playwright

Đã update với:
- ✅ Setup instructions (fixtures, helpers)
- ✅ Chi tiết code cho từng test step
- ✅ Data-testid selectors
- ✅ Assertions
- ✅ Run commands

### Test Coverage trong Tasks

**6.5.0**: Setup fixtures and helpers
**6.5.1**: Manager creates single-size item (chi tiết đầy đủ)
**6.5.2**: Manager creates multi-size item (chi tiết đầy đủ)
**6.5.3**: Waiter orders single-size item (chi tiết đầy đủ)
**6.5.4**: Waiter orders multi-size item (chi tiết đầy đủ)
**6.5.5**: Mixed order (single + multi-size) (chi tiết đầy đủ)
**6.5.6-6.5.10**: Các tests khác (cần update tương tự)

## Cách Sử Dụng

### Bước 1: Add Data-testid vào Components

Khi implement UI, thêm `data-testid` attributes:

```vue
<!-- MenuView.vue -->
<template>
  <div>
    <button data-testid="add-menu-item">Thêm món</button>
    
    <input 
      data-testid="item-name"
      v-model="form.name"
      placeholder="Tên món"
    />
    
    <input 
      type="checkbox"
      data-testid="has-variants"
      v-model="form.hasVariants"
    />
    
    <div v-if="form.hasVariants" data-testid="variants-section">
      <button data-testid="add-variant">Thêm size</button>
      
      <div v-for="(variant, index) in form.variants" :key="index">
        <input 
          :data-testid="`variant-${index}-id`"
          v-model="variant.id"
        />
        <input 
          :data-testid="`variant-${index}-name`"
          v-model="variant.name"
        />
        <input 
          :data-testid="`variant-${index}-price`"
          v-model="variant.price"
        />
        <input 
          type="checkbox"
          :data-testid="`variant-${index}-is-default`"
          v-model="variant.isDefault"
        />
      </div>
    </div>
    
    <button data-testid="save-menu-item">Lưu</button>
  </div>
</template>
```

### Bước 2: Update Selectors

Update `frontend/tests/helpers/selectors.js` nếu cần thay đổi selectors.

### Bước 3: Run Tests

```bash
cd frontend

# Run all E2E tests
npm run test:e2e

# Run specific test file
npm run test:e2e tests/menu-single-size.spec.js

# Run with UI mode (recommended)
npm run test:e2e:ui

# Debug mode
npm run test:e2e:debug
```

### Bước 4: Generate More Tests với Codegen

```bash
npm run test:e2e:codegen
```

Interact với app → Code tự động generate → Copy vào test file

## Test Structure

### Using Fixtures

```javascript
import { test, expect } from './fixtures/auth';

test('test name', async ({ managerPage: page }) => {
  // Already logged in as manager!
  await page.click('[data-testid="menu-nav"]');
  // ... rest of test
});
```

### Using Selectors

```javascript
import { selectors } from './helpers/selectors';

test('test name', async ({ page }) => {
  await page.fill(selectors.menu.name, 'Test Item');
  await page.click(selectors.menu.save);
});
```

## Example: Complete Test Flow

```javascript
// frontend/tests/menu-complete-flow.spec.js
import { test, expect } from './fixtures/auth';
import { selectors } from './helpers/selectors';

test('complete flow: create item → order → verify', async ({ page }) => {
  // Login as manager
  await page.goto('/');
  await page.fill(selectors.auth.username, 'admin');
  await page.fill(selectors.auth.password, 'password123');
  await page.click(selectors.auth.loginButton);

  // Create multi-size item
  await page.click(selectors.menu.nav);
  await page.click(selectors.menu.addButton);
  await page.fill(selectors.menu.name, 'Cà phê');
  await page.check(selectors.menu.hasVariants);
  await page.click(selectors.menu.addVariant);
  await page.fill(selectors.menu.variantId(0), 'M');
  await page.fill(selectors.menu.variantName(0), 'Size M');
  await page.fill(selectors.menu.variantPrice(0), '25000');
  await page.check(selectors.menu.variantIsDefault(0));
  await page.click(selectors.menu.save);

  // Logout
  await page.click('[data-testid="logout"]');

  // Login as waiter
  await page.fill(selectors.auth.username, 'waiter1');
  await page.fill(selectors.auth.password, 'password123');
  await page.click(selectors.auth.loginButton);

  // Create order
  await page.click(selectors.order.newOrder);
  await page.click('[data-testid="menu-item-ca-phe"]');
  await page.click(selectors.order.variantOption('M'));
  await expect(page.locator('text=Cà phê (Size M)')).toBeVisible();
  await page.click(selectors.order.completeOrder);

  // Verify order created
  await expect(page.locator('text=Đơn hàng đã tạo')).toBeVisible();
});
```

## Best Practices

### 1. Use Data-testid (Không dùng text/class)

```javascript
// ❌ Tránh
await page.click('button.btn-primary');
await page.click('text=Thêm món');

// ✅ Tốt
await page.click('[data-testid="add-menu-item"]');
```

### 2. Use Fixtures cho Common Setup

```javascript
// ✅ Tốt - Reuse login logic
test('test name', async ({ managerPage: page }) => {
  // Already logged in!
});
```

### 3. Use Selectors Helper

```javascript
// ✅ Tốt - Centralized selectors
import { selectors } from './helpers/selectors';
await page.fill(selectors.menu.name, 'Test');
```

### 4. Add Meaningful Assertions

```javascript
// ❌ Không đủ
await page.click('[data-testid="save"]');

// ✅ Tốt
await page.click('[data-testid="save"]');
await expect(page.locator('text=Đã lưu thành công')).toBeVisible();
await expect(page.locator('[data-testid="menu-item"]')).toHaveCount(1);
```

## Integration với CI/CD

### GitHub Actions Example

```yaml
# .github/workflows/e2e-tests.yml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      
      - name: Install Playwright
        run: |
          cd frontend
          npx playwright install --with-deps chromium
      
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

## Next Steps

1. **Implement UI với data-testid attributes**
2. **Run tests sau khi implement mỗi feature**
3. **Generate thêm tests với Codegen**
4. **Update selectors nếu UI thay đổi**
5. **Add tests vào CI/CD pipeline**

## Summary

✅ Playwright E2E tests đã được tích hợp vào tasks.md
✅ Test fixtures và helpers đã được tạo
✅ Example tests đã sẵn sàng
✅ Chi tiết code cho từng test step
✅ Ready to implement và test!

**Workflow**:
1. Implement UI feature
2. Add data-testid attributes
3. Run Playwright tests
4. Fix bugs nếu có
5. Generate thêm tests với Codegen

🚀 E2E testing với Playwright giúp đảm bảo chất lượng cao!

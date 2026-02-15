import { test, expect } from '@playwright/test';

/**
 * E2E Test: Batch Integration with Menu and Orders
 * 
 * This test validates batch integration with the menu system:
 * 1. Create a batch
 * 2. Add batch to a menu item recipe
 * 3. Create an order using the menu item
 * 4. Verify batch quantity is deducted
 * 5. Verify cost calculation includes batch cost
 * 
 * Requirements tested: 5.1-5.6, 3.1-3.5
 */

test.describe('Batch Menu Integration E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Login as manager
    await page.goto('/');
    await page.fill('input[type="text"]', 'manager@test.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    
    await page.waitForSelector('text=Xin chào');
  });

  test('use batch in menu item recipe', async ({ page }) => {
    // Step 1: Create a batch first
    await page.click('text=Batch');
    await page.click('text=Records');
    await page.click('text=Ghi Nhận Batch');
    
    // Select batch definition and create batch
    await page.selectOption('select', { index: 1 });
    await page.fill('input[type="number"]', '1000');
    await page.click('text=Ghi Nhận');
    await page.click('text=Xác nhận');
    
    // Wait for success
    await expect(page.locator('text=thành công')).toBeVisible({ timeout: 5000 });
    
    // Step 2: Navigate to menu
    await page.click('text=Menu');
    await page.waitForURL('**/menu');
    
    // Select a menu item to edit
    const menuItem = page.locator('[class*="menu-item"]').first();
    await menuItem.click();
    
    // Click edit recipe
    await page.click('text=Sửa công thức');
    
    // Step 3: Add batch to recipe
    await page.click('text=Thêm thành phần');
    
    // Toggle to batch mode
    await page.click('text=Batch');
    
    // Select batch
    await page.selectOption('select[name="ingredient_type"]', 'batch');
    await page.selectOption('select[name="batch_id"]', { index: 1 });
    
    // Enter quantity
    await page.fill('input[name="quantity"]', '30');
    
    // Save recipe
    await page.click('text=Lưu');
    
    // Verify success
    await expect(page.locator('text=Cập nhật thành công')).toBeVisible({ timeout: 5000 });
    
    // Step 4: Verify batch appears in recipe
    await expect(page.locator('text=Batch')).toBeVisible();
    await expect(page.locator('text=30')).toBeVisible();
  });

  test('batch deduction when order is created', async ({ page }) => {
    // Navigate to batch records and note initial quantity
    await page.click('text=Batch');
    await page.click('text=Records');
    
    const initialQuantity = await page.locator('[class*="quantity"]').first().textContent();
    
    // Navigate to orders
    await page.click('text=Orders');
    await page.waitForURL('**/orders');
    
    // Create new order
    await page.click('text=Tạo order');
    
    // Add menu item that uses batch
    await page.click('[class*="menu-item"]').first();
    
    // Complete order
    await page.click('text=Hoàn tất');
    await page.click('text=Xác nhận');
    
    // Wait for order creation
    await expect(page.locator('text=Order đã tạo')).toBeVisible({ timeout: 5000 });
    
    // Navigate back to batch records
    await page.click('text=Batch');
    await page.click('text=Records');
    
    // Verify quantity decreased
    const newQuantity = await page.locator('[class*="quantity"]').first().textContent();
    expect(newQuantity).not.toBe(initialQuantity);
  });

  test('batch cost reflected in menu item cost', async ({ page }) => {
    // Navigate to menu
    await page.click('text=Menu');
    await page.waitForURL('**/menu');
    
    // Select menu item
    const menuItem = page.locator('[class*="menu-item"]').first();
    await menuItem.click();
    
    // View cost breakdown
    await page.click('text=Chi phí');
    
    // Verify batch cost is included
    await expect(page.locator('text=Batch')).toBeVisible();
    await expect(page.locator('text=Chi phí batch')).toBeVisible();
    
    // Verify total cost calculation
    await expect(page.locator('text=Tổng chi phí')).toBeVisible();
  });

  test('insufficient batch quantity prevents order', async ({ page }) => {
    // Navigate to orders
    await page.click('text=Orders');
    
    // Try to create order with insufficient batch
    await page.click('text=Tạo order');
    
    // Add menu item
    await page.click('[class*="menu-item"]').first();
    
    // Try to add large quantity
    await page.fill('input[name="quantity"]', '1000');
    
    // Try to complete order
    await page.click('text=Hoàn tất');
    
    // Should show error
    await expect(page.locator('text=Không đủ batch')).toBeVisible({ timeout: 5000 });
  });

  test('FIFO batch usage in orders', async ({ page }) => {
    // Create two batches with different times
    await page.click('text=Batch');
    await page.click('text=Records');
    
    // Create first batch
    await page.click('text=Ghi Nhận Batch');
    await page.selectOption('select', { index: 1 });
    await page.fill('input[type="number"]', '500');
    await page.click('text=Ghi Nhận');
    await page.click('text=Xác nhận');
    await page.waitForTimeout(2000);
    
    // Create second batch
    await page.click('text=Ghi Nhận Batch');
    await page.selectOption('select', { index: 1 });
    await page.fill('input[type="number"]', '500');
    await page.click('text=Ghi Nhận');
    await page.click('text=Xác nhận');
    
    // Note the first batch ID
    await page.click('text=Danh sách');
    const firstBatchId = await page.locator('[class*="batch-record"]').first().getAttribute('data-id');
    
    // Create order
    await page.click('text=Orders');
    await page.click('text=Tạo order');
    await page.click('[class*="menu-item"]').first();
    await page.click('text=Hoàn tất');
    await page.click('text=Xác nhận');
    
    // Check batch usage log
    await page.click('text=Batch');
    await page.click('text=Lịch sử');
    
    // Verify first batch was used (FIFO)
    const usedBatchId = await page.locator('[class*="usage-log"]').first().getAttribute('data-batch-id');
    expect(usedBatchId).toBe(firstBatchId);
  });

  test('batch availability warning in menu editor', async ({ page }) => {
    // Navigate to menu
    await page.click('text=Menu');
    
    // Edit menu item
    const menuItem = page.locator('[class*="menu-item"]').first();
    await menuItem.click();
    await page.click('text=Sửa công thức');
    
    // Add batch with high quantity
    await page.click('text=Thêm thành phần');
    await page.click('text=Batch');
    await page.selectOption('select[name="batch_id"]', { index: 1 });
    await page.fill('input[name="quantity"]', '10000'); // Very high quantity
    
    // Should show warning
    await expect(page.locator('text=Cảnh báo')).toBeVisible();
    await expect(page.locator('text=không đủ')).toBeVisible();
  });
});

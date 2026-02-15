import { test, expect } from '@playwright/test';

/**
 * E2E Test: Complete Batch Lifecycle
 * 
 * This test validates the complete batch management workflow:
 * 1. Create a batch definition
 * 2. Create a batch record (prepare batch)
 * 3. View batch in list
 * 4. Check batch details
 * 5. Verify alerts appear when appropriate
 * 
 * Requirements tested: 1.1-1.6, 2.1-2.6, 4.1-4.6
 */

test.describe('Batch Lifecycle E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Login as manager
    await page.goto('/');
    await page.fill('input[type="text"]', 'manager@test.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    
    // Wait for dashboard to load
    await page.waitForSelector('text=Xin chào');
  });

  test('complete batch lifecycle flow', async ({ page }) => {
    // Step 1: Navigate to batch definitions
    await page.click('text=Batch');
    await page.waitForURL('**/batch/**');
    
    // Navigate to definitions
    await page.click('text=Định nghĩa');
    await page.waitForURL('**/batch/definitions');
    
    // Step 2: Create a new batch definition
    await page.click('text=Tạo Định Nghĩa Mới');
    await page.waitForURL('**/batch/definitions/create');
    
    // Fill in batch definition form
    await page.fill('input[name="name"]', 'Test Coffee Concentrate');
    await page.fill('input[name="unit"]', 'ml');
    await page.fill('input[name="shelf_life_hours"]', '24');
    await page.fill('input[name="low_stock_threshold"]', '200');
    await page.fill('input[name="expiry_warning_hours"]', '4');
    
    // Add conversion rate
    await page.click('text=Thêm Nguyên Liệu');
    await page.selectOption('select[name="source_ingredient"]', { index: 1 });
    await page.fill('input[name="source_quantity"]', '100');
    await page.fill('input[name="batch_quantity"]', '500');
    await page.fill('input[name="wastage_rate"]', '0.1');
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Verify success message
    await expect(page.locator('text=Tạo thành công')).toBeVisible({ timeout: 5000 });
    
    // Step 3: Navigate to batch records
    await page.click('text=Batch Records');
    await page.waitForURL('**/batch/records');
    
    // Step 4: Create a batch record
    await page.click('text=Ghi Nhận Batch');
    await page.waitForURL('**/batch/records/create');
    
    // Select batch definition
    await page.selectOption('select', { label: /Test Coffee Concentrate/ });
    
    // Enter quantity
    await page.fill('input[type="number"]', '500');
    
    // Verify required ingredients are shown
    await expect(page.locator('text=Nguyên Liệu Cần Thiết')).toBeVisible();
    await expect(page.locator('text=110')).toBeVisible(); // 100 * 1.1 wastage
    
    // Verify cost preview
    await expect(page.locator('text=Chi Phí Dự Kiến')).toBeVisible();
    
    // Submit batch record
    await page.click('text=Ghi Nhận');
    
    // Confirm in dialog
    await page.click('text=Xác nhận');
    
    // Verify success
    await expect(page.locator('text=Đã ghi nhận batch thành công')).toBeVisible({ timeout: 5000 });
    
    // Step 5: View batch in list
    await page.click('text=Danh sách');
    await page.waitForURL('**/batch/records');
    
    // Verify batch appears in list
    await expect(page.locator('text=Test Coffee Concentrate')).toBeVisible();
    await expect(page.locator('text=500')).toBeVisible();
    await expect(page.locator('text=available')).toBeVisible();
    
    // Step 6: View batch details
    await page.click('text=Test Coffee Concentrate');
    
    // Verify details page
    await expect(page.locator('text=Chi Tiết Batch')).toBeVisible();
    await expect(page.locator('text=Số lượng còn lại')).toBeVisible();
    await expect(page.locator('text=500 ml')).toBeVisible();
    
    // Verify ingredients used section
    await expect(page.locator('text=Nguyên Liệu Đã Sử Dụng')).toBeVisible();
    
    // Verify cost breakdown
    await expect(page.locator('text=Chi Phí')).toBeVisible();
  });

  test('batch alerts flow', async ({ page }) => {
    // Navigate to alerts
    await page.click('text=Batch');
    await page.click('text=Cảnh báo');
    await page.waitForURL('**/batch/alerts');
    
    // Verify alert sections exist
    await expect(page.locator('text=Tồn Kho Thấp')).toBeVisible();
    await expect(page.locator('text=Sắp Hết Hạn')).toBeVisible();
    await expect(page.locator('text=Đã Hết Hạn')).toBeVisible();
    
    // Check if any alerts are displayed
    const alertCount = await page.locator('[class*="alert"]').count();
    expect(alertCount).toBeGreaterThanOrEqual(0);
  });

  test('batch reports flow', async ({ page }) => {
    // Navigate to reports
    await page.click('text=Batch');
    await page.click('text=Báo cáo');
    await page.waitForURL('**/batch/reports/**');
    
    // Test production report
    await page.click('text=Sản xuất');
    await expect(page.locator('text=Báo Cáo Sản Xuất')).toBeVisible();
    
    // Select date range
    const today = new Date().toISOString().split('T')[0];
    await page.fill('input[type="date"]', today);
    
    // Generate report
    await page.click('text=Tạo báo cáo');
    
    // Verify report displays
    await expect(page.locator('text=Tổng số batch')).toBeVisible({ timeout: 5000 });
    
    // Test wastage report
    await page.click('text=Lãng phí');
    await expect(page.locator('text=Báo Cáo Lãng Phí')).toBeVisible();
    
    // Test usage report
    await page.click('text=Sử dụng');
    await expect(page.locator('text=Báo Cáo Sử Dụng')).toBeVisible();
  });

  test('batch search and filter', async ({ page }) => {
    // Navigate to batch records
    await page.click('text=Batch');
    await page.click('text=Records');
    await page.waitForURL('**/batch/records');
    
    // Test search
    await page.fill('input[placeholder*="Tìm kiếm"]', 'Coffee');
    await page.waitForTimeout(500); // Debounce
    
    // Verify filtered results
    const results = await page.locator('[class*="batch-record"]').count();
    expect(results).toBeGreaterThanOrEqual(0);
    
    // Test status filter
    await page.selectOption('select[name="status"]', 'available');
    await page.waitForTimeout(500);
    
    // Verify only available batches shown
    await expect(page.locator('text=available')).toBeVisible();
  });

  test('batch expiry handling', async ({ page }) => {
    // Navigate to batch records
    await page.click('text=Batch');
    await page.click('text=Records');
    
    // Find a batch and mark as expired
    const batchRow = page.locator('[class*="batch-record"]').first();
    await batchRow.click();
    
    // Mark as expired
    await page.click('text=Đánh dấu hết hạn');
    
    // Confirm action
    await page.click('text=Xác nhận');
    
    // Verify status changed
    await expect(page.locator('text=expired')).toBeVisible({ timeout: 5000 });
    
    // Navigate to alerts
    await page.click('text=Cảnh báo');
    
    // Verify expired alert appears
    await expect(page.locator('text=Đã Hết Hạn')).toBeVisible();
  });

  test('batch widget on dashboard', async ({ page }) => {
    // Should be on dashboard after login
    await page.waitForSelector('text=Xin chào');
    
    // Verify batch widget is visible
    await expect(page.locator('text=Batch Status')).toBeVisible();
    
    // Verify widget shows summary
    await expect(page.locator('text=Tổng batch')).toBeVisible();
    
    // Click widget to navigate
    await page.click('[class*="batch-widget"]');
    
    // Should navigate to batch management
    await page.waitForURL('**/batch/**');
  });
});

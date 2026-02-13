// E2E Test: Manager creates single-size menu item
import { test, expect } from './fixtures/auth';
import { selectors } from './helpers/selectors';

test.describe('Menu Management - Single-Size Items', () => {
  test('should create single-size menu item', async ({ managerPage: page }) => {
    // Navigate to menu
    await page.click(selectors.menu.nav);
    await expect(page).toHaveURL('/manager/menu');

    // Click "Thêm món"
    await page.click(selectors.menu.addButton);
    await expect(page.locator(selectors.menu.form)).toBeVisible();

    // Fill basic info
    await page.fill(selectors.menu.name, 'Bánh mì thịt');
    await page.selectOption(selectors.menu.category, 'Món ăn');
    await page.fill(selectors.menu.description, 'Bánh mì Việt Nam truyền thống');

    // Ensure has_variants is unchecked
    const checkbox = page.locator(selectors.menu.hasVariants);
    if (await checkbox.isChecked()) {
      await checkbox.uncheck();
    }
    await expect(checkbox).not.toBeChecked();

    // Enter price
    await page.fill(selectors.menu.price, '20000');

    // Select ingredients (update selector to match your app)
    await page.click(selectors.menu.selectIngredients);
    await page.click('[data-testid="ingredient-banh-mi"]');
    await page.fill(selectors.menu.ingredientQuantity, '1');
    await page.click(selectors.menu.addIngredient);

    // Save
    await page.click(selectors.menu.save);

    // Verify item appears in list
    await expect(page.locator('text=Bánh mì thịt')).toBeVisible();
    await expect(page.locator('text=20,000đ')).toBeVisible();
  });

  test('should edit single-size menu item', async ({ managerPage: page }) => {
    await page.click(selectors.menu.nav);
    
    // Click edit on first item
    await page.click(`${selectors.menu.edit}:first-child`);
    
    // Update price
    await page.fill(selectors.menu.price, '25000');
    
    // Save
    await page.click(selectors.menu.save);
    
    // Verify updated price
    await expect(page.locator('text=25,000đ')).toBeVisible();
  });

  test('should delete single-size menu item', async ({ managerPage: page }) => {
    await page.click(selectors.menu.nav);
    
    // Get item name before delete
    const itemName = await page.locator('[data-testid="menu-item-name"]:first-child').textContent();
    
    // Click delete
    await page.click(`${selectors.menu.delete}:first-child`);
    
    // Confirm delete
    await page.click('button:has-text("Xác nhận")');
    
    // Verify item removed
    await expect(page.locator(`text=${itemName}`)).not.toBeVisible();
  });
});

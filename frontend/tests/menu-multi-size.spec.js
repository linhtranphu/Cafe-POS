// E2E Test: Manager creates multi-size menu item
import { test, expect } from './fixtures/auth';
import { selectors } from './helpers/selectors';

test.describe('Menu Management - Multi-Size Items', () => {
  test('should create multi-size menu item with variants', async ({ managerPage: page }) => {
    // Navigate to menu
    await page.click(selectors.menu.nav);
    await expect(page).toHaveURL('/manager/menu');

    // Click "Thêm món"
    await page.click(selectors.menu.addButton);

    // Fill basic info
    await page.fill(selectors.menu.name, 'Cà phê sữa đá');
    await page.selectOption(selectors.menu.category, 'Cà phê');
    await page.fill(selectors.menu.description, 'Cà phê phin truyền thống với sữa đá');

    // Check "Món có nhiều size"
    await page.check(selectors.menu.hasVariants);
    await expect(page.locator(selectors.menu.variantsSection)).toBeVisible();

    // Add Size M
    await page.click(selectors.menu.addVariant);
    await page.fill(selectors.menu.variantId(0), 'M');
    await page.fill(selectors.menu.variantName(0), 'Size M');
    await page.fill(selectors.menu.variantPrice(0), '25000');
    await page.check(selectors.menu.variantIsDefault(0));

    // Select ingredients for Size M
    await page.click(selectors.menu.variantSelectIngredients(0));
    await page.click('[data-testid="ingredient-ca-phe"]');
    await page.fill(selectors.menu.ingredientQuantity, '20');
    await page.selectOption(selectors.menu.ingredientUnit, 'g');
    await page.click(selectors.menu.addIngredient);

    // Add Size L
    await page.click(selectors.menu.addVariant);
    await page.fill(selectors.menu.variantId(1), 'L');
    await page.fill(selectors.menu.variantName(1), 'Size L');
    await page.fill(selectors.menu.variantPrice(1), '30000');

    // Select ingredients for Size L
    await page.click(selectors.menu.variantSelectIngredients(1));
    await page.click('[data-testid="ingredient-ca-phe"]');
    await page.fill(selectors.menu.ingredientQuantity, '30');
    await page.click(selectors.menu.addIngredient);

    // Save
    await page.click(selectors.menu.save);

    // Verify item appears with variants
    await expect(page.locator('text=Cà phê sữa đá')).toBeVisible();
    await expect(page.locator('text=Size M')).toBeVisible();
    await expect(page.locator('text=Size L')).toBeVisible();
    await expect(page.locator('text=25,000đ')).toBeVisible();
    await expect(page.locator('text=30,000đ')).toBeVisible();
  });

  test('should toggle single-size to multi-size', async ({ managerPage: page }) => {
    await page.click(selectors.menu.nav);

    // Create single-size item first
    await page.click(selectors.menu.addButton);
    await page.fill(selectors.menu.name, 'Test Item');
    await page.selectOption(selectors.menu.category, 'Test');
    await page.fill(selectors.menu.price, '20000');
    await page.click(selectors.menu.save);

    // Edit item
    await page.click(`${selectors.menu.edit}:has-text("Test Item")`);

    // Toggle to multi-size
    await page.check(selectors.menu.hasVariants);

    // Verify price field disappears
    await expect(page.locator(selectors.menu.price)).not.toBeVisible();

    // Add variant
    await page.click(selectors.menu.addVariant);
    await page.fill(selectors.menu.variantId(0), 'M');
    await page.fill(selectors.menu.variantName(0), 'Size M');
    await page.fill(selectors.menu.variantPrice(0), '20000');
    await page.check(selectors.menu.variantIsDefault(0));

    // Save
    await page.click(selectors.menu.save);

    // Verify converted to multi-size
    await expect(page.locator('text=Size M')).toBeVisible();
  });
});

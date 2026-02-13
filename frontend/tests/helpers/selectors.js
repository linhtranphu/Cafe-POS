// Common selectors for Playwright tests
// Update these to match your actual app's data-testid attributes

export const selectors = {
  // Auth
  auth: {
    username: '[data-testid="username"]',
    password: '[data-testid="password"]',
    loginButton: '[data-testid="login-button"]',
  },

  // Menu
  menu: {
    nav: '[data-testid="menu-nav"]',
    addButton: '[data-testid="add-menu-item"]',
    form: '[data-testid="menu-form"]',
    
    // Form fields
    name: '[data-testid="item-name"]',
    category: '[data-testid="item-category"]',
    description: '[data-testid="item-description"]',
    price: '[data-testid="item-price"]',
    hasVariants: '[data-testid="has-variants"]',
    
    // Variants
    variantsSection: '[data-testid="variants-section"]',
    addVariant: '[data-testid="add-variant"]',
    variantId: (index) => `[data-testid="variant-${index}-id"]`,
    variantName: (index) => `[data-testid="variant-${index}-name"]`,
    variantPrice: (index) => `[data-testid="variant-${index}-price"]`,
    variantIsDefault: (index) => `[data-testid="variant-${index}-is-default"]`,
    variantSelectIngredients: (index) => `[data-testid="variant-${index}-select-ingredients"]`,
    
    // Ingredients
    selectIngredients: '[data-testid="select-ingredients"]',
    ingredientQuantity: '[data-testid="ingredient-quantity"]',
    ingredientUnit: '[data-testid="ingredient-unit"]',
    addIngredient: '[data-testid="add-ingredient"]',
    
    // Actions
    save: '[data-testid="save-menu-item"]',
    cancel: '[data-testid="cancel"]',
    edit: '[data-testid="menu-item-edit"]',
    delete: '[data-testid="menu-item-delete"]',
  },

  // Order
  order: {
    newOrder: '[data-testid="new-order"]',
    form: '[data-testid="order-form"]',
    
    // Menu items
    menuItem: (id) => `[data-testid="menu-item-${id}"]`,
    
    // Variant selector
    variantSelector: '[data-testid="variant-selector"]',
    variantOption: (id) => `[data-testid="variant-option-${id}"]`,
    
    // Order items
    orderItem: '[data-testid="order-item"]',
    orderItemPrice: '[data-testid="order-item-price"]',
    itemQuantity: '[data-testid="item-quantity"]',
    itemSubtotal: '[data-testid="item-subtotal"]',
    increaseQuantity: '[data-testid="increase-quantity"]',
    decreaseQuantity: '[data-testid="decrease-quantity"]',
    
    // Order summary
    orderTotal: '[data-testid="order-total"]',
    orderDetails: '[data-testid="order-details"]',
    
    // Actions
    completeOrder: '[data-testid="complete-order"]',
    cancelOrder: '[data-testid="cancel-order"]',
  },
};

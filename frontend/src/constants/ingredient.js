// Ingredient Constants
// Must match backend: backend/domain/ingredient/ingredient.go

// Unit Types
export const INGREDIENT_UNITS = {
  // Mass units (kg base)
  KILOGRAM: 'kg',
  GRAM: 'g',
  
  // Volume units (L base)
  LITER: 'L',
  MILLILITER: 'ml',
  
  // Count units
  PIECE: 'piece',
  BOX: 'box',
  PACK: 'pack'
}

export const INGREDIENT_UNIT_OPTIONS = [
  { value: INGREDIENT_UNITS.KILOGRAM, label: 'Kilogram (kg)' },
  { value: INGREDIENT_UNITS.GRAM, label: 'Gram (g)' },
  { value: INGREDIENT_UNITS.LITER, label: 'Lít (L)' },
  { value: INGREDIENT_UNITS.MILLILITER, label: 'Milliliter (ml)' },
  { value: INGREDIENT_UNITS.PIECE, label: 'Cái' },
  { value: INGREDIENT_UNITS.BOX, label: 'Hộp' },
  { value: INGREDIENT_UNITS.PACK, label: 'Gói' }
]

// Stock Adjustment Types
export const ADJUSTMENT_TYPES = {
  ADD: 'add',
  REMOVE: 'remove',
  ADJUST: 'adjust'
}

export const ADJUSTMENT_TYPE_OPTIONS = [
  { value: ADJUSTMENT_TYPES.ADD, label: 'Nhập Hàng', color: 'green' },
  { value: ADJUSTMENT_TYPES.REMOVE, label: 'Xuất Hàng', color: 'red' },
  { value: ADJUSTMENT_TYPES.ADJUST, label: 'Điều Chỉnh', color: 'blue' }
]

export const ADJUSTMENT_TYPE_CLASSES = {
  [ADJUSTMENT_TYPES.ADD]: 'bg-green-100 text-green-800',
  [ADJUSTMENT_TYPES.REMOVE]: 'bg-red-100 text-red-800',
  [ADJUSTMENT_TYPES.ADJUST]: 'bg-blue-100 text-blue-800'
}

// Stock Status
export const STOCK_STATUS = {
  IN_STOCK: 'in_stock',
  LOW_STOCK: 'low_stock',
  OUT_OF_STOCK: 'out_of_stock'
}

export const STOCK_STATUS_CLASSES = {
  [STOCK_STATUS.IN_STOCK]: 'bg-green-100 text-green-800',
  [STOCK_STATUS.LOW_STOCK]: 'bg-yellow-100 text-yellow-800',
  [STOCK_STATUS.OUT_OF_STOCK]: 'bg-red-100 text-red-800'
}

// Helper functions
export function getStockStatus(ingredient) {
  if (!ingredient) return STOCK_STATUS.OUT_OF_STOCK
  
  if (ingredient.quantity === 0) {
    return STOCK_STATUS.OUT_OF_STOCK
  } else if (ingredient.quantity <= ingredient.min_stock) {
    return STOCK_STATUS.LOW_STOCK
  }
  return STOCK_STATUS.IN_STOCK
}

export function getStockStatusClass(ingredient) {
  const status = getStockStatus(ingredient)
  return STOCK_STATUS_CLASSES[status] || 'bg-gray-100 text-gray-800'
}

export function getStockStatusText(ingredient) {
  const status = getStockStatus(ingredient)
  const texts = {
    [STOCK_STATUS.IN_STOCK]: 'Đủ Hàng',
    [STOCK_STATUS.LOW_STOCK]: 'Sắp Hết',
    [STOCK_STATUS.OUT_OF_STOCK]: 'Hết Hàng'
  }
  return texts[status] || 'Không xác định'
}

export function getAdjustmentTypeClass(type) {
  return ADJUSTMENT_TYPE_CLASSES[type] || 'bg-gray-100 text-gray-800'
}

export function getAdjustmentTypeText(type) {
  const option = ADJUSTMENT_TYPE_OPTIONS.find(opt => opt.value === type)
  return option ? option.label : type
}

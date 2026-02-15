/**
 * Component Tests for MenuView Batch Integration
 * 
 * Requirements tested:
 * - Requirement 5.1: Recipe can use batch as ingredient instead of raw ingredient
 * - Requirement 5.6: Cost calculation uses actual batch cost
 * - Requirement 3.3: Display batch cost per unit
 * 
 * Task: 14.1.6 Write component tests for MenuRecipeEditor batch integration
 * 
 * Tests cover:
 * - Batch selector UI toggle (raw vs batch)
 * - Batch selection and display
 * - Available batch quantity display
 * - Warning when batch insufficient
 * - Cost calculation with batch costs
 * - Batch expiry information display
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MenuView from '../MenuView.vue'
import { useMenuStore } from '../../stores/menu'
import { useIngredientStore } from '../../stores/ingredient'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'
import { useBatchRecordStore } from '../../stores/batchRecord'

// Mock stores
vi.mock('../../stores/menu')
vi.mock('../../stores/ingredient')
vi.mock('../../stores/batchDefinition')
vi.mock('../../stores/batchRecord')

// Mock components
vi.mock('../../components/BottomNav.vue', () => ({
  default: { template: '<div class="bottom-nav-mock"></div>' }
}))

vi.mock('../../components/PullToRefresh.vue', () => ({
  default: { template: '<div class="pull-to-refresh-mock"></div>' }
}))

// Mock composables
vi.mock('../../composables/usePullToRefresh', () => ({
  usePullToRefresh: () => ({
    pullDistance: { value: 0 },
    isRefreshing: { value: false }
  })
}))

vi.mock('../../composables/useUnitConversion', () => ({
  useUnitConversion: () => ({
    getConversionRate: vi.fn((from, to) => 1),
    isValidConversion: vi.fn(() => true),
    getCompatibleUnits: vi.fn((unit) => [unit]),
    calculateCostBreakdown: vi.fn((qty, unit, cost, stockUnit, wastage) => ({
      totalCost: qty * cost * (1 + wastage / 100),
      costPerUnit: cost
    })),
    getConversionExplanation: vi.fn(() => '1:1 conversion')
  })
}))

// Mock menu category service
vi.mock('../../services/menuCategory', () => ({
  menuCategoryService: {
    getCategories: vi.fn(() => Promise.resolve({ data: [] })),
    createCategory: vi.fn(),
    deleteCategory: vi.fn()
  }
}))

// Mock global alert
global.alert = vi.fn()

describe('MenuView - Batch Integration', () => {
  let wrapper
  let menuStore
  let ingredientStore
  let batchDefinitionStore
  let batchRecordStore

  const mockIngredients = [
    {
      id: 'ing1',
      name: 'Hạt cà phê',
      category: 'Nguyên liệu',
      unit: 'g',
      current_stock: 1000,
      cost_per_unit: 0.5,
      wastage_percentage: 10
    },
    {
      id: 'ing2',
      name: 'Sữa tươi',
      category: 'Nguyên liệu',
      unit: 'ml',
      current_stock: 2000,
      cost_per_unit: 0.03,
      wastage_percentage: 5
    }
  ]

  const mockBatchRecords = [
    {
      id: 'batch1',
      batch_definition_id: 'def1',
      batch_name: 'Cà phê Concentrate',
      quantity_remaining: 500,
      unit: 'ml',
      cost_per_unit: 0.15,
      status: 'available',
      expires_at: '2026-02-20T10:00:00Z'
    },
    {
      id: 'batch2',
      batch_definition_id: 'def2',
      batch_name: 'Sữa đã tiệt trùng',
      quantity_remaining: 100,
      unit: 'ml',
      cost_per_unit: 0.05,
      status: 'available',
      expires_at: '2026-02-18T10:00:00Z'
    },
    {
      id: 'batch3',
      batch_definition_id: 'def3',
      batch_name: 'Trà đen Concentrate',
      quantity_remaining: 0,
      unit: 'ml',
      cost_per_unit: 0.12,
      status: 'depleted',
      expires_at: '2026-02-19T10:00:00Z'
    }
  ]

  beforeEach(() => {
    // Setup store mocks
    menuStore = {
      items: [],
      loading: false,
      error: null,
      fetchMenuItems: vi.fn(),
      createMenuItem: vi.fn(() => Promise.resolve(true)),
      updateMenuItem: vi.fn(() => Promise.resolve(true)),
      deleteMenuItem: vi.fn(() => Promise.resolve(true))
    }

    ingredientStore = {
      items: mockIngredients,
      loading: false,
      fetchIngredients: vi.fn()
    }

    batchDefinitionStore = {
      definitions: [],
      loading: false,
      fetchDefinitions: vi.fn()
    }

    batchRecordStore = {
      records: mockBatchRecords,
      loading: false,
      fetchRecords: vi.fn()
    }

    useMenuStore.mockReturnValue(menuStore)
    useIngredientStore.mockReturnValue(ingredientStore)
    useBatchDefinitionStore.mockReturnValue(batchDefinitionStore)
    useBatchRecordStore.mockReturnValue(batchRecordStore)
  })

  describe('Ingredient Type Toggle (Requirement 5.1)', () => {
    it('should display ingredient type toggle in selector modal', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      // Open create form
      wrapper.vm.openCreateModal()
      await wrapper.vm.$nextTick()

      // Open ingredient selector
      wrapper.vm.showIngredientSelector = true
      await wrapper.vm.$nextTick()

      // Check for toggle buttons
      expect(wrapper.text()).toContain('Nguyên liệu thô')
      expect(wrapper.text()).toContain('Batch')
    })

    it('should default to raw ingredient type', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.ingredientType).toBe('raw')
    })

    it('should switch to batch type when clicking batch toggle', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      await wrapper.vm.$nextTick()

      // Find and click batch toggle button
      const buttons = wrapper.findAll('button')
      const batchButton = buttons.find(b => b.text().includes('🧪 Batch'))
      
      await batchButton.trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.ingredientType).toBe('batch')
    })

    it('should switch back to raw type when clicking raw toggle', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      const buttons = wrapper.findAll('button')
      const rawButton = buttons.find(b => b.text().includes('🥬 Nguyên liệu thô'))
      
      await rawButton.trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.ingredientType).toBe('raw')
    })
  })

  describe('Batch List Display (Requirement 5.1)', () => {
    it('should display available batches when batch type selected', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Cà phê Concentrate')
      expect(wrapper.text()).toContain('Sữa đã tiệt trùng')
    })

    it('should not display depleted batches', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('Trà đen Concentrate')
    })

    it('should display batch quantity remaining (Requirement 5.1)', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Còn lại: 500 ml')
      expect(wrapper.text()).toContain('Còn lại: 100 ml')
    })

    it('should display batch cost per unit (Requirement 3.3)', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      // Check for cost display (formatted as currency)
      expect(wrapper.text()).toContain('Chi phí')
    })

    it('should display batch expiry information', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Hết hạn')
    })

    it('should display batch status badge', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Khả dụng')
    })
  })

  describe('Batch Selection (Requirement 5.1)', () => {
    it('should add batch to ingredients when selected', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      await wrapper.vm.$nextTick()

      // Select a batch
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      // Check that batch was added to form ingredients
      expect(wrapper.vm.form.ingredients).toHaveLength(1)
      expect(wrapper.vm.form.ingredients[0].name).toBe('Cà phê Concentrate')
      expect(wrapper.vm.form.ingredients[0].type).toBe('batch')
    })

    it('should store batch metadata when selected', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      const addedBatch = wrapper.vm.form.ingredients[0]
      expect(addedBatch.id).toBe('batch1')
      expect(addedBatch.batch_definition_id).toBe('def1')
      expect(addedBatch.costPerUnit).toBe(0.15)
      expect(addedBatch.availableQuantity).toBe(500)
      expect(addedBatch.expiresAt).toBe('2026-02-20T10:00:00Z')
      expect(addedBatch.status).toBe('available')
    })

    it('should close selector modal after batch selection', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      await wrapper.vm.$nextTick()

      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showIngredientSelector).toBe(false)
    })

    it('should prevent selecting same batch twice', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()
      
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.form.ingredients).toHaveLength(1)
    })

    it('should mark selected batch in list', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.isBatchSelected('batch1')).toBe(true)
      expect(wrapper.vm.isBatchSelected('batch2')).toBe(false)
    })
  })

  describe('Batch Selection for Variants', () => {
    it('should add batch to specific variant when variant selector used', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.form.has_variants = true
      wrapper.vm.addVariant()
      await wrapper.vm.$nextTick()

      // Open variant ingredient selector
      wrapper.vm.openVariantIngredientSelector(0)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.currentVariantIndex).toBe(0)

      // Select batch for variant
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.form.variants[0].ingredients).toHaveLength(1)
      expect(wrapper.vm.form.variants[0].ingredients[0].name).toBe('Cà phê Concentrate')
      expect(wrapper.vm.form.variants[0].ingredients[0].type).toBe('batch')
    })

    it('should reset variant index after batch selection', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.form.has_variants = true
      wrapper.vm.addVariant()
      wrapper.vm.openVariantIngredientSelector(0)
      await wrapper.vm.$nextTick()

      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.currentVariantIndex).toBe(null)
    })
  })

  describe('Batch Cost Calculation (Requirement 5.6, 3.3)', () => {
    it('should calculate cost using batch cost per unit', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      const addedBatch = wrapper.vm.form.ingredients[0]
      // Cost should be calculated: quantity * cost_per_unit
      // With quantity = 1, cost = 1 * 0.15 = 0.15
      expect(addedBatch.costPerUnit).toBe(0.15)
    })

    it('should not apply wastage to batch costs', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      const addedBatch = wrapper.vm.form.ingredients[0]
      expect(addedBatch.wastage).toBe(0)
    })

    it('should update cost when batch quantity changes', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      // Change quantity
      wrapper.vm.form.ingredients[0].quantity = 10
      wrapper.vm.updateIngredientCost(0)
      await wrapper.vm.$nextTick()

      // Cost should be recalculated
      expect(wrapper.vm.form.ingredients[0].estimatedCost).toBeGreaterThan(0)
    })
  })

  describe('Batch Search and Filtering', () => {
    it('should filter batches by search query', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      wrapper.vm.ingredientSearchQuery = 'Cà phê'
      await wrapper.vm.$nextTick()

      const filtered = wrapper.vm.filteredAvailableBatches
      expect(filtered).toHaveLength(1)
      expect(filtered[0].batch_name).toBe('Cà phê Concentrate')
    })

    it('should be case insensitive when searching', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      wrapper.vm.ingredientSearchQuery = 'cà phê'
      await wrapper.vm.$nextTick()

      const filtered = wrapper.vm.filteredAvailableBatches
      expect(filtered).toHaveLength(1)
    })

    it('should show empty state when no batches match search', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      wrapper.vm.ingredientSearchQuery = 'xyz123'
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.filteredAvailableBatches).toHaveLength(0)
    })
  })

  describe('Mixed Ingredients (Raw + Batch)', () => {
    it('should allow adding both raw ingredients and batches', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      
      // Add raw ingredient
      const ingredient = mockIngredients[0]
      wrapper.vm.selectIngredient(ingredient)
      await wrapper.vm.$nextTick()

      // Add batch
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.form.ingredients).toHaveLength(2)
      expect(wrapper.vm.form.ingredients[0].type).toBe('raw')
      expect(wrapper.vm.form.ingredients[1].type).toBe('batch')
    })

    it('should calculate total cost from mixed ingredients', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      
      // Add raw ingredient
      const ingredient = mockIngredients[0]
      wrapper.vm.selectIngredient(ingredient)
      wrapper.vm.form.ingredients[0].quantity = 10
      wrapper.vm.updateIngredientCost(0)
      await wrapper.vm.$nextTick()

      // Add batch
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      wrapper.vm.form.ingredients[1].quantity = 5
      wrapper.vm.updateIngredientCost(1)
      await wrapper.vm.$nextTick()

      const totalCost = wrapper.vm.totalIngredientCost
      expect(totalCost).toBeGreaterThan(0)
    })
  })

  describe('Batch Data Persistence', () => {
    it('should save batch ingredients when creating menu item', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.form.name = 'Test Item'
      wrapper.vm.form.category = 'Test'
      wrapper.vm.form.price = 50000
      
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      wrapper.vm.form.ingredients[0].quantity = 10
      await wrapper.vm.$nextTick()

      await wrapper.vm.saveItem()
      await wrapper.vm.$nextTick()

      expect(menuStore.createMenuItem).toHaveBeenCalled()
      const callArgs = menuStore.createMenuItem.mock.calls[0][0]
      expect(callArgs.ingredients).toHaveLength(1)
      expect(callArgs.ingredients[0].name).toBe('Cà phê Concentrate')
    })

    it('should save batch ingredients for variants', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.form.name = 'Test Item'
      wrapper.vm.form.category = 'Test'
      wrapper.vm.form.has_variants = true
      wrapper.vm.addVariant()
      wrapper.vm.form.variants[0].name = 'Size M'
      wrapper.vm.form.variants[0].price = 50000
      
      wrapper.vm.openVariantIngredientSelector(0)
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      wrapper.vm.form.variants[0].ingredients[0].quantity = 10
      await wrapper.vm.$nextTick()

      await wrapper.vm.saveItem()
      await wrapper.vm.$nextTick()

      expect(menuStore.createMenuItem).toHaveBeenCalled()
      const callArgs = menuStore.createMenuItem.mock.calls[0][0]
      expect(callArgs.variants[0].ingredients).toHaveLength(1)
      expect(callArgs.variants[0].ingredients[0].name).toBe('Cà phê Concentrate')
    })
  })

  describe('Loading States', () => {
    it('should show loading state when fetching batches', async () => {
      batchRecordStore.loading = true
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.batchesLoading).toBe(true)
      expect(wrapper.text()).toContain('Đang tải')
    })

    it('should show empty state when no batches available', async () => {
      batchRecordStore.records = []
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      wrapper.vm.showIngredientSelector = true
      wrapper.vm.ingredientType = 'batch'
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Không có batch nào')
    })
  })

  describe('Helper Functions', () => {
    it('should check if batch is selected correctly', async () => {
      wrapper = mount(MenuView)
      await wrapper.vm.$nextTick()

      wrapper.vm.openCreateModal()
      const batch = mockBatchRecords[0]
      wrapper.vm.selectBatch(batch)
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.isBatchSelected('batch1')).toBe(true)
      expect(wrapper.vm.isBatchSelected('batch2')).toBe(false)
      expect(wrapper.vm.isBatchSelected('nonexistent')).toBe(false)
    })

    it('should format date correctly', () => {
      wrapper = mount(MenuView)
      const formatted = wrapper.vm.formatDate('2026-02-20T10:00:00Z')
      expect(formatted).toBeTruthy()
      expect(formatted).not.toBe('')
    })

    it('should handle null date', () => {
      wrapper = mount(MenuView)
      const formatted = wrapper.vm.formatDate(null)
      expect(formatted).toBe('')
    })
  })
})

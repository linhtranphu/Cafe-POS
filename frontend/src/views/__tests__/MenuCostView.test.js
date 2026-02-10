/**
 * Unit Tests for MenuCostView Component
 * 
 * NOTE: These tests require a testing framework (vitest + @vue/test-utils) to be installed.
 * To run these tests:
 * 1. Install dependencies: npm install -D vitest @vue/test-utils happy-dom
 * 2. Add test script to package.json: "test": "vitest --run"
 * 3. Run tests: npm test
 * 
 * Requirements tested:
 * - 4.1: Menu Item Cost Report API integration
 * - 4.3: Category filtering
 * - 4.4: Sorting functionality
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import MenuCostView from '../MenuCostView.vue'
import { menuCostService } from '../../services/menuCost'

// Mock the menu cost service
vi.mock('../../services/menuCost', () => ({
  menuCostService: {
    getMenuCosts: vi.fn(),
    getMenuCostDetail: vi.fn()
  }
}))

// Mock the composables
vi.mock('../../composables/usePullToRefresh', () => ({
  usePullToRefresh: () => ({
    pullDistance: { value: 0 },
    isRefreshing: { value: false }
  })
}))

describe('MenuCostView', () => {
  let wrapper
  
  const mockMenuCosts = {
    items: [
      {
        menu_item_id: '1',
        name: 'Cappuccino',
        category: 'Coffee',
        price: 45000,
        current_cost: 15000,
        profit_margin: 66.67,
        absolute_profit: 30000,
        cost_status: 'FINAL',
        warning_status: 'none'
      },
      {
        menu_item_id: '2',
        name: 'Latte',
        category: 'Coffee',
        price: 50000,
        current_cost: 20000,
        profit_margin: 60.00,
        absolute_profit: 30000,
        cost_status: 'FINAL',
        warning_status: 'none'
      },
      {
        menu_item_id: '3',
        name: 'Green Tea',
        category: 'Tea',
        price: 30000,
        current_cost: 25000,
        profit_margin: 16.67,
        absolute_profit: 5000,
        cost_status: 'FINAL',
        warning_status: 'low_margin'
      },
      {
        menu_item_id: '4',
        name: 'Special Promo',
        category: 'Coffee',
        price: 20000,
        current_cost: 25000,
        profit_margin: -25.00,
        absolute_profit: -5000,
        cost_status: 'FINAL',
        warning_status: 'loss'
      }
    ],
    summary: {
      total_items: 4,
      loss_count: 1,
      low_margin_count: 1,
      average_profit_margin: 54.59
    },
    recalculation_status: {
      in_progress: false,
      queued_items: 0,
      processed_items: 4,
      last_updated: '2024-01-15T10:00:00Z'
    }
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    menuCostService.getMenuCosts.mockResolvedValue(mockMenuCosts)
  })

  describe('Component Rendering', () => {
    it('should render the component with menu items', async () => {
      wrapper = mount(MenuCostView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true
          }
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.find('h1').text()).toContain('Chi phí món')
      expect(wrapper.findAll('.menu-item').length).toBe(4)
    })

    it('should display summary statistics', async () => {
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const summary = wrapper.find('.summary-statistics')
      expect(summary.text()).toContain('4') // total items
      expect(summary.text()).toContain('1') // loss count
      expect(summary.text()).toContain('1') // low margin count
    })

    it('should show loading state initially', () => {
      wrapper = mount(MenuCostView)
      expect(wrapper.text()).toContain('Đang tải dữ liệu')
    })

    it('should show error state when API fails', async () => {
      menuCostService.getMenuCosts.mockRejectedValue(new Error('API Error'))
      
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải dữ liệu')
    })
  })

  describe('Category Filtering (Requirement 4.3)', () => {
    beforeEach(async () => {
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display all items when no category filter is selected', () => {
      const items = wrapper.findAll('.menu-item')
      expect(items.length).toBe(4)
    })

    it('should filter items by Coffee category', async () => {
      await wrapper.vm.categoryFilter = 'Coffee'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(3)
      expect(items.every(item => item.category === 'Coffee')).toBe(true)
    })

    it('should filter items by Tea category', async () => {
      await wrapper.vm.categoryFilter = 'Tea'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(1)
      expect(items[0].category).toBe('Tea')
    })

    it('should show all items when category filter is cleared', async () => {
      await wrapper.vm.categoryFilter = 'Coffee'
      await wrapper.vm.$nextTick()
      
      await wrapper.vm.categoryFilter = ''
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(4)
    })
  })

  describe('Sorting Functionality (Requirement 4.4)', () => {
    beforeEach(async () => {
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should sort by profit_margin in descending order by default', () => {
      const items = wrapper.vm.filteredMenuItems
      expect(items[0].profit_margin).toBeGreaterThan(items[1].profit_margin)
    })

    it('should sort by profit_margin in ascending order', async () => {
      wrapper.vm.sortBy = 'profit_margin'
      wrapper.vm.sortOrder = 'asc'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items[0].profit_margin).toBeLessThan(items[items.length - 1].profit_margin)
    })

    it('should sort by absolute_profit in descending order', async () => {
      wrapper.vm.sortBy = 'absolute_profit'
      wrapper.vm.sortOrder = 'desc'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items[0].absolute_profit).toBeGreaterThanOrEqual(items[1].absolute_profit)
    })

    it('should sort by name in ascending order', async () => {
      wrapper.vm.sortBy = 'name'
      wrapper.vm.sortOrder = 'asc'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items[0].name.localeCompare(items[1].name)).toBeLessThanOrEqual(0)
    })

    it('should toggle sort order', async () => {
      const initialOrder = wrapper.vm.sortOrder
      wrapper.vm.toggleSortOrder()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.sortOrder).not.toBe(initialOrder)
    })
  })

  describe('Warning Color Coding (Requirement 7.2, 7.3)', () => {
    beforeEach(async () => {
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display green border for profitable items', () => {
      const profitableItem = mockMenuCosts.items[0]
      const borderClass = wrapper.vm.getItemBorderClass(profitableItem.warning_status)
      expect(borderClass).toContain('border-green-500')
    })

    it('should display yellow border for low margin items', () => {
      const lowMarginItem = mockMenuCosts.items[2]
      const borderClass = wrapper.vm.getItemBorderClass(lowMarginItem.warning_status)
      expect(borderClass).toContain('border-yellow-500')
    })

    it('should display red border for loss items', () => {
      const lossItem = mockMenuCosts.items[3]
      const borderClass = wrapper.vm.getItemBorderClass(lossItem.warning_status)
      expect(borderClass).toContain('border-red-500')
    })

    it('should display red text for negative profit margin', () => {
      const color = wrapper.vm.getProfitMarginColor(-25, 'loss')
      expect(color).toContain('text-red-600')
    })

    it('should display yellow text for low margin', () => {
      const color = wrapper.vm.getProfitMarginColor(16.67, 'low_margin')
      expect(color).toContain('text-yellow-600')
    })

    it('should display green text for good profit margin', () => {
      const color = wrapper.vm.getProfitMarginColor(66.67, 'none')
      expect(color).toContain('text-green-600')
    })
  })

  describe('Search Functionality', () => {
    beforeEach(async () => {
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should filter items by name search', async () => {
      wrapper.vm.searchQuery = 'Cappuccino'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(1)
      expect(items[0].name).toBe('Cappuccino')
    })

    it('should filter items by category search', async () => {
      wrapper.vm.searchQuery = 'Coffee'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(3)
      expect(items.every(item => item.category === 'Coffee')).toBe(true)
    })

    it('should be case insensitive', async () => {
      wrapper.vm.searchQuery = 'cappuccino'
      await wrapper.vm.$nextTick()

      const items = wrapper.vm.filteredMenuItems
      expect(items.length).toBe(1)
    })
  })

  describe('Cost Breakdown Modal', () => {
    const mockCostBreakdown = {
      menu_item: {
        id: '1',
        name: 'Cappuccino',
        price: 45000,
        current_cost: 15000
      },
      ingredients: [
        {
          name: 'Espresso',
          quantity: 30,
          unit: 'ml',
          cost_per_unit: 200,
          conversion_rate: 1.0,
          wastage_percentage: 5.0,
          total_cost: 6300
        },
        {
          name: 'Milk',
          quantity: 150,
          unit: 'ml',
          cost_per_unit: 50,
          conversion_rate: 1.0,
          wastage_percentage: 10.0,
          total_cost: 8250
        }
      ],
      total_cost: 15000
    }

    beforeEach(async () => {
      menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdown)
      wrapper = mount(MenuCostView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should open cost breakdown modal when item is clicked', async () => {
      await wrapper.vm.openCostBreakdown(mockMenuCosts.items[0])
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showCostBreakdown).toBe(true)
      expect(menuCostService.getMenuCostDetail).toHaveBeenCalledWith('1')
    })

    it('should display ingredient breakdown in modal', async () => {
      await wrapper.vm.openCostBreakdown(mockMenuCosts.items[0])
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.costBreakdown).toBeTruthy()
      expect(wrapper.vm.costBreakdown.ingredients.length).toBe(2)
    })

    it('should close modal when close button is clicked', async () => {
      await wrapper.vm.openCostBreakdown(mockMenuCosts.items[0])
      await wrapper.vm.$nextTick()

      wrapper.vm.closeCostBreakdown()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showCostBreakdown).toBe(false)
    })
  })

  describe('Helper Functions', () => {
    beforeEach(() => {
      wrapper = mount(MenuCostView)
    })

    it('should format percentage correctly', () => {
      expect(wrapper.vm.formatPercentage(66.67)).toBe('66.67%')
      expect(wrapper.vm.formatPercentage(0)).toBe('0.00%')
      expect(wrapper.vm.formatPercentage(null)).toBe('N/A')
    })

    it('should get correct cost status label', () => {
      expect(wrapper.vm.getCostStatusLabel('FINAL')).toContain('Chính thức')
      expect(wrapper.vm.getCostStatusLabel('ESTIMATED')).toContain('Ước tính')
      expect(wrapper.vm.getCostStatusLabel('INCOMPLETE')).toContain('Thiếu dữ liệu')
    })

    it('should get correct warning message', () => {
      expect(wrapper.vm.getWarningMessage('loss')).toContain('Bán lỗ')
      expect(wrapper.vm.getWarningMessage('low_margin')).toContain('Lợi nhuận thấp')
    })
  })
})

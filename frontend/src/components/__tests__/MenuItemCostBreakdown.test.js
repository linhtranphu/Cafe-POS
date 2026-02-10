/**
 * Unit Tests for MenuItemCostBreakdown Component
 * 
 * Requirements tested:
 * - 8.1: Display cost breakdown modal
 * - 8.2: Display ingredient breakdown table
 * - 8.3: Display total cost summary and warnings
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MenuItemCostBreakdown from '../MenuItemCostBreakdown.vue'
import { menuCostService } from '../../services/menuCost'

// Mock the menu cost service
vi.mock('../../services/menuCost', () => ({
  menuCostService: {
    getMenuCostDetail: vi.fn()
  }
}))

describe('MenuItemCostBreakdown', () => {
  let wrapper

  const mockCostBreakdown = {
    menu_item: {
      id: '1',
      name: 'Cappuccino',
      price: 45000,
      current_cost: 15000,
      cost_status: 'FINAL'
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

  const mockCostBreakdownWithMissingCost = {
    menu_item: {
      id: '2',
      name: 'Latte',
      price: 50000,
      current_cost: 0,
      cost_status: 'INCOMPLETE'
    },
    ingredients: [
      {
        name: 'Espresso',
        quantity: 30,
        unit: 'ml',
        cost_per_unit: 200,
        conversion_rate: 1.0,
        wastage_percentage: 0,
        total_cost: 6000
      },
      {
        name: 'Milk',
        quantity: 200,
        unit: 'ml',
        cost_per_unit: 0, // Missing cost
        conversion_rate: 1.0,
        wastage_percentage: 0,
        total_cost: 0
      }
    ],
    total_cost: 6000
  }

  const mockCostBreakdownWithConversion = {
    menu_item: {
      id: '3',
      name: 'Special Coffee',
      price: 60000,
      current_cost: 20000,
      cost_status: 'FINAL'
    },
    ingredients: [
      {
        name: 'Coffee Beans',
        quantity: 20,
        unit: 'g',
        cost_per_unit: 500,
        conversion_rate: 0.001, // kg to g conversion
        wastage_percentage: 15.0,
        total_cost: 11500
      },
      {
        name: 'Sugar',
        quantity: 10,
        unit: 'g',
        cost_per_unit: 100,
        conversion_rate: 1.0,
        wastage_percentage: 0,
        total_cost: 1000
      }
    ],
    total_cost: 20000
  }

  beforeEach(() => {
    menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdown)
  })

  describe('Component Rendering (Requirement 8.1)', () => {
    it('should not render when isOpen is false', () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: false,
          menuItemId: '1'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(false)
    })

    it('should render modal when isOpen is true', () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(true)
      expect(wrapper.text()).toContain('Chi tiết chi phí')
    })

    it('should show loading state initially', () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      expect(wrapper.text()).toContain('Đang tải')
    })

    it('should display menu item info after loading', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Cappuccino')
      expect(wrapper.text()).toContain('45,000')
      expect(wrapper.text()).toContain('15,000')
    })

    it('should show error state when API fails', async () => {
      menuCostService.getMenuCostDetail.mockRejectedValue(new Error('API Error'))

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải chi tiết chi phí')
    })

    it('should emit close event when close button is clicked', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const closeButton = wrapper.find('button')
      await closeButton.trigger('click')

      expect(wrapper.emitted('close')).toBeTruthy()
    })

    it('should emit close event when clicking outside modal', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()

      const backdrop = wrapper.find('.fixed')
      await backdrop.trigger('click')

      expect(wrapper.emitted('close')).toBeTruthy()
    })
  })

  describe('Ingredient Breakdown Table (Requirement 8.2, 8.3, 8.4)', () => {
    it('should display all ingredients with complete data', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Espresso')
      expect(wrapper.text()).toContain('30 ml')
      expect(wrapper.text()).toContain('200')
      expect(wrapper.text()).toContain('6,300')

      expect(wrapper.text()).toContain('Milk')
      expect(wrapper.text()).toContain('150 ml')
      expect(wrapper.text()).toContain('50')
      expect(wrapper.text()).toContain('8,250')
    })

    it('should display conversion rate when non-default', async () => {
      menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdownWithConversion)

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '3'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('0.001')
      expect(wrapper.text()).toContain('quy đổi')
    })

    it('should display wastage percentage when present', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('5%')
      expect(wrapper.text()).toContain('hao hụt')
      expect(wrapper.text()).toContain('10%')
    })

    it('should not display conversion rate when default (1.0)', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const milkIngredient = wrapper.vm.breakdown.ingredients.find(i => i.name === 'Milk')
      expect(wrapper.vm.hasNonDefaultConversion(milkIngredient)).toBe(false)
    })

    it('should not display wastage when zero', async () => {
      menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdownWithConversion)

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '3'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const sugarIngredient = wrapper.vm.breakdown.ingredients.find(i => i.name === 'Sugar')
      expect(wrapper.vm.hasWastage(sugarIngredient)).toBe(false)
    })

    it('should highlight ingredients with missing cost_per_unit', async () => {
      menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdownWithMissingCost)

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Thiếu giá nguyên liệu')
      
      const milkIngredient = wrapper.vm.breakdown.ingredients.find(i => i.name === 'Milk')
      expect(wrapper.vm.hasIncompleteCost(milkIngredient)).toBe(true)
    })
  })

  describe('Total Cost Summary (Requirement 8.3)', () => {
    it('should display total cost at bottom', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Tổng chi phí')
      expect(wrapper.text()).toContain('15,000')
    })

    it('should show warning when any ingredient has incomplete cost', async () => {
      menuCostService.getMenuCostDetail.mockResolvedValue(mockCostBreakdownWithMissingCost)

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.hasAnyIncompleteCost).toBe(true)
      expect(wrapper.text()).toContain('Một số nguyên liệu thiếu giá')
    })

    it('should not show warning when all ingredients have complete cost', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.hasAnyIncompleteCost).toBe(false)
    })
  })

  describe('Helper Functions', () => {
    beforeEach(() => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })
    })

    it('should detect incomplete cost correctly', () => {
      const ingredientWithCost = { cost_per_unit: 200 }
      const ingredientWithoutCost = { cost_per_unit: 0 }
      const ingredientWithNullCost = { cost_per_unit: null }

      expect(wrapper.vm.hasIncompleteCost(ingredientWithCost)).toBe(false)
      expect(wrapper.vm.hasIncompleteCost(ingredientWithoutCost)).toBe(true)
      expect(wrapper.vm.hasIncompleteCost(ingredientWithNullCost)).toBe(true)
    })

    it('should detect non-default conversion rate', () => {
      const defaultConversion = { conversion_rate: 1.0 }
      const nonDefaultConversion = { conversion_rate: 0.001 }
      const noConversion = { conversion_rate: null }

      expect(wrapper.vm.hasNonDefaultConversion(defaultConversion)).toBe(false)
      expect(wrapper.vm.hasNonDefaultConversion(nonDefaultConversion)).toBe(true)
      expect(wrapper.vm.hasNonDefaultConversion(noConversion)).toBe(false)
    })

    it('should detect wastage correctly', () => {
      const noWastage = { wastage_percentage: 0 }
      const withWastage = { wastage_percentage: 10 }
      const nullWastage = { wastage_percentage: null }

      expect(wrapper.vm.hasWastage(noWastage)).toBe(false)
      expect(wrapper.vm.hasWastage(withWastage)).toBe(true)
      expect(wrapper.vm.hasWastage(nullWastage)).toBe(false)
    })
  })

  describe('Data Fetching', () => {
    it('should fetch cost breakdown when component opens', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuCostService.getMenuCostDetail).toHaveBeenCalledWith('1')
    })

    it('should not fetch when menuItemId is null', async () => {
      menuCostService.getMenuCostDetail.mockClear()

      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: null
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuCostService.getMenuCostDetail).not.toHaveBeenCalled()
    })

    it('should refetch when menuItemId changes', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuCostService.getMenuCostDetail).toHaveBeenCalledWith('1')

      await wrapper.setProps({ menuItemId: '2' })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuCostService.getMenuCostDetail).toHaveBeenCalledWith('2')
    })

    it('should refetch when modal reopens', async () => {
      wrapper = mount(MenuItemCostBreakdown, {
        props: {
          isOpen: false,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()

      await wrapper.setProps({ isOpen: true })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuCostService.getMenuCostDetail).toHaveBeenCalledWith('1')
    })
  })
})

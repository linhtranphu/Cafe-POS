/**
 * Unit Tests for CostAnalysisView Component
 * 
 * Requirements tested:
 * - AC-10.1: Each variant displays current_cost
 * - AC-10.2: Each variant displays cost_status (FINAL/ESTIMATED/INCOMPLETE)
 * - AC-10.3: Each variant displays cost_last_calculated_at
 * - AC-10.5: Can see profit margin per variant (price - cost)
 * - AC-12.1: Can view all variants with their costs in one view
 * 
 * Task: 11a.4 Write unit tests for cost analysis components
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import CostAnalysisView from '../CostAnalysisView.vue'
import { menuService } from '../../services/menu'

// Mock components
vi.mock('../../components/BottomNav.vue', () => ({
  default: { template: '<div class="bottom-nav-mock"></div>' }
}))

vi.mock('../../components/PullToRefresh.vue', () => ({
  default: { template: '<div class="pull-to-refresh-mock"></div>' }
}))

vi.mock('../../components/SkeletonLoader.vue', () => ({
  default: { template: '<div class="skeleton-loader-mock"></div>' }
}))

vi.mock('../../components/CostBreakdownModal.vue', () => ({
  default: { 
    template: '<div class="cost-breakdown-modal-mock"></div>',
    props: ['isOpen', 'menuItemId', 'defaultVariantId'],
    emits: ['close']
  }
}))

vi.mock('../../components/ProfitComparisonModal.vue', () => ({
  default: { 
    template: '<div class="profit-comparison-modal-mock"></div>',
    props: ['isOpen', 'menuItemId'],
    emits: ['close']
  }
}))

// Mock composables
vi.mock('../../composables/usePullToRefresh', () => ({
  usePullToRefresh: () => ({
    pullDistance: { value: 0 },
    isRefreshing: { value: false }
  })
}))

// Mock the menu service
vi.mock('../../services/menu', () => ({
  menuService: {
    getMenuItems: vi.fn()
  }
}))

describe('CostAnalysisView', () => {
  let wrapper

  const mockSingleSizeItems = {
    data: [
      {
        id: '1',
        name: 'Bánh mì thịt',
        category: 'Món ăn',
        has_variants: false,
        price: 20000,
        current_cost: 12000,
        cost_status: 'FINAL',
        cost_last_calculated_at: '2026-02-13T10:00:00Z',
        available: true
      },
      {
        id: '2',
        name: 'Bánh mì chả',
        category: 'Món ăn',
        has_variants: false,
        price: 18000,
        current_cost: 10000,
        cost_status: 'ESTIMATED',
        cost_last_calculated_at: '2026-02-13T09:00:00Z',
        available: true
      }
    ]
  }

  const mockMultiSizeItems = {
    data: [
      {
        id: '3',
        name: 'Cà phê sữa đá',
        category: 'Cà phê',
        has_variants: true,
        available: true,
        variants: [
          {
            id: 'M',
            name: 'Size M',
            price: 25000,
            current_cost: 15000,
            cost_status: 'FINAL',
            cost_last_calculated_at: '2026-02-13T10:00:00Z',
            is_default: true,
            available: true
          },
          {
            id: 'L',
            name: 'Size L',
            price: 30000,
            current_cost: 18000,
            cost_status: 'FINAL',
            cost_last_calculated_at: '2026-02-13T10:00:00Z',
            is_default: false,
            available: true
          },
          {
            id: 'XL',
            name: 'Size XL',
            price: 35000,
            current_cost: 20000,
            cost_status: 'INCOMPLETE',
            cost_last_calculated_at: '2026-02-13T10:00:00Z',
            is_default: false,
            available: true
          }
        ]
      }
    ]
  }

  const mockMixedItems = {
    data: [
      ...mockSingleSizeItems.data,
      ...mockMultiSizeItems.data
    ]
  }

  beforeEach(() => {
    menuService.getMenuItems.mockResolvedValue(mockMixedItems)
  })

  describe('Component Rendering', () => {
    it('should render the view with header', async () => {
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()

      expect(wrapper.find('h1').text()).toContain('Phân tích chi phí theo size')
    })

    it('should render search input', async () => {
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()

      const searchInput = wrapper.find('input[type="text"]')
      expect(searchInput.exists()).toBe(true)
      expect(searchInput.attributes('placeholder')).toContain('Tìm kiếm')
    })

    it('should render cost status filter buttons', async () => {
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Tất cả')
      expect(wrapper.text()).toContain('Chính thức')
      expect(wrapper.text()).toContain('Ước tính')
      expect(wrapper.text()).toContain('Thiếu dữ liệu')
    })
  })

  describe('Single-Size Item Display (AC-10.1, AC-10.2, AC-10.3, AC-10.5)', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockSingleSizeItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display single-size items', () => {
      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Bánh mì chả')
    })

    it('should display item type as "Một size"', () => {
      expect(wrapper.text()).toContain('Một size')
    })

    it('should display price for single-size items', () => {
      expect(wrapper.text()).toContain('20,000')
      expect(wrapper.text()).toContain('18,000')
    })

    it('should display current_cost for single-size items (AC-10.1)', () => {
      expect(wrapper.text()).toContain('12,000')
      expect(wrapper.text()).toContain('10,000')
    })

    it('should display cost_status for single-size items (AC-10.2)', () => {
      expect(wrapper.text()).toContain('✓ Chính thức')
      expect(wrapper.text()).toContain('~ Ước tính')
    })

    it('should display cost_last_calculated_at for single-size items (AC-10.3)', () => {
      expect(wrapper.text()).toContain('Cập nhật')
    })

    it('should display profit for single-size items (AC-10.5)', () => {
      // Bánh mì thịt: 20000 - 12000 = 8000
      // Bánh mì chả: 18000 - 10000 = 8000
      expect(wrapper.text()).toContain('8,000')
    })

    it('should display profit margin for single-size items (AC-10.5)', () => {
      // Bánh mì thịt: (8000/20000) * 100 = 40%
      // Bánh mì chả: (8000/18000) * 100 = 44.44%
      expect(wrapper.text()).toContain('40')
      expect(wrapper.text()).toContain('44')
    })

    it('should have cost breakdown button for single-size items', () => {
      expect(wrapper.text()).toContain('Xem chi tiết chi phí')
    })
  })

  describe('Multi-Size Item Display (AC-10.1, AC-10.2, AC-10.3, AC-10.5, AC-12.1)', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockMultiSizeItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display multi-size items', () => {
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should display item type as "Nhiều size"', () => {
      expect(wrapper.text()).toContain('Nhiều size')
    })

    it('should display variant count', () => {
      expect(wrapper.text()).toContain('3 variants')
    })

    it('should display all variants (AC-12.1)', () => {
      expect(wrapper.text()).toContain('Size M')
      expect(wrapper.text()).toContain('Size L')
      expect(wrapper.text()).toContain('Size XL')
    })

    it('should mark default variant', () => {
      expect(wrapper.text()).toContain('Mặc định')
    })

    it('should display price for each variant (AC-12.1)', () => {
      expect(wrapper.text()).toContain('25,000')
      expect(wrapper.text()).toContain('30,000')
      expect(wrapper.text()).toContain('35,000')
    })

    it('should display current_cost for each variant (AC-10.1, AC-12.1)', () => {
      expect(wrapper.text()).toContain('15,000')
      expect(wrapper.text()).toContain('18,000')
      expect(wrapper.text()).toContain('20,000')
    })

    it('should display cost_status for each variant (AC-10.2)', () => {
      expect(wrapper.text()).toContain('✓ Chính thức')
      expect(wrapper.text()).toContain('⚠ Thiếu dữ liệu')
    })

    it('should display cost_last_calculated_at for each variant (AC-10.3)', () => {
      const updateTexts = wrapper.text().match(/Cập nhật/g)
      expect(updateTexts).toBeTruthy()
      expect(updateTexts.length).toBeGreaterThan(0)
    })

    it('should display profit for each variant (AC-10.5)', () => {
      // Size M: 25000 - 15000 = 10000
      // Size L: 30000 - 18000 = 12000
      // Size XL: 35000 - 20000 = 15000
      expect(wrapper.text()).toContain('10,000')
      expect(wrapper.text()).toContain('12,000')
      expect(wrapper.text()).toContain('15,000')
    })

    it('should display profit margin for each variant (AC-10.5)', () => {
      // Size M: (10000/25000) * 100 = 40%
      // Size L: (12000/30000) * 100 = 40%
      // Size XL: (15000/35000) * 100 = 42.86%
      expect(wrapper.text()).toContain('40')
      expect(wrapper.text()).toContain('42')
    })

    it('should have cost breakdown button for each variant', () => {
      const buttons = wrapper.findAll('button')
      const costBreakdownButtons = buttons.filter(b => b.text().includes('Xem chi tiết chi phí'))
      expect(costBreakdownButtons.length).toBeGreaterThan(0)
    })

    it('should have profit comparison button for multi-size items', () => {
      expect(wrapper.text()).toContain('So sánh lợi nhuận các size')
    })
  })

  describe('Filtering by Cost Status', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockMixedItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should show all items by default', () => {
      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Bánh mì chả')
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should filter by FINAL status', async () => {
      const buttons = wrapper.findAll('button')
      const finalButton = buttons.find(b => b.text().includes('✓ Chính thức'))
      await finalButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Should show items with FINAL status
      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should filter by ESTIMATED status', async () => {
      const buttons = wrapper.findAll('button')
      const estimatedButton = buttons.find(b => b.text().includes('~ Ước tính'))
      await estimatedButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Should show items with ESTIMATED status
      expect(wrapper.text()).toContain('Bánh mì chả')
    })

    it('should filter by INCOMPLETE status', async () => {
      const buttons = wrapper.findAll('button')
      const incompleteButton = buttons.find(b => b.text().includes('⚠ Thiếu dữ liệu'))
      await incompleteButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Should show multi-size item with INCOMPLETE variant
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should reset filter when clicking "Tất cả"', async () => {
      // First apply a filter
      const buttons = wrapper.findAll('button')
      const finalButton = buttons.find(b => b.text().includes('✓ Chính thức'))
      await finalButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Then reset
      const allButton = buttons.find(b => b.text() === 'Tất cả')
      await allButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Should show all items again
      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Bánh mì chả')
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })
  })

  describe('Search Functionality', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockMixedItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should filter items by name', async () => {
      const searchInput = wrapper.find('input[type="text"]')
      await searchInput.setValue('Cà phê')
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Cà phê sữa đá')
      expect(wrapper.text()).not.toContain('Bánh mì thịt')
    })

    it('should filter items by category', async () => {
      const searchInput = wrapper.find('input[type="text"]')
      await searchInput.setValue('Món ăn')
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Bánh mì chả')
      expect(wrapper.text()).not.toContain('Cà phê sữa đá')
    })

    it('should be case insensitive', async () => {
      const searchInput = wrapper.find('input[type="text"]')
      await searchInput.setValue('cà phê')
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should show empty state when no results', async () => {
      const searchInput = wrapper.find('input[type="text"]')
      await searchInput.setValue('xyz123')
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Không tìm thấy món nào')
    })
  })

  describe('Modal Interactions', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockMixedItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should open cost breakdown modal when clicking detail button', async () => {
      const buttons = wrapper.findAll('button')
      const detailButton = buttons.find(b => b.text().includes('Xem chi tiết chi phí'))
      
      expect(wrapper.vm.showCostBreakdownModal).toBe(false)
      
      await detailButton.trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showCostBreakdownModal).toBe(true)
      expect(wrapper.vm.selectedMenuItemId).toBeTruthy()
    })

    it('should open profit comparison modal when clicking comparison button', async () => {
      const buttons = wrapper.findAll('button')
      const comparisonButton = buttons.find(b => b.text().includes('So sánh lợi nhuận các size'))
      
      expect(wrapper.vm.showProfitComparisonModal).toBe(false)
      
      await comparisonButton.trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showProfitComparisonModal).toBe(true)
      expect(wrapper.vm.selectedMenuItemId).toBeTruthy()
    })

    it('should close cost breakdown modal', async () => {
      wrapper.vm.showCostBreakdownModal = true
      wrapper.vm.selectedMenuItemId = '1'
      await wrapper.vm.$nextTick()

      wrapper.vm.closeCostBreakdownModal()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showCostBreakdownModal).toBe(false)
      expect(wrapper.vm.selectedMenuItemId).toBe(null)
    })

    it('should close profit comparison modal', async () => {
      wrapper.vm.showProfitComparisonModal = true
      wrapper.vm.selectedMenuItemId = '3'
      await wrapper.vm.$nextTick()

      wrapper.vm.closeProfitComparisonModal()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.showProfitComparisonModal).toBe(false)
      expect(wrapper.vm.selectedMenuItemId).toBe(null)
    })
  })

  describe('Loading and Error States', () => {
    it('should show loading state initially', async () => {
      wrapper = mount(CostAnalysisView)
      
      // Check that loading is true initially
      expect(wrapper.vm.loading).toBe(true)
    })

    it('should show error state when API fails', async () => {
      menuService.getMenuItems.mockRejectedValue({
        response: { data: { error: 'API Error' } }
      })

      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Lỗi tải dữ liệu')
      expect(wrapper.text()).toContain('API Error')
    })

    it('should show generic error when no error message provided', async () => {
      menuService.getMenuItems.mockRejectedValue(new Error('Network error'))

      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải dữ liệu menu')
    })

    it('should have retry button on error', async () => {
      menuService.getMenuItems.mockRejectedValue({
        response: { data: { error: 'API Error' } }
      })

      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const buttons = wrapper.findAll('button')
      const retryButton = buttons.find(b => b.text().includes('Thử lại'))
      expect(retryButton).toBeTruthy()
    })

    it('should retry fetching data when clicking retry button', async () => {
      menuService.getMenuItems.mockClear()
      menuService.getMenuItems.mockRejectedValueOnce({
        response: { data: { error: 'API Error' } }
      })
      menuService.getMenuItems.mockResolvedValueOnce(mockMixedItems)

      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItems).toHaveBeenCalledTimes(1)

      const buttons = wrapper.findAll('button')
      const retryButton = buttons.find(b => b.text().includes('Thử lại'))
      await retryButton.trigger('click')
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItems).toHaveBeenCalledTimes(2)
    })

    it('should show empty state when no items', async () => {
      menuService.getMenuItems.mockResolvedValue({ data: [] })

      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.menuItems).toEqual([])
      expect(wrapper.vm.filteredMenuItems).toEqual([])
      expect(wrapper.text()).toContain('Không tìm thấy món nào')
      expect(wrapper.text()).toContain('Chưa có món nào trong hệ thống')
    })
  })

  describe('Helper Functions', () => {
    beforeEach(async () => {
      menuService.getMenuItems.mockResolvedValue(mockMixedItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should format prices correctly', () => {
      const { formatPrice } = wrapper.vm
      expect(formatPrice(20000)).toContain('20')
      expect(formatPrice(20000)).toContain('₫')
    })

    it('should format percentages correctly', () => {
      const { formatPercentage } = wrapper.vm
      expect(formatPercentage(40.5)).toContain('40')
      expect(formatPercentage(40.5)).toContain('%')
    })

    it('should calculate profit margin correctly', () => {
      const { calculateProfitMargin } = wrapper.vm
      expect(calculateProfitMargin(25000, 15000)).toBe(40)
      expect(calculateProfitMargin(30000, 18000)).toBe(40)
      expect(calculateProfitMargin(0, 15000)).toBe(0)
    })

    it('should get profit color correctly', () => {
      const { getProfitColor } = wrapper.vm
      expect(getProfitColor(10000)).toBe('text-green-600')
      expect(getProfitColor(0)).toBe('text-gray-600')
      expect(getProfitColor(-5000)).toBe('text-red-600')
    })

    it('should get profit margin color correctly', () => {
      const { getProfitMarginColor } = wrapper.vm
      expect(getProfitMarginColor(40)).toBe('text-green-600')
      expect(getProfitMarginColor(15)).toBe('text-yellow-600')
      expect(getProfitMarginColor(-10)).toBe('text-red-600')
    })

    it('should get cost status label correctly', () => {
      const { getCostStatusLabel } = wrapper.vm
      expect(getCostStatusLabel('FINAL')).toBe('✓ Chính thức')
      expect(getCostStatusLabel('ESTIMATED')).toBe('~ Ước tính')
      expect(getCostStatusLabel('INCOMPLETE')).toBe('⚠ Thiếu dữ liệu')
      expect(getCostStatusLabel('')).toBe('Chưa tính')
    })

    it('should get cost status class correctly', () => {
      const { getCostStatusClass } = wrapper.vm
      expect(getCostStatusClass('FINAL')).toBe('bg-green-100 text-green-700')
      expect(getCostStatusClass('ESTIMATED')).toBe('bg-yellow-100 text-yellow-700')
      expect(getCostStatusClass('INCOMPLETE')).toBe('bg-red-100 text-red-700')
      expect(getCostStatusClass('')).toBe('bg-gray-100 text-gray-700')
    })

    it('should format date correctly', () => {
      const { formatDate } = wrapper.vm
      const result = formatDate('2026-02-13T10:00:00Z')
      expect(result).toBeTruthy()
      expect(result).not.toBe('Chưa có')
    })

    it('should handle null date', () => {
      const { formatDate } = wrapper.vm
      expect(formatDate(null)).toBe('Chưa có')
      expect(formatDate('')).toBe('Chưa có')
    })

    it('should get default variant ID correctly', () => {
      const { getDefaultVariantId } = wrapper.vm
      const item = mockMultiSizeItems.data[0]
      expect(getDefaultVariantId(item)).toBe('M')
    })

    it('should return null for single-size items', () => {
      const { getDefaultVariantId } = wrapper.vm
      const item = mockSingleSizeItems.data[0]
      expect(getDefaultVariantId(item)).toBe(null)
    })
  })

  describe('Data Fetching', () => {
    it('should fetch menu items on mount', async () => {
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItems).toHaveBeenCalled()
    })

    it('should store fetched items in state', async () => {
      menuService.getMenuItems.mockResolvedValue(mockMixedItems)
      wrapper = mount(CostAnalysisView)
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.menuItems).toEqual(mockMixedItems.data)
    })
  })
})

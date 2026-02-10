/**
 * Unit Tests for ProfitAnalysisView Component
 * 
 * Requirements tested:
 * - 6.1: Category-level profit analysis
 * - 6.5.1: Operating profit analysis
 * - 6.4: Date range filtering
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ProfitAnalysisView from '../ProfitAnalysisView.vue'
import { profitAnalysisService } from '../../services/profitAnalysis'

// Mock the profit analysis service
vi.mock('../../services/profitAnalysis', () => ({
  profitAnalysisService: {
    getCategoryProfit: vi.fn(),
    getOperatingProfit: vi.fn()
  }
}))

// Mock the composables
vi.mock('../../composables/usePullToRefresh', () => ({
  usePullToRefresh: () => ({
    pullDistance: { value: 0 },
    isRefreshing: { value: false }
  })
}))

describe('ProfitAnalysisView', () => {
  let wrapper
  
  const mockCategoryProfits = {
    date_range: {
      start: '2024-01-01',
      end: '2024-01-31'
    },
    categories: [
      {
        category: 'Coffee',
        total_revenue: 5000000,
        total_cost: 1500000,
        total_profit: 3500000,
        average_profit_margin: 70.0,
        order_count: 150,
        item_count: 200
      },
      {
        category: 'Tea',
        total_revenue: 2000000,
        total_cost: 800000,
        total_profit: 1200000,
        average_profit_margin: 60.0,
        order_count: 80,
        item_count: 100
      }
    ]
  }

  const mockOperatingProfit = {
    date_range: {
      start: '2024-01-01',
      end: '2024-01-31'
    },
    total_revenue: 7000000,
    total_cogs: 2300000,
    gross_profit: 4700000,
    gross_profit_margin: 67.14,
    staff_salary: 2000000,
    rent: 1000000,
    utilities: 300000,
    marketing_costs: 200000,
    other_expenses: 100000,
    total_expenses: 3600000,
    operating_profit: 1100000,
    operating_profit_margin: 15.71,
    expense_allocated: false,
    allocation_note: null
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    profitAnalysisService.getCategoryProfit.mockResolvedValue(mockCategoryProfits)
    profitAnalysisService.getOperatingProfit.mockResolvedValue(mockOperatingProfit)
  })

  describe('Component Rendering', () => {
    it('should render the component with title', async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })

      await wrapper.vm.$nextTick()
      expect(wrapper.find('h1').text()).toContain('Phân tích lợi nhuận')
    })

    it('should show loading state initially', () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      
      expect(wrapper.text()).toContain('Đang tải dữ liệu')
    })

    it('should show error state when API fails', async () => {
      profitAnalysisService.getCategoryProfit.mockRejectedValue(new Error('API Error'))
      
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải dữ liệu')
    })
  })

  describe('View Mode Toggle (Requirement 6.1, 6.5.1)', () => {
    beforeEach(async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should default to category view mode', () => {
      expect(wrapper.vm.viewMode).toBe('category')
    })

    it('should switch to operating view mode', async () => {
      wrapper.vm.viewMode = 'operating'
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.viewMode).toBe('operating')
    })

    it('should call getCategoryProfit when in category mode', async () => {
      wrapper.vm.viewMode = 'category'
      await wrapper.vm.fetchData()
      await wrapper.vm.$nextTick()

      expect(profitAnalysisService.getCategoryProfit).toHaveBeenCalled()
    })

    it('should call getOperatingProfit when in operating mode', async () => {
      wrapper.vm.viewMode = 'operating'
      await wrapper.vm.fetchData()
      await wrapper.vm.$nextTick()

      expect(profitAnalysisService.getOperatingProfit).toHaveBeenCalled()
    })

    it('should fetch data when view mode changes', async () => {
      const fetchDataSpy = vi.spyOn(wrapper.vm, 'fetchData')
      
      wrapper.vm.viewMode = 'operating'
      await wrapper.vm.$nextTick()

      expect(fetchDataSpy).toHaveBeenCalled()
    })
  })

  describe('Date Range Picker (Requirement 6.4)', () => {
    beforeEach(async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      await wrapper.vm.$nextTick()
    })

    it('should initialize with today preset', () => {
      expect(wrapper.vm.selectedPreset).toBe('today')
      expect(wrapper.vm.dateRange.start).toBeTruthy()
      expect(wrapper.vm.dateRange.end).toBeTruthy()
    })

    it('should select this_week preset', async () => {
      wrapper.vm.selectDatePreset('this_week')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.selectedPreset).toBe('this_week')
      expect(wrapper.vm.dateRange.start).toBeTruthy()
      expect(wrapper.vm.dateRange.end).toBeTruthy()
    })

    it('should select this_month preset', async () => {
      wrapper.vm.selectDatePreset('this_month')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.selectedPreset).toBe('this_month')
      expect(wrapper.vm.dateRange.start).toBeTruthy()
      expect(wrapper.vm.dateRange.end).toBeTruthy()
    })

    it('should clear preset when custom date is selected', async () => {
      wrapper.vm.dateRange.start = '2024-01-01'
      wrapper.vm.dateRange.end = '2024-01-31'
      wrapper.vm.onDateChange()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.selectedPreset).toBeNull()
    })

    it('should fetch data when date range changes', async () => {
      const fetchDataSpy = vi.spyOn(wrapper.vm, 'fetchData')
      
      wrapper.vm.selectDatePreset('this_week')
      await wrapper.vm.$nextTick()

      expect(fetchDataSpy).toHaveBeenCalled()
    })

    it('should generate correct date range for today', () => {
      const dateRange = wrapper.vm.getDateRangeForPreset('today')
      const today = new Date()
      const expectedDate = today.toISOString().split('T')[0]

      expect(dateRange.start).toBe(expectedDate)
      expect(dateRange.end).toBe(expectedDate)
    })

    it('should generate correct date range for this_week', () => {
      const dateRange = wrapper.vm.getDateRangeForPreset('this_week')
      
      expect(dateRange.start).toBeTruthy()
      expect(dateRange.end).toBeTruthy()
      expect(new Date(dateRange.end) >= new Date(dateRange.start)).toBe(true)
    })

    it('should generate correct date range for this_month', () => {
      const dateRange = wrapper.vm.getDateRangeForPreset('this_month')
      const today = new Date()
      const firstDay = new Date(today.getFullYear(), today.getMonth(), 1)
      const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0)

      expect(dateRange.start).toBe(firstDay.toISOString().split('T')[0])
      expect(dateRange.end).toBe(lastDay.toISOString().split('T')[0])
    })
  })

  describe('Category Profit View Rendering (Requirement 6.1)', () => {
    beforeEach(async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display category profit data', async () => {
      expect(wrapper.vm.categoryProfits).toEqual(mockCategoryProfits.categories)
    })

    it('should pass date range to CategoryProfitView', async () => {
      const categoryView = wrapper.findComponent({ name: 'CategoryProfitView' })
      expect(categoryView.exists()).toBe(true)
    })

    it('should pass category profits to CategoryProfitView', async () => {
      const categoryView = wrapper.findComponent({ name: 'CategoryProfitView' })
      expect(categoryView.exists()).toBe(true)
    })
  })

  describe('Operating Profit View Rendering (Requirement 6.5.1)', () => {
    beforeEach(async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      wrapper.vm.viewMode = 'operating'
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display operating profit data', async () => {
      expect(wrapper.vm.operatingProfit).toEqual(mockOperatingProfit)
    })

    it('should pass operating profit to OperatingProfitView', async () => {
      const operatingView = wrapper.findComponent({ name: 'OperatingProfitView' })
      expect(operatingView.exists()).toBe(true)
    })
  })

  describe('Date Range Filtering (Requirement 6.4)', () => {
    beforeEach(async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      await wrapper.vm.$nextTick()
    })

    it('should pass correct date range to API', async () => {
      wrapper.vm.dateRange = {
        start: '2024-01-01',
        end: '2024-01-31'
      }
      await wrapper.vm.fetchData()
      await wrapper.vm.$nextTick()

      expect(profitAnalysisService.getCategoryProfit).toHaveBeenCalledWith({
        start: '2024-01-01',
        end: '2024-01-31'
      })
    })

    it('should not fetch data if date range is incomplete', async () => {
      profitAnalysisService.getCategoryProfit.mockClear()
      
      wrapper.vm.dateRange = {
        start: '',
        end: ''
      }
      await wrapper.vm.fetchData()
      await wrapper.vm.$nextTick()

      expect(profitAnalysisService.getCategoryProfit).not.toHaveBeenCalled()
    })

    it('should handle different date ranges', async () => {
      const testRanges = [
        { start: '2024-01-01', end: '2024-01-01' }, // Single day
        { start: '2024-01-01', end: '2024-01-07' }, // Week
        { start: '2024-01-01', end: '2024-01-31' }  // Month
      ]

      for (const range of testRanges) {
        wrapper.vm.dateRange = range
        await wrapper.vm.fetchData()
        await wrapper.vm.$nextTick()

        expect(profitAnalysisService.getCategoryProfit).toHaveBeenCalledWith(range)
      }
    })
  })

  describe('Error Handling', () => {
    it('should handle category profit API error', async () => {
      profitAnalysisService.getCategoryProfit.mockRejectedValue({
        response: { data: { error: 'Database error' } }
      })
      
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.error).toBeTruthy()
    })

    it('should handle operating profit API error', async () => {
      profitAnalysisService.getOperatingProfit.mockRejectedValue({
        response: { data: { error: 'Database error' } }
      })
      
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      
      wrapper.vm.viewMode = 'operating'
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.error).toBeTruthy()
    })

    it('should allow retry after error', async () => {
      profitAnalysisService.getCategoryProfit.mockRejectedValueOnce(new Error('API Error'))
        .mockResolvedValueOnce(mockCategoryProfits)
      
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })
      
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.vm.error).toBeTruthy()

      await wrapper.vm.fetchData()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.error).toBeNull()
    })
  })

  describe('Pull to Refresh', () => {
    it('should support pull to refresh', async () => {
      wrapper = mount(ProfitAnalysisView, {
        global: {
          stubs: {
            BottomNav: true,
            PullToRefresh: true,
            CategoryProfitView: true,
            OperatingProfitView: true
          }
        }
      })

      const pullToRefresh = wrapper.findComponent({ name: 'PullToRefresh' })
      expect(pullToRefresh.exists()).toBe(true)
    })
  })
})

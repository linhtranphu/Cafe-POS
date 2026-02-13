/**
 * Unit Tests for ProfitComparisonModal Component
 * 
 * Requirements tested:
 * - AC-12.2: Show cost difference between sizes
 * - AC-12.3: Show profit margin difference between sizes
 * - AC-12.4: Highlight most profitable variant
 * 
 * Task: 11a.3 Implement profit comparison view
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ProfitComparisonModal from '../ProfitComparisonModal.vue'
import { menuService } from '../../services/menu'

// Mock the menu service
vi.mock('../../services/menu', () => ({
  menuService: {
    getMenuItem: vi.fn()
  }
}))

describe('ProfitComparisonModal', () => {
  let wrapper

  const mockSingleSizeItem = {
    data: {
      id: '1',
      name: 'Bánh mì thịt',
      category: 'Món ăn',
      has_variants: false,
      price: 20000,
      current_cost: 12000,
      cost_status: 'FINAL'
    }
  }

  const mockMultiSizeItem = {
    data: {
      id: '2',
      name: 'Cà phê sữa đá',
      category: 'Cà phê',
      has_variants: true,
      variants: [
        {
          id: 'M',
          name: 'Size M',
          price: 25000,
          current_cost: 15000,
          cost_status: 'FINAL',
          is_default: true,
          available: true
        },
        {
          id: 'L',
          name: 'Size L',
          price: 30000,
          current_cost: 18000,
          cost_status: 'FINAL',
          is_default: false,
          available: true
        },
        {
          id: 'XL',
          name: 'Size XL',
          price: 35000,
          current_cost: 20000,
          cost_status: 'FINAL',
          is_default: false,
          available: true
        }
      ]
    }
  }

  const mockUnprofitableItem = {
    data: {
      id: '3',
      name: 'Món lỗ',
      category: 'Test',
      has_variants: true,
      variants: [
        {
          id: 'M',
          name: 'Size M',
          price: 10000,
          current_cost: 15000,
          cost_status: 'FINAL',
          is_default: true,
          available: true
        },
        {
          id: 'L',
          name: 'Size L',
          price: 15000,
          current_cost: 20000,
          cost_status: 'FINAL',
          is_default: false,
          available: true
        }
      ]
    }
  }

  beforeEach(() => {
    menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
  })

  describe('Modal Visibility', () => {
    it('should not render when isOpen is false', () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: false,
          menuItemId: '2'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(false)
    })

    it('should render modal when isOpen is true', () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(true)
      expect(wrapper.text()).toContain('So sánh lợi nhuận theo size')
    })

    it('should emit close event when close button is clicked', async () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const closeButtons = wrapper.findAll('button')
      const headerCloseButton = closeButtons[0]
      await headerCloseButton.trigger('click')

      expect(wrapper.emitted('close')).toBeTruthy()
    })

    it('should emit close event when clicking backdrop', async () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()

      const backdrop = wrapper.find('.fixed')
      await backdrop.trigger('click')

      expect(wrapper.emitted('close')).toBeTruthy()
    })
  })

  describe('Multi-Size Item Display (AC-12.2, AC-12.3, AC-12.4)', () => {
    beforeEach(async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display menu item name', () => {
      expect(wrapper.text()).toContain('Cà phê sữa đá')
    })

    it('should display summary stats', () => {
      expect(wrapper.text()).toContain('Tổng số size')
      expect(wrapper.text()).toContain('3')
      expect(wrapper.text()).toContain('Size lời nhất')
    })

    it('should display all variants in comparison table', () => {
      expect(wrapper.text()).toContain('Size M')
      expect(wrapper.text()).toContain('Size L')
      expect(wrapper.text()).toContain('Size XL')
    })

    it('should display price for each variant', () => {
      expect(wrapper.text()).toContain('25,000')
      expect(wrapper.text()).toContain('30,000')
      expect(wrapper.text()).toContain('35,000')
    })

    it('should display cost for each variant', () => {
      expect(wrapper.text()).toContain('15,000')
      expect(wrapper.text()).toContain('18,000')
      expect(wrapper.text()).toContain('20,000')
    })

    it('should display profit for each variant (AC-12.2)', () => {
      // Size M: 25000 - 15000 = 10000
      // Size L: 30000 - 18000 = 12000
      // Size XL: 35000 - 20000 = 15000
      expect(wrapper.text()).toContain('10,000')
      expect(wrapper.text()).toContain('12,000')
      expect(wrapper.text()).toContain('15,000')
    })

    it('should display profit margin for each variant (AC-12.3)', () => {
      // Size M: (10000/25000) * 100 = 40%
      // Size L: (12000/30000) * 100 = 40%
      // Size XL: (15000/35000) * 100 = 42.86%
      expect(wrapper.text()).toContain('40')
      expect(wrapper.text()).toContain('42')
    })

    it('should highlight most profitable variant (AC-12.4)', () => {
      // Size XL has highest profit (15000)
      expect(wrapper.text()).toContain('🏆 Lời nhất')
    })

    it('should mark default variant', () => {
      expect(wrapper.text()).toContain('Mặc định')
    })

    it('should display cost status for each variant', () => {
      expect(wrapper.text()).toContain('✓ Chính thức')
    })
  })

  describe('Cost Difference Analysis (AC-12.2, AC-12.3)', () => {
    beforeEach(async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display cost difference section', () => {
      expect(wrapper.text()).toContain('Phân tích chênh lệch')
    })

    it('should show price difference between sizes (AC-12.2)', () => {
      expect(wrapper.text()).toContain('Chênh lệch giá')
      // M to L: 30000 - 25000 = +5000
      // L to XL: 35000 - 30000 = +5000
      expect(wrapper.text()).toContain('+5,000')
    })

    it('should show cost difference between sizes (AC-12.2)', () => {
      expect(wrapper.text()).toContain('Chênh lệch chi phí')
      // M to L: 18000 - 15000 = +3000
      // L to XL: 20000 - 18000 = +2000
      expect(wrapper.text()).toContain('+3,000')
      expect(wrapper.text()).toContain('+2,000')
    })

    it('should show profit difference between sizes (AC-12.2)', () => {
      expect(wrapper.text()).toContain('Chênh lệch lợi nhuận')
      // M to L: 12000 - 10000 = +2000
      // L to XL: 15000 - 12000 = +3000
      expect(wrapper.text()).toContain('+2,000')
      expect(wrapper.text()).toContain('+3,000')
    })

    it('should show profit margin difference between sizes (AC-12.3)', () => {
      expect(wrapper.text()).toContain('Chênh lệch tỷ suất LN')
    })

    it('should display comparison pairs', () => {
      expect(wrapper.text()).toContain('Size M → Size L')
      expect(wrapper.text()).toContain('Size L → Size XL')
    })
  })

  describe('Insights Generation', () => {
    it('should show insights for profitable items', async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Nhận xét')
      expect(wrapper.text()).toContain('✅ Tất cả các size đều có lợi nhuận')
    })

    it('should warn about unprofitable variants', async () => {
      menuService.getMenuItem.mockResolvedValue(mockUnprofitableItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '3'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('⚠️')
      expect(wrapper.text()).toContain('size đang bị lỗ')
    })

    it('should identify most profitable variant in insights', async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Size XL')
      expect(wrapper.text()).toContain('có lợi nhuận cao nhất')
    })
  })

  describe('Single-Size Item Handling', () => {
    it('should show message for single-size items', async () => {
      menuService.getMenuItem.mockResolvedValue(mockSingleSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Món này không có nhiều size')
    })
  })

  describe('Loading and Error States', () => {
    it('should show loading state initially', () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      expect(wrapper.find('.animate-spin').exists()).toBe(true)
    })

    it('should show error state when API fails', async () => {
      menuService.getMenuItem.mockRejectedValue({
        response: { data: { error: 'API Error' } }
      })

      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('API Error')
    })

    it('should show generic error when no error message provided', async () => {
      menuService.getMenuItem.mockRejectedValue(new Error('Network error'))

      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải dữ liệu món')
    })
  })

  describe('Data Fetching', () => {
    it('should fetch menu item when modal opens', async () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItem).toHaveBeenCalledWith('2')
    })

    it('should not fetch when menuItemId is null', async () => {
      menuService.getMenuItem.mockClear()

      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: null
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItem).not.toHaveBeenCalled()
    })

    it('should refetch when modal reopens', async () => {
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: false,
          menuItemId: '2'
        }
      })

      await wrapper.vm.$nextTick()

      menuService.getMenuItem.mockClear()
      await wrapper.setProps({ isOpen: true })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getMenuItem).toHaveBeenCalledWith('2')
    })
  })

  describe('Variant Sorting', () => {
    it('should sort variants by profit (descending)', async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const { sortedVariants } = wrapper.vm
      expect(sortedVariants[0].id).toBe('XL') // Highest profit: 15000
      expect(sortedVariants[1].id).toBe('L')  // Middle profit: 12000
      expect(sortedVariants[2].id).toBe('M')  // Lowest profit: 10000
    })
  })

  describe('Helper Functions', () => {
    beforeEach(async () => {
      menuService.getMenuItem.mockResolvedValue(mockMultiSizeItem)
      wrapper = mount(ProfitComparisonModal, {
        props: {
          isOpen: true,
          menuItemId: '2'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should calculate profit correctly', () => {
      const { calculateProfit } = wrapper.vm
      const variant = { price: 25000, current_cost: 15000 }
      expect(calculateProfit(variant)).toBe(10000)
    })

    it('should calculate profit margin correctly', () => {
      const { calculateProfitMargin } = wrapper.vm
      const variant = { price: 25000, current_cost: 15000 }
      expect(calculateProfitMargin(variant)).toBe(40)
    })

    it('should handle zero price in profit margin calculation', () => {
      const { calculateProfitMargin } = wrapper.vm
      const variant = { price: 0, current_cost: 15000 }
      expect(calculateProfitMargin(variant)).toBe(0)
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
  })
})

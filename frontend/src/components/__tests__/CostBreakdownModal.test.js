/**
 * Unit Tests for CostBreakdownModal Component
 * 
 * Requirements tested:
 * - AC-10.4: Display detailed ingredient costs per variant
 * - AC-11.1-AC-11.5: Display conversion rates, wastage percentages, and formula breakdown
 * 
 * Task: 11a.2 Implement cost breakdown modal
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import CostBreakdownModal from '../CostBreakdownModal.vue'
import { menuService } from '../../services/menu'

// Mock the menu service
vi.mock('../../services/menu', () => ({
  menuService: {
    getCostBreakdown: vi.fn()
  }
}))

describe('CostBreakdownModal', () => {
  let wrapper

  const mockSingleSizeCostBreakdown = {
    data: {
      menu_item_id: '1',
      menu_item_name: 'Bánh mì thịt',
      has_variants: false,
      price: 20000,
      total_cost: 12000,
      cost_status: 'FINAL',
      ingredients: [
        {
          name: 'Bánh mì',
          quantity: 1,
          unit: 'cái',
          cost_per_unit: 5000,
          conversion_rate: 1.0,
          wastage_percentage: 5,
          total_cost: 5250
        },
        {
          name: 'Thịt',
          quantity: 50,
          unit: 'g',
          cost_per_unit: 120,
          conversion_rate: 0.001,
          wastage_percentage: 10,
          total_cost: 6600
        }
      ]
    }
  }

  const mockMultiSizeCostBreakdown = {
    data: {
      menu_item_id: '2',
      menu_item_name: 'Cà phê sữa đá',
      has_variants: true,
      variants: [
        {
          variant_id: 'M',
          variant_name: 'Size M',
          price: 25000,
          total_cost: 15000,
          cost_status: 'FINAL',
          cost_last_calculated_at: '2026-02-13T10:00:00Z',
          ingredients: [
            {
              name: 'Cà phê',
              quantity: 20,
              unit: 'g',
              cost_per_unit: 500,
              conversion_rate: 0.001,
              wastage_percentage: 5,
              total_cost: 10500
            },
            {
              name: 'Sữa đặc',
              quantity: 30,
              unit: 'ml',
              cost_per_unit: 150,
              conversion_rate: 1.0,
              wastage_percentage: 0,
              total_cost: 4500
            }
          ]
        },
        {
          variant_id: 'L',
          variant_name: 'Size L',
          price: 30000,
          total_cost: 22000,
          cost_status: 'FINAL',
          cost_last_calculated_at: '2026-02-13T10:00:00Z',
          ingredients: [
            {
              name: 'Cà phê',
              quantity: 30,
              unit: 'g',
              cost_per_unit: 500,
              conversion_rate: 0.001,
              wastage_percentage: 5,
              total_cost: 15750
            },
            {
              name: 'Sữa đặc',
              quantity: 45,
              unit: 'ml',
              cost_per_unit: 150,
              conversion_rate: 1.0,
              wastage_percentage: 0,
              total_cost: 6750
            }
          ]
        }
      ]
    }
  }

  const mockIncompleteCostBreakdown = {
    data: {
      menu_item_id: '3',
      menu_item_name: 'Món thiếu giá',
      has_variants: false,
      price: 30000,
      total_cost: 5000,
      cost_status: 'INCOMPLETE',
      ingredients: [
        {
          name: 'Nguyên liệu A',
          quantity: 10,
          unit: 'g',
          cost_per_unit: 500,
          conversion_rate: 1.0,
          wastage_percentage: 0,
          total_cost: 5000
        },
        {
          name: 'Nguyên liệu B',
          quantity: 20,
          unit: 'g',
          cost_per_unit: 0,
          conversion_rate: 1.0,
          wastage_percentage: 0,
          total_cost: 0
        }
      ]
    }
  }

  beforeEach(() => {
    menuService.getCostBreakdown.mockResolvedValue(mockSingleSizeCostBreakdown)
  })

  describe('Modal Visibility', () => {
    it('should not render when isOpen is false', () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: false,
          menuItemId: '1'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(false)
    })

    it('should render modal when isOpen is true', () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      expect(wrapper.find('.fixed').exists()).toBe(true)
      expect(wrapper.text()).toContain('Chi tiết chi phí')
    })

    it('should emit close event when close button is clicked', async () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
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
      wrapper = mount(CostBreakdownModal, {
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

  describe('Single-Size Item Display (AC-10.4, AC-11.1-AC-11.5)', () => {
    beforeEach(async () => {
      menuService.getCostBreakdown.mockResolvedValue(mockSingleSizeCostBreakdown)
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display menu item name and type', () => {
      expect(wrapper.text()).toContain('Bánh mì thịt')
      expect(wrapper.text()).toContain('Một size')
    })

    it('should display price and total cost', () => {
      expect(wrapper.text()).toContain('20,000')
      expect(wrapper.text()).toContain('12,000')
    })

    it('should display all ingredients', () => {
      expect(wrapper.text()).toContain('Bánh mì')
      expect(wrapper.text()).toContain('Thịt')
    })

    it('should display ingredient quantities and units', () => {
      expect(wrapper.text()).toContain('1 cái')
      expect(wrapper.text()).toContain('50 g')
    })

    it('should display cost per unit for each ingredient', () => {
      expect(wrapper.text()).toContain('5,000')
      expect(wrapper.text()).toContain('120')
    })

    it('should display conversion rates', () => {
      expect(wrapper.text()).toContain('1.0000')
      expect(wrapper.text()).toContain('0.0010')
    })

    it('should display wastage percentages', () => {
      expect(wrapper.text()).toContain('5%')
      expect(wrapper.text()).toContain('10%')
    })

    it('should display formula breakdown for each ingredient', () => {
      const text = wrapper.text()
      expect(text).toContain('Công thức')
      expect(text).toContain('5,000')
      expect(text).toContain('120')
      expect(text).toContain('1.0000')
      expect(text).toContain('0.0010')
      expect(text).toContain('(1 + 5/100)')
      expect(text).toContain('(1 + 10/100)')
    })

    it('should display total cost for each ingredient', () => {
      expect(wrapper.text()).toContain('5,250')
      expect(wrapper.text()).toContain('6,600')
    })
  })

  describe('Multi-Size Item Display (AC-10.4, AC-11.1-AC-11.5)', () => {
    beforeEach(async () => {
      menuService.getCostBreakdown.mockResolvedValue(mockMultiSizeCostBreakdown)
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '2',
          defaultVariantId: 'M'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should display menu item name and type', () => {
      expect(wrapper.text()).toContain('Cà phê sữa đá')
      expect(wrapper.text()).toContain('Nhiều size')
    })

    it('should display all variants', () => {
      expect(wrapper.text()).toContain('Size M')
      expect(wrapper.text()).toContain('Size L')
    })

    it('should highlight default variant', () => {
      expect(wrapper.text()).toContain('Mặc định')
    })

    it('should display price and cost for each variant', () => {
      expect(wrapper.text()).toContain('25,000')
      expect(wrapper.text()).toContain('15,000')
      expect(wrapper.text()).toContain('30,000')
      expect(wrapper.text()).toContain('22,000')
    })

    it('should display ingredients for each variant', () => {
      expect(wrapper.text()).toContain('Cà phê')
      expect(wrapper.text()).toContain('Sữa đặc')
    })

    it('should display different quantities for different variants', () => {
      expect(wrapper.text()).toContain('20 g')
      expect(wrapper.text()).toContain('30 g')
      expect(wrapper.text()).toContain('30 ml')
      expect(wrapper.text()).toContain('45 ml')
    })

    it('should display formula breakdown for each variant ingredient', () => {
      const text = wrapper.text()
      expect(text).toContain('Công thức')
      expect(text).toContain('500')
      expect(text).toContain('150')
      expect(text).toContain('0.0010')
      expect(text).toContain('1.0000')
      expect(text).toContain('(1 + 5/100)')
      expect(text).toContain('(1 + 0/100)')
    })

    it('should display cost status for each variant', () => {
      expect(wrapper.text()).toContain('✓ Chính thức')
    })
  })

  describe('Cost Status Display', () => {
    it('should display FINAL status correctly', async () => {
      menuService.getCostBreakdown.mockResolvedValue(mockSingleSizeCostBreakdown)
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('✓ Chính thức')
    })

    it('should display INCOMPLETE status correctly', async () => {
      menuService.getCostBreakdown.mockResolvedValue(mockIncompleteCostBreakdown)
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '3'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('⚠ Thiếu dữ liệu')
    })
  })

  describe('Loading and Error States', () => {
    it('should show loading state initially', () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      expect(wrapper.find('.animate-pulse').exists()).toBe(true)
    })

    it('should show error state when API fails', async () => {
      menuService.getCostBreakdown.mockRejectedValue({
        response: { data: { error: 'API Error' } }
      })

      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('API Error')
    })

    it('should show generic error when no error message provided', async () => {
      menuService.getCostBreakdown.mockRejectedValue(new Error('Network error'))

      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(wrapper.text()).toContain('Không thể tải chi tiết chi phí')
    })
  })

  describe('Data Fetching', () => {
    it('should fetch cost breakdown when modal opens', async () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getCostBreakdown).toHaveBeenCalledWith('1')
    })

    it('should not fetch when menuItemId is null', async () => {
      menuService.getCostBreakdown.mockClear()

      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: null
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getCostBreakdown).not.toHaveBeenCalled()
    })

    it('should refetch when menuItemId changes', async () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getCostBreakdown).toHaveBeenCalledWith('1')

      menuService.getCostBreakdown.mockClear()
      await wrapper.setProps({ menuItemId: '2' })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getCostBreakdown).toHaveBeenCalledWith('2')
    })

    it('should refetch when modal reopens', async () => {
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: false,
          menuItemId: '1'
        }
      })

      await wrapper.vm.$nextTick()

      menuService.getCostBreakdown.mockClear()
      await wrapper.setProps({ isOpen: true })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      expect(menuService.getCostBreakdown).toHaveBeenCalledWith('1')
    })
  })

  describe('Helper Functions', () => {
    beforeEach(async () => {
      menuService.getCostBreakdown.mockResolvedValue(mockSingleSizeCostBreakdown)
      wrapper = mount(CostBreakdownModal, {
        props: {
          isOpen: true,
          menuItemId: '1'
        }
      })
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
    })

    it('should format prices correctly', () => {
      const { formatPrice } = wrapper.vm
      // Just check that the function returns a string with the number and currency symbol
      expect(formatPrice(20000)).toContain('20')
      expect(formatPrice(20000)).toContain('₫')
      expect(formatPrice(12500)).toContain('12')
      expect(formatPrice(0)).toContain('0')
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
  })
})

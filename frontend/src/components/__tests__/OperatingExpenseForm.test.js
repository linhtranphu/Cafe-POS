/**
 * Unit Tests for OperatingExpenseForm Component
 * 
 * Requirements tested:
 * - 6.5.2: Operating expense form with validation and save functionality
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import OperatingExpenseForm from '../OperatingExpenseForm.vue'
import { profitAnalysisService } from '../../services/profitAnalysis'

// Mock the profit analysis service
vi.mock('../../services/profitAnalysis', () => ({
  profitAnalysisService: {
    createOperatingExpense: vi.fn()
  }
}))

// Mock formatters
vi.mock('../../utils/formatters', () => ({
  formatPrice: (value) => {
    if (value === undefined || value === null || isNaN(value)) return '0'
    return value.toLocaleString('vi-VN')
  }
}))

describe('OperatingExpenseForm', () => {
  let wrapper

  const mockExpenseData = {
    id: '1',
    period_start: '2024-01-01',
    period_end: '2024-01-31',
    staff_salary: 10000000,
    rent: 5000000,
    utilities: 1000000,
    marketing_costs: 500000,
    other_expenses: 300000,
    total_expenses: 16800000
  }

  beforeEach(() => {
    profitAnalysisService.createOperatingExpense.mockResolvedValue(mockExpenseData)
  })

  describe('Component Rendering', () => {
    it('should render form with all input fields', () => {
      wrapper = mount(OperatingExpenseForm)

      expect(wrapper.text()).toContain('Chi phí vận hành')
      expect(wrapper.text()).toContain('Khoảng thời gian')
      expect(wrapper.text()).toContain('Ngày bắt đầu')
      expect(wrapper.text()).toContain('Ngày kết thúc')
      expect(wrapper.text()).toContain('Lương nhân viên')
      expect(wrapper.text()).toContain('Tiền thuê mặt bằng')
      expect(wrapper.text()).toContain('Điện nước')
      expect(wrapper.text()).toContain('Chi phí marketing')
      expect(wrapper.text()).toContain('Chi phí khác')
      expect(wrapper.text()).toContain('Tổng chi phí vận hành')
    })

    it('should render save and cancel buttons', () => {
      wrapper = mount(OperatingExpenseForm)

      const buttons = wrapper.findAll('button[type="button"], button[type="submit"]')
      expect(buttons.length).toBeGreaterThanOrEqual(2)
      expect(wrapper.text()).toContain('Lưu')
      expect(wrapper.text()).toContain('Hủy')
    })

    it('should initialize with empty form when no initialData provided', () => {
      wrapper = mount(OperatingExpenseForm)

      expect(wrapper.vm.formData.period_start).toBe('')
      expect(wrapper.vm.formData.period_end).toBe('')
      expect(wrapper.vm.formData.staff_salary).toBe(0)
      expect(wrapper.vm.formData.rent).toBe(0)
      expect(wrapper.vm.formData.utilities).toBe(0)
      expect(wrapper.vm.formData.marketing_costs).toBe(0)
      expect(wrapper.vm.formData.other_expenses).toBe(0)
    })

    it('should initialize with provided initialData', () => {
      wrapper = mount(OperatingExpenseForm, {
        props: {
          initialData: mockExpenseData
        }
      })

      expect(wrapper.vm.formData.period_start).toBe('2024-01-01')
      expect(wrapper.vm.formData.period_end).toBe('2024-01-31')
      expect(wrapper.vm.formData.staff_salary).toBe(10000000)
      expect(wrapper.vm.formData.rent).toBe(5000000)
      expect(wrapper.vm.formData.utilities).toBe(1000000)
      expect(wrapper.vm.formData.marketing_costs).toBe(500000)
      expect(wrapper.vm.formData.other_expenses).toBe(300000)
    })
  })

  describe('Total Calculation (Requirement 6.5.2)', () => {
    it('should auto-calculate total expenses', async () => {
      wrapper = mount(OperatingExpenseForm)

      // Set expense values
      wrapper.vm.formData.staff_salary = 10000000
      wrapper.vm.formData.rent = 5000000
      wrapper.vm.formData.utilities = 1000000
      wrapper.vm.formData.marketing_costs = 500000
      wrapper.vm.formData.other_expenses = 300000

      await wrapper.vm.$nextTick()

      expect(wrapper.vm.totalExpenses).toBe(16800000)
    })

    it('should update total when any expense changes', async () => {
      wrapper = mount(OperatingExpenseForm)

      wrapper.vm.formData.staff_salary = 5000000
      await wrapper.vm.$nextTick()
      expect(wrapper.vm.totalExpenses).toBe(5000000)

      wrapper.vm.formData.rent = 3000000
      await wrapper.vm.$nextTick()
      expect(wrapper.vm.totalExpenses).toBe(8000000)

      wrapper.vm.formData.utilities = 500000
      await wrapper.vm.$nextTick()
      expect(wrapper.vm.totalExpenses).toBe(8500000)
    })

    it('should handle zero values in total calculation', async () => {
      wrapper = mount(OperatingExpenseForm)

      wrapper.vm.formData.staff_salary = 0
      wrapper.vm.formData.rent = 0
      wrapper.vm.formData.utilities = 0
      wrapper.vm.formData.marketing_costs = 0
      wrapper.vm.formData.other_expenses = 0

      await wrapper.vm.$nextTick()

      expect(wrapper.vm.totalExpenses).toBe(0)
    })

    it('should display formatted total expenses', async () => {
      wrapper = mount(OperatingExpenseForm)

      wrapper.vm.formData.staff_salary = 10000000
      wrapper.vm.formData.rent = 5000000

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('15,000,000')
    })
  })

  describe('Form Validation (Requirement 6.5.2)', () => {
    beforeEach(() => {
      wrapper = mount(OperatingExpenseForm)
    })

    it('should validate period_start is required', async () => {
      wrapper.vm.formData.period_start = ''
      wrapper.vm.formData.period_end = '2024-01-31'

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(false)
      expect(wrapper.vm.errors.period_start).toBe('Vui lòng chọn ngày bắt đầu')
    })

    it('should validate period_end is required', async () => {
      wrapper.vm.formData.period_start = '2024-01-01'
      wrapper.vm.formData.period_end = ''

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(false)
      expect(wrapper.vm.errors.period_end).toBe('Vui lòng chọn ngày kết thúc')
    })

    it('should validate period_start <= period_end', async () => {
      wrapper.vm.formData.period_start = '2024-01-31'
      wrapper.vm.formData.period_end = '2024-01-01'

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(false)
      expect(wrapper.vm.errors.period_end).toBe('Ngày kết thúc phải sau ngày bắt đầu')
    })

    it('should accept equal start and end dates', async () => {
      wrapper.vm.formData.period_start = '2024-01-15'
      wrapper.vm.formData.period_end = '2024-01-15'
      wrapper.vm.formData.staff_salary = 1000000

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(true)
      expect(wrapper.vm.errors.period_end).toBeUndefined()
    })

    it('should validate all amounts >= 0', async () => {
      wrapper.vm.formData.period_start = '2024-01-01'
      wrapper.vm.formData.period_end = '2024-01-31'
      wrapper.vm.formData.staff_salary = -1000

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(false)
      expect(wrapper.vm.errors.staff_salary).toBe('Số tiền không được âm')
    })

    it('should validate multiple negative amounts', async () => {
      wrapper.vm.formData.period_start = '2024-01-01'
      wrapper.vm.formData.period_end = '2024-01-31'
      wrapper.vm.formData.staff_salary = -1000
      wrapper.vm.formData.rent = -500
      wrapper.vm.formData.utilities = -100

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(false)
      expect(wrapper.vm.errors.staff_salary).toBe('Số tiền không được âm')
      expect(wrapper.vm.errors.rent).toBe('Số tiền không được âm')
      expect(wrapper.vm.errors.utilities).toBe('Số tiền không được âm')
    })

    it('should pass validation with valid data', async () => {
      wrapper.vm.formData.period_start = '2024-01-01'
      wrapper.vm.formData.period_end = '2024-01-31'
      wrapper.vm.formData.staff_salary = 10000000
      wrapper.vm.formData.rent = 5000000
      wrapper.vm.formData.utilities = 1000000
      wrapper.vm.formData.marketing_costs = 500000
      wrapper.vm.formData.other_expenses = 300000

      const isValid = wrapper.vm.validateForm()

      expect(isValid).toBe(true)
      expect(Object.keys(wrapper.vm.errors).length).toBe(0)
    })

    it('should display validation errors in UI', async () => {
      wrapper.vm.formData.period_start = ''
      wrapper.vm.formData.period_end = '2024-01-01'
      wrapper.vm.formData.staff_salary = -1000

      wrapper.vm.validateForm()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Vui lòng chọn ngày bắt đầu')
      expect(wrapper.text()).toContain('Số tiền không được âm')
    })

    it('should clear errors when form data changes', async () => {
      wrapper.vm.errors = { period_start: 'Error message' }
      await wrapper.vm.$nextTick()

      wrapper.vm.formData.period_start = '2024-01-01'
      await wrapper.vm.$nextTick()

      expect(Object.keys(wrapper.vm.errors).length).toBe(0)
    })
  })

  describe('Save Action (Requirement 6.5.2)', () => {
    beforeEach(() => {
      wrapper = mount(OperatingExpenseForm)
    })

    it('should call API with correct data on submit', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()

      expect(profitAnalysisService.createOperatingExpense).toHaveBeenCalledWith({
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      })
    })

    it('should emit save event with result on success', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toEqual(mockExpenseData)
    })

    it('should not submit if validation fails', async () => {
      wrapper.vm.formData.period_start = ''
      wrapper.vm.formData.period_end = '2024-01-31'

      await wrapper.vm.handleSubmit()

      expect(profitAnalysisService.createOperatingExpense).not.toHaveBeenCalled()
      expect(wrapper.emitted('save')).toBeFalsy()
    })

    it('should show loading state while saving', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      const savePromise = wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.saving).toBe(true)
      expect(wrapper.text()).toContain('Đang lưu...')

      await savePromise
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.saving).toBe(false)
    })

    it('should disable submit button while saving', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      const savePromise = wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.attributes('disabled')).toBeDefined()

      await savePromise
    })

    it('should handle zero values correctly', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 0,
        rent: 0,
        utilities: 0,
        marketing_costs: 0,
        other_expenses: 0
      }

      await wrapper.vm.handleSubmit()

      expect(profitAnalysisService.createOperatingExpense).toHaveBeenCalledWith({
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 0,
        rent: 0,
        utilities: 0,
        marketing_costs: 0,
        other_expenses: 0
      })
    })
  })

  describe('Error Handling (Requirement 6.5.2)', () => {
    beforeEach(() => {
      wrapper = mount(OperatingExpenseForm)
    })

    it('should display error message when API fails', async () => {
      const errorMessage = 'Failed to save expense'
      profitAnalysisService.createOperatingExpense.mockRejectedValue({
        response: {
          data: {
            error: errorMessage
          }
        }
      })

      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.errors.submit).toBe(errorMessage)
      expect(wrapper.text()).toContain(errorMessage)
    })

    it('should display generic error when no error message provided', async () => {
      profitAnalysisService.createOperatingExpense.mockRejectedValue(new Error('Network error'))

      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.errors.submit).toBe('Không thể lưu chi phí vận hành')
    })

    it('should not emit save event on error', async () => {
      profitAnalysisService.createOperatingExpense.mockRejectedValue(new Error('API Error'))

      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.emitted('save')).toBeFalsy()
    })

    it('should reset saving state after error', async () => {
      profitAnalysisService.createOperatingExpense.mockRejectedValue(new Error('API Error'))

      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      await wrapper.vm.handleSubmit()
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.saving).toBe(false)
    })
  })

  describe('Cancel Action', () => {
    beforeEach(() => {
      wrapper = mount(OperatingExpenseForm)
    })

    it('should emit cancel event when cancel button clicked', async () => {
      const cancelButtons = wrapper.findAll('button[type="button"]')
      const cancelButton = cancelButtons.find(btn => btn.text().includes('Hủy'))
      
      await cancelButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('should emit cancel event when close button clicked', async () => {
      const closeButton = wrapper.find('button.text-gray-400')
      
      await closeButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('should not save data when cancel is clicked', async () => {
      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      const cancelButtons = wrapper.findAll('button[type="button"]')
      const cancelButton = cancelButtons.find(btn => btn.text().includes('Hủy'))
      
      await cancelButton.trigger('click')

      expect(profitAnalysisService.createOperatingExpense).not.toHaveBeenCalled()
      expect(wrapper.emitted('save')).toBeFalsy()
    })
  })

  describe('Form Submission via Enter Key', () => {
    it('should submit form when pressing enter in form', async () => {
      wrapper = mount(OperatingExpenseForm)

      wrapper.vm.formData = {
        period_start: '2024-01-01',
        period_end: '2024-01-31',
        staff_salary: 10000000,
        rent: 5000000,
        utilities: 1000000,
        marketing_costs: 500000,
        other_expenses: 300000
      }

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(profitAnalysisService.createOperatingExpense).toHaveBeenCalled()
    })
  })
})

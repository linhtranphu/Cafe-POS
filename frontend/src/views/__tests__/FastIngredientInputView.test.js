import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import FastIngredientInputView from '../FastIngredientInputView.vue'

vi.mock('../../services/ingredient', () => ({
  ingredientService: {
    getRecentRestocked: vi.fn(),
    getIngredients: vi.fn()
  }
}))

vi.mock('../../services/fundIngredientService', () => ({
  fundIngredientService: {
    restockIngredientFromFund: vi.fn()
  }
}))

import { ingredientService } from '../../services/ingredient'
import { fundIngredientService } from '../../services/fundIngredientService'

const mockIngredient = {
  id: 'abc123',
  name: 'Sữa tươi',
  unit: 'L',
  quantity: 5.0,
  cost_per_unit: 28000,
  created_at: '2026-04-01T00:00:00Z',
  last_restock: {
    quantity: 10,
    cost_per_unit: 27500,
    created_at: '2026-04-21T08:00:00Z'
  }
}

const mountView = () => mount(FastIngredientInputView, {
  global: {
    stubs: { RouterLink: true },
    mocks: { $router: { back: vi.fn(), push: vi.fn() } }
  }
})

describe('FastIngredientInputView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    ingredientService.getRecentRestocked.mockResolvedValue([mockIngredient])
    ingredientService.getIngredients.mockResolvedValue([mockIngredient])
    fundIngredientService.restockIngredientFromFund.mockResolvedValue({ success: true })
  })

  it('renders the recent-restocked list', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    expect(ingredientService.getRecentRestocked).toHaveBeenCalledWith(30)
    expect(wrapper.text()).toContain('Sữa tươi')
  })

  it('fetches all ingredients when Hiện tất cả is toggled', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const toggleBtn = wrapper.findAll('button').find(b => b.text().includes('Hiện tất cả'))
    await toggleBtn.trigger('click')
    await nextTick()
    await nextTick()
    expect(ingredientService.getIngredients).toHaveBeenCalled()
  })

  it('filters display list by search query', async () => {
    ingredientService.getRecentRestocked.mockResolvedValue([
      mockIngredient,
      { ...mockIngredient, id: 'xyz', name: 'Đường', last_restock: null }
    ])
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    await wrapper.find('input[placeholder*="Tìm"]').setValue('Sữa')
    await nextTick()
    expect(wrapper.text()).toContain('Sữa tươi')
    expect(wrapper.text()).not.toContain('Đường')
  })

  it('initializes quantity=0 and cost from last_restock when ingredient selected', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const inputs = wrapper.findAll('input[type="number"]')
    expect(Number(inputs[0].element.value)).toBe(0)
    expect(Number(inputs[1].element.value)).toBe(27500)
  })

  it('increments quantity by last_restock.quantity per tap', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const plusBtn = wrapper.findAll('button').find(b => b.text() === '+')
    await plusBtn.trigger('click')
    await nextTick()
    const qtyInput = wrapper.find('input[type="number"]')
    expect(Number(qtyInput.element.value)).toBe(10)
  })

  it('decrement does not go below 0', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const minusBtn = wrapper.findAll('button').find(b => b.text() === '−')
    await minusBtn.trigger('click')
    await nextTick()
    expect(Number(wrapper.find('input[type="number"]').element.value)).toBe(0)
  })

  it('submit calls restockIngredientFromFund with correct payload', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === '+').trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === 'Lưu').trigger('click')
    await nextTick()
    expect(fundIngredientService.restockIngredientFromFund).toHaveBeenCalledWith(
      'abc123',
      expect.objectContaining({
        quantity: 10,
        cost_per_unit: 27500,
        reason: 'Nhập nhanh',
        money_type: 'cash'
      })
    )
  })

  it('submit success clears active selection and refreshes list', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === '+').trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === 'Lưu').trigger('click')
    await nextTick()
    await nextTick()
    expect(ingredientService.getRecentRestocked).toHaveBeenCalledTimes(2)
    expect(wrapper.findAll('input[type="number"]').length).toBe(0)
  })

  it('Lưu button is disabled when quantity is 0', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const saveBtn = wrapper.findAll('button').find(b => b.text() === 'Lưu')
    expect(saveBtn.attributes('disabled')).toBeDefined()
    await saveBtn.trigger('click')
    await nextTick()
    expect(fundIngredientService.restockIngredientFromFund).not.toHaveBeenCalled()
  })
})

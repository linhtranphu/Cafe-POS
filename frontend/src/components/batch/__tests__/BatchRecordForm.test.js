import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import BatchRecordForm from '../BatchRecordForm.vue'
import { useBatchRecordStore } from '../../../stores/batchRecord'
import { useBatchDefinitionStore } from '../../../stores/batchDefinition'
import { useIngredientStore } from '../../../stores/ingredient'
import { useAuthStore } from '../../../stores/auth'

// Mock router
const mockPush = vi.fn()
const mockBack = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush,
    back: mockBack
  }),
  useRoute: () => ({
    path: '/batch/records/create',
    params: {},
    query: {}
  })
}))

describe('BatchRecordForm', () => {
  let wrapper
  let batchRecordStore
  let batchDefinitionStore
  let ingredientStore
  let authStore

  beforeEach(() => {
    setActivePinia(createPinia())
    batchRecordStore = useBatchRecordStore()
    batchDefinitionStore = useBatchDefinitionStore()
    ingredientStore = useIngredientStore()
    authStore = useAuthStore()
    
    mockPush.mockClear()
    mockBack.mockClear()

    // Setup mock data
    authStore.user = { id: 'user123', name: 'Test User' }
    
    batchDefinitionStore.definitions = [
      {
        id: 'def1',
        name: 'Cà Phê Concentrate',
        unit: 'ml',
        shelf_life_hours: 24,
        conversion_rates: [
          {
            source_ingredient_id: 'ing1',
            source_ingredient_name: 'Hạt Cà Phê',
            source_quantity: 100,
            source_unit: 'g',
            batch_quantity: 500,
            wastage_rate: 0.1
          }
        ]
      }
    ]

    ingredientStore.items = [
      {
        id: 'ing1',
        name: 'Hạt Cà Phê',
        cost_per_unit: 0.5,
        unit: 'g'
      }
    ]
  })

  it('renders form fields', () => {
    wrapper = mount(BatchRecordForm)
    
    expect(wrapper.find('select').exists()).toBe(true)
    expect(wrapper.find('input[type="number"]').exists()).toBe(true)
  })

  it('displays batch definitions in dropdown', () => {
    wrapper = mount(BatchRecordForm)
    
    const options = wrapper.findAll('option')
    expect(options.length).toBeGreaterThan(1)
    expect(wrapper.text()).toContain('Cà Phê Concentrate')
  })

  it('calculates required ingredients when quantity is entered', async () => {
    wrapper = mount(BatchRecordForm)
    
    // Select batch definition
    const select = wrapper.find('select')
    await select.setValue('def1')
    
    // Enter quantity
    const input = wrapper.find('input[type="number"]')
    await input.setValue(500)
    await flushPromises()
    
    // Should show required ingredients
    expect(wrapper.text()).toContain('Nguyên Liệu Cần Thiết')
    expect(wrapper.text()).toContain('Hạt Cà Phê')
    expect(wrapper.text()).toContain('110') // 100 * 1.1 (with 10% wastage)
  })

  it('calculates expected cost', async () => {
    wrapper = mount(BatchRecordForm)
    
    const select = wrapper.find('select')
    await select.setValue('def1')
    
    const input = wrapper.find('input[type="number"]')
    await input.setValue(500)
    await flushPromises()
    
    // Should show cost: 110g * 0.5 = 55
    expect(wrapper.text()).toContain('Chi Phí Dự Kiến')
    expect(wrapper.text()).toContain('55')
  })

  it('shows confirmation dialog when submitting', async () => {
    wrapper = mount(BatchRecordForm)
    
    const select = wrapper.find('select')
    await select.setValue('def1')
    
    const input = wrapper.find('input[type="number"]')
    await input.setValue(500)
    await flushPromises()
    
    const submitButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Ghi Nhận')
    )
    await submitButton.trigger('click')
    await flushPromises()
    
    expect(wrapper.text()).toContain('Xác Nhận Ghi Nhận Batch')
  })

  it('submits form with correct data', async () => {
    const createRecordSpy = vi.spyOn(batchRecordStore, 'createRecord')
    createRecordSpy.mockResolvedValue({
      id: 'record1',
      expires_at: new Date().toISOString()
    })
    
    wrapper = mount(BatchRecordForm)
    
    const select = wrapper.find('select')
    await select.setValue('def1')
    
    const input = wrapper.find('input[type="number"]')
    await input.setValue(500)
    await flushPromises()
    
    // Click submit
    const submitButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Ghi Nhận')
    )
    await submitButton.trigger('click')
    await flushPromises()
    
    // Confirm in dialog
    const confirmButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Xác nhận')
    )
    await confirmButton.trigger('click')
    await flushPromises()
    
    expect(createRecordSpy).toHaveBeenCalledWith({
      batch_definition_id: 'def1',
      quantity_produced: 500,
      prepared_by: 'user123'
    })
  })

  it('displays error message on submission failure', async () => {
    const createRecordSpy = vi.spyOn(batchRecordStore, 'createRecord')
    createRecordSpy.mockRejectedValue(new Error('Insufficient ingredients'))
    batchRecordStore.error = 'Không đủ nguyên liệu'
    
    wrapper = mount(BatchRecordForm)
    
    const select = wrapper.find('select')
    await select.setValue('def1')
    
    const input = wrapper.find('input[type="number"]')
    await input.setValue(500)
    await flushPromises()
    
    const submitButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Ghi Nhận')
    )
    await submitButton.trigger('click')
    await flushPromises()
    
    const confirmButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Xác nhận')
    )
    await confirmButton.trigger('click')
    await flushPromises()
    
    expect(wrapper.text()).toContain('Không đủ nguyên liệu')
  })

  it('navigates back when cancel is clicked', async () => {
    wrapper = mount(BatchRecordForm)
    
    const cancelButton = wrapper.findAll('button').find(btn => 
      btn.text().includes('Hủy')
    )
    await cancelButton.trigger('click')
    
    expect(mockBack).toHaveBeenCalled()
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import BatchDefinitionList from '../BatchDefinitionList.vue'
import { useBatchDefinitionStore } from '../../../stores/batchDefinition'

// Mock router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({
    path: '/batch/definitions',
    params: {},
    query: {}
  })
}))

describe('BatchDefinitionList', () => {
  let wrapper
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useBatchDefinitionStore()
    mockPush.mockClear()
  })

  const mockDefinitions = [
    {
      id: '1',
      name: 'Cà Phê Concentrate',
      unit: 'ml',
      shelf_life_hours: 24,
      low_stock_threshold: 200,
      conversion_rates: [
        {
          source_ingredient_name: 'Hạt Cà Phê',
          source_quantity: 100,
          source_unit: 'g',
          batch_quantity: 500,
          wastage_rate: 0.1
        }
      ]
    },
    {
      id: '2',
      name: 'Trà Đen Concentrate',
      unit: 'ml',
      shelf_life_hours: 12,
      low_stock_threshold: 150,
      conversion_rates: []
    }
  ]

  it('renders loading state', () => {
    store.loading = true
    wrapper = mount(BatchDefinitionList)
    
    expect(wrapper.text()).toContain('Đang tải')
  })

  it('renders empty state when no definitions', () => {
    store.definitions = []
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    expect(wrapper.text()).toContain('Chưa có định nghĩa batch nào')
  })

  it('renders list of definitions', () => {
    store.definitions = mockDefinitions
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    expect(wrapper.text()).toContain('Cà Phê Concentrate')
    expect(wrapper.text()).toContain('Trà Đen Concentrate')
    expect(wrapper.text()).toContain('ml')
  })

  it('navigates to create form when clicking create button', async () => {
    store.definitions = []
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    const createButton = wrapper.find('button')
    await createButton.trigger('click')
    
    expect(mockPush).toHaveBeenCalledWith('/batch/definitions/create')
  })

  it('filters definitions by search term', async () => {
    store.definitions = mockDefinitions
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    const searchInput = wrapper.find('input[type="text"]')
    await searchInput.setValue('Cà Phê')
    
    expect(wrapper.text()).toContain('Cà Phê Concentrate')
    expect(wrapper.text()).not.toContain('Trà Đen Concentrate')
  })

  it('shows shelf life in hours', () => {
    store.definitions = mockDefinitions
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    expect(wrapper.text()).toContain('24')
    expect(wrapper.text()).toContain('12')
  })

  it('displays conversion rates count', () => {
    store.definitions = mockDefinitions
    store.loading = false
    wrapper = mount(BatchDefinitionList)
    
    expect(wrapper.text()).toContain('1') // First definition has 1 conversion rate
  })
})

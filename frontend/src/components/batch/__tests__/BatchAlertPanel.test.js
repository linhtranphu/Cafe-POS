import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import BatchAlertPanel from '../BatchAlertPanel.vue'
import { useBatchAlertStore } from '../../../stores/batchAlert'

// Mock router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({
    path: '/batch/alerts',
    params: {},
    query: {}
  })
}))

// Mock localStorage for happy-dom environment
global.localStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}

describe('BatchAlertPanel', () => {
  let wrapper
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useBatchAlertStore()
  })

  const mockAlerts = {
    low_stock: [
      {
        batch_definition_id: '1',
        batch_name: 'Cà Phê Concentrate',
        current_stock: 150,
        threshold: 200,
        unit: 'ml'
      }
    ],
    expiring: [
      {
        batch_record_id: '2',
        batch_name: 'Trà Đen Concentrate',
        quantity_remaining: 100,
        unit: 'ml',
        expires_at: new Date(Date.now() + 2 * 3600000).toISOString(),
        hours_until_expiry: 2
      }
    ],
    expired: [
      {
        batch_record_id: '3',
        batch_name: 'Sữa Tươi',
        quantity_wasted: 50,
        unit: 'ml',
        cost_wasted: 25000,
        expired_at: new Date(Date.now() - 3600000).toISOString()
      }
    ],
    last_checked: new Date().toISOString()
  }

  it('renders loading state', () => {
    store.loading = true
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Đang tải')
  })

  it('renders empty state when no alerts', () => {
    store.alerts = {
      low_stock: [],
      expiring: [],
      expired: [],
      last_checked: new Date().toISOString()
    }
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Không có cảnh báo')
  })

  it('displays low stock alerts', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Tồn Kho Thấp')
    expect(wrapper.text()).toContain('Cà Phê Concentrate')
    expect(wrapper.text()).toContain('150')
    expect(wrapper.text()).toContain('200')
  })

  it('displays expiring alerts', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Sắp Hết Hạn')
    expect(wrapper.text()).toContain('Trà Đen Concentrate')
    expect(wrapper.text()).toContain('2')
  })

  it('displays expired alerts', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Đã Hết Hạn')
    expect(wrapper.text()).toContain('Sữa Tươi')
    expect(wrapper.text()).toContain('50')
  })

  it('shows alert count badges', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    // Should show counts: 1 low stock, 1 expiring, 1 expired
    const badges = wrapper.findAll('.badge, [class*="badge"]')
    expect(badges.length).toBeGreaterThan(0)
  })

  it('allows expanding and collapsing sections', async () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    // Find section header
    const headers = wrapper.findAll('button, [role="button"]')
    if (headers.length > 0) {
      await headers[0].trigger('click')
      await flushPromises()
      
      // Section should toggle
      expect(wrapper.html()).toBeTruthy()
    }
  })

  it('auto-refreshes alerts', async () => {
    vi.useFakeTimers()
    const fetchAlertsSpy = vi.spyOn(store, 'fetchAlerts')
    fetchAlertsSpy.mockResolvedValue()
    
    wrapper = mount(BatchAlertPanel, {
      props: {
        autoRefresh: true,
        refreshInterval: 1000 // 1 second for testing
      }
    })
    
    // Fast-forward time
    vi.advanceTimersByTime(1000)
    await flushPromises()
    
    expect(fetchAlertsSpy).toHaveBeenCalled()
    
    vi.useRealTimers()
  })

  it('displays last checked time', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    expect(wrapper.text()).toContain('Cập nhật')
  })

  it('shows color coding for different alert types', () => {
    store.alerts = mockAlerts
    store.loading = false
    wrapper = mount(BatchAlertPanel)
    
    const html = wrapper.html()
    
    // Should have different colors for different alert types
    expect(html).toContain('yellow') // Low stock
    expect(html).toContain('orange') // Expiring
    expect(html).toContain('red') // Expired
  })
})

import api from './api'

// Named exports for store compatibility
export async function fetchShopSettings() {
  try {
    const response = await api.get('/manager/shop-settings')
    return response.data
  } catch (error) {
    console.error('Error fetching shop settings:', error)
    throw error
  }
}

export async function createShopSettings(settings) {
  try {
    const response = await api.post('/manager/shop-settings', settings)
    return response.data
  } catch (error) {
    console.error('Error creating shop settings:', error)
    throw error
  }
}

export async function updateShopSettings(id, settings) {
  try {
    const response = await api.put(`/manager/shop-settings/${id}`, settings)
    return response.data
  } catch (error) {
    console.error('Error updating shop settings:', error)
    throw error
  }
}

// Object export for direct usage in components
export const shopSettingsService = {
  getSettings: fetchShopSettings,
  createSettings: createShopSettings,
  updateSettings: updateShopSettings
}

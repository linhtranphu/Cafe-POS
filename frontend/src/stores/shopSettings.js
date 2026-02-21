import { defineStore } from 'pinia'
import { fetchShopSettings, updateShopSettings } from '../services/shopSettings'

export const useShopSettingsStore = defineStore('shopSettings', {
  state: () => ({
    settings: null,
    loading: false,
    error: null
  }),

  getters: {
    hasSettings: (state) => state.settings !== null,
    
    paperWidthOptions: () => [
      { value: 58, label: '58mm' },
      { value: 80, label: '80mm' }
    ],
    
    labelSizeOptions: () => [
      { value: '40x30', label: '40x30mm' },
      { value: '50x30', label: '50x30mm' },
      { value: '60x40', label: '60x40mm' }
    ]
  },

  actions: {
    async fetchSettings() {
      this.loading = true
      this.error = null
      try {
        this.settings = await fetchShopSettings()
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch settings'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateSettings(settingsData) {
      if (!this.settings?.id) {
        throw new Error('Settings not loaded')
      }

      this.loading = true
      this.error = null
      try {
        this.settings = await updateShopSettings(this.settings.id, settingsData)
        return this.settings
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update settings'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})

import { defineStore } from 'pinia'
import { fetchShopSettings, updateShopSettings } from '../services/shopSettings'

export const useShopSettingsStore = defineStore('shopSettings', {
  state: () => ({
    settings: null,
    loading: false,
    error: null
  }),

  getters: {
    hasSettings: (state) => state.settings !== null
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

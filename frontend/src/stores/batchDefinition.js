import { defineStore } from 'pinia'
import { batchDefinitionService } from '../services/batchDefinition'

export const useBatchDefinitionStore = defineStore('batchDefinition', {
  state: () => ({
    definitions: [],
    loading: false,
    error: null
  }),

  getters: {
    getDefinitionById: (state) => (id) => {
      return state.definitions.find(d => d.id === id) || null
    },

    activeDefinitions: (state) => {
      return state.definitions.filter(d => d.id)
    }
  },

  actions: {
    async fetchDefinitions() {
      this.loading = true
      this.error = null
      try {
        this.definitions = await batchDefinitionService.getDefinitions() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải định nghĩa batch'
        this.definitions = []
      } finally {
        this.loading = false
      }
    },

    async fetchDefinitionById(id) {
      this.loading = true
      this.error = null
      try {
        const definition = await batchDefinitionService.getDefinitionById(id)
        const index = this.definitions.findIndex(d => d.id === id)
        if (index !== -1) {
          this.definitions[index] = definition
        } else {
          this.definitions.push(definition)
        }
        return definition
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải định nghĩa batch'
        throw error
      } finally {
        this.loading = false
      }
    },

    async createDefinition(definition) {
      this.error = null
      try {
        const newDefinition = await batchDefinitionService.createDefinition(definition)
        this.definitions.push(newDefinition)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo định nghĩa batch'
        return false
      }
    },

    async updateDefinition(id, definition) {
      this.error = null
      try {
        const updatedDefinition = await batchDefinitionService.updateDefinition(id, definition)
        const index = this.definitions.findIndex(d => d.id === id)
        if (index !== -1) {
          this.definitions[index] = updatedDefinition
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật định nghĩa batch'
        return false
      }
    },

    async deleteDefinition(id) {
      this.error = null
      try {
        await batchDefinitionService.deleteDefinition(id)
        this.definitions = this.definitions.filter(d => d.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa định nghĩa batch'
        return false
      }
    }
  }
})

import { defineStore } from 'pinia'
import { printTemplateService } from '../services/printTemplate'

export const usePrintTemplateStore = defineStore('printTemplate', {
  state: () => ({
    templates: [],
    loading: false,
    error: null,
    previewLoading: false,
    previewResult: null
  }),

  getters: {
    defaultTemplates: (state) => {
      return state.templates.filter(t => t.is_default)
    },

    templatesByType: (state) => (type) => {
      return state.templates.filter(t => t.type === type)
    },

    defaultBillTemplate: (state) => {
      return state.templates.find(t => t.type === 'BILL' && t.is_default) || null
    },

    defaultLabelTemplate: (state) => {
      return state.templates.find(t => t.type === 'LABEL' && t.is_default) || null
    },

    billTemplates: (state) => {
      return state.templates.filter(t => t.type === 'BILL')
    },

    labelTemplates: (state) => {
      return state.templates.filter(t => t.type === 'LABEL')
    },

    getTemplateById: (state) => (id) => {
      return state.templates.find(t => t.id === id) || null
    }
  },

  actions: {
    async fetchTemplates(params = {}) {
      this.loading = true
      this.error = null
      try {
        this.templates = await printTemplateService.fetchTemplates(params)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh sách templates'
        this.templates = []
      } finally {
        this.loading = false
      }
    },

    async fetchTemplateById(id) {
      this.loading = true
      this.error = null
      try {
        const template = await printTemplateService.fetchTemplateById(id)
        const index = this.templates.findIndex(t => t.id === id)
        if (index !== -1) {
          this.templates[index] = template
        } else {
          this.templates.push(template)
        }
        return template
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải template'
        throw error
      } finally {
        this.loading = false
      }
    },

    async saveTemplate(template) {
      this.error = null
      try {
        if (template.id) {
          // Update existing template
          const updated = await printTemplateService.updateTemplate(template.id, template)
          const index = this.templates.findIndex(t => t.id === template.id)
          if (index !== -1) {
            this.templates[index] = updated
          }
          return updated
        } else {
          // Create new template
          const created = await printTemplateService.createTemplate(template)
          this.templates.push(created)
          return created
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi lưu template'
        throw error
      }
    },

    async deleteTemplate(id) {
      this.error = null
      try {
        await printTemplateService.deleteTemplate(id)
        this.templates = this.templates.filter(t => t.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa template'
        return false
      }
    },

    async previewTemplate(id, sampleData = {}) {
      this.previewLoading = true
      this.previewResult = null
      this.error = null
      try {
        const result = await printTemplateService.previewTemplate(id, sampleData)
        this.previewResult = result
        return result
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi preview template'
        this.previewResult = null
        throw error
      } finally {
        this.previewLoading = false
      }
    },

    clearPreview() {
      this.previewResult = null
    },

    clearError() {
      this.error = null
    }
  }
})

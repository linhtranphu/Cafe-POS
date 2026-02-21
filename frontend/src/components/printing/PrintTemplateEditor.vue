<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">📝 Templates</h2>
        <button
          @click="createNewTemplate"
          class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform"
        >
          ➕ Tạo mới
        </button>
      </div>

      <!-- Type Filter -->
      <div class="flex gap-2">
        <button
          v-for="filter in typeFilters"
          :key="filter.value"
          @click="selectedType = filter.value"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors',
            selectedType === filter.value
              ? 'bg-blue-500 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          {{ filter.label }}
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải...</p>
      </div>

      <!-- Template List (when no template selected) -->
      <div v-else-if="!editingTemplate" class="px-4 py-4">
        <!-- Empty State -->
        <div v-if="filteredTemplates.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📝</div>
          <p class="text-gray-500 mb-2">Chưa có template nào</p>
          <button
            @click="createNewTemplate"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 text-sm"
          >
            Tạo template đầu tiên
          </button>
        </div>

        <!-- Templates List -->
        <div v-else class="space-y-3">
          <div
            v-for="template in filteredTemplates"
            :key="template.id"
            class="bg-white rounded-2xl p-4 shadow-sm cursor-pointer hover:shadow-md transition-shadow"
            @click="editTemplate(template)"
          >
            <div class="flex justify-between items-start mb-2">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-2xl">{{ getTemplateIcon(template.type) }}</span>
                  <h3 class="font-bold text-lg">{{ template.name }}</h3>
                  <span
                    v-if="template.is_default"
                    class="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs font-bold rounded-full"
                  >
                    ⭐ Mặc định
                  </span>
                </div>
                <p class="text-sm text-gray-600">{{ getTemplateTypeLabel(template.type) }}</p>
              </div>
            </div>
            <div class="text-xs text-gray-500">
              Cập nhật: {{ formatDateTime(template.updated_at) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Template Editor (when template selected) -->
      <div v-else class="h-full flex flex-col">
        <!-- Editor Header -->
        <div class="bg-white border-b px-4 py-3 flex items-center justify-between">
          <button
            @click="cancelEdit"
            class="text-gray-600 hover:text-gray-800 flex items-center gap-2"
          >
            <span class="text-xl">←</span>
            <span class="font-medium">Quay lại</span>
          </button>
          <div class="flex gap-2">
            <button
              @click="handlePreview"
              :disabled="previewLoading"
              class="bg-purple-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-purple-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <span v-if="previewLoading" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <span v-else>👁️</span>
              <span>Preview</span>
            </button>
            <button
              @click="handleSave"
              :disabled="saving"
              class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <span v-if="saving" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <span>💾 Lưu</span>
            </button>
          </div>
        </div>

        <!-- Editor Content -->
        <div class="flex-1 overflow-hidden flex flex-col lg:flex-row">
          <!-- Left: Editor -->
          <div class="flex-1 flex flex-col border-r">
            <div class="bg-gray-100 px-4 py-2 border-b">
              <h3 class="font-bold text-sm text-gray-700">📝 Template Content</h3>
            </div>
            <div class="flex-1 overflow-y-auto p-4 space-y-4">
              <!-- Template Info -->
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">Tên Template *</label>
                <input
                  v-model="editingTemplate.name"
                  type="text"
                  placeholder="VD: Bill Template 1"
                  class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">Loại Template *</label>
                <select
                  v-model="editingTemplate.type"
                  class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="BILL">🧾 Bill</option>
                  <option value="LABEL">🏷️ Label</option>
                </select>
              </div>

              <div>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input
                    v-model="editingTemplate.is_default"
                    type="checkbox"
                    class="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                  />
                  <span class="text-sm font-medium text-gray-700">⭐ Đặt làm template mặc định</span>
                </label>
              </div>

              <!-- Template Content Editor -->
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">Nội dung Template *</label>
                <textarea
                  v-model="editingTemplate.content"
                  rows="20"
                  placeholder="Nhập nội dung template (Go template syntax)..."
                  class="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
                ></textarea>
                <p class="text-xs text-gray-500 mt-2">
                  Sử dụng Go template syntax. VD: &#123;&#123;.ShopName&#125;&#125;, &#123;&#123;.Order.OrderNumber&#125;&#125;
                </p>
              </div>

              <!-- Template Variables Help -->
              <div class="bg-blue-50 border-2 border-blue-200 rounded-lg p-3">
                <h4 class="font-bold text-sm text-blue-900 mb-2">📚 Biến có sẵn:</h4>
                <div class="text-xs text-blue-800 space-y-1 font-mono">
                  <div v-if="editingTemplate.type === 'BILL'">
                    <p>• &#123;&#123;.ShopName&#125;&#125; - Tên quán</p>
                    <p>• &#123;&#123;.ShopAddress&#125;&#125; - Địa chỉ</p>
                    <p>• &#123;&#123;.ShopPhone&#125;&#125; - Số điện thoại</p>
                    <p>• &#123;&#123;.Order.OrderNumber&#125;&#125; - Số order</p>
                    <p>• &#123;&#123;.Order.Total&#125;&#125; - Tổng tiền</p>
                    <p>• &#123;&#123;range .Order.Items&#125;&#125;...&#123;&#123;end&#125;&#125; - Lặp items</p>
                  </div>
                  <div v-else>
                    <p>• &#123;&#123;.Order.OrderNumber&#125;&#125; - Số order</p>
                    <p>• &#123;&#123;.ItemIndex&#125;&#125; - Số thứ tự item</p>
                    <p>• &#123;&#123;.TotalItems&#125;&#125; - Tổng số items</p>
                    <p>• &#123;&#123;.Item.Name&#125;&#125; - Tên món</p>
                    <p>• &#123;&#123;.Item.VariantName&#125;&#125; - Tên variant</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right: Preview -->
          <div class="flex-1 flex flex-col bg-gray-100">
            <div class="bg-gray-200 px-4 py-2 border-b">
              <h3 class="font-bold text-sm text-gray-700">👁️ Preview</h3>
            </div>
            <div class="flex-1 overflow-y-auto p-4">
              <div v-if="previewResult" class="bg-white rounded-lg p-4 shadow-sm">
                <pre class="font-mono text-xs whitespace-pre-wrap">{{ previewResult.content }}</pre>
              </div>
              <div v-else class="text-center py-16 text-gray-500">
                <div class="text-4xl mb-2">👁️</div>
                <p class="text-sm">Click "Preview" để xem kết quả</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="bg-red-50 border-t-2 border-red-200 px-4 py-3">
          <div class="flex items-center gap-2 text-red-800 text-sm">
            <span>⚠️</span>
            <span class="font-bold">{{ error }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePrintTemplateStore } from '../../stores/printTemplate'

const templateStore = usePrintTemplateStore()

const selectedType = ref('all')
const editingTemplate = ref(null)
const saving = ref(false)

const typeFilters = [
  { label: 'Tất cả', value: 'all' },
  { label: '🧾 Bill', value: 'BILL' },
  { label: '🏷️ Tem', value: 'LABEL' }
]

const loading = computed(() => templateStore.loading)
const error = computed(() => templateStore.error)
const previewLoading = computed(() => templateStore.previewLoading)
const previewResult = computed(() => templateStore.previewResult)

const filteredTemplates = computed(() => {
  if (selectedType.value === 'all') {
    return templateStore.templates
  }
  return templateStore.templatesByType(selectedType.value)
})

const getTemplateIcon = (type) => {
  return type === 'BILL' ? '🧾' : '🏷️'
}

const getTemplateTypeLabel = (type) => {
  return type === 'BILL' ? 'Bill Template' : 'Label Template'
}

const formatDateTime = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const createNewTemplate = () => {
  editingTemplate.value = {
    name: '',
    type: 'BILL',
    content: '',
    is_default: false
  }
  templateStore.clearPreview()
  templateStore.clearError()
}

const editTemplate = (template) => {
  editingTemplate.value = { ...template }
  templateStore.clearPreview()
  templateStore.clearError()
}

const cancelEdit = () => {
  if (editingTemplate.value && !editingTemplate.value.id) {
    // New template, just close
    editingTemplate.value = null
  } else if (confirm('Bạn có chắc muốn hủy? Các thay đổi chưa lưu sẽ bị mất.')) {
    editingTemplate.value = null
  }
  templateStore.clearPreview()
}

const handlePreview = async () => {
  if (!editingTemplate.value) return

  try {
    // For preview, we need to save first if it's a new template
    if (!editingTemplate.value.id) {
      alert('Vui lòng lưu template trước khi preview')
      return
    }

    await templateStore.previewTemplate(editingTemplate.value.id)
  } catch (err) {
    console.error('Preview error:', err)
  }
}

const handleSave = async () => {
  if (!editingTemplate.value || saving.value) return

  if (!editingTemplate.value.name || !editingTemplate.value.type || !editingTemplate.value.content) {
    alert('Vui lòng điền đầy đủ thông tin')
    return
  }

  saving.value = true
  templateStore.clearError()

  try {
    await templateStore.saveTemplate(editingTemplate.value)
    alert('Đã lưu template')
    editingTemplate.value = null
    await templateStore.fetchTemplates()
  } catch (err) {
    console.error('Save error:', err)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  templateStore.fetchTemplates()
})
</script>

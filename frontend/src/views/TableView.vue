<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">🪑 Quản lý Bàn</h2>
        <button v-if="isManager" @click="showCreateForm = true" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium">
          + Thêm bàn
        </button>
      </div>

      <!-- Status Filter -->
      <div class="flex gap-2 mb-4">
        <button @click="filterStatus = 'ALL'" :class="filterStatus === 'ALL' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
          class="px-4 py-2 rounded-lg font-medium">
          Tất cả ({{ tables.length }})
        </button>
        <button @click="filterStatus = 'EMPTY'" :class="filterStatus === 'EMPTY' ? 'bg-green-500 text-white' : 'bg-white text-gray-700'"
          class="px-4 py-2 rounded-lg font-medium">
          Trống ({{ emptyTables.length }})
        </button>
        <button @click="filterStatus = 'OCCUPIED'" :class="filterStatus === 'OCCUPIED' ? 'bg-red-500 text-white' : 'bg-white text-gray-700'"
          class="px-4 py-2 rounded-lg font-medium">
          Đang dùng ({{ occupiedTables.length }})
        </button>
      </div>

      <!-- Tables Grid -->
      <div v-if="loading" class="text-center py-10">Đang tải...</div>
      <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div v-for="table in filteredTables" :key="table.id" 
          :class="getTableColor(table.status)"
          class="rounded-xl p-6 shadow-sm cursor-pointer hover:shadow-md transition-shadow">
          <div class="text-center">
            <div class="text-4xl mb-2">🪑</div>
            <h3 class="font-bold text-lg">{{ table.name }}</h3>
            <p class="text-sm text-gray-600">{{ table.capacity }} chỗ</p>
            <p v-if="table.area" class="text-xs text-gray-500 mt-1">{{ table.area }}</p>
            <span :class="getStatusBadge(table.status)" class="inline-block mt-2 px-3 py-1 rounded-full text-xs font-medium">
              {{ getStatusText(table.status) }}
            </span>
          </div>
          
          <div v-if="isManager" class="mt-4 grid grid-cols-2 gap-2">
            <button @click="editTable(table)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded text-sm">
              Sửa
            </button>
            <button @click="deleteTable(table.id)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm">
              Xóa
            </button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <div v-if="showCreateForm || editingTable" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">{{ editingTable ? 'Sửa bàn' : 'Thêm bàn mới' }}</h3>
          <form @submit.prevent="saveTable" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tên bàn *</label>
              <input v-model="form.name" type="text" required placeholder="Ví dụ: Bàn 1" class="w-full p-3 border rounded-lg">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Số chỗ ngồi *</label>
              <input v-model.number="form.capacity" type="number" min="1" required class="w-full p-3 border rounded-lg">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Khu vực</label>
              <input v-model="form.area" type="text" placeholder="Ví dụ: Tầng 1, Ngoài trời" class="w-full p-3 border rounded-lg">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="cancelEdit" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                {{ editingTable ? 'Cập nhật' : 'Thêm' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTableStore } from '../stores/table'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const tableStore = useTableStore()
const authStore = useAuthStore()

const showCreateForm = ref(false)
const editingTable = ref(null)
const filterStatus = ref('ALL')

const form = ref({
  name: '',
  capacity: 2,
  area: ''
})

const loading = computed(() => tableStore.loading)
const tables = computed(() => tableStore.tables)
const emptyTables = computed(() => tableStore.emptyTables)
const occupiedTables = computed(() => tableStore.occupiedTables)
const isManager = computed(() => authStore.user?.role === 'manager')

const filteredTables = computed(() => {
  if (filterStatus.value === 'ALL') return tables.value
  return tables.value.filter(t => t.status === filterStatus.value)
})

onMounted(() => {
  tableStore.fetchTables()
})

const editTable = (table) => {
  editingTable.value = table
  form.value = {
    name: table.name,
    capacity: table.capacity,
    area: table.area || ''
  }
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingTable.value = null
  form.value = { name: '', capacity: 2, area: '' }
}

const saveTable = async () => {
  try {
    if (editingTable.value) {
      await tableStore.updateTable(editingTable.value.id, form.value)
    } else {
      await tableStore.createTable(form.value)
    }
    cancelEdit()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const deleteTable = async (id) => {
  if (confirm('Bạn có chắc muốn xóa bàn này?')) {
    try {
      await tableStore.deleteTable(id)
    } catch (error) {
      alert('Lỗi: ' + (error.response?.data?.error || error.message))
    }
  }
}

const getTableColor = (status) => {
  return status === 'EMPTY' ? 'bg-green-50 border-2 border-green-200' : 'bg-red-50 border-2 border-red-200'
}

const getStatusBadge = (status) => {
  return status === 'EMPTY' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
}

const getStatusText = (status) => {
  return status === 'EMPTY' ? 'Trống' : 'Đang dùng'
}
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>

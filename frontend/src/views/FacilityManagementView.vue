<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">🏢 Cơ sở vật chất</h1>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm thiết bị..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Stats Cards - Single Row -->
      <div class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="text-xs opacity-90 mb-2">Tổng quan</div>
        <div class="grid grid-cols-4 gap-1.5">
          <div class="text-center">
            <div class="text-lg font-bold">{{ facilities.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Tổng</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ operationalCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Hoạt động</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ maintenanceCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Bảo trì</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ brokenCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Hỏng</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal"
            class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Tạo thiết bị</div>
          </button>
          <button @click="showCategoryModal = true"
            class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📁</div>
            <div class="text-sm font-bold">Quản lý danh mục</div>
          </button>
          <button @click="showMaintenanceSchedule = true"
            class="bg-gradient-to-br from-yellow-500 to-orange-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📅</div>
            <div class="text-sm font-bold">Lịch bảo trì</div>
          </button>
          <button @click="showIssueReports = true"
            class="bg-gradient-to-br from-red-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">⚠️</div>
            <div class="text-sm font-bold">Sự cố</div>
          </button>
        </div>
      </div>

      <!-- Facilities List -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">📋 Danh sách thiết bị</h2>
        </div>
        
        <div v-if="filteredFacilities.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-500">Không có thiết bị nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="facility in filteredFacilities" :key="facility.id"
            @click="viewDetails(facility)"
            class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform">
            
            <!-- Facility Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ facility.name }}</h3>
                <p class="text-sm text-gray-600">{{ facility.type }}</p>
                <p class="text-xs text-gray-400">📍 {{ facility.area }}</p>
              </div>
              <span :class="getStatusClass(facility.status)" class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getStatusText(facility.status) }}
              </span>
            </div>

            <!-- Facility Info -->
            <div class="mb-3 space-y-1 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600">Số lượng:</span>
                <span class="font-medium">{{ facility.quantity }}</span>
              </div>
              <div v-if="facility.cost" class="flex justify-between">
                <span class="text-gray-600">Giá trị:</span>
                <span class="font-medium text-green-600">{{ formatPrice(facility.cost) }}</span>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="flex gap-2 pt-3 border-t">
              <button @click.stop="openEditModal(facility)"
                class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-sm font-medium active:bg-blue-600">
                ✏️ Sửa
              </button>
              <button @click.stop="deleteFacility(facility)"
                class="flex-1 bg-red-500 text-white py-2 rounded-lg text-sm font-medium active:bg-red-600">
                🗑️ Xóa
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Maintenance Schedule Modal - Mobile Optimized -->
    <transition name="slide-up">
      <div v-if="showMaintenanceSchedule" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
            <h3 class="text-lg font-bold">📅 Lịch bảo trì</h3>
            <button @click="showMaintenanceSchedule = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
            <div v-if="maintenanceSchedule.length === 0" class="text-center py-16">
              <div class="text-6xl mb-4">📭</div>
              <p class="text-gray-500">Không có lịch bảo trì nào</p>
            </div>
            
            <div v-else class="space-y-3 pb-4">
              <div v-for="item in maintenanceSchedule" :key="item.id" 
                class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
                <div class="flex justify-between items-start gap-3 mb-3">
                  <div class="flex-1 min-w-0">
                    <h4 class="font-bold text-gray-900 truncate">{{ item.facility_name }}</h4>
                    <p class="text-sm text-gray-600 truncate">📍 {{ item.location }}</p>
                  </div>
                  <span :class="item.is_overdue ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'" 
                    class="px-3 py-1 text-xs font-medium rounded-full whitespace-nowrap flex-shrink-0">
                    {{ item.is_overdue ? 'Quá hạn' : 'Sắp tới' }}
                  </span>
                </div>
                <p class="text-sm text-gray-500">📅 {{ formatDate(item.scheduled_date) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Issue Reports Modal - Mobile Optimized -->
    <transition name="slide-up">
      <div v-if="showIssueReports" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
            <h3 class="text-lg font-bold">⚠️ Báo cáo sự cố</h3>
            <button @click="showIssueReports = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
            <div v-if="issueReports.length === 0" class="text-center py-16">
              <div class="text-6xl mb-4">✅</div>
              <p class="text-gray-500">Không có sự cố nào</p>
            </div>
            
            <div v-else class="space-y-3 pb-4">
              <div v-for="issue in issueReports" :key="issue.id" 
                class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
                <div class="flex justify-between items-start gap-3 mb-3">
                  <div class="flex-1 min-w-0">
                    <h4 class="font-bold text-gray-900 truncate">{{ issue.facility_name }}</h4>
                    <p class="text-sm text-gray-700 mt-2 line-clamp-2">{{ issue.description }}</p>
                  </div>
                  <span :class="getIssueStatusClassLocal(issue.status)" 
                    class="px-3 py-1 text-xs font-medium rounded-full whitespace-nowrap flex-shrink-0">
                    {{ getIssueStatusText(issue.status) }}
                  </span>
                </div>
                <p class="text-xs text-gray-500 mt-3">
                  👤 {{ issue.reported_by }} • {{ formatDate(issue.reported_at) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Category Management Modal -->
    <transition name="slide-up">
      <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
            <h3 class="text-lg font-bold">📁 Quản lý danh mục thiết bị</h3>
            <button @click="showCategoryModal = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
            <!-- Add New Category -->
            <div class="bg-gray-50 rounded-xl p-4 mb-4 flex-shrink-0">
              <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
              <div class="flex flex-col sm:flex-row gap-2">
                <input v-model="newCategoryName" type="text" placeholder="Tên danh mục..." 
                  class="flex-1 px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500" />
                <button @click="addCategory" class="bg-purple-500 text-white px-6 py-3 rounded-lg font-medium text-base active:bg-purple-600 whitespace-nowrap">
                  Thêm
                </button>
              </div>
            </div>

            <!-- Category List -->
            <div class="space-y-3 pb-4">
              <div v-for="cat in facilityCategories" :key="cat" 
                class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                <div class="flex items-center gap-3 flex-1 min-w-0">
                  <div class="w-12 h-12 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center text-2xl flex-shrink-0">
                    🏢
                  </div>
                  <div class="min-w-0">
                    <div class="font-medium text-gray-800 truncate">{{ cat }}</div>
                    <div class="text-xs text-gray-500">{{ getCategoryCount(cat) }} thiết bị</div>
                  </div>
                </div>
                <button @click="deleteCategory(cat)" class="text-red-500 hover:text-red-700 p-2 flex-shrink-0 ml-2">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Facility Form Modal - Slide from Right -->
    <transition name="slide-right">
      <div v-if="showFacilityForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <!-- Mobile Header - Fixed -->
          <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
            <div class="px-4 py-3">
              <div class="flex items-center justify-between">
                <button @click="showFacilityForm = false" class="text-2xl text-gray-600">←</button>
                <h1 class="text-xl font-bold text-gray-800">{{ editingFacility ? '✏️ Cập nhật thiết bị' : '➕ Thêm thiết bị mới' }}</h1>
                <div class="w-8"></div>
              </div>
            </div>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
            <!-- Tên thiết bị -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Tên thiết bị *</label>
              <input v-model="formData.name" type="text" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Loại & Số lượng - Responsive Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Loại *</label>
                <select v-model="formData.type" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option value="">Chọn loại</option>
                  <option v-for="cat in facilityCategories" :key="cat" :value="cat">{{ cat }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Số lượng *</label>
                <input v-model.number="formData.quantity" type="number" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
            </div>

            <!-- Khu vực & Trạng thái - Responsive Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Khu vực</label>
                <input v-model="formData.area" type="text" placeholder="Mặc định"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Trạng thái *</label>
                <select v-model="formData.status" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option v-for="option in FACILITY_STATUS_OPTIONS" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Ngày mua & Giá trị - Responsive Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Ngày mua</label>
                <input v-model="formData.purchase_date" type="date" 
                  class="w-full px-3 py-3 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Giá trị (VND)</label>
                <input v-model.number="formData.cost" type="number" placeholder="0"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
            </div>

            <!-- Nhà cung cấp -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Nhà cung cấp</label>
              <input v-model="formData.supplier" type="text" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Ghi chú -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Ghi chú</label>
              <textarea v-model="formData.notes" rows="3" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"></textarea>
            </div>

            <!-- Auto-Expense Indicator -->
            <div v-if="!editingFacility && formData.cost > 0" 
              class="bg-green-50 border border-green-200 rounded-lg p-4 mt-6">
              <div class="flex items-start gap-3">
                <span class="text-green-600 text-2xl flex-shrink-0">✅</span>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-green-800">Tự động ghi nhận chi phí</p>
                  <p class="text-xs text-green-600 mt-2 break-words">
                    Hệ thống sẽ tự động tạo chi phí: <span class="font-semibold">{{ formatPrice(formData.cost) }}</span>
                  </p>
                  <p class="text-xs text-green-600 mt-1">Danh mục: Cơ sở vật chất</p>
                </div>
              </div>
            </div>

            <!-- Spacer for bottom buttons -->
            <div class="h-24"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="showFacilityForm = false" 
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors">
              Hủy
            </button>
            <button @click="saveFacility" 
              class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 transition-colors">
              {{ editingFacility ? 'Cập nhật' : 'Thêm mới' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useFacilityStore } from '../stores/facility'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { 
  FACILITY_STATUS, 
  FACILITY_STATUS_OPTIONS,
  FACILITY_TYPE_OPTIONS,
  getFacilityStatusClass,
  getIssueStatusClass
} from '../constants/facility'
import { 
  formatDate, 
  formatPrice
} from '../utils/formatters'
const facilityStore = useFacilityStore()

const searchQuery = ref('')
const showCategoryModal = ref(false)
const showMaintenanceSchedule = ref(false)
const showIssueReports = ref(false)
const showFacilityForm = ref(false)
const editingFacility = ref(null)
const newCategoryName = ref('')

const formData = ref({
  name: '',
  type: '',
  area: 'Mặc định',
  quantity: 1,
  status: FACILITY_STATUS.IN_USE,
  purchase_date: '',
  cost: 0,
  supplier: '',
  notes: ''
})

const facilities = computed(() => facilityStore.items || [])
const maintenanceSchedule = ref([])
const issueReports = ref([])

// Facility categories from constants + custom categories from backend
const facilityCategories = computed(() => {
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  const backendTypes = facilityStore.types.map(t => t.name)
  return [...new Set([...defaultCategories, ...backendTypes])]
})

const filteredFacilities = computed(() => {
  if (!searchQuery.value) return facilities.value
  const query = searchQuery.value.toLowerCase()
  return facilities.value.filter(f => 
    f.name?.toLowerCase().includes(query) ||
    f.type?.toLowerCase().includes(query) ||
    f.area?.toLowerCase().includes(query)
  )
})

const operationalCount = computed(() => 
  facilities.value.filter(f => f.status === FACILITY_STATUS.IN_USE).length
)
const maintenanceCount = computed(() => 
  facilities.value.filter(f => f.status === FACILITY_STATUS.REPAIRING).length
)
const brokenCount = computed(() => 
  facilities.value.filter(f => f.status === FACILITY_STATUS.BROKEN).length
)

const getStatusClass = (status) => {
  return getFacilityStatusClass(status)
}

const getStatusText = (status) => {
  return status
}

const getIssueStatusClassLocal = (status) => {
  return getIssueStatusClass(status)
}

const getIssueStatusText = (status) => {
  const texts = {
    open: 'Chờ xử lý',
    in_progress: 'Đang xử lý',
    resolved: 'Đã giải quyết'
  }
  return texts[status] || status
}

const getCategoryCount = (categoryName) => {
  return facilities.value.filter(f => f.type === categoryName).length
}

const addCategory = async () => {
  if (!newCategoryName.value.trim()) return
  
  if (facilityCategories.value.includes(newCategoryName.value.trim())) {
    alert('Danh mục đã tồn tại!')
    return
  }
  
  const success = await facilityStore.createFacilityType({ name: newCategoryName.value.trim() })
  if (success) {
    newCategoryName.value = ''
    alert('Thêm danh mục thành công')
  } else {
    alert(facilityStore.error || 'Lỗi thêm danh mục')
  }
}

const deleteCategory = async (categoryName) => {
  // Check if category is in use
  const hasItems = facilities.value.some(f => f.type === categoryName)
  if (hasItems) {
    alert('Không thể xóa danh mục đã có thiết bị!')
    return
  }
  
  // Check if it's a default category
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  if (defaultCategories.includes(categoryName)) {
    alert('Không thể xóa danh mục mặc định!')
    return
  }
  
  if (confirm(`Bạn có chắc muốn xóa danh mục "${categoryName}"?`)) {
    const type = facilityStore.types.find(t => t.name === categoryName)
    if (type) {
      const success = await facilityStore.deleteFacilityType(type.id)
      if (success) {
        alert('Xóa danh mục thành công')
      } else {
        alert(facilityStore.error || 'Lỗi xóa danh mục')
      }
    }
  }
}

// Utility functions imported from formatters.js

const openCreateModal = () => {
  const today = new Date()
  const year = today.getFullYear()
  const month = String(today.getMonth() + 1).padStart(2, '0')
  const day = String(today.getDate()).padStart(2, '0')
  const todayDate = `${year}-${month}-${day}`
  
  formData.value = {
    name: '',
    type: '',
    area: 'Mặc định',
    quantity: 1,
    status: FACILITY_STATUS.IN_USE,
    purchase_date: todayDate,
    cost: 0,
    supplier: '',
    notes: ''
  }
  editingFacility.value = null
  showFacilityForm.value = true
}

const openEditModal = (facility) => {
  editingFacility.value = facility
  formData.value = {
    name: facility.name || '',
    type: facility.type || '',
    area: facility.area || 'Mặc định',
    quantity: facility.quantity || 1,
    status: facility.status || FACILITY_STATUS.IN_USE,
    purchase_date: facility.purchase_date || '',
    cost: facility.cost || 0,
    supplier: facility.supplier || '',
    notes: facility.notes || ''
  }
  showFacilityForm.value = true
}

const deleteFacility = async (facility) => {
  if (confirm(`Bạn có chắc muốn xóa thiết bị "${facility.name}"?`)) {
    try {
      await facilityStore.deleteFacility(facility.id)
      alert('Xóa thiết bị thành công')
    } catch (error) {
      console.error('Error deleting facility:', error)
      const errorMessage = error.response?.data?.error || 'Có lỗi xảy ra khi xóa thiết bị'
      alert(errorMessage)
    }
  }
}

const saveFacility = async () => {
  try {
    // Prepare data - convert date format
    const dataToSend = { ...formData.value }
    
    // Remove empty purchase_date or convert to ISO format
    if (!dataToSend.purchase_date) {
      delete dataToSend.purchase_date
    } else {
      // Convert YYYY-MM-DD to ISO format YYYY-MM-DDT00:00:00Z
      dataToSend.purchase_date = dataToSend.purchase_date + 'T00:00:00Z'
    }
    
    console.log('Sending facility data:', dataToSend)
    
    if (editingFacility.value) {
      await facilityStore.updateFacility(editingFacility.value.id, dataToSend)
      alert('Cập nhật thiết bị thành công')
    } else {
      await facilityStore.createFacility(dataToSend)
      alert('Thêm thiết bị thành công')
    }
    
    // Reset form and close modal
    formData.value = {
      name: '',
      type: '',
      area: 'Mặc định',
      quantity: 1,
      status: FACILITY_STATUS.IN_USE,
      purchase_date: '',
      cost: 0,
      supplier: '',
      notes: ''
    }
    showFacilityForm.value = false
    editingFacility.value = null
    
    // Refresh facilities list
    await facilityStore.fetchFacilities()
  } catch (error) {
    console.error('Error saving facility:', error)
    const errorMessage = error.response?.data?.error || 'Có lỗi xảy ra khi lưu thiết bị'
    alert(errorMessage)
  }
}

const viewDetails = (facility) => {
  // Show detail modal or navigate to detail page
  console.log('View details:', facility)
  openEditModal(facility)
}

// Refresh data function
const refreshData = async () => {
  await facilityStore.fetchFacilities()
  await facilityStore.fetchFacilityTypes()
  maintenanceSchedule.value = await facilityStore.fetchScheduledMaintenance()
  issueReports.value = await facilityStore.fetchIssueReports()
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
}

.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-right-enter-from {
  transform: translateX(100%);
}

.slide-right-leave-to {
  transform: translateX(100%);
}

.pb-safe {
  padding-bottom: max(1rem, env(safe-area-inset-bottom));
}
</style>

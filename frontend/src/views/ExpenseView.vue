<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Header -->
      <div class="flex flex-col lg:flex-row justify-between items-center mb-6">
        <h2 class="text-xl lg:text-2xl font-semibold text-gray-800 mb-4 lg:mb-0">
          💰 Quản lý Chi phí
        </h2>
        <div class="flex flex-wrap gap-2">
          <button @click="showReportView = true" class="bg-teal-500 hover:bg-teal-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📊 Báo cáo
          </button>
          <button @click="showPrepaidForm = true" class="bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            💳 Chi phí trả trước
          </button>
          <button @click="showRecurringForm = true" class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            🔁 Chi phí định kỳ
          </button>
          <button @click="showCategoryForm = true" class="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📁 Quản lý loại
          </button>
          <button @click="showCreateForm = true" class="btn-primary text-sm px-4 py-2">
            + Thêm chi phí
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-10 text-gray-600 text-lg">Đang tải...</div>
      <div v-if="error" class="text-center py-10 text-red-600 bg-red-50 border border-red-200 rounded-lg">{{ error }}</div>

      <!-- Recurring Expenses Alert -->
      <div v-if="recurringDueReminders.length > 0" class="bg-orange-50 border border-orange-200 rounded-xl p-4 mb-4">
        <h3 class="text-orange-800 font-semibold mb-2">🔔 Nhắc nhập chi phí định kỳ ({{ recurringDueReminders.length }})</h3>
        <div class="space-y-2">
          <div v-for="reminder in recurringDueReminders" :key="reminder.id" class="flex items-center justify-between bg-white p-3 rounded-lg">
            <div>
              <div class="font-medium text-gray-800">{{ reminder.name }}</div>
              <div class="text-sm text-gray-500">Dự kiến: {{ formatPrice(reminder.estimated_amount) }} - Tháng {{ formatMonth(reminder.due_month) }}</div>
            </div>
            <button @click="recordRecurringExpense(reminder)" class="bg-orange-500 hover:bg-orange-600 text-white px-3 py-2 rounded-lg text-sm">
              Ghi nhận
            </button>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-xl p-4 mb-4 shadow-sm">
        <div class="grid grid-cols-1 gap-3">
          <select v-model="filterMonth" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500" @change="fetchExpenses">
            <option value="">Tất cả tháng</option>
            <option v-for="month in months" :key="month.value" :value="month.value">{{ month.label }}</option>
          </select>
          <select v-model="filterCategory" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500" @change="fetchExpenses">
            <option value="">Tất cả nhóm</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
      </div>

      <!-- Summary Card -->
      <div class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl p-6 mb-4 text-white shadow-lg">
        <div class="text-sm opacity-90 mb-1">Tổng chi phí {{ filterMonth ? 'tháng này' : 'tất cả' }}</div>
        <div class="text-3xl font-bold">{{ formatPrice(totalExpenses) }}</div>
      </div>

      <!-- Expense List -->
      <div class="grid grid-cols-1 gap-4">
        <div v-for="expense in filteredExpenses" :key="expense.id" class="bg-white rounded-xl p-4 shadow-md">
          <!-- Expense Header -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" :class="getCategoryColor(expense.category)">
                {{ getCategoryIcon(expense.category) }}
              </div>
              <div>
                <h4 class="font-bold text-gray-800">{{ expense.name }}</h4>
                <p class="text-sm text-gray-500">{{ expense.category }}</p>
              </div>
            </div>
            <div class="text-right">
              <div class="text-lg font-bold text-gray-800">{{ formatPrice(expense.amount) }}</div>
              <div class="text-xs text-gray-500">{{ formatMonth(expense.month) }}</div>
            </div>
          </div>

          <!-- Expense Details -->
          <div class="grid grid-cols-2 gap-3 text-sm">
            <div class="bg-gray-50 rounded-lg p-3">
              <div class="text-gray-500 text-xs mb-1">Loại</div>
              <div class="font-medium text-gray-800">{{ getExpenseType(expense.type) }}</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <div class="text-gray-500 text-xs mb-1">Ngày</div>
              <div class="font-medium text-gray-800">{{ formatDate(expense.date) }}</div>
            </div>
            <div v-if="expense.is_allocated" class="col-span-2 bg-indigo-50 rounded-lg p-3">
              <div class="text-indigo-600 text-xs mb-1">💳 Chi phí phân bổ</div>
              <div class="font-medium text-indigo-800">{{ expense.allocated_month_index }}/{{ expense.total_months }} tháng - Gốc: {{ formatPrice(expense.original_amount) }}</div>
            </div>
            <div v-if="expense.attachments && expense.attachments.length > 0" class="col-span-2 bg-blue-50 rounded-lg p-3">
              <div class="text-blue-600 text-xs mb-1">📎 Hóa đơn đính kèm</div>
              <div class="font-medium text-blue-800">{{ expense.attachments.length }} file</div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-2 mt-3">
            <button @click="editExpense(expense)" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              📝 Sửa
            </button>
            <button @click="deleteExpense(expense.id)" class="flex-1 bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              🗑️ Xóa
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredExpenses.length === 0" class="text-center py-20">
        <div class="text-6xl mb-4">💰</div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">Chưa có chi phí</h3>
        <p class="text-gray-600 mb-4">Hãy thêm chi phí đầu tiên của bạn</p>
        <button @click="showCreateForm = true" class="btn-primary">
          + Thêm chi phí
        </button>
      </div>

      <!-- Expense Report Modal -->
      <div v-if="showReportView" class="modal">
        <div class="modal-content max-w-4xl">
          <h3>📊 Báo cáo Chi phí</h3>
          
          <!-- Report Filters -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Tháng</label>
                <select v-model="reportMonth" class="w-full p-3 border border-gray-300 rounded-lg">
                  <option value="">Tất cả</option>
                  <option v-for="month in months" :key="month.value" :value="month.value">{{ month.label }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Loại chi phí</label>
                <select v-model="reportType" class="w-full p-3 border border-gray-300 rounded-lg">
                  <option value="">Tất cả</option>
                  <option value="fixed">Cố định</option>
                  <option value="variable">Biến đổi</option>
                  <option value="one-time">Một lần</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Summary Cards -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
            <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl p-4 text-white">
              <div class="text-sm opacity-90 mb-1">Tổng chi phí</div>
              <div class="text-2xl font-bold">{{ formatPrice(reportData.total) }}</div>
            </div>
            <div class="bg-gradient-to-br from-green-500 to-green-600 rounded-xl p-4 text-white">
              <div class="text-sm opacity-90 mb-1">Dự kiến</div>
              <div class="text-2xl font-bold">{{ formatPrice(reportData.estimated) }}</div>
            </div>
            <div class="bg-gradient-to-br rounded-xl p-4 text-white" :class="reportData.variance >= 0 ? 'from-red-500 to-red-600' : 'from-teal-500 to-teal-600'">
              <div class="text-sm opacity-90 mb-1">Chênh lệch</div>
              <div class="text-2xl font-bold">{{ reportData.variance >= 0 ? '+' : '' }}{{ formatPrice(reportData.variance) }}</div>
            </div>
          </div>

          <!-- By Category -->
          <div class="bg-white border border-gray-200 rounded-xl p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Theo nhóm chi phí</h4>
            <div class="space-y-2">
              <div v-for="cat in reportData.byCategory" :key="cat.category" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 rounded-lg flex items-center justify-center text-xl" :class="getCategoryColor(cat.category)">
                    {{ getCategoryIcon(cat.category) }}
                  </div>
                  <div>
                    <div class="font-medium text-gray-800">{{ cat.category }}</div>
                    <div class="text-xs text-gray-500">{{ cat.count }} chi phí</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-bold text-gray-800">{{ formatPrice(cat.total) }}</div>
                  <div class="text-xs text-gray-500">{{ cat.percentage }}%</div>
                </div>
              </div>
            </div>
          </div>

          <!-- By Month -->
          <div v-if="!reportMonth" class="bg-white border border-gray-200 rounded-xl p-4">
            <h4 class="font-semibold text-gray-800 mb-3">Theo tháng</h4>
            <div class="space-y-2">
              <div v-for="mon in reportData.byMonth" :key="mon.month" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">{{ formatMonth(mon.month) }}</div>
                  <div class="text-xs text-gray-500">{{ mon.count }} chi phí</div>
                </div>
                <div class="text-right">
                  <div class="font-bold text-gray-800">{{ formatPrice(mon.total) }}</div>
                </div>
              </div>
            </div>
          </div>

          <div class="form-actions mt-4">
            <button type="button" @click="showReportView = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Prepaid/Allocated Expense Modal -->
      <div v-if="showPrepaidForm" class="modal">
        <div class="modal-content">
          <h3>💳 Quản lý Chi phí Trả trước / Phân bổ</h3>
          
          <!-- Add New Prepaid Expense -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Thiết lập chi phí trả trước</h4>
            <form @submit.prevent="addPrepaidExpense">
              <div class="form-group">
                <label>Tên chi phí *</label>
                <input v-model="prepaidForm.name" type="text" required placeholder="Ví dụ: Tiền cọc nhà" />
              </div>
              <div class="form-group">
                <label>Nhóm chi phí *</label>
                <select v-model="prepaidForm.category" required>
                  <option value="">Chọn nhóm</option>
                  <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
                </select>
              </div>
              <div class="form-group">
                <label>Tổng số tiền (VNĐ) *</label>
                <input v-model.number="prepaidForm.total_amount" type="number" min="0" step="1000" required />
              </div>
              <div class="form-group">
                <label>Tháng bắt đầu *</label>
                <input v-model="prepaidForm.start_month" type="month" required />
              </div>
              <div class="form-group">
                <label>Số tháng phân bổ *</label>
                <input v-model.number="prepaidForm.allocation_months" type="number" min="1" max="60" required />
              </div>
              <div class="form-group">
                <label>Ngày thanh toán</label>
                <input v-model="prepaidForm.payment_date" type="date" />
              </div>
              <div class="form-group">
                <label>Ghi chú</label>
                <textarea v-model="prepaidForm.notes" rows="2" placeholder="Ghi chú thêm..."></textarea>
              </div>
              <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-3">
                <div class="text-sm text-blue-800">
                  <strong>Phân bổ:</strong> {{ formatPrice(prepaidForm.total_amount / (prepaidForm.allocation_months || 1)) }}/tháng
                </div>
              </div>
              <button type="submit" class="btn-primary w-full">+ Tạo chi phí trả trước</button>
            </form>
          </div>

          <!-- Prepaid Expense List -->
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div v-for="prep in prepaidExpenses" :key="prep.id" class="bg-white border border-gray-200 rounded-lg p-3">
              <div class="flex items-center justify-between mb-2">
                <div>
                  <div class="font-medium text-gray-800">{{ prep.name }}</div>
                  <div class="text-sm text-gray-500">{{ prep.category }} - {{ formatPrice(prep.total_amount) }}</div>
                </div>
                <button @click="deletePrepaidExpense(prep.id)" class="text-red-500 hover:text-red-700 p-2">
                  🗑️
                </button>
              </div>
              <div class="flex items-center justify-between text-xs">
                <span class="px-2 py-1 bg-indigo-100 text-indigo-800 rounded">{{ prep.allocation_months }} tháng</span>
                <span class="text-gray-500">{{ formatPrice(prep.monthly_amount) }}/tháng</span>
              </div>
              <div class="mt-2 text-xs text-gray-500">
                Từ {{ formatMonth(prep.start_month) }} - {{ formatMonth(prep.end_month) }}
              </div>
            </div>
          </div>

          <div class="form-actions mt-4">
            <button type="button" @click="showPrepaidForm = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Recurring Expense Management Modal -->
      <div v-if="showRecurringForm" class="modal">
        <div class="modal-content">
          <h3>🔁 Quản lý Chi phí Định kỳ</h3>
          
          <!-- Add New Recurring Expense -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Thiết lập chi phí định kỳ</h4>
            <form @submit.prevent="addRecurringExpense">
              <div class="form-group">
                <label>Tên chi phí *</label>
                <input v-model="recurringForm.name" type="text" required placeholder="Ví dụ: Tiền thuê nhà" />
              </div>
              <div class="form-group">
                <label>Nhóm chi phí *</label>
                <select v-model="recurringForm.category" required>
                  <option value="">Chọn nhóm</option>
                  <option v-for="cat in recurringCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
                </select>
              </div>
              <div class="form-group">
                <label>Số tiền dự kiến (VNĐ) *</label>
                <input v-model.number="recurringForm.estimated_amount" type="number" min="0" step="1000" required />
              </div>
              <div class="form-group">
                <label>Chu kỳ *</label>
                <select v-model="recurringForm.cycle" required>
                  <option value="monthly">Hàng tháng</option>
                  <option value="quarterly">Hàng quý</option>
                  <option value="yearly">Hàng năm</option>
                </select>
              </div>
              <div class="form-group">
                <label>Ngày nhắc nhập</label>
                <input v-model.number="recurringForm.reminder_day" type="number" min="1" max="31" placeholder="Ngày trong tháng (1-31)" />
              </div>
              <button type="submit" class="btn-primary w-full">+ Thêm chi phí định kỳ</button>
            </form>
          </div>

          <!-- Recurring Expense List -->
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div v-for="rec in recurringExpenses" :key="rec.id" class="bg-white border border-gray-200 rounded-lg p-3">
              <div class="flex items-center justify-between mb-2">
                <div>
                  <div class="font-medium text-gray-800">{{ rec.name }}</div>
                  <div class="text-sm text-gray-500">{{ rec.category }} - {{ formatPrice(rec.estimated_amount) }}</div>
                </div>
                <button @click="deleteRecurringExpense(rec.id)" class="text-red-500 hover:text-red-700 p-2">
                  🗑️
                </button>
              </div>
              <div class="flex items-center justify-between text-xs">
                <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded">{{ getCycleText(rec.cycle) }}</span>
                <span class="text-gray-500">Nhắc: Ngày {{ rec.reminder_day || 1 }}</span>
              </div>
            </div>
          </div>

          <div class="form-actions mt-4">
            <button type="button" @click="showRecurringForm = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Category Management Modal -->
      <div v-if="showCategoryForm" class="modal">
        <div class="modal-content">
          <h3>📁 Quản lý Loại Chi phí</h3>
          
          <!-- Add New Category -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Thêm loại mới</h4>
            <form @submit.prevent="addCategory">
              <div class="form-group">
                <label>Tên loại chi phí *</label>
                <input v-model="categoryForm.name" type="text" required placeholder="Ví dụ: Thuê nhà" />
              </div>
              <div class="form-group">
                <label>Loại *</label>
                <select v-model="categoryForm.type" required>
                  <option value="fixed">Cố định</option>
                  <option value="variable">Biến đổi</option>
                  <option value="one-time">Một lần</option>
                </select>
              </div>
              <div class="form-group">
                <label class="flex items-center space-x-2">
                  <input v-model="categoryForm.is_recurring" type="checkbox" class="w-4 h-4" />
                  <span>Chi phí định kỳ</span>
                </label>
              </div>
              <button type="submit" class="btn-primary w-full">+ Thêm loại</button>
            </form>
          </div>

          <!-- Category List -->
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div v-for="cat in expenseCategories" :key="cat.id" class="bg-white border border-gray-200 rounded-lg p-3 flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-lg flex items-center justify-center text-xl" :class="getCategoryColor(cat.name)">
                  {{ getCategoryIcon(cat.name) }}
                </div>
                <div>
                  <div class="font-medium text-gray-800">{{ cat.name }}</div>
                  <div class="text-xs text-gray-500">
                    {{ getExpenseType(cat.type) }}
                    <span v-if="cat.is_recurring" class="ml-2 px-2 py-0.5 bg-blue-100 text-blue-800 rounded">Định kỳ</span>
                  </div>
                </div>
              </div>
              <button @click="deleteCategory(cat.id)" class="text-red-500 hover:text-red-700 p-2">
                🗑️
              </button>
            </div>
          </div>

          <div class="form-actions mt-4">
            <button type="button" @click="showCategoryForm = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <div v-if="showCreateForm || editingExpense" class="modal">
        <div class="modal-content">
          <h3>{{ editingExpense ? 'Sửa chi phí' : 'Thêm chi phí mới' }}</h3>
          <form @submit.prevent="saveExpense">
            <div class="form-group">
              <label>Tên chi phí *</label>
              <input v-model="form.name" type="text" required placeholder="Ví dụ: Tiền thuê nhà" />
            </div>
            <div class="form-group">
              <label>Nhóm chi phí *</label>
              <select v-model="form.category" required>
                <option value="">Chọn nhóm</option>
                <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Số tiền (VNĐ) *</label>
              <input v-model.number="form.amount" type="number" min="0" step="1000" required />
            </div>
            <div class="form-group">
              <label>Tháng áp dụng *</label>
              <input v-model="form.month" type="month" required />
            </div>
            <div class="form-group">
              <label>Ngày phát sinh</label>
              <input v-model="form.date" type="date" />
            </div>
            <div class="form-group">
              <label>Loại chi phí</label>
              <select v-model="form.type">
                <option value="fixed">Cố định</option>
                <option value="variable">Biến đổi</option>
                <option value="one-time">Một lần</option>
              </select>
            </div>
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="form.notes" rows="3" placeholder="Ghi chú thêm..."></textarea>
            </div>
            <div class="form-group">
              <label>Đính kèm hóa đơn</label>
              <input type="file" @change="handleFileUpload" accept="image/*,.pdf" multiple class="w-full" />
              <small class="text-gray-500">Hỗ trợ: JPG, PNG, PDF (Tối đa 5 file)</small>
              
              <!-- Preview uploaded files -->
              <div v-if="form.attachments && form.attachments.length > 0" class="mt-3 space-y-2">
                <div v-for="(file, index) in form.attachments" :key="index" class="flex items-center justify-between bg-gray-50 p-2 rounded">
                  <div class="flex items-center space-x-2">
                    <span class="text-2xl">{{ getFileIcon(file.type) }}</span>
                    <span class="text-sm text-gray-700">{{ file.name }}</span>
                    <span class="text-xs text-gray-500">({{ formatFileSize(file.size) }})</span>
                  </div>
                  <button type="button" @click="removeFile(index)" class="text-red-500 hover:text-red-700">
                    ✕
                  </button>
                </div>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" @click="cancelEdit" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">{{ editingExpense ? 'Cập nhật' : 'Thêm' }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Navigation from '../components/Navigation.vue'

const showCreateForm = ref(false)
const showCategoryForm = ref(false)
const showRecurringForm = ref(false)
const showPrepaidForm = ref(false)
const showReportView = ref(false)
const editingExpense = ref(null)
const loading = ref(false)
const error = ref('')
const filterMonth = ref('')
const filterCategory = ref('')
const reportMonth = ref('')
const reportType = ref('')

const form = ref({
  name: '',
  category: '',
  amount: 0,
  month: '',
  date: '',
  type: 'fixed',
  notes: '',
  attachments: []
})

const categoryForm = ref({
  name: '',
  type: 'fixed',
  is_recurring: false
})

const recurringForm = ref({
  name: '',
  category: '',
  estimated_amount: 0,
  cycle: 'monthly',
  reminder_day: 1
})

const prepaidForm = ref({
  name: '',
  category: '',
  total_amount: 0,
  start_month: '',
  allocation_months: 12,
  payment_date: '',
  notes: ''
})

// Recurring Expenses (FR-EX-04)
const recurringExpenses = ref([
  { id: '1', name: 'Tiền thuê nhà', category: 'Thuê nhà', estimated_amount: 15000000, cycle: 'monthly', reminder_day: 5 },
  { id: '2', name: 'Tiền điện', category: 'Điện', estimated_amount: 2500000, cycle: 'monthly', reminder_day: 10 },
  { id: '3', name: 'Tiền nước', category: 'Nước', estimated_amount: 500000, cycle: 'monthly', reminder_day: 10 }
])

// Prepaid Expenses (FR-EX-05)
const prepaidExpenses = ref([
  { 
    id: '1', 
    name: 'Tiền cọc nhà', 
    category: 'Thuê nhà', 
    total_amount: 30000000, 
    start_month: '2024-01',
    allocation_months: 12,
    monthly_amount: 2500000,
    end_month: '2024-12',
    payment_date: '2024-01-01',
    notes: 'Cọc nhà 2 tháng'
  }
])

// Expense Categories (FR-EX-02)
const expenseCategories = ref([
  { id: '1', name: 'Thuê nhà', type: 'fixed', is_recurring: true },
  { id: '2', name: 'Điện', type: 'variable', is_recurring: true },
  { id: '3', name: 'Nước', type: 'variable', is_recurring: true },
  { id: '4', name: 'Internet', type: 'fixed', is_recurring: true },
  { id: '5', name: 'Lương', type: 'fixed', is_recurring: true },
  { id: '6', name: 'Nguyên liệu', type: 'variable', is_recurring: false },
  { id: '7', name: 'Bảo trì', type: 'variable', is_recurring: false },
  { id: '8', name: 'Khác', type: 'one-time', is_recurring: false }
])

// Mock data for now
const expenses = ref([
  {
    id: '1',
    name: 'Tiền thuê nhà',
    category: 'Thuê nhà',
    amount: 15000000,
    month: '2024-01',
    date: '2024-01-05',
    type: 'fixed',
    notes: 'Thuê nhà tháng 1'
  },
  {
    id: '2',
    name: 'Tiền điện',
    category: 'Điện',
    amount: 2500000,
    month: '2024-01',
    date: '2024-01-10',
    type: 'variable',
    notes: ''
  },
  {
    id: '3',
    name: 'Tiền nước',
    category: 'Nước',
    amount: 500000,
    month: '2024-01',
    date: '2024-01-10',
    type: 'variable',
    notes: ''
  }
])

const categories = computed(() => {
  return expenseCategories.value.map(c => c.name)
})

const recurringCategories = computed(() => {
  return expenseCategories.value.filter(c => c.is_recurring)
})

const recurringDueReminders = computed(() => {
  const today = new Date()
  const currentMonth = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}`
  const currentDay = today.getDate()
  
  return recurringExpenses.value.filter(rec => {
    // Check if already recorded this month
    const alreadyRecorded = expenses.value.some(e => 
      e.name === rec.name && e.month === currentMonth
    )
    
    // Show reminder if not recorded and past reminder day
    return !alreadyRecorded && currentDay >= rec.reminder_day
  }).map(rec => ({
    ...rec,
    due_month: currentMonth
  }))
})

const months = computed(() => {
  const monthSet = new Set(expenses.value.map(e => e.month))
  return Array.from(monthSet).map(m => ({
    value: m,
    label: formatMonth(m)
  }))
})

const filteredExpenses = computed(() => {
  let filtered = expenses.value
  
  if (filterMonth.value) {
    filtered = filtered.filter(e => e.month === filterMonth.value)
  }
  
  if (filterCategory.value) {
    filtered = filtered.filter(e => e.category === filterCategory.value)
  }
  
  return filtered.sort((a, b) => new Date(b.date) - new Date(a.date))
})

const totalExpenses = computed(() => {
  return filteredExpenses.value.reduce((sum, e) => sum + e.amount, 0)
})

const reportData = computed(() => {
  let filtered = expenses.value
  
  if (reportMonth.value) {
    filtered = filtered.filter(e => e.month === reportMonth.value)
  }
  
  if (reportType.value) {
    filtered = filtered.filter(e => e.type === reportType.value)
  }
  
  const total = filtered.reduce((sum, e) => sum + e.amount, 0)
  
  // Calculate estimated from recurring expenses
  let estimated = 0
  if (reportMonth.value) {
    recurringExpenses.value.forEach(rec => {
      const cat = expenseCategories.value.find(c => c.name === rec.category)
      if (!reportType.value || (cat && cat.type === reportType.value)) {
        estimated += rec.estimated_amount
      }
    })
  }
  
  // By category
  const categoryMap = {}
  filtered.forEach(e => {
    if (!categoryMap[e.category]) {
      categoryMap[e.category] = { category: e.category, total: 0, count: 0 }
    }
    categoryMap[e.category].total += e.amount
    categoryMap[e.category].count++
  })
  
  const byCategory = Object.values(categoryMap)
    .map(c => ({
      ...c,
      percentage: total > 0 ? ((c.total / total) * 100).toFixed(1) : 0
    }))
    .sort((a, b) => b.total - a.total)
  
  // By month
  const monthMap = {}
  expenses.value.forEach(e => {
    if (!reportType.value || e.type === reportType.value) {
      if (!monthMap[e.month]) {
        monthMap[e.month] = { month: e.month, total: 0, count: 0 }
      }
      monthMap[e.month].total += e.amount
      monthMap[e.month].count++
    }
  })
  
  const byMonth = Object.values(monthMap).sort((a, b) => b.month.localeCompare(a.month))
  
  return {
    total,
    estimated,
    variance: total - estimated,
    byCategory,
    byMonth
  }
})

const fetchExpenses = async () => {
  // TODO: Implement API call
  loading.value = true
  setTimeout(() => {
    loading.value = false
  }, 500)
}

const saveExpense = () => {
  if (editingExpense.value) {
    const index = expenses.value.findIndex(e => e.id === editingExpense.value.id)
    expenses.value[index] = { ...form.value, id: editingExpense.value.id }
  } else {
    expenses.value.push({
      ...form.value,
      id: Date.now().toString()
    })
  }
  cancelEdit()
}

const handleFileUpload = (event) => {
  const files = Array.from(event.target.files)
  
  if (files.length + (form.value.attachments?.length || 0) > 5) {
    alert('Chỉ được đính kèm tối đa 5 file!')
    return
  }
  
  files.forEach(file => {
    if (file.size > 5 * 1024 * 1024) {
      alert(`File ${file.name} quá lớn! Tối đa 5MB`)
      return
    }
    
    const reader = new FileReader()
    reader.onload = (e) => {
      if (!form.value.attachments) {
        form.value.attachments = []
      }
      form.value.attachments.push({
        name: file.name,
        type: file.type,
        size: file.size,
        data: e.target.result
      })
    }
    reader.readAsDataURL(file)
  })
  
  event.target.value = ''
}

const removeFile = (index) => {
  form.value.attachments.splice(index, 1)
}

const getFileIcon = (type) => {
  if (type.includes('pdf')) return '📄'
  if (type.includes('image')) return '🖼️'
  return '📎'
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const editExpense = (expense) => {
  editingExpense.value = expense
  form.value = { ...expense }
}

const deleteExpense = (id) => {
  if (confirm('Bạn có chắc muốn xóa chi phí này?')) {
    expenses.value = expenses.value.filter(e => e.id !== id)
  }
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingExpense.value = null
  form.value = {
    name: '',
    category: '',
    amount: 0,
    month: '',
    date: '',
    type: 'fixed',
    notes: '',
    attachments: []
  }
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('vi-VN')
}

const formatMonth = (month) => {
  if (!month) return ''
  const [year, m] = month.split('-')
  return `Tháng ${parseInt(m)}/${year}`
}

const getCategoryIcon = (category) => {
  const iconMap = {
    'Thuê nhà': '🏠',
    'Điện': '⚡',
    'Nước': '💧',
    'Internet': '🌐',
    'Lương': '💵',
    'Nguyên liệu': '🥬',
    'Bảo trì': '🔧',
    'Khác': '💰'
  }
  return iconMap[category] || '💰'
}

const getCategoryColor = (category) => {
  const colorMap = {
    'Thuê nhà': 'bg-blue-100 text-blue-600',
    'Điện': 'bg-yellow-100 text-yellow-600',
    'Nước': 'bg-cyan-100 text-cyan-600',
    'Internet': 'bg-purple-100 text-purple-600',
    'Lương': 'bg-green-100 text-green-600',
    'Nguyên liệu': 'bg-orange-100 text-orange-600',
    'Bảo trì': 'bg-red-100 text-red-600',
    'Khác': 'bg-gray-100 text-gray-600'
  }
  return colorMap[category] || 'bg-gray-100 text-gray-600'
}

const getExpenseType = (type) => {
  const typeMap = {
    'fixed': 'Cố định',
    'variable': 'Biến đổi',
    'one-time': 'Một lần'
  }
  return typeMap[type] || type
}

const addCategory = () => {
  if (!categoryForm.value.name) return
  
  expenseCategories.value.push({
    id: Date.now().toString(),
    name: categoryForm.value.name,
    type: categoryForm.value.type,
    is_recurring: categoryForm.value.is_recurring
  })
  
  categoryForm.value = {
    name: '',
    type: 'fixed',
    is_recurring: false
  }
}

const deleteCategory = (id) => {
  const category = expenseCategories.value.find(c => c.id === id)
  const hasExpenses = expenses.value.some(e => e.category === category.name)
  
  if (hasExpenses) {
    alert('Không thể xóa loại chi phí đã có giao dịch!')
    return
  }
  
  if (confirm(`Bạn có chắc muốn xóa loại "${category.name}"?`)) {
    expenseCategories.value = expenseCategories.value.filter(c => c.id !== id)
  }
}

const addRecurringExpense = () => {
  if (!recurringForm.value.name || !recurringForm.value.category) return
  
  recurringExpenses.value.push({
    id: Date.now().toString(),
    name: recurringForm.value.name,
    category: recurringForm.value.category,
    estimated_amount: recurringForm.value.estimated_amount,
    cycle: recurringForm.value.cycle,
    reminder_day: recurringForm.value.reminder_day || 1
  })
  
  recurringForm.value = {
    name: '',
    category: '',
    estimated_amount: 0,
    cycle: 'monthly',
    reminder_day: 1
  }
}

const deleteRecurringExpense = (id) => {
  if (confirm('Bạn có chắc muốn xóa chi phí định kỳ này?')) {
    recurringExpenses.value = recurringExpenses.value.filter(r => r.id !== id)
  }
}

const recordRecurringExpense = (reminder) => {
  form.value = {
    name: reminder.name,
    category: reminder.category,
    amount: reminder.estimated_amount,
    month: reminder.due_month,
    date: new Date().toISOString().split('T')[0],
    type: 'fixed',
    notes: 'Chi phí định kỳ',
    attachments: []
  }
  showCreateForm.value = true
}

const getCycleText = (cycle) => {
  const cycleMap = {
    'monthly': 'Hàng tháng',
    'quarterly': 'Hàng quý',
    'yearly': 'Hàng năm'
  }
  return cycleMap[cycle] || cycle
}

const addPrepaidExpense = () => {
  if (!prepaidForm.value.name || !prepaidForm.value.category || !prepaidForm.value.start_month) return
  
  const monthlyAmount = prepaidForm.value.total_amount / prepaidForm.value.allocation_months
  const startDate = new Date(prepaidForm.value.start_month + '-01')
  const endDate = new Date(startDate)
  endDate.setMonth(endDate.getMonth() + prepaidForm.value.allocation_months - 1)
  const endMonth = `${endDate.getFullYear()}-${String(endDate.getMonth() + 1).padStart(2, '0')}`
  
  const prepaidId = Date.now().toString()
  
  prepaidExpenses.value.push({
    id: prepaidId,
    name: prepaidForm.value.name,
    category: prepaidForm.value.category,
    total_amount: prepaidForm.value.total_amount,
    start_month: prepaidForm.value.start_month,
    allocation_months: prepaidForm.value.allocation_months,
    monthly_amount: monthlyAmount,
    end_month: endMonth,
    payment_date: prepaidForm.value.payment_date,
    notes: prepaidForm.value.notes
  })
  
  // Auto-generate allocated expenses for each month
  for (let i = 0; i < prepaidForm.value.allocation_months; i++) {
    const allocDate = new Date(startDate)
    allocDate.setMonth(allocDate.getMonth() + i)
    const allocMonth = `${allocDate.getFullYear()}-${String(allocDate.getMonth() + 1).padStart(2, '0')}`
    
    expenses.value.push({
      id: `${prepaidId}-${i}`,
      name: `${prepaidForm.value.name} (Phân bổ ${i + 1}/${prepaidForm.value.allocation_months})`,
      category: prepaidForm.value.category,
      amount: monthlyAmount,
      month: allocMonth,
      date: `${allocMonth}-01`,
      type: 'one-time',
      notes: prepaidForm.value.notes,
      is_allocated: true,
      prepaid_id: prepaidId,
      original_amount: prepaidForm.value.total_amount,
      allocated_month_index: i + 1,
      total_months: prepaidForm.value.allocation_months
    })
  }
  
  prepaidForm.value = {
    name: '',
    category: '',
    total_amount: 0,
    start_month: '',
    allocation_months: 12,
    payment_date: '',
    notes: ''
  }
}

const deletePrepaidExpense = (id) => {
  const prepaid = prepaidExpenses.value.find(p => p.id === id)
  
  if (confirm(`Bạn có chắc muốn xóa chi phí trả trước "${prepaid.name}"? Tất cả chi phí phân bổ sẽ bị xóa.`)) {
    // Remove all allocated expenses
    expenses.value = expenses.value.filter(e => e.prepaid_id !== id)
    // Remove prepaid expense
    prepaidExpenses.value = prepaidExpenses.value.filter(p => p.id !== id)
  }
}

onMounted(() => {
  fetchExpenses()
})
</script>

<style scoped>
.modal {
  @apply fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50 p-4;
}

.modal-content {
  @apply bg-white p-6 lg:p-8 rounded-xl w-full max-w-lg max-h-screen overflow-y-auto;
}

.modal-content.max-w-4xl {
  max-width: 56rem;
}

.modal-content h3 {
  @apply text-lg lg:text-xl font-semibold text-gray-800 mb-5 text-center;
}

.form-group {
  @apply mb-4;
}

.form-group label {
  @apply block mb-2 font-medium text-gray-700 text-sm;
}

.form-group input,
.form-group select,
.form-group textarea {
  @apply w-full p-3 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
}

.form-actions {
  @apply flex flex-col lg:flex-row gap-3 justify-end mt-6;
}

.form-actions button {
  @apply w-full lg:w-auto px-5 py-2 rounded-lg font-medium transition-colors duration-200;
}

.btn-save {
  @apply bg-green-600 text-white hover:bg-green-700;
}

.btn-cancel {
  @apply bg-gray-600 text-white hover:bg-gray-700;
}
</style>
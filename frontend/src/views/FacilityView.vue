<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Header -->
      <div class="flex flex-col lg:flex-row justify-between items-center mb-6">
        <h2 class="text-xl lg:text-2xl font-semibold text-gray-800 mb-4 lg:mb-0">
          🏢 Quản lý Cơ sở vật chất
        </h2>
        <div class="flex flex-wrap gap-2">
          <button @click="showMaintenanceSchedule = true" class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📅 Lịch bảo trì
          </button>
          <button @click="generateAssetReport" class="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📊 Báo cáo
          </button>
          <button @click="showCreateForm = true" class="btn-primary text-sm px-4 py-2">
            + Thêm tài sản
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-10 text-gray-600 text-lg">Đang tải...</div>
      <div v-if="error" class="text-center py-10 text-red-600 bg-red-50 border border-red-200 rounded-lg">{{ error }}</div>

      <!-- Filters -->
      <div class="bg-white rounded-xl p-4 mb-4 shadow-sm">
        <div class="grid grid-cols-1 gap-3">
          <input v-model="searchQuery" type="text" placeholder="Tìm kiếm tài sản..." class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500" />
          <select v-model="filterType" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500">
            <option value="">Tất cả loại</option>
            <option value="Bàn ghế">Bàn ghế</option>
            <option value="Máy móc">Máy móc</option>
            <option value="Dụng cụ">Dụng cụ</option>
            <option value="Điện tử">Điện tử</option>
            <option value="Khác">Khác</option>
          </select>
          <select v-model="filterStatus" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500">
            <option value="">Tất cả trạng thái</option>
            <option value="Đang sử dụng">Đang sử dụng</option>
            <option value="Hỏng">Hỏng</option>
            <option value="Đang sửa">Đang sửa</option>
            <option value="Ngừng sử dụng">Ngừng sử dụng</option>
          </select>
          <select v-model="filterArea" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500">
            <option value="">Tất cả khu vực</option>
            <option value="Phòng khách">Phòng khách</option>
            <option value="Bếp">Bếp</option>
            <option value="Quầy bar">Quầy bar</option>
            <option value="Kho">Kho</option>
            <option value="Văn phòng">Văn phòng</option>
          </select>
        </div>
      </div>

      <!-- Summary Card -->
      <div class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl p-6 mb-4 text-white shadow-lg">
        <div class="text-sm opacity-90 mb-1">Tổng số tài sản</div>
        <div class="text-3xl font-bold">{{ filteredItems.length }}</div>
      </div>

      <!-- Facility List -->
      <div class="grid grid-cols-1 gap-4">
        <div v-for="item in filteredItems" :key="item.id" class="bg-white rounded-xl p-4 shadow-md">
          <!-- Facility Header -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" :class="getTypeColor(item.type)">
                {{ getTypeIcon(item.type) }}
              </div>
              <div>
                <h4 class="font-bold text-gray-800">{{ item.name }}</h4>
                <p class="text-sm text-gray-500">{{ item.area }}</p>
              </div>
            </div>
            <div class="text-right">
              <span class="px-3 py-1 rounded-full text-xs font-medium" :class="getStatusBadge(item.status)">
                {{ item.status }}
              </span>
            </div>
          </div>

          <!-- Facility Details -->
          <div class="grid grid-cols-2 gap-3 text-sm mb-3">
            <div class="bg-gray-50 rounded-lg p-3">
              <div class="text-gray-500 text-xs mb-1">Loại</div>
              <div class="font-medium text-gray-800">{{ item.type }}</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <div class="text-gray-500 text-xs mb-1">Số lượng</div>
              <div class="font-medium text-gray-800">{{ item.quantity }}</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <div class="text-gray-500 text-xs mb-1">Ngày mua</div>
              <div class="font-medium text-gray-800">{{ formatDate(item.purchase_date) }}</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3" v-if="item.cost">
              <div class="text-gray-500 text-xs mb-1">Giá trị</div>
              <div class="font-medium text-gray-800">{{ formatPrice(item.cost) }}</div>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-2">
            <button @click="editItem(item)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              📝 Sửa
            </button>
            <button @click="showStatusUpdate(item)" class="bg-purple-500 hover:bg-purple-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              🔄 Trạng thái
            </button>
            <button @click="showMaintenance(item)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              🔧 Bảo trì
            </button>
            <button @click="moveAsset(item)" class="bg-indigo-500 hover:bg-indigo-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              🚚 Di chuyển
            </button>
            <button @click="reportIssue(item)" class="bg-orange-500 hover:bg-orange-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              ⚠️ Báo hỏng
            </button>
            <button @click="showHistory(item)" class="bg-teal-500 hover:bg-teal-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              📜 Lịch sử
            </button>
            <button @click="deleteItem(item.id)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
              🗑️ Xóa
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredItems.length === 0" class="text-center py-20">
        <div class="text-6xl mb-4">🏢</div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">Chưa có tài sản</h3>
        <p class="text-gray-600 mb-4">Hãy thêm tài sản đầu tiên của bạn</p>
        <button @click="showCreateForm = true" class="btn-primary">+ Thêm tài sản</button>
      </div>

      <!-- Create/Edit Modal -->
      <div v-if="showCreateForm || editingItem" class="modal">
        <div class="modal-content">
          <h3>{{ editingItem ? 'Sửa tài sản' : 'Thêm tài sản mới' }}</h3>
          <form @submit.prevent="saveItem">
            <div class="form-group">
              <label>Tên tài sản *</label>
              <input v-model="form.name" type="text" required placeholder="Ví dụ: Bàn gỗ" />
            </div>
            <div class="form-group">
              <label>Loại *</label>
              <select v-model="form.type" required>
                <option value="">Chọn loại</option>
                <option value="Bàn ghế">Bàn ghế</option>
                <option value="Máy móc">Máy móc</option>
                <option value="Dụng cụ">Dụng cụ</option>
                <option value="Điện tử">Điện tử</option>
                <option value="Khác">Khác</option>
              </select>
            </div>
            <div class="form-group">
              <label>Khu vực *</label>
              <select v-model="form.area" required>
                <option value="">Chọn khu vực</option>
                <option value="Phòng khách">Phòng khách</option>
                <option value="Bếp">Bếp</option>
                <option value="Quầy bar">Quầy bar</option>
                <option value="Kho">Kho</option>
                <option value="Văn phòng">Văn phòng</option>
                <option value="Khác">Khác</option>
              </select>
            </div>
            <div class="form-group">
              <label>Số lượng *</label>
              <input v-model.number="form.quantity" type="number" min="1" required />
            </div>
            <div class="form-group">
              <label>Trạng thái *</label>
              <select v-model="form.status" required>
                <option value="Đang sử dụng">Đang sử dụng</option>
                <option value="Hỏng">Hỏng</option>
                <option value="Đang sửa">Đang sửa</option>
                <option value="Ngừng sử dụng">Ngừng sử dụng</option>
              </select>
            </div>
            <div class="form-group">
              <label>Ngày mua</label>
              <input v-model="form.purchase_date" type="date" />
            </div>
            <div class="form-group">
              <label>Chi phí (VNĐ)</label>
              <input v-model.number="form.cost" type="number" min="0" step="1000" />
            </div>
            <div class="form-group">
              <label>Nhà cung cấp</label>
              <input v-model="form.supplier" type="text" placeholder="Tên nhà cung cấp" />
            </div>
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="form.notes" rows="3" placeholder="Ghi chú thêm..."></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="cancelEdit" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">{{ editingItem ? 'Cập nhật' : 'Thêm' }}</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Maintenance Modal -->
      <div v-if="maintenanceItem" class="modal">
        <div class="modal-content max-w-4xl">
          <h3>🔧 Bảo trì: {{ maintenanceItem.name }}</h3>
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <button @click="showMaintenanceForm = true" class="btn-primary w-full">+ Thêm bảo trì mới</button>
          </div>
          <div v-if="maintenanceRecords.length === 0" class="text-center py-10 text-gray-600">Chưa có lịch sử bảo trì</div>
          <div v-else class="space-y-3 max-h-96 overflow-y-auto">
            <div v-for="record in maintenanceRecords" :key="record.id" class="bg-white border border-gray-200 rounded-lg p-4">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">{{ getMaintenanceTypeText(record.type) }}</span>
                  <span class="text-sm text-gray-500 ml-2">{{ formatDate(record.date) }}</span>
                </div>
                <div class="text-lg font-bold text-gray-800">{{ formatPrice(record.cost) }}</div>
              </div>
              <p class="text-sm text-gray-700 mb-1"><strong>Mô tả:</strong> {{ record.description }}</p>
              <p v-if="record.vendor" class="text-sm text-gray-600"><strong>Đơn vị:</strong> {{ record.vendor }}</p>
            </div>
          </div>
          <div class="form-actions mt-4">
            <button @click="maintenanceItem = null" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Move Asset Modal -->
      <div v-if="movingAsset" class="modal">
        <div class="modal-content">
          <h3>🚚 Di chuyển tài sản: {{ movingAsset.name }}</h3>
          <p class="text-sm text-gray-600 mb-4">Khu vực hiện tại: <strong>{{ movingAsset.area }}</strong></p>
          <form @submit.prevent="saveMoveAsset">
            <div class="form-group">
              <label>Khu vực mới *</label>
              <select v-model="moveForm.new_area" required>
                <option value="">Chọn khu vực mới</option>
                <option value="Phòng khách">Phòng khách</option>
                <option value="Bếp">Bếp</option>
                <option value="Quầy bar">Quầy bar</option>
                <option value="Kho">Kho</option>
                <option value="Văn phòng">Văn phòng</option>
                <option value="Khác">Khác</option>
              </select>
            </div>
            <div class="form-group">
              <label>Lý do di chuyển *</label>
              <textarea v-model="moveForm.reason" required rows="3" placeholder="Nhập lý do di chuyển tài sản..."></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="movingAsset = null" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Di chuyển</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Maintenance Schedule Modal -->
      <div v-if="showMaintenanceSchedule" class="modal">
        <div class="modal-content max-w-4xl">
          <h3>📅 Lịch Bảo trì</h3>
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <button @click="showScheduleForm = true" class="btn-primary w-full">+ Lên lịch bảo trì mới</button>
          </div>
          <div v-if="scheduledTasks.length === 0" class="text-center py-10 text-gray-600">Chưa có lịch bảo trì</div>
          <div v-else class="space-y-3 max-h-96 overflow-y-auto">
            <div v-for="task in scheduledTasks" :key="task.id" class="bg-white border border-gray-200 rounded-lg p-4">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h5 class="font-bold text-gray-800">{{ task.facility_name }}</h5>
                  <p class="text-sm text-gray-600">{{ task.description }}</p>
                </div>
                <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">{{ formatDate(task.scheduled_date) }}</span>
              </div>
              <div class="flex gap-2 mt-3">
                <button @click="completeTask(task)" class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm">✓ Hoàn thành</button>
              </div>
            </div>
          </div>
          <div class="form-actions mt-4">
            <button @click="showMaintenanceSchedule = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Schedule Form Modal -->
      <div v-if="showScheduleForm" class="modal">
        <div class="modal-content">
          <h3>Lên lịch bảo trì</h3>
          <form @submit.prevent="saveScheduledTask">
            <div class="form-group">
              <label>Tài sản *</label>
              <select v-model="scheduleForm.facility_id" required>
                <option value="">Chọn tài sản</option>
                <option v-for="item in items" :key="item.id" :value="item.id">{{ item.name }} - {{ item.area }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mô tả công việc *</label>
              <textarea v-model="scheduleForm.description" required rows="3"></textarea>
            </div>
            <div class="form-group">
              <label>Ngày dự kiến *</label>
              <input v-model="scheduleForm.scheduled_date" type="date" required />
            </div>
            <div class="form-group">
              <label>Chi phí dự kiến (VNĐ)</label>
              <input v-model.number="scheduleForm.estimated_cost" type="number" min="0" step="1000" />
            </div>
            <div class="form-actions">
              <button type="button" @click="showScheduleForm = false" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Lên lịch</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Asset Report Modal -->
      <div v-if="showAssetReport" class="modal">
        <div class="modal-content max-w-4xl">
          <h3>📊 Báo cáo Tài sản</h3>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
            <div class="bg-blue-50 rounded-lg p-4">
              <div class="text-sm text-blue-600 mb-1">Tổng số tài sản</div>
              <div class="text-2xl font-bold text-blue-800">{{ assetReport.total }}</div>
            </div>
            <div class="bg-green-50 rounded-lg p-4">
              <div class="text-sm text-green-600 mb-1">Đang sử dụng</div>
              <div class="text-2xl font-bold text-green-800">{{ assetReport.active }}</div>
            </div>
            <div class="bg-red-50 rounded-lg p-4">
              <div class="text-sm text-red-600 mb-1">Hỏng hóc</div>
              <div class="text-2xl font-bold text-red-800">{{ assetReport.broken }}</div>
            </div>
          </div>
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div class="bg-white border border-gray-200 rounded-lg p-4">
              <h4 class="font-semibold text-gray-800 mb-3">Theo khu vực</h4>
              <div v-for="(count, area) in assetReport.byArea" :key="area" class="flex justify-between py-1">
                <span class="text-gray-600">{{ area }}</span>
                <span class="font-medium text-gray-800">{{ count }}</span>
              </div>
            </div>
            <div class="bg-white border border-gray-200 rounded-lg p-4">
              <h4 class="font-semibold text-gray-800 mb-3">Theo loại</h4>
              <div v-for="(count, type) in assetReport.byType" :key="type" class="flex justify-between py-1">
                <span class="text-gray-600">{{ type }}</span>
                <span class="font-medium text-gray-800">{{ count }}</span>
              </div>
            </div>
          </div>
          <div class="form-actions mt-4">
            <button @click="showAssetReport = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Status Update Modal (FR-FM-04) -->
      <div v-if="statusUpdateItem" class="modal">
        <div class="modal-content">
          <h3>🔄 Cập nhật trạng thái: {{ statusUpdateItem.name }}</h3>
          <form @submit.prevent="saveStatusUpdate">
            <div class="form-group">
              <label>Trạng thái mới *</label>
              <select v-model="statusForm.status" required>
                <option value="Đang sử dụng">Đang sử dụng</option>
                <option value="Hỏng">Hỏng</option>
                <option value="Đang sửa">Đang sửa</option>
                <option value="Ngừng sử dụng">Ngừng sử dụng</option>
                <option value="Thanh lý">Thanh lý</option>
              </select>
            </div>
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="statusForm.notes" rows="3" placeholder="Lý do thay đổi trạng thái..."></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="statusUpdateItem = null" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Cập nhật</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Issue Report Modal (FR-FM-05) -->
      <div v-if="reportingIssue" class="modal">
        <div class="modal-content">
          <h3>⚠️ Báo cáo sự cố: {{ reportingIssue.name }}</h3>
          <form @submit.prevent="saveIssueReport">
            <div class="form-group">
              <label>Loại sự cố *</label>
              <select v-model="issueForm.issue_type" required>
                <option value="">Chọn loại sự cố</option>
                <option value="hư hỏng">Hư hỏng</option>
                <option value="không hoạt động">Không hoạt động</option>
                <option value="hoạt động bất thường">Hoạt động bất thường</option>
                <option value="an toàn">Vấn đề an toàn</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mức độ ảnh hưởng *</label>
              <select v-model="issueForm.severity" required>
                <option value="">Chọn mức độ</option>
                <option value="thấp">Thấp</option>
                <option value="trung bình">Trung bình</option>
                <option value="cao">Cao</option>
                <option value="khẩn cấp">Khẩn cấp</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mô tả sự cố *</label>
              <textarea v-model="issueForm.description" required rows="4" placeholder="Mô tả chi tiết sự cố..."></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="reportingIssue = null" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Gửi báo cáo</button>
            </div>
          </form>
        </div>
      </div>

      <!-- History Modal (FR-FM-07) -->
      <div v-if="historyItem" class="modal">
        <div class="modal-content max-w-4xl">
          <h3>📜 Lịch sử: {{ historyItem.name }}</h3>
          <div v-if="facilityHistory.length === 0" class="text-center py-10 text-gray-600">Chưa có lịch sử</div>
          <div v-else class="space-y-3 max-h-96 overflow-y-auto">
            <div v-for="history in facilityHistory" :key="history.id" class="bg-white border border-gray-200 rounded-lg p-4">
              <div class="flex justify-between items-start mb-2">
                <span class="px-2 py-1 bg-gray-100 text-gray-800 rounded text-xs font-medium">{{ history.action_type }}</span>
                <span class="text-sm text-gray-500">{{ formatDateTime(history.created_at) }}</span>
              </div>
              <p class="text-sm text-gray-700 mb-1"><strong>Mô tả:</strong> {{ history.description }}</p>
              <p class="text-sm text-gray-600"><strong>Người thực hiện:</strong> {{ history.username }}</p>
            </div>
          </div>
          <div class="form-actions mt-4">
            <button @click="historyItem = null" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Maintenance Form Modal -->
      <div v-if="showMaintenanceForm" class="modal">
        <div class="modal-content">
          <h3>Thêm bảo trì mới</h3>
          <form @submit.prevent="saveMaintenanceRecord">
            <div class="form-group">
              <label>Loại bảo trì *</label>
              <select v-model="maintenanceForm.type" required>
                <option value="scheduled">Định kỳ</option>
                <option value="preventive">Phòng ngừa</option>
                <option value="corrective">Sửa chữa</option>
                <option value="emergency">Khẩn cấp</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mô tả công việc *</label>
              <textarea v-model="maintenanceForm.description" required rows="3"></textarea>
            </div>
            <div class="form-group">
              <label>Ngày thực hiện *</label>
              <input v-model="maintenanceForm.date" type="date" required />
            </div>
            <div class="form-group">
              <label>Chi phí (VNĐ)</label>
              <input v-model.number="maintenanceForm.cost" type="number" min="0" step="1000" />
            </div>
            <div class="form-group">
              <label>Đơn vị thực hiện</label>
              <input v-model="maintenanceForm.vendor" type="text" />
            </div>
            <div class="form-actions">
              <button type="button" @click="showMaintenanceForm = false" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Lưu</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useFacilityStore } from '../stores/facility'
import Navigation from '../components/Navigation.vue'

const facilityStore = useFacilityStore()

const showCreateForm = ref(false)
const editingItem = ref(null)
const maintenanceItem = ref(null)
const maintenanceRecords = ref([])
const showMaintenanceForm = ref(false)
const searchQuery = ref('')
const filterType = ref('')
const filterStatus = ref('')
const filterArea = ref('')
const statusUpdateItem = ref(null)
const reportingIssue = ref(null)
const historyItem = ref(null)
const facilityHistory = ref([])
const movingAsset = ref(null)
const showMaintenanceSchedule = ref(false)
const showScheduleForm = ref(false)
const scheduledTasks = ref([])
const showAssetReport = ref(false)
const assetReport = ref({})

const form = ref({
  name: '', type: '', area: '', quantity: 1, status: 'Đang sử dụng',
  purchase_date: '', cost: 0, supplier: '', notes: ''
})

const maintenanceForm = ref({
  type: 'scheduled', description: '', date: '', cost: 0, vendor: ''
})

const statusForm = ref({
  status: '', notes: ''
})

const issueForm = ref({
  issue_type: '', severity: '', description: ''
})

const moveForm = ref({
  new_area: '', reason: ''
})

const scheduleForm = ref({
  facility_id: '', description: '', scheduled_date: '', estimated_cost: 0
})

const items = computed(() => facilityStore.items || [])
const loading = computed(() => facilityStore.loading)
const error = computed(() => facilityStore.error)

const filteredItems = computed(() => {
  let filtered = items.value
  if (searchQuery.value) {
    filtered = filtered.filter(item =>
      item.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.type.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.area.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  if (filterType.value) {
    filtered = filtered.filter(item => item.type === filterType.value)
  }
  if (filterStatus.value) {
    filtered = filtered.filter(item => item.status === filterStatus.value)
  }
  if (filterArea.value) {
    filtered = filtered.filter(item => item.area === filterArea.value)
  }
  return filtered
})

onMounted(async () => {
  await facilityStore.fetchFacilities()
  try {
    scheduledTasks.value = await facilityStore.fetchScheduledMaintenance()
  } catch (error) {
    scheduledTasks.value = []
  }
})

const saveItem = async () => {
  const success = editingItem.value
    ? await facilityStore.updateFacility(editingItem.value.id, form.value)
    : await facilityStore.createFacility(form.value)
  if (success) cancelEdit()
}

const editItem = (item) => {
  editingItem.value = item
  form.value = { ...item }
}

const deleteItem = async (id) => {
  if (confirm('Bạn có chắc muốn xóa tài sản này?')) {
    await facilityStore.deleteFacility(id)
  }
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingItem.value = null
  form.value = {
    name: '', type: '', area: '', quantity: 1, status: 'Đang sử dụng',
    purchase_date: '', cost: 0, supplier: '', notes: ''
  }
}

const showMaintenance = async (item) => {
  maintenanceItem.value = item
  try {
    maintenanceRecords.value = await facilityStore.fetchMaintenanceHistory(item.id)
  } catch (error) {
    maintenanceRecords.value = []
  }
}

const saveMaintenanceRecord = async () => {
  const record = {
    ...maintenanceForm.value,
    facility_id: maintenanceItem.value.id
  }
  const success = await facilityStore.createMaintenanceRecord(record)
  if (success) {
    showMaintenanceForm.value = false
    maintenanceForm.value = { type: 'scheduled', description: '', date: '', cost: 0, vendor: '' }
    maintenanceRecords.value = await facilityStore.fetchMaintenanceHistory(maintenanceItem.value.id)
  }
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('vi-VN')
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const getTypeIcon = (type) => {
  const iconMap = {
    'Bàn ghế': '🪑',
    'Máy móc': '⚙️',
    'Dụng cụ': '🔧',
    'Điện tử': '💻',
    'Khác': '📦'
  }
  return iconMap[type] || '📦'
}

const getTypeColor = (type) => {
  const colorMap = {
    'Bàn ghế': 'bg-blue-100 text-blue-600',
    'Máy móc': 'bg-green-100 text-green-600',
    'Dụng cụ': 'bg-yellow-100 text-yellow-600',
    'Điện tử': 'bg-purple-100 text-purple-600',
    'Khác': 'bg-gray-100 text-gray-600'
  }
  return colorMap[type] || 'bg-gray-100 text-gray-600'
}

const getStatusBadge = (status) => {
  const badgeMap = {
    'Đang sử dụng': 'bg-green-100 text-green-800',
    'Hỏng': 'bg-red-100 text-red-800',
    'Đang sửa': 'bg-yellow-100 text-yellow-800',
    'Ngừng sử dụng': 'bg-gray-100 text-gray-800'
  }
  return badgeMap[status] || 'bg-gray-100 text-gray-600'
}

const getMaintenanceTypeText = (type) => {
  const typeMap = {
    'scheduled': 'Định kỳ',
    'preventive': 'Phòng ngừa',
    'corrective': 'Sửa chữa',
    'emergency': 'Khẩn cấp'
  }
  return typeMap[type] || type
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleString('vi-VN')
}

// FR-FM-04: Update status
const showStatusUpdate = (item) => {
  statusUpdateItem.value = item
  statusForm.value = { status: item.status, notes: '' }
}

const saveStatusUpdate = async () => {
  const success = await facilityStore.updateFacilityStatus(statusUpdateItem.value.id, statusForm.value)
  if (success) {
    statusUpdateItem.value = null
    statusForm.value = { status: '', notes: '' }
  }
}

// FR-FM-05: Report issue
const reportIssue = (item) => {
  reportingIssue.value = item
  issueForm.value = { issue_type: '', severity: '', description: '' }
}

const saveIssueReport = async () => {
  const reportData = {
    ...issueForm.value,
    facility_id: reportingIssue.value.id,
    facility_name: reportingIssue.value.name
  }
  const success = await facilityStore.createIssueReport(reportData)
  if (success) {
    reportingIssue.value = null
    issueForm.value = { issue_type: '', severity: '', description: '' }
    alert('Báo cáo sự cố đã được gửi thành công!')
  }
}

// FR-FM-07: Show history
const showHistory = async (item) => {
  historyItem.value = item
  try {
    facilityHistory.value = await facilityStore.fetchFacilityHistory(item.id)
  } catch (error) {
    facilityHistory.value = []
  }
}

// Move asset
const moveAsset = (item) => {
  movingAsset.value = item
  moveForm.value = { new_area: '', reason: '' }
}

const saveMoveAsset = async () => {
  const moveData = {
    area: moveForm.value.new_area,
    change_reason: `Di chuyển từ ${movingAsset.value.area} đến ${moveForm.value.new_area}: ${moveForm.value.reason}`
  }
  const success = await facilityStore.updateFacility(movingAsset.value.id, moveData)
  if (success) {
    movingAsset.value = null
    moveForm.value = { new_area: '', reason: '' }
  }
}

// Maintenance schedule
const saveScheduledTask = async () => {
  const success = await facilityStore.scheduleMaintenanceTask(scheduleForm.value)
  if (success) {
    showScheduleForm.value = false
    scheduleForm.value = { facility_id: '', description: '', scheduled_date: '', estimated_cost: 0 }
    scheduledTasks.value = await facilityStore.fetchScheduledMaintenance()
  }
}

const completeTask = async (task) => {
  if (confirm('Đánh dấu nhiệm vụ này đã hoàn thành?')) {
    const success = await facilityStore.updateMaintenanceTask(task.id, { status: 'completed' })
    if (success) {
      scheduledTasks.value = await facilityStore.fetchScheduledMaintenance()
    }
  }
}

// Asset report
const generateAssetReport = () => {
  const report = {
    total: items.value.length,
    active: items.value.filter(i => i.status === 'Đang sử dụng').length,
    needMaintenance: items.value.filter(i => i.status === 'Đang sửa').length,
    broken: items.value.filter(i => i.status === 'Hỏng').length,
    byArea: {},
    byType: {}
  }
  items.value.forEach(item => {
    report.byArea[item.area] = (report.byArea[item.area] || 0) + 1
    report.byType[item.type] = (report.byType[item.type] || 0) + 1
  })
  assetReport.value = report
  showAssetReport.value = true
}
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

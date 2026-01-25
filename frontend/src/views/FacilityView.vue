<template>
  <div class="facility-management">
    <Navigation />
    <div class="content">
      <div class="header">
        <h2>🏢 Quản lý Cơ sở vật chất</h2>
        <div class="header-actions">
          <button @click="showMaintenanceSchedule = true" class="btn-schedule">📅 Lịch bảo trì</button>
          <button @click="showCreateForm = true" class="btn-primary">+ Thêm tài sản</button>
        </div>
      </div>

      <div v-if="loading" class="loading">Đang tải...</div>
      <div v-if="error" class="error">{{ error }}</div>

      <div class="filters-section">
        <input v-model="searchQuery" type="text" placeholder="Tìm kiếm..." class="search-input" @input="searchWithFilters" />
        <select v-model="filterType" class="filter-select" @change="searchWithFilters">
          <option value="">Tất cả loại</option>
          <option value="Bàn ghế">Bàn ghế</option>
          <option value="Máy móc">Máy móc</option>
          <option value="Dụng cụ">Dụng cụ</option>
          <option value="Điện tử">Điện tử</option>
          <option value="Khác">Khác</option>
        </select>
        <select v-model="filterStatus" class="filter-select" @change="searchWithFilters">
          <option value="">Tất cả trạng thái</option>
          <option value="Đang sử dụng">Đang sử dụng</option>
          <option value="Hỏng">Hỏng</option>
          <option value="Đang sửa">Đang sửa</option>
          <option value="Ngừng sử dụng">Ngừng sử dụng</option>
          <option value="Thanh lý">Thanh lý</option>
        </select>
        <select v-model="filterArea" class="filter-select" @change="searchWithFilters">
          <option value="">Tất cả khu vực</option>
          <option value="Phòng khách">Phòng khách</option>
          <option value="Bếp">Bếp</option>
          <option value="Quầy bar">Quầy bar</option>
          <option value="Kho">Kho</option>
          <option value="Văn phòng">Văn phòng</option>
          <option value="Khác">Khác</option>
        </select>
        <button @click="searchWithFilters" class="btn-search">🔍 Tìm kiếm</button>
        <button @click="generateAssetReport" class="btn-report">📈 Báo cáo</button>
        <button @click="showMaintenanceDue" class="btn-due">⚠️ Đến hạn</button>
        <button @click="showStatusDashboard" class="btn-status-dashboard">📊 Trạng thái</button>
        <button @click="showIssueReports" class="btn-issues">📝 Báo cáo sự cố</button>
      </div>

      <div class="facility-grid">
        <div v-for="item in filteredItems" :key="item.id" class="facility-card">
          <div class="facility-info">
            <h4>{{ item.name }}</h4>
            <p><strong>Loại:</strong> {{ item.type }}</p>
            <p><strong>Khu vực:</strong> {{ item.area }}</p>
            <p><strong>Số lượng:</strong> {{ item.quantity }}</p>
            <p><strong>Trạng thái:</strong> 
              <span :class="'status-' + getStatusClass(item.status)">{{ item.status }}</span>
              <span v-if="getStatusAge(item.last_status_change)" class="status-age">
                ({{ getStatusAge(item.last_status_change) }})
              </span>
            </p>
            <p><strong>Ngày mua:</strong> {{ formatDate(item.purchase_date) }}</p>
            <p v-if="item.notes"><strong>Ghi chú:</strong> {{ item.notes }}</p>
          </div>
          <div class="facility-actions">
            <button @click="showHistory(item)" class="btn-history">📈 Lịch sử</button>
            <button @click="showMaintenance(item)" class="btn-maintenance">🔧 Bảo trì</button>
            <button @click="scheduleMaintenanceForItem(item)" class="btn-schedule-item">📅 Lên lịch</button>
            <button @click="showStatusUpdate(item)" class="btn-status">🔄 Trạng thái</button>
            <button @click="showStatusHistory(item)" class="btn-status-history">📅 Lịch sử TT</button>
            <button @click="reportIssue(item)" class="btn-report-issue">⚠️ Báo hỏng</button>
            <button @click="moveAsset(item)" class="btn-move">🚚 Di chuyển</button>
            <button @click="editItem(item)" class="btn-edit">📝 Sửa</button>
            <button @click="deleteItem(item.id)" class="btn-delete" :disabled="item.has_maintenance_history">🗑️ Xóa</button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <div v-if="showCreateForm || editingItem" class="modal">
        <div class="modal-content">
          <h3>{{ editingItem ? 'Sửa tài sản' : 'Thêm tài sản mới' }}</h3>
          <form @submit.prevent="saveItem">
            <div class="form-group">
              <label>Tên tài sản *</label>
              <input v-model="form.name" type="text" required />
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
                <option value="Thanh lý">Thanh lý</option>
              </select>
            </div>
            <div class="form-group">
              <label>Ngày mua</label>
              <input v-model="form.purchase_date" type="date" />
            </div>
            <div class="form-group">
              <label>Chi phí (VNĐ)</label>
              <input v-model.number="form.cost" type="number" min="0" />
            </div>
            <div class="form-group">
              <label>Nhà cung cấp</label>
              <input v-model="form.supplier" type="text" />
            </div>
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="form.notes" rows="3"></textarea>
            </div>
            <div class="form-group" v-if="editingItem">
              <label>Lý do thay đổi</label>
              <textarea v-model="form.change_reason" rows="2" placeholder="Nhập lý do cập nhật thông tin..."></textarea>
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
        <div class="modal-content maintenance-modal">
          <h3>🔧 Bảo trì: {{ maintenanceItem.name }}</h3>
          
          <div class="maintenance-tabs">
            <button @click="maintenanceTab = 'history'" :class="{active: maintenanceTab === 'history'}">Lịch sử</button>
            <button @click="maintenanceTab = 'schedule'" :class="{active: maintenanceTab === 'schedule'}">Lên lịch</button>
            <button @click="maintenanceTab = 'costs'" :class="{active: maintenanceTab === 'costs'}">Chi phí</button>
          </div>
          
          <!-- Maintenance History Tab -->
          <div v-if="maintenanceTab === 'history'" class="tab-content">
            <div class="maintenance-summary">
              <div class="summary-stats">
                <div class="stat-item">
                  <span class="stat-value">{{ maintenanceStats.total }}</span>
                  <span class="stat-label">Tổng số lần</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ formatPrice(maintenanceStats.totalCost) }}</span>
                  <span class="stat-label">Tổng chi phí</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ maintenanceStats.avgInterval }}</span>
                  <span class="stat-label">Chu kỳ TB (ngày)</span>
                </div>
              </div>
            </div>
            
            <div v-if="maintenanceRecords.length === 0" class="no-data">Chưa có lịch sử bảo trì</div>
            <div v-else class="maintenance-list">
              <div v-for="record in maintenanceRecords" :key="record.id" class="maintenance-item" 
                   :class="'type-' + record.type">
                <div class="maintenance-header">
                  <div class="maintenance-type">
                    <span class="type-badge" :class="'type-' + record.type">
                      {{ getMaintenanceTypeText(record.type) }}
                    </span>
                    <span class="maintenance-date">{{ formatDate(record.date) }}</span>
                  </div>
                  <div class="maintenance-cost">{{ formatPrice(record.cost) }}</div>
                </div>
                
                <div class="maintenance-content">
                  <p><strong>Mô tả:</strong> {{ record.description }}</p>
                  <p v-if="record.vendor"><strong>Đơn vị:</strong> {{ record.vendor }}</p>
                  <p v-if="record.technician"><strong>Kỹ thuật viên:</strong> {{ record.technician }}</p>
                  <p><strong>Người thực hiện:</strong> {{ record.performed_by }}</p>
                  <p v-if="record.duration"><strong>Thời gian:</strong> {{ record.duration }} giờ</p>
                  <p v-if="record.parts_used"><strong>Liệu kiện:</strong> {{ record.parts_used }}</p>
                </div>
                
                <div v-if="record.notes" class="maintenance-notes">
                  <strong>Ghi chú:</strong> {{ record.notes }}
                </div>
                
                <div class="maintenance-actions">
                  <button @click="editMaintenanceRecord(record)" class="btn-edit-maintenance">📝 Sửa</button>
                  <button @click="duplicateMaintenanceRecord(record)" class="btn-duplicate">📋 Sao chép</button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Schedule Tab -->
          <div v-if="maintenanceTab === 'schedule'" class="tab-content">
            <div class="next-maintenance-info">
              <h4>Bảo trì tiếp theo</h4>
              <p v-if="nextMaintenanceDate">
                Dự kiến: <strong>{{ formatDate(nextMaintenanceDate) }}</strong>
                (còn {{ getDaysUntil(nextMaintenanceDate) }} ngày)
              </p>
              <p v-else>Chưa có lịch bảo trì tiếp theo</p>
            </div>
            
            <button @click="showMaintenanceForm = true" class="btn-primary">+ Lên lịch bảo trì mới</button>
          </div>
          
          <!-- Costs Tab -->
          <div v-if="maintenanceTab === 'costs'" class="tab-content">
            <div class="cost-analysis">
              <div class="cost-chart-placeholder">
                <h4>Phân tích chi phí bảo trì</h4>
                <div class="cost-breakdown">
                  <div class="cost-item">
                    <span>Bảo trì định kỳ:</span>
                    <span>{{ formatPrice(maintenanceStats.scheduledCost) }}</span>
                  </div>
                  <div class="cost-item">
                    <span>Sửa chữa phát sinh:</span>
                    <span>{{ formatPrice(maintenanceStats.emergencyCost) }}</span>
                  </div>
                  <div class="cost-item total">
                    <span>Tổng cộng:</span>
                    <span>{{ formatPrice(maintenanceStats.totalCost) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <div class="form-actions">
            <button @click="closeMaintenance" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Enhanced Maintenance Form Modal -->
      <div v-if="showMaintenanceForm" class="modal">
        <div class="modal-content">
          <h3>{{ editingMaintenance ? 'Sửa bảo trì' : 'Thêm bảo trì mới' }}</h3>
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
              <textarea v-model="maintenanceForm.description" required rows="3" 
                       placeholder="Mô tả chi tiết công việc bảo trì..."></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label>Ngày thực hiện *</label>
                <input v-model="maintenanceForm.date" type="date" required />
              </div>
              <div class="form-group">
                <label>Thời gian (giờ)</label>
                <input v-model.number="maintenanceForm.duration" type="number" min="0.5" step="0.5" />
              </div>
            </div>
            
            <div class="form-group">
              <label>Đơn vị thực hiện</label>
              <select v-model="maintenanceForm.vendor_type">
                <option value="internal">Nội bộ</option>
                <option value="external">Bên ngoài</option>
              </select>
            </div>
            
            <div v-if="maintenanceForm.vendor_type === 'external'" class="form-group">
              <label>Tên đơn vị</label>
              <input v-model="maintenanceForm.vendor" type="text" placeholder="Tên công ty/thợ sửa chữa" />
            </div>
            
            <div v-if="maintenanceForm.vendor_type === 'internal'" class="form-group">
              <label>Nhân viên thực hiện</label>
              <input v-model="maintenanceForm.technician" type="text" placeholder="Tên nhân viên" />
            </div>
            
            <div class="form-group">
              <label>Chi phí (VNĐ)</label>
              <input v-model.number="maintenanceForm.cost" type="number" min="0" />
            </div>
            
            <div class="form-group">
              <label>Liệu kiện sử dụng</label>
              <textarea v-model="maintenanceForm.parts_used" rows="2" 
                       placeholder="Danh sách liệu kiện, phụ tùng đã thay thế..."></textarea>
            </div>
            
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="maintenanceForm.notes" rows="2" 
                       placeholder="Ghi chú thêm về quá trình bảo trì..."></textarea>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="closeMaintenanceForm" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">{{ editingMaintenance ? 'Cập nhật' : 'Lưu' }}</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Add Maintenance Modal -->
      <div v-if="showAddMaintenance" class="modal">
        <div class="modal-content">
          <h3>Thêm bảo trì</h3>
          <form @submit.prevent="saveMaintenance">
            <div class="form-group">
              <label>Loại *</label>
              <select v-model="maintenanceForm.type" required>
                <option value="scheduled">Định kỳ</option>
                <option value="emergency">Phát sinh</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mô tả *</label>
              <textarea v-model="maintenanceForm.description" required rows="3"></textarea>
            </div>
            <div class="form-group">
              <label>Chi phí (VNĐ)</label>
              <input v-model.number="maintenanceForm.cost" type="number" min="0" />
            </div>
            <div class="form-group">
              <label>Đơn vị sửa chữa</label>
              <input v-model="maintenanceForm.vendor" type="text" />
            </div>
            <div class="form-group">
              <label>Ngày thực hiện</label>
              <input v-model="maintenanceForm.date" type="date" />
            </div>
            <div class="form-actions">
              <button type="button" @click="showAddMaintenance = false" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Lưu</button>
            </div>
          </form>
        </div>
      </div>
      <!-- Maintenance Schedule Modal -->
      <div v-if="showMaintenanceSchedule" class="modal">
        <div class="modal-content schedule-modal">
          <h3>📅 Lịch Bảo trì</h3>
          <div class="schedule-tabs">
            <button @click="scheduleTab = 'scheduled'" :class="{active: scheduleTab === 'scheduled'}">Lịch hẹn</button>
            <button @click="scheduleTab = 'due'" :class="{active: scheduleTab === 'due'}">Đến hạn</button>
            <button @click="scheduleTab = 'overdue'" :class="{active: scheduleTab === 'overdue'}">Quá hạn</button>
          </div>
          
          <div v-if="scheduleTab === 'scheduled'" class="schedule-content">
            <div class="schedule-header">
              <h4>Các lịch bảo trì đã lên</h4>
              <button @click="showScheduleForm = true" class="btn-primary">+ Lên lịch mới</button>
            </div>
            <div v-if="scheduledTasks.length === 0" class="no-data">Chưa có lịch bảo trì</div>
            <div v-else class="task-list">
              <div v-for="task in scheduledTasks" :key="task.id" class="task-item">
                <div class="task-info">
                  <h5>{{ task.facility_name }}</h5>
                  <p><strong>Loại:</strong> {{ task.type === 'scheduled' ? 'Định kỳ' : 'Phát sinh' }}</p>
                  <p><strong>Mô tả:</strong> {{ task.description }}</p>
                  <p><strong>Ngày dự kiến:</strong> {{ formatDate(task.scheduled_date) }}</p>
                  <p><strong>Trạng thái:</strong> <span :class="'status-' + task.status">{{ getTaskStatusText(task.status) }}</span></p>
                </div>
                <div class="task-actions">
                  <button @click="completeTask(task)" class="btn-complete">✓ Hoàn thành</button>
                  <button @click="editTask(task)" class="btn-edit">📝 Sửa</button>
                </div>
              </div>
            </div>
          </div>
          
          <div v-if="scheduleTab === 'due'" class="schedule-content">
            <h4>Bảo trì đến hạn</h4>
            <div v-if="dueTasks.length === 0" class="no-data">Không có bảo trì đến hạn</div>
            <div v-else class="task-list">
              <div v-for="task in dueTasks" :key="task.id" class="task-item due">
                <div class="task-info">
                  <h5>{{ task.facility_name }}</h5>
                  <p><strong>Mô tả:</strong> {{ task.description }}</p>
                  <p><strong>Ngày đến hạn:</strong> {{ formatDate(task.scheduled_date) }}</p>
                </div>
                <div class="task-actions">
                  <button @click="completeTask(task)" class="btn-complete">✓ Hoàn thành</button>
                </div>
              </div>
            </div>
          </div>
          
          <div class="form-actions">
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
              <label>Loại bảo trì *</label>
              <select v-model="scheduleForm.type" required>
                <option value="scheduled">Định kỳ</option>
                <option value="preventive">Phòng ngừa</option>
                <option value="corrective">Sửa chữa</option>
              </select>
            </div>
            <div class="form-group">
              <label>Mô tả công việc *</label>
              <textarea v-model="scheduleForm.description" required rows="3" placeholder="Mô tả chi tiết công việc bảo trì..."></textarea>
            </div>
            <div class="form-group">
              <label>Ngày dự kiến *</label>
              <input v-model="scheduleForm.scheduled_date" type="date" required />
            </div>
            <div class="form-group">
              <label>Thời gian dự kiến (giờ)</label>
              <input v-model.number="scheduleForm.estimated_hours" type="number" min="0.5" step="0.5" />
            </div>
            <div class="form-group">
              <label>Chi phí dự kiến (VNĐ)</label>
              <input v-model.number="scheduleForm.estimated_cost" type="number" min="0" />
            </div>
            <div class="form-group">
              <label>Đơn vị thực hiện</label>
              <input v-model="scheduleForm.assigned_to" type="text" placeholder="Tên nhân viên hoặc đơn vị ngoài" />
            </div>
            <div class="form-group">
              <label>Ghi chú</label>
              <textarea v-model="scheduleForm.notes" rows="2"></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="showScheduleForm = false" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Lên lịch</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Move Asset Modal -->
      <div v-if="movingAsset" class="modal">
        <div class="modal-content">
          <h3>🚚 Di chuyển tài sản: {{ movingAsset.name }}</h3>
          <p>Khu vực hiện tại: <strong>{{ movingAsset.area }}</strong></p>
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
              <button type="button" @click="closeMoveAsset" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Di chuyển</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Issue Report Modal -->
      <div v-if="reportingIssue" class="modal">
        <div class="modal-content">
          <h3>⚠️ Báo cáo sự cố: {{ reportingIssue.name }}</h3>
          <form @submit.prevent="saveIssueReport">
            <div class="form-group">
              <label>Loại sự cố *</label>
              <select v-model="issueForm.issue_type" required>
                <option value="">Chọn loại sự cố</option>
                <option value="hư hỏng">Đồ bị hư hỏng</option>
                <option value="không hoạt động">Không hoạt động</option>
                <option value="hoạt động bất thường">Hoạt động bất thường</option>
                <option value="an toàn">Vấn đề an toàn</option>
                <option value="khác">Khác</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>Mức độ ảnh hưởng *</label>
              <select v-model="issueForm.severity" required>
                <option value="">Chọn mức độ</option>
                <option value="thấp">Thấp - Không ảnh hưởng hoạt động</option>
                <option value="trung bình">Trung bình - Ảnh hưởng một phần</option>
                <option value="cao">Cao - Ảnh hưởng nghiêm trọng</option>
                <option value="khẩn cấp">Khẩn cấp - Cần xử lý ngay</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>Mô tả sự cố *</label>
              <textarea v-model="issueForm.description" required rows="4" 
                       placeholder="Mô tả chi tiết sự cố, triệu chứng, thời điểm xảy ra..."></textarea>
            </div>
            
            <div class="form-group">
              <label>Vị trí cụ thể</label>
              <input v-model="issueForm.location" type="text" 
                     placeholder="Ví dụ: Góc phải quầy bar, cạnh cửa sổ..." />
            </div>
            
            <div class="form-group">
              <label>Hình ảnh miễn họa</label>
              <input type="file" @change="handleImageUpload" accept="image/*" multiple />
              <small>Có thể chọn nhiều hình ảnh</small>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="closeIssueReport" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Gửi báo cáo</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Issue Reports List Modal -->
      <div v-if="showIssueReportsModal" class="modal">
        <div class="modal-content reports-modal">
          <h3>📝 Báo cáo Sự cố</h3>
          
          <div v-if="issueReports.length === 0" class="no-data">Không có báo cáo sự cố</div>
          <div v-else class="issue-reports-list">
            <div v-for="report in issueReports" :key="report.id" class="issue-report-item">
              <div class="report-header">
                <h5>{{ report.facility_name }}</h5>
                <span class="severity-badge" :class="'severity-' + report.severity">
                  {{ getSeverityText(report.severity) }}
                </span>
              </div>
              
              <div class="report-content">
                <p><strong>Loại:</strong> {{ report.issue_type }}</p>
                <p><strong>Mô tả:</strong> {{ report.description }}</p>
                <p><strong>Người báo cáo:</strong> {{ report.reported_by }}</p>
                <p><strong>Thời gian:</strong> {{ formatDateTime(report.reported_at) }}</p>
              </div>
              
              <div class="report-actions">
                <button @click="updateReportStatus(report.id, 'resolved')" class="btn-resolve">Giải quyết</button>
              </div>
            </div>
          </div>
          
          <div class="form-actions">
            <button @click="showIssueReportsModal = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Status Update Modal -->
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
              <button type="button" @click="closeStatusUpdate" class="btn-cancel">Hủy</button>
              <button type="submit" class="btn-save">Cập nhật</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Asset Report Modal -->
      <div v-if="showAssetReport" class="modal">
        <div class="modal-content report-modal">
          <h3>📈 Báo cáo Tài sản</h3>
          <div class="report-summary">
            <div class="summary-card">
              <h4>Tổng quan</h4>
              <p>Tổng số tài sản: <strong>{{ assetReport.total }}</strong></p>
              <p>Đang sử dụng: <strong>{{ assetReport.active }}</strong></p>
              <p>Cần bảo trì: <strong>{{ assetReport.needMaintenance }}</strong></p>
              <p>Hỏng hóc: <strong>{{ assetReport.broken }}</strong></p>
            </div>
            <div class="summary-card">
              <h4>Theo khu vực</h4>
              <div v-for="(count, area) in assetReport.byArea" :key="area">
                {{ area }}: <strong>{{ count }}</strong>
              </div>
            </div>
            <div class="summary-card">
              <h4>Theo loại</h4>
              <div v-for="(count, type) in assetReport.byType" :key="type">
                {{ type }}: <strong>{{ count }}</strong>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button @click="exportAssetReport" class="btn-export">📎 Xuất Excel</button>
            <button @click="showAssetReport = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Status History Modal -->
      <div v-if="statusHistoryItem" class="modal">
        <div class="modal-content">
          <h3>📅 Lịch sử Trạng thái: {{ statusHistoryItem.name }}</h3>
          
          <div v-if="statusHistory.length === 0" class="no-data">Chưa có lịch sử thay đổi trạng thái</div>
          <div v-else class="status-timeline">
            <div v-for="history in statusHistory" :key="history.id" class="timeline-item">
              <div class="timeline-marker" :style="{backgroundColor: getStatusColor(history.new_value)}"></div>
              <div class="timeline-content">
                <div class="status-change">
                  <span class="old-status" :class="'status-' + getStatusClass(history.old_value)">{{ history.old_value }}</span>
                  <span class="arrow">→</span>
                  <span class="new-status" :class="'status-' + getStatusClass(history.new_value)">{{ history.new_value }}</span>
                </div>
                <div class="change-details">
                  <p><strong>Mô tả:</strong> {{ history.description }}</p>
                  <p><strong>Người thực hiện:</strong> {{ history.username }}</p>
                  <p><strong>Thời gian:</strong> {{ formatDateTime(history.created_at) }}</p>
                </div>
              </div>
            </div>
          </div>
          
          <div class="form-actions">
            <button @click="closeStatusHistory" class="btn-cancel">Đóng</button>
          </div>
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
const showAddMaintenance = ref(false)
const searchQuery = ref('')
const filterType = ref('')
const filterStatus = ref('')
const filterArea = ref('')
const statusUpdateItem = ref(null)
const showAssetReport = ref(false)
const assetReport = ref({})

const statusHistoryItem = ref(null)
const statusHistory = ref([])
const showStatusDashboardModal = ref(false)
const statusSummary = ref({})
const statusAlerts = ref([])

const statusForm = ref({
  status: '', notes: '', damage_level: '', expected_completion: '', disposal_reason: ''
})

const showMaintenanceSchedule = ref(false)
const showScheduleForm = ref(false)
const scheduleTab = ref('scheduled')
const scheduledTasks = ref([])
const dueTasks = ref([])

const movingAsset = ref(null)
const moveForm = ref({ new_area: '', reason: '' })
const reportingIssue = ref(null)
const showIssueReportsModal = ref(false)
const issueReports = ref([])
const scheduleForm = ref({
  facility_id: '', type: 'scheduled', description: '', scheduled_date: '',
  estimated_hours: 0, estimated_cost: 0, assigned_to: '', notes: ''
})

const issueForm = ref({
  issue_type: '', severity: '', description: '', location: '', 
  occurred_at: '', images: [], actions_taken: ''
})

const form = ref({
  name: '', type: '', area: '', quantity: 1, status: 'Đang sử dụng',
  purchase_date: '', cost: 0, supplier: '', notes: '', change_reason: ''
})

const maintenanceTab = ref('history')
const showMaintenanceForm = ref(false)
const editingMaintenance = ref(null)
const maintenanceStats = ref({})
const nextMaintenanceDate = ref(null)

const maintenanceForm = ref({
  type: 'scheduled', description: '', date: '', duration: 0, 
  vendor_type: 'internal', vendor: '', technician: '', 
  cost: 0, parts_used: '', notes: ''
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

// FR-FM-08: Search & Filter functionality
const searchWithFilters = async () => {
  const filters = {
    name: searchQuery.value || undefined,
    type: filterType.value || undefined,
    area: filterArea.value || undefined,
    status: filterStatus.value || undefined,
    limit: 50
  }
  
  const result = await facilityStore.searchFacilities(filters)
  // Update items with search results
  facilityStore.items = result.data || []
}

onMounted(async () => {
  await facilityStore.fetchFacilities()
  await loadMaintenanceData()
})

const loadMaintenanceData = async () => {
  try {
    scheduledTasks.value = await facilityStore.fetchScheduledMaintenance()
  } catch (error) {
    scheduledTasks.value = []
  }
  
  try {
    dueTasks.value = await facilityStore.fetchMaintenanceDue()
  } catch (error) {
    dueTasks.value = []
  }
}

const resetForm = () => {
  form.value = {
    name: '', type: '', area: '', quantity: 1, status: 'Đang sử dụng',
    purchase_date: '', cost: 0, supplier: '', notes: '', change_reason: ''
  }
}

const saveItem = async () => {
  const itemData = { ...form.value }
  if (editingItem.value && form.value.change_reason) {
    itemData.change_log = {
      reason: form.value.change_reason,
      timestamp: new Date().toISOString(),
      user: 'current_user' // TODO: get from auth store
    }
  }
  
  const success = editingItem.value 
    ? await facilityStore.updateFacility(editingItem.value.id, itemData)
    : await facilityStore.createFacility(itemData)
  
  if (success) {
    cancelEdit()
  }
}

const editItem = (item) => {
  editingItem.value = item
  form.value = { ...item }
  showCreateForm.value = false
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingItem.value = null
  resetForm()
}

const deleteItem = async (id) => {
  const item = items.value.find(i => i.id === id)
  if (item?.has_maintenance_history) {
    alert('Không thể xóa tài sản đã có lịch sử bảo trì. Chỉ có thể chuyển sang trạng thái "Ngừng sử dụng".')
    return
  }
  
  if (confirm('Bạn có chắc muốn xóa tài sản này?')) {
    await facilityStore.deleteFacility(id)
  }
}

const showStatusUpdate = (item) => {
  statusUpdateItem.value = item
  statusForm.value = { status: item.status, notes: '' }
}

const closeStatusUpdate = () => {
  statusUpdateItem.value = null
  statusForm.value = { status: '', notes: '', damage_level: '', expected_completion: '', disposal_reason: '' }
}

const saveStatusUpdate = async () => {
  const statusData = {
    status: statusForm.value.status,
    notes: statusForm.value.notes,
    damage_level: statusForm.value.damage_level,
    expected_completion: statusForm.value.expected_completion,
    disposal_reason: statusForm.value.disposal_reason
  }
  
  const success = await facilityStore.updateStatusWithDetails(statusUpdateItem.value.id, statusData)
  if (success) {
    closeStatusUpdate()
  }
}

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

const scheduleMaintenanceForItem = (item) => {
  scheduleForm.value.facility_id = item.id
  showScheduleForm.value = true
}

const saveScheduledTask = async () => {
  const success = await facilityStore.scheduleMaintenanceTask(scheduleForm.value)
  if (success) {
    showScheduleForm.value = false
    scheduleForm.value = {
      facility_id: '', type: 'scheduled', description: '', scheduled_date: '',
      estimated_hours: 0, estimated_cost: 0, assigned_to: '', notes: ''
    }
    await loadMaintenanceData()
  }
}

const completeTask = async (task) => {
  if (confirm('Đánh dấu nhiệm vụ này đã hoàn thành?')) {
    const success = await facilityStore.updateMaintenanceTask(task.id, { status: 'completed' })
    if (success) {
      await loadMaintenanceData()
    }
  }
}

const editTask = (task) => {
  scheduleForm.value = { ...task }
  showScheduleForm.value = true
}

const showMaintenanceDue = async () => {
  await loadMaintenanceData()
  scheduleTab.value = 'due'
  showMaintenanceSchedule.value = true
}

const moveAsset = (item) => {
  movingAsset.value = item
  moveForm.value = { new_area: '', reason: '' }
}

const closeMoveAsset = () => {
  movingAsset.value = null
  moveForm.value = { new_area: '', reason: '' }
}

const saveMoveAsset = async () => {
  const moveData = {
    area: moveForm.value.new_area,
    change_reason: `Di chuyển từ ${movingAsset.value.area} đến ${moveForm.value.new_area}: ${moveForm.value.reason}`
  }
  
  const success = await facilityStore.updateFacility(movingAsset.value.id, moveData)
  if (success) {
    closeMoveAsset()
  }
}
const getTaskStatusText = (status) => {
  const statusMap = {
    'pending': 'Chờ thực hiện',
    'in_progress': 'Đang thực hiện', 
    'completed': 'Hoàn thành',
    'cancelled': 'Đã hủy'
  }
  return statusMap[status] || status
}

const showStatusDashboard = async () => {
  statusSummary.value = calculateStatusSummary()
  statusAlerts.value = await facilityStore.fetchStatusAlerts()
  showStatusDashboardModal.value = true
}

const calculateStatusSummary = () => {
  const summary = {}
  items.value.forEach(item => {
    summary[item.status] = (summary[item.status] || 0) + 1
  })
  return summary
}

const getStatusPercentage = (count) => {
  return Math.round((count / items.value.length) * 100)
}

const getStatusIcon = (status) => {
  const icons = {
    'Đang sử dụng': '✅',
    'Hỏng': '❌',
    'Đang sửa': '🔧',
    'Ngừng sử dụng': '⏸️',
    'Thanh lý': '🗑️'
  }
  return icons[status] || '❓'
}

const showStatusHistory = async (item) => {
  console.log('showStatusHistory called with item:', item)
  statusHistoryItem.value = item
  try {
    const history = await facilityStore.fetchStatusHistory(item.id)
    console.log('Status history received:', history)
    statusHistory.value = history
  } catch (error) {
    console.error('Error fetching status history:', error)
    statusHistory.value = []
  }
}

const closeStatusHistory = () => {
  statusHistoryItem.value = null
  statusHistory.value = []
}

const onStatusChange = () => {
  statusForm.value.damage_level = ''
  statusForm.value.expected_completion = ''
  statusForm.value.disposal_reason = ''
}

const getNotesPlaceholder = (status) => {
  const placeholders = {
    'Hỏng': 'Mô tả chi tiết tình trạng hư hỏng...',
    'Đang sửa': 'Thông tin về quá trình sửa chữa...',
    'Ngừng sử dụng': 'Lý do ngừng sử dụng...',
    'Thanh lý': 'Chi tiết về việc thanh lý...'
  }
  return placeholders[status] || 'Lý do thay đổi trạng thái...'
}

const getStatusImpact = (status) => {
  const impacts = {
    'Hỏng': 'Tài sản sẽ không thể sử dụng cho đến khi được sửa chữa',
    'Đang sửa': 'Tài sản tạm thời không khả dụng',
    'Ngừng sử dụng': 'Tài sản sẽ được loại khỏi hoạt động',
    'Thanh lý': 'Tài sản sẽ bị xóa vĩnh viễn khỏi hệ thống'
  }
  return impacts[status]
}

const getStatusAge = (lastChange) => {
  if (!lastChange) return ''
  const days = Math.floor((new Date() - new Date(lastChange)) / (1000 * 60 * 60 * 24))
  if (days === 0) return 'Hôm nay'
  if (days === 1) return '1 ngày'
  return `${days} ngày`
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleString('vi-VN')
}

const reportIssue = (item) => {
  reportingIssue.value = item
  issueForm.value = {
    issue_type: '', severity: '', description: '', location: '', 
    occurred_at: '', images: [], actions_taken: ''
  }
}

const closeIssueReport = () => {
  reportingIssue.value = null
  issueForm.value = {
    issue_type: '', severity: '', description: '', location: '', 
    occurred_at: '', images: [], actions_taken: ''
  }
}

const saveIssueReport = async () => {
  const reportData = {
    ...issueForm.value,
    facility_id: reportingIssue.value.id,
    facility_name: reportingIssue.value.name
  }
  
  const success = await facilityStore.createIssueReport(reportData)
  if (success) {
    closeIssueReport()
    alert('Báo cáo sự cố đã được gửi thành công!')
  }
}

const handleImageUpload = (event) => {
  const files = Array.from(event.target.files)
  files.forEach(file => {
    if (file.type.startsWith('image/')) {
      const reader = new FileReader()
      reader.onload = (e) => {
        issueForm.value.images.push({
          file: file,
          preview: e.target.result,
          name: file.name
        })
      }
      reader.readAsDataURL(file)
    }
  })
}

const removeImage = (index) => {
  issueForm.value.images.splice(index, 1)
}

const showIssueReports = async () => {
  issueReports.value = await facilityStore.fetchIssueReports()
  showIssueReportsModal.value = true
}

const updateReportStatus = async (reportId, status) => {
  const success = await facilityStore.updateIssueReportStatus(reportId, status)
  if (success) {
    issueReports.value = await facilityStore.fetchIssueReports()
  }
}

const getSeverityText = (severity) => {
  const severityMap = {
    'thấp': 'Thấp',
    'trung bình': 'Trung bình',
    'cao': 'Cao',
    'khẩn cấp': 'Khẩn cấp'
  }
  return severityMap[severity] || severity
}

const getUrgencyText = (severity) => {
  const urgencyMap = {
    'thấp': 'Không cần giải quyết ngay',
    'trung bình': 'Nên giải quyết trong ngày',
    'cao': 'Cần giải quyết sớm',
    'khẩn cấp': 'Cần giải quyết ngay lập tức'
  }
  return urgencyMap[severity] || ''
}

const showMaintenance = async (item) => {
  maintenanceItem.value = item
  try {
    maintenanceRecords.value = await facilityStore.fetchMaintenanceHistory(item.id)
  } catch (error) {
    maintenanceRecords.value = []
  }
  
  try {
    maintenanceStats.value = await facilityStore.fetchMaintenanceStats(item.id)
  } catch (error) {
    maintenanceStats.value = { total: 0, totalCost: 0, avgInterval: 0, scheduledCost: 0, emergencyCost: 0 }
  }
  
  try {
    nextMaintenanceDate.value = await facilityStore.fetchNextMaintenanceDate(item.id)
  } catch (error) {
    nextMaintenanceDate.value = null
  }
  
  maintenanceTab.value = 'history'
}

const closeMaintenance = () => {
  maintenanceItem.value = null
  maintenanceRecords.value = []
  showMaintenanceForm.value = false
  editingMaintenance.value = null
}

const closeMaintenanceForm = () => {
  showMaintenanceForm.value = false
  editingMaintenance.value = null
  maintenanceForm.value = {
    type: 'scheduled', description: '', date: '', duration: 0, 
    vendor_type: 'internal', vendor: '', technician: '', 
    cost: 0, parts_used: '', notes: ''
  }
}

const saveMaintenanceRecord = async () => {
  const record = {
    ...maintenanceForm.value,
    facility_id: maintenanceItem.value.id
  }
  
  const success = editingMaintenance.value
    ? await facilityStore.updateMaintenanceRecord(editingMaintenance.value.id, record)
    : await facilityStore.createMaintenanceRecord(record)
    
  if (success) {
    closeMaintenanceForm()
    maintenanceRecords.value = await facilityStore.fetchMaintenanceHistory(maintenanceItem.value.id)
    maintenanceStats.value = await facilityStore.fetchMaintenanceStats(maintenanceItem.value.id)
  }
}

const editMaintenanceRecord = (record) => {
  editingMaintenance.value = record
  maintenanceForm.value = { ...record }
  showMaintenanceForm.value = true
}

const duplicateMaintenanceRecord = (record) => {
  maintenanceForm.value = {
    ...record,
    date: '',
    cost: record.cost || 0
  }
  delete maintenanceForm.value.id
  showMaintenanceForm.value = true
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

const getDaysUntil = (date) => {
  if (!date) return 0
  const days = Math.ceil((new Date(date) - new Date()) / (1000 * 60 * 60 * 24))
  return Math.max(0, days)
}

const getStatusClass = (status) => {
  const statusMap = {
    'Đang sử dụng': 'active',
    'Hỏng': 'broken',
    'Đang sửa': 'repair',
    'Ngừng sử dụng': 'inactive',
    'Thanh lý': 'disposed'
  }
  return statusMap[status] || 'default'
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('vi-VN')
}

const formatPrice = (price) => {
  if (!price) return '0 VNĐ'
  return new Intl.NumberFormat('vi-VN').format(price) + ' VNĐ'
}

const showHistory = async (item) => {
  const history = await facilityStore.fetchFacilityHistory(item.id)
  // TODO: Show history modal
  console.log('History for', item.name, history)
}

const getStatusColor = (status) => {
  const colorMap = {
    'Đang sử dụng': '#27ae60',
    'Hỏng': '#e74c3c',
    'Đang sửa': '#f39c12',
    'Ngừng sử dụng': '#95a5a6',
    'Thanh lý': '#7f8c8d'
  }
  return colorMap[status] || '#3498db'
}
</script>

<style scoped>
.facility-management {
  min-height: 100vh;
  background: #f5f5f5;
}

.content {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h2 {
  color: #2c3e50;
  margin: 0;
  font-size: 24px;
}

.header-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.filters-section {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-input, .filter-select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.facility-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.facility-card {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  transition: transform 0.2s;
}

.facility-card:hover {
  transform: translateY(-2px);
}

.facility-info h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 18px;
}

.facility-info p {
  margin: 8px 0;
  color: #666;
  font-size: 14px;
}

.status-active { color: #27ae60; font-weight: bold; }
.status-broken { color: #e74c3c; font-weight: bold; }
.status-repair { color: #f39c12; font-weight: bold; }
.status-inactive { color: #95a5a6; font-weight: bold; }
.status-disposed { color: #7f8c8d; font-weight: bold; }

.status-age {
  font-size: 12px;
  color: #999;
  font-weight: normal;
}

.facility-actions {
  display: flex;
  gap: 8px;
  margin-top: 15px;
  flex-wrap: wrap;
}

.facility-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background-color 0.2s;
  white-space: nowrap;
}

.btn-history { background: #3498db; color: white; }
.btn-maintenance { background: #f39c12; color: white; }
.btn-edit { background: #2ecc71; color: white; }
.btn-delete { background: #e74c3c; color: white; }
.btn-schedule { background: #17a2b8; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }
.btn-schedule-item { background: #6c757d; color: white; }
.btn-due { background: #dc3545; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }
.btn-complete { background: #28a745; color: white; }
.btn-status { background: #6f42c1; color: white; }
.btn-status-history { background: #20c997; color: white; }
.btn-report-issue { background: #dc3545; color: white; }
.btn-move { background: #fd7e14; color: white; }
.btn-issues { background: #fd7e14; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }
.btn-status-dashboard { background: #6f42c1; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }
.btn-report { background: #17a2b8; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }
.btn-search { background: #007bff; color: white; padding: 10px 15px; border: none; border-radius: 5px; cursor: pointer; }

.btn-primary {
  background: #3498db;
  color: white;
  padding: 12px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn-primary:hover { background: #2980b9; }

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
  box-sizing: border-box;
}

.modal-content {
  background: white;
  padding: 30px;
  border-radius: 10px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-content h3 {
  margin: 0 0 20px 0;
  color: #2c3e50;
  font-size: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
  font-size: 14px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.btn-save {
  background: #27ae60;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-cancel {
  background: #95a5a6;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

/* Maintenance Modal Styles */
.maintenance-modal { max-width: 900px; }
.maintenance-tabs { display: flex; gap: 5px; margin-bottom: 20px; }
.maintenance-tabs button { padding: 8px 16px; border: 1px solid #ddd; background: #f8f9fa; cursor: pointer; border-radius: 4px 4px 0 0; }
.maintenance-tabs button.active { background: #007bff; color: white; border-color: #007bff; }

.maintenance-summary { margin-bottom: 20px; }
.summary-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; }
.stat-item { text-align: center; padding: 15px; background: #f8f9fa; border-radius: 8px; }
.stat-value { display: block; font-size: 24px; font-weight: bold; color: #007bff; }
.stat-label { font-size: 12px; color: #666; }

.maintenance-list { max-height: 400px; overflow-y: auto; }
.maintenance-item { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 15px; }
.maintenance-item.type-scheduled { border-left: 4px solid #28a745; }
.maintenance-item.type-preventive { border-left: 4px solid #17a2b8; }
.maintenance-item.type-corrective { border-left: 4px solid #ffc107; }
.maintenance-item.type-emergency { border-left: 4px solid #dc3545; }

.maintenance-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.maintenance-type { display: flex; align-items: center; gap: 10px; }
.type-badge { padding: 3px 8px; border-radius: 12px; font-size: 11px; font-weight: bold; }
.type-badge.type-scheduled { background: #d4edda; color: #155724; }
.type-badge.type-preventive { background: #d1ecf1; color: #0c5460; }
.type-badge.type-corrective { background: #fff3cd; color: #856404; }
.type-badge.type-emergency { background: #f8d7da; color: #721c24; }

.maintenance-date { font-size: 12px; color: #666; }
.maintenance-cost { font-weight: bold; color: #007bff; }
.maintenance-content p { margin: 5px 0; font-size: 14px; color: #666; }
.maintenance-notes { background: #e9ecef; padding: 10px; border-radius: 4px; margin-top: 10px; font-style: italic; }
.maintenance-actions { display: flex; gap: 8px; margin-top: 10px; }
.btn-edit-maintenance, .btn-duplicate { padding: 4px 8px; border: none; border-radius: 3px; cursor: pointer; font-size: 11px; }
.btn-edit-maintenance { background: #ffc107; color: #212529; }
.btn-duplicate { background: #6c757d; color: white; }

/* Schedule Modal Styles */
.schedule-modal { max-width: 900px; }
.schedule-tabs { display: flex; gap: 5px; margin-bottom: 20px; }
.schedule-tabs button { padding: 8px 16px; border: 1px solid #ddd; background: #f8f9fa; cursor: pointer; }
.schedule-tabs button.active { background: #007bff; color: white; border-color: #007bff; }

.schedule-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
.task-list { max-height: 400px; overflow-y: auto; }
.task-item { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }
.task-item.due { border-left: 4px solid #dc3545; }
.task-info h5 { margin: 0 0 8px 0; color: #333; }
.task-info p { margin: 3px 0; font-size: 14px; color: #666; }
.task-actions { display: flex; gap: 8px; }
.task-actions button { padding: 6px 12px; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; }

/* Status Timeline Styles */
.status-timeline { max-height: 400px; overflow-y: auto; }
.timeline-item { display: flex; gap: 15px; margin-bottom: 20px; }
.timeline-marker { width: 12px; height: 12px; border-radius: 50%; margin-top: 5px; flex-shrink: 0; }
.timeline-content { flex: 1; }
.status-change { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.old-status, .new-status { padding: 3px 8px; border-radius: 12px; font-size: 12px; font-weight: bold; }
.arrow { color: #666; }
.change-details p { margin: 3px 0; font-size: 14px; color: #666; }

/* Report Modal Styles */
.reports-modal { max-width: 800px; }
.issue-reports-list { max-height: 500px; overflow-y: auto; }
.issue-report-item { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 15px; border-left: 4px solid #dc3545; }
.report-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.report-header h5 { margin: 0; color: #333; }
.severity-badge { padding: 3px 8px; border-radius: 12px; font-size: 11px; font-weight: bold; }
.severity-thấp { background: #d4edda; color: #155724; }
.severity-trung-bình { background: #fff3cd; color: #856404; }
.severity-cao { background: #f8d7da; color: #721c24; }
.severity-khẩn-cấp { background: #f5c6cb; color: #721c24; animation: pulse 2s infinite; }
.report-content p { margin: 5px 0; font-size: 14px; color: #666; }
.report-actions { display: flex; gap: 8px; margin-top: 10px; }
.report-actions button { padding: 6px 12px; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; }
.btn-resolve { background: #28a745; color: white; }

/* Form Row Styles */
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 15px; }

/* Cost Analysis Styles */
.cost-analysis { background: #f8f9fa; padding: 20px; border-radius: 8px; }
.cost-breakdown { margin-top: 15px; }
.cost-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #dee2e6; }
.cost-item.total { font-weight: bold; border-top: 2px solid #007bff; border-bottom: none; }

/* Next Maintenance Info */
.next-maintenance-info { background: #e7f3ff; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
.next-maintenance-info h4 { margin: 0 0 10px 0; color: #0066cc; }

/* Report Summary Styles */
.report-modal { max-width: 800px; }
.report-summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin-bottom: 20px; }
.summary-card { background: #f8f9fa; padding: 15px; border-radius: 8px; border-left: 4px solid #3498db; }
.summary-card h4 { margin: 0 0 10px 0; color: #2c3e50; }
.summary-card p, .summary-card div { margin: 5px 0; color: #666; }

.loading, .error, .no-data {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 16px;
}

.error {
  color: #e74c3c;
  background: #fdf2f2;
  border: 1px solid #f5c6cb;
  border-radius: 5px;
}

/* Mobile Responsive Styles */
@media (max-width: 768px) {
  .content {
    padding: 15px;
  }
  
  .header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
    margin-bottom: 20px;
  }
  
  .header h2 {
    font-size: 20px;
  }
  
  .header-actions {
    justify-content: center;
    width: 100%;
  }
  
  .header-actions button {
    flex: 1;
    min-width: 120px;
    padding: 12px 8px;
    font-size: 12px;
  }
  
  .filters-section {
    flex-direction: column;
    gap: 10px;
  }
  
  .search-input, .filter-select {
    width: 100%;
    min-width: auto;
  }
  
  .facility-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .facility-card {
    padding: 15px;
  }
  
  .facility-info h4 {
    font-size: 16px;
    margin-bottom: 12px;
  }
  
  .facility-info p {
    font-size: 13px;
    margin: 6px 0;
  }
  
  .facility-actions {
    justify-content: center;
    gap: 6px;
  }
  
  .facility-actions button {
    padding: 8px 10px;
    font-size: 11px;
    flex: 1;
    min-width: 60px;
  }
  
  .modal {
    padding: 10px;
  }
  
  .modal-content {
    padding: 20px;
    max-width: 100%;
    width: 100%;
    max-height: 90vh;
  }
  
  .modal-content h3 {
    font-size: 18px;
    text-align: center;
  }
  
  .form-row {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .form-actions {
    flex-direction: column;
    gap: 10px;
  }
  
  .form-actions button {
    width: 100%;
    padding: 12px;
  }
  
  /* Maintenance Modal Mobile */
  .maintenance-modal {
    max-width: 100%;
  }
  
  .maintenance-tabs {
    flex-wrap: wrap;
    gap: 3px;
  }
  
  .maintenance-tabs button {
    flex: 1;
    min-width: 80px;
    padding: 8px 4px;
    font-size: 12px;
  }
  
  .summary-stats {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .stat-item {
    padding: 12px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .maintenance-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .maintenance-type {
    flex-wrap: wrap;
    gap: 5px;
  }
  
  .maintenance-actions {
    justify-content: center;
  }
  
  /* Schedule Modal Mobile */
  .schedule-modal {
    max-width: 100%;
  }
  
  .schedule-header {
    flex-direction: column;
    gap: 10px;
    align-items: stretch;
  }
  
  .schedule-header h4 {
    text-align: center;
    margin: 0;
  }
  
  .task-item {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  
  .task-actions {
    justify-content: center;
  }
  
  /* Timeline Mobile */
  .timeline-item {
    gap: 10px;
  }
  
  .status-change {
    flex-wrap: wrap;
    gap: 5px;
  }
  
  /* Reports Modal Mobile */
  .reports-modal {
    max-width: 100%;
  }
  
  .report-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .report-actions {
    justify-content: center;
  }
  
  /* Report Summary Mobile */
  .report-summary {
    grid-template-columns: 1fr;
    gap: 15px;
  }
}

@media (max-width: 480px) {
  .content {
    padding: 10px;
  }
  
  .header h2 {
    font-size: 18px;
  }
  
  .header-actions button {
    padding: 10px 6px;
    font-size: 11px;
  }
  
  .facility-card {
    padding: 12px;
  }
  
  .facility-info h4 {
    font-size: 15px;
  }
  
  .facility-info p {
    font-size: 12px;
  }
  
  .facility-actions button {
    padding: 6px 8px;
    font-size: 10px;
  }
  
  .modal-content {
    padding: 15px;
  }
  
  .modal-content h3 {
    font-size: 16px;
  }
  
  .form-group label {
    font-size: 13px;
  }
  
  .form-group input,
  .form-group select,
  .form-group textarea {
    font-size: 13px;
    padding: 8px;
  }
  
  .maintenance-tabs button {
    padding: 6px 2px;
    font-size: 11px;
  }
  
  .stat-value {
    font-size: 18px;
  }
  
  .stat-label {
    font-size: 11px;
  }
}
</style>

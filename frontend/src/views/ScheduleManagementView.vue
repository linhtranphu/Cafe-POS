<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">

    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">📅 Quản lý lịch ca</h1>

        <!-- Tabs -->
        <div class="flex gap-2">
          <button
            @click="tab = 'periods'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              tab === 'periods'
                ? 'bg-blue-500 text-white shadow-lg'
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📆</span>
              <span>Kỳ đăng ký</span>
            </div>
          </button>
          <button
            @click="tab = 'templates'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              tab === 'templates'
                ? 'bg-blue-500 text-white shadow-lg'
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>🗂️</span>
              <span>Ca mẫu</span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-3 pb-24">

      <!-- ── Tab: Kỳ đăng ký ── -->
      <div v-if="tab === 'periods'">
        <button
          @click="openPeriodForm()"
          class="w-full mb-4 py-3 rounded-2xl bg-blue-500 text-white font-bold text-base shadow-sm active:bg-blue-600 transition-all touch-manipulation">
          + Tạo kỳ đăng ký mới
        </button>

        <div v-if="loadingPeriods" class="flex justify-center py-12">
          <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <div v-else-if="periods.length === 0" class="text-center py-12 text-gray-400">
          <div class="text-4xl mb-2">📭</div>
          <p>Chưa có kỳ đăng ký nào.</p>
        </div>

        <div v-for="p in periods" :key="p.id" class="mb-4">
          <!-- Period card -->
          <div :class="[
            'bg-white rounded-2xl shadow-sm overflow-hidden border-2',
            p.status === 'open' ? 'border-green-300' :
            p.status === 'published' ? 'border-blue-300' :
            p.status === 'closed' ? 'border-orange-200' :
            'border-gray-100'
          ]">
            <!-- Header -->
            <div :class="[
              'px-4 py-3',
              p.status === 'open' ? 'bg-green-50' :
              p.status === 'published' ? 'bg-blue-50' :
              p.status === 'closed' ? 'bg-orange-50' :
              p.status === 'cancelled' ? 'bg-red-50' :
              'bg-gray-50'
            ]">
              <div class="flex justify-between items-start">
                <div>
                  <div class="font-bold text-gray-900">{{ p.name }}</div>
                  <div class="text-xs text-gray-500 mt-0.5">
                    {{ fmtDate(p.start_date) }} → {{ fmtDate(p.end_date) }}
                  </div>
                  <div class="text-xs text-gray-400 mt-0.5">
                    Hạn ĐK: {{ fmtDate(p.register_deadline) }}
                  </div>
                </div>
                <span :class="[
                  'px-2.5 py-1 rounded-full text-xs font-bold',
                  p.status === 'open' ? 'bg-green-500 text-white' :
                  p.status === 'published' ? 'bg-blue-500 text-white' :
                  p.status === 'closed' ? 'bg-orange-400 text-white' :
                  p.status === 'cancelled' ? 'bg-red-500 text-white' :
                  'bg-gray-300 text-gray-700'
                ]">{{ statusLabel(p.status) }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="px-4 py-3 flex flex-wrap gap-2 border-t border-gray-100">
              <button v-if="p.status === 'draft'"
                @click="changeStatus(p, 'open')"
                class="flex-1 py-2 rounded-xl bg-green-500 text-white text-sm font-semibold active:bg-green-600 touch-manipulation">
                Mở đăng ký
              </button>
              <button v-if="p.status === 'open'"
                @click="changeStatus(p, 'closed')"
                class="flex-1 py-2 rounded-xl bg-orange-400 text-white text-sm font-semibold active:bg-orange-500 touch-manipulation">
                Đóng đăng ký
              </button>
              <button v-if="p.status === 'closed'"
                @click="changeStatus(p, 'published')"
                class="flex-1 py-2 rounded-xl bg-blue-500 text-white text-sm font-semibold active:bg-blue-600 touch-manipulation">
                Công bố lịch
              </button>
              <button
                v-if="p.status !== 'cancelled'"
                @click="selectPeriod(p)"
                class="flex-1 py-2 rounded-xl bg-gray-100 text-gray-700 text-sm font-semibold active:bg-gray-200 touch-manipulation">
                {{ selectedPeriod?.id === p.id ? 'Ẩn ca' : 'Xem & thêm ca' }}
              </button>
              <button
                v-if="p.status !== 'cancelled'"
                @click="changeStatus(p, 'cancelled')"
                class="py-2 px-4 rounded-xl border border-red-300 text-red-500 text-sm font-semibold active:bg-red-50 touch-manipulation">
                Huỷ lịch
              </button>
            </div>

            <!-- Slot management (expanded) -->
            <div v-if="selectedPeriod?.id === p.id" class="border-t border-gray-100">
              <!-- Add slot form (only for draft/open) -->
              <div v-if="p.status === 'draft' || p.status === 'open'" class="px-4 py-3 bg-gray-50 border-b border-gray-100">
                <div class="flex items-center justify-between mb-2">
                  <div class="text-sm font-semibold text-gray-700">Thêm ca vào kỳ</div>
                  <div class="flex gap-2 text-xs">
                    <button
                      @click="slotAddMode = 'single'"
                      :class="slotAddMode === 'single' ? 'bg-blue-500 text-white' : 'bg-white text-gray-600 border border-gray-200'"
                      class="px-3 py-1 rounded-lg font-semibold touch-manipulation">
                      Một ngày
                    </button>
                    <button
                      @click="openWeeklyModal(p)"
                      :class="slotAddMode === 'weekly' ? 'bg-blue-500 text-white' : 'bg-white text-gray-600 border border-gray-200'"
                      class="px-3 py-1 rounded-lg font-semibold touch-manipulation">
                      📅 Cả tuần
                    </button>
                  </div>
                </div>
                <select v-model="newSlot.templateId"
                  class="w-full mb-2 px-3 py-2.5 border border-gray-200 rounded-xl text-sm bg-white">
                  <option value="">Chọn ca mẫu...</option>
                  <option v-for="t in templates" :key="t.id" :value="t.id">
                    {{ t.name }} ({{ t.start_time }}–{{ t.end_time }})
                  </option>
                </select>
                <input type="date" v-model="newSlot.date"
                  class="w-full mb-2 px-3 py-2.5 border border-gray-200 rounded-xl text-sm bg-white" />
                <button
                  @click="addSlot(p)"
                  :disabled="!newSlot.templateId || !newSlot.date"
                  class="w-full py-2.5 rounded-xl bg-blue-500 text-white text-sm font-semibold active:bg-blue-600 touch-manipulation disabled:opacity-40">
                  + Thêm ca
                </button>
              </div>

              <!-- Slot list -->
              <div v-if="loadingDetail" class="flex justify-center py-6">
                <div class="w-6 h-6 border-3 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
              </div>
              <div v-else-if="!periodDetail?.slots?.length" class="text-center py-6 text-gray-400 text-sm">
                Chưa có ca nào trong kỳ này.
              </div>
              <div v-else class="divide-y divide-gray-100">
                <div v-for="sl in periodDetail.slots" :key="sl.id" class="px-4 py-3">
                  <div class="flex justify-between items-start">
                    <div>
                      <div class="font-semibold text-sm text-gray-900">{{ sl.template_name }}</div>
                      <div class="text-xs text-gray-500 mt-0.5">
                        {{ fmtDateShort(sl.date) }} · {{ sl.start_time }}–{{ sl.end_time }}
                      </div>
                      <div class="flex flex-wrap gap-1 mt-1.5">
                        <span
                          v-for="r in sl.registrations" :key="r.id"
                          class="inline-flex items-center gap-1 bg-blue-100 text-blue-700 text-xs px-2 py-0.5 rounded-full font-medium"
                        >
                          {{ r.user_name }}
                          <button
                            @click="cancelReg(p, sl, r)"
                            class="text-blue-400 hover:text-red-500 font-bold leading-none touch-manipulation"
                            title="Xóa đăng ký"
                          >×</button>
                        </span>
                        <span v-if="!sl.registrations?.length" class="text-xs text-gray-400 italic">Chưa có ai</span>
                      </div>
                    </div>
                    <div class="flex flex-col items-end gap-1.5 ml-2 flex-shrink-0">
                      <span :class="[
                        'text-xs font-bold px-2 py-0.5 rounded-full',
                        sl.registered_count >= sl.max_slots ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-700'
                      ]">
                        {{ sl.registered_count }}/{{ sl.max_slots }}
                      </span>
                      <button
                        v-if="p.status === 'draft' || p.status === 'open'"
                        @click="deleteSlot(p, sl)"
                        class="text-xs text-red-400 underline touch-manipulation">
                        Xóa
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Tab: Ca mẫu ── -->
      <div v-if="tab === 'templates'">
        <button
          @click="openTemplateForm()"
          class="w-full mb-4 py-3 rounded-2xl bg-blue-500 text-white font-bold text-base shadow-sm active:bg-blue-600 transition-all touch-manipulation">
          + Thêm ca mẫu mới
        </button>

        <div v-if="loadingTemplates" class="flex justify-center py-12">
          <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <div v-else-if="templates.length === 0" class="text-center py-12 text-gray-400">
          <div class="text-4xl mb-2">📭</div>
          <p>Chưa có ca mẫu nào.</p>
        </div>

        <div class="flex flex-col gap-3">
          <div v-for="t in templates" :key="t.id"
            class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100">
            <div class="flex justify-between items-start">
              <div>
                <div class="font-bold text-gray-900">{{ t.name }}</div>
                <div class="text-sm text-gray-500 mt-0.5">🕐 {{ t.start_time }} – {{ t.end_time }}</div>
                <div class="text-xs text-gray-400 mt-0.5">
                  Tối đa {{ t.max_slots }} người ·
                  {{ (t.roles || []).includes('all') ? 'Tất cả vai trò' : (t.roles || []).join(', ') }}
                </div>
              </div>
              <div class="flex flex-col items-end gap-2 ml-2">
                <span :class="t.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'"
                  class="text-xs font-semibold px-2 py-0.5 rounded-full">
                  {{ t.is_active ? 'Hoạt động' : 'Tắt' }}
                </span>
                <div class="flex gap-2">
                  <button @click="openTemplateForm(t)"
                    class="text-xs text-blue-500 underline touch-manipulation">Sửa</button>
                  <button @click="confirmDeleteTemplate(t)"
                    class="text-xs text-red-400 underline touch-manipulation">Xóa</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Template Modal ── -->
    <div v-if="showTemplateModal"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showTemplateModal = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[85vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-5"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-5">
          {{ editingTemplate ? 'Sửa ca mẫu' : 'Thêm ca mẫu' }}
        </h3>

        <div class="flex flex-col gap-4">
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Tên ca</label>
            <input v-model="templateForm.name" placeholder="VD: Ca sáng"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Giờ bắt đầu</label>
              <input type="time" v-model="templateForm.start_time"
                class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
            </div>
            <div class="flex-1">
              <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Giờ kết thúc</label>
              <input type="time" v-model="templateForm.end_time"
                class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
            </div>
          </div>
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Số người tối đa</label>
            <input type="number" v-model.number="templateForm.max_slots" min="1"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Vai trò được đăng ký</label>
            <select v-model="templateForm.roleMode"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base bg-white">
              <option value="all">Tất cả vai trò</option>
              <option value="specific">Chỉ định vai trò</option>
            </select>
          </div>
          <div v-if="templateForm.roleMode === 'specific'" class="flex gap-4 flex-wrap">
            <label v-for="r in allRoles" :key="r" class="flex items-center gap-2 text-sm">
              <input type="checkbox" :value="r" v-model="templateForm.roles"
                class="w-4 h-4 rounded" />
              {{ r }}
            </label>
          </div>
          <label class="flex items-center gap-3 text-sm font-semibold text-gray-700">
            <input type="checkbox" v-model="templateForm.is_active" class="w-4 h-4 rounded" />
            Đang hoạt động
          </label>
        </div>

        <p v-if="templateError" class="mt-3 text-red-500 text-sm">{{ templateError }}</p>

        <div class="flex gap-3 mt-6">
          <button @click="showTemplateModal = false"
            class="flex-1 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
            Hủy
          </button>
          <button @click="saveTemplate" :disabled="savingTemplate"
            class="flex-1 py-3 rounded-xl bg-blue-500 text-white font-semibold active:bg-blue-600 touch-manipulation disabled:opacity-50">
            {{ savingTemplate ? 'Đang lưu...' : 'Lưu' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ── Period Modal ── -->
    <div v-if="showPeriodModal"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showPeriodModal = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[85vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-5"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-5">Tạo kỳ đăng ký mới</h3>

        <div class="flex flex-col gap-4">
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Tên kỳ</label>
            <input v-model="periodForm.name" placeholder="VD: Tuần 1 tháng 4"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Từ ngày</label>
              <input type="date" v-model="periodForm.start_date"
                class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
            </div>
            <div class="flex-1">
              <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Đến ngày</label>
              <input type="date" v-model="periodForm.end_date"
                class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
            </div>
          </div>
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Hạn đăng ký</label>
            <input type="datetime-local" v-model="periodForm.register_deadline"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
        </div>

        <p v-if="periodError" class="mt-3 text-red-500 text-sm">{{ periodError }}</p>

        <div class="flex gap-3 mt-6">
          <button @click="showPeriodModal = false"
            class="flex-1 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
            Hủy
          </button>
          <button @click="savePeriod" :disabled="savingPeriod"
            class="flex-1 py-3 rounded-xl bg-blue-500 text-white font-semibold active:bg-blue-600 touch-manipulation disabled:opacity-50">
            {{ savingPeriod ? 'Đang tạo...' : 'Tạo kỳ' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ── Weekly Slot Modal ── -->
    <div v-if="showWeeklyModal"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showWeeklyModal = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[90vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-5"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-1">📅 Thêm ca theo tuần</h3>
        <p class="text-xs text-gray-400 mb-4">Gán ca mẫu cho từng ngày trong tuần, rồi chọn các tuần cần áp dụng.</p>

        <!-- Day-template assignment -->
        <div class="mb-4">
          <div class="text-sm font-semibold text-gray-700 mb-2">Ca mẫu theo ngày</div>
          <div class="flex flex-col gap-2">
            <div v-for="day in weekDays" :key="day.key" class="flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 w-8 shrink-0">{{ day.label }}</span>
              <select v-model="weekPattern[day.key]"
                class="flex-1 px-3 py-2 border border-gray-200 rounded-xl text-sm bg-white">
                <option value="">Nghỉ</option>
                <option v-for="t in templates" :key="t.id" :value="t.id">
                  {{ t.name }} ({{ t.start_time }}–{{ t.end_time }})
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- Week range selection -->
        <div class="mb-4">
          <div class="text-sm font-semibold text-gray-700 mb-2">Áp dụng cho tuần</div>
          <div class="flex gap-2 mb-2">
            <div class="flex-1">
              <label class="text-xs text-gray-500 mb-1 block">Từ tuần (ngày T2)</label>
              <input type="date" v-model="weekRange.start"
                class="w-full px-3 py-2 border border-gray-200 rounded-xl text-sm bg-white" />
            </div>
            <div class="flex-1">
              <label class="text-xs text-gray-500 mb-1 block">Đến tuần (ngày T2)</label>
              <input type="date" v-model="weekRange.end"
                class="w-full px-3 py-2 border border-gray-200 rounded-xl text-sm bg-white" />
            </div>
          </div>
          <!-- Preview -->
          <div v-if="weekPreview.length" class="bg-gray-50 rounded-xl p-3 max-h-40 overflow-y-auto">
            <div class="text-xs text-gray-500 mb-1.5 font-semibold">Xem trước ({{ weekPreview.length }} ca):</div>
            <div v-for="item in weekPreview" :key="item.date + item.templateId"
              class="text-xs text-gray-700 py-0.5">
              {{ item.dateLabel }} — {{ item.templateName }}
            </div>
          </div>
          <div v-else-if="weekRange.start" class="text-xs text-gray-400 italic mt-1">
            Không có ca nào được gán cho khoảng tuần này.
          </div>
        </div>

        <p v-if="weeklyError" class="text-red-500 text-sm mb-3">{{ weeklyError }}</p>

        <div class="flex gap-3">
          <button @click="showWeeklyModal = false"
            class="flex-1 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
            Hủy
          </button>
          <button @click="applyWeeklyPattern"
            :disabled="!weekPreview.length || weeklyAdding"
            class="flex-1 py-3 rounded-xl bg-blue-500 text-white font-semibold active:bg-blue-600 touch-manipulation disabled:opacity-50">
            {{ weeklyAdding ? `Đang tạo... (${weeklyProgress}/${weekPreview.length})` : `Tạo ${weekPreview.length} ca` }}
          </button>
        </div>
      </div>
    </div>

    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import * as scheduleApi from '../services/schedule.js'

const tab = ref('periods')
const allRoles = ['waiter', 'barista', 'cashier']
const slotAddMode = ref('single')

// ── Templates ──────────────────────────────────────────
const templates = ref([])
const loadingTemplates = ref(false)
const showTemplateModal = ref(false)
const editingTemplate = ref(null)
const savingTemplate = ref(false)
const templateError = ref('')

const defaultTemplateForm = () => ({
  name: '', start_time: '07:00', end_time: '15:00',
  max_slots: 5, is_active: true, roleMode: 'all', roles: [],
})
const templateForm = ref(defaultTemplateForm())

async function loadTemplates() {
  loadingTemplates.value = true
  try { templates.value = await scheduleApi.getTemplates() || [] }
  finally { loadingTemplates.value = false }
}

function openTemplateForm(t = null) {
  editingTemplate.value = t
  templateError.value = ''
  if (t) {
    templateForm.value = {
      name: t.name, start_time: t.start_time, end_time: t.end_time,
      max_slots: t.max_slots, is_active: t.is_active,
      roleMode: (t.roles?.includes('all') || !t.roles?.length) ? 'all' : 'specific',
      roles: (t.roles || []).filter(r => r !== 'all'),
    }
  } else {
    templateForm.value = defaultTemplateForm()
  }
  showTemplateModal.value = true
}

async function saveTemplate() {
  templateError.value = ''
  savingTemplate.value = true
  try {
    const roles = templateForm.value.roleMode === 'all' ? ['all'] : templateForm.value.roles
    const payload = {
      name: templateForm.value.name, start_time: templateForm.value.start_time,
      end_time: templateForm.value.end_time, max_slots: templateForm.value.max_slots,
      is_active: templateForm.value.is_active, roles,
    }
    if (editingTemplate.value) {
      await scheduleApi.updateTemplate(editingTemplate.value.id, payload)
    } else {
      await scheduleApi.createTemplate(payload)
    }
    showTemplateModal.value = false
    await loadTemplates()
  } catch (e) {
    templateError.value = e.response?.data?.error || e.message
  } finally {
    savingTemplate.value = false
  }
}

async function confirmDeleteTemplate(t) {
  if (!confirm(`Xóa ca mẫu "${t.name}"?`)) return
  await scheduleApi.deleteTemplate(t.id)
  await loadTemplates()
}

// ── Periods ────────────────────────────────────────────
const periods = ref([])
const loadingPeriods = ref(false)
const showPeriodModal = ref(false)
const savingPeriod = ref(false)
const periodError = ref('')
const periodForm = ref({ name: '', start_date: '', end_date: '', register_deadline: '' })

const selectedPeriod = ref(null)
const periodDetail = ref(null)
const loadingDetail = ref(false)
const newSlot = ref({ templateId: '', date: '' })

async function loadPeriods() {
  loadingPeriods.value = true
  try { periods.value = await scheduleApi.getPeriods() || [] }
  finally { loadingPeriods.value = false }
}

function openPeriodForm() {
  periodError.value = ''
  periodForm.value = { name: '', start_date: '', end_date: '', register_deadline: '' }
  showPeriodModal.value = true
}

async function savePeriod() {
  periodError.value = ''
  const f = periodForm.value
  if (!f.name.trim()) { periodError.value = 'Tên kỳ không được để trống'; return }
  if (!f.start_date) { periodError.value = 'Vui lòng chọn ngày bắt đầu'; return }
  if (!f.end_date) { periodError.value = 'Vui lòng chọn ngày kết thúc'; return }
  if (f.end_date < f.start_date) { periodError.value = 'Ngày kết thúc phải sau ngày bắt đầu'; return }

  savingPeriod.value = true
  try {
    const payload = {
      name: f.name.trim(),
      start_date: new Date(f.start_date + 'T00:00:00Z').toISOString(),
      end_date: new Date(f.end_date + 'T00:00:00Z').toISOString(),
    }
    if (f.register_deadline) {
      payload.register_deadline = new Date(f.register_deadline).toISOString()
    }
    await scheduleApi.createPeriod(payload)
    showPeriodModal.value = false
    await loadPeriods()
  } catch (e) {
    periodError.value = e.response?.data?.error || e.message
  } finally {
    savingPeriod.value = false
  }
}

async function changeStatus(p, status) {
  const labels = { open: 'mở đăng ký', closed: 'đóng đăng ký', published: 'công bố lịch', cancelled: 'huỷ lịch' }
  if (!confirm(`Xác nhận ${labels[status]} kỳ "${p.name}"?`)) return
  try {
    await scheduleApi.setPeriodStatus(p.id, status)
    await loadPeriods()
    if (selectedPeriod.value?.id === p.id) await loadPeriodDetail(p.id)
  } catch (e) {
    alert(e.response?.data?.error || e.message)
  }
}

async function selectPeriod(p) {
  if (selectedPeriod.value?.id === p.id) {
    selectedPeriod.value = null; periodDetail.value = null; return
  }
  selectedPeriod.value = p
  await loadPeriodDetail(p.id)
}

async function loadPeriodDetail(id) {
  loadingDetail.value = true
  try { periodDetail.value = await scheduleApi.getPeriodDetail(id) }
  finally { loadingDetail.value = false }
}

async function addSlot(p) {
  try {
    await scheduleApi.addSlot(p.id, newSlot.value.templateId, new Date(newSlot.value.date).toISOString())
    newSlot.value = { templateId: '', date: '' }
    await loadPeriodDetail(p.id)
  } catch (e) { alert(e.response?.data?.error || e.message) }
}

async function deleteSlot(p, sl) {
  if (!confirm('Xóa ca này?')) return
  try {
    await scheduleApi.removeSlot(p.id, sl.id)
    await loadPeriodDetail(p.id)
  } catch (e) { alert(e.response?.data?.error || e.message) }
}

async function cancelReg(p, sl, reg) {
  if (!confirm(`Xóa đăng ký của ${reg.user_name} khỏi ca "${sl.template_name}"?`)) return
  try {
    await scheduleApi.cancelRegistration(reg.id)
    await loadPeriodDetail(p.id)
  } catch (e) { alert(e.response?.data?.error || e.message) }
}

// ── Weekly slot creation ────────────────────────────────
const weekDays = [
  { key: 1, label: 'T2' },
  { key: 2, label: 'T3' },
  { key: 3, label: 'T4' },
  { key: 4, label: 'T5' },
  { key: 5, label: 'T6' },
  { key: 6, label: 'T7' },
  { key: 0, label: 'CN' },
]
const weekPattern = ref({ 0: '', 1: '', 2: '', 3: '', 4: '', 5: '', 6: '' })
const weekRange = ref({ start: '', end: '' })
const showWeeklyModal = ref(false)
const weeklyAdding = ref(false)
const weeklyProgress = ref(0)
const weeklyError = ref('')
let weeklyPeriod = null

// Returns Monday of the week containing date d
function toMonday(d) {
  const day = d.getDay() // 0=Sun
  const diff = (day === 0 ? -6 : 1 - day)
  const m = new Date(d)
  m.setDate(m.getDate() + diff)
  m.setHours(0, 0, 0, 0)
  return m
}

const weekPreview = computed(() => {
  if (!weekRange.value.start) return []
  const startMonday = toMonday(new Date(weekRange.value.start))
  const endMonday = weekRange.value.end
    ? toMonday(new Date(weekRange.value.end))
    : startMonday
  if (endMonday < startMonday) return []

  const items = []
  const cur = new Date(startMonday)
  while (cur <= endMonday) {
    // iterate 7 days of this week
    for (let i = 0; i < 7; i++) {
      const d = new Date(cur)
      d.setDate(d.getDate() + i)
      const dow = d.getDay() // 0=Sun,1=Mon,...
      const templateId = weekPattern.value[dow]
      if (!templateId) continue
      const tmpl = templates.value.find(t => t.id === templateId)
      if (!tmpl) continue
      items.push({
        date: d.toISOString().split('T')[0],
        dateLabel: d.toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit' }),
        templateId,
        templateName: tmpl.name,
      })
    }
    cur.setDate(cur.getDate() + 7)
  }
  return items
})

function openWeeklyModal(p) {
  weeklyPeriod = p
  weeklyError.value = ''
  weeklyAdding.value = false
  weeklyProgress.value = 0
  // default: week containing period start date
  if (p.start_date) {
    const monday = toMonday(new Date(p.start_date))
    weekRange.value.start = monday.toISOString().split('T')[0]
    const endMonday = toMonday(new Date(p.end_date || p.start_date))
    weekRange.value.end = endMonday.toISOString().split('T')[0]
  } else {
    weekRange.value = { start: '', end: '' }
  }
  showWeeklyModal.value = true
}

async function applyWeeklyPattern() {
  if (!weekPreview.value.length || !weeklyPeriod) return
  weeklyAdding.value = true
  weeklyProgress.value = 0
  weeklyError.value = ''
  const errors = []
  for (const item of weekPreview.value) {
    try {
      await scheduleApi.addSlot(weeklyPeriod.id, item.templateId, new Date(item.date).toISOString())
    } catch (e) {
      errors.push(`${item.dateLabel}: ${e.response?.data?.error || e.message}`)
    }
    weeklyProgress.value++
  }
  weeklyAdding.value = false
  if (errors.length) {
    weeklyError.value = `Lỗi ${errors.length} ca: ${errors.slice(0, 3).join('; ')}${errors.length > 3 ? '...' : ''}`
  } else {
    showWeeklyModal.value = false
  }
  await loadPeriodDetail(weeklyPeriod.id)
}

// ── Helpers ────────────────────────────────────────────
function fmtDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function fmtDateShort(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit' })
}

function statusLabel(s) {
  return { draft: 'Nháp', open: 'Đang mở', closed: 'Đã đóng', published: 'Đã công bố', cancelled: 'Đã huỷ' }[s] || s
}

onMounted(() => { loadTemplates(); loadPeriods() })
</script>

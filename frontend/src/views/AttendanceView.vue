<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">

    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 pb-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <h1 class="text-2xl font-bold text-gray-900 mb-3">⏱️ Giờ công thực tế</h1>

        <!-- Week navigation -->
        <div class="flex items-center gap-2">
          <button @click="changeWeek(-1)"
            class="w-10 h-10 flex items-center justify-center rounded-xl bg-gray-100 active:bg-gray-200 touch-manipulation text-gray-600 font-bold text-lg">‹</button>
          <div class="flex-1 text-center text-sm font-semibold text-gray-700">
            {{ fmtDate(weekStart) }} — {{ fmtDate(weekEnd) }}
          </div>
          <button @click="changeWeek(1)"
            class="w-10 h-10 flex items-center justify-center rounded-xl bg-gray-100 active:bg-gray-200 touch-manipulation text-gray-600 font-bold text-lg">›</button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-3 pb-24">

      <!-- Action buttons -->
      <div class="flex gap-2 mb-4">
        <button @click="openAddForm"
          class="flex-1 py-3 rounded-2xl bg-blue-500 text-white font-bold text-base shadow-sm active:bg-blue-600 touch-manipulation">
          + Thêm giờ công
        </button>
        <button @click="openGenerateModal"
          class="flex-1 py-3 rounded-2xl bg-emerald-500 text-white font-bold text-base shadow-sm active:bg-emerald-600 touch-manipulation">
          ⚡ Tạo từ lịch
        </button>
        <button @click="openSummaryModal"
          class="flex-1 py-3 rounded-2xl bg-violet-500 text-white font-bold text-base shadow-sm active:bg-violet-600 touch-manipulation">
          📊 Tổng kết
        </button>
      </div>

      <!-- Staff summary -->
      <div v-if="staffSummary.length" class="mb-4">
        <div class="text-xs font-semibold text-gray-500 mb-2 uppercase tracking-wide">Tổng tuần này</div>
        <div class="flex flex-wrap gap-2">
          <div v-for="s in staffSummary" :key="s.name"
            :class="[
              'px-3 py-1.5 rounded-xl text-sm font-semibold',
              s.total > 0 ? 'bg-green-100 text-green-700' :
              s.total < 0 ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-600'
            ]">
            {{ s.name }}: {{ s.total > 0 ? '+' : '' }}{{ s.total }}h
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
      </div>

      <!-- Empty -->
      <div v-else-if="!entries.length" class="text-center py-12 text-gray-400">
        <div class="text-4xl mb-2">📋</div>
        <p>Chưa có giờ công nào trong tuần này.</p>
      </div>

      <!-- Entries grouped by date -->
      <div v-else class="flex flex-col gap-3">
        <div v-for="group in groupedEntries" :key="group.date">
          <div class="text-xs font-semibold text-gray-500 mb-1.5 uppercase tracking-wide">
            {{ group.dateLabel }}
          </div>
          <div class="bg-white rounded-2xl shadow-sm overflow-hidden border border-gray-100">
            <div v-for="(e, idx) in group.entries" :key="e.id"
              :class="['px-4 py-3 flex items-center justify-between', idx > 0 ? 'border-t border-gray-100' : '']">
              <div class="flex-1 min-w-0">
                <div class="font-semibold text-gray-900 text-sm">{{ e.staff_name }}</div>
                <div v-if="e.note" class="text-xs text-gray-400 mt-0.5 truncate">{{ e.note }}</div>
                <div class="text-xs text-gray-300 mt-0.5">bởi {{ e.created_by_name }}</div>
              </div>
              <div class="flex items-center gap-2 ml-2 shrink-0">
                <span :class="['text-lg font-bold', e.hours > 0 ? 'text-green-600' : 'text-red-500']">
                  {{ e.hours > 0 ? '+' : '' }}{{ e.hours }}h
                </span>
                <button @click="openEditForm(e)"
                  class="w-7 h-7 flex items-center justify-center rounded-full bg-gray-100 text-gray-500 active:bg-blue-100 active:text-blue-500 touch-manipulation text-sm">
                  ✏️
                </button>
                <button @click="removeEntry(e)"
                  class="w-7 h-7 flex items-center justify-center rounded-full bg-gray-100 text-gray-400 active:bg-red-100 active:text-red-500 touch-manipulation text-base font-bold">
                  ×
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add / Edit Entry Modal -->
    <div v-if="showForm"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showForm = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[85vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-5"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-5">{{ editingId ? 'Sửa giờ công' : 'Thêm giờ công' }}</h3>

        <div class="flex flex-col gap-4">
          <div v-if="!editingId">
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Ngày</label>
            <input type="date" v-model="form.date"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base bg-white" />
          </div>
          <!-- Scheduled staff for this date (only when adding) -->
          <div v-if="!editingId">
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Tên nhân viên</label>
            <div v-if="loadingSchedule" class="text-xs text-gray-400 mb-2">Đang tải ca đăng ký...</div>
            <div v-else-if="scheduledStaff.length" class="flex flex-wrap gap-1.5 mb-2">
              <button
                v-for="s in scheduledStaff" :key="s.staff_name + s.template_name"
                @click="selectScheduledStaff(s)"
                :class="[
                  'px-3 py-1.5 rounded-xl text-xs font-semibold touch-manipulation border transition-all',
                  form.staff_name === s.staff_name
                    ? 'bg-blue-500 text-white border-blue-500'
                    : 'bg-blue-50 text-blue-700 border-blue-200 active:bg-blue-100'
                ]">
                {{ s.staff_name }}
                <span class="opacity-70 ml-1">{{ s.scheduled_hours }}h</span>
              </button>
            </div>
            <div v-else-if="form.date" class="text-xs text-gray-400 italic mb-2">Không có nhân viên đăng ký ca ngày này.</div>
            <input v-model="form.staff_name" placeholder="Nhập tên nhân viên (hoặc chọn ở trên)"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
          <div v-else>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Nhân viên</label>
            <div class="px-4 py-3 bg-gray-50 rounded-xl text-base text-gray-700 font-semibold">{{ form.staff_name }}</div>
          </div>
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">
              Số giờ
              <span class="font-normal text-gray-400">(dương = cộng thêm, âm = trừ đi)</span>
            </label>
            <div class="flex items-center gap-3">
              <button @click="form.hours = +(((form.hours || 0) - 0.5).toFixed(1))"
                class="w-11 h-11 flex items-center justify-center rounded-xl bg-red-50 text-red-500 font-bold text-xl active:bg-red-100 touch-manipulation">−</button>
              <input type="number" v-model.number="form.hours" step="0.5" min="-24" max="24"
                class="flex-1 px-4 py-3 border border-gray-200 rounded-xl text-base text-center font-bold" />
              <button @click="form.hours = +(((form.hours || 0) + 0.5).toFixed(1))"
                class="w-11 h-11 flex items-center justify-center rounded-xl bg-green-50 text-green-500 font-bold text-xl active:bg-green-100 touch-manipulation">+</button>
            </div>
            <div class="flex gap-2 mt-2 flex-wrap">
              <button v-for="h in quickHours" :key="h"
                @click="form.hours = h"
                :class="[
                  'px-3 py-1 rounded-lg text-xs font-semibold touch-manipulation',
                  form.hours === h
                    ? (h > 0 ? 'bg-green-500 text-white' : 'bg-red-500 text-white')
                    : 'bg-gray-100 text-gray-600 active:bg-gray-200'
                ]">
                {{ h > 0 ? '+' : '' }}{{ h }}h
              </button>
            </div>
          </div>
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">
              Ghi chú <span class="font-normal text-gray-400">(tuỳ chọn)</span>
            </label>
            <input v-model="form.note" placeholder="VD: Tăng ca cuối tuần, về sớm..."
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
          </div>
        </div>

        <p v-if="formError" class="mt-3 text-red-500 text-sm">{{ formError }}</p>

        <div class="flex gap-3 mt-6">
          <button @click="showForm = false"
            class="flex-1 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
            Hủy
          </button>
          <button @click="saveEntry" :disabled="saving"
            class="flex-1 py-3 rounded-xl bg-blue-500 text-white font-semibold active:bg-blue-600 touch-manipulation disabled:opacity-50">
            {{ saving ? 'Đang lưu...' : 'Lưu' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Generate from Schedule Modal -->
    <div v-if="showGenerateModal"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showGenerateModal = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[90vh] flex flex-col">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-4"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-1">Tạo giờ công từ lịch đăng ký</h3>
        <p v-if="activePeriod" class="text-sm text-gray-500 mb-4">
          {{ activePeriod.name }} ·
          {{ fmtDateShort(activePeriod.start_date.split('T')[0]) }} — {{ fmtDateShort(activePeriod.end_date.split('T')[0]) }}
        </p>
        <p v-else class="text-sm text-gray-400 mb-4">Đang tải kỳ đăng ký...</p>

        <div v-if="loadingGenerate" class="flex justify-center py-8">
          <div class="w-8 h-8 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
        </div>

        <div v-else-if="generateError && !generateItems.length" class="text-center py-8 text-gray-400">
          <div class="text-3xl mb-2">📭</div>
          <p class="text-red-500 text-sm">{{ generateError }}</p>
        </div>

        <div v-else-if="!generateItems.length" class="text-center py-8 text-gray-400">
          <div class="text-3xl mb-2">📭</div>
          <p>Không có nhân viên nào đăng ký ca trong kỳ này.</p>
        </div>

        <div v-else class="flex-1 overflow-y-auto flex flex-col gap-2 mb-4">
          <div v-for="(item, idx) in generateItems" :key="idx"
            :class="[
              'flex items-center gap-3 px-4 py-3 rounded-xl border transition-all',
              item.selected ? 'bg-emerald-50 border-emerald-200' : 'bg-gray-50 border-gray-200 opacity-60'
            ]">
            <button @click="item.selected = !item.selected"
              :class="[
                'w-6 h-6 rounded-full flex-shrink-0 border-2 flex items-center justify-center text-xs font-bold',
                item.selected ? 'bg-emerald-500 border-emerald-500 text-white' : 'border-gray-300'
              ]">
              {{ item.selected ? '✓' : '' }}
            </button>
            <div class="flex-1 min-w-0">
              <div class="font-semibold text-gray-900 text-sm">{{ item.staff_name }}</div>
              <div class="text-xs text-gray-500">
                {{ fmtDateShort(item.date) }} · {{ item.template_name }} ({{ item.start_time }}–{{ item.end_time }})
              </div>
            </div>
            <div class="flex items-center gap-1 shrink-0">
              <button @click="adjustGenerateHours(item, -0.5)"
                class="w-6 h-6 flex items-center justify-center rounded-lg bg-red-50 text-red-500 font-bold text-sm active:bg-red-100 touch-manipulation">−</button>
              <span class="w-12 text-center text-sm font-bold text-emerald-700">{{ item.hours }}h</span>
              <button @click="adjustGenerateHours(item, 0.5)"
                class="w-6 h-6 flex items-center justify-center rounded-lg bg-green-50 text-green-600 font-bold text-sm active:bg-green-100 touch-manipulation">+</button>
            </div>
          </div>
        </div>

        <p v-if="generateError && generateItems.length" class="mb-3 text-red-500 text-sm">{{ generateError }}</p>

        <div v-if="generateItems.length" class="flex gap-3">
          <button @click="showGenerateModal = false"
            class="flex-1 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
            Hủy
          </button>
          <button @click="confirmGenerate" :disabled="generating || !selectedGenerateItems.length"
            class="flex-2 px-6 py-3 rounded-xl bg-emerald-500 text-white font-semibold active:bg-emerald-600 touch-manipulation disabled:opacity-50 whitespace-nowrap">
            {{ generating ? 'Đang tạo...' : `Tạo ${selectedGenerateItems.length} giờ công` }}
          </button>
        </div>
      </div>
    </div>

    <!-- Summary Modal -->
    <div v-if="showSummaryModal"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showSummaryModal = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[90vh] flex flex-col">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-4"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-4">📊 Tổng kết giờ làm</h3>

        <!-- Filters -->
        <div class="flex flex-col gap-3 mb-4">
          <div>
            <label class="text-xs font-semibold text-gray-500 mb-1 block">Nhân viên</label>
            <input v-model="summaryForm.staffName" list="staff-names-list"
              placeholder="Nhập tên nhân viên..."
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base" />
            <datalist id="staff-names-list">
              <option v-for="name in knownStaffNames" :key="name" :value="name" />
            </datalist>
          </div>
          <div class="flex gap-2">
            <div class="flex-1">
              <label class="text-xs font-semibold text-gray-500 mb-1 block">Từ ngày</label>
              <input type="date" v-model="summaryForm.from"
                class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm" />
            </div>
            <div class="flex-1">
              <label class="text-xs font-semibold text-gray-500 mb-1 block">Đến ngày</label>
              <input type="date" v-model="summaryForm.to"
                class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm" />
            </div>
          </div>
          <button @click="loadSummary" :disabled="loadingSummary || !summaryForm.staffName || !summaryForm.from || !summaryForm.to"
            class="w-full py-3 rounded-xl bg-violet-500 text-white font-semibold active:bg-violet-600 touch-manipulation disabled:opacity-50">
            {{ loadingSummary ? 'Đang tải...' : 'Xem tổng kết' }}
          </button>
        </div>

        <!-- Result -->
        <div v-if="summaryData" class="flex-1 overflow-y-auto">
          <!-- Total banner -->
          <div class="flex items-center justify-between bg-violet-50 rounded-2xl px-4 py-3 mb-3">
            <div>
              <div class="font-bold text-violet-800 text-base">{{ summaryData.staff_name }}</div>
              <div class="text-xs text-violet-500 mt-0.5">
                {{ fmtDateShort(summaryForm.from) }} — {{ fmtDateShort(summaryForm.to) }}
              </div>
            </div>
            <div :class="['text-2xl font-bold', summaryData.total_hours >= 0 ? 'text-green-600' : 'text-red-500']">
              {{ summaryData.total_hours > 0 ? '+' : '' }}{{ summaryData.total_hours }}h
            </div>
          </div>

          <!-- Entries table -->
          <div v-if="!summaryData.entries.length" class="text-center py-6 text-gray-400 text-sm">
            Không có giờ công nào trong khoảng thời gian này.
          </div>
          <div v-else class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
            <div v-for="(e, idx) in summaryData.entries" :key="e.id"
              :class="['px-4 py-3 flex items-center justify-between', idx > 0 ? 'border-t border-gray-100' : '']">
              <div class="flex-1 min-w-0">
                <div class="text-sm font-semibold text-gray-800">
                  {{ new Date(e.date + 'T12:00:00Z').toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit' }) }}
                </div>
                <div v-if="e.note" class="text-xs text-gray-400 mt-0.5 truncate">{{ e.note }}</div>
              </div>
              <span :class="['text-base font-bold ml-3 shrink-0', e.hours > 0 ? 'text-green-600' : 'text-red-500']">
                {{ e.hours > 0 ? '+' : '' }}{{ e.hours }}h
              </span>
            </div>
          </div>
        </div>

        <p v-if="summaryError" class="mt-3 text-red-500 text-sm">{{ summaryError }}</p>

        <button @click="showSummaryModal = false"
          class="mt-4 w-full py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold active:bg-gray-50 touch-manipulation">
          Đóng
        </button>
      </div>
    </div>

    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import * as attendanceApi from '../services/attendance.js'
const { getStaffSummary } = attendanceApi
import { getRegistrationsByDate, getRegistrationsByPeriod, getCurrentPeriod } from '../services/schedule.js'

const quickHours = [-4, -2, -1, -0.5, 0.5, 1, 2, 4, 8]

// ── Week navigation ──────────────────────────────────
const weekOffset = ref(0)

function getMonday(d) {
  const day = d.getDay()
  const diff = day === 0 ? -6 : 1 - day
  const m = new Date(d)
  m.setDate(m.getDate() + diff)
  m.setHours(0, 0, 0, 0)
  return m
}

// Dùng ngày LOCAL (không phải UTC) để tránh lệch múi giờ UTC+7
function toLocalDateStr(d) {
  return [
    d.getFullYear(),
    String(d.getMonth() + 1).padStart(2, '0'),
    String(d.getDate()).padStart(2, '0'),
  ].join('-')
}

const weekStart = computed(() => {
  const m = getMonday(new Date())
  m.setDate(m.getDate() + weekOffset.value * 7)
  return m
})

const weekEnd = computed(() => {
  const d = new Date(weekStart.value)
  d.setDate(d.getDate() + 6)
  return d
})

function changeWeek(delta) {
  weekOffset.value += delta
  loadEntries()
}

// ── Entries ──────────────────────────────────────────
const entries = ref([])
const loading = ref(false)

async function loadEntries() {
  loading.value = true
  try {
    const from = toLocalDateStr(weekStart.value)
    const to = toLocalDateStr(weekEnd.value)
    entries.value = await attendanceApi.getEntries(from, to) || []
  } finally {
    loading.value = false
  }
}

const groupedEntries = computed(() => {
  const map = {}
  for (const e of entries.value) {
    const date = e.date.split('T')[0]
    if (!map[date]) map[date] = []
    map[date].push(e)
  }
  return Object.entries(map)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([date, list]) => ({
      date,
      dateLabel: new Date(date + 'T12:00:00Z').toLocaleDateString('vi-VN', {
        weekday: 'long', day: '2-digit', month: '2-digit',
      }),
      entries: list,
    }))
})

const staffSummary = computed(() => {
  const map = {}
  for (const e of entries.value) {
    if (!map[e.staff_name]) map[e.staff_name] = 0
    map[e.staff_name] += e.hours
  }
  return Object.entries(map)
    .map(([name, total]) => ({ name, total: Math.round(total * 10) / 10 }))
    .sort((a, b) => b.total - a.total)
})

// ── Add / Edit Form ───────────────────────────────────
const showForm = ref(false)
const editingId = ref(null)
const saving = ref(false)
const formError = ref('')
const form = ref({ date: '', staff_name: '', hours: 1, note: '' })
const scheduledStaff = ref([])
const loadingSchedule = ref(false)

async function fetchScheduledStaff(date) {
  if (!date) { scheduledStaff.value = []; return }
  loadingSchedule.value = true
  try {
    scheduledStaff.value = await getRegistrationsByDate(date) || []
  } catch {
    scheduledStaff.value = []
  } finally {
    loadingSchedule.value = false
  }
}

function selectScheduledStaff(s) {
  form.value.staff_name = s.staff_name
  form.value.hours = s.scheduled_hours
  form.value.note = `${s.template_name} (${s.start_time}–${s.end_time})`
}

function openAddForm() {
  editingId.value = null
  formError.value = ''
  const today = new Date().toISOString().split('T')[0]
  form.value = { date: today, staff_name: '', hours: 1, note: '' }
  fetchScheduledStaff(today)
  showForm.value = true
}

function openEditForm(entry) {
  editingId.value = entry.id
  formError.value = ''
  form.value = {
    date: entry.date.split('T')[0],
    staff_name: entry.staff_name,
    hours: entry.hours,
    note: entry.note || '',
  }
  showForm.value = true
}

watch(() => form.value.date, (date) => {
  if (!editingId.value) fetchScheduledStaff(date)
})

async function saveEntry() {
  formError.value = ''
  if (!form.value.hours) { formError.value = 'Số giờ không được bằng 0'; return }

  saving.value = true
  try {
    if (editingId.value) {
      await attendanceApi.updateEntry(editingId.value, {
        hours: form.value.hours,
        note: form.value.note,
      })
    } else {
      if (!form.value.staff_name.trim()) { formError.value = 'Vui lòng nhập tên nhân viên'; saving.value = false; return }
      if (!form.value.date) { formError.value = 'Vui lòng chọn ngày'; saving.value = false; return }
      await attendanceApi.addEntry({
        date: new Date(form.value.date + 'T00:00:00Z').toISOString(),
        staff_name: form.value.staff_name.trim(),
        hours: form.value.hours,
        note: form.value.note,
      })
    }
    showForm.value = false
    await loadEntries()
  } catch (e) {
    formError.value = e.response?.data?.error || e.message
  } finally {
    saving.value = false
  }
}

async function removeEntry(e) {
  if (!confirm(`Xóa giờ công của ${e.staff_name} (${e.hours > 0 ? '+' : ''}${e.hours}h)?`)) return
  try {
    await attendanceApi.deleteEntry(e.id)
    await loadEntries()
  } catch (err) {
    alert(err.response?.data?.error || err.message)
  }
}

// ── Generate from Schedule ────────────────────────────
const showGenerateModal = ref(false)
const generateItems = ref([])
const loadingGenerate = ref(false)
const generating = ref(false)
const generateError = ref('')
const activePeriod = ref(null)

const selectedGenerateItems = computed(() => generateItems.value.filter(i => i.selected))

async function openGenerateModal() {
  generateError.value = ''
  generateItems.value = []
  activePeriod.value = null
  showGenerateModal.value = true
  loadingGenerate.value = true
  try {
    const period = await getCurrentPeriod()
    if (!period) {
      generateError.value = 'Không có kỳ đăng ký nào đang mở hoặc đã công bố.'
      return
    }
    activePeriod.value = period
    const list = await getRegistrationsByPeriod(period.id) || []
    generateItems.value = list.map(item => ({
      ...item,
      hours: item.scheduled_hours,
      selected: true,
    }))
  } catch (e) {
    generateError.value = e.response?.data?.error || e.message
  } finally {
    loadingGenerate.value = false
  }
}

function adjustGenerateHours(item, delta) {
  item.hours = Math.round((item.hours + delta) * 10) / 10
  if (item.hours === 0) item.hours = delta > 0 ? 0.5 : -0.5
}

async function confirmGenerate() {
  generateError.value = ''
  const selected = selectedGenerateItems.value
  if (!selected.length) return
  generating.value = true
  try {
    const payload = selected.map(item => ({
      date: new Date(item.date + 'T00:00:00Z').toISOString(),
      staff_name: item.staff_name,
      hours: item.hours,
      note: `${item.template_name} (${item.start_time}–${item.end_time})`,
    }))
    await attendanceApi.bulkAddEntries(payload)
    showGenerateModal.value = false
    await loadEntries()
  } catch (e) {
    generateError.value = e.response?.data?.error || e.message
  } finally {
    generating.value = false
  }
}

// ── Summary ──────────────────────────────────────────
const showSummaryModal = ref(false)
const loadingSummary = ref(false)
const summaryData = ref(null)
const summaryError = ref('')
const summaryForm = ref({ staffName: '', from: '', to: '' })

const knownStaffNames = computed(() => {
  const names = new Set(entries.value.map(e => e.staff_name).filter(Boolean))
  return [...names].sort()
})

function openSummaryModal() {
  summaryData.value = null
  summaryError.value = ''
  // Default: current week
  summaryForm.value = {
    staffName: '',
    from: toLocalDateStr(weekStart.value),
    to: toLocalDateStr(weekEnd.value),
  }
  showSummaryModal.value = true
}

async function loadSummary() {
  summaryError.value = ''
  summaryData.value = null
  loadingSummary.value = true
  try {
    summaryData.value = await getStaffSummary(summaryForm.value.staffName, summaryForm.value.from, summaryForm.value.to)
  } catch (e) {
    summaryError.value = e.response?.data?.error || e.message
  } finally {
    loadingSummary.value = false
  }
}

// ── Helpers ──────────────────────────────────────────
function fmtDate(d) {
  return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function fmtDateShort(dateStr) {
  return new Date(dateStr + 'T12:00:00Z').toLocaleDateString('vi-VN', {
    weekday: 'short', day: '2-digit', month: '2-digit',
  })
}

onMounted(loadEntries)
</script>

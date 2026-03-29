<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">

    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
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

      <!-- Add button -->
      <button @click="openAddForm"
        class="w-full mb-4 py-3 rounded-2xl bg-blue-500 text-white font-bold text-base shadow-sm active:bg-blue-600 touch-manipulation">
        + Thêm giờ công
      </button>

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
              <div class="flex items-center gap-3 ml-2 shrink-0">
                <span :class="['text-lg font-bold', e.hours > 0 ? 'text-green-600' : 'text-red-500']">
                  {{ e.hours > 0 ? '+' : '' }}{{ e.hours }}h
                </span>
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

    <!-- Add Entry Modal -->
    <div v-if="showForm"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-40"
      @click.self="showForm = false">
      <div class="bg-white rounded-t-3xl w-full max-w-lg p-6 pb-8 max-h-[85vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-gray-200 rounded-full mx-auto mb-5"></div>
        <h3 class="text-xl font-bold text-gray-900 mb-5">Thêm giờ công</h3>

        <div class="flex flex-col gap-4">
          <div>
            <label class="text-sm font-semibold text-gray-700 mb-1.5 block">Ngày</label>
            <input type="date" v-model="form.date"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl text-base bg-white" />
          </div>
          <!-- Scheduled staff for this date -->
          <div>
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

    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import * as attendanceApi from '../services/attendance.js'
import { getRegistrationsByDate } from '../services/schedule.js'

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
    const from = weekStart.value.toISOString().split('T')[0]
    const to = weekEnd.value.toISOString().split('T')[0]
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

// ── Form ─────────────────────────────────────────────
const showForm = ref(false)
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
  formError.value = ''
  const today = new Date().toISOString().split('T')[0]
  form.value = { date: today, staff_name: '', hours: 1, note: '' }
  fetchScheduledStaff(today)
  showForm.value = true
}

watch(() => form.value.date, (date) => {
  fetchScheduledStaff(date)
})

async function saveEntry() {
  formError.value = ''
  if (!form.value.staff_name.trim()) { formError.value = 'Vui lòng nhập tên nhân viên'; return }
  if (!form.value.date) { formError.value = 'Vui lòng chọn ngày'; return }
  if (!form.value.hours) { formError.value = 'Số giờ không được bằng 0'; return }

  saving.value = true
  try {
    await attendanceApi.addEntry({
      date: new Date(form.value.date + 'T00:00:00Z').toISOString(),
      staff_name: form.value.staff_name.trim(),
      hours: form.value.hours,
      note: form.value.note,
    })
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

// ── Helpers ──────────────────────────────────────────
function fmtDate(d) {
  return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(loadEntries)
</script>

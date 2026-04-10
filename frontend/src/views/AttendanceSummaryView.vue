<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Top bar -->
    <div class="flex items-center gap-3 px-4 pb-3 bg-white border-b border-gray-100" style="padding-top: max(1.25rem, env(safe-area-inset-top))">
      <button @click="$router.back()"
        class="w-9 h-9 flex items-center justify-center rounded-xl bg-gray-100 active:bg-gray-200 touch-manipulation text-gray-600 font-bold text-lg">
        ‹
      </button>
      <h1 class="text-lg font-bold text-gray-900">Báo cáo giờ công</h1>
    </div>

    <!-- Filters -->
    <div class="px-4 py-3 bg-white border-b border-gray-100 flex flex-col gap-2">
      <input v-model="form.staffName" list="staff-names-list"
        placeholder="Nhân viên..."
        class="w-full px-3 py-2 border border-gray-200 rounded-xl text-sm" />
      <datalist id="staff-names-list">
        <option v-for="name in staffNames" :key="name" :value="name" />
      </datalist>
      <div class="flex gap-2">
        <input type="date" v-model="form.from"
          class="flex-1 px-3 py-2 border border-gray-200 rounded-xl text-sm" />
        <input type="date" v-model="form.to"
          class="flex-1 px-3 py-2 border border-gray-200 rounded-xl text-sm" />
        <button @click="load" :disabled="loading || !form.staffName || !form.from || !form.to"
          class="shrink-0 px-4 py-2 rounded-xl bg-violet-500 text-white text-sm font-semibold active:bg-violet-600 touch-manipulation disabled:opacity-50">
          {{ loading ? '...' : 'Xem' }}
        </button>
      </div>
      <p v-if="error" class="text-red-500 text-xs">{{ error }}</p>
    </div>

    <!-- Report -->
    <div v-if="data" class="flex-1 px-4 py-4">
      <div class="border border-gray-200 rounded-2xl overflow-hidden">
        <!-- Header -->
        <div class="bg-gradient-to-br from-violet-600 to-violet-500 px-4 py-3 text-white flex items-center justify-between gap-3">
          <div class="min-w-0">
            <div class="text-[10px] font-bold uppercase tracking-widest opacity-60">Báo cáo giờ công</div>
            <div class="text-base font-bold leading-tight truncate">{{ data.staff_name }}</div>
            <div class="text-xs opacity-75 mt-0.5">{{ fmtDateShort(form.from) }} — {{ fmtDateShort(form.to) }}</div>
          </div>
          <div class="shrink-0 text-right">
            <div class="text-[10px] opacity-60 uppercase tracking-widest">Tổng</div>
            <div :class="['text-2xl font-bold', data.total_hours >= 0 ? 'text-green-300' : 'text-red-300']">
              {{ data.total_hours > 0 ? '+' : '' }}{{ data.total_hours }}h
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-if="!data.entries.length" class="text-center py-10 text-gray-400 text-sm">
          Không có giờ công nào trong khoảng thời gian này.
        </div>

        <!-- Grid -->
        <div v-else class="grid grid-cols-3 gap-px bg-gray-100">
          <div v-for="e in data.entries" :key="e.id"
            class="bg-white p-2.5 h-[76px] flex flex-col justify-between">
            <div>
              <div class="text-sm font-bold text-gray-800 leading-none">
                {{ fmtWeekday(e.date.split('T')[0]) }} {{ fmtDayMonth(e.date.split('T')[0]) }}
              </div>
              <div v-if="e.note" class="text-[10px] text-gray-400 leading-tight mt-0.5 line-clamp-2">{{ e.note }}</div>
            </div>
            <div :class="['text-sm font-bold', e.hours > 0 ? 'text-green-600' : 'text-red-500']">
              {{ e.hours > 0 ? '+' : '' }}{{ e.hours }}h
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Placeholder before search -->
    <div v-else-if="!loading" class="flex-1 flex items-center justify-center text-gray-300 text-sm">
      Chọn nhân viên và khoảng ngày để xem báo cáo
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getStaffSummary, getEntries } from '../services/attendance.js'

const route = useRoute()

const form = ref({
  staffName: '',
  from: route.query.from || '',
  to: route.query.to || '',
})

const data = ref(null)
const loading = ref(false)
const error = ref('')
const staffNames = ref([])

async function load() {
  if (!form.value.staffName || !form.value.from || !form.value.to) return
  error.value = ''
  data.value = null
  loading.value = true
  try {
    data.value = await getStaffSummary(form.value.staffName, form.value.from, form.value.to)
  } catch (e) {
    error.value = e.response?.data?.error || e.message
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // Load staff name suggestions from current week entries
  try {
    const entries = await getEntries(form.value.from, form.value.to)
    const names = new Set((entries || []).map(e => e.staff_name).filter(Boolean))
    staffNames.value = [...names].sort()
  } catch {}
})

function fmtDateShort(dateStr) {
  return new Date(dateStr + 'T12:00:00Z').toLocaleDateString('vi-VN', {
    weekday: 'short', day: '2-digit', month: '2-digit',
  })
}

function fmtWeekday(dateStr) {
  return new Date(dateStr + 'T12:00:00Z').toLocaleDateString('vi-VN', { weekday: 'short' })
}

function fmtDayMonth(dateStr) {
  return new Date(dateStr + 'T12:00:00Z').toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })
}
</script>

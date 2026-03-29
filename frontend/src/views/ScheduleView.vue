<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">

    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">📅 Lịch ca làm việc</h1>

        <!-- Tabs -->
        <div class="flex gap-2">
          <button
            @click="tab = 'register'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              tab === 'register'
                ? 'bg-blue-500 text-white shadow-lg'
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📋</span>
              <span>Đăng ký ca</span>
            </div>
          </button>
          <button
            @click="tab = 'mine'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              tab === 'mine'
                ? 'bg-blue-500 text-white shadow-lg'
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>✅</span>
              <span>Ca của tôi</span>
              <span v-if="mySlots.length > 0" class="bg-white bg-opacity-30 text-xs px-1.5 py-0.5 rounded-full">
                {{ mySlots.length }}
              </span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-3 pb-24">

      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-16 gap-3">
        <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 text-sm">Đang tải...</p>
      </div>

      <!-- No active period -->
      <div v-else-if="!period" class="flex flex-col items-center justify-center py-16 gap-3">
        <div class="text-5xl">📭</div>
        <p class="text-gray-600 font-semibold">Chưa có kỳ đăng ký nào</p>
        <p class="text-gray-400 text-sm text-center">Admin chưa mở kỳ đăng ký ca.<br>Vui lòng chờ thông báo.</p>
      </div>

      <div v-else>
        <!-- Period banner -->
        <div :class="[
          'rounded-2xl p-4 mb-4 shadow-sm',
          period.status === 'open' ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white' :
          period.status === 'published' ? 'bg-gradient-to-r from-blue-500 to-indigo-500 text-white' :
          'bg-gray-100 text-gray-700'
        ]">
          <div class="flex justify-between items-start">
            <div>
              <h2 class="font-bold text-lg">{{ period.name }}</h2>
              <p class="text-sm opacity-90 mt-0.5">
                {{ fmtDate(period.start_date) }} – {{ fmtDate(period.end_date) }}
              </p>
            </div>
            <span class="bg-white bg-opacity-25 px-3 py-1 rounded-full text-sm font-semibold">
              {{ statusLabel(period.status) }}
            </span>
          </div>
          <div v-if="period.status === 'open'" class="mt-2 text-sm opacity-90 flex items-center gap-1">
            <span>⏰</span>
            <span>Hạn đăng ký: {{ fmtDateTime(period.register_deadline) }}</span>
          </div>
          <div v-if="isPastDeadline && period.status === 'open'" class="mt-1 text-sm font-semibold">
            ⚠️ Đã quá hạn đăng ký
          </div>
        </div>

        <!-- ── Tab: Đăng ký ca ── -->
        <div v-if="tab === 'register'">
          <div v-if="Object.keys(slotsByDate).length === 0" class="text-center py-12 text-gray-400">
            <div class="text-4xl mb-2">📭</div>
            <p>Kỳ này chưa có ca nào được tạo.</p>
          </div>

          <div v-for="(daySlots, dateKey) in slotsByDate" :key="dateKey" class="mb-4">
            <!-- Day header -->
            <div class="flex items-center gap-2 mb-2">
              <div class="h-px flex-1 bg-gray-200"></div>
              <span class="text-xs font-semibold text-gray-500 uppercase tracking-wide">
                {{ fmtDateHeader(dateKey) }}
              </span>
              <div class="h-px flex-1 bg-gray-200"></div>
            </div>

            <!-- Slot cards -->
            <div class="flex flex-col gap-3">
              <div
                v-for="sl in daySlots"
                :key="sl.id"
                :class="[
                  'bg-white rounded-2xl p-4 shadow-sm border-2 transition-all',
                  sl.my_registration ? 'border-blue-400 bg-blue-50' :
                  isFull(sl) ? 'border-gray-100 opacity-60' :
                  'border-transparent'
                ]"
              >
                <div class="flex justify-between items-start mb-2">
                  <div>
                    <div class="font-bold text-gray-900">{{ sl.template_name }}</div>
                    <div class="text-sm text-gray-500 flex items-center gap-1 mt-0.5">
                      <span>🕐</span>
                      <span>{{ sl.start_time }} – {{ sl.end_time }}</span>
                    </div>
                  </div>
                  <div class="text-right">
                    <div :class="[
                      'text-sm font-semibold px-2.5 py-1 rounded-full',
                      isFull(sl) ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-700'
                    ]">
                      {{ sl.registered_count }}/{{ sl.max_slots }}
                    </div>
                    <div class="text-xs text-gray-400 mt-0.5">
                      {{ isFull(sl) ? 'Hết chỗ' : `Còn ${sl.max_slots - sl.registered_count} chỗ` }}
                    </div>
                  </div>
                </div>

                <!-- Registrant chips -->
                <div v-if="sl.registrations && sl.registrations.length > 0" class="flex flex-wrap gap-1.5 mb-3">
                  <span
                    v-for="r in sl.registrations"
                    :key="r.id"
                    :class="[
                      'text-xs px-2 py-0.5 rounded-full font-medium',
                      r.user_id === myUserId ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-600'
                    ]"
                  >
                    {{ r.user_name }}
                  </span>
                </div>

                <!-- Action button -->
                <div v-if="period.status === 'open' && !isPastDeadline">
                  <button
                    v-if="sl.my_registration"
                    @click="unregister(sl)"
                    :disabled="actioningSlot === sl.id"
                    class="w-full py-2.5 rounded-xl border-2 border-red-300 text-red-500 font-semibold text-sm active:bg-red-50 transition-all touch-manipulation disabled:opacity-50"
                  >
                    {{ actioningSlot === sl.id ? 'Đang hủy...' : '✕ Hủy đăng ký' }}
                  </button>
                  <button
                    v-else-if="!isFull(sl)"
                    @click="register(sl)"
                    :disabled="actioningSlot === sl.id"
                    class="w-full py-2.5 rounded-xl bg-blue-500 text-white font-semibold text-sm active:bg-blue-600 transition-all touch-manipulation disabled:opacity-50 shadow-sm"
                  >
                    {{ actioningSlot === sl.id ? 'Đang đăng ký...' : '+ Đăng ký ca này' }}
                  </button>
                  <div v-else class="w-full py-2.5 text-center text-gray-400 text-sm">
                    Ca đã đủ người
                  </div>
                </div>
                <div v-else>
                  <div v-if="sl.my_registration" class="flex items-center gap-1.5 text-blue-600 text-sm font-semibold">
                    <span>✅</span><span>Đã đăng ký</span>
                  </div>
                  <div v-else-if="period.status !== 'open'" class="text-xs text-gray-400 italic">
                    {{ period.status === 'published' ? '📌 Lịch đã công bố' : '🔒 Kỳ chưa mở đăng ký' }}
                  </div>
                  <div v-else-if="isPastDeadline" class="text-xs text-orange-500 font-medium">
                    ⏰ Đã hết hạn đăng ký ({{ fmtDateTime(period.register_deadline) }})
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- ── Tab: Ca của tôi ── -->
        <div v-if="tab === 'mine'">
          <div v-if="mySlots.length === 0" class="flex flex-col items-center justify-center py-16 gap-3">
            <div class="text-4xl">📭</div>
            <p class="text-gray-600 font-semibold">Bạn chưa đăng ký ca nào</p>
            <p class="text-gray-400 text-sm">Chuyển sang tab "Đăng ký ca" để đăng ký.</p>
          </div>

          <div v-else class="flex flex-col gap-3">
            <div
              v-for="sl in mySlots"
              :key="sl.id"
              class="bg-white rounded-2xl p-4 shadow-sm border-2 border-blue-300"
            >
              <div class="flex justify-between items-center">
                <div>
                  <div class="font-bold text-gray-900">{{ sl.template_name }}</div>
                  <div class="text-sm text-gray-500 mt-0.5">{{ fmtDateFull(sl.date) }}</div>
                  <div class="text-sm text-blue-600 font-medium mt-0.5">
                    🕐 {{ sl.start_time }} – {{ sl.end_time }}
                  </div>
                </div>
                <div class="flex flex-col items-end gap-2">
                  <span class="bg-blue-100 text-blue-700 text-xs font-semibold px-2.5 py-1 rounded-full">
                    ✅ Đã đăng ký
                  </span>
                  <button
                    v-if="period.status === 'open' && !isPastDeadline"
                    @click="unregister(sl)"
                    :disabled="actioningSlot === sl.id"
                    class="text-xs text-red-400 underline touch-manipulation disabled:opacity-50"
                  >
                    {{ actioningSlot === sl.id ? 'Đang hủy...' : 'Hủy' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import BottomNav from '../components/BottomNav.vue'
import * as scheduleApi from '../services/schedule.js'

const tab = ref('register')
const loading = ref(false)
const period = ref(null)
const actioningSlot = ref(null)

const authStore = useAuthStore()
const myUserId = computed(() => authStore.user?.id || authStore.user?._id)

async function load() {
  loading.value = true
  try {
    period.value = await scheduleApi.getCurrentPeriod()
  } catch {
    period.value = null
  } finally {
    loading.value = false
  }
}

const isPastDeadline = computed(() => {
  if (!period.value?.register_deadline) return false
  return new Date() > new Date(period.value.register_deadline)
})

const mySlots = computed(() => {
  if (!period.value?.slots) return []
  return period.value.slots.filter(sl => sl.my_registration)
})

const slotsByDate = computed(() => {
  if (!period.value?.slots) return {}
  const map = {}
  for (const sl of period.value.slots) {
    const key = sl.date ? sl.date.substring(0, 10) : 'unknown'
    if (!map[key]) map[key] = []
    map[key].push(sl)
  }
  return map
})

function isFull(sl) {
  return sl.max_slots > 0 && sl.registered_count >= sl.max_slots
}

async function register(sl) {
  actioningSlot.value = sl.id
  try {
    await scheduleApi.registerSlot(sl.id)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || e.message)
  } finally {
    actioningSlot.value = null
  }
}

async function unregister(sl) {
  if (!confirm('Hủy đăng ký ca này?')) return
  actioningSlot.value = sl.id
  try {
    await scheduleApi.unregisterSlot(sl.id)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || e.message)
  } finally {
    actioningSlot.value = null
  }
}

function fmtDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function fmtDateTime(d) {
  if (!d) return '—'
  return new Date(d).toLocaleString('vi-VN', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function fmtDateHeader(key) {
  if (!key || key === 'unknown') return '—'
  const d = new Date(key)
  return d.toLocaleDateString('vi-VN', { weekday: 'long', day: '2-digit', month: '2-digit' })
}

function fmtDateFull(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('vi-VN', { weekday: 'long', day: '2-digit', month: '2-digit', year: 'numeric' })
}

function statusLabel(s) {
  return { draft: 'Nháp', open: 'Đang mở', closed: 'Đã đóng', published: 'Đã công bố' }[s] || s
}

onMounted(load)
</script>

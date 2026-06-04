<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center gap-3">
          <button @click="$router.back()" class="text-gray-600 text-xl leading-none">←</button>
          <div>
            <h1 class="text-xl font-bold text-gray-800">🤖 AI Command</h1>
            <p class="text-xs text-gray-500">Trợ lý AI thông minh</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Chat thread -->
    <div ref="threadRef" class="flex-1 overflow-y-auto px-4 py-4 pb-4 flex flex-col gap-3">
      <div v-if="loading" class="text-center text-sm text-gray-400 py-8">Đang tải lịch sử...</div>

      <div v-if="!loading && messages.length === 0" class="flex flex-col items-center justify-center py-12 gap-3">
        <div class="text-5xl">🤖</div>
        <p class="text-gray-500 text-sm text-center">Xin chào! Tôi có thể giúp bạn<br>quản lý nguyên liệu, chi phí và nhiều hơn nữa.</p>
      </div>

      <template v-for="(msg, i) in messages" :key="i">
        <div v-if="msg.role === 'user'" class="flex justify-end">
          <div class="bg-blue-500 text-white px-4 py-3 rounded-2xl rounded-br-sm text-sm max-w-[75%] shadow-sm">
            {{ msg.message }}
          </div>
        </div>
        <div v-else class="flex flex-col gap-2 items-start">
          <div class="flex items-start gap-2 max-w-[85%]">
            <div class="w-7 h-7 bg-emerald-100 rounded-full flex items-center justify-center text-sm flex-shrink-0 mt-0.5">🤖</div>
            <div class="bg-white border border-gray-200 px-4 py-3 rounded-2xl rounded-bl-sm text-sm text-gray-800 shadow-sm">
              {{ msg.message }}
            </div>
          </div>
          <div class="ml-9">
            <IngredientActionCard
              v-if="msg.action_type === 'add_ingredient' && msg.action"
              :fields="msg.action.fields"
              @confirmed="onConfirmed"
            />
            <RestockActionCard
              v-else-if="msg.action_type === 'restock_ingredient' && msg.action"
              :fields="msg.action.fields"
              @confirmed="onConfirmed"
            />
            <ExpenseActionCard
              v-else-if="msg.action_type === 'add_expense' && msg.action"
              :fields="msg.action.fields"
              :categories="categories"
              @confirmed="onConfirmed"
            />
          </div>
        </div>
      </template>

      <div v-if="thinking" class="flex items-center gap-2">
        <div class="w-7 h-7 bg-emerald-100 rounded-full flex items-center justify-center text-sm flex-shrink-0">🤖</div>
        <div class="bg-white border border-gray-200 px-4 py-3 rounded-2xl rounded-bl-sm shadow-sm flex items-center gap-1">
          <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></span>
          <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay:0.15s"></span>
          <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay:0.3s"></span>
        </div>
      </div>
    </div>

    <!-- Input bar -->
    <div class="bg-white border-t border-gray-200 px-4 py-3 flex gap-2 items-center flex-shrink-0">
      <input
        v-model="input"
        @keydown.enter.prevent="send"
        class="flex-1 bg-gray-100 rounded-full px-4 py-3 text-sm outline-none focus:bg-white focus:ring-2 focus:ring-emerald-400 transition-all"
        placeholder="Nhập lệnh cho AI..."
        :disabled="thinking"
      />
      <button
        @click="send"
        :disabled="!input.trim() || thinking"
        class="w-10 h-10 bg-emerald-500 disabled:bg-gray-300 rounded-full text-white flex items-center justify-center flex-shrink-0 shadow-sm active:scale-95 transition-transform"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="19" x2="12" y2="5"/><polyline points="5 12 12 5 19 12"/>
        </svg>
      </button>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { aiService } from '../services/aiService'
import { expenseService } from '../services/expense'
import IngredientActionCard from '../components/ai/IngredientActionCard.vue'
import RestockActionCard from '../components/ai/RestockActionCard.vue'
import ExpenseActionCard from '../components/ai/ExpenseActionCard.vue'

const messages = ref([])
const input = ref('')
const thinking = ref(false)
const loading = ref(true)
const threadRef = ref(null)
const categories = ref([])

onMounted(async () => {
  try {
    const [logs, cats] = await Promise.all([
      aiService.getHistory(),
      expenseService.getCategories(),
    ])
    messages.value = logs.map(l => ({ ...l, action: null }))
    categories.value = cats || []
  } catch (e) {
    console.error('Failed to load AI history or categories', e)
  } finally {
    loading.value = false
    scrollToBottom()
  }
})

async function send() {
  const text = input.value.trim()
  if (!text || thinking.value) return
  input.value = ''
  messages.value.push({ role: 'user', message: text, action_type: null, action: null })
  thinking.value = true
  scrollToBottom()
  try {
    const history = messages.value.slice(-20).map(m => ({ role: m.role, message: m.message }))
    const resp = await aiService.parseCommand(text, history)
    messages.value.push({
      role: 'agent',
      message: resp.reply_text,
      action_type: resp.action?.type || null,
      action: resp.action || null,
    })
  } catch {
    messages.value.push({ role: 'agent', message: 'Có lỗi xảy ra. Vui lòng thử lại.', action_type: null, action: null })
  } finally {
    thinking.value = false
    scrollToBottom()
  }
}

function onConfirmed(payload) {
  const label = payload.type === 'add_expense' ? 'ghi nhận chi phí' : 'cập nhật nguyên liệu'
  messages.value.push({ role: 'agent', message: `✅ Đã ${label} thành công! Bạn cần làm gì thêm không?`, action_type: null, action: null })
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => { if (threadRef.value) threadRef.value.scrollTop = threadRef.value.scrollHeight })
}
</script>

<template>
  <div class="flex flex-col h-screen bg-gray-50">
    <div class="bg-white border-b border-gray-200 px-4 py-3 flex items-center gap-2 shrink-0">
      <span class="text-lg">🤖</span>
      <h1 class="font-semibold text-gray-900">AI Command</h1>
    </div>

    <div ref="threadRef" class="flex-1 overflow-y-auto px-4 py-4 flex flex-col gap-3">
      <div v-if="loading" class="text-center text-sm text-gray-400 py-8">Đang tải lịch sử...</div>

      <template v-for="(msg, i) in messages" :key="i">
        <div v-if="msg.role === 'user'" class="flex justify-end">
          <div class="bg-blue-500 text-white px-4 py-2 rounded-2xl rounded-br-sm text-sm max-w-xs">
            {{ msg.message }}
          </div>
        </div>
        <div v-else class="flex flex-col gap-2 items-start">
          <div class="bg-white border border-gray-200 px-4 py-2 rounded-2xl rounded-bl-sm text-sm max-w-xs text-gray-800">
            {{ msg.message }}
          </div>
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
      </template>

      <div v-if="thinking" class="flex items-center gap-1 text-gray-400 text-sm">
        <span class="animate-bounce">●</span><span class="animate-bounce" style="animation-delay:0.1s">●</span><span class="animate-bounce" style="animation-delay:0.2s">●</span>
      </div>
    </div>

    <div class="bg-white border-t border-gray-200 px-4 py-3 flex gap-2 items-center shrink-0">
      <input
        v-model="input"
        @keydown.enter.prevent="send"
        class="flex-1 bg-gray-50 border border-gray-200 rounded-full px-4 py-2 text-sm outline-none focus:border-blue-400"
        placeholder="Nhập lệnh cho AI..."
        :disabled="thinking"
      />
      <button
        @click="send"
        :disabled="!input.trim() || thinking"
        class="w-10 h-10 bg-blue-500 disabled:bg-gray-300 rounded-full text-white text-lg flex items-center justify-center shrink-0"
      >↑</button>
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

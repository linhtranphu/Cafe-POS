<template>
  <transition name="slide-down">
    <div
      v-if="visible"
      class="fixed top-0 left-0 right-0 z-50 mx-4 mt-4"
      style="padding-top: env(safe-area-inset-top)"
    >
      <div
        :class="[
          'rounded-2xl shadow-2xl p-4 border-2',
          notificationClass
        ]"
      >
        <!-- Header -->
        <div class="flex items-start justify-between mb-2">
          <div class="flex items-center gap-2">
            <span class="text-2xl">{{ icon }}</span>
            <h3 class="font-bold text-lg">{{ title }}</h3>
          </div>
          <button
            @click="dismiss"
            class="text-gray-400 hover:text-gray-600 text-xl"
          >
            ×
          </button>
        </div>

        <!-- Message -->
        <p class="text-sm mb-3 pl-8">{{ message }}</p>

        <!-- Error Details (if provided) -->
        <div v-if="errorDetails" class="mb-3 pl-8">
          <details class="text-xs">
            <summary class="cursor-pointer text-gray-600 hover:text-gray-800">
              Chi tiết lỗi
            </summary>
            <pre class="mt-2 p-2 bg-gray-100 rounded text-xs overflow-x-auto">{{ errorDetails }}</pre>
          </details>
        </div>

        <!-- Actions -->
        <div v-if="showActions" class="flex gap-2 pl-8">
          <button
            v-if="onRetry"
            @click="handleRetry"
            :disabled="retrying"
            class="flex-1 bg-blue-500 text-white py-2 px-4 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="retrying" class="inline-block w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span v-else>🔄</span>
            <span>{{ retrying ? 'Đang thử...' : 'Thử lại' }}</span>
          </button>
          <button
            @click="dismiss"
            class="flex-1 bg-gray-200 text-gray-800 py-2 px-4 rounded-lg text-sm font-bold hover:bg-gray-300 active:scale-95 transition-transform"
          >
            Đóng
          </button>
        </div>

        <!-- Auto-dismiss Progress Bar -->
        <div v-if="autoDismiss && !paused" class="mt-3">
          <div class="h-1 bg-gray-200 rounded-full overflow-hidden">
            <div
              class="h-full bg-current transition-all"
              :style="{ width: `${progress}%` }"
            ></div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'error',
    validator: (value) => ['error', 'warning', 'info', 'success'].includes(value)
  },
  title: {
    type: String,
    default: 'Lỗi in ấn'
  },
  message: {
    type: String,
    required: true
  },
  errorDetails: {
    type: String,
    default: null
  },
  autoDismiss: {
    type: Boolean,
    default: true
  },
  dismissTimeout: {
    type: Number,
    default: 5000 // 5 seconds
  },
  onRetry: {
    type: Function,
    default: null
  },
  onDismiss: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['dismiss', 'retry'])

const visible = ref(false)
const paused = ref(false)
const retrying = ref(false)
const progress = ref(100)
let dismissTimer = null
let progressInterval = null

const showActions = computed(() => {
  return props.onRetry !== null
})

const icon = computed(() => {
  const icons = {
    error: '❌',
    warning: '⚠️',
    info: 'ℹ️',
    success: '✅'
  }
  return icons[props.type] || icons.error
})

const notificationClass = computed(() => {
  const classes = {
    error: 'bg-red-50 border-red-300 text-red-800',
    warning: 'bg-yellow-50 border-yellow-300 text-yellow-800',
    info: 'bg-blue-50 border-blue-300 text-blue-800',
    success: 'bg-green-50 border-green-300 text-green-800'
  }
  return classes[props.type] || classes.error
})

const show = () => {
  visible.value = true
  if (props.autoDismiss) {
    startAutoDismiss()
  }
}

const dismiss = () => {
  clearTimers()
  visible.value = false
  if (props.onDismiss) {
    props.onDismiss()
  }
  emit('dismiss')
}

const handleRetry = async () => {
  if (retrying.value || !props.onRetry) return

  retrying.value = true
  pauseAutoDismiss()

  try {
    await props.onRetry()
    // If retry succeeds, dismiss the notification
    dismiss()
  } catch (error) {
    console.error('Retry failed:', error)
    // Keep notification visible if retry fails
    resumeAutoDismiss()
  } finally {
    retrying.value = false
  }
}

const startAutoDismiss = () => {
  clearTimers()
  progress.value = 100

  const startTime = Date.now()
  const duration = props.dismissTimeout

  progressInterval = setInterval(() => {
    const elapsed = Date.now() - startTime
    progress.value = Math.max(0, 100 - (elapsed / duration) * 100)
  }, 50)

  dismissTimer = setTimeout(() => {
    dismiss()
  }, duration)
}

const pauseAutoDismiss = () => {
  paused.value = true
  clearTimers()
}

const resumeAutoDismiss = () => {
  paused.value = false
  if (props.autoDismiss && visible.value) {
    startAutoDismiss()
  }
}

const clearTimers = () => {
  if (dismissTimer) {
    clearTimeout(dismissTimer)
    dismissTimer = null
  }
  if (progressInterval) {
    clearInterval(progressInterval)
    progressInterval = null
  }
}

// Pause auto-dismiss on hover
const handleMouseEnter = () => {
  if (props.autoDismiss) {
    pauseAutoDismiss()
  }
}

const handleMouseLeave = () => {
  if (props.autoDismiss) {
    resumeAutoDismiss()
  }
}

// Watch for message changes to restart timer
watch(() => props.message, () => {
  if (visible.value && props.autoDismiss) {
    startAutoDismiss()
  }
})

onMounted(() => {
  show()
})

onUnmounted(() => {
  clearTimers()
})

// Expose methods for parent components
defineExpose({
  show,
  dismiss,
  pauseAutoDismiss,
  resumeAutoDismiss
})
</script>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-down-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}

.active\:scale-95:active {
  transform: scale(0.95);
}
</style>

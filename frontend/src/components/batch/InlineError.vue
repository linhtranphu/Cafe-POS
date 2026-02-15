<template>
  <div 
    v-if="show"
    class="inline-error rounded-lg p-3 mb-4 flex items-start gap-3"
    :class="severityClass"
  >
    <!-- Icon -->
    <div class="flex-shrink-0 text-xl">
      {{ icon }}
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <!-- Message -->
      <p class="text-sm font-medium" :class="textClass">
        {{ message }}
      </p>

      <!-- Details (if provided) -->
      <p v-if="details" class="text-xs mt-1" :class="detailsTextClass">
        {{ details }}
      </p>
    </div>

    <!-- Actions -->
    <div v-if="showRetry || dismissible" class="flex-shrink-0 flex gap-2">
      <!-- Retry Button -->
      <button
        v-if="showRetry"
        @click="handleRetry"
        :disabled="retrying"
        class="text-sm underline hover:no-underline disabled:opacity-50"
        :class="textClass"
      >
        {{ retrying ? 'Đang thử...' : 'Thử lại' }}
      </button>

      <!-- Dismiss Button -->
      <button
        v-if="dismissible"
        @click="handleDismiss"
        class="text-sm hover:opacity-70"
        :class="textClass"
      >
        ✕
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  // Content
  message: {
    type: String,
    required: true
  },
  details: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: '⚠️'
  },
  
  // Severity level
  severity: {
    type: String,
    default: 'error', // 'error', 'warning', 'info'
    validator: (value) => ['error', 'warning', 'info'].includes(value)
  },
  
  // Display control
  show: {
    type: Boolean,
    default: true
  },
  dismissible: {
    type: Boolean,
    default: true
  },
  
  // Retry functionality
  showRetry: {
    type: Boolean,
    default: false
  },
  onRetry: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['dismiss', 'retry'])

const retrying = ref(false)

// Computed classes based on severity
const severityClass = computed(() => {
  const classes = {
    error: 'bg-red-50 border-2 border-red-200',
    warning: 'bg-yellow-50 border-2 border-yellow-200',
    info: 'bg-blue-50 border-2 border-blue-200'
  }
  return classes[props.severity] || classes.error
})

const textClass = computed(() => {
  const classes = {
    error: 'text-red-800',
    warning: 'text-yellow-800',
    info: 'text-blue-800'
  }
  return classes[props.severity] || classes.error
})

const detailsTextClass = computed(() => {
  const classes = {
    error: 'text-red-600',
    warning: 'text-yellow-600',
    info: 'text-blue-600'
  }
  return classes[props.severity] || classes.error
})

// Handlers
const handleRetry = async () => {
  if (!props.onRetry || retrying.value) return
  
  retrying.value = true
  try {
    await props.onRetry()
    emit('retry')
  } catch (error) {
    console.error('Retry failed:', error)
  } finally {
    retrying.value = false
  }
}

const handleDismiss = () => {
  emit('dismiss')
}
</script>

<style scoped>
.inline-error {
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

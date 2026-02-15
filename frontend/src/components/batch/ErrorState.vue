<template>
  <div 
    class="error-state rounded-2xl p-6 text-center"
    :class="[
      variant === 'inline' ? 'bg-red-50 border-2 border-red-200' : 'bg-white shadow-lg',
      sizeClass
    ]"
  >
    <!-- Error Icon -->
    <div class="error-icon text-6xl mb-4" v-if="showIcon">
      {{ icon }}
    </div>

    <!-- Error Title -->
    <h3 
      v-if="title" 
      class="font-semibold mb-2"
      :class="titleClass"
    >
      {{ title }}
    </h3>

    <!-- Error Message -->
    <p 
      class="mb-4"
      :class="messageClass"
    >
      {{ message }}
    </p>

    <!-- Error Details (expandable) -->
    <div v-if="details && showDetails" class="mb-4">
      <button
        @click="detailsExpanded = !detailsExpanded"
        class="text-sm text-gray-600 hover:text-gray-800 underline"
      >
        {{ detailsExpanded ? 'Ẩn chi tiết' : 'Xem chi tiết' }}
      </button>
      <div 
        v-if="detailsExpanded"
        class="mt-2 p-3 bg-gray-100 rounded-lg text-left text-sm text-gray-700"
      >
        {{ details }}
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex gap-3 justify-center flex-wrap">
      <!-- Retry Button -->
      <button
        v-if="showRetry && retryable"
        @click="handleRetry"
        :disabled="retrying"
        class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
      >
        <span v-if="retrying" class="animate-spin">⟳</span>
        <span>{{ retrying ? 'Đang thử lại...' : 'Thử lại' }}</span>
      </button>

      <!-- Custom Action Button -->
      <button
        v-if="actionLabel && actionHandler"
        @click="handleAction"
        class="px-6 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
      >
        {{ actionLabel }}
      </button>

      <!-- Go Back Button -->
      <button
        v-if="showGoBack"
        @click="handleGoBack"
        class="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
      >
        Quay lại
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  // Error content
  icon: {
    type: String,
    default: '❌'
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  },
  details: {
    type: String,
    default: ''
  },
  
  // Error type
  errorType: {
    type: String,
    default: 'unknown' // network, validation, permission, not_found, server, unknown
  },
  
  // Display options
  variant: {
    type: String,
    default: 'default', // 'default' or 'inline'
    validator: (value) => ['default', 'inline'].includes(value)
  },
  size: {
    type: String,
    default: 'medium', // 'small', 'medium', 'large'
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  showIcon: {
    type: Boolean,
    default: true
  },
  showDetails: {
    type: Boolean,
    default: true
  },
  
  // Retry functionality
  retryable: {
    type: Boolean,
    default: true
  },
  showRetry: {
    type: Boolean,
    default: true
  },
  onRetry: {
    type: Function,
    default: null
  },
  
  // Custom action
  actionLabel: {
    type: String,
    default: ''
  },
  actionHandler: {
    type: Function,
    default: null
  },
  
  // Navigation
  showGoBack: {
    type: Boolean,
    default: false
  },
  goBackRoute: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['retry', 'action', 'goBack'])

const router = useRouter()
const retrying = ref(false)
const detailsExpanded = ref(false)

// Computed classes
const sizeClass = computed(() => {
  const sizes = {
    small: 'py-4',
    medium: 'py-6',
    large: 'py-8'
  }
  return sizes[props.size] || sizes.medium
})

const titleClass = computed(() => {
  const sizes = {
    small: 'text-base',
    medium: 'text-lg',
    large: 'text-xl'
  }
  return `${sizes[props.size] || sizes.medium} text-red-800`
})

const messageClass = computed(() => {
  const sizes = {
    small: 'text-sm',
    medium: 'text-base',
    large: 'text-lg'
  }
  return `${sizes[props.size] || sizes.medium} text-red-600`
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

const handleAction = () => {
  if (props.actionHandler) {
    props.actionHandler()
  }
  emit('action')
}

const handleGoBack = () => {
  if (props.goBackRoute) {
    router.push(props.goBackRoute)
  } else {
    router.back()
  }
  emit('goBack')
}
</script>

<style scoped>
.error-state {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-icon {
  animation: bounce 0.5s ease-in-out;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}
</style>

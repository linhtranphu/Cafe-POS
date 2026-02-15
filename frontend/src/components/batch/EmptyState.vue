<template>
  <div 
    class="empty-state text-center py-12"
    :class="containerClass"
  >
    <!-- Icon/Illustration -->
    <div class="empty-icon text-6xl mb-4">
      {{ icon }}
    </div>

    <!-- Title -->
    <h3 
      v-if="title"
      class="text-xl font-semibold text-gray-800 mb-2"
    >
      {{ title }}
    </h3>

    <!-- Description -->
    <p class="text-gray-600 mb-6 max-w-md mx-auto">
      {{ description }}
    </p>

    <!-- Action Button -->
    <button
      v-if="actionLabel"
      @click="handleAction"
      class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors inline-flex items-center gap-2"
    >
      <span v-if="actionIcon">{{ actionIcon }}</span>
      <span>{{ actionLabel }}</span>
    </button>

    <!-- Secondary Action -->
    <button
      v-if="secondaryLabel"
      @click="handleSecondaryAction"
      class="ml-3 px-6 py-3 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
    >
      {{ secondaryLabel }}
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  // Content
  icon: {
    type: String,
    default: '📦'
  },
  title: {
    type: String,
    default: 'Không có dữ liệu'
  },
  description: {
    type: String,
    default: 'Chưa có dữ liệu để hiển thị'
  },
  
  // Primary action
  actionLabel: {
    type: String,
    default: ''
  },
  actionIcon: {
    type: String,
    default: ''
  },
  onAction: {
    type: Function,
    default: null
  },
  
  // Secondary action
  secondaryLabel: {
    type: String,
    default: ''
  },
  onSecondaryAction: {
    type: Function,
    default: null
  },
  
  // Styling
  variant: {
    type: String,
    default: 'default', // 'default', 'compact'
    validator: (value) => ['default', 'compact'].includes(value)
  }
})

const emit = defineEmits(['action', 'secondaryAction'])

const containerClass = computed(() => {
  return props.variant === 'compact' ? 'py-8' : 'py-12'
})

const handleAction = () => {
  if (props.onAction) {
    props.onAction()
  }
  emit('action')
}

const handleSecondaryAction = () => {
  if (props.onSecondaryAction) {
    props.onSecondaryAction()
  }
  emit('secondaryAction')
}
</script>

<script>
import { computed } from 'vue'

export default {
  name: 'EmptyState'
}
</script>

<style scoped>
.empty-state {
  animation: fadeIn 0.4s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.empty-icon {
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}
</style>

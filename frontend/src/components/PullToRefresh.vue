<template>
  <transition name="fade">
    <div v-if="pullDistance > 0 || isRefreshing" 
      class="fixed top-0 left-0 right-0 z-50 flex justify-center pointer-events-none"
      :style="{ transform: `translateY(${pullDistance}px)` }">
      <div class="bg-white rounded-b-2xl shadow-lg px-6 py-3 flex items-center gap-3">
        <div v-if="isRefreshing" class="animate-spin text-2xl">🔄</div>
        <div v-else-if="pullDistance >= threshold" class="text-2xl">🎯</div>
        <div v-else class="text-2xl">⬇️</div>
        
        <span class="font-medium text-gray-700">
          {{ statusText }}
        </span>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  pullDistance: {
    type: Number,
    default: 0
  },
  isRefreshing: {
    type: Boolean,
    default: false
  },
  threshold: {
    type: Number,
    default: 80
  }
})

const statusText = computed(() => {
  if (props.isRefreshing) return 'Đang tải...'
  if (props.pullDistance >= props.threshold) return 'Thả để làm mới'
  if (props.pullDistance > 0) return 'Kéo xuống để làm mới'
  return ''
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>

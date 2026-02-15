<template>
  <div class="loading-skeleton">
    <!-- List Skeleton -->
    <div v-if="type === 'list'" class="space-y-4">
      <div 
        v-for="i in rows" 
        :key="i"
        class="skeleton-item bg-gray-200 rounded-lg p-4 animate-pulse"
      >
        <div class="flex items-center gap-4">
          <div class="skeleton-circle w-12 h-12 bg-gray-300 rounded-full"></div>
          <div class="flex-1 space-y-2">
            <div class="skeleton-line h-4 bg-gray-300 rounded w-3/4"></div>
            <div class="skeleton-line h-3 bg-gray-300 rounded w-1/2"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Card Skeleton -->
    <div v-else-if="type === 'card'" class="skeleton-card bg-gray-200 rounded-2xl p-6 animate-pulse">
      <div class="space-y-4">
        <div class="skeleton-line h-6 bg-gray-300 rounded w-1/2"></div>
        <div class="skeleton-line h-4 bg-gray-300 rounded w-full"></div>
        <div class="skeleton-line h-4 bg-gray-300 rounded w-3/4"></div>
        <div class="flex gap-2 mt-4">
          <div class="skeleton-button h-10 bg-gray-300 rounded-lg w-24"></div>
          <div class="skeleton-button h-10 bg-gray-300 rounded-lg w-24"></div>
        </div>
      </div>
    </div>

    <!-- Table Skeleton -->
    <div v-else-if="type === 'table'" class="skeleton-table">
      <div class="bg-gray-200 rounded-lg overflow-hidden">
        <!-- Header -->
        <div class="bg-gray-300 p-4 flex gap-4">
          <div 
            v-for="i in columns" 
            :key="i"
            class="skeleton-line h-4 bg-gray-400 rounded flex-1"
          ></div>
        </div>
        <!-- Rows -->
        <div 
          v-for="i in rows" 
          :key="i"
          class="p-4 border-b border-gray-300 flex gap-4 animate-pulse"
        >
          <div 
            v-for="j in columns" 
            :key="j"
            class="skeleton-line h-4 bg-gray-300 rounded flex-1"
          ></div>
        </div>
      </div>
    </div>

    <!-- Form Skeleton -->
    <div v-else-if="type === 'form'" class="skeleton-form space-y-6">
      <div 
        v-for="i in fields" 
        :key="i"
        class="animate-pulse"
      >
        <div class="skeleton-line h-4 bg-gray-300 rounded w-32 mb-2"></div>
        <div class="skeleton-input h-12 bg-gray-200 rounded-lg w-full"></div>
      </div>
      <div class="flex gap-3">
        <div class="skeleton-button h-12 bg-gray-300 rounded-lg w-32"></div>
        <div class="skeleton-button h-12 bg-gray-200 rounded-lg w-24"></div>
      </div>
    </div>

    <!-- Chart Skeleton -->
    <div v-else-if="type === 'chart'" class="skeleton-chart bg-gray-200 rounded-2xl p-6 animate-pulse">
      <div class="skeleton-line h-6 bg-gray-300 rounded w-1/3 mb-6"></div>
      <div class="flex items-end gap-2 h-48">
        <div 
          v-for="i in 8" 
          :key="i"
          class="skeleton-bar bg-gray-300 rounded-t flex-1"
          :style="{ height: `${Math.random() * 80 + 20}%` }"
        ></div>
      </div>
    </div>

    <!-- Text Skeleton -->
    <div v-else-if="type === 'text'" class="skeleton-text space-y-2 animate-pulse">
      <div 
        v-for="i in lines" 
        :key="i"
        class="skeleton-line h-4 bg-gray-300 rounded"
        :style="{ width: i === lines ? '60%' : '100%' }"
      ></div>
    </div>

    <!-- Default/Custom -->
    <div v-else class="skeleton-default bg-gray-200 rounded-lg animate-pulse" :style="customStyle">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Skeleton type
  type: {
    type: String,
    default: 'list', // 'list', 'card', 'table', 'form', 'chart', 'text', 'custom'
    validator: (value) => ['list', 'card', 'table', 'form', 'chart', 'text', 'custom'].includes(value)
  },
  
  // Configuration for different types
  rows: {
    type: Number,
    default: 3
  },
  columns: {
    type: Number,
    default: 4
  },
  fields: {
    type: Number,
    default: 4
  },
  lines: {
    type: Number,
    default: 3
  },
  
  // Custom styling
  height: {
    type: String,
    default: ''
  },
  width: {
    type: String,
    default: ''
  }
})

const customStyle = computed(() => {
  const style = {}
  if (props.height) style.height = props.height
  if (props.width) style.width = props.width
  return style
})
</script>

<style scoped>
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.skeleton-line,
.skeleton-circle,
.skeleton-button,
.skeleton-input,
.skeleton-bar {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.2) 20%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0)
  );
  background-size: 200% 100%;
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}
</style>

<template>
  <div class="animate-pulse">
    <!-- Card Skeleton -->
    <div v-if="type === 'card'" class="bg-white rounded-2xl p-4 shadow-sm">
      <div class="flex justify-between items-start mb-3">
        <div class="flex-1">
          <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
          <div class="h-4 bg-gray-200 rounded w-1/2"></div>
        </div>
        <div class="text-right ml-2">
          <div class="h-4 bg-gray-200 rounded w-16 mb-1"></div>
          <div class="h-6 bg-gray-200 rounded w-20"></div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-3 pt-3 border-t">
        <div>
          <div class="h-3 bg-gray-200 rounded w-20 mb-1"></div>
          <div class="h-5 bg-gray-200 rounded w-16"></div>
        </div>
        <div>
          <div class="h-3 bg-gray-200 rounded w-20 mb-1"></div>
          <div class="h-5 bg-gray-200 rounded w-16"></div>
        </div>
      </div>
    </div>

    <!-- Table Row Skeleton -->
    <tr v-else-if="type === 'table-row'">
      <td class="px-4 py-3"><div class="h-4 bg-gray-200 rounded w-32"></div></td>
      <td class="px-4 py-3"><div class="h-4 bg-gray-200 rounded w-20"></div></td>
      <td class="px-4 py-3 text-right"><div class="h-4 bg-gray-200 rounded w-24 ml-auto"></div></td>
      <td class="px-4 py-3 text-right"><div class="h-4 bg-gray-200 rounded w-24 ml-auto"></div></td>
      <td class="px-4 py-3 text-right"><div class="h-4 bg-gray-200 rounded w-16 ml-auto"></div></td>
      <td class="px-4 py-3 text-right"><div class="h-4 bg-gray-200 rounded w-24 ml-auto"></div></td>
      <td class="px-4 py-3 text-center"><div class="h-6 bg-gray-200 rounded w-20 mx-auto"></div></td>
    </tr>

    <!-- Summary Skeleton -->
    <div v-else-if="type === 'summary'" class="bg-gradient-to-br from-gray-300 to-gray-400 rounded-xl p-4 md:p-6 mb-4">
      <div class="h-4 bg-white/30 rounded w-40 mb-3"></div>
      <div class="grid grid-cols-4 gap-3 md:gap-6">
        <div class="text-center">
          <div class="h-8 bg-white/30 rounded w-12 mx-auto mb-1"></div>
          <div class="h-3 bg-white/30 rounded w-16 mx-auto"></div>
        </div>
        <div class="text-center">
          <div class="h-8 bg-white/30 rounded w-12 mx-auto mb-1"></div>
          <div class="h-3 bg-white/30 rounded w-16 mx-auto"></div>
        </div>
        <div class="text-center">
          <div class="h-8 bg-white/30 rounded w-12 mx-auto mb-1"></div>
          <div class="h-3 bg-white/30 rounded w-16 mx-auto"></div>
        </div>
        <div class="text-center">
          <div class="h-8 bg-white/30 rounded w-12 mx-auto mb-1"></div>
          <div class="h-3 bg-white/30 rounded w-16 mx-auto"></div>
        </div>
      </div>
    </div>

    <!-- Profit Section Skeleton -->
    <div v-else-if="type === 'profit-section'" class="bg-white rounded-2xl p-4 md:p-6 shadow-sm">
      <div class="h-6 bg-gray-200 rounded w-40 mb-4"></div>
      <div class="space-y-3">
        <div class="flex justify-between items-center">
          <div class="h-4 bg-gray-200 rounded w-32"></div>
          <div class="h-4 bg-gray-200 rounded w-24"></div>
        </div>
        <div class="flex justify-between items-center">
          <div class="h-4 bg-gray-200 rounded w-32"></div>
          <div class="h-4 bg-gray-200 rounded w-24"></div>
        </div>
        <div class="pt-3 border-t flex justify-between items-center">
          <div class="h-4 bg-gray-200 rounded w-32"></div>
          <div class="h-6 bg-gray-200 rounded w-32"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  type: {
    type: String,
    default: 'card',
    validator: (value) => ['card', 'table-row', 'summary', 'profit-section'].includes(value)
  }
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
</style>

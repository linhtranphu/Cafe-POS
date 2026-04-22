<template>
  <div class="flex justify-between items-center text-sm py-1">
    <span class="text-gray-500">{{ label }}</span>
    <span :class="highlight ? 'font-semibold text-green-700' : 'font-medium text-gray-900'">{{ value }}</span>
  </div>
</template>
<script setup>
defineProps({ label: String, value: [String, Number], highlight: Boolean })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="px-6 py-4">
        <div class="flex justify-between items-center">
          <h1 class="text-2xl font-bold text-gray-900">💵 Thu ngân System</h1>
          <div class="flex items-center space-x-4">
            <button
              @click="$router.push('/dashboard')"
              class="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded"
            >
              ← Quay lại
            </button>
          </div>
        </div>
        
        <!-- Navigation Tabs -->
        <div class="mt-4 border-b border-gray-200">
          <nav class="-mb-px flex space-x-8">
            <button
              @click="activeTab = 'dashboard'"
              class="py-2 px-1 border-b-2 font-medium text-sm"
              :class="activeTab === 'dashboard' 
                ? 'border-blue-500 text-blue-600' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
            >
              📊 Dashboard
            </button>
            <button
              @click="activeTab = 'reports'"
              class="py-2 px-1 border-b-2 font-medium text-sm"
              :class="activeTab === 'reports' 
                ? 'border-blue-500 text-blue-600' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
            >
              📋 Báo cáo
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="p-6">
      <slot :activeTab="activeTab" />
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'CashierLayout',
  setup() {
    const activeTab = ref('dashboard')

    return {
      activeTab
    }
  }
}
</script>
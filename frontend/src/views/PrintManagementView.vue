<template>
  <div class="h-screen flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3 border-b">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-bold text-gray-800">🖨️ Quản Lý In Ấn</h1>
          <p class="text-sm text-gray-600 mt-1">Quản lý máy in, templates và print jobs</p>
        </div>
        
        <!-- Local Bridge Status Indicator -->
        <div class="flex items-center gap-2 px-3 py-2 rounded-lg" :class="localBridgeAvailable ? 'bg-green-50' : 'bg-gray-50'">
          <div class="flex items-center gap-2">
            <div 
              class="w-2 h-2 rounded-full"
              :class="localBridgeAvailable ? 'bg-green-500 animate-pulse' : 'bg-gray-400'"
            ></div>
            <span class="text-xs font-medium" :class="localBridgeAvailable ? 'text-green-700' : 'text-gray-600'">
              {{ localBridgeAvailable ? 'Local Bridge Online' : 'Local Bridge Offline' }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 border-b">
      <div class="flex gap-1 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'px-6 py-3 text-sm font-medium whitespace-nowrap transition-colors border-b-2',
            activeTab === tab.id
              ? 'text-blue-600 border-blue-600'
              : 'text-gray-600 border-transparent hover:text-gray-800 hover:border-gray-300'
          ]"
        >
          <span class="mr-2">{{ tab.icon }}</span>
          <span>{{ tab.label }}</span>
        </button>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 overflow-hidden">
      <!-- Print Jobs Tab -->
      <div v-show="activeTab === 'jobs'" class="h-full">
        <PrintJobList />
      </div>

      <!-- Printers Tab -->
      <div v-show="activeTab === 'printers'" class="h-full">
        <PrinterConfigList />
      </div>

      <!-- Templates Tab -->
      <div v-show="activeTab === 'templates'" class="h-full">
        <PrintTemplateEditor />
      </div>

      <!-- Settings Tab -->
      <div v-show="activeTab === 'settings'" class="h-full">
        <ShopSettingsForm />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PrintJobList from '../components/printing/PrintJobList.vue'
import PrinterConfigList from '../components/printing/PrinterConfigList.vue'
import PrintTemplateEditor from '../components/printing/PrintTemplateEditor.vue'
import ShopSettingsForm from '../components/printing/ShopSettingsForm.vue'
import { usePrintJobStore } from '../stores/printJob'
import { useLocalPrint } from '../composables/useLocalPrint'

const activeTab = ref('jobs')
const printJobStore = usePrintJobStore()
const { isAvailable: localBridgeAvailable } = useLocalPrint()

const tabs = [
  { id: 'jobs', label: 'Print Jobs', icon: '📄' },
  { id: 'printers', label: 'Máy In', icon: '🖨️' },
  { id: 'templates', label: 'Templates', icon: '📝' },
  { id: 'settings', label: 'Cài Đặt', icon: '⚙️' }
]

// Initialize print job store on mount
onMounted(async () => {
  await printJobStore.initialize()
  console.log('[PrintManagement] Local bridge available:', localBridgeAvailable.value)
})
</script>

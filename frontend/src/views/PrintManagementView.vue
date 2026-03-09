<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">🖨️ Quản Lý In Ấn</h1>
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
                {{ localBridgeAvailable ? 'Online' : 'Offline' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Tabs -->
        <div class="flex gap-2 overflow-x-auto scrollbar-hide">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98 whitespace-nowrap',
              activeTab === tab.id
                ? 'bg-blue-500 text-white shadow-lg'
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]"
          >
            <span class="mr-1.5">{{ tab.icon }}</span>
            <span>{{ tab.label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 overflow-y-auto pb-24">
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

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PrintJobList from '../components/printing/PrintJobList.vue'
import PrinterConfigList from '../components/printing/PrinterConfigList.vue'
import PrintTemplateEditor from '../components/printing/PrintTemplateEditor.vue'
import ShopSettingsForm from '../components/printing/ShopSettingsForm.vue'
import BottomNav from '../components/BottomNav.vue'
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


<style scoped>
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.active\:scale-98:active {
  transform: scale(0.98);
}
</style>

<template>
  <div class="p-6 space-y-8 bg-gray-50 min-h-screen">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">Error Handling Components Demo</h1>

    <!-- ErrorState Examples -->
    <section class="space-y-4">
      <h2 class="text-2xl font-semibold text-gray-700">ErrorState Component</h2>
      
      <!-- Network Error -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Network Error (Retryable)</h3>
        <ErrorState
          icon="📡"
          title="Lỗi kết nối"
          message="Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng."
          errorType="network"
          :retryable="true"
          :onRetry="handleRetry"
          showGoBack
        />
      </div>

      <!-- Validation Error -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Validation Error</h3>
        <ErrorState
          icon="⚠️"
          title="Dữ liệu không hợp lệ"
          message="Số lượng batch phải lớn hơn 0"
          errorType="validation"
          :retryable="false"
          showGoBack
          variant="inline"
          size="small"
        />
      </div>

      <!-- Permission Error -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Permission Error</h3>
        <ErrorState
          icon="🔒"
          title="Không có quyền"
          message="Bạn không có quyền thực hiện thao tác này"
          errorType="permission"
          :retryable="false"
          actionLabel="Đăng nhập"
          :actionHandler="handleLogin"
        />
      </div>
    </section>

    <!-- InlineError Examples -->
    <section class="space-y-4">
      <h2 class="text-2xl font-semibold text-gray-700">InlineError Component</h2>
      
      <!-- Error Severity -->
      <div class="bg-white p-6 rounded-lg shadow space-y-4">
        <h3 class="text-lg font-medium mb-4">Different Severities</h3>
        
        <InlineError
          severity="error"
          message="Không thể ghi nhận batch. Không đủ nguyên liệu."
          :showRetry="true"
          :onRetry="handleRetry"
        />
        
        <InlineError
          severity="warning"
          message="Batch sắp hết hạn trong 2 giờ nữa"
          icon="⏰"
        />
        
        <InlineError
          severity="info"
          message="Batch đã được tạo thành công"
          icon="ℹ️"
          :dismissible="true"
        />
      </div>

      <!-- With Details -->
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">With Details</h3>
        <InlineError
          message="Lỗi tạo batch"
          details="Nguyên liệu 'Hạt Cà Phê' không đủ. Cần: 110g, Có: 50g"
          :showRetry="true"
          :onRetry="handleRetry"
        />
      </div>
    </section>

    <!-- EmptyState Examples -->
    <section class="space-y-4">
      <h2 class="text-2xl font-semibold text-gray-700">EmptyState Component</h2>
      
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">No Data</h3>
        <EmptyState
          icon="📦"
          title="Chưa có batch record"
          description="Bạn chưa ghi nhận batch nào. Hãy tạo batch record đầu tiên!"
          actionLabel="Ghi nhận batch mới"
          actionIcon="➕"
          :onAction="handleCreate"
          secondaryLabel="Xem hướng dẫn"
          :onSecondaryAction="handleGuide"
        />
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">No Alerts</h3>
        <EmptyState
          icon="✅"
          title="Không có cảnh báo"
          description="Tất cả batch đều ở trạng thái tốt"
          variant="compact"
        />
      </div>
    </section>

    <!-- LoadingSkeleton Examples -->
    <section class="space-y-4">
      <h2 class="text-2xl font-semibold text-gray-700">LoadingSkeleton Component</h2>
      
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">List Skeleton</h3>
        <LoadingSkeleton type="list" :rows="3" />
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Card Skeleton</h3>
        <LoadingSkeleton type="card" />
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Table Skeleton</h3>
        <LoadingSkeleton type="table" :rows="4" :columns="5" />
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Form Skeleton</h3>
        <LoadingSkeleton type="form" :fields="5" />
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-medium mb-4">Chart Skeleton</h3>
        <LoadingSkeleton type="chart" />
      </div>
    </section>

    <!-- Complete Flow Example -->
    <section class="space-y-4">
      <h2 class="text-2xl font-semibold text-gray-700">Complete Flow Example</h2>
      
      <div class="bg-white p-6 rounded-lg shadow">
        <div class="flex gap-4 mb-6">
          <button 
            @click="currentState = 'loading'"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
            Show Loading
          </button>
          <button 
            @click="currentState = 'error'"
            class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600">
            Show Error
          </button>
          <button 
            @click="currentState = 'empty'"
            class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600">
            Show Empty
          </button>
          <button 
            @click="currentState = 'data'"
            class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600">
            Show Data
          </button>
        </div>

        <!-- Loading State -->
        <LoadingSkeleton 
          v-if="currentState === 'loading'"
          type="list"
          :rows="3"
        />

        <!-- Error State -->
        <ErrorState
          v-else-if="currentState === 'error'"
          icon="❌"
          title="Lỗi tải dữ liệu"
          message="Không thể tải danh sách batch records"
          :retryable="true"
          :onRetry="() => currentState = 'loading'"
        />

        <!-- Empty State -->
        <EmptyState
          v-else-if="currentState === 'empty'"
          icon="📦"
          title="Chưa có dữ liệu"
          description="Không có batch record nào"
          actionLabel="Tạo mới"
          :onAction="() => currentState = 'data'"
        />

        <!-- Data State -->
        <div v-else class="space-y-2">
          <div v-for="i in 3" :key="i" class="p-4 border rounded-lg">
            <div class="font-medium">Batch Record {{ i }}</div>
            <div class="text-sm text-gray-600">Số lượng: {{ i * 100 }}ml</div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ErrorState from './ErrorState.vue'
import InlineError from './InlineError.vue'
import EmptyState from './EmptyState.vue'
import LoadingSkeleton from './LoadingSkeleton.vue'

const currentState = ref('loading')

const handleRetry = async () => {
  console.log('Retrying...')
  await new Promise(resolve => setTimeout(resolve, 1000))
  console.log('Retry complete')
}

const handleLogin = () => {
  console.log('Navigate to login')
}

const handleCreate = () => {
  console.log('Navigate to create form')
}

const handleGuide = () => {
  console.log('Show guide')
}
</script>

<style scoped>
/* Demo-specific styles */
</style>

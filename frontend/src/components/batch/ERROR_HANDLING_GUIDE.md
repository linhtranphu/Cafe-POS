# Batch Error Handling Guide

This guide explains how to use the error handling components and utilities in the batch management system.

## Components

### 1. ErrorState.vue

Full-page or section error display with retry functionality.

**Usage:**

```vue
<template>
  <ErrorState
    v-if="error"
    icon="⚠️"
    title="Không thể tải dữ liệu"
    :message="error"
    :retryable="true"
    :onRetry="loadData"
    showGoBack
    goBackRoute="/batch"
  />
</template>

<script setup>
import ErrorState from './ErrorState.vue'
import { ref } from 'vue'

const error = ref(null)

const loadData = async () => {
  try {
    // Load data logic
  } catch (err) {
    error.value = err.message
  }
}
</script>
```

**Props:**

- `icon` (String): Error icon emoji (default: '❌')
- `title` (String): Error title
- `message` (String, required): Error message
- `details` (String): Additional error details
- `errorType` (String): Error type for styling ('network', 'validation', 'permission', 'not_found', 'server', 'unknown')
- `variant` (String): Display variant ('default' or 'inline')
- `size` (String): Size ('small', 'medium', 'large')
- `showIcon` (Boolean): Show icon (default: true)
- `showDetails` (Boolean): Show expandable details (default: true)
- `retryable` (Boolean): Show retry button (default: true)
- `showRetry` (Boolean): Enable retry functionality (default: true)
- `onRetry` (Function): Retry handler function
- `actionLabel` (String): Custom action button label
- `actionHandler` (Function): Custom action handler
- `showGoBack` (Boolean): Show go back button
- `goBackRoute` (String): Route to navigate back to

### 2. InlineError.vue

Inline error display for forms and small sections.

**Usage:**

```vue
<template>
  <InlineError
    v-if="error"
    :message="error"
    :show="!!error"
    :showRetry="isRetryable"
    :onRetry="handleRetry"
    @dismiss="error = null"
  />
</template>

<script setup>
import InlineError from './InlineError.vue'
import { ref } from 'vue'

const error = ref(null)
const isRetryable = ref(false)

const handleRetry = async () => {
  // Retry logic
}
</script>
```

**Props:**

- `message` (String, required): Error message
- `details` (String): Additional details
- `icon` (String): Error icon (default: '⚠️')
- `severity` (String): Severity level ('error', 'warning', 'info')
- `show` (Boolean): Show/hide error (default: true)
- `dismissible` (Boolean): Allow dismissing (default: true)
- `showRetry` (Boolean): Show retry button (default: false)
- `onRetry` (Function): Retry handler

### 3. EmptyState.vue

Display when no data is available.

**Usage:**

```vue
<template>
  <EmptyState
    v-if="!loading && items.length === 0"
    icon="📦"
    title="Chưa có dữ liệu"
    description="Bạn chưa có batch nào. Hãy tạo batch đầu tiên!"
    actionLabel="Tạo batch mới"
    actionIcon="➕"
    :onAction="createBatch"
  />
</template>

<script setup>
import EmptyState from './EmptyState.vue'

const createBatch = () => {
  // Navigate to create form
}
</script>
```

**Props:**

- `icon` (String): Display icon (default: '📦')
- `title` (String): Title (default: 'Không có dữ liệu')
- `description` (String): Description
- `actionLabel` (String): Primary action button label
- `actionIcon` (String): Action button icon
- `onAction` (Function): Primary action handler
- `secondaryLabel` (String): Secondary action label
- `onSecondaryAction` (Function): Secondary action handler
- `variant` (String): Display variant ('default' or 'compact')

### 4. LoadingSkeleton.vue

Loading placeholder with various types.

**Usage:**

```vue
<template>
  <LoadingSkeleton 
    v-if="loading"
    type="list"
    :rows="5"
  />
</template>

<script setup>
import LoadingSkeleton from './LoadingSkeleton.vue'
</script>
```

**Props:**

- `type` (String): Skeleton type ('list', 'card', 'table', 'form', 'chart', 'text', 'custom')
- `rows` (Number): Number of rows for list/table (default: 3)
- `columns` (Number): Number of columns for table (default: 4)
- `fields` (Number): Number of fields for form (default: 4)
- `lines` (Number): Number of lines for text (default: 3)
- `height` (String): Custom height
- `width` (String): Custom width

## Composables

### useBatchErrors

Error handling utilities and message catalog.

**Usage:**

```vue
<script setup>
import { useBatchErrors } from '../../composables/useBatchErrors'

const { 
  ERROR_MESSAGES,
  ERROR_TYPES,
  parseError,
  getErrorMessage,
  isRetryableError,
  getErrorIcon,
  handleError,
  createErrorHandler
} = useBatchErrors()

// Parse error from API
const handleApiError = (error) => {
  const parsed = parseError(error)
  console.log(parsed.type) // 'network', 'validation', etc.
  console.log(parsed.message) // User-friendly message
  console.log(parsed.icon) // Appropriate icon
  console.log(parsed.isRetryable) // Whether error is retryable
}

// Get simple error message
const message = getErrorMessage(error, 'Default message')

// Check if retryable
if (isRetryableError(error)) {
  // Show retry button
}

// Create error handler for specific operation
const handleBatchCreate = createErrorHandler('Batch Creation')
try {
  await createBatch()
} catch (err) {
  const errorInfo = handleBatchCreate(err)
  // Use errorInfo.message, errorInfo.icon, etc.
}
</script>
```

**Error Messages Catalog:**

```javascript
ERROR_MESSAGES = {
  // Network
  NETWORK_ERROR: 'Không thể kết nối đến máy chủ...',
  TIMEOUT_ERROR: 'Yêu cầu hết thời gian chờ...',
  
  // Batch Definition
  DEFINITION_FETCH_ERROR: 'Không thể tải danh sách định nghĩa batch',
  DEFINITION_CREATE_ERROR: 'Không thể tạo định nghĩa batch',
  // ... more messages
}
```

**Error Types:**

```javascript
ERROR_TYPES = {
  NETWORK: 'network',
  VALIDATION: 'validation',
  PERMISSION: 'permission',
  NOT_FOUND: 'not_found',
  SERVER: 'server',
  UNKNOWN: 'unknown'
}
```

## Complete Example

Here's a complete example showing all error handling patterns:

```vue
<template>
  <div class="container">
    <!-- Loading State -->
    <LoadingSkeleton 
      v-if="loading && !data"
      type="list"
      :rows="5"
    />

    <!-- Error State (full page) -->
    <ErrorState
      v-else-if="error && !data"
      :icon="errorIcon"
      :title="errorTitle"
      :message="error"
      :retryable="isRetryable"
      :onRetry="loadData"
      showGoBack
      goBackRoute="/batch"
    />

    <!-- Empty State -->
    <EmptyState
      v-else-if="!loading && !data"
      icon="📦"
      title="Chưa có dữ liệu"
      description="Bạn chưa có batch nào"
      actionLabel="Tạo mới"
      :onAction="createNew"
    />

    <!-- Data Display -->
    <div v-else>
      <!-- Inline error for operations -->
      <InlineError
        v-if="operationError"
        :message="operationError"
        :showRetry="true"
        :onRetry="retryOperation"
        @dismiss="operationError = null"
      />

      <!-- Your content here -->
      <div v-for="item in data" :key="item.id">
        {{ item.name }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useBatchErrors } from '../../composables/useBatchErrors'
import ErrorState from './ErrorState.vue'
import EmptyState from './EmptyState.vue'
import InlineError from './InlineError.vue'
import LoadingSkeleton from './LoadingSkeleton.vue'

const { parseError, getErrorIcon, isRetryableError } = useBatchErrors()

const loading = ref(false)
const data = ref(null)
const error = ref(null)
const operationError = ref(null)

// Error handling computed properties
const parsedError = computed(() => {
  if (!error.value) return null
  return parseError({ response: { data: { message: error.value } } })
})

const errorIcon = computed(() => 
  parsedError.value ? getErrorIcon(parsedError.value.type) : '❌'
)

const errorTitle = computed(() => {
  if (!parsedError.value) return 'Lỗi'
  const titles = {
    network: 'Lỗi kết nối',
    validation: 'Dữ liệu không hợp lệ',
    permission: 'Không có quyền',
    not_found: 'Không tìm thấy',
    server: 'Lỗi máy chủ',
    unknown: 'Lỗi không xác định'
  }
  return titles[parsedError.value.type] || 'Lỗi'
})

const isRetryable = computed(() => {
  if (!parsedError.value) return true
  return ['network', 'server'].includes(parsedError.value.type)
})

// Load data
const loadData = async () => {
  loading.value = true
  error.value = null
  
  try {
    // API call
    const response = await fetch('/api/batch-records')
    data.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

// Retry operation
const retryOperation = async () => {
  operationError.value = null
  // Retry logic
}

// Create new
const createNew = () => {
  // Navigate to create form
}
</script>
```

## Best Practices

1. **Always show loading states**: Use LoadingSkeleton while data is loading
2. **Distinguish between error types**: Use ErrorState for critical errors, InlineError for form errors
3. **Provide retry for network errors**: Always allow users to retry network/server errors
4. **Show empty states**: Use EmptyState when there's no data, not just blank space
5. **Clear error messages**: Use the ERROR_MESSAGES catalog for consistent messaging
6. **Handle errors at the right level**: 
   - Store-level errors for data fetching
   - Component-level errors for user actions
7. **Provide context**: Include error details when helpful for debugging
8. **Allow dismissal**: Let users dismiss non-critical errors
9. **Test error states**: Ensure all error paths are tested and display correctly

## Error Flow Diagram

```
User Action
    ↓
Try Operation
    ↓
Success? ──Yes──→ Show Success State
    ↓
   No
    ↓
Parse Error (useBatchErrors)
    ↓
Determine Error Type
    ↓
Network/Server? ──Yes──→ Show ErrorState with Retry
    ↓
   No
    ↓
Validation? ──Yes──→ Show InlineError
    ↓
   No
    ↓
Permission? ──Yes──→ Show ErrorState, redirect to login
    ↓
   No
    ↓
Show Generic Error with appropriate message
```

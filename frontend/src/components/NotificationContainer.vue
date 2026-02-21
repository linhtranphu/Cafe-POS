<template>
  <div class="fixed top-4 right-4 z-50 space-y-2 max-w-md">
    <TransitionGroup name="notification">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="[
          'rounded-lg shadow-lg p-4 border-l-4 flex items-start gap-3',
          'transform transition-all duration-300',
          notificationClasses[notification.type]
        ]"
      >
        <!-- Icon -->
        <div class="flex-shrink-0 text-2xl">
          {{ notificationIcons[notification.type] }}
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <h4 v-if="notification.title" class="font-semibold mb-1">
            {{ notification.title }}
          </h4>
          <p class="text-sm">
            {{ notification.message }}
          </p>

          <!-- Action Button -->
          <button
            v-if="notification.action"
            @click="handleAction(notification)"
            class="mt-2 text-sm font-medium underline hover:no-underline"
          >
            {{ notification.action.label }}
          </button>
        </div>

        <!-- Dismiss Button -->
        <button
          v-if="notification.dismissible"
          @click="removeNotification(notification.id)"
          class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { useNotifications } from '../composables/useNotifications'

const { notifications, removeNotification } = useNotifications()

const notificationClasses = {
  success: 'bg-green-50 border-green-500 text-green-800',
  error: 'bg-red-50 border-red-500 text-red-800',
  warning: 'bg-yellow-50 border-yellow-500 text-yellow-800',
  info: 'bg-blue-50 border-blue-500 text-blue-800'
}

const notificationIcons = {
  success: '✅',
  error: '❌',
  warning: '⚠️',
  info: 'ℹ️'
}

const handleAction = (notification) => {
  if (notification.action && notification.action.handler) {
    notification.action.handler()
  }
  // Optionally dismiss after action
  removeNotification(notification.id)
}
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.8);
}

.notification-move {
  transition: transform 0.3s ease;
}
</style>

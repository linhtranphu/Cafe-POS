import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Composable for pull-to-refresh functionality
 * @param {Function} onRefresh - Callback function to execute on refresh
 * @param {Object} options - Configuration options
 * @returns {Object} - Reactive state and methods
 */
export function usePullToRefresh(onRefresh, options = {}) {
  const {
    threshold = 80, // Distance to pull before triggering refresh
    resistance = 2.5, // Resistance factor for pull distance
    maxPullDistance = 150 // Maximum pull distance
  } = options

  const isPulling = ref(false)
  const isRefreshing = ref(false)
  const pullDistance = ref(0)
  
  let startY = 0
  let currentY = 0
  let scrollTop = 0

  const handleTouchStart = (e) => {
    // Only start if at top of page
    scrollTop = window.pageYOffset || document.documentElement.scrollTop
    if (scrollTop === 0) {
      startY = e.touches[0].clientY
      isPulling.value = true
    }
  }

  const handleTouchMove = (e) => {
    if (!isPulling.value || isRefreshing.value) return

    currentY = e.touches[0].clientY
    const distance = currentY - startY

    // Only pull down (positive distance) and when at top
    if (distance > 0 && scrollTop === 0) {
      // Prevent default scroll behavior
      e.preventDefault()
      
      // Apply resistance to make it feel natural
      pullDistance.value = Math.min(
        distance / resistance,
        maxPullDistance
      )
    }
  }

  const handleTouchEnd = async () => {
    if (!isPulling.value) return

    isPulling.value = false

    // Trigger refresh if pulled beyond threshold
    if (pullDistance.value >= threshold) {
      isRefreshing.value = true
      
      try {
        await onRefresh()
      } catch (error) {
        console.error('Refresh failed:', error)
      } finally {
        // Delay to show animation
        setTimeout(() => {
          isRefreshing.value = false
          pullDistance.value = 0
        }, 500)
      }
    } else {
      // Reset if not pulled enough
      pullDistance.value = 0
    }
  }

  onMounted(() => {
    document.addEventListener('touchstart', handleTouchStart, { passive: true })
    document.addEventListener('touchmove', handleTouchMove, { passive: false })
    document.addEventListener('touchend', handleTouchEnd, { passive: true })
  })

  onUnmounted(() => {
    document.removeEventListener('touchstart', handleTouchStart)
    document.removeEventListener('touchmove', handleTouchMove)
    document.removeEventListener('touchend', handleTouchEnd)
  })

  return {
    isPulling,
    isRefreshing,
    pullDistance
  }
}

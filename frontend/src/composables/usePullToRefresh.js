import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Composable for pull-to-refresh functionality
 * Works with both page scroll and container scroll
 * @param {Function} onRefresh - Callback function to execute on refresh
 * @param {Object} options - Configuration options
 * @returns {Object} - Reactive state and methods
 */
export function usePullToRefresh(onRefresh, options = {}) {
  const {
    threshold = 80, // Distance to pull before triggering refresh
    resistance = 2.5, // Resistance factor for pull distance
    maxPullDistance = 150, // Maximum pull distance
    containerSelector = '.overflow-y-auto' // Selector for scroll container
  } = options

  const isPulling = ref(false)
  const isRefreshing = ref(false)
  const pullDistance = ref(0)
  
  let startY = 0
  let currentY = 0
  let scrollContainer = null

  /**
   * Get the scroll position - works for both page and container scroll
   */
  const getScrollTop = () => {
    // Try to find scroll container first (for container scroll)
    if (scrollContainer && scrollContainer.scrollTop !== undefined) {
      return scrollContainer.scrollTop
    }
    
    // Fallback to page scroll (for page scroll views)
    return window.pageYOffset || document.documentElement.scrollTop
  }

  /**
   * Find the scroll container element
   */
  const findScrollContainer = (element) => {
    if (!element) return null
    
    // Check if this element is scrollable
    const style = window.getComputedStyle(element)
    const isScrollable = style.overflowY === 'auto' || style.overflowY === 'scroll'
    
    if (isScrollable && element.scrollHeight > element.clientHeight) {
      return element
    }
    
    // Recursively check parent
    return findScrollContainer(element.parentElement)
  }

  const handleTouchStart = (e) => {
    // Find scroll container if not already found
    if (!scrollContainer) {
      scrollContainer = findScrollContainer(e.target)
    }
    
    // Only start if at top of scroll area
    const scrollTop = getScrollTop()
    if (scrollTop === 0) {
      startY = e.touches[0].clientY
      isPulling.value = true
    }
  }

  const handleTouchMove = (e) => {
    if (!isPulling.value || isRefreshing.value) return

    currentY = e.touches[0].clientY
    const distance = currentY - startY
    const scrollTop = getScrollTop()

    // Only pull down (positive distance) and when at top
    if (distance > 0 && scrollTop === 0) {
      // Prevent default scroll behavior
      e.preventDefault()
      
      // Apply resistance to make it feel natural
      pullDistance.value = Math.min(
        distance / resistance,
        maxPullDistance
      )
    } else {
      // Reset if scrolled away from top
      if (isPulling.value && scrollTop > 0) {
        isPulling.value = false
        pullDistance.value = 0
      }
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

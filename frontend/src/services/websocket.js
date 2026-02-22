import io from 'socket.io-client'

class WebSocketService {
  constructor() {
    this.socket = null
    this.connected = false
    this.listeners = new Map()
  }

  connect(token) {
    if (this.socket && this.connected) {
      console.log('[WebSocket] Already connected')
      return
    }

    // Get backend URL from environment or default to localhost
    const backendUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    
    console.log('[WebSocket] Connecting to:', backendUrl)

    this.socket = io(backendUrl, {
      auth: {
        token: token
      },
      transports: ['websocket', 'polling'],
      reconnection: true,
      reconnectionDelay: 1000,
      reconnectionDelayMax: 5000,
      reconnectionAttempts: 5
    })

    this.socket.on('connect', () => {
      console.log('[WebSocket] Connected')
      this.connected = true
    })

    this.socket.on('disconnect', (reason) => {
      console.log('[WebSocket] Disconnected:', reason)
      this.connected = false
    })

    this.socket.on('connect_error', (error) => {
      console.error('[WebSocket] Connection error:', error)
      this.connected = false
    })

    this.socket.on('reconnect', (attemptNumber) => {
      console.log('[WebSocket] Reconnected after', attemptNumber, 'attempts')
      this.connected = true
    })

    this.socket.on('reconnect_error', (error) => {
      console.error('[WebSocket] Reconnection error:', error)
    })

    this.socket.on('reconnect_failed', () => {
      console.error('[WebSocket] Reconnection failed')
      this.connected = false
    })
  }

  disconnect() {
    if (this.socket) {
      console.log('[WebSocket] Disconnecting')
      this.socket.disconnect()
      this.socket = null
      this.connected = false
      this.listeners.clear()
    }
  }

  on(event, callback) {
    if (!this.socket) {
      console.warn('[WebSocket] Cannot listen to event, not connected')
      return
    }

    // Store listener for cleanup
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event).push(callback)

    this.socket.on(event, callback)
    console.log('[WebSocket] Listening to event:', event)
  }

  off(event, callback) {
    if (!this.socket) {
      return
    }

    this.socket.off(event, callback)

    // Remove from listeners map
    if (this.listeners.has(event)) {
      const callbacks = this.listeners.get(event)
      const index = callbacks.indexOf(callback)
      if (index > -1) {
        callbacks.splice(index, 1)
      }
      if (callbacks.length === 0) {
        this.listeners.delete(event)
      }
    }
  }

  emit(event, data) {
    if (!this.socket || !this.connected) {
      console.warn('[WebSocket] Cannot emit event, not connected')
      return
    }

    this.socket.emit(event, data)
  }

  isConnected() {
    return this.connected
  }
}

// Export singleton instance
export const websocketService = new WebSocketService()

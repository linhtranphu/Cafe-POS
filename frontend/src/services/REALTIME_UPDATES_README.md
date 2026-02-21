# Real-Time Updates Implementation

## Overview

This implementation provides real-time updates for print jobs and printer status using WebSocket connections. The system automatically notifies users when print jobs are created, status changes occur, or printers go offline.

## Architecture

### Components

1. **WebSocket Service** (`services/websocket.js`)
   - Singleton service managing Socket.IO connection
   - Handles connection lifecycle (connect, disconnect, reconnect)
   - Provides event listener management

2. **Print Job WebSocket Composable** (`composables/usePrintJobWebSocket.js`)
   - Listens to print job events:
     - `print-job-created`: New print job created
     - `print-job-status-changed`: Job status updated
     - `print-job-failed`: Job failed
   - Updates print job store automatically
   - Triggers notifications for important events

3. **Printer Status Composable** (`composables/usePrinterStatus.js`)
   - Monitors printer status:
     - `printer-offline`: Printer disconnected
     - `printer-online`: Printer reconnected
     - `printer-error`: Hardware errors (paper out, jam, etc.)
   - Shows appropriate warnings and notifications

4. **Notification System** (`composables/useNotifications.js`)
   - Global notification state management
   - Toast-style notifications with auto-dismiss
   - Support for action buttons (e.g., "Retry")
   - Different notification types: success, error, warning, info

5. **Notification Container** (`components/NotificationContainer.vue`)
   - Visual component displaying toast notifications
   - Positioned at top-right of screen
   - Animated entrance/exit transitions
   - Dismissible notifications

## Usage

### In Components

```javascript
import { usePrintJobWebSocket } from '@/composables/usePrintJobWebSocket'
import { usePrinterStatus } from '@/composables/usePrinterStatus'
import { useNotifications } from '@/composables/useNotifications'

// Setup WebSocket listeners (automatically done in onMounted)
const { setupListeners } = usePrintJobWebSocket()

// Monitor printer status
const { printerStatuses, getPrinterStatus } = usePrinterStatus()

// Show custom notifications
const { showSuccess, showError, showWarning } = useNotifications()

// Example: Show custom notification
showSuccess('Operation completed successfully')
showError('Something went wrong', { duration: 0 }) // Don't auto-dismiss
```

### WebSocket Connection

The WebSocket connection is automatically established when the user logs in:

1. User authenticates → Token stored
2. `main.js` connects WebSocket with auth token
3. Connection maintained throughout session
4. Auto-reconnect on connection loss
5. Disconnect on logout

### Backend Events (To Be Implemented)

The backend needs to emit these events via Socket.IO:

```go
// Print job events
socket.Emit("print-job-created", map[string]interface{}{
    "job": job,
})

socket.Emit("print-job-status-changed", map[string]interface{}{
    "job_id": jobID,
    "status": status,
    "error_msg": errorMsg,
    "order_number": orderNumber,
    "type": jobType,
})

socket.Emit("print-job-failed", map[string]interface{}{
    "job_id": jobID,
    "order_number": orderNumber,
    "type": jobType,
    "error_msg": errorMsg,
})

// Printer status events
socket.Emit("printer-offline", map[string]interface{}{
    "printer_id": printerID,
    "printer_name": printerName,
})

socket.Emit("printer-online", map[string]interface{}{
    "printer_id": printerID,
    "printer_name": printerName,
})

socket.Emit("printer-error", map[string]interface{}{
    "printer_id": printerID,
    "printer_name": printerName,
    "error_type": errorType, // paper_out, paper_jam, cover_open, hardware_error
    "error_msg": errorMsg,
})
```

## Features

### Automatic Updates

- Print jobs list updates in real-time without manual refresh
- Status badges update automatically
- Pending/failed job counts update live

### Smart Notifications

- **Success**: Job completed successfully (auto-dismiss after 5s)
- **Error**: Job failed with retry button (stays until dismissed)
- **Warning**: Printer offline or hardware error (stays until dismissed)
- **Info**: General information (auto-dismiss after 5s)

### User Experience

- Non-intrusive toast notifications
- Action buttons for quick responses (e.g., "Retry")
- Visual feedback for all state changes
- Automatic cleanup of old notifications

## Configuration

### Environment Variables

```env
VITE_API_URL=http://localhost:8080
```

The WebSocket service uses this URL to connect to the backend.

### Notification Defaults

- Auto-dismiss duration: 5000ms (5 seconds)
- Error notifications: No auto-dismiss
- Warning notifications: No auto-dismiss
- Max notifications on screen: Unlimited (but auto-dismissed)

## Testing

To test the real-time updates:

1. Open the app in two browser windows
2. Create a print job in one window
3. Observe the update in the other window
4. Simulate printer offline event from backend
5. Check notification appears in all connected clients

## Future Enhancements

- [ ] Sound notifications for critical events
- [ ] Browser push notifications (when tab not active)
- [ ] Notification history/log
- [ ] User preferences for notification types
- [ ] Batch notification grouping (e.g., "5 jobs failed")
- [ ] Desktop notifications API integration

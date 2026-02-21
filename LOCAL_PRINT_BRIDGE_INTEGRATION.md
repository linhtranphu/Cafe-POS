# Local Print Bridge Integration

## Overview

The Local Print Bridge enables printing from a cloud-hosted POS system (EC2) to local thermal printers at the cafe. This document explains the architecture and integration.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      EC2 Cloud Server                        │
│  ┌──────────────┐         ┌──────────────┐                 │
│  │   Backend    │◄────────┤   Frontend   │                 │
│  │   (Go API)   │         │   (Vue.js)   │                 │
│  └──────┬───────┘         └──────┬───────┘                 │
│         │                         │                          │
│         │ WebSocket               │ HTTPS                    │
│         │ (Outbound)              │ (Outbound)               │
└─────────┼─────────────────────────┼──────────────────────────┘
          │                         │
          │                         │ Internet
          │                         │
          ▼                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Cafe (Local Network)                      │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Browser POS (Chrome/Firefox)             │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │  Vue.js Frontend (from EC2)                     │  │  │
│  │  │  - Receives WebSocket events                    │  │  │
│  │  │  - Detects local print bridge                   │  │  │
│  │  │  - Sends print requests to localhost:3001      │  │  │
│  │  └────────────────┬───────────────────────────────┘  │  │
│  └───────────────────┼──────────────────────────────────┘  │
│                      │ HTTP (localhost)                     │
│                      ▼                                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │     Local Print Bridge (Node.js Service)             │  │
│  │     Running on localhost:3001                        │  │
│  │     - Receives print requests from browser           │  │
│  │     - Converts to ESC/POS commands                   │  │
│  │     - Sends to printer via TCP                       │  │
│  │     - Updates backend status                         │  │
│  └──────────────────┬───────────────────────────────────┘  │
│                     │ TCP Socket (Port 9100)                │
│                     ▼                                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │     Thermal Printers (ESC/POS)                       │  │
│  │     - Bill Printer: 192.168.1.100:9100               │  │
│  │     - Label Printer: 192.168.1.101:9100              │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Flow Diagram

### 1. Order Creation → Auto Print

```
1. User creates order in browser
   ↓
2. Frontend sends order to Backend (EC2)
   ↓
3. Backend creates order + print jobs
   ↓
4. Backend emits WebSocket event: "print-job-created"
   ↓
5. Browser receives WebSocket event
   ↓
6. Frontend checks if local bridge is available
   ↓
7. If available: Frontend sends print request to localhost:3001
   ↓
8. Local Bridge receives request
   ↓
9. Local Bridge sends ESC/POS commands to printer via TCP
   ↓
10. Printer prints the bill/label
   ↓
11. Local Bridge updates backend: PUT /api/print-jobs/:id/status
   ↓
12. Backend updates job status to COMPLETED
   ↓
13. Backend emits WebSocket event: "print-job-status-changed"
   ↓
14. Browser updates UI
```

### 2. Manual Reprint

```
1. User clicks "Reprint" button
   ↓
2. Frontend calls backend: POST /api/orders/:id/reprint-bill
   ↓
3. Backend creates new print job
   ↓
4. Backend emits WebSocket event: "print-job-created"
   ↓
5. [Same as steps 5-14 above]
```

## Components

### 1. Backend (Go)

**New Endpoint:**
- `PUT /api/print-jobs/:id/status` - Receives status updates from local bridge

**Files Modified:**
- `backend/interfaces/http/print_job_handler.go` - Added UpdatePrintJobStatus handler
- `backend/main.go` - Registered new route

### 2. Frontend (Vue.js)

**New Files:**
- `frontend/src/services/localPrint.js` - Service to communicate with local bridge
- `frontend/src/composables/useLocalPrint.js` - Composable for managing local bridge

**Modified Files:**
- `frontend/src/stores/printJob.js` - Added WebSocket listeners and local print handling
- `frontend/src/views/PrintManagementView.vue` - Added local bridge status indicator

**Key Features:**
- Auto-detects local bridge availability
- Periodic health checks (every 30 seconds)
- Handles WebSocket events for print jobs
- Automatically sends print requests to local bridge when available

### 3. Local Print Bridge (Node.js)

**Location:** `local-print-bridge/`

**Key Files:**
- `src/index.js` - Express server
- `src/services/printerService.js` - TCP printer communication
- `src/services/backendSync.js` - Backend status updates

**Endpoints:**
- `GET /health` - Health check
- `POST /print` - Print request
- `POST /test-connection` - Test printer connection
- `GET /status` - Get statistics

## Setup Instructions

### 1. Install Local Print Bridge

On the cafe computer:

```bash
cd local-print-bridge
npm install
```

### 2. Configure Environment

Create `.env` file:

```env
PORT=3001
BACKEND_URL=https://your-ec2-domain.com
PRINTER_TIMEOUT=5000
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

### 3. Start Local Print Bridge

```bash
npm start
```

Or use PM2 for production:

```bash
npm install -g pm2
pm2 start src/index.js --name "print-bridge"
pm2 save
pm2 startup
```

### 4. Configure Printers in Backend

1. Open Print Management in the POS
2. Go to "Máy In" tab
3. Add printers with:
   - Name: "Bill Printer"
   - Type: "BILL"
   - Connection Type: "NETWORK"
   - IP: 192.168.1.100
   - Port: 9100
   - Set as default

### 5. Test the Integration

1. Open browser console (F12)
2. Check for: `[LocalPrint] Bridge available: true`
3. Create a test order
4. Verify print job appears in "Print Jobs" tab
5. Check printer output

## Network Requirements

### At the Cafe:

**No port forwarding needed!** The browser makes outbound connections only:
- WebSocket to EC2 (outbound, port 443)
- HTTP to localhost:3001 (local only)

### At EC2:

**Required open ports:**
- Port 443 (HTTPS/WSS) - For browser connections
- Port 80 (HTTP) - Optional, for redirect to HTTPS

## Troubleshooting

### Local Bridge Not Detected

1. Check if service is running:
   ```bash
   curl http://localhost:3001/health
   ```

2. Check browser console for errors

3. Verify CORS is enabled in local bridge

### Print Jobs Not Printing

1. Check printer IP is correct:
   ```bash
   ping 192.168.1.100
   ```

2. Test printer connection:
   ```bash
   cd local-print-bridge
   node src/test-printer.js 192.168.1.100 9100
   ```

3. Check local bridge logs:
   ```bash
   pm2 logs print-bridge
   ```

### Backend Status Not Updating

1. Verify BACKEND_URL in `.env`
2. Check backend logs for PUT requests
3. Verify authentication token (if required)

## Security Considerations

1. **Local Bridge runs on localhost only** - Not exposed to internet
2. **WebSocket uses WSS (encrypted)** - Secure communication from EC2
3. **No sensitive data in print content** - Already formatted ESC/POS
4. **Backend validates all status updates** - Prevents unauthorized changes

## Performance

- **Print latency:** ~500ms (WebSocket + TCP)
- **Health check interval:** 30 seconds
- **Printer timeout:** 5 seconds
- **Backend sync timeout:** 5 seconds

## Monitoring

### Local Bridge Statistics

Access via: `http://localhost:3001/status`

Returns:
```json
{
  "success": true,
  "stats": {
    "totalPrints": 150,
    "successfulPrints": 148,
    "failedPrints": 2,
    "successRate": "98.67%",
    "lastPrintTime": "2024-02-16T10:30:00Z"
  },
  "uptime": 86400,
  "timestamp": "2024-02-16T10:35:00Z"
}
```

### Frontend Status Indicator

- **Green dot + "Local Bridge Online"** - Bridge available
- **Gray dot + "Local Bridge Offline"** - Bridge not available

## Future Enhancements

1. **Printer status monitoring** - Real-time paper/error detection
2. **Print queue management** - Pause/resume printing
3. **Multiple cafe support** - Different bridges for different locations
4. **Print preview** - Show preview before printing
5. **Offline mode** - Queue prints when backend is unreachable

## Related Documentation

- [Local Print Bridge README](local-print-bridge/README.md)
- [Print System Design](.kiro/specs/order-printing/design.md)
- [WebSocket Integration](frontend/src/services/REALTIME_UPDATES_README.md)

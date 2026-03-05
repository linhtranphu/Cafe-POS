# Print System Architecture

## Overview

Hệ thống in của Cafe POS sử dụng kiến trúc WebSocket để kết nối realtime giữa Backend và Print Bridge.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         EC2 Server                               │
│                                                                   │
│  ┌──────────────┐         ┌──────────────┐                      │
│  │   Frontend   │         │   Backend    │                      │
│  │   (Nginx)    │         │   (Go)       │                      │
│  │   Port 80    │         │   Port 3000  │                      │
│  └──────┬───────┘         └──────┬───────┘                      │
│         │                        │                               │
│         │ HTTP API               │ Socket.IO                     │
│         │ /api/*                 │ /socket.io/*                  │
│         │                        │                               │
│         └────────────────────────┘                               │
│                  │                                                │
└──────────────────┼────────────────────────────────────────────────┘
                   │
                   │ WebSocket (wss://)
                   │ Event: print-job-created
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Local Network (Cafe)                          │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Print Bridge (Node.js)                       │   │
│  │              Port 3001                                    │   │
│  │                                                           │   │
│  │  ┌─────────────────────────────────────────────────┐    │   │
│  │  │  WebSocket Client                               │    │   │
│  │  │  - Connect to backend                           │    │   │
│  │  │  - Listen: print-job-created                    │    │   │
│  │  │  - Auto reconnect                               │    │   │
│  │  └─────────────────────────────────────────────────┘    │   │
│  │                                                           │   │
│  │  ┌─────────────────────────────────────────────────┐    │   │
│  │  │  Print Job Handler                              │    │   │
│  │  │  - Receive job from WebSocket                   │    │   │
│  │  │  - Send ESC/POS to printer                      │    │   │
│  │  │  - Update status to backend                     │    │   │
│  │  └─────────────────────────────────────────────────┘    │   │
│  │                                                           │   │
│  └──────────────────┬───────────────────┬────────────────────┘   │
│                     │                   │                         │
│                     │ ESC/POS           │ ESC/POS                 │
│                     │ Port 9100         │ Port 9100               │
│                     ▼                   ▼                         │
│            ┌─────────────┐     ┌─────────────┐                   │
│            │   Printer   │     │   Printer   │                   │
│            │   (Bill)    │     │   (Label)   │                   │
│            │ 192.168.1.x │     │ 192.168.1.y │                   │
│            └─────────────┘     └─────────────┘                   │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Order Creation → Print

```
User (Browser)
    │
    │ 1. POST /api/orders
    ▼
Backend (EC2)
    │
    ├─→ 2. Save order to MongoDB
    │
    ├─→ 3. Create print job
    │
    └─→ 4. Emit WebSocket event
        │
        │ Event: print-job-created
        │ Data: {
        │   job: {
        │     id: "...",
        │     content: "ESC/POS commands",
        │     printer_ip: "192.168.1.100",
        │     printer_port: 9100,
        │     type: "BILL"
        │   }
        │ }
        │
        ▼
Print Bridge (Local)
    │
    ├─→ 5. Receive WebSocket event
    │
    ├─→ 6. Send ESC/POS to printer
    │       (TCP socket to 192.168.1.100:9100)
    │
    └─→ 7. Update status to backend
        │
        │ PUT /api/print-jobs/:id/status
        │ Body: { status: "COMPLETED" }
        │
        ▼
Backend (EC2)
    │
    └─→ 8. Update print job status in MongoDB
```

### 2. Print Job Status Updates

```
Print Bridge
    │
    │ Success case:
    ├─→ PUT /api/print-jobs/:id/status
    │   Body: { status: "COMPLETED" }
    │
    │ Failure case:
    └─→ PUT /api/print-jobs/:id/status
        Body: { 
          status: "FAILED",
          error_msg: "Printer offline"
        }
        │
        ▼
Backend
    │
    ├─→ Update MongoDB
    │
    └─→ Emit WebSocket event
        │
        │ Event: print-job-status-changed
        │ Data: {
        │   job_id: "...",
        │   status: "COMPLETED",
        │   error_msg: ""
        │ }
        │
        ▼
Frontend (Browser)
    │
    └─→ Show notification to user
```

## Components

### Backend (Go)

**Location:** `backend/infrastructure/websocket/`

**Responsibilities:**
- Socket.IO server on port 3000
- Emit `print-job-created` when order created
- Emit `print-job-status-changed` when status updated
- HTTP API for status updates

**Key Files:**
- `socketio_server.go` - Socket.IO server setup
- `broadcaster.go` - Event broadcasting logic
- `socketio_handler.go` - Socket.IO event handlers

### Print Bridge (Node.js)

**Location:** `local-print-bridge/`

**Responsibilities:**
- Connect to backend via WebSocket
- Listen for print job events
- Send ESC/POS commands to printers
- Update job status back to backend

**Key Files:**
- `src/services/websocketClient.js` - WebSocket connection
- `src/services/printJobHandler.js` - Job processing
- `src/services/printerService.js` - ESC/POS printing
- `src/services/backendSync.js` - Status updates

### Frontend (Vue.js)

**Location:** `frontend/src/services/websocket.js`

**Responsibilities:**
- Connect to backend for UI updates
- Show realtime notifications
- Update order status in UI

**Note:** Frontend does NOT communicate with Print Bridge directly

## Network Requirements

### EC2 Server

**Inbound Rules:**
- Port 80 (HTTP) - Frontend access
- Port 443 (HTTPS) - Frontend access (recommended)
- Port 3000 (TCP) - Backend API & WebSocket
  - Allow from Print Bridge IP
  - Or allow from 0.0.0.0/0 if dynamic IP

**Outbound Rules:**
- All traffic (default)

### Print Bridge Machine

**Inbound Rules:**
- Port 3001 (HTTP) - Health check endpoint (optional)

**Outbound Rules:**
- Port 3000 (TCP) - Connect to backend
- Port 443 (HTTPS) - Connect to backend (if using HTTPS)
- Port 9100 (TCP) - Connect to printers

### Printers

**Inbound Rules:**
- Port 9100 (TCP) - Receive print jobs from Print Bridge

## Configuration

### Backend

```bash
# .env.ec2
PORT=3000
MONGODB_URI=mongodb://...
```

No additional config needed for WebSocket.

### Print Bridge

```bash
# local-print-bridge/.env
PORT=3001
BACKEND_URL=https://tacafe.store
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_BILL_PRINTER_PORT=9100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
DEFAULT_LABEL_PRINTER_PORT=9100
```

### Frontend

```bash
# frontend/.env.production
VITE_API_URL=https://tacafe.store
```

Frontend connects to backend, NOT to Print Bridge.

## Advantages of This Architecture

### ✅ No CORS Issues
- Print Bridge is server-side, not browser
- No "private address space" restrictions

### ✅ Reliable Connection
- WebSocket auto-reconnect
- Print Bridge keeps trying even if backend restarts

### ✅ Secure
- Print Bridge can be behind firewall
- Only needs outbound connection to backend
- Printers stay in local network

### ✅ Scalable
- Multiple Print Bridges can connect to same backend
- Each cafe location has its own Print Bridge
- Backend handles all coordination

### ✅ Simple Frontend
- Frontend doesn't need to know about printers
- No printer configuration in browser
- Works from any device

## Deployment Checklist

### Backend (EC2)

- [x] Socket.IO server running on port 3000
- [x] Nginx proxy for `/socket.io/` path
- [ ] Security group allows port 3000 from Print Bridge IP
- [ ] SSL certificate installed (recommended)

### Print Bridge (Local)

- [ ] Node.js installed
- [ ] `.env` configured with BACKEND_URL
- [ ] Printer IPs configured
- [ ] Can reach backend (test with curl)
- [ ] WebSocket connection successful
- [ ] PM2 or systemd for auto-start

### Network

- [ ] Print Bridge can reach backend:3000
- [ ] Print Bridge can reach printers:9100
- [ ] Firewall allows outbound connections

## Troubleshooting

See `PRINT_BRIDGE_WEBSOCKET_SETUP.md` for detailed troubleshooting guide.

## Testing

```bash
# Test backend WebSocket
cd local-print-bridge
node test-backend-websocket.js https://tacafe.store

# Test printer connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100"}'

# Create test order
# → Check Print Bridge logs for "print-job-created" event
```

## Related Documentation

- `PRINT_BRIDGE_WEBSOCKET_SETUP.md` - Setup guide
- `WEBSOCKET_EC2_FIX.md` - WebSocket troubleshooting
- `local-print-bridge/README.md` - Print Bridge documentation

# WebSocket Test Results

## Test Date
2026-02-21

## Setup
- Backend: Running on port 3000 ✅
- Print Bridge: Running locally with npm start ✅
- MongoDB: Running with replica set ✅

## Test Results

### ✅ Backend WebSocket Server
- Endpoint: `http://localhost:3000/socket.io/`
- Status: **WORKING**
- Evidence:
  ```
  curl "http://localhost:3000/socket.io/?EIO=4&transport=polling"
  Response: 0{"maxPayload":1000000,"pingInterval":25000,"pingTimeout":60000,"sid":"...","upgrades":["websocket"]}
  ```
- Backend logs show connections:
  ```
  [WebSocket] New connection established: fd774421-7e6a-4a70-8497-5ec1e4929aba
  [WebSocket] Client registered: fd774421-7e6a-4a70-8497-5ec1e4929aba (User: anonymous)
  [WebSocket] Total clients: 1
  [WebSocket] Sent Socket.IO CONNECT packet to client
  ```

### ❌ Print Bridge WebSocket Client
- Status: **CONNECTION TIMEOUT**
- Error: Client không nhận được 'connect' event
- Logs:
  ```
  [WebSocket] Connecting to: http://localhost:3000
  [WebSocket] Connection error (attempt 1/10): timeout
  ```

## Problem Analysis

### What's Working
1. ✅ Backend Socket.IO server starts successfully
2. ✅ Backend accepts WebSocket connections
3. ✅ Backend sends Socket.IO CONNECT packet ("40")
4. ✅ Backend registers clients

### What's Not Working
1. ❌ Client doesn't receive 'connect' event
2. ❌ Connection times out after 20 seconds
3. ❌ Client disconnects immediately after connecting

### Root Cause
**Socket.IO Protocol Mismatch or Event Handling Issue**

Backend logs show:
```
[WebSocket] Sent Socket.IO CONNECT packet to client fd774421-7e6a-4a70-8497-5ec1e4929aba
[WebSocket] Read error: websocket: close 1005 (no status)
[WebSocket] Client unregistered
```

This suggests:
- Backend sends CONNECT packet
- Client receives it but doesn't recognize it
- Client closes connection (no status code)
- Backend sees it as "Read error"

### Possible Causes

1. **Socket.IO Version Mismatch**
   - Backend: Custom Socket.IO v4 implementation in Go
   - Client: socket.io-client v4.7.2
   - Protocol might not be 100% compatible

2. **Event Format Issue**
   - Backend sends: `"40"` (CONNECT packet)
   - Client expects: Specific format that triggers 'connect' event
   - Mismatch in packet format

3. **Namespace Issue**
   - Backend uses default namespace
   - Client might be looking for different namespace

4. **Timeout Too Short**
   - 20 second timeout might not be enough
   - But backend shows connection within seconds

## Attempted Fixes

### 1. Changed Transport Order
```javascript
transports: ['polling', 'websocket']  // ❌ Failed - xhr post error
transports: ['websocket']              // ❌ Failed - timeout
```

### 2. Increased Timeout
```javascript
timeout: 10000  // Original
timeout: 20000  // ❌ Still timeout
```

### 3. Added forceNew Flag
```javascript
forceNew: true  // ❌ Still timeout
```

## Recommendations

### Option 1: Fix Socket.IO Implementation (Recommended)
Backend's custom Socket.IO implementation might not be fully compatible. Need to:
1. Review backend Socket.IO packet format
2. Ensure CONNECT packet ("40") is sent correctly
3. Check if client needs additional handshake steps

### Option 2: Use Standard Socket.IO Library
Replace custom Go implementation with standard Socket.IO library:
- Use `github.com/googollee/go-socket.io`
- This ensures full compatibility with socket.io-client

### Option 3: Use Simple WebSocket (No Socket.IO)
Remove Socket.IO layer and use raw WebSocket:
- Backend: Use standard `gorilla/websocket`
- Client: Use standard `ws` library
- Simpler, more reliable, but lose Socket.IO features

### Option 4: HTTP Polling (Fallback)
Keep current HTTP polling approach:
- Print Bridge polls `/api/manager/print-jobs/pending`
- Works but less efficient than WebSocket
- No real-time push

## Next Steps

### Immediate (Testing)
1. ✅ Verify backend Socket.IO implementation
2. ✅ Test with simple Socket.IO client
3. ⏳ Debug packet format
4. ⏳ Check namespace handling

### Short Term (Fix)
1. Fix Socket.IO CONNECT packet format
2. Or switch to standard Socket.IO library
3. Or implement simple WebSocket without Socket.IO

### Long Term (Production)
1. Decide on WebSocket vs Polling
2. If WebSocket: Ensure full Socket.IO compatibility
3. If Polling: Optimize polling interval
4. Add comprehensive error handling

## Conclusion

WebSocket infrastructure is in place but Socket.IO protocol compatibility issue prevents connection. Backend accepts connections but client doesn't recognize CONNECT event.

**Status**: 🟡 Partially Working (Backend OK, Client Timeout)

**Recommendation**: Review backend Socket.IO implementation or switch to standard library.

## Files Involved

- `backend/infrastructure/websocket/socketio_handler.go` - Backend Socket.IO implementation
- `local-print-bridge/src/services/websocketClient.js` - Client implementation
- `backend/infrastructure/websocket/hub.go` - WebSocket hub
- `backend/infrastructure/websocket/broadcaster.go` - Event broadcaster

## Test Commands

```bash
# Start Print Bridge
cd local-print-bridge && npm start

# Watch logs
tail -f backend.log | grep WebSocket

# Test backend endpoint
curl "http://localhost:3000/socket.io/?EIO=4&transport=polling"

# Test with simple client
node test-websocket-connection.js
```

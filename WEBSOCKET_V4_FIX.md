# WebSocket Socket.IO v4 Compatibility Fix

## Problem

Frontend Socket.IO client v4.8.3 was unable to connect to backend WebSocket server, showing error:

```
Error: It seems you are trying to reach a Socket.IO server in v2.x with a v3.x client, 
but they are not compatible
```

## Root Cause

Socket.IO v4 requires a two-step connection process:
1. **HTTP Polling Handshake** - Initial connection via HTTP GET request
2. **WebSocket Upgrade** - Upgrade to WebSocket after handshake

Our initial implementation only supported direct WebSocket connections (v2 style), which is incompatible with Socket.IO v4 clients.

## Solution

Updated `backend/infrastructure/websocket/socketio_handler.go` to support both:

### 1. HTTP Polling Handshake

When client first connects, it sends HTTP GET request:
```
GET /socket.io/?EIO=4&transport=polling
```

Backend responds with handshake packet:
```
0{"sid":"xxx","upgrades":["websocket"],"pingInterval":25000,"pingTimeout":60000}
```

### 2. WebSocket Upgrade

After handshake, client upgrades to WebSocket:
```
GET /socket.io/?EIO=4&transport=websocket&sid=xxx
Upgrade: websocket
```

Backend upgrades connection and starts bidirectional communication.

## Code Changes

**File:** `backend/infrastructure/websocket/socketio_handler.go`

**Changes:**
1. Split `HandleSocketIO` into two handlers:
   - `handleHTTPPolling()` - Handles initial handshake
   - `handleWebSocketUpgrade()` - Handles WebSocket upgrade

2. Added EIO version check (must be "4")

3. Added proper handshake response format

4. Updated protocol to Socket.IO v4 / Engine.IO v4

## Testing

### Before Fix:
```
[WebSocket] Connecting to: http://localhost:3000
[WebSocket] Connection error: TransportError: websocket error
```

### After Fix:
```
[WebSocket] Connecting to: http://localhost:3000
[WebSocket] Connected
[WebSocket] Client registered: xxx
```

## Socket.IO v4 Protocol Flow

```
1. Client → Server: GET /socket.io/?EIO=4&transport=polling
   ↓
2. Server → Client: 0{"sid":"xxx","upgrades":["websocket"],...}
   ↓
3. Client → Server: GET /socket.io/?EIO=4&transport=websocket&sid=xxx
                    Upgrade: websocket
   ↓
4. Server upgrades to WebSocket
   ↓
5. Server → Client: 0{"sid":"xxx",...} (Engine.IO OPEN)
   ↓
6. Server → Client: 40 (Socket.IO CONNECT)
   ↓
7. Bidirectional communication established
   ↓
8. Server → Client: 42["event-name",{data}] (Socket.IO EVENT)
```

## Compatibility

- ✅ Socket.IO v4.x clients (current: v4.8.3)
- ✅ Engine.IO v4
- ✅ Backward compatible with existing code
- ✅ No breaking changes to API

## Files Modified

1. `backend/infrastructure/websocket/socketio_handler.go`
   - Added HTTP polling support
   - Split handler into two functions
   - Added EIO version validation

## Verification

To verify the fix works:

1. **Start Backend:**
   ```bash
   cd backend
   go run main.go
   ```

2. **Start Frontend:**
   ```bash
   cd frontend
   npm run dev
   ```

3. **Check Browser Console:**
   - Should see: `[WebSocket] Connected`
   - No more version compatibility errors

4. **Test Event Broadcasting:**
   - Create an order
   - Check console for: `[WebSocket] Received event: print-job-created`

## Additional Notes

### Socket.IO Version History

- **v2.x:** Direct WebSocket connection
- **v3.x:** Introduced breaking changes, requires handshake
- **v4.x:** Current version, requires HTTP polling handshake first

### Why HTTP Polling First?

Socket.IO v4 uses HTTP polling for initial handshake to:
1. Establish session ID
2. Negotiate protocol version
3. Exchange connection parameters
4. Support environments where WebSocket is blocked

After handshake, it upgrades to WebSocket for better performance.

### Production Considerations

For production deployment:

1. **Use WSS (WebSocket Secure):**
   ```javascript
   const backendUrl = 'https://api.yourdomain.com'
   ```

2. **Configure Reverse Proxy:**
   ```nginx
   location /socket.io/ {
       proxy_pass http://backend:3000;
       proxy_http_version 1.1;
       proxy_set_header Upgrade $http_upgrade;
       proxy_set_header Connection "upgrade";
       proxy_set_header Host $host;
   }
   ```

3. **Enable CORS Properly:**
   ```go
   CheckOrigin: func(r *http.Request) bool {
       origin := r.Header.Get("Origin")
       return origin == "https://yourdomain.com"
   }
   ```

## Related Documentation

- [WebSocket Implementation Summary](WEBSOCKET_IMPLEMENTATION_SUMMARY.md)
- [Socket.IO v4 Migration Guide](https://socket.io/docs/v4/migrating-from-3-x-to-4-0/)
- [Engine.IO Protocol](https://socket.io/docs/v4/engine-io-protocol/)


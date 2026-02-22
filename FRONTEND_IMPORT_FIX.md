# Frontend Socket.IO Import Fix

## Lỗi

```
Uncaught TypeError: io is not a function
at WebSocketService.connect (websocket.js:21:19)
```

## Nguyên Nhân

Socket.IO client v2.x và v4.x có cách export khác nhau:

### v4.x (Old - Không hoạt động với v2.x)
```javascript
import { io } from 'socket.io-client'  // Named export
```

### v2.x (New - Correct)
```javascript
import io from 'socket.io-client'  // Default export
```

## ✅ Đã Sửa

File: `frontend/src/services/websocket.js`

**Before:**
```javascript
import { io } from 'socket.io-client'  // ❌ Named export
```

**After:**
```javascript
import io from 'socket.io-client'  // ✅ Default export
```

## Verify

1. **Clear browser cache** (Ctrl+Shift+R hoặc Cmd+Shift+R)
2. **Reload page**: http://localhost:5173
3. **Check console**: Không còn lỗi "io is not a function"
4. **Verify connection**: Should see `[WebSocket] Connected`

## Version Compatibility Summary

| Component | Version | Import Style | Status |
|-----------|---------|--------------|--------|
| Backend | go-socket.io v1.7.0 | N/A | ✅ |
| Print Bridge | socket.io-client v2.5.0 | `const io = require(...)` | ✅ |
| Frontend | socket.io-client v2.5.0 | `import io from ...` | ✅ |

All using Engine.IO v3 - 100% Compatible!

## Files Changed

1. `frontend/package.json` - Downgraded to v2.5.0
2. `frontend/src/services/websocket.js` - Fixed import statement
3. `frontend/node_modules/` - Updated dependencies

## Next Steps

1. ✅ Version downgraded
2. ✅ Import fixed
3. ⏳ Clear browser cache and reload
4. ⏳ Test WebSocket connection
5. ⏳ Test end-to-end order → print flow

## Testing

### Check WebSocket Connection
```javascript
// In browser console
console.log(websocketService.isConnected())
// Should return: true
```

### Test Print Job Flow
1. Login: admin/admin123
2. Create new order
3. Check Print Bridge logs:
   ```
   [WebSocket] 📨 New print job received
   [PrintJobHandler] Processing job...
   [PrintJobHandler] ✅ Job printed successfully
   ```

## Troubleshooting

### Still seeing "io is not a function"?
1. Hard refresh: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
2. Clear browser cache completely
3. Restart Vite dev server:
   ```bash
   kill $(lsof -t -i:5173)
   cd frontend && npm run dev -- --host
   ```

### WebSocket not connecting?
1. Check backend is running: `curl http://localhost:3000/health`
2. Check Print Bridge: `curl http://localhost:3001/health`
3. Check browser console for errors
4. Verify backend URL in `.env`:
   ```
   VITE_API_URL=http://localhost:3000
   ```

---

**Status**: ✅ Import fixed, ready to test after browser cache clear

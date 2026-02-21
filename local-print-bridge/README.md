# Local Print Bridge

Local print service for Cafe POS system. Bridges browser-based POS to thermal printers via TCP/IP.

## Architecture

```
Browser POS (EC2 Frontend)
         ↓
   WebSocket notification
         ↓
Local Print Bridge (localhost:3001)
         ↓
   TCP Socket (port 9100)
         ↓
Thermal Printer (192.168.x.x)
```

## Features

- ✅ HTTP REST API for print requests
- ✅ TCP socket communication with ESC/POS printers
- ✅ Automatic status sync with backend
- ✅ Connection testing
- ✅ Error handling and retry support
- ✅ Logging and statistics
- ✅ CORS enabled for browser access

## Installation

### Option 1: Docker (Recommended)

**Fastest and easiest way to deploy!**

```bash
# 1. Configure
cp .env.docker .env
nano .env  # Update BACKEND_URL and printer IPs

# 2. Start with one command
./docker-start.sh

# 3. Done! Service is running
```

See [Docker Deployment Guide](../LOCAL_PRINT_BRIDGE_DOCKER_GUIDE.md) for details.

**Advantages:**
- ✅ No Node.js installation needed
- ✅ Auto-restart on system boot
- ✅ Easy management with docker-compose
- ✅ Isolated environment

### Option 2: Node.js (Traditional)

**Prerequisites:**
- Node.js 16+ installed
- Thermal printer connected to local network
- Printer IP address and port (default: 9100)

**Setup:**

1. **Install dependencies:**
   ```bash
   cd local-print-bridge
   npm install
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Test printer connection:**
   ```bash
   npm run test -- 192.168.1.100 9100
   ```

4. **Start the service:**
   ```bash
   npm start
   ```

   For development with auto-reload:
   ```bash
   npm run dev
   ```

## Configuration

Edit `.env` file:

```env
# Server Port
PORT=3001

# Backend API URL (EC2)
BACKEND_URL=https://api.your-domain.com

# Default Printer IPs
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_BILL_PRINTER_PORT=9100

DEFAULT_LABEL_PRINTER_IP=192.168.1.101
DEFAULT_LABEL_PRINTER_PORT=9100

# Logging
LOG_LEVEL=info

# Connection Timeout (ms)
PRINTER_TIMEOUT=5000
```

## API Endpoints

### Health Check
```http
GET /health
```

Response:
```json
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0",
  "timestamp": "2024-02-16T10:30:00.000Z"
}
```

### Print
```http
POST /print
Content-Type: application/json

{
  "jobId": "65cf1234567890abcdef1234",
  "content": "ESC/POS formatted content...",
  "printerIP": "192.168.1.100",
  "printerPort": 9100,
  "type": "BILL"
}
```

Response (Success):
```json
{
  "success": true,
  "jobId": "65cf1234567890abcdef1234",
  "message": "Print completed successfully",
  "timestamp": "2024-02-16T10:30:00.000Z"
}
```

Response (Error):
```json
{
  "success": false,
  "error": "Printer offline or unreachable at 192.168.1.100:9100",
  "jobId": "65cf1234567890abcdef1234",
  "timestamp": "2024-02-16T10:30:00.000Z"
}
```

### Test Connection
```http
POST /test-connection
Content-Type: application/json

{
  "printerIP": "192.168.1.100",
  "printerPort": 9100
}
```

Response:
```json
{
  "success": true,
  "message": "Printer connection successful",
  "printer": "192.168.1.100:9100"
}
```

### Get Status
```http
GET /status
```

Response:
```json
{
  "success": true,
  "stats": {
    "totalPrints": 150,
    "successfulPrints": 148,
    "failedPrints": 2,
    "lastPrintTime": "2024-02-16T10:30:00.000Z",
    "successRate": "98.67%"
  },
  "uptime": 3600,
  "timestamp": "2024-02-16T10:30:00.000Z"
}
```

## Testing

### Test Printer Connection

```bash
npm run test -- <printer-ip> [port]

# Example:
npm run test -- 192.168.1.100 9100
```

This will:
1. Test TCP connection to printer
2. Send a test print
3. Display statistics

### Manual Test with curl

```bash
# Health check
curl http://localhost:3001/health

# Test connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP":"192.168.1.100","printerPort":9100}'

# Test print
curl -X POST http://localhost:3001/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId":"test-123",
    "content":"Test Print\n\n\n",
    "printerIP":"192.168.1.100",
    "printerPort":9100
  }'
```

## Troubleshooting

### Printer Not Found

**Error:** `Printer offline or unreachable at 192.168.1.100:9100`

**Solutions:**
1. Check if printer is powered on
2. Verify printer IP address: `ping 192.168.1.100`
3. Ensure printer is on the same network
4. Check printer network settings
5. Try accessing printer web interface: `http://192.168.1.100`

### Connection Timeout

**Error:** `Connection timeout to 192.168.1.100:9100`

**Solutions:**
1. Check network connectivity
2. Verify firewall settings
3. Increase timeout in `.env`: `PRINTER_TIMEOUT=10000`
4. Check if port 9100 is open on printer

### Port Already in Use

**Error:** `EADDRINUSE: address already in use :::3001`

**Solutions:**
1. Change port in `.env`: `PORT=3002`
2. Kill existing process: `lsof -ti:3001 | xargs kill -9`

### Backend Update Failed

**Error:** `Failed to update backend for job xxx`

**Solutions:**
1. Check `BACKEND_URL` in `.env`
2. Verify backend is accessible
3. Check backend API endpoint exists
4. Print will still succeed, only status update fails

## Deployment

### Run as Background Service (macOS/Linux)

Create a systemd service file (Linux) or launchd plist (macOS):

**Linux (systemd):**
```bash
sudo nano /etc/systemd/system/print-bridge.service
```

```ini
[Unit]
Description=Local Print Bridge
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/local-print-bridge
ExecStart=/usr/bin/node src/index.js
Restart=always
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable print-bridge
sudo systemctl start print-bridge
sudo systemctl status print-bridge
```

### Run on Startup (Windows)

Use `pm2` or `node-windows`:

```bash
npm install -g pm2
pm2 start src/index.js --name print-bridge
pm2 startup
pm2 save
```

## Integration with Frontend

See `frontend/src/services/localPrint.js` for integration example.

## License

MIT

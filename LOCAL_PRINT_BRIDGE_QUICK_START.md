# Local Print Bridge - Quick Start Guide

## What is the Local Print Bridge?

The Local Print Bridge is a small Node.js service that runs on the cafe computer and enables printing from your cloud-hosted POS system to local thermal printers.

## Prerequisites

- Node.js 16+ installed on cafe computer
- Thermal printers connected to local network (ESC/POS compatible)
- Printer IP addresses (e.g., 192.168.1.100)

## Installation (5 minutes)

### Step 1: Install Node.js

**On Windows:**
1. Download from https://nodejs.org/
2. Run installer
3. Verify: Open Command Prompt and run `node --version`

**On macOS:**
```bash
brew install node
```

**On Linux:**
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### Step 2: Install Local Print Bridge

```bash
# Navigate to the project folder
cd local-print-bridge

# Install dependencies
npm install
```

### Step 3: Configure

Create a `.env` file in `local-print-bridge/` folder:

```env
PORT=3001
BACKEND_URL=https://your-ec2-domain.com
PRINTER_TIMEOUT=5000
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

**Important:** Replace `your-ec2-domain.com` with your actual EC2 domain!

### Step 4: Test Printer Connection

```bash
# Test if printer is reachable
node src/test-printer.js 192.168.1.100 9100
```

You should see: `✅ Printer 192.168.1.100:9100 is online`

### Step 5: Start the Service

**For testing:**
```bash
npm start
```

**For production (auto-restart):**
```bash
# Install PM2
npm install -g pm2

# Start service
pm2 start src/index.js --name "print-bridge"

# Save configuration
pm2 save

# Setup auto-start on boot
pm2 startup
```

### Step 6: Verify

1. Open browser: http://localhost:3001/health
2. You should see: `{"status":"ok","service":"Local Print Bridge",...}`

## Configure Printers in POS

1. Open your POS in browser
2. Go to **Print Management** (🖨️ icon in menu)
3. Click **Máy In** tab
4. Click **+ Thêm Máy In**
5. Fill in:
   - **Tên:** Bill Printer
   - **Loại:** BILL
   - **Kết nối:** NETWORK
   - **IP:** 192.168.1.100
   - **Port:** 9100
   - **Độ rộng giấy:** 80mm
   - **Mặc định:** ✓ (check this)
6. Click **Test Connection** to verify
7. Click **Lưu**

Repeat for Label Printer if needed.

## Test the Integration

1. Open POS in browser
2. Check header in Print Management - should show **"Local Bridge Online"** with green dot
3. Create a test order
4. Go to **Print Jobs** tab
5. You should see the print job with status "COMPLETED"
6. Check your printer - bill should be printed!

## Troubleshooting

### "Local Bridge Offline" in POS

**Check if service is running:**
```bash
# If using npm start
# Check terminal - should show "Ready to accept print requests!"

# If using PM2
pm2 status
pm2 logs print-bridge
```

**Restart service:**
```bash
# If using npm start
# Press Ctrl+C and run: npm start

# If using PM2
pm2 restart print-bridge
```

### Printer Not Printing

**1. Check printer IP:**
```bash
ping 192.168.1.100
```

**2. Test printer connection:**
```bash
node src/test-printer.js 192.168.1.100 9100
```

**3. Check printer is on and has paper**

**4. Check printer port (usually 9100 for ESC/POS)**

### Print Jobs Stuck in "PENDING"

**Check local bridge logs:**
```bash
# If using PM2
pm2 logs print-bridge --lines 50

# If using npm start
# Check terminal output
```

**Common issues:**
- Wrong printer IP
- Printer offline
- Network firewall blocking port 9100
- Local bridge not running

### Backend Status Not Updating

**Check BACKEND_URL in .env:**
```bash
cat local-print-bridge/.env
```

Should be: `BACKEND_URL=https://your-actual-domain.com`

**Test backend connection:**
```bash
curl https://your-actual-domain.com/health
```

## Daily Operations

### Starting the Service

**If using PM2 (recommended):**
- Service starts automatically on boot
- No manual action needed

**If using npm start:**
```bash
cd local-print-bridge
npm start
```

### Stopping the Service

**If using PM2:**
```bash
pm2 stop print-bridge
```

**If using npm start:**
- Press `Ctrl+C` in terminal

### Checking Status

**If using PM2:**
```bash
pm2 status
pm2 logs print-bridge
```

**If using npm start:**
- Check terminal output

### Viewing Statistics

Open in browser: http://localhost:3001/status

Shows:
- Total prints
- Success rate
- Last print time
- Uptime

## Network Setup

### Required Network Configuration

**At the Cafe:**
- ✅ Printers on same network as computer (e.g., 192.168.1.x)
- ✅ Computer can access internet (for backend communication)
- ❌ NO port forwarding needed
- ❌ NO firewall changes needed

**At EC2:**
- ✅ Port 443 open (HTTPS/WebSocket)
- ❌ NO special configuration needed

### Network Diagram

```
Internet
   │
   │ HTTPS/WSS (Outbound from browser)
   ▼
EC2 Server
   │
   │ WebSocket Events
   ▼
Browser (Cafe)
   │
   │ HTTP (localhost only)
   ▼
Local Print Bridge (localhost:3001)
   │
   │ TCP Port 9100
   ▼
Thermal Printers (192.168.1.x)
```

## Security Notes

1. **Local Bridge runs on localhost only** - Not accessible from internet
2. **No sensitive data** - Only ESC/POS print commands
3. **WebSocket is encrypted** - Uses WSS (HTTPS)
4. **No authentication needed** - Local bridge trusts localhost

## Performance

- **Print latency:** ~500ms from order creation to print
- **Health check:** Every 30 seconds
- **Printer timeout:** 5 seconds
- **Concurrent prints:** Supported

## Getting Help

### Check Logs

**PM2:**
```bash
pm2 logs print-bridge --lines 100
```

**npm start:**
- Check terminal output

### Run Integration Test

```bash
./test-local-print-integration.sh
```

### Common Log Messages

**✅ Good:**
```
✅ Printer 192.168.1.100:9100 is online
Print successful - Job ID: 507f1f77bcf86cd799439011
Backend updated - Job ID: 507f1f77bcf86cd799439011 -> COMPLETED
```

**⚠️ Warning:**
```
BACKEND_URL not configured, skipping status update
Failed to update backend for job 507f1f77bcf86cd799439011
```
→ Check BACKEND_URL in .env

**❌ Error:**
```
Printer offline or unreachable at 192.168.1.100:9100
Connection timeout to 192.168.1.100:9100
```
→ Check printer IP, power, and network

## Advanced Configuration

### Change Port

Edit `.env`:
```env
PORT=3002
```

Restart service.

### Multiple Printers

Add more printers in POS Print Management UI. Each printer can have different IP.

### Custom Timeout

Edit `.env`:
```env
PRINTER_TIMEOUT=10000  # 10 seconds
```

### Auto-Start on Windows

1. Install PM2: `npm install -g pm2`
2. Start service: `pm2 start src/index.js --name "print-bridge"`
3. Save: `pm2 save`
4. Setup startup: `pm2 startup`
5. Follow the instructions shown

### Auto-Start on macOS/Linux

Same as Windows, but PM2 will create appropriate startup scripts.

## Maintenance

### Update Local Print Bridge

```bash
cd local-print-bridge
git pull
npm install
pm2 restart print-bridge
```

### View Statistics

```bash
curl http://localhost:3001/status
```

### Clear PM2 Logs

```bash
pm2 flush
```

## Support

For detailed technical documentation, see: [LOCAL_PRINT_BRIDGE_INTEGRATION.md](LOCAL_PRINT_BRIDGE_INTEGRATION.md)

For local bridge code, see: [local-print-bridge/README.md](local-print-bridge/README.md)

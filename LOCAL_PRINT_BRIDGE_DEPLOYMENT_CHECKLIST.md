# Local Print Bridge - Deployment Checklist

## Pre-Deployment Checklist

### 1. Hardware Requirements ✓

- [ ] Cafe computer with Node.js 16+ installed
- [ ] Thermal printers (ESC/POS compatible)
- [ ] Printers connected to local network
- [ ] Network switch/router
- [ ] Stable internet connection

### 2. Network Information Gathered ✓

- [ ] Printer IP addresses noted (e.g., 192.168.1.100)
- [ ] Printer ports confirmed (usually 9100)
- [ ] EC2 backend URL noted (e.g., https://your-domain.com)
- [ ] Network allows outbound HTTPS (port 443)

### 3. Software Installed ✓

- [ ] Node.js installed and verified (`node --version`)
- [ ] npm installed and verified (`npm --version`)
- [ ] Git installed (optional, for updates)

## Deployment Steps

### Step 1: Install Local Print Bridge

```bash
# Navigate to project
cd local-print-bridge

# Install dependencies
npm install
```

**Verification:**
- [ ] No errors during `npm install`
- [ ] `node_modules/` folder created

### Step 2: Configure Environment

```bash
# Copy example config
cp .env.example .env

# Edit configuration
nano .env  # or use any text editor
```

**Required settings:**
```env
PORT=3001
BACKEND_URL=https://your-actual-domain.com
PRINTER_TIMEOUT=5000
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

**Verification:**
- [ ] `.env` file created
- [ ] `BACKEND_URL` set to actual EC2 domain
- [ ] Printer IPs match actual printer addresses

### Step 3: Test Printer Connections

```bash
# Test bill printer
node src/test-printer.js 192.168.1.100 9100

# Test label printer
node src/test-printer.js 192.168.1.101 9100
```

**Expected output:**
```
✅ Printer 192.168.1.100:9100 is online
```

**Verification:**
- [ ] Bill printer responds
- [ ] Label printer responds
- [ ] No timeout errors

### Step 4: Start Service (Testing)

```bash
npm start
```

**Expected output:**
```
==================================================
🖨️  Local Print Bridge Started
==================================================
Server running on: http://localhost:3001
Backend URL: https://your-domain.com
Default Bill Printer: 192.168.1.100
Default Label Printer: 192.168.1.101
==================================================
Ready to accept print requests!
==================================================
```

**Verification:**
- [ ] Service starts without errors
- [ ] Port 3001 is listening
- [ ] Backend URL is correct
- [ ] Printer IPs are correct

### Step 5: Test Health Endpoint

Open browser or use curl:
```bash
curl http://localhost:3001/health
```

**Expected response:**
```json
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0",
  "timestamp": "2024-02-16T10:00:00.000Z"
}
```

**Verification:**
- [ ] Health endpoint responds
- [ ] Status is "ok"
- [ ] Service name is correct

### Step 6: Run Integration Test

```bash
# Stop the service (Ctrl+C)
# Run test script
./test-local-print-integration.sh
```

**Verification:**
- [ ] All tests pass
- [ ] Local bridge health check passes
- [ ] Backend health check passes
- [ ] Printer connection test passes (or acceptable failure)

### Step 7: Setup Production Mode (PM2)

```bash
# Install PM2 globally
npm install -g pm2

# Start service with PM2
pm2 start src/index.js --name "print-bridge"

# Save PM2 configuration
pm2 save

# Setup auto-start on boot
pm2 startup
# Follow the instructions shown
```

**Verification:**
- [ ] PM2 installed successfully
- [ ] Service running in PM2 (`pm2 status`)
- [ ] Service shows "online" status
- [ ] PM2 configuration saved
- [ ] Auto-start configured

### Step 8: Configure Printers in POS

1. Open POS in browser
2. Navigate to **Print Management** (🖨️ menu)
3. Click **Máy In** tab
4. Click **+ Thêm Máy In**

**For Bill Printer:**
- [ ] Name: "Bill Printer"
- [ ] Type: "BILL"
- [ ] Connection Type: "NETWORK"
- [ ] IP: 192.168.1.100
- [ ] Port: 9100
- [ ] Paper Width: 80mm (or 58mm)
- [ ] Set as Default: ✓
- [ ] Click "Test Connection" - should succeed
- [ ] Click "Lưu"

**For Label Printer:**
- [ ] Name: "Label Printer"
- [ ] Type: "LABEL"
- [ ] Connection Type: "NETWORK"
- [ ] IP: 192.168.1.101
- [ ] Port: 9100
- [ ] Label Size: 40x30mm (or appropriate size)
- [ ] Set as Default: ✓
- [ ] Click "Test Connection" - should succeed
- [ ] Click "Lưu"

### Step 9: Verify Local Bridge Status in POS

**In Print Management header:**
- [ ] Shows "Local Bridge Online" with green dot
- [ ] If offline, check service is running: `pm2 status`

### Step 10: Test End-to-End Printing

**Create Test Order:**
1. Go to POS main screen
2. Add items to order
3. Complete payment
4. Mark as PAID

**Verify:**
- [ ] Bill prints automatically
- [ ] Labels print automatically (if items have labels)
- [ ] Print jobs appear in "Print Jobs" tab
- [ ] Job status shows "COMPLETED"
- [ ] Print quality is acceptable

**Test Manual Reprint:**
1. Open order detail
2. Click "In lại Bill"
3. Verify bill reprints

**Verify:**
- [ ] Reprint button works
- [ ] New print job created
- [ ] Bill reprints successfully

## Post-Deployment Verification

### Functional Tests ✓

- [ ] Auto-print on order creation works
- [ ] Manual reprint works
- [ ] Print jobs show correct status
- [ ] Failed prints show error messages
- [ ] Multiple orders print in sequence
- [ ] Print quality is acceptable

### Performance Tests ✓

- [ ] Print latency < 1 second
- [ ] No delays in order creation
- [ ] UI remains responsive during printing
- [ ] Multiple concurrent prints work

### Error Handling Tests ✓

**Test 1: Printer Offline**
1. Turn off printer
2. Create order
3. Verify: Job shows "FAILED" with error message
4. Turn on printer
5. Click "Retry" on failed job
6. Verify: Job prints successfully

**Test 2: Local Bridge Offline**
1. Stop local bridge: `pm2 stop print-bridge`
2. Refresh POS
3. Verify: Header shows "Local Bridge Offline"
4. Create order
5. Verify: Job stays in "PENDING" (no error)
6. Start bridge: `pm2 start print-bridge`
7. Verify: Header shows "Local Bridge Online"

**Test 3: Network Issues**
1. Disconnect internet briefly
2. Create order
3. Verify: System handles gracefully
4. Reconnect internet
5. Verify: System recovers

### Monitoring Setup ✓

**PM2 Monitoring:**
```bash
# View status
pm2 status

# View logs
pm2 logs print-bridge

# View detailed info
pm2 info print-bridge
```

**Verification:**
- [ ] Can view service status
- [ ] Can view logs
- [ ] Logs show print activity
- [ ] No error messages in logs

**Statistics Endpoint:**
```bash
curl http://localhost:3001/status
```

**Verification:**
- [ ] Shows total prints
- [ ] Shows success rate
- [ ] Shows last print time

## Maintenance Procedures

### Daily Checks ✓

- [ ] Service is running: `pm2 status`
- [ ] No errors in logs: `pm2 logs print-bridge --lines 50`
- [ ] Printers are online and have paper
- [ ] Print quality is good

### Weekly Checks ✓

- [ ] Review statistics: `curl http://localhost:3001/status`
- [ ] Check success rate (should be >95%)
- [ ] Clear old logs: `pm2 flush`
- [ ] Test printer connections

### Monthly Checks ✓

- [ ] Update local bridge if new version available
- [ ] Review and optimize printer settings
- [ ] Check for PM2 updates: `npm update -g pm2`
- [ ] Backup configuration files

## Troubleshooting Guide

### Issue: "Local Bridge Offline"

**Check:**
```bash
pm2 status
```

**Fix:**
```bash
pm2 restart print-bridge
```

### Issue: Printer Not Printing

**Check:**
```bash
node src/test-printer.js 192.168.1.100 9100
```

**Fix:**
- Verify printer is on
- Check printer has paper
- Verify IP address is correct
- Check network connection

### Issue: Print Jobs Stuck in PENDING

**Check logs:**
```bash
pm2 logs print-bridge --lines 50
```

**Common causes:**
- Wrong printer IP
- Printer offline
- Network firewall blocking port 9100
- Local bridge not running

### Issue: Backend Status Not Updating

**Check BACKEND_URL:**
```bash
cat .env | grep BACKEND_URL
```

**Test backend:**
```bash
curl https://your-domain.com/health
```

**Fix:**
- Update BACKEND_URL in .env
- Restart service: `pm2 restart print-bridge`

## Rollback Procedure

If deployment fails:

1. **Stop service:**
   ```bash
   pm2 stop print-bridge
   pm2 delete print-bridge
   ```

2. **Restore previous version:**
   ```bash
   git checkout previous-version
   npm install
   ```

3. **Restart service:**
   ```bash
   pm2 start src/index.js --name "print-bridge"
   ```

## Success Criteria

Deployment is successful when:

- ✅ Local bridge service is running
- ✅ "Local Bridge Online" shows in POS
- ✅ Auto-print works on order creation
- ✅ Manual reprint works
- ✅ Print jobs show correct status
- ✅ Error handling works correctly
- ✅ Print quality is acceptable
- ✅ Success rate >95%
- ✅ Print latency <1 second
- ✅ No errors in logs

## Sign-Off

**Deployed by:** ___________________

**Date:** ___________________

**Verified by:** ___________________

**Date:** ___________________

**Notes:**
_____________________________________________
_____________________________________________
_____________________________________________

## Support Contacts

**Technical Issues:**
- Check documentation: LOCAL_PRINT_BRIDGE_INTEGRATION.md
- Check quick start: LOCAL_PRINT_BRIDGE_QUICK_START.md
- Run test: ./test-local-print-integration.sh

**Emergency Contacts:**
- System Administrator: ___________________
- Network Administrator: ___________________
- Printer Vendor: ___________________

## Related Documentation

- [Integration Guide](LOCAL_PRINT_BRIDGE_INTEGRATION.md)
- [Quick Start Guide](LOCAL_PRINT_BRIDGE_QUICK_START.md)
- [Implementation Summary](LOCAL_PRINT_BRIDGE_IMPLEMENTATION_SUMMARY.md)
- [Local Bridge README](local-print-bridge/README.md)

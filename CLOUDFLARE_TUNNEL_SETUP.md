# Cloudflare Tunnel Setup - Print Bridge

## Tổng quan

Sử dụng Cloudflare Tunnel để expose Print Bridge (chạy trong LAN) ra internet, cho phép Backend trên EC2 gửi lệnh in trực tiếp.

## Kiến trúc

```
Backend (EC2)
    │
    │ HTTPS POST /print
    ▼
Cloudflare Tunnel
(print.tacafe.store)
    │
    │ HTTP (local)
    ▼
Print Bridge (LAN)
    │
    │ ESC/POS
    ▼
Thermal Printers
```

## Ưu điểm

✅ Backend chủ động gửi lệnh in (không cần Print Bridge kết nối ra ngoài)
✅ Không cần mở port trên router
✅ Bảo mật với Cloudflare
✅ HTTPS miễn phí
✅ Không cần static IP
✅ Dễ setup và maintain

## Bước 1: Cài đặt cloudflared

### macOS
```bash
brew install cloudflared
```

### Linux
```bash
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

### Windows
Download từ: https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/

## Bước 2: Login Cloudflare

```bash
cloudflared tunnel login
```

Browser sẽ mở, chọn domain `tacafe.store` và authorize.

## Bước 3: Tạo Tunnel

```bash
# Tạo tunnel
cloudflared tunnel create print-bridge

# Lưu lại Tunnel ID (sẽ hiển thị sau khi tạo)
# Example: a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

## Bước 4: Cấu hình Tunnel

Tạo file `~/.cloudflared/config.yml`:

```yaml
tunnel: <TUNNEL_ID>
credentials-file: /Users/<username>/.cloudflared/<TUNNEL_ID>.json

ingress:
  # Route print.tacafe.store to local Print Bridge
  - hostname: print.tacafe.store
    service: http://localhost:3001
  
  # Catch-all rule
  - service: http_status:404
```

**Thay thế:**
- `<TUNNEL_ID>` bằng tunnel ID từ bước 3
- `<username>` bằng username của bạn

## Bước 5: Thêm DNS Record

### Option A: Sử dụng CLI (Tự động)

```bash
cloudflared tunnel route dns print-bridge print.tacafe.store
```

### Option B: Manual qua Dashboard

1. Vào https://dash.cloudflare.com
2. Chọn domain `tacafe.store`
3. DNS → Add record:
   - Type: `CNAME`
   - Name: `print`
   - Target: `<TUNNEL_ID>.cfargotunnel.com`
   - Proxy status: Proxied (orange cloud)
   - TTL: Auto

## Bước 6: Start Print Bridge

```bash
cd local-print-bridge

# Start Print Bridge với Docker
docker-compose up -d

# Verify
curl http://localhost:3001/health
```

## Bước 7: Start Tunnel

### Option A: Foreground (Testing)

```bash
cloudflared tunnel run print-bridge
```

Logs sẽ hiển thị:
```
INF Connection registered connIndex=0 location=SJC
INF Connection registered connIndex=1 location=LAX
```

### Option B: Background Service (Production)

**macOS/Linux:**
```bash
# Install as service
sudo cloudflared service install

# Start service
sudo systemctl start cloudflared
sudo systemctl enable cloudflared

# Check status
sudo systemctl status cloudflared

# View logs
sudo journalctl -u cloudflared -f
```

**Windows:**
```powershell
# Install as service
cloudflared service install

# Start service
sc start cloudflared

# Check status
sc query cloudflared
```

## Bước 8: Test Tunnel

```bash
# Test từ bên ngoài
curl https://print.tacafe.store/health

# Expected response:
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0"
}
```

## Bước 9: Cấu hình Backend

### 9.1 Thêm biến môi trường

File `.env.ec2`:
```bash
# Print Bridge URL (via Cloudflare Tunnel)
PRINT_BRIDGE_URL=https://print.tacafe.store
```

### 9.2 Update Backend Code

Backend sẽ gọi Print Bridge qua HTTP thay vì WebSocket:

```go
// Khi tạo print job
printBridgeClient.SendPrintJob(printbridge.PrintRequest{
    JobID:       job.ID.Hex(),
    Content:     job.Content,
    PrinterIP:   printer.IPAddress,
    PrinterPort: printer.Port,
    Type:        string(job.Type),
})
```

## Bước 10: Test End-to-End

1. Tạo order trên Frontend
2. Backend tạo print job
3. Backend gửi HTTP POST tới `https://print.tacafe.store/print`
4. Cloudflare Tunnel forward tới Print Bridge local
5. Print Bridge in ra máy in
6. Print Bridge update status về Backend

## Monitoring

### Check Tunnel Status

```bash
# List tunnels
cloudflared tunnel list

# Check tunnel info
cloudflared tunnel info print-bridge

# View logs
sudo journalctl -u cloudflared -f  # Linux
tail -f /var/log/cloudflared.log   # macOS
```

### Check Print Bridge Logs

```bash
cd local-print-bridge
docker-compose logs -f
```

### Test Print Endpoint

```bash
curl -X POST https://print.tacafe.store/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-123",
    "content": "Test print content",
    "printerIP": "192.168.1.100",
    "printerPort": 9100,
    "type": "bill"
  }'
```

## Troubleshooting

### Tunnel không kết nối

```bash
# Check cloudflared service
sudo systemctl status cloudflared

# Restart service
sudo systemctl restart cloudflared

# Check logs
sudo journalctl -u cloudflared -f
```

### DNS không resolve

```bash
# Check DNS
nslookup print.tacafe.store

# Should return Cloudflare IPs (104.x.x.x or 172.x.x.x)
```

### Print Bridge không accessible

```bash
# Check Print Bridge running
docker-compose ps

# Check Print Bridge logs
docker-compose logs

# Test local
curl http://localhost:3001/health
```

### Backend không gửi được request

```bash
# Test từ EC2
curl https://print.tacafe.store/health

# Check backend logs
docker logs cafe-pos-backend | grep PrintBridge
```

## Security

### 1. Authentication (Optional)

Thêm authentication token:

**Print Bridge (.env):**
```bash
BRIDGE_AUTH_TOKEN=your-secret-token-here
```

**Backend (.env.ec2):**
```bash
PRINT_BRIDGE_AUTH_TOKEN=your-secret-token-here
```

### 2. IP Whitelist (Cloudflare)

1. Vào Cloudflare Dashboard
2. Security → WAF
3. Create rule:
   - If: Hostname equals `print.tacafe.store`
   - And: IP Source Address not in `<EC2-IP>`
   - Then: Block

### 3. Rate Limiting

Cloudflare tự động có rate limiting, nhưng có thể tùy chỉnh:
1. Security → Rate Limiting
2. Create rule cho `print.tacafe.store`

## Cost

✅ **Cloudflare Tunnel: MIỄN PHÍ**
- Unlimited bandwidth
- Unlimited requests
- HTTPS included
- DDoS protection included

## Backup Plan

Nếu Cloudflare Tunnel down, có thể fallback về WebSocket:

1. Print Bridge vẫn có WebSocket client
2. Backend có thể emit qua WebSocket
3. Hoặc dùng ngrok/localtunnel tạm thời

## Quick Start Script

Sử dụng script tự động:

```bash
cd local-print-bridge
./cloudflare-tunnel-setup.sh
```

Script sẽ:
1. Install cloudflared (nếu chưa có)
2. Login Cloudflare
3. Create tunnel
4. Generate config
5. Hiển thị hướng dẫn tiếp theo

## Related Documentation

- `local-print-bridge/README.md` - Print Bridge overview
- `PRINT_ARCHITECTURE.md` - System architecture
- `PRINT_BRIDGE_COMPLETE_GUIDE.md` - Complete setup guide

## Summary

**Setup:**
1. Install cloudflared
2. Create tunnel
3. Add DNS record
4. Start Print Bridge
5. Start tunnel
6. Configure backend

**Flow:**
Backend → HTTPS → Cloudflare → Tunnel → Print Bridge → Printer

**Benefits:**
- No port forwarding needed
- Free HTTPS
- Secure
- Easy to setup

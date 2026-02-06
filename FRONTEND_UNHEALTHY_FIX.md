# 🔧 Fix Frontend Container Unhealthy Issue

## 🐛 Vấn Đề

Frontend container hiển thị status `unhealthy` trong `docker ps`:
```
Up 23 hours (unhealthy)   0.0.0.0:80->80/tcp
```

## 🔍 Root Cause Analysis

### Triệu Chứng:
1. ✅ nginx đang chạy bình thường
2. ✅ Port 80 đang listen
3. ✅ Files tồn tại trong `/usr/share/nginx/html/`
4. ✅ HTTP request từ HOST thành công (200 OK)
5. ❌ Healthcheck từ INSIDE container fail: "Connection refused"

### Debug Output:
```bash
# Healthcheck command từ inside container
$ docker exec cafe-pos-frontend wget --quiet --tries=1 --spider http://localhost
wget: can't connect to remote host: Connection refused

# Nhưng từ host thì OK
$ curl -I http://localhost
HTTP/1.1 200 OK
```

### Root Cause: IPv6 vs IPv4 Conflict

**Vấn đề:**
- nginx config chỉ listen trên IPv4: `listen 80;` → `0.0.0.0:80`
- wget trong healthcheck resolve `localhost` thành IPv6: `::1`
- Không có connection giữa IPv6 client và IPv4 server

**Evidence:**
```bash
# wget trying to connect to IPv6
Connecting to localhost ([::1]:80)
wget: can't connect to remote host: Connection refused

# nginx only listening on IPv4
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN
```

## ✅ Giải Pháp

### Option 1: Add IPv6 Support to nginx (RECOMMENDED)

**File:** `frontend/nginx.conf`

```nginx
server {
    listen 80;           # IPv4
    listen [::]:80;      # IPv6 - ADD THIS LINE
    server_name localhost;
    # ... rest of config
}
```

**Ưu điểm:**
- ✅ Support cả IPv4 và IPv6
- ✅ Modern best practice
- ✅ Healthcheck sẽ work với cả 2 protocols

### Option 2: Change Healthcheck to Use IPv4

**File:** `docker-compose.hub.yml`

```yaml
frontend:
  healthcheck:
    # Change from http://localhost to http://127.0.0.1
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://127.0.0.1"]
    interval: 30s
    timeout: 10s
    retries: 3
```

**Ưu điểm:**
- ✅ Quick fix, không cần rebuild image
- ✅ Force IPv4

**Nhược điểm:**
- ❌ Không support IPv6
- ❌ Cần update docker-compose

### Option 3: Disable IPv6 in Container

**File:** `docker-compose.hub.yml`

```yaml
frontend:
  sysctls:
    - net.ipv6.conf.all.disable_ipv6=1
```

**Nhược điểm:**
- ❌ Không recommended
- ❌ Disable toàn bộ IPv6

## 🚀 Recommended Fix Steps

### Step 1: Update nginx.conf ✅

```bash
# File: frontend/nginx.conf
server {
    listen 80;
    listen [::]:80;  # Add this line
    server_name localhost;
    # ... rest
}
```

### Step 2: Rebuild và Push Image

```bash
# Build new image
cd frontend
docker build -t linhtranphu/cafe-pos-frontend:latest .

# Push to Docker Hub
docker push linhtranphu/cafe-pos-frontend:latest
```

### Step 3: Update Container trên EC2

**SSH vào EC2:**
```bash
ssh -i /path/to/key.pem ubuntu@13.212.27.222
```

**Pull image mới và restart:**
```bash
# Pull latest image
docker pull linhtranphu/cafe-pos-frontend:latest

# Restart container
docker-compose -f docker-compose.hub.yml up -d --force-recreate frontend

# Or restart all
docker-compose -f docker-compose.hub.yml restart frontend
```

### Step 4: Verify Health

```bash
# Wait 30 seconds for healthcheck
sleep 30

# Check status
docker ps | grep frontend

# Should show: Up X minutes (healthy)
```

## 🧪 Testing

### Test 1: Healthcheck từ Inside Container
```bash
docker exec cafe-pos-frontend wget --quiet --tries=1 --spider http://localhost
# Should succeed (no output)
echo $?
# Should return 0
```

### Test 2: Check nginx Listening
```bash
docker exec cafe-pos-frontend netstat -tuln | grep :80
# Should show both IPv4 and IPv6:
# tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN
# tcp6       0      0 :::80                   :::*                    LISTEN
```

### Test 3: HTTP Request từ Host
```bash
curl -I http://localhost
# Should return: HTTP/1.1 200 OK
```

### Test 4: Docker Health Status
```bash
docker inspect --format='{{.State.Health.Status}}' cafe-pos-frontend
# Should return: healthy
```

## 📊 Before vs After

### Before (Unhealthy):
```
nginx listen:     0.0.0.0:80 (IPv4 only)
wget resolves:    localhost → ::1 (IPv6)
Connection:       ❌ REFUSED (IPv6 → IPv4)
Health Status:    unhealthy
```

### After (Healthy):
```
nginx listen:     0.0.0.0:80 (IPv4) + :::80 (IPv6)
wget resolves:    localhost → ::1 (IPv6)
Connection:       ✅ SUCCESS (IPv6 → IPv6)
Health Status:    healthy
```

## 🔍 Debug Commands

### Check Container Health
```bash
# Full health status
docker inspect --format='{{json .State.Health}}' cafe-pos-frontend | jq '.'

# Just status
docker inspect --format='{{.State.Health.Status}}' cafe-pos-frontend
```

### Check nginx Process
```bash
docker exec cafe-pos-frontend ps aux | grep nginx
```

### Check Listening Ports
```bash
docker exec cafe-pos-frontend netstat -tuln | grep :80
# or
docker exec cafe-pos-frontend ss -tuln | grep :80
```

### Test Healthcheck Manually
```bash
# IPv4
docker exec cafe-pos-frontend wget --quiet --tries=1 --spider http://127.0.0.1

# IPv6
docker exec cafe-pos-frontend wget --quiet --tries=1 --spider http://[::1]

# localhost (auto-resolve)
docker exec cafe-pos-frontend wget --quiet --tries=1 --spider http://localhost
```

### View Container Logs
```bash
# Last 50 lines
docker logs --tail 50 cafe-pos-frontend

# Follow logs
docker logs -f cafe-pos-frontend

# Since 1 hour ago
docker logs --since 1h cafe-pos-frontend
```

## 📝 Notes

### Why This Happens:
1. Modern Linux containers often prefer IPv6 when available
2. `localhost` can resolve to either `127.0.0.1` (IPv4) or `::1` (IPv6)
3. nginx default `listen 80` only binds to IPv4
4. wget in Alpine Linux (base image) prefers IPv6

### Best Practice:
Always configure nginx to listen on both IPv4 and IPv6:
```nginx
listen 80;
listen [::]:80;
```

### Alternative Healthcheck:
If you want to force IPv4 in healthcheck:
```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://127.0.0.1"]
```

## 🎯 Quick Fix Summary

**1 line change in `frontend/nginx.conf`:**
```diff
server {
    listen 80;
+   listen [::]:80;
    server_name localhost;
```

**Then rebuild and redeploy:**
```bash
docker build -t linhtranphu/cafe-pos-frontend:latest frontend/
docker push linhtranphu/cafe-pos-frontend:latest
# On EC2: docker-compose up -d --force-recreate frontend
```

---

**Date:** February 6, 2026  
**Issue:** Frontend container unhealthy  
**Root Cause:** IPv6/IPv4 mismatch  
**Solution:** Add IPv6 support to nginx  
**Status:** ✅ Fixed

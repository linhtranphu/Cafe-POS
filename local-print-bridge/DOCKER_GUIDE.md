# Docker Deployment Guide - Local Print Bridge

## Quick Start

### 1. Cấu hình

```bash
cd local-print-bridge

# Copy .env (nếu chưa có)
cp .env.example .env

# Edit .env và cấu hình printer IPs
nano .env
```

### 2. Build và Run

```bash
# Build image
docker-compose build

# Start service
docker-compose up -d

# Check logs
docker-compose logs -f

# Check status
docker-compose ps
```

### 3. Test

```bash
# Health check
curl http://localhost:3001/health

# Test printer connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100"}'
```

## Commands

### Start/Stop

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Restart
docker-compose restart

# View logs
docker-compose logs -f

# View logs (last 100 lines)
docker-compose logs --tail=100
```

### Build

```bash
# Build image
docker-compose build

# Build without cache
docker-compose build --no-cache

# Pull latest base images
docker-compose build --pull
```

### Maintenance

```bash
# Remove container
docker-compose down

# Remove container and volumes
docker-compose down -v

# Remove container, volumes, and images
docker-compose down -v --rmi all

# Prune unused Docker resources
docker system prune -a
```

## Configuration

### Environment Variables

Edit `.env` file:

```bash
# Server
PORT=3001

# Printers (configure your local network IPs)
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_BILL_PRINTER_PORT=9100
DEFAULT_LABEL_PRINTER_IP=192.168.1.116
DEFAULT_LABEL_PRINTER_PORT=9100

# Logging
LOG_LEVEL=info

# Timeouts
PRINTER_TIMEOUT=5
```

### Network Mode

Docker compose sử dụng `network_mode: host` để:
- Access printers trên local network
- Không cần port mapping phức tạp
- Direct access như native application

**Lưu ý:** `network_mode: host` chỉ hoạt động trên Linux. Trên macOS/Windows, cần dùng bridge network và expose ports.

### For macOS/Windows

Nếu dùng macOS/Windows, thay đổi `docker-compose.yml`:

```yaml
services:
  print-bridge:
    # Remove network_mode: host
    # Add ports mapping
    ports:
      - "3001:3001"
    # Add networks
    networks:
      - print-network

networks:
  print-network:
    driver: bridge
```

## Troubleshooting

### Container không start

```bash
# Check logs
docker-compose logs

# Check container status
docker-compose ps

# Inspect container
docker inspect local-print-bridge
```

### Chromium errors

```bash
# Check if Chromium is installed
docker-compose exec print-bridge which chromium-browser

# Check Chromium version
docker-compose exec print-bridge chromium-browser --version

# Test Chromium
docker-compose exec print-bridge chromium-browser --headless --dump-dom about:blank
```

### Cannot connect to printer

```bash
# Check network connectivity from container
docker-compose exec print-bridge ping 192.168.1.100

# Check if port is open
docker-compose exec print-bridge nc -zv 192.168.1.100 9100

# Check container network mode
docker inspect local-print-bridge | grep NetworkMode
```

### High memory usage

```bash
# Check memory usage
docker stats local-print-bridge

# Increase shared memory if needed
# Edit docker-compose.yml:
shm_size: '512mb'  # Increase from 256mb
```

## Health Check

Container có health check tự động:

```bash
# Check health status
docker-compose ps

# Manual health check
curl http://localhost:3001/health
```

Health check runs every 30 seconds:
- ✅ healthy: Service is running
- ⚠️ unhealthy: Service has issues
- 🔄 starting: Service is starting

## Logs

### View logs

```bash
# All logs
docker-compose logs

# Follow logs (real-time)
docker-compose logs -f

# Last 100 lines
docker-compose logs --tail=100

# Logs since 1 hour ago
docker-compose logs --since 1h

# Logs with timestamps
docker-compose logs -t
```

### Log levels

Set in `.env`:
```bash
LOG_LEVEL=info  # debug, info, warn, error
```

## Performance

### Resource Limits

Add to `docker-compose.yml`:

```yaml
services:
  print-bridge:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### Optimization

1. **Use multi-stage build** (already implemented)
2. **Minimize image size** (Alpine base)
3. **Cache Go modules** (in Dockerfile)
4. **Use .dockerignore** (already created)

## Security

### Best Practices

1. **Run as non-root user** (optional):

```dockerfile
# Add to Dockerfile
RUN addgroup -g 1000 printbridge && \
    adduser -D -u 1000 -G printbridge printbridge
USER printbridge
```

2. **Read-only filesystem** (optional):

```yaml
services:
  print-bridge:
    read_only: true
    tmpfs:
      - /tmp
```

3. **Drop capabilities**:

```yaml
services:
  print-bridge:
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

## Backup & Restore

### Backup configuration

```bash
# Backup .env
cp .env .env.backup

# Backup docker-compose.yml
cp docker-compose.yml docker-compose.yml.backup
```

### Export/Import image

```bash
# Export image
docker save local-print-bridge:latest | gzip > print-bridge.tar.gz

# Import image
docker load < print-bridge.tar.gz
```

## Updates

### Update to latest version

```bash
# Pull latest code
git pull

# Rebuild image
docker-compose build --no-cache

# Restart service
docker-compose up -d

# Check logs
docker-compose logs -f
```

## Integration with Backend

Backend cần config `PRINT_BRIDGE_URL`:

```bash
# If print bridge runs on same machine as backend
PRINT_BRIDGE_URL=http://localhost:3001

# If print bridge runs on different machine
PRINT_BRIDGE_URL=http://192.168.1.X:3001
```

## Monitoring

### Check service status

```bash
# Container status
docker-compose ps

# Health check
curl http://localhost:3001/health

# Service status
curl http://localhost:3001/status

# Docker stats
docker stats local-print-bridge
```

### Metrics

```bash
# Memory usage
docker stats --no-stream local-print-bridge | awk '{print $4}'

# CPU usage
docker stats --no-stream local-print-bridge | awk '{print $3}'

# Uptime
docker inspect local-print-bridge | grep StartedAt
```

## Production Deployment

### Systemd service

Create `/etc/systemd/system/print-bridge.service`:

```ini
[Unit]
Description=Local Print Bridge (Docker)
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/local-print-bridge
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable:

```bash
sudo systemctl enable print-bridge
sudo systemctl start print-bridge
```

### Auto-restart on boot

Already configured in `docker-compose.yml`:

```yaml
restart: unless-stopped
```

## FAQ

### Q: Tại sao dùng network_mode: host?
**A:** Để access printers trên local network dễ dàng hơn.

### Q: Có thể chạy nhiều instances không?
**A:** Có, nhưng phải dùng port khác nhau.

### Q: Memory usage bao nhiêu?
**A:** ~50MB idle, ~150MB khi rendering.

### Q: Có cần GPU không?
**A:** Không, Chromium chạy headless mode.

### Q: Có thể chạy trên Raspberry Pi không?
**A:** Có, build cho ARM architecture.

---

**Version:** 1.0.0  
**Last Updated:** 2024

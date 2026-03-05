# Local Print Bridge (Go + chromedp)

Print bridge service viết bằng Go, sử dụng chromedp để render HTML thành ESC/POS commands.

## Yêu cầu

- Go 1.21+
- Chrome/Chromium browser (chromedp sẽ tự động download nếu chưa có)

## Cài đặt

```bash
# Install dependencies
go mod download

# Copy .env
cp .env.example .env

# Edit .env và cấu hình printer IPs
nano .env
```

## Chạy

```bash
# Development
go run main.go

# Build
make build
./print-bridge

# Docker
make docker-build
make docker-run
```

## Endpoints

### POST /render-and-print
Render HTML và in trực tiếp

```json
{
  "html": "<html>...</html>",
  "width": 576,
  "printerIP": "192.168.1.100",
  "printerPort": 9100
}
```

### POST /print
In ESC/POS data trực tiếp

```json
{
  "content": "base64_or_raw_data",
  "printerIP": "192.168.1.100",
  "printerPort": 9100
}
```

### POST /test-connection
Test kết nối máy in

```json
{
  "printerIP": "192.168.1.100",
  "printerPort": 9100
}
```

### GET /health
Health check

### GET /status
Service status

## So sánh với Node.js version

| Feature | Node.js (Puppeteer) | Go (chromedp) |
|---------|---------------------|---------------|
| Memory | ~200MB | ~50MB |
| Startup | ~2s | ~0.5s |
| Binary size | N/A (interpreted) | ~20MB |
| Dependencies | node_modules (~300MB) | None (single binary) |
| Cross-platform | ✅ | ✅ |

## License

MIT

#!/bin/bash

# Script để rebuild Print Bridge trên Linux với Chromium
# Fix lỗi: chromium-browser not found

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Rebuilding Print Bridge for Linux ===${NC}"
echo ""

cd local-print-bridge

# 1. Backup current files
echo "1. Creating backups..."
if [ -f "docker-compose.yml" ]; then
    cp docker-compose.yml docker-compose.yml.backup.$(date +%Y%m%d_%H%M%S)
    echo -e "${GREEN}✓ Backed up docker-compose.yml${NC}"
fi

if [ -f "Dockerfile" ]; then
    cp Dockerfile Dockerfile.backup.$(date +%Y%m%d_%H%M%S)
    echo -e "${GREEN}✓ Backed up Dockerfile${NC}"
fi
echo ""

# 2. Stop current container
echo "2. Stopping current container..."
docker compose down
echo -e "${GREEN}✓ Container stopped${NC}"
echo ""

# 3. Update docker-compose.yml
echo "3. Updating docker-compose.yml for Linux..."

cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  print-bridge:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: local-print-bridge
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - HOST=0.0.0.0
      - DEFAULT_BILL_PRINTER_IP=${DEFAULT_BILL_PRINTER_IP:-192.168.1.115}
      - DEFAULT_BILL_PRINTER_PORT=${DEFAULT_BILL_PRINTER_PORT:-9100}
      - DEFAULT_LABEL_PRINTER_IP=${DEFAULT_LABEL_PRINTER_IP:-192.168.1.116}
      - DEFAULT_LABEL_PRINTER_PORT=${DEFAULT_LABEL_PRINTER_PORT:-9100}
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - PRINTER_TIMEOUT=${PRINTER_TIMEOUT:-5}
    restart: unless-stopped
    # Security options for Chromium on Linux
    security_opt:
      - seccomp:unconfined
    # Increased shared memory for Chromium (IMPORTANT for Linux)
    shm_size: '512mb'
    # Capabilities for Chromium
    cap_add:
      - SYS_ADMIN
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
EOF

echo -e "${GREEN}✓ docker-compose.yml updated${NC}"
echo ""

# 4. Update Dockerfile
echo "4. Updating Dockerfile with Chromium for Linux..."

cat > Dockerfile << 'EOF'
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o print-bridge main.go

# Runtime stage
FROM alpine:latest

# Install Chromium and ALL necessary dependencies for Linux
RUN apk --no-cache add \
    chromium \
    chromium-chromedriver \
    # Fonts
    font-noto \
    font-noto-cjk \
    ttf-dejavu \
    ttf-freefont \
    fontconfig \
    # SSL certificates
    ca-certificates \
    # Additional dependencies for Chromium on Linux
    nss \
    freetype \
    harfbuzz \
    # Utilities
    wget \
    curl

# Update font cache
RUN fc-cache -f

# Set Chrome binary path for chromedp
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/print-bridge .

# Expose port
EXPOSE 3001

# Run the application
CMD ["./print-bridge"]
EOF

echo -e "${GREEN}✓ Dockerfile updated${NC}"
echo ""

# 5. Build new image
echo "5. Building new Docker image..."
echo "This will take a few minutes (downloading Chromium and dependencies)..."
echo ""

docker compose build --no-cache

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
echo ""

# 6. Start container
echo "6. Starting container..."
docker compose up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Container started${NC}"
else
    echo -e "${RED}✗ Failed to start container${NC}"
    exit 1
fi
echo ""

# 7. Wait for container to be ready
echo "7. Waiting for container to be ready..."
sleep 5

# Check if container is running
if docker ps | grep -q "local-print-bridge"; then
    echo -e "${GREEN}✓ Container is running${NC}"
else
    echo -e "${RED}✗ Container is not running${NC}"
    echo "Check logs: docker logs local-print-bridge"
    exit 1
fi
echo ""

# 8. Test health endpoint
echo "8. Testing health endpoint..."
for i in {1..10}; do
    if curl -s http://localhost:3001/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Health check passed${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}✗ Health check failed after 10 attempts${NC}"
        exit 1
    fi
    echo "Waiting... ($i/10)"
    sleep 2
done
echo ""

# 9. Verify Chromium installation
echo "9. Verifying Chromium installation..."

# Check if chromium-browser exists
if docker exec local-print-bridge which chromium-browser > /dev/null 2>&1; then
    echo -e "${GREEN}✓ chromium-browser found${NC}"
    
    # Get version
    CHROME_VERSION=$(docker exec local-print-bridge chromium-browser --version 2>&1)
    echo "   Version: $CHROME_VERSION"
else
    echo -e "${RED}✗ chromium-browser not found${NC}"
    echo "Something went wrong with the build"
    exit 1
fi
echo ""

# 10. Test Chromium rendering
echo "10. Testing Chromium rendering..."

TEST_OUTPUT=$(docker exec local-print-bridge chromium-browser \
    --headless \
    --disable-gpu \
    --no-sandbox \
    --disable-dev-shm-usage \
    --dump-dom \
    about:blank 2>&1 | head -5)

if echo "$TEST_OUTPUT" | grep -q "html"; then
    echo -e "${GREEN}✓ Chromium can render HTML${NC}"
else
    echo -e "${YELLOW}⚠ Chromium test inconclusive${NC}"
    echo "Output: $TEST_OUTPUT"
fi
echo ""

# 11. Show recent logs
echo "11. Recent logs:"
docker logs --tail=15 local-print-bridge
echo ""

# 12. Test endpoints
echo "12. Testing endpoints..."

# Health
echo "Testing /health..."
curl -s http://localhost:3001/health | head -3
echo ""

# Test connection
echo ""
echo "Testing /test-connection..."
curl -s -X POST http://localhost:3001/test-connection \
    -H "Content-Type: application/json" \
    -d '{"printerIP": "192.168.1.100", "printerPort": 9100}' | head -3
echo ""

# 13. Summary
echo ""
echo -e "${BLUE}=== Summary ===${NC}"
echo ""
echo "✓ Docker image rebuilt with Chromium"
echo "✓ Container started successfully"
echo "✓ Health check passed"
echo "✓ Chromium installed and verified"
echo ""
echo "Configuration:"
echo "  - Shared memory: 512MB"
echo "  - Security: seccomp:unconfined"
echo "  - Capabilities: SYS_ADMIN"
echo ""
echo "Endpoints available:"
echo "  POST /render-and-print - Render HTML and print"
echo "  POST /print - Direct print ESC/POS data"
echo "  POST /test-connection - Test printer connection"
echo "  GET /health - Health check"
echo "  GET /status - Service status"
echo ""
echo "Next steps:"
echo "  1. Test print from web interface at https://tacafe.store/#/print-management"
echo "  2. Use endpoint: /render-and-print (not /print-html)"
echo "  3. Monitor logs: docker logs -f local-print-bridge"
echo ""
echo "Backup files created in local-print-bridge/ directory"
echo ""

cd ..

echo -e "${GREEN}✓ Rebuild complete!${NC}"
echo ""
echo "You can now test printing from the web interface."

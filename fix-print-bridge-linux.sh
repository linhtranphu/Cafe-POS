#!/bin/bash

# Script để fix Print Bridge trên Linux
# Thêm các flags cần thiết cho Chromium

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Fixing Print Bridge for Linux ===${NC}"
echo ""

# Check if on Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo -e "${YELLOW}This script is for Linux systems${NC}"
    echo "Current OS: $OSTYPE"
    exit 0
fi

# 1. Update docker-compose.yml with proper settings
echo "1. Updating docker-compose.yml..."

cd local-print-bridge

# Backup original
cp docker-compose.yml docker-compose.yml.backup

# Create updated docker-compose.yml
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
      # Chromium flags for Linux
      - CHROME_BIN=/usr/bin/chromium-browser
      - CHROME_PATH=/usr/bin/chromium-browser
    restart: unless-stopped
    # Security options for Chromium on Linux
    security_opt:
      - seccomp:unconfined
    # Increased shared memory for Chromium
    shm_size: '512mb'
    # Capabilities needed for Chromium
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

# 2. Update Dockerfile with better Linux support
echo "2. Updating Dockerfile..."

# Backup original
cp Dockerfile Dockerfile.backup

# Create updated Dockerfile
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

# Install Chromium and all necessary dependencies for Linux
RUN apk --no-cache add \
    chromium \
    chromium-chromedriver \
    font-noto \
    font-noto-cjk \
    ttf-dejavu \
    fontconfig \
    ca-certificates \
    # Additional dependencies for Chromium on Linux
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    # For debugging
    wget \
    curl

# Update font cache
RUN fc-cache -f

# Set Chrome binary path and flags for chromedp
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser
# Chromium flags for headless mode on Linux
ENV CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage --disable-gpu --headless"

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

# 3. Rebuild container
echo "3. Rebuilding container..."
echo "This may take a few minutes..."

docker compose down
docker compose build --no-cache
docker compose up -d

echo -e "${GREEN}✓ Container rebuilt and started${NC}"
echo ""

# 4. Wait for container to be ready
echo "4. Waiting for container to be ready..."
sleep 5

# 5. Check health
echo "5. Checking health..."
for i in {1..10}; do
    if curl -s http://localhost:3001/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Print Bridge is healthy${NC}"
        break
    fi
    echo "Waiting... ($i/10)"
    sleep 2
done
echo ""

# 6. Test Chromium
echo "6. Testing Chromium rendering..."

TEST_RESULT=$(docker exec local-print-bridge chromium-browser \
    --headless \
    --disable-gpu \
    --no-sandbox \
    --disable-dev-shm-usage \
    --dump-dom \
    about:blank 2>&1 | head -5)

if echo "$TEST_RESULT" | grep -q "html"; then
    echo -e "${GREEN}✓ Chromium is working${NC}"
else
    echo -e "${RED}✗ Chromium test failed${NC}"
    echo "Output: $TEST_RESULT"
fi
echo ""

# 7. View logs
echo "7. Recent logs:"
docker logs --tail=20 local-print-bridge
echo ""

# 8. Summary
echo -e "${BLUE}=== Summary ===${NC}"
echo ""
echo "Changes made:"
echo "  ✓ Updated docker-compose.yml with:"
echo "    - Increased shared memory (512MB)"
echo "    - Added SYS_ADMIN capability"
echo "    - Added security options"
echo "  ✓ Updated Dockerfile with:"
echo "    - Additional Linux dependencies"
echo "    - Chromium flags for headless mode"
echo "  ✓ Rebuilt container"
echo ""
echo "Next steps:"
echo "  1. Test print from web interface"
echo "  2. Check logs: docker logs -f local-print-bridge"
echo "  3. If still issues, run: ./diagnose-print-bridge-linux.sh"
echo ""
echo "Backup files created:"
echo "  - docker-compose.yml.backup"
echo "  - Dockerfile.backup"
echo ""

cd ..

echo -e "${GREEN}✓ Fix complete!${NC}"

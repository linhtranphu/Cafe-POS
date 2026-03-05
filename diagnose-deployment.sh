#!/bin/bash

# Script để CHẨN ĐOÁN nguyên nhân server bị chết
# KHÔNG FIX - CHỈ THU THẬP THÔNG TIN

set -e

echo "=========================================="
echo "🔍 CHẨN ĐOÁN DEPLOYMENT ISSUES"
echo "=========================================="
echo ""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

LOG_FILE="deployment-diagnosis-$(date +%Y%m%d-%H%M%S).log"

log() {
    echo -e "$1" | tee -a "$LOG_FILE"
}

log "${BLUE}📝 Log file: $LOG_FILE${NC}"
log ""

# 1. Kiểm tra system resources
log "${YELLOW}=== 1. SYSTEM RESOURCES ===${NC}"
log "CPU Info:"
sysctl -n machdep.cpu.brand_string | tee -a "$LOG_FILE"
sysctl -n hw.ncpu | xargs echo "CPU Cores:" | tee -a "$LOG_FILE"
log ""

log "Memory Info:"
vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages\s+([^:]+)[^\d]+(\d+)/ and printf("%-16s % 16.2f Mi\n", "$1:", $2 * $size / 1048576);' | tee -a "$LOG_FILE"
log ""

log "Disk Space:"
df -h / | tee -a "$LOG_FILE"
log ""

# 2. Kiểm tra Docker resources
log "${YELLOW}=== 2. DOCKER CONFIGURATION ===${NC}"
log "Docker Version:"
docker version --format '{{.Server.Version}}' | tee -a "$LOG_FILE"
log ""

log "Docker Info:"
docker info 2>&1 | grep -E "CPUs|Total Memory|Docker Root Dir" | tee -a "$LOG_FILE"
log ""

# 3. Kiểm tra containers hiện tại
log "${YELLOW}=== 3. CURRENT CONTAINERS ===${NC}"
docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | tee -a "$LOG_FILE"
log ""

# 4. Kiểm tra resource usage của containers đang chạy
log "${YELLOW}=== 4. CONTAINER RESOURCE USAGE ===${NC}"
if [ $(docker ps -q | wc -l) -gt 0 ]; then
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}" | tee -a "$LOG_FILE"
else
    log "No running containers"
fi
log ""

# 5. Kiểm tra logs của containers (nếu có)
log "${YELLOW}=== 5. CONTAINER LOGS (Last 50 lines) ===${NC}"
for container in cafe-pos-backend cafe-pos-frontend cafe-pos-mongodb; do
    if docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
        log "${BLUE}--- $container ---${NC}"
        docker logs --tail 50 "$container" 2>&1 | tee -a "$LOG_FILE"
        log ""
    fi
done

# 6. Kiểm tra network
log "${YELLOW}=== 6. DOCKER NETWORKS ===${NC}"
docker network ls | tee -a "$LOG_FILE"
log ""

# 7. Kiểm tra volumes
log "${YELLOW}=== 7. DOCKER VOLUMES ===${NC}"
docker volume ls | tee -a "$LOG_FILE"
log ""

# 8. Kiểm tra images
log "${YELLOW}=== 8. DOCKER IMAGES ===${NC}"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | grep -E "cafe-pos|REPOSITORY" | tee -a "$LOG_FILE"
log ""

# 9. Kiểm tra .env configuration
log "${YELLOW}=== 9. ENVIRONMENT CONFIGURATION ===${NC}"
if [ -f .env ]; then
    log "✅ .env exists"
    log "MongoDB URI format:"
    grep "MONGODB_URI" .env | sed 's/password[^@]*/password:***/' | tee -a "$LOG_FILE"
    log ""
    log "Port configuration:"
    grep "PORT" .env | tee -a "$LOG_FILE"
else
    log "${RED}❌ .env file not found${NC}"
fi
log ""

# 10. Kiểm tra docker-compose files
log "${YELLOW}=== 10. DOCKER-COMPOSE FILES ===${NC}"
for file in docker-compose.yml docker-compose.local.yml docker-compose.prod.yml; do
    if [ -f "$file" ]; then
        log "✅ $file exists"
    else
        log "❌ $file not found"
    fi
done
log ""

# 11. Test MongoDB connection
log "${YELLOW}=== 11. MONGODB CONNECTION TEST ===${NC}"
if docker ps --format '{{.Names}}' | grep -q "cafe-pos-mongodb"; then
    log "Testing MongoDB connection..."
    docker exec cafe-pos-mongodb mongosh --eval "db.adminCommand('ping')" 2>&1 | tee -a "$LOG_FILE"
else
    log "MongoDB container not running"
fi
log ""

# 12. Kiểm tra port conflicts
log "${YELLOW}=== 12. PORT USAGE CHECK ===${NC}"
log "Port 80 (Frontend):"
lsof -i :80 2>&1 | tee -a "$LOG_FILE" || log "Port 80 is free"
log ""
log "Port 3000 (Backend):"
lsof -i :3000 2>&1 | tee -a "$LOG_FILE" || log "Port 3000 is free"
log ""
log "Port 27017 (MongoDB):"
lsof -i :27017 2>&1 | tee -a "$LOG_FILE" || log "Port 27017 is free"
log ""

# Summary
log "${GREEN}=========================================="
log "✅ CHẨN ĐOÁN HOÀN TẤT"
log "==========================================${NC}"
log ""
log "📋 Kết quả đã được lưu vào: ${BLUE}$LOG_FILE${NC}"
log ""
log "${YELLOW}BƯỚC TIẾP THEO:${NC}"
log "1. Xem file log: cat $LOG_FILE"
log "2. Chạy deploy: ./deploy-local-monitor.sh"
log "3. Quan sát resource usage trong 5-10 phút"
log "4. Tìm pattern: Container nào chết? Khi nào? Tại sao?"
log ""

#!/bin/bash

echo "=== Checking Timezone Configuration ==="
echo ""

echo "1. System Timezone:"
timedatectl 2>/dev/null || date

echo ""
echo "2. Backend Container Timezone:"
docker exec cafe-pos-backend date 2>/dev/null || echo "Backend container not running"

echo ""
echo "3. Backend Container TZ Environment:"
docker exec cafe-pos-backend printenv TZ 2>/dev/null || echo "Backend container not running"

echo ""
echo "4. MongoDB Container Timezone:"
docker exec cafe-pos-mongodb date 2>/dev/null || echo "MongoDB container not running"

echo ""
echo "5. Current time in Vietnam (should match):"
TZ=Asia/Ho_Chi_Minh date

echo ""
echo "=== Recommendations ==="
echo "- All times should show UTC+7 (Vietnam time)"
echo "- If system timezone is wrong, run: sudo timedatectl set-timezone Asia/Ho_Chi_Minh"
echo "- If container timezone is wrong, rebuild with updated Dockerfile"

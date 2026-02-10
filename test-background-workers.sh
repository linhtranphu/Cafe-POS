#!/bin/bash

# Script kiểm tra Background Workers hoạt động
# Sử dụng: ./test-background-workers.sh

set -e

echo "🧪 Testing Background Workers..."
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Backend URL
BACKEND_URL="http://localhost:3000"

# Step 1: Login
echo "📝 Step 1: Đăng nhập..."
LOGIN_RESPONSE=$(curl -s -X POST "$BACKEND_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ Đăng nhập thất bại!${NC}"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo -e "${GREEN}✅ Đăng nhập thành công${NC}"
echo ""

# Step 2: Check health status
echo "📊 Step 2: Kiểm tra health status..."
HEALTH_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BACKEND_URL/api/manager/monitoring/health")

echo "$HEALTH_RESPONSE" | jq '.'

STATUS=$(echo $HEALTH_RESPONSE | jq -r '.status')
if [ "$STATUS" == "healthy" ]; then
  echo -e "${GREEN}✅ System health: HEALTHY${NC}"
elif [ "$STATUS" == "degraded" ]; then
  echo -e "${YELLOW}⚠️  System health: DEGRADED${NC}"
else
  echo -e "${RED}❌ System health: CRITICAL${NC}"
fi
echo ""

# Step 3: Check aggregated metrics
echo "📈 Step 3: Kiểm tra aggregated metrics..."
METRICS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BACKEND_URL/api/manager/monitoring/metrics/aggregated")

echo "$METRICS_RESPONSE" | jq '.'

TOTAL_RECALC=$(echo $METRICS_RESPONSE | jq -r '.total_recalc_jobs')
SUCCESS_RECALC=$(echo $METRICS_RESPONSE | jq -r '.successful_recalc_jobs')
FAILED_RECALC=$(echo $METRICS_RESPONSE | jq -r '.failed_recalc_jobs')

echo ""
echo "📊 Recalculation Jobs Statistics:"
echo "   Total: $TOTAL_RECALC"
echo "   Success: $SUCCESS_RECALC"
echo "   Failed: $FAILED_RECALC"

if [ "$TOTAL_RECALC" -gt 0 ]; then
  SUCCESS_RATE=$(echo "scale=2; $SUCCESS_RECALC * 100 / $TOTAL_RECALC" | bc)
  echo "   Success Rate: ${SUCCESS_RATE}%"
  
  if (( $(echo "$SUCCESS_RATE >= 95" | bc -l) )); then
    echo -e "${GREEN}✅ Workers hoạt động tốt (success rate >= 95%)${NC}"
  elif (( $(echo "$SUCCESS_RATE >= 80" | bc -l) )); then
    echo -e "${YELLOW}⚠️  Workers hoạt động khá (success rate >= 80%)${NC}"
  else
    echo -e "${RED}❌ Workers có vấn đề (success rate < 80%)${NC}"
  fi
else
  echo -e "${YELLOW}⚠️  Chưa có recalculation jobs nào được xử lý${NC}"
  echo "   Tip: Thử cập nhật cost_per_unit của một ingredient để trigger workers"
fi
echo ""

# Step 4: Check recent recalculation metrics
echo "🔍 Step 4: Kiểm tra recent recalculation jobs..."
RECENT_JOBS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BACKEND_URL/api/manager/monitoring/metrics?type=recalculation_job&limit=5")

JOB_COUNT=$(echo $RECENT_JOBS | jq '.count')
echo "Recent jobs: $JOB_COUNT"

if [ "$JOB_COUNT" -gt 0 ]; then
  echo ""
  echo "Latest jobs:"
  echo "$RECENT_JOBS" | jq -r '.metrics[] | "  - Status: \(.status), Duration: \(.duration / 1000000)ms, Time: \(.timestamp)"'
  echo -e "${GREEN}✅ Workers đang xử lý jobs${NC}"
else
  echo -e "${YELLOW}⚠️  Không có recent jobs${NC}"
  echo "   Tip: Workers sẽ tự động xử lý khi có ingredient cost updates"
fi
echo ""

# Step 5: Check alerts
echo "🚨 Step 5: Kiểm tra alerts..."
ALERTS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BACKEND_URL/api/manager/monitoring/alerts?limit=5")

ALERT_COUNT=$(echo $ALERTS_RESPONSE | jq '.count')
echo "Total alerts: $ALERT_COUNT"

if [ "$ALERT_COUNT" -gt 0 ]; then
  echo ""
  echo "Recent alerts:"
  echo "$ALERTS_RESPONSE" | jq -r '.alerts[] | "  - [\(.level | ascii_upcase)] \(.message)"'
  
  CRITICAL_COUNT=$(echo $ALERTS_RESPONSE | jq '[.alerts[] | select(.level == "critical")] | length')
  if [ "$CRITICAL_COUNT" -gt 0 ]; then
    echo -e "${RED}❌ Có $CRITICAL_COUNT critical alerts!${NC}"
  else
    echo -e "${YELLOW}⚠️  Có alerts nhưng không critical${NC}"
  fi
else
  echo -e "${GREEN}✅ Không có alerts${NC}"
fi
echo ""

# Step 6: Get ingredients to test
echo "🧪 Step 6: Test trigger recalculation (optional)..."
echo "Lấy danh sách ingredients..."

INGREDIENTS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BACKEND_URL/api/manager/ingredients")

INGREDIENT_COUNT=$(echo $INGREDIENTS_RESPONSE | jq '. | length')
echo "Tìm thấy $INGREDIENT_COUNT ingredients"

if [ "$INGREDIENT_COUNT" -gt 0 ]; then
  FIRST_INGREDIENT_ID=$(echo $INGREDIENTS_RESPONSE | jq -r '.[0]._id')
  FIRST_INGREDIENT_NAME=$(echo $INGREDIENTS_RESPONSE | jq -r '.[0].name')
  CURRENT_COST=$(echo $INGREDIENTS_RESPONSE | jq -r '.[0].cost_per_unit')
  
  echo ""
  echo "Để test workers, bạn có thể chạy:"
  echo ""
  echo -e "${YELLOW}curl -X PUT \"$BACKEND_URL/api/manager/ingredients/$FIRST_INGREDIENT_ID\" \\${NC}"
  echo -e "${YELLOW}  -H \"Authorization: Bearer $TOKEN\" \\${NC}"
  echo -e "${YELLOW}  -H \"Content-Type: application/json\" \\${NC}"
  echo -e "${YELLOW}  -d '{\"name\":\"$FIRST_INGREDIENT_NAME\",\"cost_per_unit\":$((CURRENT_COST + 1000))}'${NC}"
  echo ""
  echo "Sau đó chạy lại script này để xem workers xử lý job!"
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "System Status: $STATUS"
echo "Total Recalculation Jobs: $TOTAL_RECALC"
echo "Success Rate: ${SUCCESS_RATE:-N/A}%"
echo "Active Alerts: $ALERT_COUNT"
echo ""

if [ "$STATUS" == "healthy" ] && [ "$ALERT_COUNT" -eq 0 ]; then
  echo -e "${GREEN}🎉 Background workers đang hoạt động tốt!${NC}"
  echo ""
  echo "✅ Workers tự động start khi server khởi động"
  echo "✅ Workers sẽ xử lý jobs khi cập nhật ingredient costs"
  echo "✅ Monitoring và alerts đang hoạt động"
elif [ "$STATUS" == "degraded" ]; then
  echo -e "${YELLOW}⚠️  System đang degraded, cần kiểm tra${NC}"
  echo ""
  echo "Xem chi tiết alerts:"
  echo "curl -H \"Authorization: Bearer $TOKEN\" \\"
  echo "  \"$BACKEND_URL/api/manager/monitoring/alerts?level=warning\""
else
  echo -e "${RED}❌ System có vấn đề, cần xử lý ngay!${NC}"
  echo ""
  echo "Xem chi tiết alerts:"
  echo "curl -H \"Authorization: Bearer $TOKEN\" \\"
  echo "  \"$BACKEND_URL/api/manager/monitoring/alerts?level=critical\""
fi
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

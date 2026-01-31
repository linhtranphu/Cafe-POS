#!/bin/bash

echo "=== Restarting Backend ==="
echo ""

echo "1. Stopping old backend..."
pkill -f cafe-pos-server
sleep 1

OLD_PROCESS=$(ps aux | grep cafe-pos-server | grep -v grep)
if [ -z "$OLD_PROCESS" ]; then
  echo "   ✅ Old backend stopped"
else
  echo "   ⚠️  Backend still running, force killing..."
  pkill -9 -f cafe-pos-server
  sleep 1
fi
echo ""

echo "2. Building new backend..."
cd backend
go build -o cafe-pos-server
BUILD_STATUS=$?

if [ $BUILD_STATUS -eq 0 ]; then
  echo "   ✅ Build successful"
else
  echo "   ❌ Build failed"
  exit 1
fi
echo ""

echo "3. Starting new backend..."
nohup ./cafe-pos-server > server.log 2>&1 &
NEW_PID=$!
cd ..

sleep 2
echo ""

echo "4. Verifying backend is running..."
PROCESS_INFO=$(ps aux | grep cafe-pos-server | grep -v grep)

if [ -z "$PROCESS_INFO" ]; then
  echo "   ❌ Backend failed to start"
  echo "   Check logs: tail -f backend/server.log"
  exit 1
else
  echo "   ✅ Backend is running"
  echo "   PID: $NEW_PID"
  echo "   Process: $PROCESS_INFO"
fi
echo ""

echo "5. Testing API..."
sleep 1
HEALTH_CHECK=$(curl -s http://localhost:8080/api/auth/login -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"invalid","password":"invalid"}' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$HEALTH_CHECK" | tail -n1)

if [ "$HTTP_CODE" == "401" ] || [ "$HTTP_CODE" == "400" ]; then
  echo "   ✅ API is responding (HTTP $HTTP_CODE)"
else
  echo "   ⚠️  API response: HTTP $HTTP_CODE"
fi
echo ""

echo "=== Backend Restarted Successfully! ==="
echo ""
echo "Logs: tail -f backend/server.log"
echo "Stop: pkill -f cafe-pos-server"

#!/bin/bash

echo "=== Setting EC2 Timezone to Asia/Ho_Chi_Minh ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root or with sudo"
    exit 1
fi

echo "1. Setting system timezone..."
timedatectl set-timezone Asia/Ho_Chi_Minh

echo ""
echo "2. Verifying system timezone..."
timedatectl

echo ""
echo "3. Current system time:"
date

echo ""
echo "4. Restarting backend container to apply timezone..."
docker restart cafe-pos-backend

echo ""
echo "5. Waiting for backend to start..."
sleep 5

echo ""
echo "6. Checking backend container timezone:"
docker exec cafe-pos-backend date

echo ""
echo "7. Checking backend TZ environment:"
docker exec cafe-pos-backend printenv TZ

echo ""
echo "=== Done! ==="
echo "Backend should now use Vietnam timezone (UTC+7)"
echo "Test by creating a new order and checking the printed time"

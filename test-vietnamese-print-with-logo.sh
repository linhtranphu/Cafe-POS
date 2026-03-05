#!/bin/bash

echo "🖨️  Test In Tiếng Việt với Logo"
echo "================================"
echo ""
echo "Đang compile chương trình test..."

cd backend/cmd/test-vietnamese-print
go build -o test-vietnamese-print

if [ $? -ne 0 ]; then
    echo "❌ Lỗi compile"
    exit 1
fi

echo "✓ Compile thành công"
echo ""
echo "Đang gửi lệnh in đến máy in 192.168.1.115:9100..."
echo ""

./test-vietnamese-print

echo ""
echo "Hoàn tất! Kiểm tra máy in để xem kết quả."

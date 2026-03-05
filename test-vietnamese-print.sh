#!/bin/bash

echo "=== TEST IN TIẾNG VIỆT ==="
echo ""
echo "Chương trình này sẽ test in tiếng Việt với:"
echo "  - Code Page 255 (Zywell/Xprinter internal Vietnamese)"
echo "  - IP máy in: 192.168.1.115:9100"
echo ""

cd backend/cmd/test-vietnamese-print

echo "1. Building test program..."
go build -o test-vietnamese-print

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✓ Build successful!"
echo ""
echo "2. Running test print..."
./test-vietnamese-print

echo ""
echo "=== HƯỚNG DẪN ==="
echo "Kiểm tra máy in xem các ký tự tiếng Việt có in đúng không:"
echo "  - Chữ có dấu: à, á, ả, ã, ạ, ă, ằ, ắ, ẳ, ẵ, ặ, â, ầ, ấ, ẩ, ẫ, ậ"
echo "  - Chữ ê: è, é, ẻ, ẽ, ẹ, ê, ề, ế, ể, ễ, ệ"
echo "  - Chữ ô, ơ: ò, ó, ỏ, õ, ọ, ô, ồ, ố, ổ, ỗ, ộ, ơ, ờ, ớ, ở, ỡ, ợ"
echo "  - Chữ ư: ù, ú, ủ, ũ, ụ, ư, ừ, ứ, ử, ữ, ự"
echo "  - Chữ y: ỳ, ý, ỷ, ỹ, ỵ"
echo ""
echo "Nếu in đúng → Code page 255 hoạt động tốt!"
echo "Nếu in sai → Thử code page 31 (Windows-1258) hoặc 30"

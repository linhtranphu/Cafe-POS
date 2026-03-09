#!/bin/bash

# Script để so sánh 2 luồng in:
# 1. Test Print từ Print Management (THÀNH CÔNG)
# 2. Auto Print khi Collect Payment (THẤT BẠI)

echo "=========================================="
echo "SO SÁNH 2 LUỒNG IN"
echo "=========================================="
echo ""

# Màu sắc
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}1. KIỂM TRA PRINT BRIDGE URL TRONG SETTINGS${NC}"
echo "----------------------------------------"
SETTINGS=$(curl -s http://localhost:3000/api/settings)
BRIDGE_URL=$(echo $SETTINGS | jq -r '.print_bridge_url // "NOT SET"')
AUTO_PRINT=$(echo $SETTINGS | jq -r '.auto_print_enabled // false')

echo "Print Bridge URL: $BRIDGE_URL"
echo "Auto Print Enabled: $AUTO_PRINT"
echo ""

if [ "$BRIDGE_URL" == "NOT SET" ] || [ "$BRIDGE_URL" == "null" ]; then
    echo -e "${RED}❌ PRINT BRIDGE URL CHƯA ĐƯỢC CẤU HÌNH!${NC}"
    echo "Vào Settings và cấu hình Print Bridge URL"
    exit 1
fi

if [ "$AUTO_PRINT" != "true" ]; then
    echo -e "${RED}❌ AUTO PRINT BỊ TẮT!${NC}"
    echo "Vào Settings và bật 'Tự động in khi thu tiền'"
    exit 1
fi

echo -e "${GREEN}✅ Settings OK${NC}"
echo ""

echo -e "${YELLOW}2. KIỂM TRA PRINT BRIDGE ĐANG CHẠY${NC}"
echo "----------------------------------------"
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" $BRIDGE_URL/health)

if [ "$HEALTH" == "200" ]; then
    echo -e "${GREEN}✅ Print Bridge đang chạy: $BRIDGE_URL${NC}"
else
    echo -e "${RED}❌ Print Bridge KHÔNG phản hồi: $BRIDGE_URL${NC}"
    echo "HTTP Status: $HEALTH"
    exit 1
fi
echo ""

echo -e "${YELLOW}3. KIỂM TRA CẤU HÌNH MÁY IN${NC}"
echo "----------------------------------------"
PRINTERS=$(curl -s http://localhost:3000/api/manager/printers)
BILL_PRINTER=$(echo $PRINTERS | jq -r '.printers[] | select(.type == "bill" and .is_default == true)')

if [ -z "$BILL_PRINTER" ]; then
    echo -e "${RED}❌ KHÔNG TÌM THẤY MÁY IN BILL MẶC ĐỊNH!${NC}"
    exit 1
fi

PRINTER_IP=$(echo $BILL_PRINTER | jq -r '.ip_address')
PRINTER_PORT=$(echo $BILL_PRINTER | jq -r '.port')
PRINTER_NAME=$(echo $BILL_PRINTER | jq -r '.name')

echo "Máy in mặc định: $PRINTER_NAME"
echo "IP: $PRINTER_IP:$PRINTER_PORT"
echo -e "${GREEN}✅ Máy in OK${NC}"
echo ""

echo -e "${YELLOW}4. PHÂN TÍCH LUỒNG TEST PRINT (THÀNH CÔNG)${NC}"
echo "----------------------------------------"
echo "Luồng:"
echo "  Frontend → POST /api/manager/html-templates/test-print"
echo "  Backend → html_template_handler_bridge.go::TestPrintHTMLTemplate()"
echo "  Backend → printBridgeClient.RenderAndPrint()"
echo "  Backend → POST $BRIDGE_URL/render-and-print"
echo "  Print Bridge → In ra máy in"
echo ""
echo "URL được sử dụng:"
echo "  - Lấy từ: shopSettings.PrintBridgeURL"
echo "  - Giá trị: $BRIDGE_URL"
echo ""

echo -e "${YELLOW}5. PHÂN TÍCH LUỒNG AUTO PRINT (THẤT BẠI)${NC}"
echo "----------------------------------------"
echo "Luồng:"
echo "  Frontend → POST /api/waiter/orders/:id/payment"
echo "  Backend → order_handler.go::CollectPayment()"
echo "  Backend → order_service.go::CollectPayment()"
echo "  Backend → print_service.go::CreatePrintJobsForOrder()"
echo "  Backend → print_worker.go::processPrintJob()"
echo "  Backend → bridge_printer.go::Print()"
echo "  Backend → printBridgeClient.RenderAndPrint()"
echo "  Backend → POST ???/render-and-print"
echo ""

echo -e "${YELLOW}6. KIỂM TRA PRINT BRIDGE CLIENT TRONG MAIN.GO${NC}"
echo "----------------------------------------"
echo "Cần kiểm tra xem print bridge client được khởi tạo như thế nào:"
echo ""
grep -A 10 "printbridge.NewClient" backend/main.go | head -20
echo ""

echo -e "${YELLOW}7. KIỂM TRA PRINT SERVICE KHỞI TẠO${NC}"
echo "----------------------------------------"
echo "Cần kiểm tra xem print service có nhận đúng bridge client không:"
echo ""
grep -A 5 "NewPrintService" backend/main.go | head -10
echo ""

echo -e "${YELLOW}8. KIỂM TRA BRIDGE PRINTER KHỞI TẠO${NC}"
echo "----------------------------------------"
echo "Cần kiểm tra xem bridge printer được tạo với client nào:"
echo ""
grep -B 5 -A 10 "NewBridgePrinter" backend/application/services/print_service.go | head -20
echo ""

echo "=========================================="
echo "KẾT LUẬN"
echo "=========================================="
echo ""
echo "Có thể có 2 vấn đề:"
echo ""
echo "1. Print Bridge Client được khởi tạo 2 lần:"
echo "   - Lần 1: Trong main.go khi start app (dùng cho test print)"
echo "   - Lần 2: Trong print service (dùng cho auto print)"
echo "   → Nếu lần 2 dùng URL khác hoặc không được khởi tạo đúng"
echo ""
echo "2. Print Service không nhận được bridge client:"
echo "   - Print service có thể đang dùng direct printer thay vì bridge printer"
echo "   - Hoặc bridge client không được inject vào print service"
echo ""
echo "Để debug chi tiết hơn, cần:"
echo "1. Thêm log trong bridge_printer.go::Print() để xem URL được gọi"
echo "2. Thêm log trong print_service.go để xem printer type"
echo "3. Check backend logs khi collect payment"
echo ""

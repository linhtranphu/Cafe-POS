# Auto Print via Print Bridge - Implementation Summary

## Vấn Đề

Khi deploy trên EC2, backend không thể kết nối trực tiếp đến máy in local (192.168.1.115) vì:
- Backend EC2 ở cloud
- Máy in ở mạng local
- Không cùng network

## Giải Pháp

Sử dụng **Print Bridge** làm proxy giữa EC2 và máy in local:

```
Backend EC2 → Print Bridge (via Cloudflare Tunnel) → Máy in Local
```

## Thay Đổi Code

### 1. Tạo BridgePrinter (backend/infrastructure/printing/bridge_printer.go)

Wrapper printer mới sử dụng print bridge thay vì kết nối trực tiếp:

```go
type BridgePrinter struct {
    config       *printing.PrinterConfig
    bridgeClient *printbridge.Client
    innerPrinter Printer
}
```

**Chức năng:**
- `Connect()`: Kiểm tra print bridge có available không
- `Print()`: Gửi ESC/POS data đến print bridge qua HTTP
- `Disconnect()`: No-op (không cần đóng connection)

### 2. Update PrinterManager (backend/application/services/printer_manager.go)

Thêm khả năng inject print bridge client:

```go
type printerManager struct {
    printBridgeClient *printbridge.Client
}

func (pm *printerManager) SetPrintBridgeClient(client *printbridge.Client)
```

**Logic:**
- Nếu có print bridge client và available → wrap printer với BridgePrinter
- Nếu không → sử dụng direct printer (như cũ)

### 3. Update main.go

Khởi tạo print bridge client từ shop settings:

```go
// Initialize print bridge client from shop settings
go func() {
    time.Sleep(2 * time.Second)
    
    settings, err := shopSettingsRepo.GetSettings(ctx)
    if err == nil && settings != nil && settings.PrintBridgeURL != "" {
        bridgeClient := printbridge.NewClient(settings.PrintBridgeURL, 30*time.Second)
        
        if bridgeClient.IsAvailable() {
            printerManager.SetPrintBridgeClient(bridgeClient)
            log.Printf("[PRINT BRIDGE] ✅ Print bridge configured and available")
        }
    }
}()
```

## Flow Hoàn Chỉnh

### Khi Collect Payment:

```
1. Frontend: POST /api/waiter/orders/:id/payment
   ↓
2. OrderService.CollectPayment()
   ↓
3. Check: order.Status == PAID? ✅
   Check: auto_print_enabled? ✅
   ↓
4. CreatePrintJobsForOrder() (async)
   ↓
5. Tạo print jobs trong database
   ↓
6. Print Worker (chạy mỗi 10 giây)
   ↓
7. Lấy pending jobs
   ↓
8. PrinterManager.GetPrinter()
   ├─ Có print bridge? → BridgePrinter
   └─ Không? → ESCPOSPrinter/LabelPrinter
   ↓
9. BridgePrinter.Print()
   ↓
10. printBridgeClient.PrintESCPOS()
    ↓
11. POST {print_bridge_url}/print
    Body: {
      "content": "base64_escpos_data",
      "printerIP": "192.168.1.115",
      "printerPort": 9100
    }
    ↓
12. Print Bridge nhận request
    ↓
13. Print Bridge gửi đến máy in local
    ↓
14. ✅ PRINTED!
```

## Cấu Hình Cần Thiết

### 1. Shop Settings

Vào **Print Management** → **Cài Đặt**:
- Print Bridge URL: `https://your-tunnel.trycloudflare.com`
- Auto Print Enabled: ✅

### 2. Printers

Cần có 2 printers:
- **BILL printer**: Type=BILL, IP=192.168.1.115, Port=9100, Paper Width=80mm
- **LABEL printer**: Type=LABEL, IP=192.168.1.115, Port=9100, Paper Width=58mm

### 3. Print Bridge

Chạy local print bridge với Cloudflare Tunnel:

```bash
cd local-print-bridge
./start-tunnel.sh
```

## Kiểm Tra

### 1. Kiểm tra Print Bridge

```bash
curl https://your-tunnel.trycloudflare.com/health
```

Kết quả:
```json
{
  "status": "ok",
  "timestamp": "2026-03-01T..."
}
```

### 2. Kiểm tra Backend Logs

```bash
# Trên EC2
tail -f backend.log | grep "PRINT BRIDGE"
```

Logs mong đợi:
```
[PRINT BRIDGE] Configuring print bridge: https://...
[PRINT BRIDGE] ✅ Print bridge configured and available
[BRIDGE PRINTER] Print bridge is available for printer: Máy in Bill
[BRIDGE PRINTER] Sending to print bridge - printer: Máy in Bill, IP: 192.168.1.115:9100
[BRIDGE PRINTER] Print successful via bridge - printer: Máy in Bill
```

### 3. Test Auto Print

1. Tạo order mới
2. Collect payment (tiền mặt)
3. Order status → PAID
4. Kiểm tra logs:
   ```bash
   tail -f backend.log | grep -E "PRINT|BRIDGE"
   ```
5. Bill và label sẽ tự động in

## Backward Compatibility

Code vẫn hoạt động với direct connection nếu:
- Không có print bridge URL trong settings
- Print bridge không available
- Đang chạy local (không cần bridge)

## Troubleshooting

### Print Bridge Not Available

**Triệu chứng:**
```
[PRINT BRIDGE] ⚠️  Print bridge configured but not available
```

**Giải pháp:**
1. Kiểm tra print bridge có chạy không
2. Kiểm tra Cloudflare tunnel có active không
3. Test URL: `curl https://your-tunnel.trycloudflare.com/health`

### Print Jobs Failed

**Triệu chứng:**
```
[PRINT ERROR] Print command failed
```

**Giải pháp:**
1. Kiểm tra print bridge logs
2. Kiểm tra máy in có bật không
3. Kiểm tra IP/Port đúng không
4. Test manual print: Vào Print Management → In Thử Bill Mẫu

### No Print Jobs Created

**Triệu chứng:**
- Collect payment nhưng không có print jobs

**Giải pháp:**
1. Kiểm tra `auto_print_enabled = true`
2. Kiểm tra có BILL và LABEL printer không
3. Kiểm tra backend logs: `grep "CreatePrintJobsForOrder" backend.log`

## Lợi Ích

✅ **Hoạt động trên EC2**: Backend cloud có thể in đến máy in local  
✅ **Backward compatible**: Vẫn hoạt động với direct connection  
✅ **Tự động fallback**: Nếu bridge không available, dùng direct  
✅ **Dễ debug**: Logs rõ ràng cho từng bước  
✅ **Secure**: Sử dụng Cloudflare Tunnel (HTTPS)  

## Next Steps

1. Deploy code lên EC2
2. Restart backend
3. Kiểm tra logs xem print bridge có được configure không
4. Test auto-print bằng cách tạo order và collect payment
5. Kiểm tra bill và label có in ra không

## Files Changed

- `backend/infrastructure/printing/bridge_printer.go` (NEW) - Wrapper printer sử dụng print bridge
- `backend/infrastructure/printing/helpers.go` (NEW) - Helper functions chung
- `backend/application/services/printer_manager.go` (MODIFIED) - Thêm print bridge support
- `backend/main.go` (MODIFIED) - Auto-initialize print bridge từ settings
- `frontend/src/components/printing/PrinterConfigForm.vue` (MODIFIED) - Thêm paper_width field cho LABEL printer

## Build Status

✅ **Compilation successful** - No errors

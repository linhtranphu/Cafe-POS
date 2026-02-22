# Hướng dẫn cấu hình In tự động

## Tổng quan

Khi waiter thu tiền (order chuyển sang trạng thái PAID), hệ thống sẽ tự động in:
- 1 hóa đơn (bill) - cho khách hàng
- N nhãn món (labels) - cho từng món trong order, gửi cho barista

## Yêu cầu để In tự động hoạt động

### 1. Templates (✅ Đã có)

- ✅ Default Bill Template: "Hóa đơn mặc định"
- ✅ Default Label Template: "Nhãn món mặc định"

### 2. Printer Configuration (❌ Cần cấu hình)

Cần cấu hình 2 máy in:

#### A. Bill Printer (Máy in hóa đơn)
- **Loại**: BILL
- **Đánh dấu**: Default (is_default = true)
- **Kết nối**: Network (TCP/IP)
- **Thông tin**: IP address + Port

#### B. Label Printer (Máy in nhãn món)
- **Loại**: LABEL  
- **Đánh dấu**: Default (is_default = true)
- **Kết nối**: Network (TCP/IP)
- **Thông tin**: IP address + Port

### 3. Auto-Print Setting (✅ Mặc định bật)

Setting `auto_print_enabled` mặc định = `true`

## Cách cấu hình Printer

### Bước 1: Truy cập Print Management

1. Login với tài khoản Manager
2. Vào menu: **Print Management** (hoặc **Quản lý In**)
3. Chọn tab **Printers**

### Bước 2: Thêm Bill Printer

1. Click **Add Printer** (hoặc **Thêm máy in**)
2. Điền thông tin:
   ```
   Name: Máy in hóa đơn
   Type: BILL
   Connection Type: NETWORK
   IP Address: [IP của máy in hóa đơn]
   Port: 9100 (hoặc port của máy in)
   Paper Width: 80 (hoặc 58 nếu dùng giấy 58mm)
   Is Default: ✓ (check vào)
   ```
3. Click **Test Connection** để kiểm tra kết nối
4. Click **Save**

### Bước 3: Thêm Label Printer

1. Click **Add Printer** lần nữa
2. Điền thông tin:
   ```
   Name: Máy in nhãn món
   Type: LABEL
   Connection Type: NETWORK
   IP Address: [IP của máy in nhãn]
   Port: 9100 (hoặc port của máy in)
   Paper Width: 58 (thường dùng giấy nhỏ cho nhãn)
   Is Default: ✓ (check vào)
   ```
3. Click **Test Connection** để kiểm tra kết nối
4. Click **Save**

### Bước 4: Kiểm tra Local Print Bridge

Print Bridge phải đang chạy để xử lý print jobs:

```bash
# Kiểm tra Print Bridge đang chạy
docker ps | grep print-bridge

# Nếu chưa chạy, start Print Bridge
cd local-print-bridge
docker-compose up -d
```

## Kiểm tra In tự động

### Test Flow

1. Tạo order mới với vai trò Waiter
2. Thêm món vào order
3. Thu tiền (Collect Payment)
4. Kiểm tra:
   - ✅ Hóa đơn in ra từ Bill Printer
   - ✅ Nhãn món in ra từ Label Printer (1 nhãn cho mỗi món)

### Xem Print Jobs

Vào **Print Management** > **Print Jobs** để xem:
- Danh sách print jobs
- Trạng thái: PENDING, PROCESSING, COMPLETED, FAILED
- Lỗi (nếu có)

## Troubleshooting

### Vấn đề 1: Không in tự động

**Nguyên nhân có thể:**
- Chưa cấu hình printer
- Printer không được đánh dấu default
- Auto-print setting bị tắt
- Print Bridge không chạy

**Giải pháp:**
1. Kiểm tra có printer default chưa:
   ```bash
   # Vào MongoDB
   use cafe_pos
   db.printer_configs.find({is_default: true})
   # Phải có 2 printers: 1 BILL, 1 LABEL
   ```

2. Kiểm tra auto-print setting:
   ```bash
   db.shop_settings.find({}, {auto_print_enabled: 1})
   # Phải là true
   ```

3. Kiểm tra Print Bridge:
   ```bash
   docker ps | grep print-bridge
   docker logs print-bridge
   ```

### Vấn đề 2: Print job bị FAILED

**Nguyên nhân:**
- Máy in offline
- IP/Port sai
- Máy in không hỗ trợ ESC/POS

**Giải pháp:**
1. Test connection từ Print Management
2. Kiểm tra IP/Port của máy in
3. Ping máy in từ server:
   ```bash
   ping [IP_máy_in]
   telnet [IP_máy_in] 9100
   ```

### Vấn đề 3: In ra chữ lỗi (tiếng Việt)

**Nguyên nhân:**
- Font chưa được cài đặt trên server

**Giải pháp:**
- Hệ thống đã chuyển sang in bằng image, không cần font đặc biệt
- Đảm bảo có 1 trong các font: Arial Unicode MS, Roboto, DejaVu Sans

### Vấn đề 4: Font size quá nhỏ/lớn

**Giải pháp:**
- Chỉnh font size trong `backend/infrastructure/printing/escpos_printer.go`
- Hiện tại: 26pt (có thể tăng lên 28pt hoặc 30pt)

## Tắt In tự động

Nếu muốn tắt in tự động:

1. Vào **Settings** (hoặc **Cài đặt**)
2. Tìm **Auto Print Enabled**
3. Bỏ check
4. Save

Khi tắt, waiter phải in thủ công bằng nút **Reprint Bill** / **Reprint Label**

## Cấu trúc Print Flow

```
Order Payment (PAID)
    ↓
Check auto_print_enabled = true?
    ↓ Yes
Get Default Printers (BILL + LABEL)
    ↓
Get Default Templates
    ↓
Create Print Jobs (1 bill + N labels)
    ↓
Print Worker picks up jobs
    ↓
Send to Print Bridge via WebSocket
    ↓
Print Bridge sends to physical printer
    ↓
Update job status: COMPLETED
```

## Lưu ý quan trọng

1. **Mỗi loại printer chỉ có 1 default**: Nếu đánh dấu printer mới là default, printer cũ sẽ tự động bỏ default

2. **Print jobs không block order**: Nếu in lỗi, order vẫn được tạo thành công. Có thể in lại sau.

3. **Label cho từng món**: Mỗi món trong order sẽ có 1 nhãn riêng, giúp barista dễ theo dõi

4. **Template có thể chỉnh sửa**: Vào Print Management > Templates để chỉnh sửa nội dung bill/label

5. **Font size**: Hiện tại 26pt, đủ lớn để đọc. Có thể tăng nếu cần.

## Liên hệ hỗ trợ

Nếu gặp vấn đề, cung cấp thông tin:
- Backend logs
- Print Bridge logs  
- Screenshot Print Jobs page
- Thông tin máy in (model, IP, port)

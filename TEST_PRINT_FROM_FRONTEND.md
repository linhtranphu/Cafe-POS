# Test Print từ Frontend

## Đã Thực Hiện

### Backend Changes
Đã sửa endpoint `/api/manager/printers/:id/test` để in ra một test receipt thay vì chỉ test connection.

**File**: `backend/interfaces/http/printer_config_handler.go`

**Changes**:
1. Thêm imports: `fmt`, `time`
2. Sửa `TestConnection` handler để:
   - Connect đến máy in
   - In ra test receipt với thông tin máy in
   - Trả về kết quả

**Test Receipt Content**:
```
================================
       TEST PRINT
================================
Printer: [Tên máy in]
Type: [BILL/LABEL]
Connection: [NETWORK/USB]
IP: [IP]:[Port]
================================
This is a test print
Đây là bản in thử nghiệm
================================
Date: [Timestamp]
================================
```

### Backend Status
✅ Backend đã restart và running on port 3000
✅ Endpoint `/api/manager/printers/:id/test` ready

## Cách Test

### 1. Mở Print Management
```
http://localhost:5173/#/print-management
```

### 2. Chuyển sang tab "Máy In"
Click vào tab "🖨️ Máy In"

### 3. Bấm nút "Test"
- Tìm máy in bạn muốn test (ví dụ: 192.168.1.115)
- Bấm nút "🔍 Test"
- Chờ vài giây

### 4. Kết Quả Mong Đợi

**Trên màn hình**:
- Hiển thị "✅ Test print successful - Check your printer"
- Hoặc lỗi nếu không kết nối được

**Trên máy in**:
- In ra một receipt test với thông tin máy in
- Có timestamp hiện tại
- Có text tiếng Việt

## Troubleshooting

### Lỗi: Connection Failed
```
Connection failed: dial tcp 192.168.1.115:9100: i/o timeout
```

**Nguyên nhân**:
- Máy in offline
- IP address sai
- Port sai
- Firewall block

**Giải pháp**:
```bash
# Test kết nối
nc -zv 192.168.1.115 9100

# Ping máy in
ping 192.168.1.115
```

### Lỗi: Print Test Failed
```
Print test failed: write error
```

**Nguyên nhân**:
- Máy in không hỗ trợ ESC/POS
- Máy in đang bận
- Giấy hết

**Giải pháp**:
- Kiểm tra máy in có giấy
- Kiểm tra máy in ready
- Thử lại

### Lỗi: Printer Not Found
```
Printer not found
```

**Nguyên nhân**:
- Chưa có máy in trong database
- ID sai

**Giải pháp**:
- Thêm máy in mới từ tab "Máy In"
- Click "➕ Thêm"

## API Endpoint

### POST /api/manager/printers/:id/test

**Request**:
```bash
curl -X POST http://localhost:3000/api/manager/printers/{printer_id}/test \
  -H "Authorization: Bearer {token}"
```

**Response (Success)**:
```json
{
  "success": true,
  "message": "Test print successful - Check your printer"
}
```

**Response (Error)**:
```json
{
  "success": false,
  "error": "Connection failed: ..."
}
```

## Files Changed

- `backend/interfaces/http/printer_config_handler.go`
  - Added `fmt` and `time` imports
  - Modified `TestConnection` handler to print test receipt

## Next Steps

1. ✅ Backend updated
2. ⏳ Test từ frontend
3. ⏳ Verify máy in thực sự in ra
4. ⏳ Test với nhiều loại máy in khác nhau

## Notes

- Test receipt sử dụng ESC/POS commands
- Hoạt động với cả BILL và LABEL printers
- Content có thể customize trong code nếu cần
- Timestamp dùng format: `2006-01-02 15:04:05`

---

**Status**: ✅ Ready to test from frontend!

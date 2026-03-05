# 🔍 Phân tích Backup Database

## File backup được kiểm tra
- **File**: `/Volumes/Linh-DAT/TaCafePOS/backups_from_ec2/mongodb-backup-2026-02-28_12-36-31.gz`
- **Kích thước**: 24KB
- **Ngày tạo**: 2026-02-28 12:36:31

## ⚠️ KẾT LUẬN QUAN TRỌNG

**Backup này là backup CŨ và THIẾU RẤT NHIỀU data quan trọng!**

### Collections có trong backup (13 collections):
✅ batch_definitions (7 docs)
✅ cashier_shifts (1 doc)
✅ expense_categories (5 docs)
✅ expenses (134 docs)
✅ facilities (41 docs)
✅ facility_history (81 docs)
✅ ingredient_categories (7 docs)
✅ ingredients (64 docs)
✅ menu_categories (4 docs)
✅ menu_items (15 docs)
✅ shifts (2 docs)
✅ stock_history (80 docs)
✅ users (5 docs)

### ❌ Collections THIẾU trong backup (13 collections):

#### 🚨 CRITICAL - Thiếu data kinh doanh:
- ❌ **orders** (87 docs) - THIẾU TẤT CẢ ĐƠN HÀNG!
- ❌ **order_items** (38 docs) - THIẾU CHI TIẾT ĐƠN HÀNG!
- ❌ **batch_records** (14 docs) - Thiếu lịch sử batch
- ❌ **batch_usage_logs** (77 docs) - Thiếu log sử dụng batch

#### 💰 Thiếu data tài chính:
- ❌ **cash_discrepancies** (3 docs)
- ❌ **cash_handovers** (29 docs)
- ❌ **fund_handovers** (6 docs)
- ❌ **fund_transactions** (4 docs)

#### 🖨️ Thiếu data in ấn (v2.0):
- ❌ **shop_settings** (1 doc) - CẦN TẠO BỞI MIGRATION
- ❌ **print_jobs** (166 docs)
- ❌ **print_notifications** (134 docs)
- ❌ **print_templates** (3 docs)
- ❌ **printer_configs** (2 docs)

## 📊 So sánh với database hiện tại

| Collection | Backup (EC2 cũ) | Current (Local) | Chênh lệch |
|------------|-----------------|-----------------|------------|
| orders | ❌ 0 | ✅ 87 | -87 |
| order_items | ❌ 0 | ✅ 38 | -38 |
| batch_records | ❌ 0 | ✅ 14 | -14 |
| batch_usage_logs | ❌ 0 | ✅ 77 | -77 |
| print_jobs | ❌ 0 | ✅ 166 | -166 |
| shop_settings | ❌ 0 | ✅ 1 | -1 |

## 🎯 KHUYẾN NGHỊ

### ⚠️ KHÔNG NÊN restore backup này vì:

1. **Mất toàn bộ đơn hàng** (87 orders)
2. **Mất toàn bộ dữ liệu tài chính** (cash handovers, fund transactions)
3. **Mất lịch sử batch** (batch records, usage logs)
4. **Thiếu cấu trúc v2.0** (shop_settings, print configs)

### ✅ NÊN LÀM:

#### Option 1: Tạo backup MỚI từ EC2 hiện tại
```bash
# Trên máy local
ssh -i TaCafePOS.pem ubuntu@47.128.65.142
cd ~/cafe-pos
docker exec cafe-pos-mongodb mongodump \
  --username admin \
  --password <password> \
  --authenticationDatabase admin \
  --db cafe_pos \
  --archive=/tmp/full-backup.gz \
  --gzip

# Download về
scp -i TaCafePOS.pem ubuntu@47.128.65.142:/tmp/full-backup.gz \
  /Volumes/Linh-DAT/TaCafePOS/backups_from_ec2/
```

#### Option 2: Nếu EC2 đã có data đầy đủ
- Kiểm tra EC2 có đầy đủ collections không
- Nếu có, tạo backup mới
- Nếu không, KHÔNG restore backup cũ này

#### Option 3: Nếu muốn deploy fresh
- Deploy code mới lên EC2
- Chạy migration để tạo cấu trúc v2.0
- Seed data mới (users, menu, ingredients)
- KHÔNG restore backup cũ

## 🔧 Cần sửa trong migrate-v2.0-simple.sh

Script migration hiện tại CHỈ tạo:
1. ✅ shop_settings collection
2. ✅ Indexes cho print_jobs, printer_configs, print_templates

**KHÔNG CẦN sửa gì thêm** vì:
- Các collections khác (orders, batch_records, etc.) sẽ được tạo tự động bởi backend khi có data
- Migration chỉ cần tạo shop_settings và indexes là đủ

## 🚨 CẢNH BÁO

**KHÔNG restore backup `mongodb-backup-2026-02-28_12-36-31.gz` này lên EC2 production!**

Backup này thiếu quá nhiều data quan trọng và sẽ gây mất dữ liệu nghiêm trọng.

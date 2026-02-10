# Hướng Dẫn Migration Menu Cost & Profit Analysis trên EC2

## Tổng Quan

Document này hướng dẫn chi tiết cách chạy migration cho tính năng Menu Cost & Profit Analysis trên EC2 server sau khi deploy Docker.

## Yêu Cầu

- ✅ Docker containers đã được deploy và đang chạy trên EC2
- ✅ Backend container có thể kết nối đến MongoDB
- ✅ SSH access vào EC2 instance
- ✅ Code mới nhất đã được pull/deploy

## Bước 1: Kết Nối SSH vào EC2

```bash
# Kết nối vào EC2 instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Hoặc nếu đã config SSH
ssh your-ec2-alias
```

## Bước 2: Kiểm Tra Docker Containers

```bash
# Kiểm tra containers đang chạy
docker ps

# Bạn sẽ thấy:
# - backend container
# - frontend container  
# - mongodb container (hoặc external MongoDB)
```

## Bước 3: Vào Backend Container

```bash
# Tìm backend container ID hoặc name
docker ps | grep backend

# Vào backend container
docker exec -it <backend-container-id> /bin/sh

# Hoặc nếu container name là "backend"
docker exec -it backend /bin/sh
```

## Bước 4: Kiểm Tra Môi Trường

Trong backend container:

```bash
# Kiểm tra Go đã được cài đặt
go version

# Kiểm tra cấu trúc thư mục
ls -la /app
ls -la /app/cmd/migrate

# Kiểm tra biến môi trường MongoDB
echo $MONGODB_URI
echo $MONGODB_DATABASE
```

## Bước 5: Chạy Migration Scripts

### 5.1. Schema Migration (Bắt buộc)

```bash
# Trong backend container
cd /app

# Chạy schema migration
go run cmd/migrate/run_all_menu_cost_migrations.go
```

**Output mong đợi:**
```
🔌 Connecting to MongoDB: mongodb://...
✅ Connected to database: cafe_pos

============================================================
  MENU COST & PROFIT ANALYSIS - SCHEMA MIGRATION
============================================================

📝 [1/5] Migrating menu_items collection...
   ✅ Added cost tracking fields to 50 menu items
   
📝 [2/5] Migrating ingredients collection...
   ✅ Added conversion and wastage fields to 30 ingredients
   
📝 [3/5] Creating order_items collection...
   ✅ Created order_items collection
   
📝 [4/5] Creating operating_expenses collection...
   ✅ Created operating_expenses collection
   
📝 [5/5] Migrating shop_settings collection...
   ✅ Added low_margin_threshold to 1 shop settings

📝 Creating indexes for all collections...
   ✅ All indexes created

============================================================
  ✅ ALL MIGRATIONS COMPLETED SUCCESSFULLY
============================================================
```

### 5.2. Backfill Menu Item Costs (Bắt buộc)

```bash
# Tính toán current_cost cho tất cả menu items
go run cmd/migrate/backfill_menu_item_costs.go
```

**Output mong đợi:**
```
============================================================
  BACKFILL CURRENT_COST FOR MENU ITEMS
============================================================

📝 Fetching all menu items...
   Found 50 menu items

🔄 Calculating costs for menu items...
   Progress: 10/50 items processed
   Progress: 20/50 items processed
   Progress: 30/50 items processed
   Progress: 40/50 items processed
   Progress: 50/50 items processed

============================================================
  ✅ BACKFILL COMPLETED
============================================================

📊 Backfill Summary:
   Total menu items:           50
   ✅ Successfully calculated:  42
   📦 No ingredients (cost=0):  3
   ⚠️  Incomplete (missing cost): 5
```

### 5.3. Backfill Order Item Costs (Tùy chọn)

**Chỉ chạy nếu bạn có dữ liệu lịch sử (closed shifts)**

```bash
# Tính toán accounting_cost cho orders trong closed shifts
go run cmd/migrate/backfill_order_item_costs.go
```

**Output mong đợi:**
```
============================================================
  BACKFILL ACCOUNTING_COST FOR HISTORICAL ORDERS
============================================================

📝 Fetching closed shifts...
   Found 25 closed shifts

🔄 Processing shifts and calculating costs...
   Progress: 10/25 shifts processed (150 orders, 450 items)
   Progress: 20/25 shifts processed (300 orders, 900 items)
   Progress: 25/25 shifts processed (375 orders, 1125 items)

============================================================
  ✅ BACKFILL COMPLETED
============================================================

ℹ️  Important Notes:
   • All backfilled costs are marked as ESTIMATED (not FINAL)
   • Future shift closures will use FINAL status
```

### 5.4. Verify Migration (Bắt buộc)

```bash
# Kiểm tra migration đã hoàn thành đúng
go run cmd/migrate/verify_menu_cost_migration.go
```

**Output mong đợi:**
```
============================================================
  MENU COST & PROFIT ANALYSIS - MIGRATION VERIFICATION
============================================================

📝 [1/8] Verifying menu_items schema...
   ✅ All 50 menu items have cost tracking fields

📝 [2/8] Verifying ingredients schema...
   ✅ All 30 ingredients have conversion and wastage fields

📝 [3/8] Verifying order_items collection...
   ✅ order_items collection exists with 1125 documents

📝 [4/8] Verifying operating_expenses collection...
   ✅ operating_expenses collection exists with 0 documents

📝 [5/8] Verifying shop_settings schema...
   ✅ All 1 shop settings have low_margin_threshold field

📝 [6/8] Verifying indexes...
   ✅ All required indexes exist

📝 [7/8] Verifying menu item costs...
   ✅ All menu items have valid cost_status

📝 [8/8] Verifying order item costs...
   ✅ All order items have valid cost_status

============================================================
  ✅ ALL VERIFICATION CHECKS PASSED
============================================================
```

## Bước 6: Thoát Container và Restart Services

```bash
# Thoát khỏi backend container
exit

# Restart backend container để áp dụng thay đổi
docker restart <backend-container-id>

# Hoặc restart tất cả services
docker-compose restart
```

## Bước 7: Kiểm Tra Application

```bash
# Kiểm tra backend logs
docker logs -f <backend-container-id>

# Kiểm tra API endpoint mới
curl http://localhost:8080/api/menu/costs

# Hoặc từ máy local
curl http://your-ec2-ip:8080/api/menu/costs
```

## Script Tự Động (Khuyến Nghị)

Tạo script để chạy tất cả migrations một lần:

```bash
# Trên EC2, tạo file migrate.sh
cat > migrate.sh << 'EOF'
#!/bin/bash
set -e

echo "🚀 Starting Menu Cost Migration..."

# Vào backend container và chạy migrations
docker exec backend sh -c "
  cd /app && \
  echo '📝 Step 1: Schema Migration...' && \
  go run cmd/migrate/run_all_menu_cost_migrations.go && \
  echo '' && \
  echo '📝 Step 2: Backfill Menu Item Costs...' && \
  go run cmd/migrate/backfill_menu_item_costs.go && \
  echo '' && \
  echo '📝 Step 3: Backfill Order Item Costs...' && \
  go run cmd/migrate/backfill_order_item_costs.go && \
  echo '' && \
  echo '📝 Step 4: Verify Migration...' && \
  go run cmd/migrate/verify_menu_cost_migration.go
"

if [ $? -eq 0 ]; then
  echo ""
  echo "✅ Migration completed successfully!"
  echo "🔄 Restarting backend container..."
  docker restart backend
  echo "✅ Done!"
else
  echo ""
  echo "❌ Migration failed! Please check the logs above."
  exit 1
fi
EOF

# Cho phép execute
chmod +x migrate.sh

# Chạy migration
./migrate.sh
```

## Xử Lý Lỗi Thường Gặp

### Lỗi 1: "go: command not found"

**Nguyên nhân**: Go chưa được cài trong container

**Giải pháp**: Đảm bảo Dockerfile có cài Go:
```dockerfile
FROM golang:1.21-alpine
# ... rest of Dockerfile
```

### Lỗi 2: "Failed to connect to MongoDB"

**Nguyên nhân**: Backend không kết nối được MongoDB

**Giải pháp**:
```bash
# Kiểm tra MongoDB URI
docker exec backend env | grep MONGODB

# Kiểm tra MongoDB container
docker ps | grep mongo

# Test kết nối
docker exec backend sh -c "nc -zv mongodb 27017"
```

### Lỗi 3: "Some menu items have INCOMPLETE status"

**Nguyên nhân**: Một số ingredients thiếu cost_per_unit

**Giải pháp**: 
- Đây là warning, không phải lỗi
- Cập nhật ingredient costs trong application
- Chạy lại backfill nếu cần:
```bash
docker exec backend go run cmd/migrate/backfill_menu_item_costs.go
```

### Lỗi 4: "No closed shifts found"

**Nguyên nhân**: Chưa có shifts nào được đóng

**Giải pháp**:
- Đây là normal nếu hệ thống mới
- Bỏ qua bước backfill order costs
- Costs sẽ được tính tự động khi đóng shifts trong tương lai

## Rollback (Nếu Cần)

Nếu cần rollback migration:

```bash
# Vào backend container
docker exec -it backend /bin/sh

# Kết nối MongoDB shell
# (Nếu có mongo client trong container)
mongo $MONGODB_URI

# Hoặc từ EC2 host
docker exec -it mongodb mongo cafe_pos

# Xóa các fields mới
db.menu_items.updateMany({}, {
  $unset: {
    current_cost: "",
    cost_last_calculated_at: "",
    cost_status: ""
  }
})

db.ingredients.updateMany({}, {
  $unset: {
    conversion_rate: "",
    wastage_percentage: ""
  }
})

db.shop_settings.updateMany({}, {
  $unset: {
    low_margin_threshold: ""
  }
})

# Xóa collections mới
db.order_items.drop()
db.operating_expenses.drop()
```

## Backup Trước Khi Migration

**QUAN TRỌNG**: Luôn backup database trước khi chạy migration!

```bash
# Backup MongoDB từ EC2
docker exec mongodb mongodump \
  --db cafe_pos \
  --out /backup/before-menu-cost-migration

# Copy backup ra host
docker cp mongodb:/backup/before-menu-cost-migration ./backup/

# Hoặc backup toàn bộ
docker exec mongodb mongodump --out /backup/full-backup-$(date +%Y%m%d)
```

## Restore Backup (Nếu Cần)

```bash
# Restore từ backup
docker exec mongodb mongorestore \
  --db cafe_pos \
  /backup/before-menu-cost-migration/cafe_pos
```

## Monitoring Sau Migration

```bash
# Theo dõi backend logs
docker logs -f backend

# Kiểm tra MongoDB performance
docker exec mongodb mongo --eval "db.stats()"

# Kiểm tra disk usage
docker exec mongodb df -h
```

## Checklist Hoàn Thành

- [ ] Đã backup database
- [ ] Đã chạy schema migration thành công
- [ ] Đã backfill menu item costs
- [ ] Đã backfill order item costs (nếu có data)
- [ ] Đã verify migration thành công
- [ ] Đã restart backend container
- [ ] Đã test API endpoints mới
- [ ] Đã kiểm tra manager interface
- [ ] Đã thông báo team về tính năng mới

## Thời Gian Ước Tính

Cho database trung bình:
- Schema migration: ~5 giây
- Menu item backfill (50 items): ~1 giây
- Order item backfill (1000 orders): ~20 giây
- Verification: ~2 giây
- **Tổng: ~30 giây**

## Liên Hệ Support

Nếu gặp vấn đề:
1. Kiểm tra logs: `docker logs backend`
2. Kiểm tra MongoDB: `docker logs mongodb`
3. Xem chi tiết trong các README files:
   - `backend/cmd/migrate/README_TASK_19.md`
   - `backend/cmd/migrate/TASK_19_IMPLEMENTATION_SUMMARY.md`

## Tính Năng Mới Sau Migration

Sau khi migration thành công, các tính năng sau sẽ có sẵn:

✅ **Menu Cost View** - Xem chi phí và lợi nhuận từng món
✅ **Cost Breakdown** - Chi tiết giá nguyên liệu từng món
✅ **Profit Analysis** - Phân tích lợi nhuận theo category
✅ **Operating Profit** - Lợi nhuận sau trừ chi phí vận hành
✅ **Warning Detection** - Cảnh báo món bán lỗ hoặc lợi nhuận thấp
✅ **Automatic Cost Calculation** - Tự động tính cost khi đóng ca

## Notes Quan Trọng

⚠️ **Migration là idempotent** - Có thể chạy nhiều lần an toàn
⚠️ **Không downtime** - Application vẫn chạy bình thường trong khi migration
⚠️ **Backfilled costs = ESTIMATED** - Costs lịch sử dùng giá hiện tại, không phải giá thực tế lúc đó
⚠️ **Future costs = FINAL** - Costs từ giờ trở đi sẽ chính xác 100%

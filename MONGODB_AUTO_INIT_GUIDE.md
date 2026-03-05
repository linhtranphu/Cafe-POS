# MongoDB Auto-Initialization Guide

## Overview
MongoDB container sẽ tự động khởi tạo database structure khi start lần đầu tiên, không cần chạy migration script thủ công.

## Cách hoạt động

### 1. Init Script
File `mongodb-init/init-db.js` chứa script khởi tạo:
- Tạo tất cả collections cần thiết
- Tạo indexes cho performance
- Tạo shop_settings với giá trị mặc định (bao gồm `print_bridge_url`)

### 2. Docker Volume Mount
Docker Compose mount folder `mongodb-init` vào `/docker-entrypoint-initdb.d/`:
```yaml
volumes:
  - ./mongodb-init:/docker-entrypoint-initdb.d:ro
```

MongoDB sẽ tự động chạy tất cả `.js` và `.sh` files trong folder này khi container start lần đầu.

## Sử dụng

### Fresh Start (Database mới)

```bash
# 1. Xóa volume cũ (nếu muốn reset hoàn toàn)
docker-compose down -v

# 2. Start lại - init script sẽ tự động chạy
docker-compose up -d mongodb

# 3. Kiểm tra logs
docker logs cafe-pos-mongodb

# 4. Verify database structure
docker exec -it cafe-pos-mongodb mongosh -u admin -p password123 --authenticationDatabase admin cafe_pos --eval "db.getCollectionNames()"
```

### Update Existing Database

Nếu database đã tồn tại, init script sẽ KHÔNG chạy lại. Có 2 options:

**Option 1: Chạy migration script (Recommended cho production)**
```bash
sudo bash sync-db-structure.sh
```

**Option 2: Reset database (CHỈ dùng cho development)**
```bash
# ⚠️ CẢNH BÁO: Sẽ XÓA TẤT CẢ DỮ LIỆU
docker-compose down -v
docker-compose up -d
```

## Database Structure

### Collections Created
- users
- menu_items, menu_categories
- orders, order_items
- ingredients, ingredient_categories
- batch_definitions, batch_records, batch_usage_logs
- expenses, expense_categories
- printer_configs, print_templates, print_jobs, print_notifications
- shop_settings
- shifts, cashier_shifts
- cash_handovers, cash_discrepancies
- fund_transactions, fund_handovers
- stock_history

### Indexes Created
- orders: created_at, status, order_number (unique)
- print_jobs: status+created_at, order_id, printer_id
- printer_configs: type+is_default
- print_templates: type+is_default
- users: username (unique), role
- menu_items: name, category_id
- ingredients: name
- batch_records: ingredient_id, created_at
- expenses: date, category_id

### Default Shop Settings
```javascript
{
  shop_name: "Cafe POS",
  shop_address: "",
  shop_phone: "",
  logo_url: "",
  custom_message: "Cảm ơn quý khách! Hẹn gặp lại!",
  print_bridge_url: "http://localhost:3001",
  show_logo: false,
  show_address: true,
  show_phone: true,
  show_custom_message: true,
  low_margin_threshold: 20.0,
  auto_print_enabled: true
}
```

## Deployment

### Development
```bash
# Clone repo
git clone <repo-url>
cd cafe-pos

# Start services (init script runs automatically)
docker-compose up -d

# Check initialization
docker logs cafe-pos-mongodb | grep "initialization"
```

### Production (EC2)
```bash
# 1. Copy init script to server
scp -r mongodb-init ubuntu@your-ec2:/home/ubuntu/cafe-pos/

# 2. Update docker-compose.yml on server
# (Ensure it has the volume mount)

# 3. For existing database, run migration
sudo bash sync-db-structure.sh

# 4. For fresh database
docker-compose down -v
docker-compose up -d
```

## Troubleshooting

### Init script không chạy
**Nguyên nhân**: Database volume đã tồn tại
**Giải pháp**: 
```bash
# Check if volume exists
docker volume ls | grep mongodb

# Remove volume (⚠️ loses data)
docker-compose down -v

# Or run migration script instead
sudo bash sync-db-structure.sh
```

### Kiểm tra init script đã chạy chưa
```bash
# Check logs
docker logs cafe-pos-mongodb | grep "initialization"

# Check collections
docker exec -it cafe-pos-mongodb mongosh -u admin -p password123 --authenticationDatabase admin cafe_pos --eval "db.getCollectionNames().length"

# Should return 24 collections
```

### Update init script
```bash
# 1. Edit mongodb-init/init-db.js
# 2. Remove old volume
docker-compose down -v
# 3. Restart
docker-compose up -d
```

## Benefits

✅ Không cần chạy migration script thủ công
✅ Database structure nhất quán giữa dev/staging/production
✅ Dễ dàng setup môi trường mới
✅ Version control cho database schema
✅ Tự động tạo indexes cho performance
✅ Default data sẵn sàng sử dụng

## Notes

- Init script CHỈ chạy khi database chưa tồn tại
- Để update existing database, dùng migration scripts
- Backup data trước khi xóa volumes
- Init script chạy với admin privileges
- Có thể thêm nhiều `.js` files vào `mongodb-init/` folder

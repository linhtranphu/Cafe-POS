# Đóng tất cả ca - Hướng dẫn nhanh

## Cách 1: Sử dụng Python Script (Khuyến nghị)

### Bước 1: Cài pymongo
```bash
pip3 install pymongo
```

### Bước 2: Chạy script
```bash
python3 close_all_shifts.py
```

Script sẽ:
- Tìm tất cả ca OPEN
- Đóng từng ca
- Tính tổng doanh thu
- Lock orders đã hoàn thành
- Hiển thị summary

---

## Cách 2: Update trực tiếp MongoDB (Nhanh nhất)

### Option A: Sử dụng MongoDB Compass

1. Mở MongoDB Compass
2. Connect: `mongodb://admin:password123@localhost:27017/?authSource=admin`
3. Chọn database `cafe_pos`
4. Chọn collection `shifts`
5. Filter: `{status: "OPEN"}`
6. Chọn tất cả documents
7. Update:
```json
{
  "$set": {
    "status": "CLOSED",
    "ended_at": "$$NOW",
    "updated_at": "$$NOW"
  }
}
```

### Option B: Sử dụng mongosh/mongo command

```bash
# Nếu có mongosh
mongosh cafe_pos -u admin -p password123 --authenticationDatabase admin

# Hoặc mongo
mongo cafe_pos -u admin -p password123 --authenticationDatabase admin
```

Sau đó chạy:
```javascript
// Đóng tất cả ca
db.shifts.updateMany(
  { status: "OPEN" },
  { 
    $set: { 
      status: "CLOSED",
      ended_at: new Date(),
      updated_at: new Date()
    }
  }
)

// Kiểm tra
db.shifts.find({ status: "OPEN" }).count()
// Nên trả về 0
```

### Option C: Sử dụng Docker (nếu MongoDB chạy trong Docker)

```bash
# Vào container MongoDB
docker exec -it <mongodb_container_name> mongosh -u admin -p password123 --authenticationDatabase admin cafe_pos

# Chạy lệnh update
db.shifts.updateMany(
  { status: "OPEN" },
  { 
    $set: { 
      status: "CLOSED",
      ended_at: new Date(),
      updated_at: new Date()
    }
  }
)
```

---

## Cách 3: Sử dụng MongoDB Script File

### Bước 1: Tạo file update.js
```javascript
// update.js
db = db.getSiblingDB('cafe_pos');

const result = db.shifts.updateMany(
  { status: "OPEN" },
  { 
    $set: { 
      status: "CLOSED",
      ended_at: new Date(),
      updated_at: new Date()
    }
  }
);

print(`Updated ${result.modifiedCount} shift(s)`);
```

### Bước 2: Chạy script
```bash
mongosh cafe_pos -u admin -p password123 --authenticationDatabase admin < update.js
```

---

## Kiểm tra sau khi đóng

### Xem số ca đang mở:
```javascript
db.shifts.find({ status: "OPEN" }).count()
```

### Xem ca vừa đóng:
```javascript
db.shifts.find({ 
  status: "CLOSED",
  ended_at: { $gte: new Date(Date.now() - 3600000) } // Last 1 hour
}).pretty()
```

### Xem chi tiết một ca:
```javascript
db.shifts.findOne({ _id: ObjectId("your_shift_id") })
```

---

## Lệnh MongoDB hữu ích

### Xem tất cả ca:
```javascript
db.shifts.find().sort({ started_at: -1 }).limit(10)
```

### Đếm ca theo status:
```javascript
db.shifts.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])
```

### Xem ca có vấn đề (tiền âm):
```javascript
db.shifts.find({
  $or: [
    { remaining_cash: { $lt: 0 } },
    { remaining_transfer: { $lt: 0 } }
  ]
})
```

### Sửa ca có tiền âm:
```javascript
db.shifts.updateOne(
  { _id: ObjectId("shift_id") },
  { 
    $set: { 
      remaining_cash: 0,
      remaining_transfer: 0
    }
  }
)
```

---

## Troubleshooting

### Lỗi: Authentication failed
**Giải pháp:** Kiểm tra username/password trong `.env`

### Lỗi: Connection refused
**Giải pháp:** 
1. Kiểm tra MongoDB đang chạy: `docker ps` hoặc `brew services list`
2. Kiểm tra port: `lsof -i :27017`

### Lỗi: pymongo not found
**Giải pháp:** `pip3 install pymongo`

### Ca đã đóng nhưng UI vẫn hiển thị OPEN
**Giải pháp:** 
1. Refresh browser (Cmd+R hoặc F5)
2. Clear cache
3. Kiểm tra lại trong MongoDB

---

## Backup trước khi đóng (Khuyến nghị)

```bash
# Backup toàn bộ database
mongodump --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" --out=backup_$(date +%Y%m%d_%H%M%S)

# Restore nếu cần
mongorestore --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" backup_20240223_120000/cafe_pos
```

---

## Lệnh nhanh (Copy & Paste)

### Đóng tất cả ca (MongoDB):
```javascript
db.shifts.updateMany({status:"OPEN"},{$set:{status:"CLOSED",ended_at:new Date(),updated_at:new Date()}})
```

### Kiểm tra:
```javascript
db.shifts.find({status:"OPEN"}).count()
```

### Xem ca vừa đóng:
```javascript
db.shifts.find({status:"CLOSED"}).sort({ended_at:-1}).limit(5)
```

# MongoDB Graceful Handling - Implementation Summary

## ✅ Đã hoàn thành

### 1. Created Helper Functions (`query_helpers.go`)
- `WithQueryTimeout()` - Thêm 5s timeout cho mọi query
- `IsCollectionNotFoundError()` - Check collection không tồn tại
- `SafeFindAll()` - Generic Find với timeout và error handling
- `SafeFindOne()` - Generic FindOne với timeout
- `SafeCount()` - Count với timeout
- `SafeInsertOne()`, `SafeUpdateOne()`, `SafeDeleteOne()` - CRUD với timeout

### 2. Updated Repositories

#### ✅ print_job_repository.go
- `FindPending()` - Thêm timeout 5s, return empty array nếu collection không tồn tại
- `FindFailed()` - Thêm timeout 5s, return empty array nếu collection không tồn tại  
- `FindByOrderID()` - Thêm timeout 5s, return empty array nếu collection không tồn tại

#### ✅ order_repository.go (HOÀN TẤT)
- `Create()` - Thêm timeout 5s
- `FindByID()` - Thêm timeout 5s
- `Update()` - Thêm timeout 5s
- `FindByShiftID()` - Thêm timeout 5s + graceful handling
- `FindByWaiterID()` - Thêm timeout 5s + graceful handling
- `FindByStatus()` - Thêm timeout 5s + graceful handling
- `FindAll()` - Thêm timeout 5s + graceful handling
- `FindByOrderNumber()` - Thêm timeout 5s
- `Delete()` - Thêm timeout 5s

## 📋 Pattern để apply cho các repository còn lại

### Pattern 1: Find Methods (trả về array)
```go
// BEFORE
func (r *Repository) FindXXX(ctx context.Context, ...) ([]*Type, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*Type
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// AFTER
func (r *Repository) FindXXX(ctx context.Context, ...) ([]*Type, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*Type{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*Type
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*Type{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*Type{}
	}
	return items, nil
}
```

### Pattern 2: FindOne Methods
```go
// BEFORE
func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Type, error) {
	var item Type
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// AFTER
func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Type, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	var item Type
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
```

### Pattern 3: Create/Insert Methods
```go
// BEFORE
func (r *Repository) Create(ctx context.Context, item *Type) error {
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// AFTER
func (r *Repository) Create(ctx context.Context, item *Type) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
```

### Pattern 4: Update Methods
```go
// BEFORE
func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, item *Type) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}

// AFTER
func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, item *Type) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}
```

### Pattern 5: Delete Methods
```go
// BEFORE
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AFTER
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
```

## 🎯 Repositories cần update (theo độ ưu tiên)

### Priority 1 - Critical (được gọi nhiều nhất)
- ✅ order_repository.go - DONE
- ✅ print_job_repository.go - DONE
- ⏳ menu_repository.go
- ⏳ ingredient_repository.go
- ⏳ user_repository.go
- ⏳ shift_repository.go
- ⏳ cashier_shift_repository.go

### Priority 2 - Important
- ⏳ batch_definition_repository.go
- ⏳ batch_record_repository.go
- ⏳ batch_usage_log_repository.go
- ⏳ order_item_repository.go
- ⏳ printer_config_repository.go
- ⏳ print_template_repository.go
- ⏳ shop_settings_repository.go

### Priority 3 - Normal
- ⏳ expense_repository.go
- ⏳ fund_transaction_repository.go
- ⏳ cash_handover_repository.go
- ⏳ stock_history_repository.go
- ⏳ menu_category_repository.go
- ⏳ facility_repository.go
- ⏳ operating_expense_repository.go

### Priority 4 - Low (ít được gọi)
- ⏳ print_notification_repository.go
- ⏳ cash_discrepancy_repository.go
- ⏳ fund_handover_repository.go
- ⏳ payment_audit_repository.go
- ⏳ payment_discrepancy_repository.go
- ⏳ cash_reconciliation_repository.go

## 🚀 Benefits

1. **Timeout Protection**: Mọi query có timeout 5s, tránh hang forever
2. **Graceful Degradation**: Trả về empty array thay vì error khi collection không tồn tại
3. **Better UX**: Frontend không bị crash, hiển thị "No data" thay vì error
4. **MongoDB Protection**: Tránh overload database với long-running queries
5. **Consistent Error Handling**: Tất cả repositories xử lý lỗi giống nhau

## 📝 Next Steps

Để hoàn thành việc update tất cả repositories:

1. Apply Pattern 1-5 cho từng repository theo priority
2. Test từng repository sau khi update
3. Build và deploy backend
4. Monitor logs để đảm bảo không có regression

## 🔧 Quick Apply Script

Có thể tạo script để tự động apply patterns:

```bash
# Find all Find methods
grep -r "func.*Find.*context.Context" backend/infrastructure/mongodb/*_repository.go

# Find all Create methods  
grep -r "func.*Create.*context.Context" backend/infrastructure/mongodb/*_repository.go

# Find all Update methods
grep -r "func.*Update.*context.Context" backend/infrastructure/mongodb/*_repository.go
```

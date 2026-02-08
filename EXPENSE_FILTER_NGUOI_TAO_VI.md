# ✅ Expense: Filter Theo Người Tạo

**Ngày:** 6 tháng 2, 2026

---

## 🎯 Yêu Cầu

> "edit màn hình expense, bên dưới phần 'Search Bar', thay vì hiển thị source type filter. hiển thị danh sách người tạo expense"

---

## ✅ Đã Thay Đổi

### Trước:

```
Filter theo Source Type:
[Tất cả] [✍️ Thủ công] [🥬 Nguyên liệu] 
[🏢 Cơ sở vật chất] [🔧 Bảo trì]
```

### Sau:

```
Filter theo Người Tạo:
[👥 Tất cả] [👤 Admin] [👤 Manager] 
[👤 Cashier] [👤 Hệ thống]
```

---

## 🔧 Cách Hoạt Động

### 1. Danh Sách Người Tạo (Dynamic):

```javascript
// Tự động lấy danh sách người tạo từ expenses
const uniqueCreators = computed(() => {
  const creators = expenses.value
    .map(e => e.created_by || 'Hệ thống')
    .filter((value, index, self) => self.indexOf(value) === index)
  return creators.sort() // Sort A-Z
})
```

**Ví dụ:**
```
Expenses có:
- Admin (3 expenses)
- Manager (2 expenses)
- Hệ thống (5 expenses)
- Cashier (1 expense)

→ Hiển thị buttons:
[👥 Tất cả] [👤 Admin] [👤 Cashier] 
[👤 Hệ thống] [👤 Manager]
```

### 2. Filter Logic:

```javascript
// Click vào người tạo → Chỉ hiển thị expenses của người đó
if (creatorFilter.value) {
  filtered = filtered.filter(e => {
    const creator = e.created_by || 'Hệ thống'
    return creator === creatorFilter.value
  })
}
```

---

## 📱 Giao Diện

### Filter Buttons:

```
┌─────────────────────────────────────────┐
│ Search: [Tìm kiếm chi phí...]           │
├─────────────────────────────────────────┤
│ [👥 Tất cả] [👤 Admin] [👤 Cashier]    │
│ [👤 Hệ thống] [👤 Manager]              │
└─────────────────────────────────────────┘
     ↑ Active         ↑ Inactive
   (purple)          (white)
```

### Khi Click [👤 Admin]:

```
┌─────────────────────────────────┐
│ Mua cà phê                      │
│ 👤 Admin                        │
│ -500,000 ₫                      │
├─────────────────────────────────┤
│ Sửa máy pha                     │
│ 👤 Admin                        │
│ -1,200,000 ₫                    │
└─────────────────────────────────┘
Chỉ hiển thị expenses của Admin
```

### Khi Click [👥 Tất cả]:

```
┌─────────────────────────────────┐
│ Mua cà phê                      │
│ 👤 Admin                        │
├─────────────────────────────────┤
│ Nhập nguyên liệu                │
│ 👤 Hệ thống                     │
├─────────────────────────────────┤
│ Tiền điện                       │
│ 👤 Manager                      │
└─────────────────────────────────┘
Hiển thị tất cả expenses
```

---

## 🎨 Tính Năng

### Dynamic (Tự Động):

- ✅ Danh sách người tạo tự động cập nhật
- ✅ Không cần config thủ công
- ✅ Thêm expense mới → Người tạo mới tự động xuất hiện

### Unique (Không Trùng):

- ✅ Mỗi người chỉ hiện 1 lần
- ✅ Dù có tạo nhiều expenses

### Sorted (Sắp Xếp):

- ✅ Sort theo alphabet (A-Z)
- ✅ Dễ tìm

### Fallback:

- ✅ Nếu `created_by` = null → Hiển thị "Hệ thống"
- ✅ Không bị lỗi

---

## 🧪 Test

### Cần Test:

**1. Hiển Thị Danh Sách:**
```
1. Mở màn hình expense
2. Check filter buttons
3. ✅ Phải hiển thị tất cả người tạo (unique)
4. ✅ Sort A-Z
5. ✅ Có button "Tất cả"
```

**2. Filter Theo Người:**
```
1. Click vào 1 người (vd: Admin)
2. Check danh sách
3. ✅ Chỉ hiển thị expenses của người đó
4. ✅ Expenses của người khác bị ẩn
```

**3. Hiển Thị Tất Cả:**
```
1. Click "Tất cả"
2. Check danh sách
3. ✅ Hiển thị tất cả expenses
4. ✅ Không filter
```

**4. Tạo Expense Mới:**
```
1. Tạo expense mới
2. Check filter buttons
3. ✅ Người tạo mới xuất hiện (nếu chưa có)
```

---

## 💡 Lợi Ích

### Trước (Source Type):

- ❌ Filter theo loại (manual, ingredient, facility...)
- ❌ Không biết ai tạo
- ❌ Fixed categories

### Sau (Creator):

- ✅ Filter theo người tạo
- ✅ Dễ track chi phí của từng người
- ✅ Dynamic, tự động cập nhật
- ✅ Tốt cho accountability

### Use Cases:

1. **Manager muốn xem chi phí của nhân viên:**
   - Click tên nhân viên → Xem expenses của họ

2. **Check expenses tự động:**
   - Click "Hệ thống" → Xem expenses do hệ thống tạo

3. **Audit chi phí:**
   - Filter từng người → Review chi tiêu

---

## 📁 File Đã Sửa

**frontend/src/views/ExpenseManagementView.vue**

**Thay đổi:**
- Template: Thay source type filter → creator filter
- Script: Thêm `uniqueCreators` computed
- Script: Đổi `sourceFilter` → `creatorFilter`
- Script: Update filter logic

---

## ✅ Kết Quả

**Trước:**
- ❌ Filter theo source type (manual, ingredient...)
- ❌ Không track được người tạo

**Sau:**
- ✅ Filter theo người tạo
- ✅ Dynamic list (tự động)
- ✅ Sort A-Z
- ✅ Không trùng lặp
- ✅ Dễ track chi phí

---

## 🎯 Summary

**Vấn đề:** Source type filter không hữu ích cho tracking người tạo

**Giải pháp:** Thay bằng filter theo người tạo (dynamic)

**Kết quả:**
- ✅ Hiển thị tất cả người tạo (unique)
- ✅ Filter expenses theo người
- ✅ Tự động cập nhật
- ✅ Sort A-Z

**Status:** ✅ Hoàn thành, sẵn sàng test!

---

**Ngày:** 6 tháng 2, 2026  
**Developer:** Kiro AI  
**Trạng thái:** ✅ Done

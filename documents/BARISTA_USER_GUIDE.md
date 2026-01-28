# 🍹 Hướng dẫn sử dụng cho Barista

## 🔐 Đăng nhập

### Tài khoản mặc định
```
Username: barista1
Password: barista123

Username: barista2
Password: barista123
```

## 📱 Giao diện Barista

### Bottom Navigation
```
🏠 Trang chủ | 🍹 Barista | ⏰ Ca làm | 👤 Cá nhân
```

### 3 Tabs chính

#### 1. ⏳ Queue (Hàng đợi)
- Hiển thị orders đã được waiter gửi đến
- Status: QUEUED
- Chờ barista nhận

#### 2. 🍹 Đang pha
- Orders mà bạn đã nhận
- Status: IN_PROGRESS
- Đang trong quá trình pha chế

#### 3. ✅ Sẵn sàng
- Orders đã pha xong
- Status: READY
- Chờ waiter giao cho khách

## 🔄 Workflow

### Bước 1: Xem Queue
```
1. Mở app → Tap "🍹 Barista"
2. Tab "⏳ Queue" hiển thị orders chờ pha
3. Xem thông tin:
   - Order number
   - Tên khách
   - Danh sách món
   - Ghi chú (nếu có)
   - Thời gian vào queue
```

### Bước 2: Nhận Order
```
1. Chọn order từ queue
2. Đọc kỹ danh sách món và ghi chú
3. Tap "👍 Nhận order"
4. Order chuyển sang tab "🍹 Đang pha"
```

### Bước 3: Pha chế
```
1. Chuyển sang tab "🍹 Đang pha"
2. Xem orders đang làm
3. Hiển thị thời gian đã pha (real-time)
4. Pha chế theo đúng công thức
```

### Bước 4: Hoàn tất
```
1. Khi pha xong, tap "✅ Hoàn tất"
2. Order chuyển sang tab "✅ Sẵn sàng"
3. Waiter sẽ đến lấy và giao cho khách
```

## 🎯 Tips & Best Practices

### Quản lý Queue hiệu quả
1. **Ưu tiên FIFO**: Nhận order theo thứ tự vào queue
2. **Kiểm tra ghi chú**: Đọc kỹ yêu cầu đặc biệt
3. **Nhận vừa phải**: Không nhận quá nhiều cùng lúc
4. **Communicate**: Nếu thiếu nguyên liệu, báo ngay

### Pha chế chất lượng
1. **Đúng công thức**: Theo standard recipe
2. **Đúng nhiệt độ**: Espresso 90-96°C, sữa 60-65°C
3. **Đúng tỷ lệ**: Cân đo chính xác
4. **Presentation**: Latte art, garnish đẹp mắt

### Tốc độ
1. **Espresso**: 25-30 giây
2. **Cappuccino/Latte**: 2-3 phút
3. **Trà**: 3-5 phút
4. **Smoothie**: 2-3 phút

### Vệ sinh
1. **Sau mỗi shot**: Lau group head
2. **Sau mỗi ly sữa**: Xả steam wand
3. **Mỗi giờ**: Backflush máy
4. **Cuối ca**: Vệ sinh toàn bộ

## 📊 Hiển thị thông tin

### Order Card
```
┌─────────────────────────────┐
│ #ORD-001        [⏳ Chờ pha]│
│ Nguyễn Văn A                │
│ 14:30                       │
│                             │
│ ☕ Cappuccino x2             │
│ 🍵 Trà đào x1                │
│                             │
│ 📝 Ít đường, nhiều đá       │
│                             │
│ [👍 Nhận order]             │
└─────────────────────────────┘
```

### Working Card
```
┌─────────────────────────────┐
│ #ORD-001        [🍹 Đang pha]│
│ Nguyễn Văn A                │
│ Bắt đầu: 14:30              │
│ ⏱️ 2m 15s                    │
│                             │
│ ☕ Cappuccino x2             │
│ 🍵 Trà đào x1                │
│                             │
│ 📝 Ít đường, nhiều đá       │
│                             │
│ [✅ Hoàn tất]               │
└─────────────────────────────┘
```

### Ready Card
```
┌─────────────────────────────┐
│ #ORD-001        [✅ Sẵn sàng]│
│ Nguyễn Văn A                │
│ Hoàn tất: 14:35             │
│ ⏱️ Chờ giao: 2 phút         │
│                             │
│ ☕ Cappuccino x2             │
│ 🍵 Trà đào x1                │
│                             │
│ 🎉 Chờ waiter giao cho khách│
└─────────────────────────────┘
```

## 🔔 Auto Refresh

App tự động refresh mỗi 10 giây để cập nhật:
- Orders mới vào queue
- Thời gian pha chế
- Status changes

Bạn cũng có thể tap nút 🔄 để refresh thủ công.

## ⚠️ Xử lý tình huống

### Thiếu nguyên liệu
```
1. Kiểm tra kho ngay
2. Báo cho manager
3. Thông báo waiter để inform khách
4. Suggest món thay thế
```

### Máy móc hỏng
```
1. Báo ngay cho manager
2. Chuyển sang máy backup (nếu có)
3. Điều chỉnh workflow
4. Update waiter về thời gian chờ
```

### Order phức tạp
```
1. Đọc kỹ ghi chú
2. Hỏi waiter nếu không rõ
3. Pha từng món một
4. Double check trước khi mark ready
```

### Rush hour
```
1. Ưu tiên orders đơn giản
2. Batch similar drinks
3. Communicate với team
4. Maintain quality, không rush
```

## 📈 Performance Metrics

### Cá nhân
- Số orders hoàn thành/ca
- Thời gian pha trung bình
- Tỷ lệ remake (nếu có)
- Customer feedback

### Team
- Tổng orders/ca
- Queue time trung bình
- Peak hour performance
- Waste percentage

## 🎓 Training

### Espresso Basics
1. Grind size adjustment
2. Tamping pressure (15kg)
3. Extraction time (25-30s)
4. Crema quality check

### Milk Steaming
1. Purge steam wand
2. Aeration phase (0-5s)
3. Texturing phase (5-15s)
4. Final temperature (60-65°C)

### Latte Art
1. Heart
2. Tulip
3. Rosetta
4. Swan (advanced)

## 🔧 Troubleshooting

### App không load orders
```
1. Check internet connection
2. Tap 🔄 để refresh
3. Logout và login lại
4. Báo IT support
```

### Không nhận được order
```
1. Check role (phải là barista)
2. Check ca đã mở chưa
3. Refresh app
4. Báo manager
```

### Order bị stuck
```
1. Check status trong app
2. Báo manager để investigate
3. Có thể cần manual intervention
```

## 📞 Support

### Liên hệ
- 📱 Manager: [số điện thoại]
- 💬 Team chat: [link]
- 🆘 Emergency: [số điện thoại]

### Báo lỗi
1. Chụp màn hình
2. Ghi lại order number
3. Mô tả vấn đề
4. Gửi cho IT support

## ✨ Shortcuts

### Keyboard (nếu dùng tablet)
- `Q`: Chuyển sang Queue tab
- `W`: Chuyển sang Working tab
- `R`: Chuyển sang Ready tab
- `F5`: Refresh

### Gestures
- Swipe left/right: Chuyển tabs
- Pull down: Refresh
- Long press: View order details

## 🎯 Goals

### Daily
- [ ] Hoàn thành 100% orders assigned
- [ ] Maintain quality standards
- [ ] Zero waste
- [ ] Clean workspace

### Weekly
- [ ] Improve average prep time
- [ ] Learn new recipe
- [ ] Zero customer complaints
- [ ] Help train new barista

### Monthly
- [ ] Master latte art
- [ ] Reduce waste by 10%
- [ ] Improve efficiency
- [ ] Contribute to menu development

## 🌟 Excellence Standards

### Quality
- Consistent taste
- Perfect temperature
- Beautiful presentation
- Fresh ingredients

### Speed
- Meet time targets
- Efficient workflow
- Minimal waste
- Quick recovery

### Service
- Professional attitude
- Team collaboration
- Customer focus
- Continuous improvement

---

**Remember**: Quality first, speed second. A perfect drink takes time! ☕✨

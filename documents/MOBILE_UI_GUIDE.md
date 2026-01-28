# 📱 Hướng dẫn sử dụng Mobile UI cho Waiter

## 🎯 Tổng quan

Hệ thống đã được redesign với UI tối ưu cho mobile, giúp waiter thao tác nhanh chóng và thuận tiện hơn.

## 🚀 Routes

### Cho tất cả users (Mobile-first)
- `/orders` - Quản lý orders với UI mobile-optimized
- `/mobile` - Dashboard mobile tối ưu
- `/shifts` - Quản lý ca làm việc
- `/profile` - Thông tin cá nhân

### Cho Manager (Desktop)
- `/dashboard` - Dashboard desktop
- `/menu`, `/ingredients`, `/facilities`, `/expenses` - Quản lý
- `/users` - Quản lý người dùng

### Cho Cashier (Desktop)
- `/cashier` - Thu ngân
- `/cashier/reports` - Báo cáo

## 📱 Tính năng Mobile Dashboard

### 1. Header thông minh
- Hiển thị tên user
- Thời gian real-time
- Ngày hiện tại

### 2. Trạng thái ca làm việc
- ✅ **Ca đang mở**: Hiển thị thời gian đã làm
- ⚠️ **Chưa mở ca**: Nút mở ca nhanh

### 3. Thống kê nhanh
- 📋 **Orders hôm nay**: Tổng số orders
- 💰 **Doanh thu**: Tổng tiền thu được
- 🍹 **Đang pha chế**: Orders đang xử lý
- ⏳ **Chờ thanh toán**: Orders chưa thu tiền

### 4. Thao tác nhanh
- **Tạo Order**: Tạo order mới nhanh chóng
- **Thu tiền**: Xem orders cần thu tiền
- **Quản lý ca**: Mở/đóng ca làm việc
- **Cá nhân**: Xem thông tin cá nhân

### 5. Orders gần đây
- Hiển thị 3 orders mới nhất
- Tap để xem chi tiết
- Link "Xem tất cả" để vào trang orders

## 📋 Tính năng Order Management

### 1. Tạo Order mới (Full Screen)

#### Bước 1: Tap nút FAB (➕)
- Nút tròn màu xanh ở góc dưới bên phải

#### Bước 2: Nhập thông tin
- Tên khách hàng (tùy chọn)
- Chọn category: Tất cả, Cà phê, Trà, Nước ép, Đồ ăn

#### Bước 3: Chọn món
- Grid layout hiển thị tất cả món
- Tap vào món để thêm vào cart
- Badge hiển thị số lượng đã chọn

#### Bước 4: Điều chỉnh cart
- **[+]**: Tăng số lượng
- **[-]**: Giảm số lượng
- **[×]**: Xóa món khỏi cart
- Xem tổng tiền real-time

#### Bước 5: Xác nhận
- Tap "Xác nhận" ở góc trên bên phải
- Order được tạo và hiển thị trong danh sách

### 2. Quản lý Orders

#### Filter theo trạng thái
- **📋 Tất cả**: Xem tất cả orders
- **🆕 Mới**: Orders chưa thanh toán
- **💰 Đã thu**: Orders đã thanh toán
- **🍹 Đang pha**: Orders đang pha chế
- **✅ Hoàn tất**: Orders đã phục vụ

#### Quick Actions trên card
- **💰 Thu tiền**: Thu tiền cho order mới
- **🍹 Gửi bar**: Gửi order đã thanh toán đến quầy bar
- **✅ Phục vụ**: Đánh dấu order đã phục vụ xong

#### Xem chi tiết order
- Tap vào order card
- Bottom sheet hiển thị đầy đủ thông tin
- Các actions có sẵn

### 3. Thu tiền nhanh

#### Cách 1: Quick Payment từ card
- Tap "💰 Thu tiền" trên order card
- Modal thu tiền hiện lên
- Số tiền tự động điền
- Chọn phương thức: Tiền mặt / QR / Chuyển khoản
- Tap "Xác nhận"

#### Cách 2: Từ chi tiết order
- Tap vào order card
- Tap "💰 Thu tiền" trong bottom sheet
- Thực hiện tương tự cách 1

### 4. Workflow hoàn chỉnh

```
1. Mở ca làm việc
   ↓
2. Tạo order mới (FAB ➕)
   ↓
3. Chọn món từ menu
   ↓
4. Xác nhận order
   ↓
5. Thu tiền (💰)
   ↓
6. Gửi quầy bar (🍹)
   ↓
7. Đánh dấu đã phục vụ (✅)
   ↓
8. Đóng ca khi hết giờ
```

## 🎨 UI Components

### Bottom Navigation
- 🏠 **Trang chủ**: Dashboard
- 📋 **Orders**: Quản lý orders
- ⏰ **Ca làm**: Quản lý ca
- 👤 **Cá nhân**: Profile

### Status Badges
- 🆕 **Mới tạo**: Màu xám
- 💰 **Đã thu**: Màu xanh lá
- 🍹 **Đang pha**: Màu xanh dương
- ✅ **Hoàn tất**: Màu tím
- ❌ **Đã hủy**: Màu đỏ

### Touch Targets
- Tất cả buttons ≥ 44px (dễ tap)
- Active states với scale animation
- Smooth transitions

## 💡 Tips & Tricks

### Tạo order nhanh nhất
1. Tap FAB (➕)
2. Bỏ qua tên khách (nếu không cần)
3. Tap nhanh các món thường dùng
4. Tap "Xác nhận"
⏱️ **Mục tiêu: < 20 giây**

### Thu tiền nhanh nhất
1. Tìm order trong list (dùng filter nếu cần)
2. Tap "💰 Thu tiền" ngay trên card
3. Kiểm tra số tiền
4. Chọn phương thức
5. Tap "Xác nhận"
⏱️ **Mục tiêu: < 10 giây**

### Sử dụng filters hiệu quả
- **Sáng**: Filter "🆕 Mới" để thu tiền
- **Trưa**: Filter "🍹 Đang pha" để theo dõi
- **Chiều**: Filter "✅ Hoàn tất" để kiểm tra
- **Cuối ca**: Filter "📋 Tất cả" để tổng kết

### Làm việc offline
- Orders được cache local
- Tạo order khi mất mạng
- Tự động sync khi có mạng trở lại

## 🔧 Troubleshooting

### Không tạo được order
- ✅ Kiểm tra đã mở ca chưa
- ✅ Kiểm tra kết nối mạng
- ✅ Kiểm tra đã chọn món chưa

### Không thu được tiền
- ✅ Kiểm tra order đã tạo chưa
- ✅ Kiểm tra số tiền nhập đúng chưa
- ✅ Kiểm tra phương thức thanh toán

### UI bị lag
- 🔄 Pull to refresh để làm mới
- 🔄 Đóng các apps khác
- 🔄 Restart app nếu cần

## 📊 Performance Tips

### Tối ưu tốc độ
- Sử dụng WiFi thay vì 4G khi có thể
- Đóng các tabs không dùng
- Clear cache định kỳ

### Tiết kiệm pin
- Giảm độ sáng màn hình
- Tắt các tính năng không cần thiết
- Sạc đầy trước ca làm

## 🎯 Best Practices

### Cho Waiter
1. ✅ Luôn mở ca trước khi làm việc
2. ✅ Kiểm tra orders pending định kỳ
3. ✅ Thu tiền ngay sau khi khách order
4. ✅ Gửi bar ngay sau khi thu tiền
5. ✅ Đánh dấu phục vụ khi xong
6. ✅ Đóng ca đúng giờ

### Cho Manager
1. ✅ Theo dõi thống kê hàng ngày
2. ✅ Kiểm tra orders bất thường
3. ✅ Review performance của team
4. ✅ Cập nhật menu khi cần

## 🆘 Support

### Liên hệ
- 📧 Email: support@cafepos.com
- 📱 Hotline: 1900-xxxx
- 💬 Chat: Trong app

### Báo lỗi
1. Chụp màn hình lỗi
2. Ghi lại các bước tái hiện
3. Gửi qua email hoặc chat
4. Đợi phản hồi từ team

## 🔄 Updates

### Version 2.0 (Current)
- ✅ Mobile-first UI
- ✅ Full-screen order creation
- ✅ Quick actions
- ✅ Bottom navigation
- ✅ Real-time stats

### Coming Soon
- 🔜 Voice input
- 🔜 QR scan table
- 🔜 Split bill
- 🔜 Offline mode
- 🔜 Push notifications

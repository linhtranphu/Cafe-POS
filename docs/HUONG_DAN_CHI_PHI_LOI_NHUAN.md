# Hướng Dẫn Sử Dụng - Chi Phí & Lợi Nhuận Menu

## Tổng Quan

Tính năng Chi Phí & Lợi Nhuận giúp bạn:

- 📊 Theo dõi giá vốn của từng món
- 💰 Phân tích lợi nhuận và tỷ lệ lời
- ⚠️ Phát hiện món bán lỗ hoặc lợi nhuận thấp
- 📈 Xem báo cáo theo category và thời gian
- 💼 Tính lợi nhuận sau chi phí vận hành

## Mục Lục

1. [Xem Chi Phí Menu](#xem-chi-phí-menu)
2. [Hiểu Các Cảnh Báo](#hiểu-các-cảnh-báo)
3. [Phân Tích Lợi Nhuận](#phân-tích-lợi-nhuận)
4. [Nhập Chi Phí Vận Hành](#nhập-chi-phí-vận-hành)
5. [Phân Biệt Các Loại Chi Phí](#phân-biệt-các-loại-chi-phí)
6. [Thực Hành Tốt Nhất](#thực-hành-tốt-nhất)

---

## Xem Chi Phí Menu

### Truy Cập

1. Đăng nhập với tài khoản Manager
2. Menu chính → **"Chi phí & Lợi nhuận"** → **"Chi phí món"**
3. Mobile: Tap icon menu → **"Chi phí món"**

### Bảng Hiển Thị

| Cột | Ý Nghĩa |
|-----|---------|
| **Tên món** | Tên menu item |
| **Loại** | Category (Coffee, Tea, Food...) |
| **Giá bán** | Giá bán cho khách (VND) |
| **Giá vốn** | Chi phí nguyên liệu (VND) |
| **Lợi nhuận %** | Tỷ lệ lời = ((Giá bán - Giá vốn) / Giá bán) × 100 |
| **Lợi nhuận** | Tiền lời = Giá bán - Giá vốn (VND) |
| **Trạng thái** | Trạng thái chi phí |

### Lọc và Sắp Xếp

**Lọc theo loại:**
- Click dropdown "Tất cả categories"
- Chọn loại muốn xem

**Sắp xếp:**
- Click "Sắp xếp theo"
- Chọn: Lợi nhuận %, Lợi nhuận tiền, hoặc Tên
- Click mũi tên để đổi thứ tự

### Xem Chi Tiết Nguyên Liệu

1. Click vào món muốn xem
2. Hiển thị bảng nguyên liệu:
   - Tên nguyên liệu
   - Số lượng và đơn vị
   - Giá/đơn vị
   - Thành tiền
3. Tổng giá vốn ở cuối

⚠️ **Lưu ý:** Nếu thiếu giá nguyên liệu sẽ có cảnh báo màu đỏ


### Thống Kê Tổng Quan

Phần đầu trang hiển thị:

- **Tổng số món**: Tổng menu items
- **Món bán lỗ**: Số món giá vốn > giá bán (đỏ)
- **Lợi nhuận thấp**: Số món lợi nhuận < ngưỡng (vàng)
- **Lợi nhuận TB**: Trung bình tất cả món

---

## Hiểu Các Cảnh Báo

### 🔴 Màu Đỏ - Bán Lỗ

**Nghĩa:** Giá vốn > Giá bán → Bán càng nhiều càng lỗ

**Ví dụ:**
```
Giá bán:    20,000 đ
Giá vốn:    25,000 đ
Lỗ:         -5,000 đ (-25%)
```

**Cần làm gì:**
1. ✅ Kiểm tra recipe - dùng quá nhiều nguyên liệu?
2. ✅ Tăng giá bán
3. ✅ Giảm cost:
   - Thay đổi recipe
   - Tìm supplier rẻ hơn
   - Giảm hao hụt

### 🟡 Màu Vàng - Lợi Nhuận Thấp

**Nghĩa:** Lợi nhuận < ngưỡng cảnh báo (mặc định 20%)

**Ví dụ:**
```
Giá bán:    30,000 đ
Giá vốn:    25,000 đ
Lời:         5,000 đ (16.67%)
```

**Nên làm:**
- Xem xét tăng giá
- Tối ưu recipe
- Món này có thể không đủ lời để cover chi phí vận hành

### 🟢 Màu Xanh - Lợi Nhuận Tốt

**Nghĩa:** Lợi nhuận >= ngưỡng cảnh báo

**Ví dụ:**
```
Giá bán:    45,000 đ
Giá vốn:    15,000 đ
Lời:        30,000 đ (66.67%)
```

**Trạng thái:** Tốt! Tiếp tục duy trì

### ⚪ Màu Xám - Thiếu Dữ Liệu

**Nghĩa:** Không tính được vì thiếu giá nguyên liệu

**Cần làm:**
1. Click vào món xem chi tiết
2. Xác định nguyên liệu nào thiếu giá
3. Vào Quản lý nguyên liệu → Cập nhật giá
4. Hệ thống tự động tính lại

---

## Phân Tích Lợi Nhuận

### Lợi Nhuận Theo Category

Xem lợi nhuận từng loại món (Coffee, Tea, Food...)

**Các bước:**
1. Menu → **"Chi phí & Lợi nhuận"** → **"Phân tích lợi nhuận"**
2. Tab **"Theo category"**
3. Chọn thời gian:
   - Hôm nay
   - Tuần này
   - Tháng này
   - Tùy chỉnh

**Bảng hiển thị:**

| Cột | Ý Nghĩa |
|-----|---------|
| **Loại** | Category |
| **Doanh thu** | Tổng tiền bán |
| **Giá vốn** | Tổng chi phí nguyên liệu |
| **Lợi nhuận** | Doanh thu - Giá vốn |
| **Lợi nhuận %** | Tỷ lệ lời trung bình |
| **Số đơn** | Số orders |
| **Số món** | Số items bán ra |

**Dùng để:**
- Biết category nào lời nhất
- So sánh giữa các categories
- Quyết định focus vào category nào

### Lợi Nhuận Vận Hành

Xem lợi nhuận thực tế sau khi trừ chi phí (lương, mặt bằng, điện nước...)

**Các bước:**
1. Menu → **"Chi phí & Lợi nhuận"** → **"Phân tích lợi nhuận"**
2. Tab **"Operating Profit"**
3. Chọn thời gian

**Báo cáo gồm 3 phần:**

#### 1. Lợi Nhuận Gộp
```
Doanh thu:           10,000,000 đ
Giá vốn hàng bán:    -3,000,000 đ
─────────────────────────────────
Lợi nhuận gộp:        7,000,000 đ (70%)
```

#### 2. Chi Phí Vận Hành
```
Lương nhân viên:     -2,000,000 đ
Tiền thuê mặt bằng:  -1,500,000 đ
Điện nước:             -500,000 đ
Marketing:             -300,000 đ
Chi phí khác:          -200,000 đ
─────────────────────────────────
Tổng chi phí:        -4,500,000 đ
```

#### 3. Lợi Nhuận Vận Hành
```
Lợi nhuận vận hành:   2,500,000 đ (25%)
```

**Hiểu con số:**
- **Lợi nhuận gộp** = Lời từ bán hàng (chưa trừ chi phí vận hành)
- **Lợi nhuận vận hành** = Lời thực tế (đã trừ tất cả chi phí)

⚠️ **Lưu ý quan trọng:**

1. **Phân bổ chi phí:**
   - Nhập chi phí theo tháng, xem báo cáo theo ngày
   - Hệ thống tự động chia: `chi phí ngày = chi phí tháng / số ngày`
   - Có note: "Chi phí được phân bổ từ tháng"

2. **Chưa có chi phí:**
   - Nếu chưa nhập chi phí cho kỳ này
   - Chỉ hiển thị Lợi nhuận gộp
   - Note: "Chưa nhập chi phí vận hành"

---

## Nhập Chi Phí Vận Hành

### Truy Cập Form

**Cách 1: Từ Settings**
1. Vào **Settings** (icon bánh răng)
2. Kéo xuống **"Chi phí vận hành"**
3. Click **"Thêm chi phí mới"**

**Cách 2: Từ Phân Tích**
1. Vào **"Phân tích lợi nhuận"** → **"Operating Profit"**
2. Nếu chưa có chi phí, click **"Nhập chi phí"**

### Điền Form

**1. Kỳ chi phí**
- **Từ ngày**: Ngày bắt đầu (VD: 01/01/2024)
- **Đến ngày**: Ngày kết thúc (VD: 31/01/2024)
- End date phải >= Start date

💡 **Mẹo:** Nhập theo tháng để dễ quản lý

**2. Các loại chi phí**

| Mục | Bao gồm | Ví dụ |
|-----|---------|-------|
| **Lương nhân viên** | Lương + BHXH + phụ cấp | 15,000,000 đ |
| **Tiền thuê** | Rent + phí quản lý | 10,000,000 đ |
| **Điện nước** | Điện + nước + internet | 2,000,000 đ |
| **Marketing** | Quảng cáo + khuyến mãi | 1,500,000 đ |
| **Chi phí khác** | Bảo trì, vệ sinh... | 1,000,000 đ |

**3. Tổng chi phí**
- Tự động tính khi bạn nhập
- Hiển thị real-time

**4. Lưu**
- Click **"Lưu"** để save
- Hoặc **"Hủy"** để bỏ qua

### Sửa Chi Phí Đã Có

1. Settings → **"Chi phí vận hành"**
2. Danh sách chi phí theo kỳ
3. Click vào chi phí muốn sửa
4. Form mở với data cũ
5. Sửa và **"Lưu"**

**Lưu ý:** Nếu kỳ đã tồn tại, lưu sẽ cập nhật (không tạo mới)

### Thực Hành Tốt

**1. Nhập theo tháng**
- Dễ quản lý
- Phù hợp chu kỳ thanh toán

**2. Nhập đầy đủ**
- Đừng bỏ sót mục nào
- Không có thì nhập 0

**3. Nhập đúng kỳ**
- Lương tháng 1 → kỳ 01/01 - 31/01
- Đừng nhập sai tháng

**4. Review định kỳ**
- Cuối tháng review lại
- So sánh với tháng trước

---

## Phân Biệt Các Loại Chi Phí

### Giá Vốn Hiện Tại (Current Cost)

**Là gì:** Giá vốn tính từ giá nguyên liệu **hiện tại**

**Khi nào update:**
- Tự động khi giá nguyên liệu đổi
- Background job chạy tự động

**Dùng để:**
- Quyết định giá bán
- Phân tích real-time
- Xem trong Menu Cost View

**Ví dụ:**
```
Hôm nay: Sữa = 50đ/ml → Cappuccino = 15,000đ
Ngày mai: Sữa = 60đ/ml → Cappuccino = 17,000đ
```

### Giá Vốn Kế Toán (Accounting Cost)

**Là gì:** Giá vốn **chính thức** lưu khi kết ca

**Khi nào tính:**
- Chỉ tính 1 lần khi kết ca
- Dùng giá nguyên liệu lúc kết ca
- **Không đổi** sau khi lưu

**Dùng để:**
- Báo cáo kế toán
- Phân tích lợi nhuận theo category
- Phân tích operating profit

**Ví dụ:**
```
Ca 15/01: Kết ca 10pm → Lưu cost = 15,000đ
16/01: Sữa tăng giá → Cost vẫn là 15,000đ (không đổi)
17/01: Xem báo cáo 15/01 → Dùng cost = 15,000đ
```

**Tại sao không đổi?**
- Đúng quy trình kế toán
- Cost phản ánh giá vốn lúc chốt sổ
- Không bị ảnh hưởng giá sau đó

### Trạng Thái Chi Phí

| Trạng thái | Nghĩa | Màu |
|------------|-------|-----|
| **FINAL** | Đã tính và lưu chính thức | Xanh |
| **ESTIMATED** | Tạm tính (ca chưa đóng) | Vàng |
| **INCOMPLETE** | Thiếu giá nguyên liệu | Đỏ |


---

## Thực Hành Tốt Nhất

### 1. Giữ Giá Nguyên Liệu Chính Xác

**Tại sao quan trọng:**
- Tính chi phí dựa hoàn toàn vào giá nguyên liệu
- Sai giá → sai chi phí → quyết định sai

**Cách làm:**
- Cập nhật giá mỗi khi nhập hàng
- Hệ thống tự tính weighted average
- Review định kỳ

**Kiểm tra thiếu giá:**
1. Vào Menu Cost View
2. Lọc INCOMPLETE
3. Click xem nguyên liệu nào thiếu
4. Cập nhật giá

### 2. Kết Ca Đều Đặn

**Tại sao quan trọng:**
- Accounting cost chỉ tính khi kết ca
- Không kết ca → không có data báo cáo

**Thực hành:**
- Kết ca mỗi ngày (hoặc mỗi ca)
- Đừng để ca mở quá lâu
- Kết ca trước khi đổi giá nguyên liệu

### 3. Nhập Chi Phí Hàng Tháng

**Tại sao quan trọng:**
- Operating profit = Lợi nhuận gộp - Chi phí
- Không nhập → không biết lời thực tế

**Thực hành:**
- Nhập cuối mỗi tháng
- Nhập đầy đủ tất cả mục
- So sánh với tháng trước

💡 **Mẹo:** Đặt nhắc nhở cuối tháng

### 4. Review Cảnh Báo Hàng Tuần

**Tại sao quan trọng:**
- Phát hiện sớm món có vấn đề
- Điều chỉnh kịp thời

**Checklist hàng tuần:**
- [ ] Check món đỏ (lỗ) → Xử lý ngay
- [ ] Check món vàng (lời thấp) → Xem xét điều chỉnh
- [ ] Check món xám (thiếu data) → Cập nhật giá
- [ ] So sánh tuần trước → Tìm xu hướng

### 5. Phân Tích Xu Hướng

**Tại sao quan trọng:**
- Hiểu xu hướng chi phí và lợi nhuận
- Dự đoán và plan tương lai

**Cách làm:**
- Xem báo cáo theo tháng
- So sánh tháng này vs tháng trước
- Tìm patterns:
  - Category nào tăng/giảm lời?
  - Chi phí có xu hướng tăng?
  - Chi phí vận hành bất thường?

### 6. Đặt Ngưỡng Cảnh Báo Phù Hợp

**Khuyến nghị:**
- **Cafe cao cấp**: 30-40% (chi phí cao)
- **Cafe bình dân**: 15-25% (cạnh tranh giá)
- **Mặc định**: 20% (phù hợp đa số)

**Cách điều chỉnh:**
1. Settings
2. Tìm "Ngưỡng cảnh báo lợi nhuận thấp"
3. Nhập giá trị mới
4. Lưu

### 7. Xử Lý Quy Đổi Đơn Vị

**Tại sao quan trọng:**
- Kho theo kg, recipe theo gram → Cần quy đổi
- Không đúng → Chi phí sai

**Cách làm:**
1. Quản lý nguyên liệu
2. Edit nguyên liệu
3. Đặt **conversion_rate**:
   - Kho: kg, Recipe: gram → rate = 1000
   - Kho: lít, Recipe: ml → rate = 1000
4. Đặt **wastage_percentage** nếu có hao hụt

**Công thức:**
```
chi phí = (số lượng × quy đổi × giá) × (1 + hao hụt%/100)
```

---

## Xử Lý Sự Cố

### Vấn đề: Chi phí không update sau khi đổi giá

**Nguyên nhân:**
1. Background job đang chạy
2. Lỗi trong quá trình tính

**Giải quyết:**
1. Đợi 5-10 giây và click "Refresh"
2. Check trạng thái recalculation
3. Vẫn không được → Liên hệ support

### Vấn đề: Món hiển thị "Thiếu giá nguyên liệu"

**Nguyên nhân:** Nguyên liệu không có giá

**Giải quyết:**
1. Click vào món xem chi tiết
2. Xác định nguyên liệu nào thiếu (màu đỏ)
3. Quản lý nguyên liệu → Cập nhật giá
4. Hệ thống tự tính lại

### Vấn đề: Operating Profit không có data

**Nguyên nhân:**
1. Chưa nhập chi phí
2. Không có orders
3. Ca chưa đóng

**Giải quyết:**
1. Check đã nhập chi phí chưa
2. Thử date range khác
3. Đảm bảo ca đã kết

### Vấn đề: Lợi nhuận âm nhưng món bán chạy

**Phân tích:**
- Món đang bán lỗ
- Bán càng nhiều càng lỗ

**Xử lý:**
1. **Khẩn cấp:** Tăng giá hoặc ngưng bán
2. Review recipe
3. Tìm supplier rẻ hơn
4. Đổi recipe

---

## Tham Khảo Nhanh

### Chỉ Số Quan Trọng

| Chỉ số | Công thức | Tốt |
|--------|-----------|-----|
| **Lợi nhuận %** | ((Giá - Vốn) / Giá) × 100 | > 20% |
| **Lợi nhuận gộp** | Doanh thu - Giá vốn | > 60% doanh thu |
| **Lợi nhuận vận hành** | Lợi nhuận gộp - Chi phí | > 20% doanh thu |

### Mã Màu

| Màu | Nghĩa | Hành động |
|-----|-------|-----------|
| 🔴 Đỏ | Lỗ (vốn > giá) | Xử lý ngay |
| 🟡 Vàng | Lời thấp | Xem xét điều chỉnh |
| 🟢 Xanh | Lời tốt | Duy trì |
| ⚪ Xám | Thiếu data | Cập nhật giá |

### Đường Dẫn Nhanh

| Tính năng | Đường dẫn |
|-----------|-----------|
| Chi phí món | Menu → Chi phí & Lợi nhuận → Chi phí món |
| Phân tích | Menu → Chi phí & Lợi nhuận → Phân tích lợi nhuận |
| Chi phí vận hành | Settings → Chi phí vận hành |
| Quản lý nguyên liệu | Menu → Quản lý nguyên liệu |

---

## Hỗ Trợ

Nếu gặp vấn đề:

1. **Đọc hướng dẫn này** - Hầu hết câu hỏi đã có đáp án
2. **Xem phần Xử lý sự cố** - Vấn đề thường gặp
3. **Liên hệ support** - Nếu vẫn chưa giải quyết được

**Nhớ:** Tính năng này giúp bạn quyết định kinh doanh chính xác. Hãy dành thời gian hiểu và sử dụng thường xuyên!

---

*Cập nhật: Tháng 2/2026*
*Phiên bản: 1.0*

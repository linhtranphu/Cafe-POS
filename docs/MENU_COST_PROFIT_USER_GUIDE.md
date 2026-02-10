# Menu Cost & Profit Analysis - User Guide

## Overview

Tính năng Menu Cost & Profit Analysis giúp bạn theo dõi chi phí nguyên liệu và phân tích lợi nhuận cho từng món trong menu. Với tính năng này, bạn có thể:

- Xem giá vốn (cost) của từng món dựa trên giá nguyên liệu hiện tại
- Phân tích profit margin (tỷ lệ lợi nhuận) và absolute profit (lợi nhuận tiền mặt)
- Phát hiện các món bán lỗ hoặc có lợi nhuận thấp
- Xem báo cáo lợi nhuận theo category
- Phân tích operating profit sau khi trừ chi phí vận hành

## Table of Contents

1. [Viewing Menu Costs](#viewing-menu-costs)
2. [Understanding Warning Indicators](#understanding-warning-indicators)
3. [Analyzing Profits](#analyzing-profits)
4. [Inputting Operating Expenses](#inputting-operating-expenses)
5. [Understanding Cost Types](#understanding-cost-types)
6. [Best Practices](#best-practices)

---

## Viewing Menu Costs

### Accessing the Menu Cost View

1. Đăng nhập với tài khoản Manager
2. Từ menu chính, chọn **"Chi phí & Lợi nhuận"** → **"Chi phí món"**
3. Hoặc trên mobile, tap vào icon menu và chọn **"Chi phí món"**

### Understanding the Menu Cost Table

Bảng hiển thị các thông tin sau cho mỗi món:

| Column | Description |
|--------|-------------|
| **Tên món** | Tên menu item |
| **Category** | Loại món (Coffee, Tea, Food, etc.) |
| **Giá bán** | Giá bán cho khách (VND) |
| **Giá vốn** | Chi phí nguyên liệu hiện tại (VND) |
| **Lợi nhuận %** | Profit margin = ((Giá bán - Giá vốn) / Giá bán) × 100 |
| **Lợi nhuận** | Absolute profit = Giá bán - Giá vốn (VND) |
| **Trạng thái** | Cost status (xem phần Warning Indicators) |

### Filtering and Sorting

**Filter by Category:**
- Click vào dropdown "Tất cả categories"
- Chọn category muốn xem (Coffee, Tea, Food, etc.)

**Sort Options:**
- Click vào dropdown "Sắp xếp theo"
- Chọn tiêu chí: Lợi nhuận %, Lợi nhuận tiền mặt, hoặc Tên món
- Click icon mũi tên để đổi thứ tự tăng/giảm dần

### Viewing Cost Breakdown

Để xem chi tiết nguyên liệu của một món:

1. Click vào row của món muốn xem
2. Modal sẽ hiển thị bảng breakdown với các cột:
   - **Nguyên liệu**: Tên ingredient
   - **Số lượng**: Quantity sử dụng
   - **Đơn vị**: Unit (gram, ml, etc.)
   - **Giá/đơn vị**: Cost per unit hiện tại
   - **Thành tiền**: Total cost cho ingredient này

3. Nếu ingredient có conversion rate hoặc wastage percentage, sẽ hiển thị thêm:
   - **Quy đổi**: Conversion rate (e.g., kg → gram)
   - **Hao hụt**: Wastage percentage (%)

4. Tổng giá vốn hiển thị ở cuối bảng

**⚠️ Warning:** Nếu có ingredient thiếu giá (cost_per_unit), sẽ hiển thị cảnh báo màu đỏ.

### Summary Statistics

Phần đầu trang hiển thị tổng quan:

- **Tổng số món**: Total menu items
- **Món bán lỗ**: Số món có giá vốn > giá bán (màu đỏ)
- **Lợi nhuận thấp**: Số món có profit margin < ngưỡng cảnh báo (màu vàng)
- **Lợi nhuận TB**: Average profit margin của tất cả món

### Recalculation Status

Khi giá nguyên liệu thay đổi, hệ thống tự động tính lại giá vốn. Trong lúc đang tính:
- Hiển thị indicator "Đang cập nhật chi phí..."
- Bạn vẫn có thể xem và thao tác bình thường
- Sau khi hoàn tất, click "Refresh" để xem dữ liệu mới

---

## Understanding Warning Indicators

Hệ thống sử dụng màu sắc để cảnh báo các món có vấn đề về lợi nhuận:

### 🔴 Red - Loss (Bán lỗ)

**Meaning:** Giá vốn > Giá bán → Bán càng nhiều càng lỗ

**Example:**
- Giá bán: 20,000 VND
- Giá vốn: 25,000 VND
- Lợi nhuận: -5,000 VND (-25%)

**Action Required:**
1. Kiểm tra recipe - có thể dùng quá nhiều nguyên liệu?
2. Tăng giá bán
3. Hoặc giảm cost bằng cách:
   - Thay đổi recipe (dùng ít nguyên liệu hơn)
   - Tìm supplier rẻ hơn
   - Giảm wastage

### 🟡 Yellow - Low Margin (Lợi nhuận thấp)

**Meaning:** Profit margin < ngưỡng cảnh báo (default: 20%)

**Example:**
- Giá bán: 30,000 VND
- Giá vốn: 25,000 VND
- Lợi nhuận: 5,000 VND (16.67%)

**Action Recommended:**
- Xem xét tăng giá bán
- Hoặc tối ưu recipe để giảm cost
- Món này có thể không đủ lợi nhuận để cover chi phí vận hành

**Note:** Bạn có thể điều chỉnh ngưỡng cảnh báo trong Settings (mặc định 20%)

### 🟢 Green - Profitable (Lợi nhuận tốt)

**Meaning:** Profit margin >= ngưỡng cảnh báo

**Example:**
- Giá bán: 45,000 VND
- Giá vốn: 15,000 VND
- Lợi nhuận: 30,000 VND (66.67%)

**Status:** Món này có lợi nhuận tốt, tiếp tục duy trì!

### ⚪ Gray - Incomplete Data (Thiếu dữ liệu)

**Meaning:** Không thể tính giá vốn vì thiếu giá nguyên liệu

**Icon:** ⚠ Thiếu giá nguyên liệu

**Action Required:**
1. Click vào món để xem cost breakdown
2. Xác định ingredient nào thiếu giá
3. Vào **Quản lý nguyên liệu** → Cập nhật cost_per_unit cho ingredient đó
4. Hệ thống sẽ tự động tính lại giá vốn

---

## Analyzing Profits

### Category Profit Analysis

Xem lợi nhuận theo từng category (Coffee, Tea, Food, etc.)

**Steps:**
1. Từ menu chính, chọn **"Chi phí & Lợi nhuận"** → **"Phân tích lợi nhuận"**
2. Chọn tab **"Theo category"**
3. Chọn date range:
   - **Hôm nay**: Chỉ hôm nay
   - **Tuần này**: 7 ngày gần nhất
   - **Tháng này**: Tháng hiện tại
   - **Tùy chỉnh**: Chọn start date và end date

**Table Columns:**

| Column | Description |
|--------|-------------|
| **Category** | Loại món |
| **Doanh thu** | Total revenue từ orders |
| **Giá vốn** | Total cost of goods sold (COGS) |
| **Lợi nhuận** | Total profit = Doanh thu - Giá vốn |
| **Lợi nhuận %** | Average profit margin |
| **Số đơn** | Number of orders |
| **Số món** | Number of items sold |

**Use Cases:**
- Xác định category nào có lợi nhuận cao nhất
- So sánh performance giữa các categories
- Quyết định focus vào category nào để tăng doanh thu


### Operating Profit Analysis

Xem lợi nhuận thực tế sau khi trừ chi phí vận hành (lương, mặt bằng, điện nước, marketing)

**Steps:**
1. Từ menu chính, chọn **"Chi phí & Lợi nhuận"** → **"Phân tích lợi nhuận"**
2. Chọn tab **"Operating Profit"**
3. Chọn date range (tương tự Category Profit)

**Report Sections:**

#### 1. Gross Profit (Lợi nhuận gộp)
- **Doanh thu**: Total revenue từ orders
- **Giá vốn hàng bán**: Total COGS
- **Lợi nhuận gộp**: Revenue - COGS
- **Lợi nhuận gộp %**: (Gross Profit / Revenue) × 100

#### 2. Operating Expenses (Chi phí vận hành)
- **Lương nhân viên**: Staff salary
- **Tiền thuê mặt bằng**: Rent
- **Điện nước**: Utilities
- **Marketing**: Marketing costs
- **Chi phí khác**: Other expenses
- **Tổng chi phí**: Sum of all expenses

#### 3. Operating Profit (Lợi nhuận vận hành)
- **Lợi nhuận vận hành**: Gross Profit - Total Expenses
- **Lợi nhuận vận hành %**: (Operating Profit / Revenue) × 100

**Understanding the Numbers:**

- **Gross Profit** = Lợi nhuận từ bán hàng trước khi trừ chi phí vận hành
- **Operating Profit** = Lợi nhuận thực tế sau khi trừ tất cả chi phí

**Example:**
```
Doanh thu:           10,000,000 VND
Giá vốn:             -3,000,000 VND
─────────────────────────────────
Lợi nhuận gộp:        7,000,000 VND (70%)

Chi phí vận hành:
  - Lương:           -2,000,000 VND
  - Mặt bằng:        -1,500,000 VND
  - Điện nước:         -500,000 VND
  - Marketing:         -300,000 VND
  - Khác:              -200,000 VND
─────────────────────────────────
Tổng chi phí:        -4,500,000 VND

Lợi nhuận vận hành:   2,500,000 VND (25%)
```

**⚠️ Important Notes:**

1. **Expense Allocation Indicator:**
   - Nếu bạn nhập expenses theo tháng nhưng xem report theo ngày
   - Hệ thống sẽ phân bổ tự động: `daily_expense = monthly_expense / days_in_month`
   - Hiển thị note: "Chi phí được phân bổ từ tháng"

2. **Missing Expenses:**
   - Nếu chưa nhập expenses cho period này
   - Chỉ hiển thị Gross Profit với note: "Chưa nhập chi phí vận hành"
   - Click "Nhập chi phí" để thêm expenses

---

## Inputting Operating Expenses

Operating expenses là chi phí vận hành cần nhập thủ công (lương, mặt bằng, điện nước, marketing)

### Accessing the Form

**Option 1: From Settings**
1. Vào **Settings** (icon gear)
2. Scroll xuống section **"Chi phí vận hành"**
3. Click **"Thêm chi phí mới"**

**Option 2: From Profit Analysis**
1. Vào **"Phân tích lợi nhuận"** → **"Operating Profit"**
2. Nếu chưa có expenses, click **"Nhập chi phí"**

### Filling the Form

**1. Period (Kỳ chi phí)**
- **Từ ngày**: Start date (e.g., 2024-01-01)
- **Đến ngày**: End date (e.g., 2024-01-31)
- **Validation**: End date phải >= Start date

**Tip:** Nhập theo tháng để dễ quản lý (e.g., 01/01 → 31/01)

**2. Expense Categories**

| Field | Description | Example |
|-------|-------------|---------|
| **Lương nhân viên** | Tổng lương + BHXH + phụ cấp | 15,000,000 VND |
| **Tiền thuê mặt bằng** | Rent + phí quản lý | 10,000,000 VND |
| **Điện nước** | Electricity + water + internet | 2,000,000 VND |
| **Marketing** | Quảng cáo + promotion | 1,500,000 VND |
| **Chi phí khác** | Bảo trì, vệ sinh, etc. | 1,000,000 VND |

**3. Total Expenses**
- Tự động tính = Sum of all categories
- Hiển thị real-time khi bạn nhập

**4. Save**
- Click **"Lưu"** để save
- Hoặc **"Hủy"** để cancel

### Editing Existing Expenses

1. Vào **Settings** → **"Chi phí vận hành"**
2. Danh sách expenses hiển thị theo period
3. Click vào expense muốn edit
4. Form sẽ mở với data đã có
5. Sửa và click **"Lưu"**

**Note:** Nếu period đã tồn tại, save sẽ update (không tạo mới)

### Best Practices for Expenses

**1. Nhập theo tháng:**
- Dễ quản lý và theo dõi
- Phù hợp với chu kỳ thanh toán (lương, rent thường theo tháng)

**2. Nhập đầy đủ:**
- Đừng bỏ sót category nào
- Nếu không có, nhập 0 (để rõ ràng)

**3. Nhập đúng kỳ:**
- Lương tháng 1 → nhập period 01/01 - 31/01
- Đừng nhập trước hoặc sau

**4. Review định kỳ:**
- Cuối mỗi tháng, review lại expenses
- So sánh với tháng trước để phát hiện bất thường

---

## Understanding Cost Types

Hệ thống có 2 loại cost khác nhau, quan trọng phải hiểu rõ:

### Current Cost (Giá vốn hiện tại)

**Definition:** Giá vốn tính từ cost_per_unit **hiện tại** của ingredients

**When Updated:**
- Tự động update khi giá nguyên liệu thay đổi
- Background job chạy asynchronously (eventual consistency)

**Used For:**
- Pricing decisions (quyết định giá bán)
- Real-time cost analysis
- Menu cost view

**Example:**
- Hôm nay: Milk cost = 50 VND/ml → Cappuccino current_cost = 15,000 VND
- Ngày mai: Milk cost tăng lên 60 VND/ml → Cappuccino current_cost = 17,000 VND

### Accounting Cost (Giá vốn kế toán)

**Definition:** Giá vốn **chính thức** được lưu khi kết ca (shift closure)

**When Calculated:**
- Chỉ tính 1 lần khi manager kết ca
- Sử dụng cost_per_unit tại thời điểm kết ca
- **Immutable** - không thay đổi sau khi lưu

**Used For:**
- Profit/loss reports (báo cáo kế toán)
- Category profit analysis
- Operating profit analysis

**Example:**
- Ca ngày 15/01: Kết ca lúc 10pm → Lưu accounting_cost = 15,000 VND
- Ngày 16/01: Milk cost tăng → accounting_cost vẫn là 15,000 VND (không đổi)
- Ngày 17/01: Xem report ngày 15/01 → Dùng accounting_cost = 15,000 VND

**Why Immutable?**
- Phù hợp với quy trình kế toán
- Cost phản ánh giá vốn tại thời điểm chốt sổ
- Không bị ảnh hưởng bởi thay đổi giá sau đó

### Cost Status

Mỗi cost có status để biết độ chính xác:

| Status | Meaning | Color |
|--------|---------|-------|
| **FINAL** | Cost đã được tính và lưu chính thức | Green |
| **ESTIMATED** | Cost tạm tính (shift chưa đóng) | Yellow |
| **INCOMPLETE** | Thiếu giá nguyên liệu | Red |

**In Reports:**
- Orders trong ca đã đóng → cost_status = FINAL
- Orders trong ca chưa đóng → cost_status = ESTIMATED (hiển thị indicator)
- Menu items thiếu giá → cost_status = INCOMPLETE (không tính vào reports)


---

## Best Practices

### 1. Maintain Accurate Ingredient Costs

**Why Important:**
- Cost calculations phụ thuộc hoàn toàn vào cost_per_unit của ingredients
- Sai giá nguyên liệu → sai giá vốn → quyết định sai

**How To:**
- Cập nhật cost_per_unit mỗi khi nhập hàng
- Hệ thống tự động tính weighted average cost
- Review định kỳ để đảm bảo chính xác

**Check for Missing Costs:**
1. Vào **Menu Cost View**
2. Filter by cost_status = INCOMPLETE
3. Click vào từng món để xem ingredient nào thiếu giá
4. Cập nhật giá cho ingredients đó

### 2. Close Shifts Regularly

**Why Important:**
- Accounting cost chỉ được tính khi kết ca
- Không kết ca → không có data cho profit reports

**Best Practice:**
- Kết ca mỗi ngày (hoặc mỗi shift)
- Đừng để ca mở quá lâu
- Kết ca trước khi cập nhật giá nguyên liệu (để accounting cost chính xác)

### 3. Input Operating Expenses Monthly

**Why Important:**
- Operating profit = Gross profit - Expenses
- Không nhập expenses → không biết lợi nhuận thực tế

**Best Practice:**
- Nhập expenses cuối mỗi tháng
- Nhập đầy đủ tất cả categories
- Review và compare với tháng trước

**Tip:** Set reminder cuối tháng để nhập expenses

### 4. Review Warnings Weekly

**Why Important:**
- Phát hiện sớm các món có vấn đề
- Điều chỉnh kịp thời để tránh lỗ

**Weekly Review Checklist:**
- [ ] Check **Loss Items** (red) → Action required ngay
- [ ] Check **Low Margin Items** (yellow) → Consider adjusting
- [ ] Review **Incomplete Data** (gray) → Update missing costs
- [ ] Compare với tuần trước → Identify trends

### 5. Analyze Trends Over Time

**Why Important:**
- Hiểu được xu hướng cost và profit
- Dự đoán và plan cho tương lai

**How To:**
- Xem Category Profit report theo tháng
- So sánh tháng này vs tháng trước
- Identify patterns:
  - Category nào đang tăng/giảm profit?
  - Cost có xu hướng tăng không?
  - Operating expenses có tăng bất thường không?

### 6. Set Appropriate Margin Threshold

**Why Important:**
- Threshold quá cao → Quá nhiều warnings (noise)
- Threshold quá thấp → Bỏ sót món có vấn đề

**Recommended Thresholds:**
- **Cafe cao cấp**: 30-40% (chi phí vận hành cao)
- **Cafe bình dân**: 15-25% (cạnh tranh về giá)
- **Default**: 20% (phù hợp đa số)

**How To Adjust:**
1. Vào **Settings**
2. Tìm **"Ngưỡng cảnh báo lợi nhuận thấp"**
3. Nhập giá trị mới (e.g., 25)
4. Click **"Lưu"**

### 7. Handle Unit Conversions Properly

**Why Important:**
- Stock theo kg nhưng recipe theo gram → Cần conversion
- Không convert đúng → Cost sai

**How To:**
1. Vào **Quản lý nguyên liệu**
2. Edit ingredient
3. Set **conversion_rate**:
   - Stock: kg, Recipe: gram → conversion_rate = 1000
   - Stock: liter, Recipe: ml → conversion_rate = 1000
4. Set **wastage_percentage** nếu có hao hụt (e.g., 10%)

**Formula:**
```
cost = (quantity × conversion_rate × cost_per_unit) × (1 + wastage_percentage/100)
```

### 8. Understand Report Limitations

**Important Notes:**

1. **Estimated Costs:**
   - Orders trong ca chưa đóng có cost_status = ESTIMATED
   - Chỉ là ước tính, chưa chính thức
   - Kết ca để có accounting cost chính xác

2. **Historical Data:**
   - Orders trước khi có feature này không có accounting cost
   - Hệ thống dùng current cost làm fallback (với indicator "Estimated")

3. **Expense Allocation:**
   - Xem daily report nhưng expenses nhập monthly → Phân bổ tự động
   - Chỉ là ước tính, không phản ánh chi phí thực tế từng ngày

4. **Incomplete Data:**
   - Menu items thiếu cost không tính vào reports
   - Đảm bảo tất cả ingredients có giá để reports chính xác

---

## Troubleshooting

### Problem: Cost không update sau khi đổi giá nguyên liệu

**Possible Causes:**
1. Background job đang chạy (eventual consistency)
2. Lỗi trong recalculation process

**Solutions:**
1. Đợi 5-10 giây và click **"Refresh"**
2. Check recalculation status indicator
3. Nếu vẫn không update, liên hệ support

### Problem: Món hiển thị "Thiếu giá nguyên liệu"

**Cause:** Một hoặc nhiều ingredients không có cost_per_unit

**Solution:**
1. Click vào món để xem cost breakdown
2. Xác định ingredient nào thiếu giá (highlight màu đỏ)
3. Vào **Quản lý nguyên liệu** → Cập nhật cost_per_unit
4. Hệ thống sẽ tự động tính lại

### Problem: Operating Profit report không có data

**Possible Causes:**
1. Chưa nhập operating expenses cho period này
2. Không có orders trong date range
3. Shifts chưa được đóng

**Solutions:**
1. Check xem đã nhập expenses chưa → Nếu chưa, click "Nhập chi phí"
2. Thử date range khác
3. Đảm bảo shifts đã được kết ca

### Problem: Profit margin âm nhưng món vẫn bán chạy

**Analysis:**
- Món này đang bán lỗ
- Bán càng nhiều càng lỗ tiền

**Actions:**
1. **Urgent:** Tăng giá bán hoặc tạm ngưng bán
2. Review recipe - có thể dùng quá nhiều nguyên liệu?
3. Tìm supplier rẻ hơn
4. Xem xét thay đổi recipe

### Problem: Category profit khác với tổng profit của các món

**Explanation:**
- Category profit dùng **accounting_cost** (từ shift closure)
- Menu cost view dùng **current_cost** (real-time)
- Nếu giá nguyên liệu thay đổi giữa 2 thời điểm → Số liệu khác nhau

**This is Normal:** Đây là design đúng để đảm bảo accounting accuracy

---

## Quick Reference

### Key Metrics

| Metric | Formula | Good Range |
|--------|---------|------------|
| **Profit Margin** | ((Price - Cost) / Price) × 100 | > 20% |
| **Gross Profit** | Revenue - COGS | > 60% of revenue |
| **Operating Profit** | Gross Profit - Expenses | > 20% of revenue |

### Color Codes

| Color | Meaning | Action |
|-------|---------|--------|
| 🔴 Red | Loss (cost > price) | Urgent action required |
| 🟡 Yellow | Low margin (< threshold) | Review and consider adjusting |
| 🟢 Green | Profitable (>= threshold) | Good, maintain |
| ⚪ Gray | Incomplete data | Update missing costs |

### Navigation Quick Links

| Feature | Path |
|---------|------|
| Menu Cost View | Menu → Chi phí & Lợi nhuận → Chi phí món |
| Profit Analysis | Menu → Chi phí & Lợi nhuận → Phân tích lợi nhuận |
| Operating Expenses | Settings → Chi phí vận hành |
| Ingredient Management | Menu → Quản lý nguyên liệu |

---

## Support

Nếu bạn gặp vấn đề hoặc có câu hỏi:

1. **Check this guide** - Hầu hết câu hỏi đã được trả lời ở đây
2. **Check Troubleshooting section** - Common problems và solutions
3. **Contact support** - Nếu vẫn không giải quyết được

**Remember:** Tính năng này giúp bạn đưa ra quyết định kinh doanh chính xác. Hãy dành thời gian hiểu rõ và sử dụng thường xuyên!

---

*Last updated: February 2026*
*Version: 1.0*
